package main

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/WakkeWang/EmbedFlow/pkg/protocol"
	"github.com/WakkeWang/EmbedFlow/server/internal/session"
	"github.com/WakkeWang/EmbedFlow/server/internal/sessionlog"
	"github.com/WakkeWang/EmbedFlow/server/internal/store"
	"github.com/WakkeWang/EmbedFlow/server/internal/transport"
)

var errAuthFailed = errors.New("invalid token")

// coordinator is the M1 session executor: it implements transport.Hub,
// translating kernel events into persistence, session logs and wire side
// effects. T1+T2 scope: bootstrap auth stub, connection registry, session
// state, byte-accurate session logs.
type coordinator struct {
	kernel  *session.Kernel
	store   *store.Store
	dataDir string

	mu      sync.Mutex
	conns   map[*transport.ClientConn]struct{}
	loggers map[int64]*sessionlog.Logger // session ID -> open log writer
}

func (h *coordinator) validateToken(token string) (string, error) {
	// T1 stub: single bootstrap token; real auth (server-side sessions,
	// two-level roles, requirement 1.5) lands with the login slice.
	if token == "bootstrap" {
		return "admin", nil
	}
	return "", errAuthFailed
}

// --- transport.Hub ---

func (h *coordinator) OnControl(c *transport.ClientConn, f *protocol.Frame) {
	switch body := f.Body.(type) {
	case *protocol.ShareRequestFrame:
		// Serial-port sharing: bind this client connection to a device.
		// T5 wires the full open-session flow; T1 records the binding.
		res := h.kernel.Open(session.OpenRequest{
			User:   c.User(),
			Device: session.DeviceID(body.DeviceID),
			Kind:   session.KindManual,
			Now:    time.Now(),
		})
		if res.Rejected {
			// Device-busy is a session outcome, not an auth failure (4A
			// frame semantics): report it on the session-state channel.
			_ = c.Send(protocol.Frame{Type: protocol.FrameSessionState, Body: &protocol.SessionStateFrame{
				DeviceID: body.DeviceID,
				State:    "rejected",
				Detail:   res.Reason,
			}})
			return
		}
		// Open a session log writer for the new session (T2). Failure marks
		// the session log-incomplete from birth -- loud, not silent.
		if err := h.openLogger(int64(res.SessionID)); err != nil {
			slog.Error("open session log", "session", res.SessionID, "err", err)
		}
		_ = c.Send(protocol.Frame{Type: protocol.FrameSessionState, Body: &protocol.SessionStateFrame{
			SessionID: int64(res.SessionID),
			DeviceID:  body.DeviceID,
			State:     res.State,
		}})
	default:
		slog.Warn("unhandled control frame", "type", f.Type, "user", c.User())
	}
}

func (h *coordinator) OnBinary(c *transport.ClientConn, data []byte) {
	// T2 interim: connection-to-session binding arrives with T5. Until then
	// binary frames CANNOT be attributed to a session, so they are counted
	// and dropped -- writing them to every open logger would cross-contaminate
	// concurrent sessions (device A's bytes landing in device B's log).
	h.mu.Lock()
	_, tracked := h.conns[c]
	h.mu.Unlock()
	if !tracked {
		return
	}
}

// openLogger starts (or replaces) the log writer for a session.
func (h *coordinator) openLogger(sessionID int64) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.loggers[sessionID]; ok {
		return nil
	}
	l, err := sessionlog.Open(h.dataDir, sessionID)
	if err != nil {
		return err
	}
	h.loggers[sessionID] = l
	return nil
}

// closeLogger flushes, closes and forgets a session's log writer, marking
// the session record when the log turned out incomplete (design doc Failure
// modes: the disk-full gap -- the flag must outlive the logger).
func (h *coordinator) closeLogger(sessionID int64) {
	h.mu.Lock()
	l, ok := h.loggers[sessionID]
	delete(h.loggers, sessionID)
	h.mu.Unlock()
	if ok && l != nil {
		if err := l.Close(); err != nil {
			slog.Error("close session log", "session", sessionID, "err", err)
		}
		if l.Incomplete() {
			if err := h.store.MarkLogIncomplete(context.Background(), sessionID); err != nil {
				slog.Error("mark log incomplete", "session", sessionID, "err", err)
			}
			slog.Warn("session log incomplete (write failures during session)", "session", sessionID)
		}
	}
}

// loggerFor returns the open logger for a session, if any.
func (h *coordinator) loggerFor(sessionID int64) *sessionlog.Logger {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.loggers[sessionID]
}

func (h *coordinator) OnHeartbeat(c *transport.ClientConn, seq uint64) {
	// Feed the kernel's liveness clock (CEO-16A). Session attribution to a
	// specific bound session arrives with T5; T1 tracks it per connection.
	h.mu.Lock()
	h.conns[c] = struct{}{}
	h.mu.Unlock()
	h.kernel.ClientHeartbeat(0, time.Now()) // no-op until sessions bind (T5)
}

func (h *coordinator) OnClose(c *transport.ClientConn) {
	h.mu.Lock()
	delete(h.conns, c)
	h.mu.Unlock()
	// T5: disconnect semantics route through the kernel once connections
	// carry bound sessions.
}

// --- kernel side-effect application ---

func (h *coordinator) tick(now time.Time) {
	for _, ev := range h.kernel.Tick(now) {
		h.apply(ev)
	}
}

func (h *coordinator) sweepStartup() error {
	ctx := context.Background()
	actives, err := h.store.ActiveSessions(ctx)
	if err != nil {
		return err
	}
	if len(actives) > 0 {
		slog.Warn("startup sweep: failing leftover active sessions", "count", len(actives))
	}
	for _, s := range actives {
		h.kernel.Seed(session.Session{
			ID:        session.SessionID(s.ID),
			DeviceID:  session.DeviceID(s.DeviceID),
			Kind:      s.Kind,
			Owner:     s.Owner,
			State:     session.StateActive,
			StartedAt: s.StartedAt,
		})
	}
	for _, ev := range h.kernel.SweepStartup("server restart", time.Now()) {
		h.apply(ev)
	}
	return nil
}

// apply persists a kernel event. The kernel is pure; this is where the
// executor turns transitions into rows (design doc: events drive executors).
func (h *coordinator) apply(ev session.Event) {
	ctx := context.Background()
	switch ev.Kind {
	case session.EventSessionClosed, session.EventSessionFailed, session.EventSweepFailed:
		s, ok := h.kernel.Session(ev.SessionID)
		if !ok {
			return
		}
		if err := h.store.UpdateSession(ctx, store.Session{
			ID:        int64(s.ID),
			DeviceID:  int64(s.DeviceID),
			Kind:      s.Kind,
			Owner:     s.Owner,
			State:     s.State,
			StartedAt: s.StartedAt,
			EndedAt:   s.EndedAt,
			EndReason: s.EndReason,
		}); err != nil {
			slog.Error("persist session end", "session", ev.SessionID, "err", err)
		}
		// Flush and release the session's log writer (T2).
		h.closeLogger(int64(ev.SessionID))
		slog.Info("session ended", "session", ev.SessionID, "state", s.State, "reason", s.EndReason)
	case session.EventConnDown:
		slog.Warn("client connection down (detection)", "session", ev.SessionID)
	case session.EventConnWarning:
		slog.Warn("task session in reconnect grace window", "session", ev.SessionID, "remaining", ev.Detail)
	case session.EventConnRecovered:
		slog.Info("client connection recovered within grace window", "session", ev.SessionID)
	}
}

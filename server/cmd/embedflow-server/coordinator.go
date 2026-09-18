package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/WakkeWang/EmbedFlow/pkg/protocol"
	"github.com/WakkeWang/EmbedFlow/server/internal/session"
	"github.com/WakkeWang/EmbedFlow/server/internal/store"
	"github.com/WakkeWang/EmbedFlow/server/internal/transport"
)

// coordinator is the M1 session executor: it implements transport.Hub,
// translating kernel events into persistence and wire side effects. T1 scope
// keeps it minimal: auth stub, connection registry, session state broadcast.
type coordinator struct {
	kernel *session.Kernel
	store  *store.Store

	conns map[*transport.ClientConn]struct{}
}

func (h *coordinator) validateToken(token string) (string, error) {
	// T1 stub: single bootstrap token; real auth (server-side sessions,
	// two-level roles, requirement 1.5) lands with the login slice.
	if token == "bootstrap" {
		return "admin", nil
	}
	return "", errAuthFailed
}

var errAuthFailed = errAuth{}

type errAuth struct{}

func (errAuth) Error() string { return "invalid token" }

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
			_ = c.Send(protocol.Frame{Type: protocol.FrameAuthFail, Body: &protocol.AuthFailFrame{Reason: res.Reason}})
			return
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
	// T1 skeleton: no device binding yet, so binary payload is dropped.
	// T3/T4 route it through the session log and back out to subscribers.
}

func (h *coordinator) OnClose(c *transport.ClientConn) {
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
		slog.Info("session ended", "session", ev.SessionID, "state", s.State, "reason", s.EndReason)
	case session.EventConnDown:
		slog.Warn("client connection down (detection)", "session", ev.SessionID)
	case session.EventConnWarning:
		slog.Warn("task session in reconnect grace window", "session", ev.SessionID, "remaining", ev.Detail)
	case session.EventConnRecovered:
		slog.Info("client connection recovered within grace window", "session", ev.SessionID)
	}
}

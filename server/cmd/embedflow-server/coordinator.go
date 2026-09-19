package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/WakkeWang/EmbedFlow/pkg/protocol"
	"github.com/WakkeWang/EmbedFlow/server/internal/expect"
	"github.com/WakkeWang/EmbedFlow/server/internal/serialport"
	"github.com/WakkeWang/EmbedFlow/server/internal/session"
	"github.com/WakkeWang/EmbedFlow/server/internal/sessionlog"
	"github.com/WakkeWang/EmbedFlow/server/internal/store"
	"github.com/WakkeWang/EmbedFlow/server/internal/transport"
)

// tokenRing issues and validates opaque bearer tokens (requirement 1.5:
// server-side sessions; a restart forces re-login, so memory-only is right).
type tokenRing struct {
	mu      sync.Mutex
	byToken map[string]tokenInfo
}

type tokenInfo struct {
	user    string
	role    string
	expires time.Time
}

func newTokenRing() *tokenRing {
	return &tokenRing{byToken: map[string]tokenInfo{}}
}

const tokenTTL = 24 * time.Hour

// issue mints a token for an authenticated user.
func (t *tokenRing) issue(user, role string) (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	tok := hex.EncodeToString(buf)
	t.mu.Lock()
	t.byToken[tok] = tokenInfo{user: user, role: role, expires: time.Now().Add(tokenTTL)}
	t.mu.Unlock()
	return tok, nil
}

// issueForTesting mints a token without going through login (tests only).
// The role rides the token so role-gate tests exercise the real path.
func (t *tokenRing) issueForTesting(user string) string {
	return t.issueForRole(user, "admin")
}

// issueForRole mints a test token with an explicit role.
func (t *tokenRing) issueForRole(user, role string) string {
	tok, err := t.issue(user, role)
	if err != nil {
		panic(err)
	}
	return tok
}

// validate maps a token to (user, role).
func (t *tokenRing) validate(tok string) (string, string, bool) {
	t.mu.Lock()
	info, ok := t.byToken[tok]
	t.mu.Unlock()
	if !ok || time.Now().After(info.expires) {
		return "", "", false
	}
	return info.user, info.role, true
}

// connState is the per-connection binding the coordinator tracks.
type connState struct {
	// boundDevice is the shared serial device behind this connection
	// (serial clients only). Zero for browser-only connections.
	boundDevice int64
	boundPort   string
	// sessionID is the session this connection owns (browser or client).
	sessionID int64
	// readonly marks a second viewer (CEO-15A follow mode).
	readonly bool
}

// coordinator is the M1 executor: it implements transport.Hub, owns
// connection<->session bindings, runs expect sequences, and translates
// kernel events into persistence, logs and wire effects.
type coordinator struct {
	kernel  *session.Kernel
	store   *store.Store
	dataDir string
	tokens  *tokenRing

	mu      sync.Mutex
	conns   map[*transport.ClientConn]*connState
	loggers map[int64]*sessionlog.Logger
	// expectRuns tracks the live expect runner per session (for abort).
	expectRuns map[int64]*expect.Runner
	// sharedBySession links a task session to its serial transport bridge.
	sharedBySession map[int64]*sharedTransport
	// virtual is the demo device (CEO-18A); nil without -demo.
	virtual *virtualDevice
	// demoRuleID is the seeded demo rule (0 without -demo).
	demoRuleID int64
	// virtualSession tracks the demo device's live session for fan-out.
	virtualSession map[int64]int64
	// virtualPorts are the live engine ports for virtual runs.
	virtualPorts []*virtualPort
	// nextConfirmID numbers confirmation cards monotonically.
	nextConfirmID int64
}

func newCoordinator(kernel *session.Kernel, st *store.Store, dataDir string) *coordinator {
	return &coordinator{
		kernel:          kernel,
		store:           st,
		dataDir:         dataDir,
		tokens:          newTokenRing(),
		conns:           map[*transport.ClientConn]*connState{},
		loggers:         map[int64]*sessionlog.Logger{},
		expectRuns:      map[int64]*expect.Runner{},
		sharedBySession: map[int64]*sharedTransport{},
		virtualSession:  map[int64]int64{},
	}
}

// validateToken satisfies transport.AuthConfig; the role rides AuthOK (1.5).
func (h *coordinator) validateToken(token string) (string, string, error) {
	user, role, ok := h.tokens.validate(token)
	if !ok {
		return "", "", errAuthFailed
	}
	return user, role, nil
}

// login authenticates over HTTP and mints a token.
func (h *coordinator) login(username, password string) (string, string, error) {
	u, err := h.store.AuthenticateUser(context.Background(), username, password)
	if err != nil {
		return "", "", err
	}
	tok, err := h.tokens.issue(u.Username, u.Role)
	if err != nil {
		return "", "", err
	}
	return tok, u.Role, nil
}

// --- connection bookkeeping ---

// OnAuth registers a fresh connection (transport fires it post-handshake).
func (h *coordinator) OnAuth(c *transport.ClientConn) {
	h.mu.Lock()
	h.conns[c] = &connState{}
	h.mu.Unlock()
}

func (h *coordinator) state(c *transport.ClientConn) *connState {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.conns[c]
}

// --- transport.Hub ---

func (h *coordinator) OnControl(c *transport.ClientConn, f *protocol.Frame) {
	switch body := f.Body.(type) {
	case *protocol.ShareRequestFrame:
		h.onShare(c, body)
	case *protocol.SessionCtrlFrame:
		h.onSessionCtrl(c, body)
	case *protocol.ConfirmFrame:
		h.onConfirm(c, body)
	case *protocol.DeviceStatusFrame:
		// Requirement 3.5 dual-layer detection: the client reports COM
		// errors (cable pulled, port seized) explicitly. Surface loudly.
		slog.Error("serial client port error", "user", c.User(), "port", body.Port, "err", body.Error)
	default:
		slog.Warn("unhandled control frame", "type", f.Type, "user", c.User())
	}
}

// onShare binds a serial-client connection to a device (requirement 3.4:
// share a port, choose the owning device).
func (h *coordinator) onShare(c *transport.ClientConn, req *protocol.ShareRequestFrame) {
	st := h.state(c)
	if st == nil {
		return
	}
	h.mu.Lock()
	st.boundDevice = req.DeviceID
	st.boundPort = req.Port
	h.mu.Unlock()
	_ = c.Send(protocol.Frame{Type: protocol.FrameSessionState, Body: &protocol.SessionStateFrame{
		SessionID: 0,
		DeviceID:  req.DeviceID,
		State:     "shared",
		Detail:    req.Port,
	}})
	slog.Info("serial port shared", "user", c.User(), "device", req.DeviceID, "port", req.Port)
}

// onSessionCtrl dispatches browser session commands (open/close/follow/run).
func (h *coordinator) onSessionCtrl(c *transport.ClientConn, cmd *protocol.SessionCtrlFrame) {
	switch cmd.Command {
	case "open":
		h.openSession(c, cmd.DeviceID, session.KindManual, 0)
	case "follow":
		h.followSession(c, cmd.DeviceID)
	case "close":
		h.closeSession(c, cmd.SessionID)
	case "run":
		h.runExpect(c, cmd.DeviceID, cmd.RuleID)
	case "abort":
		h.abortExpect(c, cmd.SessionID)
	default:
		slog.Warn("unknown session command", "command", cmd.Command)
	}
}

// openSession starts a session owned by this user (CEO-15A: idempotency
// keyed on user+device+kind inside the kernel; busy rejections carry the
// occupier for the DS-2A popup).
func (h *coordinator) openSession(c *transport.ClientConn, deviceID int64, kind string, ruleID int64) {
	user := c.User()
	res := h.kernel.Open(session.OpenRequest{
		User:   user,
		Device: session.DeviceID(deviceID),
		Kind:   kind,
		Now:    time.Now(),
	})
	if res.Rejected {
		if res.Occupier != nil && res.Occupier.User == user {
			// Same user, different kind (or restart race): report their
			// existing session so the UI can navigate to it.
			_ = c.Send(protocol.Frame{Type: protocol.FrameSessionState, Body: &protocol.SessionStateFrame{
				SessionID:  int64(res.Occupier.SessionID),
				DeviceID:   deviceID,
				State:      "busy",
				Detail:     res.Reason,
				OwnerIsYou: true,
			}})
			return
		}
		fr := &protocol.SessionStateFrame{
			DeviceID: deviceID,
			State:    "busy",
			Detail:   res.Reason,
		}
		if res.Occupier != nil {
			fr.OccupierUser = res.Occupier.User
			fr.OccupierKind = res.Occupier.Kind
			fr.OccupierSince = res.Occupier.Since.Format(time.RFC3339)
		}
		_ = c.Send(protocol.Frame{Type: protocol.FrameSessionState, Body: fr})
		return
	}

	sid := int64(res.SessionID)
	// Persist only genuinely new sessions: the kernel's idempotent open
	// (CEO-15A) returns an existing session ID with Existed=true, and
	// re-inserting it would be a UNIQUE violation.
	if !res.Existed {
		if err := h.persistNewSession(sid, deviceID, kind, user); err != nil {
			slog.Error("persist new session", "err", err)
		}
		if err := h.openLogger(sid); err != nil {
			slog.Error("open session log", "session", sid, "err", err)
		}
	}

	h.mu.Lock()
	if st := h.conns[c]; st != nil {
		st.sessionID = sid
	}
	h.mu.Unlock()

	_ = c.Send(protocol.Frame{Type: protocol.FrameSessionState, Body: &protocol.SessionStateFrame{
		SessionID:  sid,
		DeviceID:   deviceID,
		State:      res.State,
		OwnerIsYou: true,
	}})
	slog.Info("session opened", "session", sid, "device", deviceID, "kind", kind, "user", user, "rule", ruleID)
}

// persistNewSession inserts the session row; kernel remains the source of
// truth for state.
func (h *coordinator) persistNewSession(sid, deviceID int64, kind, owner string) error {
	s, ok := h.kernel.Session(session.SessionID(sid))
	if !ok {
		return fmt.Errorf("session %d vanished", sid)
	}
	_, err := h.store.CreateSession(context.Background(), store.Session{
		ID:        sid,
		DeviceID:  deviceID,
		Kind:      kind,
		Owner:     owner,
		State:     s.State,
		StartedAt: s.StartedAt,
	})
	return err
}

// followSession attaches a second viewer read-only (CEO-15A).
func (h *coordinator) followSession(c *transport.ClientConn, deviceID int64) {
	dev := h.kernel.DeviceState(session.DeviceID(deviceID))
	if dev.SessionID == 0 {
		_ = c.Send(protocol.Frame{Type: protocol.FrameSessionState, Body: &protocol.SessionStateFrame{
			DeviceID: deviceID, State: "idle", Detail: "no active session",
		}})
		return
	}
	h.mu.Lock()
	if st := h.conns[c]; st != nil {
		st.sessionID = int64(dev.SessionID)
		st.readonly = true
	}
	h.mu.Unlock()
	_ = c.Send(protocol.Frame{Type: protocol.FrameSessionState, Body: &protocol.SessionStateFrame{
		SessionID: int64(dev.SessionID),
		DeviceID:  deviceID,
		State:     "readonly",
		Detail:    "read-only follow",
	}})
}

// closeSession ends a session at the owner's request.
func (h *coordinator) closeSession(c *transport.ClientConn, sessionID int64) {
	st := h.state(c)
	if st == nil || st.sessionID != sessionID {
		_ = c.Send(protocol.Frame{Type: protocol.FrameSessionState, Body: &protocol.SessionStateFrame{
			SessionID: sessionID, State: "rejected", Detail: "not your session",
		}})
		return
	}
	for _, ev := range h.kernel.Close(session.SessionID(sessionID), time.Now()) {
		h.apply(ev)
	}
}

// runExpect creates a task session and executes the rule's steps against
// the bound serial transport (CEO-17A: terminal run button; human typing is
// blocked by the UI while the run is live). In demo mode the virtual device
// stands in for the physical client.
func (h *coordinator) runExpect(c *transport.ClientConn, deviceID, ruleID int64) {
	// The device needs a byte source: a shared serial client, or the
	// built-in virtual device (demo mode).
	h.mu.Lock()
	var transport_ *transport.ClientConn
	for cc, st := range h.conns {
		if st.boundDevice == deviceID && st.sessionID == 0 {
			transport_ = cc
			break
		}
	}
	h.mu.Unlock()
	if transport_ == nil {
		h.mu.Lock()
		vd := h.virtual
		h.mu.Unlock()
		if vd == nil || vd.deviceID != deviceID {
			_ = c.Send(protocol.Frame{Type: protocol.FrameSessionState, Body: &protocol.SessionStateFrame{
				DeviceID: deviceID, State: "rejected", Detail: "device has no shared serial port",
			}})
			return
		}
	}

	rule, err := h.store.GetExpectRule(context.Background(), ruleID)
	if err != nil {
		_ = c.Send(protocol.Frame{Type: protocol.FrameSessionState, Body: &protocol.SessionStateFrame{
			DeviceID: deviceID, State: "rejected", Detail: "no such rule",
		}})
		return
	}
	steps, err := decodeSteps(rule.StepsJSON)
	if err != nil {
		_ = c.Send(protocol.Frame{Type: protocol.FrameSessionState, Body: &protocol.SessionStateFrame{
			DeviceID: deviceID, State: "rejected", Detail: fmt.Sprintf("bad rule: %v", err),
		}})
		return
	}

	// Open a task session for the run.
	res := h.kernel.Open(session.OpenRequest{
		User:   c.User(),
		Device: session.DeviceID(deviceID),
		Kind:   session.KindTask,
		Now:    time.Now(),
	})
	if res.Rejected {
		fr := &protocol.SessionStateFrame{DeviceID: deviceID, State: "busy", Detail: res.Reason}
		if res.Occupier != nil {
			fr.OccupierUser = res.Occupier.User
			fr.OccupierKind = res.Occupier.Kind
			fr.OccupierSince = res.Occupier.Since.Format(time.RFC3339)
		}
		_ = c.Send(protocol.Frame{Type: protocol.FrameSessionState, Body: fr})
		return
	}
	sid := int64(res.SessionID)
	if err := h.persistNewSession(sid, deviceID, session.KindTask, c.User()); err != nil {
		slog.Error("persist task session", "err", err)
	}
	if err := h.openLogger(sid); err != nil {
		slog.Error("open session log", "session", sid, "err", err)
	}

	// Track the session on the requesting connection and remember the
	// virtual-device session for fan-out.
	h.mu.Lock()
	if st := h.conns[c]; st != nil {
		st.sessionID = sid
	}
	h.virtualSession[deviceID] = sid
	h.expectRuns[sid] = nil
	h.mu.Unlock()

	// Engine transport: over the shared serial client when one is bound,
	// otherwise straight against the virtual device (demo mode).
	var enginePort serialport.Port
	if transport_ != nil {
		shared := newSharedTransport(sid, transport_, h)
		h.mu.Lock()
		h.sharedBySession[sid] = shared
		h.mu.Unlock()
		enginePort = shared
	} else {
		vp := newVirtualPort(h, deviceID, sid)
		h.mu.Lock()
		h.virtualPorts = append(h.virtualPorts, vp)
		h.mu.Unlock()
		enginePort = vp
	}

	// writeTX appends a TX line: masked for secret sends (issue #14),
	// escaped content otherwise.
	writeTX := func(data []byte, secret bool) {
		if secret {
			_ = h.writeLogMasked(sid, "TX")
		} else {
			_ = h.writeLog(sid, "TX", data)
		}
	}

	cfg := expect.Config{
		Steps:         steps,
		DefaultWait:   15 * time.Second,
		DefaultOnFail: expect.FailAbort,
		// TX logging goes through OnSendBytes so secret sends (issue #14)
		// land masked while the device still receives real bytes.
		OnSendBytes: writeTX,
		// Live progress (issue #8, DS-3A): each step start pushes a running
		// frame so the terminal UI's progress bar and abort button appear.
		// Fanned out to every session subscriber (the owner and any
		// read-only followers -- CEO-15A followers see the same progress).
		OnStep: func(index, total int) {
			desc := describeStep(steps[index])
			h.notifySession(sid, protocol.Frame{Type: protocol.FrameExpectProgress, Body: &protocol.ExpectProgressFrame{
				SessionID: sid,
				StepIndex: index,
				StepTotal: total,
				StepDesc:  desc,
				Phase:     "running",
			}})
		},
	}
	runner := expect.New(cfg, enginePort)
	h.mu.Lock()
	h.expectRuns[sid] = runner
	h.mu.Unlock()

	go func() {
		result := runner.Run()
		phase := "completed"
		detail := result.Detail
		switch {
		case result.Aborted:
			phase = "aborted"
		case !result.OK():
			phase = "failed"
		}
		h.notifySession(sid, protocol.Frame{Type: protocol.FrameExpectProgress, Body: &protocol.ExpectProgressFrame{
			SessionID: sid,
			StepIndex: result.StepIndex,
			StepTotal: len(steps),
			Phase:     phase,
			Detail:    detail,
		}})
		// Task session ends with the run: completed -> closed, else failed.
		now := time.Now()
		if phase == "completed" {
			for _, ev := range h.kernel.Close(session.SessionID(sid), now) {
				h.apply(ev)
			}
		} else {
			for _, ev := range h.kernel.ClientDisconnected(session.SessionID(sid), now) {
				h.apply(ev)
			}
		}
	}()
}

// describeStep renders a step for progress display (issue #8: "is it
// waiting, or stuck").
func describeStep(s expect.Step) string {
	switch {
	case s.Delay > 0:
		return fmt.Sprintf("delay %s", s.Delay)
	case s.Await != "":
		return fmt.Sprintf("await %q", s.Await)
	case s.Send != "":
		return "send"
	default:
		return "step"
	}
}

// abortExpect stops a running sequence (requirement 3.6: the operator's
// abort button; DS-7B confirm happens in the UI).
func (h *coordinator) abortExpect(c *transport.ClientConn, sessionID int64) {
	h.mu.Lock()
	runner := h.expectRuns[sessionID]
	h.mu.Unlock()
	if runner != nil {
		runner.Abort()
	}
}

// onConfirm resolves a human confirmation card (issue #9).
func (h *coordinator) onConfirm(c *transport.ClientConn, cf *protocol.ConfirmFrame) {
	st := h.state(c)
	if st == nil || st.sessionID != cf.SessionID {
		_ = c.Send(protocol.Frame{Type: protocol.FrameSessionState, Body: &protocol.SessionStateFrame{
			SessionID: cf.SessionID, State: "rejected", Detail: "not your session",
		}})
		return
	}
	if err := h.store.ResolveConfirmation(context.Background(), cf.ConfirmID, cf.Result, cf.Note); err != nil {
		slog.Error("resolve confirmation", "err", err)
		return
	}
	resolved := *cf
	resolved.State = "resolved"
	// Fan out so read-only followers see the card resolve too (CEO-15A).
	h.notifySession(cf.SessionID, protocol.Frame{Type: protocol.FrameConfirm, Body: &resolved})
}

// insertConfirmation creates a card: HTTP path (from the terminal page).
func (h *coordinator) insertConfirmation(sessionID int64, prompt string) (int64, error) {
	id, err := h.store.InsertConfirmation(context.Background(), sessionID, prompt)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// OnHeartbeat feeds the kernel liveness clock for the connection's session.
func (h *coordinator) OnHeartbeat(c *transport.ClientConn, seq uint64) {
	h.mu.Lock()
	st := h.conns[c]
	h.mu.Unlock()
	if st != nil && st.sessionID != 0 {
		h.kernel.ClientHeartbeat(session.SessionID(st.sessionID), time.Now())
	}
	_ = c.Send(protocol.Frame{Type: protocol.FrameHeartbeat, Body: &protocol.HeartbeatFrame{Seq: seq}})
}

// OnBinary routes serial bytes. Direction: from the serial client = RX
// (device -> server); from a browser session owner = TX.
func (h *coordinator) OnBinary(c *transport.ClientConn, data []byte) {
	st := h.state(c)
	if st == nil {
		return
	}
	if st.boundDevice != 0 {
		// Serial client: RX bytes. The device's live session (kernel device
		// state) is the attribution key -- the client itself owns no session.
		dev := h.kernel.DeviceState(session.DeviceID(st.boundDevice))
		sid := int64(dev.SessionID)
		_ = h.writeLog(sid, "RX", data)
		h.mu.Lock()
		shared := h.sharedBySession[sid]
		targets := make([]*transport.ClientConn, 0, 2)
		for cc, cst := range h.conns {
			if cst.sessionID == sid && cst.boundDevice == 0 {
				targets = append(targets, cc)
			}
		}
		h.mu.Unlock()
		if shared != nil {
			shared.push(data)
		}
		for _, cc := range targets {
			_ = cc.SendBinary(data)
		}
		return
	}
	// Browser owner keystrokes: input gating first (CEO-17A -- the UI hint
	// alone must not be the only gate), then pure passthrough to the
	// device's client. Any input also resets the manual idle timer
	// (decision 11A: "any input resets the timer").
	if !st.readonly {
		sess, ok := h.kernel.Session(session.SessionID(st.sessionID))
		if !ok || sess.State != session.StateActive {
			return // session already ended: no device input
		}
		if h.expectRunActive(st.sessionID) {
			return // an expect run owns the device input stream
		}
		h.kernel.UserInput(session.SessionID(st.sessionID), time.Now())
		// Manual TX logging (requirement 3.3: both directions with
		// timestamps). The expect-run path logs its own sends through
		// Config.OnSendBytes (secret-aware, issue #14); this is the
		// keystroke path and it is never secret. Logged only when the bytes
		// actually reached a byte sink: a TX line claims delivery.
		if h.deliverToDevice(st.sessionID, data) {
			_ = h.writeLog(st.sessionID, "TX", data)
		}
	}
}

// expectRunActive reports whether an expect sequence is running for the
// session (its keystrokes must not interleave with automatic sends).
func (h *coordinator) expectRunActive(sessionID int64) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	runner := h.expectRuns[sessionID]
	return runner != nil && !runner.Done()
}

// writeLog appends to a session's logger.
func (h *coordinator) writeLog(sessionID int64, dir sessionlog.Dir, data []byte) error {
	h.mu.Lock()
	l := h.loggers[sessionID]
	h.mu.Unlock()
	if l == nil {
		return nil
	}
	return l.Write(time.Now(), dir, data)
}

// writeLogMasked appends a masked TX line (issue #14: password sends).
func (h *coordinator) writeLogMasked(sessionID int64, dir sessionlog.Dir) error {
	h.mu.Lock()
	l := h.loggers[sessionID]
	h.mu.Unlock()
	if l == nil {
		return nil
	}
	return l.WriteMasked(time.Now(), dir)
}

// sendToDevice pushes browser TX bytes to the serial client bound to the
// session's device (requirement 3.4: pure passthrough); the virtual device
// receives them directly (CEO-18A).
func (h *coordinator) sendToDevice(sessionID int64, data []byte) {
	_ = h.deliverToDevice(sessionID, data)
}

// deliverToDevice routes TX bytes to the device's transport and reports
// whether they reached any byte sink. Callers log TX only on true: a TX log
// line asserts delivery.
func (h *coordinator) deliverToDevice(sessionID int64, data []byte) bool {
	dev, ok := h.kernel.Session(session.SessionID(sessionID))
	if !ok || dev.DeviceID == 0 {
		return false
	}
	h.mu.Lock()
	var target *transport.ClientConn
	for cc, st := range h.conns {
		if st.boundDevice == int64(dev.DeviceID) && st.sessionID == 0 {
			target = cc
			break
		}
	}
	h.mu.Unlock()
	if target != nil {
		err := target.SendBinary(data)
		return err == nil
	}
	// No physical client: the virtual device consumes the bytes.
	h.mu.Lock()
	vd := h.virtual
	h.mu.Unlock()
	if vd != nil && vd.deviceID == int64(dev.DeviceID) {
		vd.deliver(data)
		return true
	}
	return false
}

// OnClose tears down the connection: unbind, end owned session per kind
// (decision 1A disconnect branches).
func (h *coordinator) OnClose(c *transport.ClientConn) {
	h.mu.Lock()
	st := h.conns[c]
	delete(h.conns, c)
	h.mu.Unlock()
	if st == nil {
		return
	}
	if st.sessionID != 0 {
		for _, ev := range h.kernel.ClientDisconnected(session.SessionID(st.sessionID), time.Now()) {
			h.apply(ev)
		}
	}
}

// --- kernel side-effect application ---

func (h *coordinator) tick(now time.Time) {
	for _, ev := range h.kernel.Tick(now) {
		h.apply(ev)
	}
}

func (h *coordinator) sweepStartup() error {
	// Seed the kernel counter from the highest persisted ID (not just active
	// rows) so new session IDs never collide with finished history.
	if max, err := h.store.MaxSessionID(context.Background()); err == nil && max > 0 {
		h.kernel.SeedCounter(session.SessionID(max))
	}
	actives, err := h.store.ActiveSessions(context.Background())
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

// apply persists a kernel event and notifies bound connections.
func (h *coordinator) apply(ev session.Event) {
	switch ev.Kind {
	case session.EventSessionClosed, session.EventSessionFailed, session.EventSweepFailed:
		s, ok := h.kernel.Session(ev.SessionID)
		if !ok {
			return
		}
		st := store.Session{
			ID:        int64(s.ID),
			DeviceID:  int64(s.DeviceID),
			Kind:      s.Kind,
			Owner:     s.Owner,
			State:     s.State,
			StartedAt: s.StartedAt,
			EndedAt:   s.EndedAt,
			EndReason: s.EndReason,
		}
		if err := h.store.UpdateSession(context.Background(), st); err != nil {
			slog.Error("persist session end", "session", ev.SessionID, "err", err)
		}
		// Flush, close, flag incomplete (design doc disk-full gap).
		h.closeLogger(int64(ev.SessionID))

		// Drop the virtual port for this session (demo runs) and the shared
		// transport bridge -- it holds a dead client connection reference,
		// so leaving it would leak one entry per task session.
		h.mu.Lock()
		kept := h.virtualPorts[:0]
		for _, p := range h.virtualPorts {
			if p.sid != int64(ev.SessionID) {
				kept = append(kept, p)
			}
		}
		h.virtualPorts = kept
		delete(h.virtualSession, int64(ev.DeviceID))
		delete(h.sharedBySession, int64(ev.SessionID))
		h.mu.Unlock()

		// Notify subscribers on the session's device.
		h.notifyDevice(int64(s.DeviceID), &protocol.SessionStateFrame{
			SessionID: int64(s.ID),
			DeviceID:  int64(s.DeviceID),
			State:     s.State,
			Detail:    s.EndReason,
		})
		slog.Info("session ended", "session", ev.SessionID, "state", s.State, "reason", s.EndReason)
	case session.EventConnDown:
		slog.Warn("client connection down (detection)", "session", ev.SessionID)
		h.notifyConnState(ev, "conn_down", ev.Detail)
	case session.EventConnWarning:
		slog.Warn("task session in reconnect grace window", "session", ev.SessionID, "remaining", ev.Detail)
		h.notifyConnState(ev, "conn_warning", ev.Detail)
	case session.EventConnRecovered:
		slog.Info("client connection recovered within grace window", "session", ev.SessionID)
		h.notifyConnState(ev, "conn_recovered", "")
	}
}

// notifyConnState broadcasts a liveness transition to the session's
// subscribers: one shape for conn_down / conn_warning / conn_recovered.
func (h *coordinator) notifyConnState(ev session.Event, state, detail string) {
	h.notifySession(int64(ev.SessionID), protocol.Frame{Type: protocol.FrameSessionState, Body: &protocol.SessionStateFrame{
		SessionID: int64(ev.SessionID),
		DeviceID:  int64(ev.DeviceID),
		State:     state,
		Detail:    detail,
	}})
}

// notifyDevice pushes a state frame to all connections bound to a device.
func (h *coordinator) notifyDevice(deviceID int64, fr *protocol.SessionStateFrame) {
	h.mu.Lock()
	targets := make([]*transport.ClientConn, 0, 2)
	for cc, st := range h.conns {
		if st.boundDevice == deviceID || st.sessionID == fr.SessionID && st.boundDevice == 0 {
			targets = append(targets, cc)
		}
	}
	h.mu.Unlock()
	for _, cc := range targets {
		_ = cc.Send(protocol.Frame{Type: protocol.FrameSessionState, Body: fr})
	}
}

// notifySession pushes a frame to every connection subscribed to a session
// (owner + read-only followers): progress and confirmations are session
// events, not private to the creator's socket (CEO-15A).
func (h *coordinator) notifySession(sessionID int64, fr protocol.Frame) {
	h.mu.Lock()
	targets := make([]*transport.ClientConn, 0, 2)
	for cc, st := range h.conns {
		if st.sessionID == sessionID {
			targets = append(targets, cc)
		}
	}
	h.mu.Unlock()
	for _, cc := range targets {
		_ = cc.Send(fr)
	}
}

// --- loggers ---

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

func (h *coordinator) closeLogger(sessionID int64) {
	h.mu.Lock()
	l, ok := h.loggers[sessionID]
	delete(h.loggers, sessionID)
	delete(h.expectRuns, sessionID)
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

// decodeSteps parses stored rule JSON into engine steps.
func decodeSteps(jsonStr string) ([]expect.Step, error) {
	return expect.DecodeSteps([]byte(jsonStr))
}

var errAuthFailed = errors.New("invalid token")

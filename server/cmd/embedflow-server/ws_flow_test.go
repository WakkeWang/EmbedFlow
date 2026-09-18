package main

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"

	"github.com/WakkeWang/EmbedFlow/pkg/protocol"
	"github.com/WakkeWang/EmbedFlow/server/internal/paths"
	"github.com/WakkeWang/EmbedFlow/server/internal/session"
	"github.com/WakkeWang/EmbedFlow/server/internal/store"
)

// Integration tests for the full WS session flow: share -> open -> binary
// passthrough -> close, over real WebSockets against the coordinator, with
// store and sessionlog on temp dirs (issues #3, #5, #6).

type srv struct {
	*httptest.Server
	st     *store.Store
	kernel *session.Kernel
	hub    *coordinator
}

func newTestServer(t *testing.T) *srv {
	t.Helper()
	dir := t.TempDir()
	if err := paths.Ensure(dir); err != nil {
		t.Fatalf("paths: %v", err)
	}
	st, err := store.Open(dir + "/test.db")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { st.Close() })

	kernel := session.New(session.Options{IdleTimeout: time.Hour})
	hub := newCoordinator(kernel, st, dir)

	if err := st.CreateUser(context.Background(), "admin", "pw", "admin"); err != nil {
		t.Fatalf("seed admin: %v", err)
	}
	if err := st.CreateUser(context.Background(), "alice", "pw", "member"); err != nil {
		t.Fatalf("seed alice: %v", err)
	}

	s := &srv{Server: httptest.NewServer(hub.mux("")), st: st, kernel: kernel, hub: hub}
	t.Cleanup(s.Close)
	return s
}

type wsClient struct {
	t  *testing.T
	ws *websocket.Conn
}

func dial(t *testing.T, s *srv, user, _ string) *wsClient {
	t.Helper()
	tok := s.hub.tokens.issueForTesting(user)
	ws, _, err := websocket.Dial(context.Background(), s.URL+"/ws/client", nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	c := &wsClient{t: t, ws: ws}
	c.send(protocol.FrameAuth, &protocol.AuthFrame{
		Token: tok, ProtocolVersion: protocol.ProtocolVersion,
	})
	f := c.recv(2 * time.Second)
	if f.Type != protocol.FrameAuthOK {
		t.Fatalf("want AuthOK, got type %d", f.Type)
	}
	return c
}

func (c *wsClient) send(typ uint16, body any) {
	c.t.Helper()
	wire, err := protocol.Encode(body)
	if err != nil {
		c.t.Fatalf("encode: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := c.ws.Write(ctx, websocket.MessageText, wire); err != nil {
		c.t.Fatalf("write: %v", err)
	}
}

func (c *wsClient) sendBinary(data []byte) {
	c.t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := c.ws.Write(ctx, websocket.MessageBinary, data); err != nil {
		c.t.Fatalf("write binary: %v", err)
	}
}

func (c *wsClient) recv(timeout time.Duration) *protocol.Frame {
	c.t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		typ, data, err := c.ws.Read(ctx)
		if err != nil {
			c.t.Fatalf("read: %v", err)
		}
		if typ != websocket.MessageText {
			continue
		}
		f, err := protocol.ParseControl(data)
		if err != nil {
			c.t.Fatalf("parse: %v", err)
		}
		return f
	}
}

// recvUntil reads frames until pred matches, discarding others.
func (c *wsClient) recvUntil(timeout time.Duration, pred func(*protocol.Frame) bool) *protocol.Frame {
	c.t.Helper()
	deadline := time.Now().Add(timeout)
	for {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			c.t.Fatal("recvUntil: timed out")
		}
		f := c.recv(remaining)
		if pred(f) {
			return f
		}
	}
}

func (c *wsClient) close() { c.ws.CloseNow() }

func stateOf(f *protocol.Frame) *protocol.SessionStateFrame {
	b, _ := f.Body.(*protocol.SessionStateFrame)
	return b
}

// --- flows ---

func TestFlow_ShareOpenBinaryClose(t *testing.T) {
	s := newTestServer(t)
	devID, err := s.st.CreateDevice(context.Background(), "board", "demo")
	if err != nil {
		t.Fatalf("seed device: %v", err)
	}

	client := dial(t, s, "admin", "pw")
	browser := dial(t, s, "alice", "pw")

	// Client shares its serial port under the device.
	client.send(protocol.FrameShareRequest, &protocol.ShareRequestFrame{DeviceID: devID, Port: "COM3"})
	client.recvUntil(2*time.Second, func(f *protocol.Frame) bool {
		b := stateOf(f)
		return b != nil && b.DeviceID == devID && b.State == "shared"
	})

	// Browser opens a session.
	browser.send(protocol.FrameSessionCtrl, &protocol.SessionCtrlFrame{Command: "open", DeviceID: devID})
	f := browser.recvUntil(2*time.Second, func(f *protocol.Frame) bool {
		b := stateOf(f)
		return b != nil && b.DeviceID == devID && b.State == "active"
	})
	sess := stateOf(f)
	if sess.SessionID == 0 {
		t.Fatalf("open returned no session id: %+v", sess)
	}

	// Client sends serial bytes; the browser sees them (RX binary frame).
	client.sendBinary([]byte("device says hi"))
	bctx, bcancel := context.WithTimeout(context.Background(), 2*time.Second)
	btyp, bdata, err := browser.ws.Read(bctx)
	bcancel()
	if err != nil {
		t.Fatalf("browser binary read: %v", err)
	}
	if btyp != websocket.MessageBinary || string(bdata) != "device says hi" {
		t.Fatalf("browser got type=%v data=%q", btyp, bdata)
	}

	// Browser keystrokes reach the client (TX passthrough).
	browser.sendBinary([]byte("ls\r"))
	cctx, ccancel := context.WithTimeout(context.Background(), 2*time.Second)
	ctyp, cdata, err := client.ws.Read(cctx)
	ccancel()
	if err != nil {
		t.Fatalf("client binary read: %v", err)
	}
	if ctyp != websocket.MessageBinary || string(cdata) != "ls\r" {
		t.Fatalf("client got type=%v data=%q", ctyp, cdata)
	}

	// Close releases the device.
	browser.send(protocol.FrameSessionCtrl, &protocol.SessionCtrlFrame{Command: "close", SessionID: sess.SessionID})
	browser.recvUntil(2*time.Second, func(f *protocol.Frame) bool {
		b := stateOf(f)
		return b != nil && b.State == "closed"
	})
	client.close()
	browser.close()
}

func TestFlow_SecondViewerReadOnlyAndBusy(t *testing.T) {
	s := newTestServer(t)
	devID, _ := s.st.CreateDevice(context.Background(), "d", "p")

	alice := dial(t, s, "alice", "pw")
	bob := dial(t, s, "admin", "pw")

	// Alice opens (no client sharing needed for a manual session).
	alice.send(protocol.FrameSessionCtrl, &protocol.SessionCtrlFrame{Command: "open", DeviceID: devID})
	alice.recvUntil(2*time.Second, func(f *protocol.Frame) bool {
		b := stateOf(f)
		return b != nil && b.State == "active"
	})

	// Bob's open is busy-rejected with occupier info (DS-2A).
	bob.send(protocol.FrameSessionCtrl, &protocol.SessionCtrlFrame{Command: "open", DeviceID: devID})
	f := bob.recvUntil(2*time.Second, func(f *protocol.Frame) bool {
		b := stateOf(f)
		return b != nil && b.State == "busy"
	})
	b := stateOf(f)
	if b.OccupierUser != "alice" {
		t.Fatalf("occupier = %+v", b)
	}

	// Bob follows read-only (CEO-15A second viewer).
	bob.send(protocol.FrameSessionCtrl, &protocol.SessionCtrlFrame{Command: "follow", DeviceID: devID})
	bob.recvUntil(2*time.Second, func(f *protocol.Frame) bool {
		st := stateOf(f)
		return st != nil && st.State == "readonly"
	})

	// Alice's repeat open is idempotent (CEO-15A) and marked as hers.
	alice.send(protocol.FrameSessionCtrl, &protocol.SessionCtrlFrame{Command: "open", DeviceID: devID})
	f2 := alice.recvUntil(2*time.Second, func(f *protocol.Frame) bool {
		b := stateOf(f)
		return b != nil && b.State == "active"
	})
	if !stateOf(f2).OwnerIsYou {
		t.Fatalf("idempotent open should mark owner: %+v", stateOf(f2))
	}

	alice.close()
	bob.close()
}

func TestFlow_SessionLifecyclePersistsAndSweeps(t *testing.T) {
	s := newTestServer(t)
	devID, _ := s.st.CreateDevice(context.Background(), "d", "p")

	alice := dial(t, s, "alice", "pw")
	alice.send(protocol.FrameSessionCtrl, &protocol.SessionCtrlFrame{Command: "open", DeviceID: devID})
	f := alice.recvUntil(2*time.Second, func(f *protocol.Frame) bool {
		b := stateOf(f)
		return b != nil && b.State == "active"
	})
	sid := stateOf(f).SessionID

	// The session row exists and the log file was created.
	rec, err := s.st.GetSession(context.Background(), sid)
	if err != nil {
		t.Fatalf("session not persisted: %v", err)
	}
	if rec.Owner != "alice" || rec.Kind != "manual" || rec.State != "active" {
		t.Fatalf("record = %+v", rec)
	}

	// History endpoint source: sessions-for-device lists it.
	hist, _ := s.st.SessionsForDevice(context.Background(), devID)
	if len(hist) != 1 || hist[0].ID != sid {
		t.Fatalf("history = %+v", hist)
	}

	// Abrupt disconnect: manual session closes and releases the device
	// (decision 1A disconnect branch).
	alice.close()
	deadline := time.Now().Add(2 * time.Second)
	for {
		got, _ := s.st.GetSession(context.Background(), sid)
		if got.State == "closed" {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("session not closed after disconnect: %+v", got)
		}
		time.Sleep(20 * time.Millisecond)
	}
	if st := s.kernel.DeviceState(session.DeviceID(devID)); st.SessionID != 0 {
		t.Fatalf("device still occupied: %+v", st)
	}
}

func TestFlow_ExpectRunOnVirtualDevice(t *testing.T) {
	s := newTestServer(t)
	if err := s.hub.startDemoDevice(); err != nil {
		t.Fatalf("demo: %v", err)
	}

	rules, err := s.st.ListExpectRules(context.Background())
	if err != nil || len(rules) == 0 {
		t.Fatalf("demo rule missing: %v %v", rules, err)
	}

	admin := dial(t, s, "admin", "pw")
	admin.send(protocol.FrameSessionCtrl, &protocol.SessionCtrlFrame{
		Command: "run", DeviceID: s.hub.virtual.deviceID, RuleID: rules[0].ID,
	})

	// The run must reach completion and close the task session.
	f := admin.recvUntil(15*time.Second, func(fr *protocol.Frame) bool {
		p, ok := fr.Body.(*protocol.ExpectProgressFrame)
		return ok && (p.Phase == "completed" || p.Phase == "failed")
	})
	p := f.Body.(*protocol.ExpectProgressFrame)
	if p.Phase != "completed" {
		t.Fatalf("expect run failed: %+v", p)
	}
	if p.StepTotal != 12 {
		t.Fatalf("step total = %d, want 12", p.StepTotal)
	}

	// The session log captured the whole conversation.
	rec, _ := s.st.GetSession(context.Background(), 1)
	if rec.State != "closed" && rec.State != "active" {
		t.Fatalf("session state = %q", rec.State)
	}
	admin.close()
}

func TestFlow_ConfirmationCards(t *testing.T) {
	s := newTestServer(t)
	devID, _ := s.st.CreateDevice(context.Background(), "d", "p")

	alice := dial(t, s, "alice", "pw")
	alice.send(protocol.FrameSessionCtrl, &protocol.SessionCtrlFrame{Command: "open", DeviceID: devID})
	f := alice.recvUntil(2*time.Second, func(fr *protocol.Frame) bool {
		b := stateOf(fr)
		return b != nil && b.State == "active"
	})
	sid := stateOf(f).SessionID

	// Insert a card over HTTP, resolve it over WS (issue #9).
	cid, err := s.hub.insertConfirmation(sid, "LED on?")
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	alice.send(protocol.FrameConfirm, &protocol.ConfirmFrame{
		SessionID: sid, ConfirmID: cid, Prompt: "LED on?", State: "pending", Result: "pass", Note: "both green",
	})
	fr := alice.recvUntil(2*time.Second, func(x *protocol.Frame) bool {
		c, ok := x.Body.(*protocol.ConfirmFrame)
		return ok && c.ConfirmID == cid && c.State == "resolved"
	})
	got := fr.Body.(*protocol.ConfirmFrame)
	if got.Result != "pass" || got.Note != "both green" {
		t.Fatalf("resolved = %+v", got)
	}

	// Persisted for the report (requirement 4.3: confirmations enter the record).
	confs, _ := s.st.ConfirmationsForSession(context.Background(), sid)
	if len(confs) != 1 || confs[0].Result != "pass" || confs[0].State != "resolved" {
		t.Fatalf("confirmations = %+v", confs)
	}

	// Close cleanly so the logger releases the file before TempDir cleanup
	// (Windows file locks).
	for _, ev := range s.hub.kernel.Close(session.SessionID(sid), time.Now()) {
		s.hub.apply(ev)
	}
	time.Sleep(100 * time.Millisecond)
	alice.close()
}

var _ = strings.Contains

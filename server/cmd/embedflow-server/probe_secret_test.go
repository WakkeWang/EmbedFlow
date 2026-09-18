package main

import (
	"context"
	"encoding/binary"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"

	"github.com/WakkeWang/EmbedFlow/pkg/protocol"
	"github.com/WakkeWang/EmbedFlow/server/internal/session"
	"github.com/WakkeWang/EmbedFlow/server/internal/sessionlog"
	"github.com/WakkeWang/EmbedFlow/server/internal/store"
)

// Secret-send masking (issue #14): a step flagged secret writes a masked TX
// line to the session log while the device still receives the real bytes.

func readFile(t *testing.T, path string) ([]byte, error) {
	t.Helper()
	return os.ReadFile(path)
}

func TestExpectRun_SecretSendMaskedInLog(t *testing.T) {
	s := newTestServer(t)
	if err := s.hub.startDemoDevice(); err != nil {
		t.Fatalf("demo: %v", err)
	}
	rules, err := s.st.ListExpectRules(context.Background())
	if err != nil || len(rules) == 0 {
		t.Fatalf("demo rule missing: %v %v", rules, err)
	}

	// Overwrite the demo rule's password step with a secret flag: the demo
	// sequence sends `demo\r` at index 5 (after "Password:").
	raw := `[{"await":"Press ENTER"},{"send":"\\r"},{"await":"demo login:"},{"send":"demo\\r"},{"await":"Password:"},{"send":"demo\\r","secret":true},{"await":"demo#"}]`
	if err := s.st.UpdateExpectRule(context.Background(), store.ExpectRule{ID: rules[0].ID, Name: rules[0].Name, StepsJSON: raw}); err != nil {
		t.Fatalf("update rule: %v", err)
	}

	admin := dial(t, s, "admin", "pw")
	admin.send(protocol.FrameSessionCtrl, &protocol.SessionCtrlFrame{
		Command: "run", DeviceID: s.hub.virtual.deviceID, RuleID: rules[0].ID,
	})
	admin.recvUntil(15*time.Second, func(fr *protocol.Frame) bool {
		p, ok := fr.Body.(*protocol.ExpectProgressFrame)
		return ok && (p.Phase == "completed" || p.Phase == "failed")
	})

	// Poll for the buffered logger to flush (fixed sleeps flake on slow CI).
	waitForLog(t, s.hub.dataDir, 1, func(text string) bool {
		return strings.Contains(text, "Password: ") && strings.Contains(text, "TX | ***")
	})
	logData, err := readFile(t, sessionlog.SessionFilePath(s.hub.dataDir, 1))
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	text := string(logData)
	if !strings.Contains(text, "Password: ") {
		t.Fatalf("log missing device prompt (RX): %.200s", text)
	}
	if !strings.Contains(text, "TX | ***") {
		t.Fatalf("log missing masked TX line: %s", text)
	}
	// Only the password step is masked; the username send (also "demo") must
	// still be there. Scope the leak check to after the password prompt.
	after := text[strings.Index(text, "Password: "):]
	if strings.Contains(after, "TX | demo") {
		t.Fatalf("password leaked in TX log line: %s", after)
	}
	admin.close()
}

func TestExpectRun_PlainSendStillLogged(t *testing.T) {
	s := newTestServer(t)
	if err := s.hub.startDemoDevice(); err != nil {
		t.Fatalf("demo: %v", err)
	}
	rules, _ := s.st.ListExpectRules(context.Background())

	admin := dial(t, s, "admin", "pw")
	admin.send(protocol.FrameSessionCtrl, &protocol.SessionCtrlFrame{
		Command: "run", DeviceID: s.hub.virtual.deviceID, RuleID: rules[0].ID,
	})
	admin.recvUntil(15*time.Second, func(fr *protocol.Frame) bool {
		p, ok := fr.Body.(*protocol.ExpectProgressFrame)
		return ok && p.Phase == "completed"
	})
	waitForLog(t, s.hub.dataDir, 1, func(text string) bool {
		return strings.Contains(text, "TX | install")
	})
	logData, err := readFile(t, sessionlog.SessionFilePath(s.hub.dataDir, 1))
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	// Non-secret sends stay fully logged.
	if !strings.Contains(string(logData), "TX | install") {
		t.Fatalf("plain TX line missing: %s", logData)
	}
	admin.close()
}

// waitForLog polls until pred matches the log file, with a deadline (the
// logger flushes on a 200ms ticker; fixed sleeps flake on slow CI).
func waitForLog(t *testing.T, dataDir string, sessionID int64, pred func(string) bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for {
		logData, err := readFile(t, sessionlog.SessionFilePath(dataDir, sessionID))
		if err == nil && pred(string(logData)) {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("log condition not met within deadline")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// probe-echo (issue #2): binary frames come back unchanged.
func TestProbeEcho_BinaryRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		runProbeEchoHandler(w, r)
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ws, _, err := websocket.Dial(ctx, srv.URL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer ws.CloseNow()

	// Handshake.
	auth, _ := protocol.Encode(&protocol.AuthFrame{ProtocolVersion: protocol.ProtocolVersion})
	if err := ws.Write(ctx, websocket.MessageText, auth); err != nil {
		t.Fatalf("auth write: %v", err)
	}
	_, data, err := ws.Read(ctx)
	if err != nil {
		t.Fatalf("auth read: %v", err)
	}
	f, err := protocol.ParseControl(data)
	if err != nil || f.Type != protocol.FrameAuthOK {
		t.Fatalf("want AuthOK, got %v %v", f, err)
	}

	// Binary round trip.
	payload := make([]byte, 8)
	binary.BigEndian.PutUint64(payload, 42)
	if err := ws.Write(ctx, websocket.MessageBinary, payload); err != nil {
		t.Fatalf("write: %v", err)
	}
	for {
		typ, data, err := ws.Read(ctx)
		if err != nil {
			t.Fatalf("echo read: %v", err)
		}
		if typ == websocket.MessageBinary {
			if string(data) != string(payload) {
				t.Fatalf("echo mismatch: %q", data)
			}
			return
		}
	}
}

func TestProbeEcho_HeartbeatEcho(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		runProbeEchoHandler(w, r)
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ws, _, err := websocket.Dial(ctx, srv.URL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer ws.CloseNow()

	hb, _ := protocol.Encode(&protocol.HeartbeatFrame{Seq: 7})
	if err := ws.Write(ctx, websocket.MessageText, hb); err != nil {
		t.Fatalf("write: %v", err)
	}
	for {
		typ, data, err := ws.Read(ctx)
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		if typ != websocket.MessageText {
			continue
		}
		f, err := protocol.ParseControl(data)
		if err != nil {
			t.Fatalf("parse: %v", err)
		}
		if h, ok := f.Body.(*protocol.HeartbeatFrame); ok && h.Seq == 7 {
			return
		}
	}
}

// Manual keystroke TX logging (requirement 3.3): browser input over a shared
// physical client lands in the log as a TX line.
func TestManualKeystrokes_LoggedAsTX(t *testing.T) {
	s := newTestServer(t)
	devID, _ := s.st.CreateDevice(context.Background(), "d", "p")

	sess := dial(t, s, "alice", "pw")
	browser := dial(t, s, "admin", "pw")

	sess.send(protocol.FrameShareRequest, &protocol.ShareRequestFrame{DeviceID: devID, Port: "COM3"})
	sess.recvUntil(2*time.Second, func(f *protocol.Frame) bool {
		b := stateOf(f)
		return b != nil && b.DeviceID == devID && b.State == "shared"
	})

	browser.send(protocol.FrameSessionCtrl, &protocol.SessionCtrlFrame{Command: "open", DeviceID: devID})
	f := browser.recvUntil(2*time.Second, func(fr *protocol.Frame) bool {
		b := stateOf(fr)
		return b != nil && b.DeviceID == devID && b.State == "active"
	})
	sid := stateOf(f).SessionID

	browser.sendBinary([]byte("ls\r"))
	// Drain the echo frame on the client side so the write is consumed.
	_, _, _ = sess.ws.Read(context.Background())

	deadline := time.Now().Add(2 * time.Second)
	for {
		logData, err := readFile(t, sessionlog.SessionFilePath(s.hub.dataDir, sid))
		if err == nil && strings.Contains(string(logData), "TX | ls") {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("manual TX never logged")
		}
		time.Sleep(20 * time.Millisecond)
	}

	// Clean close (Windows file locks).
	for _, ev := range s.hub.kernel.Close(session.SessionID(sid), time.Now()) {
		s.hub.apply(ev)
	}
	time.Sleep(100 * time.Millisecond)
	sess.close()
	browser.close()
}

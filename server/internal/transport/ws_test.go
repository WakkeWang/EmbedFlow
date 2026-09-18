package transport

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"

	"github.com/WakkeWang/EmbedFlow/pkg/protocol"
)

// Integration tests over real WebSocket connections against httptest:
// the client-auth flow, version handshake, unknown-frame tolerance and
// binary frame routing (decision 4A, CEO-11A, CEO-3A).

type testHub struct {
	frames     chan *protocol.Frame
	binary     chan []byte
	heartbeats chan heartbeatEvt
}

type heartbeatEvt struct {
	c   *ClientConn
	seq uint64
}

func (h *testHub) OnControl(c *ClientConn, f *protocol.Frame) {
	if h.frames != nil {
		h.frames <- f
	}
}

func (h *testHub) OnBinary(c *ClientConn, data []byte) {
	if h.binary != nil {
		cp := make([]byte, len(data))
		copy(cp, data)
		h.binary <- cp
	}
}

func (h *testHub) OnHeartbeat(c *ClientConn, seq uint64) {
	if h.heartbeats != nil {
		h.heartbeats <- heartbeatEvt{c: c, seq: seq}
	}
}

func (h *testHub) OnClose(c *ClientConn) {}

func startServer(t *testing.T, hub Hub, auth *AuthConfig) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(Handler(hub, auth))
	t.Cleanup(srv.Close)
	return srv
}

func defaultAuth() *AuthConfig {
	return &AuthConfig{
		ValidateToken: func(token string) (string, error) {
			if token == "good-token" {
				return "alice", nil
			}
			return "", errBadToken
		},
		ProtocolVersion: protocol.ProtocolVersion,
	}
}

func dialAndAuth(t *testing.T, srv *httptest.Server, token string, version uint16) (*testConn, *httptest.Server) {
	t.Helper()
	ws, _, err := websocket.Dial(context.Background(), wsURL(srv.URL), nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	tc := &testConn{ws: ws}
	authWire, _ := protocol.Encode(&protocol.AuthFrame{Token: token, ProtocolVersion: version})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := ws.Write(ctx, websocket.MessageText, authWire); err != nil {
		t.Fatalf("write auth: %v", err)
	}
	return tc, srv
}

type testConn struct {
	ws *websocket.Conn
}

func (tc *testConn) readFrame(t *testing.T, timeout time.Duration) (websocket.MessageType, []byte) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	typ, data, err := tc.ws.Read(ctx)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	return typ, data
}

func (tc *testConn) close(t *testing.T) {
	t.Helper()
	tc.ws.CloseNow()
}

func wsURL(base string) string {
	// httptest gives http://; websocket scheme maps directly for coder/websocket.
	return base + "/ws/client"
}

func TestWS_AuthOK_ReceivesAuthOK(t *testing.T) {
	hub := &testHub{frames: make(chan *protocol.Frame, 8)}
	srv := startServer(t, hub, defaultAuth())
	tc, _ := dialAndAuth(t, srv, "good-token", protocol.ProtocolVersion)
	defer tc.close(t)

	typ, data := tc.readFrame(t, 2*time.Second)
	if typ != websocket.MessageText {
		t.Fatalf("frame type = %v, want text", typ)
	}
	f, err := protocol.ParseControl(data)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if f.Type != protocol.FrameAuthOK {
		t.Fatalf("type = %d, want AuthOK(%d): %s", f.Type, protocol.FrameAuthOK, data)
	}
}

func TestWS_BadToken_RejectedAndClosed(t *testing.T) {
	hub := &testHub{}
	srv := startServer(t, hub, defaultAuth())
	tc, _ := dialAndAuth(t, srv, "wrong", protocol.ProtocolVersion)
	defer tc.close(t)

	_, data := tc.readFrame(t, 2*time.Second)
	f, err := protocol.ParseControl(data)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if f.Type != protocol.FrameAuthFail {
		t.Fatalf("type = %d, want AuthFail", f.Type)
	}
	// Connection must be closed by the server after AuthFail.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	for {
		_, _, err = tc.ws.Read(ctx)
		if err != nil {
			break // expected
		}
	}
}

func TestWS_VersionMismatch_RejectedWithUpdateHint(t *testing.T) {
	// CEO-11A: a client speaking a different protocol version is refused with
	// an explicit "client needs update" reason.
	hub := &testHub{}
	srv := startServer(t, hub, defaultAuth())
	tc, _ := dialAndAuth(t, srv, "good-token", protocol.ProtocolVersion+10)
	defer tc.close(t)

	_, data := tc.readFrame(t, 2*time.Second)
	f, err := protocol.ParseControl(data)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	fail, ok := f.Body.(*protocol.AuthFailFrame)
	if f.Type != protocol.FrameAuthFail || !ok {
		t.Fatalf("want AuthFail frame, got type %d body %T", f.Type, f.Body)
	}
	if !strings.Contains(fail.Reason, "update") {
		t.Fatalf("reason = %q, want it to mention update", fail.Reason)
	}
}

func TestWS_UnknownFrameType_ConnectionStaysUp(t *testing.T) {
	// CEO-3A: unknown control frame = ignore + warn, never disconnect.
	hub := &testHub{frames: make(chan *protocol.Frame, 8)}
	srv := startServer(t, hub, defaultAuth())
	tc, _ := dialAndAuth(t, srv, "good-token", protocol.ProtocolVersion)
	defer tc.close(t)
	tc.readFrame(t, 2*time.Second) // AuthOK

	// Send an unknown-type control frame.
	tc.ws.Write(context.Background(), websocket.MessageText, []byte(`{"type":9999,"body":{}}`))

	// Follow with a heartbeat: if the server had disconnected us, this fails.
	beat, _ := protocol.Encode(&protocol.HeartbeatFrame{Seq: 1})
	if err := tc.ws.Write(context.Background(), websocket.MessageText, beat); err != nil {
		t.Fatalf("connection dropped after unknown frame: %v", err)
	}
	deadline := time.After(2 * time.Second)
	for {
		typ, data := tc.readFrame(t, 2*time.Second)
		if typ != websocket.MessageText {
			continue
		}
		f, err := protocol.ParseControl(data)
		if err != nil {
			t.Fatalf("parse: %v", err)
		}
		if f.Type == protocol.FrameHeartbeat {
			return // heartbeat echoed: connection alive
		}
		select {
		case <-deadline:
			t.Fatal("no heartbeat echo before deadline")
		default:
		}
	}
}

func TestWS_BinaryFrames_ForwardedToHub(t *testing.T) {
	hub := &testHub{binary: make(chan []byte, 8)}
	srv := startServer(t, hub, defaultAuth())
	tc, _ := dialAndAuth(t, srv, "good-token", protocol.ProtocolVersion)
	defer tc.close(t)
	tc.readFrame(t, 2*time.Second) // AuthOK

	payload := []byte{0x01, 0x02, 0xFF, 0x00, 0x55}
	if err := tc.ws.Write(context.Background(), websocket.MessageBinary, payload); err != nil {
		t.Fatalf("write binary: %v", err)
	}
	select {
	case got := <-hub.binary:
		if string(got) != string(payload) {
			t.Fatalf("binary = %v, want %v", got, payload)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("binary frame never reached the hub")
	}
}

func TestWS_HeartbeatEchoed(t *testing.T) {
	hub := &testHub{frames: make(chan *protocol.Frame, 8)}
	srv := startServer(t, hub, defaultAuth())
	tc, _ := dialAndAuth(t, srv, "good-token", protocol.ProtocolVersion)
	defer tc.close(t)
	tc.readFrame(t, 2*time.Second) // AuthOK

	beat, _ := protocol.Encode(&protocol.HeartbeatFrame{Seq: 42})
	if err := tc.ws.Write(context.Background(), websocket.MessageText, beat); err != nil {
		t.Fatalf("write: %v", err)
	}
	for {
		typ, data := tc.readFrame(t, 2*time.Second)
		if typ != websocket.MessageText {
			continue
		}
		f, err := protocol.ParseControl(data)
		if err != nil {
			t.Fatalf("parse: %v", err)
		}
		if f.Type == protocol.FrameHeartbeat {
			hb := f.Body.(*protocol.HeartbeatFrame)
			if hb.Seq != 42 {
				t.Fatalf("seq = %d, want 42", hb.Seq)
			}
			return
		}
	}
}

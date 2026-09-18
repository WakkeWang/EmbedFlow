package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coder/websocket"

	"github.com/WakkeWang/EmbedFlow/pkg/protocol"
	"github.com/WakkeWang/EmbedFlow/serial-client"
)

// fakeServer implements the server half of the protocol for unit-level
// integration tests of the client core (the full server flow is covered by
// the server package's WS integration tests).

type fakeServer struct {
	authed   chan struct{}
	shared   chan *protocol.ShareRequestFrame
	toDevice chan []byte
}

func newFakeServer(t *testing.T) (*fakeServer, *httptest.Server) {
	fs := &fakeServer{
		authed:   make(chan struct{}, 1),
		shared:   make(chan *protocol.ShareRequestFrame, 1),
		toDevice: make(chan []byte, 16),
	}
	upgrader := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer ws.CloseNow()
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Handshake: expect auth, answer AuthOK.
		typ, data, err := ws.Read(ctx)
		if err != nil || typ != websocket.MessageText {
			return
		}
		f, err := protocol.ParseControl(data)
		if err != nil || f.Type != protocol.FrameAuth {
			return
		}
		_ = ws.Write(ctx, websocket.MessageText, mustEncode(t, &protocol.AuthOKFrame{ServerVersion: "test"}))
		fs.authed <- struct{}{}

		for {
			typ, data, err := ws.Read(ctx)
			if err != nil {
				return
			}
			switch typ {
			case websocket.MessageText:
				f, err := protocol.ParseControl(data)
				if err != nil {
					continue
				}
				if sf, ok := f.Body.(*protocol.ShareRequestFrame); ok {
					fs.shared <- sf
					_ = ws.Write(ctx, websocket.MessageText,
						mustEncode(t, &protocol.SessionStateFrame{DeviceID: sf.DeviceID, State: "shared"}))
				}
			case websocket.MessageBinary:
				select {
				case fs.toDevice <- data:
				default:
				}
			}
		}
	})
	return fs, httptest.NewServer(upgrader)
}

func mustEncode(t *testing.T, body any) []byte {
	t.Helper()
	wire, err := protocol.Encode(body)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	return wire
}

func TestLogin_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/login" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"token":"tok-1","role":"member"}`))
	}))
	defer srv.Close()

	tok, err := client.Login(srv.URL, "u", "p")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if tok != "tok-1" {
		t.Fatalf("token = %q", tok)
	}
}

func TestLogin_BadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()
	if _, err := client.Login(srv.URL, "u", "p"); err == nil {
		t.Fatal("401 should error")
	}
}

func TestConnect_AuthHandshake(t *testing.T) {
	fs, srv := newFakeServer(t)
	defer srv.Close()

	cl := client.New(srv.URL, "tok")
	if err := cl.Connect(context.Background()); err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer cl.Close()

	select {
	case <-fs.authed:
	case <-time.After(2 * time.Second):
		t.Fatal("client never authenticated")
	}
}

func TestShare_RequestReachesServer(t *testing.T) {
	fs, srv := newFakeServer(t)
	defer srv.Close()

	cl := client.New(srv.URL, "tok")
	if err := cl.Connect(context.Background()); err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer cl.Close()

	port := &loopPort{}
	cl.SetPort(port)
	if err := cl.Share(7, "COM3"); err != nil {
		t.Fatalf("share: %v", err)
	}

	select {
	case sf := <-fs.shared:
		if sf.DeviceID != 7 || sf.Port != "COM3" {
			t.Fatalf("share = %+v", sf)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("share request never arrived")
	}
}

func TestDoubleShare_Rejected(t *testing.T) {
	_, srv := newFakeServer(t)
	defer srv.Close()

	cl := client.New(srv.URL, "tok")
	if err := cl.Connect(context.Background()); err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer cl.Close()

	if err := cl.Share(1, "COM3"); err != nil {
		t.Fatalf("first share: %v", err)
	}
	if err := cl.Share(2, "COM4"); err == nil {
		t.Fatal("second share should be rejected")
	}
}

// loopPort is the local serial port stand-in.
type loopPort struct{}

func (l *loopPort) Read(p []byte) (int, error) {
	time.Sleep(50 * time.Millisecond)
	return 0, nil
}

func (l *loopPort) Write(p []byte) (int, error) { return len(p), nil }
func (l *loopPort) Close() error                { return nil }

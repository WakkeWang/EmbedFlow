// Package transport is the WebSocket execution layer for the serial tunnel:
// it accepts client/browser connections, runs the auth + version handshake,
// decodes control frames, routes binary frames, and hands events to a Hub
// (decision 4A). Liveness detection lives with the session kernel; this
// package owns the wire.
package transport

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"

	"github.com/WakkeWang/EmbedFlow/pkg/protocol"
)

// errBadToken marks token validation failures without leaking why.
var errBadToken = errors.New("invalid token")

// AuthConfig carries the handshake policy.
type AuthConfig struct {
	// ValidateToken maps a token to its owning user, or errors.
	ValidateToken func(token string) (string, error)
	// ProtocolVersion is the version this server speaks (CEO-11A).
	ProtocolVersion uint16
	// ServerVersion is reported back in AuthOK for client-side checks.
	ServerVersion string
}

// Hub receives decoded transport events. The session executor implements it.
type Hub interface {
	// OnAuth fires once per connection after a successful handshake; the
	// executor registers its per-connection state here.
	OnAuth(c *ClientConn)
	// OnControl receives every known control frame after auth.
	OnControl(c *ClientConn, f *protocol.Frame)
	// OnBinary receives serial byte-stream payloads (data frames).
	OnBinary(c *ClientConn, data []byte)
	// OnHeartbeat receives client liveness pings so the executor can drive
	// the kernel's dual-threshold check (CEO-16A). The transport still echoes
	// them itself.
	OnHeartbeat(c *ClientConn, seq uint64)
	// OnClose fires when a connection goes away (cleanly or not).
	OnClose(c *ClientConn)
}

// ClientConn is one live WebSocket connection with its authenticated identity.
type ClientConn struct {
	Conn *websocket.Conn

	mu     sync.Mutex
	user   string
	kind   string // "client" (browser vs device-client distinction rides the hub's connState)
	authed bool
	closed bool
}

// User returns the authenticated username.
func (c *ClientConn) User() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.user
}

// Kind returns the connection kind ("client" or "browser").
func (c *ClientConn) Kind() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.kind
}

// Authed reports whether the handshake completed.
func (c *ClientConn) Authed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.authed
}

// Send enqueues a control frame to this connection.
func (c *ClientConn) Send(f protocol.Frame) error {
	body, err := protocol.Encode(f.Body)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return errors.New("transport: connection closed")
	}
	return c.Conn.Write(ctx, websocket.MessageText, body)
}

// SendBinary writes raw serial bytes as a binary data frame.
func (c *ClientConn) SendBinary(data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return errors.New("transport: connection closed")
	}
	return c.Conn.Write(ctx, websocket.MessageBinary, data)
}

const (
	writeTimeout = 5 * time.Second
	readTimeout  = 60 * time.Second // heartbeats arrive well within this
)

// Handler builds the WebSocket HTTP handler.
func Handler(hub Hub, auth *AuthConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			// Cross-origin browser upgrades are refused by the SameOrigin
			// middleware wrapping this handler (decision 3A). Here we accept
			// whatever reached us and rely on token auth.
			OriginPatterns: []string{"*"},
		})
		if err != nil {
			slog.Warn("ws accept failed", "err", err)
			return
		}
		ws.SetReadLimit(1 << 16) // 64KB max frame (decision 8A cap)
		c := &ClientConn{Conn: ws}
		defer func() {
			c.mu.Lock()
			c.closed = true
			c.mu.Unlock()
			hub.OnClose(c)
			ws.CloseNow()
		}()

		if !handshake(c, hub, auth) {
			return
		}
		readLoop(c, hub)
	})
}

// handshake runs the auth exchange; it returns false when the connection was
// refused (bad token or version mismatch) and already closed.
func handshake(c *ClientConn, hub Hub, auth *AuthConfig) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	typ, data, err := c.Conn.Read(ctx)
	if err != nil {
		slog.Warn("ws: no auth frame", "err", err)
		return false
	}
	if typ != websocket.MessageText {
		sendAuthFail(c, "first frame must be a control (text) frame")
		return false
	}
	f, err := protocol.ParseControl(data)
	if err != nil {
		sendAuthFail(c, "malformed auth frame")
		return false
	}
	af, ok := f.Body.(*protocol.AuthFrame)
	if !ok {
		sendAuthFail(c, "expected auth frame")
		return false
	}
	if af.ProtocolVersion != auth.ProtocolVersion {
		// CEO-11A: explicit update hint -- the client cannot be trusted to
		// interpret anything else.
		sendAuthFail(c, fmt.Sprintf("protocol version mismatch: client %d, server %d; client needs update", af.ProtocolVersion, auth.ProtocolVersion))
		return false
	}
	user, err := auth.ValidateToken(af.Token)
	if err != nil {
		sendAuthFail(c, "authentication failed")
		return false
	}
	c.mu.Lock()
	c.user = user
	c.kind = "client"
	c.authed = true
	c.mu.Unlock()

	if err := c.Send(protocol.Frame{Type: protocol.FrameAuthOK, Body: &protocol.AuthOKFrame{ServerVersion: auth.ServerVersion}}); err != nil {
		return false
	}
	hub.OnAuth(c)
	return true
}

func sendAuthFail(c *ClientConn, reason string) {
	_ = c.Send(protocol.Frame{Type: protocol.FrameAuthFail, Body: &protocol.AuthFailFrame{Reason: reason}})
	slog.Warn("ws: handshake refused", "reason", reason)
}

// readLoop decodes frames until the connection dies. Unknown control types
// are logged and ignored -- never fatal (CEO-3A).
func readLoop(c *ClientConn, hub Hub) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
		typ, data, err := c.Conn.Read(ctx)
		cancel()
		if err != nil {
			return
		}
		switch typ {
		case websocket.MessageBinary:
			// Data frame: raw serial bytes; route untouched (decision 4A).
			if c.Authed() {
				hub.OnBinary(c, data)
			}
		case websocket.MessageText:
			f, err := protocol.ParseControl(data)
			if err != nil {
				// Malformed control frame: log and drop the frame, not the
				// connection; the design doc's "disconnect on malformed" only
				// applies to auth-time garbage.
				slog.Warn("ws: bad control frame", "err", err)
				continue
			}
			if _, unknown := f.Body.(*protocol.UnknownBody); unknown {
				slog.Warn("ws: unknown frame type ignored", "type", f.Type)
				continue
			}
			if !c.Authed() {
				continue
			}
			// Heartbeats are transport-internal: echo for RTT inspection,
			// and hand to the hub so the kernel's liveness clock advances
			// (CEO-16A). Without this, lastHeartbeat freezes at open time
			// and live task sessions get killed by the grace window.
			if hb, ok := f.Body.(*protocol.HeartbeatFrame); ok {
				_ = c.Send(protocol.Frame{Type: protocol.FrameHeartbeat, Body: hb})
				hub.OnHeartbeat(c, hb.Seq)
				continue
			}
			hub.OnControl(c, f)
		}
	}
}

// Package client implements the serial-client core: login, port sharing
// and the passthrough loop over the shared WS protocol (requirement 3.4:
// pure passthrough, no business logic; decision 12A: reconnect restores
// the shared port, never a task session). A thin GUI (Wails) wraps this
// later; the CLI exercises the same core.
package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"

	"github.com/WakkeWang/EmbedFlow/pkg/protocol"
)

// Client is the serial-client core: one WS connection to the server, one
// bound serial port, bytes through both ways.
type Client struct {
	serverURL string
	token     string

	mu        sync.Mutex
	conn      *websocket.Conn
	shared    bool
	deviceID  int64
	port      string
	localPort PortOpener

	stop chan struct{}
	done sync.WaitGroup
}

// New creates a client pointed at a server (e.g. http://192.168.0.80:8420).
func New(serverURL, token string) *Client {
	return &Client{
		serverURL: strings.TrimRight(serverURL, "/"),
		token:     token,
		stop:      make(chan struct{}),
	}
}

// Login exchanges username+password for a token over the REST API
// (requirement 1.5). The token is stored by the caller (DPAPI in the GUI,
// decision 13A).
func Login(serverURL, username, password string) (string, error) {
	body := fmt.Sprintf(`{"username":%q,"password":%q}`, username, password)
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(serverURL, "/")+"/api/login", strings.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := new(http.Client).Do(req)
	if err != nil {
		return "", fmt.Errorf("client: login: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("client: login: status %d", resp.StatusCode)
	}
	var out struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("client: login decode: %w", err)
	}
	return out.Token, nil
}

// Connect dials the WS endpoint and authenticates (CEO-11A version field).
func (c *Client) Connect(ctx context.Context) error {
	wsURL := strings.Replace(c.serverURL, "http", "ws", 1) + "/ws/client"
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	ws, _, err := websocket.Dial(dialCtx, wsURL, &websocket.DialOptions{
		HTTPClient: &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: false}}},
	})
	if err != nil {
		return fmt.Errorf("client: dial: %w", err)
	}
	c.mu.Lock()
	c.conn = ws
	c.mu.Unlock()

	authWire, err := protocol.Encode(&protocol.AuthFrame{
		Token:           c.token,
		ProtocolVersion: protocol.ProtocolVersion,
		ClientVersion:   "serial-cli-dev",
	})
	if err != nil {
		return err
	}
	wctx, wcancel := context.WithTimeout(ctx, 5*time.Second)
	defer wcancel()
	if err := ws.Write(wctx, websocket.MessageText, authWire); err != nil {
		return fmt.Errorf("client: auth write: %w", err)
	}
	typ, data, err := ws.Read(wctx)
	if err != nil {
		return fmt.Errorf("client: auth read: %w", err)
	}
	if typ != websocket.MessageText {
		return errors.New("client: auth: expected control frame")
	}
	f, err := protocol.ParseControl(data)
	if err != nil {
		return err
	}
	if f.Type != protocol.FrameAuthOK {
		return fmt.Errorf("client: auth refused: %s", f.Body)
	}

	c.done.Add(1)
	go c.readLoop()
	return nil
}

// Share binds the local serial port to a registered device and starts the
// passthrough loop (requirement 3.4).
func (c *Client) Share(deviceID int64, portName string) error {
	c.mu.Lock()
	if c.shared {
		c.mu.Unlock()
		return errors.New("client: already sharing")
	}
	c.shared = true
	c.deviceID = deviceID
	c.port = portName
	c.mu.Unlock()

	c.send(&protocol.ShareRequestFrame{DeviceID: deviceID, Port: portName})
	return nil
}

// PortOpener abstracts the OS serial port (go.bug.st/serial in the real
// build; in-memory in tests). Byte-level, pure passthrough.
type PortOpener interface {
	// Read blocks for device bytes; Write sends operator bytes.
	Read(p []byte) (int, error)
	Write(p []byte) (int, error)
	Close() error
}

// SetPort attaches the local serial port and starts the passthrough
// pumps: port bytes -> server, server bytes -> port (decision 12A: on
// reconnect, Share + SetPort restore the shared state; never a session).
func (c *Client) SetPort(p PortOpener) {
	c.SetLocalPort(p)
	c.done.Add(1)
	go c.portToServer(p)
}

// portToServer pumps local port bytes to the server (binary data frames).
func (c *Client) portToServer(p PortOpener) {
	defer c.done.Done()
	buf := make([]byte, 4096)
	for {
		select {
		case <-c.stop:
			return
		default:
		}
		n, err := p.Read(buf)
		if n > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			c.mu.Lock()
			conn := c.conn
			c.mu.Unlock()
			if conn != nil {
				_ = conn.Write(ctx, websocket.MessageBinary, buf[:n])
			}
			cancel()
		}
		if err != nil {
			// Port error: report to the server (requirement 3.5 dual-layer
			// detection: client-side COM errors surface explicitly).
			c.send(&protocol.DeviceStatusFrame{Port: c.port, Open: false, Error: err.Error()})
			return
		}
	}
}

// readLoop decodes server frames: control (state/progress) logged; binary
// frames are server->device bytes written into the local port.
func (c *Client) readLoop() {
	defer c.done.Done()
	portRef := func() PortOpener {
		c.mu.Lock()
		defer c.mu.Unlock()
		return c.localPort
	}
	for {
		select {
		case <-c.stop:
			return
		default:
		}
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		typ, data, err := c.conn.Read(ctx)
		cancel()
		if err != nil {
			slog.Warn("server connection lost", "err", err)
			return
		}
		switch typ {
		case websocket.MessageBinary:
			if p := portRef(); p != nil {
				if _, err := p.Write(data); err != nil {
					slog.Error("port write failed", "err", err)
				}
			}
		case websocket.MessageText:
			if f, err := protocol.ParseControl(data); err == nil {
				c.onFrame(f)
			}
		}
	}
}

// onFrame reacts to control frames the client cares about.
func (c *Client) onFrame(f *protocol.Frame) {
	switch b := f.Body.(type) {
	case *protocol.SessionStateFrame:
		slog.Info("session state", "device", b.DeviceID, "state", b.State, "detail", b.Detail)
	case *protocol.ExpectProgressFrame:
		slog.Info("expect progress", "session", b.SessionID, "step", b.StepIndex, "phase", b.Phase)
	case *protocol.DeviceStatusFrame:
		// Echo of our own status report; ignore.
	default:
		slog.Debug("control frame", "type", f.Type)
	}
}

// send marshals and writes a control frame best-effort.
func (c *Client) send(body any) {
	wire, err := protocol.Encode(body)
	if err != nil {
		return
	}
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()
	if conn == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = conn.Write(ctx, websocket.MessageText, wire)
}

// SetLocalPort stores the port reference for the read loop.
func (c *Client) SetLocalPort(p PortOpener) {
	c.mu.Lock()
	c.localPort = p
	c.mu.Unlock()
}

// Close tears down the connection and goroutines.
func (c *Client) Close() {
	close(c.stop)
	c.done.Wait()
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()
	if conn != nil {
		conn.CloseNow()
	}
}

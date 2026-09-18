package main

import (
	"context"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/coder/websocket"

	"github.com/WakkeWang/EmbedFlow/pkg/protocol"
)

// runProbeEcho is the server half of the M1 latency spike (issue #2,
// CEO-12A): it echoes every client binary frame back unchanged and every
// heartbeat with its seq, nothing else -- no store, no sessions, no logs.
// The serial-client's probe mode measures the full-path round trip against
// it: device bytes -> COM -> client -> WS -> this echo -> WS -> client ->
// COM -> device.
func runProbeEcho(addr string) {
	slog.Info("probe-echo listening", "addr", addr)
	srv := &http.Server{Addr: addr, Handler: probeEchoHandler(), ReadHeaderTimeout: 10 * time.Second}
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("probe server", "err", err)
	}
}

func probeEchoHandler() http.Handler {
	// Concurrent-connection cap: the probe mode is unauthenticated by design,
	// but it must not become a free goroutine farm on the production port.
	var conns atomic.Int64
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if conns.Add(1) > maxProbeConns {
			conns.Add(-1)
			http.Error(w, "probe busy", http.StatusServiceUnavailable)
			return
		}
		defer conns.Add(-1)
		runProbeEchoHandler(w, r)
	})
}

const maxProbeConns = 16

func runProbeEchoHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer ws.CloseNow()
	ws.SetReadLimit(1 << 16)
	ctx := r.Context()
	// No token validation in probe mode: the measurement path must add
	// nothing the production path could avoid (the handshake itself is one
	// frame and happens before timing starts, so it does not skew numbers).
	for {
		typ, data, err := ws.Read(ctx)
		if err != nil {
			return
		}
		switch typ {
		case websocket.MessageBinary:
			writeProbe(ws, websocket.MessageBinary, data)
		case websocket.MessageText:
			f, err := protocol.ParseControl(data)
			if err != nil {
				continue
			}
			switch body := f.Body.(type) {
			case *protocol.AuthFrame:
				writeProbe(ws, websocket.MessageText, mustEncodeWire(&protocol.AuthOKFrame{}))
			case *protocol.HeartbeatFrame:
				writeProbe(ws, websocket.MessageText, mustEncodeWire(body))
			}
		}
	}
}

// writeProbe sends one frame with the probe's bounded write timeout; a peer
// that stops reading cannot pin the goroutine past 5s.
func writeProbe(ws *websocket.Conn, typ websocket.MessageType, data []byte) {
	wctx, wcancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer wcancel()
	_ = ws.Write(wctx, typ, data)
}

func mustEncodeWire(body any) []byte {
	wire, err := protocol.Encode(body)
	if err != nil {
		panic(err)
	}
	return wire
}

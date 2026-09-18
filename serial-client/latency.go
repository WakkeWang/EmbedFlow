package client

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/coder/websocket"

	"github.com/WakkeWang/EmbedFlow/pkg/protocol"
)

// Latency probe (issue #2, CEO-12A): measure the full-path serial round trip
// "device bytes -> COM -> client -> WS -> server echo -> WS -> client ->
// COM -> device" against the server's -probe-echo mode. Payloads carry a
// monotonic send counter; the echo returns them untouched, so latency =
// recv time - send time of the matching counter. The report prints p50/p90/
// p99/max: the client-side COM legs are inside the number, the server-side
// matcher leg is not (the echo is the matcher stand-in).
type LatencyResult struct {
	N        int
	P50      time.Duration
	P90      time.Duration
	P99      time.Duration
	Max      time.Duration
	Mean     time.Duration
	Samples  []time.Duration
	TimedOut int
}

// MeasureLatency runs n echo round trips through an already-connected
// client. The serial port must be open: bytes go out the COM port and the
// echo comes back through the same WS the passthrough loop uses.
//
// On a bare loopback plug (TX wired to RX) this measures the whole tunnel;
// against a live uboot console it doubles as the go/no-go signal for the
// server-side expect matcher (issue #2 acceptance).
func MeasureLatency(ctx context.Context, c *Client, p PortOpener, n int, timeout time.Duration) (*LatencyResult, error) {
	if n <= 0 {
		n = 200
	}
	res := &LatencyResult{N: n, Samples: make([]time.Duration, 0, n)}

	// pending maps the 8-byte big-endian counter to its send time. The echo
	// path is ordered per-connection, so FIFO matching would do; the map
	// keeps it correct even if a frame is lost.
	pending := make(map[uint64]time.Time)
	arrivals := make(chan []byte, 256)

	// Tap the inbound bytes: the client's read loop delivers server bytes to
	// the port; we also observe them here via the port read side. To avoid
	// racing the passthrough pump, the probe owns the port exclusively: it
	// reads what comes back itself.
	stop := make(chan struct{})
	defer close(stop)

	go func() {
		buf := make([]byte, 4096)
		for {
			select {
			case <-stop:
				return
			default:
			}
			n, err := p.Read(buf)
			if n > 0 {
				chunk := make([]byte, n)
				copy(chunk, buf[:n])
				select {
				case arrivals <- chunk:
				default:
				}
			}
			if err != nil {
				return
			}
		}
	}()

	var acc []byte
	send := func(counter uint64) error {
		payload := make([]byte, 8)
		binary.BigEndian.PutUint64(payload, counter)
		pending[counter] = time.Now()
		_, err := p.Write(payload)
		return err
	}

	for i := 0; i < n; i++ {
		if err := ctx.Err(); err != nil {
			break
		}
		counter := uint64(i + 1)
		if err := send(counter); err != nil {
			return nil, fmt.Errorf("client: probe write: %w", err)
		}
		deadline := time.After(timeout)
		got := false
		for !got {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-deadline:
				res.TimedOut++
				got = true
			case chunk := <-arrivals:
				acc = append(acc, chunk...)
				for len(acc) >= 8 {
					seq := binary.BigEndian.Uint64(acc[:8])
					acc = acc[8:]
					if t0, ok := pending[seq]; ok {
						rtt := time.Since(t0)
						delete(pending, seq)
						if seq == counter {
							res.Samples = append(res.Samples, rtt)
							got = true
						}
					}
				}
			}
		}
		// Pace the probe so the link never sees back-to-back saturation.
		select {
		case <-time.After(10 * time.Millisecond):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	if len(res.Samples) == 0 {
		return res, nil
	}
	sorted := append([]time.Duration(nil), res.Samples...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	pct := func(q float64) time.Duration {
		idx := int(math.Ceil(q*float64(len(sorted)))) - 1
		if idx < 0 {
			idx = 0
		}
		if idx >= len(sorted) {
			idx = len(sorted) - 1
		}
		return sorted[idx]
	}
	var sum time.Duration
	for _, d := range sorted {
		sum += d
	}
	res.P50, res.P90, res.P99, res.Max = pct(0.50), pct(0.90), pct(0.99), sorted[len(sorted)-1]
	res.Mean = sum / time.Duration(len(sorted))
	return res, nil
}

// ProbeDial connects a raw probe socket to a -probe-echo server: auth-free,
// one dial, returns a minimal send/recv surface for latency runs.
func ProbeDial(ctx context.Context, serverURL string) (*ProbeConn, error) {
	wsURL := replaceScheme(serverURL) + "/ws/client"
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	ws, _, err := websocket.Dial(dialCtx, wsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("client: probe dial: %w", err)
	}
	return &ProbeConn{ws: ws}, nil
}

// ProbeConn is one raw echo connection: it skips the production handshake
// except for the version-free auth ping the echo answers.
type ProbeConn struct {
	ws *websocket.Conn
}

// Handshake exchanges the auth frame so the echo treats the connection as
// established (probe mode skips token validation server-side).
func (p *ProbeConn) Handshake(ctx context.Context) error {
	wire, err := protocol.Encode(&protocol.AuthFrame{ProtocolVersion: protocol.ProtocolVersion, ClientVersion: "latency-probe"})
	if err != nil {
		return err
	}
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := p.ws.Write(ctx2, websocket.MessageText, wire); err != nil {
		return err
	}
	_, _, err = p.ws.Read(ctx2)
	return err
}

// RTT sends one 8-byte payload and waits for its echo, returning the
// round-trip time. Pure WS path: no serial port involved (the WS leg of the
// split, for isolating where the delay lives).
//
// Timing uses a wall-clock deadline, not a ctx timeout on Read: cancelling
// a coder/websocket Read closes the connection, so one slow sample would
// kill every sample after it. Late echoes from an earlier timed-out sample
// are matched by their payload counter and discarded, never mis-attributed
// to the current sample.
func (p *ProbeConn) RTT(ctx context.Context, payload []byte, timeout time.Duration) (time.Duration, error) {
	wctx, wcancel := context.WithTimeout(ctx, 5*time.Second)
	defer wcancel()
	t0 := time.Now()
	if err := p.ws.Write(wctx, websocket.MessageBinary, payload); err != nil {
		return 0, err
	}
	deadline := time.After(timeout)
	for {
		rctx, rcancel := context.WithTimeout(ctx, 500*time.Millisecond)
		typ, data, err := p.ws.Read(rctx)
		rcancel()
		if err != nil {
			select {
			case <-deadline:
				return 0, fmt.Errorf("probe: echo timeout for counter %d", binary.BigEndian.Uint64(payload))
			case <-ctx.Done():
				return 0, ctx.Err()
			default:
				continue // read pacing window elapsed; keep waiting to the deadline
			}
		}
		if typ == websocket.MessageBinary && len(data) == len(payload) && binary.BigEndian.Uint64(data) == binary.BigEndian.Uint64(payload) {
			return time.Since(t0), nil
		}
		// Heartbeat or stale echo: ignore, keep reading.
		select {
		case <-deadline:
			return 0, fmt.Errorf("probe: echo timeout for counter %d", binary.BigEndian.Uint64(payload))
		default:
		}
	}
}

func (p *ProbeConn) Close() { _ = p.ws.CloseNow() }

func replaceScheme(u string) string {
	for i := 0; i+1 < len(u); i++ {
		if u[i] == ':' {
			scheme := u[:i]
			switch scheme {
			case "http":
				return "ws" + u[i:]
			case "https":
				return "wss" + u[i:]
			default:
				return u
			}
		}
	}
	return u
}

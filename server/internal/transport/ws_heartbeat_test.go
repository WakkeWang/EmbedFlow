package transport

import (
	"context"
	"testing"
	"time"

	"github.com/coder/websocket"

	"github.com/WakkeWang/EmbedFlow/pkg/protocol"
)

// Heartbeats must reach the hub: the session kernel's dual-threshold
// liveness (CEO-16A) is driven by them, so an echo-only implementation
// freezes lastHeartbeat and murders live task sessions (review finding).
func TestWS_HeartbeatReachesHub(t *testing.T) {
	hub := &testHub{heartbeats: make(chan heartbeatEvt, 4)}
	srv := startServer(t, hub, defaultAuth())
	tc, _ := dialAndAuth(t, srv, "good-token", protocol.ProtocolVersion)
	defer tc.close(t)
	tc.readFrame(t, 2*time.Second) // AuthOK

	beat, _ := protocol.Encode(&protocol.HeartbeatFrame{Seq: 9})
	if err := tc.ws.Write(context.Background(), websocket.MessageText, beat); err != nil {
		t.Fatalf("write: %v", err)
	}
	select {
	case evt := <-hub.heartbeats:
		if evt.seq != 9 {
			t.Fatalf("seq = %d, want 9", evt.seq)
		}
		if evt.c == nil {
			t.Fatal("heartbeat arrived without connection")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("heartbeat never reached the hub")
	}
}

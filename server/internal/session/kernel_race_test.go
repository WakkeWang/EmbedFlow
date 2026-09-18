package session

import (
	"sync"
	"testing"
	"time"
)

// Concurrent open requests on the same device: exactly one wins, everyone
// else gets a clean busy rejection (design doc Failure modes row 1: the
// device-mutex race must be safe under the kernel lock).
func TestOpen_ConcurrentSameDevice_ExactlyOneWinner(t *testing.T) {
	const users = 16
	k := New(Options{IdleTimeout: defaultIdleTimeout})

	var wg sync.WaitGroup
	results := make(chan OpenResult, users)
	for i := 0; i < users; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results <- k.Open(openReq{
				User:   "u",
				Device: 1,
				Kind:   KindManual,
				Now:    t0.Add(time.Duration(i) * time.Millisecond),
			})
		}(i)
	}
	wg.Wait()
	close(results)

	// Same user is idempotent: all requests map to the same active session.
	// Exactly one distinct session ID may exist, and the device must end up
	// with exactly one active session.
	winners := map[SessionID]bool{}
	for r := range results {
		if r.Rejected {
			t.Fatalf("same-user open must be idempotent, got rejection: %s", r.Reason)
		}
		winners[r.SessionID] = true
	}
	active := 0
	for _, s := range k.ActiveSessions() {
		if s.DeviceID == 1 {
			active++
		}
	}
	if active != 1 {
		t.Fatalf("active sessions on device 1 = %d, want 1", active)
	}
	if len(winners) != 1 {
		t.Fatalf("distinct session IDs = %d, want 1", len(winners))
	}
}

package session

import (
	"strings"
	"testing"
	"time"
)

var t0 = time.Date(2026, 9, 18, 10, 0, 0, 0, time.UTC)

const defaultIdleTimeout = 30 * time.Minute

const idle = "idle"

// ---- result helpers ----

func (r OpenResult) wantState(t *testing.T, state string) {
	t.Helper()
	if r.Rejected {
		t.Fatalf("open rejected: %s", r.Reason)
	}
	if r.State != state {
		t.Fatalf("state = %q, want %q", r.State, state)
	}
}

func (r OpenResult) wantKind(t *testing.T, kind string) {
	t.Helper()
	if r.Kind != kind {
		t.Fatalf("kind = %q, want %q", r.Kind, kind)
	}
}

func (r OpenResult) wantOwner(t *testing.T, user string) {
	t.Helper()
	if r.Owner != user {
		t.Fatalf("owner = %q, want %q", r.Owner, user)
	}
}

func (r OpenResult) wantRejected(t *testing.T, reasonPart string) {
	t.Helper()
	if !r.Rejected {
		t.Fatalf("want rejection, got %+v", r)
	}
	if reasonPart != "" && !strings.Contains(r.Reason, reasonPart) {
		t.Fatalf("reason = %q, want it to contain %q", r.Reason, reasonPart)
	}
}

// ---- kernel helpers ----

func newTestKernel(opts ...Option) *Kernel {
	o := Options{IdleTimeout: defaultIdleTimeout}
	for _, fn := range opts {
		fn(&o)
	}
	return New(o)
}

func withIdleTimeout(d time.Duration) Option {
	return func(o *Options) { o.IdleTimeout = d }
}

func (k *Kernel) wantSession(t *testing.T, id SessionID, state string) {
	t.Helper()
	s, ok := k.Session(id)
	if !ok {
		t.Fatalf("session %d not found", id)
	}
	if s.State != state {
		t.Fatalf("session %d state = %q, want %q", id, s.State, state)
	}
}

func (k *Kernel) wantDevice(t *testing.T, dev DeviceID, wantState string) {
	t.Helper()
	st := k.DeviceState(dev)
	if wantState == idle && st != (DeviceState{}) {
		t.Fatalf("device %d state = %+v, want unoccupied", dev, st)
	}
	if wantState != idle && st != (DeviceState{SessionID: st.SessionID}) {
		t.Fatalf("device %d occupied: %+v, want %q", dev, st, wantState)
	}
}

// ---- event helpers ----

func hasEvent(evts []Event, want Event) bool {
	for _, e := range evts {
		if e.Kind == want.Kind && e.SessionID == want.SessionID {
			return true
		}
	}
	return false
}

func evtSessionClosed(id SessionID) Event {
	return Event{Kind: EventSessionClosed, SessionID: id}
}

func evtSweepFailed(id SessionID) Event {
	return Event{Kind: EventSweepFailed, SessionID: id}
}

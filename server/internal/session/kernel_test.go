package session

import (
	"testing"
	"time"
)

// The kernel state machine is tested as a pure object: events in, state out,
// no IO (decision 1A/6A). Times are injected so no test sleeps.

// --- device mutex & open ---

func TestOpen_FirstManualSession_OccupiesDevice(t *testing.T) {
	k := newTestKernel()
	r := k.Open(openReq{User: "alice", Device: 1, Kind: KindManual, Now: t0})
	r.wantState(t, StateActive)
	r.wantKind(t, KindManual)
	r.wantOwner(t, "alice")
}

func TestOpen_SecondUserWhileBusy_RejectedWithOccupier(t *testing.T) {
	k := newTestKernel()
	k.Open(openReq{User: "alice", Device: 1, Kind: KindManual, Now: t0})
	r := k.Open(openReq{User: "bob", Device: 1, Kind: KindManual, Now: t0.Add(time.Second)})
	r.wantRejected(t, "device busy")
}

func TestOpen_SameUserIdempotent_ReturnsExistingSession(t *testing.T) {
	k := newTestKernel()
	r1 := k.Open(openReq{User: "alice", Device: 1, Kind: KindManual, Now: t0})
	r2 := k.Open(openReq{User: "alice", Device: 1, Kind: KindManual, Now: t0.Add(time.Second)})
	if r2.Rejected || r2.SessionID != r1.SessionID {
		t.Fatalf("idempotent open: got %+v, want session %d", r2, r1.SessionID)
	}
}

func TestOpen_SameUserDifferentKind_NotIdempotent(t *testing.T) {
	// A manual session does not satisfy a task-session request: different
	// lifecycle semantics (decision CEO-15A keys on user+device+kind).
	k := newTestKernel()
	k.Open(openReq{User: "alice", Device: 1, Kind: KindManual, Now: t0})
	r := k.Open(openReq{User: "alice", Device: 1, Kind: KindTask, Now: t0.Add(time.Second)})
	r.wantRejected(t, "device busy")
}

func TestOpen_IndependentDevices_BothSucceed(t *testing.T) {
	k := newTestKernel()
	r1 := k.Open(openReq{User: "alice", Device: 1, Kind: KindManual, Now: t0})
	r2 := k.Open(openReq{User: "bob", Device: 2, Kind: KindManual, Now: t0})
	if r1.Rejected || r2.Rejected {
		t.Fatalf("independent devices should both open: %+v %+v", r1, r2)
	}
}

// --- disconnect branches (decision 1A state diagram) ---

func TestDisconnect_TaskSession_FailsTerminated(t *testing.T) {
	k := newTestKernel()
	r := k.Open(openReq{User: "alice", Device: 1, Kind: KindTask, Now: t0})
	k.ClientDisconnected(r.SessionID, t0.Add(time.Minute))
	k.wantSession(t, r.SessionID, StateFailed)
	k.wantDevice(t, 1, idle) // mutex released
}

func TestDisconnect_ManualSession_ClosesAndReleases(t *testing.T) {
	k := newTestKernel()
	r := k.Open(openReq{User: "alice", Device: 1, Kind: KindManual, Now: t0})
	k.ClientDisconnected(r.SessionID, t0.Add(time.Minute))
	k.wantSession(t, r.SessionID, StateClosed)
	k.wantDevice(t, 1, idle)
}

func TestDisconnect_UnknownSession_Noop(t *testing.T) {
	k := newTestKernel()
	k.ClientDisconnected(999, t0) // must not panic
}

// --- manual session idle timeout (decision 11A, CEO-1A) ---

func TestTick_ManualSessionIdle_ExceedsTimeoutCloses(t *testing.T) {
	k := newTestKernel(withIdleTimeout(30 * time.Minute))
	r := k.Open(openReq{User: "alice", Device: 1, Kind: KindManual, Now: t0})
	evts := k.Tick(t0.Add(31 * time.Minute))
	k.wantSession(t, r.SessionID, StateClosed)
	k.wantDevice(t, 1, idle)
	if !hasEvent(evts, evtSessionClosed(r.SessionID)) {
		t.Fatalf("want evtSessionClosed in %+v", evts)
	}
}

func TestTick_ManualSessionIdle_UnderTimeoutStays(t *testing.T) {
	k := newTestKernel(withIdleTimeout(30 * time.Minute))
	r := k.Open(openReq{User: "alice", Device: 1, Kind: KindManual, Now: t0})
	k.Tick(t0.Add(29 * time.Minute))
	k.wantSession(t, r.SessionID, StateActive)
}

func TestTick_ManualSessionAnyInputResetsIdleTimer(t *testing.T) {
	k := newTestKernel(withIdleTimeout(30 * time.Minute))
	r := k.Open(openReq{User: "alice", Device: 1, Kind: KindManual, Now: t0})
	k.UserInput(r.SessionID, t0.Add(20*time.Minute)) // keepalive
	k.Tick(t0.Add(40 * time.Minute))                 // 20+30 exceeded, but input at 20 reset
	k.wantSession(t, r.SessionID, StateActive)
	k.Tick(t0.Add(51 * time.Minute)) // 20+30 < 51
	k.wantSession(t, r.SessionID, StateClosed)
}

func TestTick_TaskSession_NoIdleTimeout(t *testing.T) {
	// CEO-1A: task sessions have no global idle timer -- flash install can be
	// silent for 10+ minutes; only per-step expect timeouts apply (M3+).
	k := newTestKernel(withIdleTimeout(30 * time.Minute))
	r := k.Open(openReq{User: "alice", Device: 1, Kind: KindTask, Now: t0})
	k.Tick(t0.Add(24 * time.Hour))
	k.wantSession(t, r.SessionID, StateActive)
}

// --- client connection dual-threshold liveness (decision CEO-16A) ---

func TestTick_HeartbeatLost5s_ClientConnectionMarkedDown(t *testing.T) {
	k := newTestKernel()
	r := k.Open(openReq{User: "alice", Device: 1, Kind: KindManual, Now: t0})
	k.ClientHeartbeat(r.SessionID, t0)
	k.Tick(t0.Add(6 * time.Second))
	k.wantSession(t, r.SessionID, StateActive) // detection only: manual session not yet torn down
}

func TestTick_TaskSession_HeartbeatLostPastGrace_Fails(t *testing.T) {
	k := newTestKernel()
	r := k.Open(openReq{User: "alice", Device: 1, Kind: KindTask, Now: t0})
	k.ClientHeartbeat(r.SessionID, t0)
	// 30s grace window (reconnect tolerance) after the 5s detection threshold.
	k.Tick(t0.Add(6 * time.Second))
	k.wantSession(t, r.SessionID, StateActive)
	k.Tick(t0.Add(36 * time.Second))
	k.wantSession(t, r.SessionID, StateFailed)
	k.wantDevice(t, 1, idle)
}

func TestTick_TaskSession_HeartbeatRecoversWithinGrace_Continues(t *testing.T) {
	k := newTestKernel()
	r := k.Open(openReq{User: "alice", Device: 1, Kind: KindTask, Now: t0})
	k.ClientHeartbeat(r.SessionID, t0)
	k.Tick(t0.Add(6 * time.Second))                        // detected down
	k.ClientHeartbeat(r.SessionID, t0.Add(20*time.Second)) // back within 30s grace
	k.Tick(t0.Add(40 * time.Second))
	k.wantSession(t, r.SessionID, StateActive)
}

// --- startup sweep (decision CEO-7A) ---

func TestSweepStartup_LegacyActiveSessions_FailedWithReason(t *testing.T) {
	k := newTestKernel()
	mr := k.Open(openReq{User: "alice", Device: 1, Kind: KindManual, Now: t0})
	tr := k.Open(openReq{User: "bob", Device: 2, Kind: KindTask, Now: t0})
	m, _ := k.Session(mr.SessionID)
	task, _ := k.Session(tr.SessionID)
	k2 := New(Options{IdleTimeout: defaultIdleTimeout})
	k2.Seed(m, task)
	evts := k2.SweepStartup("server restart", t0.Add(time.Hour))
	k2.wantSession(t, mr.SessionID, StateFailed)
	k2.wantSession(t, tr.SessionID, StateFailed)
	if !hasEvent(evts, evtSweepFailed(mr.SessionID)) {
		t.Fatalf("want sweep event for %d in %+v", mr.SessionID, evts)
	}
	k2.wantDevice(t, 1, idle)
	k2.wantDevice(t, 2, idle)
}

func TestSweepStartup_FinishedSessions_Untouched(t *testing.T) {
	k := newTestKernel()
	mr := k.Open(openReq{User: "alice", Device: 1, Kind: KindManual, Now: t0})
	k.Close(mr.SessionID, t0.Add(time.Minute))
	m, _ := k.Session(mr.SessionID)
	k2 := New(Options{IdleTimeout: defaultIdleTimeout})
	k2.Seed(m)
	k2.SweepStartup("server restart", t0.Add(time.Hour))
	k2.wantSession(t, mr.SessionID, StateClosed) // stays closed, not overwritten
}

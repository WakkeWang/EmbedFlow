// Package session implements the EmbedFlow session domain kernel: data model
// and state machine, pure Go with no IO (decisions 1A/6A). All times are
// injected by the caller; the kernel never reads the wall clock. Side effects
// (logging, persistence, notifications) are the executor layer's job, driven
// by the events the kernel returns.
package session

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// SessionID identifies a session record.
type SessionID int64

// DeviceID identifies a registered device.
type DeviceID int64

// Session kinds (decision 1A: sessions are the serial persistent entity;
// SSH command executions are modeled elsewhere as short transactions).
const (
	KindManual = "manual" // tester opens a terminal to poke the device
	KindTask   = "task"   // flash/test run drives the device via expect
)

// Session terminal and active states.
const (
	StateActive = "active"
	StateClosed = "closed" // manual session ended normally
	StateFailed = "failed" // task session failed / aborted / swept
)

// Failure reasons carry the why in session records.
const (
	ReasonDisconnect    = "connection interrupted"
	ReasonHeartbeatGone = "client connection lost"
	ReasonIdleTimeout   = "idle timeout"
	ReasonManualClose   = "closed by user"
)

// IdleTimeoutDefault is the default manual-session idle timeout (decision 11A).
const IdleTimeoutDefault = 30 * time.Minute

// Liveness threshold defaults (decision CEO-16A): heartbeat lost 5s marks the
// client connection down (detection); a task session is killed only after a
// further 30s reconnect grace window so a WiFi blip does not murder a flash
// run. Both are configurable via Options.
const (
	HeartbeatLostAfterDefault = 5 * time.Second
	TaskGraceWindowDefault    = 30 * time.Second
)

// OpenResult is the outcome of an open request.
type OpenResult struct {
	SessionID SessionID
	State     string
	Kind      string
	Owner     string
	Rejected  bool
	Reason    string
	// Existed is true when Open returned via the idempotent path: the
	// session already existed (same user+device+kind) and nothing new was
	// created -- executors must not persist it again.
	Existed bool
	// Occupier is set on busy-rejection: who holds the device, what kind of
	// session, and since when (decision DS-2A wants this in the popup).
	Occupier *OccupierInfo
}

// OccupierInfo describes the current holder of a busy device.
type OccupierInfo struct {
	User      string    `json:"user"`
	Kind      string    `json:"kind"`
	Since     time.Time `json:"since"`
	SessionID SessionID `json:"session_id"`
}

// Session is the kernel's session record.
type Session struct {
	ID        SessionID
	DeviceID  DeviceID
	Kind      string
	Owner     string
	State     string
	StartedAt time.Time
	EndedAt   time.Time
	EndReason string

	// lastActivity resets the manual idle timer on any input (issue 5).
	lastActivity time.Time
	// lastHeartbeat tracks client liveness for the dual-threshold check.
	lastHeartbeat time.Time
	// connDownAt is when the heartbeat was first seen as lost; zero while the
	// connection is considered up. Past TaskGraceWindow, a task session dies.
	connDownAt time.Time
}

// DeviceState is what the device list shows.
type DeviceState struct {
	SessionID SessionID
	Kind      string
	Owner     string
	Since     time.Time
}

// Event is a side-effect the executor layer must perform after a transition.
type Event struct {
	Kind      string
	SessionID SessionID
	DeviceID  DeviceID
	Detail    string
	At        time.Time
}

// Event kinds emitted by the kernel.
const (
	EventSessionOpened = "session_opened"
	EventSessionClosed = "session_closed"
	EventSessionFailed = "session_failed"
	EventConnDown      = "conn_down"    // heartbeat lost 5s (detection level)
	EventConnWarning   = "conn_warning" // task grace window countdown
	EventConnRecovered = "conn_recovered"
	EventSweepFailed   = "sweep_failed"
)

// Options configures the kernel. Zero-value fields take the design-doc
// defaults; the server exposes them as flags (CEO-16A: both liveness
// parameters configurable).
type Options struct {
	IdleTimeout        time.Duration
	HeartbeatLostAfter time.Duration
	TaskGraceWindow    time.Duration
}

// Option mutates Options (used by tests for targeted overrides).
type Option func(*Options)

// WithIdleTimeout overrides the manual-session idle timeout.
func WithIdleTimeout(d time.Duration) Option { return func(o *Options) { o.IdleTimeout = d } }

// Kernel is the pure session state machine. It is safe for concurrent use:
// device mutex decisions happen under the kernel's own lock (decisions CEO-2A,
// CEO-15A: lock-then-check-then-create).
type Kernel struct {
	mu                 sync.Mutex
	idleTimeout        time.Duration
	heartbeatLostAfter time.Duration
	taskGraceWindow    time.Duration
	nextID             SessionID
	sessions           map[SessionID]*Session
	devices            map[DeviceID]*DeviceState
}

// New creates an empty kernel.
func New(opts Options) *Kernel {
	if opts.IdleTimeout <= 0 {
		opts.IdleTimeout = IdleTimeoutDefault
	}
	if opts.HeartbeatLostAfter <= 0 {
		opts.HeartbeatLostAfter = HeartbeatLostAfterDefault
	}
	if opts.TaskGraceWindow <= 0 {
		opts.TaskGraceWindow = TaskGraceWindowDefault
	}
	return &Kernel{
		idleTimeout:        opts.IdleTimeout,
		heartbeatLostAfter: opts.HeartbeatLostAfter,
		taskGraceWindow:    opts.TaskGraceWindow,
		sessions:           map[SessionID]*Session{},
		devices:            map[DeviceID]*DeviceState{},
	}
}

// OpenRequest asks the kernel to start a session.
type OpenRequest struct {
	User   string
	Device DeviceID
	Kind   string
	Now    time.Time
}

// Open attempts to start a session. Idempotency key = user+device+kind
// (decision CEO-15A as amended: manual and task sessions differ in lifecycle
// semantics, so one kind never satisfies the other; repeats of the same kind
// return the active session; a different user is rejected busy with occupier
// details).
func (k *Kernel) Open(req OpenRequest) OpenResult {
	k.mu.Lock()
	defer k.mu.Unlock()

	// Idempotent path: same user, same device, same kind, still active.
	for id, s := range k.sessions {
		if s.DeviceID == req.Device && s.Owner == req.User && s.Kind == req.Kind && s.State == StateActive {
			return OpenResult{SessionID: id, State: s.State, Kind: s.Kind, Owner: s.Owner, Existed: true}
		}
	}

	// Device mutex: one active session per device.
	if occ, busy := k.devices[req.Device]; busy {
		return OpenResult{
			Rejected: true,
			Reason:   "device busy",
			Occupier: &OccupierInfo{User: occ.Owner, Kind: occ.Kind, Since: occ.Since, SessionID: occ.SessionID},
		}
	}

	k.nextID++
	id := k.nextID
	s := &Session{
		ID:            id,
		DeviceID:      req.Device,
		Kind:          req.Kind,
		Owner:         req.User,
		State:         StateActive,
		StartedAt:     req.Now,
		lastActivity:  req.Now,
		lastHeartbeat: req.Now,
	}
	k.sessions[id] = s
	k.devices[req.Device] = &DeviceState{SessionID: id, Kind: req.Kind, Owner: req.User, Since: req.Now}
	return OpenResult{SessionID: id, State: s.State, Kind: s.Kind, Owner: s.Owner}
}

// Session returns a copy of a session record.
func (k *Kernel) Session(id SessionID) (Session, bool) {
	k.mu.Lock()
	defer k.mu.Unlock()
	s, ok := k.sessions[id]
	if !ok {
		return Session{}, false
	}
	return *s, true
}

// DeviceState returns the active occupier of a device, if any; empty zero
// value means the device is free.
func (k *Kernel) DeviceState(dev DeviceID) DeviceState {
	k.mu.Lock()
	defer k.mu.Unlock()
	st, ok := k.devices[dev]
	if !ok {
		return DeviceState{}
	}
	return *st
}

// UserInput records activity on a session, resetting the manual idle timer.
func (k *Kernel) UserInput(id SessionID, now time.Time) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if s, ok := k.sessions[id]; ok && s.State == StateActive {
		s.lastActivity = now
	}
}

// ClientHeartbeat records a liveness ping, restoring a detected-down
// connection (within the task grace window this continues the session).
func (k *Kernel) ClientHeartbeat(id SessionID, now time.Time) {
	k.mu.Lock()
	defer k.mu.Unlock()
	s, ok := k.sessions[id]
	if !ok || s.State != StateActive {
		return
	}
	if !s.connDownAt.IsZero() {
		s.connDownAt = time.Time{}
	}
	s.lastHeartbeat = now
}

// ClientDisconnected handles an explicit disconnect: task sessions fail
// terminated, manual sessions close; both release the device mutex (decision
// 1A state diagram, requirement 3.5).
func (k *Kernel) ClientDisconnected(id SessionID, now time.Time) []Event {
	k.mu.Lock()
	defer k.mu.Unlock()
	s, ok := k.sessions[id]
	if !ok || s.State != StateActive {
		return nil
	}
	return k.terminate(s, StateForKind(s.Kind), ReasonDisconnect, now)
}

// StateForKind maps kind to the terminal state an abnormal end produces:
// task sessions fail, manual sessions close (they end, not fail).
func StateForKind(kind string) string {
	if kind == KindTask {
		return StateFailed
	}
	return StateClosed
}

// Close ends an active session at the user's request.
func (k *Kernel) Close(id SessionID, now time.Time) []Event {
	k.mu.Lock()
	defer k.mu.Unlock()
	s, ok := k.sessions[id]
	if !ok || s.State != StateActive {
		return nil
	}
	return k.terminate(s, StateClosed, ReasonManualClose, now)
}

// Tick advances time-driven transitions: manual idle timeout and the
// dual-threshold client-liveness check. Returns the side-effect events.
func (k *Kernel) Tick(now time.Time) []Event {
	k.mu.Lock()
	defer k.mu.Unlock()

	var evts []Event
	for _, id := range k.activeIDsLocked() {
		s := k.sessions[id]

		// Manual idle timeout (decision 11A): no input for the timeout closes
		// the session and frees the device. Task sessions never idle out
		// (decision CEO-1A).
		if s.Kind == KindManual && now.Sub(s.lastActivity) >= k.idleTimeout {
			evts = append(evts, k.terminate(s, StateClosed, ReasonIdleTimeout, now)...)
			continue
		}

		// Dual-threshold liveness (decision CEO-16A).
		hbDown := now.Sub(s.lastHeartbeat) >= k.heartbeatLostAfter
		switch {
		case hbDown && s.connDownAt.IsZero():
			// Detection level: mark down, notify; do not tear anything yet.
			s.connDownAt = now
			evts = append(evts, Event{Kind: EventConnDown, SessionID: s.ID, DeviceID: s.DeviceID, Detail: ReasonHeartbeatGone, At: now})

		case !hbDown && !s.connDownAt.IsZero():
			// Recovered within the grace window.
			s.connDownAt = time.Time{}
			evts = append(evts, Event{Kind: EventConnRecovered, SessionID: s.ID, DeviceID: s.DeviceID, At: now})

		case hbDown && !s.connDownAt.IsZero():
			// Still down: warn inside the grace window, kill past it -- but
			// only task sessions die from heartbeat loss; a manual session's
			// transport drop is handled by explicit disconnect or idle timeout.
			elapsed := now.Sub(s.connDownAt)
			if s.Kind == KindTask {
				if elapsed >= k.taskGraceWindow {
					evts = append(evts, k.terminate(s, StateFailed, ReasonHeartbeatGone, now)...)
				} else {
					evts = append(evts, Event{Kind: EventConnWarning, SessionID: s.ID, DeviceID: s.DeviceID, Detail: (k.taskGraceWindow - elapsed).String(), At: now})
				}
			}
		}
	}
	return evts
}

// activeIDsLocked lists active session IDs in deterministic order. Callers
// hold k.mu.
func (k *Kernel) activeIDsLocked() []SessionID {
	ids := make([]SessionID, 0, len(k.sessions))
	for id, s := range k.sessions {
		if s.State == StateActive {
			ids = append(ids, id)
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}

// SweepStartup marks every active session found in storage as failed with the
// given reason (decision CEO-7A: server restart sweep), releasing all devices.
func (k *Kernel) SweepStartup(reason string, now time.Time) []Event {
	k.mu.Lock()
	defer k.mu.Unlock()

	var evts []Event
	for _, id := range k.activeIDsLocked() {
		s := k.sessions[id]
		evts = append(evts, k.terminate(s, StateFailed, reason, now)...)
		evts = append(evts, Event{Kind: EventSweepFailed, SessionID: s.ID, DeviceID: s.DeviceID, Detail: reason, At: now})
	}
	return evts
}

// terminate ends an active session under lock, releasing the device mutex.
func (k *Kernel) terminate(s *Session, state, reason string, now time.Time) []Event {
	s.State = state
	s.EndedAt = now
	s.EndReason = reason
	delete(k.devices, s.DeviceID)
	kind := EventSessionClosed
	if state == StateFailed {
		kind = EventSessionFailed
	}
	return []Event{{Kind: kind, SessionID: s.ID, DeviceID: s.DeviceID, Detail: reason, At: now}}
}

// Seed loads session records from storage at startup. Active ones keep their
// state in the kernel so SweepStartup can find and fail them (CEO-7A); the
// caller must have applied the startup sweep to storage first.
func (k *Kernel) Seed(sessions ...Session) {
	k.mu.Lock()
	defer k.mu.Unlock()
	for i := range sessions {
		s := sessions[i]
		if s.ID > k.nextID {
			k.nextID = s.ID
		}
		cp := s
		k.sessions[s.ID] = &cp
		if s.State == StateActive {
			k.devices[s.DeviceID] = &DeviceState{SessionID: s.ID, Kind: s.Kind, Owner: s.Owner, Since: s.StartedAt}
		}
	}
}

// SeedCounter advances the ID counter past n without registering any
// session -- the startup sweep uses it so new IDs never collide with
// persisted history rows.
func (k *Kernel) SeedCounter(n SessionID) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if n > k.nextID {
		k.nextID = n
	}
}

// ActiveSessions returns all active sessions (for persistence sync).
func (k *Kernel) ActiveSessions() []Session {
	k.mu.Lock()
	defer k.mu.Unlock()
	out := make([]Session, 0, len(k.sessions))
	for _, s := range k.sessions {
		if s.State == StateActive {
			out = append(out, *s)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// String renders an event compactly for logs.
func (e Event) String() string {
	return fmt.Sprintf("%s session=%d device=%d %s", e.Kind, e.SessionID, e.DeviceID, e.Detail)
}

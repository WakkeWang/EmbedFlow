// Package batch implements the build-batch scheduling kernel: data model
// and state machine, pure Go with no IO (the session kernel's discipline,
// decisions 1A/6A applied to M2). Callers inject time; the kernel returns
// events and the executor layer performs them.
//
// Semantics (requirement 2.3, as pinned in the M2 plan):
//   - Up to N batches run in parallel (system setting, default 4); further
//     batches queue and are admitted when a slot frees.
//   - Within a batch, items build one at a time in topological order
//     (user decision 2026-09-19): global concurrency comes from parallel
//     batches, not intra-batch parallelism.
//   - The build set is the selected items plus their prerequisite closure
//     (cross-batch reuse never happens, so prerequisites must build in this
//     batch; shared prerequisites build once).
//   - Prereq groups: any item in a group succeeding satisfies the group
//     (OR); all groups must be satisfied (AND). A dependent whose group has
//     no success and no pending/building item left is skipped, recursively.
package batch

import (
	"fmt"
	"sort"
	"sync"
)

// ID types are distinct so session and batch ids can never be confused.
type (
	BatchID     int64
	ItemID      int64
	RecordID    int64
)

// Batch statuses.
const (
	BatchQueued    = "queued"
	BatchRunning   = "running"
	BatchCompleted = "completed"
	BatchFailed    = "failed"
	BatchCanceled  = "canceled"
)

// Record statuses (requirement 2.4's seven states; queued = pending).
const (
	RecPending   = "pending"
	RecBuilding  = "building"
	RecSucceeded = "succeeded"
	RecFailed    = "failed"
	RecSkipped   = "skipped" // skipped - prerequisite failure
	RecCanceled  = "canceled"
)

// Outcome is what the executor reports when a build finishes.
type Outcome string

const (
	OutcomeSucceeded Outcome = "succeeded"
	OutcomeFailed    Outcome = "failed"
)

// Item is the kernel's view of a build item: identity plus prerequisite
// groups. PrereqGroups[[a,b],[c]] = some of {a,b} AND some of {c}.
type Item struct {
	ID           ItemID
	PrereqGroups [][]ItemID
}

// Plan is the batch input: which items were selected plus the full item
// universe (the kernel expands the prerequisite closure itself).
type Plan struct {
	Selected []ItemID
	Items    []Item
}

// Record is one build execution inside a batch.
type Record struct {
	RecordID RecordID
	ItemID   ItemID
	Status   string
}

// Batch is the kernel's batch record.
type Batch struct {
	ID      BatchID
	Status  string
	Records []Record // in planning order (topological)
}

// Event kinds emitted by the kernel.
const (
	EventBatchAdmitted = "batch_admitted" // a queued batch got a run slot
	EventRecordStart   = "record_start"   // executor must run this build
	EventRecordSkip    = "record_skip"    // executor must mark skipped
	EventBatchDone     = "batch_done"     // batch reached a terminal state
)

// Event is a side effect the executor layer must perform.
type Event struct {
	Kind     string
	BatchID  BatchID
	RecordID RecordID
	ItemID   ItemID
	// Terminal is the batch's final status on EventBatchDone.
	Terminal string
	Detail   string
}

// Options configures the kernel.
type Options struct {
	// MaxParallel caps concurrently running batches (default 4).
	MaxParallel int
}

// Kernel is the pure batch scheduler. Safe for concurrent use; all
// transitions happen under one lock (the session kernel's pattern).
type Kernel struct {
	mu          sync.Mutex
	maxParallel int
	running     int

	nextBatch  BatchID
	nextRecord RecordID

	// queue holds admitted-order batch ids waiting for a run slot.
	queue []BatchID
	// batches holds every live (non-terminal) batch the kernel knows.
	batches map[BatchID]*Batch
	// plans holds each batch's item graph for closure/skip computation.
	plans map[BatchID]map[ItemID]Item
	// selected holds each batch's originally selected item ids (for
	// reporting which records were explicit vs pulled in as prerequisites).
	selected map[BatchID]map[ItemID]bool
}

// New creates an empty kernel.
func New(opts Options) *Kernel {
	if opts.MaxParallel <= 0 {
		opts.MaxParallel = 4
	}
	return &Kernel{
		maxParallel: opts.MaxParallel,
		batches:     map[BatchID]*Batch{},
		plans:       map[BatchID]map[ItemID]Item{},
		selected:    map[BatchID]map[ItemID]bool{},
	}
}

// NextIDs seeds the counters from persisted state (startup sweep parity
// with session.SeedCounter: kernel ids must not collide with history rows).
func (k *Kernel) NextIDs(batch, record int64) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if BatchID(batch) > k.nextBatch {
		k.nextBatch = BatchID(batch)
	}
	if RecordID(record) > k.nextRecord {
		k.nextRecord = RecordID(record)
	}
}

// MaxParallel returns the current concurrent-batch cap (the admin page's
// effective-value display).
func (k *Kernel) MaxParallel() int {
	k.mu.Lock()
	defer k.mu.Unlock()
	return k.maxParallel
}

// SetMaxParallel adjusts the concurrent-batch cap at runtime (the system
// setting "max_parallel_batches": changes apply to admissions from now on;
// already-running batches are not affected). Values <= 0 are ignored so a
// corrupt setting cannot zero the scheduler. Raising the cap may admit
// queued batches immediately; the returned events must be applied by the
// caller (EventBatchAdmitted / EventRecordStart).
func (k *Kernel) SetMaxParallel(n int) []Event {
	if n <= 0 {
		return nil
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	k.maxParallel = n
	return k.pumpLocked()
}

// enqueuePlan computes the build set for a batch: selected items plus the
// prerequisite closure, deduplicated, ordered so prerequisites come before
// dependents (Kahn's algorithm with insertion-order tie-breaking).
func enqueuePlan(plan Plan, selected []ItemID) []ItemID {
	byID := make(map[ItemID]Item, len(plan.Items))
	for _, it := range plan.Items {
		byID[it.ID] = it
	}

	// Closure: BFS over prerequisites from the selected set.
	inClosure := make(map[ItemID]bool)
	var walk func(id ItemID)
	walk = func(id ItemID) {
		if inClosure[id] {
			return
		}
		inClosure[id] = true
		if it, ok := byID[id]; ok {
			for _, group := range it.PrereqGroups {
				for _, pid := range group {
					walk(pid)
				}
			}
		}
	}
	for _, id := range selected {
		walk(id)
	}

	// Prerequisite edges: item -> each prereq id that is in the closure.
	// (A prereq id outside the closure cannot happen -- closure includes
	// them all -- but a prereq id not in the universe is ignored here; the
	// executor surfaces it as a failed lookup. Selection validation happens
	// before the kernel sees the plan.)
	deps := make(map[ItemID][]ItemID, len(inClosure)) // item -> prerequisites
	rdeps := make(map[ItemID][]ItemID)                // prereq -> dependents
	for id := range inClosure {
		it, ok := byID[id]
		if !ok {
			continue
		}
		seen := make(map[ItemID]bool)
		for _, group := range it.PrereqGroups {
			for _, pid := range group {
				if !inClosure[pid] || seen[pid] {
					continue
				}
				seen[pid] = true
				deps[id] = append(deps[id], pid)
				rdeps[pid] = append(rdeps[pid], id)
			}
		}
	}

	// Kahn's algorithm; cycle members never get in-degree 0 and are dropped
	// (validation catches cycles earlier; the scheduler just stays total).
	var order []ItemID
	// Deterministic tie-break: ascending item id.
	ids := make([]ItemID, 0, len(inClosure))
	for id := range inClosure {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	// Seed with items that have no prerequisites inside the closure.
	degree := make(map[ItemID]int, len(ids))
	for _, id := range ids {
		degree[id] = len(deps[id])
	}
	ready := make([]ItemID, 0)
	for _, id := range ids {
		if degree[id] == 0 {
			ready = append(ready, id)
		}
	}
	for len(ready) > 0 {
		// Pop the smallest ready id (stable, deterministic order).
		sort.Slice(ready, func(i, j int) bool { return ready[i] < ready[j] })
		id := ready[0]
		ready = ready[1:]
		order = append(order, id)
		for _, dep := range rdeps[id] {
			degree[dep]--
			if degree[dep] == 0 {
				ready = append(ready, dep)
			}
		}
	}
	return order
}

// Enqueue registers a batch and either admits it immediately (a run slot
// is free) or parks it in the queue. The returned events tell the executor
// which records to persist/start: EventRecordStart for the first buildable
// item when admitted.
func (k *Kernel) Enqueue(plan Plan, selected []ItemID) (BatchID, []Event) {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.nextBatch++
	id := k.nextBatch
	order := enqueuePlan(plan, selected)

	b := &Batch{ID: id, Status: BatchQueued}
	selSet := make(map[ItemID]bool, len(selected))
	for _, s := range selected {
		selSet[s] = true
	}
	k.selected[id] = selSet

	pgraph := make(map[ItemID]Item, len(plan.Items))
	for _, it := range plan.Items {
		pgraph[it.ID] = it
	}
	k.plans[id] = pgraph

	for _, itemID := range order {
		k.nextRecord++
		b.Records = append(b.Records, Record{RecordID: k.nextRecord, ItemID: itemID, Status: RecPending})
	}
	k.batches[id] = b

	// Admission is implicit: an empty queue + free slot = this batch runs.
	// (Queued batches ahead of this one, if any, get their slot first via
	// pump below.)
	k.queue = append(k.queue, id)
	events := k.pumpLocked()
	return id, events
}

// PlanRecords returns the batch's planned records in execution order
// (record id, item id, pending status). The executor persists these rows
// before applying any Enqueue events: RecordStart refers to a record the
// store must already contain.
func (k *Kernel) PlanRecords(id BatchID) []Record {
	k.mu.Lock()
	defer k.mu.Unlock()
	b, ok := k.batches[id]
	if !ok {
		return nil
	}
	out := make([]Record, len(b.Records))
	copy(out, b.Records)
	return out
}

// pumpLocked admits queued batches while slots are free and returns the
// cumulative events (admission + first record start). Caller holds k.mu.
func (k *Kernel) pumpLocked() []Event {
	var events []Event
	for k.running < k.maxParallel && len(k.queue) > 0 {
		id := k.queue[0]
		k.queue = k.queue[1:]
		b := k.batches[id]
		b.Status = BatchRunning
		k.running++
		events = append(events, Event{Kind: EventBatchAdmitted, BatchID: id})
		if start := k.firstStartableLocked(b, k.plans[id]); start != nil {
			events = append(events, *start)
		}
	}
	return events
}

// firstStartableLocked picks the batch's next buildable record: the first
// pending record whose prerequisites are all satisfied, or -- if some
// prerequisite has reached a terminal-unsuccessful state -- the first
// skippable record. Returns nil when the batch has nothing to start right
// now (something is building, or only skipped/pending-behind items remain).
// Caller holds k.mu.
func (k *Kernel) firstStartableLocked(b *Batch, graph map[ItemID]Item) *Event {
	for i := range b.Records {
		rec := &b.Records[i]
		if rec.Status != RecPending {
			continue
		}
		item, ok := graph[rec.ItemID]
		if !ok {
			// Universe lost the item (deleted mid-batch): fail it rather
			// than stall the batch.
			return k.skipRecordLocked(b, rec, "build item deleted")
		}
		ready := true
		for _, group := range item.PrereqGroups {
			satisfied := false
			livePending := false
			for _, pid := range group {
				st := k.statusOfLocked(b, pid)
				switch st {
				case RecSucceeded:
					satisfied = true
				case RecPending, RecBuilding:
					livePending = true
				}
			}
			if !satisfied && !livePending {
				// The group is dead: nothing in it succeeded and nothing
				// will. This item can never run -> skipped (requirement
				// 2.3: prerequisite failure skips the dependent).
				return k.skipRecordLocked(b, rec, "prerequisite group failed")
			}
			if !satisfied {
				ready = false // group still has live hope; wait behind it
			}
		}
		if ready {
			rec.Status = RecBuilding
			return &Event{Kind: EventRecordStart, BatchID: b.ID, RecordID: rec.RecordID, ItemID: rec.ItemID}
		}
	}
	return nil
}

// statusOfLocked looks up an item's record status in the batch. An id not
// in the closure maps to RecSkipped -- for prereq gating that is exactly
// the "dead" state: it can never succeed and never block an OR group.
// Caller holds k.mu.
func (k *Kernel) statusOfLocked(b *Batch, id ItemID) string {
	for i := range b.Records {
		if b.Records[i].ItemID == id {
			return b.Records[i].Status
		}
	}
	return RecSkipped // item not in the closure: treat as dead for gating
}

// skipRecordLocked marks a record skipped and emits the event.
// Caller holds k.mu.
func (k *Kernel) skipRecordLocked(b *Batch, rec *Record, reason string) *Event {
	rec.Status = RecSkipped
	return &Event{Kind: EventRecordSkip, BatchID: b.ID, RecordID: rec.RecordID, ItemID: rec.ItemID, Detail: reason}
}

// NotifyRecordDone reports a build's outcome and advances the batch: skip
// propagation, next record start, batch terminal state, and queue pumping
// when a slot frees. Events describe everything the executor must apply.
func (k *Kernel) NotifyRecordDone(batchID BatchID, recordID RecordID, outcome Outcome) []Event {
	k.mu.Lock()
	defer k.mu.Unlock()

	b, ok := k.batches[batchID]
	if !ok {
		return nil
	}
	for i := range b.Records {
		if b.Records[i].RecordID == recordID {
			switch outcome {
			case OutcomeSucceeded:
				b.Records[i].Status = RecSucceeded
			default:
				b.Records[i].Status = RecFailed
			}
			break
		}
	}
	return k.advanceLocked(b)
}

// advanceLocked runs the post-completion progression: skip cascades and
// next starts until stable (a skip can unlock the next skip down the chain),
// then batch terminal state and slot pump. Caller holds k.mu.
func (k *Kernel) advanceLocked(b *Batch) []Event {
	var events []Event
	for {
		start := k.firstStartableLocked(b, k.plans[b.ID])
		if start == nil {
			break
		}
		events = append(events, *start)
		if start.Kind != EventRecordSkip {
			break // a build is running; the batch advances on its completion
		}
	}
	if done, terminal := batchTerminalLocked(b); done {
		b.Status = terminal
		events = append(events, Event{Kind: EventBatchDone, BatchID: b.ID, Terminal: terminal})
		// Terminal batches leave kernel memory immediately: a late build
		// result (the cancel race) must find no batch to write into.
		k.forgetLocked(b.ID)
		// Free the slot and admit whoever waits.
		k.running--
		events = append(events, k.pumpLocked()...)
	}
	return events
}

// forgetLocked drops a batch from kernel memory. Caller holds k.mu.
func (k *Kernel) forgetLocked(id BatchID) {
	delete(k.batches, id)
	delete(k.plans, id)
	delete(k.selected, id)
}

// batchTerminalLocked reports whether every record reached a terminal state
// and what the batch's final status is (any failed record -> failed).
// Caller holds k.mu.
func batchTerminalLocked(b *Batch) (bool, string) {
	failed := false
	for i := range b.Records {
		switch b.Records[i].Status {
		case RecSucceeded, RecFailed, RecSkipped, RecCanceled:
			if b.Records[i].Status == RecFailed {
				failed = true
			}
		default:
			return false, ""
		}
	}
	if failed {
		return true, BatchFailed
	}
	return true, BatchCompleted
}

// Cancel stops a batch: the running record is reported for cancellation,
// everything still pending becomes canceled, and the batch terminates.
// Records already succeeded/failed/skipped keep their state (requirement:
// 已跑完的保留原状态).
func (k *Kernel) Cancel(batchID BatchID) []Event {
	k.mu.Lock()
	defer k.mu.Unlock()

	b, ok := k.batches[batchID]
	if !ok {
		return nil
	}
	wasRunning := b.Status == BatchRunning
	var events []Event
	var runningRecord RecordID
	for i := range b.Records {
		switch b.Records[i].Status {
		case RecBuilding:
			b.Records[i].Status = RecCanceled
			runningRecord = b.Records[i].RecordID
			events = append(events, Event{Kind: EventRecordSkip, BatchID: batchID, RecordID: b.Records[i].RecordID, ItemID: b.Records[i].ItemID, Detail: RecCanceled})
		case RecPending:
			b.Records[i].Status = RecCanceled
			events = append(events, Event{Kind: EventRecordSkip, BatchID: batchID, RecordID: b.Records[i].RecordID, ItemID: b.Records[i].ItemID, Detail: RecCanceled})
		}
	}
	// runningRecord marks the record whose build was executing; the
	// executor's cancel flag (set before this call) is what actually stops
	// it, and the orchestrator drops its late result via kernel.Active.
	_ = runningRecord
	b.Status = BatchCanceled
	events = append(events, Event{Kind: EventBatchDone, BatchID: batchID, Terminal: BatchCanceled})
	k.forgetLocked(batchID)
	if wasRunning {
		k.running--
		events = append(events, k.pumpLocked()...)
	} else {
		// Queued batch: remove from the wait queue.
		for i, id := range k.queue {
			if id == batchID {
				k.queue = append(k.queue[:i], k.queue[i+1:]...)
				break
			}
		}
	}
	return events
}

// Active reports whether the kernel still tracks the batch (used by the
// executor to decide if a late build result should be applied).
func (k *Kernel) Active(batchID BatchID) bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	_, ok := k.batches[batchID]
	return ok
}

// String renders an event compactly for logs.
func (e Event) String() string {
	return fmt.Sprintf("%s batch=%d record=%d item=%d %s", e.Kind, e.BatchID, e.RecordID, e.ItemID, e.Detail)
}

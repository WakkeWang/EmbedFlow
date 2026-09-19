package batch

import (
	"testing"
)

// The batch kernel is tested as a pure object: calls in, events out, no IO
// (the session kernel's test discipline). No test sleeps.

func item(id int64, groups ...[]ItemID) Item {
	it := Item{ID: ItemID(id)}
	it.PrereqGroups = append(it.PrereqGroups, groups...)
	return it
}

func ids(vs ...int64) []ItemID {
	out := make([]ItemID, 0, len(vs))
	for _, v := range vs {
		out = append(out, ItemID(v))
	}
	return out
}

// eventsOf flattens helper: run and collect event kinds in order.
func kinds(evts []Event) []string {
	out := make([]string, 0, len(evts))
	for _, e := range evts {
		out = append(out, e.Kind)
	}
	return out
}

func wantKinds(t *testing.T, got []string, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("event kinds = %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("event kinds = %v, want %v", got, want)
		}
	}
}

// --- topology & closure ---

func TestEnqueue_TopologicalOrder(t *testing.T) {
	// C depends on B depends on A: order must be A, B, C.
	k := New(Options{})
	id, evts := k.Enqueue(Plan{
		Selected: ids(3),
		Items: []Item{
			item(1), item(2, ids(1)), item(3, ids(2)),
		},
	}, ids(3))
	_ = id
	wantKinds(t, kinds(evts), EventBatchAdmitted, EventRecordStart)
	if evts[1].ItemID != 1 {
		t.Fatalf("first start = item %d, want 1", evts[1].ItemID)
	}
}

func TestEnqueue_ClosurePullsPrerequisites(t *testing.T) {
	// Only D selected; D needs C, C needs A and B: all four build.
	k := New(Options{})
	_, evts := k.Enqueue(Plan{
		Selected: ids(4),
		Items: []Item{
			item(1), item(2), item(3, ids(1), ids(2)), item(4, ids(3)),
		},
	}, ids(4))
	if evts[1].ItemID != 1 {
		t.Fatalf("closure did not pull prerequisites: first = %d", evts[1].ItemID)
	}
}

func TestEnqueue_SharedPrerequisiteBuildsOnce(t *testing.T) {
	// B and C both depend on A; selecting B and C must plan A once.
	k := New(Options{})
	batchID, evts := k.Enqueue(Plan{
		Selected: ids(2, 3),
		Items:    []Item{item(1), item(2, ids(1)), item(3, ids(1))},
	}, ids(2, 3))
	_ = evts
	b := k.batches[batchID]
	if len(b.Records) != 3 {
		t.Fatalf("planned %d records, want 3 (A, B, C)", len(b.Records))
	}
	seen := map[ItemID]int{}
	for _, r := range b.Records {
		seen[r.ItemID]++
	}
	for id, n := range seen {
		if n != 1 {
			t.Fatalf("item %d planned %d times", id, n)
		}
	}
}

func TestEnqueue_DiamondOrder(t *testing.T) {
	// Diamond: A; B,C need A; D needs B and C. A before B/C before D.
	k := New(Options{})
	batchID, _ := k.Enqueue(Plan{
		Selected: ids(4),
		Items: []Item{
			item(1), item(2, ids(1)), item(3, ids(1)), item(4, ids(2), ids(3)),
		},
	}, ids(4))
	b := k.batches[batchID]
	pos := map[ItemID]int{}
	for i, r := range b.Records {
		pos[r.ItemID] = i
	}
	if !(pos[1] < pos[2] && pos[1] < pos[3] && pos[2] < pos[4] && pos[3] < pos[4]) {
		t.Fatalf("diamond order broken: %v", pos)
	}
}

// --- serial advance ---

func TestSerial_OneAtATime_AdvanceOnDone(t *testing.T) {
	k := New(Options{})
	batchID, evts := k.Enqueue(Plan{
		Selected: ids(1, 2),
		Items:    []Item{item(1), item(2)},
	}, ids(1, 2))
	// Admitted + first start; nothing else may start yet.
	wantKinds(t, kinds(evts), EventBatchAdmitted, EventRecordStart)

	done := k.NotifyRecordDone(batchID, evts[1].RecordID, OutcomeSucceeded)
	// Next start only: batch not finished.
	wantKinds(t, kinds(done), EventRecordStart)

	second := done[0]
	done2 := k.NotifyRecordDone(batchID, second.RecordID, OutcomeSucceeded)
	wantKinds(t, kinds(done2), EventBatchDone)
	if done2[0].Terminal != BatchCompleted {
		t.Fatalf("terminal = %s, want completed", done2[0].Terminal)
	}
}

// --- parallel cap & queueing ---

func TestParallelCap_QueueAndAdmit(t *testing.T) {
	k := New(Options{MaxParallel: 2})
	b1, evts1 := k.Enqueue(Plan{Selected: ids(1), Items: []Item{item(1)}}, ids(1))
	b2, _ := k.Enqueue(Plan{Selected: ids(2), Items: []Item{item(2)}}, ids(2))
	b3, _ := k.Enqueue(Plan{Selected: ids(3), Items: []Item{item(3)}}, ids(3))

	// Batch 3 queued: admission events only for b1, b2.
	if k.batches[b3].Status != BatchQueued {
		t.Fatalf("b3 status = %s, want queued", k.batches[b3].Status)
	}
	if k.batches[b1].Status != BatchRunning || k.batches[b2].Status != BatchRunning {
		t.Fatalf("b1=%s b2=%s, both want running", k.batches[b1].Status, k.batches[b2].Status)
	}

	// Finishing b1's record frees the slot for b3.
	k.NotifyRecordDone(b1, evts1[len(evts1)-1].RecordID, OutcomeSucceeded)
	if k.batches[b3].Status != BatchRunning {
		t.Fatalf("b3 after slot free = %s, want running", k.batches[b3].Status)
	}
}

// --- prereq groups ---

func TestGroupOR_OneSucceedsOneFails_DependentRuns(t *testing.T) {
	// Group {1,2} is OR: 1 fails, 2 succeeds -> 3 still builds.
	k := New(Options{})
	batchID, evts := k.Enqueue(Plan{
		Selected: ids(3),
		Items:    []Item{item(1), item(2), item(3, ids(1, 2))},
	}, ids(3))
	first := evts[1] // item 1 starts
	done := k.NotifyRecordDone(batchID, first.RecordID, OutcomeFailed)
	// Item 2 starts next (the OR group still has a live hope).
	wantKinds(t, kinds(done), EventRecordStart)
	if done[0].ItemID != 2 {
		t.Fatalf("after OR-member failure, item %d started, want 2", done[0].ItemID)
	}
	done2 := k.NotifyRecordDone(batchID, done[0].RecordID, OutcomeSucceeded)
	// Item 3 now starts.
	wantKinds(t, kinds(done2), EventRecordStart)
	if done2[0].ItemID != 3 {
		t.Fatalf("dependent item %d started, want 3", done2[0].ItemID)
	}
}

func TestGroupOR_AllFail_DependentSkipped(t *testing.T) {
	k := New(Options{})
	batchID, evts := k.Enqueue(Plan{
		Selected: ids(3),
		Items:    []Item{item(1), item(2), item(3, ids(1, 2))},
	}, ids(3))
	r1 := evts[1]
	done := k.NotifyRecordDone(batchID, r1.RecordID, OutcomeFailed)
	r2 := done[0]
	done2 := k.NotifyRecordDone(batchID, r2.RecordID, OutcomeFailed)
	// Item 3 skipped, batch failed.
	wantKinds(t, kinds(done2), EventRecordSkip, EventBatchDone)
	if done2[0].Detail != "prerequisite group failed" {
		t.Fatalf("skip detail = %q", done2[0].Detail)
	}
	if done2[1].Terminal != BatchFailed {
		t.Fatalf("terminal = %s, want failed", done2[1].Terminal)
	}
}

func TestGroupAND_OneGroupFails_DependentSkipped(t *testing.T) {
	// Groups {1,2} AND {3}: group 1 fully fails -> dependent skipped even
	// though group 2 succeeds.
	k := New(Options{})
	batchID, evts := k.Enqueue(Plan{
		Selected: ids(4),
		Items:    []Item{item(1), item(2), item(3), item(4, ids(1, 2), ids(3))},
	}, ids(4))
	r1 := evts[1] // item 1
	done := k.NotifyRecordDone(batchID, r1.RecordID, OutcomeFailed)
	if done[0].ItemID != 2 {
		t.Fatalf("want item 2 next, got %d", done[0].ItemID)
	}
	done2 := k.NotifyRecordDone(batchID, done[0].RecordID, OutcomeFailed)
	if done[0].RecordID == r1.RecordID {
		t.Fatal("sanity")
	}
	// Group 1 dead; item 3 still builds (independent), then 4 skips.
	if len(done2) == 0 || done2[0].Kind != EventRecordStart || done2[0].ItemID != 3 {
		t.Fatalf("after AND group death, events = %v", kinds(done2))
	}
	done3 := k.NotifyRecordDone(batchID, done2[0].RecordID, OutcomeSucceeded)
	wantKinds(t, kinds(done3), EventRecordSkip, EventBatchDone)
	if done3[1].Terminal != BatchFailed {
		t.Fatalf("terminal = %s, want failed", done3[1].Terminal)
	}
}

func TestSkip_PropagatesRecursively(t *testing.T) {
	// A fails; B skipped (needs A); C skipped (needs B).
	k := New(Options{})
	batchID, evts := k.Enqueue(Plan{
		Selected: ids(3),
		Items:    []Item{item(1), item(2, ids(1)), item(3, ids(2))},
	}, ids(3))
	r1 := evts[1]
	done := k.NotifyRecordDone(batchID, r1.RecordID, OutcomeFailed)
	wantKinds(t, kinds(done), EventRecordSkip, EventRecordSkip, EventBatchDone)
	if done[2].Terminal != BatchFailed {
		t.Fatalf("terminal = %s, want failed", done[2].Terminal)
	}
}

// --- cancel ---

func TestCancel_RunningBatch(t *testing.T) {
	k := New(Options{})
	batchID, evts := k.Enqueue(Plan{
		Selected: ids(1, 2, 3),
		Items:    []Item{item(1), item(2), item(3)},
	}, ids(1, 2, 3))
	started := evts[1] // item 1 building

	// Finish item 1, item 2 starts; cancel during item 2.
	k.NotifyRecordDone(batchID, started.RecordID, OutcomeSucceeded)
	cancelEvts := k.Cancel(batchID)
	// Canceled: the running record + pending record cancel, batch done.
	n := 0
	for _, e := range cancelEvts {
		if e.Kind == EventRecordSkip {
			n++
		}
	}
	if n != 2 {
		t.Fatalf("cancel produced %d record events, want 2", n)
	}
	var terminal string
	for _, e := range cancelEvts {
		if e.Kind == EventBatchDone {
			terminal = e.Terminal
		}
	}
	if terminal != BatchCanceled {
		t.Fatalf("terminal = %s, want canceled", terminal)
	}
	if k.Active(batchID) {
		t.Fatal("canceled batch still active")
	}
}

func TestCancel_QueuedBatch(t *testing.T) {
	k := New(Options{MaxParallel: 1})
	b1, _ := k.Enqueue(Plan{Selected: ids(1), Items: []Item{item(1)}}, ids(1))
	b2, _ := k.Enqueue(Plan{Selected: ids(2), Items: []Item{item(2)}}, ids(2))

	evts := k.Cancel(b2)
	if len(evts) == 0 || evts[len(evts)-1].Terminal != BatchCanceled {
		t.Fatalf("queued cancel events = %v", evts)
	}
	if k.Active(b2) {
		t.Fatal("queued batch still active after cancel")
	}
	// Slot accounting: finishing b1 must still admit cleanly.
	done := k.NotifyRecordDone(b1, k.batches[b1].Records[0].RecordID, OutcomeSucceeded)
	_ = done
}

// --- cross-batch isolation ---

func TestCrossBatch_NoReuse(t *testing.T) {
	// Two batches selecting the same item each plan their own build.
	k := New(Options{MaxParallel: 4})
	plan := Plan{Selected: ids(1), Items: []Item{item(1)}}
	b1, _ := k.Enqueue(plan, ids(1))
	b2, _ := k.Enqueue(plan, ids(1))

	r1 := k.batches[b1].Records[0]
	r2 := k.batches[b2].Records[0]
	if r1.RecordID == r2.RecordID {
		t.Fatal("same record id across batches: reuse leaked in")
	}
}

// --- terminal accounting ---

func TestTerminal_FailedRecordMeansFailedBatch(t *testing.T) {
	k := New(Options{})
	batchID, evts := k.Enqueue(Plan{Selected: ids(1), Items: []Item{item(1)}}, ids(1))
	done := k.NotifyRecordDone(batchID, evts[1].RecordID, OutcomeFailed)
	wantKinds(t, kinds(done), EventBatchDone)
	if done[0].Terminal != BatchFailed {
		t.Fatalf("terminal = %s, want failed", done[0].Terminal)
	}
	// Terminal batches leave kernel memory immediately: a late result must
	// find nothing to write into (cancel/completion race).
	if k.Active(batchID) {
		t.Fatal("terminal batch still active")
	}
}

// --- id seeding ---

func TestNextIDs_Seed(t *testing.T) {
	k := New(Options{})
	k.NextIDs(100, 200)
	batchID, evts := k.Enqueue(Plan{Selected: ids(1), Items: []Item{item(1)}}, ids(1))
	if batchID != 101 {
		t.Fatalf("batch id = %d, want 101", batchID)
	}
	if evts[1].RecordID != 201 {
		t.Fatalf("record id = %d, want 201", evts[1].RecordID)
	}
}

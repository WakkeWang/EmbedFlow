package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/WakkeWang/EmbedFlow/pkg/protocol"
	"github.com/WakkeWang/EmbedFlow/server/internal/batch"
	"github.com/WakkeWang/EmbedFlow/server/internal/builder"
	"github.com/WakkeWang/EmbedFlow/server/internal/store"
	"github.com/WakkeWang/EmbedFlow/server/internal/transport"
)

// System setting keys (requirement 1.6: four admin-managed items).
const (
	SettingMaxParallelBatches = "max_parallel_batches"
	SettingChecksum           = "checksum"
	SettingPublishWhitelist   = "publish_whitelist"
	SettingTmpDir             = "build_tmp_dir"
)

// buildOrchestrator wires the batch kernel, the build executor and
// persistence together (the M2 analog of what the coordinator does for the
// session kernel). It owns live subscriptions: browser connections register
// per-batch and receive BuildEvent frames. The real subscriber type is the
// transport's ClientConn; tests inject fakes via the same interface.
type buildOrchestrator struct {
	kernel *batch.Kernel
	store  *store.Store
	hub    *coordinator // for WS fan-out + data dir

	mu       sync.Mutex
	subs     map[int64]map[*transport.ClientConn]struct{} // batchID -> subscribers
	cancels  map[int64]bool                               // recordID -> cancel requested
	items    map[int64]map[int64]store.BuildItem          // batchID -> itemID -> item
	itemName map[int64]string                             // itemID -> name (for events)
}

func newBuildOrchestrator(k *batch.Kernel, st *store.Store, hub *coordinator) *buildOrchestrator {
	return &buildOrchestrator{
		kernel:  k,
		store:   st,
		hub:     hub,
		subs:    map[int64]map[*transport.ClientConn]struct{}{},
		cancels: map[int64]bool{},
		items:   map[int64]map[int64]store.BuildItem{},
		itemName: map[int64]string{},
	}
}

// sweepStartup cancels batches left non-terminal by a server restart (M2
// mirrors CEO-7A: honest terminal states beat zombie "running" rows).
// Kernel counters seed from persisted max ids so new ids never collide.
func (b *buildOrchestrator) sweepStartup() error {
	ctx := context.Background()
	var maxBatch, maxRecord int64
	batchRows, err := b.store.NonterminalBatches(ctx)
	if err != nil {
		return err
	}
	recordRows, err := b.store.NonterminalRecords(ctx)
	if err != nil {
		return err
	}
	for _, bt := range batchRows {
		if bt.ID > maxBatch {
			maxBatch = bt.ID
		}
		_ = b.store.UpdateBatchStatus(ctx, bt.ID, batch.BatchCanceled)
	}
	for _, rec := range recordRows {
		if rec.ID > maxRecord {
			maxRecord = rec.ID
		}
		rec.Status = batch.RecCanceled
		rec.EndedAt = time.Now().Format(time.RFC3339)
		_ = b.store.UpdateBuildRecord(ctx, rec)
	}
	if maxBatch > 0 || maxRecord > 0 {
		b.kernel.NextIDs(maxBatch, maxRecord)
		slog.Warn("build startup sweep: canceled in-flight batches", "batches", len(batchRows), "records", len(recordRows))
	}
	return nil
}

// --- batch creation (from the HTTP handler) ---

// CreateBatch validates the selection, expands the plan, persists the batch
// plus pending records, and enqueues it in the kernel. Returns the batch id.
func (b *buildOrchestrator) CreateBatch(ctx context.Context, projectID int64, selected []int64, user string) (int64, error) {
	if len(selected) == 0 {
		return 0, fmt.Errorf("no items selected")
	}
	items, err := b.store.ListBuildItemsByIDs(ctx, selected)
	if err != nil {
		return 0, err
	}
	if len(items) != len(selected) {
		return 0, fmt.Errorf("some selected items do not exist")
	}
	// The kernel needs the full item universe for the closure; the project's
	// items are the universe (prerequisites may reference cross-project ids
	// only when M3 lands cross-project refs; M2 keeps closure within the
	// project by construction -- a foreign id resolves to no record and the
	// kernel treats it as skipped).
	universe, err := b.store.ListBuildItemsByProject(ctx, projectID)
	if err != nil {
		return 0, err
	}

	batchID, err := b.store.CreateBatch(ctx, store.Batch{
		ProjectID: projectID,
		Status:    batch.BatchQueued,
		ItemsJSON: encodeJSONInts(selected),
		CreatedBy: user,
	})
	if err != nil {
		return 0, err
	}

	// Plan into the kernel: the kernel assigns record ids and may already
	// admit the batch. Persist every planned record BEFORE applying events
	// (RecordStart handlers read the row from the store).
	plan := batch.Plan{Selected: toBatchIDs(selected), Items: toBatchItems(universe)}
	_, events := b.kernel.Enqueue(plan, toBatchIDs(selected))

	b.mu.Lock()
	b.items[batchID] = indexItems(items, universe)
	for _, it := range universe {
		b.itemName[it.ID] = it.Name
	}
	b.mu.Unlock()

	for _, rec := range b.kernel.PlanRecords(batch.BatchID(batchID)) {
		if _, err := b.store.CreateBuildRecord(ctx, store.BuildRecord{
			BatchID:   batchID,
			ItemID:    int64(rec.ItemID),
			ProjectID: projectID,
			Status:    batch.RecPending,
			Executor:  user,
		}); err != nil {
			return 0, err
		}
	}

	for _, ev := range events {
		b.applyEvent(ev, batchID, projectID, user)
	}
	return batchID, nil
}

// applyEvent persists kernel events and drives the executor.
func (b *buildOrchestrator) applyEvent(ev batch.Event, batchID, projectID int64, user string) {
	switch ev.Kind {
	case batch.EventBatchAdmitted:
		_ = b.store.UpdateBatchStatus(context.Background(), batchID, batch.BatchRunning)
		b.emit(batchID, 0, "", "batch_admitted", "")

	case batch.EventRecordStart:
		// Pending -> building, then run the build in a goroutine.
		rec, err := b.store.GetBuildRecord(context.Background(), int64(ev.RecordID))
		if err != nil {
			slog.Error("load record for start", "err", err)
			return
		}
		rec.Status = batch.RecBuilding
		rec.Executor = user
		rec.StartedAt = time.Now().Format(time.RFC3339)
		if err := b.store.UpdateBuildRecord(context.Background(), rec); err != nil {
			slog.Error("persist record start", "err", err)
		}
		b.emit(batchID, int64(ev.RecordID), b.nameOf(int64(ev.ItemID)), "started", "")
		go b.runRecord(batchID, projectID, int64(ev.RecordID), int64(ev.ItemID), user)

	case batch.EventRecordSkip:
		rec, err := b.store.GetBuildRecord(context.Background(), int64(ev.RecordID))
		if err != nil {
			slog.Error("load record for skip", "err", err)
			return
		}
		if ev.Detail == batch.RecCanceled {
			rec.Status = batch.RecCanceled
		} else {
			rec.Status = batch.RecSkipped
			rec.EndedAt = time.Now().Format(time.RFC3339)
		}
		if err := b.store.UpdateBuildRecord(context.Background(), rec); err != nil {
			slog.Error("persist record skip", "err", err)
		}
		phase := "skipped"
		if ev.Detail == batch.RecCanceled {
			phase = "canceled"
		}
		b.emit(batchID, int64(ev.RecordID), b.nameOf(int64(ev.ItemID)), phase, ev.Detail)

	case batch.EventBatchDone:
		_ = b.store.UpdateBatchStatus(context.Background(), batchID, ev.Terminal)
		b.emit(batchID, 0, "", "batch_done", ev.Terminal)
	}
}

// runRecord executes one build and reports the outcome back to the kernel.
func (b *buildOrchestrator) runRecord(batchID, projectID, recordID, itemID int64, user string) {
	item, err := b.store.GetBuildItem(context.Background(), itemID)
	if err != nil {
		slog.Error("load build item", "err", err)
		b.kernel.NotifyRecordDone(batch.BatchID(batchID), batch.RecordID(recordID), batch.OutcomeFailed)
		return
	}

	checksumSetting, _ := b.store.GetSetting(context.Background(), SettingChecksum)
	tmpRoot, err := b.tmpRoot()
	if err != nil {
		slog.Error("resolve tmp root", "err", err)
		b.kernel.NotifyRecordDone(batch.BatchID(batchID), batch.RecordID(recordID), batch.OutcomeFailed)
		return
	}

	ex := &builder.Executor{
		DataDir:  b.hub.dataDir,
		TmpRoot:  tmpRoot,
		Notify:   b.notifyBuild,
		Canceled: b.isCanceled,
	}
	res := ex.Run(recordID, item, checksumSetting)

	rec, err := b.store.GetBuildRecord(context.Background(), recordID)
	if err != nil {
		slog.Error("load record for done", "err", err)
		return
	}
	rec.CommitSHA = res.CommitSHA
	rec.VersionInfo = res.VersionInfo
	rec.EndedAt = time.Now().Format(time.RFC3339)
	if res.ExitCode != nil {
		ec := *res.ExitCode
		rec.ExitCode = &ec
	}
	if res.OK {
		rec.Status = batch.RecSucceeded
	} else {
		rec.Status = batch.RecFailed
	}
	if err := b.store.UpdateBuildRecord(context.Background(), rec); err != nil {
		slog.Error("persist record done", "err", err)
	}
	if len(res.Artifacts) > 0 {
		if err := b.store.InsertArtifacts(context.Background(), res.Artifacts); err != nil {
			slog.Error("insert artifacts", "err", err)
		}
	}

	phase := "succeeded"
	if !res.OK {
		phase = "failed"
	}
	b.emit(batchID, recordID, item.Name, phase, res.Detail)

	outcome := batch.OutcomeSucceeded
	if !res.OK {
		outcome = batch.OutcomeFailed
	}
	for _, ev := range b.kernel.NotifyRecordDone(batch.BatchID(batchID), batch.RecordID(recordID), outcome) {
		b.applyEvent(ev, batchID, projectID, user)
	}
}

// notifyBuild pushes streamed build log lines: the executor already wrote
// the file; here we forward to WS subscribers. Phases ride in front of the
// line so the UI can group without a second channel.
func (b *buildOrchestrator) notifyBuild(recordID int64, phase, line string) {
	rec, err := b.store.GetBuildRecord(context.Background(), recordID)
	if err != nil {
		return
	}
	b.emit(rec.BatchID, recordID, "", "log", phase+": "+line)
}

func (b *buildOrchestrator) isCanceled(recordID int64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.cancels[recordID]
}

func (b *buildOrchestrator) nameOf(itemID int64) string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.itemName[itemID]
}

// tmpRoot resolves the configured build temp dir (requirement 1.6 item 4).
func (b *buildOrchestrator) tmpRoot() (string, error) {
	if v, err := b.store.GetSetting(context.Background(), SettingTmpDir); err != nil {
		return "", err
	} else if strings.TrimSpace(v) != "" {
		return v, nil
	}
	return filepath.Join(b.hub.dataDir, "tmp"), nil
}

// CancelBatch requests cancellation: kernel-level cancel events + per-record
// cancel flags the executor checks between phases and during the command.
func (b *buildOrchestrator) CancelBatch(batchID int64) error {
	b.mu.Lock()
	for _, rec := range b.recordsOf(batchID) {
		b.cancels[rec] = true
	}
	b.mu.Unlock()
	for _, ev := range b.kernel.Cancel(batch.BatchID(batchID)) {
		b.applyEvent(ev, batchID, 0, "")
	}
	return nil
}

func (b *buildOrchestrator) recordsOf(batchID int64) []int64 {
	recs, err := b.store.RecordsForBatch(context.Background(), batchID)
	if err != nil {
		return nil
	}
	out := make([]int64, 0, len(recs))
	for _, r := range recs {
		if r.Status == batch.RecBuilding {
			out = append(out, r.ID)
		}
	}
	return out
}

// --- subscriptions ---

// subscribe registers a subscriber and immediately replays the batch's
// current snapshot (one event per record + the batch status): a viewer that
// opens the batch page after the run started still sees the full picture;
// the live stream continues from there.
func (b *buildOrchestrator) subscribe(batchID int64, c *transport.ClientConn) {
	b.mu.Lock()
	if b.subs[batchID] == nil {
		b.subs[batchID] = map[*transport.ClientConn]struct{}{}
	}
	b.subs[batchID][c] = struct{}{}
	b.mu.Unlock()

	records, err := b.store.RecordsForBatch(context.Background(), batchID)
	if err != nil {
		return
	}
	batchStatus := ""
	if bt, err := b.store.GetBatch(context.Background(), batchID); err == nil {
		batchStatus = bt.Status
	}
	for _, rec := range records {
		phase := "started"
		switch rec.Status {
		case batch.RecSucceeded:
			phase = "succeeded"
		case batch.RecFailed:
			phase = "failed"
		case batch.RecSkipped:
			phase = "skipped"
		case batch.RecCanceled:
			phase = "canceled"
		}
		_ = c.Send(protocol.Frame{Type: protocol.FrameBuildEvent, Body: &protocol.BuildEventFrame{
			BatchID: batchID, RecordID: rec.ID, ItemName: b.nameOf(rec.ItemID), Phase: phase,
		}})
	}
	if batchStatus != "" {
		_ = c.Send(protocol.Frame{Type: protocol.FrameBuildEvent, Body: &protocol.BuildEventFrame{
			BatchID: batchID, Phase: "batch_status", Detail: batchStatus,
		}})
	}
}

func (b *buildOrchestrator) unsubscribe(batchID int64, c *transport.ClientConn) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.subs[batchID], c)
}

func (b *buildOrchestrator) unsubscribeAll(c *transport.ClientConn) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for bid, m := range b.subs {
		delete(m, c)
		if len(m) == 0 {
			delete(b.subs, bid)
		}
	}
}

func (b *buildOrchestrator) emit(batchID int64, recordID int64, itemName, phase, detail string) {
	b.mu.Lock()
	targets := make([]*transport.ClientConn, 0, 2)
	for c := range b.subs[batchID] {
		targets = append(targets, c)
	}
	b.mu.Unlock()
	for _, c := range targets {
		_ = c.Send(protocol.Frame{Type: protocol.FrameBuildEvent, Body: &protocol.BuildEventFrame{
			BatchID:  batchID,
			RecordID: recordID,
			ItemName: itemName,
			Phase:    phase,
			Detail:   detail,
		}})
	}
}

// --- small helpers ---

func indexItems(selected, universe []store.BuildItem) map[int64]store.BuildItem {
	m := make(map[int64]store.BuildItem, len(universe))
	for _, it := range universe {
		m[it.ID] = it
	}
	return m
}

func toBatchIDs(ids []int64) []batch.ItemID {
	out := make([]batch.ItemID, 0, len(ids))
	for _, id := range ids {
		out = append(out, batch.ItemID(id))
	}
	return out
}

func toBatchItems(items []store.BuildItem) []batch.Item {
	out := make([]batch.Item, 0, len(items))
	for _, it := range items {
		out = append(out, batch.Item{ID: batch.ItemID(it.ID), PrereqGroups: decodePrereqGroups(it.PrereqJSON)})
	}
	return out
}

// decodePrereqGroups parses the stored prerequisite-group JSON
// ([[1,2],[3]] -- group-internal OR, group-to-group AND).
func decodePrereqGroups(raw string) [][]batch.ItemID {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var groups [][]int64
	if err := json.Unmarshal([]byte(raw), &groups); err != nil {
		return nil
	}
	out := make([][]batch.ItemID, 0, len(groups))
	for _, g := range groups {
		gids := make([]batch.ItemID, 0, len(g))
		for _, v := range g {
			gids = append(gids, batch.ItemID(v))
		}
		out = append(out, gids)
	}
	return out
}

func encodeJSONInts(ids []int64) string {
	if len(ids) == 0 {
		return "[]"
	}
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, fmt.Sprint(id))
	}
	return "[" + strings.Join(parts, ",") + "]"
}

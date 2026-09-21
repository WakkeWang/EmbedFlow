package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/WakkeWang/EmbedFlow/pkg/protocol"
	"github.com/WakkeWang/EmbedFlow/server/internal/batch"
	"github.com/WakkeWang/EmbedFlow/server/internal/store"
)

// Build-module end-to-end tests: REST batch creation + WS BuildEvent
// subscription, over real WebSockets against the coordinator (the
// ws_flow_test.go pattern applied to M2).

func seedGitFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	run := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	run("init", "-b", "main")
	run("config", "user.email", "t@t")
	run("config", "user.name", "t")
	script := "#!/bin/sh\necho payload-of-bootfw > bootfw-p0133.bin\n"
	if runtime.GOOS == "windows" {
		script = "@echo off\r\necho payload-of-bootfw> bootfw-p0133.bin\r\n"
	}
	name := "build.sh"
	if runtime.GOOS == "windows" {
		name = "build.bat"
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(script), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}
	run("add", "-A")
	run("commit", "-m", "x")
	return dir
}

func buildCmdForTest() string {
	if runtime.GOOS == "windows" {
		return ".\\build.bat"
	}
	return "sh build.sh"
}

// runGit runs a git command in dir; the caller decides what failure means.
func runGit(dir string, args ...string) error {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	_, err := cmd.CombinedOutput()
	return err
}

func TestBuildE2E_BatchRunsAndSubscribes(t *testing.T) {
	s := newTestServer(t)
	ctx := context.Background()
	adminTok := s.hub.tokens.issueForTesting("admin")

	// Project + build item pointing at the git fixture.
	presp := apiCall(t, s, adminTok, "POST", "/api/projects", map[string]string{"name": "P-b2b"})
	var proj struct {
		ID int64 `json:"id"`
	}
	_ = json.NewDecoder(presp.Body).Decode(&proj)

	fixture := seedGitFixture(t)
	itemReq := map[string]any{
		"name":        "bootfw",
		"source_type": "local",
		"local_path":  fixture,
		"command":     buildCmdForTest(),
		"artifacts":   []string{"bootfw*.bin"},
		"version_cmd": "cat VERSION",
	}
	if err := os.WriteFile(filepath.Join(fixture, "VERSION"), []byte("9.9.9-e2e\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := runGit(fixture, "add", "-A"); err != nil {
		t.Fatal(err)
	}
	if err := runGit(fixture, "commit", "-m", "version"); err != nil {
		t.Fatal(err)
	}

	iresp := apiCall(t, s, adminTok, "POST", "/api/projects/"+itoa(proj.ID)+"/build-items", itemReq)
	if iresp.StatusCode != 200 {
		t.Fatalf("create build item = %d", iresp.StatusCode)
	}
	var item struct {
		ID int64 `json:"id"`
	}
	_ = json.NewDecoder(iresp.Body).Decode(&item)

	// Subscribe over WS, then trigger the batch.
	browser := dial(t, s, "admin", "pw")

	bresp := apiCall(t, s, adminTok, "POST", "/api/batches", map[string]any{
		"project_id": proj.ID, "item_ids": []int64{item.ID},
	})
	if bresp.StatusCode != 200 {
		t.Fatalf("create batch = %d", bresp.StatusCode)
	}
	var batch struct {
		ID int64 `json:"id"`
	}
	_ = json.NewDecoder(bresp.Body).Decode(&batch)
	if batch.ID == 0 {
		t.Fatal("no batch id")
	}

	// Subscribe with the real batch id. The run may already be live: the
	// subscription replays a snapshot, so we still see every record state.
	browser.send(protocol.FrameBuildCtrl, &protocol.BuildCtrlFrame{Command: "subscribe", BatchID: batch.ID})

	// Collect events until batch_done (snapshot replay + live stream).
	var phases []string
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		f := browser.recvUntil(time.Until(deadline), func(f *protocol.Frame) bool {
			b, ok := f.Body.(*protocol.BuildEventFrame)
			if !ok {
				return false
			}
			return b.Phase == "batch_done" || b.Phase == "batch_status" || b.Phase == "started" ||
				b.Phase == "succeeded" || b.Phase == "failed" || b.Phase == "skipped"
		})
		b := f.Body.(*protocol.BuildEventFrame)
		phases = append(phases, b.Phase)
		if b.Phase == "batch_done" {
			if b.Detail != "completed" {
				t.Fatalf("batch terminal = %s, want completed", b.Detail)
			}
			break
		}
	}
	if len(phases) < 3 {
		t.Fatalf("events = %v, want snapshot + live events", phases)
	}

	// Record state in the store.
	records, err := s.st.RecordsForBatch(ctx, batch.ID)
	if err != nil || len(records) != 1 {
		t.Fatalf("records: %v", err)
	}
	if records[0].Status != "succeeded" {
		t.Fatalf("record status = %s detail: ended=%q exit=%v version=%q", records[0].Status, records[0].EndedAt, records[0].ExitCode, records[0].VersionInfo)
	}
	if records[0].CommitSHA == "" {
		t.Fatal("no commit snapshot on record")
	}
	if !strings.Contains(records[0].VersionInfo, "9.9.9-e2e") {
		t.Fatalf("version info = %q", records[0].VersionInfo)
	}

	// Artifacts present with checksum.
	arts, err := s.st.ArtifactsForRecord(ctx, records[0].ID)
	if err != nil || len(arts) != 1 {
		t.Fatalf("artifacts: %v (n=%d)", err, len(arts))
	}
	if len(arts[0].Checksum) != 64 {
		t.Fatalf("checksum = %q, want sha256 hex", arts[0].Checksum)
	}

	// Batch listing for the project.
	batches, err := s.st.ListBatchesByProject(ctx, proj.ID)
	if err != nil || len(batches) != 1 {
		t.Fatalf("batches: %v", err)
	}
	if batches[0].Status != "completed" {
		t.Fatalf("batch status = %s", batches[0].Status)
	}
}

func TestBuildE2E_BatchCancel(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("cancel timing needs sh sleep semantics")
	}
	s := newTestServer(t)
	adminTok := s.hub.tokens.issueForTesting("admin")

	presp := apiCall(t, s, adminTok, "POST", "/api/projects", map[string]string{"name": "P-cancel"})
	var proj struct {
		ID int64 `json:"id"`
	}
	_ = json.NewDecoder(presp.Body).Decode(&proj)

	fixture := seedGitFixture(t)
	itemReq := map[string]any{
		"name":        "slow",
		"source_type": "local",
		"local_path":  fixture,
		"command":     "sleep 30",
		"artifacts":   []string{"*.bin"},
	}
	iresp := apiCall(t, s, adminTok, "POST", "/api/projects/"+itoa(proj.ID)+"/build-items", itemReq)
	var item struct {
		ID int64 `json:"id"`
	}
	_ = json.NewDecoder(iresp.Body).Decode(&item)

	bresp := apiCall(t, s, adminTok, "POST", "/api/batches", map[string]any{
		"project_id": proj.ID, "item_ids": []int64{item.ID},
	})
	var batch struct {
		ID int64 `json:"id"`
	}
	_ = json.NewDecoder(bresp.Body).Decode(&batch)

	// Wait for the build to start, then cancel.
	time.Sleep(2 * time.Second)
	cresp := apiCall(t, s, adminTok, "POST", "/api/batches/"+itoa(batch.ID)+"/cancel", nil)
	if cresp.StatusCode != 204 {
		t.Fatalf("cancel = %d", cresp.StatusCode)
	}

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		recs, _ := s.st.RecordsForBatch(context.Background(), batch.ID)
		if len(recs) == 1 && recs[0].Status == "canceled" {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatal("record never reached canceled")
}

func TestBuildE2E_SettingsAPI(t *testing.T) {
	s := newTestServer(t)
	adminTok := s.hub.tokens.issueForTesting("admin")
	memberTok := s.hub.tokens.issueForRole("alice", "member")

	// Defaults surface.
	gresp := apiCall(t, s, adminTok, "GET", "/api/settings", nil)
	var settings map[string]string
	_ = json.NewDecoder(gresp.Body).Decode(&settings)
	if settings["max_parallel_batches"] != "4" || settings["checksum"] != "sha256" {
		t.Fatalf("defaults = %v", settings)
	}

	// Members read settings (effective values, plan B4: GET login-only) but
	// cannot write them.
	if resp := apiCall(t, s, memberTok, "GET", "/api/settings", nil); resp.StatusCode != 200 {
		t.Fatalf("member GET settings = %d, want 200", resp.StatusCode)
	}
	if resp := apiCall(t, s, memberTok, "PUT", "/api/settings", map[string]string{"checksum": "md5"}); resp.StatusCode != 403 {
		t.Fatalf("member PUT settings = %d, want 403", resp.StatusCode)
	}

	// Admin updates; bad checksum refused.
	if resp := apiCall(t, s, adminTok, "PUT", "/api/settings", map[string]string{"max_parallel_batches": "2", "checksum": "md5"}); resp.StatusCode != 204 {
		t.Fatalf("PUT settings = %d, want 204", resp.StatusCode)
	}
	if resp := apiCall(t, s, adminTok, "PUT", "/api/settings", map[string]string{"checksum": "crc32"}); resp.StatusCode != 400 {
		t.Fatalf("bad checksum = %d, want 400", resp.StatusCode)
	}

	v, _ := s.st.GetSetting(context.Background(), "max_parallel_batches")
	if v != "2" {
		t.Fatalf("setting persisted = %q", v)
	}
	// The setting is live: the scheduler cap follows it (requirement 1.6:
	// changes apply to admissions from now on).
	if got := s.hub.builds.kernel.MaxParallel(); got != 2 {
		t.Fatalf("kernel max parallel = %d, want 2", got)
	}
}

func TestBuildE2E_DeleteRecordModes(t *testing.T) {
	s := newTestServer(t)
	ctx := context.Background()
	adminTok := s.hub.tokens.issueForTesting("admin")

	presp := apiCall(t, s, adminTok, "POST", "/api/projects", map[string]string{"name": "P-del"})
	var proj struct {
		ID int64 `json:"id"`
	}
	_ = json.NewDecoder(presp.Body).Decode(&proj)

	fixture := seedGitFixture(t)
	itemReq := map[string]any{
		"name": "bootfw", "source_type": "local", "local_path": fixture,
		"command": buildCmdForTest(), "artifacts": []string{"bootfw*.bin"},
	}
	iresp := apiCall(t, s, adminTok, "POST", "/api/projects/"+itoa(proj.ID)+"/build-items", itemReq)
	var item struct {
		ID int64 `json:"id"`
	}
	_ = json.NewDecoder(iresp.Body).Decode(&item)

	bresp := apiCall(t, s, adminTok, "POST", "/api/batches", map[string]any{
		"project_id": proj.ID, "item_ids": []int64{item.ID},
	})
	var batch struct {
		ID int64 `json:"id"`
	}
	_ = json.NewDecoder(bresp.Body).Decode(&batch)

	// Wait for completion.
	deadline := time.Now().Add(30 * time.Second)
	var recID int64
	for time.Now().Before(deadline) {
		recs, _ := s.st.RecordsForBatch(ctx, batch.ID)
		if len(recs) == 1 && (recs[0].Status == "succeeded" || recs[0].Status == "failed") {
			recID = recs[0].ID
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if recID == 0 {
		t.Fatal("batch never finished")
	}

	// Deleting a live record is refused -- n/a here (finished), so exercise
	// mode=record (artifacts stay on disk) then mode=artifacts on a second
	// record is unnecessary: verify artifact download then delete.
	arts, _ := s.st.ArtifactsForRecord(ctx, recID)
	if len(arts) != 1 {
		t.Fatalf("artifacts = %d", len(arts))
	}
	dresp := apiCall(t, s, adminTok, "GET", "/api/artifacts/"+itoa(arts[0].ID)+"/download", nil)
	if dresp.StatusCode != 200 {
		t.Fatalf("artifact download = %d", dresp.StatusCode)
	}

	// mode=record: row gone, artifact file still on disk.
	if resp := apiCall(t, s, adminTok, "DELETE", "/api/build-records/"+itoa(recID)+"?mode=record", nil); resp.StatusCode != 204 {
		t.Fatalf("delete record = %d", resp.StatusCode)
	}
	if _, err := os.Stat(filepath.Join(s.hub.dataDir, "artifacts", itoa(recID), "bootfw-p0133.bin")); err != nil {
		t.Fatalf("mode=record must keep the artifact file: %v", err)
	}
}

func TestBuildE2E_BatchValidation(t *testing.T) {
	s := newTestServer(t)
	adminTok := s.hub.tokens.issueForTesting("admin")

	// Empty selection refused.
	resp := apiCall(t, s, adminTok, "POST", "/api/batches", map[string]any{"project_id": 1, "item_ids": []int64{}})
	if resp.StatusCode != 400 {
		t.Fatalf("empty selection = %d, want 400", resp.StatusCode)
	}
	// Nonexistent items refused.
	resp = apiCall(t, s, adminTok, "POST", "/api/batches", map[string]any{"project_id": 1, "item_ids": []int64{999}})
	if resp.StatusCode != 400 {
		t.Fatalf("bad items = %d, want 400", resp.StatusCode)
	}
}

// The batch name is an optional trigger-time label ("发布前测试"): carried
// through to the row verbatim when set, empty when not, truncated (never a
// 400) when over the 80-rune cap.
func TestBuildE2E_BatchName(t *testing.T) {
	s := newTestServer(t)
	ctx := context.Background()
	adminTok := s.hub.tokens.issueForTesting("admin")

	presp := apiCall(t, s, adminTok, "POST", "/api/projects", map[string]string{"name": "P-name"})
	proj := decodeBody[struct {
		ID int64 `json:"id"`
	}](t, readRespBody(t, presp))
	fixture := seedGitFixture(t)
	iresp := apiCall(t, s, adminTok, "POST", "/api/projects/"+itoa(proj.ID)+"/build-items", map[string]any{
		"name": "bootfw", "source_type": "local", "local_path": fixture,
		"command": buildCmdForTest(), "artifacts": []string{"bootfw*.bin"},
	})
	item := decodeBody[struct {
		ID int64 `json:"id"`
	}](t, readRespBody(t, iresp))

	// Named batch: the label round-trips.
	cresp := apiCall(t, s, adminTok, "POST", "/api/batches", map[string]any{
		"project_id": proj.ID, "item_ids": []int64{item.ID}, "name": "发布前测试",
	})
	if cresp.StatusCode != 200 {
		t.Fatalf("named create = %d", cresp.StatusCode)
	}
	var created struct {
		ID int64 `json:"id"`
	}
	_ = json.NewDecoder(cresp.Body).Decode(&created)
	b, err := s.st.GetBatch(ctx, created.ID)
	if err != nil {
		t.Fatalf("get named batch: %v", err)
	}
	if b.Name != "发布前测试" {
		t.Fatalf("batch name = %q", b.Name)
	}

	// Unnamed batch: empty string, not "undefined"/null artifacts.
	cresp = apiCall(t, s, adminTok, "POST", "/api/batches", map[string]any{
		"project_id": proj.ID, "item_ids": []int64{item.ID},
	})
	if cresp.StatusCode != 200 {
		t.Fatalf("unnamed create = %d", cresp.StatusCode)
	}
	_ = json.NewDecoder(cresp.Body).Decode(&created)
	b, err = s.st.GetBatch(ctx, created.ID)
	if err != nil {
		t.Fatalf("get unnamed batch: %v", err)
	}
	if b.Name != "" {
		t.Fatalf("unnamed batch name = %q, want empty", b.Name)
	}

	// An 200-rune name is truncated to 80, still 200.
	long := strings.Repeat("名", 200)
	cresp = apiCall(t, s, adminTok, "POST", "/api/batches", map[string]any{
		"project_id": proj.ID, "item_ids": []int64{item.ID}, "name": long,
	})
	if cresp.StatusCode != 200 {
		t.Fatalf("long-name create = %d", cresp.StatusCode)
	}
	_ = json.NewDecoder(cresp.Body).Decode(&created)
	b, err = s.st.GetBatch(ctx, created.ID)
	if err != nil {
		t.Fatalf("get long-named batch: %v", err)
	}
	if got := []rune(b.Name); len(got) != 80 {
		t.Fatalf("long name len = %d, want 80", len(got))
	}

	// Whitespace-only name stores as empty.
	cresp = apiCall(t, s, adminTok, "POST", "/api/batches", map[string]any{
		"project_id": proj.ID, "item_ids": []int64{item.ID}, "name": "   ",
	})
	_ = json.NewDecoder(cresp.Body).Decode(&created)
	b, _ = s.st.GetBatch(ctx, created.ID)
	if b.Name != "" {
		t.Fatalf("whitespace name = %q, want empty", b.Name)
	}
}

// The kernel's event contract: events carry their own batch id, and the
// orchestrator must apply them under THAT id. Finishing (or canceling) a
// record frees a run slot; pumpLocked may admit ANOTHER queued batch and
// return its Admitted/RecordStart events. Applying those under the wrong
// batch id ran the other batch's build under this batch's id, whose
// kernel.Active check then failed -- the other batch's record stuck
// "building" forever, its card "running", one kernel slot leaked.
func TestBuildE2E_QueuedBatchEventsCarryOwnBatchID(t *testing.T) {
	s := newTestServer(t)
	ctx := context.Background()
	adminTok := s.hub.tokens.issueForTesting("admin")

	presp := apiCall(t, s, adminTok, "POST", "/api/projects", map[string]string{"name": "P-pump"})
	proj := decodeBody[struct {
		ID int64 `json:"id"`
	}](t, readRespBody(t, presp))
	fixture := seedGitFixture(t)
	iresp := apiCall(t, s, adminTok, "POST", "/api/projects/"+itoa(proj.ID)+"/build-items", map[string]any{
		"name": "bootfw", "source_type": "local", "local_path": fixture,
		"command": buildCmdForTest(), "artifacts": []string{"bootfw*.bin"},
	})
	item := decodeBody[struct {
		ID int64 `json:"id"`
	}](t, readRespBody(t, iresp))

	// Cap parallel batches at 1 so the second batch queues behind the first.
	if resp := apiCall(t, s, adminTok, "PUT", "/api/settings", map[string]string{"max_parallel_batches": "1"}); resp.StatusCode != 204 {
		t.Fatalf("PUT settings = %d", resp.StatusCode)
	}

	b1, err := s.hub.builds.CreateBatch(ctx, proj.ID, []int64{item.ID}, "", "admin")
	if err != nil {
		t.Fatalf("batch 1: %v", err)
	}
	b2, err := s.hub.builds.CreateBatch(ctx, proj.ID, []int64{item.ID}, "", "admin")
	if err != nil {
		t.Fatalf("batch 2: %v", err)
	}

	// Both batches must reach a terminal state: batch 2 is admitted when
	// batch 1 finishes, and its RecordStart must run under batch 2's id.
	terminal := func(id int64) bool {
		b, _ := s.st.GetBatch(ctx, id)
		return b.Status == "completed" || b.Status == "failed" || b.Status == "canceled"
	}
	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		if terminal(b1) && terminal(b2) {
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	if !terminal(b1) {
		t.Fatalf("batch 1 stuck: %s", func() string { bb, _ := s.st.GetBatch(ctx, b1); return bb.Status }())
	}
	if !terminal(b2) {
		t.Fatalf("batch 2 stuck (queued batch never ran under its own batch id): %s",
			func() string { bb, _ := s.st.GetBatch(ctx, b2); return bb.Status }())
	}
	recs, _ := s.st.RecordsForBatch(ctx, b2)
	for _, r := range recs {
		if r.Status == "building" || r.Status == "pending" {
			t.Fatalf("batch 2 record %d stuck in %s", r.ID, r.Status)
		}
	}
}

// The 0.80 deployment bug: after a restart with an all-terminal history,
// the kernel's record counter restarted at zero and the next batch's
// CreateBuildRecord collided with historical ids (primary-key error, the
// batch row stayed behind as a zombie "running"). sweepStartup must seed
// from the absolute max ids, not only the nonterminal rows.
func TestBuildE2E_KernelCountersSurviveRestart(t *testing.T) {
	s := newTestServer(t)
	adminTok := s.hub.tokens.issueForTesting("admin")

	// History: one project, one item, one batch left terminal (failed).
	presp := apiCall(t, s, adminTok, "POST", "/api/projects", map[string]string{"name": "P-restart"})
	proj := decodeBody[struct {
		ID int64 `json:"id"`
	}](t, readRespBody(t, presp))
	fixture := seedGitFixture(t)
	iresp := apiCall(t, s, adminTok, "POST", "/api/projects/"+itoa(proj.ID)+"/build-items", map[string]any{
		"name": "bootfw", "source_type": "local", "local_path": fixture,
		"command": buildCmdForTest(), "artifacts": []string{"bootfw*.bin"},
	})
	item := decodeBody[struct {
		ID int64 `json:"id"`
	}](t, readRespBody(t, iresp))

	// Simulate the historical rows: a failed batch + failed record (id 1, 1).
	// Direct SQL: CreateBatch would mint fresh ids; the seed needs exact ids.
	if _, err := s.st.CreateBatch(context.Background(), store.Batch{
		ProjectID: proj.ID, Status: "failed", ItemsJSON: "[" + itoa(item.ID) + "]", CreatedBy: "seed",
	}); err != nil {
		t.Fatalf("seed history batch: %v", err)
	}
	if _, err := s.st.CreateBuildRecord(context.Background(), store.BuildRecord{
		BatchID: 1, ItemID: item.ID, ProjectID: proj.ID, Status: "failed",
	}); err != nil {
		t.Fatalf("seed history record: %v", err)
	}
	// The seeds are the first rows, so both carry id 1; force the batch row
	// terminal (it was created queued by the store's default).
	if err := s.st.UpdateBatchStatus(context.Background(), 1, "failed"); err != nil {
		t.Fatalf("terminal history batch: %v", err)
	}

	// Restart sweep on a coordinator whose kernel has never seen these ids.
	hub2 := newCoordinator(s.kernel, s.st, s.hub.dataDir)
	hub2.builds = newBuildOrchestrator(batch.New(batch.Options{}), s.st, hub2)
	if err := hub2.builds.sweepStartup(); err != nil {
		t.Fatalf("sweep: %v", err)
	}

	// A new batch must NOT reuse record id 1: its records start above the
	// historical max.
	batchID, err := hub2.builds.CreateBatch(context.Background(), proj.ID, []int64{item.ID}, "", "admin")
	if err != nil {
		t.Fatalf("create batch after restart: %v", err)
	}
	if batchID <= 1 {
		t.Fatalf("batch id %d collides with history", batchID)
	}
	records, err := s.st.RecordsForBatch(context.Background(), batchID)
	if err != nil || len(records) == 0 {
		t.Fatalf("records for new batch: %v", err)
	}
	for _, r := range records {
		if r.ID <= 1 {
			t.Fatalf("record id %d collides with the historical row", r.ID)
		}
	}
}

// The 0.80 production bug (second fork): after running a batch and deleting
// the whole batch history, a server restart left sqlite_sequence high
// (AUTOINCREMENT never reuses ids) while the kernel's counters seeded from
// MAX(id)=0. The next batch was kernel-id 1 but sqlite row 14;
// PlanRecords(14) found nothing, no records were ever persisted, and the
// batch sat "running" forever with zero build records. The fix persists rows
// under the kernel-minted ids, so this test asserts id parity structurally:
// a fresh kernel (never seeded, the restart condition) must still produce a
// batch whose records exist and complete.
func TestBuildE2E_BatchIDsSurviveHistoryDelete(t *testing.T) {
	s := newTestServer(t)
	ctx := context.Background()
	adminTok := s.hub.tokens.issueForTesting("admin")

	presp := apiCall(t, s, adminTok, "POST", "/api/projects", map[string]string{"name": "P-fork"})
	proj := decodeBody[struct {
		ID int64 `json:"id"`
	}](t, readRespBody(t, presp))
	fixture := seedGitFixture(t)
	iresp := apiCall(t, s, adminTok, "POST", "/api/projects/"+itoa(proj.ID)+"/build-items", map[string]any{
		"name": "bootfw", "source_type": "local", "local_path": fixture,
		"command": buildCmdForTest(), "artifacts": []string{"bootfw*.bin"},
	})
	item := decodeBody[struct {
		ID int64 `json:"id"`
	}](t, readRespBody(t, iresp))

	// Phase 1: run one batch to completion so sqlite_sequence is high.
	batchID, err := s.hub.builds.CreateBatch(ctx, proj.ID, []int64{item.ID}, "", "admin")
	if err != nil {
		t.Fatalf("first batch: %v", err)
	}
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		recs, _ := s.st.RecordsForBatch(ctx, batchID)
		if len(recs) == 1 && (recs[0].Status == "succeeded" || recs[0].Status == "failed") {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	b1, _ := s.st.GetBatch(ctx, batchID)
	if b1.Status != "completed" && b1.Status != "failed" {
		t.Fatalf("first batch never finished: %s", b1.Status)
	}

	// Delete the whole batch history (the 0.80 operator action: artifacts
	// rows go first -- FK artifacts.build_record_id -- then records, batch).
	recs1, _ := s.st.RecordsForBatch(ctx, batchID)
	for _, r := range recs1 {
		if err := s.st.DeleteArtifactsForRecord(ctx, r.ID); err != nil {
			t.Fatalf("delete history artifacts: %v", err)
		}
		if err := s.st.DeleteBuildRecord(ctx, r.ID); err != nil {
			t.Fatalf("delete history record: %v", err)
		}
	}
	if err := s.st.DeleteBatch(ctx, batchID); err != nil {
		t.Fatalf("delete history batch: %v", err)
	}
	if m, _ := s.st.MaxBatchID(ctx); m != 0 {
		t.Fatalf("history not empty: max batch id %d", m)
	}

	// Phase 2: a fresh kernel + sweepStartup WITHOUT history to seed from --
	// exactly the post-restart state that forked the id spaces.
	hub2 := newCoordinator(s.kernel, s.st, s.hub.dataDir)
	hub2.builds = newBuildOrchestrator(batch.New(batch.Options{}), s.st, hub2)
	if err := hub2.builds.sweepStartup(); err != nil {
		t.Fatalf("sweep: %v", err)
	}

	// The new batch must get a row under the kernel's id, with records that
	// actually persist and run -- not a zombie "running" with zero records.
	batch2, err := hub2.builds.CreateBatch(ctx, proj.ID, []int64{item.ID}, "", "admin")
	if err != nil {
		t.Fatalf("batch after history delete: %v", err)
	}
	recs, err := s.st.RecordsForBatch(ctx, batch2)
	if err != nil || len(recs) == 0 {
		t.Fatalf("records for new batch: %v, %d", err, len(recs))
	}
	b2, _ := s.st.GetBatch(ctx, batch2)
	deadline = time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		b2, _ = s.st.GetBatch(ctx, batch2)
		if b2.Status == "completed" || b2.Status == "failed" {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if b2.Status != "completed" && b2.Status != "failed" {
		t.Fatalf("second batch stuck: %s", b2.Status)
	}
	for _, r := range recs {
		if r.ID != recs[0].ID && r.BatchID != batch2 {
			t.Fatalf("record %d not under batch %d", r.ID, batch2)
		}
	}
}

func itoa(v int64) string {
	return strconv.FormatInt(v, 10)
}

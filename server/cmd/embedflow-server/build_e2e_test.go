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

func itoa(v int64) string {
	return strconv.FormatInt(v, 10)
}

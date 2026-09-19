package store

import (
	"context"
	"testing"
)

// --- build items ---

func TestBuildItems_CRUDRoundTrip(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	pid, err := s.CreateProject(ctx, "P0001", "demo line")
	if err != nil {
		t.Fatalf("create project: %v", err)
	}

	id, err := s.CreateBuildItem(ctx, BuildItem{
		ProjectID:  pid,
		Name:       "bootfw",
		SourceType: "local",
		LocalPath:  "/root/git/bootfw",
		Command:    "./bb.sh board_defconfig && ./bb.sh",
		Artifacts:  []string{"bootfw*.bin"},
		TimeoutSec: 600,
		VersionCmd: "ver.sh -a",
		PrereqJSON: "[[2]]",
	})
	if err != nil {
		t.Fatalf("create build item: %v", err)
	}

	got, err := s.GetBuildItem(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name != "bootfw" || got.SourceType != "local" || got.Command != "./bb.sh board_defconfig && ./bb.sh" {
		t.Fatalf("round trip mismatch: %+v", got)
	}
	if len(got.Artifacts) != 1 || got.Artifacts[0] != "bootfw*.bin" {
		t.Fatalf("artifacts = %v", got.Artifacts)
	}
	if got.TimeoutSec != 600 || got.VersionCmd != "ver.sh -a" || got.PrereqJSON != "[[2]]" {
		t.Fatalf("scalar fields = %+v", got)
	}
	if got.CheckLatest {
		t.Fatal("check_latest not set should read back false")
	}

	got.Artifacts = []string{"out/*.bin", "map/*.map"}
	got.CheckLatest = true
	if err := s.UpdateBuildItem(ctx, got); err != nil {
		t.Fatalf("update: %v", err)
	}
	again, _ := s.GetBuildItem(ctx, id)
	if len(again.Artifacts) != 2 || !again.CheckLatest {
		t.Fatalf("update not persisted: %+v", again)
	}

	items, err := s.ListBuildItemsByProject(ctx, pid)
	if err != nil || len(items) != 1 {
		t.Fatalf("list by project: %v (n=%d)", err, len(items))
	}

	if err := s.DeleteBuildItem(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := s.GetBuildItem(ctx, id); err == nil {
		t.Fatal("get after delete should fail")
	}
}

func TestBuildItems_ListByIDs(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	pid, _ := s.CreateProject(ctx, "P", "")
	id1, _ := s.CreateBuildItem(ctx, BuildItem{ProjectID: pid, Name: "a", SourceType: "local", Command: "x"})
	id2, _ := s.CreateBuildItem(ctx, BuildItem{ProjectID: pid, Name: "b", SourceType: "local", Command: "x"})

	got, err := s.ListBuildItemsByIDs(ctx, []int64{id2, id1})
	if err != nil {
		t.Fatalf("by ids: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 items, got %d", len(got))
	}
	if got[0].Name != "a" || got[1].Name != "b" {
		t.Fatalf("unexpected order: %s, %s", got[0].Name, got[1].Name)
	}
}

// --- batches & records ---

func TestBatchAndRecords_Lifecycle(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	pid, _ := s.CreateProject(ctx, "P", "")
	itemID, _ := s.CreateBuildItem(ctx, BuildItem{ProjectID: pid, Name: "u-boot", SourceType: "git", GitURL: "http://example/git", Command: "make"})

	batchID, err := s.CreateBatch(ctx, Batch{ProjectID: pid, Status: "queued", ItemsJSON: "[1]", CreatedBy: "alice"})
	if err != nil {
		t.Fatalf("create batch: %v", err)
	}
	recID, err := s.CreateBuildRecord(ctx, BuildRecord{BatchID: batchID, ItemID: itemID, ProjectID: pid, Status: "pending", Executor: "alice"})
	if err != nil {
		t.Fatalf("create record: %v", err)
	}

	if err := s.UpdateBatchStatus(ctx, batchID, "running"); err != nil {
		t.Fatalf("batch status: %v", err)
	}

	rec, err := s.GetBuildRecord(ctx, recID)
	if err != nil {
		t.Fatalf("get record: %v", err)
	}
	rec.Status = "building"
	rec.StartedAt = "2026-09-19T10:00:00Z"
	if err := s.UpdateBuildRecord(ctx, rec); err != nil {
		t.Fatalf("update record: %v", err)
	}
	rec.CommitSHA = "abc123"
	rec.Status = "succeeded"
	rec.VersionInfo = "os 1.2.3\nkernel 4.19"
	ec := 0
	rec.ExitCode = &ec
	rec.EndedAt = "2026-09-19T10:05:00Z"
	if err := s.UpdateBuildRecord(ctx, rec); err != nil {
		t.Fatalf("update record final: %v", err)
	}

	got, err := s.GetBuildRecord(ctx, recID)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if got.Status != "succeeded" || got.CommitSHA != "abc123" || got.ExitCode == nil || *got.ExitCode != 0 {
		t.Fatalf("record fields: %+v", got)
	}
	if got.VersionInfo != "os 1.2.3\nkernel 4.19" {
		t.Fatalf("version_info = %q", got.VersionInfo)
	}

	nonterm, err := s.NonterminalRecords(ctx)
	if err != nil || len(nonterm) != 0 {
		t.Fatalf("nonterminal after success: %v (n=%d)", err, len(nonterm))
	}

	records, err := s.RecordsForBatch(ctx, batchID)
	if err != nil || len(records) != 1 {
		t.Fatalf("records for batch: %v", err)
	}
	byProject, err := s.RecordsForProject(ctx, pid)
	if err != nil || len(byProject) != 1 {
		t.Fatalf("records for project: %v", err)
	}
}

func TestNonterminalBatches_FindsQueuedAndRunning(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	pid, _ := s.CreateProject(ctx, "P", "")
	s.CreateBatch(ctx, Batch{ProjectID: pid, Status: "queued", CreatedBy: "a"})
	s.CreateBatch(ctx, Batch{ProjectID: pid, Status: "running", CreatedBy: "a"})
	s.CreateBatch(ctx, Batch{ProjectID: pid, Status: "completed", CreatedBy: "a"})

	got, err := s.NonterminalBatches(ctx)
	if err != nil {
		t.Fatalf("nonterminal: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 nonterminal, got %d", len(got))
	}
}

// --- artifacts ---

func TestArtifacts_InsertListDelete(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	pid, _ := s.CreateProject(ctx, "P", "")
	itemID, _ := s.CreateBuildItem(ctx, BuildItem{ProjectID: pid, Name: "x", SourceType: "local", Command: "x"})
	batchID, _ := s.CreateBatch(ctx, Batch{ProjectID: pid, Status: "completed", CreatedBy: "a"})
	recID, _ := s.CreateBuildRecord(ctx, BuildRecord{BatchID: batchID, ItemID: itemID, ProjectID: pid, Status: "succeeded"})

	if err := s.InsertArtifacts(ctx, []Artifact{
		{BuildRecordID: recID, Name: "bootfw-v2.bin", Size: 1024, Checksum: "sha256:abc"},
		{BuildRecordID: recID, Name: "u-boot.bin", Size: 512, Checksum: "sha256:def"},
	}); err != nil {
		t.Fatalf("insert artifacts: %v", err)
	}

	got, err := s.ArtifactsForRecord(ctx, recID)
	if err != nil || len(got) != 2 {
		t.Fatalf("artifacts for record: %v (n=%d)", err, len(got))
	}
	if got[0].Name != "bootfw-v2.bin" || got[0].Size != 1024 {
		t.Fatalf("artifact 0: %+v", got[0])
	}

	a, err := s.GetArtifact(ctx, got[1].ID)
	if err != nil || a.Name != "u-boot.bin" {
		t.Fatalf("get artifact: %v %+v", err, a)
	}

	if err := s.DeleteArtifactsForRecord(ctx, recID); err != nil {
		t.Fatalf("delete artifacts: %v", err)
	}
	got, _ = s.ArtifactsForRecord(ctx, recID)
	if len(got) != 0 {
		t.Fatalf("artifacts after delete: %d", len(got))
	}
}

// --- settings ---

func TestSettings_GetSetDefaults(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	// Unset keys read as empty; callers own the defaults.
	v, err := s.GetSetting(ctx, "max_parallel_batches")
	if err != nil || v != "" {
		t.Fatalf("unset setting: %v %q", err, v)
	}

	if err := s.SetSetting(ctx, "max_parallel_batches", "2"); err != nil {
		t.Fatalf("set: %v", err)
	}
	if err := s.SetSetting(ctx, "max_parallel_batches", "6"); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	if v, _ := s.GetSetting(ctx, "max_parallel_batches"); v != "6" {
		t.Fatalf("upsert not applied: %q", v)
	}

	all, err := s.AllSettings(ctx)
	if err != nil || all["max_parallel_batches"] != "6" {
		t.Fatalf("all settings: %v %+v", err, all)
	}
}

// --- users ---

func TestUsers_ListAndUpdatePassword(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	s.CreateUser(ctx, "admin", "pw1", "admin")
	s.CreateUser(ctx, "alice", "pw2", "member")

	users, err := s.ListUsers(ctx)
	if err != nil || len(users) != 2 {
		t.Fatalf("list users: %v (n=%d)", err, len(users))
	}
	if users[0].Username != "admin" || users[0].Role != "admin" {
		t.Fatalf("user 0: %+v", users[0])
	}

	if err := s.ChangeOwnPassword(ctx, users[1].ID, "wrong", "new"); err == nil {
		t.Fatal("wrong old password must fail")
	}
	if err := s.ChangeOwnPassword(ctx, users[1].ID, "pw2", "newpw"); err != nil {
		t.Fatalf("change own: %v", err)
	}
	if _, err := s.AuthenticateUser(ctx, "alice", "newpw"); err != nil {
		t.Fatalf("auth with new password: %v", err)
	}

	if err := s.UpdatePassword(ctx, users[1].ID, "resetpw"); err != nil {
		t.Fatalf("admin reset: %v", err)
	}
	if _, err := s.AuthenticateUser(ctx, "alice", "resetpw"); err != nil {
		t.Fatalf("auth after reset: %v", err)
	}
}

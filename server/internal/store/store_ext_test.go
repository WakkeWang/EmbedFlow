package store

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func openTestStore2(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, err := Open(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

// --- users (requirement 1.5: two-level roles, admin creates accounts) ---

func TestUsers_CreateAndAuthenticate(t *testing.T) {
	s := openTestStore2(t)
	ctx := context.Background()

	if err := s.CreateUser(ctx, "admin", "pw-admin", "admin"); err != nil {
		t.Fatalf("create admin: %v", err)
	}
	if err := s.CreateUser(ctx, "alice", "pw-alice", "member"); err != nil {
		t.Fatalf("create alice: %v", err)
	}

	u, err := s.AuthenticateUser(ctx, "alice", "pw-alice")
	if err != nil {
		t.Fatalf("auth: %v", err)
	}
	if u.Role != "member" {
		t.Fatalf("role = %q", u.Role)
	}
	if _, err := s.AuthenticateUser(ctx, "alice", "wrong"); err == nil {
		t.Fatal("wrong password should fail")
	}
	if _, err := s.AuthenticateUser(ctx, "nobody", "x"); err == nil {
		t.Fatal("unknown user should fail")
	}
}

// --- devices ---

func TestDevices_ListAndDelete(t *testing.T) {
	s := openTestStore2(t)
	ctx := context.Background()

	id1, _ := s.CreateDevice(ctx, "d2000-lab", "yy0133")
	id2, _ := s.CreateDevice(ctx, "d2000-rack", "yy0133")

	devs, err := s.ListDevices(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(devs) != 2 || devs[0].Name != "d2000-lab" || devs[1].Name != "d2000-rack" {
		t.Fatalf("devices = %+v", devs)
	}
	_ = id1
	if err := s.DeleteDevice(ctx, id2); err != nil {
		t.Fatalf("delete: %v", err)
	}
	devs, _ = s.ListDevices(ctx)
	if len(devs) != 1 {
		t.Fatalf("after delete: %d devices", len(devs))
	}
}

// --- expect rules (issue #7: shared, reusable, reviewable config) ---

func TestExpectRules_CRUD(t *testing.T) {
	s := openTestStore2(t)
	ctx := context.Background()

	id, err := s.CreateExpectRule(ctx, "yy0133-flash", `[{"await":"Login:","send":"root\r"}]`)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := s.GetExpectRule(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name != "yy0133-flash" || got.StepsJSON == "" {
		t.Fatalf("rule = %+v", got)
	}

	got.Name = "renamed"
	if err := s.UpdateExpectRule(ctx, got); err != nil {
		t.Fatalf("update: %v", err)
	}
	after, _ := s.GetExpectRule(ctx, id)
	if after.Name != "renamed" {
		t.Fatalf("after update: %+v", after)
	}

	list, err := s.ListExpectRules(ctx)
	if err != nil || len(list) != 1 {
		t.Fatalf("list = %+v err %v", list, err)
	}
}

// --- confirmations (issue #9: human cards persisted in the session record) ---

func TestConfirmations_InsertAndResolve(t *testing.T) {
	s := openTestStore2(t)
	ctx := context.Background()
	devID, _ := s.CreateDevice(ctx, "dev", "proj")
	sessID, _ := s.CreateSession(ctx, Session{ID: 100, DeviceID: devID, Kind: "manual", Owner: "a", State: "active", StartedAt: time.Now()})

	cid, err := s.InsertConfirmation(ctx, sessID, "LED on?")
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if err := s.ResolveConfirmation(ctx, cid, "pass", "both LEDs green"); err != nil {
		t.Fatalf("resolve: %v", err)
	}
	confs, err := s.ConfirmationsForSession(ctx, sessID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(confs) != 1 || confs[0].Result != "pass" || confs[0].Note != "both LEDs green" || confs[0].State != "resolved" {
		t.Fatalf("confirmations = %+v", confs)
	}
}

// --- session history per device (issue #6) ---

func TestSessionHistory_ForDevice(t *testing.T) {
	s := openTestStore2(t)
	ctx := context.Background()
	devID, _ := s.CreateDevice(ctx, "dev", "proj")
	now := time.Date(2026, 9, 18, 12, 0, 0, 0, time.UTC)

	id1, _ := s.CreateSession(ctx, Session{ID: 1, DeviceID: devID, Kind: "manual", Owner: "a", State: "closed", StartedAt: now})
	id2, _ := s.CreateSession(ctx, Session{ID: 2, DeviceID: devID, Kind: "task", Owner: "b", State: "failed", StartedAt: now.Add(time.Hour)})

	hist, err := s.SessionsForDevice(ctx, devID)
	if err != nil {
		t.Fatalf("history: %v", err)
	}
	if len(hist) != 2 || hist[0].ID != id1 || hist[1].ID != id2 {
		t.Fatalf("history = %+v", hist)
	}
	if hist[1].LogIncomplete {
		t.Fatal("default should not be incomplete")
	}
}

// --- log_incomplete flag (design doc Failure modes) ---

func TestMarkLogIncomplete_Persists(t *testing.T) {
	s := openTestStore2(t)
	ctx := context.Background()
	devID, _ := s.CreateDevice(ctx, "dev", "proj")
	sessID, _ := s.CreateSession(ctx, Session{ID: 100, DeviceID: devID, Kind: "manual", Owner: "a", State: "active", StartedAt: time.Now()})

	if err := s.MarkLogIncomplete(ctx, sessID); err != nil {
		t.Fatalf("mark: %v", err)
	}
	got, _ := s.GetSession(ctx, sessID)
	if !got.LogIncomplete {
		t.Fatalf("log_incomplete not persisted: %+v", got)
	}
}

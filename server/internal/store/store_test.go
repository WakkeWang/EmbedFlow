package store

import (
	"context"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func openTestStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, err := Open(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestOpen_WALModeAndBusyTimeout(t *testing.T) {
	s := openTestStore(sT(t))
	for _, want := range [][2]string{
		{"journal_mode", "wal"},
		{"busy_timeout", "5000"},
	} {
		var k, v string
		if err := s.db.QueryRow("PRAGMA " + want[0]).Scan(&v); err != nil {
			t.Fatalf("pragma %s: %v", want[0], err)
		}
		_ = k
		if v != want[1] {
			t.Fatalf("pragma %s = %q, want %q", want[0], v, want[1])
		}
	}
}

func TestCreateDevice_AndGet(t *testing.T) {
	s := openTestStore(sT(t))
	id, err := s.CreateDevice(context.Background(), "d2000-lab-01", "yy0133")
	if err != nil {
		t.Fatalf("create device: %v", err)
	}
	if id <= 0 {
		t.Fatalf("id = %d, want positive", id)
	}
	dev, err := s.GetDevice(context.Background(), id)
	if err != nil {
		t.Fatalf("get device: %v", err)
	}
	if dev.Name != "d2000-lab-01" || dev.Project != "yy0133" {
		t.Fatalf("device = %+v", dev)
	}
}

func TestSessionRoundTrip_WithIndexes(t *testing.T) {
	s := openTestStore(sT(t))
	ctx := context.Background()
	devID, _ := s.CreateDevice(ctx, "dev", "proj")

	now := time.Date(2026, 9, 18, 10, 0, 0, 0, time.UTC)
	id, err := s.CreateSession(ctx, Session{
		DeviceID:  devID,
		Kind:      "manual",
		Owner:     "alice",
		State:     "active",
		StartedAt: now,
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	got, err := s.GetSession(ctx, id)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}
	if got.Kind != "manual" || got.Owner != "alice" || got.State != "active" {
		t.Fatalf("session = %+v", got)
	}

	// Finish it and verify the update persists.
	got.State = "closed"
	got.EndedAt = now.Add(time.Minute)
	got.EndReason = "closed by user"
	if err := s.UpdateSession(ctx, got); err != nil {
		t.Fatalf("update session: %v", err)
	}
	after, _ := s.GetSession(ctx, id)
	if after.State != "closed" || after.EndReason != "closed by user" {
		t.Fatalf("after update: %+v", after)
	}

	// Active index query: finished session must not appear.
	actives, err := s.ActiveSessions(ctx)
	if err != nil {
		t.Fatalf("active sessions: %v", err)
	}
	if len(actives) != 0 {
		t.Fatalf("active sessions = %d, want 0", len(actives))
	}
}

func TestActiveSessions_ReturnsActiveOnly(t *testing.T) {
	s := openTestStore(sT(t))
	ctx := context.Background()
	devID, _ := s.CreateDevice(ctx, "dev", "proj")
	now := time.Date(2026, 9, 18, 10, 0, 0, 0, time.UTC)

	activeID, _ := s.CreateSession(ctx, Session{DeviceID: devID, Kind: "task", Owner: "bob", State: "active", StartedAt: now})
	closedID, _ := s.CreateSession(ctx, Session{DeviceID: devID, Kind: "manual", Owner: "alice", State: "closed", StartedAt: now})

	actives, err := s.ActiveSessions(ctx)
	if err != nil {
		t.Fatalf("active sessions: %v", err)
	}
	if len(actives) != 1 || actives[0].ID != activeID {
		t.Fatalf("active = %+v, want only session %d (closed %d excluded)", actives, activeID, closedID)
	}
}

// Concurrent writers: all writes funnel through the single-writer queue
// (decision 10A) and none are lost.
func TestConcurrentWrites_SerialQueue_NoLostWrites(t *testing.T) {
	s := openTestStore(sT(t))
	ctx := context.Background()
	devID, _ := s.CreateDevice(ctx, "dev", "proj")

	const n = 32
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := s.CreateSession(ctx, Session{
				DeviceID:  devID,
				Kind:      "manual",
				Owner:     "u",
				State:     "closed",
				StartedAt: time.Now(),
			})
			if err != nil {
				t.Errorf("write %d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()

	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM sessions").Scan(&count); err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != n {
		t.Fatalf("rows = %d, want %d", count, n)
	}
}

func TestIndexes_ExistAtCreation(t *testing.T) {
	// CEO-8A: the sessions(device_id, started_at) composite index and the
	// sessions(status) index are created with the table, not later.
	s := openTestStore(sT(t))
	rows, err := s.db.Query("SELECT name FROM sqlite_master WHERE type='index' AND tbl_name='sessions' AND name LIKE 'idx_%'")
	if err != nil {
		t.Fatalf("query indexes: %v", err)
	}
	defer rows.Close()
	names := map[string]bool{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan: %v", err)
		}
		names[name] = true
	}
	for _, want := range []string{"idx_sessions_device_started", "idx_sessions_state"} {
		if !names[want] {
			t.Fatalf("missing index %q, have %v", want, names)
		}
	}
}

package store

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

// Reopen idempotency: every launch of an M3+ database re-runs migrate()
// (the ALTER-based device column upgrades), so a second Open must succeed
// -- the historical bug was "duplicate column name: ssh_host" killing the
// server on restart (requirement 3.2 upgrade path).
func TestMigrate_ReopenWithDeviceColumns(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "reopen.db")

	s1, err := Open(path)
	if err != nil {
		t.Fatalf("first open: %v", err)
	}
	if _, err := s1.CreateDevice(t.Context(), "board-a", "P0133"); err != nil {
		t.Fatalf("seed device: %v", err)
	}
	s1.Close()

	// Second open runs the ALTERs again: must not fail on existing columns.
	s2, err := Open(path)
	if err != nil {
		t.Fatalf("reopen (duplicate column regression): %v", err)
	}
	defer s2.Close()

	// Data survives; SSH columns usable through the normal surface.
	d, err := s2.GetDeviceFull(t.Context(), 1)
	if err != nil {
		t.Fatalf("get device after reopen: %v", err)
	}
	if d.Name != "board-a" || d.SSHPort != 22 || d.SshSet {
		t.Fatalf("device after reopen = %+v", d)
	}
}

// Pre-M3 databases upgrade in place: a devices table without the SSH
// columns gains them on Open (the original M3 ALTER intent).
func TestMigrate_UpgradeFromPreSSHSchema(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "legacy-devices.db")

	// Shape the devices table by hand WITHOUT the SSH columns (M1 era) --
	// Open() then must add them.
	raw, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open raw: %v", err)
	}
	if _, err := raw.Exec(`CREATE TABLE devices (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		project TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	)`); err != nil {
		t.Fatalf("create legacy devices: %v", err)
	}
	if _, err := raw.Exec(`INSERT INTO devices (name, project) VALUES ('old-board', 'P-old')`); err != nil {
		t.Fatalf("seed legacy device: %v", err)
	}
	raw.Close()

	s, err := Open(path)
	if err != nil {
		t.Fatalf("open over legacy schema: %v", err)
	}
	defer s.Close()

	d, err := s.GetDeviceFull(t.Context(), 1)
	if err != nil {
		t.Fatalf("get upgraded device: %v", err)
	}
	if d.Name != "old-board" || d.SSHHost != "" || d.SshSet {
		t.Fatalf("upgraded device = %+v", d)
	}
}

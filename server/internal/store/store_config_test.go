package store

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
)

// The unified config_objects migration (requirement 1.6: four rule
// families, one storage shape). Legacy-shape tests build a pre-M3 database
// by hand and verify Open() moves rows into config_objects with payloads
// the compat methods read back identically.

// openLegacyStore creates a store file with the pre-M3 schema and seed
// rows, WITHOUT running the new migrations (raw sql.DB, no Store).
func openLegacyStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "legacy.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open legacy: %v", err)
	}
	defer db.Close()
	const ddl = `
CREATE TABLE projects (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	name       TEXT NOT NULL UNIQUE,
	note       TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE TABLE build_items (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	project_id   INTEGER NOT NULL REFERENCES projects(id),
	name         TEXT NOT NULL,
	source_type  TEXT NOT NULL,
	git_url      TEXT NOT NULL DEFAULT '',
	git_branch   TEXT NOT NULL DEFAULT '',
	git_commit   TEXT NOT NULL DEFAULT '',
	check_latest INTEGER NOT NULL DEFAULT 0,
	local_path   TEXT NOT NULL DEFAULT '',
	command      TEXT NOT NULL,
	artifacts_json TEXT NOT NULL DEFAULT '[]',
	timeout_sec  INTEGER NOT NULL DEFAULT 0,
	version_cmd  TEXT NOT NULL DEFAULT '',
	prereq_json  TEXT NOT NULL DEFAULT '[]',
	created_at   TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
	updated_at   TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE TABLE expect_rules (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	name       TEXT NOT NULL,
	steps_json TEXT NOT NULL,
	updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
INSERT INTO projects (name, note) VALUES ('P0133', 'legacy line');
INSERT INTO build_items (project_id, name, source_type, local_path, command, artifacts_json, timeout_sec, version_cmd, prereq_json)
	VALUES (1, 'bootfw', 'local', '/root/git/bootfw', './bb.sh', '["bootfw*.bin"]', 600, 'ver.sh -a', '[[2]]');
INSERT INTO build_items (project_id, name, source_type, git_url, git_branch, command, artifacts_json)
	VALUES (1, 'u-boot', 'git', 'http://git.local/uboot', 'release', 'make', '["u-boot.bin"]');
INSERT INTO expect_rules (name, steps_json) VALUES ('demo-flash', '[{"await":"login:","send":"root\r"}]');
`
	if _, err := db.Exec(ddl); err != nil {
		t.Fatalf("legacy ddl: %v", err)
	}
	return path
}

func TestMigrate_LegacyTablesIntoConfigObjects(t *testing.T) {
	path := openLegacyStore(t)
	s, err := Open(path)
	if err != nil {
		t.Fatalf("open migrated: %v", err)
	}
	defer s.Close()
	ctx := context.Background()

	// Build items read back through the compat surface with all fields.
	items, err := s.ListBuildItemsByProject(ctx, 1)
	if err != nil {
		t.Fatalf("list migrated items: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("migrated items = %d, want 2", len(items))
	}
	var bootfw, uboot *BuildItem
	for i := range items {
		switch items[i].Name {
		case "bootfw":
			bootfw = &items[i]
		case "u-boot":
			uboot = &items[i]
		}
	}
	if bootfw == nil || uboot == nil {
		t.Fatalf("missing migrated items: %+v", items)
	}
	if bootfw.SourceType != "local" || bootfw.LocalPath != "/root/git/bootfw" ||
		bootfw.Command != "./bb.sh" || len(bootfw.Artifacts) != 1 ||
		bootfw.Artifacts[0] != "bootfw*.bin" || bootfw.TimeoutSec != 600 ||
		bootfw.VersionCmd != "ver.sh -a" || bootfw.PrereqJSON != "[[2]]" {
		t.Fatalf("bootfw payload mismatch: %+v", bootfw)
	}
	if uboot.SourceType != "git" || uboot.GitURL != "http://git.local/uboot" || uboot.GitBranch != "release" {
		t.Fatalf("uboot payload mismatch: %+v", uboot)
	}

	// Expect rules land as flash rules attached to a project (the demo
	// project is created as their home when absent).
	rules, err := s.ListExpectRules(ctx)
	if err != nil {
		t.Fatalf("list migrated rules: %v", err)
	}
	if len(rules) != 1 || rules[0].Name != "demo-flash" || rules[0].StepsJSON == "" {
		t.Fatalf("migrated rules = %+v", rules)
	}
	obj, err := s.GetConfigObject(ctx, rules[0].ID)
	if err != nil {
		t.Fatalf("get rule object: %v", err)
	}
	if obj.Kind != KindFlashRule || obj.ProjectID == 0 {
		t.Fatalf("flash rule object = %+v", obj)
	}

	// New ids continue after the migrated ones (no AUTOINCREMENT reset).
	newID, err := s.CreateBuildItem(ctx, BuildItem{ProjectID: 1, Name: "post-mig", SourceType: "local", Command: "x"})
	if err != nil {
		t.Fatalf("post-migration insert: %v", err)
	}
	if newID <= 2 {
		t.Fatalf("new id %d reuses the migrated id range (1..2)", newID)
	}

	// Legacy tables stay as backups, untouched.
	var n int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM build_items").Scan(&n); err != nil || n != 2 {
		t.Fatalf("legacy build_items backup rows = %d, err=%v, want 2", n, err)
	}
}

func TestConfigObjects_CRUDRoundTrip(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	pid, _ := s.CreateProject(ctx, "P1", "")

	id, err := s.CreateConfigObject(ctx, ConfigObject{
		ProjectID: pid, Kind: KindDeployRule, Name: "flash-via-tftp", PayloadJSON: `{"mode":"flash"}`,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := s.GetConfigObject(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Kind != KindDeployRule || got.PayloadJSON != `{"mode":"flash"}` {
		t.Fatalf("round trip = %+v", got)
	}

	got.PayloadJSON = `{"mode":"flash","tftp":"/srv/tftp"}`
	if err := s.UpdateConfigObject(ctx, got); err != nil {
		t.Fatalf("update: %v", err)
	}
	after, _ := s.GetConfigObject(ctx, id)
	if after.PayloadJSON != `{"mode":"flash","tftp":"/srv/tftp"}` {
		t.Fatalf("update lost: %+v", after)
	}

	byKind, err := s.ListConfigObjectsByKind(ctx, KindDeployRule)
	if err != nil || len(byKind) != 1 {
		t.Fatalf("list by kind = %v, err %v", byKind, err)
	}
	byIDs, err := s.ListConfigObjectsByIDs(ctx, KindDeployRule, []int64{id, 999})
	if err != nil || len(byIDs) != 1 {
		t.Fatalf("list by ids = %v, err %v", byIDs, err)
	}

	if err := s.DeleteConfigObject(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := s.GetConfigObject(ctx, id); err == nil {
		t.Fatal("deleted object still readable")
	}
}

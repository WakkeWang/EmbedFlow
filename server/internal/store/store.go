// Package store is the persistence layer: SQLite with WAL mode, a
// busy_timeout, and all writes funneled through a single writer goroutine
// (decision 10A) so concurrent callers never fight over the write lock.
package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// Session is a persisted session record. Field vocabulary matches the
// session kernel (CONTEXT.md: 会话).
type Session struct {
	ID            int64     `json:"id"`
	DeviceID      int64     `json:"device_id"`
	Kind          string    `json:"kind"` // manual | task
	Owner         string    `json:"owner"`
	State         string    `json:"state"` // active | closed | failed
	StartedAt     time.Time `json:"started_at"`
	EndedAt       time.Time `json:"ended_at"`
	EndReason     string    `json:"end_reason"`
	LogIncomplete bool      `json:"log_incomplete"`
}

// ExpectRule is the M1 view of a flash rule (now stored as a config
// object, kind=flash_rule -- see store_config.go).
type ExpectRule struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	StepsJSON string    `json:"steps_json"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Device is a registered target machine (CONTEXT.md: 设备) -- the SSH-field
// variant lives in store_ext.go (M3).

// Store wraps the SQLite handle plus the serializing write queue.
type Store struct {
	db      *sql.DB
	writeCh chan writeOp
	done    chan struct{}
}

type writeOp struct {
	fn   func() error
	errc chan error
}

// Open creates/opens the database file, applies pragmas, runs migrations,
// and starts the single-writer goroutine.
func Open(path string) (*Store, error) {
	// modernc driver: pass pragmas via DSN is unreliable; apply after open.
	db, err := sql.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)")
	if err != nil {
		return nil, fmt.Errorf("store: open %s: %w", path, err)
	}
	// Single writer means the write pool needs exactly one connection.
	db.SetMaxOpenConns(2) // 1 writer + 1 reader (WAL allows concurrent reads)
	db.SetMaxIdleConns(2)

	s := &Store{
		db:      db,
		writeCh: make(chan writeOp, 256),
		done:    make(chan struct{}),
	}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	go s.writerLoop()
	return s, nil
}

// migrate creates tables and indexes in one shot (CEO-8A: indexes are born
// with the tables), then runs the data migrations.
func (s *Store) migrate() error {
	const ddl = `
CREATE TABLE IF NOT EXISTS devices (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	name       TEXT NOT NULL,
	project    TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- M3 (requirement 3.2): SSH credentials, encrypted at rest (requirement
-- 6.2, secretbox AES-256-GCM). ALTER-based so existing databases upgrade.
ALTER TABLE devices ADD COLUMN ssh_host TEXT NOT NULL DEFAULT '';
ALTER TABLE devices ADD COLUMN ssh_port INTEGER NOT NULL DEFAULT 22;
ALTER TABLE devices ADD COLUMN ssh_user TEXT NOT NULL DEFAULT '';
ALTER TABLE devices ADD COLUMN ssh_pass_enc TEXT NOT NULL DEFAULT '';
ALTER TABLE devices ADD COLUMN note TEXT NOT NULL DEFAULT '';

CREATE TABLE IF NOT EXISTS projects (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	name       TEXT NOT NULL UNIQUE,
	note       TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE TABLE IF NOT EXISTS sessions (
	id             INTEGER PRIMARY KEY AUTOINCREMENT,
	device_id      INTEGER NOT NULL REFERENCES devices(id),
	kind           TEXT NOT NULL CHECK (kind IN ('manual','task')),
	owner          TEXT NOT NULL,
	state          TEXT NOT NULL CHECK (state IN ('active','closed','failed')),
	started_at     TEXT NOT NULL,
	ended_at       TEXT NOT NULL DEFAULT '',
	end_reason     TEXT NOT NULL DEFAULT '',
	log_incomplete INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_sessions_device_started ON sessions(device_id, started_at);
CREATE INDEX IF NOT EXISTS idx_sessions_state ON sessions(state);

CREATE TABLE IF NOT EXISTS users (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	username     TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	role         TEXT NOT NULL CHECK (role IN ('admin','member')),
	created_at   TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE TABLE IF NOT EXISTS config_objects (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	project_id   INTEGER NOT NULL REFERENCES projects(id),
	kind         TEXT NOT NULL CHECK (kind IN ('build_item','deploy_rule','flash_rule','test_unit','test_item','test_suite','release_rule')),
	name         TEXT NOT NULL,
	payload_json TEXT NOT NULL DEFAULT '{}',
	created_at   TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
	updated_at   TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE INDEX IF NOT EXISTS idx_config_objects_project_kind ON config_objects(project_id, kind);

CREATE TABLE IF NOT EXISTS confirmations (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	session_id INTEGER NOT NULL REFERENCES sessions(id),
	prompt     TEXT NOT NULL,
	state      TEXT NOT NULL DEFAULT 'pending' CHECK (state IN ('pending','resolved')),
	result     TEXT NOT NULL DEFAULT '',
	note       TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
	resolved_at TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_confirmations_session ON confirmations(session_id);

CREATE TABLE IF NOT EXISTS batches (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	project_id INTEGER NOT NULL REFERENCES projects(id),
	status     TEXT NOT NULL CHECK (status IN ('queued','running','completed','failed','canceled')),
	items_json TEXT NOT NULL DEFAULT '[]',
	created_by TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
	started_at TEXT NOT NULL DEFAULT '',
	ended_at   TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_batches_project ON batches(project_id);

CREATE TABLE IF NOT EXISTS build_records (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	batch_id     INTEGER NOT NULL REFERENCES batches(id),
	item_id      INTEGER NOT NULL,
	project_id   INTEGER NOT NULL,
	status       TEXT NOT NULL CHECK (status IN ('pending','building','succeeded','failed','canceled','skipped')),
	commit_sha   TEXT NOT NULL DEFAULT '',
	version_info TEXT NOT NULL DEFAULT '',
	exit_code    INTEGER,
	executor     TEXT NOT NULL DEFAULT '',
	started_at   TEXT NOT NULL DEFAULT '',
	ended_at     TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_build_records_batch ON build_records(batch_id);
CREATE INDEX IF NOT EXISTS idx_build_records_project ON build_records(project_id);

CREATE TABLE IF NOT EXISTS artifacts (
	id              INTEGER PRIMARY KEY AUTOINCREMENT,
	build_record_id INTEGER NOT NULL REFERENCES build_records(id),
	name            TEXT NOT NULL,
	size            INTEGER NOT NULL DEFAULT 0,
	checksum        TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_artifacts_record ON artifacts(build_record_id);

CREATE TABLE IF NOT EXISTS settings (
	key   TEXT PRIMARY KEY,
	value TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS deploy_records (
	id              INTEGER PRIMARY KEY AUTOINCREMENT,
	project_id      INTEGER NOT NULL,
	rule_id         INTEGER NOT NULL,
	device_id       INTEGER NOT NULL,
	build_record_id INTEGER NOT NULL DEFAULT 0,
	executor        TEXT NOT NULL DEFAULT '',
	status          TEXT NOT NULL CHECK (status IN ('running','succeeded','failed','canceled')),
	detail          TEXT NOT NULL DEFAULT '',
	exec_record_id  INTEGER NOT NULL DEFAULT 0,
	session_id      INTEGER NOT NULL DEFAULT 0,
	started_at      TEXT NOT NULL DEFAULT '',
	ended_at        TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_deploy_records_project ON deploy_records(project_id);

CREATE TABLE IF NOT EXISTS exec_records (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	device_id  INTEGER NOT NULL,
	kind       TEXT NOT NULL DEFAULT 'deploy',
	command    TEXT NOT NULL,
	status     TEXT NOT NULL CHECK (status IN ('running','succeeded','failed')),
	exit_code  INTEGER,
	detail     TEXT NOT NULL DEFAULT '',
	started_at TEXT NOT NULL,
	ended_at   TEXT NOT NULL DEFAULT ''
);
`
	_, err := s.db.Exec(ddl)
	if err != nil {
		return fmt.Errorf("store: migrate: %w", err)
	}
	// Data migrations: legacy tables into config_objects. The legacy
	// build_items / expect_rules tables stay as untouched backups after the
	// move (no drop -- a rollback-friendly copy).
	if err := s.migrateLegacyTables(); err != nil {
		return fmt.Errorf("store: legacy migration: %w", err)
	}
	return nil
}

// migrateLegacyTables moves rows from the pre-M3 build_items / expect_rules
// tables into config_objects. Runs once per database (settings marker); the
// legacy tables stay as untouched backups after the move (no drop -- a
// rollback-friendly copy).
func (s *Store) migrateLegacyTables() error {
	// One-shot guard: direct SQL (migrate runs before the writer loop).
	var done int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM settings WHERE key = 'legacy_migrated_v1'").Scan(&done); err == nil && done > 0 {
		return nil
	}

	// build_items -> config_objects (kind=build_item, ids preserved: batch
	// history's items_json and build_records.item_id reference them).
	has, err := s.tableExists("build_items")
	if err != nil {
		return err
	}
	if has {
		rows, err := s.db.Query("SELECT id, project_id, name, source_type, git_url, git_branch, git_commit, check_latest, local_path, command, artifacts_json, timeout_sec, version_cmd, prereq_json FROM build_items")
		if err != nil {
			return err
		}
		type itemRow struct {
			id                                int64
			projectID                         int64
			name, sourceType, gitURL          string
			gitBranch, gitCommit              string
			checkLatest                       bool
			localPath, command, artifactsJSON string
			timeoutSec                        int
			versionCmd, prereqJSON            string
		}
		var legacyItems []itemRow
		for rows.Next() {
			var r itemRow
			if err := rows.Scan(&r.id, &r.projectID, &r.name, &r.sourceType, &r.gitURL, &r.gitBranch,
				&r.gitCommit, &r.checkLatest, &r.localPath, &r.command, &r.artifactsJSON,
				&r.timeoutSec, &r.versionCmd, &r.prereqJSON); err == nil {
				legacyItems = append(legacyItems, r)
			}
		}
		rows.Close()
		for _, r := range legacyItems {
			payload, _ := json.Marshal(BuildItemPayload{
				SourceType: r.sourceType, GitURL: r.gitURL, GitBranch: r.gitBranch,
				GitCommit: r.gitCommit, CheckLatest: r.checkLatest, LocalPath: r.localPath,
				Command: r.command, Artifacts: decodeJSONStrings(r.artifactsJSON),
				TimeoutSec: r.timeoutSec, VersionCmd: r.versionCmd, PrereqJSON: r.prereqJSON,
			})
			if _, err := s.db.Exec(
				"INSERT OR IGNORE INTO config_objects (id, project_id, kind, name, payload_json) VALUES (?, ?, ?, ?, ?)",
				r.id, r.projectID, KindBuildItem, r.name, string(payload)); err != nil {
				return err
			}
		}
	}

	// expect_rules -> config_objects (kind=flash_rule, attached to the demo
	// project; M1 rules were global, demo is their de-facto home). Ids are
	// NOT preserved: nothing references rule ids persistently, and the
	// legacy id space collides with build items'.
	hasRules, err := s.tableExists("expect_rules")
	if err != nil {
		return err
	}
	if hasRules {
		rows, err := s.db.Query("SELECT name, steps_json FROM expect_rules")
		if err != nil {
			return err
		}
		type ruleRow struct {
			name      string
			stepsJSON string
		}
		var legacyRules []ruleRow
		for rows.Next() {
			var r ruleRow
			if err := rows.Scan(&r.name, &r.stepsJSON); err == nil {
				legacyRules = append(legacyRules, r)
			}
		}
		rows.Close()
		for _, r := range legacyRules {
			payload, _ := json.Marshal(FlashRulePayload{StepsJSON: r.stepsJSON})
			// Idempotent within a run AND across a partial run: skip when
			// this rule name already exists as a flash rule.
			if _, err := s.db.Exec(
				"INSERT INTO config_objects (project_id, kind, name, payload_json) SELECT ?, ?, ?, ? WHERE NOT EXISTS (SELECT 1 FROM config_objects WHERE kind = ? AND name = ?)",
				s.legacyFlashProjectID(), KindFlashRule, r.name, string(payload), KindFlashRule, r.name); err != nil {
				return err
			}
		}
	}

	// Keep AUTOINCREMENT above every migrated id (new rows never reuse one).
	if _, err := s.db.Exec(
		"UPDATE sqlite_sequence SET seq = (SELECT MAX(id) FROM config_objects) WHERE name = 'config_objects' AND seq < (SELECT MAX(id) FROM config_objects)"); err != nil {
		return err
	}
	_, err = s.db.Exec("INSERT INTO settings (key, value) VALUES ('legacy_migrated_v1', '1')")
	return err
}

// legacyFlashProjectID resolves the project flash rules attach to during
// migration (the demo project; created if absent).
func (s *Store) legacyFlashProjectID() int64 {
	var id int64
	err := s.db.QueryRow("SELECT id FROM projects WHERE name = 'demo'").Scan(&id)
	if err == nil {
		return id
	}
	res, err := s.db.Exec("INSERT INTO projects (name, note) VALUES ('demo', 'built-in demo project')")
	if err != nil {
		return 1 // config_objects.project_id has no hard FK enforcement at migration time
	}
	id, _ = res.LastInsertId()
	return id
}

// tableExists reports whether a table exists (legacy migration guard).
func (s *Store) tableExists(name string) (bool, error) {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name = ?", name).Scan(&n)
	return n > 0, err
}

// writerLoop is the single writer goroutine: it owns all mutations so SQLite
// never sees concurrent write transactions (decision 10A).
func (s *Store) writerLoop() {
	defer close(s.done)
	for op := range s.writeCh {
		op.errc <- op.fn()
	}
}

// enqueue runs a write on the writer goroutine.
func (s *Store) enqueue(ctx context.Context, fn func() error) error {
	op := writeOp{fn: fn, errc: make(chan error, 1)}
	select {
	case s.writeCh <- op:
	case <-ctx.Done():
		return ctx.Err()
	}
	select {
	case err := <-op.errc:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Close stops the writer and closes the database.
func (s *Store) Close() error {
	if s.writeCh != nil {
		// Drain in-flight ops, then stop the loop.
		close(s.writeCh)
		<-s.done
		s.writeCh = nil
	}
	return s.db.Close()
}

// CreateDevice registers a device (name + project only; SSH fields ride
// UpdateDeviceSSH -- M3).
func (s *Store) CreateDevice(ctx context.Context, name, project string) (int64, error) {
	var id int64
	err := s.enqueue(ctx, func() error {
		res, err := s.db.ExecContext(ctx, "INSERT INTO devices (name, project) VALUES (?, ?)", name, project)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		return err
	})
	return id, err
}

// GetDevice fetches one device (SSH columns included in store_ext.go's
// full-scan variant; this row-level one stays minimal for the session kernel).
func (s *Store) GetDevice(ctx context.Context, id int64) (Device, error) {
	var d Device
	err := s.db.QueryRowContext(ctx, "SELECT id, name, project FROM devices WHERE id = ?", id).
		Scan(&d.ID, &d.Name, &d.Project)
	if err != nil {
		return Device{}, fmt.Errorf("store: get device %d: %w", id, err)
	}
	return d, nil
}

// CreateSession inserts a session record. The kernel owns ID assignment;
// the store must persist exactly that ID so log paths, HTTP actions and
// restart sweeps all address the same row (a plain autoincrement drifts
// from the kernel's counter after any restart).
func (s *Store) CreateSession(ctx context.Context, sess Session) (int64, error) {
	err := s.enqueue(ctx, func() error {
		_, err := s.db.ExecContext(ctx,
			"INSERT INTO sessions (id, device_id, kind, owner, state, started_at, ended_at, end_reason) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			sess.ID,
			sess.DeviceID, sess.Kind, sess.Owner, sess.State,
			sess.StartedAt.Format(time.RFC3339Nano),
			sess.EndedAt.Format(time.RFC3339Nano),
			sess.EndReason)
		if err != nil {
			return err
		}
		return nil
	})
	return sess.ID, err
}

// GetSession fetches one session record.
func (s *Store) GetSession(ctx context.Context, id int64) (Session, error) {
	var sess Session
	var started, ended string
	err := s.db.QueryRowContext(ctx,
		"SELECT id, device_id, kind, owner, state, started_at, ended_at, end_reason, log_incomplete FROM sessions WHERE id = ?", id).
		Scan(&sess.ID, &sess.DeviceID, &sess.Kind, &sess.Owner, &sess.State, &started, &ended, &sess.EndReason, &sess.LogIncomplete)
	if err != nil {
		return Session{}, fmt.Errorf("store: get session %d: %w", id, err)
	}
	sess.StartedAt = parseTime(started)
	sess.EndedAt = parseTime(ended)
	return sess, nil
}

// UpdateSession persists mutable session fields (state, end time, reason,
// log-incomplete flag).
func (s *Store) UpdateSession(ctx context.Context, sess Session) error {
	return s.enqueue(ctx, func() error {
		_, err := s.db.ExecContext(ctx,
			"UPDATE sessions SET state = ?, ended_at = ?, end_reason = ?, log_incomplete = ? WHERE id = ?",
			sess.State,
			sess.EndedAt.Format(time.RFC3339Nano),
			sess.EndReason,
			boolToInt(sess.LogIncomplete),
			sess.ID)
		return err
	})
}

// MarkLogIncomplete flags a session's log as having had write failures
// (design doc Failure modes: the disk-full gap). Callable mid-session, the
// moment a write breaks -- not only at session end.
func (s *Store) MarkLogIncomplete(ctx context.Context, id int64) error {
	return s.enqueue(ctx, func() error {
		_, err := s.db.ExecContext(ctx, "UPDATE sessions SET log_incomplete = 1 WHERE id = ?", id)
		return err
	})
}

// ActiveSessions lists sessions in the active state -- the startup sweep's
// input (CEO-7A). Served by idx_sessions_state.
func (s *Store) ActiveSessions(ctx context.Context) ([]Session, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, device_id, kind, owner, state, started_at, ended_at, end_reason, log_incomplete FROM sessions WHERE state = 'active' ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("store: active sessions: %w", err)
	}
	defer rows.Close()
	var out []Session
	for rows.Next() {
		var sess Session
		var started, ended string
		if err := rows.Scan(&sess.ID, &sess.DeviceID, &sess.Kind, &sess.Owner, &sess.State, &started, &ended, &sess.EndReason, &sess.LogIncomplete); err != nil {
			return nil, err
		}
		sess.StartedAt = parseTime(started)
		sess.EndedAt = parseTime(ended)
		out = append(out, sess)
	}
	return out, rows.Err()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

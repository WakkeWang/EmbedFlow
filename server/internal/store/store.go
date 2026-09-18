// Package store is the persistence layer: SQLite with WAL mode, a
// busy_timeout, and all writes funneled through a single writer goroutine
// (decision 10A) so concurrent callers never fight over the write lock.
package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// Session is a persisted session record. Field vocabulary matches the
// session kernel (CONTEXT.md: 会话).
type Session struct {
	ID        int64
	DeviceID  int64
	Kind      string // manual | task
	Owner     string
	State     string // active | closed | failed
	StartedAt time.Time
	EndedAt   time.Time
	EndReason string
	// LogIncomplete marks a session whose log had write failures (design doc
	// Failure modes: the disk-full gap) -- surfaced in UI and reports.
	LogIncomplete bool
}

// Device is a registered target machine (CONTEXT.md: 设备).
type Device struct {
	ID      int64
	Name    string
	Project string
}

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
// with the tables).
func (s *Store) migrate() error {
	const ddl = `
CREATE TABLE IF NOT EXISTS devices (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	name       TEXT NOT NULL,
	project    TEXT NOT NULL,
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

CREATE TABLE IF NOT EXISTS expect_rules (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	name       TEXT NOT NULL,
	steps_json TEXT NOT NULL,
	updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

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
`
	_, err := s.db.Exec(ddl)
	if err != nil {
		return fmt.Errorf("store: migrate: %w", err)
	}
	return nil
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

// CreateDevice registers a device, returning its id.
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

// GetDevice fetches one device.
func (s *Store) GetDevice(ctx context.Context, id int64) (Device, error) {
	var d Device
	err := s.db.QueryRowContext(ctx, "SELECT id, name, project FROM devices WHERE id = ?", id).
		Scan(&d.ID, &d.Name, &d.Project)
	if err != nil {
		return Device{}, fmt.Errorf("store: get device %d: %w", id, err)
	}
	return d, nil
}

// CreateSession inserts a session record.
func (s *Store) CreateSession(ctx context.Context, sess Session) (int64, error) {
	var id int64
	err := s.enqueue(ctx, func() error {
		res, err := s.db.ExecContext(ctx,
			"INSERT INTO sessions (device_id, kind, owner, state, started_at, ended_at, end_reason) VALUES (?, ?, ?, ?, ?, ?, ?)",
			sess.DeviceID, sess.Kind, sess.Owner, sess.State,
			sess.StartedAt.Format(time.RFC3339Nano),
			sess.EndedAt.Format(time.RFC3339Nano),
			sess.EndReason)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		return err
	})
	return id, err
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

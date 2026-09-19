package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User is an account (requirement 1.5): admin manages config, members
// execute operations. Passwords are bcrypt-hashed, never stored raw.
type User struct {
	ID       int64
	Username string
	Role     string // admin | member
}

// Project is the top-level configuration container (requirement 1.6):
// build items, deploy rules, devices, tests and release rules all live
// inside one. M1 seeds the entity; its sections land M2+.
type Project struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateProject inserts a project, returning its id.
func (s *Store) CreateProject(ctx context.Context, name, note string) (int64, error) {
	var id int64
	err := s.enqueue(ctx, func() error {
		res, err := s.db.ExecContext(ctx,
			"INSERT INTO projects (name, note) VALUES (?, ?)", name, note)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		return err
	})
	return id, err
}

// ListProjects returns all projects in insertion order.
func (s *Store) ListProjects(ctx context.Context) ([]Project, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, note, created_at FROM projects ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("store: list projects: %w", err)
	}
	defer rows.Close()
	var out []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Note, &p.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// DeleteProject removes a project record (sections M2+ cascade by design
// when they hang entities off project_id).
func (s *Store) DeleteProject(ctx context.Context, id int64) error {
	return s.enqueue(ctx, func() error {
		_, err := s.db.ExecContext(ctx, "DELETE FROM projects WHERE id = ?", id)
		return err
	})
}

// ExpectRule is a stored expect sequence (issue #7): team-shared,
// reusable, reviewable configuration.
type ExpectRule struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	StepsJSON string    `json:"steps_json"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Confirmation is one human confirmation card (issue #9).
type Confirmation struct {
	ID         int64     `json:"id"`
	SessionID  int64     `json:"session_id"`
	Prompt     string    `json:"prompt"`
	State      string    `json:"state"`  // pending | resolved
	Result     string    `json:"result"` // pass | fail
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
	ResolvedAt time.Time `json:"resolved_at"`
}

// CreateUser inserts an account with a bcrypt-hashed password.
func (s *Store) CreateUser(ctx context.Context, username, password, role string) error {
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}
	return s.enqueue(ctx, func() error {
		_, err := s.db.ExecContext(ctx,
			"INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)", username, string(hash), role)
		return err
	})
}

// ErrBadCredentials is returned for unknown users or wrong passwords
// (indistinguishable on purpose; CEO-4A: no lockout, README warns deployers).
var ErrBadCredentials = errors.New("store: bad credentials")

// AuthenticateUser verifies credentials and returns the account.
func (s *Store) AuthenticateUser(ctx context.Context, username, password string) (User, error) {
	var (
		u    User
		hash string
	)
	err := s.db.QueryRowContext(ctx,
		"SELECT id, username, role, password_hash FROM users WHERE username = ?", username).
		Scan(&u.ID, &u.Username, &u.Role, &hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrBadCredentials
		}
		return User{}, fmt.Errorf("store: find user: %w", err)
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return User{}, ErrBadCredentials
	}
	return u, nil
}

// ListDevices returns all registered devices in insertion order.
func (s *Store) ListDevices(ctx context.Context) ([]Device, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, project FROM devices ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("store: list devices: %w", err)
	}
	defer rows.Close()
	var out []Device
	for rows.Next() {
		var d Device
		if err := rows.Scan(&d.ID, &d.Name, &d.Project); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// DeleteDevice removes a device registration.
func (s *Store) DeleteDevice(ctx context.Context, id int64) error {
	return s.enqueue(ctx, func() error {
		_, err := s.db.ExecContext(ctx, "DELETE FROM devices WHERE id = ?", id)
		return err
	})
}

// CreateExpectRule stores a rule; steps_json carries the serialized steps.
func (s *Store) CreateExpectRule(ctx context.Context, name, stepsJSON string) (int64, error) {
	var id int64
	err := s.enqueue(ctx, func() error {
		res, err := s.db.ExecContext(ctx,
			"INSERT INTO expect_rules (name, steps_json) VALUES (?, ?)", name, stepsJSON)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		return err
	})
	return id, err
}

// GetExpectRule fetches one rule.
func (s *Store) GetExpectRule(ctx context.Context, id int64) (ExpectRule, error) {
	var r ExpectRule
	var updated string
	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, steps_json, updated_at FROM expect_rules WHERE id = ?", id).
		Scan(&r.ID, &r.Name, &r.StepsJSON, &updated)
	if err != nil {
		return ExpectRule{}, fmt.Errorf("store: get rule %d: %w", id, err)
	}
	r.UpdatedAt = parseTime(updated)
	return r, nil
}

// UpdateExpectRule overwrites name and steps.
func (s *Store) UpdateExpectRule(ctx context.Context, r ExpectRule) error {
	return s.enqueue(ctx, func() error {
		_, err := s.db.ExecContext(ctx,
			"UPDATE expect_rules SET name = ?, steps_json = ?, updated_at = strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id = ?",
			r.Name, r.StepsJSON, r.ID)
		return err
	})
}

// ListExpectRules returns all rules in insertion order.
func (s *Store) ListExpectRules(ctx context.Context) ([]ExpectRule, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, steps_json, updated_at FROM expect_rules ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("store: list rules: %w", err)
	}
	defer rows.Close()
	var out []ExpectRule
	for rows.Next() {
		var r ExpectRule
		var updated string
		if err := rows.Scan(&r.ID, &r.Name, &r.StepsJSON, &updated); err != nil {
			return nil, err
		}
		r.UpdatedAt = parseTime(updated)
		out = append(out, r)
	}
	return out, rows.Err()
}

// DeleteExpectRule removes a rule (the editor's delete button).
func (s *Store) DeleteExpectRule(ctx context.Context, id int64) error {
	return s.enqueue(ctx, func() error {
		_, err := s.db.ExecContext(ctx, "DELETE FROM expect_rules WHERE id = ?", id)
		return err
	})
}

// InsertConfirmation adds a pending confirmation card to a session.
func (s *Store) InsertConfirmation(ctx context.Context, sessionID int64, prompt string) (int64, error) {
	var id int64
	err := s.enqueue(ctx, func() error {
		res, err := s.db.ExecContext(ctx,
			"INSERT INTO confirmations (session_id, prompt) VALUES (?, ?)", sessionID, prompt)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		return err
	})
	return id, err
}

// ResolveConfirmation records the human verdict on a card (issue #9: cards
// never auto-dismiss; only an explicit resolve closes them).
func (s *Store) ResolveConfirmation(ctx context.Context, id int64, result, note string) error {
	return s.enqueue(ctx, func() error {
		_, err := s.db.ExecContext(ctx,
			"UPDATE confirmations SET state = 'resolved', result = ?, note = ?, resolved_at = strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id = ?",
			result, note, id)
		return err
	})
}

// ConfirmationsForSession lists a session's cards in creation order.
func (s *Store) ConfirmationsForSession(ctx context.Context, sessionID int64) ([]Confirmation, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, session_id, prompt, state, result, note, created_at, resolved_at FROM confirmations WHERE session_id = ? ORDER BY id", sessionID)
	if err != nil {
		return nil, fmt.Errorf("store: confirmations: %w", err)
	}
	defer rows.Close()
	var out []Confirmation
	for rows.Next() {
		var c Confirmation
		var created, resolved string
		if err := rows.Scan(&c.ID, &c.SessionID, &c.Prompt, &c.State, &c.Result, &c.Note, &created, &resolved); err != nil {
			return nil, err
		}
		c.CreatedAt = parseTime(created)
		c.ResolvedAt = parseTime(resolved)
		out = append(out, c)
	}
	return out, rows.Err()
}

// MaxSessionID returns the highest session ID ever persisted -- the startup
// sweep seeds the kernel counter from it so new IDs never collide with
// finished rows (plain autoincrement inside the kernel would drift).
func (s *Store) MaxSessionID(ctx context.Context) (int64, error) {
	var max sql.NullInt64
	err := s.db.QueryRowContext(ctx, "SELECT MAX(id) FROM sessions").Scan(&max)
	if err != nil {
		return 0, fmt.Errorf("store: max session id: %w", err)
	}
	return max.Int64, nil
}

// SessionsForDevice lists a device's session history, oldest first
// (requirement 3.3: history keyed by device). Served by
// idx_sessions_device_started.
func (s *Store) SessionsForDevice(ctx context.Context, deviceID int64) ([]Session, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, device_id, kind, owner, state, started_at, ended_at, end_reason, log_incomplete FROM sessions WHERE device_id = ? ORDER BY started_at, id", deviceID)
	if err != nil {
		return nil, fmt.Errorf("store: sessions for device: %w", err)
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

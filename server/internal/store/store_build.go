// Build-module persistence (requirement 2): build items (the M2 member of
// the four rule families), batches, build records, artifacts and the
// system settings key-value table. Same write-queue discipline as the rest
// of the store: mutations go through enqueue, reads hit WAL directly.
package store

import (
	"context"
	"database/sql"
	"fmt"
)

// BuildItem is one buildable unit inside a project (CONTEXT.md: 构建项目).
// SourceType picks which source fields apply; PrereqGroups holds prerequisite
// groupings -- group-internal OR, group-to-group AND (requirement 2.1).
// Persistence rides config_objects (kind=build_item); this struct stays the
// API/exec surface.
type BuildItem struct {
	ID          int64    `json:"id"`
	ProjectID   int64    `json:"project_id"`
	Name        string   `json:"name"`
	SourceType  string   `json:"source_type"` // git | local
	GitURL      string   `json:"git_url,omitempty"`
	GitBranch   string   `json:"git_branch,omitempty"`
	GitCommit   string   `json:"git_commit,omitempty"`
	CheckLatest bool     `json:"check_latest,omitempty"`
	LocalPath   string   `json:"local_path,omitempty"`
	Command     string   `json:"command"`
	Artifacts   []string `json:"artifacts"` // relative globs against the build dir
	TimeoutSec  int      `json:"timeout_sec,omitempty"`
	VersionCmd  string   `json:"version_cmd,omitempty"`
	// PrereqGroups is a JSON-encoded array of id arrays: [[1,2],[3]] means
	// "some of {1,2} AND some of {3}". Kept as raw JSON so the client edits
	// the structure directly (M3 unification may refine the shape).
	PrereqJSON string `json:"prereq_json"`
}

// Batch is one "build selected items" action (CONTEXT.md: 批次).
type Batch struct {
	ID        int64  `json:"id"`
	ProjectID int64  `json:"project_id"`
	Status    string `json:"status"` // queued | running | completed | failed | canceled
	ItemsJSON string `json:"items_json"`
	CreatedBy string `json:"created_by"`
}

// BuildRecord is one item's execution inside a batch (CONTEXT.md: 构建记录).
type BuildRecord struct {
	ID          int64   `json:"id"`
	BatchID     int64   `json:"batch_id"`
	ItemID      int64   `json:"item_id"`
	ProjectID   int64   `json:"project_id"`
	Status      string  `json:"status"`
	CommitSHA   string  `json:"commit_sha"`
	VersionInfo string  `json:"version_info"`
	ExitCode    *int    `json:"exit_code"`
	Executor    string  `json:"executor"`
	StartedAt   string  `json:"started_at"`
	EndedAt     string  `json:"ended_at"`
}

// Artifact is one declared output of a build record (CONTEXT.md: 产物),
// archived with its checksum.
type Artifact struct {
	ID            int64  `json:"id"`
	BuildRecordID int64  `json:"build_record_id"`
	Name          string `json:"name"`
	Size          int64  `json:"size"`
	Checksum      string `json:"checksum"`
}

const buildItemCols = "id, project_id, name, source_type, git_url, git_branch, git_commit, check_latest, local_path, command, artifacts_json, timeout_sec, version_cmd, prereq_json"

func scanBuildItem(sc interface{ Scan(dest ...any) error }) (BuildItem, error) {
	var it BuildItem
	var artifacts string
	err := sc.Scan(&it.ID, &it.ProjectID, &it.Name, &it.SourceType, &it.GitURL, &it.GitBranch,
		&it.GitCommit, &it.CheckLatest, &it.LocalPath, &it.Command, &artifacts, &it.TimeoutSec,
		&it.VersionCmd, &it.PrereqJSON)
	if err != nil {
		return BuildItem{}, err
	}
	it.Artifacts = decodeJSONStrings(artifacts)
	return it, nil
}

// Build-item CRUD (CreateBuildItem / GetBuildItem / ListBuildItems* /
// UpdateBuildItem / DeleteBuildItem) lives in store_config.go -- build items
// ride the unified config_objects table (kind=build_item).

// --- batches ---

// CreateBatch inserts a queued batch, returning its id.
func (s *Store) CreateBatch(ctx context.Context, b Batch) (int64, error) {
	var id int64
	err := s.enqueue(ctx, func() error {
		res, err := s.db.ExecContext(ctx,
			"INSERT INTO batches (project_id, status, items_json, created_by) VALUES (?, ?, ?, ?)",
			b.ProjectID, b.Status, b.ItemsJSON, b.CreatedBy)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		return err
	})
	return id, err
}

// GetBatch fetches one batch.
func (s *Store) GetBatch(ctx context.Context, id int64) (Batch, error) {
	var b Batch
	err := s.db.QueryRowContext(ctx,
		"SELECT id, project_id, status, items_json, created_by FROM batches WHERE id = ?", id).
		Scan(&b.ID, &b.ProjectID, &b.Status, &b.ItemsJSON, &b.CreatedBy)
	if err != nil {
		return Batch{}, fmt.Errorf("store: get batch %d: %w", id, err)
	}
	return b, nil
}

// ListBatchesByProject returns a project's batches, newest first.
func (s *Store) ListBatchesByProject(ctx context.Context, projectID int64) ([]Batch, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, project_id, status, items_json, created_by FROM batches WHERE project_id = ? ORDER BY id DESC", projectID)
	if err != nil {
		return nil, fmt.Errorf("store: list batches: %w", err)
	}
	defer rows.Close()
	var out []Batch
	for rows.Next() {
		var b Batch
		if err := rows.Scan(&b.ID, &b.ProjectID, &b.Status, &b.ItemsJSON, &b.CreatedBy); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// UpdateBatchStatus persists a batch status transition.
func (s *Store) UpdateBatchStatus(ctx context.Context, id int64, status string) error {
	return s.enqueue(ctx, func() error {
		_, err := s.db.ExecContext(ctx, "UPDATE batches SET status = ? WHERE id = ?", status, id)
		return err
	})
}

// NonterminalBatches lists batches still queued or running -- the startup
// sweep's input (M2 mirrors CEO-7A: a restart cancels in-flight builds).
func (s *Store) NonterminalBatches(ctx context.Context) ([]Batch, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, project_id, status, items_json, created_by FROM batches WHERE status IN ('queued','running')")
	if err != nil {
		return nil, fmt.Errorf("store: nonterminal batches: %w", err)
	}
	defer rows.Close()
	var out []Batch
	for rows.Next() {
		var b Batch
		if err := rows.Scan(&b.ID, &b.ProjectID, &b.Status, &b.ItemsJSON, &b.CreatedBy); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// MaxBatchID returns the highest batch id ever persisted (the kernel
// counter seed: terminal history still occupies the id space).
func (s *Store) MaxBatchID(ctx context.Context) (int64, error) {
	var m int64
	err := s.db.QueryRowContext(ctx, "SELECT COALESCE(MAX(id), 0) FROM batches").Scan(&m)
	return m, err
}

// MaxBuildRecordID returns the highest build-record id ever persisted.
func (s *Store) MaxBuildRecordID(ctx context.Context) (int64, error) {
	var m int64
	err := s.db.QueryRowContext(ctx, "SELECT COALESCE(MAX(id), 0) FROM build_records").Scan(&m)
	return m, err
}

// --- build records ---

// CreateBuildRecord inserts a pending record, returning its id.
func (s *Store) CreateBuildRecord(ctx context.Context, r BuildRecord) (int64, error) {
	var id int64
	err := s.enqueue(ctx, func() error {
		res, err := s.db.ExecContext(ctx,
			"INSERT INTO build_records (batch_id, item_id, project_id, status, executor) VALUES (?, ?, ?, ?, ?)",
			r.BatchID, r.ItemID, r.ProjectID, r.Status, r.Executor)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		return err
	})
	return id, err
}

// GetBuildRecord fetches one record.
func (s *Store) GetBuildRecord(ctx context.Context, id int64) (BuildRecord, error) {
	var r BuildRecord
	var exit sql.NullInt64
	err := s.db.QueryRowContext(ctx,
		"SELECT id, batch_id, item_id, project_id, status, commit_sha, version_info, exit_code, executor, started_at, ended_at FROM build_records WHERE id = ?", id).
		Scan(&r.ID, &r.BatchID, &r.ItemID, &r.ProjectID, &r.Status, &r.CommitSHA, &r.VersionInfo, &exit, &r.Executor, &r.StartedAt, &r.EndedAt)
	if err != nil {
		return BuildRecord{}, fmt.Errorf("store: get build record %d: %w", id, err)
	}
	if exit.Valid {
		v := int(exit.Int64)
		r.ExitCode = &v
	}
	return r, nil
}

// RecordsForBatch lists a batch's records oldest first.
func (s *Store) RecordsForBatch(ctx context.Context, batchID int64) ([]BuildRecord, error) {
	rows, err := s.db.QueryContext(ctx, recordSelect+" WHERE batch_id = ? ORDER BY id", batchID)
	if err != nil {
		return nil, fmt.Errorf("store: records for batch: %w", err)
	}
	return scanRecords(rows)
}

// RecordsForProject lists a project's records newest first (build history
// across batches, requirement 2.4). Served by idx_build_records_project.
func (s *Store) RecordsForProject(ctx context.Context, projectID int64) ([]BuildRecord, error) {
	rows, err := s.db.QueryContext(ctx, recordSelect+" WHERE project_id = ? ORDER BY id DESC", projectID)
	if err != nil {
		return nil, fmt.Errorf("store: records for project: %w", err)
	}
	return scanRecords(rows)
}

// NonterminalRecords lists records still pending or building -- the startup
// sweep cancels them along with their batches.
func (s *Store) NonterminalRecords(ctx context.Context) ([]BuildRecord, error) {
	rows, err := s.db.QueryContext(ctx, recordSelect+" WHERE status IN ('pending','building')")
	if err != nil {
		return nil, fmt.Errorf("store: nonterminal records: %w", err)
	}
	return scanRecords(rows)
}

const recordSelect = "SELECT id, batch_id, item_id, project_id, status, commit_sha, version_info, exit_code, executor, started_at, ended_at FROM build_records"

func scanRecords(rows *sql.Rows) ([]BuildRecord, error) {
	defer rows.Close()
	var out []BuildRecord
	for rows.Next() {
		var r BuildRecord
		var exit sql.NullInt64
		if err := rows.Scan(&r.ID, &r.BatchID, &r.ItemID, &r.ProjectID, &r.Status, &r.CommitSHA,
			&r.VersionInfo, &exit, &r.Executor, &r.StartedAt, &r.EndedAt); err != nil {
			return nil, err
		}
		if exit.Valid {
			v := int(exit.Int64)
			r.ExitCode = &v
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// UpdateBuildRecord persists a record's mutable fields. Timestamps ride as
// strings (RFC3339) -- the executor owns clock reads.
func (s *Store) UpdateBuildRecord(ctx context.Context, r BuildRecord) error {
	return s.enqueue(ctx, func() error {
		var exit any
		if r.ExitCode != nil {
			exit = *r.ExitCode
		}
		_, err := s.db.ExecContext(ctx,
			"UPDATE build_records SET status = ?, commit_sha = ?, version_info = ?, exit_code = ?, executor = ?, started_at = ?, ended_at = ? WHERE id = ?",
			r.Status, r.CommitSHA, r.VersionInfo, exit, r.Executor, r.StartedAt, r.EndedAt, r.ID)
		return err
	})
}

// DeleteBuildRecord removes a record row (mode=record; with mode=artifacts
// the handler deletes artifacts first).
func (s *Store) DeleteBuildRecord(ctx context.Context, id int64) error {
	return s.enqueue(ctx, func() error {
		_, err := s.db.ExecContext(ctx, "DELETE FROM build_records WHERE id = ?", id)
		return err
	})
}

// --- artifacts ---

// InsertArtifacts stores the archived artifact list for a record.
func (s *Store) InsertArtifacts(ctx context.Context, arts []Artifact) error {
	if len(arts) == 0 {
		return nil
	}
	return s.enqueue(ctx, func() error {
		for _, a := range arts {
			if _, err := s.db.ExecContext(ctx,
				"INSERT INTO artifacts (build_record_id, name, size, checksum) VALUES (?, ?, ?, ?)",
				a.BuildRecordID, a.Name, a.Size, a.Checksum); err != nil {
				return err
			}
		}
		return nil
	})
}

// ArtifactsForRecord lists a record's artifacts.
func (s *Store) ArtifactsForRecord(ctx context.Context, recordID int64) ([]Artifact, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, build_record_id, name, size, checksum FROM artifacts WHERE build_record_id = ? ORDER BY id", recordID)
	if err != nil {
		return nil, fmt.Errorf("store: artifacts for record: %w", err)
	}
	defer rows.Close()
	var out []Artifact
	for rows.Next() {
		var a Artifact
		if err := rows.Scan(&a.ID, &a.BuildRecordID, &a.Name, &a.Size, &a.Checksum); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// DeleteArtifactsForRecord removes artifact rows (mode=artifacts delete).
func (s *Store) DeleteArtifactsForRecord(ctx context.Context, recordID int64) error {
	return s.enqueue(ctx, func() error {
		_, err := s.db.ExecContext(ctx, "DELETE FROM artifacts WHERE build_record_id = ?", recordID)
		return err
	})
}

// GetArtifact fetches one artifact row (download path).
func (s *Store) GetArtifact(ctx context.Context, id int64) (Artifact, error) {
	var a Artifact
	err := s.db.QueryRowContext(ctx,
		"SELECT id, build_record_id, name, size, checksum FROM artifacts WHERE id = ?", id).
		Scan(&a.ID, &a.BuildRecordID, &a.Name, &a.Size, &a.Checksum)
	if err != nil {
		return Artifact{}, fmt.Errorf("store: get artifact %d: %w", id, err)
	}
	return a, nil
}

// --- settings ---

// GetSetting returns one setting; empty string when unset (callers apply
// defaults: max_parallel_batches=4, checksum=sha256, tmp dir falls back).
func (s *Store) GetSetting(ctx context.Context, key string) (string, error) {
	var v string
	err := s.db.QueryRowContext(ctx, "SELECT value FROM settings WHERE key = ?", key).Scan(&v)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("store: get setting %s: %w", key, err)
	}
	return v, nil
}

// SetSetting upserts one setting.
func (s *Store) SetSetting(ctx context.Context, key, value string) error {
	return s.enqueue(ctx, func() error {
		_, err := s.db.ExecContext(ctx,
			"INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value",
			key, value)
		return err
	})
}

// AllSettings returns every setting (the admin page's initial load).
func (s *Store) AllSettings(ctx context.Context) (map[string]string, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT key, value FROM settings")
	if err != nil {
		return nil, fmt.Errorf("store: all settings: %w", err)
	}
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, rows.Err()
}

// --- users (account management, requirement 1.5) ---

// UserWithID adds the id to the kernel-agnostic user view for the admin page.
type UserWithID struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// ListUsers returns all accounts for the admin page.
func (s *Store) ListUsers(ctx context.Context) ([]UserWithID, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, username, role FROM users ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("store: list users: %w", err)
	}
	defer rows.Close()
	var out []UserWithID
	for rows.Next() {
		var u UserWithID
		if err := rows.Scan(&u.ID, &u.Username, &u.Role); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// GetUserByID fetches one account.
func (s *Store) GetUserByID(ctx context.Context, id int64) (UserWithID, error) {
	var u UserWithID
	err := s.db.QueryRowContext(ctx, "SELECT id, username, role FROM users WHERE id = ?", id).
		Scan(&u.ID, &u.Username, &u.Role)
	if err != nil {
		return UserWithID{}, fmt.Errorf("store: get user %d: %w", id, err)
	}
	return u, nil
}

// AuthenticateUserByName resolves an already-authenticated username to its
// account row (the token ring stores names, the self-password path needs ids).
func (s *Store) AuthenticateUserByName(ctx context.Context, username string) (UserWithID, error) {
	var u UserWithID
	err := s.db.QueryRowContext(ctx, "SELECT id, username, role FROM users WHERE username = ?", username).
		Scan(&u.ID, &u.Username, &u.Role)
	if err != nil {
		return UserWithID{}, fmt.Errorf("store: find user %q: %w", username, err)
	}
	return u, nil
}

// UpdatePassword overwrites a user's password hash (self-change or admin
// reset). Hashing happens here so every password path goes through bcrypt.
func (s *Store) UpdatePassword(ctx context.Context, id int64, password string) error {
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}
	return s.enqueue(ctx, func() error {
		_, err := s.db.ExecContext(ctx, "UPDATE users SET password_hash = ? WHERE id = ?", hash, id)
		return err
	})
}

// ChangeOwnPassword verifies the old password before replacing it.
func (s *Store) ChangeOwnPassword(ctx context.Context, id int64, oldPassword, newPassword string) error {
	var hash string
	err := s.db.QueryRowContext(ctx, "SELECT password_hash FROM users WHERE id = ?", id).Scan(&hash)
	if err != nil {
		return fmt.Errorf("store: change password: %w", err)
	}
	if !verifyPassword(hash, oldPassword) {
		return ErrBadCredentials
	}
	return s.UpdatePassword(ctx, id, newPassword)
}

// Unified configuration objects (requirement 1.6: the four rule families
// share one storage shape). One table, one CRUD surface, one future
// copy/paste + JSON export/import path; execution semantics stay with the
// executors (design doc: unified persistence, not unified behavior).
package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// ConfigObject kinds (requirement 1.6's families; more land with M4/M5).
const (
	KindBuildItem  = "build_item"
	KindDeployRule = "deploy_rule"
	KindFlashRule  = "flash_rule"
)

// ConfigObject is one row of the unified config table. Payload is the
// kind-specific body as raw JSON -- each executor package owns its shape.
type ConfigObject struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"project_id"`
	Kind        string    `json:"kind"`
	Name        string    `json:"name"`
	PayloadJSON string    `json:"payload_json"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const configObjectCols = "id, project_id, kind, name, payload_json, created_at, updated_at"

func scanConfigObject(sc interface{ Scan(dest ...any) error }) (ConfigObject, error) {
	var o ConfigObject
	var created, updated string
	err := sc.Scan(&o.ID, &o.ProjectID, &o.Kind, &o.Name, &o.PayloadJSON, &created, &updated)
	if err != nil {
		return ConfigObject{}, err
	}
	o.CreatedAt = parseTime(created)
	o.UpdatedAt = parseTime(updated)
	return o, nil
}

// CreateConfigObject inserts one object, returning its id.
func (s *Store) CreateConfigObject(ctx context.Context, o ConfigObject) (int64, error) {
	var id int64
	err := s.enqueue(ctx, func() error {
		res, err := s.db.ExecContext(ctx,
			"INSERT INTO config_objects (project_id, kind, name, payload_json) VALUES (?, ?, ?, ?)",
			o.ProjectID, o.Kind, o.Name, o.PayloadJSON)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		return err
	})
	return id, err
}

// GetConfigObject fetches one row.
func (s *Store) GetConfigObject(ctx context.Context, id int64) (ConfigObject, error) {
	o, err := scanConfigObject(s.db.QueryRowContext(ctx,
		"SELECT "+configObjectCols+" FROM config_objects WHERE id = ?", id))
	if err != nil {
		return ConfigObject{}, fmt.Errorf("store: get config object %d: %w", id, err)
	}
	return o, nil
}

// ListConfigObjects returns one kind's objects for a project, insertion
// order (served by idx_config_objects_project_kind).
func (s *Store) ListConfigObjects(ctx context.Context, projectID int64, kind string) ([]ConfigObject, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT "+configObjectCols+" FROM config_objects WHERE project_id = ? AND kind = ? ORDER BY id",
		projectID, kind)
	if err != nil {
		return nil, fmt.Errorf("store: list config objects: %w", err)
	}
	defer rows.Close()
	var out []ConfigObject
	for rows.Next() {
		o, err := scanConfigObject(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

// ListConfigObjectsByKind returns every object of a kind across projects
// (the cross-project picker's source; requirement 2.1 cross-project prereqs).
func (s *Store) ListConfigObjectsByKind(ctx context.Context, kind string) ([]ConfigObject, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT "+configObjectCols+" FROM config_objects WHERE kind = ? ORDER BY id", kind)
	if err != nil {
		return nil, fmt.Errorf("store: list config objects by kind: %w", err)
	}
	defer rows.Close()
	var out []ConfigObject
	for rows.Next() {
		o, err := scanConfigObject(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

// ListConfigObjectsByIDs fetches a specific set of one kind (batch planning
// input). Missing ids are silently absent -- callers check the count.
func (s *Store) ListConfigObjectsByIDs(ctx context.Context, kind string, ids []int64) ([]ConfigObject, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := "SELECT " + configObjectCols + " FROM config_objects WHERE kind = ? AND id IN ("
	args := make([]any, 0, len(ids)+1)
	args = append(args, kind)
	for i, id := range ids {
		if i > 0 {
			q += ","
		}
		q += "?"
		args = append(args, id)
	}
	rows, err := s.db.QueryContext(ctx, q+")", args...)
	if err != nil {
		return nil, fmt.Errorf("store: list config objects by ids: %w", err)
	}
	defer rows.Close()
	var out []ConfigObject
	for rows.Next() {
		o, err := scanConfigObject(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

// UpdateConfigObject overwrites name and payload (project and kind are
// immutable: moving an object between projects/kinds is a delete+create).
func (s *Store) UpdateConfigObject(ctx context.Context, o ConfigObject) error {
	return s.enqueue(ctx, func() error {
		_, err := s.db.ExecContext(ctx,
			"UPDATE config_objects SET name = ?, payload_json = ?, updated_at = strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id = ?",
			o.Name, o.PayloadJSON, o.ID)
		return err
	})
}

// DeleteConfigObject removes one row.
func (s *Store) DeleteConfigObject(ctx context.Context, id int64) error {
	return s.enqueue(ctx, func() error {
		_, err := s.db.ExecContext(ctx, "DELETE FROM config_objects WHERE id = ?", id)
		return err
	})
}

// CopyConfigObject duplicates one object into another project (requirement
// 1.6 跨工程复制粘贴): same kind, name and payload, fresh id. The cross-
// project reference inside build items' prereq JSON is copied verbatim --
// references are logical ids and stay valid across projects (batch planning
// reads them cross-project anyway).
func (s *Store) CopyConfigObject(ctx context.Context, id, targetProjectID int64) (int64, error) {
	o, err := s.GetConfigObject(ctx, id)
	if err != nil {
		return 0, err
	}
	return s.CreateConfigObject(ctx, ConfigObject{
		ProjectID:   targetProjectID,
		Kind:        o.Kind,
		Name:        o.Name,
		PayloadJSON: o.PayloadJSON,
	})
}

// --- payload shapes for the migrated kinds ---

// BuildItemPayload is the config_objects body for kind=build_item. Field
// meanings match the old build_items columns (store.BuildItem); it is the
// persistence shape -- the API surface keeps returning store.BuildItem.
type BuildItemPayload struct {
	SourceType  string   `json:"source_type"`
	GitURL      string   `json:"git_url,omitempty"`
	GitBranch   string   `json:"git_branch,omitempty"`
	GitCommit   string   `json:"git_commit,omitempty"`
	CheckLatest bool     `json:"check_latest,omitempty"`
	LocalPath   string   `json:"local_path,omitempty"`
	Command     string   `json:"command"`
	Artifacts   []string `json:"artifacts"`
	TimeoutSec  int      `json:"timeout_sec,omitempty"`
	VersionCmd  string   `json:"version_cmd,omitempty"`
	PrereqJSON  string   `json:"prereq_json"`
}

// FlashRulePayload is the body for kind=flash_rule (the M1 expect rules,
// migrated from expect_rules).
type FlashRulePayload struct {
	StepsJSON string `json:"steps_json"`
}

// --- translation helpers (the compat layer's engines) ---

func buildItemToPayload(it BuildItem) string {
	raw, err := json.Marshal(BuildItemPayload{
		SourceType:  it.SourceType,
		GitURL:      it.GitURL,
		GitBranch:   it.GitBranch,
		GitCommit:   it.GitCommit,
		CheckLatest: it.CheckLatest,
		LocalPath:   it.LocalPath,
		Command:     it.Command,
		Artifacts:   it.Artifacts,
		TimeoutSec:  it.TimeoutSec,
		VersionCmd:  it.VersionCmd,
		PrereqJSON:  it.PrereqJSON,
	})
	if err != nil {
		return "{}" // strings and bools always marshal; unreachable
	}
	return string(raw)
}

func buildItemFromPayload(o ConfigObject) (BuildItem, error) {
	var p BuildItemPayload
	if err := json.Unmarshal([]byte(o.PayloadJSON), &p); err != nil {
		return BuildItem{}, fmt.Errorf("store: build item payload %d: %w", o.ID, err)
	}
	return BuildItem{
		ID:          o.ID,
		ProjectID:   o.ProjectID,
		Name:        o.Name,
		SourceType:  p.SourceType,
		GitURL:      p.GitURL,
		GitBranch:   p.GitBranch,
		GitCommit:   p.GitCommit,
		CheckLatest: p.CheckLatest,
		LocalPath:   p.LocalPath,
		Command:     p.Command,
		Artifacts:   p.Artifacts,
		TimeoutSec:  p.TimeoutSec,
		VersionCmd:  p.VersionCmd,
		PrereqJSON:  p.PrereqJSON,
	}, nil
}

// --- compat: build items ride config_objects now (same method surface) ---

// CreateBuildItem inserts an item, returning its id.
func (s *Store) CreateBuildItem(ctx context.Context, it BuildItem) (int64, error) {
	return s.CreateConfigObject(ctx, ConfigObject{
		ProjectID:   it.ProjectID,
		Kind:        KindBuildItem,
		Name:        it.Name,
		PayloadJSON: buildItemToPayload(it),
	})
}

// GetBuildItem fetches one item.
func (s *Store) GetBuildItem(ctx context.Context, id int64) (BuildItem, error) {
	o, err := s.GetConfigObject(ctx, id)
	if err != nil {
		return BuildItem{}, err
	}
	if o.Kind != KindBuildItem {
		return BuildItem{}, fmt.Errorf("store: object %d is %s, not a build item", id, o.Kind)
	}
	return buildItemFromPayload(o)
}

// ListBuildItemsByProject returns a project's items in insertion order.
func (s *Store) ListBuildItemsByProject(ctx context.Context, projectID int64) ([]BuildItem, error) {
	objects, err := s.ListConfigObjects(ctx, projectID, KindBuildItem)
	if err != nil {
		return nil, err
	}
	return buildItemsFromObjects(objects)
}

// ListAllBuildItems returns every item across projects (batch-planning
// universe: prerequisites may reference cross-project ids, requirement 2.1).
func (s *Store) ListAllBuildItems(ctx context.Context) ([]BuildItem, error) {
	objects, err := s.ListConfigObjectsByKind(ctx, KindBuildItem)
	if err != nil {
		return nil, err
	}
	return buildItemsFromObjects(objects)
}

// ListBuildItemsByIDs fetches a specific set (batch planning input).
func (s *Store) ListBuildItemsByIDs(ctx context.Context, ids []int64) ([]BuildItem, error) {
	objects, err := s.ListConfigObjectsByIDs(ctx, KindBuildItem, ids)
	if err != nil {
		return nil, err
	}
	return buildItemsFromObjects(objects)
}

// UpdateBuildItem overwrites the editable fields.
func (s *Store) UpdateBuildItem(ctx context.Context, it BuildItem) error {
	return s.UpdateConfigObject(ctx, ConfigObject{
		ID:          it.ID,
		Name:        it.Name,
		PayloadJSON: buildItemToPayload(it),
	})
}

// DeleteBuildItem removes an item. Existing build records keep item_id for
// history display (now a logical reference to config_objects).
func (s *Store) DeleteBuildItem(ctx context.Context, id int64) error {
	return s.DeleteConfigObject(ctx, id)
}

func buildItemsFromObjects(objects []ConfigObject) ([]BuildItem, error) {
	out := make([]BuildItem, 0, len(objects))
	for _, o := range objects {
		it, err := buildItemFromPayload(o)
		if err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, nil
}

// --- compat: expect rules ride config_objects now (kind=flash_rule) ---

// CreateExpectRule stores a rule; steps_json carries the serialized steps.
func (s *Store) CreateExpectRule(ctx context.Context, projectID int64, name, stepsJSON string) (int64, error) {
	var id int64
	err := s.enqueue(ctx, func() error {
		res, err := s.db.ExecContext(ctx,
			"INSERT INTO config_objects (project_id, kind, name, payload_json) VALUES (?, ?, ?, ?)",
			projectID, KindFlashRule, name, flashRuleToPayload(stepsJSON))
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
	o, err := s.GetConfigObject(ctx, id)
	if err != nil {
		return ExpectRule{}, err
	}
	return flashRuleFromLegacy(o)
}

// UpdateExpectRule overwrites name and steps.
func (s *Store) UpdateExpectRule(ctx context.Context, r ExpectRule) error {
	return s.UpdateConfigObject(ctx, ConfigObject{
		ID:          r.ID,
		Name:        r.Name,
		PayloadJSON: flashRuleToPayload(r.StepsJSON),
	})
}

// ListExpectRules returns all rules in insertion order.
func (s *Store) ListExpectRules(ctx context.Context) ([]ExpectRule, error) {
	objects, err := s.ListConfigObjectsByKind(ctx, KindFlashRule)
	if err != nil {
		return nil, err
	}
	out := make([]ExpectRule, 0, len(objects))
	for _, o := range objects {
		r, err := flashRuleFromLegacy(o)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

// DeleteExpectRule removes a rule (the editor's delete button).
func (s *Store) DeleteExpectRule(ctx context.Context, id int64) error {
	return s.DeleteConfigObject(ctx, id)
}

func flashRuleToPayload(stepsJSON string) string {
	raw, err := json.Marshal(FlashRulePayload{StepsJSON: stepsJSON})
	if err != nil {
		return "{}"
	}
	return string(raw)
}

// flashRuleFromLegacy translates a config object into the M1 ExpectRule
// view (used by the M1 endpoints and the demo seed).
func flashRuleFromLegacy(o ConfigObject) (ExpectRule, error) {
	var p FlashRulePayload
	if err := json.Unmarshal([]byte(o.PayloadJSON), &p); err != nil {
		return ExpectRule{}, fmt.Errorf("store: flash rule payload %d: %w", o.ID, err)
	}
	return ExpectRule{
		ID:        o.ID,
		Name:      o.Name,
		StepsJSON: p.StepsJSON,
		UpdatedAt: o.UpdatedAt,
	}, nil
}

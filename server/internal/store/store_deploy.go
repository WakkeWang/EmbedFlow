// Deploy-module persistence (requirement 3): deploy rules (three forms),
// deploy records, and command-execution records (design decision 1A: SSH
// runs are short transactions outside the session state machine).
package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// DeployRule is one deployment recipe (CONTEXT.md: 部署规则). Mode picks the
// form; the mode-specific fields ride DeployPayload (one JSON body in
// config_objects, kind=deploy_rule).
type DeployRule struct {
	ID        int64     `json:"id"`
	ProjectID int64     `json:"project_id"`
	Name      string    `json:"name"`
	Mode      string    `json:"mode"` // manual | ssh | flash
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DeployPayload is the config-object body for kind=deploy_rule. Exactly one
// mode's fields apply; JSON keeps the others empty.
type DeployPayload struct {
	Mode string `json:"mode"`
	// manual: the step description shown to the operator (requirement
	// 3.1.1; artifact references are by name, the UI links downloads).
	StepsMD string `json:"steps_md,omitempty"`
	// ssh: one command template per line (requirement 3.1.2), variables
	// ${artifact.*} / ${device.*} / ${params.*} rendered at run time.
	SSHDeviceID  int64             `json:"ssh_device_id,omitempty"`
	SSHCommands  []string          `json:"ssh_commands,omitempty"`
	SSHParams    []DeployParamDef  `json:"ssh_params,omitempty"`
	SSHTimeoutSec int              `json:"ssh_timeout_sec,omitempty"`
	// flash: an expect step sequence (requirement 3.1.3) executed through
	// the M1 engine on the device's serial transport.
	FlashDeviceID int64           `json:"flash_device_id,omitempty"`
	FlashStepsJSON string         `json:"flash_steps_json,omitempty"`
	FlashTimeoutSec int            `json:"flash_timeout_sec,omitempty"`
}

// DeployParamDef declares one operator-supplied parameter (name + default)
// the SSH template references as ${params.<name>}.
type DeployParamDef struct {
	Name    string `json:"name"`
	Default string `json:"default,omitempty"`
}

// DeployRecord is one deployment execution (CONTEXT.md: 部署记录).
type DeployRecord struct {
	ID            int64  `json:"id"`
	ProjectID     int64  `json:"project_id"`
	RuleID        int64  `json:"rule_id"`
	DeviceID      int64  `json:"device_id"`
	BuildRecordID int64  `json:"build_record_id"`
	Executor      string `json:"executor"`
	Status        string `json:"status"` // running | succeeded | failed | canceled
	Detail        string `json:"detail"`
	// ExecRecordID links the SSH command-execution record (decision 1A);
	// SessionID links the flash task session.
	ExecRecordID int64  `json:"exec_record_id,omitempty"`
	SessionID    int64  `json:"session_id,omitempty"`
	StartedAt    string `json:"started_at"`
	EndedAt      string `json:"ended_at"`
}

// Deploy statuses (keep the string set closed; CHECK in the DDL mirrors it).
const (
	DeployRunning   = "running"
	DeploySucceeded = "succeeded"
	DeployFailed    = "failed"
	DeployCanceled  = "canceled"
)

// ExecRecord is one SSH command execution (decision 1A: short transaction,
// own log file, referenced by deploy records and later test runs).
type ExecRecord struct {
	ID        int64     `json:"id"`
	DeviceID  int64     `json:"device_id"`
	Kind      string    `json:"kind"` // deploy | test
	Command   string    `json:"command"`
	Status    string    `json:"status"` // running | succeeded | failed
	ExitCode  *int      `json:"exit_code"`
	Detail    string    `json:"detail"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
}

const deployRecordCols = "id, project_id, rule_id, device_id, build_record_id, executor, status, detail, exec_record_id, session_id, started_at, ended_at"

func scanDeployRecord(sc interface{ Scan(dest ...any) error }) (DeployRecord, error) {
	var r DeployRecord
	err := sc.Scan(&r.ID, &r.ProjectID, &r.RuleID, &r.DeviceID, &r.BuildRecordID, &r.Executor,
		&r.Status, &r.Detail, &r.ExecRecordID, &r.SessionID, &r.StartedAt, &r.EndedAt)
	return r, err
}

// CreateDeployRule stores a rule into config_objects, returning its id.
func (s *Store) CreateDeployRule(ctx context.Context, projectID int64, name string, p DeployPayload) (int64, error) {
	raw, err := json.Marshal(p)
	if err != nil {
		return 0, fmt.Errorf("store: deploy payload: %w", err)
	}
	return s.CreateConfigObject(ctx, ConfigObject{
		ProjectID: projectID, Kind: KindDeployRule, Name: name, PayloadJSON: string(raw),
	})
}

// GetDeployRule fetches one rule.
func (s *Store) GetDeployRule(ctx context.Context, id int64) (DeployRule, DeployPayload, error) {
	o, err := s.GetConfigObject(ctx, id)
	if err != nil {
		return DeployRule{}, DeployPayload{}, err
	}
	if o.Kind != KindDeployRule {
		return DeployRule{}, DeployPayload{}, fmt.Errorf("store: object %d is %s, not a deploy rule", id, o.Kind)
	}
	var p DeployPayload
	if err := json.Unmarshal([]byte(o.PayloadJSON), &p); err != nil {
		return DeployRule{}, DeployPayload{}, fmt.Errorf("store: deploy rule payload %d: %w", id, err)
	}
	return DeployRule{ID: o.ID, ProjectID: o.ProjectID, Name: o.Name, Mode: p.Mode, CreatedAt: o.CreatedAt, UpdatedAt: o.UpdatedAt}, p, nil
}

// ListDeployRules returns a project's rules, oldest first.
func (s *Store) ListDeployRules(ctx context.Context, projectID int64) ([]DeployRule, []DeployPayload, error) {
	objects, err := s.ListConfigObjects(ctx, projectID, KindDeployRule)
	if err != nil {
		return nil, nil, err
	}
	rules := make([]DeployRule, 0, len(objects))
	payloads := make([]DeployPayload, 0, len(objects))
	for _, o := range objects {
		var p DeployPayload
		if err := json.Unmarshal([]byte(o.PayloadJSON), &p); err != nil {
			return nil, nil, fmt.Errorf("store: deploy rule payload %d: %w", o.ID, err)
		}
		rules = append(rules, DeployRule{ID: o.ID, ProjectID: o.ProjectID, Name: o.Name, Mode: p.Mode, CreatedAt: o.CreatedAt, UpdatedAt: o.UpdatedAt})
		payloads = append(payloads, p)
	}
	return rules, payloads, nil
}

// UpdateDeployRule overwrites name and payload.
func (s *Store) UpdateDeployRule(ctx context.Context, r DeployRule, p DeployPayload) error {
	raw, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("store: deploy payload: %w", err)
	}
	return s.UpdateConfigObject(ctx, ConfigObject{ID: r.ID, Name: r.Name, PayloadJSON: string(raw)})
}

// DeleteDeployRule removes one rule.
func (s *Store) DeleteDeployRule(ctx context.Context, id int64) error {
	return s.DeleteConfigObject(ctx, id)
}

// CreateDeployRecord inserts a running deploy record, returning its id.
func (s *Store) CreateDeployRecord(ctx context.Context, r DeployRecord) (int64, error) {
	var id int64
	err := s.enqueue(ctx, func() error {
		res, err := s.db.ExecContext(ctx,
			"INSERT INTO deploy_records (project_id, rule_id, device_id, build_record_id, executor, status, detail, started_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			r.ProjectID, r.RuleID, r.DeviceID, r.BuildRecordID, r.Executor, r.Status, r.Detail, r.StartedAt)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		return err
	})
	return id, err
}

// GetDeployRecord fetches one record.
func (s *Store) GetDeployRecord(ctx context.Context, id int64) (DeployRecord, error) {
	r, err := scanDeployRecord(s.db.QueryRowContext(ctx,
		"SELECT "+deployRecordCols+" FROM deploy_records WHERE id = ?", id))
	if err != nil {
		return DeployRecord{}, fmt.Errorf("store: get deploy record %d: %w", id, err)
	}
	return r, nil
}

// ListDeployRecords lists a project's deploy records, newest first.
func (s *Store) ListDeployRecords(ctx context.Context, projectID int64) ([]DeployRecord, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT "+deployRecordCols+" FROM deploy_records WHERE project_id = ? ORDER BY id DESC", projectID)
	if err != nil {
		return nil, fmt.Errorf("store: list deploy records: %w", err)
	}
	defer rows.Close()
	var out []DeployRecord
	for rows.Next() {
		r, err := scanDeployRecord(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// NonterminalDeployRecords lists records still running (startup sweep input).
func (s *Store) NonterminalDeployRecords(ctx context.Context) ([]DeployRecord, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT "+deployRecordCols+" FROM deploy_records WHERE status = ?", DeployRunning)
	if err != nil {
		return nil, fmt.Errorf("store: nonterminal deploy records: %w", err)
	}
	defer rows.Close()
	var out []DeployRecord
	for rows.Next() {
		r, err := scanDeployRecord(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// UpdateDeployRecord persists mutable fields (status, detail, links, end).
func (s *Store) UpdateDeployRecord(ctx context.Context, r DeployRecord) error {
	return s.enqueue(ctx, func() error {
		_, err := s.db.ExecContext(ctx,
			"UPDATE deploy_records SET status = ?, detail = ?, exec_record_id = ?, session_id = ?, ended_at = ? WHERE id = ?",
			r.Status, r.Detail, r.ExecRecordID, r.SessionID, r.EndedAt, r.ID)
		return err
	})
}

// CreateExecRecord inserts a running command-execution record, returning its id.
func (s *Store) CreateExecRecord(ctx context.Context, r ExecRecord) (int64, error) {
	var id int64
	err := s.enqueue(ctx, func() error {
		res, err := s.db.ExecContext(ctx,
			"INSERT INTO exec_records (device_id, kind, command, status, started_at) VALUES (?, ?, ?, ?, ?)",
			r.DeviceID, r.Kind, r.Command, r.Status, r.StartedAt.Format(time.RFC3339Nano))
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		return err
	})
	return id, err
}

// UpdateExecRecord persists the execution outcome.
func (s *Store) UpdateExecRecord(ctx context.Context, r ExecRecord) error {
	return s.enqueue(ctx, func() error {
		var exit any
		if r.ExitCode != nil {
			exit = *r.ExitCode
		}
		_, err := s.db.ExecContext(ctx,
			"UPDATE exec_records SET status = ?, exit_code = ?, detail = ?, ended_at = ? WHERE id = ?",
			r.Status, exit, r.Detail, r.EndedAt.Format(time.RFC3339Nano), r.ID)
		return err
	})
}

// GetExecRecord fetches one execution record.
func (s *Store) GetExecRecord(ctx context.Context, id int64) (ExecRecord, error) {
	var r ExecRecord
	var exit sql.NullInt64
	var started, ended string
	err := s.db.QueryRowContext(ctx,
		"SELECT id, device_id, kind, command, status, exit_code, detail, started_at, ended_at FROM exec_records WHERE id = ?", id).
		Scan(&r.ID, &r.DeviceID, &r.Kind, &r.Command, &r.Status, &exit, &r.Detail, &started, &ended)
	if err != nil {
		return ExecRecord{}, fmt.Errorf("store: get exec record %d: %w", id, err)
	}
	if exit.Valid {
		v := int(exit.Int64)
		r.ExitCode = &v
	}
	r.StartedAt = parseTime(started)
	r.EndedAt = parseTime(ended)
	return r, nil
}

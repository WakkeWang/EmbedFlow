package store

import (
	"context"
	"testing"
	"time"
)

func timeNow() time.Time { return time.Now() }

func TestDeployRules_CRUD(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	pid, _ := s.CreateProject(ctx, "P-dep", "")

	id, err := s.CreateDeployRule(ctx, pid, "flash-over-tftp", DeployPayload{
		Mode: "ssh", SSHDeviceID: 7, SSHCommands: []string{"cd ${params.dir} && ls"},
		SSHParams: []DeployParamDef{{Name: "dir", Default: "/tmp"}}, SSHTimeoutSec: 120,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	rule, payload, err := s.GetDeployRule(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if rule.Name != "flash-over-tftp" || rule.Mode != "ssh" || rule.ProjectID != pid {
		t.Fatalf("rule = %+v", rule)
	}
	if len(payload.SSHCommands) != 1 || payload.SSHCommands[0] != "cd ${params.dir} && ls" || payload.SSHDeviceID != 7 {
		t.Fatalf("payload = %+v", payload)
	}

	payload.SSHTimeoutSec = 300
	if err := s.UpdateDeployRule(ctx, rule, payload); err != nil {
		t.Fatalf("update: %v", err)
	}
	_, after, _ := s.GetDeployRule(ctx, id)
	if after.SSHTimeoutSec != 300 {
		t.Fatalf("after update = %+v", after)
	}

	rules, payloads, err := s.ListDeployRules(ctx, pid)
	if err != nil || len(rules) != 1 || len(payloads) != 1 {
		t.Fatalf("list = %v err %v", rules, err)
	}
	if err := s.DeleteDeployRule(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, _, err := s.GetDeployRule(ctx, id); err == nil {
		t.Fatal("deleted rule still readable")
	}
}

func TestDeployRecords_RoundTrip(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	pid, _ := s.CreateProject(ctx, "P-dep2", "")

	id, err := s.CreateDeployRecord(ctx, DeployRecord{
		ProjectID: pid, RuleID: 3, DeviceID: 4, BuildRecordID: 5,
		Executor: "alice", Status: DeployRunning, StartedAt: "2026-09-19T10:00:00+08:00",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	rec, err := s.GetDeployRecord(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if rec.Status != DeployRunning || rec.Executor != "alice" || rec.RuleID != 3 || rec.DeviceID != 4 {
		t.Fatalf("record = %+v", rec)
	}

	rec.Status = DeploySucceeded
	rec.Detail = "2 commands completed"
	rec.ExecRecordID = 9
	rec.SessionID = 11
	rec.EndedAt = "2026-09-19T10:01:00+08:00"
	if err := s.UpdateDeployRecord(ctx, rec); err != nil {
		t.Fatalf("update: %v", err)
	}
	after, _ := s.GetDeployRecord(ctx, id)
	if after.Status != DeploySucceeded || after.ExecRecordID != 9 || after.SessionID != 11 {
		t.Fatalf("after = %+v", after)
	}

	list, err := s.ListDeployRecords(ctx, pid)
	if err != nil || len(list) != 1 {
		t.Fatalf("list = %v err %v", list, err)
	}

	// Nonterminal sweep input: only the running one surfaces.
	if _, err := s.CreateDeployRecord(ctx, DeployRecord{ProjectID: pid, RuleID: 1, DeviceID: 1, Status: DeployRunning, StartedAt: "t"}); err != nil {
		t.Fatalf("create second: %v", err)
	}
	non, err := s.NonterminalDeployRecords(ctx)
	if err != nil || len(non) != 1 {
		t.Fatalf("nonterminal = %v err %v", non, err)
	}
}

func TestExecRecords_RoundTrip(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	id, err := s.CreateExecRecord(ctx, ExecRecord{
		DeviceID: 2, Kind: "deploy", Command: "reboot", Status: "running",
		StartedAt: timeNow(),
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	rec, err := s.GetExecRecord(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if rec.Command != "reboot" || rec.Status != "running" {
		t.Fatalf("record = %+v", rec)
	}

	code := 0
	rec.Status = "succeeded"
	rec.ExitCode = &code
	rec.Detail = "ok"
	rec.EndedAt = timeNow()
	if err := s.UpdateExecRecord(ctx, rec); err != nil {
		t.Fatalf("update: %v", err)
	}
	after, _ := s.GetExecRecord(ctx, id)
	if after.Status != "succeeded" || after.ExitCode == nil || *after.ExitCode != 0 {
		t.Fatalf("after = %+v", after)
	}
}

func TestDeviceSSH_FieldsRoundTrip(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	id, err := s.CreateDeviceFull(ctx, Device{
		Name: "board-01", Project: "P0133", SSHHost: "10.0.0.8", SSHPort: 2222,
		SSHUser: "root", Note: "lab bench",
	}, "enc:v1:sealed")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	d, err := s.GetDeviceFull(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if d.SSHHost != "10.0.0.8" || d.SSHPort != 2222 || d.SSHUser != "root" || d.Note != "lab bench" {
		t.Fatalf("device = %+v", d)
	}
	if !d.SshSet || d.SshPassEnc != "enc:v1:sealed" {
		t.Fatalf("password fields = %+v", d)
	}

	// Update without a password keeps the stored credential.
	if err := s.UpdateDeviceSSH(ctx, Device{ID: id, Name: "board-01", Project: "P0133", SSHHost: "10.0.0.9", SSHPort: 22, SSHUser: "root"}, ""); err != nil {
		t.Fatalf("update: %v", err)
	}
	d2, _ := s.GetDeviceFull(ctx, id)
	if d2.SSHHost != "10.0.0.9" || d2.SshPassEnc != "enc:v1:sealed" || !d2.SshSet {
		t.Fatalf("after update = %+v", d2)
	}

	// Update with a password replaces it.
	if err := s.UpdateDeviceSSH(ctx, Device{ID: id, Name: "board-01", SSHHost: "10.0.0.9"}, "enc:v1:new"); err != nil {
		t.Fatalf("update pass: %v", err)
	}
	d3, _ := s.GetDeviceFull(ctx, id)
	if d3.SshPassEnc != "enc:v1:new" {
		t.Fatalf("password not replaced: %+v", d3)
	}

	// The full-device list carries SSH fields too.
	devs, err := s.ListDevices(ctx)
	if err != nil || len(devs) != 1 || devs[0].SSHHost != "10.0.0.9" {
		t.Fatalf("list = %v err %v", devs, err)
	}
}

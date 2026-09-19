package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/WakkeWang/EmbedFlow/pkg/protocol"
	"github.com/WakkeWang/EmbedFlow/server/internal/deployer"
	"github.com/WakkeWang/EmbedFlow/server/internal/expect"
	"github.com/WakkeWang/EmbedFlow/server/internal/session"
	"github.com/WakkeWang/EmbedFlow/server/internal/store"
	"github.com/WakkeWang/EmbedFlow/server/internal/transport"
)

// deployOrchestrator wires deploy rules, the SSH/flash executors and
// persistence together (the M3 analog of buildOrchestrator). Subscriptions
// mirror the build module: browser connections register per deploy record
// and receive DeployEvent frames.
type deployOrchestrator struct {
	store *store.Store
	hub   *coordinator // for WS fan-out + data dir + kernel + secrets

	mu      sync.Mutex
	subs    map[int64]map[*transport.ClientConn]struct{} // deployRecordID -> subscribers
	sshRuns map[int64]*deployer.SSHRunner                // deployRecordID -> running SSH batch
}

func newDeployOrchestrator(st *store.Store, hub *coordinator) *deployOrchestrator {
	return &deployOrchestrator{
		store:   st,
		hub:     hub,
		subs:    map[int64]map[*transport.ClientConn]struct{}{},
		sshRuns: map[int64]*deployer.SSHRunner{},
	}
}

// sweepStartup cancels deploy records left running by a restart (M3 mirrors
// the build sweep: a restart's half-run deploys are not resumable).
func (d *deployOrchestrator) sweepStartup() error {
	ctx := context.Background()
	rows, err := d.store.NonterminalDeployRecords(ctx)
	if err != nil {
		return err
	}
	for _, r := range rows {
		r.Status = store.DeployCanceled
		r.EndedAt = time.Now().Format(time.RFC3339)
		r.Detail = "server restart"
		if err := d.store.UpdateDeployRecord(ctx, r); err != nil {
			return err
		}
	}
	if len(rows) > 0 {
		slog.Warn("deploy startup sweep: canceled in-flight deployments", "count", len(rows))
	}
	return nil
}

// TriggerRequest is the deploy trigger payload (requirement 5.2's shape,
// applied to deployments): a rule, a build record, operator params.
type DeployTrigger struct {
	RuleID        int64
	BuildRecordID int64
	Params        map[string]string
	User          string
}

// Trigger creates the deploy record and starts the execution. Returns the
// record id. Manual deploys complete immediately (their "execution" is the
// UI showing steps + downloads).
func (d *deployOrchestrator) Trigger(ctx context.Context, t DeployTrigger) (int64, error) {
	rule, payload, err := d.store.GetDeployRule(ctx, t.RuleID)
	if err != nil {
		return 0, fmt.Errorf("no such rule: %w", err)
	}

	// The device comes from the rule (ssh/flash); manual has none.
	deviceID := payload.SSHDeviceID
	if payload.Mode == "flash" {
		deviceID = payload.FlashDeviceID
	}
	var device store.Device
	if deviceID != 0 {
		device, err = d.store.GetDeviceFull(ctx, deviceID)
		if err != nil {
			return 0, fmt.Errorf("no such device: %w", err)
		}
	}

	// Artifacts of the chosen build record (requirement 5.2's selection
	// model, applied to deploys: deploy what a build produced).
	var artifacts []deployer.ArtifactView
	if t.BuildRecordID != 0 {
		arts, err := d.store.ArtifactsForRecord(ctx, t.BuildRecordID)
		if err != nil {
			return 0, err
		}
		for _, a := range arts {
			artifacts = append(artifacts, deployer.ArtifactView{ID: a.ID, Name: a.Name, Size: a.Size, Checksum: a.Checksum})
		}
	}

	now := time.Now().Format(time.RFC3339)
	recID, err := d.store.CreateDeployRecord(ctx, store.DeployRecord{
		ProjectID:     rule.ProjectID,
		RuleID:        rule.ID,
		DeviceID:      deviceID,
		BuildRecordID: t.BuildRecordID,
		Executor:      t.User,
		Status:        store.DeployRunning,
		StartedAt:     now,
	})
	if err != nil {
		return 0, err
	}

	switch payload.Mode {
	case "manual":
		// No execution: the record completes at once; the UI shows the
		// markdown steps + the artifact download list (requirement 3.1.1).
		rec, _ := d.store.GetDeployRecord(ctx, recID)
		rec.Status = store.DeploySucceeded
		rec.Detail = "manual deployment: follow the steps, artifacts listed"
		rec.EndedAt = time.Now().Format(time.RFC3339)
		if err := d.store.UpdateDeployRecord(ctx, rec); err != nil {
			return 0, err
		}
		d.emit(recID, "succeeded", rec.Detail)
		return recID, nil

	case "ssh":
		return recID, d.runSSH(recID, rule, payload, device, artifacts, t)

	case "flash":
		return recID, d.runFlash(recID, rule, payload, device, t)

	default:
		return recID, fmt.Errorf("unknown deploy mode %q", payload.Mode)
	}
}

// runSSH executes the SSH form: render every template line, then run the
// batch on one connection (decision 1A: an exec record rides along).
func (d *deployOrchestrator) runSSH(recID int64, rule store.DeployRule, payload store.DeployPayload, device store.Device, artifacts []deployer.ArtifactView, t DeployTrigger) error {
	ctx := context.Background()

	if device.SSHHost == "" {
		d.finishDeploy(recID, store.DeployFailed, "device has no SSH configuration")
		return nil
	}
	password := d.hub.decryptPassword(device.SshPassEnc)
	if password == "" {
		d.finishDeploy(recID, store.DeployFailed, "device credential unreadable -- re-enter the SSH password")
		return nil
	}

	// Command-execution record (decision 1A: short transaction outside the
	// session state machine).
	commands, missing := d.renderCommands(payload, artifacts, device, t.Params)
	for _, m := range missing {
		d.log(recID, deployer.PhaseRender, "warning: unresolved variable ${"+m+"} renders empty")
	}
	execID, err := d.store.CreateExecRecord(ctx, store.ExecRecord{
		DeviceID:  device.ID,
		Kind:      "deploy",
		Command:   strings.Join(commands, " && "),
		Status:    "running",
		StartedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	rec, _ := d.store.GetDeployRecord(ctx, recID)
	rec.ExecRecordID = execID
	_ = d.store.UpdateDeployRecord(ctx, rec)

	runner := &deployer.SSHRunner{
		Host:     device.SSHHost,
		Port:     orDefault(device.SSHPort, 22),
		User:     device.SSHUser,
		Password: password,
		Timeout:  time.Duration(payload.SSHTimeoutSec) * time.Second,
		Notify:   d.notifyWithLog(execID),
		Record:   recID,
	}
	d.mu.Lock()
	d.sshRuns[recID] = runner
	d.mu.Unlock()
	d.emit(recID, "running", "")

	go func() {
		defer func() {
			d.mu.Lock()
			delete(d.sshRuns, recID)
			d.mu.Unlock()
		}()
		// The HTTP request's ctx dies with the handler; the deploy outlives
		// it, so background + the runner's own timeout govern this goroutine.
		bg := context.Background()
		res := runner.Run(bg, commands)

		exec, _ := d.store.GetExecRecord(bg, execID)
		exec.EndedAt = time.Now()
		if res.Canceled {
			exec.Status = "failed"
			exec.Detail = "canceled"
			_ = d.store.UpdateExecRecord(bg, exec)
			d.finishDeploy(recID, store.DeployCanceled, "canceled")
			return
		}
		exec.ExitCode = &res.ExitCode
		if res.OK {
			exec.Status = "succeeded"
		} else {
			exec.Status = "failed"
		}
		exec.Detail = res.Detail
		_ = d.store.UpdateExecRecord(bg, exec)

		if res.OK {
			d.finishDeploy(recID, store.DeploySucceeded, res.Detail)
		} else {
			d.finishDeploy(recID, store.DeployFailed, res.Detail)
		}
	}()
	return nil
}

// renderCommands expands the three variable families (requirement 3.1.2).
func (d *deployOrchestrator) renderCommands(payload store.DeployPayload, artifacts []deployer.ArtifactView, device store.Device, params map[string]string) ([]string, []string) {
	artifactVars := deployer.ArtifactVariables(artifacts)
	deviceVars := map[string]string{
		"ip":   device.SSHHost,
		"port": strconv.Itoa(device.SSHPort),
		"name": device.Name,
	}
	if params == nil {
		params = map[string]string{}
	}
	// Declared defaults fill gaps in the operator's params.
	for _, def := range payload.SSHParams {
		if _, ok := params[def.Name]; !ok {
			params[def.Name] = def.Default
		}
	}
	var missing []string
	out := make([]string, 0, len(payload.SSHCommands))
	for _, tpl := range payload.SSHCommands {
		rendered, m := deployer.Render(tpl, deployer.Variables{
			Artifact: artifactVars, Device: deviceVars, Params: params,
		})
		missing = append(missing, m...)
		out = append(out, rendered)
	}
	return out, missing
}

// runFlash executes the flash form through the M1 expect engine on a task
// session (requirement 3.1.3). The heavy lifting reuses the coordinator's
// expect wiring; this method links the session to the deploy record.
func (d *deployOrchestrator) runFlash(recID int64, rule store.DeployRule, payload store.DeployPayload, device store.Device, t DeployTrigger) error {
	if payload.FlashDeviceID == 0 {
		d.finishDeploy(recID, store.DeployFailed, "flash rule has no device")
		return nil
	}
	steps, err := expect.DecodeSteps([]byte(payload.FlashStepsJSON))
	if err != nil {
		d.finishDeploy(recID, store.DeployFailed, "bad flash steps: "+err.Error())
		return nil
	}

	// Open the task session through the kernel (busy devices reject with
	// the occupier info -- requirement 3.3).
	res := d.hub.kernel.Open(session.OpenRequest{
		User:   t.User,
		Device: session.DeviceID(device.ID),
		Kind:   session.KindTask,
		Now:    time.Now(),
	})
	if res.Rejected {
		reason := res.Reason
		if res.Occupier != nil {
			reason = "device busy: " + res.Occupier.User + " (" + res.Occupier.Kind + ")"
		}
		d.finishDeploy(recID, store.DeployFailed, reason)
		return nil
	}
	sid := int64(res.SessionID)
	if err := d.hub.persistNewSession(sid, device.ID, session.KindTask, t.User); err != nil {
		slog.Error("persist flash session", "err", err)
	}
	if err := d.hub.openLogger(sid); err != nil {
		slog.Error("open flash session log", "err", err)
	}
	hub := d.hub
	hub.mu.Lock()
	hub.virtualSession[device.ID] = sid
	hub.expectRuns[sid] = nil
	hub.mu.Unlock()

	rec, _ := d.store.GetDeployRecord(context.Background(), recID)
	rec.SessionID = sid
	_ = d.store.UpdateDeployRecord(context.Background(), rec)

	// Engine transport: shared serial client, or the virtual device in demo
	// mode (the same dispatch as the M1 runExpect path).
	enginePort, done := hub.flashTransport(device.ID, sid)
	if enginePort == nil {
		for _, ev := range hub.kernel.ClientDisconnected(session.SessionID(sid), time.Now()) {
			hub.apply(ev)
		}
		d.finishDeploy(recID, store.DeployFailed, "device has no shared serial port")
		return nil
	}

	writeTX := func(data []byte, secret bool) {
		if secret {
			_ = hub.writeLogMasked(sid, "TX")
		} else {
			_ = hub.writeLog(sid, "TX", data)
		}
	}
	cfg := expect.Config{
		Steps:         steps,
		DefaultWait:   15 * time.Second,
		DefaultOnFail: expect.FailAbort,
		OnSendBytes:   writeTX,
		OnStep: func(index, total int) {
			hub.notifySession(sid, protocol.Frame{Type: protocol.FrameExpectProgress, Body: &protocol.ExpectProgressFrame{
				SessionID: sid, StepIndex: index, StepTotal: total,
				StepDesc: describeStep(steps[index]), Phase: "running",
			}})
			d.emit(recID, "running", fmt.Sprintf("step %d/%d", index+1, total))
		},
	}
	runner := expect.New(cfg, enginePort)
	hub.mu.Lock()
	hub.expectRuns[sid] = runner
	hub.mu.Unlock()

	go func() {
		result := runner.Run()
		_ = done()

		// Task session ends with the run (completed -> closed, else failed),
		// mirroring the M1 expect path.
		now := time.Now()
		if result.OK() {
			for _, ev := range hub.kernel.Close(session.SessionID(sid), now) {
				hub.apply(ev)
			}
		} else {
			for _, ev := range hub.kernel.ClientDisconnected(session.SessionID(sid), now) {
				hub.apply(ev)
			}
		}

		switch {
		case result.Aborted:
			d.finishDeploy(recID, store.DeployCanceled, "aborted at step "+strconv.Itoa(result.StepIndex+1))
		case result.OK():
			d.finishDeploy(recID, store.DeploySucceeded, "flash sequence completed")
		default:
			d.finishDeploy(recID, store.DeployFailed, result.Detail)
		}
	}()
	return nil
}

// CancelDeploy stops a running deployment: SSH closes the connection, flash
// aborts the expect runner (requirement 3.6; the flash double-confirm is a
// UI concern).
func (d *deployOrchestrator) CancelDeploy(recordID int64) error {
	d.mu.Lock()
	runner := d.sshRuns[recordID]
	d.mu.Unlock()
	if runner != nil {
		runner.Cancel()
		return nil
	}
	// Flash: abort via the session's expect runner (the M1 abort path).
	rec, err := d.store.GetDeployRecord(context.Background(), recordID)
	if err != nil {
		return err
	}
	if rec.SessionID != 0 {
		d.hub.abortExpect(nil, rec.SessionID)
	}
	return nil
}

// deployLogPath is the per-deploy-record log.
func deployLogPath(dataDir string, recordID int64) string {
	return filepath.Join(dataDir, "deployments", fmt.Sprint(recordID), "log.txt")
}

// notifyWithLog wraps Notify so every line also lands in the exec log file.
func (d *deployOrchestrator) notifyWithLog(execID int64) deployer.Notify {
	return func(recordID int64, phase, line string) {
		d.log(recordID, phase, line)
		_ = deployer.WriteExecLog(d.hub.dataDir, execID, phase, line)
	}
}

// log appends to the deploy record's log file and pushes to subscribers.
func (d *deployOrchestrator) log(recordID int64, phase, line string) {
	dir := filepath.Dir(deployLogPath(d.hub.dataDir, recordID))
	if err := mkdirAll(dir); err == nil {
		f, err := openAppend(deployLogPath(d.hub.dataDir, recordID))
		if err == nil {
			fmt.Fprintf(f, "%s [%s] %s\n", time.Now().Format("2006-01-02 15:04:05.000"), phase, line)
			f.Close()
		}
	}
	d.emit(recordID, "log", phase+": "+line)
}

func (d *deployOrchestrator) finishDeploy(recID int64, status, detail string) {
	rec, err := d.store.GetDeployRecord(context.Background(), recID)
	if err != nil {
		return
	}
	// Canceled is terminal here too: a late success must not overwrite it.
	if rec.Status == store.DeployCanceled {
		return
	}
	rec.Status = status
	rec.Detail = detail
	rec.EndedAt = time.Now().Format(time.RFC3339)
	if err := d.store.UpdateDeployRecord(context.Background(), rec); err != nil {
		slog.Error("persist deploy end", "err", err)
	}
	d.emit(recID, status, detail)
	slog.Info("deploy finished", "record", recID, "status", status, "detail", detail)
}

// --- subscriptions (the build module's shape) ---

func (d *deployOrchestrator) subscribe(recordID int64, c *transport.ClientConn) {
	d.mu.Lock()
	if d.subs[recordID] == nil {
		d.subs[recordID] = map[*transport.ClientConn]struct{}{}
	}
	d.subs[recordID][c] = struct{}{}
	d.mu.Unlock()

	// Snapshot replay so a late viewer sees the current state.
	rec, err := d.store.GetDeployRecord(context.Background(), recordID)
	if err == nil {
		_ = c.Send(protocol.Frame{Type: protocol.FrameDeployEvent, Body: &protocol.DeployEventFrame{
			DeployID: recordID, Phase: rec.Status, Detail: rec.Detail,
		}})
	}
}

func (d *deployOrchestrator) unsubscribe(recordID int64, c *transport.ClientConn) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.subs[recordID], c)
}

func (d *deployOrchestrator) unsubscribeAll(c *transport.ClientConn) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for id, m := range d.subs {
		delete(m, c)
		if len(m) == 0 {
			delete(d.subs, id)
		}
	}
}

func (d *deployOrchestrator) emit(recordID int64, phase, detail string) {
	d.mu.Lock()
	targets := make([]*transport.ClientConn, 0, 2)
	for c := range d.subs[recordID] {
		targets = append(targets, c)
	}
	d.mu.Unlock()
	for _, c := range targets {
		_ = c.Send(protocol.Frame{Type: protocol.FrameDeployEvent, Body: &protocol.DeployEventFrame{
			DeployID: recordID, Phase: phase, Detail: detail,
		}})
	}
}

// marshalJSON is a tiny indirection for tests.
func marshalJSON(v any) string {
	raw, _ := json.Marshal(v)
	return string(raw)
}

// orDefault returns v when non-zero, else def.
func orDefault(v, def int) int {
	if v != 0 {
		return v
	}
	return def
}

// testSSH probes one device's SSH reachability (requirement 3.2: manual
// Test button; no automatic probing). Dial-only, no command runs.
func (d *deployOrchestrator) testSSH(host string, port int, user, password string) (bool, string) {
	if port == 0 {
		port = 22
	}
	runner := &deployer.SSHRunner{
		Host: host, Port: port, User: user, Password: password,
		Timeout: 10 * time.Second,
	}
	res := runner.Run(context.Background(), []string{"true"})
	if res.OK {
		return true, "ssh ok"
	}
	return false, res.Detail
}

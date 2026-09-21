package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/WakkeWang/EmbedFlow/server/internal/deployer"
	"github.com/WakkeWang/EmbedFlow/server/internal/sessionlog"
	"github.com/WakkeWang/EmbedFlow/server/internal/store"
)

// Deploy-module HTTP handlers (requirement 3). REST surface only; live
// progress rides the WS DeployEvent frames.

// deployRuleBody is the create/update request shape for deploy rules.
type deployRuleBody struct {
	Name            string                 `json:"name"`
	Mode            string                 `json:"mode"`
	StepsMD         string                 `json:"steps_md"`
	SSHDeviceID     int64                  `json:"ssh_device_id"`
	SSHCommands     []string               `json:"ssh_commands"`
	SSHParams       []store.DeployParamDef `json:"ssh_params"`
	SSHTimeoutSec   int                    `json:"ssh_timeout_sec"`
	FlashDeviceID   int64                  `json:"flash_device_id"`
	FlashSteps      json.RawMessage        `json:"flash_steps"`
	FlashTimeoutSec int                    `json:"flash_timeout_sec"`
	USBDisk         *store.USBDisk         `json:"usb_disk"`
}

// toPayload maps the request body into the persistence payload.
func (b deployRuleBody) toPayload() store.DeployPayload {
	p := store.DeployPayload{Mode: b.Mode, USBDisk: b.USBDisk}
	switch b.Mode {
	case "manual":
		p.StepsMD = b.StepsMD
	case "ssh":
		p.SSHDeviceID = b.SSHDeviceID
		p.SSHCommands = b.SSHCommands
		p.SSHParams = b.SSHParams
		p.SSHTimeoutSec = b.SSHTimeoutSec
	case "flash":
		p.FlashDeviceID = b.FlashDeviceID
		if len(b.FlashSteps) > 0 {
			p.FlashStepsJSON = string(b.FlashSteps)
		}
		p.FlashTimeoutSec = b.FlashTimeoutSec
	}
	return p
}

func (h *coordinator) handleCreateDeployRule(w http.ResponseWriter, r *http.Request) {
	pid, err := pathPIDValue(r, "pid")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad project id")
		return
	}
	var body deployRuleBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		writeErr(w, http.StatusBadRequest, "name required")
		return
	}
	if body.Mode != "manual" && body.Mode != "ssh" && body.Mode != "flash" {
		writeErr(w, http.StatusBadRequest, "mode must be manual, ssh or flash")
		return
	}
	id, err := h.store.CreateDeployRule(r.Context(), pid, body.Name, body.toPayload())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"id": id})
}

func (h *coordinator) handleListDeployRules(w http.ResponseWriter, r *http.Request) {
	pid, err := pathPIDValue(r, "pid")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad project id")
		return
	}
	rules, payloads, err := h.store.ListDeployRules(r.Context(), pid)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	type ruleView struct {
		store.DeployRule
		Payload store.DeployPayload `json:"payload"`
	}
	out := make([]ruleView, 0, len(rules))
	for i := range rules {
		out = append(out, ruleView{DeployRule: rules[i], Payload: payloads[i]})
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *coordinator) handleGetDeployRule(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	rule, payload, err := h.store.GetDeployRule(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no such rule")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"rule": rule, "payload": payload})
}

func (h *coordinator) handleUpdateDeployRule(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	rule, _, err := h.store.GetDeployRule(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no such rule")
		return
	}
	var body deployRuleBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		writeErr(w, http.StatusBadRequest, "name required")
		return
	}
	rule.Name = body.Name
	if err := h.store.UpdateDeployRule(r.Context(), rule, body.toPayload()); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *coordinator) handleDeleteDeployRule(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	if err := h.store.DeleteDeployRule(r.Context(), id); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// handleTriggerDeploy starts a deployment (requirement 3.6: deploy =
// rule + build record + operator params).
func (h *coordinator) handleTriggerDeploy(w http.ResponseWriter, r *http.Request) {
	user, _, ok := h.tokens.validate(bearerToken(r))
	if !ok {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var body struct {
		RuleID        int64            `json:"rule_id"`
		BuildRecordID int64            `json:"build_record_id"`
		Params        map[string]string `json:"params"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RuleID == 0 {
		writeErr(w, http.StatusBadRequest, "rule_id required")
		return
	}
	recID, err := h.deploys.Trigger(r.Context(), DeployTrigger{
		RuleID:        body.RuleID,
		BuildRecordID: body.BuildRecordID,
		Params:        body.Params,
		User:          user,
	})
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"id": recID})
}

func (h *coordinator) handleListDeployments(w http.ResponseWriter, r *http.Request) {
	pid, err := pathPIDValue(r, "pid")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad project id")
		return
	}
	recs, err := h.store.ListDeployRecords(r.Context(), pid)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if recs == nil {
		recs = []store.DeployRecord{}
	}
	writeJSON(w, http.StatusOK, recs)
}

func (h *coordinator) handleGetDeployment(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	rec, err := h.store.GetDeployRecord(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no such deploy record")
		return
	}
	writeJSON(w, http.StatusOK, rec)
}

func (h *coordinator) handleCancelDeployment(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	rec, err := h.store.GetDeployRecord(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no such deploy record")
		return
	}
	if rec.Status != store.DeployRunning {
		writeErr(w, http.StatusConflict, "deploy is not running")
		return
	}
	if err := h.deploys.CancelDeploy(id); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// handleDeleteDeployRecord removes a terminal deploy record and its log
// dir. A running record is not deletable -- cancel it first (the build
// records' delete semantics, applied to deploys).
func (h *coordinator) handleDeleteDeployRecord(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	rec, err := h.store.GetDeployRecord(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no such deploy record")
		return
	}
	if rec.Status == store.DeployRunning {
		writeErr(w, http.StatusConflict, "deploy is running: cancel it first")
		return
	}
	if err := h.store.DeleteDeployRecord(r.Context(), id); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = os.RemoveAll(filepath.Dir(deployLogPath(h.dataDir, id)))
	w.WriteHeader(http.StatusNoContent)
}

// handleDeployLogTail serves the tail of a deploy record's log.
func (h *coordinator) handleDeployLogTail(w http.ResponseWriter, r *http.Request) {	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	n := 1000
	if v := r.URL.Query().Get("n"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 && parsed <= 10000 {
			n = parsed
		}
	}
	tail, err := sessionlog.Tail(deployLogPath(h.dataDir, id), n)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no log")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"tail": tail})
}

// --- devices: SSH fields + test (requirement 3.2) ---

// devicePayload is the device create/update body. Password rides write-only:
// it is encrypted for storage and never echoed back (SshSet says whether one
// is stored).
type devicePayload struct {
	Name     string `json:"name"`
	Project  string `json:"project"`
	SSHHost  string `json:"ssh_host"`
	SSHPort  int    `json:"ssh_port"`
	SSHUser  string `json:"ssh_user"`
	SSHPassword string `json:"ssh_password"`
	Note     string `json:"note"`
}

func (h *coordinator) handleCreateDeviceV2(w http.ResponseWriter, r *http.Request) {
	var p devicePayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil || p.Name == "" {
		writeErr(w, http.StatusBadRequest, "name required")
		return
	}
	passEnc, err := h.encryptPassword(p.SSHPassword)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	id, err := h.store.CreateDeviceFull(r.Context(), store.Device{
		Name: p.Name, Project: p.Project, SSHHost: p.SSHHost,
		SSHPort: p.SSHPort, SSHUser: p.SSHUser, Note: p.Note,
	}, passEnc)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"id": id})
}

func (h *coordinator) handleUpdateDevice(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	if _, err := h.store.GetDevice(r.Context(), id); err != nil {
		writeErr(w, http.StatusNotFound, "no such device")
		return
	}
	var p devicePayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil || p.Name == "" {
		writeErr(w, http.StatusBadRequest, "name required")
		return
	}
	passEnc, err := h.encryptPassword(p.SSHPassword)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.store.UpdateDeviceSSH(r.Context(), store.Device{
		ID: id, Name: p.Name, Project: p.Project, SSHHost: p.SSHHost,
		SSHPort: p.SSHPort, SSHUser: p.SSHUser, Note: p.Note,
	}, passEnc); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// handleSSHTest probes a device's SSH reachability on demand (requirement
// 3.2: no automatic probing; the Test button runs one check).
func (h *coordinator) handleSSHTest(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	device, err := h.store.GetDeviceFull(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no such device")
		return
	}
	if device.SSHHost == "" {
		writeErr(w, http.StatusBadRequest, "device has no SSH configuration")
		return
	}
	password := h.decryptPassword(device.SshPassEnc)
	ok, detail := h.deploys.testSSH(device.SSHHost, device.SSHPort, device.SSHUser, password)
	writeJSON(w, http.StatusOK, map[string]any{"ok": ok, "detail": detail})
}

// handleDeployUSBZip streams the rule's USB stick zip for one build record
// (requirement 3.1.3 second half: 目录名对齐实际路径 + 操作指引, generated
// on request, nothing cached on disk). Every configured glob must match
// exactly one artifact or the whole zip fails with a 422 listing them.
func (h *coordinator) handleDeployUSBZip(w http.ResponseWriter, r *http.Request) {
	ruleID, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad rule id")
		return
	}
	buildRecordID, err := strconv.ParseInt(r.URL.Query().Get("build_record_id"), 10, 64)
	if err != nil || buildRecordID <= 0 {
		writeErr(w, http.StatusBadRequest, "build_record_id required")
		return
	}
	rule, payload, err := h.store.GetDeployRule(r.Context(), ruleID)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no such rule")
		return
	}
	if payload.USBDisk == nil {
		writeErr(w, http.StatusBadRequest, "rule has no USB stick configuration")
		return
	}
	// The record must exist (its artifacts are the zip's content); the
	// lookup doubles as the existence check -- no row, no zip.
	if _, err := h.store.GetBuildRecord(r.Context(), buildRecordID); err != nil {
		writeErr(w, http.StatusNotFound, "no such build record")
		return
	}
	arts, err := h.store.ArtifactsForRecord(r.Context(), buildRecordID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	views := make([]deployer.ArtifactView, 0, len(arts))
	for _, a := range arts {
		views = append(views, deployer.ArtifactView{ID: a.ID, Name: a.Name, Size: a.Size, Checksum: a.Checksum})
	}

	// The root folder name is a template so it tracks the build's version.
	rootName, _ := deployer.Render(payload.USBDisk.RootName, deployer.Variables{
		Artifact: deployer.ArtifactVariables(views),
		Device:   map[string]string{},
		Params:   map[string]string{},
	})

	// Resolve before touching the response: a bad glob or a missing
	// artifact fails as a clean 422, never as a half-downloaded zip.
	if err := deployer.ResolveUSB(*payload.USBDisk, views); err != nil {
		var ze *deployer.USBZipError
		missing := []string{}
		if errors.As(err, &ze) {
			missing = ze.Missing
		}
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": err.Error(), "missing": missing})
		return
	}

	open := func(a deployer.ArtifactView) (io.ReadCloser, error) {
		return os.Open(artifactFilePath(h.dataDir, buildRecordID, a.Name))
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment; filename=%q", "usb-deploy-rule"+strconv.FormatInt(rule.ID, 10)+"-build"+strconv.FormatInt(buildRecordID, 10)+".zip"))
	if err := deployer.BuildUSBZip(w, rootName, *payload.USBDisk, views, open, time.Now()); err != nil {
		// Past this point headers are out; log it -- the client sees a
		// truncated zip and the server log tells why.
		slog.Error("usb zip stream failed", "rule", rule.ID, "build", buildRecordID, "err", err)
	}
}

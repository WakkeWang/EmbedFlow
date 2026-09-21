package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/WakkeWang/EmbedFlow/server/internal/batch"
	"github.com/WakkeWang/EmbedFlow/server/internal/sessionlog"
	"github.com/WakkeWang/EmbedFlow/server/internal/store"
)

// Build-module HTTP handlers (requirement 2). REST surface only -- live
// progress rides the WS BuildEvent frames (see coordinator.onBuildCtrl).

func pathPID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("pid"), 10, 64)
}

// decodeBuildItem reads the item payload; artifacts ride as a JSON array,
// prereq groups as raw JSON ([[1,2],[3]]).
func decodeBuildItem(r *http.Request, projectID int64) (store.BuildItem, error) {
	var req struct {
		Name        string          `json:"name"`
		SourceType  string          `json:"source_type"`
		GitURL      string          `json:"git_url"`
		GitBranch   string          `json:"git_branch"`
		GitCommit   string          `json:"git_commit"`
		CheckLatest bool            `json:"check_latest"`
		LocalPath   string          `json:"local_path"`
		Command     string          `json:"command"`
		Artifacts   []string        `json:"artifacts"`
		TimeoutSec  int             `json:"timeout_sec"`
		VersionCmd  string          `json:"version_cmd"`
		PrereqJSON  json.RawMessage `json:"prereq_json"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return store.BuildItem{}, err
	}
	if req.Name == "" || req.Command == "" {
		return store.BuildItem{}, fmt.Errorf("name and command required")
	}
	if req.SourceType != "git" && req.SourceType != "local" {
		return store.BuildItem{}, fmt.Errorf("source_type must be git or local")
	}
	if req.SourceType == "git" && req.GitURL == "" {
		return store.BuildItem{}, fmt.Errorf("git source requires git_url")
	}
	if req.SourceType == "local" && req.LocalPath == "" {
		return store.BuildItem{}, fmt.Errorf("local source requires local_path")
	}
	// Prereq groups normalize to the canonical "[[1,2],[3]]" array form.
	// A client that stringified twice (historical bug) sends a JSON string
	// literal; unwrap one level so the stored payload is always an array.
	prereq := normalizePrereqJSON(req.PrereqJSON)
	it := store.BuildItem{
		ProjectID:   projectID,
		Name:        req.Name,
		SourceType:  req.SourceType,
		GitURL:      req.GitURL,
		GitBranch:   req.GitBranch,
		GitCommit:   req.GitCommit,
		CheckLatest: req.CheckLatest,
		LocalPath:   req.LocalPath,
		Command:     req.Command,
		Artifacts:   req.Artifacts,
		TimeoutSec:  req.TimeoutSec,
		VersionCmd:  req.VersionCmd,
		PrereqJSON:  prereq,
	}
	return it, nil
}

// normalizePrereqJSON coerces the request's prereq groups into the canonical
// array-of-arrays JSON. Accepts the raw array, or one level of JSON-string
// encoding (a historical client bug stringified twice); anything else
// stores as "[]" (empty groups).
func normalizePrereqJSON(raw json.RawMessage) string {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == `""` {
		return "[]"
	}
	var groups [][]int64
	if err := json.Unmarshal(raw, &groups); err == nil {
		out, err := json.Marshal(groups)
		if err == nil {
			return string(out)
		}
		return "[]"
	}
	// Not an array: try one unwrap (the raw value is a JSON string).
	var inner string
	if err := json.Unmarshal(raw, &inner); err == nil {
		inner = strings.TrimSpace(inner)
		if err := json.Unmarshal([]byte(inner), &groups); err == nil {
			out, err := json.Marshal(groups)
			if err == nil {
				return string(out)
			}
		}
	}
	return "[]"
}

func (h *coordinator) handleCreateBuildItem(w http.ResponseWriter, r *http.Request) {
	pid, err := pathPID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad project id")
		return
	}
	it, err := decodeBuildItem(r, pid)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.store.CreateBuildItem(r.Context(), it)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"id": id})
}

func (h *coordinator) handleListBuildItems(w http.ResponseWriter, r *http.Request) {
	pid, err := pathPID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad project id")
		return
	}
	items, err := h.store.ListBuildItemsByProject(r.Context(), pid)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []store.BuildItem{}
	}
	writeJSON(w, http.StatusOK, items)
}

// handleListAllBuildItems returns every item across projects (the
// prerequisite picker's source: prereq refs may cross projects, 2.1).
func (h *coordinator) handleListAllBuildItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.store.ListAllBuildItems(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []store.BuildItem{}
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *coordinator) handleUpdateBuildItem(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	existing, err := h.store.GetBuildItem(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no such build item")
		return
	}
	it, err := decodeBuildItem(r, existing.ProjectID)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	it.ID = id
	if err := h.store.UpdateBuildItem(r.Context(), it); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *coordinator) handleDeleteBuildItem(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	if err := h.store.DeleteBuildItem(r.Context(), id); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- batches ---

func (h *coordinator) handleCreateBatch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProjectID int64   `json:"project_id"`
		ItemIDs   []int64 `json:"item_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ProjectID == 0 {
		writeErr(w, http.StatusBadRequest, "project_id and item_ids required")
		return
	}
	user, _, ok := h.tokens.validate(bearerToken(r))
	if !ok {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := h.builds.CreateBatch(r.Context(), req.ProjectID, req.ItemIDs, user)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"id": id})
}

func (h *coordinator) handleListBatches(w http.ResponseWriter, r *http.Request) {
	pid, err := strconv.ParseInt(r.URL.Query().Get("project_id"), 10, 64)
	if err != nil || pid == 0 {
		writeErr(w, http.StatusBadRequest, "project_id required")
		return
	}
	batches, err := h.store.ListBatchesByProject(r.Context(), pid)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if batches == nil {
		batches = []store.Batch{}
	}
	writeJSON(w, http.StatusOK, batches)
}

// handleListBuildRecords returns a project's records newest first
// (requirement 2.4 build history); one query instead of the per-batch walk
// the frontend did before.
func (h *coordinator) handleListBuildRecords(w http.ResponseWriter, r *http.Request) {
	pid, err := strconv.ParseInt(r.URL.Query().Get("project_id"), 10, 64)
	if err != nil || pid == 0 {
		writeErr(w, http.StatusBadRequest, "project_id required")
		return
	}
	records, err := h.store.RecordsForProject(r.Context(), pid)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if records == nil {
		records = []store.BuildRecord{}
	}
	writeJSON(w, http.StatusOK, records)
}

func (h *coordinator) handleGetBatch(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	b, err := h.store.GetBatch(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no such batch")
		return
	}
	records, err := h.store.RecordsForBatch(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if records == nil {
		records = []store.BuildRecord{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"batch": b, "records": records})
}

func (h *coordinator) handleCancelBatch(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	if err := h.builds.CancelBatch(id); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- build records ---

func (h *coordinator) handleGetBuildRecord(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	rec, err := h.store.GetBuildRecord(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no such record")
		return
	}
	writeJSON(w, http.StatusOK, rec)
}

func (h *coordinator) handleRecordArtifacts(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	arts, err := h.store.ArtifactsForRecord(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if arts == nil {
		arts = []store.Artifact{}
	}
	writeJSON(w, http.StatusOK, arts)
}

func (h *coordinator) handleBuildLogTail(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
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
	// sessionlog.Tail reads backward chunk-wise: huge logs never load fully
	// into memory (the same discipline as the session log tail, 3.3).
	tail, err := sessionlog.Tail(buildLogPath(h.dataDir, id), n)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no build log")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"tail": tail})
}

// deleteBuildRecordFiles removes the record's artifacts from disk
// (mode=artifacts; mode=record keeps them as orphans per requirement 2.4's
// "keep artifacts, delete record" option).
func deleteBuildRecordFiles(dataDir string, recordID int64) {
	_ = os.RemoveAll(artifactDirPath(dataDir, recordID))
}

func (h *coordinator) handleDeleteBuildRecord(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	rec, err := h.store.GetBuildRecord(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no such record")
		return
	}
	// A record still pending/building belongs to a live batch: cancel the
	// batch instead (the UI's third delete option maps there).
	if rec.Status == batch.RecPending || rec.Status == batch.RecBuilding {
		writeErr(w, http.StatusConflict, "record is live: cancel its batch instead")
		return
	}
	mode := r.URL.Query().Get("mode")
	if mode == "artifacts" {
		deleteBuildRecordFiles(h.dataDir, id)
	}
	// Artifact rows always go with the record (foreign key); the mode
	// decides only whether the FILES on disk survive as orphans
	// (requirement 2.4: 仅删记录时产物成孤儿、占盘，直至手动清理).
	if err := h.store.DeleteArtifactsForRecord(r.Context(), id); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.store.DeleteBuildRecord(r.Context(), id); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *coordinator) handleDeleteBuildRecordsAll(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProjectID int64  `json:"project_id"`
		Mode      string `json:"mode"` // record | artifacts
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ProjectID == 0 {
		writeErr(w, http.StatusBadRequest, "project_id required")
		return
	}
	records, err := h.store.RecordsForProject(r.Context(), req.ProjectID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, rec := range records {
		if rec.Status == batch.RecPending || rec.Status == batch.RecBuilding {
			continue // live work is not deletable
		}
		if req.Mode == "artifacts" {
			deleteBuildRecordFiles(h.dataDir, rec.ID)
		}
		if err := h.store.DeleteArtifactsForRecord(r.Context(), rec.ID); err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := h.store.DeleteBuildRecord(r.Context(), rec.ID); err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *coordinator) handleArtifactDownload(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	art, err := h.store.GetArtifact(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no such artifact")
		return
	}
	data, err := os.ReadFile(artifactFilePath(h.dataDir, art.BuildRecordID, art.Name))
	if err != nil {
		writeErr(w, http.StatusNotFound, "artifact file missing")
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", art.Name))
	w.Header().Set("X-Checksum", art.Checksum)
	_, _ = w.Write(data)
}

// --- settings ---

func (h *coordinator) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	all, err := h.store.AllSettings(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Defaults surface explicitly so the admin page shows effective values.
	if all[SettingMaxParallelBatches] == "" {
		all[SettingMaxParallelBatches] = "4"
	}
	if all[SettingChecksum] == "" {
		all[SettingChecksum] = "sha256"
	}
	writeJSON(w, http.StatusOK, all)
}

func (h *coordinator) handlePutSettings(w http.ResponseWriter, r *http.Request) {
	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "bad body")
		return
	}
	allowed := map[string]bool{
		SettingMaxParallelBatches: true,
		SettingChecksum:           true,
		SettingPublishWhitelist:   true,
		SettingTmpDir:             true,
	}
	for k, v := range req {
		if !allowed[k] {
			writeErr(w, http.StatusBadRequest, "unknown setting: "+k)
			return
		}
		if k == SettingChecksum && v != "sha256" && v != "md5" {
			writeErr(w, http.StatusBadRequest, "checksum must be sha256 or md5")
			return
		}
	}
	for k, v := range req {
		if err := h.store.SetSetting(r.Context(), k, v); err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	// Live settings take effect immediately (requirement 1.6: the parallel
	// cap applies to admissions from now on).
	h.builds.applyMaxParallelSetting()
	w.WriteHeader(http.StatusNoContent)
}

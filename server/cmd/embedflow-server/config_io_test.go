package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/WakkeWang/EmbedFlow/server/internal/secretbox"
	"github.com/WakkeWang/EmbedFlow/server/internal/store"
)

// Config copy/export/import end-to-end (requirement 1.6: 跨工程复制粘贴 +
// JSON 导出/导入): REST calls against the coordinator over the real mux,
// the apiCall pattern from roles_test.go.

func createProjectForTest(t *testing.T, s *srv, tok, name string) int64 {
	t.Helper()
	resp := apiCall(t, s, tok, "POST", "/api/projects", map[string]string{"name": name})
	if resp.StatusCode != 200 {
		t.Fatalf("create project = %d", resp.StatusCode)
	}
	var out struct {
		ID int64 `json:"id"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return out.ID
}

func bodyBytes(t *testing.T, v any) []byte {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return raw
}

func decodeBody[T any](t *testing.T, data []byte) T {
	t.Helper()
	var out T
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("decode %s: %v", string(data[:min(len(data), 120)]), err)
	}
	return out
}

func contains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}

func TestConfigIO_CopyAndExportImport(t *testing.T) {
	s := newTestServer(t)
	adminTok := s.hub.tokens.issueForTesting("admin")

	src := createProjectForTest(t, s, adminTok, "P-io-src")
	dst := createProjectForTest(t, s, adminTok, "P-io-dst")

	// Seed one rule in src via the deploy-rule endpoint.
	resp := apiCall(t, s, adminTok, "POST", "/api/projects/"+itoa(src)+"/deploy-rules", map[string]any{
		"name":     "flash-usb",
		"mode":     "manual",
		"steps_md": "1. flash it",
		"usb_disk": map[string]any{
			"root_name": "ver-${artifact.name}",
			"layout":    []map[string]string{{"glob": "*.itb", "dir": "."}},
			"checksums": []map[string]string{{"name": "itbfile", "glob": "board.itb"}},
		},
	})
	if resp.StatusCode != 200 {
		t.Fatalf("create deploy rule = %d", resp.StatusCode)
	}
	created := decodeBody[struct {
		ID int64 `json:"id"`
	}](t, readRespBody(t, resp))

	// Member cannot copy (write-config is admin-only).
	memberTok := s.hub.tokens.issueForRole("alice", "member")
	if resp := apiCall(t, s, memberTok, "POST", "/api/config-objects/"+itoa(created.ID)+"/copy",
		map[string]any{"target_project_id": dst}); resp.StatusCode != 403 {
		t.Fatalf("member copy = %d, want 403", resp.StatusCode)
	}

	// Admin copies the rule into dst.
	resp = apiCall(t, s, adminTok, "POST", "/api/config-objects/"+itoa(created.ID)+"/copy",
		map[string]any{"target_project_id": dst})
	if resp.StatusCode != 200 {
		t.Fatalf("admin copy = %d", resp.StatusCode)
	}
	copied := decodeBody[struct {
		ID int64 `json:"id"`
	}](t, readRespBody(t, resp))
	if copied.ID == 0 || copied.ID == created.ID {
		t.Fatalf("copy id = %d (source %d)", copied.ID, created.ID)
	}

	// The copy is in dst with the same payload (usb_disk rides along).
	_, payload, err := s.st.GetDeployRule(t.Context(), copied.ID)
	if err != nil {
		t.Fatalf("get copied rule: %v", err)
	}
	if payload.USBDisk == nil || payload.USBDisk.RootName != "ver-${artifact.name}" ||
		len(payload.USBDisk.Layout) != 1 || len(payload.USBDisk.Checksums) != 1 {
		t.Fatalf("copied usb_disk = %+v", payload.USBDisk)
	}

	// Export dst as JSON: one deploy_rule item with the usb_disk payload.
	resp = apiCall(t, s, adminTok, "GET", "/api/projects/"+itoa(dst)+"/export", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("export = %d", resp.StatusCode)
	}
	raw, _ := io.ReadAll(resp.Body)
	if contains(string(raw), `"password"`) {
		t.Fatal("rule export should never carry credentials")
	}
	bundle := decodeBody[struct {
		Version int `json:"version"`
		Project string `json:"project"`
		Items []struct {
			Kind    string          `json:"kind"`
			Name    string          `json:"name"`
			Payload json.RawMessage `json:"payload"`
		} `json:"items"`
	}](t, raw)
	if bundle.Version != 1 || len(bundle.Items) != 1 || bundle.Items[0].Kind != "deploy_rule" ||
		bundle.Items[0].Name != "flash-usb" {
		t.Fatalf("export bundle = %+v", bundle)
	}
	if !contains(string(bundle.Items[0].Payload), "usb_disk") {
		t.Fatalf("export payload lost usb_disk: %s", bundle.Items[0].Payload)
	}

	// Import the same bundle into a third project.
	third := createProjectForTest(t, s, adminTok, "P-io-third")
	rebundle := bodyBytes(t, map[string]any{"version": 1, "project": bundle.Project, "items": bundle.Items})
	importResp, err := postRaw(t, s, adminTok, "/api/projects/"+itoa(third)+"/import", rebundle)
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if importResp.StatusCode != 200 {
		t.Fatalf("import = %d", importResp.StatusCode)
	}
	tally := decodeBody[struct {
		Imported int `json:"imported"`
		Skipped  int `json:"skipped"`
	}](t, readRespBody(t, importResp))
	if tally.Imported != 1 || tally.Skipped != 0 {
		t.Fatalf("import tally = %+v", tally)
	}

	// The imported rule lists in the third project.
	listResp := apiCall(t, s, adminTok, "GET", "/api/projects/"+itoa(third)+"/deploy-rules", nil)
	rules := decodeBody[[]store.DeployRule](t, readRespBody(t, listResp))
	if len(rules) != 1 || rules[0].Name != "flash-usb" {
		t.Fatalf("imported rules = %+v", rules)
	}
}

func TestConfigIO_ExportDevices_ExcludesCredential(t *testing.T) {
	s := newTestServer(t)
	// The device create path encrypts the password; a key-less test hub
	// returns 500. Wire a deterministic test key.
	box, err := secretbox.New(bytes.Repeat([]byte("k"), 32))
	if err != nil {
		t.Fatalf("secretbox: %v", err)
	}
	s.hub.secrets = box
	adminTok := s.hub.tokens.issueForTesting("admin")
	proj := createProjectForTest(t, s, adminTok, "P-dev")
	// Device in the project's name (devices carry the project string).
	if resp := apiCall(t, s, adminTok, "POST", "/api/devices/v2", map[string]any{
		"name": "board-a", "project": "P-dev", "ssh_host": "10.0.0.9",
		"ssh_port": 22, "ssh_user": "root", "ssh_password": "sekret", "note": "bench 1",
	}); resp.StatusCode != 200 {
		t.Fatalf("create device = %d", resp.StatusCode)
	}
	resp := apiCall(t, s, adminTok, "GET", "/api/projects/"+itoa(proj)+"/export?kinds=device", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("export = %d", resp.StatusCode)
	}
	raw := readRespBody(t, resp)
	if contains(string(raw), "sekret") {
		t.Fatal("export leaks the SSH password")
	}
	bundle := decodeBody[struct {
		Items []struct {
			Kind    string `json:"kind"`
			Name    string `json:"name"`
			SSHHost string `json:"ssh_host"`
			SSHUser string `json:"ssh_user"`
		} `json:"items"`
	}](t, raw)
	if len(bundle.Items) != 1 || bundle.Items[0].Kind != "device" || bundle.Items[0].SSHHost != "10.0.0.9" {
		t.Fatalf("device export = %+v", bundle.Items)
	}
}

func TestConfigIO_ImportDevice_NoPassword(t *testing.T) {
	s := newTestServer(t)
	adminTok := s.hub.tokens.issueForTesting("admin")
	proj := createProjectForTest(t, s, adminTok, "P-imp-dev")
	bundle := bodyBytes(t, map[string]any{
		"version": 1, "project": "wherever",
		"items": []map[string]any{{
			"kind": "device", "name": "board-imp",
			"ssh_host": "10.0.0.10", "ssh_port": 22, "ssh_user": "root", "note": "n",
		}},
	})
	importResp, err := postRaw(t, s, adminTok, "/api/projects/"+itoa(proj)+"/import", bundle)
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if importResp.StatusCode != 200 {
		t.Fatalf("import = %d", importResp.StatusCode)
	}
	// The device exists without a stored password (ssh_set false).
	devs, err := s.st.ListDevices(t.Context())
	if err != nil {
		t.Fatalf("list devices: %v", err)
	}
	var found *store.Device
	for i := range devs {
		if devs[i].Name == "board-imp" {
			found = &devs[i]
		}
	}
	if found == nil {
		t.Fatal("imported device missing")
	}
	if found.SshSet {
		t.Fatal("imported device has a credential")
	}
}

func TestConfigIO_ExportUnknownKind(t *testing.T) {
	s := newTestServer(t)
	adminTok := s.hub.tokens.issueForTesting("admin")
	proj := createProjectForTest(t, s, adminTok, "P-bad")
	resp := apiCall(t, s, adminTok, "GET", "/api/projects/"+itoa(proj)+"/export?kinds=release_note", nil)
	if resp.StatusCode != 400 {
		t.Fatalf("unknown kind = %d, want 400", resp.StatusCode)
	}
}

// readRespBody drains an apiCall response into bytes.
func readRespBody(t *testing.T, resp *http.Response) []byte {
	t.Helper()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return raw
}

func TestDeployRecord_DeleteTerminal(t *testing.T) {
	s := newTestServer(t)
	adminTok := s.hub.tokens.issueForTesting("admin")
	proj := createProjectForTest(t, s, adminTok, "P-del")

	// A manual deploy completes immediately: terminal, deletable.
	resp := apiCall(t, s, adminTok, "POST", "/api/deployments", map[string]any{
		"rule_id": seedManualDeployRule(t, s, adminTok, proj), "build_record_id": 0,
	})
	if resp.StatusCode != 200 {
		t.Fatalf("trigger = %d", resp.StatusCode)
	}
	rec := decodeBody[struct {
		ID int64 `json:"id"`
	}](t, readRespBody(t, resp))

	// A running record cannot be deleted: trigger another one whose rule
	// points at a device without SSH config (stays running? no -- it fails
	// fast). Instead simulate: the sweep-free way is to create a record
	// through the store directly.
	runningID, err := s.st.CreateDeployRecord(t.Context(), store.DeployRecord{
		ProjectID: proj, RuleID: 1, DeviceID: 1, Executor: "t",
		Status: store.DeployRunning, StartedAt: "2026-09-21T00:00:00Z",
	})
	if err != nil {
		t.Fatalf("seed running record: %v", err)
	}
	if resp := apiCall(t, s, adminTok, "DELETE", "/api/deployments/"+itoa(runningID), nil); resp.StatusCode != 409 {
		t.Fatalf("delete running = %d, want 409", resp.StatusCode)
	}

	// The terminal record deletes with 204 and disappears.
	if resp := apiCall(t, s, adminTok, "DELETE", "/api/deployments/"+itoa(rec.ID), nil); resp.StatusCode != 204 {
		t.Fatalf("delete terminal = %d, want 204", resp.StatusCode)
	}
	if resp := apiCall(t, s, adminTok, "GET", "/api/deployments/"+itoa(rec.ID), nil); resp.StatusCode != 404 {
		t.Fatalf("get deleted = %d, want 404", resp.StatusCode)
	}
	// Deleting an unknown id is 404, not a panic.
	if resp := apiCall(t, s, adminTok, "DELETE", "/api/deployments/99999", nil); resp.StatusCode != 404 {
		t.Fatalf("delete unknown = %d, want 404", resp.StatusCode)
	}
}

// seedManualDeployRule creates a manual deploy rule and returns its id.
func seedManualDeployRule(t *testing.T, s *srv, tok string, pid int64) int64 {
	t.Helper()
	resp := apiCall(t, s, tok, "POST", "/api/projects/"+itoa(pid)+"/deploy-rules", map[string]any{
		"name": "manual-del", "mode": "manual", "steps_md": "go",
	})
	if resp.StatusCode != 200 {
		t.Fatalf("seed rule = %d", resp.StatusCode)
	}
	return decodeBody[struct {
		ID int64 `json:"id"`
	}](t, readRespBody(t, resp)).ID
}

// postRaw issues a POST with raw bytes (the import path re-sends an
// exported bundle verbatim).
func postRaw(t *testing.T, s *srv, tok, path string, body []byte) (*http.Response, error) {
	t.Helper()
	req, err := http.NewRequest("POST", s.URL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	return http.DefaultClient.Do(req)
}

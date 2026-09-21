package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/WakkeWang/EmbedFlow/server/internal/store"
)

// Cross-project copy + JSON export/import (requirement 1.6: 跨工程复制粘贴，
// JSON 导出/导入 -- 开源后跨团队分享配置). One generic surface over the
// unified config_objects table; devices ride the same export shape with
// their own endpoints.

// copyBody is the cross-project copy request.
type copyBody struct {
	TargetProjectID int64 `json:"target_project_id"`
}

// handleCopyConfigObject duplicates one config object into the target
// project. Valid kinds: the config_objects families (build_item,
// deploy_rule, flash_rule) -- records and devices are not config objects.
func (h *coordinator) handleCopyConfigObject(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	var body copyBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.TargetProjectID == 0 {
		writeErr(w, http.StatusBadRequest, "target_project_id required")
		return
	}
	if _, err := h.store.GetProject(r.Context(), body.TargetProjectID); err != nil {
		writeErr(w, http.StatusNotFound, "no such target project")
		return
	}
	newID, err := h.store.CopyConfigObject(r.Context(), id, body.TargetProjectID)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"id": newID})
}

// configExportItem is one exported rule (or device) in the JSON bundle.
type configExportItem struct {
	Kind        string          `json:"kind"`
	Name        string          `json:"name"`
	Payload     json.RawMessage `json:"payload"`
	// Device-only fields (kind=device): SSH fields without the credential.
	SSHHost     string          `json:"ssh_host,omitempty"`
	SSHPort     int             `json:"ssh_port,omitempty"`
	SSHUser     string          `json:"ssh_user,omitempty"`
	Note        string          `json:"note,omitempty"`
}

// configExport is the whole bundle for one project.
type configExport struct {
	Version int                `json:"version"` // 1
	Project string             `json:"project"`
	Items   []configExportItem `json:"items"`
}

var exportableKinds = map[string]bool{
	store.KindBuildItem:  true,
	store.KindDeployRule: true,
	store.KindFlashRule:  true,
}

// handleExportProject serves the project's rules (+ optionally devices) as
// one JSON bundle. kinds= limits the families; kinds=device includes
// devices (SSH host/port/user + note -- the encrypted password never
// leaves the server).
func (h *coordinator) handleExportProject(w http.ResponseWriter, r *http.Request) {
	pid, err := pathPIDValue(r, "pid")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad project id")
		return
	}
	proj, err := h.store.GetProject(r.Context(), pid)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no such project")
		return
	}
	want := map[string]bool{}
	withDevices := false
	if v := r.URL.Query().Get("kinds"); v != "" {
		for _, k := range strings.Split(v, ",") {
			k = strings.TrimSpace(k)
			if k == "device" {
				withDevices = true
				continue
			}
			if !exportableKinds[k] {
				writeErr(w, http.StatusBadRequest, "unknown kind: "+k)
				return
			}
			want[k] = true
		}
	} else {
		for k := range exportableKinds {
			want[k] = true
		}
	}

	out := configExport{Version: 1, Project: proj.Name, Items: []configExportItem{}}
	for kind := range want {
		objects, err := h.store.ListConfigObjects(r.Context(), pid, kind)
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		for _, o := range objects {
			out.Items = append(out.Items, configExportItem{
				Kind: kind, Name: o.Name, Payload: json.RawMessage(o.PayloadJSON),
			})
		}
	}
	if withDevices {
		devices, err := h.store.ListDevices(r.Context())
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		for _, d := range devices {
			if d.Project != proj.Name {
				continue
			}
			out.Items = append(out.Items, configExportItem{
				Kind: "device", Name: d.Name,
				SSHHost: d.SSHHost, SSHPort: d.SSHPort, SSHUser: d.SSHUser, Note: d.Note,
			})
		}
	}
	w.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment; filename=%q", "embedflow-export-"+proj.Name+".json"))
	writeJSON(w, http.StatusOK, out)
}

// handleImportProject loads an export bundle into the target project.
// Every item is inserted as a new row (no name dedup -- the caller sees
// both and deletes by hand; config is cheap, silence is not). Devices are
// created without a stored password (the operator re-enters it -- the
// export never carries credentials).
func (h *coordinator) handleImportProject(w http.ResponseWriter, r *http.Request) {
	pid, err := pathPIDValue(r, "pid")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad project id")
		return
	}
	if _, err := h.store.GetProject(r.Context(), pid); err != nil {
		writeErr(w, http.StatusNotFound, "no such project")
		return
	}
	var bundle configExport
	if err := json.NewDecoder(r.Body).Decode(&bundle); err != nil || bundle.Version != 1 {
		writeErr(w, http.StatusBadRequest, "bad bundle (want version 1)")
		return
	}
	imported, skipped := 0, 0
	for _, it := range bundle.Items {
		switch it.Kind {
		case "device":
			if _, err := h.store.CreateDeviceFull(r.Context(), store.Device{
				Name: it.Name, Project: "", SSHHost: it.SSHHost, SSHPort: it.SSHPort,
				SSHUser: it.SSHUser, Note: it.Note,
			}, ""); err != nil {
				skipped++
			} else {
				imported++
			}
		case store.KindBuildItem, store.KindDeployRule, store.KindFlashRule:
			if _, err := h.store.CreateConfigObject(r.Context(), store.ConfigObject{
				ProjectID: pid, Kind: it.Kind, Name: it.Name, PayloadJSON: string(it.Payload),
			}); err != nil {
				skipped++
			} else {
				imported++
			}
		default:
			skipped++
		}
	}
	writeJSON(w, http.StatusOK, map[string]int{"imported": imported, "skipped": skipped})
}

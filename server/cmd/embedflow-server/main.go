// Command embedflow-server is the EmbedFlow server binary. M1 scope: the
// session layer end to end -- store + session kernel + WebSocket transport +
// session logs + expect runs + human confirmations + HTTP API, with the
// standard data layout (CEO-10A), startup sweep (CEO-7A) and tick loop.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/WakkeWang/EmbedFlow/pkg/protocol"
	"github.com/WakkeWang/EmbedFlow/server/internal/authmw"
	"github.com/WakkeWang/EmbedFlow/server/internal/expect"
	"github.com/WakkeWang/EmbedFlow/server/internal/paths"
	"github.com/WakkeWang/EmbedFlow/server/internal/session"
	"github.com/WakkeWang/EmbedFlow/server/internal/sessionlog"
	"github.com/WakkeWang/EmbedFlow/server/internal/store"
	"github.com/WakkeWang/EmbedFlow/server/internal/transport"
	"github.com/WakkeWang/EmbedFlow/server/internal/webui"
)

var version = "dev"

func main() {
	var (
		addr      = flag.String("addr", ":8420", "listen address")
		data      = flag.String("data", "", "data directory (default ./data)")
		demo      = flag.Bool("demo", false, "run with the built-in virtual device (CEO-18A)")
		probeEcho = flag.Bool("probe-echo", false, "latency-probe mode (issue #2): echo every client binary frame back unchanged, no store, no sessions")
	)
	flag.Parse()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	if *probeEcho {
		runProbeEcho(*addr)
		return
	}

	dataDir := *data
	if dataDir == "" {
		dataDir = "data"
	}
	if err := paths.Ensure(dataDir); err != nil {
		slog.Error("create data dir layout", "err", err)
		os.Exit(1)
	}

	st, err := store.Open(filepath.Join(dataDir, "embedflow.db"))
	if err != nil {
		slog.Error("open store", "err", err)
		os.Exit(1)
	}
	defer st.Close()

	kernel := session.New(session.Options{IdleTimeout: session.IdleTimeoutDefault})
	hub := newCoordinator(kernel, st, dataDir)

	bootstrapAdmin(st)

	// Startup sweep (CEO-7A).
	if err := hub.sweepStartup(); err != nil {
		slog.Error("startup sweep", "err", err)
		os.Exit(1)
	}

	// Build-module startup sweep (M2 mirrors CEO-7A): a restart cancels
	// in-flight batches -- their tmp dirs and half-run states are not
	// resumable, so honest terminal states beat zombie "running" rows.
	if err := hub.builds.sweepStartup(); err != nil {
		slog.Error("build startup sweep", "err", err)
		os.Exit(1)
	}

	// Liveness tick loop.
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			hub.tick(time.Now())
		}
	}()
	defer ticker.Stop()

	// Demo mode (CEO-18A): zero-hardware virtual device for GitHub visitors.
	if *demo {
		if err := hub.startDemoDevice(); err != nil {
			slog.Error("demo device", "err", err)
			os.Exit(1)
		}
	}

	mux := hub.mux("")

	slog.Info("embedflow server listening", "addr", *addr, "data", dataDir, "version", version, "demo", *demo)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		slog.Error("http server", "err", err)
		os.Exit(1)
	}
}

// bootstrapAdmin creates the default admin account on first run
// (requirement 1.5: admins create accounts; there must be a first admin).
func bootstrapAdmin(st *store.Store) {
	if _, err := st.AuthenticateUser(context.Background(), "admin", "admin"); err == nil {
		return // exists
	}
	if err := st.CreateUser(context.Background(), "admin", "admin", "admin"); err == nil {
		slog.Warn("created default admin account (admin/admin) -- change the password")
	}
}

// mux builds the HTTP surface. frontDir serves the built web assets; empty
// means API-only (dev mode uses the Vite proxy).
func (h *coordinator) mux(frontDir string) http.Handler {
	mux := http.NewServeMux()

	// WebSocket endpoint: same-origin upgrade check (decision 3A) + token
	// auth inside the handshake.
	mux.Handle("/ws/client", authmw.SameOrigin(transport.Handler(h, &transport.AuthConfig{
		ValidateToken:   h.validateToken,
		ProtocolVersion: protocol.ProtocolVersion,
		ServerVersion:   version,
	})))

	// REST API: token-gated, except /api/login (the gate itself).
	// Write-config endpoints additionally require the admin role
	// (requirement 1.5: admins manage configuration, members execute).
	//
	// adminAPI holds the admin-only subtree; Go's ServeMux matches the more
	// specific pattern first, so a POST /api/projects lands in adminAPI even
	// though api also has a GET for the same path shape.
	admin := func(tok string) (string, string, bool) { return h.tokens.validate(tok) }
	api := http.NewServeMux()
	api.Handle("POST /api/projects", authmw.RequireAdmin(admin, http.HandlerFunc(h.handleCreateProject)))
	api.Handle("DELETE /api/projects/{id}", authmw.RequireAdmin(admin, http.HandlerFunc(h.handleDeleteProject)))
	api.Handle("POST /api/devices", authmw.RequireAdmin(admin, http.HandlerFunc(h.handleCreateDevice)))
	api.Handle("DELETE /api/devices/{id}", authmw.RequireAdmin(admin, http.HandlerFunc(h.handleDeleteDevice)))
	api.Handle("POST /api/expect-rules", authmw.RequireAdmin(admin, http.HandlerFunc(h.handleCreateRule)))
	api.Handle("PUT /api/expect-rules/{id}", authmw.RequireAdmin(admin, http.HandlerFunc(h.handleUpdateRule)))
	api.Handle("DELETE /api/expect-rules/{id}", authmw.RequireAdmin(admin, http.HandlerFunc(h.handleDeleteRule)))

	// Account management (requirement 1.5: admins create accounts; every
	// user may change their own password).
	api.Handle("GET /api/users", authmw.RequireAdmin(admin, http.HandlerFunc(h.handleListUsers)))
	api.Handle("POST /api/users", authmw.RequireAdmin(admin, http.HandlerFunc(h.handleCreateUser)))
	api.HandleFunc("POST /api/users/self/password", h.handleChangeOwnPassword)
	api.Handle("POST /api/users/{id}/password", authmw.RequireAdmin(admin, http.HandlerFunc(h.handleResetPassword)))

	// Build module (requirement 2). Writes are admin-only; triggering a
	// batch and reading records/downloads are member-executable.
	api.Handle("POST /api/projects/{pid}/build-items", authmw.RequireAdmin(admin, http.HandlerFunc(h.handleCreateBuildItem)))
	api.HandleFunc("GET /api/projects/{pid}/build-items", h.handleListBuildItems)
	api.Handle("PUT /api/build-items/{id}", authmw.RequireAdmin(admin, http.HandlerFunc(h.handleUpdateBuildItem)))
	api.Handle("DELETE /api/build-items/{id}", authmw.RequireAdmin(admin, http.HandlerFunc(h.handleDeleteBuildItem)))
	api.HandleFunc("POST /api/batches", h.handleCreateBatch)
	api.HandleFunc("GET /api/batches", h.handleListBatches)
	api.HandleFunc("GET /api/batches/{id}", h.handleGetBatch)
	api.HandleFunc("POST /api/batches/{id}/cancel", h.handleCancelBatch)
	api.HandleFunc("GET /api/build-records/{id}", h.handleGetBuildRecord)
	api.HandleFunc("GET /api/build-records/{id}/artifacts", h.handleRecordArtifacts)
	api.HandleFunc("GET /api/build-records/{id}/log/tail", h.handleBuildLogTail)
	api.HandleFunc("DELETE /api/build-records/{id}", h.handleDeleteBuildRecord)
	api.HandleFunc("POST /api/build-records/delete-all", h.handleDeleteBuildRecordsAll)
	api.HandleFunc("GET /api/artifacts/{id}/download", h.handleArtifactDownload)
	api.Handle("GET /api/settings", authmw.RequireAdmin(admin, http.HandlerFunc(h.handleGetSettings)))
	api.Handle("PUT /api/settings", authmw.RequireAdmin(admin, http.HandlerFunc(h.handlePutSettings)))

	api.HandleFunc("GET /api/me", h.handleMe)
	api.HandleFunc("GET /api/projects", h.handleListProjects)
	api.HandleFunc("GET /api/devices", h.handleListDevices)
	api.HandleFunc("GET /api/devices/{id}/sessions", h.handleDeviceSessions)
	api.HandleFunc("GET /api/expect-rules", h.handleListRules)
	api.HandleFunc("GET /api/expect-rules/{id}", h.handleGetRule)
	api.HandleFunc("GET /api/sessions/{id}", h.handleGetSession)
	api.HandleFunc("POST /api/sessions/{id}/close", h.handleCloseSession)
	api.HandleFunc("POST /api/sessions/{id}/confirm", h.handleInsertConfirm)
	api.HandleFunc("GET /api/sessions/{id}/confirmations", h.handleListConfirmations)
	api.HandleFunc("GET /api/sessions/{id}/log/tail", h.handleLogTail)
	api.Handle("GET /api/sessions/{id}/log/download", sessionlog.DownloadHandler(h.dataDir))

	mux.HandleFunc("POST /api/login", h.handleLogin)
	mux.Handle("/api/", authmw.Token("", func(tok string) bool {
		_, _, ok := h.tokens.validate(tok)
		return ok
	}, api))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
	// Embedded SPA (single-binary distribution); in dev the Vite proxy
	// serves the UI instead.
	if ui, err := webui.Handler(); err == nil {
		mux.Handle("/", ui)
	}
	return mux
}

// --- HTTP handlers (JSON over the REST surface) ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func pathID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

func (h *coordinator) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "bad body")
		return
	}
	token, role, err := h.login(req.Username, req.Password)
	if err != nil {
		writeErr(w, http.StatusUnauthorized, "bad credentials")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": token, "role": role})
}

func (h *coordinator) handleMe(w http.ResponseWriter, r *http.Request) {
	user, _, ok := h.tokens.validate(bearerToken(r))
	if !ok {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"user": user})
}

func bearerToken(r *http.Request) string {
	const p = "Bearer "
	h := r.Header.Get("Authorization")
	if len(h) > len(p) && h[:len(p)] == p {
		return h[len(p):]
	}
	return r.URL.Query().Get("token")
}

func (h *coordinator) handleListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.store.ListProjects(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if projects == nil {
		projects = []store.Project{}
	}
	writeJSON(w, http.StatusOK, projects)
}

func (h *coordinator) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		Note string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		writeErr(w, http.StatusBadRequest, "name required")
		return
	}
	id, err := h.store.CreateProject(r.Context(), req.Name, req.Note)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"id": id})
}

func (h *coordinator) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	if err := h.store.DeleteProject(r.Context(), id); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *coordinator) handleListDevices(w http.ResponseWriter, r *http.Request) {
	devs, err := h.store.ListDevices(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Enrich with live kernel state (idle vs occupied, DS-1A badge).
	type deviceView struct {
		store.Device
		Busy      bool   `json:"busy"`
		Owner     string `json:"owner,omitempty"`
		Kind      string `json:"kind,omitempty"`
		SessionID int64  `json:"session_id,omitempty"`
		Shared    bool   `json:"shared"`
		Since     string `json:"since,omitempty"`
	}
	out := make([]deviceView, 0, len(devs))
	for _, d := range devs {
		v := deviceView{Device: d}
		if st := h.kernel.DeviceState(session.DeviceID(d.ID)); st.SessionID != 0 {
			v.Busy = true
			v.Owner = st.Owner
			v.Kind = st.Kind
			v.SessionID = int64(st.SessionID)
			v.Since = st.Since.Format(time.RFC3339)
		}
		v.Shared = h.deviceShared(d.ID)
		out = append(out, v)
	}
	writeJSON(w, http.StatusOK, out)
}

// deviceShared reports whether a serial client is bound to the device.
func (h *coordinator) deviceShared(deviceID int64) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, st := range h.conns {
		if st.boundDevice == deviceID {
			return true
		}
	}
	return false
}

func (h *coordinator) handleCreateDevice(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		Project string `json:"project"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		writeErr(w, http.StatusBadRequest, "name required")
		return
	}
	id, err := h.store.CreateDevice(r.Context(), req.Name, req.Project)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"id": id})
}

func (h *coordinator) handleDeleteDevice(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	if err := h.store.DeleteDevice(r.Context(), id); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *coordinator) handleDeviceSessions(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	sess, err := h.store.SessionsForDevice(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sess)
}

func (h *coordinator) handleListRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.store.ListExpectRules(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rules)
}

func (h *coordinator) handleCreateRule(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string          `json:"name"`
		Steps json.RawMessage `json:"steps"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		writeErr(w, http.StatusBadRequest, "name and steps required")
		return
	}
	steps, err := expect.DecodeSteps(req.Steps)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad steps: "+err.Error())
		return
	}
	if errs := expect.Validate(steps); len(errs) > 0 {
		writeErr(w, http.StatusBadRequest, errs[0].Error())
		return
	}
	id, err := h.store.CreateExpectRule(r.Context(), req.Name, string(req.Steps))
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"id": id})
}

func (h *coordinator) handleGetRule(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	rule, err := h.store.GetExpectRule(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no such rule")
		return
	}
	writeJSON(w, http.StatusOK, rule)
}

func (h *coordinator) handleUpdateRule(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	var req struct {
		Name  string          `json:"name"`
		Steps json.RawMessage `json:"steps"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "bad body")
		return
	}
	steps, err := expect.DecodeSteps(req.Steps)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad steps: "+err.Error())
		return
	}
	if errs := expect.Validate(steps); len(errs) > 0 {
		writeErr(w, http.StatusBadRequest, errs[0].Error())
		return
	}
	if err := h.store.UpdateExpectRule(r.Context(), store.ExpectRule{ID: id, Name: req.Name, StepsJSON: string(req.Steps)}); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *coordinator) handleDeleteRule(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	if err := h.store.DeleteExpectRule(r.Context(), id); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *coordinator) handleGetSession(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	sess, err := h.store.GetSession(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no such session")
		return
	}
	writeJSON(w, http.StatusOK, sess)
}

func (h *coordinator) handleCloseSession(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	// Issue #6: owner or admin may close a session.
	user, role, _ := h.tokens.validate(bearerToken(r))
	sess, ok := h.kernel.Session(session.SessionID(id))
	if !ok {
		writeErr(w, http.StatusNotFound, "no such session")
		return
	}
	if role != "admin" && sess.Owner != user {
		writeErr(w, http.StatusForbidden, "not your session")
		return
	}
	for _, ev := range h.kernel.Close(session.SessionID(id), time.Now()) {
		h.apply(ev)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *coordinator) handleInsertConfirm(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	var req struct {
		Prompt string `json:"prompt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Prompt == "" {
		writeErr(w, http.StatusBadRequest, "prompt required")
		return
	}
	cid, err := h.insertConfirmation(id, req.Prompt)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"id": cid})
}

func (h *coordinator) handleListConfirmations(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	confs, err := h.store.ConfirmationsForSession(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, confs)
}

func (h *coordinator) handleLogTail(w http.ResponseWriter, r *http.Request) {
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
	tail, err := sessionlog.Tail(sessionlog.SessionFilePath(h.dataDir, id), n)
	if err != nil {
		writeErr(w, http.StatusNotFound, "no log")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"tail": tail})
}

// --- account management (requirement 1.5: admins open accounts) ---

func (h *coordinator) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.ListUsers(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if users == nil {
		users = []store.UserWithID{}
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *coordinator) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" || req.Password == "" {
		writeErr(w, http.StatusBadRequest, "username and password required")
		return
	}
	if req.Role != "admin" && req.Role != "member" {
		writeErr(w, http.StatusBadRequest, "role must be admin or member")
		return
	}
	if err := h.store.CreateUser(r.Context(), req.Username, req.Password, req.Role); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// meID resolves the authenticated caller's user row (for self password
// change); ok=false means the token was bad.
func (h *coordinator) meID(r *http.Request) (int64, bool) {
	user, _, ok := h.tokens.validate(bearerToken(r))
	if !ok {
		return 0, false
	}
	// tokens map to usernames; the id lookup goes through the users table.
	u, err := h.store.AuthenticateUserByName(r.Context(), user)
	if err != nil {
		return 0, false
	}
	return u.ID, true
}

func (h *coordinator) handleChangeOwnPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Old string `json:"old"`
		New string `json:"new"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.New == "" {
		writeErr(w, http.StatusBadRequest, "old and new required")
		return
	}
	id, ok := h.meID(r)
	if !ok {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := h.store.ChangeOwnPassword(r.Context(), id, req.Old, req.New); err != nil {
		if errors.Is(err, store.ErrBadCredentials) {
			writeErr(w, http.StatusForbidden, "wrong old password")
			return
		}
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *coordinator) handleResetPassword(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Password == "" {
		writeErr(w, http.StatusBadRequest, "password required")
		return
	}
	if _, err := h.store.GetUserByID(r.Context(), id); err != nil {
		writeErr(w, http.StatusNotFound, "no such user")
		return
	}
	if err := h.store.UpdatePassword(r.Context(), id, req.Password); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

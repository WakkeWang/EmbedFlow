// Command embedflow-server is the EmbedFlow server binary. M1 scope: the
// session layer end to end -- store + session kernel + WebSocket transport +
// session logs + expect runs + human confirmations + HTTP API, with the
// standard data layout (CEO-10A), startup sweep (CEO-7A) and tick loop.
package main

import (
	"context"
	"encoding/json"
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
)

var version = "dev"

func main() {
	var (
		addr = flag.String("addr", ":8420", "listen address")
		data = flag.String("data", "", "data directory (default ./data)")
		demo = flag.Bool("demo", false, "run with the built-in virtual device (CEO-18A)")
	)
	flag.Parse()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

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

	// REST API: token-gated.
	api := http.NewServeMux()
	api.HandleFunc("POST /api/login", h.handleLogin)
	api.HandleFunc("GET /api/me", h.handleMe)
	api.HandleFunc("GET /api/devices", h.handleListDevices)
	api.HandleFunc("POST /api/devices", h.handleCreateDevice)
	api.HandleFunc("DELETE /api/devices/{id}", h.handleDeleteDevice)
	api.HandleFunc("GET /api/devices/{id}/sessions", h.handleDeviceSessions)
	api.HandleFunc("GET /api/expect-rules", h.handleListRules)
	api.HandleFunc("POST /api/expect-rules", h.handleCreateRule)
	api.HandleFunc("GET /api/expect-rules/{id}", h.handleGetRule)
	api.HandleFunc("PUT /api/expect-rules/{id}", h.handleUpdateRule)
	api.HandleFunc("GET /api/sessions/{id}", h.handleGetSession)
	api.HandleFunc("POST /api/sessions/{id}/close", h.handleCloseSession)
	api.HandleFunc("POST /api/sessions/{id}/confirm", h.handleInsertConfirm)
	api.HandleFunc("GET /api/sessions/{id}/confirmations", h.handleListConfirmations)
	api.HandleFunc("GET /api/sessions/{id}/log/tail", h.handleLogTail)
	api.Handle("GET /api/sessions/{id}/log/download", sessionlog.DownloadHandler(h.dataDir))

	mux.Handle("/api/", authmw.Token("", func(tok string) bool {
		_, _, ok := h.tokens.validate(tok)
		return ok
	}, api))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
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

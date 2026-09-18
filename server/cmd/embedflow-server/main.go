// Command embedflow-server is the EmbedFlow server binary. T1+T2 scope: the
// session-layer skeleton -- store + session kernel + WebSocket transport +
// session logs, wired together with the startup sweep (CEO-7A), the standard
// data layout (CEO-10A) and a liveness tick loop.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/WakkeWang/EmbedFlow/pkg/protocol"
	"github.com/WakkeWang/EmbedFlow/server/internal/authmw"
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

	hub := &coordinator{
		kernel:  kernel,
		store:   st,
		dataDir: dataDir,
		conns:   map[*transport.ClientConn]struct{}{},
		loggers: map[int64]*sessionlog.Logger{},
	}

	// Startup sweep (CEO-7A): any session left "active" in the DB by a
	// previous run is failed with the restart reason, and the kernel is
	// seeded with history so IDs continue monotonically.
	if err := hub.sweepStartup(); err != nil {
		slog.Error("startup sweep", "err", err)
		os.Exit(1)
	}

	// Liveness tick loop: drives heartbeat-loss detection and manual idle
	// timeouts through the kernel, applying side effects to store/transport.
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			hub.tick(time.Now())
		}
	}()
	defer ticker.Stop()

	// M1 bootstrap token: gates the API (log downloads) until the real
	// login slice replaces the validator inside authmw.
	bootstrapToken := "bootstrap"

	mux := http.NewServeMux()
	mux.Handle("/ws/client", authmw.SameOrigin(transport.Handler(hub, &transport.AuthConfig{
		ValidateToken:   hub.validateToken,
		ProtocolVersion: protocol.ProtocolVersion,
		ServerVersion:   version,
	})))
	mux.Handle("/api/sessions/", authmw.Token(bootstrapToken, sessionlog.DownloadHandler(dataDir)))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	slog.Info("embedflow server listening", "addr", *addr, "data", dataDir, "version", version)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		slog.Error("http server", "err", err)
		os.Exit(1)
	}
}

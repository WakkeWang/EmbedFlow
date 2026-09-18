// Command serial-client is the EmbedFlow serial tunnel client (T4 minimal
// core): login, share a serial port to the server, pure passthrough. A
// Wails GUI wraps this core later; today the CLI proves the loop and
// serves as the developer harness.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/WakkeWang/EmbedFlow/serial-client"
)

func main() {
	var (
		server   = flag.String("server", "http://127.0.0.1:8420", "server URL")
		username = flag.String("user", "admin", "username")
		password = flag.String("pass", "", "password (required)")
		device   = flag.Int64("device", 0, "device id to share as (required)")
		port     = flag.String("port", "COM3", "local COM port name")
	)
	flag.Parse()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	if *password == "" || *device == 0 {
		fmt.Fprintln(os.Stderr, "usage: serial-client -server URL -user U -pass P -device N [-port COM3]")
		os.Exit(2)
	}

	ctx := signalContext()

	tok, err := client.Login(*server, *username, *password)
	if err != nil {
		slog.Error("login failed", "err", err)
		os.Exit(1)
	}
	slog.Info("logged in")

	cl := client.New(*server, tok)
	if err := cl.Connect(ctx); err != nil {
		slog.Error("connect failed", "err", err)
		os.Exit(1)
	}
	defer cl.Close()

	// Real COM port wiring arrives with the Wails build (go.bug.st/serial);
	// without it the client still shares and reports state.
	if err := cl.Share(*device, *port); err != nil {
		slog.Error("share failed", "err", err)
		os.Exit(1)
	}
	slog.Info("sharing", "device", *device, "port", *port)

	<-ctx.Done()
	slog.Info("bye")
}

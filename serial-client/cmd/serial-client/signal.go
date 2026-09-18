package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// signalContext returns a context cancelled on Ctrl+C / SIGTERM.
func signalContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		cancel()
	}()
	return ctx
}

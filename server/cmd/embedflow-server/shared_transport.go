package main

import (
	"errors"
	"sync"
	"time"

	"github.com/WakkeWang/EmbedFlow/server/internal/transport"
)

// sharedTransport bridges the expect engine and the serial client's
// connection: engine Reads drain device bytes (RX pushed from the client),
// engine Writes go straight back to the client (TX). It satisfies
// serialport.Port so the engine needs no changes for the live path.
type sharedTransport struct {
	sid    int64
	client *transport.ClientConn
	hub    *coordinator

	mu         sync.Mutex
	buf        []byte // device bytes awaiting the engine
	notify     chan struct{}
	notifyOnce sync.Once

	closedCh  chan struct{}
	closeOnce sync.Once
}

func newSharedTransport(sid int64, client *transport.ClientConn, hub *coordinator) *sharedTransport {
	return &sharedTransport{
		sid:      sid,
		client:   client,
		hub:      hub,
		notify:   make(chan struct{}),
		closedCh: make(chan struct{}),
	}
}

// push appends device bytes (called from the coordinator when the serial
// client sends binary frames).
func (s *sharedTransport) push(data []byte) {
	s.mu.Lock()
	s.buf = append(s.buf, data...)
	s.mu.Unlock()
	s.notifyOnce.Do(func() { close(s.notify) })
}

// Read drains accumulated device bytes; bounded-wait so aborts stay live.
func (s *sharedTransport) Read(p []byte) (int, error) {
	for {
		s.mu.Lock()
		if len(s.buf) > 0 {
			n := copy(p, s.buf)
			s.buf = s.buf[n:]
			s.mu.Unlock()
			return n, nil
		}
		s.mu.Unlock()
		select {
		case <-s.notify:
			// re-check; a fresh notify channel arms after each drain burst
			s.mu.Lock()
			if len(s.buf) == 0 {
				s.notify = make(chan struct{})
				s.notifyOnce = sync.Once{}
			}
			s.mu.Unlock()
		case <-s.closedCh:
			return 0, errTransportClosed
		case <-time.After(100 * time.Millisecond):
			return 0, errTransportTimeout
		}
	}
}

// Write sends TX bytes to the serial client and logs them.
func (s *sharedTransport) Write(p []byte) (int, error) {
	_ = s.hub.writeLog(s.sid, "TX", p)
	_ = s.client.SendBinary(append([]byte(nil), p...))
	return len(p), nil
}

// Close tears the bridge down.
func (s *sharedTransport) Close() error {
	s.closeOnce.Do(func() { close(s.closedCh) })
	return nil
}

// SetReadDeadline is accepted for the Port interface; Read's fixed 100ms
// wait keeps abort responsiveness without deadline plumbing.
func (s *sharedTransport) SetReadDeadline(t time.Time) error { return nil }

type transportError string

func (e transportError) Error() string { return string(e) }

const (
	errTransportTimeout = transportError("sharedtransport: read timeout")
	errTransportClosed  = transportError("sharedtransport: closed")
)

var _ = errors.Is

package client

import (
	"fmt"
	"time"

	"go.bug.st/serial"
)

// realPort adapts a go.bug.st/serial port to the PortOpener interface.
// Byte-level pure passthrough (requirement 3.4): no line buffering, no
// business logic.
type realPort struct {
	p serial.Port
}

// OpenPort opens a COM port with the fixed settings the flash flow assumes
// (8N1). Baud varies per device; the caller passes what the operator chose.
func OpenPort(name string, baud int) (PortOpener, error) {
	mode := &serial.Mode{BaudRate: baud, DataBits: 8, Parity: serial.NoParity, StopBits: serial.OneStopBit}
	p, err := serial.Open(name, mode)
	if err != nil {
		return nil, fmt.Errorf("client: open %s: %w", name, err)
	}
	// A serial link never blocks indefinitely like a pipe can: a modest read
	// timeout keeps the passthrough pump responsive to shutdown and to
	// server-side aborts without burning CPU.
	if err := p.SetReadTimeout(100 * time.Millisecond); err != nil {
		p.Close()
		return nil, fmt.Errorf("client: %s read timeout: %w", name, err)
	}
	return &realPort{p: p}, nil
}

func (r *realPort) Read(p []byte) (int, error)  { return r.p.Read(p) }
func (r *realPort) Write(p []byte) (int, error) { return r.p.Write(p) }
func (r *realPort) Close() error                { return r.p.Close() }

// ListPorts returns the system's serial port names (requirement 3.4 scan).
func ListPorts() ([]string, error) {
	return serial.GetPortsList()
}

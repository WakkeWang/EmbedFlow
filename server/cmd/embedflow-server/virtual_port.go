package main

import (
	"errors"
	"time"
)

// virtualPort connects the expect engine directly to the virtual device
// (demo mode): engine Reads drain device output, engine Writes deliver
// operator bytes into the device. Satisfies serialport.Port.
type virtualPort struct {
	hub *coordinator
	dev *virtualDevice
	sid int64
	ch  chan []byte
}

func newVirtualPort(hub *coordinator, deviceID, sessionID int64) *virtualPort {
	return &virtualPort{
		hub: hub,
		dev: hub.virtual,
		sid: sessionID,
		ch:  make(chan []byte, 64),
	}
}

// Read blocks for device output pushed by the device's emit path.
func (p *virtualPort) Read(b []byte) (int, error) {
	select {
	case data := <-p.ch:
		n := copy(b, data)
		return n, nil
	case <-time.After(100 * time.Millisecond):
		return 0, errTransportTimeout
	}
}

// Write delivers operator bytes into the device. Logging is the caller's
// job: the expect engine routes TX through Config.OnSendBytes (so secret
// sends land masked, issue #14), and manual keystrokes log in sendToDevice.
func (p *virtualPort) Write(b []byte) (int, error) {
	p.dev.deliver(b)
	return len(b), nil
}

func (p *virtualPort) Close() error { return nil }

func (p *virtualPort) SetReadDeadline(t time.Time) error { return nil }

var _ = errors.New

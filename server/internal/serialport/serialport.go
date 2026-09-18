// Package serialport abstracts the serial transport behind a small interface
// (decision 7A): the CI suite runs against in-memory pipes, com0com stays a
// local-manual-test concern, and the real COM driver plugs in for T4.
package serialport

import (
	"errors"
	"io"
	"sync"
	"time"
)

// Port is the minimal surface the session layer needs from a serial device.
type Port interface {
	io.ReadWriteCloser
	// SetReadDeadline bounds a pending/future Read so liveness and shutdown
	// logic never blocks forever. Timeouts return an error; zero clears.
	SetReadDeadline(t time.Time) error
}

// NewPipe returns a connected (client, device) pair: bytes written to one
// end arrive at the other, both binary safe. It stands in for a real serial
// link in CI.
func NewPipe() (client Port, device Port) {
	c, d := newPipePair()
	return c, d
}

// pipeEnd is one half of an in-memory pipe: its Read drains bytes the other
// end Wrote. Waiters get a fresh closed-channel broadcast on every state
// change (data arrived, close, deadline change).
type pipeEnd struct {
	mu        sync.Mutex
	buf       []byte
	peer      *pipeEnd
	closed    bool
	deadline  time.Time
	waiters   []chan struct{}
	waitersMu sync.Mutex
}

func (p *pipeEnd) Write(b []byte) (int, error) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return 0, errors.New("serialport: write on closed end")
	}
	peer := p.peer
	peer.mu.Lock()
	peer.buf = append(peer.buf, b...)
	peer.mu.Unlock()
	p.mu.Unlock()
	peer.broadcast()
	return len(b), nil
}

func (p *pipeEnd) Read(b []byte) (int, error) {
	for {
		p.mu.Lock()
		if len(p.buf) > 0 {
			n := copy(b, p.buf)
			p.buf = p.buf[n:]
			p.mu.Unlock()
			return n, nil
		}
		if p.closed {
			p.mu.Unlock()
			return 0, errors.New("serialport: read on closed end")
		}
		deadline := p.deadline
		p.mu.Unlock()

		// Wait for the next broadcast (data/close/deadline) or the deadline.
		wake := make(chan struct{})
		p.waitersMu.Lock()
		p.waiters = append(p.waiters, wake)
		p.waitersMu.Unlock()

		var timer <-chan time.Time
		if !deadline.IsZero() {
			d := time.Until(deadline)
			if d <= 0 {
				p.removeWaiter(wake)
				return 0, errors.New("serialport: read deadline exceeded")
			}
			t := time.NewTimer(d)
			timer = t.C
			defer t.Stop()
		}

		select {
		case <-wake:
			// Loop re-checks buf/closed; nothing to do here.
		case <-timer:
			p.removeWaiter(wake)
			return 0, errors.New("serialport: read deadline exceeded")
		}
	}
}

func (p *pipeEnd) removeWaiter(ch chan struct{}) {
	p.waitersMu.Lock()
	for i, w := range p.waiters {
		if w == ch {
			p.waiters = append(p.waiters[:i], p.waiters[i+1:]...)
			break
		}
	}
	p.waitersMu.Unlock()
}

func (p *pipeEnd) Close() error {
	p.mu.Lock()
	alreadyClosed := p.closed
	peer := p.peer
	p.closed = true
	if peer != nil {
		// One end closing tears down the whole link: a real serial line
		// unplugged on one side reads as an error on the other.
		peer.closed = true
	}
	p.mu.Unlock()
	if alreadyClosed {
		return nil
	}
	p.broadcast()
	if peer != nil {
		peer.broadcast()
	}
	return nil
}

func (p *pipeEnd) SetReadDeadline(t time.Time) error {
	p.mu.Lock()
	p.deadline = t
	p.mu.Unlock()
	p.broadcast()
	return nil
}

// broadcast closes all wait channels; the next waiter re-arms fresh.
func (p *pipeEnd) broadcast() {
	p.waitersMu.Lock()
	ws := p.waiters
	p.waiters = nil
	p.waitersMu.Unlock()
	for _, w := range ws {
		close(w)
	}
}

func newPipePair() (*pipeEnd, *pipeEnd) {
	a := &pipeEnd{}
	b := &pipeEnd{}
	a.peer = b
	b.peer = a
	return a, b
}

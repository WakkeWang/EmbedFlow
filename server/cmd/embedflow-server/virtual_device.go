package main

import (
	"context"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/WakkeWang/EmbedFlow/server/internal/expect"
	"github.com/WakkeWang/EmbedFlow/server/internal/transport"
)

// virtualDevice is the built-in demo device (CEO-18A): zero hardware, so a
// GitHub visitor can run the full loop -- edit an expect rule, run it, watch
// the terminal, get a session log, insert a human confirmation. It speaks
// uboot-ish on an in-memory link: what the engine sends arrives in inbuf,
// what the device "emits" goes out through the hub's normal RX fan-out.
type virtualDevice struct {
	deviceID int64

	mu     sync.Mutex
	inbuf  []byte
	closed bool
}

// startDemoDevice registers the virtual device (idempotently: exactly one
// per data dir) and seeds a demo expect rule.
func (h *coordinator) startDemoDevice() error {
	// Reuse an existing demo device across restarts.
	devs, err := h.store.ListDevices(context.Background())
	if err != nil {
		return err
	}
	var devID int64
	for _, d := range devs {
		if d.Name == "demo-virtual" {
			devID = d.ID
			break
		}
	}
	if devID == 0 {
		devID, err = h.store.CreateDevice(context.Background(), "demo-virtual", "demo")
		if err != nil {
			return err
		}
	}
	vd := &virtualDevice{deviceID: devID}
	// Idempotent demo rule exercising every engine feature: await, send,
	// delay (silent install), and a long-running sequence.
	steps := []expect.Step{
		{Await: "Press ENTER", Timeout: 10 * time.Second},
		{Send: `\r`},
		{Await: "demo login:", Timeout: 10 * time.Second},
		{Send: `demo\r`},
		{Await: "Password:", Timeout: 10 * time.Second},
		{Send: `demo\r`},
		{Await: `demo#`, Timeout: 10 * time.Second},
		{Send: `install\r`},
		{Delay: 1500 * time.Millisecond},
		{Await: `install over`, Timeout: 15 * time.Second},
		{Send: `reboot\r`},
		{Await: `demo-virtual ready`, Timeout: 10 * time.Second},
	}
	raw, err := expect.EncodeSteps(steps)
	if err != nil {
		return err
	}
	// Seed the demo rule once per data dir.
	rules, _ := h.store.ListExpectRules(context.Background())
	for _, r := range rules {
		if r.Name == "demo-flash" {
			h.mu.Lock()
			h.demoRuleID = r.ID
			h.mu.Unlock()
			go h.attachVirtual(vd, devID)
			return nil
		}
	}
	ruleID, err := h.store.CreateExpectRule(context.Background(), "demo-flash", string(raw))
	if err != nil {
		return err
	}
	h.mu.Lock()
	h.demoRuleID = ruleID
	h.mu.Unlock()

	h.attachVirtual(vd, devID)
	return nil
}

// attachVirtual wires the device and starts its script loop.
func (h *coordinator) attachVirtual(vd *virtualDevice, devID int64) {
	h.mu.Lock()
	h.virtual = vd
	h.mu.Unlock()
	go vd.run(h)
	slog.Info("demo virtual device ready", "device", devID)
}

// run plays the device script forever: emit prompts, await operator bytes.
// Each pass waits for an operator session (virtual port) before starting,
// so the banner reaches the engine rather than firing into the void.
func (v *virtualDevice) run(h *coordinator) {
	for !v.isClosed() {
		if !v.awaitOperator(h, 24*time.Hour) {
			return
		}
		v.emit(h, "U-Boot 2026.01-demo\r\ndemo-virtual board\r\nPress ENTER to activate this console.")
		if _, ok := v.awaitInput("\r", 60*time.Second); !ok {
			continue
		}
		v.emit(h, "\r\ndemo login: ")
		if _, ok := v.awaitInput("demo\r", 60*time.Second); !ok {
			continue
		}
		v.emit(h, "Password: ")
		if _, ok := v.awaitInput("demo\r", 60*time.Second); !ok {
			continue
		}
		v.emit(h, "\r\ndemo# ")
		if _, ok := v.awaitInput("install\r", 60*time.Second); !ok {
			continue
		}
		// Silent install phase: the CEO-1A case (no global idle timeout on
		// task sessions, so 1.5s of silence is safe).
		time.Sleep(1200 * time.Millisecond)
		v.emit(h, "writing partition 1... ok\r\n")
		time.Sleep(300 * time.Millisecond)
		v.emit(h, "writing partition 2... ok\r\n")
		time.Sleep(300 * time.Millisecond)
		v.emit(h, "install over\r\n")
		if _, ok := v.awaitInput("reboot\r", 60*time.Second); !ok {
			continue
		}
		v.emit(h, "rebooting...\r\n")
		time.Sleep(500 * time.Millisecond)
		v.emit(h, "demo-virtual ready\r\n")
		time.Sleep(2 * time.Second)
	}
}

// awaitOperator blocks until a virtual port registers for this device
// (an expect run started) or the timeout elapses. Returns false on shutdown.
func (v *virtualDevice) awaitOperator(h *coordinator, timeout time.Duration) bool {
	deadline := time.After(timeout)
	tick := time.NewTicker(20 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-deadline:
			return false
		case <-tick.C:
			h.mu.Lock()
			n := len(h.virtualPorts)
			h.mu.Unlock()
			if n > 0 {
				return true
			}
		}
	}
}

// emit publishes device output: log RX + deliver into the live virtual port
// (expect engine) + fan out to subscribed browsers.
func (v *virtualDevice) emit(h *coordinator, s string) {
	data := []byte(s)
	h.mu.Lock()
	sid := h.virtualSession[v.deviceID]
	var ports []*virtualPort
	for _, p := range h.virtualPorts {
		if p.sid == sid {
			ports = append(ports, p)
		}
	}
	targets := make([]*transport.ClientConn, 0, 2)
	for cc, st := range h.conns {
		if st.sessionID == sid && st.boundDevice == 0 {
			targets = append(targets, cc)
		}
	}
	h.mu.Unlock()
	_ = h.writeLog(sid, "RX", data)
	for _, p := range ports {
		select {
		case p.ch <- data:
		default:
		}
	}
	for _, cc := range targets {
		_ = cc.SendBinary(data)
	}
}

// deliver feeds operator TX bytes into the device (called by sendToDevice
// when the target device is virtual).
func (v *virtualDevice) deliver(data []byte) {
	v.mu.Lock()
	v.inbuf = append(v.inbuf, data...)
	v.mu.Unlock()
}

// awaitInput waits for the operator/engine to send a line containing want.
func (v *virtualDevice) awaitInput(want string, timeout time.Duration) ([]byte, bool) {
	deadline := time.After(timeout)
	tick := time.NewTicker(20 * time.Millisecond)
	defer tick.Stop()
	var acc []byte
	for {
		select {
		case <-deadline:
			return nil, false
		case <-tick.C:
			acc = append(acc, v.takeInput()...)
			if strings.Contains(string(acc), want) {
				return acc, true
			}
		}
	}
}

func (v *virtualDevice) takeInput() []byte {
	v.mu.Lock()
	defer v.mu.Unlock()
	b := v.inbuf
	v.inbuf = nil
	return b
}

func (v *virtualDevice) isClosed() bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.closed
}

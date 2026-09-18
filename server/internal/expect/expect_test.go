package expect

import (
	"testing"
	"time"

	"github.com/WakkeWang/EmbedFlow/server/internal/serialport"
)

// The engine is tested against in-memory pipes (decision 7A): the "device"
// goroutine writes scripted output with controlled timing, the engine runs
// on the client end. Table-driven per the T3 verify line.

// --- matching ---

func TestEngine_Contains_Match(t *testing.T) {
	runner, dev := newFixture(t, []Step{
		{Await: "Login:"},
	})
	dev.script([]scriptLine{
		{after: 10 * time.Millisecond, data: "Welcome\r\nLogin:"},
	})
	res := runner.run()
	res.wantOK(t)
}

func TestEngine_Regex_Match(t *testing.T) {
	runner, dev := newFixture(t, []Step{
		{Await: `IP: \d+\.\d+`, Match: MatchRegex},
	})
	dev.script([]scriptLine{
		{after: 5 * time.Millisecond, data: "IP: 192.168.0.80"},
	})
	res := runner.run()
	res.wantOK(t)
}

func TestEngine_Exact_Match(t *testing.T) {
	runner, dev := newFixture(t, []Step{
		{Await: "OK", Match: MatchExact},
	})
	dev.script([]scriptLine{
		{after: 5 * time.Millisecond, data: "OK"},
	})
	res := runner.run()
	res.wantOK(t)
}

func TestEngine_Exact_RejectsSubstring(t *testing.T) {
	runner, dev := newFixture(t, []Step{
		{Await: "OK", Match: MatchExact},
	})
	dev.script([]scriptLine{
		{after: 5 * time.Millisecond, data: "NOK"},
	})
	res := runner.run()
	res.wantFailed(t, 0)
}

// --- send between awaits ---

func TestEngine_SendBetweenAwaits(t *testing.T) {
	runner, dev := newFixture(t, []Step{
		{Await: "Login:"},
		{Send: "root\r"},
		{Await: "Password:"},
		{Send: "secret\r"},
		{Await: "#"},
	})
	dev.script([]scriptLine{
		{after: 5 * time.Millisecond, data: "Login:"},
		{thenExpect: "root\r", data: "Password:"},
		{thenExpect: "secret\r", data: "D2000#"},
	})
	res := runner.run()
	res.wantOK(t)
}

// --- timeout & failure actions (VBS three-element model) ---

func TestEngine_Timeout_FailActionAbort(t *testing.T) {
	runner, dev := newFixture(t, []Step{
		{Await: "never-comes", Timeout: 200 * time.Millisecond, OnFail: FailAbort},
		{Await: "also-never"},
	})
	dev.script(nil)
	res := runner.run()
	res.wantFailed(t, 0)
	res.wantStoppedAfter(t, 0) // abort: later steps never run
}

func TestEngine_Timeout_FailActionContinue(t *testing.T) {
	runner, dev := newFixture(t, []Step{
		{Await: "never-comes", Timeout: 100 * time.Millisecond, OnFail: FailContinue},
		{Await: "second-step"},
	})
	dev.script([]scriptLine{
		{after: 150 * time.Millisecond, data: "second-step"},
	})
	res := runner.run()
	res.wantOK(t) // continue: step 0 failed but execution went on
}

func TestEngine_Timeout_FailActionRetry(t *testing.T) {
	runner, dev := newFixture(t, []Step{
		{Await: "late", Timeout: 100 * time.Millisecond, OnFail: FailRetry, MaxRetries: 2},
	})
	dev.script([]scriptLine{
		{after: 150 * time.Millisecond, data: "late"}, // first timeout, then arrives
	})
	res := runner.run()
	res.wantOK(t)
}

func TestEngine_Retry_Exhausted_Fails(t *testing.T) {
	runner, dev := newFixture(t, []Step{
		{Await: "never", Timeout: 50 * time.Millisecond, OnFail: FailRetry, MaxRetries: 2},
	})
	dev.script(nil)
	res := runner.run()
	res.wantFailed(t, 0)
}

// --- delay steps (CEO-13A: blind-send / wait windows) ---

func TestEngine_DelayStep(t *testing.T) {
	runner, dev := newFixture(t, []Step{
		{Send: "reboot\r"},
		{Delay: 150 * time.Millisecond},
		{Await: "Loading"},
	})
	dev.script([]scriptLine{
		{thenExpect: "reboot\r", data: ""},
		{after: 100 * time.Millisecond, data: "Loading Environment..."},
	})
	res := runner.run()
	res.wantOK(t)
}

// --- abort (requirement 3.6: cancel = stop remaining steps) ---

func TestEngine_Abort_StopsSequence(t *testing.T) {
	runner, dev := newFixture(t, []Step{
		{Await: "step-one"},
		{Await: "step-two", Timeout: 5 * time.Second},
	})
	dev.script([]scriptLine{
		{after: 5 * time.Millisecond, data: "step-one"},
	})
	go func() {
		// Wait for the engine to be mid-sequence, then abort.
		for runner.eng == nil {
			time.Sleep(time.Millisecond)
		}
		time.Sleep(20 * time.Millisecond)
		runner.eng.Abort()
	}()
	res := runner.run()
	res.wantAborted(t)
}

// --- send encoding (CEO-13A: \r, Ctrl-C 0x03, HEX escape) ---

func TestEncode_CR(t *testing.T) {
	got, err := EncodeSend("hello\\r")
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	want := []byte("hello\r")
	if string(got) != string(want) {
		t.Fatalf("got % x, want % x", got, want)
	}
}

func TestEncode_CtrlC(t *testing.T) {
	got, err := EncodeSend("\\C")
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if len(got) != 1 || got[0] != 0x03 {
		t.Fatalf("got % x, want [0x03]", got)
	}
}

func TestEncode_HexEscape(t *testing.T) {
	got, err := EncodeSend("\\x1b[0m")
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	want := []byte{0x1b, '[', '0', 'm'}
	if string(got) != string(want) {
		t.Fatalf("got % x, want % x", got, want)
	}
}

func TestEncode_InvalidHex(t *testing.T) {
	if _, err := EncodeSend("\\xZZ"); err == nil {
		t.Fatal("invalid hex should error")
	}
}

func TestEncode_PlainTextPassthrough(t *testing.T) {
	got, err := EncodeSend("plain text 123")
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if string(got) != "plain text 123" {
		t.Fatalf("got %q", got)
	}
}

// --- validation ---

func TestValidate_EmptyAwaitAndNoSend(t *testing.T) {
	errs := Validate([]Step{
		{}, // neither await nor delay nor send: meaningless
	})
	if len(errs) != 1 {
		t.Fatalf("errs = %v, want 1", errs)
	}
}

func TestValidate_NegativeTimeout(t *testing.T) {
	errs := Validate([]Step{
		{Await: "x", Timeout: -time.Second},
	})
	if len(errs) != 1 {
		t.Fatalf("errs = %v, want 1", errs)
	}
}

func TestValidate_ValidSequence(t *testing.T) {
	errs := Validate([]Step{
		{Await: "Login:"},
		{Send: "root\\r"},
		{Delay: time.Second},
	})
	if len(errs) != 0 {
		t.Fatalf("errs = %v, want 0", errs)
	}
}

// --- fixtures ---

type scriptLine struct {
	after      time.Duration // delay from the previous line
	data       string        // bytes the "device" emits
	thenExpect string        // wait for these client bytes before emitting data
}

type fakeDevice struct {
	t    *testing.T
	port serialport.Port
}

func (d *fakeDevice) script(lines []scriptLine) {
	go func() {
		for _, ln := range lines {
			if ln.thenExpect != "" {
				d.waitInput(ln.thenExpect)
			}
			if ln.after > 0 {
				time.Sleep(ln.after)
			}
			if ln.data != "" {
				d.port.Write([]byte(ln.data))
			}
		}
	}()
}

func (d *fakeDevice) waitInput(want string) {
	buf := make([]byte, 256)
	deadline := time.Now().Add(2 * time.Second)
	d.port.SetReadDeadline(deadline)
	var acc []byte
	for {
		n, err := d.port.Read(buf)
		if err != nil {
			return // timeout/close: script proceeds
		}
		acc = append(acc, buf[:n]...)
		if len(acc) >= len(want) && string(acc[len(acc)-len(want):]) == want {
			return
		}
	}
}

func newFixture(t *testing.T, steps []Step) (*fixtureRunner, *fakeDevice) {
	t.Helper()
	client, device := serialport.NewPipe()
	cfg := Config{
		Steps:         steps,
		DefaultWait:   5 * time.Second,
		DefaultOnFail: FailAbort,
	}
	return &fixtureRunner{t: t, cfg: cfg, port: client}, &fakeDevice{t: t, port: device}
}

type fixtureRunner struct {
	t    *testing.T
	cfg  Config
	port serialport.Port
	res  *Result
	eng  *Runner
}

func (r *fixtureRunner) run() *Result {
	r.eng = New(r.cfg, r.port)
	if r.res == nil {
		r.res = r.eng.Run()
	}
	return r.res
}

func (r *Result) wantOK(t *testing.T) {
	t.Helper()
	if !r.OK() {
		t.Fatalf("want completed sequence, got %+v", r)
	}
}

func (r *Result) wantFailed(t *testing.T, stepIdx int) {
	t.Helper()
	if r.OK() {
		t.Fatalf("want failure at step %d, got %+v", stepIdx, r)
	}
	if r.StepIndex != stepIdx {
		t.Fatalf("failed step = %d, want %d (%s)", r.StepIndex, stepIdx, r.Detail)
	}
}

func (r *Result) wantStoppedAfter(t *testing.T, stepIdx int) {
	t.Helper()
	if r.LastStep() != stepIdx {
		t.Fatalf("last executed step = %d, want %d", r.LastStep(), stepIdx)
	}
}

func (r *Result) wantAborted(t *testing.T) {
	t.Helper()
	if !r.Aborted {
		t.Fatalf("want aborted, got %+v", r)
	}
}

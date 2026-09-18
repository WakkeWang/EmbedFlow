// Package expect implements the expect engine: a sequence of steps --
// await (wait for a string), optional send, timeout, failure action --
// executed against a serial port (requirement 3.1.3, the VBS three-element
// model). M1 ships the await mode only; verdict (PASS/FAIL semantics)
// lands with M4 (decision CEO-14A).
package expect

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/WakkeWang/EmbedFlow/server/internal/serialport"
)

// Match is how the await string is compared against device output.
type Match string

const (
	MatchContains Match = "contains" // default (VBS semantics)
	MatchRegex    Match = "regex"
	MatchExact    Match = "exact" // whole accumulated output equals the string
)

// OnFail is the configured action when a step times out.
type OnFail string

const (
	FailAbort    OnFail = "abort"    // stop the sequence (default)
	FailContinue OnFail = "continue" // log and move to the next step
	FailRetry    OnFail = "retry"    // retry up to MaxRetries, then abort
)

// Step is one element of an expect sequence.
type Step struct {
	// Await is the string to wait for. Empty means a non-await step (send
	// only or delay only).
	Await string `json:"await,omitempty"`
	// Match selects the comparison mode for Await.
	Match Match `json:"match,omitempty"`
	// Send is the encoded send content (control chars via EncodeSend).
	Send string `json:"send,omitempty"`
	// SendBeforeWait, when true, sends before waiting for Await (blind send
	// pattern; otherwise Send fires after the previous step matched).
	SendBeforeWait bool `json:"send_before_wait,omitempty"`
	// Timeout bounds the await; zero uses the config default.
	Timeout time.Duration `json:"timeout,omitempty"`
	// OnFail selects the failure action; zero uses the config default.
	OnFail OnFail `json:"on_fail,omitempty"`
	// MaxRetries bounds FailRetry; zero means one attempt only.
	MaxRetries int `json:"max_retries,omitempty"`
	// Delay makes this a delay step: wait this long, matching nothing.
	Delay time.Duration `json:"delay,omitempty"`
	// Secret marks the send content as a password (issue #14): the device
	// still receives the real bytes, but the TX log line is fully masked.
	Secret bool `json:"secret,omitempty"`
}

// Config drives one Run.
type Config struct {
	Steps         []Step
	DefaultWait   time.Duration
	DefaultOnFail OnFail
	// OnStep, when set, fires as each step starts: (index, total). Live
	// progress for the terminal UI (issue #8, DS-3A).
	OnStep func(index, total int)
	// OnSendBytes, when set, fires for every write the engine makes:
	// (bytes, secret). The executor logs TX through it so secret sends land
	// masked in the session log while the device still gets real bytes
	// (issue #14).
	OnSendBytes func(data []byte, secret bool)
}

// Result reports how the sequence ended.
type Result struct {
	// StepIndex is the step the sequence stopped on (0-based): the failed
	// step, the aborted-at step, or the last step on success.
	StepIndex int
	// Detail explains the end: matched / timeout / abort / error.
	Detail string
	// Aborted is true when Abort() was called (requirement 3.6 second
	// confirmation lands in the UI; the engine just reports).
	Aborted bool
}

// OK reports whether the whole sequence completed.
func (r *Result) OK() bool { return !r.Aborted && strings.HasPrefix(r.Detail, "completed") }

// LastStep returns the highest step index that started executing.
func (r *Result) LastStep() int { return r.StepIndex }

// Runner executes a sequence and supports abort from another goroutine
// (requirement 3.6: the UI's abort button lands here in T5).
type Runner struct {
	cfg   Config
	port  serialport.Port
	abort chan struct{}

	doneMu sync.Mutex
	done   bool

	outMu sync.Mutex
	out   []byte // device output since the current await began

	readQuit chan struct{} // closed by stopRead to end the read loop
	readDone chan struct{} // closed by the read loop on exit
}

// New creates a runner for the given sequence.
func New(cfg Config, port serialport.Port) *Runner {
	return &Runner{cfg: cfg, port: port, abort: make(chan struct{})}
}

// Run executes the sequence synchronously: it returns when the sequence
// completes, fails, or is aborted. One Runner runs at most one sequence.
// On return the internal read loop is stopped: a second Runner on the same
// port never races a zombie reader.
func (r *Runner) Run() *Result {
	res := r.exec()
	r.stopRead()
	r.doneMu.Lock()
	r.done = true
	r.doneMu.Unlock()
	return res
}

// Abort stops the sequence at the next check.
func (r *Runner) Abort() { close(r.abort) }

// Done reports whether the sequence has finished running (result returned).
// Input gating uses it: while a run is live, keystrokes must not interleave
// with automatic sends (CEO-17A).
func (r *Runner) Done() bool {
	r.doneMu.Lock()
	defer r.doneMu.Unlock()
	return r.done
}

func (r *Runner) exec() *Result {
	readErr := make(chan error, 1)
	r.readQuit = make(chan struct{})
	r.readDone = make(chan struct{})
	go r.readLoop(readErr)

	for i := range r.cfg.Steps {
		step := &r.cfg.Steps[i]
		select {
		case <-r.abort:
			return &Result{StepIndex: i, Detail: "aborted", Aborted: true}
		default:
		}
		if r.cfg.OnStep != nil {
			r.cfg.OnStep(i, len(r.cfg.Steps))
		}

		// Pure delay step (CEO-13A: wait windows).
		if step.Delay > 0 {
			select {
			case <-time.After(step.Delay):
			case <-r.abort:
				return &Result{StepIndex: i, Detail: "aborted", Aborted: true}
			}
			continue
		}

		// Send (after previous match by default; before wait when flagged).
		if step.Send != "" && step.SendBeforeWait {
			if res := r.send(i, step); res != nil {
				return res
			}
		}

		// Await step: wait for the string (send-only steps without Await
		// have nothing to match and fall through).
		if step.Await != "" {
			if res := r.await(i, step, readErr); res != nil {
				return res
			}
		}

		// Send (after wait).
		if step.Send != "" && !step.SendBeforeWait {
			if res := r.send(i, step); res != nil {
				return res
			}
		}
	}
	return &Result{StepIndex: len(r.cfg.Steps) - 1, Detail: "completed"}
}

// onSendBytes notifies the executor of an engine write (TX log hook).
func (r *Runner) onSendBytes(step *Step, data []byte) {
	if r.cfg.OnSendBytes != nil {
		r.cfg.OnSendBytes(data, step.Secret)
	}
}

// send writes the step's decoded send content to the port. The TX log hook
// fires only after a successful write: a log line claims the device
// received the bytes, so a failed write must not produce one.
func (r *Runner) send(stepIdx int, step *Step) *Result {
	data, err := EncodeSend(step.Send)
	if err != nil {
		return &Result{StepIndex: stepIdx, Detail: fmt.Sprintf("send encode: %v", err)}
	}
	if _, err := r.port.Write(data); err != nil {
		return &Result{StepIndex: stepIdx, Detail: fmt.Sprintf("send: %v", err)}
	}
	r.onSendBytes(step, data)
	return nil
}

// await waits for step.Await, applying timeout + failure action. A nil
// return means the step succeeded (or the failure was tolerated).
func (r *Runner) await(stepIdx int, step *Step, readErr <-chan error) *Result {
	timeout := step.Timeout
	if timeout <= 0 {
		timeout = r.cfg.DefaultWait
	}
	onFail := step.OnFail
	if onFail == "" {
		onFail = r.cfg.DefaultOnFail
		if onFail == "" {
			onFail = FailAbort
		}
	}

	attempts := 1
	if onFail == FailRetry && step.MaxRetries > 0 {
		attempts = step.MaxRetries
	}

	m := step.Match
	if m == "" {
		m = MatchContains
	}

	// Compile the regex once for the whole step (was per 5ms poll tick).
	var re *regexp.Regexp
	if m == MatchRegex {
		var err error
		re, err = regexp.Compile(step.Await)
		if err != nil {
			return &Result{StepIndex: stepIdx, Detail: fmt.Sprintf("step %d bad regex: %v", stepIdx, err)}
		}
	}

	for attempt := 0; attempt < attempts; attempt++ {
		// A retry is a fresh attempt: re-send (VBS re-runs the whole step)
		// and discard stale output so old bytes can never satisfy attempt 2.
		if attempt > 0 {
			r.resetOut()
			if step.Send != "" {
				if res := r.send(stepIdx, step); res != nil {
					return res
				}
			}
		}
		ok, detail := r.waitMatch(step.Await, m, re, timeout, readErr)
		if ok {
			return nil
		}
		if detail == "aborted" {
			return &Result{StepIndex: stepIdx, Detail: "aborted", Aborted: true}
		}
		switch onFail {
		case FailContinue:
			return nil // failure tolerated: log detail and move on
		case FailRetry:
			if attempt == attempts-1 {
				return &Result{StepIndex: stepIdx, Detail: fmt.Sprintf("step %d timeout after %d attempts: %s", stepIdx, attempts, detail)}
			}
		default: // FailAbort
			return &Result{StepIndex: stepIdx, Detail: fmt.Sprintf("step %d timeout: %s", stepIdx, detail)}
		}
	}
	// Unreachable: every branch above returns.
	return &Result{StepIndex: stepIdx, Detail: "exhausted"}
}

// waitMatch blocks until Await matches accumulated output, the deadline
// passes, the port errors, or the run is aborted.
func (r *Runner) waitMatch(await string, m Match, re *regexp.Regexp, timeout time.Duration, readErr <-chan error) (bool, string) {
	deadline := time.After(timeout)
	tick := time.NewTicker(5 * time.Millisecond)
	defer tick.Stop()

	if r.match(await, m, re) {
		return true, ""
	}
	for {
		select {
		case <-r.abort:
			return false, "aborted"
		case err := <-readErr:
			return false, fmt.Sprintf("port read: %v", err)
		case <-deadline:
			return false, fmt.Sprintf("timeout after %s", timeout)
		case <-tick.C:
			if r.match(await, m, re) {
				return true, ""
			}
		}
	}
}

// match tests accumulated output; on success it clears the buffer so the
// next await starts from fresh output (VBS semantics: each wait sees new
// data only).
func (r *Runner) match(await string, m Match, re *regexp.Regexp) bool {
	r.outMu.Lock()
	defer r.outMu.Unlock()
	if !matchOutput(await, m, re, r.out) {
		return false
	}
	r.out = nil
	return true
}

func (r *Runner) resetOut() {
	r.outMu.Lock()
	r.out = nil
	r.outMu.Unlock()
}

func (r *Runner) appendOut(data []byte) {
	r.outMu.Lock()
	r.out = append(r.out, data...)
	r.outMu.Unlock()
}

func (r *Runner) readLoop(errCh chan<- error) {
	defer close(r.readDone)
	buf := make([]byte, 4096)
	for {
		// Bound each Read so the stop signal is honored promptly even when
		// no data ever arrives (the pipe honors deadlines; real COM too).
		r.port.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		n, err := r.port.Read(buf)
		if n > 0 {
			r.appendOut(buf[:n])
		}
		select {
		case <-r.readQuit:
			return
		default:
		}
		if err != nil {
			// Bounded-read timeouts are normal pacing (pipes, the shared
			// transport and the virtual port all return timeout errors);
			// anything else really ends the loop.
			msg := err.Error()
			if strings.Contains(msg, "deadline") || strings.Contains(msg, "timeout") {
				continue
			}
			errCh <- err
			return
		}
	}
}

// stopRead ends the read loop and waits for it, so the port is exclusively
// ours while running and unclaimed after Run returns.
func (r *Runner) stopRead() {
	if r.readQuit == nil {
		return
	}
	close(r.readQuit)
	select {
	case <-r.readDone:
	case <-time.After(time.Second):
		// Port Read overruns its deadline by a bounded amount at worst.
	}
	r.readQuit = nil
}

func matchOutput(await string, m Match, re *regexp.Regexp, out []byte) bool {
	switch m {
	case MatchRegex:
		return re.Match(out)
	case MatchExact:
		return string(out) == await
	default:
		return strings.Contains(string(out), await)
	}
}

// EncodeSend expands the send-content escapes (CEO-13A): \r -> CR,
// \n -> LF, \t -> TAB, \C -> Ctrl-C (0x03), \xHH -> raw byte, everything
// else passes through literally.
func EncodeSend(s string) ([]byte, error) {
	var out []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '\\' {
			out = append(out, c)
			continue
		}
		i++
		if i >= len(s) {
			return nil, errors.New("expect: dangling escape in send content")
		}
		switch s[i] {
		case 'r':
			out = append(out, '\r')
		case 'n':
			out = append(out, '\n')
		case 't':
			out = append(out, '\t')
		case 'C':
			out = append(out, 0x03)
		case 'x':
			if i+2 >= len(s) {
				return nil, errors.New("expect: \\x escape needs two hex digits")
			}
			v, err := strconv.ParseUint(s[i+1:i+3], 16, 8)
			if err != nil {
				return nil, fmt.Errorf("expect: bad hex escape %q: %w", s[i:i+3], err)
			}
			out = append(out, byte(v))
			i += 2
		case '\\':
			out = append(out, '\\')
		default:
			return nil, fmt.Errorf("expect: unknown escape \\%c", s[i])
		}
	}
	return out, nil
}

// DecodeSteps parses stored rule JSON into engine steps (issue #7: rules
// persist as JSON).
func DecodeSteps(data []byte) ([]Step, error) {
	var steps []Step
	if err := json.Unmarshal(data, &steps); err != nil {
		return nil, fmt.Errorf("expect: decode steps: %w", err)
	}
	return steps, nil
}

// EncodeSteps serializes steps for storage.
func EncodeSteps(steps []Step) ([]byte, error) {
	return json.Marshal(steps)
}

// Validate checks a sequence before it is stored or run: each step must do
// something (await, send, or delay), timeouts are non-negative, regexes
// compile, escapes decode.
func Validate(steps []Step) []error {
	var errs []error
	for i := range steps {
		s := &steps[i]
		if s.Await == "" && s.Send == "" && s.Delay == 0 {
			errs = append(errs, fmt.Errorf("step %d: empty (no await, send, or delay)", i))
		}
		if s.Timeout < 0 {
			errs = append(errs, fmt.Errorf("step %d: negative timeout", i))
		}
		if s.Delay < 0 {
			errs = append(errs, fmt.Errorf("step %d: negative delay", i))
		}
		if s.Await != "" && s.Match == MatchRegex {
			if _, err := regexp.Compile(s.Await); err != nil {
				errs = append(errs, fmt.Errorf("step %d: bad regex: %w", i, err))
			}
		}
		if s.Send != "" {
			if _, err := EncodeSend(s.Send); err != nil {
				errs = append(errs, fmt.Errorf("step %d: bad send content: %w", i, err))
			}
		}
	}
	return errs
}

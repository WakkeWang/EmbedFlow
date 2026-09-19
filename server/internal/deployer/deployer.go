// Package deployer executes deploy rules (requirement 3.1). Three forms:
// manual (a step description for the operator -- no execution), ssh
// (rendered command templates over one SSH connection, decision 1A's short
// transaction), and flash (an expect sequence through the M1 engine -- the
// orchestrator wires that to a task session, this package only renders).
//
// Logs are phase-tagged like the builder's: "<ts> [phase] line" with phase
// in render|ssh|system, streamed live through Notify.
package deployer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// Phases tag log lines (the web UI groups by these).
const (
	PhaseRender = "render"
	PhaseSSH    = "ssh"
	PhaseSystem = "system"
)

// Notify is the live-push hook: (deployRecordID, phase, line).
type Notify func(deployRecordID int64, phase, line string)

// Variables carries the three variable families (requirement 3.1.2):
// artifact.* from the build record, device.* from the device row,
// params.* from the operator's form.
type Variables struct {
	Artifact map[string]string
	Device   map[string]string
	Params   map[string]string
}

// Render expands ${family.name} references. Unknown references render as
// empty and are reported (the command template is configuration; a typo
// must be visible in the log, not silently swallowed).
func Render(tpl string, v Variables) (string, []string) {
	var missing []string
	out := varPattern.ReplaceAllStringFunc(tpl, func(m string) string {
		name := strings.TrimSuffix(strings.TrimPrefix(m, "${"), "}")
		switch {
		case strings.HasPrefix(name, "artifact."):
			if val, ok := v.Artifact[strings.TrimPrefix(name, "artifact.")]; ok {
				return val
			}
		case strings.HasPrefix(name, "device."):
			if val, ok := v.Device[strings.TrimPrefix(name, "device.")]; ok {
				return val
			}
		case strings.HasPrefix(name, "params."):
			if val, ok := v.Params[strings.TrimPrefix(name, "params.")]; ok {
				return val
			}
		}
		missing = append(missing, name)
		return ""
	})
	return out, missing
}

var varPattern = regexp.MustCompile(`\$\{[a-zA-Z][a-zA-Z0-9_.]*\}`)

// ExecResult reports how one SSH command batch ended.
type ExecResult struct {
	OK       bool
	Detail   string
	ExitCode int
	// Canceled is true when the run ended because of a cancel request: the
	// caller must not overwrite the record's canceled state.
	Canceled bool
}

// SSHRunner executes one deploy rule's command list against a device over
// SSH. Each Deploy() call is one short transaction: dial, run every
// command in order (stop at the first failure), close.
type SSHRunner struct {
	Host    string
	Port    int
	User    string
	Password string
	// Timeout bounds the whole batch (0 = no limit).
	Timeout time.Duration
	Notify  Notify
	Record  int64

	mu      sync.Mutex
	session *ssh.Session
	client  *ssh.Client
	canceled bool
}

// Cancel aborts the running batch: the SSH session and connection close,
// the running command dies with the connection (requirement 3.6: cancel =
// close the session + mark failed).
func (r *SSHRunner) Cancel() {
	r.mu.Lock()
	r.canceled = true
	sess, client := r.session, r.client
	r.mu.Unlock()
	if sess != nil {
		_ = sess.Close()
	}
	if client != nil {
		_ = client.Close()
	}
}

func (r *SSHRunner) wasCanceled() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.canceled
}

// Run executes the commands in order. Each command's combined output is
// logged line by line; a non-zero exit fails the batch (requirement 3.1.2:
// the deploy result and record must give feedback).
func (r *SSHRunner) Run(ctx context.Context, commands []string) ExecResult {
	if r.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.Timeout)
		defer cancel()
	}

	cfg := &ssh.ClientConfig{
		User:            r.User,
		Auth:            []ssh.AuthMethod{ssh.Password(r.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // internal trusted network (requirement 6.2 note)
		Timeout:         10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", r.Host, r.Port)
	r.log(PhaseSSH, "dialing "+addr)
	client, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		return r.fail("ssh dial: " + err.Error())
	}
	r.mu.Lock()
	r.client = client
	r.mu.Unlock()
	defer client.Close()

	for i, command := range commands {
		if r.wasCanceled() {
			return ExecResult{Canceled: true, Detail: "canceled"}
		}
		if ctx.Err() != nil {
			return r.fail(fmt.Sprintf("timeout before command %d", i+1))
		}
		r.log(PhaseSSH, fmt.Sprintf("$ %s", command))

		sess, err := client.NewSession()
		if err != nil {
			if r.wasCanceled() {
				return ExecResult{Canceled: true, Detail: "canceled"}
			}
			return r.fail("ssh session: " + err.Error())
		}
		r.mu.Lock()
		r.session = sess
		r.mu.Unlock()

		out, err := sess.CombinedOutput(command)
		r.mu.Lock()
		r.session = nil
		r.mu.Unlock()
		sess.Close()

		for _, l := range strings.Split(string(out), "\n") {
			if l = strings.TrimRight(l, "\r"); l != "" {
				r.log(PhaseSSH, l)
			}
		}
		if err != nil {
			if r.wasCanceled() {
				return ExecResult{Canceled: true, Detail: "canceled"}
			}
			code := 255
			var ee *ssh.ExitError
			if errors.As(err, &ee) {
				code = ee.ExitStatus()
			}
			return r.fail(fmt.Sprintf("command %d failed (exit %d): %s", i+1, code, command))
		}
	}
	r.log(PhaseSSH, fmt.Sprintf("%d command(s) completed", len(commands)))
	return ExecResult{OK: true, Detail: fmt.Sprintf("%d command(s) completed", len(commands))}
}

func (r *SSHRunner) fail(detail string) ExecResult {
	r.log(PhaseSystem, detail)
	return ExecResult{Detail: detail, ExitCode: 1}
}

func (r *SSHRunner) log(phase, line string) {
	if r.Notify != nil {
		r.Notify(r.Record, phase, line)
	}
}

// WriteExecLog appends phase-tagged lines to <data>/execs/<id>/log.txt.
// Buffered like the builder's log (small files, short lives -- a plain
// file append per line is enough at deploy scale).
func WriteExecLog(dataDir string, execID int64, phase, line string) error {
	dir := filepath.Join(dataDir, "execs", fmt.Sprint(execID))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(dir, "log.txt"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s [%s] %s\n", time.Now().Format("2006-01-02 15:04:05.000"), phase, line)
	return err
}

// ArtifactView is the deploy trigger's view of one build artifact (name and
// checksum ride the UI; size formats for display).
type ArtifactView struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Checksum string `json:"checksum"`
}

// ArtifactVariables maps artifacts into template variables (requirement
// 3.1.2): artifact.name = the first artifact's name, artifact.url = the
// download path, artifact.checksum = its checksum. With several artifacts
// the JSON list rides artifact.list.
func ArtifactVariables(arts []ArtifactView) map[string]string {
	m := map[string]string{}
	if len(arts) == 0 {
		return m
	}
	first := arts[0]
	m["name"] = first.Name
	m["checksum"] = first.Checksum
	m["url"] = fmt.Sprintf("/api/artifacts/%d/download", first.ID)
	if raw, err := json.Marshal(arts); err == nil {
		m["list"] = string(raw)
	}
	return m
}

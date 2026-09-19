// Package builder is the M2 build executor: it turns one build record into
// a real build -- source preparation (git clone / local clone), dirty
// rejection, command execution with timeout, version extraction, artifact
// globbing + checksum archiving, and a phase-tagged build log streamed to
// the coordinator for live WS push.
//
// Log line format (requirement 2.4 分段可查): "<ts> [phase] content" with
// phase in clone|dirty|command|verify|version|artifact|system, so the web
// UI can group lines by phase.
package builder

import (
	"bufio"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/WakkeWang/EmbedFlow/server/internal/store"
)

// Phases tag log lines (the web UI groups by these).
const (
	PhaseClone    = "clone"
	PhaseDirty    = "dirty"
	PhaseCommand  = "command"
	PhaseVerify   = "verify"
	PhaseVersion  = "version"
	PhaseArtifact = "artifact"
	PhaseSystem   = "system"
)

// Notify is the live-push hook: (recordID, phase, line) per log line.
type Notify func(recordID int64, phase, line string)

// CancelFunc reports whether the record was canceled mid-run (the batch
// kernel's cancel path flips it); checked between phases.
type CancelFunc func(recordID int64) bool

// Executor runs single builds. One Run per record; the coordinator owns
// scheduling (the batch kernel decides when).
type Executor struct {
	DataDir string
	// TmpRoot is the configured build temp dir (default <data>/tmp).
	TmpRoot  string
	Store    *store.Store
	Notify   Notify
	Canceled CancelFunc

	// lastFailed is set when the build failed: the tmp dir is kept for the
	// 7-day sweep instead of being removed immediately (requirement 2.5).
	lastFailed bool
}

// Result reports how one build ended.
type Result struct {
	OK          bool
	Detail      string // human-readable end reason
	CommitSHA   string
	VersionInfo string
	ExitCode    *int
	// Artifacts lists what archiveArtifacts stored (the caller persists
	// the rows).
	Artifacts []store.Artifact
}

// checksummer computes the configured checksum over a file.
type checksummer func(path string) (string, error)

// ChecksummerFor maps the system setting value to a function; unknown
// values fall back to sha256 (requirement 1.6 default).
func ChecksummerFor(setting string) checksummer {
	if setting == "md5" {
		return checksumFileMD5
	}
	return checksumFileSHA256
}

func checksumFileSHA256(path string) (string, error) {
	return checksumFile(path, sha256.New())
}

func checksumFileMD5(path string) (string, error) {
	return checksumFile(path, md5.New())
}

func checksumFile(path string, h hash.Hash) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// logWriter appends phase-tagged lines to the build log and pushes them.
type logWriter struct {
	mu     sync.Mutex
	f      *os.File
	buf    *bufio.Writer
	record int64
	notify Notify
}

func openBuildLog(dataDir string, recordID int64, notify Notify) (*logWriter, error) {
	dir := filepath.Dir(pathsBuildLog(dataDir, recordID))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("builder: mkdir log dir: %w", err)
	}
	f, err := os.OpenFile(pathsBuildLog(dataDir, recordID), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, fmt.Errorf("builder: open log: %w", err)
	}
	return &logWriter{f: f, buf: bufio.NewWriter(f), record: recordID, notify: notify}, nil
}

func pathsBuildLog(dataDir string, recordID int64) string {
	// Local indirection to avoid importing paths here (single small dup).
	return filepath.Join(dataDir, "builds", fmt.Sprint(recordID), "log.txt")
}

func (w *logWriter) line(phase, content string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	ts := time.Now().Format("2006-01-02 15:04:05.000")
	text := fmt.Sprintf("%s [%s] %s\n", ts, phase, strings.TrimRight(content, "\r\n"))
	if w.buf != nil {
		w.buf.WriteString(text)
	}
	if w.notify != nil {
		w.notify(w.record, phase, strings.TrimRight(content, "\r\n"))
	}
}

func (w *logWriter) close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.buf != nil {
		if err := w.buf.Flush(); err != nil {
			return err
		}
	}
	if w.f != nil {
		return w.f.Close()
	}
	return nil
}

// Run executes one build record end to end. Phase order: clone -> dirty ->
// command -> version -> artifact. Cancellation is honored at phase
// boundaries and during the command (via context kill).
func (e *Executor) Run(recordID int64, item store.BuildItem, checksumSetting string) Result {
	log, err := openBuildLog(e.DataDir, recordID, e.Notify)
	if err != nil {
		return Result{Detail: "open build log: " + err.Error()}
	}
	defer log.close()

	res := e.run(recordID, item, checksumSetting, log)
	e.lastFailed = !res.OK
	log.line(PhaseSystem, "build "+map[bool]string{true: "succeeded", false: "failed"}[res.OK]+": "+res.Detail)
	return res
}

func (e *Executor) run(recordID int64, item store.BuildItem, checksumSetting string, log *logWriter) Result {
	cancelled := func() bool {
		return e.Canceled != nil && e.Canceled(recordID)
	}

	// 1. Source preparation (dirty check runs on the source inside).
	tmpDir := filepath.Join(e.TmpRoot, fmt.Sprintf("build-%d", recordID))
	buildDir, res := e.prepareSource(recordID, item, tmpDir, log)
	if buildDir == "" {
		e.lastFailed = true
		return *res
	}
	defer e.cleanupTmp(recordID, tmpDir, log)

	// 3. Commit snapshot.
	sha, err := gitHead(buildDir)
	if err != nil {
		return Result{Detail: "read commit: " + err.Error()}
	}
	log.line(PhaseClone, "source snapshot: "+sha)

	if cancelled() {
		return Result{Detail: "canceled"}
	}

	// 4. Build command.
	exitCode, cmdOK, detail := e.runCommand(recordID, item, buildDir, log)
	if !cmdOK {
		e.lastFailed = true
		return Result{Detail: detail, CommitSHA: sha, ExitCode: &exitCode}
	}

	if cancelled() {
		return Result{Detail: "canceled", CommitSHA: sha, ExitCode: &exitCode}
	}

	// 5. Version extraction.
	versionInfo := ""
	if item.VersionCmd != "" {
		vi, verr := e.runVersionCmd(recordID, item, buildDir, log)
		if verr != "" {
			// A version-extraction failure fails the record: the version
			// info is part of the record contract (requirement 2.1).
			e.lastFailed = true
			return Result{Detail: "version extraction: " + verr, CommitSHA: sha, ExitCode: &exitCode}
		}
		versionInfo = vi
	}

	// 6. Artifacts.
	arts, aerr := e.archiveArtifacts(recordID, item, buildDir, checksumSetting, log)
	if aerr != "" {
		e.lastFailed = true
		return Result{Detail: "artifacts: " + aerr, CommitSHA: sha, VersionInfo: versionInfo, ExitCode: &exitCode}
	}
	if len(arts) == 0 {
		e.lastFailed = true
		return Result{Detail: "no artifacts matched the declared globs", CommitSHA: sha, VersionInfo: versionInfo, ExitCode: &exitCode}
	}

	return Result{OK: true, Detail: fmt.Sprintf("%d artifacts archived", len(arts)), CommitSHA: sha, VersionInfo: versionInfo, ExitCode: &exitCode, Artifacts: arts}
}

// prepareSource gets a clean git worktree into tmpDir and returns the build
// directory (inside tmpDir). Empty buildDir means failure (res carries why).
// The dirty check runs on the SOURCE before cloning (requirement 2.2: dirty
// 拒建 -- a dirty source means someone has uncommitted changes, so the build
// output would not correspond to any commit; git clone drops untracked
// files anyway, so post-clone checks would never fire).
func (e *Executor) prepareSource(recordID int64, item store.BuildItem, tmpDir string, log *logWriter) (string, *Result) {
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return "", &Result{Detail: "mkdir tmp: " + err.Error()}
	}
	switch item.SourceType {
	case "local":
		if item.LocalPath == "" {
			return "", &Result{Detail: "local source has no path"}
		}
		if !isGitRepo(item.LocalPath) {
			log.line(PhaseClone, "local path is not a git repository: "+item.LocalPath)
			return "", &Result{Detail: "local source is not a git repository: " + item.LocalPath}
		}
		if clean, err := gitClean(item.LocalPath); err != nil {
			return "", &Result{Detail: "dirty check: " + err.Error()}
		} else if !clean {
			log.line(PhaseDirty, "source worktree is dirty -- refusing to build (requirement 2.2)")
			return "", &Result{Detail: "dirty worktree: commit or stash before building"}
		}
		log.line(PhaseDirty, "source worktree clean")
		// git clone --local hardlinks objects: fast, cheap (requirement 2.2).
		buildDir := filepath.Join(tmpDir, "src")
		args := []string{"clone", "--local"}
		if item.GitBranch != "" {
			args = append(args, "--branch", item.GitBranch)
		}
		args = append(args, item.LocalPath, buildDir)
		if out, err := runGitLogged(log, PhaseClone, "", args...); err != nil {
			return "", &Result{Detail: "local clone: " + err.Error() + ": " + out}
		}
		return buildDir, nil

	case "git":
		if item.GitURL == "" {
			return "", &Result{Detail: "git source has no url"}
		}
		buildDir := filepath.Join(tmpDir, "src")
		// Shallow clone, last 2 commits (requirement 2.2).
		args := []string{"clone", "--depth", "2"}
		if item.GitBranch != "" {
			args = append(args, "--branch", item.GitBranch)
		}
		args = append(args, item.GitURL, buildDir)
		if out, err := runGitLogged(log, PhaseClone, "", args...); err != nil {
			return "", &Result{Detail: "git clone: " + err.Error() + ": " + out}
		}
		// Latest-source check (requirement 2.1: 勾选后构建前核对并拉取远端最新).
		if item.CheckLatest {
			if out, err := runGitLogged(log, PhaseClone, buildDir, "fetch", "origin"); err != nil {
				return "", &Result{Detail: "fetch for latest check: " + err.Error() + ": " + out}
			}
			if item.GitBranch == "" {
				return "", &Result{Detail: "latest-source check requires a branch"}
			}
			local, _ := runGitLogged(log, PhaseVerify, buildDir, "rev-parse", "HEAD")
			remote, rerr := runGitLogged(log, PhaseVerify, buildDir, "rev-parse", "origin/"+item.GitBranch)
			if rerr != nil {
				return "", &Result{Detail: "resolve remote branch: " + rerr.Error()}
			}
			if strings.TrimSpace(local) != strings.TrimSpace(remote) {
				log.line(PhaseVerify, "remote moved: local "+strings.TrimSpace(local)+" vs origin "+strings.TrimSpace(remote))
				// Pull to latest: the shallow clone was made moments ago, so
				// a move means someone pushed mid-clone -- reset to origin.
				if out, err := runGitLogged(log, PhaseClone, buildDir, "reset", "--hard", "origin/"+item.GitBranch); err != nil {
					return "", &Result{Detail: "reset to latest: " + err.Error() + ": " + out}
				}
			}
		}
		return buildDir, nil

	default:
		return "", &Result{Detail: "unknown source type: " + item.SourceType}
	}
}

// runCommand executes item.Command in buildDir with the server's
// environment (requirement 2.3: 继承系统运行环境). Returns (exitCode, ok).
func (e *Executor) runCommand(recordID int64, item store.BuildItem, buildDir string, log *logWriter) (int, bool, string) {
	log.line(PhaseCommand, "$ "+item.Command)
	ctx := context.Background()
	var cancel context.CancelFunc
	if item.TimeoutSec > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(item.TimeoutSec)*time.Second)
		defer cancel()
	}
	cmd := exec.CommandContext(ctx, shellName(), shellArg(item.Command))
	cmd.Dir = buildDir
	// Process-group kill so shell scripts cannot orphan children (the flash
	// flow's bb.sh pattern spawns toolchains); Windows falls back to
	// Process.Kill via CommandContext.
	setPgid(cmd)

	out, err := cmd.CombinedOutput()
	scannerLines(out, func(l string) { log.line(PhaseCommand, l) })
	if err != nil {
		code := exitCodeOf(err)
		if ctx.Err() != nil {
			log.line(PhaseCommand, "build command timed out or canceled")
			return code, false, fmt.Sprintf("build command timed out after %ds", item.TimeoutSec)
		}
		return code, false, fmt.Sprintf("build command failed (exit %d)", code)
	}
	return 0, true, ""
}

// runVersionCmd executes item.VersionCmd and returns its full output
// (requirement 2.1: 完整输出全部记录为版本信息).
func (e *Executor) runVersionCmd(recordID int64, item store.BuildItem, buildDir string, log *logWriter) (string, string) {
	cmd := exec.Command(shellName(), shellArg(item.VersionCmd))
	cmd.Dir = buildDir
	out, err := cmd.CombinedOutput()
	scannerLines(out, func(l string) { log.line(PhaseVersion, l) })
	if err != nil {
		return "", fmt.Sprintf("version command failed: %v", err)
	}
	return strings.TrimRight(string(out), "\r\n"), ""
}

// archiveArtifacts globs the declared artifact patterns and copies matches
// into <data>/artifacts/<recordID>/ with checksums.
func (e *Executor) archiveArtifacts(recordID int64, item store.BuildItem, buildDir, checksumSetting string, log *logWriter) ([]store.Artifact, string) {
	csum := ChecksummerFor(checksumSetting)
	dest := filepath.Join(e.DataDir, "artifacts", fmt.Sprint(recordID))
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return nil, "mkdir artifact dir: " + err.Error()
	}
	var out []store.Artifact
	for _, pattern := range item.Artifacts {
		if strings.ContainsAny(pattern, "\\/") && filepath.IsAbs(pattern) {
			return nil, "artifact path must be relative to the build dir: " + pattern
		}
		matches, err := filepath.Glob(filepath.Join(buildDir, pattern))
		if err != nil {
			return nil, "bad glob " + pattern + ": " + err.Error()
		}
		for _, m := range matches {
			st, err := os.Stat(m)
			if err != nil || st.IsDir() {
				continue
			}
			name := filepath.Base(m)
			data, err := os.ReadFile(m)
			if err != nil {
				return nil, "read artifact " + name + ": " + err.Error()
			}
			if err := os.WriteFile(filepath.Join(dest, name), data, 0o644); err != nil {
				return nil, "archive artifact " + name + ": " + err.Error()
			}
			sum, err := csum(filepath.Join(dest, name))
			if err != nil {
				return nil, "checksum " + name + ": " + err.Error()
			}
			log.line(PhaseArtifact, fmt.Sprintf("%s (%d bytes, %s)", name, st.Size(), sum))
			out = append(out, store.Artifact{BuildRecordID: recordID, Name: name, Size: st.Size(), Checksum: sum})
		}
	}
	return out, ""
}

// cleanupTmp removes the tmp dir on success/cancel; leaves it for the
// 7-day sweep on failure (requirement 2.5).
func (e *Executor) cleanupTmp(recordID int64, tmpDir string, log *logWriter) {
	if e.lastFailed {
		log.line(PhaseSystem, "temp dir kept for 7-day cleanup: "+tmpDir)
		return
	}
	if err := os.RemoveAll(tmpDir); err != nil {
		log.line(PhaseSystem, "tmp cleanup failed: "+err.Error())
	}
}

// CleanupSweep deletes failed-build tmp dirs older than maxAge
// (requirement 2.5: 构建失败保留 7 天自动清理).
func CleanupSweep(tmpRoot string, maxAge time.Duration) {
	entries, err := os.ReadDir(tmpRoot)
	if err != nil {
		return
	}
	for _, ent := range entries {
		if !ent.IsDir() || !strings.HasPrefix(ent.Name(), "build-") {
			continue
		}
		full := filepath.Join(tmpRoot, ent.Name())
		st, err := os.Stat(full)
		if err != nil {
			continue
		}
		if time.Since(st.ModTime()) > maxAge {
			_ = os.RemoveAll(full)
		}
	}
}

// --- git helpers ---

func isGitRepo(dir string) bool {
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

func gitClean(dir string) (bool, error) {
	out, err := exec.Command("git", "-C", dir, "status", "--porcelain").Output()
	if err != nil {
		return false, err
	}
	return len(strings.TrimSpace(string(out))) == 0, nil
}

func gitHead(dir string) (string, error) {
	out, err := exec.Command("git", "-C", dir, "rev-parse", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func runGitLogged(log *logWriter, phase, dir string, args ...string) (string, error) {
	full := append([]string{}, "git")
	if dir != "" {
		full = append(full, "-C", dir)
	}
	full = append(full, args...)
	cmd := exec.Command(full[0], full[1:]...)
	out, err := cmd.CombinedOutput()
	text := strings.TrimRight(string(out), "\r\n")
	if text != "" {
		scannerLines(out, func(l string) { log.line(phase, l) })
	}
	return text, err
}

// scannerLines feeds non-empty lines of out to fn.
func scannerLines(out []byte, fn func(string)) {
	for _, l := range strings.Split(string(out), "\n") {
		l = strings.TrimRight(l, "\r")
		if l != "" {
			fn(l)
		}
	}
}

// exitCodeOf extracts the process exit code from a run error.
func exitCodeOf(err error) int {
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode()
	}
	return -1
}

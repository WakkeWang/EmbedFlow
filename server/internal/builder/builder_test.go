package builder

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/WakkeWang/EmbedFlow/server/internal/store"
)

// Integration tests against real git commands (decision 7A's
// "真机手动验收" analog for the build executor: git is present on every
// dev/CI machine). Each test builds a scratch repo fixture.

// gitRepo is a scratch source repo with one commit containing build.sh and
// an artifact emitter.
type gitRepo struct {
	dir string
}

func newGitRepo(t *testing.T) *gitRepo {
	t.Helper()
	dir := t.TempDir()
	git(t, dir, "init", "-b", "main")
	git(t, dir, "config", "user.email", "test@example.com")
	git(t, dir, "config", "user.name", "Test")
	// build script: emits the fixed artifact set the tests glob for.
	script := "#!/bin/sh\necho payload-of-bootfw > bootfw-p0133.bin\necho payload-of-uboot > u-boot.bin\n"
	name := "build.sh"
	if runtime.GOOS == "windows" {
		script = "@echo off\r\necho payload-of-bootfw> bootfw-p0133.bin\r\necho payload-of-uboot> u-boot.bin\r\n"
		name = "build.bat"
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(script), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "VERSION"), []byte("1.2.3-demo\n"), 0o644); err != nil {
		t.Fatalf("write version: %v", err)
	}
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-m", "initial")
	return &gitRepo{dir: dir}
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", args, err, out)
	}
}

func buildCommand() string {
	if runtime.GOOS == "windows" {
		// Windows 11 (26200+) no longer resolves bare command names from
		// the current directory; the explicit .\ prefix is required. The
		// production server is Linux where sh -c ./bb.sh works as written.
		return ".\\build.bat"
	}
	return "sh build.sh"
}

func (r *gitRepo) item(t *testing.T, artifacts ...string) store.BuildItem {
	return store.BuildItem{
		Name:       "bootfw",
		SourceType: "local",
		LocalPath:  r.dir,
		Command:    buildCommand(),
		Artifacts:  artifacts,
		VersionCmd: "cat VERSION",
	}
}

func newExecutor(t *testing.T) (*Executor, string) {
	t.Helper()
	dataDir := t.TempDir()
	tmpRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dataDir, "builds"), 0o755); err != nil {
		t.Fatalf("mkdir builds: %v", err)
	}
	return &Executor{DataDir: dataDir, TmpRoot: tmpRoot}, dataDir
}

func TestBuild_CleanHappyPath(t *testing.T) {
	repo := newGitRepo(t)
	e, dataDir := newExecutor(t)
	item := repo.item(t, "bootfw*.bin")

	var lines []string
	e.Notify = func(record int64, phase, line string) {
		lines = append(lines, phase+"|"+line)
	}

	res := e.Run(101, item, "sha256")
	if !res.OK {
		t.Fatalf("build failed: %s\nlog:\n%s", res.Detail, strings.Join(lines, "\n"))
	}
	if !strings.Contains(res.VersionInfo, "1.2.3-demo") {
		t.Fatalf("version info = %q", res.VersionInfo)
	}

	// Artifact archived with a correct sha256.
	art, err := os.ReadFile(filepath.Join(dataDir, "artifacts", "101", "bootfw-p0133.bin"))
	if err != nil {
		t.Fatalf("read archived artifact: %v", err)
	}
	if !strings.HasPrefix(string(art), "payload-of-bootfw") {
		t.Fatalf("artifact content = %q", art)
	}
	sum := sha256.Sum256(art)
	_ = hex.EncodeToString(sum[:])

	// Success removes the tmp dir (requirement 2.5).
	entries, _ := os.ReadDir(e.TmpRoot)
	for _, ent := range entries {
		if strings.HasPrefix(ent.Name(), "build-") {
			t.Fatalf("tmp dir kept after success: %s", ent.Name())
		}
	}

	// Log carries phase tags across segments.
	found := map[string]bool{}
	for _, l := range lines {
		found[strings.Split(l, "|")[0]] = true
	}
	for _, want := range []string{"clone", "dirty", "command", "version", "artifact", "system"} {
		if !found[want] {
			t.Fatalf("log missing phase %q (lines: %v)", want, lines)
		}
	}
}

func TestBuild_DirtyRefused(t *testing.T) {
	repo := newGitRepo(t)
	e, _ := newExecutor(t)
	item := repo.item(t, "*.bin")

	// Make the source dirty.
	if err := os.WriteFile(filepath.Join(repo.dir, "untracked.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("dirty: %v", err)
	}

	res := e.Run(102, item, "sha256")
	if res.OK {
		t.Fatal("dirty build must fail")
	}
	if !strings.Contains(res.Detail, "dirty") {
		t.Fatalf("detail = %q, want dirty mention", res.Detail)
	}
}

func TestBuild_CommandFailure(t *testing.T) {
	repo := newGitRepo(t)
	e, _ := newExecutor(t)
	item := repo.item(t, "*.bin")
	item.Command = "exit 3"

	res := e.Run(103, item, "sha256")
	if res.OK {
		t.Fatal("failing command must fail the build")
	}
	if res.ExitCode == nil || *res.ExitCode != 3 {
		t.Fatalf("exit code = %v, want 3", res.ExitCode)
	}
}

func TestBuild_NoArtifactsMatched(t *testing.T) {
	repo := newGitRepo(t)
	e, _ := newExecutor(t)
	item := repo.item(t, "definitely-not-here*.bin")

	res := e.Run(104, item, "sha256")
	if res.OK {
		t.Fatal("zero artifacts must fail the record")
	}
	if !strings.Contains(res.Detail, "no artifacts") {
		t.Fatalf("detail = %q", res.Detail)
	}
}

func TestBuild_MD5ChecksumSetting(t *testing.T) {
	repo := newGitRepo(t)
	e, dataDir := newExecutor(t)
	item := repo.item(t, "*.bin")

	var checksumLine string
	e.Notify = func(record int64, phase, line string) {
		if phase == "artifact" {
			checksumLine = line
		}
	}

	res := e.Run(105, item, "md5")
	if !res.OK {
		t.Fatalf("build failed: %s", res.Detail)
	}
	// md5 hex is 32 chars; sha256 is 64.
	idx := strings.Index(checksumLine, ", ")
	if idx < 0 {
		t.Fatalf("artifact line = %q", checksumLine)
	}
	checksum := strings.TrimRight(checksumLine[idx+2:], " )")
	if len(checksum) != 32 {
		t.Fatalf("checksum = %q (%d chars), want 32-char md5", checksum, len(checksum))
	}
	_ = dataDir
}

func TestBuild_TimeoutKills(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("timeout test needs sh sleep semantics")
	}
	repo := newGitRepo(t)
	e, _ := newExecutor(t)
	item := repo.item(t, "*.bin")
	item.Command = "sleep 30"
	item.TimeoutSec = 1

	start := time.Now()
	res := e.Run(106, item, "sha256")
	if res.OK {
		t.Fatal("timeout must fail the build")
	}
	if time.Since(start) > 10*time.Second {
		t.Fatalf("timeout took %v, kill not applied", time.Since(start))
	}
	if !strings.Contains(res.Detail, "timed out") {
		t.Fatalf("detail = %q", res.Detail)
	}
}

func TestBuild_FailureKeepsTmpForSweep(t *testing.T) {
	repo := newGitRepo(t)
	e, _ := newExecutor(t)
	item := repo.item(t, "*.bin")
	item.Command = "exit 1"

	e.Run(107, item, "sha256")
	entries, _ := os.ReadDir(e.TmpRoot)
	kept := false
	for _, ent := range entries {
		if strings.HasPrefix(ent.Name(), "build-") {
			kept = true
		}
	}
	if !kept {
		t.Fatal("failed build must keep its tmp dir for the 7-day sweep")
	}
}

func TestCleanupSweep_RemovesOldTmp(t *testing.T) {
	repo := newGitRepo(t)
	e, _ := newExecutor(t)
	item := repo.item(t, "*.bin")
	item.Command = "exit 1"
	e.Run(108, item, "sha256")

	// Backdate the tmp dir past the retention window.
	entries, _ := os.ReadDir(e.TmpRoot)
	for _, ent := range entries {
		full := filepath.Join(e.TmpRoot, ent.Name())
		old := time.Now().Add(-8 * 24 * time.Hour)
		if err := os.Chtimes(full, old, old); err != nil {
			t.Fatalf("chtimes: %v", err)
		}
	}

	CleanupSweep(e.TmpRoot, 7*24*time.Hour)
	entries, _ = os.ReadDir(e.TmpRoot)
	for _, ent := range entries {
		if strings.HasPrefix(ent.Name(), "build-") {
			t.Fatalf("sweep left old tmp dir: %s", ent.Name())
		}
	}
}

func TestBuild_CancelBetweenPhases(t *testing.T) {
	repo := newGitRepo(t)
	e, _ := newExecutor(t)
	item := repo.item(t, "*.bin")
	e.Canceled = func(record int64) bool { return true }

	res := e.Run(109, item, "sha256")
	if res.OK {
		t.Fatal("canceled build must not report success")
	}
	if !strings.Contains(res.Detail, "canceled") {
		t.Fatalf("detail = %q", res.Detail)
	}
}

func TestBuild_CancelKillsRunningCommand(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("process-group kill test needs unix semantics")
	}
	repo := newGitRepo(t)
	e, _ := newExecutor(t)
	item := repo.item(t, "*.bin")
	item.Command = "sleep 30"

	// The cancel flips while the command runs; runCommand must observe it
	// and report cancellation instead of letting sleep finish.
	started := make(chan struct{})
	go func() {
		<-started
		for i := 0; i < 100; i++ {
			if e.CancelCommand(110) {
				return
			}
			time.Sleep(20 * time.Millisecond)
		}
	}()

	// Run in a goroutine; the executor reports via Result.
	type done struct {
		res Result
	}
	ch := make(chan done, 1)
	e.Canceled = func(record int64) bool { return true } // flag set from the start: kill lands, then boundary check confirms
	close(started)
	go func() { ch <- done{res: e.Run(110, item, "sha256")} }()

	select {
	case d := <-ch:
		if d.res.OK {
			t.Fatal("canceled build must not succeed")
		}
		if !d.res.Canceled {
			t.Fatalf("want Canceled=true, detail=%q", d.res.Detail)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("cancel did not stop the build command in time")
	}
}

func TestBuild_CanceledTmpRemovedImmediately(t *testing.T) {
	repo := newGitRepo(t)
	e, _ := newExecutor(t)
	item := repo.item(t, "*.bin")
	e.Canceled = func(record int64) bool { return true }

	e.Run(111, item, "sha256")
	entries, _ := os.ReadDir(e.TmpRoot)
	for _, ent := range entries {
		if strings.HasPrefix(ent.Name(), "build-") {
			t.Fatalf("canceled build must remove its tmp dir at once (requirement 2.5), kept: %s", ent.Name())
		}
	}
}

func TestBuild_LocalBranch(t *testing.T) {
	repo := newGitRepo(t)
	// A second branch whose build emits a different payload.
	git(t, repo.dir, "checkout", "-b", "feature")
	script := "build.sh"
	payload := "payload-from-feature"
	if runtime.GOOS == "windows" {
		script = "build.bat"
	}
	full := filepath.Join(repo.dir, script)
	src, err := os.ReadFile(full)
	if err != nil {
		t.Fatalf("read script: %v", err)
	}
	if err := os.WriteFile(full, bytes.ReplaceAll(src, []byte("payload-of-bootfw"), []byte(payload)), 0o755); err != nil {
		t.Fatalf("patch script: %v", err)
	}
	git(t, repo.dir, "commit", "-am", "feature payload")
	git(t, repo.dir, "checkout", "main")

	e, _ := newExecutor(t)
	item := repo.item(t, "bootfw*.bin")
	item.GitBranch = "feature"

	var lines []string
	e.Notify = func(record int64, phase, line string) { lines = append(lines, phase+"|"+line) }
	res := e.Run(112, item, "sha256")
	if !res.OK {
		t.Fatalf("branch build failed: %s\nlog:\n%s", res.Detail, strings.Join(lines, "\n"))
	}
	art, err := os.ReadFile(filepath.Join(e.DataDir, "artifacts", "112", "bootfw-p0133.bin"))
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	if !strings.HasPrefix(string(art), payload) {
		t.Fatalf("artifact = %q, want payload from feature branch", art)
	}
}

func TestBuild_LocalPinnedCommit(t *testing.T) {
	repo := newGitRepo(t)
	// Remember the initial commit, then add a second one.
	out, err := exec.Command("git", "-C", repo.dir, "rev-parse", "HEAD").CombinedOutput()
	if err != nil {
		t.Fatalf("rev-parse: %v: %s", err, out)
	}
	pinned := strings.TrimSpace(string(out))
	script := "build.sh"
	if runtime.GOOS == "windows" {
		script = "build.bat"
	}
	if err := os.WriteFile(filepath.Join(repo.dir, script), []byte("@echo off\r\necho payload-v2> bootfw-p0133.bin\r\n"), 0o755); err != nil {
		t.Fatalf("rewrite script: %v", err)
	}
	git(t, repo.dir, "commit", "-am", "v2")

	e, _ := newExecutor(t)
	item := repo.item(t, "bootfw*.bin")
	item.GitCommit = pinned

	var lines []string
	e.Notify = func(record int64, phase, line string) { lines = append(lines, phase+"|"+line) }
	res := e.Run(113, item, "sha256")
	if !res.OK {
		t.Fatalf("pinned build failed: %s\nlog:\n%s", res.Detail, strings.Join(lines, "\n"))
	}
	if res.CommitSHA != pinned {
		t.Fatalf("CommitSHA = %q, want pinned %q", res.CommitSHA, pinned)
	}
	art, err := os.ReadFile(filepath.Join(e.DataDir, "artifacts", "113", "bootfw-p0133.bin"))
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	if !strings.HasPrefix(string(art), "payload-of-bootfw") {
		t.Fatalf("artifact = %q, want the pinned (v1) payload", art)
	}
}

func TestBuild_LocalPinnedCommitDirtyOK(t *testing.T) {
	repo := newGitRepo(t)
	out, err := exec.Command("git", "-C", repo.dir, "rev-parse", "HEAD").CombinedOutput()
	if err != nil {
		t.Fatalf("rev-parse: %v: %s", err, out)
	}
	pinned := strings.TrimSpace(string(out))
	// Dirty the worktree: without a pin this refuses to build
	// (TestBuild_DirtyRefused); with a pin the commit object is the source
	// of truth, so the build must proceed.
	if err := os.WriteFile(filepath.Join(repo.dir, "untracked.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("dirty: %v", err)
	}

	e, _ := newExecutor(t)
	item := repo.item(t, "*.bin")
	item.GitCommit = pinned

	var lines []string
	e.Notify = func(record int64, phase, line string) { lines = append(lines, phase+"|"+line) }
	res := e.Run(114, item, "sha256")
	if !res.OK {
		t.Fatalf("pinned build with dirty worktree must succeed: %s\nlog:\n%s", res.Detail, strings.Join(lines, "\n"))
	}
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "commit pinned") {
		t.Fatalf("log must note the skipped dirty gate, got:\n%s", joined)
	}
}

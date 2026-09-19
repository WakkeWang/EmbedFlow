//go:build windows

package builder

import (
	"os/exec"
	"strconv"
	"syscall"
)

// setPgid is a no-op on Windows: there are no process groups in the unix
// sense.
func setPgid(cmd *exec.Cmd) {}

// killProcessTree kills the whole process tree on Windows via taskkill
// (Go's Process.Kill only kills the shell; the command's children -- ping,
// toolchains -- would survive, keep file handles open, and block the
// tmp-dir cleanup).
func killProcessTree(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}
	kill := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
	kill.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_ = kill.Run()
}

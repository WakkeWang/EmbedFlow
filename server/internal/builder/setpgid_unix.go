//go:build !windows

package builder

import (
	"os/exec"
	"syscall"
)

// setPgid puts the build process in its own process group so a timeout or
// user cancel kills the whole tree (bb.sh spawning toolchains must not
// orphan work).
func setPgid(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// killProcessTree kills the process group of a still-running command.
// On unix the Setpgid'd group id equals the child's pid, so one signal
// reaps the shell and every child it spawned.
func killProcessTree(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}

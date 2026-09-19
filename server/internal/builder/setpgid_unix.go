//go:build !windows

package builder

import (
	"os/exec"
	"syscall"
)

// setPgid puts the build process in its own process group so a timeout
// kills the whole tree (bb.sh spawning toolchains must not orphan work).
func setPgid(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

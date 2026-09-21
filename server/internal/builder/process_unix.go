// Platform process plumbing: shell selection. The process-group kill lives
// in setpgid_unix.go / setpgid_windows.go.
package builder

import (
	"context"
	"os/exec"
	"runtime"
)

// shellName is the shell the build command runs through (requirement 2.1:
// 构建命令 like ./bb.sh ... are shell scripts; the server runs on Linux).
func shellName() string {
	if runtime.GOOS == "windows" {
		return "cmd"
	}
	return "sh"
}

// shellArgs wraps a command string for the platform shell. Unix needs the
// command as its own argv entry ("sh" "-c" COMMAND); Windows cmd accepts
// the whole line as one argument after /C.
func shellArgs(command string) []string {
	if runtime.GOOS == "windows" {
		return []string{"/C " + command}
	}
	return []string{"-c", command}
}

// runShell builds the exec.Cmd for one shell command line. The caller
// swaps in CommandContext when a timeout/cancel governs the run.
func runShell(command string) *exec.Cmd {
	return exec.Command(shellName(), shellArgs(command)...)
}

// runShellContext is runShell bound to a context (timeout + cancel kill the
// process group via cmd.Cancel, see runCommand).
func runShellContext(ctx context.Context, command string) *exec.Cmd {
	return exec.CommandContext(ctx, shellName(), shellArgs(command)...)
}

// Platform process plumbing: shell selection. The process-group kill lives
// in setpgid_unix.go / setpgid_windows.go.
package builder

import "runtime"

// shellName is the shell the build command runs through (requirement 2.1:
// 构建命令 like ./bb.sh ... are shell scripts; the server runs on Linux).
func shellName() string {
	if runtime.GOOS == "windows" {
		return "cmd"
	}
	return "sh"
}

// shellArg wraps a command string for the platform shell.
func shellArg(command string) string {
	if runtime.GOOS == "windows" {
		return "/C " + command
	}
	return "-c"
}

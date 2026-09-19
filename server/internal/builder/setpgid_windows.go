//go:build windows

package builder

import "os/exec"

// setPgid is a no-op on Windows: CommandContext's Process.Kill is the
// process-tree fallback there.
func setPgid(cmd *exec.Cmd) {}

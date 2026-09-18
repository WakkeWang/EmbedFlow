// Package paths owns the data-directory standard layout (decision CEO-10A):
//
//	<data>/embedflow.db            SQLite database
//	<data>/sessions/<id>/log.txt   per-session full log
//	<data>/artifacts/<build-id>/   build artifacts
//	<data>/tmp/                    scratch space
//
// Backup = copy the whole directory.
package paths

import (
	"os"
	"path/filepath"
	"strconv"
)

// Subdirectories created under the data root.
var subdirs = []string{"sessions", "artifacts", "tmp"}

// Ensure creates the standard layout under dir, idempotently.
func Ensure(dir string) error {
	for _, sub := range append([]string{""}, subdirs...) {
		if err := os.MkdirAll(filepath.Join(dir, sub), 0o755); err != nil {
			return err
		}
	}
	return nil
}

// SessionLog is the log file path for a session ID.
func SessionLog(dataDir string, sessionID int64) string {
	return filepath.Join(dataDir, "sessions", strconv.FormatInt(sessionID, 10), "log.txt")
}

// SessionDir is the per-session directory.
func SessionDir(dataDir string, sessionID int64) string {
	return filepath.Join(dataDir, "sessions", strconv.FormatInt(sessionID, 10))
}

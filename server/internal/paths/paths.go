// Package paths owns the data-directory standard layout (decision CEO-10A):
//
//	<data>/embedflow.db            SQLite database
//	<data>/sessions/<id>/log.txt   per-session full log
//	<data>/artifacts/<build-id>/   build artifacts
//	<data>/tmp/                    scratch space
//	<data>/builds/<record-id>/log.txt  per-build-record log
//
// Backup = copy the whole directory.
package paths

import (
	"os"
	"path/filepath"
	"strconv"
)

// Subdirectories created under the data root.
var subdirs = []string{"sessions", "artifacts", "tmp", "builds"}

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

// ArtifactDir is the archive directory for one build record's artifacts.
func ArtifactDir(dataDir string, recordID int64) string {
	return filepath.Join(dataDir, "artifacts", strconv.FormatInt(recordID, 10))
}

// BuildLog is the log file path for one build record.
func BuildLog(dataDir string, recordID int64) string {
	return filepath.Join(dataDir, "builds", strconv.FormatInt(recordID, 10), "log.txt")
}

// BuildLogDir is the per-record log directory.
func BuildLogDir(dataDir string, recordID int64) string {
	return filepath.Join(dataDir, "builds", strconv.FormatInt(recordID, 10))
}

// BuildTmp is the scratch clone/build directory for one build record.
// It lives under the configurable tmp root, not the data dir, when the
// admin points build_tmp_dir elsewhere.
func BuildTmp(tmpRoot string, recordID int64) string {
	return filepath.Join(tmpRoot, "build-"+strconv.FormatInt(recordID, 10))
}

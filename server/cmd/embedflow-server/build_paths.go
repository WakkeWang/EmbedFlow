package main

import (
	"path/filepath"
	"strconv"
	"strings"
)

// Build-module filesystem path helpers. The paths package owns the layout;
// these wrap it for the cmd layer (builds/<id>/log.txt and artifacts/<id>/).

func buildLogPath(dataDir string, recordID int64) string {
	return filepath.Join(dataDir, "builds", strconv.FormatInt(recordID, 10), "log.txt")
}

func artifactDirPath(dataDir string, recordID int64) string {
	return filepath.Join(dataDir, "artifacts", strconv.FormatInt(recordID, 10))
}

func artifactFilePath(dataDir string, recordID int64, name string) string {
	return filepath.Join(artifactDirPath(dataDir, recordID), name)
}

// splitLines splits keeping trailing empty lines out.
func splitLines(s string) []string {
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

func joinLines(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}

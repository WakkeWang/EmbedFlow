package main

import (
	"github.com/WakkeWang/EmbedFlow/server/internal/paths"
)

// Build-module filesystem path helpers delegate to the paths package (the
// single owner of the data-dir layout). The cmd layer just binds record ids.

func buildLogPath(dataDir string, recordID int64) string {
	return paths.BuildLog(dataDir, recordID)
}

func artifactDirPath(dataDir string, recordID int64) string {
	return paths.ArtifactDir(dataDir, recordID)
}

func artifactFilePath(dataDir string, recordID int64, name string) string {
	// paths.ArtifactDir + the artifact's file name; the name comes from the
	// store (never user input at this point) and is a bare file name.
	return paths.ArtifactDir(dataDir, recordID) + string('/') + name
}

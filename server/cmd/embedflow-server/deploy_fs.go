package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// Small filesystem/parse helpers for the deploy module (single-line call
// sites; error handling stays at the caller).

func mkdirAll(dir string) error {
	return os.MkdirAll(dir, 0o755)
}

func openAppend(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
}

// pathPIDValue parses a path segment value (the {pid} wildcard).
func pathPIDValue(r *http.Request, name string) (int64, error) {
	return strconv.ParseInt(r.PathValue(name), 10, 64)
}

// guard keeps filepath referenced when helpers above change shape.
var _ = filepath.Join

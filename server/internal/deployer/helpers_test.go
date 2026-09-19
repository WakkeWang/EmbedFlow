package deployer

import (
	"os"
	"strings"
	"testing"
)

func readTail(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(raw)
}

// unused guard for strings in the helper file.
var _ = strings.Contains

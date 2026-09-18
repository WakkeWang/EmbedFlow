package paths

import (
	"os"
	"path/filepath"
	"testing"
)

// CEO-10A: the data directory has a standard layout -- backup means copying
// the whole directory.
func TestEnsure_CreatesStandardLayout(t *testing.T) {
	dir := t.TempDir()
	if err := Ensure(dir); err != nil {
		t.Fatalf("ensure: %v", err)
	}
	for _, want := range []string{
		filepath.Join(dir, "sessions"),
		filepath.Join(dir, "artifacts"),
		filepath.Join(dir, "tmp"),
	} {
		st, err := os.Stat(want)
		if err != nil || !st.IsDir() {
			t.Fatalf("%s missing or not a dir (err %v)", want, err)
		}
	}
}

func TestSessionLogPath_Layout(t *testing.T) {
	got := SessionLog("data", 42)
	want := filepath.Join("data", "sessions", "42", "log.txt")
	if got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}
}

func TestEnsure_Idempotent(t *testing.T) {
	dir := t.TempDir()
	if err := Ensure(dir); err != nil {
		t.Fatalf("first: %v", err)
	}
	if err := Ensure(dir); err != nil {
		t.Fatalf("second: %v", err)
	}
}

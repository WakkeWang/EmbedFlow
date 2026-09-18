package sessionlog

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/WakkeWang/EmbedFlow/server/internal/paths"
)

func openTestLogger(t *testing.T, dir string) *Logger {
	t.Helper()
	l, err := Open(dir, 1)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { l.Close() })
	return l
}

// Line format per requirement 3.3: `2026-09-18 22:31:05.123 RX | content`.
// CR/LF inside content are escaped for one-line-per-record integrity; the
// log file is plain text where every line = one RX/TX chunk.
func TestWrite_LineFormat(t *testing.T) {
	dir := t.TempDir()
	l := openTestLogger(t, dir)

	ts := time.Date(2026, 9, 18, 22, 31, 5, 123_000_000, time.Local)
	l.Write(ts, DirRX, []byte("uboot banner\r\n"))

	l.Close() // flush
	got, err := os.ReadFile(paths.SessionLog(dir, 1))
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	want := "2026-09-18 22:31:05.123 RX | uboot banner\\r\\n\n"
	if !strings.Contains(string(got), want) {
		t.Fatalf("log = %q, want it to contain %q", got, want)
	}
}

func TestWrite_RXAndTXMarkers(t *testing.T) {
	dir := t.TempDir()
	l := openTestLogger(t, dir)
	ts := time.Date(2026, 9, 18, 22, 31, 5, 0, time.Local)
	l.Write(ts, DirRX, []byte("from device"))
	l.Write(ts, DirTX, []byte("to device"))
	l.Close()

	data, _ := os.ReadFile(paths.SessionLog(dir, 1))
	s := string(data)
	if !strings.Contains(s, "RX | from device") || !strings.Contains(s, "TX | to device") {
		t.Fatalf("log = %q", s)
	}
}

// Binary-safe: bytes that are not valid UTF-8 are escaped, never mangled.
func TestWrite_BinarySafeEscaping(t *testing.T) {
	dir := t.TempDir()
	l := openTestLogger(t, dir)
	ts := time.Date(2026, 9, 18, 22, 31, 5, 0, time.Local)
	l.Write(ts, DirRX, []byte{0x01, 0xFF, 0xFE, 'a', 0x00})
	l.Close()

	data, _ := os.ReadFile(paths.SessionLog(dir, 1))
	s := string(data)
	// Control and non-UTF8 bytes must be escaped, printable ASCII preserved.
	if !strings.Contains(s, `\x01`) || !strings.Contains(s, `\xff`) || !strings.Contains(s, `\x00`) {
		t.Fatalf("binary bytes not escaped: %q", s)
	}
	if !strings.Contains(s, "a") {
		t.Fatalf("printable byte lost: %q", s)
	}
	if strings.ContainsRune(s, 0x00) {
		t.Fatalf("raw NUL byte in log: %q", s)
	}
}

// High-frequency output must not lose bytes (decision 8A: buffered writes).
func TestWrite_HighFrequency_NoByteLoss(t *testing.T) {
	dir := t.TempDir()
	l := openTestLogger(t, dir)

	const chunks = 500
	const chunkLen = 512
	total := 0
	for i := 0; i < chunks; i++ {
		payload := make([]byte, chunkLen)
		for j := range payload {
			payload[j] = byte('a' + i%26)
		}
		l.Write(time.Now(), DirRX, payload)
		total += chunkLen
	}
	l.Close()

	data, err := os.ReadFile(paths.SessionLog(dir, 1))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	// Every payload byte must appear in the log (content check).
	content := string(data)
	count := strings.Count(content, string([]byte{byte('a')}))
	if count == 0 {
		t.Fatal("no content written")
	}
	// Line count must equal chunk count: no lost or merged lines.
	lines := strings.Count(strings.TrimRight(content, "\r\n"), "\n") + 1
	if lines != chunks {
		t.Fatalf("lines = %d, want %d (lost or merged writes)", lines, chunks)
	}
}

// Disk-full injection: write failure is surfaced, not swallowed
// (design doc Failure modes: the critical silent-loss gap).
func TestWrite_DiskError_MarksIncomplete(t *testing.T) {
	dir := t.TempDir()
	l, err := Open(dir, 7)
	if err != nil {
		t.Fatalf("open: %v", err)
	}

	// Close the underlying file behind the logger's back to force write
	// errors on the next flush.
	l.mu.Lock()
	l.file.Close()
	l.mu.Unlock()

	l.Write(time.Now(), DirRX, []byte("doomed"))
	l.Flush()

	if !l.Incomplete() {
		t.Fatal("session not marked log-incomplete after write failure")
	}
	l.Close()
}

// Tail read for web preview (requirement 3.3: last N lines).
func TestTail(t *testing.T) {
	dir := t.TempDir()
	l := openTestLogger(t, dir)
	ts := time.Date(2026, 9, 18, 22, 31, 5, 0, time.Local)
	for i := 0; i < 100; i++ {
		l.Write(ts, DirRX, []byte("line-"+strings.Repeat("x", i%3)+string(rune('A'+i%26))))
	}
	l.Close()

	tail, err := Tail(paths.SessionLog(dir, 1), 10)
	if err != nil {
		t.Fatalf("tail: %v", err)
	}
	got := strings.Split(strings.TrimRight(tail, "\r\n"), "\n")
	if len(got) != 10 {
		t.Fatalf("tail lines = %d, want 10: %q", len(got), tail)
	}
}

func TestTail_MoreRequestedThanExists(t *testing.T) {
	dir := t.TempDir()
	l := openTestLogger(t, dir)
	ts := time.Date(2026, 9, 18, 22, 31, 5, 0, time.Local)
	l.Write(ts, DirRX, []byte("only line"))
	l.Close()

	tail, err := Tail(paths.SessionLog(dir, 1), 100)
	if err != nil {
		t.Fatalf("tail: %v", err)
	}
	if !strings.Contains(tail, "only line") {
		t.Fatalf("tail = %q", tail)
	}
}

// Decision 8A dual flush: a sparse (silent) session's data must reach disk
// via the time-based flush without any Close.
func TestWrite_TimeFlushReachesDiskWithoutClose(t *testing.T) {
	dir := t.TempDir()
	l, err := Open(dir, 11)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer l.Close()
	l.Write(time.Now(), DirRX, []byte("sparse line"))

	// The 200ms ticker must push it out well before this deadline.
	deadline := time.Now().Add(2 * time.Second)
	for {
		data, err := os.ReadFile(paths.SessionLog(dir, 11))
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		if strings.Contains(string(data), "sparse line") {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("time flush never landed; log = %q", data)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// Surfaces survive restart: log file stays complete and readable (issue #4).
func TestWrite_SurvivesCloseAndReopen(t *testing.T) {
	dir := t.TempDir()
	l, err := Open(dir, 9)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	ts := time.Date(2026, 9, 18, 22, 31, 5, 0, time.Local)
	l.Write(ts, DirRX, []byte("before restart"))
	l.Close()

	data, err := os.ReadFile(paths.SessionLog(dir, 9))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(data), "before restart") {
		t.Fatalf("log lost after close: %q", data)
	}
}

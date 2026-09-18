// Package sessionlog persists session byte streams line by line with
// timestamps and RX/TX direction markers (requirement 3.3), buffered for
// high-frequency output (decision 8A). Write failures are surfaced as a
// "log incomplete" flag on the session -- never silently swallowed (design
// doc Failure modes: the disk-full gap).
package sessionlog

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/WakkeWang/EmbedFlow/server/internal/paths"
)

// Dir is the direction marker of a log line.
type Dir string

const (
	DirRX Dir = "RX" // device -> server
	DirTX Dir = "TX" // server -> device
)

// Default buffer sizes: 64KB flush threshold (aligned with the 64KB WS
// frame cap, decision 8A).
const (
	flushBytes    = 64 * 1024
	flushInterval = 200 * time.Millisecond
)

// Logger writes one session's log file.
type Logger struct {
	mu     sync.Mutex
	file   *os.File
	buf    *bufio.Writer
	wrote  bool
	broken bool // a write failed: session is "log incomplete"
}

// Open creates/opens the log file for a session under the data directory.
func Open(dataDir string, sessionID int64) (*Logger, error) {
	if err := os.MkdirAll(paths.SessionDir(dataDir, sessionID), 0o755); err != nil {
		return nil, fmt.Errorf("sessionlog: mkdir: %w", err)
	}
	f, err := os.OpenFile(paths.SessionLog(dataDir, sessionID), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("sessionlog: open: %w", err)
	}
	return &Logger{file: f, buf: bufio.NewWriterSize(f, flushBytes)}, nil
}

// Write appends one chunk of session data as a log line. Binary input is
// escaped printable-safely; the log is always valid text.
func (l *Logger) Write(ts time.Time, dir Dir, data []byte) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.broken {
		return errBroken
	}
	line := fmt.Sprintf("%s %s | %s\n", ts.Format("2006-01-02 15:04:05.000"), dir, escape(data))
	if _, err := l.buf.WriteString(line); err != nil {
		l.broken = true
		return fmt.Errorf("sessionlog: buffer: %w", err)
	}
	l.wrote = true
	if l.buf.Buffered() >= flushBytes {
		return l.flushLocked()
	}
	return nil
}

// Flush pushes buffered lines to the OS.
func (l *Logger) Flush() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.flushLocked()
}

// Close flushes and releases the file.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if err := l.flushLocked(); err != nil {
		return err
	}
	if l.file != nil {
		err := l.file.Close()
		l.file = nil
		return err
	}
	return nil
}

// Incomplete reports whether any write to this log has failed. The session
// executor turns this into the "log incomplete" warning (WS alert to the
// operator) at the transport layer.
func (l *Logger) Incomplete() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.broken
}

func (l *Logger) flushLocked() error {
	if l.buf == nil {
		return nil
	}
	if err := l.buf.Flush(); err != nil {
		l.broken = true
		return fmt.Errorf("sessionlog: flush: %w", err)
	}
	if l.file != nil {
		if err := l.file.Sync(); err != nil {
			l.broken = true
			return fmt.Errorf("sessionlog: sync: %w", err)
		}
	}
	return nil
}

var errBroken = fmt.Errorf("sessionlog: broken (write already failed)")

// escape renders bytes printable-safely: printable ASCII passes through,
// \r \n \t are kept readable, everything else becomes \xHH.
func escape(data []byte) string {
	var b strings.Builder
	for _, c := range data {
		switch {
		case c == '\r':
			b.WriteString(`\r`)
		case c == '\n':
			b.WriteString(`\n`)
		case c == '\t':
			b.WriteString(`\t`)
		case c >= 0x20 && c < 0x7f:
			b.WriteByte(c)
		default:
			fmt.Fprintf(&b, `\x%02x`, c)
		}
	}
	return b.String()
}

// Tail returns the last n lines of a log file for web preview
// (requirement 3.3). Reads backward from the end in chunks so huge logs
// never load fully into memory.
func Tail(path string, n int) (string, error) {
	if n <= 0 {
		return "", nil
	}
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("sessionlog: tail open: %w", err)
	}
	defer f.Close()

	st, err := f.Stat()
	if err != nil {
		return "", err
	}
	end := st.Size()
	if end == 0 {
		return "", nil
	}

	const chunk = 8 * 1024
	// Walk back chunk by chunk until we hold more than n newlines (or hit
	// the file start), keeping the bytes we have seen.
	var seen []byte
	newlines := 0
	for end > 0 && newlines <= n {
		start := end - chunk
		if start < 0 {
			start = 0
		}
		buf := make([]byte, end-start)
		if _, err := f.ReadAt(buf, start); err != nil {
			return "", err
		}
		for _, c := range buf {
			if c == '\n' {
				newlines++
			}
		}
		seen = append(make([]byte, 0, len(buf)+len(seen)), append(buf, seen...)...)
		end = start
	}

	text := strings.TrimRight(string(seen), "\n")
	lines := strings.Split(text, "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	for i, l := range lines {
		lines[i] = strings.TrimSuffix(l, "\r")
	}
	return strings.Join(lines, "\n") + "\n", nil
}

package sessionlog

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/WakkeWang/EmbedFlow/server/internal/paths"
)

func serveTestLog(t *testing.T, dir string, sessionID int64, lines int) *httptest.Server {
	t.Helper()
	l, err := Open(dir, sessionID)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	for i := 0; i < lines; i++ {
		fmt.Fprintf(l.buf, "line %04d padding padding padding\n", i)
	}
	if err := l.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	srv := httptest.NewServer(DownloadHandler(dir))
	t.Cleanup(srv.Close)
	return srv
}

func TestDownload_Full(t *testing.T) {
	dir := t.TempDir()
	srv := serveTestLog(t, dir, 1, 50)

	resp, err := http.Get(srv.URL + "/api/sessions/1/log/download")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "line 0049") {
		t.Fatalf("body missing last line: %q", body[:100])
	}
	if !strings.Contains(resp.Header.Get("Content-Disposition"), "attachment") {
		t.Fatal("missing attachment disposition")
	}
}

func TestDownload_RangeRequest(t *testing.T) {
	// Decision 9A: HTTP Range streaming -- browsers get resume for free.
	dir := t.TempDir()
	srv := serveTestLog(t, dir, 2, 100)

	req, _ := http.NewRequest("GET", srv.URL+"/api/sessions/2/log/download", nil)
	req.Header.Set("Range", "bytes=0-99")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusPartialContent {
		t.Fatalf("status = %d, want 206", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if len(body) != 100 {
		t.Fatalf("range bytes = %d, want 100 (bytes=0-99)", len(body))
	}
}

func TestDownload_RangeResume(t *testing.T) {
	// Resume semantics: a mid-file Range returns exactly the requested span.
	dir := t.TempDir()
	srv := serveTestLog(t, dir, 3, 100)

	st, _ := os.Stat(paths.SessionLog(dir, 3))
	req, _ := http.NewRequest("GET", srv.URL+"/api/sessions/3/log/download", nil)
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", st.Size()-10))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if len(body) != 10 {
		t.Fatalf("resume bytes = %d, want 10", len(body))
	}
}

func TestDownload_MissingSession_404(t *testing.T) {
	srv := httptest.NewServer(DownloadHandler(t.TempDir()))
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/api/sessions/999/log/download")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", resp.StatusCode)
	}
}

func TestDownload_RejectsBadSessionID(t *testing.T) {
	srv := httptest.NewServer(DownloadHandler(t.TempDir()))
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/api/sessions/notanumber/log/download")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", resp.StatusCode)
	}
}

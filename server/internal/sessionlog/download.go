package sessionlog

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/WakkeWang/EmbedFlow/server/internal/paths"
)

// DownloadHandler serves full-session log downloads with HTTP Range support
// (decision 9A): Go's http.ServeContent implements 206/If-Range/resume
// natively, so browsers get seek and resume for free.
func DownloadHandler(dataDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := sessionIDFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "bad session id", http.StatusBadRequest)
			return
		}
		f, err := os.Open(paths.SessionLog(dataDir, sessionID))
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, "no such session log", http.StatusNotFound)
				return
			}
			http.Error(w, "open log", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		st, err := f.Stat()
		if err != nil {
			http.Error(w, "stat log", http.StatusInternalServerError)
			return
		}
		name := fmt.Sprintf("session-%d-log.txt", sessionID)
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%q`, name))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		// ServeContent implements Range, If-Range, Content-Range, HEAD.
		http.ServeContent(w, r, name, st.ModTime(), f)
	})
}

// sessionIDFromPath parses "/api/sessions/<id>/log/download".
func sessionIDFromPath(p string) (int64, error) {
	const prefix = "/api/sessions/"
	const suffix = "/log/download"
	if len(p) < len(prefix)+len(suffix)+1 || p[:len(prefix)] != prefix || p[len(p)-len(suffix):] != suffix {
		return 0, errors.New("sessionlog: unexpected path shape")
	}
	return strconv.ParseInt(p[len(prefix):len(p)-len(suffix)], 10, 64)
}

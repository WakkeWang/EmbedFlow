package main

import (
	"fmt"
	"net/http"
	"os"
)

// buildLogDownloadHandler serves full build-log downloads with HTTP Range
// support via http.ServeContent (the sessionlog.DownloadHandler pattern,
// decision 9A): browsers get seek and resume for free.
func buildLogDownloadHandler(dataDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := pathID(r)
		if err != nil {
			http.Error(w, "bad record id", http.StatusBadRequest)
			return
		}
		path := buildLogPath(dataDir, id)
		f, err := os.Open(path)
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, "no such build log", http.StatusNotFound)
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
		name := fmt.Sprintf("build-%d-log.txt", id)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", name))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.ServeContent(w, r, name, st.ModTime(), f)
	})
}

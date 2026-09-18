// Package authmw holds the M1 HTTP middleware: token gate for API routes
// and the same-origin check for WebSocket upgrades (decision 3A).
package authmw

import (
	"crypto/subtle"
	"net/http"
	"net/url"
	"strings"
)

// Token wraps a handler so it only serves requests whose token validates
// (Authorization: Bearer, or ?token= for download links). M1 runs a single
// bootstrap token when static is non-empty; otherwise validate decides.
func Token(static string, validate func(string) bool, next http.Handler) http.Handler {
	check := func(tok string) bool {
		if static != "" {
			return subtle.ConstantTimeCompare([]byte(tok), []byte(static)) == 1
		}
		return validate != nil && validate(tok)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := bearer(r)
		if got == "" || !check(got) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func bearer(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return r.URL.Query().Get("token")
}

// SameOrigin allows only same-origin WebSocket upgrades (decision 3A):
// the Origin header must be absent (non-browser clients) or match the
// Host. Cross-origin browser upgrades are refused with 403.
func SameOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			u, err := url.Parse(origin)
			if err != nil || u.Host == "" || u.Host != r.Host {
				http.Error(w, "cross-origin websocket refused", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

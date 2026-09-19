package authmw

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestToken_AcceptsBearer(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	srv := httptest.NewServer(Token("secret", nil, next))
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL, nil)
	req.Header.Set("Authorization", "Bearer secret")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
}

func TestToken_AcceptsQueryParam(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	srv := httptest.NewServer(Token("secret", nil, next))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "?token=secret")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
}

func TestToken_RejectsMissingAndWrong(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	srv := httptest.NewServer(Token("secret", nil, next))
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("missing token status = %d", resp.StatusCode)
	}

	resp, err = http.Get(srv.URL + "?token=wrong")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("wrong token status = %d", resp.StatusCode)
	}
}

func TestToken_ValidatorMode(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	valid := map[string]bool{"tok-1": true}
	srv := httptest.NewServer(Token("", func(tok string) bool { return valid[tok] }, next))
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL, nil)
	req.Header.Set("Authorization", "Bearer tok-1")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("valid token status = %d", resp.StatusCode)
	}

	req, _ = http.NewRequest("GET", srv.URL, nil)
	req.Header.Set("Authorization", "Bearer tok-2")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("invalid token status = %d", resp.StatusCode)
	}
}

func TestSameOrigin_SameOriginAllowed(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	srv := httptest.NewServer(SameOrigin(next))
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL, nil)
	req.Header.Set("Origin", srv.URL) // same host as the server
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
}

func TestSameOrigin_CrossOriginRefused(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	srv := httptest.NewServer(SameOrigin(next))
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL, nil)
	req.Header.Set("Origin", "http://evil.example.com:8080")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", resp.StatusCode)
	}
}

func TestSameOrigin_NoOriginAllowed(t *testing.T) {
	// Non-browser clients (the serial client) send no Origin header.
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	srv := httptest.NewServer(SameOrigin(next))
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
}

func TestRequireAdmin_RoleGate(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	validate := func(tok string) (string, string, bool) {
		switch tok {
		case "tok-admin":
			return "root", "admin", true
		case "tok-member":
			return "alice", "member", true
		}
		return "", "", false
	}
	srv := httptest.NewServer(RequireAdmin(validate, next))
	defer srv.Close()

	do := func(tok string) int {
		req, _ := http.NewRequest("GET", srv.URL, nil)
		if tok != "" {
			req.Header.Set("Authorization", "Bearer "+tok)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("do: %v", err)
		}
		resp.Body.Close()
		return resp.StatusCode
	}

	if got := do("tok-admin"); got != http.StatusOK {
		t.Fatalf("admin status = %d, want 200", got)
	}
	if got := do("tok-member"); got != http.StatusForbidden {
		t.Fatalf("member status = %d, want 403", got)
	}
	if got := do("tok-bogus"); got != http.StatusUnauthorized {
		t.Fatalf("bogus status = %d, want 401", got)
	}
	if got := do(""); got != http.StatusUnauthorized {
		t.Fatalf("no-token status = %d, want 401", got)
	}
}

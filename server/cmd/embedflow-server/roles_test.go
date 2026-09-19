package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/WakkeWang/EmbedFlow/server/internal/store"
)

// Role-gate integration (requirement 1.5): admin manages configuration,
// members execute. M2 turns the previously decorative two-level roles into
// enforced HTTP behavior.

func apiCall(t *testing.T, s *srv, tok, method, path string, body any) *http.Response {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req, err := http.NewRequest(method, s.URL+path, &buf)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do %s %s: %v", method, path, err)
	}
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

func TestRoles_MemberCannotWriteConfig(t *testing.T) {
	s := newTestServer(t)
	adminTok := s.hub.tokens.issueForTesting("admin")
	memberTok := s.hub.tokens.issueForRole("alice", "member")

	// Member reads are fine.
	if resp := apiCall(t, s, memberTok, "GET", "/api/projects", nil); resp.StatusCode != http.StatusOK {
		t.Fatalf("member GET projects = %d, want 200", resp.StatusCode)
	}
	if resp := apiCall(t, s, memberTok, "GET", "/api/devices", nil); resp.StatusCode != http.StatusOK {
		t.Fatalf("member GET devices = %d, want 200", resp.StatusCode)
	}

	// Member config writes are refused with 403.
	for _, tc := range []struct{ method, path string }{
		{"POST", "/api/projects"},
		{"DELETE", "/api/projects/1"},
		{"POST", "/api/devices"},
		{"DELETE", "/api/devices/1"},
		{"POST", "/api/expect-rules"},
		{"PUT", "/api/expect-rules/1"},
		{"DELETE", "/api/expect-rules/1"},
	} {
		if resp := apiCall(t, s, memberTok, tc.method, tc.path, map[string]string{"name": "x"}); resp.StatusCode != http.StatusForbidden {
			t.Fatalf("member %s %s = %d, want 403", tc.method, tc.path, resp.StatusCode)
		}
	}

	// The same writes succeed for admin (except deletes of nonexistent ids,
	// which are 204/404 but never 403 -- the gate passed).
	if resp := apiCall(t, s, adminTok, "POST", "/api/projects", map[string]string{"name": "P-x"}); resp.StatusCode != http.StatusOK {
		t.Fatalf("admin POST projects = %d, want 200", resp.StatusCode)
	}
	if resp := apiCall(t, s, adminTok, "POST", "/api/devices", map[string]string{"name": "board", "project": "P-x"}); resp.StatusCode != http.StatusOK {
		t.Fatalf("admin POST devices = %d, want 200", resp.StatusCode)
	}
}

func TestRoles_UserManagement(t *testing.T) {
	s := newTestServer(t)
	ctx := context.Background()
	adminTok := s.hub.tokens.issueForTesting("admin")
	memberTok := s.hub.tokens.issueForRole("alice", "member")

	// Admin creates a member account.
	resp := apiCall(t, s, adminTok, "POST", "/api/users", map[string]string{
		"username": "bob", "password": "pw-bob", "role": "member",
	})
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("create user = %d, want 204", resp.StatusCode)
	}
	if _, err := s.st.AuthenticateUser(ctx, "bob", "pw-bob"); err != nil {
		t.Fatalf("bob cannot log in: %v", err)
	}

	// Bad role refused.
	if resp := apiCall(t, s, adminTok, "POST", "/api/users", map[string]string{
		"username": "eve", "password": "pw", "role": "superuser",
	}); resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("bad role = %d, want 400", resp.StatusCode)
	}

	// Member cannot list users or reset passwords.
	if resp := apiCall(t, s, memberTok, "GET", "/api/users", nil); resp.StatusCode != http.StatusForbidden {
		t.Fatalf("member list users = %d, want 403", resp.StatusCode)
	}
	if resp := apiCall(t, s, memberTok, "POST", "/api/users/1/password", map[string]string{"password": "x"}); resp.StatusCode != http.StatusForbidden {
		t.Fatalf("member reset password = %d, want 403", resp.StatusCode)
	}

	// A user changes their own password (old must verify).
	selfTok := s.hub.tokens.issueForTesting("bob")
	if resp := apiCall(t, s, selfTok, "POST", "/api/users/self/password", map[string]string{
		"old": "wrong", "new": "pw2",
	}); resp.StatusCode != http.StatusForbidden {
		t.Fatalf("wrong old password = %d, want 403", resp.StatusCode)
	}
	if resp := apiCall(t, s, selfTok, "POST", "/api/users/self/password", map[string]string{
		"old": "pw-bob", "new": "pw2",
	}); resp.StatusCode != http.StatusNoContent {
		t.Fatalf("self password change = %d, want 204", resp.StatusCode)
	}
	if _, err := s.st.AuthenticateUser(ctx, "bob", "pw2"); err != nil {
		t.Fatalf("auth after change: %v", err)
	}
}

func TestRoles_UnknownUserToken(t *testing.T) {
	// A username deleted from the store (or never present) cannot resolve
	// for the self-password path: 401, not a panic.
	s := newTestServer(t)
	tok := s.hub.tokens.issueForTesting("ghost")
	resp := apiCall(t, s, tok, "POST", "/api/users/self/password", map[string]string{"old": "a", "new": "b"})
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("ghost self-password = %d, want 401", resp.StatusCode)
	}
}

// compile-time: keep the store import used even if the tests above change.
var _ = store.ErrBadCredentials

func TestRoles_ResponseBodyShape(t *testing.T) {
	// The users list returns ids (the admin reset path targets them).
	s := newTestServer(t)
	adminTok := s.hub.tokens.issueForTesting("admin")
	resp := apiCall(t, s, adminTok, "GET", "/api/users", nil)
	var users []store.UserWithID
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		t.Fatalf("decode users: %v", err)
	}
	if len(users) == 0 || !strings.HasPrefix(users[0].Username, "admin") {
		t.Fatalf("users = %+v", users)
	}
	if users[0].ID == 0 {
		t.Fatal("user ids missing")
	}
}

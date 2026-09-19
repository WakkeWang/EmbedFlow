// Shared helpers for the build-module tables: small JSON list codecs for
// the artifact-glob and prerequisite-group columns, and the bcrypt wrapper
// every password path goes through.
package store

import (
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// encodeJSONStrings serializes a string list for the artifacts_json column.
func encodeJSONStrings(list []string) string {
	if len(list) == 0 {
		return "[]"
	}
	raw, err := json.Marshal(list)
	if err != nil {
		return "[]" // strings always marshal; unreachable in practice
	}
	return string(raw)
}

// decodeJSONStrings parses a string list column; corrupt or empty -> nil.
func decodeJSONStrings(raw string) []string {
	if raw == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}

// hashPassword wraps bcrypt for every password-writing path.
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("store: hash: %w", err)
	}
	return string(hash), nil
}

// verifyPassword checks a plaintext against a stored bcrypt hash.
func verifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Package secretbox encrypts device credentials at rest (requirement 6.2:
// symmetric encryption keyed from the startup config, so a leaked database
// file alone does not expose SSH passwords). AES-256-GCM; ciphertext format
// "enc:v1:<base64(nonce|ciphertext)>" so future versions can re-key.
package secretbox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

const prefix = "enc:v1:"

// ErrNoKey marks a decrypt against an uninitialized Box (operator removed
// the key, or the row predates encryption). Callers surface it as "credential
// unreadable -- re-enter the password", never as a crash.
var ErrNoKey = errors.New("secretbox: no key configured")

// ErrNotEncrypted marks a value that does not carry the enc:v1 prefix
// (plaintext from a pre-encryption row).
var ErrNotEncrypted = errors.New("secretbox: value is not encrypted")

// Box encrypts and decrypts with one key. A key-less Box refuses both
// directions: reading encrypted rows must fail loudly (not silently return
// the ciphertext), and writing plaintext silently would break the
// "encrypted at rest" contract.
type Box struct {
	aead cipher.AEAD
}

// New builds a Box from a 32-byte key.
func New(key []byte) (*Box, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("secretbox: key must be 32 bytes, got %d", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("secretbox: cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("secretbox: gcm: %w", err)
	}
	return &Box{aead: aead}, nil
}

// Encrypt seals plaintext into the enc:v1 wire format.
func (b *Box) Encrypt(plaintext string) (string, error) {
	if b == nil || b.aead == nil {
		return "", ErrNoKey
	}
	nonce := make([]byte, b.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("secretbox: nonce: %w", err)
	}
	sealed := b.aead.Seal(nonce, nonce, []byte(plaintext), nil)
	return prefix + base64.StdEncoding.EncodeToString(sealed), nil
}

// Decrypt opens an enc:v1 value. Tampering fails GCM authentication.
func (b *Box) Decrypt(value string) (string, error) {
	if b == nil || b.aead == nil {
		return "", ErrNoKey
	}
	if !strings.HasPrefix(value, prefix) {
		return "", ErrNotEncrypted
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, prefix))
	if err != nil {
		return "", fmt.Errorf("secretbox: base64: %w", err)
	}
	ns := b.aead.NonceSize()
	if len(raw) < ns {
		return "", errors.New("secretbox: ciphertext too short")
	}
	plain, err := b.aead.Open(nil, raw[:ns], raw[ns:], nil)
	if err != nil {
		return "", fmt.Errorf("secretbox: open: %w", err)
	}
	return string(plain), nil
}

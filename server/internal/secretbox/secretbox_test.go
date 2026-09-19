package secretbox

import (
	"bytes"
	"errors"
	"testing"
)

func testKey() []byte {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i + 1)
	}
	return key
}

func TestRoundTrip(t *testing.T) {
	box, err := New(testKey())
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	enc, err := box.Encrypt("s3cret-pass")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if len(enc) < len("s3cret-pass") || enc == "s3cret-pass" {
		t.Fatalf("ciphertext = %q, want sealed form", enc)
	}
	if got, err := box.Decrypt(enc); err != nil || got != "s3cret-pass" {
		t.Fatalf("decrypt = %q, err %v", got, err)
	}

	// Same plaintext twice must produce different ciphertexts (random nonce).
	enc2, _ := box.Encrypt("s3cret-pass")
	if enc2 == enc {
		t.Fatal("nonce reuse: identical ciphertexts for identical plaintexts")
	}
}

func TestDecrypt_TamperFails(t *testing.T) {
	box, _ := New(testKey())
	enc, _ := box.Encrypt("s3cret-pass")
	raw := []byte(enc)
	raw[len(raw)-1] ^= 0xFF
	if _, err := box.Decrypt(string(raw)); err == nil {
		t.Fatal("tampered ciphertext must fail GCM auth")
	}
}

func TestDecrypt_WrongKeyFails(t *testing.T) {
	box1, _ := New(testKey())
	other := bytes.Clone(testKey())
	other[0] ^= 0xFF
	box2, _ := New(other)
	enc, _ := box1.Encrypt("s3cret-pass")
	if _, err := box2.Decrypt(enc); err == nil {
		t.Fatal("wrong-key decrypt must fail")
	}
}

func TestNoKey_RefusesBothDirections(t *testing.T) {
	var box *Box
	if _, err := box.Encrypt("x"); !errors.Is(err, ErrNoKey) {
		t.Fatalf("key-less encrypt err = %v, want ErrNoKey", err)
	}
	if _, err := box.Decrypt("enc:v1:AAAA"); !errors.Is(err, ErrNoKey) {
		t.Fatalf("key-less decrypt err = %v, want ErrNoKey", err)
	}

	keyless, err := New([]byte("short"))
	if err == nil {
		t.Fatal("short key must be rejected at construction")
	}
	_ = keyless
}

func TestDecrypt_PlaintextPassthrough(t *testing.T) {
	box, _ := New(testKey())
	if _, err := box.Decrypt("legacy-plaintext"); !errors.Is(err, ErrNotEncrypted) {
		t.Fatalf("plaintext err = %v, want ErrNotEncrypted", err)
	}
}

package nix

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

// PublicKey is a named Ed25519 public key used to verify NAR signatures.
// The zero value is not a valid key; use [ParsePublicKey] or [GenerateKey].
//
// Keys are encoded in "name:base64" format where the base64 value is the
// 32-byte Ed25519 public key. See the Nix C++ reference implementation:
// https://github.com/NixOS/nix/blob/master/src/libutil/signature/local-keys.cc
//
//	pub, _ := nix.ParsePublicKey("cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY=")
//	// pub.Name() == "cache.nixos.org-1"
type PublicKey struct {
	name string
	key  ed25519.PublicKey
}

// Name returns the key name, or "" if k is nil.
func (k *PublicKey) Name() string {
	if k == nil {
		return ""
	}
	return k.name
}

// String returns the key in "name:base64" format, or "<nil>" if k is nil.
func (k *PublicKey) String() string {
	if k == nil {
		return "<nil>"
	}
	return k.name + ":" + base64.StdEncoding.EncodeToString(k.key)
}

// MarshalText implements [encoding.TextMarshaler].
func (k *PublicKey) MarshalText() ([]byte, error) {
	return []byte(k.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (k *PublicKey) UnmarshalText(text []byte) error {
	parsed, err := ParsePublicKey(string(text))
	if err != nil {
		return err
	}
	*k = *parsed
	return nil
}

// ParsePublicKey parses a Nix public key in "name:base64" format.
// The base64-encoded value must decode to exactly 32 bytes.
//
//	pub, _ := nix.ParsePublicKey("cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY=")
//	// pub.Name() == "cache.nixos.org-1"
func ParsePublicKey(s string) (*PublicKey, error) {
	name, b64, err := splitKeyString(s, "public key")
	if err != nil {
		return nil, err
	}
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("nix: cannot parse public key %q: invalid base64: %w", s, err)
	}
	if len(b) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("nix: cannot parse public key %q: expected %d bytes, got %d", s, ed25519.PublicKeySize, len(b))
	}
	return &PublicKey{name: name, key: ed25519.PublicKey(b)}, nil
}

// PrivateKey is a named Ed25519 private key used to sign NARInfo fingerprints.
// The zero value is not a valid key; use [ParsePrivateKey] or [GenerateKey].
//
// Keys are encoded in "name:base64" format where the base64 value is the
// 64-byte Ed25519 private key (seed || public key). See the Nix C++ reference
// implementation:
// https://github.com/NixOS/nix/blob/master/src/libutil/signature/local-keys.cc
//
//	_, priv, _ := nix.GenerateKey("mycache-1", nil)
//	// priv.Name() == "mycache-1"
type PrivateKey struct {
	name string
	key  ed25519.PrivateKey
}

// Name returns the key name, or "" if k is nil.
func (k *PrivateKey) Name() string {
	if k == nil {
		return ""
	}
	return k.name
}

// PublicKey returns the corresponding public key, or nil if k is nil.
func (k *PrivateKey) PublicKey() *PublicKey {
	if k == nil {
		return nil
	}
	return &PublicKey{name: k.name, key: k.key.Public().(ed25519.PublicKey)}
}

// String returns the key in "name:base64" format, or "<nil>" if k is nil.
func (k *PrivateKey) String() string {
	if k == nil {
		return "<nil>"
	}
	return k.name + ":" + base64.StdEncoding.EncodeToString(k.key)
}

// MarshalText implements [encoding.TextMarshaler].
func (k *PrivateKey) MarshalText() ([]byte, error) {
	return []byte(k.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (k *PrivateKey) UnmarshalText(text []byte) error {
	parsed, err := ParsePrivateKey(string(text))
	if err != nil {
		return err
	}
	*k = *parsed
	return nil
}

// ParsePrivateKey parses a Nix private key in "name:base64" format.
// The base64-encoded value must decode to exactly 64 bytes.
//
//	priv, _ := nix.ParsePrivateKey("mycache-1:<base64>")
//	// priv.PublicKey().Name() == "mycache-1"
func ParsePrivateKey(s string) (*PrivateKey, error) {
	name, b64, err := splitKeyString(s, "private key")
	if err != nil {
		return nil, err
	}
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		// Use name rather than s to avoid logging private key material.
		return nil, fmt.Errorf("nix: cannot parse private key %q: invalid base64: %w", name, err)
	}
	if len(b) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("nix: cannot parse private key %q: expected %d bytes, got %d", name, ed25519.PrivateKeySize, len(b))
	}
	return &PrivateKey{name: name, key: ed25519.PrivateKey(b)}, nil
}

// GenerateKey generates a new named Ed25519 key pair. A nil r defaults to
// [crypto/rand.Reader].
//
//	pub, priv, _ := nix.GenerateKey("mycache-1", nil)
//	// pub.Name() == "mycache-1"
//	// priv.PublicKey().String() == pub.String()
func GenerateKey(name string, r io.Reader) (*PublicKey, *PrivateKey, error) {
	if r == nil {
		r = rand.Reader
	}
	pub, priv, err := ed25519.GenerateKey(r)
	if err != nil {
		return nil, nil, fmt.Errorf("nix: key generation failed: %w", err)
	}
	return &PublicKey{name: name, key: pub}, &PrivateKey{name: name, key: priv}, nil
}

// splitKeyString splits a "name:base64" key string at the first colon,
// returning an error if the separator is missing or the name is empty.
func splitKeyString(s, kind string) (name, b64 string, err error) {
	name, b64, ok := strings.Cut(s, ":")
	if !ok {
		return "", "", fmt.Errorf("nix: cannot parse %s: missing ':' separator", kind)
	}
	if name == "" {
		return "", "", fmt.Errorf("nix: cannot parse %s: empty name", kind)
	}
	return name, b64, nil
}

package nix

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

// Signature is a named Ed25519 signature over a NARInfo fingerprint.
// The zero value is not a valid signature; use [ParseSignature] or [SignNARInfo].
//
// Signatures are encoded in "name:base64" format where the base64 value is
// the 64-byte Ed25519 signature. See the Nix C++ reference implementation:
// https://github.com/NixOS/nix/blob/master/src/libutil/signature/local-keys.cc
//
//	sig, _ := nix.ParseSignature("cache.nixos.org-1:sDcIalq8XS+YQKtGsDi+...")
//	// sig.Name() == "cache.nixos.org-1"
type Signature struct {
	name string
	sig  [ed25519.SignatureSize]byte
}

// Name returns the name of the key that produced this signature, or "" if s is nil.
func (s *Signature) Name() string {
	if s == nil {
		return ""
	}
	return s.name
}

// String returns the signature in "name:base64" format, or "<nil>" if s is nil.
func (s *Signature) String() string {
	if s == nil {
		return "<nil>"
	}
	return s.name + ":" + base64.StdEncoding.EncodeToString(s.sig[:])
}

// MarshalText implements [encoding.TextMarshaler].
func (s *Signature) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (s *Signature) UnmarshalText(text []byte) error {
	parsed, err := ParseSignature(string(text))
	if err != nil {
		return err
	}
	*s = *parsed
	return nil
}

// ParseSignature parses a Nix NAR signature in "name:base64" format.
// The base64-encoded value must decode to exactly 64 bytes.
//
//	sig, _ := nix.ParseSignature("cache.nixos.org-1:sDcIalq8XS+YQKtGsDi+WlTXk59tMwPkBrgT/LNqJqXEIBMG4AXaYOnaB5dmAerq10TRShEz2C+ynvvL2sINCA==")
//	// sig.Name() == "cache.nixos.org-1"
func ParseSignature(s string) (*Signature, error) {
	name, b64, err := splitKeyString(s, "signature")
	if err != nil {
		return nil, err
	}
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("nix: cannot parse signature %q: invalid base64: %w", s, err)
	}
	if len(b) != ed25519.SignatureSize {
		return nil, fmt.Errorf("nix: cannot parse signature %q: expected %d bytes, got %d", s, ed25519.SignatureSize, len(b))
	}
	var sig Signature
	sig.name = name
	copy(sig.sig[:], b)
	return &sig, nil
}

// SignNARInfo signs the fingerprint of info using pk and returns the resulting
// [Signature].
//
//	_, priv, _ := nix.GenerateKey("mycache-1", nil)
//	sig, _ := nix.SignNARInfo(priv, &ni)
//	// sig.Name() == "mycache-1"
func SignNARInfo(pk *PrivateKey, info *NARInfo) (*Signature, error) {
	raw := ed25519.Sign(pk.key, []byte(info.Fingerprint()))
	var sig Signature
	sig.name = pk.name
	copy(sig.sig[:], raw)
	return &sig, nil
}

// VerifyNARInfo verifies sig against the fingerprint of info using the
// matching key from trusted (matched by key name). It returns an error if no
// trusted key matches the signature name or if the signature is
// cryptographically invalid.
//
//	pub, _ := nix.ParsePublicKey("cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY=")
//	err := nix.VerifyNARInfo([]*nix.PublicKey{pub}, &ni, ni.Sig[0])
func VerifyNARInfo(trusted []*PublicKey, info *NARInfo, sig *Signature) error {
	var pub *PublicKey
	for _, k := range trusted {
		if k.Name() == sig.Name() {
			pub = k
			break
		}
	}
	if pub == nil {
		return fmt.Errorf("nix: no trusted key matches signature name %q", sig.Name())
	}
	if !ed25519.Verify(pub.key, []byte(info.Fingerprint()), sig.sig[:]) {
		return fmt.Errorf("nix: signature verification failed for key %q", sig.Name())
	}
	return nil
}

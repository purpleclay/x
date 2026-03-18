package nix

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/purpleclay/x/nix/base32"
)

// HashType identifies a cryptographic hash algorithm used by Nix.
type HashType int

const (
	MD5 HashType = iota + 1
	SHA1
	SHA256
	SHA512
)

// Size returns the number of bytes produced by the hash algorithm.
func (t HashType) Size() int {
	switch t {
	case MD5:
		return 16
	case SHA1:
		return 20
	case SHA256:
		return 32
	case SHA512:
		return 64
	default:
		return 0
	}
}

// String returns the Nix canonical name for the hash algorithm
// (e.g. "sha256").
func (t HashType) String() string {
	switch t {
	case MD5:
		return "md5"
	case SHA1:
		return "sha1"
	case SHA256:
		return "sha256"
	case SHA512:
		return "sha512"
	default:
		return "unknown"
	}
}

// IsValid reports whether t is a recognised hash type.
func (t HashType) IsValid() bool {
	return t >= MD5 && t <= SHA512
}

// ParseHashType parses a Nix hash type name. Recognised values are
// "md5", "sha1", "sha256", and "sha512".
func ParseHashType(s string) (HashType, error) {
	switch s {
	case "md5":
		return MD5, nil
	case "sha1":
		return SHA1, nil
	case "sha256":
		return SHA256, nil
	case "sha512":
		return SHA512, nil
	default:
		return 0, fmt.Errorf("nix: unknown hash type %q", s)
	}
}

// Hash is an opaque representation of a Nix cryptographic hash.
// The zero value is not a valid hash; use [NewHash] or [ParseHash]
// to construct one.
type Hash struct {
	typ HashType
	// sized for SHA-512; unused bytes are always zero
	data [64]byte
}

// NewHash constructs a Hash from a type and raw bytes. It panics if
// the length of b does not match the size required by typ.
func NewHash(typ HashType, b []byte) Hash {
	if !typ.IsValid() {
		panic("nix: invalid hash type")
	}
	if len(b) != typ.Size() {
		panic(fmt.Sprintf("nix: hash size mismatch: %s requires %d bytes, got %d", typ, typ.Size(), len(b)))
	}
	var h Hash
	h.typ = typ
	copy(h.data[:], b)
	return h
}

// Type returns the hash algorithm used to produce this hash.
func (h Hash) Type() HashType {
	return h.typ
}

// IsZero reports whether h is the zero value (i.e. not a valid hash).
func (h Hash) IsZero() bool {
	return h.typ == 0
}

// Equal reports whether h and other represent the same hash value.
func (h Hash) Equal(other Hash) bool {
	if h.typ != other.typ {
		return false
	}
	for i := range h.typ.Size() {
		if h.data[i] != other.data[i] {
			return false
		}
	}
	return true
}

// Bytes returns a copy of the raw hash bytes.
func (h Hash) Bytes() []byte {
	b := make([]byte, h.typ.Size())
	copy(b, h.data[:h.typ.Size()])
	return b
}

// String returns the hash in prefixed nix32 format (e.g. "sha256:abc...").
func (h Hash) String() string {
	return h.Base32()
}

// RawHex returns the hash encoded as a lowercase hexadecimal string,
// without a type prefix.
func (h Hash) RawHex() string {
	return hex.EncodeToString(h.data[:h.typ.Size()])
}

// Hex returns the hash in prefixed hex format: "type:hex".
func (h Hash) Hex() string {
	return h.typ.String() + ":" + h.RawHex()
}

// RawBase32 returns the hash encoded as a nix base32 string,
// without a type prefix.
func (h Hash) RawBase32() string {
	return base32.StdEncoding.EncodeToString(h.data[:h.typ.Size()])
}

// Base32 returns the hash in prefixed nix32 format: "type:nix32".
func (h Hash) Base32() string {
	return h.typ.String() + ":" + h.RawBase32()
}

// RawBase64 returns the hash encoded as a standard base64 string
// (with padding), without a type prefix.
func (h Hash) RawBase64() string {
	return base64.StdEncoding.EncodeToString(h.data[:h.typ.Size()])
}

// Base64 returns the hash in prefixed base64 format: "type:base64".
func (h Hash) Base64() string {
	return h.typ.String() + ":" + h.RawBase64()
}

// SRI returns the hash in Subresource Integrity format: "type-base64".
func (h Hash) SRI() string {
	typ := h.typ.String()
	n := h.typ.Size()
	buf := make([]byte, len(typ)+1+base64.StdEncoding.EncodedLen(n))
	i := copy(buf, typ)
	buf[i] = '-'
	base64.StdEncoding.Encode(buf[i+1:], h.data[:n])
	return string(buf)
}

// MarshalText implements [encoding.TextMarshaler] using SRI format.
func (h Hash) MarshalText() ([]byte, error) {
	return []byte(h.SRI()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It accepts any
// format understood by [ParseHash].
func (h *Hash) UnmarshalText(text []byte) error {
	parsed, err := ParseHash(string(text))
	if err != nil {
		return err
	}
	*h = parsed
	return nil
}

// ParseHashError is returned by [ParseHash] when a hash string cannot be parsed.
type ParseHashError struct {
	// Input is the string that could not be parsed.
	Input string
	// Err is the underlying cause of the failure.
	Err error
}

func (e *ParseHashError) Error() string {
	return fmt.Sprintf("nix: cannot parse hash %q: %s", e.Input, e.Err)
}

func (e *ParseHashError) Unwrap() error {
	return e.Err
}

// ParseHash parses a hash string in any of the formats Nix uses:
//
//   - Prefixed hex:    type:hex
//   - Prefixed nix32:  type:nix32
//   - Prefixed base64: type:base64
//   - SRI:             type-base64
//
// The encoding is inferred from the length of the encoded part.
// Parse failures return a [*ParseHashError].
func ParseHash(s string) (Hash, error) {
	parseErr := func(cause error) error {
		return &ParseHashError{Input: s, Err: cause}
	}

	if s == "" {
		return Hash{}, parseErr(errors.New("empty string"))
	}

	// stack-allocated decode buffer, sized for SHA-512
	var buf [64]byte

	// SRI format uses a hyphen separator: type-base64.
	// Hash type names never contain a hyphen, so the first hyphen is the separator.
	if idx := strings.Index(s, "-"); idx > 0 {
		typ, err := ParseHashType(s[:idx])
		if err == nil {
			n, err := base64.StdEncoding.Decode(buf[:], []byte(s[idx+1:]))
			if err != nil {
				return Hash{}, parseErr(fmt.Errorf("invalid base64: %w", err))
			}
			if n != typ.Size() {
				return Hash{}, parseErr(fmt.Errorf("%s requires %d bytes, got %d", typ, typ.Size(), n))
			}
			return NewHash(typ, buf[:n]), nil
		}
	}

	// Prefixed format: type:encoded
	before, after, ok := strings.Cut(s, ":")
	if !ok {
		return Hash{}, parseErr(errors.New("missing type prefix"))
	}

	typ, err := ParseHashType(before)
	if err != nil {
		return Hash{}, parseErr(err)
	}

	size := typ.Size()

	switch len(after) {
	case size * 2:
		n, err := hex.Decode(buf[:size], []byte(after))
		if err != nil {
			return Hash{}, parseErr(fmt.Errorf("invalid hex: %w", err))
		}
		return NewHash(typ, buf[:n]), nil
	case base32.StdEncoding.EncodedLen(size):
		n, err := base32.StdEncoding.Decode(buf[:], []byte(after))
		if err != nil {
			return Hash{}, parseErr(fmt.Errorf("invalid nix32: %w", err))
		}
		return NewHash(typ, buf[:n]), nil
	case base64.StdEncoding.EncodedLen(size):
		n, err := base64.StdEncoding.Decode(buf[:], []byte(after))
		if err != nil {
			return Hash{}, parseErr(fmt.Errorf("invalid base64: %w", err))
		}
		return NewHash(typ, buf[:n]), nil
	default:
		return Hash{}, parseErr(fmt.Errorf("cannot determine encoding (length %d)", len(after)))
	}
}

// CompressHash compresses src into dst by cyclic XOR, truncating the hash
// to len(dst) bytes. It is used during Nix store path derivation to reduce
// a SHA-256 digest to 20 bytes before encoding as nix base32.
func CompressHash(dst, src []byte) {
	for i := range dst {
		dst[i] = 0
	}
	for i, b := range src {
		dst[i%len(dst)] ^= b
	}
}

package nix

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/purpleclay/x/nix/base32"
)

const digestLen = 32

// StoreDirectory is the root directory of a Nix store.
type StoreDirectory string

// DefaultStoreDirectory is the conventional Nix store root.
const DefaultStoreDirectory StoreDirectory = "/nix/store"

// Join joins path elements to the store directory.
func (d StoreDirectory) Join(elem ...string) string {
	parts := make([]string, 0, 1+len(elem))
	parts = append(parts, string(d))
	parts = append(parts, elem...)
	return path.Join(parts...)
}

// Object constructs a StorePath from a base name of the form
// <32-char-nix-base32-digest>-<name>, validating the digest and name.
//
// The name must include the digest — it is the full base component of the
// store path, not just the human-readable package name:
//
//	sp, _ := nix.DefaultStoreDirectory.Object("v6szabc66jn5r6vq8p7g7m5b2ks1xwif-hello-2.12.1")
//	// sp == "/nix/store/v6szabc66jn5r6vq8p7g7m5b2ks1xwif-hello-2.12.1"
func (d StoreDirectory) Object(name string) (StorePath, error) {
	input := string(d) + "/" + name
	if err := validateBase(name, input); err != nil {
		return "", err
	}
	return StorePath(input), nil
}

// ParsePath parses an absolute path rooted within this store directory.
// It returns the StorePath for the store object and any trailing sub-path.
//
//	sp, sub, _ := nix.DefaultStoreDirectory.ParsePath(
//		"/nix/store/v6szabc66jn5r6vq8p7g7m5b2ks1xwif-hello-2.12.1/bin/hello",
//	)
//	// sp  == "/nix/store/v6szabc66jn5r6vq8p7g7m5b2ks1xwif-hello-2.12.1"
//	// sub == "/bin/hello"
func (d StoreDirectory) ParsePath(absPath string) (StorePath, string, error) {
	prefix := string(d) + "/"
	if !strings.HasPrefix(absPath, prefix) {
		return "", "", &ParseStorePathError{
			Input: absPath,
			Err:   fmt.Errorf("path is not within store directory %q", string(d)),
		}
	}
	rest := absPath[len(prefix):]
	var baseName, sub string
	if idx := strings.IndexByte(rest, '/'); idx >= 0 {
		baseName = rest[:idx]
		sub = rest[idx:]
	} else {
		baseName = rest
	}
	if err := validateBase(baseName, absPath); err != nil {
		return "", "", err
	}
	return StorePath(string(d) + "/" + baseName), sub, nil
}

// StorePath is a validated Nix store path of the form
// <store-dir>/<32-char-nix-base32-digest>-<name>.
type StorePath string

// ParseStorePath parses and validates a Nix store path using DefaultStoreDirectory.
// It rejects paths with sub-path components; use [StoreDirectory.ParsePath] for those.
//
//	sp, _ := nix.ParseStorePath("/nix/store/v6szabc66jn5r6vq8p7g7m5b2ks1xwif-hello-2.12.1")
//	// sp.Digest() == "v6szabc66jn5r6vq8p7g7m5b2ks1xwif"
//	// sp.Name()   == "hello-2.12.1"
func ParseStorePath(s string) (StorePath, error) {
	prefix := string(DefaultStoreDirectory) + "/"
	if !strings.HasPrefix(s, prefix) {
		return "", &ParseStorePathError{
			Input: s,
			Err:   fmt.Errorf("path is not within store directory %q", string(DefaultStoreDirectory)),
		}
	}
	baseName := s[len(prefix):]
	if strings.IndexByte(baseName, '/') >= 0 {
		return "", &ParseStorePathError{Input: s, Err: errors.New("path contains sub-path components")}
	}
	if err := validateBase(baseName, s); err != nil {
		return "", err
	}
	return StorePath(s), nil
}

// Base returns the last element of the store path: <digest>-<name>.
//
//	sp := nix.StorePath("/nix/store/v6szabc66jn5r6vq8p7g7m5b2ks1xwif-hello-2.12.1")
//	// sp.Base() == "v6szabc66jn5r6vq8p7g7m5b2ks1xwif-hello-2.12.1"
func (p StorePath) Base() string {
	s := string(p)
	if idx := strings.LastIndexByte(s, '/'); idx >= 0 {
		return s[idx+1:]
	}
	return s
}

// Dir returns the store directory component of the path.
//
//	sp := nix.StorePath("/nix/store/v6szabc66jn5r6vq8p7g7m5b2ks1xwif-hello-2.12.1")
//	// sp.Dir() == "/nix/store"
func (p StorePath) Dir() StoreDirectory {
	s := string(p)
	if idx := strings.LastIndexByte(s, '/'); idx >= 0 {
		return StoreDirectory(s[:idx])
	}
	return ""
}

// Digest returns the 32-character nix base32 digest portion of the base name.
//
//	sp := nix.StorePath("/nix/store/v6szabc66jn5r6vq8p7g7m5b2ks1xwif-hello-2.12.1")
//	// sp.Digest() == "v6szabc66jn5r6vq8p7g7m5b2ks1xwif"
func (p StorePath) Digest() string {
	return p.Base()[:digestLen]
}

// Name returns the human-readable name portion of the store path (after the digest and dash).
//
//	sp := nix.StorePath("/nix/store/v6szabc66jn5r6vq8p7g7m5b2ks1xwif-hello-2.12.1")
//	// sp.Name() == "hello-2.12.1"
func (p StorePath) Name() string {
	return p.Base()[digestLen+1:]
}

// IsDerivation reports whether the store path refers to a derivation (.drv file).
//
//	sp := nix.StorePath("/nix/store/v6szabc66jn5r6vq8p7g7m5b2ks1xwif-hello-2.12.1.drv")
//	// sp.IsDerivation() == true
func (p StorePath) IsDerivation() bool {
	return strings.HasSuffix(string(p), ".drv")
}

// String returns the full store path string.
func (p StorePath) String() string {
	return string(p)
}

// MarshalText implements [encoding.TextMarshaler].
func (p StorePath) MarshalText() ([]byte, error) {
	return []byte(p), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (p *StorePath) UnmarshalText(text []byte) error {
	parsed, err := ParseStorePath(string(text))
	if err != nil {
		return err
	}
	*p = parsed
	return nil
}

// ParseStorePathError is returned by [ParseStorePath] when a store path
// string cannot be parsed.
type ParseStorePathError struct {
	// Input is the string that could not be parsed.
	Input string
	// Err is the underlying cause of the failure.
	Err error
}

func (e *ParseStorePathError) Error() string {
	return fmt.Sprintf("nix: cannot parse store path %q: %s", e.Input, e.Err)
}

func (e *ParseStorePathError) Unwrap() error {
	return e.Err
}

// validateBase validates the base name component (<digest>-<name>) of a store path.
// input is the original string passed by the caller, used for error reporting.
func validateBase(baseName, input string) error {
	parseErr := func(cause error) error {
		return &ParseStorePathError{Input: input, Err: cause}
	}
	if len(baseName) < digestLen+2 {
		return parseErr(fmt.Errorf("digest must be %d characters", digestLen))
	}
	digest := baseName[:digestLen]
	var buf [20]byte
	if _, err := base32.StdEncoding.Decode(buf[:], []byte(digest)); err != nil {
		return parseErr(fmt.Errorf("invalid digest: %w", err))
	}
	if baseName[digestLen] != '-' {
		return parseErr(errors.New("digest must be followed by '-'"))
	}
	return nil
}

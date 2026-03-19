package nix

import (
	"fmt"
	"strconv"
	"strings"
)

// NARInfoMIMEType is the MIME type for .narinfo resources.
const NARInfoMIMEType = "text/x-nix-narinfo"

// NARInfoName returns the resource path for a store path's narinfo file.
// It is the 32-character digest followed by ".narinfo".
//
//	sp, _ := nix.ParseStorePath("/nix/store/s66mzxpvicwk07gjbjfw9izjfa797vsw-hello-2.12.1")
//	// nix.NARInfoName(sp) == "s66mzxpvicwk07gjbjfw9izjfa797vsw.narinfo"
func NARInfoName(sp StorePath) string {
	return sp.Digest() + ".narinfo"
}

// Compression identifies the compression algorithm applied to a NAR archive.
type Compression int

const (
	// Bzip2 is the bzip2 compression algorithm, identified as "bzip2" in narinfo
	// files. It is the Nix default when the Compression field is absent.
	// See https://sourceware.org/bzip2/.
	Bzip2 Compression = iota + 1

	// Gzip is the gzip compression algorithm, identified as "gzip" in narinfo files.
	// See https://www.gnu.org/software/gzip/.
	Gzip

	// XZ is the xz/LZMA2 compression algorithm, identified as "xz" in narinfo files.
	// See https://tukaani.org/xz/.
	XZ

	// Zstd is the Zstandard compression algorithm, identified as "zstd" in narinfo files.
	// See https://github.com/facebook/zstd.
	Zstd

	// Brotli is the Brotli compression algorithm, identified as "br" in narinfo files.
	// See https://github.com/google/brotli.
	Brotli

	// NoCompression indicates the NAR archive is stored uncompressed, identified
	// as "none" in narinfo files.
	NoCompression
)

// String returns the Nix canonical name for the compression algorithm.
//
//	// nix.Zstd.String() == "zstd"
func (c Compression) String() string {
	switch c {
	case Bzip2:
		return "bzip2"
	case Gzip:
		return "gzip"
	case XZ:
		return "xz"
	case Zstd:
		return "zstd"
	case Brotli:
		return "br"
	case NoCompression:
		return "none"
	default:
		return "unknown"
	}
}

// ParseCompression parses a Nix compression algorithm name. Recognised values
// are "bzip2", "gzip", "xz", "zstd", "br", and "none".
//
//	c, _ := nix.ParseCompression("zstd")
//	// c == nix.Zstd
func ParseCompression(s string) (Compression, error) {
	switch s {
	case "bzip2":
		return Bzip2, nil
	case "gzip":
		return Gzip, nil
	case "xz":
		return XZ, nil
	case "zstd":
		return Zstd, nil
	case "br":
		return Brotli, nil
	case "none":
		return NoCompression, nil
	default:
		return 0, fmt.Errorf("nix: unknown compression %q", s)
	}
}

// NARInfo describes a Nix store object within a binary cache. Every store
// object is described by a .narinfo file that Nix fetches before downloading
// the NAR archive itself. See the Nix binary cache protocol documentation at
// https://nixos.org/manual/nix/stable/protocols/http-binary-cache-store.
//
//	ni := nix.NARInfo{
//		StorePath:   "/nix/store/s66mzxpvicwk07gjbjfw9izjfa797vsw-hello-2.12.1",
//		URL:         "nar/0vgdasb2zswyxqb2g6q3rabg1rj9ijj21gzgwvkz8bbyb5nji7mk.nar.zst",
//		Compression: nix.Zstd,
//		NARSize:     54416,
//	}
type NARInfo struct {
	// StorePath is the full Nix store path of the object.
	StorePath StorePath

	// URL is the path to the NAR archive relative to the binary cache root.
	URL string

	// Compression is the algorithm used to compress the NAR archive.
	// Defaults to [Bzip2] when absent in the narinfo file.
	Compression Compression

	// FileHash is the hash of the compressed NAR file (typically SRI format).
	FileHash Hash

	// FileSize is the size of the compressed NAR file in bytes.
	FileSize int64

	// NARHash is the hash of the uncompressed NAR archive (typically SRI format).
	NARHash Hash

	// NARSize is the size of the uncompressed NAR archive in bytes.
	NARSize int64

	// References holds the base names (<digest>-<name>) of store objects
	// that this object depends on at runtime. Empty means no dependencies.
	References []string

	// Deriver is the base name of the .drv file that produced this object.
	// Empty if not recorded in the narinfo.
	Deriver string

	// Sig holds the raw "name:base64" signatures attached to this narinfo.
	// Multiple signatures may be present.
	Sig []string
}

// MarshalText implements [encoding.TextMarshaler]. It serialises NARInfo in
// the Nix Key: Value line format that Nix's narinfo parser expects. A zero
// Compression is marshaled as [Bzip2]. Deriver is omitted when empty. Each
// signature produces a separate Sig line.
//
//	ni := nix.NARInfo{
//		StorePath:   "/nix/store/s66mzxpvicwk07gjbjfw9izjfa797vsw-hello-2.12.1",
//		URL:         "nar/0vgdasb2zswyxqb2g6q3rabg1rj9ijj21gzgwvkz8bbyb5nji7mk.nar.zst",
//		Compression: nix.Zstd,
//		NARSize:     54416,
//	}
//	text, _ := ni.MarshalText()
//	// "StorePath: /nix/store/s66mzxpvicwk07gjbjfw9izjfa797vsw-hello-2.12.1\nURL: nar/...nar.zst\n..."
func (n NARInfo) MarshalText() ([]byte, error) {
	compression := n.Compression
	if compression == 0 {
		compression = Bzip2
	}
	var b []byte
	b = fmt.Appendf(b, "StorePath: %s\n", n.StorePath)
	b = fmt.Appendf(b, "URL: %s\n", n.URL)
	b = fmt.Appendf(b, "Compression: %s\n", compression)
	b = fmt.Appendf(b, "FileHash: %s\n", n.FileHash.SRI())
	b = fmt.Appendf(b, "FileSize: %d\n", n.FileSize)
	b = fmt.Appendf(b, "NARHash: %s\n", n.NARHash.SRI())
	b = fmt.Appendf(b, "NARSize: %d\n", n.NARSize)
	b = fmt.Appendf(b, "References: %s\n", strings.Join(n.References, " "))
	if n.Deriver != "" {
		b = fmt.Appendf(b, "Deriver: %s\n", n.Deriver)
	}
	for _, sig := range n.Sig {
		b = fmt.Appendf(b, "Sig: %s\n", sig)
	}
	return b, nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It parses NARInfo from
// the Nix Key: Value line format. A missing Compression field defaults to
// [Bzip2]. Multiple Sig lines are accumulated in order. Unknown keys are
// silently ignored for forward compatibility.
//
//	var ni nix.NARInfo
//	_ = ni.UnmarshalText([]byte("StorePath: /nix/store/s66mzxpvicwk07gjbjfw9izjfa797vsw-hello-2.12.1\n..."))
//	// ni.StorePath.Digest() == "s66mzxpvicwk07gjbjfw9izjfa797vsw"
func (n *NARInfo) UnmarshalText(text []byte) error {
	for line := range strings.SplitSeq(string(text), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		key, value, ok := strings.Cut(line, ": ")
		if !ok {
			continue
		}
		switch key {
		case "StorePath":
			sp, err := ParseStorePath(value)
			if err != nil {
				return fmt.Errorf("nix: NARInfo: invalid StorePath: %w", err)
			}
			n.StorePath = sp
		case "URL":
			n.URL = value
		case "Compression":
			c, err := ParseCompression(value)
			if err != nil {
				return fmt.Errorf("nix: NARInfo: invalid Compression %q", value)
			}
			n.Compression = c
		case "FileHash":
			h, err := ParseHash(value)
			if err != nil {
				return fmt.Errorf("nix: NARInfo: invalid FileHash: %w", err)
			}
			n.FileHash = h
		case "FileSize":
			sz, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("nix: NARInfo: invalid FileSize %q", value)
			}
			n.FileSize = sz
		case "NARHash":
			h, err := ParseHash(value)
			if err != nil {
				return fmt.Errorf("nix: NARInfo: invalid NARHash: %w", err)
			}
			n.NARHash = h
		case "NARSize":
			sz, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("nix: NARInfo: invalid NARSize %q", value)
			}
			n.NARSize = sz
		case "References":
			n.References = strings.Fields(value)
		case "Deriver":
			n.Deriver = value
		case "Sig":
			n.Sig = append(n.Sig, value)
		}
	}
	if n.Compression == 0 {
		n.Compression = Bzip2
	}
	return nil
}

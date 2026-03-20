package nix_test

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/purpleclay/x/nix"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type hashVector struct {
	name  string
	typ   nix.HashType
	hex   string
	nix32 string
}

// Vectors based on: encoding_test.go
var vectors = []hashVector{
	{
		name:  "MD5",
		typ:   nix.MD5,
		hex:   "47b2d8f260c2d48116044bc43fe3de0f",
		nix32: "0gvvikzi2b0hb83m62c3rdicj7",
	},
	{
		name:  "SHA1",
		typ:   nix.SHA1,
		hex:   "1f74d74729abdc08f4f84e8f7f8c808c8ed92ee5",
		nix32: "wlpdk3lch267z3sfz3s0ip5b553xfx0z",
	},
	{
		name:  "SHA256",
		typ:   nix.SHA256,
		hex:   "a315ab26a0c4829321730c44a26f4497f7da0631402669caa4e24bdcd9db7c87",
		nix32: "11vwvgcxqjz2lk56j9j0643dmxwp8ips4i0cfchr70n4l0kan5d3",
	},
	{
		name:  "SHA512",
		typ:   nix.SHA512,
		hex:   "296a445bfa5e1990af299ec74582468ab7a77e495861691ff79ab21234e514b64fc72b294d7305ecdd0febaa13b1bc1a3f359a711bb93dfb2b82804c64354dab",
		nix32: "2mlsdb49j084azv7nwinwcs6lzimg5i2fmfn3yxxh2p6k995g3lzdhlwls15clsywgnjqaq95zagdwa8s14biwy56pr06ayz9dl8si9",
	},
}

func rawBytes(t *testing.T, v hashVector) []byte {
	t.Helper()
	b, err := hex.DecodeString(v.hex)
	require.NoError(t, err)
	return b
}

func rawBytesBase64(t *testing.T, v hashVector) string {
	t.Helper()
	return base64.StdEncoding.EncodeToString(rawBytes(t, v))
}

func TestHashType(t *testing.T) {
	tests := []struct {
		name  string
		typ   nix.HashType
		str   string
		size  int
		valid bool
	}{
		{"MD5", nix.MD5, "md5", 16, true},
		{"SHA1", nix.SHA1, "sha1", 20, true},
		{"SHA256", nix.SHA256, "sha256", 32, true},
		{"SHA512", nix.SHA512, "sha512", 64, true},
		{"Unknown", nix.HashType(0), "unknown", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.str, tt.typ.String())
			assert.Equal(t, tt.size, tt.typ.Size())
			assert.Equal(t, tt.valid, tt.typ.IsValid())
		})
	}
}

func TestParseHashType(t *testing.T) {
	for _, tt := range vectors {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nix.ParseHashType(tt.typ.String())
			require.NoError(t, err)
			assert.Equal(t, tt.typ, got)
		})
	}
}

func TestParseHashTypeUnknown(t *testing.T) {
	_, err := nix.ParseHashType("blake3")
	assert.ErrorContains(t, err, "blake3")
}

func TestHashOutputMethods(t *testing.T) {
	for _, tt := range vectors {
		t.Run(tt.name, func(t *testing.T) {
			raw := rawBytes(t, tt)
			b64 := rawBytesBase64(t, tt)

			h := nix.NewHash(tt.typ, raw)

			assert.Equal(t, tt.hex, h.RawHex())
			assert.Equal(t, tt.typ.String()+":"+tt.hex, h.Hex())
			assert.Equal(t, tt.nix32, h.RawBase32())
			assert.Equal(t, tt.typ.String()+":"+tt.nix32, h.Base32())
			assert.Equal(t, b64, h.RawBase64())
			assert.Equal(t, tt.typ.String()+":"+b64, h.Base64())
			assert.Equal(t, tt.typ.String()+"-"+b64, h.SRI())
			assert.Equal(t, h.Base32(), h.String())
		})
	}
}

func TestParseHash(t *testing.T) {
	for _, tt := range vectors {
		raw := rawBytes(t, tt)
		b64 := rawBytesBase64(t, tt)

		t.Run(tt.name+"/Hex", func(t *testing.T) {
			h, err := nix.ParseHash(tt.typ.String() + ":" + tt.hex)
			require.NoError(t, err)
			assert.Equal(t, tt.typ, h.Type())
			assert.Equal(t, raw, h.Bytes())
		})

		t.Run(tt.name+"/Nix32", func(t *testing.T) {
			h, err := nix.ParseHash(tt.typ.String() + ":" + tt.nix32)
			require.NoError(t, err)
			assert.Equal(t, tt.typ, h.Type())
			assert.Equal(t, raw, h.Bytes())
		})

		t.Run(tt.name+"/Base64", func(t *testing.T) {
			h, err := nix.ParseHash(tt.typ.String() + ":" + b64)
			require.NoError(t, err)
			assert.Equal(t, tt.typ, h.Type())
			assert.Equal(t, raw, h.Bytes())
		})

		t.Run(tt.name+"/SRI", func(t *testing.T) {
			h, err := nix.ParseHash(tt.typ.String() + "-" + b64)
			require.NoError(t, err)
			assert.Equal(t, tt.typ, h.Type())
			assert.Equal(t, raw, h.Bytes())
		})
	}
}

func TestParseHashError(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Empty", ""},
		{"MissingPrefix", "abc123"},
		{"UnknownType", "blake3:" + strings.Repeat("a", 64)},
		{"WrongLength", "sha256:abc"},
		{"InvalidHex", "sha256:" + strings.Repeat("z", 64)},
		{"InvalidNix32", "sha256:" + strings.Repeat("e", 52)},
		{"InvalidBase64", "sha256:" + strings.Repeat("!", 44)},
		{"InvalidSRIBase64", "sha256-!!!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := nix.ParseHash(tt.input)
			require.Error(t, err)

			var parseErr *nix.ParseHashError
			require.ErrorAs(t, err, &parseErr)
			assert.Equal(t, tt.input, parseErr.Input)
			assert.NotNil(t, parseErr.Err)
		})
	}
}

func TestHashIsZero(t *testing.T) {
	var zero nix.Hash
	assert.True(t, zero.IsZero())

	raw := rawBytes(t, vectors[0])
	h := nix.NewHash(vectors[0].typ, raw)
	assert.False(t, h.IsZero())
}

func TestHashEqual(t *testing.T) {
	sha256Raw := rawBytes(t, vectors[2])

	mutated := make([]byte, len(sha256Raw))
	copy(mutated, sha256Raw)
	mutated[0] ^= 0xff

	tests := []struct {
		name  string
		a     nix.Hash
		b     nix.Hash
		equal bool
	}{
		{
			name:  "SameTypeSameBytes",
			a:     nix.NewHash(nix.SHA256, sha256Raw),
			b:     nix.NewHash(nix.SHA256, sha256Raw),
			equal: true,
		},
		{
			name:  "SameTypeDifferentBytes",
			a:     nix.NewHash(nix.SHA256, sha256Raw),
			b:     nix.NewHash(nix.SHA256, mutated),
			equal: false,
		},
		{
			name:  "DifferentType",
			a:     nix.NewHash(nix.MD5, rawBytes(t, vectors[0])),
			b:     nix.NewHash(nix.SHA1, rawBytes(t, vectors[1])),
			equal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.equal, tt.a.Equal(tt.b))
		})
	}
}

func TestHashTextMarshaler(t *testing.T) {
	for _, tt := range vectors {
		t.Run(tt.name, func(t *testing.T) {
			original := nix.NewHash(tt.typ, rawBytes(t, tt))

			text, err := original.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, original.SRI(), string(text))

			var restored nix.Hash
			require.NoError(t, restored.UnmarshalText(text))
			assert.True(t, original.Equal(restored))
		})
	}
}

func TestCompressHash(t *testing.T) {
	t.Run("IdentityWhenSameLength", func(t *testing.T) {
		src := []byte{0x01, 0x02, 0x03, 0x04}
		dst := make([]byte, len(src))
		nix.CompressHash(dst, src)
		assert.Equal(t, src, dst)
	})

	t.Run("CyclicXOR", func(t *testing.T) {
		src := []byte{0x01, 0x02, 0x03, 0x04}
		dst := make([]byte, 2)
		nix.CompressHash(dst, src)
		// dst[0] = 0x01 ^ 0x03 = 0x02
		// dst[1] = 0x02 ^ 0x04 = 0x06
		assert.Equal(t, []byte{0x02, 0x06}, dst)
	})

	t.Run("ZeroesDestinationFirst", func(t *testing.T) {
		src := []byte{0x00, 0x00}
		dst := []byte{0xff, 0xff}
		nix.CompressHash(dst, src)
		assert.Equal(t, []byte{0x00, 0x00}, dst)
	})
}

var sha256Vector = vectors[2] // SHA256

func BenchmarkParseHash(b *testing.B) {
	raw, _ := hex.DecodeString(sha256Vector.hex)
	b64 := base64.StdEncoding.EncodeToString(raw)

	inputs := map[string]string{
		"Hex":    sha256Vector.typ.String() + ":" + sha256Vector.hex,
		"Nix32":  sha256Vector.typ.String() + ":" + sha256Vector.nix32,
		"Base64": sha256Vector.typ.String() + ":" + b64,
		"SRI":    sha256Vector.typ.String() + "-" + b64,
	}

	for name, input := range inputs {
		b.Run(name, func(b *testing.B) {
			for b.Loop() {
				nix.ParseHash(input)
			}
		})
	}
}

func BenchmarkHashOutput(b *testing.B) {
	raw, _ := hex.DecodeString(sha256Vector.hex)
	h := nix.NewHash(sha256Vector.typ, raw)

	b.Run("RawHex", func(b *testing.B) {
		for b.Loop() {
			h.RawHex()
		}
	})

	b.Run("RawBase32", func(b *testing.B) {
		for b.Loop() {
			h.RawBase32()
		}
	})

	b.Run("RawBase64", func(b *testing.B) {
		for b.Loop() {
			h.RawBase64()
		}
	})

	b.Run("SRI", func(b *testing.B) {
		for b.Loop() {
			h.SRI()
		}
	})
}

func BenchmarkCompressHash(b *testing.B) {
	raw, _ := hex.DecodeString(sha256Vector.hex)
	dst := make([]byte, 20)

	for b.Loop() {
		nix.CompressHash(dst, raw)
	}
}

func FuzzParseHash(f *testing.F) {
	// One valid nix32 seed per hash type.
	f.Add("sha256:11vwvgcxqjz2lk56j9j0643dmxwp8ips4i0cfchr70n4l0kan5d3")
	f.Add("sha256:a315ab26a0c4829321730c44a26f4497f7da0631402669caa4e24bdcd9db7c87")
	// SRI format (sha256 only in practice).
	f.Add("sha256-oxWrJqDEgpMhcwxEom9El/faBjFAJmnKpOJL3NnbfIc=")
	// Shorter hash types — exercises the length-dispatch in ParseHash.
	f.Add("sha1:wlpdk3lch267z3sfz3s0ip5b553xfx0z")
	f.Add("md5:0gvvikzi2b0hb83m62c3rdicj7")
	f.Add("sha512:2mlsdb49j084azv7nwinwcs6lzimg5i2fmfn3yxxh2p6k995g3lzdhlwls15clsywgnjqaq95zagdwa8s14biwy56pr06ayz9dl8si9")

	f.Fuzz(func(_ *testing.T, s string) {
		nix.ParseHash(s)
	})
}

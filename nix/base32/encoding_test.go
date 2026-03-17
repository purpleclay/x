package base32_test

import (
	"encoding/hex"
	"testing"

	"github.com/purpleclay/x/nix/base32"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testVector struct {
	name   string
	hex    string
	base32 string
}

var (
	emptyVector    = testVector{name: "Empty", hex: "", base32: ""}
	singleByteZero = testVector{name: "SingleByteZero", hex: "00", base32: "00"}
	singleByteFF   = testVector{name: "SingleByteFF", hex: "ff", base32: "7z"}

	// Vectors taken from: github.com/kolloch/nix-base32 and verified against the Nix C++ impl
	md5Nar         = testVector{name: "MD5Nar", hex: "47b2d8f260c2d48116044bc43fe3de0f", base32: "0gvvikzi2b0hb83m62c3rdicj7"}
	sha1Nar        = testVector{name: "SHA1Nar", hex: "1f74d74729abdc08f4f84e8f7f8c808c8ed92ee5", base32: "wlpdk3lch267z3sfz3s0ip5b553xfx0z"}
	sha256Nar      = testVector{name: "SHA256Nar", hex: "a315ab26a0c4829321730c44a26f4497f7da0631402669caa4e24bdcd9db7c87", base32: "11vwvgcxqjz2lk56j9j0643dmxwp8ips4i0cfchr70n4l0kan5d3"}
	sha512Nar      = testVector{name: "SHA512Nar", hex: "296a445bfa5e1990af299ec74582468ab7a77e495861691ff79ab21234e514b64fc72b294d7305ecdd0febaa13b1bc1a3f359a711bb93dfb2b82804c64354dab", base32: "2mlsdb49j084azv7nwinwcs6lzimg5i2fmfn3yxxh2p6k995g3lzdhlwls15clsywgnjqaq95zagdwa8s14biwy56pr06ayz9dl8si9"}
	sha256CacheNar = testVector{name: "SHA256CacheNar", hex: "99a2da84cec54d17325bcee0a079669c1b15eb7ead32246514b75b97862f1e00", base32: "000y5y39fnxp2ijj8cmdgvmia6wwcrws1q6fbcr1fkf5rs2dm8lr"}

	// Vector taken from: https://nixos.org/manual/nix/stable/command-ref/nix-hash
	// 	nix-hash --type sha1 --base32 test/  =>  nvd61k9nalji1zl9rrdfmsmvyyjqpzg4
	// 	nix-hash --type sha1 --base16 test/  =>  e4fd8ba5f7bbeaea5ace89fe10255536cd60dab6
	nixManualSHA1Dir = testVector{name: "NixHashSHA1", hex: "e4fd8ba5f7bbeaea5ace89fe10255536cd60dab6", base32: "nvd61k9nalji1zl9rrdfmsmvyyjqpzg4"}
)

var vectors = []testVector{
	emptyVector,
	singleByteZero,
	singleByteFF,
	md5Nar,
	sha1Nar,
	sha256Nar,
	sha512Nar,
	sha256CacheNar,
	nixManualSHA1Dir,
}

func TestEncode(t *testing.T) {
	for _, tt := range vectors {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := hex.DecodeString(tt.hex)
			require.NoError(t, err)

			got := base32.StdEncoding.EncodeToString(raw)
			assert.Equal(t, tt.base32, got)
		})
	}
}

func TestDecode(t *testing.T) {
	for _, tt := range vectors {
		t.Run(tt.name, func(t *testing.T) {
			decoded, err := base32.StdEncoding.DecodeString(tt.base32)
			require.NoError(t, err)
			assert.Equal(t, tt.hex, hex.EncodeToString(decoded))
		})
	}
}

func TestDecodeInvalid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		badByte  byte
		position int64
	}{
		{"ContainsLowercaseE", "egvvikzi2b0hb83m62c3rdicj7", 'e', 0},
		{"ContainsLowercaseO", "wlpdk3lch267z3sfo3s0ip5b553xfx0z", 'o', 16},
		{"ContainsLowercaseT", "nvd61k9nalji1zltrrdfmsmvyyjqpzg4", 't', 15},
		{"ContainsLowercaseU", "2mlsdb49j084azv7nwinwcs6lzimg5i2fmfn3yxxh2p6k995g3uzdhlwls15clsywgnjqaq95zagdwa8s14biwy56pr06ayz9dl8si9", 'u', 50},
		{"ContainsUppercase", "000y5y39fnxp2ijj8cmdgvmia6wwcrws1q6Fbcr1fkf5rs2dm8lr", 'F', 35},
		{"ContainsSpecialChar", "wlpdk3lch2@7z3sfz3s0ip5b553xfx0z", '@', 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := base32.StdEncoding.DecodeString(tt.input)
			require.ErrorAs(t, err, &base32.InvalidByteError{})
			var invalid base32.InvalidByteError
			require.ErrorAs(t, err, &invalid)
			assert.Equal(t, tt.badByte, invalid.Byte)
			assert.Equal(t, tt.position, invalid.Offset)
		})
	}
}

func TestEncodedLen(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int
		expected int
	}{
		{"MD5_16_Bytes", 16, 26},
		{"SHA1_20_Bytes", 20, 32},
		{"SHA256_32_Bytes", 32, 52},
		{"SHA512_64_Bytes", 64, 103},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, base32.StdEncoding.EncodedLen(tt.bytes))
		})
	}
}

func TestDecodedLen(t *testing.T) {
	tests := []struct {
		name     string
		chars    int
		expected int
	}{
		{"MD5_26_Chars", 26, 16},
		{"SHA1_32_Chars", 32, 20},
		{"SHA256_52_Chars", 52, 32},
		{"SHA512_103_Chars", 103, 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, base32.StdEncoding.DecodedLen(tt.chars))
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	raw, _ := hex.DecodeString(sha256Nar.hex)
	dst := make([]byte, base32.StdEncoding.EncodedLen(len(raw)))

	for b.Loop() {
		base32.StdEncoding.Encode(dst, raw)
	}
}

func BenchmarkDecode(b *testing.B) {
	src := []byte(sha256Nar.base32)
	dst := make([]byte, base32.StdEncoding.DecodedLen(len(src)))

	for b.Loop() {
		base32.StdEncoding.Decode(dst, src)
	}
}

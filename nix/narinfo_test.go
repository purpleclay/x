package nix_test

import (
	"testing"

	"github.com/purpleclay/x/nix"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Real hello-2.12.1 narinfo values from cache.nixos.org.
const (
	// testFakeSigB64 is a valid base64 encoding of 64 zero bytes used only
	// to test signature parsing mechanics, not cryptographic validity.
	testFakeSigB64 = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="

	testNARStorePath = "/nix/store/s66mzxpvicwk07gjbjfw9izjfa797vsw-hello-2.12.1"
	testNARURL       = "nar/1nhgq6wcggx0plpy4991h3ginj6hipsdslv4fd4zml1n707j26yq.nar.xz"
	testNARFileHash  = "sha256:1nhgq6wcggx0plpy4991h3ginj6hipsdslv4fd4zml1n707j26yq"
	testNARNARHash   = "sha256:0yzhigwjl6bws649vcs2asa4lbs8hg93hyix187gc7s7a74w5h80"
	testNARDeriver   = "ib3sh3pcz10wsmavxvkdbayhqivbghlq-hello-2.12.1.drv"
	testNARSig       = "cache.nixos.org-1:8ijECciSFzWHwwGVOIVYdp2fOIOJAfmzGHPQVwpktfTQJF6kMPPDre7UtFw3o+VqenC5P8RikKOAAfN7CvPEAg=="

	testNARInfoText = "StorePath: /nix/store/s66mzxpvicwk07gjbjfw9izjfa797vsw-hello-2.12.1\n" +
		"URL: nar/1nhgq6wcggx0plpy4991h3ginj6hipsdslv4fd4zml1n707j26yq.nar.xz\n" +
		"Compression: xz\n" +
		"FileHash: sha256:1nhgq6wcggx0plpy4991h3ginj6hipsdslv4fd4zml1n707j26yq\n" +
		"FileSize: 50088\n" +
		"NarHash: sha256:0yzhigwjl6bws649vcs2asa4lbs8hg93hyix187gc7s7a74w5h80\n" +
		"NarSize: 226488\n" +
		"References: 3n58xw4373jp0ljirf06d8077j15pc4j-glibc-2.37-8 s66mzxpvicwk07gjbjfw9izjfa797vsw-hello-2.12.1\n" +
		"Deriver: ib3sh3pcz10wsmavxvkdbayhqivbghlq-hello-2.12.1.drv\n" +
		"Sig: cache.nixos.org-1:8ijECciSFzWHwwGVOIVYdp2fOIOJAfmzGHPQVwpktfTQJF6kMPPDre7UtFw3o+VqenC5P8RikKOAAfN7CvPEAg==\n"
)

func TestNARInfoName(t *testing.T) {
	sp, err := nix.ParseStorePath(testNARStorePath)
	require.NoError(t, err)
	assert.Equal(t, "s66mzxpvicwk07gjbjfw9izjfa797vsw.narinfo", nix.NARInfoName(sp))
}

func TestCompressionString(t *testing.T) {
	tests := []struct {
		compression nix.Compression
		want        string
	}{
		{nix.Bzip2, "bzip2"},
		{nix.Gzip, "gzip"},
		{nix.XZ, "xz"},
		{nix.Zstd, "zstd"},
		{nix.Brotli, "br"},
		{nix.NoCompression, "none"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.compression.String())
		})
	}
}

func TestParseCompression(t *testing.T) {
	tests := []struct {
		input string
		want  nix.Compression
	}{
		{"bzip2", nix.Bzip2},
		{"gzip", nix.Gzip},
		{"xz", nix.XZ},
		{"zstd", nix.Zstd},
		{"br", nix.Brotli},
		{"none", nix.NoCompression},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			c, err := nix.ParseCompression(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.want, c)
		})
	}
}

func TestParseCompressionError(t *testing.T) {
	_, err := nix.ParseCompression("lz4")
	require.Error(t, err)
}

func TestNARInfoMarshalText(t *testing.T) {
	sp, err := nix.ParseStorePath(testNARStorePath)
	require.NoError(t, err)

	fileHash, err := nix.ParseHash(testNARFileHash)
	require.NoError(t, err)

	narHash, err := nix.ParseHash(testNARNARHash)
	require.NoError(t, err)

	sig, err := nix.ParseSignature(testNARSig)
	require.NoError(t, err)

	ni := nix.NARInfo{
		StorePath:   sp,
		URL:         testNARURL,
		Compression: nix.XZ,
		FileHash:    fileHash,
		FileSize:    50088,
		NARHash:     narHash,
		NARSize:     226488,
		References: []string{
			"3n58xw4373jp0ljirf06d8077j15pc4j-glibc-2.37-8",
			"s66mzxpvicwk07gjbjfw9izjfa797vsw-hello-2.12.1",
		},
		Deriver: testNARDeriver,
		Sig:     []*nix.Signature{sig},
	}

	text, err := ni.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, testNARInfoText, string(text))
}

func TestNARInfoMarshalTextDefaultCompression(t *testing.T) {
	ni := nix.NARInfo{
		StorePath: nix.StorePath(testNARStorePath),
	}

	text, err := ni.MarshalText()
	require.NoError(t, err)
	assert.Contains(t, string(text), "Compression: bzip2\n")
}

func TestNARInfoMarshalTextEmptyReferences(t *testing.T) {
	ni := nix.NARInfo{
		StorePath: nix.StorePath(testNARStorePath),
	}

	text, err := ni.MarshalText()
	require.NoError(t, err)
	assert.Contains(t, string(text), "References: \n")
}

func TestNARInfoUnmarshalText(t *testing.T) {
	var ni nix.NARInfo
	require.NoError(t, ni.UnmarshalText([]byte(testNARInfoText)))

	assert.Equal(t, nix.StorePath(testNARStorePath), ni.StorePath)
	assert.Equal(t, testNARURL, ni.URL)
	assert.Equal(t, nix.XZ, ni.Compression)
	assert.Equal(t, int64(50088), ni.FileSize)
	assert.Equal(t, int64(226488), ni.NARSize)
	assert.Equal(t, []string{
		"3n58xw4373jp0ljirf06d8077j15pc4j-glibc-2.37-8",
		"s66mzxpvicwk07gjbjfw9izjfa797vsw-hello-2.12.1",
	}, ni.References)
	assert.Equal(t, testNARDeriver, ni.Deriver)
	require.Len(t, ni.Sig, 1)
	assert.Equal(t, testNARSig, ni.Sig[0].String())
}

func TestNARInfoUnmarshalTextDefaultCompression(t *testing.T) {
	input := "StorePath: " + testNARStorePath + "\n" +
		"URL: " + testNARURL + "\n" +
		"NarSize: 226488\n"

	var ni nix.NARInfo
	require.NoError(t, ni.UnmarshalText([]byte(input)))
	assert.Equal(t, nix.Bzip2, ni.Compression)
}

func TestNARInfoUnmarshalTextEmptyReferences(t *testing.T) {
	input := "StorePath: " + testNARStorePath + "\n" +
		"References: \n"

	var ni nix.NARInfo
	require.NoError(t, ni.UnmarshalText([]byte(input)))
	assert.Empty(t, ni.References)
}

func TestNARInfoUnmarshalTextMultipleSigs(t *testing.T) {
	input := "StorePath: " + testNARStorePath + "\n" +
		"Sig: cache.nixos.org-1:" + testFakeSigB64 + "\n" +
		"Sig: mycache-1:" + testFakeSigB64 + "\n"

	var ni nix.NARInfo
	require.NoError(t, ni.UnmarshalText([]byte(input)))
	require.Len(t, ni.Sig, 2)
	assert.Equal(t, "cache.nixos.org-1", ni.Sig[0].Name())
	assert.Equal(t, "mycache-1", ni.Sig[1].Name())
}

func TestNARInfoUnknownKeysIgnored(t *testing.T) {
	input := "StorePath: " + testNARStorePath + "\n" +
		"UnknownField: somevalue\n" +
		"NarSize: 226488\n"

	var ni nix.NARInfo
	require.NoError(t, ni.UnmarshalText([]byte(input)))
	assert.Equal(t, int64(226488), ni.NARSize)
}

func TestNARInfoUnmarshalTextError(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"InvalidStorePath", "StorePath: /usr/local/not-a-store-path\n"},
		{"InvalidCompression", "StorePath: " + testNARStorePath + "\nCompression: lz4\n"},
		{"InvalidFileHash", "StorePath: " + testNARStorePath + "\nFileHash: notahash\n"},
		{"InvalidFileSize", "StorePath: " + testNARStorePath + "\nFileSize: big\n"},
		{"InvalidNarHash", "StorePath: " + testNARStorePath + "\nNarHash: notahash\n"},
		{"InvalidNarSize", "StorePath: " + testNARStorePath + "\nNarSize: big\n"},
		{"InvalidSig", "StorePath: " + testNARStorePath + "\nSig: cache.nixos.org-1:notbase64!!\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ni nix.NARInfo
			require.Error(t, ni.UnmarshalText([]byte(tt.input)))
		})
	}
}

func TestNARInfoAddSignatures(t *testing.T) {
	sig1, err := nix.ParseSignature("cache.nixos.org-1:" + testFakeSigB64)
	require.NoError(t, err)

	sig2, err := nix.ParseSignature("mycache-1:" + testFakeSigB64)
	require.NoError(t, err)

	var ni nix.NARInfo

	// Add two distinct signatures.
	ni.AddSignatures(sig1, sig2)
	require.Len(t, ni.Sig, 2)

	// Adding the same key name again is a no-op.
	ni.AddSignatures(sig1)
	assert.Len(t, ni.Sig, 2)

	// Nil signatures are ignored.
	ni.AddSignatures(nil)
	assert.Len(t, ni.Sig, 2)
}

func BenchmarkNARInfoMarshalText(b *testing.B) {
	var ni nix.NARInfo
	if err := ni.UnmarshalText([]byte(testNARInfoText)); err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		ni.MarshalText()
	}
}

func BenchmarkNARInfoUnmarshalText(b *testing.B) {
	text := []byte(testNARInfoText)

	for b.Loop() {
		var ni nix.NARInfo
		ni.UnmarshalText(text)
	}
}

func FuzzNARInfoUnmarshalText(f *testing.F) {
	// Seed with a structurally valid narinfo so the fuzzer begins from a
	// meaningful input and mutates toward edge cases (truncation, bad field
	// names, malformed hash strings, etc.).
	f.Add(testNARInfoText)
	f.Fuzz(func(_ *testing.T, s string) {
		var ni nix.NARInfo
		ni.UnmarshalText([]byte(s))
	})
}

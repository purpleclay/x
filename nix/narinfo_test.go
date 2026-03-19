package nix_test

import (
	"testing"

	"github.com/purpleclay/x/nix"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Real hello-2.12.1 narinfo values from cache.nixos.org.
const (
	testNARStorePath = "/nix/store/s66mzxpvicwk07gjbjfw9izjfa797vsw-hello-2.12.1"
	testNARURL       = "nar/0vgdasb2zswyxqb2g6q3rabg1rj9ijj21gzgwvkz8bbyb5nji7mk.nar.zst"
	testNARFileHash  = "sha256-xHkGMFLJnBymOQoG6VHjKH/2/wRVNBRg3PUOyBmsMtY="
	testNARNARHash   = "sha256-Dz+UMEzSkFv+Ao5JxWBjyOBnJpokBxCPqJfGjHt+CBg="
	testNARDeriver   = "11cldfnpspd2a1rqihf7xhsjfb0bwj7m-hello-2.12.1.drv"
	testNARSig       = "cache.nixos.org-1:sDcIalq8XS+YQKtGsDi+WlTXk59tMwPkBrgT/LNqJqXEIBMG4AXaYOnaB5dmAerq10TRShEz2C+ynvvL2sINCA=="

	testNARInfoText = "StorePath: /nix/store/s66mzxpvicwk07gjbjfw9izjfa797vsw-hello-2.12.1\n" +
		"URL: nar/0vgdasb2zswyxqb2g6q3rabg1rj9ijj21gzgwvkz8bbyb5nji7mk.nar.zst\n" +
		"Compression: zstd\n" +
		"FileHash: sha256-xHkGMFLJnBymOQoG6VHjKH/2/wRVNBRg3PUOyBmsMtY=\n" +
		"FileSize: 30752\n" +
		"NARHash: sha256-Dz+UMEzSkFv+Ao5JxWBjyOBnJpokBxCPqJfGjHt+CBg=\n" +
		"NARSize: 54416\n" +
		"References: s66mzxpvicwk07gjbjfw9izjfa797vsw-hello-2.12.1 aasvr6g9965bh3df9wgvmjmmapds7b18-glibc-2.40-36\n" +
		"Deriver: 11cldfnpspd2a1rqihf7xhsjfb0bwj7m-hello-2.12.1.drv\n" +
		"Sig: cache.nixos.org-1:sDcIalq8XS+YQKtGsDi+WlTXk59tMwPkBrgT/LNqJqXEIBMG4AXaYOnaB5dmAerq10TRShEz2C+ynvvL2sINCA==\n"
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

	ni := nix.NARInfo{
		StorePath:   sp,
		URL:         testNARURL,
		Compression: nix.Zstd,
		FileHash:    fileHash,
		FileSize:    30752,
		NARHash:     narHash,
		NARSize:     54416,
		References: []string{
			"s66mzxpvicwk07gjbjfw9izjfa797vsw-hello-2.12.1",
			"aasvr6g9965bh3df9wgvmjmmapds7b18-glibc-2.40-36",
		},
		Deriver: testNARDeriver,
		Sig:     []string{testNARSig},
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
	assert.Equal(t, nix.Zstd, ni.Compression)
	assert.Equal(t, int64(30752), ni.FileSize)
	assert.Equal(t, int64(54416), ni.NARSize)
	assert.Equal(t, []string{
		"s66mzxpvicwk07gjbjfw9izjfa797vsw-hello-2.12.1",
		"aasvr6g9965bh3df9wgvmjmmapds7b18-glibc-2.40-36",
	}, ni.References)
	assert.Equal(t, testNARDeriver, ni.Deriver)
	assert.Equal(t, []string{testNARSig}, ni.Sig)
}

func TestNARInfoUnmarshalTextDefaultCompression(t *testing.T) {
	input := "StorePath: " + testNARStorePath + "\n" +
		"URL: " + testNARURL + "\n" +
		"NARSize: 54416\n"

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
		"Sig: cache.nixos.org-1:abc123==\n" +
		"Sig: mycache-1:xyz789==\n"

	var ni nix.NARInfo
	require.NoError(t, ni.UnmarshalText([]byte(input)))
	assert.Equal(t, []string{"cache.nixos.org-1:abc123==", "mycache-1:xyz789=="}, ni.Sig)
}

func TestNARInfoUnknownKeysIgnored(t *testing.T) {
	input := "StorePath: " + testNARStorePath + "\n" +
		"UnknownField: somevalue\n" +
		"NARSize: 54416\n"

	var ni nix.NARInfo
	require.NoError(t, ni.UnmarshalText([]byte(input)))
	assert.Equal(t, int64(54416), ni.NARSize)
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
		{"InvalidNARHash", "StorePath: " + testNARStorePath + "\nNARHash: notahash\n"},
		{"InvalidNARSize", "StorePath: " + testNARStorePath + "\nNARSize: big\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ni nix.NARInfo
			require.Error(t, ni.UnmarshalText([]byte(tt.input)))
		})
	}
}

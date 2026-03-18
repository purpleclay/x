package nix_test

import (
	"testing"

	"github.com/purpleclay/x/nix"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testDigest is the real store path digest for hello-2.12.1 on NixOS.
// The full output path /nix/store/v6szabc66jn5r6vq8p7g7m5b2ks1xwif-hello-2.12.1
// is the canonical Nix hello example referenced throughout the Nix documentation.
const (
	testDigest    = "v6szabc66jn5r6vq8p7g7m5b2ks1xwif"
	testName      = "hello-2.12.1"
	testStorePath = "/nix/store/" + testDigest + "-" + testName
)

// testDrvPath reuses testDigest for simplicity. In practice a .drv file and its
// output have distinct digests, since the derivation hash is computed from inputs
// rather than content.
const testDrvPath = "/nix/store/" + testDigest + "-" + testName + ".drv"

func TestDefaultStoreDirectory(t *testing.T) {
	assert.Equal(t, nix.StoreDirectory("/nix/store"), nix.DefaultStoreDirectory)
}

func TestStoreDirectoryJoin(t *testing.T) {
	tests := []struct {
		name string
		elem []string
		want string
	}{
		{"NoElements", []string{}, "/nix/store"},
		{"Single", []string{"foo"}, "/nix/store/foo"},
		{"Multiple", []string{"foo", "bar", "baz"}, "/nix/store/foo/bar/baz"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, nix.DefaultStoreDirectory.Join(tt.elem...))
		})
	}
}

func TestStoreDirectoryObject(t *testing.T) {
	sp, err := nix.DefaultStoreDirectory.Object(testDigest + "-" + testName)
	require.NoError(t, err)
	assert.Equal(t, nix.StorePath(testStorePath), sp)
}

func TestStoreDirectoryParsePath(t *testing.T) {
	t.Run("ExactStorePath", func(t *testing.T) {
		sp, sub, err := nix.DefaultStoreDirectory.ParsePath(testStorePath)
		require.NoError(t, err)
		assert.Equal(t, nix.StorePath(testStorePath), sp)
		assert.Equal(t, "", sub)
	})

	t.Run("WithSubPath", func(t *testing.T) {
		sp, sub, err := nix.DefaultStoreDirectory.ParsePath(testStorePath + "/bin/hello")
		require.NoError(t, err)
		assert.Equal(t, nix.StorePath(testStorePath), sp)
		assert.Equal(t, "/bin/hello", sub)
	})

	t.Run("OutsideStoreDirectory", func(t *testing.T) {
		_, _, err := nix.DefaultStoreDirectory.ParsePath("/usr/local/bin/hello")
		require.Error(t, err)
		var parseErr *nix.ParseStorePathError
		require.ErrorAs(t, err, &parseErr)
		assert.Equal(t, "/usr/local/bin/hello", parseErr.Input)
		assert.NotNil(t, parseErr.Err)
	})
}

func TestParseStorePath(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		sp, err := nix.ParseStorePath(testStorePath)
		require.NoError(t, err)
		assert.Equal(t, nix.StorePath(testStorePath), sp)
	})

	t.Run("Derivation", func(t *testing.T) {
		sp, err := nix.ParseStorePath(testDrvPath)
		require.NoError(t, err)
		assert.Equal(t, nix.StorePath(testDrvPath), sp)
	})
}

func TestParseStorePathError(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"WrongDirectory", "/usr/local/" + testDigest + "-hello"},
		{"SubPath", testStorePath + "/bin/hello"},
		{"ShortDigest", "/nix/store/abc-hello"},
		{"InvalidDigestChar", "/nix/store/eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee-hello"},
		{"MissingSeparator", "/nix/store/" + testDigest + "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := nix.ParseStorePath(tt.input)
			require.Error(t, err)
			var parseErr *nix.ParseStorePathError
			require.ErrorAs(t, err, &parseErr)
			assert.Equal(t, tt.input, parseErr.Input)
			assert.NotNil(t, parseErr.Err)
		})
	}
}

func TestStorePathAccessors(t *testing.T) {
	sp := nix.StorePath(testStorePath)

	assert.Equal(t, testDigest+"-"+testName, sp.Base())
	assert.Equal(t, nix.DefaultStoreDirectory, sp.Dir())
	assert.Equal(t, testDigest, sp.Digest())
	assert.Equal(t, testName, sp.Name())
	assert.False(t, sp.IsDerivation())
	assert.Equal(t, testStorePath, sp.String())
}

func TestStorePathIsDerivation(t *testing.T) {
	sp := nix.StorePath(testDrvPath)
	assert.True(t, sp.IsDerivation())
}

func TestStorePathTextMarshaler(t *testing.T) {
	original := nix.StorePath(testStorePath)

	text, err := original.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, testStorePath, string(text))

	var restored nix.StorePath
	require.NoError(t, restored.UnmarshalText(text))
	assert.Equal(t, original, restored)
}

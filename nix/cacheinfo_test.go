package nix_test

import (
	"testing"

	"github.com/purpleclay/x/nix"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCacheInfoMarshalText(t *testing.T) {
	c := nix.CacheInfo{
		StoreDir:      "/custom/store",
		WantMassQuery: true,
		Priority:      39,
	}

	text, err := c.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "StoreDir: /custom/store\nWantMassQuery: 1\nPriority: 39\n", string(text))
}

func TestCacheInfoMarshalTextDefaultStoreDir(t *testing.T) {
	c := nix.CacheInfo{WantMassQuery: false, Priority: 39}

	text, err := c.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "StoreDir: /nix/store\nWantMassQuery: 0\nPriority: 39\n", string(text))
}

func TestCacheInfoUnmarshalText(t *testing.T) {
	var c nix.CacheInfo
	err := c.UnmarshalText([]byte("StoreDir: /nix/store\nWantMassQuery: 1\nPriority: 39\n"))
	require.NoError(t, err)
	assert.Equal(t, nix.StoreDirectory("/nix/store"), c.StoreDir)
	assert.True(t, c.WantMassQuery)
	assert.Equal(t, 39, c.Priority)
}

func TestCacheInfoUnknownKeysIgnored(t *testing.T) {
	var c nix.CacheInfo
	err := c.UnmarshalText([]byte("StoreDir: /nix/store\nUnknownKey: somevalue\nWantMassQuery: 0\nPriority: 10\n"))
	require.NoError(t, err)
	assert.Equal(t, nix.StoreDirectory("/nix/store"), c.StoreDir)
	assert.False(t, c.WantMassQuery)
	assert.Equal(t, 10, c.Priority)
}

func TestCacheInfoUnmarshalTextError(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"InvalidWantMassQuery", "StoreDir: /nix/store\nWantMassQuery: yes\nPriority: 40\n"},
		{"InvalidPriority", "StoreDir: /nix/store\nWantMassQuery: 1\nPriority: high\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c nix.CacheInfo
			err := c.UnmarshalText([]byte(tt.input))
			require.Error(t, err)
		})
	}
}

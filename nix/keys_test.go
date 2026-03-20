package nix_test

import (
	"testing"

	"github.com/purpleclay/x/nix"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateKey(t *testing.T) {
	pub, priv, err := nix.GenerateKey("mycache-1", nil)
	require.NoError(t, err)
	assert.Equal(t, "mycache-1", pub.Name())
	assert.Equal(t, "mycache-1", priv.Name())
	assert.Equal(t, pub.String(), priv.PublicKey().String())
}

func TestPublicKeyNilSafety(t *testing.T) {
	var k *nix.PublicKey
	assert.Equal(t, "", k.Name())
	assert.Equal(t, "<nil>", k.String())
}

func TestPrivateKeyNilSafety(t *testing.T) {
	var k *nix.PrivateKey
	assert.Equal(t, "", k.Name())
	assert.Equal(t, "<nil>", k.String())
	assert.Nil(t, k.PublicKey())
}

func TestParsePublicKey(t *testing.T) {
	const s = "cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY="

	pub, err := nix.ParsePublicKey(s)
	require.NoError(t, err)
	assert.Equal(t, "cache.nixos.org-1", pub.Name())
	assert.Equal(t, s, pub.String())
}

func TestParsePublicKeyError(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"MissingSeparator", "cache.nixos.org-16NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY="},
		{"EmptyName", ":6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY="},
		{"InvalidBase64", "mycache-1:not!valid"},
		{"WrongLength", "mycache-1:dG9vc2hvcnQ="},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := nix.ParsePublicKey(tt.input)
			require.Error(t, err)
		})
	}
}

func TestParsePrivateKey(t *testing.T) {
	pub, priv, err := nix.GenerateKey("mycache-1", nil)
	require.NoError(t, err)

	s := priv.String()
	parsed, err := nix.ParsePrivateKey(s)
	require.NoError(t, err)
	assert.Equal(t, "mycache-1", parsed.Name())
	assert.Equal(t, pub.String(), parsed.PublicKey().String())
}

func TestParsePrivateKeyError(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"MissingSeparator", "mycache1AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="},
		{"EmptyName", ":AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="},
		{"InvalidBase64", "mycache-1:not!valid"},
		{"WrongLength", "mycache-1:dG9vc2hvcnQ="},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := nix.ParsePrivateKey(tt.input)
			require.Error(t, err)
		})
	}
}

func TestKeyRoundTrip(t *testing.T) {
	pub, priv, err := nix.GenerateKey("mycache-1", nil)
	require.NoError(t, err)

	pubText, err := pub.MarshalText()
	require.NoError(t, err)

	var restoredPub nix.PublicKey
	require.NoError(t, restoredPub.UnmarshalText(pubText))
	assert.Equal(t, pub.String(), restoredPub.String())

	privText, err := priv.MarshalText()
	require.NoError(t, err)

	var restoredPriv nix.PrivateKey
	require.NoError(t, restoredPriv.UnmarshalText(privText))
	assert.Equal(t, priv.String(), restoredPriv.String())
}

func FuzzParsePublicKey(f *testing.F) {
	// Real cache.nixos.org-1 public key.
	f.Add("cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY=")
	f.Fuzz(func(_ *testing.T, s string) {
		nix.ParsePublicKey(s) // must not panic; errors are acceptable
	})
}

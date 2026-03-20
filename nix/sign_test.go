package nix_test

import (
	"testing"

	"github.com/purpleclay/x/nix"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Real cache.nixos.org-1 public key, well-known and published at
// https://cache.nixos.org/nix-cache-info.
const cacheNixosOrgPubKey = "cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY="

func TestSignatureNilSafety(t *testing.T) {
	var s *nix.Signature
	assert.Equal(t, "", s.Name())
	assert.Equal(t, "<nil>", s.String())
}

func TestParseSignature(t *testing.T) {
	sig, err := nix.ParseSignature(testNARSig)
	require.NoError(t, err)
	assert.Equal(t, "cache.nixos.org-1", sig.Name())
	assert.Equal(t, testNARSig, sig.String())
}

func TestParseSignatureError(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"MissingSeparator", "cache.nixos.org-1sDcIalq8XS+YQKtGsDi+WlTXk59tMwPkBrgT/LNqJqXEIBMG4AXaYOnaB5dmAerq10TRShEz2C+ynvvL2sINCA=="},
		{"EmptyName", ":" + testNARSig},
		{"InvalidBase64", "mycache-1:not!valid=="},
		{"WrongLength", "mycache-1:dG9vc2hvcnQ="},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := nix.ParseSignature(tt.input)
			require.Error(t, err)
		})
	}
}

func TestSignatureRoundTrip(t *testing.T) {
	pub, priv, err := nix.GenerateKey("mycache-1", nil)
	require.NoError(t, err)

	var ni nix.NARInfo
	require.NoError(t, ni.UnmarshalText([]byte(testNARInfoText)))

	sig, err := nix.SignNARInfo(priv, &ni)
	require.NoError(t, err)

	// Marshal → parse → verify
	text, err := sig.MarshalText()
	require.NoError(t, err)

	var restored nix.Signature
	require.NoError(t, restored.UnmarshalText(text))

	require.NoError(t, nix.VerifyNARInfo([]*nix.PublicKey{pub}, &ni, &restored))
}

func TestVerifyNARInfoCacheNixosOrg(t *testing.T) {
	pub, err := nix.ParsePublicKey(cacheNixosOrgPubKey)
	require.NoError(t, err)

	var ni nix.NARInfo
	require.NoError(t, ni.UnmarshalText([]byte(testNARInfoText)))
	require.Len(t, ni.Sig, 1)

	require.NoError(t, nix.VerifyNARInfo([]*nix.PublicKey{pub}, &ni, ni.Sig[0]))
}

func TestVerifyNARInfoNoTrustedKey(t *testing.T) {
	var ni nix.NARInfo
	require.NoError(t, ni.UnmarshalText([]byte(testNARInfoText)))
	require.Len(t, ni.Sig, 1)

	err := nix.VerifyNARInfo([]*nix.PublicKey{}, &ni, ni.Sig[0])
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no trusted key matches")
}

func TestVerifyNARInfoInvalidSignature(t *testing.T) {
	pub, priv, err := nix.GenerateKey("mycache-1", nil)
	require.NoError(t, err)

	var ni nix.NARInfo
	require.NoError(t, ni.UnmarshalText([]byte(testNARInfoText)))

	sig, err := nix.SignNARInfo(priv, &ni)
	require.NoError(t, err)

	// Tamper: change NARSize so the fingerprint no longer matches the signature.
	ni.NARSize++

	err = nix.VerifyNARInfo([]*nix.PublicKey{pub}, &ni, sig)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "verification failed")
}

func BenchmarkSignNARInfo(b *testing.B) {
	_, priv, err := nix.GenerateKey("mycache-1", nil)
	if err != nil {
		b.Fatal(err)
	}

	var ni nix.NARInfo
	if err := ni.UnmarshalText([]byte(testNARInfoText)); err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		nix.SignNARInfo(priv, &ni)
	}
}

func BenchmarkVerifyNARInfo(b *testing.B) {
	pub, priv, err := nix.GenerateKey("mycache-1", nil)
	if err != nil {
		b.Fatal(err)
	}

	var ni nix.NARInfo
	if err := ni.UnmarshalText([]byte(testNARInfoText)); err != nil {
		b.Fatal(err)
	}

	sig, err := nix.SignNARInfo(priv, &ni)
	if err != nil {
		b.Fatal(err)
	}

	trusted := []*nix.PublicKey{pub}

	for b.Loop() {
		nix.VerifyNARInfo(trusted, &ni, sig)
	}
}

func FuzzParseSignature(f *testing.F) {
	f.Add(testNARSig)
	f.Fuzz(func(_ *testing.T, s string) {
		nix.ParseSignature(s)
	})
}

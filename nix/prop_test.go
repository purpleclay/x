package nix_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/purpleclay/x/nix"
	"github.com/purpleclay/x/nix/base32"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

// drawHash generates a random valid Hash of any of the four hash types.
// Both the type and the raw bytes are drawn uniformly at random, so the
// generator covers MD5/SHA1/SHA256/SHA512 with equal probability.
func drawHash(t *rapid.T, label string) nix.Hash {
	typ := rapid.SampledFrom([]nix.HashType{nix.MD5, nix.SHA1, nix.SHA256, nix.SHA512}).Draw(t, label+".type")
	raw := rapid.SliceOfN(rapid.Byte(), typ.Size(), typ.Size()).Draw(t, label+".bytes")
	return nix.NewHash(typ, raw)
}

// drawStorePath generates a random valid Nix store path rooted at the default
// store directory. The 32-character nix-base32 digest is derived from 20 random
// bytes — encoding them guarantees a structurally valid digest without needing
// to filter-reject invalid strings.
func drawStorePath(t *rapid.T, label string) nix.StorePath {
	raw := rapid.SliceOfN(rapid.Byte(), 20, 20).Draw(t, label+".digest.raw")
	digest := base32.StdEncoding.EncodeToString(raw)
	name := rapid.StringMatching(`[a-z][a-z0-9+_.-]{0,19}`).Draw(t, label+".name")
	sp, err := nix.ParseStorePath("/nix/store/" + digest + "-" + name)
	require.NoError(t, err)
	return sp
}

// drawNARInfo generates a random but structurally valid NARInfo. The Sig field
// is always empty — signing properties are covered separately in TestPropSignVerify.
func drawNARInfo(t *rapid.T) nix.NARInfo {
	sp := drawStorePath(t, "narinfo.storepath")

	// Generate between 0 and 4 references, each a valid store-path base name.
	refCount := rapid.IntRange(0, 4).Draw(t, "narinfo.refcount")
	refs := make([]string, refCount)
	for i := range refCount {
		raw := rapid.SliceOfN(rapid.Byte(), 20, 20).Draw(t, fmt.Sprintf("ref[%d].digest.raw", i))
		digest := base32.StdEncoding.EncodeToString(raw)
		name := rapid.StringMatching(`[a-z][a-z0-9+_.-]{0,9}`).Draw(t, fmt.Sprintf("ref[%d].name", i))
		refs[i] = digest + "-" + name
	}

	// Deriver is optional: 50% chance of being set.
	deriver := ""
	if rapid.Bool().Draw(t, "narinfo.hasderiver") {
		raw := rapid.SliceOfN(rapid.Byte(), 20, 20).Draw(t, "deriver.digest.raw")
		digest := base32.StdEncoding.EncodeToString(raw)
		name := rapid.StringMatching(`[a-z][a-z0-9+_.-]{0,9}`).Draw(t, "deriver.name")
		deriver = digest + "-" + name + ".drv"
	}

	return nix.NARInfo{
		StorePath:   sp,
		URL:         "nar/" + sp.Digest() + ".nar",
		Compression: rapid.SampledFrom([]nix.Compression{nix.Bzip2, nix.Gzip, nix.XZ, nix.Zstd, nix.Brotli, nix.NoCompression}).Draw(t, "narinfo.compression"),
		FileHash:    drawHash(t, "narinfo.filehash"),
		FileSize:    rapid.Int64Range(0, 1<<40).Draw(t, "narinfo.filesize"),
		NARHash:     drawHash(t, "narinfo.narhash"),
		NARSize:     rapid.Int64Range(0, 1<<40).Draw(t, "narinfo.narsize"),
		References:  refs,
		Deriver:     deriver,
	}
}

func TestPropHashRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		h := drawHash(t, "hash")

		for _, encoded := range []string{h.Base32(), h.SRI(), h.Hex()} {
			parsed, err := nix.ParseHash(encoded)
			require.NoErrorf(t, err, "ParseHash(%q)", encoded)
			require.Truef(t, parsed.Equal(h), "round-trip mismatch for %q: got %v, want %v", encoded, parsed, h)
		}
	})
}

func TestPropNARInfoRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		original := drawNARInfo(t)

		text, err := original.MarshalText()
		require.NoError(t, err)

		var got nix.NARInfo
		require.NoErrorf(t, got.UnmarshalText(text), "UnmarshalText\nwire text:\n%s", text)

		require.Equal(t, original.StorePath, got.StorePath)
		require.Equal(t, original.URL, got.URL)
		require.Equal(t, original.Compression, got.Compression)
		require.Truef(t, original.FileHash.Equal(got.FileHash), "FileHash: got %v, want %v", got.FileHash, original.FileHash)
		require.Equal(t, original.FileSize, got.FileSize)
		require.Truef(t, original.NARHash.Equal(got.NARHash), "NARHash: got %v, want %v", got.NARHash, original.NARHash)
		require.Equal(t, original.NARSize, got.NARSize)
		// slices.Equal treats nil and []string{} as equal; require.Equal does not.
		require.Truef(t, slices.Equal(original.References, got.References), "References: got %v, want %v", got.References, original.References)
		require.Equal(t, original.Deriver, got.Deriver)
	})
}

func TestPropKeyRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Key names follow the Nix convention "name-N" (e.g. "mycache-1").
		keyName := rapid.StringMatching(`[a-z][a-z0-9-]{1,10}-[0-9]`).Draw(t, "keyname")
		pub, priv, err := nix.GenerateKey(keyName, nil)
		require.NoError(t, err)

		// Public key: string → parse → string must be identity.
		parsedPub, err := nix.ParsePublicKey(pub.String())
		require.NoErrorf(t, err, "ParsePublicKey(%q)", pub.String())
		require.Equal(t, pub.String(), parsedPub.String())

		// Private key: string → parse → string must be identity.
		parsedPriv, err := nix.ParsePrivateKey(priv.String())
		require.NoError(t, err)
		require.Equal(t, priv.String(), parsedPriv.String())

		// The public key derived from the private key must match the generated public key.
		require.Equal(t, pub.String(), priv.PublicKey().String())
	})
}

func TestPropSignVerify(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		keyName := rapid.StringMatching(`[a-z][a-z0-9-]{1,10}-[0-9]`).Draw(t, "keyname")
		pub, priv, err := nix.GenerateKey(keyName, nil)
		require.NoError(t, err)

		ni := drawNARInfo(t)

		sig, err := nix.SignNARInfo(priv, &ni)
		require.NoError(t, err)

		// Verification must succeed with the matching public key.
		require.NoError(t, nix.VerifyNARInfo([]*nix.PublicKey{pub}, &ni, sig))

		// Signature round-trip: the wire representation must parse back to an
		// equivalent signature that still verifies.
		parsed, err := nix.ParseSignature(sig.String())
		require.NoErrorf(t, err, "ParseSignature(%q)", sig.String())
		require.Equal(t, sig.String(), parsed.String())
		require.NoError(t, nix.VerifyNARInfo([]*nix.PublicKey{pub}, &ni, parsed))

		// Mutating any field covered by the fingerprint must break verification.
		// The fingerprint covers StorePath, NARHash, NARSize, and References.
		ni.NARSize++
		require.Error(t, nix.VerifyNARInfo([]*nix.PublicKey{pub}, &ni, sig), "VerifyNARInfo should fail after mutating NARSize")
	})
}

// Package nix provides types and functions for interoperating with the
// [Nix] package manager, specifically its binary cache protocol.
//
// This package implements the core primitives needed to serve and push
// to a Nix binary cache: store paths, NAR info, cache info, content
// hashing (including the non-standard nix base32 encoding), and
// Ed25519-based NAR signing and verification.
//
// # Nix Base32 Encoding
//
// Nix uses a non-standard base32 encoding with the alphabet
// "0123456789abcdfghijklmnpqrsvwxyz" (note: e, o, t, u are missing).
// Additionally, bytes are read in reverse order compared to RFC 4648.
// The [base32] sub-package exposes [base32.StdEncoding] for direct use.
//
// [Nix]: https://nixos.org/
package nix

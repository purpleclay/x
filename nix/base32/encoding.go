package base32

import "strconv"

const (
	alphabet     = "0123456789abcdfghijklmnpqrsvwxyz"
	invalidIndex = '\xff'
)

// InvalidByteError is returned by [Encoding.Decode] and [Encoding.DecodeString]
// when an invalid nix base32 byte is encountered.
type InvalidByteError struct {
	// Byte is the offending character.
	Byte byte
	// Offset is the byte position of the offending character within the input.
	Offset int64
}

func (e InvalidByteError) Error() string {
	return "invalid nix base32 character " + strconv.QuoteRune(rune(e.Byte)) + " at input byte " + strconv.FormatInt(e.Offset, 10)
}

// nixDecodeMap is a pre-computed lookup table mapping ASCII bytes to their
// base32 digit value, or 0xff if the byte is not a valid alphabet character.
var nixDecodeMap = func() [256]byte {
	var m [256]byte
	for i := range m {
		m[i] = 0xff
	}
	for i, c := range []byte(alphabet) {
		m[c] = byte(i)
	}
	return m
}()

// StdEncoding is the standard Nix base32 encoding, as defined in the Nix C++ library.
var StdEncoding = newEncoding()

// An Encoding is a Nix base32 encoding/decoding scheme. It uses a custom
// 32-character alphabet that omits e, o, t, and u to reduce visual ambiguity.
// Unlike RFC 4648 base32, Nix base32 does not pad output and writes encoded
// characters in reverse bit order.
//
// See the Nix C++ reference implementation:
// https://github.com/NixOS/nix/blob/master/src/libutil/base-nix-32.cc
type Encoding struct {
	encode    [32]byte
	decodeMap [256]byte
}

func newEncoding() *Encoding {
	e := &Encoding{}
	copy(e.encode[:], alphabet)
	e.decodeMap = nixDecodeMap
	return e
}

// EncodedLen returns the length of the base32 encoding of n bytes.
//
//	// base32.StdEncoding.EncodedLen(20) == 32
func (enc *Encoding) EncodedLen(n int) int {
	if n == 0 {
		return 0
	}
	return (n*8 + 4) / 5
}

// DecodedLen returns the maximum length in bytes of the decoded
// data corresponding to n bytes of nix base32 encoded data.
//
//	// base32.StdEncoding.DecodedLen(32) == 20
func (enc *Encoding) DecodedLen(n int) int {
	return n * 5 / 8
}

// EncodeToString returns the nix base32 encoding of src.
//
//	src, _ := hex.DecodeString("1f74d74729abdc08f4f84e8f7f8c808c8ed92ee5")
//	// base32.StdEncoding.EncodeToString(src) == "wlpdk3lch267z3sfz3s0ip5b553xfx0z"
func (enc *Encoding) EncodeToString(src []byte) string {
	dst := make([]byte, enc.EncodedLen(len(src)))
	enc.Encode(dst, src)
	return string(dst)
}

// Encode encodes src into [Encoding.EncodedLen](len(src)) bytes of dst.
//
// Nix base32 uses a custom alphabet and reverses the output character
// order compared to a naive base32. This matches the Nix C++ reference
// implementation which iterates from the most significant output
// position down to 0.
//
//	src, _ := hex.DecodeString("1f74d74729abdc08f4f84e8f7f8c808c8ed92ee5")
//	dst := make([]byte, base32.StdEncoding.EncodedLen(len(src)))
//	base32.StdEncoding.Encode(dst, src)
//	// string(dst) == "wlpdk3lch267z3sfz3s0ip5b553xfx0z"
func (enc *Encoding) Encode(dst, src []byte) {
	n := len(src)
	if n == 0 {
		return
	}
	outLen := enc.EncodedLen(n)

	for i := range outLen {
		b := uint(i) * 5
		j := b / 8
		k := b % 8

		// Widen to uint before shifting to avoid byte truncation.
		var c uint
		if int(j) < n {
			c = uint(src[j]) >> k
		}
		if k > 3 && int(j+1) < n {
			c |= uint(src[j+1]) << (8 - k)
		}
		dst[outLen-1-i] = enc.encode[c&0x1f]
	}
}

// DecodeString returns the bytes represented by the nix base32 string s.
//
//	b, _ := base32.StdEncoding.DecodeString("wlpdk3lch267z3sfz3s0ip5b553xfx0z")
//	// hex.EncodeToString(b) == "1f74d74729abdc08f4f84e8f7f8c808c8ed92ee5"
func (enc *Encoding) DecodeString(s string) ([]byte, error) {
	dst := make([]byte, enc.DecodedLen(len(s)))
	n, err := enc.Decode(dst, []byte(s))
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}

// Decode decodes src using nix base32. It writes at most
// [Encoding.DecodedLen](len(src)) bytes to dst and returns the number of
// bytes written.
func (enc *Encoding) Decode(dst, src []byte) (int, error) {
	if len(src) == 0 {
		return 0, nil
	}

	outLen := enc.DecodedLen(len(src))
	for i := range outLen {
		dst[i] = 0
	}

	for i := range len(src) {
		// Read src in reverse order to undo the encoding reversal.
		c := src[len(src)-1-i]
		digit := enc.decodeMap[c]
		if digit == invalidIndex {
			return 0, InvalidByteError{Byte: c, Offset: int64(len(src) - 1 - i)}
		}

		b := uint(i) * 5
		j := b / 8
		k := b % 8

		// Widen to uint before shifting to avoid byte truncation.
		if int(j) < outLen {
			dst[j] |= byte(uint(digit) << k)
		}
		if k > 3 && int(j+1) < outLen {
			dst[j+1] |= byte(uint(digit) >> (8 - k))
		}
	}

	// Verify that any padding bits are zero.
	totalBits := uint(len(src)) * 5
	dataBits := uint(outLen) * 8
	if totalBits > dataBits {
		extraBits := totalBits - dataBits
		firstDigit := enc.decodeMap[src[0]]
		if firstDigit>>(5-extraBits) != 0 {
			return 0, InvalidByteError{Byte: src[0], Offset: 0}
		}
	}

	return outLen, nil
}

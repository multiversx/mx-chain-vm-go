package math

import "math/big"

// UnsignedByteLength returns the number of bytes that the absolute value of the
// given big int occupies, i.e. len(bi.Bytes()), without allocating those bytes.
//
// Note that, just like big.Int.Bytes, this disregards the sign.
func UnsignedByteLength(bi *big.Int) int {
	if bi == nil {
		return 0
	}
	return (bi.BitLen() + 7) / 8
}

// TwosComplementByteLength returns the number of bytes that the two's complement
// representation of the given value occupies, i.e. len(twos.ToBytes(bi)),
// without allocating or building that representation.
//
// The length is derived from the bit length alone:
//   - zero is encoded on 0 bytes;
//   - a positive value needs BitLen() bits plus one sign bit that must be 0,
//     hence BitLen()/8 + 1 bytes;
//   - a negative value -a is encoded as the bit-negation of a-1 plus a sign bit
//     that must be 1, hence BitLen(a-1)/8 + 1 bytes. BitLen(a-1) equals
//     BitLen(a), except when a is an exact power of two, when it is one less.
//
// TODO: move to mx-components-big-int/twos-complement after the next VM release.
func TwosComplementByteLength(bi *big.Int) int {
	if bi == nil {
		return 0
	}
	sign := bi.Sign()
	if sign == 0 {
		return 0
	}

	// big.Int.BitLen and big.Int.TrailingZeroBits both refer to the absolute value
	bitLen := bi.BitLen()
	if sign < 0 && bi.TrailingZeroBits() == uint(bitLen-1) {
		// the absolute value is an exact power of two, so subtracting 1 from it
		// (as the two's complement encoding does) drops the leading bit
		bitLen--
	}

	return bitLen/8 + 1
}

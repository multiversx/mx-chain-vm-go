package math

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	twos "github.com/multiversx/mx-components-big-int/twos-complement"
	"github.com/stretchr/testify/require"
)

func TestUnsignedByteLength_KnownValues(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		value    int64
		expected int
	}{
		{0, 0},
		{1, 1},
		{255, 1},
		{256, 2},
		{65535, 2},
		{65536, 3},
		{-1, 1},
		{-255, 1},
		{-256, 2},
		{-65536, 3},
	} {
		require.Equal(t, tc.expected, UnsignedByteLength(big.NewInt(tc.value)),
			"unexpected length for %d", tc.value)
	}
}

// the length must always agree with big.Int.Bytes, which is what
// BigIntGetUnsignedBytes writes out; the sign is disregarded by both
func TestUnsignedByteLength_MatchesBytes(t *testing.T) {
	t.Parallel()

	forEachTestValue(t, func(value *big.Int) {
		require.Equal(t, len(value.Bytes()), UnsignedByteLength(value),
			"unexpected length for %s", value.String())
	})
}

// testing.AllocsPerRun cannot be used in a parallel test, hence no t.Parallel() here
func TestUnsignedByteLength_DoesNotAllocate(t *testing.T) {
	positive := new(big.Int).Lsh(big.NewInt(1), 1<<20)
	negative := new(big.Int).Neg(positive)

	allocations := testing.AllocsPerRun(100, func() {
		_ = UnsignedByteLength(positive)
		_ = UnsignedByteLength(negative)
	})
	require.Zero(t, allocations)
}

func TestTwosComplementByteLength_KnownValues(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		value    int64
		expected int
	}{
		{0, 0},
		{1, 1},
		{127, 1},
		{128, 2},
		{255, 2},
		{256, 2},
		{32767, 2},
		{32768, 3},
		{-1, 1},
		{-127, 1},
		{-128, 1},
		{-129, 2},
		{-255, 2},
		{-256, 2},
		{-257, 2},
		{-32768, 2},
		{-32769, 3},
	} {
		require.Equal(t, tc.expected, TwosComplementByteLength(big.NewInt(tc.value)),
			"unexpected length for %d", tc.value)
	}
}

// the length must always agree with the actual two's complement encoding
func TestTwosComplementByteLength_MatchesEncoding(t *testing.T) {
	t.Parallel()

	forEachTestValue(t, func(value *big.Int) {
		require.Equal(t, len(twos.ToBytes(value)), TwosComplementByteLength(value),
			"unexpected length for %s", value.String())
	})
}

// forEachTestValue covers small values exhaustively, then the boundaries where the
// byte length changes (powers of two and their neighbours), then random values of
// every bit width; each value is checked with both signs.
func forEachTestValue(t *testing.T, assertMatches func(value *big.Int)) {
	t.Run("small values, exhaustive", func(t *testing.T) {
		for i := int64(-100000); i <= 100000; i++ {
			assertMatches(big.NewInt(i))
		}
	})

	t.Run("powers of two and their neighbours", func(t *testing.T) {
		for exponent := 0; exponent < 1024; exponent++ {
			powerOfTwo := new(big.Int).Lsh(big.NewInt(1), uint(exponent))
			for _, delta := range []int64{-2, -1, 0, 1, 2} {
				value := new(big.Int).Add(powerOfTwo, big.NewInt(delta))
				assertMatches(value)
				assertMatches(new(big.Int).Neg(value))
			}
		}
	})

	t.Run("random values of all sizes", func(t *testing.T) {
		rng := rand.New(rand.NewSource(42))
		for bits := 1; bits <= 1024; bits++ {
			limit := new(big.Int).Lsh(big.NewInt(1), uint(bits))
			for i := 0; i < 10; i++ {
				value := new(big.Int).Rand(rng, limit)
				assertMatches(value)
				assertMatches(new(big.Int).Neg(value))
			}
		}
	})
}

// the whole point of the helper: no allocation, whatever the size of the value
// testing.AllocsPerRun cannot be used in a parallel test, hence no t.Parallel() here
func TestTwosComplementByteLength_DoesNotAllocate(t *testing.T) {
	for _, bits := range []uint{64, 1 << 10, 1 << 16, 1 << 20} {
		t.Run(fmt.Sprintf("%d bits", bits), func(t *testing.T) {
			powerOfTwo := new(big.Int).Lsh(big.NewInt(1), bits)
			// -2^bits is the worst case: all low words are zero, so the
			// power-of-two check has to scan the entire value
			values := []*big.Int{
				powerOfTwo,
				new(big.Int).Neg(powerOfTwo),
				new(big.Int).Sub(powerOfTwo, big.NewInt(12345)),
				new(big.Int).Neg(new(big.Int).Sub(powerOfTwo, big.NewInt(12345))),
			}

			allocations := testing.AllocsPerRun(100, func() {
				for _, value := range values {
					_ = TwosComplementByteLength(value)
				}
			})
			require.Zero(t, allocations)
		})
	}
}

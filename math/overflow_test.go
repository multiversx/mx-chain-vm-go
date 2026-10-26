package math

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

var maxEvenInt64 = int64(math.MaxInt64 - 1)
var maxEvenUint64 = uint64(math.MaxUint64 - 1)

func TestAddUint64WithErr(t *testing.T) {
	a := uint64(1)
	b := uint64(math.MaxUint64 - 2)
	sum, err := AddUint64WithErr(a, b)
	require.Nil(t, err)
	require.Equal(t, uint64(math.MaxUint64-1), sum)

	c := uint64(3)
	sum, err = AddUint64WithErr(b, c)
	require.Equal(t, ErrAdditionOverflow, err)
	require.Equal(t, uint64(math.MaxUint64), sum)

	d := uint64(0)
	sum, err = AddUint64WithErr(b, d)
	require.Nil(t, err)
	require.Equal(t, uint64(math.MaxUint64-2), sum)
}

func TestAddUint64(t *testing.T) {
	a := uint64(1)
	b := uint64(math.MaxUint64 - 2)
	sum := AddUint64(a, b)
	require.Equal(t, uint64(math.MaxUint64-1), sum)

	c := uint64(3)
	sum = AddUint64(b, c)
	require.Equal(t, uint64(math.MaxUint64), sum)

	d := uint64(0)
	sum = AddUint64(b, d)
	require.Equal(t, uint64(math.MaxUint64-2), sum)
}

func TestAddInt32WithErr(t *testing.T) {
	a := int32(1)
	b := int32(math.MaxInt32 - 2)
	sum, err := AddInt32WithErr(a, b)
	require.Nil(t, err)
	require.Equal(t, int32(math.MaxInt32-1), sum)

	c := int32(3)
	sum, err = AddInt32WithErr(b, c)
	require.Equal(t, ErrAdditionOverflow, err)
	require.Equal(t, int32(math.MaxInt32), sum)

	d := int32(0)
	sum, err = AddInt32WithErr(b, d)
	require.Nil(t, err)
	require.Equal(t, int32(math.MaxInt32-2), sum)
}

func TestAddInt32(t *testing.T) {
	a := int32(1)
	b := int32(math.MaxInt32 - 2)
	sum := AddInt32(a, b)
	require.Equal(t, int32(math.MaxInt32-1), sum)

	c := int32(3)
	sum = AddInt32(b, c)
	require.Equal(t, int32(math.MaxInt32), sum)

	d := int32(0)
	sum = AddInt32(b, d)
	require.Equal(t, int32(math.MaxInt32-2), sum)
}

func TestSubUint64(t *testing.T) {
	a := uint64(2)
	b := uint64(1)
	diff := SubUint64(a, b)
	require.Equal(t, uint64(1), diff)

	a = uint64(2)
	b = uint64(2)
	diff = SubUint64(a, b)
	require.Equal(t, uint64(0), diff)

	a = uint64(2)
	b = uint64(3)
	diff = SubUint64(a, b)
	require.Equal(t, uint64(0), diff)

	a = uint64(2)
	b = uint64(math.MaxUint64)
	diff = SubUint64(a, b)
	require.Equal(t, uint64(0), diff)
}

func TestSubInt(t *testing.T) {
	require.Equal(t, math.MinInt, math.MinInt64)

	a := math.MinInt
	b := 1
	diff := SubInt(a, b)
	require.Equal(t, math.MinInt, diff)

	a = math.MinInt + 10
	b = 1
	diff = SubInt(a, b)
	require.Equal(t, math.MinInt+9, diff)

	a = math.MinInt + 10
	b = 10
	diff = SubInt(a, b)
	require.Equal(t, math.MinInt, diff)
}

func TestMulUint64WithErr(t *testing.T) {
	a := maxEvenUint64 / 2
	b := uint64(2)
	product, err := MulUint64WithErr(a, b)
	require.Nil(t, err)
	require.Equal(t, maxEvenUint64, product)

	a = maxEvenUint64 / 2
	b = uint64(1)
	product, err = MulUint64WithErr(a, b)
	require.Nil(t, err)
	require.Equal(t, a, product)

	a = maxEvenUint64 / 2
	b = uint64(3)
	product, err = MulUint64WithErr(a, b)
	require.NotNil(t, err)
	require.Equal(t, uint64(math.MaxUint64), product)
}

func TestMulUint64(t *testing.T) {
	a := maxEvenUint64 / 2
	b := uint64(2)
	product := MulUint64(a, b)
	require.Equal(t, maxEvenUint64, product)

	a = maxEvenUint64 / 2
	b = uint64(1)
	product = MulUint64(a, b)
	require.Equal(t, a, product)

	a = maxEvenUint64 / 2
	b = uint64(3)
	product = MulUint64(a, b)
	require.Equal(t, uint64(math.MaxUint64), product)
}

func TestSubUint64WithErr(t *testing.T) {
	diff, err := SubUint64WithErr(0, 0)
	require.Nil(t, err)
	require.Equal(t, uint64(0), diff)

	diff, err = SubUint64WithErr(10, 4)
	require.Nil(t, err)
	require.Equal(t, uint64(6), diff)

	diff, err = SubUint64WithErr(10, 10)
	require.Nil(t, err)
	require.Equal(t, uint64(0), diff)

	diff, err = SubUint64WithErr(math.MaxUint64, 1)
	require.Nil(t, err)
	require.Equal(t, uint64(math.MaxUint64-1), diff)

	diff, err = SubUint64WithErr(0, 1)
	require.Equal(t, ErrSubtractionUnderflow, err)
	require.Equal(t, uint64(0), diff)

	diff, err = SubUint64WithErr(5, 10)
	require.Equal(t, ErrSubtractionUnderflow, err)
	require.Equal(t, uint64(0), diff)

	diff, err = SubUint64WithErr(0, math.MaxUint64)
	require.Equal(t, ErrSubtractionUnderflow, err)
	require.Equal(t, uint64(0), diff)
}

func TestSubInt32(t *testing.T) {
	require.Equal(t, int32(0), SubInt32(0, 0))
	require.Equal(t, int32(6), SubInt32(10, 4))
	require.Equal(t, int32(-14), SubInt32(-10, 4))
	require.Equal(t, int32(math.MinInt32), SubInt32(math.MinInt32, 1))
	require.Equal(t, int32(math.MinInt32), SubInt32(math.MinInt32, 100))
	require.Equal(t, int32(math.MaxInt32), SubInt32(math.MaxInt32, -1))
}

func TestSubInt32WithErr(t *testing.T) {
	diff, err := SubInt32WithErr(0, 0)
	require.Nil(t, err)
	require.Equal(t, int32(0), diff)

	diff, err = SubInt32WithErr(10, 4)
	require.Nil(t, err)
	require.Equal(t, int32(6), diff)

	diff, err = SubInt32WithErr(-10, -4)
	require.Nil(t, err)
	require.Equal(t, int32(-6), diff)

	diff, err = SubInt32WithErr(math.MinInt32, 0)
	require.Nil(t, err)
	require.Equal(t, int32(math.MinInt32), diff)

	diff, err = SubInt32WithErr(math.MaxInt32, 0)
	require.Nil(t, err)
	require.Equal(t, int32(math.MaxInt32), diff)

	diff, err = SubInt32WithErr(math.MinInt32, 1)
	require.Equal(t, ErrSubtractionUnderflow, err)
	require.Equal(t, int32(math.MinInt32), diff)

	diff, err = SubInt32WithErr(math.MaxInt32, -1)
	require.Equal(t, ErrSubtractionUnderflow, err)
	require.Equal(t, int32(math.MaxInt32), diff)

	diff, err = SubInt32WithErr(math.MinInt32, math.MaxInt32)
	require.Equal(t, ErrSubtractionUnderflow, err)
	require.Equal(t, int32(math.MinInt32), diff)

	diff, err = SubInt32WithErr(math.MaxInt32, math.MinInt32)
	require.Equal(t, ErrSubtractionUnderflow, err)
	require.Equal(t, int32(math.MaxInt32), diff)
}

func TestCheckBoundsInt32(t *testing.T) {
	// valid boundary cases
	require.Nil(t, CheckBoundsInt32(0, 0, 0))
	require.Nil(t, CheckBoundsInt32(0, 10, 10))
	require.Nil(t, CheckBoundsInt32(5, 5, 10))
	require.Nil(t, CheckBoundsInt32(10, 0, 10))
	require.Nil(t, CheckBoundsInt32(0, math.MaxInt32, math.MaxInt32))

	// negative start or length
	require.Equal(t, ErrOutOfBounds, CheckBoundsInt32(-1, 5, 10))
	require.Equal(t, ErrOutOfBounds, CheckBoundsInt32(0, -1, 10))
	require.Equal(t, ErrOutOfBounds, CheckBoundsInt32(-5, -5, 10))
	require.Equal(t, ErrOutOfBounds, CheckBoundsInt32(math.MinInt32, 1, 10))
	require.Equal(t, ErrOutOfBounds, CheckBoundsInt32(0, math.MinInt32, 10))

	// exceeding bounds
	require.Equal(t, ErrOutOfBounds, CheckBoundsInt32(5, 6, 10))
	require.Equal(t, ErrOutOfBounds, CheckBoundsInt32(11, 0, 10))
	require.Equal(t, ErrOutOfBounds, CheckBoundsInt32(0, 1, 0))
	require.Equal(t, ErrOutOfBounds, CheckBoundsInt32(0, 0, -1))
	require.Equal(t, ErrOutOfBounds, CheckBoundsInt32(1, math.MaxInt32, math.MaxInt32))

	// int32 overflow cases (wrapping prevention)
	require.Equal(t, ErrOutOfBounds, CheckBoundsInt32(math.MaxInt32, 1, 100))
	require.Equal(t, ErrOutOfBounds, CheckBoundsInt32(1, math.MaxInt32, 100))
	require.Equal(t, ErrOutOfBounds, CheckBoundsInt32(math.MaxInt32, math.MaxInt32, 100))
}

package math

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

var MaxEvenInt64 = int64(math.MaxInt64 - 1)
var MaxEvenUint64 = uint64(math.MaxUint64 - 1)

func TestAddUint64WithErr(t *testing.T) {
	a := uint64(1)
	b := uint64(math.MaxUint64 - 2)
	sum, err := AddUint64WithErr(a, b)
	require.Nil(t, err)
	require.Equal(t, uint64(math.MaxUint64-1), sum)

	c := uint64(3)
	sum, err = AddUint64WithErr(b, c)
	require.Equal(t, ErrAdditionOverflow, err)
	require.Equal(t, b+c, sum)

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

func TestAddInt64WithErr(t *testing.T) {
	a := int64(1)
	b := int64(math.MaxInt64 - 2)
	sum, err := AddInt64WithErr(a, b)
	require.Nil(t, err)
	require.Equal(t, int64(math.MaxInt64-1), sum)

	c := int64(3)
	sum, err = AddInt64WithErr(b, c)
	require.Equal(t, ErrAdditionOverflow, err)
	require.Equal(t, b+c, sum)

	d := int64(0)
	sum, err = AddInt64WithErr(b, d)
	require.Nil(t, err)
	require.Equal(t, int64(math.MaxInt64-2), sum)
}

func TestAddInt64(t *testing.T) {
	a := int64(1)
	b := int64(math.MaxInt64 - 2)
	sum := AddInt64(a, b)
	require.Equal(t, int64(math.MaxInt64-1), sum)

	c := int64(3)
	sum = AddInt64(b, c)
	require.Equal(t, int64(math.MaxInt64), sum)

	d := int64(0)
	sum = AddInt64(b, d)
	require.Equal(t, int64(math.MaxInt64-2), sum)
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
	require.Equal(t, b+c, sum)

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

func TestMulInt64WithErr(t *testing.T) {
	a := int64(MaxEvenInt64 / 2)
	b := int64(2)
	product, err := MulInt64WithErr(a, b)
	require.Nil(t, err)
	require.Equal(t, MaxEvenInt64, product)

	a = int64(MaxEvenInt64 / 2)
	b = int64(1)
	product, err = MulInt64WithErr(a, b)
	require.Nil(t, err)
	require.Equal(t, a, product)

	a = int64(MaxEvenInt64 / 2)
	b = int64(3)
	product, err = MulInt64WithErr(a, b)
	require.NotNil(t, err)
	require.Equal(t, int64(0), product)
}

func TestMulInt64(t *testing.T) {
	a := int64(MaxEvenInt64 / 2)
	b := int64(2)
	product := MulInt64(a, b)
	require.Equal(t, MaxEvenInt64, product)

	a = int64(MaxEvenInt64 / 2)
	b = int64(1)
	product = MulInt64(a, b)
	require.Equal(t, a, product)

	a = int64(MaxEvenInt64 / 2)
	b = int64(3)
	product = MulInt64(a, b)
	require.Equal(t, int64(math.MaxInt64), product)
}

func TestMulUint64WithErr(t *testing.T) {
	a := uint64(MaxEvenUint64 / 2)
	b := uint64(2)
	product, err := MulUint64WithErr(a, b)
	require.Nil(t, err)
	require.Equal(t, MaxEvenUint64, product)

	a = uint64(MaxEvenUint64 / 2)
	b = uint64(1)
	product, err = MulUint64WithErr(a, b)
	require.Nil(t, err)
	require.Equal(t, a, product)

	a = uint64(MaxEvenUint64 / 2)
	b = uint64(3)
	product, err = MulUint64WithErr(a, b)
	require.NotNil(t, err)
	require.Equal(t, uint64(0), product)
}

func TestMulUint64(t *testing.T) {
	a := uint64(MaxEvenUint64 / 2)
	b := uint64(2)
	product := MulUint64(a, b)
	require.Equal(t, MaxEvenUint64, product)

	a = uint64(MaxEvenUint64 / 2)
	b = uint64(1)
	product = MulUint64(a, b)
	require.Equal(t, a, product)

	a = uint64(MaxEvenUint64 / 2)
	b = uint64(3)
	product = MulUint64(a, b)
	require.Equal(t, uint64(math.MaxUint64), product)
}

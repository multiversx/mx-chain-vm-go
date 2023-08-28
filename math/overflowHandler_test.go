package math

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

var expectedErr = errors.New("expected error")

func TestOverflowHandler_AddInt64(t *testing.T) {
	t.Parallel()

	t.Run("at the limit, does not cause an overflow", func(t *testing.T) {
		handler := NewOverflowHandler()

		sum := handler.AddInt64(1, math.MaxInt64-2)
		require.Nil(t, handler.Error())
		require.Equal(t, int64(math.MaxInt64-1), sum)
	})
	t.Run("over the limit, does cause an overflow", func(t *testing.T) {
		handler := NewOverflowHandler()

		sum := handler.AddInt64(math.MaxInt64-2, 3)
		require.ErrorIs(t, handler.Error(), ErrAdditionOverflow)
		require.Equal(t, int64(math.MaxInt64), sum)
	})
	t.Run("addition with 0 at the limit, does not cause an overflow", func(t *testing.T) {
		handler := NewOverflowHandler()

		sum := handler.AddInt64(math.MaxInt64-2, 0)
		require.Nil(t, handler.Error())
		require.Equal(t, int64(math.MaxInt64-2), sum)
	})
	t.Run("addition with negative numbers, does not cause an overflow", func(t *testing.T) {
		handler := NewOverflowHandler()

		sum := handler.AddInt64(-5, 4)
		require.Nil(t, handler.Error())
		require.Equal(t, int64(-1), sum)
	})
	t.Run("over the negative limit, does cause an overflow", func(t *testing.T) {
		handler := NewOverflowHandler()

		sum := handler.AddInt64(-math.MaxInt64+2, -4)
		require.ErrorIs(t, handler.Error(), ErrAdditionOverflow)
		require.Equal(t, int64(math.MaxInt64), sum)
	})
	t.Run("error already set, should not perform add", func(t *testing.T) {
		handler := NewOverflowHandler()
		handler.err = expectedErr

		t.Run("overflow occurs", func(t *testing.T) {
			sum := handler.AddInt64(-math.MaxInt64+2, -4) // overflow occurs
			require.Equal(t, int64(math.MaxInt64), sum)
			require.ErrorIs(t, handler.Error(), expectedErr) // not the overflow error as the sum was not performed
		})
		t.Run("no overflow occurs", func(t *testing.T) {
			sum := handler.AddInt64(0, 0)
			require.Equal(t, int64(math.MaxInt64), sum)
			require.ErrorIs(t, handler.Error(), expectedErr)
		})
	})
}

func TestOverflowHandler_MulInt64(t *testing.T) {
	t.Parallel()

	t.Run("multiply with no overflow should work", func(t *testing.T) {
		handler := NewOverflowHandler()

		product := handler.MulInt64(maxEvenInt64/2, 2)
		require.Nil(t, handler.Error())
		require.Equal(t, maxEvenInt64, product)

		product = handler.MulInt64(maxEvenInt64/2, 1)
		require.Nil(t, handler.Error())
		require.Equal(t, maxEvenInt64/2, product)
	})
	t.Run("multiply with overflow should error", func(t *testing.T) {
		handler := NewOverflowHandler()

		product := handler.MulInt64(maxEvenInt64/2, 3)
		require.ErrorIs(t, handler.Error(), ErrMultiplicationOverflow)
		require.Equal(t, int64(math.MaxInt64), product)
	})
	t.Run("error already set, should not perform multiply", func(t *testing.T) {
		handler := NewOverflowHandler()
		handler.err = expectedErr

		t.Run("overflow occurs", func(t *testing.T) {
			product := handler.MulInt64(maxEvenInt64/2, 3) // overflow occurs
			require.Equal(t, int64(math.MaxInt64), product)
			require.ErrorIs(t, handler.Error(), expectedErr) // not the overflow error as the sum was not performed
		})
		t.Run("no overflow occurs", func(t *testing.T) {
			product := handler.MulInt64(1, 1)
			require.Equal(t, int64(math.MaxInt64), product)
			require.ErrorIs(t, handler.Error(), expectedErr)
		})
	})
}

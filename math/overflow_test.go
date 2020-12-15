package math

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddUint64(t *testing.T) {
	a := uint64(1)
	b := uint64(math.MaxUint64 - 2)
	sum, err := AddUint64WithErr(a, b)
	require.Nil(t, err)
	require.Equal(t, uint64(math.MaxUint64-1), sum)

	c := uint64(3)
	sum, err = AddUint64WithErr(b, c)
	require.Equal(t, ErrAdditionOverflow, err)

	d := uint64(0)
	sum, err = AddUint64WithErr(b, d)
	require.Nil(t, err)
}

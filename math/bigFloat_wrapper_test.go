package math

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBigFloatSub(t *testing.T) {
	value1 := new(big.Float)
	_ = value1.GobDecode([]byte{1, 10, 0, 0, 0, 53, 0, 0, 0, 58, 0, 31, 28, 26, 150, 254, 14, 45})
	value2 := new(big.Float)
	_ = value2.GobDecode([]byte{1, 10, 0, 0, 0, 53, 0, 0, 0, 52, 222, 212, 49, 108, 64, 122, 107, 100})
	var err error
	result := new(big.Float)
	err = Sub(result, value1, value2)
	require.Equal(t, big.NewFloat(0), result)
	require.Equal(t, ErrOperationCausingPanic, err)
}

package math

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBigFloatSub_Panic(t *testing.T) {
	value1 := new(big.Float)
	_ = value1.GobDecode([]byte{1, 10, 0, 0, 0, 53, 0, 0, 0, 58, 0, 31, 28, 26, 150, 254, 14, 45})
	value2 := new(big.Float)
	_ = value2.GobDecode([]byte{1, 10, 0, 0, 0, 53, 0, 0, 0, 52, 222, 212, 49, 108, 64, 122, 107, 100})
	result, err := SubBigFloat(value1, value2)

	require.Equal(t, big.NewFloat(0), result)
	require.Equal(t, ErrBigFloatSub, err)
}

func TestBigFloatSub_Success(t *testing.T) {
	value1 := new(big.Float)
	_ = value1.GobDecode([]byte{1, 10, 0, 0, 0, 53, 0, 0, 0, 58, 0, 31, 28, 26, 150, 254, 14, 45})
	value2 := new(big.Float)
	_ = value2.GobDecode([]byte{1, 10, 0, 0, 0, 53, 0, 0, 0, 68, 222, 212, 49, 108, 64, 122, 107, 100})
	result, err := SubBigFloat(value1, value2)

	require.Nil(t, err)
	require.Equal(t, new(big.Float).Sub(value1, value2), result)

	encodedResult, _ := result.GobEncode()
	require.Equal(t, []byte{1, 19, 0, 0, 0, 53, 0, 0, 0, 68, 222, 212, 41, 165, 57, 212, 168, 0}, encodedResult)
}

func TestBigFloatAdd_Panic(t *testing.T) {
	value1 := new(big.Float)
	_ = value1.GobDecode([]byte{1, 10, 0, 0, 0, 53, 0, 0, 0, 58, 0, 31, 28, 26, 150, 254, 14, 45})
	value2 := new(big.Float)
	_ = value2.GobDecode([]byte{1, 11, 0, 0, 0, 53, 0, 0, 0, 52, 222, 212, 49, 108, 64, 122, 107, 100})
	result, err := AddBigFloat(value1, value2)

	require.Equal(t, big.NewFloat(0), result)
	require.Equal(t, ErrBigFloatAdd, err)
}

func TestBigFloatAdd_Success(t *testing.T) {
	value1 := new(big.Float)
	_ = value1.GobDecode([]byte{1, 10, 0, 0, 0, 53, 0, 0, 0, 58, 0, 31, 28, 26, 150, 254, 14, 45})
	value2 := new(big.Float)
	_ = value2.GobDecode([]byte{1, 10, 0, 0, 0, 53, 0, 0, 0, 52, 222, 212, 49, 108, 64, 122, 107, 100})
	result, err := AddBigFloat(value1, value2)

	require.Nil(t, err)
	require.Equal(t, new(big.Float).Add(value1, value2), result)

	encodedResult, _ := result.GobEncode()
	require.Equal(t, []byte{1, 18, 0, 0, 0, 53, 0, 0, 0, 52, 230, 155, 56, 17, 255, 253, 248, 0}, encodedResult)
}

func TestBigFloatQuo_Panic(t *testing.T) {
	value1, value2 := big.NewFloat(0), big.NewFloat(0)
	result, err := QuoBigFloat(value1, value2)

	require.Equal(t, big.NewFloat(0), result)
	require.Equal(t, ErrBigFloatQuo, err)
}

func TestBigFloatQuo_Success(t *testing.T) {
	value1, value2 := big.NewFloat(4), big.NewFloat(2)
	result, err := QuoBigFloat(value1, value2)

	require.Nil(t, err)
	require.Equal(t, big.NewFloat(2), result)
}

func TestBigFloatMul_Panic(t *testing.T) {
	value1, value2 := big.NewFloat(0), new(big.Float).SetInf(false)
	result, err := MulBigFloat(value1, value2)

	require.Equal(t, ErrBigFloatMul, err)
	require.Equal(t, big.NewFloat(0), result)
}

func TestBigFloatMul_Success(t *testing.T) {
	value1 := new(big.Float)
	_ = value1.GobDecode([]byte{1, 10, 0, 0, 0, 53, 0, 0, 0, 58, 0, 31, 28, 26, 150, 254, 14, 45})
	value2 := new(big.Float)
	_ = value2.GobDecode([]byte{1, 10, 0, 0, 0, 53, 0, 0, 0, 68, 222, 212, 49, 108, 64, 122, 107, 100})
	result, err := MulBigFloat(value1, value2)
	encodedResult, _ := result.GobEncode()

	require.Nil(t, err)
	require.Equal(t, []byte{1, 2, 0, 0, 0, 53, 0, 0, 0, 115, 216, 161, 66, 179, 241, 21, 160, 0}, encodedResult)
}

func TestBigFloatSqrt_Panic(t *testing.T) {
	value := big.NewFloat(-1)
	result, err := SqrtBigFloat(value)

	require.Equal(t, big.NewFloat(0), result)
	require.Equal(t, ErrBigFloatSqrt, err)
}

func TestBigFloatSqrt_Success(t *testing.T) {
	value1 := new(big.Float)
	_ = value1.GobDecode([]byte{1, 10, 0, 0, 0, 53, 0, 0, 0, 58, 0, 31, 28, 26, 150, 254, 14, 45})
	result, err := SqrtBigFloat(value1)
	encodedResult, _ := result.GobEncode()

	require.Nil(t, err)
	require.Equal(t, []byte{1, 10, 0, 0, 0, 53, 0, 0, 0, 20, 131, 226, 32, 17, 29, 166, 88, 0}, encodedResult)
}

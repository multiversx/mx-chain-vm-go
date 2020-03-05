package arwen

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/stretchr/testify/require"
)

func TestConvertReturnValue(t *testing.T) {
	t.Parallel()

	value := int32(100)

	wasmI32 := wasmer.I32(value)
	bigValue := ConvertReturnValue(wasmI32)
	require.Equal(t, big.NewInt(int64(value)).Bytes(), bigValue)

	wasmI64 := wasmer.I64(int64(value))
	bigValue = ConvertReturnValue(wasmI64)
	require.Equal(t, big.NewInt(int64(value)).Bytes(), bigValue)

	wasmVoid := wasmer.Void()
	bigValue = ConvertReturnValue(wasmVoid)
	require.Equal(t, big.NewInt(0).Bytes(), bigValue)

	defer func() {
		r := recover()
		if r == nil {
			require.Fail(t, "test should panic")
		}
	}()
	wasm := wasmer.F32(0)
	_ = ConvertReturnValue(wasm)
}

func TestGuardedMakeByteSlice2D(t *testing.T) {
	t.Parallel()

	byteSlice, err := GuardedMakeByteSlice2D(-1)
	require.Error(t, err)
	require.Nil(t, byteSlice)

	byteSlice, err = GuardedMakeByteSlice2D(0)
	require.Nil(t, err)
	require.NotNil(t, byteSlice)
}

func TestGuardedGetBytesSlice(t *testing.T) {
	t.Parallel()

	dataSlice := []byte("data1_data2_data3_data4")

	slice, err := GuardedGetBytesSlice(dataSlice, 100, 100)
	require.Nil(t, slice)
	require.NotNil(t, err)

	slice, err = GuardedGetBytesSlice(dataSlice, 5, -1)
	require.Nil(t, slice)
	require.NotNil(t, err)

	expectedResult := []byte("data1")
	slice, err = GuardedGetBytesSlice(dataSlice, 0, 5)
	require.Nil(t, err)
	require.True(t, bytes.Equal(expectedResult, slice))
}

func TestInverseBytes(t *testing.T) {
	t.Parallel()

	data := []byte("qwerty")
	expectedData := []byte("ytrewq")

	result := InverseBytes(data)
	require.Equal(t, expectedData, result)

	result = InverseBytes(nil)
	require.Equal(t, []byte{}, result)

	result = InverseBytes([]byte("a"))
	require.Equal(t, []byte("a"), result)
}

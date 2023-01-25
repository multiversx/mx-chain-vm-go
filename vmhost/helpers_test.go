package vmhost

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

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

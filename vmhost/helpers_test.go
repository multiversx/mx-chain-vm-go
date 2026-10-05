package vmhost

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

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

func TestWrappableError_UnwrapAndErrorsIs(t *testing.T) {
	t.Parallel()

	baseErr := errors.New("base error")
	targetErr := errors.New("target error")

	// 1. Single error wrapped
	wErr1 := WrapError(baseErr)
	require.True(t, errors.Is(wErr1, baseErr))
	require.False(t, errors.Is(wErr1, targetErr))
	require.Equal(t, baseErr, errors.Unwrap(wErr1))

	// 2. Multiple wrapped errors
	wErr2 := wErr1.WrapWithError(targetErr)
	require.True(t, errors.Is(wErr2, baseErr))
	require.True(t, errors.Is(wErr2, targetErr))

	// 3. Unwrap chain traversal should not panic
	current := error(wErr2)
	for current != nil {
		current = errors.Unwrap(current)
	}

	// 4. Unwrap on empty wrappableError should return nil
	emptyErr := &wrappableError{}
	require.Nil(t, emptyErr.Unwrap())
}

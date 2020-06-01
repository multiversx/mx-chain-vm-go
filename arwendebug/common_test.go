package arwendebug

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DecodeArguments(t *testing.T) {
	decoded, err := decodeArguments([]string{"64", "74657374"})
	require.Nil(t, err)
	require.Equal(t, []byte{100}, decoded[0])
	require.Equal(t, []byte("test"), decoded[1])

	decoded, err = decodeArguments([]string{"0"})
	require.Equal(t, ErrInvalidArgumentEncoding, err)

	decoded, err = decodeArguments([]string{"foo"})
	require.Equal(t, ErrInvalidArgumentEncoding, err)
}

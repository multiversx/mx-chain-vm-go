package arwendebug

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FixTestAddress(t *testing.T) {
	require.Len(t, fixTestAddress("alice"), 32)
	require.Len(t, fixTestAddress("bob"), 32)
	require.Equal(t, "000000000000000000000000000alice", string(fixTestAddress("alice")))
	require.Equal(t, "00000000000000000000000000000bob", string(fixTestAddress("bob")))
}

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

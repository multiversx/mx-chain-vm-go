package mandosjsonparse

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBigInt(t *testing.T) {
	p := Parser{}
	result, err := p.parseBigInt("5", bigIntSignedBytes)
	require.Nil(t, err)
	require.Equal(t, big.NewInt(5), result)

	result, err = p.parseBigInt("0x05", bigIntSignedBytes)
	require.Nil(t, err)
	require.Equal(t, big.NewInt(5), result)

	result, err = p.parseBigInt("-1", bigIntSignedBytes)
	require.Nil(t, err)
	require.Equal(t, big.NewInt(-1), result)

	result, err = p.parseBigInt("-1", bigIntUnsignedBytes)
	require.Nil(t, err)
	require.Equal(t, big.NewInt(255), result)

	result, err = p.parseBigInt("0x", bigIntSignedBytes)
	require.Nil(t, err)
	require.Equal(t, big.NewInt(0), result)

	result, err = p.parseBigInt("", bigIntSignedBytes)
	require.Nil(t, err)
	require.Equal(t, big.NewInt(0), result)

	result, err = p.parseBigInt("0x00", bigIntSignedBytes)
	require.Nil(t, err)
	require.True(t, big.NewInt(0).Cmp(result) == 0)
}

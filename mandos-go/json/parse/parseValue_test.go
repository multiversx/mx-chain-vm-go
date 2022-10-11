package mandosjsonparse

import (
	"math/big"
	"testing"

	oj "github.com/ElrondNetwork/wasm-vm-v1_4/mandos-go/orderedjson"
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

func TestParseBool(t *testing.T) {
	p := Parser{}

	objBool := oj.OJsonBool(false)
	valueBool, err := p.parseBool(&objBool)
	require.Nil(t, err)
	require.Equal(t, false, valueBool)

	objBool = true
	valueBool, err = p.parseBool(&objBool)
	require.Nil(t, err)
	require.Equal(t, true, valueBool)

	objStr := oj.OJsonString{Value: "my_str"}
	valueBool, err = p.parseBool(&objStr)
	require.NotNil(t, err)

	valueBool, err = p.parseBool(nil)
	require.NotNil(t, err)
}

package ethapi

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_convertToEthAddress(t *testing.T) {
	elrondAddress, _ := hex.DecodeString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	expectedResult, _ := hex.DecodeString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	result := convertToEthAddress(elrondAddress)

	assert.Equal(t, expectedResult, result)
}

func Test_converToEthU128(t *testing.T) {
	data, _ := hex.DecodeString("aa")
	expectedResult, _ := hex.DecodeString("000000000000000000000000000000aa")

	result := converToEthU128(data)
	assert.Equal(t, expectedResult, result)
}

func Test_converToEthU128_whenLargeData(t *testing.T) {
	data, _ := hex.DecodeString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	expectedResult, _ := hex.DecodeString("00000000000000000000000000000000")

	result := converToEthU128(data)
	assert.Equal(t, expectedResult, result)
}

package mandosjsonmodel

import (
	"bytes"
	"encoding/hex"
	"math/big"
)

// ResultEqual returns true if result bytes encode the same number.
func ResultEqual(expected JSONBytesFromString, actual []byte) bool {
	if bytes.Equal(expected.Value, actual) {
		return true
	}

	return big.NewInt(0).SetBytes(expected.Value).Cmp(big.NewInt(0).SetBytes(actual)) == 0
}

// ResultAsString helps create nicer error messages.
func ResultAsString(result [][]byte) string {
	str := "["
	for i, res := range result {
		str += "0x" + hex.EncodeToString(res)
		if i < len(result)-1 {
			str += ", "
		}
	}
	return str + "]"
}

// JSONBytesFromStringValues extracts values from a slice of JSONBytesFromString into a list
func JSONBytesFromStringValues(jbs []JSONBytesFromString) [][]byte {
	result := make([][]byte, len(jbs))
	for i, jb := range jbs {
		result[i] = jb.Value
	}
	return result
}

// JSONBytesFromTreeValues extracts values from a slice of JSONBytesFromTree into a list
func JSONBytesFromTreeValues(jbs []JSONBytesFromTree) [][]byte {
	result := make([][]byte, len(jbs))
	for i, jb := range jbs {
		result[i] = jb.Value
	}
	return result
}

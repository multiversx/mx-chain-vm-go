package scenjsonmodel

import (
	"bytes"
	"math/big"
)

// ResultEqual returns true if result bytes encode the same number.
func ResultEqual(expected JSONBytesFromString, actual []byte) bool {
	if bytes.Equal(expected.Value, actual) {
		return true
	}

	return big.NewInt(0).SetBytes(expected.Value).Cmp(big.NewInt(0).SetBytes(actual)) == 0
}

// JSONBytesFromTreeValues extracts values from a slice of JSONBytesFromTree into a list
func JSONBytesFromTreeValues(jbs []JSONBytesFromTree) [][]byte {
	result := make([][]byte, len(jbs))
	for i, jb := range jbs {
		result[i] = jb.Value
	}
	return result
}

func (tgs TraceGasStatus) ToInt() int {
	return int(tgs)
}

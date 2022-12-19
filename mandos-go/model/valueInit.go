package mandosjsonmodel

import (
	"math/big"

	oj "github.com/ElrondNetwork/wasm-vm-v1_4/mandos-go/orderedjson"
)

// JSONBytesFromString stores a byte slice
type JSONBytesFromString struct {
	Value       []byte
	Original    string
	Unspecified bool
}

// NewJSONBytesFromString creates a new JSONBytesFromString instance.
func NewJSONBytesFromString(value []byte, originalStr string) JSONBytesFromString {
	return JSONBytesFromString{
		Value:       value,
		Original:    originalStr,
		Unspecified: false,
	}
}

// JSONBytesEmpty creates a new JSONBytesFromString instance with default values.
func JSONBytesEmpty() JSONBytesFromString {
	return JSONBytesFromString{
		Value:       nil,
		Original:    "",
		Unspecified: true,
	}
}

// JSONBytesFromTree stores a parsed byte slice, either from a string, or from a list of strings.
// The list of strings representation can be used in storage, arguments or results,
// and it is designed to make it easier to express serialized objects.
// The strings in the list get simply concatenated to produce a value.
type JSONBytesFromTree struct {
	Value       []byte
	Original    oj.OJsonObject
	Unspecified bool
}

// OriginalEmpty returns true if the object originates from "".
func (jb JSONBytesFromTree) OriginalEmpty() bool {
	if str, isStr := jb.Original.(*oj.OJsonString); isStr {
		return len(str.Value) == 0
	}
	return false
}

// JSONBigInt stores the parsed big int value but also the original parsed string
type JSONBigInt struct {
	Value       *big.Int
	Original    string
	Unspecified bool
}

// JSONBigIntZero provides an unitialized zero value.
func JSONBigIntZero() JSONBigInt {
	return JSONBigInt{
		Value:       big.NewInt(0),
		Original:    "",
		Unspecified: true,
	}
}

// JSONUint64 stores the parsed uint64 value but also the original parsed string
type JSONUint64 struct {
	Value       uint64
	Original    string
	Unspecified bool
}

// OriginalEmpty returns true if the object originates from "".
func (ju *JSONUint64) OriginalEmpty() bool {
	return len(ju.Original) == 0
}

// JSONUint64Zero provides an unitialized zero value.
func JSONUint64Zero() JSONUint64 {
	return JSONUint64{
		Value:       0,
		Original:    "",
		Unspecified: true,
	}
}

// JSONValueList represents a list of values, as expressed in JSON
type JSONValueList struct {
	Values []JSONBytesFromString
}

// IsUnspecified yields true if the field was originally unspecified.
func (jvl JSONValueList) IsUnspecified() bool {
	return len(jvl.Values) == 0
}

// ToValues extracts values from a JSONValueList
func (jvl JSONValueList) ToValues() [][]byte {
	result := make([][]byte, len(jvl.Values))
	for i, jb := range jvl.Values {
		result[i] = jb.Value
	}
	return result
}

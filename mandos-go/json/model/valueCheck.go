package mandosjsonmodel

import (
	"bytes"
	"math/big"

	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
)

// JSONCheckBytes holds a byte slice condition.
// Values are checked for equality.
// "*" allows all values.
type JSONCheckBytes struct {
	Value    []byte
	IsStar   bool
	Original oj.OJsonObject
}

// JSONCheckBytesDefault yields JSONCheckBytes default "*" value.
func JSONCheckBytesDefault() JSONCheckBytes {
	return JSONCheckBytes{
		Value:    []byte{},
		IsStar:   true,
		Original: &oj.OJsonString{Value: ""},
	}
}

// JSONCheckBytesExplicitStar yields JSONCheckBytes explicit "*" value.
func JSONCheckBytesExplicitStar() JSONCheckBytes {
	return JSONCheckBytes{
		Value:    []byte{},
		IsStar:   true,
		Original: &oj.OJsonString{Value: "*"},
	}
}

// JSONCheckBytesReconstructed creates a JSONCheckBytes without an original JSON source.
func JSONCheckBytesReconstructed(value []byte) JSONCheckBytes {
	return JSONCheckBytes{
		Value:    value,
		IsStar:   false,
		Original: &oj.OJsonString{Value: ""},
	}
}

// OriginalEmpty returns true if original = "".
func (jcbytes JSONCheckBytes) OriginalEmpty() bool {
	if str, isStr := jcbytes.Original.(*oj.OJsonString); isStr {
		return len(str.Value) == 0
	}
	return false
}

// IsDefault yields true if the field was originally unspecified.
func (jcbytes JSONCheckBytes) IsDefault() bool {
	return jcbytes.IsStar && jcbytes.OriginalEmpty()
}

// Check returns true if condition expressed in object holds for another value.
// Explicit values are interpreted as equals assertion.
func (jcbytes JSONCheckBytes) Check(other []byte) bool {
	if jcbytes.IsStar {
		return true
	}
	return bytes.Equal(jcbytes.Value, other)
}

// JSONCheckBigInt holds a big int condition.
// Values are checked for equality.
// "*" allows all values.
type JSONCheckBigInt struct {
	Value    *big.Int
	IsStar   bool
	Original string
}

// JSONCheckBigIntDefault yields JSONCheckBigInt default "*" value.
func JSONCheckBigIntDefault() JSONCheckBigInt {
	return JSONCheckBigInt{
		Value:    nil,
		IsStar:   true,
		Original: "",
	}
}

// IsDefault yields true if the field was originally unspecified.
func (jcbi JSONCheckBigInt) IsDefault() bool {
	return jcbi.IsStar && len(jcbi.Original) == 0
}

// Check returns true if condition expressed in object holds for another value.
// Explicit values are interpreted as equals assertion.
func (jcbi JSONCheckBigInt) Check(other *big.Int) bool {
	if jcbi.IsStar {
		return true
	}
	return jcbi.Value.Cmp(other) == 0
}

// JSONCheckUint64 holds a uint64 condition.
// Values are checked for equality.
// "*" allows all values.
type JSONCheckUint64 struct {
	Value    uint64
	IsStar   bool
	Original string
}

// JSONCheckUint64Default yields JSONCheckBigInt default "*" value.
func JSONCheckUint64Default() JSONCheckUint64 {
	return JSONCheckUint64{
		Value:    0,
		IsStar:   true,
		Original: "",
	}
}

// IsDefault yields true if the field was originally unspecified.
func (jcu JSONCheckUint64) IsDefault() bool {
	return jcu.IsStar && len(jcu.Original) == 0
}

// Check returns true if condition expressed in object holds for another value.
// Explicit values are interpreted as equals assertion.
func (jcu JSONCheckUint64) Check(other uint64) bool {
	if jcu.IsStar {
		return true
	}
	return jcu.Value == other
}

// CheckBool interprets own value as bool (true = anything > 0, false = 0),
// We are using JSONCheckUint64 for bool too so we don't create another type.
func (jcu JSONCheckUint64) CheckBool(other bool) bool {
	if jcu.IsStar {
		return true
	}
	return jcu.Value > 0 == other
}

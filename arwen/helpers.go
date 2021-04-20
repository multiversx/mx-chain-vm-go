package arwen

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/math"
)

// Zero is the big integer 0
var Zero = big.NewInt(0)

// One is the big integer 1
var One = big.NewInt(1)

// CustomStorageKey appends the given key type to the given associated key
func CustomStorageKey(keyType string, associatedKey []byte) []byte {
	return append(associatedKey, []byte(keyType)...)
}

// BooleanToInt returns 1 if the given bool is true, 0 otherwise
func BooleanToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// GuardedMakeByteSlice2D creates a new two-dimensional byte slice of the given dimension.
func GuardedMakeByteSlice2D(length int32) ([][]byte, error) {
	if length < 0 {
		return nil, fmt.Errorf("GuardedMakeByteSlice2D: negative length (%d)", length)
	}

	result := make([][]byte, length)
	return result, nil
}

// GuardedGetBytesSlice returns a chunk from the given data
func GuardedGetBytesSlice(data []byte, offset int32, length int32) ([]byte, error) {
	dataLength := uint32(len(data))
	isOffsetTooSmall := offset < 0
	isOffsetTooLarge := uint32(offset) > dataLength
	requestedEnd := math.AddInt32(offset, length)
	isRequestedEndTooLarge := uint32(requestedEnd) > dataLength
	isLengthNegative := length < 0

	if isOffsetTooSmall || isOffsetTooLarge || isRequestedEndTooLarge {
		return nil, fmt.Errorf("GuardedGetBytesSlice: bad bounds")
	}

	if isLengthNegative {
		return nil, fmt.Errorf("GuardedGetBytesSlice: negative length")
	}

	result := data[offset:requestedEnd]
	return result, nil
}

// PadBytesLeft adds a padding of the given size to the left the byte slice
func PadBytesLeft(data []byte, size int) []byte {
	if data == nil {
		return nil
	}
	if len(data) == 0 {
		return []byte{}
	}
	padSize := math.SubInt(size, len(data))
	if padSize <= 0 {
		return data
	}

	paddedBytes := make([]byte, padSize)
	paddedBytes = append(paddedBytes, data...)
	return paddedBytes
}

// InverseBytes reverses the bytes of the given byte slice
func InverseBytes(data []byte) []byte {
	length := len(data)
	invBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		invBytes[length-i-1] = data[i]
	}
	return invBytes
}

// WithFault returns true if the error is not nil, and uses the remaining gas if the execution has failed
func WithFault(err error, context unsafe.Pointer, failExecution bool) bool {
	if err == nil {
		return false
	}

	if failExecution {
		runtime := GetRuntimeContext(context)
		metering := GetMeteringContext(context)

		metering.UseGas(metering.GasLeft())
		runtime.FailExecution(err)
	}

	return true
}

// GetSCCode returns the SC code from a given file
func GetSCCode(fileName string) []byte {
	code, _ := ioutil.ReadFile(filepath.Clean(fileName))
	return code
}

// U64ToLEB128 encodes an uint64 using LEB128 (Little Endian Base 128), used in WASM bytecode
// See https://en.wikipedia.org/wiki/LEB128
// Copied from https://github.com/filecoin-project/go-leb128/blob/master/leb128.go
func U64ToLEB128(n uint64) (out []byte) {
	more := true
	for more {
		b := byte(n & 0x7F)
		n >>= 7
		if n == 0 {
			more = false
		} else {
			b |= 0x80
		}
		out = append(out, b)
	}
	return
}

// IfNil tests if the provided interface pointer or underlying object is nil
func IfNil(checker nilInterfaceChecker) bool {
	if checker == nil {
		return true
	}
	return checker.IsInterfaceNil()
}

type nilInterfaceChecker interface {
	IsInterfaceNil() bool
}

package arwen

import (
	"fmt"
	"unsafe"
)

func GuardedMakeByteSlice2D(length int32) ([][]byte, error) {
	if length < 0 {
		return nil, fmt.Errorf("GuardedMakeByteSlice2D: negative length (%d)", length)
	}

	result := make([][]byte, length)
	return result, nil
}

func GuardedGetBytesSlice(data []byte, offset int32, length int32) ([]byte, error) {
	dataLength := uint32(len(data))
	isOffsetTooSmall := offset < 0
	isOffsetTooLarge := uint32(offset) > dataLength
	requestedEnd := uint32(offset + length)
	isRequestedEndTooLarge := requestedEnd > dataLength
	isLengthNegative := length < 0

	if isOffsetTooSmall || isOffsetTooLarge {
		return nil, fmt.Errorf("GuardedGetBytesSlice: bad bounds")
	}

	if isRequestedEndTooLarge {
		return nil, fmt.Errorf("GuardedGetBytesSlice: bad bounds")
	}

	if isLengthNegative {
		return nil, fmt.Errorf("GuardedGetBytesSlice: negative length")
	}

	result := data[offset : offset+length]
	return result, nil
}

func InverseBytes(data []byte) []byte {
	length := len(data)
	invBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		invBytes[length-i-1] = data[i]
	}
	return invBytes
}

func WithFault(err error, context unsafe.Pointer, failExecution bool) bool {
	if err != nil {
		runtime := GetRuntimeContext(context)
		metering := GetMeteringContext(context)

		metering.UseGas(metering.GasLeft())
		if failExecution {
			runtime.FailExecution(err)
		}

		return true
	}

	return false
}

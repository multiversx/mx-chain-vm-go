package arwen

import "C"
import (
	"github.com/wasmerio/go-ext-wasm/wasmer"
	"unsafe"
)

const addressLen = 32
const hashLen = 32
const balanceLen = 32

func goByteSlice(data *C.uint8_t, size C.size_t) []byte {
	if size == 0 {
		return []byte{}
	}
	return (*[1 << 30]byte)(unsafe.Pointer(data))[:size:size]
}

func loadBytes(from *wasmer.Memory, offset int32, length int32) []byte {
	if from.Length() < uint32(offset+length) {
		return from.Data()[offset:]
	}

	return from.Data()[offset : offset+length]
}

func storeBytes(to *wasmer.Memory, offset int32, data []byte) error {
	var memoryData = to.Data()
	length := int32(len(data))

	if to.Length() < uint32(offset+length) {
		err := to.Grow(1)
		if err != nil {
			return err
		}
	}

	copy(memoryData[offset:offset+length], data)

	return nil
}

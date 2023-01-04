package executor

import (
	"fmt"

	"github.com/ElrondNetwork/wasm-vm/math"
)

// MemPtr indicates that an argument refers to a location in WASM memory.
type MemPtr int32

// MemLength indicates that an argument refers to the length of a data section in WASM memory.
type MemLength = int32

// Offset adds to a pointer to WASM memory.
func (memPtr MemPtr) Offset(offset int32) MemPtr {
	return MemPtr(math.AddInt32(int32(memPtr), offset))
}

// Memory defines the functionality of the memory of a Wasmer instance.
// Now considered an implementation detail and will likely stop being a public interface.
type Memory interface {
	Length() uint32
	Data() []byte
	Grow(pages uint32) error
	Destroy()
	IsInterfaceNil() bool
}

// MemLoad returns the contents from the given offset of the WASM memory.
func MemLoad(memory Memory, memPtr MemPtr, length MemLength) ([]byte, error) {
	if length == 0 {
		return []byte{}, nil
	}

	memoryView := memory.Data()
	memoryLength := memory.Length()
	requestedEnd := memPtr.Offset(length)

	isOffsetTooSmall := memPtr < 0
	isOffsetTooLarge := uint32(memPtr) > memoryLength
	isRequestedEndTooLarge := uint32(requestedEnd) > memoryLength
	isLengthNegative := length < 0

	if isOffsetTooSmall || isOffsetTooLarge {
		return nil, fmt.Errorf("mem load: %w", ErrMemoryBadBounds)
	}
	if isLengthNegative {
		return nil, fmt.Errorf("mem load: %w", ErrMemoryNegativeLength)
	}

	result := make([]byte, length)
	if isRequestedEndTooLarge {
		copy(result, memoryView[memPtr:])
	} else {
		copy(result, memoryView[memPtr:requestedEnd])
	}

	return result, nil
}

// MemStore stores the given data in the WASM memory at the given offset.
func MemStore(memory Memory, offset int32, data []byte) error {
	dataLength := int32(len(data))
	if dataLength == 0 {
		return nil
	}

	memoryView := memory.Data()
	memoryLength := memory.Length()
	requestedEnd := math.AddInt32(offset, dataLength)

	isOffsetTooSmall := offset < 0
	isNewPageNecessary := uint32(requestedEnd) > memoryLength

	if isOffsetTooSmall {
		return ErrMemoryBadBoundsLower
	}
	if isNewPageNecessary {
		err := memory.Grow(1)
		if err != nil {
			return err
		}

		memoryView = memory.Data()
		memoryLength = memory.Length()
	}

	isRequestedEndTooLarge := uint32(requestedEnd) > memoryLength
	if isRequestedEndTooLarge {
		return ErrMemoryBadBoundsUpper
	}

	copy(memoryView[offset:requestedEnd], data)
	return nil
}

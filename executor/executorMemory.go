package executor

import (
	"fmt"

	"github.com/ElrondNetwork/wasm-vm/math"
)

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
func MemLoad(memory Memory, offset int32, length int32) ([]byte, error) {
	if length == 0 {
		return []byte{}, nil
	}

	memoryView := memory.Data()
	memoryLength := memory.Length()
	requestedEnd := math.AddInt32(offset, length)

	isOffsetTooSmall := offset < 0
	isOffsetTooLarge := uint32(offset) > memoryLength
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
		copy(result, memoryView[offset:])
	} else {
		copy(result, memoryView[offset:requestedEnd])
	}

	return result, nil
}

// MemLoadMultiple returns multiple byte slices loaded from the WASM memory, starting at the given offset and having the provided lengths.
func MemLoadMultiple(memory Memory, offset int32, lengths []int32) ([][]byte, error) {
	if len(lengths) == 0 {
		return [][]byte{}, nil
	}

	results := make([][]byte, len(lengths))

	for i, length := range lengths {
		result, err := MemLoad(memory, offset, length)
		if err != nil {
			return nil, err
		}

		results[i] = result
		offset += length
	}

	return results, nil
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

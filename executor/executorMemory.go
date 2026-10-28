package executor

import (
	"fmt"

	"github.com/multiversx/mx-chain-vm-go/math"
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

// MemLoadFromMemory is a bridge to the old Memory interface.
// We are moving away from that, this is to ease the transition.
func MemLoadFromMemory(memory Memory, memPtr MemPtr, length MemLength) ([]byte, error) {
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

	var result []byte
	if isRequestedEndTooLarge {
		result = make([]byte, memoryLength-uint32(memPtr))
		copy(result, memoryView[memPtr:])
	} else {
		result = make([]byte, requestedEnd-memPtr)
		copy(result, memoryView[memPtr:requestedEnd])
	}

	return result, nil
}

// MemStoreToMemory is a bridge to the old Memory interface.
// We are moving away from that, this is to ease the transition.
func MemStoreToMemory(memory Memory, memPtr MemPtr, data []byte) error {
	dataLength := int32(len(data))
	if dataLength == 0 {
		return nil
	}

	memoryView := memory.Data()
	memoryLength := memory.Length()
	requestedEnd := memPtr.Offset(dataLength)

	isOffsetTooSmall := memPtr < 0
	if isOffsetTooSmall {
		return ErrMemoryBadBoundsLower
	}

	isRequestedEndTooLarge := uint32(requestedEnd) > memoryLength
	if isRequestedEndTooLarge {
		return ErrMemoryBadBoundsUpper
	}

	copy(memoryView[memPtr:requestedEnd], data)
	return nil
}

// MemoryLoadFromMemory is the strict load: a range that does not fit entirely inside the
// memory is rejected rather than truncated.
//
// Named after InstanceState.MemoryLoad, where MemLoadFromMemory is named after the old
// Instance.MemLoad that truncates. Two functions that differ only in behaviour must not also
// agree in name.
//
// Counterpart of the Rust checked_load_range (new_traits/memory.rs). The bounds rules are
// observable by contracts, so there is one implementation of them and not one per backend.
// Two asymmetries are deliberate, and shared with the clamped load: a zero length returns
// empty before the offset is looked at, and an offset exactly at the end of the memory is a
// valid offset that nothing can be read from.
func MemoryLoadFromMemory(memory Memory, memPtr MemPtr, length MemLength) ([]byte, error) {
	if length == 0 {
		return []byte{}, nil
	}

	memoryLength := uint64(memory.Length())

	if memPtr < 0 || uint64(memPtr) > memoryLength {
		return nil, fmt.Errorf("mem load: %w", ErrMemoryBadBounds)
	}
	if length < 0 {
		return nil, fmt.Errorf("mem load: %w", ErrMemoryNegativeLength)
	}

	// Added in uint64 rather than through MemPtr.Offset, so that a length a contract is free
	// to choose cannot wrap the sum into a range that looks valid.
	requestedEnd := uint64(memPtr) + uint64(length)
	if requestedEnd > memoryLength {
		return nil, fmt.Errorf("mem load: %w", ErrMemoryBadBoundsUpper)
	}

	result := make([]byte, length)
	copy(result, memory.Data()[memPtr:requestedEnd])

	return result, nil
}

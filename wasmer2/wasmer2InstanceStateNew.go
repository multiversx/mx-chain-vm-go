package wasmer2

import (
	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ executor.InstanceState = (*Wasmer2InstanceStateNew)(nil)

// wasmPageSize is the size in bytes of a WASM linear memory page, as vmhost.WASMPageSize
// defines it. Repeated rather than imported, vmhost/hostCore importing this package.
const wasmPageSize = uint64(65536)

// Wasmer2InstanceStateNew is the inside view of a Wasmer 2.2 instance: what a VM hook can
// reach while it runs.
//
// Counterpart of the Rust WasmerProdInstanceState. It is a handle to the instance rather
// than a snapshot of it, so that memory and gas are read when the hook asks, not when the
// state was built.
//
// Wasmer 2.2 would let this live as long as the instance does, but executor.InstanceState
// promises less than that on purpose - under Wasmer 7 such a handle is valid for exactly one
// hook invocation. Nothing here depends on the longer lifetime, and a holder must not assume
// it either: an instance that has been closed is reported as detached rather than followed.
type Wasmer2InstanceStateNew struct {
	instance *Wasmer2InstanceNew
}

// liveInstance returns the instance this state belongs to, or an error if there is none to
// speak of: the state was never attached to one, or the instance has since been closed.
func (state *Wasmer2InstanceStateNew) liveInstance() (*Wasmer2InstanceNew, error) {
	if state == nil || state.instance == nil || state.instance.cgoInstance == nil {
		return nil, ErrDetachedInstanceState
	}

	return state.instance, nil
}

// GetPointsUsedNew returns the number of points (gas) used so far by the running call.
func (state *Wasmer2InstanceStateNew) GetPointsUsedNew() (uint64, error) {
	instance, err := state.liveInstance()
	if err != nil {
		return 0, err
	}

	return cWasmerInstanceGetPointsUsed(instance.cgoInstance), nil
}

// SetPointsUsedNew overwrites the points (gas) used so far, which is how a hook charges for
// its own work.
func (state *Wasmer2InstanceStateNew) SetPointsUsedNew(points uint64) error {
	instance, err := state.liveInstance()
	if err != nil {
		return err
	}

	if cWasmerInstanceSetPointsUsed(instance.cgoInstance, points) != cWasmerOk {
		return newWrappedError(ErrSetPointsUsed)
	}

	return nil
}

// MemoryLength returns the size in bytes of the instance memory.
func (state *Wasmer2InstanceStateNew) MemoryLength() (uint64, error) {
	instance, err := state.liveInstance()
	if err != nil {
		return 0, err
	}

	return uint64(instance.memory.Length()), nil
}

// MemoryGrow grows the instance memory by the given number of pages, returning the number of
// pages the memory had before growing.
//
// The C API reports only success or failure, so the previous size is read beforehand rather
// than returned by the growth itself.
func (state *Wasmer2InstanceStateNew) MemoryGrow(byNumPages uint32) (uint32, error) {
	instance, err := state.liveInstance()
	if err != nil {
		return 0, err
	}

	pagesBefore := uint32(uint64(instance.memory.Length()) / wasmPageSize)

	err = instance.memory.Grow(byNumPages)
	if err != nil {
		return 0, err
	}

	return pagesBefore, nil
}

// MemoryLoad copies data out of the instance memory, rejecting a range that does not fit
// entirely inside it.
func (state *Wasmer2InstanceStateNew) MemoryLoad(memPtr executor.MemPtr, length executor.MemLength) ([]byte, error) {
	instance, err := state.liveInstance()
	if err != nil {
		return nil, err
	}

	return executor.MemoryLoadFromMemory(&instance.memory, memPtr, length)
}

// MemoryLoadClamped copies data out of the instance memory, truncating the range at the end
// of the memory instead of rejecting it.
//
// The same helper the old Instance.MemLoad uses, which is the point: this is the behaviour
// contracts have already observed, and it has to stay bit for bit what it was.
func (state *Wasmer2InstanceStateNew) MemoryLoadClamped(memPtr executor.MemPtr, length executor.MemLength) ([]byte, error) {
	instance, err := state.liveInstance()
	if err != nil {
		return nil, err
	}

	return executor.MemLoadFromMemory(&instance.memory, memPtr, length)
}

// MemoryStore writes to the instance memory. Never truncated.
func (state *Wasmer2InstanceStateNew) MemoryStore(memPtr executor.MemPtr, data []byte) error {
	instance, err := state.liveInstance()
	if err != nil {
		return err
	}

	return executor.MemStoreToMemory(&instance.memory, memPtr, data)
}

// IsInterfaceNil returns true if underlying object is nil
func (state *Wasmer2InstanceStateNew) IsInterfaceNil() bool {
	return state == nil
}

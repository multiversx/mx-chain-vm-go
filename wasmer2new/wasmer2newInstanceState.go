package wasmer2new

// #include <stdlib.h>
// #include "./libvmexeccapi_w2.h"
import "C"

import (
	"fmt"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ executor.InstanceState = (*Wasmer2InstanceState)(nil)

// Wasmer2InstanceState is the inside view of an instance: what a VM hook can reach while
// it runs.
//
// It is the reason executor.InstanceState exists. Under Wasmer 7 an instance's memory and
// globals are only reachable through the store borrow inside the host-call frame; Wasmer 2.2
// is not under that constraint, but it is held to the same shape anyway, because the C API is
// the same for both backends and a difference between them must never be a difference in what
// the boundary offers.
//
// Where package wasmer2 has to make that promise by convention, here it is enforced: every
// operation below carries the hook context, whose call token the library checks, so a handle
// kept past its call comes back as wasmercommon.ErrStaleContext rather than as a
// use-after-free.
//
// It must therefore not be stored, and must not be touched from a goroutine spawned inside
// the hook: cgo keeps the hook itself on one OS thread, but nothing keeps a goroutine it
// starts there, and the token check is not a lock.
type Wasmer2InstanceState struct {
	ctx   *cHookContextT
	owner *instanceHooks
}

// GetPointsUsedNew returns the number of points (gas) used so far by the running call.
func (state *Wasmer2InstanceState) GetPointsUsedNew() (uint64, error) {
	var points C.uint64_t

	status := cStateGetPointsUsed(state.ctx, &points, state.owner.stateErrBuf)
	if err := state.owner.stateErrBuf.AsError(status, "could not read points used"); err != nil {
		return 0, err
	}

	return uint64(points), nil
}

// SetPointsUsedNew overwrites the points (gas) used so far, which is how a hook charges for
// its own work.
func (state *Wasmer2InstanceState) SetPointsUsedNew(points uint64) error {
	status := cStateSetPointsUsed(state.ctx, points, state.owner.stateErrBuf)

	return state.owner.stateErrBuf.AsError(status, "could not set points used")
}

// MemoryLength returns the size in bytes of the instance memory.
func (state *Wasmer2InstanceState) MemoryLength() (uint64, error) {
	var length C.uint64_t

	status := cStateMemLength(state.ctx, &length, state.owner.stateErrBuf)
	if err := state.owner.stateErrBuf.AsError(status, "could not read memory length"); err != nil {
		return 0, err
	}

	return uint64(length), nil
}

// MemoryGrow grows the instance memory by the given number of pages, returning the number of
// pages the memory had before growing.
func (state *Wasmer2InstanceState) MemoryGrow(byNumPages uint32) (uint32, error) {
	var prevPages C.uint32_t

	status := cStateMemGrow(state.ctx, byNumPages, &prevPages, state.owner.stateErrBuf)
	if err := state.owner.stateErrBuf.AsError(status, "could not grow memory"); err != nil {
		return 0, err
	}

	return uint32(prevPages), nil
}

// MemoryLoad copies data out of the instance memory, rejecting a range that does not fit
// entirely inside it.
//
// Built on the clamped load, the library offering only that one: the two differ in nothing
// but the end-of-range check, and a short read is exactly the case the strict load rejects.
// The offset and length rules are the library's, so they cannot drift from the clamped load's.
func (state *Wasmer2InstanceState) MemoryLoad(memPtr executor.MemPtr, length executor.MemLength) ([]byte, error) {
	loaded, err := state.MemoryLoadClamped(memPtr, length)
	if err != nil {
		return nil, err
	}

	if executor.MemLength(len(loaded)) != length {
		return nil, fmt.Errorf("mem load: %w", executor.ErrMemoryBadBoundsUpper)
	}

	return loaded, nil
}

// MemoryLoadClamped copies data out of the instance memory, truncating the range at the end
// of the memory instead of rejecting it.
//
// The truncation is deliberate and is where the Go VM's historical load semantics live: a
// load whose offset is valid but whose end lies past the end of the memory yields the bytes
// that are available, and no error. Contracts have observed that behaviour, so it is kept.
// An offset outside the memory, or a negative length, is still an error.
//
// The returned slice is Go-owned and may be retained; the bytes were copied out of wasm
// memory during the call.
func (state *Wasmer2InstanceState) MemoryLoadClamped(memPtr executor.MemPtr, length executor.MemLength) ([]byte, error) {
	if length == 0 {
		return []byte{}, nil
	}
	if length < 0 {
		// Checked here only because make() would panic on it. The rule itself is Rust's -
		// clamped_load_len is its single definition - and rejects the same input with the
		// same message, so this guard changes no behaviour, it just reaches the allocation
		// in one piece.
		return nil, executor.ErrMemoryNegativeLength
	}

	dest := make([]byte, length)
	loaded := cStateMemLoad(state.ctx, int32(memPtr), dest, state.owner.stateErrBuf)
	if loaded < 0 {
		return nil, state.owner.stateErrBuf.AsErrorFromCode(loaded, "memory load failed")
	}

	return dest[:loaded], nil
}

// MemoryStore writes to the instance memory. Never truncated: writing past the end of the
// memory is an error, because silently dropping a contract's write would corrupt its state
// rather than merely shorten a read.
func (state *Wasmer2InstanceState) MemoryStore(memPtr executor.MemPtr, data []byte) error {
	if len(data) == 0 {
		return nil
	}

	status := cStateMemStore(state.ctx, int32(memPtr), data, state.owner.stateErrBuf)

	return state.owner.stateErrBuf.AsError(status, "memory store failed")
}

// IsInterfaceNil returns true if underlying object is nil
func (state *Wasmer2InstanceState) IsInterfaceNil() bool {
	return state == nil
}

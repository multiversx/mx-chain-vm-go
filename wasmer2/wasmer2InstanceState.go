package wasmer2

// #include <stdlib.h>
// #include "./libvmexeccapi_w2.h"
import "C"

import (
	"github.com/multiversx/mx-chain-vm-go/executor"
)

// InstanceState is the inside view of an instance: what a VM hook can reach while it runs.
//
// It is the Go counterpart of the Rust InstanceState trait, and the reason the trait exists.
// The trait exists because Wasmer 6 can only reach an instance's memory and globals through
// the store borrow inside the host call frame. Wasmer 2.2 is not under that constraint, but it
// is held to the same shape anyway: the C API is the same for both backends, so one Go caller
// serves either, and a difference between them is never a difference in what the boundary
// offers. In Rust the borrow checker enforces the lifetime; here the enforcement is the hook
// context's call token, checked on every operation below, which turns a handle kept past its
// call into an ErrStaleContext rather than a use-after-free.
//
// An InstanceState must therefore not be stored, and must not be touched from a goroutine
// spawned inside the hook: cgo keeps the hook itself on one OS thread, but nothing keeps a
// goroutine it starts there, and the token check is not a lock.
type InstanceState struct {
	ctx   *cHookContextT
	owner *instanceHooks
}

// MemLoad reads from instance memory, truncating at the end of the memory.
//
// The truncation is deliberate and is where the Go VM's historical load semantics live: a
// load whose offset is valid but whose end lies past the end of the memory yields the bytes
// that are available, and no error. Contracts have observed that behaviour, so it is kept.
// An offset outside the memory, or a negative length, is still an error.
//
// The returned slice is Go-owned and may be retained; the bytes were copied out of wasm
// memory during the call.
func (state *InstanceState) MemLoad(
	memPtr executor.MemPtr,
	length executor.MemLength,
) ([]byte, error) {
	if length == 0 {
		return []byte{}, nil
	}
	if length < 0 {
		// Checked here only because make() would panic on it. The rule itself is Rust's —
		// clamped_load_len is its single definition — and rejects the same input with the
		// same message, so this guard changes no behaviour, it just reaches the allocation
		// in one piece.
		return nil, executor.ErrMemoryNegativeLength
	}

	dest := make([]byte, length)
	loaded := cStateMemLoad(state.ctx, int32(memPtr), dest, state.owner.stateErrBuf)
	if loaded < 0 {
		return nil, state.owner.stateErrBuf.AsErrorFromCode(int64(loaded), "memory load failed")
	}

	return dest[:loaded], nil
}

// MemStore writes to instance memory. Unlike a load, a store is never truncated: writing
// past the end of the memory is an error, because silently dropping a contract's write
// would corrupt its state rather than merely shorten a read.
func (state *InstanceState) MemStore(memPtr executor.MemPtr, data []byte) error {
	if len(data) == 0 {
		return nil
	}

	status := cStateMemStore(state.ctx, int32(memPtr), data, state.owner.stateErrBuf)

	return state.owner.stateErrBuf.AsError(status, "memory store failed")
}

// MemLength returns the size of the instance memory in bytes.
func (state *InstanceState) MemLength() (uint64, error) {
	var length C.uint64_t

	status := cStateMemLength(state.ctx, &length, state.owner.stateErrBuf)
	if err := state.owner.stateErrBuf.AsError(status, "could not read memory length"); err != nil {
		return 0, err
	}

	return uint64(length), nil
}

// MemGrow grows the instance memory by whole 64Kb pages, returning the number of pages it
// had before growing.
func (state *InstanceState) MemGrow(pages uint32) (uint32, error) {
	var prevPages C.uint32_t

	status := cStateMemGrow(state.ctx, pages, &prevPages, state.owner.stateErrBuf)
	if err := state.owner.stateErrBuf.AsError(status, "could not grow memory"); err != nil {
		return 0, err
	}

	return uint32(prevPages), nil
}

// GetPointsUsed returns the gas consumed so far by the running call.
func (state *InstanceState) GetPointsUsed() (uint64, error) {
	var points C.uint64_t

	status := cStateGetPointsUsed(state.ctx, &points, state.owner.stateErrBuf)
	if err := state.owner.stateErrBuf.AsError(status, "could not read points used"); err != nil {
		return 0, err
	}

	return uint64(points), nil
}

// SetPointsUsed overwrites the gas consumed so far, which is how a hook charges for its own
// work.
func (state *InstanceState) SetPointsUsed(points uint64) error {
	status := cStateSetPointsUsed(state.ctx, points, state.owner.stateErrBuf)

	return state.owner.stateErrBuf.AsError(status, "could not set points used")
}

package wasmer7

import "github.com/multiversx/mx-chain-vm-go/executor"

// This backend uses the call-scoped half of executor.VMHooksProvider, and only that half: the
// hook context is what makes the instance state reachable, and it is valid for exactly the
// call it was handed to. VMHooksForInstanceLifetime is never called here, because there is no
// such lifetime to offer.
//
// The asymmetry is the whole reason the second set of executor traits exists; see
// executor.InstanceState.

var _ = (executor.InstanceState)((*InstanceState)(nil))

// instanceHooks is what the instance's cgo.Handle resolves to.
//
// One of these per live instance, created with the instance and released with it.
type instanceHooks struct {
	provider executor.VMHooksProvider

	// activeState is the state of the hook currently running on this instance, or nil.
	//
	// It exists only because the VM still reaches wasm memory through executor.Instance
	// rather than through the state it was handed — see PLAN-go-ffi.md, Phase 4b item 1. Once
	// Instance is shrunk to the outside view, this goes away with it.
	//
	// It is never cleared, and does not need to be. A context used after its call has ended
	// fails the executor's token check and comes back as ErrStaleContext, which is the whole
	// point of that check; and the slot cannot be clobbered mid-call, because a nested call
	// runs on a different instance and re-entering the same one is refused outright.
	activeState *InstanceState

	// stateErrBuf is shared by every state operation this instance's hooks perform.
	//
	// Sharing is safe, and not merely convenient: a state operation is a single cgo call
	// into Rust that reads or writes wasm memory and returns. It never calls back into Go,
	// so the buffer's live range — the call, plus reading the message out of it — contains
	// no re-entry. The alternative, a buffer per state operation, would be an allocation on
	// every memory access a contract makes.
	stateErrBuf errBuf
}

package wasmer3

import "github.com/multiversx/mx-chain-vm-go/executor"

// VMHooksProvider hands out the VM hooks that service one host function call.
//
// The indirection is what the Wasmer 6 ownership model forces, and it mirrors the Rust side
// exactly: there, ExperimentalVMHooksBuilder::with_vm_hooks builds the hooks around the
// instance state for the duration of one call and lends them out. Here the provider is asked
// for hooks bound to this call's InstanceState.
//
// Why not a plain executor.VMHooks set once per instance, as the legacy executor did? Because
// the state the hooks need is only valid for one call. Binding it at construction would mean
// either keeping a mutable "current state" slot on the hooks object — which an inner call
// re-entering the same instance would clobber — or reaching for the state through a side
// channel. Taking it as a parameter makes the nesting correct by construction, which is the
// same reason the Rust side builds its hooks per call rather than per instance.
type VMHooksProvider interface {
	VMHooksForCall(state *InstanceState) executor.VMHooks
}

// instanceHooks is what the instance's cgo.Handle resolves to.
//
// One of these per live instance, created with the instance and released with it.
type instanceHooks struct {
	provider VMHooksProvider

	// stateErrBuf is shared by every state operation this instance's hooks perform.
	//
	// Sharing is safe, and not merely convenient: a state operation is a single cgo call
	// into Rust that reads or writes wasm memory and returns. It never calls back into Go,
	// so the buffer's live range — the call, plus reading the message out of it — contains
	// no re-entry. The alternative, a buffer per state operation, would be an allocation on
	// every memory access a contract makes.
	stateErrBuf errBuf
}

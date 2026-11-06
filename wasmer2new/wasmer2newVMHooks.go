package wasmer2new

import (
	"github.com/multiversx/mx-chain-vm-go/executor"
)

// instanceHooks is what an instance's cgo.Handle resolves to: everything a hook of that
// instance needs, and nothing that varies between hooks.
//
// One of these per live instance, created with the instance and released with it. It holds
// the runtime rather than a hooks object, because the hooks themselves are built per hook,
// around the state of the call that is running - see executor.ExecutorRuntime.
//
// This is the same arrangement as the legacy backend's, reached differently. There, the
// instance is identified by which adapter sits in the executor's one hooks slot; here, by
// the handle the executor hands back in every hook context. Neither keeps a record of which
// instance is running.
type instanceHooks struct {
	runtime executor.ExecutorRuntime

	// stateErrBuf is shared by every state operation this instance's hooks perform.
	//
	// Sharing is safe, and not merely convenient: a state operation is a single cgo call
	// into Rust that reads or writes wasm memory and returns. It never calls back into Go,
	// so the buffer's live range - the call, plus reading the message out of it - contains
	// no re-entry. The alternative, a buffer per state operation, would be an allocation on
	// every memory access a contract makes.
	stateErrBuf errBuf
}

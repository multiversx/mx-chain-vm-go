package executor

import "github.com/multiversx/mx-chain-core-go/core/check"

// InstanceState is the inside view of an instance: what a VM hook can reach while it runs.
//
// Counterpart of the Rust `InstanceState` trait (new_traits/instance_state.rs). There is no
// old interface of this name — today the VM hooks reach the instance through the host's
// instance stack instead, and this is what replaces that.
//
// It exists because the backends disagree about how long such a handle may live, and the
// stricter of the two has to win. Under Wasmer 2.2 an instance's memory is reachable for the
// instance's whole lifetime. Under Wasmer 7 it is only reachable through the store borrow
// inside the host-call frame, so it is valid for exactly one hook invocation and not a moment
// longer. Writing the VM against this interface means writing it against the stricter
// contract, which 2.2 satisfies for free, and which stops the VM from accidentally depending
// on a lifetime only 2.2 offers.
//
// An implementation must therefore not be stored past the hook it was handed to. A backend
// whose state is genuinely call-scoped rejects a stale handle rather than following it.
//
// All accesses are validated by the executor, which is the only side that knows the real
// extent of the instance memory. The memory methods are named after the Rust ones rather
// than after Instance.MemLoad and friends, on purpose: MemoryLoad is strict where
// Instance.MemLoad truncates, and two methods that differ only in behaviour must not also
// agree in name.
type InstanceState interface {
	check.NilInterfaceChecker

	// GetPointsUsedNew returns the number of points (gas) used so far by the running call.
	//
	// New suffix: Instance.GetPointsUsed has no error result.
	GetPointsUsedNew() (uint64, error)

	// SetPointsUsedNew overwrites the points (gas) used so far, which is how a hook charges
	// for its own work.
	//
	// New suffix: Instance.SetPointsUsed has no error result.
	SetPointsUsedNew(points uint64) error

	// MemoryLength returns the size in bytes of the instance memory.
	MemoryLength() (uint64, error)

	// MemoryGrow grows the instance memory by the given number of pages (of 64Kb each),
	// returning the number of pages the memory had before growing.
	MemoryGrow(byNumPages uint32) (uint32, error)

	// MemoryLoad copies data out of the instance memory.
	//
	// Strict: a range that does not fit entirely inside the memory is rejected.
	MemoryLoad(memPtr MemPtr, length MemLength) ([]byte, error)

	// MemoryLoadClamped copies data out of the instance memory, truncating the range at the
	// end of the memory instead of rejecting it.
	//
	// This reproduces, deliberately, the behaviour the Go VM has always had: a load whose
	// offset is valid but whose end lies past the end of the memory yields the bytes that
	// are available, and no error. Contracts have observed this, so the executor has to keep
	// offering it, distinct from the strict MemoryLoad.
	//
	// An offset outside the memory, or a negative length, is still an error.
	MemoryLoadClamped(memPtr MemPtr, length MemLength) ([]byte, error)

	// MemoryStore writes to the instance memory.
	//
	// Never truncated: writing past the end of the memory is an error, because silently
	// dropping a contract's write would corrupt its state rather than merely shorten a read.
	MemoryStore(memPtr MemPtr, data []byte) error
}

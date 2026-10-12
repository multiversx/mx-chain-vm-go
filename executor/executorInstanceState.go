package executor

// InstanceState is the inside view of an instance: what a VM hook can reach while it runs.
//
// It exists because the two backends disagree about how long such a handle may live, and the
// stricter of the two has to win. Under Wasmer 2.2 an instance's memory is reachable for the
// instance's whole lifetime. Under Wasmer 7 it is only reachable through the store borrow
// inside the host-call frame, so it is valid for exactly one hook invocation and not a moment
// longer. Writing the VM against this interface means writing it against the stricter
// contract, which 2.2 satisfies for free, and which stops the VM from accidentally depending
// on a lifetime only 2.2 offers.
//
// Only what a hook actually needs is here. Gas accounting and memory growth are reached
// through RuntimeContext and MeteringContext, whose callers sit outside a host call as often
// as inside one, so they stay on Instance.
//
// An implementation must not be stored past the hook it was handed to. A backend whose state
// is genuinely call-scoped will reject a stale handle rather than follow it.
type InstanceState interface {
	// MemLoad returns the contents of the instance memory at the given offset, truncating at
	// the end of the memory rather than failing.
	MemLoad(memPtr MemPtr, length MemLength) ([]byte, error)

	// MemStore writes to the instance memory at the given offset. Unlike a load, a store is
	// never truncated: writing past the end of the memory is an error.
	MemStore(memPtr MemPtr, data []byte) error
}

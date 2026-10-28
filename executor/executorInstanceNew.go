package executor

import "github.com/multiversx/mx-chain-core-go/core/check"

// InstanceNew is the view of a running contract from *outside* a host call: calling it, gas
// accounting, module introspection and lifecycle.
//
// Counterpart of the Rust `Instance` trait (new_traits/instance_new.rs). Named with a New
// suffix only because [Instance] is still in use; the rename is postponed.
//
// What a VM hook needs while a call is running is not here — that is [InstanceState],
// reached through the hook's own context. Keeping the two apart is what lets a backend whose
// instance state is only reachable for the duration of a call implement this interface
// honestly.
//
// Three methods carry a New suffix for a different reason: [Instance] already has a method
// of that name with a different signature, and a Go type cannot have both. They lose the
// suffix once the old interface goes.
type InstanceNew interface {
	check.NilInterfaceChecker

	// Call calls an exported function by name, with the gas it is allowed to burn.
	//
	// The gas limit is an argument rather than instance configuration, which is what
	// Instance.SetGasLimit turns into.
	Call(funcName string, pointsLimit uint64) InstanceCallResult

	// CheckSignatures checks that all public module functions (SC endpoints) have no
	// arguments or results.
	CheckSignatures() bool

	// HasFunction checks whether the SC has an endpoint with the given name.
	HasFunction(funcName string) bool

	// GetExportedFunctionNames returns the names of all the SC endpoints.
	GetExportedFunctionNames() []string

	// GetImportedFunctionNames returns the names of all the host functions (VM hooks) the
	// contract imports.
	//
	// Needed to validate a contract against the set of hooks it is allowed to use. All the
	// names are returned at once, rather than offering a per-name query, because the callers
	// check the module against whole sets of names.
	GetImportedFunctionNames() []string

	// GetPointsUsedNew returns the number of points (gas) used by the instance.
	//
	// New suffix: Instance.GetPointsUsed has no error result.
	GetPointsUsedNew() (uint64, error)

	// ResetNew returns the instance to the state it was in immediately after creation, so
	// that it can serve another execution.
	//
	// Memory contents and length, every mutable global, and anything else a call can leave
	// behind must be restored. What must *not* change is the instance's identity: the host
	// holds this instance across the reset, so a backend has to rebuild in place rather than
	// hand back a new object.
	//
	// A backend that cannot offer this returns an error, and the host is expected to discard
	// the instance and build a fresh one.
	//
	// New suffix: Instance.Reset reports failure as a bool.
	ResetNew() error

	// Cache serializes the compiled module, for [ExecutorNew.NewInstanceFromCache].
	Cache() ([]byte, error)

	// RequestStop asks a call that is running on this instance to stop at the next
	// opportunity.
	//
	// The one method here that may be called from another goroutine, and while [Call] is in
	// flight - which is the only moment it is useful. It is what an execution deadline is
	// enforced with.
	//
	// Best effort, and deliberately silent about whether it worked, for two reasons no
	// backend can do anything about: a contract already past its last checkpoint may finish
	// first, and neither backend can interrupt a contract that is inside a VM hook, where
	// the stop is only noticed once control returns to wasm. A caller that must bound the
	// wait needs its own timeout regardless.
	//
	// Safe at any time, including with no call running or after the call has finished: both
	// are no-ops. The call it stops reports [CallInterrupted].
	RequestStop()

	// ID returns an identifier for the instance, unique at runtime.
	//
	// On the interface because a host that caches instances needs to key them by something,
	// and because every trace line about an instance names it.
	ID() string

	// Close releases the instance.
	//
	// Go has no destructor, so this is the only thing that frees what a backend holds
	// outside the Go heap - C memory, and for a backend reached over cgo the handle its
	// hooks are found through. An instance that is never closed leaks both.
	//
	// The Rust Instance trait has no counterpart, and needs none: there, Drop does this.
	Close()
}

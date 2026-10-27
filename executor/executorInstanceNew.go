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
}

package executor

import "github.com/multiversx/mx-chain-core-go/core/check"

// ExecutorRuntime is the host behind an executor: one reference, for the executor's whole
// life.
//
// Counterpart of the Rust `ExecutorRuntime` trait (new_traits/executor_runtime.rs). This is
// the only thing an executor needs from its host that is not in the wasm and not in the
// compilation options. It is supplied once, at executor construction, and never replaced —
// which is what lets an executor be configured at once and then treated as immutable.
type ExecutorRuntime interface {
	check.NilInterfaceChecker

	// VMHooksForCall returns the VM hooks servicing one host call, built around the instance
	// state for that call.
	//
	// Built per call, not per instance. Under Wasmer 7 the instance state is only reachable
	// for the duration of a call — memory and globals live behind the store borrow the host
	// function is handed — so hooks that hold it cannot outlive the call. The result must
	// therefore not be retained past the hook that asked for it.
	//
	// Which execution the hooks serve is resolved here, by the host, from its own call stack
	// — not captured when the instance was created. That is the property that makes a cached
	// instance safe to reuse across transactions: the instance contributes the state, never
	// the identity.
	//
	// The state arrives as an [InstanceState] so that one runtime serves every backend.
	//
	// The Rust trait lends the hooks to a closure rather than returning them, to avoid an
	// allocation per hook invocation and to bound their lifetime. Neither reason applies
	// here, so this returns them.
	VMHooksForCall(instanceState InstanceState) VMHooks
}

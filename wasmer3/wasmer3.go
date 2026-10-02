// Package wasmer3 is the cgo bridge to the new executor C API (`libvmexeccapi2`), which
// exposes the Wasmer 6 backend through the second set of executor traits.
//
// It is deliberately *not* an implementation of executor.Executor yet. The Go VM still
// speaks the legacy interface set, and converging the two is a separate step — see
// PLAN-go-ffi.md, Phase 4b. What lives here is the boundary itself: the C calls, the hook
// context, and the call-scoped instance state, driven directly by tests rather than by the
// VM. Keeping the VM out of it means a wrong turn at the boundary costs a package instead of
// a refactor.
//
// The two halves of the boundary:
//
//   - Outward calls, Go to Rust: everything in bridge3.go. The VM creates an executor and
//     instances, calls functions, and reads gas.
//   - Inward calls, Rust to Go: the generated shims in wasmer3ImportsCgo.go. Each one is
//     handed a hook context, recovers the VM hooks for the instance from it, and may call
//     back through the context to reach the instance's memory and gas.
package wasmer3

import logger "github.com/multiversx/mx-chain-logger-go"

// Logger for this executor, matching the name wasmer2 uses so that the two are configured
// together.
var logWasmer3 = logger.GetOrCreate("vm/executor")

// Package wasmer2new is the cgo bridge to the new executor C API (`libvmexeccapi_w2`), which
// exposes the Wasmer 2.2 backend through the second set of executor interfaces.
//
// The same wasm runtime as package wasmer2, reached through the new boundary instead of the
// old one. That pairing is the point: with the runtime held constant, any behavioural
// difference between the two is attributable to the migration itself. A
// wasmer7-against-wasmer2 comparison cannot make that claim, because it moves the runtime
// and the boundary at once.
//
// Both packages implement executor.ExecutorNew, and the difference between them is where the
// instance state comes from. Package wasmer2 has to build it on the Go side, the old C API
// knowing nothing of hooks contexts; here the executor hands one to every hook, and the
// lifetime that executor.InstanceState describes is enforced by the library's own token
// check rather than by convention.
//
// The two halves of the boundary:
//
//   - Outward calls, Go to Rust: everything in wasmer2newBridge.go. The VM creates an executor and
//     instances, calls functions, and reads gas.
//   - Inward calls, Rust to Go: the generated shims in wasmer2newImportsCgo.go. Each one is
//     handed a hook context, recovers the VM hooks for the instance from it, and may call
//     back through the context to reach the instance's memory and gas.
package wasmer2new

import logger "github.com/multiversx/mx-chain-logger-go"

// Logger for this executor, matching the name wasmer2 uses so that the two are configured
// together.
var logWasmer2New = logger.GetOrCreate("vm/executor")

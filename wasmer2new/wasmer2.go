// Package wasmer2new is the cgo bridge to the new executor C API (`libvmexeccapi_w2`), which
// exposes the Wasmer 2.2 backend through the second set of executor traits.
//
// The same wasm runtime as package wasmer2, reached through the new boundary instead of
// the old one. That pairing is the point: with the runtime held constant, any behavioural
// difference between the two is attributable to the migration itself, which is what the
// hook-sequence differential in integrationTests/json measures. A wasmer7-against-wasmer2
// comparison cannot make that claim, because it moves the runtime and the boundary at once.
//
// Apart from the backend it selects, this package is a copy of wasmer7 with a different
// symbol prefix — it has to be, because a cgo //export name is a C symbol and two packages
// cannot export the same one. What does not depend on the prefix is in wasmercommon.
//
// The two halves of the boundary:
//
//   - Outward calls, Go to Rust: everything in bridge2.go. The VM creates an executor and
//     instances, calls functions, and reads gas.
//   - Inward calls, Rust to Go: the generated shims in wasmer2ImportsCgo.go. Each one is
//     handed a hook context, recovers the VM hooks for the instance from it, and may call
//     back through the context to reach the instance's memory and gas.
package wasmer2new

import logger "github.com/multiversx/mx-chain-logger-go"

// Logger for this executor, matching the name wasmer2 uses so that the two are configured
// together.
var logWasmer2 = logger.GetOrCreate("vm/executor")

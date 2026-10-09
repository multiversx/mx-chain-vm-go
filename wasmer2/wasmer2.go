// Package wasmer2 is the cgo bridge to the new executor C API (`libvmexeccapi_w2`), which
// exposes the Wasmer 2.2 backend through the second set of executor traits.
//
// The same wasm runtime as package wasmer2legacy, reached through the new boundary instead of
// the old one. That pairing is the point: with the runtime held constant, any behavioural
// difference between the two is attributable to the migration itself, which is what the
// hook-sequence differential in integrationTests/json measures. A wasmer6-against-wasmer2legacy
// comparison cannot make that claim, because it moves the runtime and the boundary at once.
//
// Apart from the backend it selects, this package is a copy of wasmer6 with a different
// symbol prefix — it has to be, because a cgo //export name is a C symbol and two packages
// cannot export the same one. Folding the shared half into a wasmercommon package is a
// follow-up, deliberately taken after the differential is green rather than before: the whole
// value of this package is as a reference, and refactoring it while it is still being brought
// up would make a failure ambiguous.
//
// The two halves of the boundary:
//
//   - Outward calls, Go to Rust: everything in bridge2.go. The VM creates an executor and
//     instances, calls functions, and reads gas.
//   - Inward calls, Rust to Go: the generated shims in wasmer2ImportsCgo.go. Each one is
//     handed a hook context, recovers the VM hooks for the instance from it, and may call
//     back through the context to reach the instance's memory and gas.
package wasmer2

import logger "github.com/multiversx/mx-chain-logger-go"

// Logger for this executor, matching the name wasmer2 uses so that the two are configured
// together.
var logWasmer2 = logger.GetOrCreate("vm/executor")

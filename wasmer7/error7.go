package wasmer7

// #include <stdlib.h>
// #include "./libvmexeccapi_w7.h"
import "C"

import (
	"github.com/multiversx/mx-chain-vm-go/wasmercommon"
)

// The error vocabulary itself is in wasmercommon, shared with the other backend. What has to
// stay here is the part that touches this library's own C types: the recommended buffer
// length, which is a symbol in this header, and the cast of the buffer to this package's
// *C.char.

// errBufLen is the error buffer size every fallible call uses.
//
// Taken from the library rather than picked here, so that the two cannot drift: the Rust
// side sizes it to hold every message the executor currently produces, and truncates
// silently beyond it.
var errBufLen = C.int(C.VM_EXEC_ERROR_BUF_RECOMMENDED_LEN)

type errBuf = wasmercommon.ErrBuf

func newErrBuf() errBuf {
	return wasmercommon.NewErrBuf(int(errBufLen))
}

// bufPtr hands the buffer to this package's C. A method would be tidier, but Go does not
// allow one on a type from another package, and the alternative — a local defined type —
// would hide wasmercommon's own methods.
func bufPtr(buf errBuf) *C.char {
	return (*C.char)(buf.Ptr())
}

// toStatus converts this library's C status to the shared one. The bridge below is the one
// place the conversion happens, so nothing above it deals in C types.
func toStatus(s C.vm_exec_status_t) wasmercommon.Status {
	return wasmercommon.Status(s)
}

// Compile-time check that wasmercommon's hand-written mirror of vm_exec_status_t still matches
// this library's header.
//
// wasmercommon has no cgo, deliberately: the recommended buffer length above is an extern
// symbol, so reading anything out of the header there would mean linking a backend library —
// the one dependency a backend-neutral package must not have. The price is a mirror, and this
// is what keeps it honest. Each difference makes one of the two subtractions negative, and a
// negative constant converted to uint does not compile.
const (
	_ = uint(int64(C.VM_EXEC_OK) - int64(wasmercommon.StatusOk))
	_ = uint(int64(wasmercommon.StatusOk) - int64(C.VM_EXEC_OK))

	_ = uint(int64(C.VM_EXEC_ERROR) - int64(wasmercommon.StatusError))
	_ = uint(int64(wasmercommon.StatusError) - int64(C.VM_EXEC_ERROR))

	_ = uint(int64(C.VM_EXEC_NULL_POINTER) - int64(wasmercommon.StatusNullPointer))
	_ = uint(int64(wasmercommon.StatusNullPointer) - int64(C.VM_EXEC_NULL_POINTER))

	_ = uint(int64(C.VM_EXEC_STALE_CONTEXT) - int64(wasmercommon.StatusStaleContext))
	_ = uint(int64(wasmercommon.StatusStaleContext) - int64(C.VM_EXEC_STALE_CONTEXT))

	_ = uint(int64(C.VM_EXEC_PANIC) - int64(wasmercommon.StatusPanic))
	_ = uint(int64(wasmercommon.StatusPanic) - int64(C.VM_EXEC_PANIC))
)

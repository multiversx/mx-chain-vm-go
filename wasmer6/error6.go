package wasmer6

// #include <stdlib.h>
// #include "./libvmexeccapi_w6.h"
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"unsafe"
)

var (
	// ErrNullPointer is reported when the executor rejects an argument that must not be
	// null. It carries no message, because the null argument may be the error buffer.
	ErrNullPointer = errors.New("null pointer passed to the executor")

	// ErrStaleContext is reported when a hook context is used outside the host call it
	// belongs to. It is the dynamic replacement for the borrow the Rust side cannot
	// express across C; see the hook context docs in libvmexeccapi_w6.h.
	ErrStaleContext = errors.New("hook context used outside its host call")

	// ErrPanic is reported when a Rust panic was caught at the boundary. A panic must
	// never unwind into C, so every entry point traps it and reports this instead.
	ErrPanic = errors.New("panic caught at the executor boundary")

	// ErrInvalidBytecode is reported for empty contract code.
	ErrInvalidBytecode = errors.New("invalid bytecode")

	// ErrBufferTooSmall is reported when a two-step string-out call rejects the buffer
	// sized by its length call, which can only mean the two disagreed.
	ErrBufferTooSmall = errors.New("destination buffer rejected by the executor")

	// ErrMissingOpcodeCost is reported for a configuration with no cost table. Rejected
	// here rather than passed on, because metering is always on and a missing table would
	// mean a contract running for free.
	ErrMissingOpcodeCost = errors.New("executor config has no opcode cost table")

	// ErrCacheUnsupported is reported for the instance cache path, which this executor does
	// not implement yet. See PLAN-go-ffi.md, Phase 6.
	ErrCacheUnsupported = errors.New("instance caching is not supported by this executor")

	// ErrInstanceReentered is reported when a hook calls back into the instance it is
	// running in. Nesting into another instance is supported; nesting into the same one is
	// not, and this is where that becomes an error rather than undefined behaviour.
	ErrInstanceReentered = errors.New("instance called while a call is already in progress")
)

// errBufLen is the error buffer size every fallible call uses.
//
// Taken from the library rather than picked here, so that the two cannot drift: the Rust
// side sizes it to hold every message the executor currently produces, and truncates
// silently beyond it.
var errBufLen = C.int(C.VM_EXEC_W6_ERROR_BUF_RECOMMENDED_LEN)

// errBuf is the caller-provided buffer every fallible call writes its message into.
//
// The legacy C API instead kept the last error in a thread-local that the caller read back
// in a second call. That does not survive the move to this interface: nothing pins a
// goroutine to one OS thread between two cgo calls, so the second call could land on a
// thread holding someone else's error.
type errBuf []byte

func newErrBuf() errBuf {
	return make(errBuf, int(errBufLen))
}

// ptr hands the buffer to C. Passing a Go pointer for the duration of a call is allowed —
// the rule cgo enforces is against C *storing* one — and this buffer holds no Go pointers.
func (buf errBuf) ptr() *C.char {
	return (*C.char)(unsafe.Pointer(&buf[0]))
}

// message reads back what the executor wrote, up to the NUL terminator.
func (buf errBuf) message() string {
	end := bytes.IndexByte(buf, 0)
	if end < 0 {
		end = len(buf)
	}
	return string(buf[:end])
}

// asError turns a status into an error, reading the message out of buf for the one status
// that produces one.
//
// Mapping the non-message statuses to distinct sentinels rather than to the message channel
// is what lets a caller use errors.Is on them: a stale context in particular is a bug in
// the caller's lifetime handling, not a contract failure, and worth telling apart.
func (buf errBuf) asError(status C.enum_vm_exec_w6_status_t, context string) error {
	switch status {
	case C.VM_EXEC_W6_OK:
		return nil
	case C.VM_EXEC_W6_NULL_POINTER:
		return fmt.Errorf("%s: %w", context, ErrNullPointer)
	case C.VM_EXEC_W6_STALE_CONTEXT:
		return fmt.Errorf("%s: %w", context, ErrStaleContext)
	case C.VM_EXEC_W6_PANIC:
		return fmt.Errorf("%s: %w", context, ErrPanic)
	case C.VM_EXEC_W6_ERROR:
		return fmt.Errorf("%s: %s", context, buf.message())
	default:
		return fmt.Errorf("%s: unknown executor status %d", context, int(status))
	}
}

// asErrorFromCode is asError for the calls that return a count, where a negative result
// carries a status in the same code space.
func (buf errBuf) asErrorFromCode(code int64, context string) error {
	return buf.asError(C.enum_vm_exec_w6_status_t(code), context)
}

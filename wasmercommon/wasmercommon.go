// Package wasmercommon is the Go side of the executor crate `c-api-common`: everything the
// new executor C APIs share, independent of which wasm runtime is underneath.
//
// There is one Go package per backend — wasmer2 and wasmer7 — because a cgo //export name is
// a C symbol and two packages cannot export the same one. They are otherwise identical, and
// what can be shared is shared here.
//
// What cannot: anything whose signature mentions a C type. cgo types are package-local, so
// `C.vm_exec_status_t` in package wasmer7 is a different Go type from the identically-named
// one in wasmer2, and a function taking either could not serve both. That is why this package
// has no cgo at all and deals in plain Go types, with each backend converting at its own
// boundary — in its bridge file, which is where the conversion belongs anyway.
//
// It is also why Status below is mirrored by hand rather than read from the header: the
// recommended error-buffer length is an `extern const int`, a symbol, so a cgo package would
// have to link one of the backend libraries to read anything from the header — which is
// exactly the dependency this package exists not to have. Each backend asserts the mirror
// against its own header instead; see TestStatusMirrorsHeader.
package wasmercommon

import (
	"bytes"
	"errors"
	"fmt"
	"unsafe"
)

// Status mirrors vm_exec_status_t: how a call to the executor ended.
//
// VM_EXEC_OK is zero and every failure is negative, so the functions that return a count (the
// memory loads) share the same code space: a non-negative result is a count, a negative one is
// one of these.
type Status int32

const (
	// StatusOk means the call succeeded.
	StatusOk Status = 0

	// StatusError means the call failed and wrote a message to the error buffer.
	StatusError Status = -1

	// StatusNullPointer means a pointer argument that must not be null was null. No message
	// is produced, because a null pointer may well be the error buffer itself.
	StatusNullPointer Status = -2

	// StatusStaleContext means a hook context was used outside the host call it belongs to.
	StatusStaleContext Status = -3

	// StatusPanic means a Rust panic was caught at the boundary.
	StatusPanic Status = -4
)

var (
	// ErrNullPointer is reported when the executor rejects an argument that must not be
	// null. It carries no message, because the null argument may be the error buffer.
	ErrNullPointer = errors.New("null pointer passed to the executor")

	// ErrStaleContext is reported when a hook context is used outside the host call it
	// belongs to. It is the dynamic replacement for the borrow the Rust side cannot
	// express across C; see the hook context docs in the generated header.
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
	// rather than passed on, because metering is always on and a missing table would mean a
	// contract running for free.
	ErrMissingOpcodeCost = errors.New("executor config has no opcode cost table")

	// ErrEmptyCacheArtifact is reported when the compiled-code path is handed nothing. Caught
	// here rather than at the boundary so that an empty cache entry reads as a cache miss the
	// caller can recover from, rather than as a deserialization failure deep in the library.
	ErrEmptyCacheArtifact = errors.New("cache artifact is empty")

	// ErrInstanceReentered is reported when a hook calls back into the instance it is
	// running in. Nesting into another instance is supported; nesting into the same one is
	// not, and this is where that becomes an error rather than undefined behaviour.
	ErrInstanceReentered = errors.New("instance called while a call is already in progress")
)

// ErrBuf is the caller-provided buffer every fallible call writes its message into.
//
// The legacy C API instead kept the last error in a thread-local that the caller read back
// in a second call. That does not survive the move to this interface: nothing pins a
// goroutine to one OS thread between two cgo calls, so the second call could land on a
// thread holding someone else's error.
type ErrBuf []byte

// NewErrBuf allocates a buffer of the length the executor recommends.
//
// The length is taken from the library rather than picked here, so that the two cannot drift:
// the Rust side sizes it to hold every message the executor currently produces, and truncates
// silently beyond it. It is passed in because it is an `extern const int` in the header, which
// only a package that links a backend library can read.
func NewErrBuf(length int) ErrBuf {
	return make(ErrBuf, length)
}

// Ptr hands the buffer to C. Passing a Go pointer for the duration of a call is allowed —
// the rule cgo enforces is against C *storing* one — and this buffer holds no Go pointers.
//
// Returned as unsafe.Pointer rather than as *C.char, because this package has no cgo; the
// caller casts it to its own *C.char.
func (buf ErrBuf) Ptr() unsafe.Pointer {
	return unsafe.Pointer(&buf[0])
}

// Message reads back what the executor wrote, up to the NUL terminator.
func (buf ErrBuf) Message() string {
	end := bytes.IndexByte(buf, 0)
	if end < 0 {
		end = len(buf)
	}

	return string(buf[:end])
}

// AsError turns a status into an error, reading the message out of the buffer for the one
// status that produces one.
//
// Mapping the non-message statuses to distinct sentinels rather than to the message channel
// is what lets a caller use errors.Is on them: a stale context in particular is a bug in
// the caller's lifetime handling, not a contract failure, and worth telling apart.
func (buf ErrBuf) AsError(status Status, context string) error {
	switch status {
	case StatusOk:
		return nil
	case StatusNullPointer:
		return fmt.Errorf("%s: %w", context, ErrNullPointer)
	case StatusStaleContext:
		return fmt.Errorf("%s: %w", context, ErrStaleContext)
	case StatusPanic:
		return fmt.Errorf("%s: %w", context, ErrPanic)
	case StatusError:
		return fmt.Errorf("%s: %s", context, buf.Message())
	default:
		return fmt.Errorf("%s: unknown executor status %d", context, int32(status))
	}
}

// AsErrorFromCode is AsError for the calls that return a count, where a negative result
// carries a status in the same code space.
func (buf ErrBuf) AsErrorFromCode(code int64, context string) error {
	return buf.AsError(Status(code), context)
}

package wasmer2

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

var ErrFailedInstantiation = errors.New("could not create wasmer instance")

var ErrFailedCacheImports = errors.New("could not cache imports")

var ErrInvalidBytecode = errors.New("invalid bytecode")

// ErrVMHooksCannotHandleErrors signals that the VM hooks given to this backend cannot absorb
// an error returned by a hook. See executor.VMHooksErrorHandler for why this backend needs
// them to.
var ErrVMHooksCannotHandleErrors = errors.New("vm hooks do not implement VMHooksErrorHandler")

var ErrCachingFailed = errors.New("instance caching failed")

// GetLastError returns the last error message if any, otherwise returns an error.
func GetLastError() (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var errorLength = cWasmerLastErrorLength()

	if errorLength == 0 {
		return "", nil
	}

	var errorMessage = make([]cChar, errorLength)
	var errorMessagePointer = (*cChar)(unsafe.Pointer(&errorMessage[0]))

	var errorResult = cWasmerLastErrorMessage(errorMessagePointer, errorLength)

	if errorResult == -1 {
		return "", errors.New("cannot read last error")
	}

	return cGoString(errorMessagePointer), nil
}

func newWrappedError(target error) error {
	var lastError string
	var err error
	lastError, err = GetLastError()

	if err != nil {
		lastError = "unknown details"
	}

	return fmt.Errorf("%w: %s", target, lastError)
}

// ErrNilExecutorRuntime signals that an ExecutorNew was asked for without the host behind it.
// Unlike the VM hooks of the old interface, the runtime is not optional: it is the only
// channel through which an instance reaches the VM.
var ErrNilExecutorRuntime = errors.New("nil executor runtime")

// ErrEmptyCacheArtifact signals an attempt to instantiate from an empty compiled-code buffer.
var ErrEmptyCacheArtifact = errors.New("cache artifact is empty")

// ErrInstanceClosed signals an operation on an instance whose C memory has been released.
var ErrInstanceClosed = errors.New("instance is closed")

// ErrDetachedInstanceState signals a use of an InstanceState whose instance is gone, or that
// was never attached to one. Counterpart of the Rust backend's BadInstancePointer.
var ErrDetachedInstanceState = errors.New("instance state is detached")

// ErrSetPointsLimit signals that the gas limit of a call could not be set.
var ErrSetPointsLimit = errors.New("could not set the points limit")

// ErrSetPointsUsed signals that the gas counter of an instance could not be written.
var ErrSetPointsUsed = errors.New("could not set the points used")

// ErrSetBreakpointValue signals that an instance's breakpoint global could not be written.
var ErrSetBreakpointValue = errors.New("could not set the breakpoint value")

// ErrInstanceReset signals that an instance could not be returned to its post-creation state.
var ErrInstanceReset = errors.New("could not reset the instance")

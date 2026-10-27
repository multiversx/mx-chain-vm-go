package wasmer2new

// #include <stdlib.h>
// #include "./libvmexeccapi_w2.h"
import "C"

import (
	"runtime/cgo"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/wasmercommon"
)

var _ executor.InstanceNew = (*Wasmer2Instance)(nil)

// Wasmer2Instance is one instantiated contract, seen from outside a host call.
//
// Everything a hook needs is on Wasmer2InstanceState instead, reached through the hook
// context, because this C API makes memory and gas reachable only from inside a host call.
// That split is executor.InstanceNew against executor.InstanceState, and this backend is
// where it is not merely a convention.
type Wasmer2Instance struct {
	cInstance   *cInstanceT
	hooksHandle cgo.Handle

	// exportedFunctions and importedFunctions are read once at creation. The library hands
	// over each whole list in a pair of calls, and neither can change for a compiled module,
	// so the repeated questions the VM asks - once per name of a blocklist, for every
	// contract it validates - cost a map lookup rather than a cgo call.
	exportedFunctions []string
	importedFunctions []string

	// inCall guards against re-entering this same instance, which the Rust side cannot
	// survive: WasmerProdInstance::call takes &mut self, so a second call while the first is
	// in flight aliases a mutable borrow. Nesting into a different instance is fine and is
	// what the VM does for executeOnDestContext.
	//
	// A plain bool, not a mutex, on purpose. An instance is not safe for concurrent use at
	// all, and a lock would imply otherwise; what this catches is the single-goroutine
	// re-entry a hook could perform, which is the mistake that is actually available to make.
	inCall bool

	// id is stable for the instance's whole life, including after it is closed.
	id string
}

// Call calls an exported function by name, with the gas it is allowed to burn.
//
// A failure of the contract is not an error here: it comes back as the result kind, which is
// what the enum is for. A failure of the boundary itself becomes CallRuntimeError, there
// being no other channel - executor.InstanceCallResult has no "the executor broke" case, and
// inventing one would mean every caller handling a case that cannot be recovered from
// anyway.
func (instance *Wasmer2Instance) Call(funcName string, pointsLimit uint64) executor.InstanceCallResult {
	if instance.cInstance == nil {
		return executor.CallResultRuntimeError(ErrInstanceClosed)
	}
	if instance.inCall {
		return executor.CallResultRuntimeError(wasmercommon.ErrInstanceReentered)
	}
	instance.inCall = true
	defer func() { instance.inCall = false }()

	buf := newErrBuf()
	var cResult C.struct_vm_exec_call_result_t

	status := cInstanceCall(instance.cInstance, funcName, pointsLimit, &cResult, buf)
	if err := buf.AsError(status, "instance call failed"); err != nil {
		return executor.CallResultRuntimeError(err)
	}

	return instance.toCallResult(&cResult, buf)
}

// toCallResult converts the C out-struct, reading the message only for the kinds that have
// one. For the rest the error buffer holds whatever the previous call left behind.
func (instance *Wasmer2Instance) toCallResult(
	cResult *C.struct_vm_exec_call_result_t,
	buf errBuf,
) executor.InstanceCallResult {
	switch cResult.kind {
	case C.VM_EXEC_CALL_OK:
		return executor.CallResultOk()

	case C.VM_EXEC_CALL_FUNCTION_NOT_FOUND:
		return executor.CallResultFunctionNotFound()

	case C.VM_EXEC_CALL_EARLY_EXIT:
		// The error itself stayed on the Go side: the hook recorded its message in the VM's
		// output state before returning, and only the code crossed the boundary. What comes
		// back here is therefore the code, and whatever text the library chose to add.
		return executor.CallResultEarlyExit(
			executor.ExitExecutionFailed(newCallError(ErrEarlyExit, buf.Message())),
			uint64(cResult.early_exit_code))

	case C.VM_EXEC_CALL_BREAKPOINT:
		breakpoint, err := executor.ParseBreakpointValueNew(uint64(cResult.breakpoint))
		if err != nil {
			return executor.CallResultRuntimeError(err)
		}

		return executor.CallResultBreakpoint(breakpoint)

	default:
		return executor.CallResultRuntimeError(newCallError(ErrInstanceCall, buf.Message()))
	}
}

// CheckSignatures checks that all public module functions (SC endpoints) have no arguments
// or results.
func (instance *Wasmer2Instance) CheckSignatures() bool {
	if instance.cInstance == nil {
		return false
	}

	return cInstanceCheckSignatures(instance.cInstance) == 1
}

// HasFunction checks whether the SC has an endpoint with the given name.
func (instance *Wasmer2Instance) HasFunction(funcName string) bool {
	if instance.cInstance == nil {
		return false
	}

	return cInstanceHasFunction(instance.cInstance, funcName) == 1
}

// GetExportedFunctionNames returns the names of all the SC endpoints.
func (instance *Wasmer2Instance) GetExportedFunctionNames() []string {
	return instance.exportedFunctions
}

// GetImportedFunctionNames returns the names of all the host functions (VM hooks) the
// contract imports.
//
// The library answers this directly, which the old C API could not: there, the only question
// available was whether one particular name is imported.
func (instance *Wasmer2Instance) GetImportedFunctionNames() []string {
	return instance.importedFunctions
}

// GetPointsUsedNew returns the number of points (gas) used by the instance.
func (instance *Wasmer2Instance) GetPointsUsedNew() (uint64, error) {
	if instance.cInstance == nil {
		return 0, ErrInstanceClosed
	}

	buf := newErrBuf()
	var points C.uint64_t

	status := cInstanceGetPointsUsed(instance.cInstance, &points, buf)
	if err := buf.AsError(status, "could not read points used"); err != nil {
		return 0, err
	}

	return uint64(points), nil
}

// ResetNew returns the instance to the state it was in immediately after creation, so that
// it can serve another execution.
//
// Memory contents and length, the contract's own mutable globals, and the metering globals
// are all restored. The handle stays valid: the instance is rebuilt behind it rather than
// replaced, so the caller keeps using the same object, which is what executor.InstanceNew
// requires.
func (instance *Wasmer2Instance) ResetNew() error {
	if instance.cInstance == nil {
		return ErrInstanceClosed
	}

	buf := newErrBuf()

	status := cInstanceReset(instance.cInstance, buf)

	return buf.AsError(status, "could not reset instance")
}

// Cache serializes the compiled module, for ExecutorNew.NewInstanceFromCache.
//
// Two calls, but the first does the work rather than measuring it: the artifact is a few
// hundred KB, so the library serializes once and holds the result for the copy. That is why
// this is not wrapped in wasmercommon.ReadNames - the length call is not idempotent here, and
// the pair must stay adjacent.
func (instance *Wasmer2Instance) Cache() ([]byte, error) {
	if instance.cInstance == nil {
		return nil, ErrInstanceClosed
	}

	buf := newErrBuf()

	var length C.uint32_t
	status := cInstanceCacheLength(instance.cInstance, &length, buf)
	if err := buf.AsError(status, "could not serialize compiled module"); err != nil {
		return nil, err
	}
	if length == 0 {
		return nil, wasmercommon.ErrEmptyCacheArtifact
	}

	artifact := make([]byte, int(length))
	status = cInstanceCache(instance.cInstance, artifact, buf)
	if err := buf.AsError(status, "could not read serialized compiled module"); err != nil {
		return nil, err
	}

	return artifact, nil
}

// ID returns an identifier for the instance, unique at runtime.
func (instance *Wasmer2Instance) ID() string {
	return instance.id
}

// Close releases the instance and the hooks handle that reaches back into Go.
//
// In that order: once the instance is gone no hook can run, so the handle is safe to drop.
// Dropping it first would leave a window in which a hook could resolve a deleted handle.
//
// Not part of executor.InstanceNew, which has no lifecycle method; the handle would leak if
// this were never called.
func (instance *Wasmer2Instance) Close() {
	if instance.cInstance == nil {
		return
	}

	cInstanceDestroy(instance.cInstance)
	instance.cInstance = nil
	instance.hooksHandle.Delete()

	logWasmer2New.Trace("closed instance", "id", instance.id)
}

// IsInterfaceNil returns true if underlying object is nil
func (instance *Wasmer2Instance) IsInterfaceNil() bool {
	return instance == nil
}

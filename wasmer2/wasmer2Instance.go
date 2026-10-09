package wasmer2

// #include <stdlib.h>
// #include "./libvmexeccapi_w2.h"
import "C"

import (
	"runtime/cgo"
	"strings"
)

// CallResultKind says how a call ended.
//
// This is the enum that replaces the legacy arrangement, where a call returned a bare
// ok/error and the VM then read a wasm global to find out what had actually happened. The
// reason now comes back with the result.
type CallResultKind int

const (
	// CallOk means the function returned normally.
	CallOk CallResultKind = iota

	// CallFunctionNotFound means there is no exported function by that name.
	CallFunctionNotFound

	// CallRuntimeError means the instance trapped for a reason that is neither a VM hook's
	// doing nor a middleware breakpoint.
	CallRuntimeError

	// CallEarlyExit means a VM hook asked for execution to stop.
	CallEarlyExit

	// CallBreakpoint means a middleware stopped execution: out of gas, or the memory limit.
	CallBreakpoint
)

// CallResult is the outcome of a call.
type CallResult struct {
	Kind CallResultKind

	// EarlyExitCode is meaningful for CallEarlyExit. It is the vmcommon.ReturnCode the hook
	// exited with.
	EarlyExitCode uint64

	// Breakpoint is meaningful for CallBreakpoint, and matches executor.BreakpointValue.
	Breakpoint uint64

	// Message is informational, and carries the trap text for CallRuntimeError.
	//
	// For CallEarlyExit it is empty when the hook was a Go one, by construction: the hook
	// recorded its own message in the VM's output state before returning, and only the code
	// came back through the boundary. That is why the text does not have to travel.
	Message string
}

// Wasmer2Instance is one instantiated contract.
//
// The methods here are the outside view — what the VM can do to an instance between calls.
// Everything a *hook* needs is on InstanceState instead, reached through the hook context,
// because this C API makes memory and gas reachable only from inside a host call.
type Wasmer2Instance struct {
	cInstance   *cInstanceT
	hooksHandle cgo.Handle

	// importedFunctions is read once at creation, because the question the VM asks of it is
	// "is this name imported?" for each entry of a 21-name blocklist, on every contract it
	// validates. The legacy C API could only answer one name at a time, so that cost 21 cgo
	// calls; this one hands over the whole list in a pair of calls, so it costs two and the
	// rest is a map lookup.
	importedFunctions map[string]struct{}

	// inCall guards against re-entering this same instance, which the Rust side cannot
	// survive: WasmerProdInstance::call takes &mut self, so a second call while the first
	// is in flight aliases a mutable borrow. Nesting into a *different* instance is fine and
	// is what the VM does for executeOnDestContext.
	//
	// A plain bool, not a mutex, on purpose. An instance is not safe for concurrent use at
	// all, and a lock would imply otherwise; what this catches is the single-goroutine
	// re-entry a hook could perform, which is the mistake that is actually available to make.
	inCall bool

	// hooks is this instance's side of the hook bridge, kept so that the legacy interface can
	// reach the state of a running hook. See wasmer2InstanceLegacy.go.
	hooks *instanceHooks

	// gasLimit is the limit for the next call, which the legacy interface sets beforehand.
	gasLimit uint64

	// breakpoint stands in for the wasm global this executor no longer has.
	breakpoint uint64

	cleaned bool

	// id is stable for the instance's whole life, including after it is closed. See ID().
	id string
}

// Call invokes an exported function with a gas limit for this call.
//
// The returned error reports a failure of the boundary itself. A failure of the *contract*
// is not an error: it comes back as the CallResult kind, which is the point of the enum.
func (instance *Wasmer2Instance) Call(funcName string, pointsLimit uint64) (CallResult, error) {
	if instance.inCall {
		return CallResult{}, ErrInstanceReentered
	}
	instance.inCall = true
	defer func() { instance.inCall = false }()

	buf := newErrBuf()
	var cResult C.struct_vm_exec_call_result_t

	status := cInstanceCall(instance.cInstance, funcName, pointsLimit, &cResult, buf)
	if status != C.VM_EXEC_OK {
		return CallResult{}, buf.asError(status, "instance call failed")
	}

	result := CallResult{
		Kind:          CallResultKind(cResult.kind),
		EarlyExitCode: uint64(cResult.early_exit_code),
		Breakpoint:    uint64(cResult.breakpoint),
	}

	// The error buffer doubles as the message channel for the two outcomes that have one.
	// For the rest it holds whatever the previous call left behind, so it is not read.
	if result.Kind == CallRuntimeError || result.Kind == CallEarlyExit {
		result.Message = buf.message()
	}

	return result, nil
}

// PointsUsed returns the gas consumed so far.
//
// Named apart from the legacy GetPointsUsed, which has nowhere to report a failure; see
// wasmer2InstanceLegacy.go.
func (instance *Wasmer2Instance) PointsUsed() (uint64, error) {
	buf := newErrBuf()
	var points C.uint64_t

	status := cInstanceGetPointsUsed(instance.cInstance, &points, buf)
	if err := buf.asError(status, "could not read points used"); err != nil {
		return 0, err
	}

	return uint64(points), nil
}

// CheckSignatures reports whether every exported function takes no arguments and returns no
// results, which is what the contract ABI requires.
func (instance *Wasmer2Instance) CheckSignatures() bool {
	return cInstanceCheckSignatures(instance.cInstance) == 1
}

// HasFunction reports whether the contract exports a function by that name.
func (instance *Wasmer2Instance) HasFunction(funcName string) bool {
	return cInstanceHasFunction(instance.cInstance, funcName) == 1
}

// ExportedFunctionNames returns the names of the contract's endpoints.
func (instance *Wasmer2Instance) ExportedFunctionNames() ([]string, error) {
	return readNames(
		func() C.int { return cInstanceExportedFunctionNamesLength(instance.cInstance) },
		func(dest []byte) C.int { return cInstanceExportedFunctionNames(instance.cInstance, dest) },
	)
}

// ImportedFunctionNames returns the names of the VM hooks the contract imports.
//
// This is what the contract validator needs, and it is new: the legacy API could only answer
// "is this particular name imported?", one name at a time against a hardcoded list.
func (instance *Wasmer2Instance) ImportedFunctionNames() ([]string, error) {
	return readNames(
		func() C.int { return cInstanceImportedFunctionNamesLength(instance.cInstance) },
		func(dest []byte) C.int { return cInstanceImportedFunctionNames(instance.cInstance, dest) },
	)
}

// IsFunctionImported reports whether the contract imports a VM hook by that name.
//
// Answered from the set captured at creation, not from the executor, so it costs nothing to
// ask repeatedly. That matters because the contract validator asks it once per name of a
// blocklist for every contract it checks.
func (instance *Wasmer2Instance) IsFunctionImported(name string) bool {
	_, imported := instance.importedFunctions[name]

	return imported
}

// Close releases the instance and the hooks handle that reaches back into Go.
//
// In that order: once the instance is gone no hook can run, so the handle is safe to drop.
// Dropping it first would leave a window in which a hook could resolve a deleted handle.
func (instance *Wasmer2Instance) Close() {
	if instance.cInstance == nil {
		return
	}

	cInstanceDestroy(instance.cInstance)
	instance.cInstance = nil
	instance.hooksHandle.Delete()
}

// readNames drives the two-step string-out convention: ask for the length, allocate, fill.
//
// Two calls rather than one because the executor does not allocate for the caller and the
// caller cannot guess the size. A zero length means there are no names at all, which is not
// an error.
func readNames(length func() C.int, fill func([]byte) C.int) ([]string, error) {
	size := length()
	if size < 0 {
		return nil, ErrNullPointer
	}
	if size == 0 {
		return nil, nil
	}

	dest := make([]byte, int(size))
	written := fill(dest)
	if written <= 0 {
		return nil, ErrBufferTooSmall
	}

	// `written` counts the NUL terminator the executor wrote; the names stop before it.
	return strings.Split(string(dest[:written-1]), "|"), nil
}

package wasmer2new

// #include <stdlib.h>
// #include "./libvmexeccapi_w2.h"
import "C"

import (
	"runtime/cgo"

	"github.com/multiversx/mx-chain-vm-go/wasmercommon"
)

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
func (instance *Wasmer2Instance) Call(funcName string, pointsLimit uint64) (wasmercommon.CallResult, error) {
	if instance.inCall {
		return wasmercommon.CallResult{}, wasmercommon.ErrInstanceReentered
	}
	instance.inCall = true
	defer func() { instance.inCall = false }()

	buf := newErrBuf()
	var cResult C.struct_vm_exec_call_result_t

	status := cInstanceCall(instance.cInstance, funcName, pointsLimit, &cResult, buf)
	if status != C.VM_EXEC_OK {
		return wasmercommon.CallResult{}, buf.AsError(status, "instance call failed")
	}

	result := wasmercommon.CallResult{
		Kind:          wasmercommon.CallResultKind(cResult.kind),
		EarlyExitCode: uint64(cResult.early_exit_code),
		Breakpoint:    uint64(cResult.breakpoint),
	}

	// The error buffer doubles as the message channel for the two outcomes that have one.
	// For the rest it holds whatever the previous call left behind, so it is not read.
	if result.Kind == wasmercommon.CallRuntimeError || result.Kind == wasmercommon.CallEarlyExit {
		result.Message = buf.Message()
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
	if err := buf.AsError(status, "could not read points used"); err != nil {
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
	return wasmercommon.ReadNames(
		func() int32 { return cInstanceExportedFunctionNamesLength(instance.cInstance) },
		func(dest []byte) int32 { return cInstanceExportedFunctionNames(instance.cInstance, dest) },
	)
}

// ImportedFunctionNames returns the names of the VM hooks the contract imports.
//
// This is what the contract validator needs, and it is new: the legacy API could only answer
// "is this particular name imported?", one name at a time against a hardcoded list.
func (instance *Wasmer2Instance) ImportedFunctionNames() ([]string, error) {
	return wasmercommon.ReadNames(
		func() int32 { return cInstanceImportedFunctionNamesLength(instance.cInstance) },
		func(dest []byte) int32 { return cInstanceImportedFunctionNames(instance.cInstance, dest) },
	)
}

// ResetInstance returns the instance to the state it had immediately after creation, so it
// can serve another execution.
//
// Memory contents and length, the contract's own mutable globals, and the metering and
// breakpoint globals are all restored. The handle stays valid: the instance is rebuilt behind
// it rather than replaced, so the caller keeps using the same object.
func (instance *Wasmer2Instance) ResetInstance() error {
	buf := newErrBuf()

	status := cInstanceReset(instance.cInstance, buf)
	if err := buf.AsError(status, "could not reset instance"); err != nil {
		return err
	}

	// The wasm globals are fresh, so the Go-side stand-in for the breakpoint has to be too —
	// it is not a wasm global on this backend, so nothing in the library cleared it.
	instance.breakpoint = 0

	return nil
}

// CacheArtifact serializes the compiled module, for NewInstanceFromCache to load back.
//
// Two calls, as for the name readers, but the first does the work rather than measuring it:
// the artifact is a few hundred KB, so the library serializes once and holds the result for
// the copy. That is why these are not wrapped in wasmercommon.ReadNames — the length call is
// not idempotent here, and the pair must stay adjacent.
func (instance *Wasmer2Instance) CacheArtifact() ([]byte, error) {
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

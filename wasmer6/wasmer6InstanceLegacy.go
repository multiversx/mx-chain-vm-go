package wasmer6

// Transitional: executor.Instance as it is today, served by an instance whose own interface is
// already the narrower one.
//
// Everything in this file exists because the Go VM has not yet been converged onto the second
// interface set (PLAN-go-ffi.md, Phase 4b item 1). It is kept apart from wasmer6Instance.go so
// that the split stays visible, and so that this file can be deleted rather than untangled.
//
// Three groups:
//
//   - Methods that forward to the real thing, losing only an error the legacy signature has
//     nowhere to put.
//   - Methods that reach wasm memory or gas from *outside* a host call, which this executor
//     cannot do. They go through the hook context that is active while a hook runs, which is
//     exactly when the VM asks.
//   - Methods the new executor has no notion of at all: the breakpoint globals the call result
//     replaces, the per-instance hooks pointer the legacy C API needed, and the reset/cache
//     path deferred to Phase 6.

import (
	"fmt"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ = (executor.Instance)((*Wasmer6Instance)(nil))

// --- Forwarding, with the error dropped ---------------------------------------

// GetPointsUsed returns the gas consumed so far.
//
// The legacy signature has nowhere to report a failure, so a failure is logged and reported as
// zero. That is the pre-existing convention, not a new one: the legacy backend ignores the
// same status code without even looking at it.
func (instance *Wasmer6Instance) GetPointsUsed() uint64 {
	points, err := instance.PointsUsed()
	if err != nil {
		logWasmer6.Error("GetPointsUsed failed", "err", err)
		return 0
	}

	return points
}

// HasFunction reports whether the contract exports a function by that name.
func (instance *Wasmer6Instance) GetFunctionNames() []string {
	names, err := instance.ExportedFunctionNames()
	if err != nil {
		logWasmer6.Error("GetFunctionNames failed", "err", err)
		return nil
	}

	return names
}

// ValidateFunctionArities checks that every endpoint takes no arguments and returns none.
func (instance *Wasmer6Instance) ValidateFunctionArities() error {
	if !instance.CheckSignatures() {
		return executor.ErrFunctionNonvoidSignature
	}

	return nil
}

// --- Reached through the active hook context ---------------------------------

// MemLoad returns the contents of the instance memory at the given offset.
//
// Only answerable while a hook is running, because that is when this executor can address the
// instance's memory at all, and it is also the only time the VM asks.
func (instance *Wasmer6Instance) MemLoad(memPtr executor.MemPtr, length executor.MemLength) ([]byte, error) {
	state, err := instance.stateInActiveCall()
	if err != nil {
		return nil, err
	}

	return state.MemLoad(memPtr, length)
}

// MemStore writes to the instance memory at the given offset. See MemLoad.
func (instance *Wasmer6Instance) MemStore(memPtr executor.MemPtr, data []byte) error {
	state, err := instance.stateInActiveCall()
	if err != nil {
		return err
	}

	return state.MemStore(memPtr, data)
}

// MemLength returns the size of the instance memory in bytes. Used only by tests.
func (instance *Wasmer6Instance) MemLength() uint32 {
	state, err := instance.stateInActiveCall()
	if err != nil {
		logWasmer6.Error("MemLength failed", "err", err)
		return 0
	}

	length, err := state.MemLength()
	if err != nil {
		logWasmer6.Error("MemLength failed", "err", err)
		return 0
	}

	return uint32(length)
}

// MemGrow grows the instance memory by whole pages. Used only by tests.
func (instance *Wasmer6Instance) MemGrow(pages uint32) error {
	state, err := instance.stateInActiveCall()
	if err != nil {
		return err
	}

	_, err = state.MemGrow(pages)

	return err
}

// SetPointsUsed overwrites the gas consumed so far.
//
// Reachable only from inside a host call: the new C API offers gas writes through the hook
// context and not on the instance handle. The VM does write it from outside one place — when
// it resets a warm instance's counter — which this executor does not yet support anyway. See
// PLAN-go-ffi.md, Phase 4b.
func (instance *Wasmer6Instance) SetPointsUsed(points uint64) {
	state, err := instance.stateInActiveCall()
	if err != nil {
		logWasmer6.Error("SetPointsUsed failed", "err", err)
		return
	}

	if err := state.SetPointsUsed(points); err != nil {
		logWasmer6.Error("SetPointsUsed failed", "err", err)
	}
}

// stateInActiveCall returns the state of the hook currently running on this instance.
func (instance *Wasmer6Instance) stateInActiveCall() (*InstanceState, error) {
	if instance.hooks == nil || instance.hooks.activeState == nil {
		return nil, fmt.Errorf("%w: no host call is running on this instance", ErrStaleContext)
	}

	return instance.hooks.activeState, nil
}

// --- Gas limit, which moved from compilation to the call ---------------------

// SetGasLimit records the gas limit for the next call.
//
// In the new interface the limit is an argument of the call rather than a property of the
// compiled module, which is where it belongs: it changes on every invocation while the module
// does not. The legacy interface sets it beforehand, so it is held here until CallFunction.
func (instance *Wasmer6Instance) SetGasLimit(gasLimit uint64) {
	instance.gasLimit = gasLimit
}

// --- The call, and the breakpoint the call result replaces -------------------

// CallFunction invokes an exported function and reports why it stopped.
//
// Unlike the legacy backend, this one knows: the reason arrives with the result instead of
// being left in a wasm global for the VM to poll afterwards. Recording it where the VM's
// existing readers look is the VM's job, not this one's — see runtimeContext.CallSCFunction.
func (instance *Wasmer6Instance) CallFunction(functionName string) executor.InstanceCallResult {
	result, err := instance.Call(functionName, instance.gasLimit)
	if err != nil {
		return executor.CallResultError(err)
	}

	switch result.Kind {
	case CallOk:
		return executor.CallResultOk()
	case CallFunctionNotFound:
		return executor.CallResultError(executor.ErrFuncNotFound)
	case CallEarlyExit:
		return executor.InstanceCallResult{
			Kind:          executor.CallEarlyExit,
			Err:           fmt.Errorf("execution stopped by VM hook, code %d", result.EarlyExitCode),
			EarlyExitCode: result.EarlyExitCode,
		}
	case CallBreakpoint:
		// The values are the same on both sides on purpose: 4 is out of gas, 5 is the memory
		// limit, in executor.BreakpointValue and in the Rust BreakpointValue alike.
		return executor.InstanceCallResult{
			Kind:       executor.CallBreakpoint,
			Err:        fmt.Errorf("execution stopped by breakpoint %d", result.Breakpoint),
			Breakpoint: result.Breakpoint,
		}
	default:
		return executor.CallResultError(fmt.Errorf("instance call error: %s", result.Message))
	}
}

// SetBreakpointValue records a breakpoint the VM wants the instance to stop on.
//
// There is no wasm global to write here: this executor reports why a call stopped in the call
// result instead of leaving the reason in a global for the VM to poll. The value is kept so
// that GetBreakpointValue answers consistently, which is all the VM uses it for.
func (instance *Wasmer6Instance) SetBreakpointValue(value uint64) {
	instance.breakpoint = value
}

// GetBreakpointValue returns the last breakpoint set, whether by the VM or by a call result.
func (instance *Wasmer6Instance) GetBreakpointValue() uint64 {
	return instance.breakpoint
}

// --- Lifecycle --------------------------------------------------------------

// Clean releases the instance.
func (instance *Wasmer6Instance) Clean() bool {
	if instance.cleaned {
		return false
	}

	instance.Close()
	instance.cleaned = true

	return true
}

// IsAlreadyCleaned reports whether Clean has already run.
func (instance *Wasmer6Instance) IsAlreadyCleaned() bool {
	return instance.cleaned
}

// ID returns an identifier unique among the live instances.
//
// Captured at creation rather than derived from the handle on each call, because the instance
// tracker keys its bookkeeping by this and still asks for it after the instance is closed —
// at which point the handle is nil and every closed instance would answer the same thing.
func (instance *Wasmer6Instance) ID() string {
	return instance.id
}

// IsInterfaceNil returns true if there is no value under the interface.
func (instance *Wasmer6Instance) IsInterfaceNil() bool {
	return instance == nil
}

// --- Not supported by this executor -----------------------------------------

// Reset is not supported yet; see PLAN-go-ffi.md, Phase 6. Returning false makes the instance
// tracker treat the instance as unusable for warm reuse, which is the correct outcome.
func (instance *Wasmer6Instance) Reset() bool {
	return false
}

// Cache is not supported yet; see PLAN-go-ffi.md, Phase 6.
func (instance *Wasmer6Instance) Cache() ([]byte, error) {
	return nil, ErrCacheUnsupported
}

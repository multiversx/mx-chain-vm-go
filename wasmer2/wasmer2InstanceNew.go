package wasmer2

// #include <stdlib.h>
import "C"
import (
	"fmt"
	"runtime"
	"sync/atomic"
	"unsafe"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ executor.InstanceNew = (*Wasmer2InstanceNew)(nil)
var _ executor.VMHooksErrorHandler = (*Wasmer2InstanceNew)(nil)

// Wasmer2InstanceNew is a Wasmer 2.2 instance seen through the new interfaces.
//
// Counterpart of the Rust WasmerProdInstance. It is the outside view - calling the contract,
// gas accounting, module introspection - while what a hook reaches from inside the call is
// Wasmer2InstanceStateNew, held here and handed to the runtime when hooks are built.
//
// It is also the executor.VMHooksErrorHandler for its own calls: a legacy import has nowhere
// to report that a hook wants execution to stop, so the failure is applied to this instance
// instead.
type Wasmer2InstanceNew struct {
	cgoInstance *cWasmerInstanceT

	// The executor that made this instance, which owns the one slot the C side reads its
	// hooks from.
	wasmerExecutor *Wasmer2ExecutorNew

	// vmHooks is this instance's own hooks, built around its own state and kept for its
	// whole life, exactly as the Rust backend wires an adapter into each instance's import
	// object. They are installed in the executor's slot while this instance runs.
	vmHooks executor.VMHooksLegacy

	memory Wasmer2Memory
	state  Wasmer2InstanceStateNew

	id string

	// earlyExit is the error of the hook that asked for this call to stop, kept from the
	// moment the hook fails until Call picks it up once the instance has unwound.
	//
	// Written only from a hook, which runs on the goroutine that is inside the C call, so it
	// needs no synchronisation of its own.
	earlyExit error
}

var instanceNewIDCounter uint64

func newInstanceNew(cgoInstance *cWasmerInstanceT, wasmerExecutor *Wasmer2ExecutorNew) *Wasmer2InstanceNew {
	id := atomic.AddUint64(&instanceNewIDCounter, 1)

	instance := &Wasmer2InstanceNew{
		cgoInstance:    cgoInstance,
		wasmerExecutor: wasmerExecutor,
		memory:         Wasmer2Memory{cgoInstance: cgoInstance},
		id:             fmt.Sprintf("%p-%d", cgoInstance, id),
	}
	instance.state.instance = instance
	instance.vmHooks = newInstanceVMHooks(wasmerExecutor.runtime, instance)

	return instance
}

// newInstanceVMHooks builds the legacy-shaped hooks of one instance: the host behind the
// executor, this instance's state, and the instance itself as the handler for an error a
// legacy hook cannot report.
//
// The same three the Rust backend wires into an instance's import object, give or take the
// handler, which Rust keeps on its InstanceState.
//
// instance may be nil, which is the hooks of no instance: every state operation reports
// itself detached, and an error has nowhere to go but the log.
func newInstanceVMHooks(runtime executor.ExecutorRuntime, instance *Wasmer2InstanceNew) executor.VMHooksLegacy {
	return executor.NewVMHooksLegacyAdapterForRuntime(
		runtime,
		instance.instanceState(),
		instanceErrorHandler{runtime: runtime, instance: instance},
	)
}

var _ executor.VMHooksErrorHandler = instanceErrorHandler{}

// instanceErrorHandler applies a failed hook's side effects on both sides.
//
// Two concerns, and both are needed: the VM records in its output state why the contract
// stopped - the return code, the message - while the instance halts the contract. Neither
// can do the other's half. The new C API path does the same pairing in hookEarlyExit; here
// it has to be assembled, because the legacy adapter takes a single error handler.
type instanceErrorHandler struct {
	runtime  executor.ExecutorRuntime
	instance *Wasmer2InstanceNew
}

// HandleHookError implements executor.VMHooksErrorHandler.
func (handler instanceErrorHandler) HandleHookError(err error) {
	if err == nil {
		return
	}

	// The VM's half. Asked for per call, like the hooks themselves, because that is the only
	// way to reach the hooks object servicing this call.
	vmHooks := handler.runtime.VMHooksForCall(handler.instance.instanceState())
	if vmErrorHandler, ok := vmHooks.(executor.VMHooksErrorHandler); ok {
		vmErrorHandler.HandleHookError(err)
	} else {
		logWasmer2.Error("VM hooks cannot record a hook error", "err", err)
	}

	handler.instance.HandleHookError(err)
}

// Call calls an exported function by name, with the gas it is allowed to burn.
func (instance *Wasmer2InstanceNew) Call(funcName string, pointsLimit uint64) executor.InstanceCallResult {
	if instance.cgoInstance == nil {
		return executor.CallResultRuntimeError(ErrInstanceClosed)
	}

	if !instance.HasFunction(funcName) {
		return executor.CallResultFunctionNotFound()
	}

	if cWasmerInstanceSetGasLimit(instance.cgoInstance, pointsLimit) != cWasmerOk {
		return executor.CallResultRuntimeError(newWrappedError(ErrSetPointsLimit))
	}

	// Clear the breakpoint global, so the instance is armed for this call rather than still
	// halted from the last one.
	//
	// Under the old interface this was the VM's job, done in runtimeContext.UseWarmInstance.
	// The new interface does not expose the global at all - it is 2.2's private mechanism,
	// and Wasmer 7 has nothing like it - so the only side left that can do this is the
	// executor. Not doing it would make an instance single-use after any early exit.
	if cWasmerInstanceSetBreakpointValue(instance.cgoInstance, legacyBreakpointNone) != cWasmerOk {
		return executor.CallResultRuntimeError(newWrappedError(ErrSetBreakpointValue))
	}

	instance.earlyExit = nil

	err := instance.callFunction(funcName)
	if err == nil {
		return executor.CallResultOk()
	}

	// A hook that asked to stop is the reason the call failed, and it knows more than the
	// breakpoint does, so it is read first.
	if earlyExit := instance.takeEarlyExit(); earlyExit != nil {
		return executor.CallResultEarlyExit(earlyExit, executor.EarlyExitCodeOf(earlyExit))
	}

	// Reached only when no hook stashed an error above, which is what makes the
	// host-requested stop distinguishable from a hook's failure; see
	// legacyBreakpointHostRequestedStop.
	switch cWasmerInstanceGetBreakpointValue(instance.cgoInstance) {
	case legacyBreakpointHostRequestedStop:
		return executor.CallResultInterrupted()
	case legacyBreakpointOutOfGas:
		return executor.CallResultBreakpoint(executor.BreakpointOutOfGasNew)
	case legacyBreakpointMemoryLimit:
		return executor.CallResultBreakpoint(executor.BreakpointMemoryLimitNew)
	default:
		return executor.CallResultRuntimeError(err)
	}
}

// callFunction is the C call itself, with the instance announced to the executor for exactly
// as long as it is running.
func (instance *Wasmer2InstanceNew) callFunction(funcName string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cFuncName := cCString(funcName)
	defer cFree(unsafe.Pointer(cFuncName))

	// The C side reaches the hooks through one slot, held by the executor, so this
	// instance's own hooks go in it for as long as the instance is running, and whatever was
	// there goes back afterwards. A contract called from inside a hook does the same with
	// its own, which is what makes nesting work: the saving and restoring is the Go call
	// stack, not a stack anything keeps.
	previous := instance.wasmerExecutor.useVMHooks(instance.vmHooks)
	defer instance.wasmerExecutor.useVMHooks(previous)

	if cWasmerInstanceCall(instance.cgoInstance, cFuncName) != cWasmerOk {
		return newWrappedError(fmt.Errorf("failed to call the `%s` exported function", funcName))
	}

	return nil
}

// HandleHookError stops this instance because one of its hooks failed.
//
// Counterpart of the Rust InstanceState::handle_early_exit_legacy, and it is on the instance
// for the same reason: a legacy import returns a plain value, so the only way to report the
// failure is to halt the contract and leave the error where the caller will find it.
func (instance *Wasmer2InstanceNew) HandleHookError(err error) {
	if err == nil {
		return
	}

	if instance == nil {
		// Only reachable if a hook fires with no call in progress, which the push around the
		// C call rules out. Logged rather than ignored, because it would mean the halt below
		// never happened and the contract kept running.
		logWasmer2.Error("VM hook failed outside of a call", "err", err)
		return
	}

	// Trip the breakpoint global first: under 2.2 that is the mechanical halt, and it is what
	// stops the contract after the current instruction.
	if instance.cgoInstance != nil {
		breakpoint := legacyBreakpointForEarlyExit(executor.EarlyExitCodeOf(err))
		if cWasmerInstanceSetBreakpointValue(instance.cgoInstance, breakpoint) != cWasmerOk {
			logWasmer2.Error("could not halt the instance after a VM hook failed",
				"id", instance.id, "err", err)
		}
	}

	// Then stash the error, for Call to turn into the call result. The first one wins: it is
	// the one that caused the halt, and anything after it is a consequence.
	if instance.earlyExit == nil {
		instance.earlyExit = err
	}
}

// instanceState returns the handle through which a hook reaches this instance.
//
// Nil-safe, and deliberately so: with no call in progress there is no instance to speak of,
// and a detached state reports that on every operation rather than letting a hook reach into
// nothing.
func (instance *Wasmer2InstanceNew) instanceState() executor.InstanceState {
	if instance == nil {
		return &Wasmer2InstanceStateNew{}
	}

	return &instance.state
}

// RequestStop asks a running call to stop at the next opportunity.
//
// Writes the breakpoint global, which the metering middleware polls at every basic-block
// boundary. Called from another goroutine while Call is in flight, which is the only moment
// it does anything: the cgo instance pointer is set at creation and never reassigned, so
// there is nothing here for the two to race over beyond the global itself, and writing that
// while wasm runs is what a halt trigger is.
func (instance *Wasmer2InstanceNew) RequestStop() {
	if instance.cgoInstance == nil {
		return
	}

	// Nothing useful to do with a failure: the call may simply have finished already, and a
	// stop is best effort by contract.
	_ = cWasmerInstanceSetBreakpointValue(instance.cgoInstance, legacyBreakpointHostRequestedStop)
}

func (instance *Wasmer2InstanceNew) takeEarlyExit() error {
	earlyExit := instance.earlyExit
	instance.earlyExit = nil

	return earlyExit
}

// CheckSignatures checks that all public module functions (SC endpoints) have no arguments
// or results.
func (instance *Wasmer2InstanceNew) CheckSignatures() bool {
	if instance.cgoInstance == nil {
		return false
	}

	return cWasmerCheckSignatures(instance.cgoInstance) == cWasmerOk
}

// HasFunction checks whether the SC has an endpoint with the given name.
func (instance *Wasmer2InstanceNew) HasFunction(funcName string) bool {
	if instance.cgoInstance == nil {
		return false
	}

	cFuncName := cCString(funcName)
	defer cFree(unsafe.Pointer(cFuncName))

	return cWasmerInstanceHasFunction(instance.cgoInstance, cFuncName) == 1
}

// GetExportedFunctionNames returns the names of all the SC endpoints.
func (instance *Wasmer2InstanceNew) GetExportedFunctionNames() []string {
	if instance.cgoInstance == nil {
		return nil
	}

	names, err := exportedFunctionNames(instance.cgoInstance)
	if err != nil {
		logWasmer2.Error("could not read the exported function names", "id", instance.id, "err", err)
		return nil
	}

	return names
}

// IsFunctionImported returns whether the contract imports the host function with the given
// name.
//
// One name per call is all the 2.2 C API offers, and all the host needs: it asks about the
// handful of hooks an epoch has not activated yet, never about the whole import section.
func (instance *Wasmer2InstanceNew) IsFunctionImported(name string) bool {
	if instance.cgoInstance == nil {
		return false
	}

	cImportName := cCString(name)
	defer cFree(unsafe.Pointer(cImportName))

	return cWasmerInstanceHasImportedFunction(instance.cgoInstance, cImportName) == 1
}

// GetPointsUsedNew returns the number of points (gas) used by the instance.
func (instance *Wasmer2InstanceNew) GetPointsUsedNew() (uint64, error) {
	if instance.cgoInstance == nil {
		return 0, ErrInstanceClosed
	}

	return cWasmerInstanceGetPointsUsed(instance.cgoInstance), nil
}

// SetPointsUsedNew overwrites the points (gas) used by the instance.
func (instance *Wasmer2InstanceNew) SetPointsUsedNew(points uint64) error {
	if instance.cgoInstance == nil {
		return ErrInstanceClosed
	}

	if cWasmerInstanceSetPointsUsed(instance.cgoInstance, points) != cWasmerOk {
		return newWrappedError(ErrSetPointsUsed)
	}

	return nil
}

// ResetNew returns the instance to the state it was in immediately after creation, so that
// it can serve another execution.
//
// Rebuilt in place, as the interface requires: the VM holds this instance across the reset.
func (instance *Wasmer2InstanceNew) ResetNew() error {
	if instance.cgoInstance == nil {
		return ErrInstanceClosed
	}

	if cWasmerInstanceReset(instance.cgoInstance) != cWasmerOk {
		logWasmer2.Trace("reset: warm instance", "id", instance.id, "ok", false)
		return newWrappedError(ErrInstanceReset)
	}

	instance.earlyExit = nil
	logWasmer2.Trace("reset: warm instance", "id", instance.id, "ok", true)

	return nil
}

// Cache serializes the compiled module, for ExecutorNew.NewInstanceFromCache.
func (instance *Wasmer2InstanceNew) Cache() ([]byte, error) {
	if instance.cgoInstance == nil {
		return nil, ErrInstanceClosed
	}

	var cacheBytes *cUchar
	var cacheLen cUint32T

	if cWasmerInstanceCache(instance.cgoInstance, &cacheBytes, &cacheLen) != cWasmerOk {
		return nil, newWrappedError(ErrCachingFailed)
	}

	goBytes := C.GoBytes(unsafe.Pointer(cacheBytes), C.int(cacheLen))
	C.free(unsafe.Pointer(cacheBytes))

	return goBytes, nil
}

// ID returns an identifier for the instance, unique at runtime.
func (instance *Wasmer2InstanceNew) ID() string {
	return instance.id
}

// Close releases the instance's C memory, and every state handle that points at it reports
// itself detached from here on.
//
// It must not run while the instance is executing -
// the VM's own abandon-on-timeout handling, which the old Wasmer2Instance carries, has no
// counterpart here because the new interface cannot express it.
func (instance *Wasmer2InstanceNew) Close() {
	if instance.cgoInstance == nil {
		return
	}

	cWasmerInstanceDestroy(instance.cgoInstance)
	instance.cgoInstance = nil
	instance.memory.cgoInstance = nil

	logWasmer2.Trace("closed instance", "id", instance.id)
}

// InstanceState returns the handle a hook reaches this instance's memory and gas through.
//
// Test-only, and it is why instanceState has an exported twin: during a call the state
// arrives at the hooks from the executor, so nothing in production ever asks an instance for
// it. A test that has to read or write an instance's memory between calls has no other way
// in, the new instance interface carrying no memory operations at all.
func (instance *Wasmer2InstanceNew) InstanceState() executor.InstanceState {
	return instance.instanceState()
}

// MemLength is the size of the instance's linear memory, in bytes.
//
// Not on any interface: memory belongs to the per-call instance state, and this reaches it
// from outside a call. It exists for tests that have to observe an instance's memory between
// calls - the warm instance memory-reset test above all - which an InstanceState cannot do.
func (instance *Wasmer2InstanceNew) MemLength() uint32 {
	return instance.memory.Length()
}

// MemDump yields the instance's linear memory as the executor sees it, without copying.
//
// Test-only, for the same reason as MemLength, and additionally because the slice aliases
// the memory rather than copying it: a test that needs the memory's base address - to check
// that a reset did not move it - cannot use a copy.
func (instance *Wasmer2InstanceNew) MemDump() []byte {
	return instance.memory.Data()
}

// IsAlreadyCleaned reports whether the instance has already been released.
//
// The C pointer is the honest answer: it is what Close and Abandon both clear.
func (instance *Wasmer2InstanceNew) IsAlreadyCleaned() bool {
	return instance.cgoInstance == nil
}

// Abandon gives up the instance without releasing it, leaking its C memory on purpose.
//
// The VM reaches for this when a contract has run past its deadline and has not stopped
// within the grace period that follows. The C thread may still be executing inside this
// instance, so destroying it would free memory out from under a live call - a SIGSEGV
// instead of a timeout. Leaking a few megabytes once per runaway contract is the cheaper
// failure.
//
// The instance reports itself cleaned from here on, so nothing closes it later.
func (instance *Wasmer2InstanceNew) Abandon() {
	if instance.cgoInstance == nil {
		return
	}

	instance.cgoInstance = nil
	instance.memory.cgoInstance = nil

	logWasmer2.Trace("abandoned instance", "id", instance.id)
}

// IsInterfaceNil returns true if underlying object is nil
func (instance *Wasmer2InstanceNew) IsInterfaceNil() bool {
	return instance == nil
}

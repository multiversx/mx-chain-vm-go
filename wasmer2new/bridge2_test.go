package wasmer2new

// Phase 4a of PLAN-go-ffi.md: prove the boundary, with the VM kept out of it.
//
// What is under test here is not the VM's behaviour but the FFI's: that the 16-byte hook
// result struct really is returned in registers across a cgo //export, that a Go hook can
// reach back into wasm memory through its hook context, that a hook's early exit stops the
// instance, and that a context kept past its call is rejected rather than followed. Each of
// those can only fail at the boundary, and none of them needs the VM to be involved.

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/wasmercommon"
	"github.com/stretchr/testify/require"
)

// The values the test hooks return, chosen to exercise the full width of each result shape:
// gasLeftValue needs all 64 bits to catch a result truncated to 32, and numArgumentsValue
// has its high bit clear so that a sign-extension bug shows up as a wrong value rather than
// an accidentally correct one.
const (
	gasLeftValue      = int64(0x0102030405060708)
	numArgumentsValue = int32(0x12345678)
)

// The bytes the getSCAddress hook writes, and where.
const (
	scAddressOffset = 128
	scAddressLen    = 32
	scAddressFirst  = 0xA1
	scAddressLast   = 0xA2
)

// hooksContract exercises one hook of each result shape, in both directions.
//
// Every endpoint takes no arguments and returns no results, as the contract ABI requires,
// so the contract reports a disagreement by calling signalError — which the test then sees
// as an early exit. "The call came back OK" therefore means the contract's own assertions
// passed as well.
const hooksContract = `
(module
    (import "env" "getGasLeft" (func $getGasLeft (result i64)))
    (import "env" "getNumArguments" (func $getNumArguments (result i32)))
    (import "env" "getSCAddress" (func $getSCAddress (param i32)))
    (import "env" "signalError" (func $signalError (param i32 i32)))
    (import "env" "storageLoadLength" (func $storageLoadLength (param i32 i32) (result i32)))

    (memory (export "memory") 1)
    (data (i32.const 0) "boundary failure")

    (func $fail (call $signalError (i32.const 0) (i32.const 16)))

    (func (export "doNothing"))

    ;; An i64 hook result, checked against all 64 bits of the promised value.
    (func (export "checkI64Result")
        (if (i64.ne (call $getGasLeft) (i64.const 0x0102030405060708))
            (then (call $fail))))

    ;; An i32 hook result, which travels in the same 64-bit struct field.
    (func (export "checkI32Result")
        (if (i32.ne (call $getNumArguments) (i32.const 0x12345678))
            (then (call $fail))))

    ;; A void hook that writes to wasm memory from the Go side, checked at both ends of the
    ;; range it was asked to write.
    (func (export "checkHookMemStore")
        (call $getSCAddress (i32.const 128))
        (if (i32.ne (i32.load8_u (i32.const 128)) (i32.const 0xA1))
            (then (call $fail)))
        (if (i32.ne (i32.load8_u (i32.const 159)) (i32.const 0xA2))
            (then (call $fail))))

    ;; A hook that reads wasm memory the contract just wrote, and returns what it made of it.
    (func (export "checkHookMemLoad")
        (i32.store8 (i32.const 64) (i32.const 11))
        (i32.store8 (i32.const 65) (i32.const 22))
        (i32.store8 (i32.const 66) (i32.const 33))
        (if (i32.ne (call $storageLoadLength (i32.const 64) (i32.const 3)) (i32.const 66))
            (then (call $fail))))

    ;; An early exit, followed by a hook that must never run.
    (func (export "earlyExitStops")
        (call $signalError (i32.const 0) (i32.const 16))
        (drop (call $getNumArguments)))

    ;; A trap that is neither a hook's doing nor a middleware breakpoint.
    (func (export "unreachableTrap") (unreachable))

    ;; Spends gas until something stops it.
    (func (export "burnGas") (local $i i32)
        (loop $l
            (local.set $i (i32.add (local.get $i) (i32.const 1)))
            (br $l)))
)
`

// --- Test VM hooks -----------------------------------------------------------

// testVMHooks implements only the hooks hooksContract imports.
//
// executor.VMHooks is embedded rather than implemented: all 277 methods are satisfied by the
// nil interface, so a hook the fixture was not supposed to call fails loudly instead of
// quietly returning a zero.
type testVMHooks struct {
	executor.VMHooks

	state *InstanceState
	log   *hookLog
}

// hookLog is shared across the per-call hooks objects of one instance, which is what lets a
// test assert on what happened across a whole call.
type hookLog struct {
	calls         []string
	handledErrors []error
	loadedMessage string
	memLength     uint64
	pointsSeen    uint64
	stateEscaped  *InstanceState
}

func (log *hookLog) countOf(name string) int {
	count := 0
	for _, call := range log.calls {
		if call == name {
			count++
		}
	}

	return count
}

func (hooks *testVMHooks) record(name string) {
	hooks.log.calls = append(hooks.log.calls, name)
}

func (hooks *testVMHooks) HandleHookError(err error) {
	hooks.log.handledErrors = append(hooks.log.handledErrors, err)
}

func (hooks *testVMHooks) GetGasLeft() (int64, error) {
	hooks.record("GetGasLeft")

	// Reading the gas and the memory size proves the state operations work from inside a
	// hook, on the path the contract takes for an ordinary i64 result.
	points, err := hooks.state.GetPointsUsed()
	if err != nil {
		return 0, err
	}
	hooks.log.pointsSeen = points

	length, err := hooks.state.MemLength()
	if err != nil {
		return 0, err
	}
	hooks.log.memLength = length

	// The context is kept past the hook's return on purpose, so that a test can confirm it
	// is rejected rather than followed.
	hooks.log.stateEscaped = hooks.state

	return gasLeftValue, nil
}

func (hooks *testVMHooks) GetNumArguments() (int32, error) {
	hooks.record("GetNumArguments")

	return numArgumentsValue, nil
}

func (hooks *testVMHooks) GetSCAddress(resultOffset executor.MemPtr) error {
	hooks.record("GetSCAddress")

	address := make([]byte, scAddressLen)
	address[0] = scAddressFirst
	address[scAddressLen-1] = scAddressLast

	return hooks.state.MemStore(resultOffset, address)
}

func (hooks *testVMHooks) StorageLoadLength(
	keyOffset executor.MemPtr,
	keyLength executor.MemLength,
) (int32, error) {
	hooks.record("StorageLoadLength")

	key, err := hooks.state.MemLoad(keyOffset, keyLength)
	if err != nil {
		return 0, err
	}

	sum := int32(0)
	for _, b := range key {
		sum += int32(b)
	}

	return sum, nil
}

func (hooks *testVMHooks) SignalError(
	messageOffset executor.MemPtr,
	messageLength executor.MemLength,
) error {
	hooks.record("SignalError")

	// The real hook records the message in the VM's output state here, which is why the
	// text never has to cross the boundary. The test does the same thing, into its log.
	message, err := hooks.state.MemLoad(messageOffset, messageLength)
	if err != nil {
		return err
	}
	hooks.log.loadedMessage = string(message)

	return &executor.EarlyExit{
		Code:    vmcommon.UserError,
		Message: string(message),
	}
}

// testHooksProvider builds a fresh hooks object for each host call, bound to that call's
// state, and shares one log across them all.
type testHooksProvider struct {
	log *hookLog
}

func (provider *testHooksProvider) VMHooksForCall(state executor.InstanceState) executor.VMHooks {
	return &testVMHooks{state: callScopedState(state), log: provider.log}
}

func (provider *testHooksProvider) VMHooksForInstanceLifetime() executor.VMHooks {
	panic("this backend has no instance-lifetime state; see wasmer2new/vmHooksProvider.go")
}

// callScopedState recovers the concrete state the hook context produced.
//
// The tests reach further than executor.InstanceState does — gas and memory size as well as
// memory access — because the boundary offers more than the VM currently asks of it, and all
// of it has to be exercised. The assertion is safe by construction: this provider is only ever
// called by this package's own shims.
func callScopedState(state executor.InstanceState) *InstanceState {
	return state.(*InstanceState)
}

// --- Fixtures ----------------------------------------------------------------

// uniformOpcodeCost gives every opcode the same cost.
//
// Reflection rather than a literal because the struct has over 500 fields, and rather than
// the config package's gas schedule because the amounts do not matter here — only that they
// are non-zero, so that metering actually counts and the out-of-gas path is reachable.
func uniformOpcodeCost(cost uint32) *executor.WASMOpcodeCost {
	opcodeCost := &executor.WASMOpcodeCost{}

	fields := reflect.ValueOf(opcodeCost).Elem()
	for i := 0; i < fields.NumField(); i++ {
		fields.Field(i).SetUint(uint64(cost))
	}

	return opcodeCost
}

func newTestExecutor(t *testing.T) *Wasmer2Executor {
	wasmerExecutor, err := NewExecutor(wasmercommon.ExecutorConfig{
		OpcodeVersion:  executor.OpcodeVersionV2,
		WASMOpcodeCost: uniformOpcodeCost(1),
	})
	require.Nil(t, err)
	t.Cleanup(wasmerExecutor.Close)

	return wasmerExecutor
}

func newTestInstance(t *testing.T, code []byte) (*Wasmer2Instance, *hookLog) {
	log := &hookLog{}
	instance, err := newTestExecutor(t).NewInstance(
		code,
		wasmercommon.CompilationOptions{MaxMemoryGrow: 10, MaxMemoryGrowDelta: 10},
		&testHooksProvider{log: log},
	)
	require.Nil(t, err)
	t.Cleanup(instance.Close)

	return instance, log
}

// newHooksInstance instantiates hooksContract.
//
// The contract is handed over as WAT text, which the executor compiles just as it does a
// wasm binary. That keeps the fixture readable in place of a committed blob; a real contract
// binary goes through the same path in TestRealContractCompiles.
func newHooksInstance(t *testing.T) (*Wasmer2Instance, *hookLog) {
	return newTestInstance(t, []byte(hooksContract))
}

// callOk calls funcName and requires that it returned normally, which for this contract also
// means its own in-wasm assertions passed.
func callOk(t *testing.T, instance *Wasmer2Instance, funcName string) {
	result, err := instance.Call(funcName, 10_000_000)
	require.Nil(t, err)
	require.Equal(t, wasmercommon.CallOk, result.Kind, "call %s: %+v", funcName, result)
}

// --- Executor and instance lifecycle -----------------------------------------

func TestNewExecutor(t *testing.T) {
	wasmerExecutor := newTestExecutor(t)
	require.NotNil(t, wasmerExecutor.cExecutor)
}

func TestNewInstanceRejectsEmptyCode(t *testing.T) {
	_, err := newTestExecutor(t).NewInstance([]byte{}, wasmercommon.CompilationOptions{}, &testHooksProvider{})
	require.ErrorIs(t, err, wasmercommon.ErrInvalidBytecode)
}

func TestNewInstanceReportsCompilationError(t *testing.T) {
	_, err := newTestExecutor(t).NewInstance(
		[]byte("(module (func (export \"broken\") (i32.const 1)))"),
		wasmercommon.CompilationOptions{},
		&testHooksProvider{},
	)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "could not create instance")
}

func TestCloseIsIdempotent(t *testing.T) {
	instance, _ := newHooksInstance(t)
	instance.Close()
	instance.Close()
}

// --- Introspection -----------------------------------------------------------

func TestInstanceIntrospection(t *testing.T) {
	instance, _ := newHooksInstance(t)

	require.True(t, instance.CheckSignatures())
	require.True(t, instance.HasFunction("doNothing"))
	require.False(t, instance.HasFunction("noSuchFunction"))

	exported, err := instance.ExportedFunctionNames()
	require.Nil(t, err)
	require.Contains(t, exported, "doNothing")
	require.Contains(t, exported, "checkHookMemLoad")
	require.NotContains(t, exported, "memory", "a memory export is not a function")

	imported, err := instance.ImportedFunctionNames()
	require.Nil(t, err)
	require.ElementsMatch(t,
		[]string{"getGasLeft", "getNumArguments", "getSCAddress", "signalError", "storageLoadLength"},
		imported,
	)
}

// --- The hook result ABI -----------------------------------------------------

// TestHookResultShapes is the test the whole phase exists for.
//
// Until it runs, the 16-byte hook result struct is proven in two halves — Rust's view by its
// own boundary tests, Go's view by this package compiling against the declared prototypes —
// and never end to end. A struct returned by value from a cgo //export arrives in registers
// under both calling conventions the VM is built for; this is where that is observed rather
// than assumed, for each result width the hooks use.
func TestHookResultShapes(t *testing.T) {
	t.Run("i64 result", func(t *testing.T) {
		instance, log := newHooksInstance(t)
		callOk(t, instance, "checkI64Result")
		require.Equal(t, 1, log.countOf("GetGasLeft"))
	})

	t.Run("i32 result", func(t *testing.T) {
		instance, log := newHooksInstance(t)
		callOk(t, instance, "checkI32Result")
		require.Equal(t, 1, log.countOf("GetNumArguments"))
	})

	t.Run("void result", func(t *testing.T) {
		instance, log := newHooksInstance(t)
		callOk(t, instance, "checkHookMemStore")
		require.Equal(t, 1, log.countOf("GetSCAddress"))
	})
}

func TestCallWithNoHooks(t *testing.T) {
	instance, log := newHooksInstance(t)
	callOk(t, instance, "doNothing")
	require.Empty(t, log.calls)
}

func TestFunctionNotFound(t *testing.T) {
	instance, _ := newHooksInstance(t)

	result, err := instance.Call("noSuchFunction", 10_000_000)
	require.Nil(t, err)
	require.Equal(t, wasmercommon.CallFunctionNotFound, result.Kind)
}

// --- State operations from inside a hook -------------------------------------

func TestHookWritesWasmMemory(t *testing.T) {
	instance, log := newHooksInstance(t)

	// The contract checks the bytes itself and signals an error if they are wrong, so a
	// plain OK is the assertion.
	callOk(t, instance, "checkHookMemStore")
	require.Equal(t, []string{"GetSCAddress"}, log.calls)
}

func TestHookReadsWasmMemory(t *testing.T) {
	instance, log := newHooksInstance(t)

	callOk(t, instance, "checkHookMemLoad")
	require.Equal(t, []string{"StorageLoadLength"}, log.calls)
}

func TestHookReadsGasAndMemoryLength(t *testing.T) {
	instance, log := newHooksInstance(t)

	callOk(t, instance, "checkI64Result")
	require.Greater(t, log.pointsSeen, uint64(0), "metering is always on, so gas was spent")
	require.Equal(t, uint64(65536), log.memLength, "one declared page")
}

// --- Early exit --------------------------------------------------------------

func TestEarlyExitStopsExecution(t *testing.T) {
	instance, log := newHooksInstance(t)

	result, err := instance.Call("earlyExitStops", 10_000_000)
	require.Nil(t, err)

	require.Equal(t, wasmercommon.CallEarlyExit, result.Kind)
	require.Equal(t, uint64(vmcommon.UserError), result.EarlyExitCode)

	// The message was recorded inside the hook, which is the whole reason it does not have
	// to travel back through the result struct.
	require.Equal(t, "boundary failure", log.loadedMessage)

	// The hook after signalError must not have run.
	require.Equal(t, []string{"SignalError"}, log.calls)

	// And the VM side was told, so it could record the return code and message.
	require.Len(t, log.handledErrors, 1)
	var earlyExit *executor.EarlyExit
	require.True(t, errors.As(log.handledErrors[0], &earlyExit))
	require.Equal(t, vmcommon.UserError, earlyExit.Code)
}

// --- Failure modes that are not hooks ----------------------------------------

func TestRuntimeError(t *testing.T) {
	instance, _ := newHooksInstance(t)

	result, err := instance.Call("unreachableTrap", 10_000_000)
	require.Nil(t, err)
	require.Equal(t, wasmercommon.CallRuntimeError, result.Kind)
	require.NotEmpty(t, result.Message)
}

func TestOutOfGasBreakpoint(t *testing.T) {
	instance, _ := newHooksInstance(t)

	result, err := instance.Call("burnGas", 1000)
	require.Nil(t, err)
	require.Equal(t, wasmercommon.CallBreakpoint, result.Kind)
	require.Equal(t, uint64(4), result.Breakpoint, "4 is out of gas")
}

func TestGasAccounting(t *testing.T) {
	instance, _ := newHooksInstance(t)

	callOk(t, instance, "checkI64Result")

	used, err := instance.PointsUsed()
	require.Nil(t, err)
	require.Greater(t, used, uint64(0))
}

// --- Lifetime enforcement ----------------------------------------------------

// TestStaleContextIsRejected is the replacement for the borrow the Rust side cannot express
// across C.
//
// The hook keeps its InstanceState past its own return. Every operation on it must then be
// refused: the context's call token no longer matches the instance's innermost active call,
// which is a check on the Rust side that exists precisely so that this reads as an error
// instead of dereferencing a borrow that has gone away.
func TestStaleContextIsRejected(t *testing.T) {
	instance, log := newHooksInstance(t)

	callOk(t, instance, "checkI64Result")
	require.NotNil(t, log.stateEscaped, "the hook was supposed to let its state escape")

	_, err := log.stateEscaped.MemLoad(0, 4)
	require.ErrorIs(t, err, wasmercommon.ErrStaleContext)

	err = log.stateEscaped.MemStore(0, []byte{1, 2, 3})
	require.ErrorIs(t, err, wasmercommon.ErrStaleContext)

	_, err = log.stateEscaped.MemLength()
	require.ErrorIs(t, err, wasmercommon.ErrStaleContext)

	_, err = log.stateEscaped.GetPointsUsed()
	require.ErrorIs(t, err, wasmercommon.ErrStaleContext)

	err = log.stateEscaped.SetPointsUsed(1)
	require.ErrorIs(t, err, wasmercommon.ErrStaleContext)

	_, err = log.stateEscaped.MemGrow(1)
	require.ErrorIs(t, err, wasmercommon.ErrStaleContext)
}

// --- A real contract binary --------------------------------------------------

// TestRealContractCompiles puts a compiled contract through the same path, so that the suite
// does not rest entirely on WAT text.
//
// It stops at introspection: calling counter's endpoints would need the storage and argument
// hooks implemented, which is the VM's job and is what Phase 4b and the differential run are
// for.
func TestRealContractCompiles(t *testing.T) {
	code, err := os.ReadFile("../test/contracts/counter/output/counter.wasm")
	require.Nil(t, err, "fixture missing")

	instance, _ := newTestInstance(t, code)

	require.True(t, instance.CheckSignatures())
	require.True(t, instance.HasFunction("increment"))

	imported, err := instance.ImportedFunctionNames()
	require.Nil(t, err)
	require.NotEmpty(t, imported)
	fmt.Printf("counter.wasm imports %d hooks: %v\n", len(imported), imported)
}

func TestNewExecutorRejectsMissingOpcodeCost(t *testing.T) {
	_, err := NewExecutor(wasmercommon.ExecutorConfig{OpcodeVersion: executor.OpcodeVersionV2})
	require.ErrorIs(t, err, wasmercommon.ErrMissingOpcodeCost)
}

// --- Nesting -----------------------------------------------------------------

// nestingVMHooks re-enters the executor from inside a hook, which is what the VM does when a
// contract calls executeOnDestContext.
//
// This is the test for the decision Phase 4a made: hooks are built per call rather than per
// instance, so each nesting level gets its own state and no level can clobber another's. The
// outer hook therefore has to keep working *after* the inner call has finished, which is the
// part a mutable per-instance state slot would get wrong.
type nestingVMHooks struct {
	executor.VMHooks

	state *InstanceState
	log   *hookLog
	inner *Wasmer2Instance
}

func (hooks *nestingVMHooks) HandleHookError(err error) {
	hooks.log.handledErrors = append(hooks.log.handledErrors, err)
}

func (hooks *nestingVMHooks) GetSCAddress(resultOffset executor.MemPtr) error {
	hooks.log.calls = append(hooks.log.calls, "outer:GetSCAddress")

	// Reaching the state before re-entering establishes that it worked to begin with.
	lengthBefore, err := hooks.state.MemLength()
	if err != nil {
		return err
	}

	// Re-enter: a whole other instance runs, and calls hooks of its own, while this hook is
	// still on the stack.
	result, callErr := hooks.inner.Call("checkI64Result", 10_000_000)
	if callErr != nil {
		return callErr
	}
	if result.Kind != wasmercommon.CallOk {
		return fmt.Errorf("inner call did not succeed: %+v", result)
	}

	// And now the part that matters: this hook's own state must still be usable.
	lengthAfter, err := hooks.state.MemLength()
	if err != nil {
		return fmt.Errorf("outer state unusable after the inner call: %w", err)
	}
	if lengthAfter != lengthBefore {
		return fmt.Errorf("outer state changed under the inner call: %d then %d",
			lengthBefore, lengthAfter)
	}

	address := make([]byte, scAddressLen)
	address[0] = scAddressFirst
	address[scAddressLen-1] = scAddressLast

	return hooks.state.MemStore(resultOffset, address)
}

type nestingHooksProvider struct {
	log   *hookLog
	inner *Wasmer2Instance
}

func (provider *nestingHooksProvider) VMHooksForCall(state executor.InstanceState) executor.VMHooks {
	return &nestingVMHooks{state: callScopedState(state), log: provider.log, inner: provider.inner}
}

func (provider *nestingHooksProvider) VMHooksForInstanceLifetime() executor.VMHooks {
	panic("this backend has no instance-lifetime state; see wasmer2new/vmHooksProvider.go")
}

func TestHookReentersTheExecutor(t *testing.T) {
	// The inner instance has ordinary hooks, and its own log, so that the two levels can be
	// told apart.
	innerLog := &hookLog{}
	inner, err := newTestExecutor(t).NewInstance(
		[]byte(hooksContract),
		wasmercommon.CompilationOptions{MaxMemoryGrow: 10, MaxMemoryGrowDelta: 10},
		&testHooksProvider{log: innerLog},
	)
	require.Nil(t, err)
	t.Cleanup(inner.Close)

	outerLog := &hookLog{}
	outer, err := newTestExecutor(t).NewInstance(
		[]byte(hooksContract),
		wasmercommon.CompilationOptions{MaxMemoryGrow: 10, MaxMemoryGrowDelta: 10},
		&nestingHooksProvider{log: outerLog, inner: inner},
	)
	require.Nil(t, err)
	t.Cleanup(outer.Close)

	// checkHookMemStore calls getSCAddress and then verifies the bytes itself, so an OK
	// result means the outer hook completed its write after the nested call returned.
	callOk(t, outer, "checkHookMemStore")

	require.Equal(t, []string{"outer:GetSCAddress"}, outerLog.calls)
	require.Equal(t, []string{"GetGasLeft"}, innerLog.calls, "the inner instance ran its own hook")
	require.Empty(t, outerLog.handledErrors)
}

// reentrantVMHooks calls back into the instance it is already running in, which is the one
// nesting shape the Rust side cannot survive.
type reentrantVMHooks struct {
	executor.VMHooks

	log  *hookLog
	self func() *Wasmer2Instance
}

func (hooks *reentrantVMHooks) HandleHookError(err error) {
	hooks.log.handledErrors = append(hooks.log.handledErrors, err)
}

func (hooks *reentrantVMHooks) GetSCAddress(_ executor.MemPtr) error {
	_, err := hooks.self().Call("doNothing", 10_000_000)

	return err
}

type reentrantHooksProvider struct {
	log  *hookLog
	self func() *Wasmer2Instance
}

func (provider *reentrantHooksProvider) VMHooksForCall(_ executor.InstanceState) executor.VMHooks {
	return &reentrantVMHooks{log: provider.log, self: provider.self}
}

func (provider *reentrantHooksProvider) VMHooksForInstanceLifetime() executor.VMHooks {
	panic("this backend has no instance-lifetime state; see wasmer2new/vmHooksProvider.go")
}

// TestSameInstanceReentryIsRejected pins down the invariant the guard exists for.
//
// ExperimentalInstance::call takes &mut self, so a call placed while another is in flight on
// the same instance would alias a mutable borrow across the C boundary — undefined behaviour
// with nothing to catch it. The guard turns it into an error the hook sees, which then travels
// out as an ordinary early exit.
func TestSameInstanceReentryIsRejected(t *testing.T) {
	log := &hookLog{}

	var instance *Wasmer2Instance
	instance, err := newTestExecutor(t).NewInstance(
		[]byte(hooksContract),
		wasmercommon.CompilationOptions{MaxMemoryGrow: 10, MaxMemoryGrowDelta: 10},
		&reentrantHooksProvider{log: log, self: func() *Wasmer2Instance { return instance }},
	)
	require.Nil(t, err)
	t.Cleanup(instance.Close)

	result, err := instance.Call("checkHookMemStore", 10_000_000)
	require.Nil(t, err)
	require.Equal(t, wasmercommon.CallEarlyExit, result.Kind)

	require.Len(t, log.handledErrors, 1)
	require.ErrorIs(t, log.handledErrors[0], wasmercommon.ErrInstanceReentered)

	// Usable again afterwards, so the guard releases on the way out.
	callOk(t, instance, "doNothing")
}

// TestIsFunctionImported covers the question the contract validator actually asks.
//
// It checks a 21-name blocklist against every contract it validates. The legacy C API could
// only answer one name per call; here the names are captured once at creation, so the
// validator's loop never reaches the executor at all.
func TestIsFunctionImported(t *testing.T) {
	instance, _ := newHooksInstance(t)

	require.True(t, instance.IsFunctionImported("getGasLeft"))
	require.True(t, instance.IsFunctionImported("storageLoadLength"))

	// Names the contract exports rather than imports must not count as imported.
	require.False(t, instance.IsFunctionImported("doNothing"))

	// A hook that exists but this contract does not use.
	require.False(t, instance.IsFunctionImported("transferValue"))

	require.False(t, instance.IsFunctionImported(""))
	require.False(t, instance.IsFunctionImported("doesNotExist"))
}

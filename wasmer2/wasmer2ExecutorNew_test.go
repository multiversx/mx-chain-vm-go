package wasmer2

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/stretchr/testify/require"
)

// stubVMHooks implements the three hooks the counter contract imports, and nothing else.
//
// The interface is embedded rather than implemented: a contract that reaches a hook outside
// this set would fail loudly, which is what a test wants.
type stubVMHooks struct {
	executor.VMHooks

	// state is replaced on every hook, by the runtime below. That it is never stale is the
	// property these tests are built around.
	state executor.InstanceState

	storage  map[string]int64
	finished []int64
	calls    []string

	// failOn names a hook that reports failWith instead of doing its work.
	failOn   string
	failWith error

	// nested, if set, runs from inside Int64storageStore, which is how a contract called
	// from a hook is simulated.
	nested func()
}

func newStubVMHooks() *stubVMHooks {
	return &stubVMHooks{storage: make(map[string]int64)}
}

func (hooks *stubVMHooks) enter(name string) error {
	hooks.calls = append(hooks.calls, name)
	if hooks.failOn == name {
		return hooks.failWith
	}

	return nil
}

func (hooks *stubVMHooks) key(keyOffset executor.MemPtr, keyLength executor.MemLength) (string, error) {
	key, err := hooks.state.MemoryLoad(keyOffset, keyLength)

	return string(key), err
}

func (hooks *stubVMHooks) Int64storageStore(keyOffset executor.MemPtr, keyLength executor.MemLength, value int64) (int32, error) {
	if err := hooks.enter("int64storageStore"); err != nil {
		return 0, err
	}

	if hooks.nested != nil {
		nested := hooks.nested
		hooks.nested = nil
		nested()
	}

	key, err := hooks.key(keyOffset, keyLength)
	if err != nil {
		return 0, err
	}
	hooks.storage[key] = value

	return 0, nil
}

func (hooks *stubVMHooks) Int64storageLoad(keyOffset executor.MemPtr, keyLength executor.MemLength) (int64, error) {
	if err := hooks.enter("int64storageLoad"); err != nil {
		return 0, err
	}

	key, err := hooks.key(keyOffset, keyLength)
	if err != nil {
		return 0, err
	}

	return hooks.storage[key], nil
}

func (hooks *stubVMHooks) Int64finish(value int64) error {
	if err := hooks.enter("int64finish"); err != nil {
		return err
	}
	hooks.finished = append(hooks.finished, value)

	return nil
}

// stubRuntime is the host behind the executor, and the only thing it does is build hooks
// around the state it is handed.
type stubRuntime struct {
	hooks      *stubVMHooks
	statesSeen []executor.InstanceState
}

func (runtime *stubRuntime) VMHooksForCall(instanceState executor.InstanceState) executor.VMHooks {
	runtime.hooks.state = instanceState
	runtime.statesSeen = append(runtime.statesSeen, instanceState)

	return runtime.hooks
}

func (runtime *stubRuntime) IsInterfaceNil() bool {
	return runtime == nil
}

func contractCode(t *testing.T, name string) []byte {
	code, err := os.ReadFile(filepath.Join("..", "test", "contracts", name, "output", name+".wasm"))
	require.Nil(t, err)

	return code
}

func testCompilationOptions(t *testing.T) *executor.CompilationOptionsNew {
	gasCost, err := config.CreateGasConfig(config.MakeGasMapForTests())
	require.Nil(t, err)

	return &executor.CompilationOptionsNew{
		UnmeteredLocals:    uint64(gasCost.WASMOpcodeCost.LocalsUnmetered),
		MaxMemoryGrow:      uint64(gasCost.WASMOpcodeCost.MaxMemoryGrow),
		MaxMemoryGrowDelta: uint64(gasCost.WASMOpcodeCost.MaxMemoryGrowDelta),
		OpcodeConfig: &executor.OpcodeConfig{
			OpcodeVersion: executor.OpcodeVersionV2,
			OpcodeCosts:   gasCost.WASMOpcodeCost,
		},
	}
}

func newTestExecutor(t *testing.T) (*Wasmer2ExecutorNew, *stubVMHooks, *stubRuntime) {
	hooks := newStubVMHooks()
	runtime := &stubRuntime{hooks: hooks}

	wasmerExecutor, err := CreateExecutorNew(runtime)
	require.Nil(t, err)

	return wasmerExecutor, hooks, runtime
}

func newCounterInstance(t *testing.T, wasmerExecutor *Wasmer2ExecutorNew) *Wasmer2InstanceNew {
	instance, err := wasmerExecutor.NewInstance(contractCode(t, "counter"), testCompilationOptions(t))
	require.Nil(t, err)

	return instance.(*Wasmer2InstanceNew)
}

const testGasLimit = uint64(1_000_000_000)

func TestWasmer2ExecutorNew_RequiresARuntime(t *testing.T) {
	wasmerExecutor, err := CreateExecutorNew(nil)

	require.Nil(t, wasmerExecutor)
	require.ErrorIs(t, err, ErrNilExecutorRuntime)
}

func TestWasmer2ExecutorNew_CallReachesTheHooks(t *testing.T) {
	wasmerExecutor, hooks, _ := newTestExecutor(t)
	instance := newCounterInstance(t, wasmerExecutor)
	defer instance.Close()

	require.Equal(t, executor.CallOk, instance.Call("init", testGasLimit).Kind)
	require.Equal(t, executor.CallOk, instance.Call("increment", testGasLimit).Kind)

	// The key the contract passed came out of its own memory, through the instance state the
	// hooks were built around.
	require.Equal(t, map[string]int64{"COUNTER": 2}, hooks.storage)
	require.Equal(t, []int64{2}, hooks.finished)
	require.Equal(t,
		[]string{"int64storageStore", "int64storageLoad", "int64storageStore", "int64finish"},
		hooks.calls)
}

func TestWasmer2ExecutorNew_HooksAreBuiltOncePerHook(t *testing.T) {
	wasmerExecutor, hooks, runtime := newTestExecutor(t)
	instance := newCounterInstance(t, wasmerExecutor)
	defer instance.Close()

	require.Equal(t, executor.CallOk, instance.Call("increment", testGasLimit).Kind)

	// Three hooks ran, so the runtime was asked three times. Resolving per hook rather than
	// per call is what the Wasmer 7 contract requires, and what keeps a reused instance from
	// serving a stale execution.
	require.Len(t, hooks.calls, 3)
	require.Len(t, runtime.statesSeen, 3)
}

func TestWasmer2ExecutorNew_FunctionNotFound(t *testing.T) {
	wasmerExecutor, _, _ := newTestExecutor(t)
	instance := newCounterInstance(t, wasmerExecutor)
	defer instance.Close()

	require.Equal(t, executor.CallFunctionNotFound, instance.Call("nonexistent", testGasLimit).Kind)
}

func TestWasmer2ExecutorNew_FailingHookStopsTheContract(t *testing.T) {
	wasmerExecutor, hooks, _ := newTestExecutor(t)
	instance := newCounterInstance(t, wasmerExecutor)
	defer instance.Close()

	hookErr := errors.New("storage is unavailable")
	hooks.failOn = "int64storageStore"
	hooks.failWith = hookErr

	result := instance.Call("init", testGasLimit)

	require.Equal(t, executor.CallEarlyExit, result.Kind)
	require.ErrorIs(t, result.Err, hookErr)

	// A plain error means "execution failed", which is the default the VM applies too.
	require.Equal(t, uint64(vmcommon.ExecutionFailed), result.EarlyExitCode)
	require.Empty(t, hooks.storage)
}

func TestWasmer2ExecutorNew_EarlyExitKeepsItsOwnCode(t *testing.T) {
	wasmerExecutor, hooks, _ := newTestExecutor(t)
	instance := newCounterInstance(t, wasmerExecutor)
	defer instance.Close()

	hooks.failOn = "int64storageStore"
	hooks.failWith = &executor.EarlyExit{
		Code:    vmcommon.UserError,
		Message: "the contract said no",
	}

	result := instance.Call("init", testGasLimit)

	require.Equal(t, executor.CallEarlyExit, result.Kind)
	require.Equal(t, uint64(vmcommon.UserError), result.EarlyExitCode)
}

func TestWasmer2ExecutorNew_OutOfGasIsABreakpoint(t *testing.T) {
	wasmerExecutor, _, _ := newTestExecutor(t)
	instance := newCounterInstance(t, wasmerExecutor)
	defer instance.Close()

	result := instance.Call("increment", 1)

	require.Equal(t, executor.CallBreakpoint, result.Kind)
	require.Equal(t, executor.BreakpointOutOfGasNew, result.Breakpoint)
}

func TestWasmer2ExecutorNew_NestedCallResolvesToTheInnerInstance(t *testing.T) {
	wasmerExecutor, hooks, runtime := newTestExecutor(t)

	outer := newCounterInstance(t, wasmerExecutor)
	defer outer.Close()
	inner := newCounterInstance(t, wasmerExecutor)
	defer inner.Close()

	var innerResult executor.InstanceCallResult
	var innerStates []executor.InstanceState
	var outerState executor.InstanceState

	// Runs from inside the outer contract's first hook, which is where a contract calling
	// another contract happens.
	hooks.nested = func() {
		outerState = hooks.state

		before := len(runtime.statesSeen)
		innerResult = inner.Call("increment", testGasLimit)
		innerStates = runtime.statesSeen[before:]

		// A real VM holds the state it was handed for the duration of its hook; the stub
		// keeps one field, which the nested call has just overwritten. Put it back, so the
		// rest of the outer hook reads the outer contract's memory.
		hooks.state = outerState
	}

	require.Equal(t, executor.CallOk, outer.Call("init", testGasLimit).Kind)
	require.Equal(t, executor.CallOk, innerResult.Kind)

	// The outer hook was served by the outer instance, every hook of the nested call by the
	// inner one.
	require.Same(t, outer, outerState.(*Wasmer2InstanceStateNew).instance)
	require.NotEmpty(t, innerStates)
	for _, state := range innerStates {
		require.Same(t, inner, state.(*Wasmer2InstanceStateNew).instance)
	}

	// The outer handle is still usable once the nested call has unwound: a nested execution
	// suspends the outer one, it does not end it.
	memoryLength, err := outerState.MemoryLength()
	require.Nil(t, err)
	require.Greater(t, memoryLength, uint64(0))

	// And the slot the C side reads is back where it was: a further call on the outer
	// instance is still served by the outer instance.
	before := len(runtime.statesSeen)
	require.Equal(t, executor.CallOk, outer.Call("increment", testGasLimit).Kind)
	require.Greater(t, len(runtime.statesSeen), before)
	for _, state := range runtime.statesSeen[before:] {
		require.Same(t, outer, state.(*Wasmer2InstanceStateNew).instance)
	}
}

func TestWasmer2InstanceNew_ExportedAndImportedNames(t *testing.T) {
	wasmerExecutor, _, _ := newTestExecutor(t)
	instance := newCounterInstance(t, wasmerExecutor)
	defer instance.Close()

	require.ElementsMatch(t,
		[]string{"init", "upgrade", "increment", "decrement", "get"},
		instance.GetExportedFunctionNames())

	// Assembled by asking about every hook this backend offers, which is complete: the module
	// only instantiated because each of its imports resolved against that same table.
	require.Equal(t,
		[]string{"int64finish", "int64storageLoad", "int64storageStore"},
		instance.GetImportedFunctionNames())

	require.True(t, instance.HasFunction("increment"))
	require.False(t, instance.HasFunction("nonexistent"))
	require.True(t, instance.CheckSignatures())
}

func TestWasmer2InstanceNew_ResetAllowsReuse(t *testing.T) {
	wasmerExecutor, hooks, _ := newTestExecutor(t)
	instance := newCounterInstance(t, wasmerExecutor)
	defer instance.Close()

	require.Equal(t, executor.CallOk, instance.Call("init", testGasLimit).Kind)

	pointsUsed, err := instance.GetPointsUsedNew()
	require.Nil(t, err)
	require.Greater(t, pointsUsed, uint64(0))

	require.Nil(t, instance.ResetNew())

	hooks.storage = make(map[string]int64)
	require.Equal(t, executor.CallOk, instance.Call("init", testGasLimit).Kind)
	require.Equal(t, map[string]int64{"COUNTER": 1}, hooks.storage)
}

func TestWasmer2InstanceNew_ResetArmsTheInstanceAfterAnEarlyExit(t *testing.T) {
	wasmerExecutor, hooks, _ := newTestExecutor(t)
	instance := newCounterInstance(t, wasmerExecutor)
	defer instance.Close()

	hooks.failOn = "int64storageStore"
	hooks.failWith = errors.New("storage is unavailable")
	require.Equal(t, executor.CallEarlyExit, instance.Call("init", testGasLimit).Kind)

	// The breakpoint the failing hook tripped is still set. Clearing it is the executor's
	// job, and it happens on the next call rather than only on a reset - otherwise an
	// instance would be single-use after any early exit.
	hooks.failOn = ""
	require.Equal(t, executor.CallOk, instance.Call("init", testGasLimit).Kind)
	require.Equal(t, map[string]int64{"COUNTER": 1}, hooks.storage)
}

func TestWasmer2InstanceNew_CacheRoundTrip(t *testing.T) {
	wasmerExecutor, hooks, _ := newTestExecutor(t)
	instance := newCounterInstance(t, wasmerExecutor)
	defer instance.Close()

	cacheBytes, err := instance.Cache()
	require.Nil(t, err)
	require.NotEmpty(t, cacheBytes)

	restored, err := wasmerExecutor.NewInstanceFromCache(cacheBytes, testCompilationOptions(t))
	require.Nil(t, err)
	defer restored.(*Wasmer2InstanceNew).Close()

	require.Equal(t, executor.CallOk, restored.Call("init", testGasLimit).Kind)
	require.Equal(t, map[string]int64{"COUNTER": 1}, hooks.storage)
}

func TestWasmer2ExecutorNew_RejectsEmptyInput(t *testing.T) {
	wasmerExecutor, _, _ := newTestExecutor(t)

	_, err := wasmerExecutor.NewInstance(nil, testCompilationOptions(t))
	require.ErrorIs(t, err, ErrInvalidBytecode)

	_, err = wasmerExecutor.NewInstanceFromCache(nil, testCompilationOptions(t))
	require.ErrorIs(t, err, ErrEmptyCacheArtifact)
}

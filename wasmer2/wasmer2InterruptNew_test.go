package wasmer2

import (
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/stretchr/testify/require"
)

// spinContract announces itself through a hook and then loops forever.
//
// The announcement is what makes the test deterministic: the stop is requested from the hook
// itself, so the contract is provably mid-call when it arrives, with no sleeping and nothing
// to race. The loop is unbounded, so only the stop - or running out of gas - ends the call.
const spinContract = `
(module
    (import "env" "getNumArguments" (func $announce (result i32)))
    (memory (export "memory") 1)

    (func (export "spin") (local $i i32)
        (drop (call $announce))
        (loop $l
            (local.set $i (i32.add (local.get $i) (i32.const 1)))
            (br $l)))
)
`

// The limit bounds a failure instead of hanging the suite: if the stop never lands, the
// contract runs out of gas and the call comes back as a breakpoint, which the assertion
// catches.
const interruptGasLimit = uint64(1_000_000_000)

func newSpinInstance(t *testing.T) (*Wasmer2InstanceNew, *stubVMHooks) {
	wasmerExecutor, hooks, _ := newTestExecutor(t)

	instance, err := wasmerExecutor.NewInstance([]byte(spinContract), testCompilationOptions(t))
	require.Nil(t, err)
	t.Cleanup(instance.(*Wasmer2InstanceNew).Close)

	return instance.(*Wasmer2InstanceNew), hooks
}

func TestWasmer2InstanceNew_RequestStopEndsARunningCall(t *testing.T) {
	instance, hooks := newSpinInstance(t)

	// Requested from inside the contract's own first hook, so the call is certainly running.
	hooks.onHook = func() { instance.RequestStop() }

	result := instance.Call("spin", interruptGasLimit)

	require.Equal(t, executor.CallInterrupted, result.Kind,
		"a contract that never returns was stopped by something other than the request")
}

// The instance survives being stopped: the next call arms it afresh, which is what lets a
// warm instance be reused after a deadline was hit.
func TestWasmer2InstanceNew_UsableAfterAStop(t *testing.T) {
	instance, hooks := newSpinInstance(t)

	hooks.onHook = func() { instance.RequestStop() }
	require.Equal(t, executor.CallInterrupted, instance.Call("spin", interruptGasLimit).Kind)

	// Let it run this time, and bound it with gas instead.
	hooks.onHook = nil
	require.Equal(t, executor.CallBreakpoint, instance.Call("spin", 100_000).Kind)
}

// Stopping when nothing is running must not poison the next call: Call arms the instance
// before it starts, so the stale request is cleared rather than honoured.
func TestWasmer2InstanceNew_StopWithNoCallRunningIsIgnored(t *testing.T) {
	instance, hooks := newSpinInstance(t)
	_ = hooks

	instance.RequestStop()

	require.Equal(t, executor.CallBreakpoint, instance.Call("spin", 100_000).Kind,
		"a stop requested before the call leaked into it")
}

func TestWasmer2InstanceNew_StopOnAClosedInstanceIsANoOp(t *testing.T) {
	instance, _ := newSpinInstance(t)

	instance.Close()

	require.NotPanics(t, instance.RequestStop)
}

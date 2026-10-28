package wasmer2new

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
const interruptGasLimit = uint64(100_000_000)

// announcingVMHooks implements the single hook spinContract imports, and runs a callback when
// the contract reaches it.
type announcingVMHooks struct {
	executor.VMHooks

	onAnnounce func()
}

func (hooks *announcingVMHooks) GetNumArguments() (int32, error) {
	if hooks.onAnnounce != nil {
		hooks.onAnnounce()
	}

	return 0, nil
}

type announcingRuntime struct {
	onAnnounce func()
}

func (runtime *announcingRuntime) VMHooksForCall(_ executor.InstanceState) executor.VMHooks {
	return &announcingVMHooks{onAnnounce: runtime.onAnnounce}
}

func (runtime *announcingRuntime) IsInterfaceNil() bool {
	return runtime == nil
}

// newSpinInstance builds an instance of spinContract whose hook runs onAnnounce.
//
// The callback is set after construction because it needs the instance it is going to stop.
func newSpinInstance(t *testing.T) (*Wasmer2Instance, *func()) {
	var onAnnounce func()
	runtime := &announcingRuntime{onAnnounce: func() {
		if onAnnounce != nil {
			onAnnounce()
		}
	}}

	instance := newInstanceOn(t, newTestExecutorFor(t, runtime), []byte(spinContract))

	return instance, &onAnnounce
}

func TestRequestStopEndsARunningCall(t *testing.T) {
	instance, onAnnounce := newSpinInstance(t)

	// Requested from inside the contract's own first hook, so the call is certainly running.
	*onAnnounce = func() { instance.RequestStop() }

	result := instance.Call("spin", interruptGasLimit)

	require.Equal(t, executor.CallInterrupted, result.Kind,
		"a contract that never returns was stopped by something other than the request")
}

// The instance survives being stopped: the next call runs normally, which is what lets a warm
// instance be reused after a deadline was hit.
func TestInstanceUsableAfterAStop(t *testing.T) {
	instance, onAnnounce := newSpinInstance(t)

	*onAnnounce = func() { instance.RequestStop() }
	require.Equal(t, executor.CallInterrupted, instance.Call("spin", interruptGasLimit).Kind)

	// Let it run this time, and bound it with gas instead.
	*onAnnounce = nil
	require.Equal(t, executor.CallBreakpoint, instance.Call("spin", 1_000_000).Kind)
}

// A stop requested when nothing is running must not leak into the next call.
func TestStopWithNoCallRunningIsIgnored(t *testing.T) {
	instance, _ := newSpinInstance(t)

	instance.RequestStop()

	require.Equal(t, executor.CallBreakpoint, instance.Call("spin", 1_000_000).Kind,
		"a stop requested before the call leaked into it")
}

// The handle is released with the instance, and asking a closed instance to stop is a no-op
// rather than a read of freed memory.
func TestStopOnAClosedInstanceIsANoOp(t *testing.T) {
	instance, _ := newSpinInstance(t)

	instance.Close()

	require.NotPanics(t, instance.RequestStop)
}

// The handle survives a reset, which rebuilds the instance behind it.
func TestStopStillWorksAfterAReset(t *testing.T) {
	instance, onAnnounce := newSpinInstance(t)

	require.Nil(t, instance.ResetNew())

	*onAnnounce = func() { instance.RequestStop() }
	require.Equal(t, executor.CallInterrupted, instance.Call("spin", interruptGasLimit).Kind)
}

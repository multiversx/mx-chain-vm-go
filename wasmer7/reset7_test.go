package wasmer7

// Instance reset, from the Go side.
//
// The backend's own tests already pin what reset clears, against globals and memory this
// layer cannot see. What is left to check here is the wiring: that Reset reaches the library,
// that the handle survives it, and that a reset instance is still a working instance — which
// is what the instance tracker relies on when it serves a warm instance.

import (
	"testing"

	"github.com/multiversx/mx-chain-vm-go/wasmercommon"
	"github.com/stretchr/testify/require"
)

func TestResetKeepsTheInstanceUsable(t *testing.T) {
	instance, _ := newTestInstance(t, []byte(hooksContract))

	before, err := instance.ExportedFunctionNames()
	require.Nil(t, err)

	require.Nil(t, instance.ResetInstance())

	after, err := instance.ExportedFunctionNames()
	require.Nil(t, err)
	require.Equal(t, before, after, "reset must not change the module behind the handle")

	// Still callable through the same handle.
	result, err := instance.Call("doNothing", 1000)
	require.Nil(t, err)
	require.Equal(t, wasmercommon.CallOk, result.Kind)
}

// Reset is the legacy-interface entry point the instance tracker calls, and it must report
// success for a warm instance to be reused at all.
func TestResetReportsSuccess(t *testing.T) {
	instance, _ := newTestInstance(t, []byte(hooksContract))

	require.True(t, instance.Reset())
	require.True(t, instance.Reset(), "reset must be repeatable")
}

// Gas is tracked in a wasm global, so a reset has to clear it — otherwise a reused instance
// would start its next execution already partly spent.
func TestResetClearsGasUsed(t *testing.T) {
	instance, _ := newTestInstance(t, []byte(hooksContract))

	_, err := instance.Call("doNothing", 1000)
	require.Nil(t, err)

	require.Nil(t, instance.ResetInstance())

	used, err := instance.PointsUsed()
	require.Nil(t, err)
	require.Zero(t, used, "the metering global survived the reset")
}

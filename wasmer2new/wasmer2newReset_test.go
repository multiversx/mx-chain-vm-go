package wasmer2new

// Instance reset, from the Go side.
//
// The backend's own tests already pin what reset clears, against globals and memory this
// layer cannot see. What is left to check here is the wiring: that ResetNew reaches the
// library, that the handle survives it, and that a reset instance is still a working instance
// — which is what the instance tracker relies on when it serves a warm instance.

import (
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/stretchr/testify/require"
)

func TestResetKeepsTheInstanceUsable(t *testing.T) {
	instance, _ := newTestInstance(t, []byte(hooksContract))

	before := instance.GetExportedFunctionNames()

	require.Nil(t, instance.ResetNew())

	require.Equal(t, before, instance.GetExportedFunctionNames(),
		"reset must not change the module behind the handle")

	// Still callable through the same handle.
	require.Equal(t, executor.CallOk, instance.Call("doNothing", 1000).Kind)
}

// Rebuilt in place, as executor.InstanceNew requires: the VM holds this instance across the
// reset, so the identity must survive it.
func TestResetIsRepeatable(t *testing.T) {
	instance, _ := newTestInstance(t, []byte(hooksContract))

	require.Nil(t, instance.ResetNew())
	require.Nil(t, instance.ResetNew())
}

// Gas is tracked in a wasm global, so a reset has to clear it — otherwise a reused instance
// would start its next execution already partly spent.
func TestResetClearsGasUsed(t *testing.T) {
	instance, _ := newTestInstance(t, []byte(hooksContract))

	require.Equal(t, executor.CallOk, instance.Call("doNothing", 1000).Kind)

	require.Nil(t, instance.ResetNew())

	used, err := instance.GetPointsUsedNew()
	require.Nil(t, err)
	require.Zero(t, used, "the metering global survived the reset")
}

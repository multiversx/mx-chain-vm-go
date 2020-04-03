package contexts

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/stretchr/testify/require"
)

func TestReservedFunctions_IsFunctionReserved(t *testing.T) {
	reserved := NewReservedFunctions([]string{}, make(ProtocolReservedFunctions, 0))

	require.False(t, reserved.IsReserved("foo"))
	require.True(t, reserved.IsReserved("claimDeveloperRewards"))
	require.True(t, reserved.IsReserved(arwen.UpgradeFunctionName))
}

func TestReservedFunctions_IsFunctionReservedExplicit(t *testing.T) {
	reserved := NewReservedFunctions([]string{"rockets"}, make(ProtocolReservedFunctions, 0))

	require.False(t, reserved.IsReserved("foo"))
	require.True(t, reserved.IsReserved("rockets"))
	require.True(t, reserved.IsReserved("claimDeveloperRewards"))
}

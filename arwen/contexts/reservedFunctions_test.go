package contexts

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/stretchr/testify/require"
)

func TestReservedFunctions_IsFunctionReserved(t *testing.T) {
	reserved := NewReservedFunctions()

	require.False(t, reserved.IsReserved("foo"))
	require.True(t, reserved.IsReserved("claimDeveloperRewards"))
	require.True(t, reserved.IsReserved(arwen.UpgradeFunctionName))
}

func TestReservedFunctions_GetReserved(t *testing.T) {
	reserved := NewReservedFunctions()
	actualReserved := []string{"claimDeveloperRewards", arwen.UpgradeFunctionName}
	require.ElementsMatch(t, actualReserved, reserved.GetReserved())
}

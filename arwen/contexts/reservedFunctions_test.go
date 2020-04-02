package contexts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReservedFunctions_IsFunctionReserved(t *testing.T) {
	reserved := NewReservedFunctions([]string{})

	require.False(t, reserved.IsReserved("foo"))
	require.True(t, reserved.IsReserved("claimDeveloperRewards"))
}

func TestReservedFunctions_IsFunctionReservedExplicit(t *testing.T) {
	reserved := NewReservedFunctions([]string{"rockets"})

	require.False(t, reserved.IsReserved("foo"))
	require.True(t, reserved.IsReserved("rockets"))
	require.True(t, reserved.IsReserved("claimDeveloperRewards"))
}

func TestReservedFunctions_GetReserved(t *testing.T) {
	reserved := NewReservedFunctions([]string{})
	require.ElementsMatch(t, []string{"claimDeveloperRewards"}, reserved.GetReserved())
}

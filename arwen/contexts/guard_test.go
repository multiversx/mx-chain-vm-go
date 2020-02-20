package contexts

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFunctionsGuard_isValidFunctionName(t *testing.T) {
	guard := newFunctionsGuard(nil)

	require.True(t, guard.isValidFunctionName("foo"))
	require.True(t, guard.isValidFunctionName("_"))
	require.True(t, guard.isValidFunctionName("a"))
	require.True(t, guard.isValidFunctionName("i"))

	require.False(t, guard.isValidFunctionName(""))
	require.False(t, guard.isValidFunctionName("â"))
	require.False(t, guard.isValidFunctionName("ș"))
	require.False(t, guard.isValidFunctionName("Ä"))

	require.False(t, guard.isValidFunctionName("claimDeveloperRewards"))

	require.True(t, guard.isValidFunctionName(strings.Repeat("_", 255)))
	require.False(t, guard.isValidFunctionName(strings.Repeat("_", 256)))
}

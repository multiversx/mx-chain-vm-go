package contexts

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFunctionsGuard_isValidFunctionName(t *testing.T) {
	validator := NewWASMValidator()

	require.True(t, validator.isValidFunctionName("foo"))
	require.True(t, validator.isValidFunctionName("_"))
	require.True(t, validator.isValidFunctionName("a"))
	require.True(t, validator.isValidFunctionName("i"))

	require.False(t, validator.isValidFunctionName(""))
	require.False(t, validator.isValidFunctionName("â"))
	require.False(t, validator.isValidFunctionName("ș"))
	require.False(t, validator.isValidFunctionName("Ä"))

	require.False(t, validator.isValidFunctionName("claimDeveloperRewards"))

	require.True(t, validator.isValidFunctionName(strings.Repeat("_", 255)))
	require.False(t, validator.isValidFunctionName(strings.Repeat("_", 256)))
}

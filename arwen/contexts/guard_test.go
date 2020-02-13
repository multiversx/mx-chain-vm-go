package contexts

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGuard_isValidFunctionName(t *testing.T) {
	require.True(t, isValidFunctionName("foo"))
	require.True(t, isValidFunctionName("_"))
	require.True(t, isValidFunctionName("a"))
	require.True(t, isValidFunctionName("i"))

	require.False(t, isValidFunctionName(""))
	require.False(t, isValidFunctionName("â"))
	require.False(t, isValidFunctionName("ș"))
	require.False(t, isValidFunctionName("Ä"))

	require.True(t, isValidFunctionName(strings.Repeat("_", 255)))
	require.False(t, isValidFunctionName(strings.Repeat("_", 256)))
}

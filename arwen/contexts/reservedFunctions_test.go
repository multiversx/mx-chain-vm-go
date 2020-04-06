package contexts

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/stretchr/testify/require"
)

func TestReservedFunctions_IsFunctionReserved(t *testing.T) {
	scAPINames := []string{
		"rockets",
	}

	fromProtocol := []string{
		"protocolFunctionFoo",
		"protocolFunctionBar",
	}

	reserved := NewReservedFunctions(scAPINames, fromProtocol)

	require.False(t, reserved.IsReserved("foo"))
	require.True(t, reserved.IsReserved("rockets"))
	require.True(t, reserved.IsReserved("protocolFunctionFoo"))
	require.True(t, reserved.IsReserved("protocolFunctionBar"))
	require.True(t, reserved.IsReserved(arwen.UpgradeFunctionName))
}

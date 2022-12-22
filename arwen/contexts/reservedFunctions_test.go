package contexts

import (
	"testing"

	"github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/builtInFunctions"
	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen"
	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen/mock"
	"github.com/stretchr/testify/require"
)

func TestReservedFunctions_IsFunctionReserved(t *testing.T) {
	scAPINames := vmcommon.FunctionNames{
		"rockets": {},
	}

	builtInFuncContainer := builtInFunctions.NewBuiltInFunctionContainer()
	_ = builtInFuncContainer.Add("protocolFunctionFoo", &mock.BuiltInFunctionStub{})
	_ = builtInFuncContainer.Add("protocolFunctionBar", &mock.BuiltInFunctionStub{})

	reserved := NewReservedFunctions(scAPINames, builtInFuncContainer)

	require.False(t, reserved.IsReserved("foo"))
	require.True(t, reserved.IsReserved("rockets"))
	require.True(t, reserved.IsReserved("protocolFunctionFoo"))
	require.True(t, reserved.IsReserved("protocolFunctionBar"))
	require.True(t, reserved.IsReserved(arwen.UpgradeFunctionName))
}

package contexts

import (
	"testing"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
	"github.com/multiversx/mx-chain-vm-v1_4-go/arwen"
	"github.com/multiversx/mx-chain-vm-v1_4-go/arwen/mock"
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

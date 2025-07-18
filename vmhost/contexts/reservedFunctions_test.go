package contexts

import (
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/mock"
	"github.com/stretchr/testify/require"
)

func TestReservedFunctions_IsFunctionReserved(t *testing.T) {
	scAPINames := vmcommon.FunctionNames{
		"rockets": struct{}{},
	}

	builtInFuncContainer := builtInFunctions.NewBuiltInFunctionContainer()
	_ = builtInFuncContainer.Add("protocolFunctionFoo", &mock.BuiltInFunctionStub{})
	_ = builtInFuncContainer.Add("protocolFunctionBar", &mock.BuiltInFunctionStub{})

	enableEpochsHandler := worldmock.EnableEpochsHandlerStubAllFlags()

	activationFlags := map[string]core.EnableEpochFlag{}

	reserved := NewReservedFunctions(scAPINames, builtInFuncContainer, activationFlags, enableEpochsHandler)

	require.False(t, reserved.IsReserved("foo"))
	require.True(t, reserved.IsReserved("rockets"))
	require.True(t, reserved.IsReserved("protocolFunctionFoo"))
	require.True(t, reserved.IsReserved("protocolFunctionBar"))
	require.True(t, reserved.IsReserved(vmhost.DeleteFunctionName))
}

func TestReservedFunctions_IsFunctionReserved_WithBarnardOpcodesFlag(t *testing.T) {
	scAPINames := vmcommon.FunctionNames{
		"rockets": struct{}{},
		"paper":   struct{}{},
	}

	builtInFuncContainer := builtInFunctions.NewBuiltInFunctionContainer()
	_ = builtInFuncContainer.Add("protocolFunctionFoo", &mock.BuiltInFunctionStub{})
	_ = builtInFuncContainer.Add("protocolFunctionBar", &mock.BuiltInFunctionStub{})

	paperFlag := core.EnableEpochFlag("paper_flag")

	activationFlags := map[string]core.EnableEpochFlag{
		"paper": paperFlag,
	}

	enableEpochsHandler := worldmock.EnableEpochsHandlerStubAllFlags()
	enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
		if flag == vmhost.BarnardOpcodesFlag {
			return true
		}

		if flag == paperFlag {
			return false
		}

		return true
	}

	reserved := NewReservedFunctions(scAPINames, builtInFuncContainer, activationFlags, enableEpochsHandler)

	require.False(t, reserved.IsReserved("foo"))
	require.True(t, reserved.IsReserved("paper"))
	require.True(t, reserved.IsReserved("rockets"))
	require.True(t, reserved.IsReserved("protocolFunctionFoo"))
	require.True(t, reserved.IsReserved("protocolFunctionBar"))
	require.True(t, reserved.IsReserved(vmhost.DeleteFunctionName))

	enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
		if flag == vmhost.BarnardOpcodesFlag {
			return true
		}

		if flag == paperFlag {
			return true
		}

		return true
	}

	require.False(t, reserved.IsReserved("foo"))
	require.True(t, reserved.IsReserved("paper"))
	require.True(t, reserved.IsReserved("rockets"))
	require.True(t, reserved.IsReserved("protocolFunctionFoo"))
	require.True(t, reserved.IsReserved("protocolFunctionBar"))
	require.True(t, reserved.IsReserved(vmhost.DeleteFunctionName))
}

func TestReservedFunctions_IsFunctionReserved_WithBarnardOpcodesFlag_Disabled(t *testing.T) {
	scAPINames := vmcommon.FunctionNames{
		"rockets": struct{}{},
		"paper":   struct{}{},
	}

	builtInFuncContainer := builtInFunctions.NewBuiltInFunctionContainer()
	_ = builtInFuncContainer.Add("protocolFunctionFoo", &mock.BuiltInFunctionStub{})
	_ = builtInFuncContainer.Add("protocolFunctionBar", &mock.BuiltInFunctionStub{})

	paperFlag := core.EnableEpochFlag("paper_flag")

	activationFlags := map[string]core.EnableEpochFlag{
		"paper": paperFlag,
	}

	enableEpochsHandler := worldmock.EnableEpochsHandlerStubAllFlags()
	enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
		if flag == vmhost.BarnardOpcodesFlag {
			return false
		}

		if flag == paperFlag {
			return false
		}

		return true
	}

	reserved := NewReservedFunctions(scAPINames, builtInFuncContainer, activationFlags, enableEpochsHandler)

	require.False(t, reserved.IsReserved("foo"))
	require.False(t, reserved.IsReserved("paper"))
	require.True(t, reserved.IsReserved("rockets"))
	require.True(t, reserved.IsReserved("protocolFunctionFoo"))
	require.True(t, reserved.IsReserved("protocolFunctionBar"))
	require.True(t, reserved.IsReserved(vmhost.DeleteFunctionName))

	enableEpochsHandler.IsFlagEnabledCalled = func(flag core.EnableEpochFlag) bool {
		if flag == vmhost.BarnardOpcodesFlag {
			return false
		}

		if flag == paperFlag {
			return true
		}

		return true
	}

	require.False(t, reserved.IsReserved("foo"))
	require.True(t, reserved.IsReserved("paper"))
	require.True(t, reserved.IsReserved("rockets"))
	require.True(t, reserved.IsReserved("protocolFunctionFoo"))
	require.True(t, reserved.IsReserved("protocolFunctionBar"))
	require.True(t, reserved.IsReserved(vmhost.DeleteFunctionName))
}

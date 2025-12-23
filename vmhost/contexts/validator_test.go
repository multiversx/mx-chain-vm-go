package contexts

import (
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/vmhost/mock"
	"github.com/stretchr/testify/require"
)

func testImportNames() vmcommon.FunctionNames {
	importNames := make(vmcommon.FunctionNames)
	var empty struct{}
	importNames["getArgument"] = empty
	importNames["asyncCall"] = empty
	return importNames
}

func TestFunctionsGuard_isValidFunctionName(t *testing.T) {
	builtInFuncContainer := builtInFunctions.NewBuiltInFunctionContainer()
	_ = builtInFuncContainer.Add("protocolFunctionFoo", &mock.BuiltInFunctionStub{})
	_ = builtInFuncContainer.Add("protocolFunctionBar", &mock.BuiltInFunctionStub{})

	validator := NewWASMValidator(testImportNames(), builtInFuncContainer, worldmock.EnableEpochsHandlerStubAllFlags())

	require.Nil(t, validator.VerifyValidFunctionName("foo"))
	require.Nil(t, validator.VerifyValidFunctionName("_"))
	require.Nil(t, validator.VerifyValidFunctionName("a"))
	require.Nil(t, validator.VerifyValidFunctionName("i"))

	require.NotNil(t, validator.VerifyValidFunctionName(""))
	require.NotNil(t, validator.VerifyValidFunctionName("3"))
	require.NotNil(t, validator.VerifyValidFunctionName("π"))
	require.NotNil(t, validator.VerifyValidFunctionName("2foo"))
	require.NotNil(t, validator.VerifyValidFunctionName("-"))
	require.NotNil(t, validator.VerifyValidFunctionName("â"))
	require.NotNil(t, validator.VerifyValidFunctionName("ș"))
	require.NotNil(t, validator.VerifyValidFunctionName("Ä"))

	require.NotNil(t, validator.VerifyValidFunctionName("protocolFunctionFoo"))
	require.NotNil(t, validator.VerifyValidFunctionName("protocolFunctionBar"))

	require.Nil(t, validator.VerifyValidFunctionName(strings.Repeat("_", 255)))
	require.NotNil(t, validator.VerifyValidFunctionName(strings.Repeat("_", 256)))

	require.NotNil(t, validator.VerifyValidFunctionName("getArgument"))
	require.NotNil(t, validator.VerifyValidFunctionName("asyncCall"))
	require.Nil(t, validator.VerifyValidFunctionName("getArgument55"))
}

func TestFunctionsProtected(t *testing.T) {
	host := InitializeVMAndWasmer()

	validator := NewWASMValidator(testImportNames(), builtInFunctions.NewBuiltInFunctionContainer(), worldmock.EnableEpochsHandlerStubAllFlags())

	world := worldmock.NewMockWorld()
	imb := contextmock.NewExecutorMock(world)
	instance := imb.CreateAndStoreInstanceMock(t, host, []byte{}, []byte{}, []byte{}, []byte{}, 0, 0, false)

	instance.AddMockMethod("transferValueOnly", func() *contextmock.InstanceMock {
		testHost := instance.Host
		testInstance := contextmock.GetMockInstance(testHost)
		return testInstance
	})

	err := validator.VerifyProtectedFunctions(instance)
	require.NotNil(t, err)
}

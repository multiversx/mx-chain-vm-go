package contexts

import (
	"strings"
	"testing"

	"github.com/ElrondNetwork/elrond-vm-common/builtInFunctions"
	"github.com/ElrondNetwork/wasm-vm/arwen/mock"
	contextmock "github.com/ElrondNetwork/wasm-vm/mock/context"
	worldmock "github.com/ElrondNetwork/wasm-vm/mock/world"
	"github.com/stretchr/testify/require"
)

func TestFunctionsGuard_isValidFunctionName(t *testing.T) {
	imports := MakeAPIImports()

	builtInFuncContainer := builtInFunctions.NewBuiltInFunctionContainer()
	_ = builtInFuncContainer.Add("protocolFunctionFoo", &mock.BuiltInFunctionStub{})
	_ = builtInFuncContainer.Add("protocolFunctionBar", &mock.BuiltInFunctionStub{})

	validator := newWASMValidator(imports.Names(), builtInFuncContainer)

	require.Nil(t, validator.verifyValidFunctionName("foo"))
	require.Nil(t, validator.verifyValidFunctionName("_"))
	require.Nil(t, validator.verifyValidFunctionName("a"))
	require.Nil(t, validator.verifyValidFunctionName("i"))

	require.NotNil(t, validator.verifyValidFunctionName(""))
	require.NotNil(t, validator.verifyValidFunctionName("3"))
	require.NotNil(t, validator.verifyValidFunctionName("π"))
	require.NotNil(t, validator.verifyValidFunctionName("2foo"))
	require.NotNil(t, validator.verifyValidFunctionName("-"))
	require.NotNil(t, validator.verifyValidFunctionName("â"))
	require.NotNil(t, validator.verifyValidFunctionName("ș"))
	require.NotNil(t, validator.verifyValidFunctionName("Ä"))

	require.NotNil(t, validator.verifyValidFunctionName("protocolFunctionFoo"))
	require.NotNil(t, validator.verifyValidFunctionName("protocolFunctionBar"))

	require.Nil(t, validator.verifyValidFunctionName(strings.Repeat("_", 255)))
	require.NotNil(t, validator.verifyValidFunctionName(strings.Repeat("_", 256)))

	require.NotNil(t, validator.verifyValidFunctionName("getArgument"))
	require.NotNil(t, validator.verifyValidFunctionName("asyncCall"))
	require.Nil(t, validator.verifyValidFunctionName("getArgument55"))
}

func TestFunctionsProtected(t *testing.T) {
	host := InitializeArwenAndWasmer()
	imports := host.SCAPIMethods

	validator := newWASMValidator(imports.Names(), builtInFunctions.NewBuiltInFunctionContainer())

	world := worldmock.NewMockWorld()
	imb := contextmock.NewInstanceBuilderMock(world)
	instance := imb.CreateAndStoreInstanceMock(t, host, []byte{}, []byte{}, []byte{}, []byte{}, 0, 0, false)

	instance.AddMockMethod("transferValueOnly", func() *contextmock.InstanceMock {
		host := instance.Host
		instance := contextmock.GetMockInstance(host)
		return instance
	})

	err := validator.verifyProtectedFunctions(instance)
	require.NotNil(t, err)
}

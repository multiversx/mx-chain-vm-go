package contexts

import (
	"strings"
	"testing"

	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen"
	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen/mock"
	contextmock "github.com/ElrondNetwork/wasm-vm-v1_4/mock/context"
	worldmock "github.com/ElrondNetwork/wasm-vm-v1_4/mock/world"
	"github.com/ElrondNetwork/wasm-vm-v1_4/wasmer"
	"github.com/ElrondNetwork/elrond-vm-common/builtInFunctions"
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

func TestFunctionsGuard_Arity(t *testing.T) {
	host := InitializeArwenAndWasmer()
	imports := host.SCAPIMethods

	validator := newWASMValidator(imports.Names(), builtInFunctions.NewBuiltInFunctionContainer())

	gasLimit := uint64(100000000)
	path := "./../../test/contracts/signatures/output/signatures.wasm"
	contractCode := arwen.GetSCCode(path)
	options := wasmer.CompilationOptions{
		GasLimit:           gasLimit,
		OpcodeTrace:        false,
		Metering:           true,
		RuntimeBreakpoints: true,
	}
	instance, err := wasmer.NewInstanceWithOptions(contractCode, options)
	require.Nil(t, err)

	inArity, _ := validator.getInputArity(instance, "goodFunction")
	require.Equal(t, 0, inArity)

	outArity, _ := validator.getOutputArity(instance, "goodFunction")
	require.Equal(t, 0, outArity)

	inArity, _ = validator.getInputArity(instance, "wrongReturn")
	require.Equal(t, 0, inArity)

	outArity, _ = validator.getOutputArity(instance, "wrongReturn")
	require.Equal(t, 1, outArity)

	inArity, _ = validator.getInputArity(instance, "wrongParams")
	require.Equal(t, 1, inArity)

	outArity, _ = validator.getOutputArity(instance, "wrongParams")
	require.Equal(t, 0, outArity)

	inArity, _ = validator.getInputArity(instance, "wrongParamsAndReturn")
	require.Equal(t, 2, inArity)

	outArity, _ = validator.getOutputArity(instance, "wrongParamsAndReturn")
	require.Equal(t, 1, outArity)

	err = validator.verifyVoidFunction(instance, "goodFunction")
	require.Nil(t, err)

	err = validator.verifyVoidFunction(instance, "wrongReturn")
	require.NotNil(t, err)

	err = validator.verifyVoidFunction(instance, "wrongParams")
	require.NotNil(t, err)

	err = validator.verifyVoidFunction(instance, "wrongParamsAndReturn")
	require.NotNil(t, err)
}

func TestFunctionsProtected(t *testing.T) {
	host := InitializeArwenAndWasmer()
	imports := host.SCAPIMethods

	validator := newWASMValidator(imports.Names(), builtInFunctions.NewBuiltInFunctionContainer())

	world := worldmock.NewMockWorld()
	imb := contextmock.NewInstanceBuilderMock(world)
	instance := imb.CreateAndStoreInstanceMock(t, host, []byte{}, []byte{}, []byte{}, []byte{}, 0, 0)

	instance.AddMockMethod("transferValueOnly", func() *contextmock.InstanceMock {
		host := instance.Host
		instance := contextmock.GetMockInstance(host)
		return instance
	})

	err := validator.verifyProtectedFunctions(instance)
	require.NotNil(t, err)
}

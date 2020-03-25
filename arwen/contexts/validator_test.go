package contexts

import (
	"strings"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/stretchr/testify/require"
)

func TestFunctionsGuard_isValidFunctionName(t *testing.T) {
	imports := MakeAPIImports()
	validator := NewWASMValidator(imports.Names())

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

	require.False(t, validator.isValidFunctionName("getArgument"))
	require.False(t, validator.isValidFunctionName("asyncCall"))
	require.True(t, validator.isValidFunctionName("getArgument55"))
}

func TestFunctionsGuard_Arity(t *testing.T) {
	imports := InitializeWasmer()
	validator := NewWASMValidator(imports.Names())

	gasLimit := uint64(100000000)
	path := "./../../test/contracts/signatures/signatures.wasm"
	contractCode := arwen.GetSCCode(path)
	instance, err := wasmer.NewMeteredInstance(contractCode, gasLimit)
	require.Nil(t, err)

	require.Equal(t, 0, validator.getInputArity(instance, "goodFunction"))
	require.Equal(t, 0, validator.getOutputArity(instance, "goodFunction"))

	require.Equal(t, 0, validator.getInputArity(instance, "wrongReturn"))
	require.Equal(t, 1, validator.getOutputArity(instance, "wrongReturn"))

	require.Equal(t, 1, validator.getInputArity(instance, "wrongParams"))
	require.Equal(t, 0, validator.getOutputArity(instance, "wrongParams"))

	require.Equal(t, 2, validator.getInputArity(instance, "wrongParamsAndReturn"))
	require.Equal(t, 1, validator.getOutputArity(instance, "wrongParamsAndReturn"))

	require.Equal(t, true, validator.isVoidFunction(instance, "goodFunction"))
	require.Equal(t, false, validator.isVoidFunction(instance, "wrongReturn"))
	require.Equal(t, false, validator.isVoidFunction(instance, "wrongParams"))
	require.Equal(t, false, validator.isVoidFunction(instance, "wrongParamsAndReturn"))
}

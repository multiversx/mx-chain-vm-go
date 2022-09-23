package wasmer

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/ElrondNetwork/wasm-vm/executor"
	"github.com/stretchr/testify/require"
)

// GetSCCode retrieves the bytecode of a WASM module from a file.
func getSCCode(fileName string) []byte {
	code, err := ioutil.ReadFile(filepath.Clean(fileName))
	if err != nil {
		panic(fmt.Sprintf("GetSCCode(): %s", fileName))
	}
	return code
}

func TestFunctionsGuard_Arity(t *testing.T) {
	// Empty imports on purpose.
	// We have currently no access to the elrondapi package here, due to cyclic imports.
	// Fortunately, imports are not necessary for this test.
	imports := executor.NewImportFunctions()
	SetImports(ConvertImports(imports))

	gasLimit := uint64(100000000)
	path := "./../test/contracts/signatures/output/signatures.wasm"
	contractCode := getSCCode(path)
	options := executor.CompilationOptions{
		GasLimit:           gasLimit,
		OpcodeTrace:        false,
		Metering:           true,
		RuntimeBreakpoints: true,
	}
	instance, err := NewInstanceWithOptions(contractCode, options)
	require.Nil(t, err)

	inArity, _ := instance.getInputArity("goodFunction")
	require.Equal(t, 0, inArity)

	outArity, _ := instance.getOutputArity("goodFunction")
	require.Equal(t, 0, outArity)

	inArity, _ = instance.getInputArity("wrongReturn")
	require.Equal(t, 0, inArity)

	outArity, _ = instance.getOutputArity("wrongReturn")
	require.Equal(t, 1, outArity)

	inArity, _ = instance.getInputArity("wrongParams")
	require.Equal(t, 1, inArity)

	outArity, _ = instance.getOutputArity("wrongParams")
	require.Equal(t, 0, outArity)

	inArity, _ = instance.getInputArity("wrongParamsAndReturn")
	require.Equal(t, 2, inArity)

	outArity, _ = instance.getOutputArity("wrongParamsAndReturn")
	require.Equal(t, 1, outArity)

	err = instance.verifyVoidFunction("goodFunction")
	require.Nil(t, err)

	err = instance.verifyVoidFunction("wrongReturn")
	require.NotNil(t, err)

	err = instance.verifyVoidFunction("wrongParams")
	require.NotNil(t, err)

	err = instance.verifyVoidFunction("wrongParamsAndReturn")
	require.NotNil(t, err)
}

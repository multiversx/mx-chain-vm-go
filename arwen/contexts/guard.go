package contexts

import (
	"fmt"
	"unicode"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
)

func (context *runtimeContext) VerifyContractCode() error {
	err := verifyMemoryDeclaration(context.instance)
	if err != nil {
		return err
	}

	err = verifyFunctionsNames(context.instance)
	if err != nil {
		return err
	}

	return nil
}

func verifyMemoryDeclaration(instance *wasmer.Instance) error {
	if !instance.HasMemory() {
		return arwen.ErrMemoryDeclarationMissing
	}

	return nil
}

func verifyFunctionsNames(instance *wasmer.Instance) error {
	for functionName := range instance.Exports {
		if !isValidFunctionName(functionName) {
			return fmt.Errorf("%v: %s", arwen.ErrInvalidFunctionName, functionName)
		}
	}

	return nil
}

func isValidFunctionName(functionName string) bool {
	const maxLengthOfFunctionName = 256

	if len(functionName) == 0 {
		return false
	}

	if len(functionName) >= maxLengthOfFunctionName {
		return false
	}

	if !isASCIIString(functionName) {
		return false
	}

	return true
}

func isASCIIString(input string) bool {
	for i := 0; i < len(input); i++ {
		if input[i] > unicode.MaxASCII {
			return false
		}
	}

	return true
}

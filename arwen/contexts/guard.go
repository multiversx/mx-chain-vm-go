package contexts

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"unicode"
)

func (context *runtimeContext) VerifyContractCode() error {
	if !context.instance.HasMemory() {
		return arwen.ErrMemoryDeclarationMissing
	}

	for functionName := range context.instance.Exports {
		if !isValidFunctionName(functionName) {
			return arwen.ErrInvalidFunctionName
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

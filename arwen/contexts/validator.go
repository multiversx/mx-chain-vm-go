package contexts

import (
	"fmt"
	"unicode"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
)

type WASMValidator struct {
	reserved *ReservedFunctions
}

func NewWASMValidator() *WASMValidator {
	return &WASMValidator{
		reserved: NewReservedFunctions(),
	}
}

func (validator *WASMValidator) verifyMemoryDeclaration(instance *wasmer.Instance) error {
	if !instance.HasMemory() {
		return arwen.ErrMemoryDeclarationMissing
	}

	return nil
}

func (validator *WASMValidator) verifyFunctionsNames(instance *wasmer.Instance) error {
	for functionName := range instance.Exports {
		if !validator.isValidFunctionName(functionName) {
			return fmt.Errorf("%w: %s", arwen.ErrInvalidFunctionName, functionName)
		}
	}

	return nil
}

func (validator *WASMValidator) isValidFunctionName(functionName string) bool {
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
	if validator.reserved.IsReserved(functionName) {
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

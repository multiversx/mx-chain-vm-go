package contexts

import (
	"fmt"
	"unicode"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
)

// WASMValidator is a validator for WASM SmartContracts
type WASMValidator struct {
	reserved *ReservedFunctions
}

// NewWASMValidator creates a new WASMValidator
func NewWASMValidator(scAPINames []string) *WASMValidator {
	return &WASMValidator{
		reserved: NewReservedFunctions(scAPINames),
	}
}

func (validator *WASMValidator) verifyMemoryDeclaration(instance *wasmer.Instance) error {
	if !instance.HasMemory() {
		return arwen.ErrMemoryDeclarationMissing
	}

	return nil
}

func (validator *WASMValidator) verifyFunctions(instance *wasmer.Instance) error {
	for functionName := range instance.Exports {
		if !validator.isValidFunctionName(functionName) {
			return fmt.Errorf("%w: %s", arwen.ErrInvalidFunctionName, functionName)
		}

		if !validator.isVoidFunction(instance, functionName) {
			return fmt.Errorf("%w: %s", arwen.ErrFunctionNonvoidSignature, functionName)
		}
	}

	return nil
}

func (validator *WASMValidator) isVoidFunction(instance *wasmer.Instance, functionName string) bool {
	inArity := validator.getInputArity(instance, functionName)
	outArity := validator.getOutputArity(instance, functionName)

	isVoid := (inArity == 0 && outArity == 0)
	return isVoid
}

func (validator *WASMValidator) getInputArity(instance *wasmer.Instance, functionName string) int {
	signature, ok := instance.Signatures[functionName]
	if !ok {
		return -1
	}
	return signature.InputArity
}

func (validator *WASMValidator) getOutputArity(instance *wasmer.Instance, functionName string) int {
	signature, ok := instance.Signatures[functionName]
	if !ok {
		return -1
	}
	return signature.OutputArity
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

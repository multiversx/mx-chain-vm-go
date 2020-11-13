package contexts

import (
	"fmt"
	"unicode"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-go/core/vm-common"
)

const NoArity = -1

// WASMValidator is a validator for WASM SmartContracts
type WASMValidator struct {
	reserved *ReservedFunctions
}

// NewWASMValidator creates a new WASMValidator
func NewWASMValidator(scAPINames vmcommon.FunctionNames, protocolBuiltinFunctions vmcommon.FunctionNames) *WASMValidator {
	return &WASMValidator{
		reserved: NewReservedFunctions(scAPINames, protocolBuiltinFunctions),
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
		err := validator.verifyValidFunctionName(functionName)
		if err != nil {
			return err
		}

		err = validator.verifyVoidFunction(instance, functionName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (validator *WASMValidator) verifyVoidFunction(instance *wasmer.Instance, functionName string) error {
	inArity, err := validator.getInputArity(instance, functionName)
	if err != nil {
		return err
	}

	outArity, err := validator.getOutputArity(instance, functionName)
	if err != nil {
		return err
	}

	isVoid := inArity == 0 && outArity == 0
	if !isVoid {
		return fmt.Errorf("%w: %s", arwen.ErrFunctionNonvoidSignature, functionName)
	}
	return nil
}

func (validator *WASMValidator) getInputArity(instance *wasmer.Instance, functionName string) (int, error) {
	signature, ok := instance.Signatures[functionName]
	if !ok {
		return NoArity, fmt.Errorf("%w: %s", arwen.ErrFuncNotFound, functionName)
	}
	return signature.InputArity, nil
}

func (validator *WASMValidator) getOutputArity(instance *wasmer.Instance, functionName string) (int, error) {
	signature, ok := instance.Signatures[functionName]
	if !ok {
		return NoArity, fmt.Errorf("%w: %s", arwen.ErrFuncNotFound, functionName)
	}
	return signature.OutputArity, nil
}

func (validator *WASMValidator) verifyValidFunctionName(functionName string) error {
	const maxLengthOfFunctionName = 256

	errInvalidName := fmt.Errorf("%w: %s", arwen.ErrInvalidFunctionName, functionName)

	if len(functionName) == 0 {
		return errInvalidName
	}
	if len(functionName) >= maxLengthOfFunctionName {
		return errInvalidName
	}
	if !isASCIIString(functionName) {
		return errInvalidName
	}
	if validator.reserved.IsReserved(functionName) {
		return errInvalidName
	}

	return nil
}

// TODO: Add more constraints (too loose currently)
func isASCIIString(input string) bool {
	for i := 0; i < len(input); i++ {
		if input[i] > unicode.MaxASCII {
			return false
		}
	}

	return true
}

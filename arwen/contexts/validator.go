package contexts

import (
	"fmt"
	"strings"

	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

const noArity = -1
const allowedCharsInFunctionName = "abcdefghijklmnopqrstuvwxyz0123456789_"

// wasmValidator is a validator for WASM SmartContracts
type wasmValidator struct {
	reserved *reservedFunctions
}

// newWASMValidator creates a new WASMValidator
func newWASMValidator(scAPINames vmcommon.FunctionNames, builtInFuncContainer vmcommon.BuiltInFunctionContainer) *wasmValidator {
	return &wasmValidator{
		reserved: NewReservedFunctions(scAPINames, builtInFuncContainer),
	}
}

func (validator *wasmValidator) verifyMemoryDeclaration(instance wasmer.InstanceHandler) error {
	if !instance.HasMemory() {
		return arwen.ErrMemoryDeclarationMissing
	}

	return nil
}

func (validator *wasmValidator) verifyFunctions(instance wasmer.InstanceHandler) error {
	for functionName := range instance.GetExports() {
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

var protectedFunctions = map[string]bool{
	"internalVMErrors":  true,
	"transferValueOnly": true,
	"writeLog":          true,
	"signalError":       true,
	"completedTxEvent":  true}

func (validator *wasmValidator) verifyProtectedFunctions(instance wasmer.InstanceHandler) error {
	for functionName := range instance.GetExports() {
		_, found := protectedFunctions[functionName]
		if found {
			return arwen.ErrContractInvalid
		}

	}

	return nil
}

func (validator *wasmValidator) verifyVoidFunction(instance wasmer.InstanceHandler, functionName string) error {
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

func (validator *wasmValidator) getInputArity(instance wasmer.InstanceHandler, functionName string) (int, error) {
	signature, ok := instance.GetSignature(functionName)
	if !ok {
		return noArity, fmt.Errorf("%w: %s", arwen.ErrFuncNotFound, functionName)
	}
	return signature.InputArity, nil
}

func (validator *wasmValidator) getOutputArity(instance wasmer.InstanceHandler, functionName string) (int, error) {
	signature, ok := instance.GetSignature(functionName)
	if !ok {
		return noArity, fmt.Errorf("%w: %s", arwen.ErrFuncNotFound, functionName)
	}
	return signature.OutputArity, nil
}

func (validator *wasmValidator) verifyValidFunctionName(functionName string) error {
	const maxLengthOfFunctionName = 256

	errInvalidName := fmt.Errorf("%w: %s", arwen.ErrInvalidFunctionName, functionName)

	if len(functionName) == 0 {
		return errInvalidName
	}
	if len(functionName) >= maxLengthOfFunctionName {
		return errInvalidName
	}
	if isFirstCharacterNumeric(functionName) {
		return errInvalidName
	}
	if !validCharactersOnly(functionName) {
		return errInvalidName
	}
	if validator.reserved.IsReserved(functionName) {
		return errInvalidName
	}

	return nil
}

func validCharactersOnly(input string) bool {
	input = strings.ToLower(input)
	for i := 0; i < len(input); i++ {
		c := string(input[i])
		if !strings.Contains(allowedCharsInFunctionName, c) {
			return false
		}
	}

	return true
}

func isFirstCharacterNumeric(name string) bool {
	return name[0] >= '0' && name[0] <= '9'
}

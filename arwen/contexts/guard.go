package contexts

import (
	"fmt"
	"unicode"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
)

func (context *runtimeContext) VerifyContractCode() error {
	memoryGuard := newMemoryGuard(context.instance)
	err := memoryGuard.verifyMemoryDeclaration()
	if err != nil {
		return err
	}

	functionsGuard := newFunctionsGuard(context.instance)
	err = functionsGuard.verifyFunctionsNames()
	if err != nil {
		return err
	}

	return nil
}

type memoryGuard struct {
	instance *wasmer.Instance
}

func newMemoryGuard(instance *wasmer.Instance) *memoryGuard {
	return &memoryGuard{instance: instance}
}

func (guard *memoryGuard) verifyMemoryDeclaration() error {
	if !guard.instance.HasMemory() {
		return arwen.ErrMemoryDeclarationMissing
	}

	return nil
}

type functionsGuard struct {
	instance *wasmer.Instance
	reserved *ReservedFunctions
}

func newFunctionsGuard(instance *wasmer.Instance) *functionsGuard {
	return &functionsGuard{instance: instance, reserved: NewReservedFunctions()}
}

func (guard *functionsGuard) verifyFunctionsNames() error {
	for functionName := range guard.instance.Exports {
		if !guard.isValidFunctionName(functionName) {
			return fmt.Errorf("%v: %s", arwen.ErrInvalidFunctionName, functionName)
		}
	}

	return nil
}

func (guard *functionsGuard) isValidFunctionName(functionName string) bool {
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
	if guard.reserved.IsReserved(functionName) {
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

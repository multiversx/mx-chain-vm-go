package wasmer

import (
	"fmt"

	"github.com/ElrondNetwork/wasm-vm/executor"
)

const noArity = -1

// getSignature returns the signature for the given functionName
func (instance *WasmerInstance) getSignature(functionName string) (*ExportedFunctionSignature, bool) {
	signature, ok := instance.Signatures[functionName]
	return signature, ok
}

func (instance *WasmerInstance) verifyVoidFunction(functionName string) error {
	inArity, err := instance.getInputArity(functionName)
	if err != nil {
		return err
	}

	outArity, err := instance.getOutputArity(functionName)
	if err != nil {
		return err
	}

	isVoid := inArity == 0 && outArity == 0
	if !isVoid {
		return fmt.Errorf("%w: %s", executor.ErrFunctionNonvoidSignature, functionName)
	}
	return nil
}

func (instance *WasmerInstance) getInputArity(functionName string) (int, error) {
	signature, ok := instance.getSignature(functionName)
	if !ok {
		return noArity, fmt.Errorf("%w: %s", executor.ErrFuncNotFound, functionName)
	}
	return signature.InputArity, nil
}

func (instance *WasmerInstance) getOutputArity(functionName string) (int, error) {
	signature, ok := instance.getSignature(functionName)
	if !ok {
		return noArity, fmt.Errorf("%w: %s", executor.ErrFuncNotFound, functionName)
	}
	return signature.OutputArity, nil
}

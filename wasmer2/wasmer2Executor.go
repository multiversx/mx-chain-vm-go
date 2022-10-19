package wasmer2

import (
	"errors"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/executor"
)

// WasmerExecutor oversees the creation of Wasmer instances and execution.
type Wasmer2Executor struct {
}

// NewExecutor creates a new wasmer executor.
func NewExecutor() (*Wasmer2Executor, error) {
	// functionNames, err := injectCgoFunctionPointers()
	// if err != nil {
	// 	return nil, err
	// }
	return &Wasmer2Executor{
		// eiFunctionNames: functionNames,
	}, nil
}

// SetOpcodeCosts sets gas costs globally inside the Wasmer executor.
func (wasmerExecutor *Wasmer2Executor) SetOpcodeCosts(opcodeCosts *[executor.OpcodeCount]uint32) {
	// SetOpcodeCosts(opcodeCosts)
}

// SetRkyvSerializationEnabled controls a Wasmer flag.
func (wasmerExecutor *Wasmer2Executor) SetRkyvSerializationEnabled(enabled bool) {
}

// SetSIGSEGVPassthrough controls a Wasmer flag.
func (wasmerExecutor *Wasmer2Executor) SetSIGSEGVPassthrough() {
}

func (wasmerExecutor *Wasmer2Executor) FunctionNames() vmcommon.FunctionNames {
	return functionNames
}

// NewInstanceWithOptions creates a new Wasmer instance from WASM bytecode,
// respecting the provided options
func (wasmerExecutor *Wasmer2Executor) NewInstanceWithOptions(
	contractCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	return NewInstanceWithOptions(contractCode, options)
}

// NewInstanceFromCompiledCodeWithOptions creates a new Wasmer instance from
// precompiled machine code, respecting the provided options
func (wasmerExecutor *Wasmer2Executor) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	// return NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
	return nil, errors.New("NewInstanceFromCompiledCodeWithOptions not implemented")
}

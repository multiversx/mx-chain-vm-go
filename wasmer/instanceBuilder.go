package wasmer

import (
	"github.com/ElrondNetwork/wasm-vm/executor"
)

// WasmerExecutor is the default instance builder, which produces real Wasmer instances from WASM bytecode.
type WasmerExecutor struct {
}

// NewExecutor creates a new wasmer executor.
func NewExecutor() *WasmerExecutor {
	return &WasmerExecutor{}
}

// SetOpcodeCosts sets gas costs globally inside the Wasmer executor.
func (builder *WasmerExecutor) SetOpcodeCosts(opcodeCosts *[executor.OpcodeCount]uint32) {
	SetOpcodeCosts(opcodeCosts)
}

// NewInstanceWithOptions creates a new Wasmer instance from WASM bytecode,
// respecting the provided options
func (builder *WasmerExecutor) NewInstanceWithOptions(
	contractCode []byte,
	options executor.CompilationOptions,
) (executor.InstanceHandler, error) {
	return NewInstanceWithOptions(contractCode, options)
}

// NewInstanceFromCompiledCodeWithOptions creates a new Wasmer instance from
// precompiled machine code, respecting the provided options
func (builder *WasmerExecutor) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (executor.InstanceHandler, error) {
	return NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
}

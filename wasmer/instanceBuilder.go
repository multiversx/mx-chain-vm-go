package wasmer

import (
	"github.com/ElrondNetwork/wasm-vm/executor"
)

// WasmerInstanceBuilder is the default instance builder, which produces real
// Wasmer instances from WASM bytecode.
// TODO: rename to WasmerExecutor.
type WasmerInstanceBuilder struct {
}

// NewWasmerInstanceBuilder creates a new wasmer executor.
func NewExecutor() *WasmerInstanceBuilder {
	return &WasmerInstanceBuilder{}
}

// SetOpcodeCosts sets gas costs globally inside the Wasmer executor.
func (builder *WasmerInstanceBuilder) SetOpcodeCosts(opcodeCosts *[executor.OPCODE_COUNT]uint32) {
	SetOpcodeCosts(opcodeCosts)
}

// NewInstanceWithOptions creates a new Wasmer instance from WASM bytecode,
// respecting the provided options
func (builder *WasmerInstanceBuilder) NewInstanceWithOptions(
	contractCode []byte,
	options executor.CompilationOptions,
) (executor.InstanceHandler, error) {
	return NewInstanceWithOptions(contractCode, options)
}

// NewInstanceFromCompiledCodeWithOptions creates a new Wasmer instance from
// precompiled machine code, respecting the provided options
func (builder *WasmerInstanceBuilder) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (executor.InstanceHandler, error) {
	return NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
}

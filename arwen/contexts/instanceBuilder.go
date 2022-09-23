package contexts

import (
	"github.com/ElrondNetwork/wasm-vm/executor"
	"github.com/ElrondNetwork/wasm-vm/wasmer"
)

// WasmerInstanceBuilder is the default instance builder, which produces real
// Wasmer instances from WASM bytecode
type WasmerInstanceBuilder struct {
}

// NewInstanceWithOptions creates a new Wasmer instance from WASM bytecode,
// respecting the provided options
func (builder *WasmerInstanceBuilder) NewInstanceWithOptions(
	contractCode []byte,
	options executor.CompilationOptions,
) (executor.InstanceHandler, error) {
	return wasmer.NewInstanceWithOptions(contractCode, options)
}

// NewInstanceFromCompiledCodeWithOptions creates a new Wasmer instance from
// precompiled machine code, respecting the provided options
func (builder *WasmerInstanceBuilder) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (executor.InstanceHandler, error) {
	return wasmer.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
}

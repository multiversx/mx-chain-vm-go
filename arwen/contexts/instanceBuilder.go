package contexts

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/wasmer"
)

type WasmerInstanceBuilder struct {
}

// NewInstanceWithOptions creates a new Wasmer instance from WASM bytecode,
// respecting the provided options
func (builder *WasmerInstanceBuilder) NewInstanceWithOptions(
	contractCode []byte,
	options wasmer.CompilationOptions,
) (wasmer.InstanceHandler, error) {
	return wasmer.NewInstanceWithOptions(contractCode, options)
}

// NewInstanceFromCompiledCodeWithOptions creates a new Wasmer instance from
// precompiled machine code, respecting the provided options
func (builder *WasmerInstanceBuilder) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options wasmer.CompilationOptions,
) (wasmer.InstanceHandler, error) {
	return wasmer.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
}

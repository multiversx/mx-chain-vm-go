package contexts

import (
	"github.com/ElrondNetwork/wasm-vm/executor"
	"github.com/ElrondNetwork/wasm-vm/wasmer"
)

// WasmerInstanceBuilder is the default instance builder, which produces real
// Wasmer instances from WASM bytecode
type WasmerInstanceBuilder struct {
	hostPointer uintptr
}

func (builder *WasmerInstanceBuilder) SetContextData(hostPointer uintptr) {
	builder.hostPointer = hostPointer
}

// NewInstanceWithOptions creates a new Wasmer instance from WASM bytecode,
// respecting the provided options
func (builder *WasmerInstanceBuilder) NewInstanceWithOptions(
	contractCode []byte,
	options executor.CompilationOptions,
) (executor.InstanceHandler, error) {
	instance, err := wasmer.NewInstanceWithOptions(contractCode, options)
	if err != nil {
		return nil, err
	}
	instance.SetContextData(builder.hostPointer)
	return instance, nil
}

// NewInstanceFromCompiledCodeWithOptions creates a new Wasmer instance from
// precompiled machine code, respecting the provided options
func (builder *WasmerInstanceBuilder) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (executor.InstanceHandler, error) {
	instance, err := wasmer.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
	if err != nil {
		return nil, err
	}
	instance.SetContextData(builder.hostPointer)
	return instance, nil
}

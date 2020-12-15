package contexts

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
)

type wasmerInstanceBuilder struct {
}

func (builder *wasmerInstanceBuilder) NewInstanceWithOptions(
	contractCode []byte,
	options wasmer.CompilationOptions,
) (*wasmer.Instance, error) {
	return wasmer.NewInstanceWithOptions(contractCode, options)
}

func (builder *wasmerInstanceBuilder) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options wasmer.CompilationOptions,
) (*wasmer.Instance, error) {
	return wasmer.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
}

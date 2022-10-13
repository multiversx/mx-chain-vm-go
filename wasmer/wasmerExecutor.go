package wasmer

import (
	"github.com/ElrondNetwork/wasm-vm/executor"
)

// WasmerExecutor oversees the creation of Wasmer instances and execution.
type WasmerExecutor struct {
}

// NewExecutor creates a new wasmer executor.
func NewExecutor() (*WasmerExecutor, error) {
	return &WasmerExecutor{}, nil
}

// SetOpcodeCosts sets gas costs globally inside the Wasmer executor.
func (wasmerExecutor *WasmerExecutor) SetOpcodeCosts(opcodeCosts *[executor.OpcodeCount]uint32) {
	SetOpcodeCosts(opcodeCosts)
}

// SetRkyvSerializationEnabled controls a Wasmer flag.
func (wasmerExecutor *WasmerExecutor) SetRkyvSerializationEnabled(enabled bool) {
	SetRkyvSerializationEnabled(enabled)
}

// SetSIGSEGVPassthrough controls a Wasmer flag.
func (wasmerExecutor *WasmerExecutor) SetSIGSEGVPassthrough() {
	SetSIGSEGVPassthrough()
}

// NewInstanceWithOptions creates a new Wasmer instance from WASM bytecode,
// respecting the provided options
func (wasmerExecutor *WasmerExecutor) NewInstanceWithOptions(
	contractCode []byte,
	options executor.CompilationOptions,
) (executor.InstanceHandler, error) {
	return NewInstanceWithOptions(contractCode, options)
}

// NewInstanceFromCompiledCodeWithOptions creates a new Wasmer instance from
// precompiled machine code, respecting the provided options
func (wasmerExecutor *WasmerExecutor) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (executor.InstanceHandler, error) {
	return NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
}

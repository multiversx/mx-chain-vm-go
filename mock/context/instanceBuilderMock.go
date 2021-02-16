package mock

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
)

// InstanceBuilderMock can be passed to RuntimeContext as an InstanceBuilder to
// create mocked Wasmer instances.
type InstanceBuilderMock struct {
	InstanceMap map[string]wasmer.InstanceHandler
}

// NewInstanceBuilderMock constructs a new InstanceBuilderMock
func NewInstanceBuilderMock() *InstanceBuilderMock {
	return &InstanceBuilderMock{
		InstanceMap: make(map[string]wasmer.InstanceHandler),
	}
}

// CreateAndStoreInstanceMock creates a real Wasmer instance using the
// DefaultCode and returns it, so that a test may alter its Exports map to
// inject new contract methods; afterwards, the RuntimeContext will call
// NewInstanceWithOptions() and obtain the instance with the injected methods.
//
// It is necessary to call CreateAndStoreInstanceMock() for any contract that is
// to be called, or at least manually populate the InstanceMap appropriately
// (real WASM contracts may be used to populate it, as well).
func (builder *InstanceBuilderMock) CreateAndStoreInstanceMock(code []byte) *InstanceMock {
	instance := NewInstanceMock(code)
	builder.InstanceMap[string(code)] = instance

	return instance
}

// GetStoredInstanc retrieves and initializes a stored Wasmer instance, or
// nil if it doesn't exist
func (builder *InstanceBuilderMock) GetStoredInstance(code []byte, gasLimit uint64) (wasmer.InstanceHandler, bool) {
	instance, ok := builder.InstanceMap[string(code)]
	if ok {
		instance.SetPointsUsed(0)
		instance.SetGasLimit(gasLimit)
		return instance, true
	}
	return nil, false
}

// NewInstanceWithOptions attempts to load a prepared instance using
// GetStoredInstance; if it doesn't exist, it creates a true Wasmer
// instance with the provided contract code.
func (builder *InstanceBuilderMock) NewInstanceWithOptions(
	contractCode []byte,
	options wasmer.CompilationOptions,
) (wasmer.InstanceHandler, error) {

	instance, ok := builder.GetStoredInstance(contractCode, options.GasLimit)
	if ok {
		return instance, nil
	}
	return wasmer.NewInstanceWithOptions(contractCode, options)
}

// NewInstanceFromCompiledCodeWithOptions attempts to load a prepared instance
// using GetStoredInstance; if it doesn't exist, it creates a true Wasmer
// instance with the provided precompiled code.
func (builder *InstanceBuilderMock) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options wasmer.CompilationOptions,
) (wasmer.InstanceHandler, error) {
	instance, ok := builder.GetStoredInstance(compiledCode, options.GasLimit)
	if ok {
		return instance, nil
	}
	return wasmer.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
}

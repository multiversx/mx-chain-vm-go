package mock

import (
	"github.com/ElrondNetwork/wasm-vm/executor"
	"github.com/ElrondNetwork/wasm-vm/wasmer"
)

// InstanceBuilderRecorderMock can be passed to RuntimeContext as an InstanceBuilder to
// create mocked Wasmer instances.
type InstanceBuilderRecorderMock struct {
	InstanceMap map[string][]executor.InstanceHandler
}

// NewInstanceBuilderRecorderMock constructs a new InstanceBuilderRecorderMock
func NewInstanceBuilderRecorderMock() *InstanceBuilderRecorderMock {
	return &InstanceBuilderRecorderMock{
		InstanceMap: make(map[string][]executor.InstanceHandler),
	}
}

// NewInstanceWithOptions - see InstanceBuilderMock.NewInstanceWithOptions()
func (builder *InstanceBuilderRecorderMock) NewInstanceWithOptions(
	contractCode []byte,
	options executor.CompilationOptions,
) (executor.InstanceHandler, error) {
	instance, err := wasmer.NewInstanceWithOptions(contractCode, options)
	if err == nil {
		builder.addContractInstanceToInstanceMap(contractCode, instance)
	}
	return instance, err
}

// NewInstanceFromCompiledCodeWithOptions - see InstanceBuilderMock.NewInstanceFromCompiledCodeWithOptions()
func (builder *InstanceBuilderRecorderMock) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (executor.InstanceHandler, error) {
	instance, err := wasmer.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
	if err == nil {
		builder.addContractInstanceToInstanceMap(compiledCode, instance)
	}
	return instance, err
}

// add contract instance to the instance map for the given code
func (builder *InstanceBuilderRecorderMock) addContractInstanceToInstanceMap(code []byte, instance executor.InstanceHandler) {
	instances, ok := builder.InstanceMap[string(code)]
	if ok {
		instances = append(instances, instance)
	} else {
		instances = []executor.InstanceHandler{instance}
	}
	builder.InstanceMap[string(code)] = instances
}

// GetContractInstances gets contract instances for code
func (builder *InstanceBuilderRecorderMock) GetContractInstances(code []byte) []executor.InstanceHandler {
	return builder.InstanceMap[string(code)]
}

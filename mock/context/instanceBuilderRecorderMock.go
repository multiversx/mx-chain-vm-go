package mock

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
)

// InstanceBuilderRecorderMock can be passed to RuntimeContext as an InstanceBuilder to
// create mocked Wasmer instances.
type InstanceBuilderRecorderMock struct {
	InstanceMap map[string][]wasmer.InstanceHandler
}

// NewInstanceBuilderRecorderMock constructs a new InstanceBuilderRecorderMock
func NewInstanceBuilderRecorderMock() *InstanceBuilderRecorderMock {
	return &InstanceBuilderRecorderMock{
		InstanceMap: make(map[string][]wasmer.InstanceHandler),
	}
}

func (builder *InstanceBuilderRecorderMock) NewInstanceWithOptions(
	contractCode []byte,
	options wasmer.CompilationOptions,
) (wasmer.InstanceHandler, error) {
	instance, err := wasmer.NewInstanceWithOptions(contractCode, options)
	if err == nil {
		builder.addInstanceToInstanceMap(contractCode, instance)
	}
	return instance, err
}

func (builder *InstanceBuilderRecorderMock) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options wasmer.CompilationOptions,
) (wasmer.InstanceHandler, error) {
	instance, err := wasmer.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
	if err == nil {
		builder.addInstanceToInstanceMap(compiledCode, instance)
	}
	return instance, err
}

// add contract instance to the instance map
func (builder *InstanceBuilderRecorderMock) addInstanceToInstanceMap(code []byte, instance wasmer.InstanceHandler) {
	instances, ok := builder.InstanceMap[string(code)]
	if ok {
		instances = append(instances, instance)
	} else {
		instances = []wasmer.InstanceHandler{instance}
	}
	builder.InstanceMap[string(code)] = instances
}

func (builder *InstanceBuilderRecorderMock) GetInstances(code []byte) []wasmer.InstanceHandler {
	return builder.InstanceMap[string(code)]
}

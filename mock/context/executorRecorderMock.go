package mock

import (
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/executor"
	"github.com/ElrondNetwork/wasm-vm/wasmer"
)

// ExecutorRecorderMock can be passed to RuntimeContext as an InstanceBuilder to
// create mocked Wasmer instances.
type ExecutorRecorderMock struct {
	InstanceMap map[string][]executor.Instance
}

// NewExecutorRecorderMock constructs a new InstanceBuilderRecorderMock
func NewExecutorRecorderMock() *ExecutorRecorderMock {
	return &ExecutorRecorderMock{
		InstanceMap: make(map[string][]executor.Instance),
	}
}

func (executorMock *ExecutorRecorderMock) SetOpcodeCosts(opcodeCosts *[executor.OpcodeCount]uint32) {
	wasmer.SetOpcodeCosts(opcodeCosts)
}

// SetRkyvSerializationEnabled controls a Wasmer flag.
func (executorMock *ExecutorRecorderMock) SetRkyvSerializationEnabled(enabled bool) {
	wasmer.SetRkyvSerializationEnabled(enabled)
}

// SetSIGSEGVPassthrough controls a Wasmer flag.
func (executorMock *ExecutorRecorderMock) SetSIGSEGVPassthrough() {
	wasmer.SetSIGSEGVPassthrough()
}

func (executorMock *ExecutorRecorderMock) FunctionNames() vmcommon.FunctionNames {
	return nil
}

// NewInstanceWithOptions - see InstanceBuilderMock.NewInstanceWithOptions()
func (executorMock *ExecutorRecorderMock) NewInstanceWithOptions(
	contractCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	instance, err := wasmer.NewInstanceWithOptions(contractCode, options)
	if err == nil {
		executorMock.addContractInstanceToInstanceMap(contractCode, instance)
	}
	return instance, err
}

// NewInstanceFromCompiledCodeWithOptions - see InstanceBuilderMock.NewInstanceFromCompiledCodeWithOptions()
func (executorMock *ExecutorRecorderMock) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	instance, err := wasmer.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
	if err == nil {
		executorMock.addContractInstanceToInstanceMap(compiledCode, instance)
	}
	return instance, err
}

// SetVMHooks mocked method
func (executorMock *ExecutorRecorderMock) SetVMHooks(vmHooks executor.VMHooks) {
}

// GetVMHooks mocked method
func (executorMock *ExecutorRecorderMock) GetVMHooks() executor.VMHooks {
	return nil
}

// add contract instance to the instance map for the given code
func (executorMock *ExecutorRecorderMock) addContractInstanceToInstanceMap(code []byte, instance executor.Instance) {
	instances, ok := executorMock.InstanceMap[string(code)]
	if ok {
		instances = append(instances, instance)
	} else {
		instances = []executor.Instance{instance}
	}
	executorMock.InstanceMap[string(code)] = instances
}

// GetContractInstances gets contract instances for code
func (executorMock *ExecutorRecorderMock) GetContractInstances(code []byte) []executor.Instance {
	return executorMock.InstanceMap[string(code)]
}

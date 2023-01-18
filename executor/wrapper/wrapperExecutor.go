package executorwrapper

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/wasm-vm/executor"
)

// WrapperExecutor is a wrapper around an executor, who additionally caches all created instances.
// It also offers support for logging the executor operations. This logging is designed for testing and debugging.
// It is meant for testing the behavior of an executor.
type WrapperExecutor struct {
	logger          ExecutorLogger
	wrappedExecutor executor.Executor

	WrappedInstances map[string][]executor.Instance
}

// SetOpcodeCosts wraps the call to the underlying executor.
func (wexec *WrapperExecutor) SetOpcodeCosts(opcodeCosts *executor.WASMOpcodeCost) {
	wexec.wrappedExecutor.SetOpcodeCosts(opcodeCosts)
}

// FunctionNames wraps the call to the underlying executor.
func (wexec *WrapperExecutor) FunctionNames() vmcommon.FunctionNames {
	functionNames := wexec.wrappedExecutor.FunctionNames()
	return functionNames
}

// NewInstanceWithOptions wraps the call to the underlying executor.
func (wexec *WrapperExecutor) NewInstanceWithOptions(
	contractCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	wrappedInstance, err := wexec.wrappedExecutor.NewInstanceWithOptions(contractCode, options)
	if err != nil {
		return nil, err
	}
	wexec.addContractInstanceToInstanceMap(contractCode, wrappedInstance)
	wexec.logger.SetCurrentInstance(wrappedInstance)
	return &WrapperInstance{
		logger:          wexec.logger,
		wrappedInstance: wrappedInstance,
	}, nil
}

// NewInstanceFromCompiledCodeWithOptions wraps the call to the underlying executor.
func (wexec *WrapperExecutor) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	wrappedInstance, err := wexec.wrappedExecutor.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
	if err != nil {
		return nil, err
	}
	wexec.addContractInstanceToInstanceMap(compiledCode, wrappedInstance)
	wexec.logger.SetCurrentInstance(wrappedInstance)
	return &WrapperInstance{
		logger:          wexec.logger,
		wrappedInstance: wrappedInstance,
	}, nil
}

// add contract instance to the instance map for the given code
func (wexec *WrapperExecutor) addContractInstanceToInstanceMap(code []byte, instance executor.Instance) {
	instances, ok := wexec.WrappedInstances[string(code)]
	if ok {
		instances = append(instances, instance)
	} else {
		instances = []executor.Instance{instance}
	}
	wexec.WrappedInstances[string(code)] = instances
}

// GetContractInstances gets contract instances for code
func (wexec *WrapperExecutor) GetContractInstances(code []byte) []executor.Instance {
	return wexec.WrappedInstances[string(code)]
}

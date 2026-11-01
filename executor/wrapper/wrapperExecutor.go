package executorwrapper

import (
	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ executor.ExecutorNew = (*WrapperExecutor)(nil)

// WrapperExecutor is a wrapper around an executor, who additionally caches all created instances.
// It also offers support for logging the executor operations. This logging is designed for testing and debugging.
// It is meant for testing the behavior of an executor.
type WrapperExecutor struct {
	logger          ExecutorLogger
	wrappedExecutor executor.ExecutorNew

	WrappedInstances map[string][]executor.InstanceNew
}

// NewInstance wraps the call to the underlying executor.
func (wexec *WrapperExecutor) NewInstance(
	contractCode []byte,
	options *executor.CompilationOptionsNew,
) (executor.InstanceNew, error) {
	wrappedInstance, err := wexec.wrappedExecutor.NewInstance(contractCode, options)
	if err != nil {
		return nil, err
	}
	wexec.addContractInstanceToInstanceMap(contractCode, wrappedInstance)
	return &WrapperInstance{
		logger:          wexec.logger,
		wrappedInstance: wrappedInstance,
	}, nil
}

// NewInstanceFromCache wraps the call to the underlying executor.
func (wexec *WrapperExecutor) NewInstanceFromCache(
	compiledCode []byte,
	options *executor.CompilationOptionsNew,
) (executor.InstanceNew, error) {
	wrappedInstance, err := wexec.wrappedExecutor.NewInstanceFromCache(compiledCode, options)
	if err != nil {
		return nil, err
	}
	wexec.addContractInstanceToInstanceMap(compiledCode, wrappedInstance)
	return &WrapperInstance{
		logger:          wexec.logger,
		wrappedInstance: wrappedInstance,
	}, nil
}

// add contract instance to the instance map for the given code
func (wexec *WrapperExecutor) addContractInstanceToInstanceMap(code []byte, instance executor.InstanceNew) {
	instances, ok := wexec.WrappedInstances[string(code)]
	if ok {
		instances = append(instances, instance)
	} else {
		instances = []executor.InstanceNew{instance}
	}
	wexec.WrappedInstances[string(code)] = instances
}

// GetContractInstances gets contract instances for code
func (wexec *WrapperExecutor) GetContractInstances(code []byte) []executor.InstanceNew {
	return wexec.WrappedInstances[string(code)]
}

// IsInterfaceNil returns true if there is no value under the interface
func (wexec *WrapperExecutor) IsInterfaceNil() bool {
	return wexec == nil
}

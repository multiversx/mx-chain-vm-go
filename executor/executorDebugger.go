package executor

import (
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/config"
)

// ExecutorDebuggerFactory is the factory for the ExecutorDebugger.
type ExecutorDebuggerFactory struct {
	wrappedFactory ExecutorAbstractFactory

	// LastCreatedExecutor gives access to the created Executor
	LastCreatedExecutor *ExecutorDebugger
}

// ExecutorFactory returns the Wasmer executor factory.
func NewExecutorDebuggerFactory(wrappedFactory ExecutorAbstractFactory) *ExecutorDebuggerFactory {
	return &ExecutorDebuggerFactory{
		wrappedFactory: wrappedFactory,
	}
}

// CreateExecutor creates a new Executor instance.
func (ed *ExecutorDebuggerFactory) CreateExecutor(args ExecutorFactoryArgs) (Executor, error) {
	wrappedExecutor, err := ed.wrappedFactory.CreateExecutor(args)
	if err != nil {
		return nil, err
	}
	executorMock := &ExecutorDebugger{
		wrappedExecutor:  wrappedExecutor,
		CreatedInstances: make(map[string][]Instance),
	}
	ed.LastCreatedExecutor = executorMock
	return executorMock, nil
}

// ExecutorDebugger is a wrapper around an executor, who additionally caches all created instances.
// It is useful for testing the behavior of an executor.
type ExecutorDebugger struct {
	wrappedExecutor  Executor
	CreatedInstances map[string][]Instance
}

// SetOpcodeCosts wraps the call to the underlying executor.
func (executorMock *ExecutorDebugger) SetOpcodeCosts(opcodeCosts *config.WASMOpcodeCost) {
	executorMock.wrappedExecutor.SetOpcodeCosts(opcodeCosts)
}

// FunctionNames wraps the call to the underlying executor.
func (executorMock *ExecutorDebugger) FunctionNames() vmcommon.FunctionNames {
	return executorMock.wrappedExecutor.FunctionNames()
}

// NewInstanceWithOptions wraps the call to the underlying executor.
func (executorMock *ExecutorDebugger) NewInstanceWithOptions(
	contractCode []byte,
	options CompilationOptions,
) (Instance, error) {
	instance, err := executorMock.wrappedExecutor.NewInstanceWithOptions(contractCode, options)
	if err == nil {
		executorMock.addContractInstanceToInstanceMap(contractCode, instance)
	}
	return instance, err
}

// NewInstanceFromCompiledCodeWithOptions wraps the call to the underlying executor.
func (executorMock *ExecutorDebugger) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options CompilationOptions,
) (Instance, error) {
	instance, err := executorMock.wrappedExecutor.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
	if err == nil {
		executorMock.addContractInstanceToInstanceMap(compiledCode, instance)
	}
	return instance, err
}

// add contract instance to the instance map for the given code
func (executorMock *ExecutorDebugger) addContractInstanceToInstanceMap(code []byte, instance Instance) {
	instances, ok := executorMock.CreatedInstances[string(code)]
	if ok {
		instances = append(instances, instance)
	} else {
		instances = []Instance{instance}
	}
	executorMock.CreatedInstances[string(code)] = instances
}

// GetContractInstances gets contract instances for code
func (executorMock *ExecutorDebugger) GetContractInstances(code []byte) []Instance {
	return executorMock.CreatedInstances[string(code)]
}

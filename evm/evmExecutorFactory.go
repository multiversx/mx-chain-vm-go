package evm

import (
	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ = (executor.ExecutorAbstractFactory)((*EVMExecutorFactory)(nil))

// EVMExecutorFactory builds EVM Executors.
type EVMExecutorFactory struct{}

// ExecutorFactory returns the EVM executor factory.
func ExecutorFactory() *EVMExecutorFactory {
	return &EVMExecutorFactory{}
}

// CreateExecutor creates a new Executor instance.
func (eef *EVMExecutorFactory) CreateExecutor(args executor.ExecutorFactoryArgs) (executor.Executor, error) {
	return CreateExecutor(args)
}

// IsInterfaceNil returns true if there is no value under the interface
func (eef *EVMExecutorFactory) IsInterfaceNil() bool {
	return eef == nil
}

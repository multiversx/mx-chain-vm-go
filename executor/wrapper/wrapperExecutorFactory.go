package executorwrappers

import (
	"github.com/ElrondNetwork/wasm-vm/executor"
)

// WrapperExecutorFactory is the factory for the WrapperExecutor.
type WrapperExecutorFactory struct {
	logger         ExecutorLogger
	wrappedFactory executor.ExecutorAbstractFactory

	// LastCreatedExecutor gives access to the created Executor
	LastCreatedExecutor *WrapperExecutor
}

// NewWrappedExecutorFactory yields a new WrapperExecutor factory.
func NewWrappedExecutorFactory(
	logger ExecutorLogger,
	wrappedFactory executor.ExecutorAbstractFactory) *WrapperExecutorFactory {
	return &WrapperExecutorFactory{
		logger:         logger,
		wrappedFactory: wrappedFactory,
	}
}

// SimpleWrappedExecutorFactory yields a WrappedExecutor factory without logging.
func SimpleWrappedExecutorFactory(wrappedFactory executor.ExecutorAbstractFactory) *WrapperExecutorFactory {
	return NewWrappedExecutorFactory(&NoLogger{}, wrappedFactory)
}

// CreateExecutor creates a new Executor instance.
func (factory *WrapperExecutorFactory) CreateExecutor(args executor.ExecutorFactoryArgs) (executor.Executor, error) {
	wrappedExecutor, err := factory.wrappedFactory.CreateExecutor(args)
	if err != nil {
		return nil, err
	}
	factory.LastCreatedExecutor = &WrapperExecutor{
		logger:          factory.logger,
		wrappedExecutor: wrappedExecutor,

		WrappedInstances: make(map[string][]executor.Instance),
	}
	return factory.LastCreatedExecutor, nil
}

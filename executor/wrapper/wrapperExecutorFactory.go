package executorwrapper

import (
	"errors"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-vm-go/executor"
)

// ErrNilWrappedExecutorFactory signals that the wrapper was built around no factory at all.
var ErrNilWrappedExecutorFactory = errors.New("nil wrapped executor factory")

// ErrNilExecutorLogger signals that the wrapper was built without a logger. SimpleWrappedExecutorFactory
// supplies NoLogger for callers that want no logging; nil is not that.
var ErrNilExecutorLogger = errors.New("nil executor logger")

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
//
// The dependencies are checked here rather than in NewWrappedExecutorFactory, which has no
// error to return: the factory struct is inert, and this is the first point at which a
// missing one would be used.
func (factory *WrapperExecutorFactory) CreateExecutor(args executor.ExecutorFactoryArgs) (executor.Executor, error) {
	if check.IfNil(factory.wrappedFactory) {
		return nil, ErrNilWrappedExecutorFactory
	}
	if factory.logger == nil {
		return nil, ErrNilExecutorLogger
	}
	// The wrapper forwards HandleHookError, so it can only wrap hooks that have it.
	// Asserted here, at construction, rather than silently dropped at call time.
	vmHooks, ok := args.VMHooks.(vmHooksWithErrorHandler)
	if !ok {
		return nil, executor.ErrVMHooksCannotHandleErrors
	}

	wrappedExecutor, err := factory.wrappedFactory.CreateExecutor(executor.ExecutorFactoryArgs{
		VMHooks: &WrapperVMHooks{
			logger:         factory.logger,
			wrappedVMHooks: vmHooks,
		},
		OpcodeVersion:            args.OpcodeVersion,
		OpcodeCosts:              args.OpcodeCosts,
		RkyvSerializationEnabled: args.RkyvSerializationEnabled,
		WasmerSIGSEGVPassthrough: args.WasmerSIGSEGVPassthrough,
	})
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

// IsInterfaceNil returns true if there is no value under the interface
func (factory *WrapperExecutorFactory) IsInterfaceNil() bool {
	return factory == nil
}

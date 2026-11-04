package executorwrapper

import (
	"errors"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ executor.ExecutorNewAbstractFactory = (*WrapperExecutorFactory)(nil)

// ErrNilWrappedExecutorFactory signals that the wrapper was built around no factory at all.
var ErrNilWrappedExecutorFactory = errors.New("nil wrapped executor factory")

// ErrNilExecutorLogger signals that the wrapper was built without a logger. SimpleWrappedExecutorFactory
// supplies NoLogger for callers that want no logging; nil is not that.
var ErrNilExecutorLogger = errors.New("nil executor logger")

// WrapperExecutorFactory is the factory for the WrapperExecutor.
type WrapperExecutorFactory struct {
	logger         ExecutorLogger
	wrappedFactory executor.ExecutorNewAbstractFactory

	// LastCreatedExecutor gives access to the created Executor
	LastCreatedExecutor *WrapperExecutor
}

// NewWrappedExecutorFactory yields a new WrapperExecutor factory.
func NewWrappedExecutorFactory(
	logger ExecutorLogger,
	wrappedFactory executor.ExecutorNewAbstractFactory) *WrapperExecutorFactory {
	return &WrapperExecutorFactory{
		logger:         logger,
		wrappedFactory: wrappedFactory,
	}
}

// SimpleWrappedExecutorFactory yields a WrappedExecutor factory without logging.
func SimpleWrappedExecutorFactory(wrappedFactory executor.ExecutorNewAbstractFactory) *WrapperExecutorFactory {
	return NewWrappedExecutorFactory(&NoLogger{}, wrappedFactory)
}

// CreateExecutorNew creates a new ExecutorNew instance.
//
// The dependencies are checked here rather than in NewWrappedExecutorFactory, which has no
// error to return: the factory struct is inert, and this is the first point at which a
// missing one would be used.
//
// The hooks are no longer handed over here - they come from the runtime, per call - so what
// is wrapped is the runtime itself, and the hooks it hands out are wrapped on the way
// through.
func (factory *WrapperExecutorFactory) CreateExecutorNew(args executor.ExecutorFactoryArgsNew) (executor.ExecutorNew, error) {
	if check.IfNil(factory.wrappedFactory) {
		return nil, ErrNilWrappedExecutorFactory
	}
	if factory.logger == nil {
		return nil, ErrNilExecutorLogger
	}

	runtime := &WrapperRuntime{
		logger:         factory.logger,
		wrappedRuntime: args.Runtime,
	}

	// The wrapper forwards HandleHookError, so it can only wrap hooks that have it.
	// Asserted once here, at construction, rather than per call where there would be no way
	// to report it: VMHooksForCall cannot fail.
	err := runtime.checkHooksCanHandleErrors()
	if err != nil {
		return nil, err
	}

	wrappedExecutor, err := factory.wrappedFactory.CreateExecutorNew(executor.ExecutorFactoryArgsNew{
		Runtime: runtime,
	})
	if err != nil {
		return nil, err
	}

	factory.LastCreatedExecutor = &WrapperExecutor{
		logger:          factory.logger,
		wrappedExecutor: wrappedExecutor,

		WrappedInstances: make(map[string][]executor.InstanceNew),
	}
	return factory.LastCreatedExecutor, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (factory *WrapperExecutorFactory) IsInterfaceNil() bool {
	return factory == nil
}

// WrapperRuntime is the host behind a WrapperExecutor: it forwards to the real runtime and
// wraps the hooks it returns, so that the backend's hook calls are logged.
type WrapperRuntime struct {
	logger         ExecutorLogger
	wrappedRuntime executor.ExecutorRuntime
}

var _ executor.ExecutorRuntime = (*WrapperRuntime)(nil)

// VMHooksForCall returns the wrapped hooks servicing one host call.
func (runtime *WrapperRuntime) VMHooksForCall(instanceState executor.InstanceState) executor.VMHooks {
	vmHooks := runtime.wrappedRuntime.VMHooksForCall(instanceState)

	hooksWithHandler, ok := vmHooks.(vmHooksWithErrorHandler)
	if !ok {
		// Unreachable: the factory established this before the executor existed. Passing
		// the hooks through unwrapped loses the logging, which is better than losing the
		// call.
		return vmHooks
	}

	return &WrapperVMHooks{
		logger:         runtime.logger,
		wrappedVMHooks: hooksWithHandler,
	}
}

// checkHooksCanHandleErrors asks the wrapped runtime for one set of hooks and checks that
// they carry an error handler.
//
// A nil instance state is deliberate and safe: nothing is called on these hooks, they are
// only examined. Every host builds the same type of hooks for every call, so one is enough
// to answer the question for all of them.
func (runtime *WrapperRuntime) checkHooksCanHandleErrors() error {
	vmHooks := runtime.wrappedRuntime.VMHooksForCall(nil)
	if _, ok := vmHooks.(vmHooksWithErrorHandler); !ok {
		return ErrVMHooksCannotHandleErrors
	}

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (runtime *WrapperRuntime) IsInterfaceNil() bool {
	return runtime == nil
}

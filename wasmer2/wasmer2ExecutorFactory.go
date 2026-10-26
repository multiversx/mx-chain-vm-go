package wasmer2

import (
	"os/signal"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ = (executor.ExecutorAbstractFactory)((*Wasmer2ExecutorFactory)(nil))

// Wasmer2ExecutorFactory builds Wasmer2 Executors.
type Wasmer2ExecutorFactory struct{}

// ExecutorFactory returns the Wasmer executor factory.
func ExecutorFactory() *Wasmer2ExecutorFactory {
	return &Wasmer2ExecutorFactory{}
}

// CreateExecutor creates a new Executor instance.
func (wef *Wasmer2ExecutorFactory) CreateExecutor(args executor.ExecutorFactoryArgs) (executor.Executor, error) {
	signal.Reset()

	// This backend reaches the hooks as plain C functions, which cannot carry an error back,
	// so it runs them through the legacy adapter. Built here rather than handed in: the VM
	// and the backends moving to the new interfaces only deal in error-returning hooks.
	//
	// Absent hooks are tolerated, as they were before the adapter: tests that only compile
	// and instantiate a module never reach one.
	var legacyVMHooks executor.VMHooksLegacy
	if args.VMHooks != nil {
		errorHandler, ok := args.VMHooks.(executor.VMHooksErrorHandler)
		if !ok {
			return nil, ErrVMHooksCannotHandleErrors
		}
		legacyVMHooks = executor.NewVMHooksLegacyAdapter(args.VMHooks, errorHandler)
	}

	executor, err := CreateExecutor()
	if err != nil {
		return nil, err
	}
	executor.initVMHooks(legacyVMHooks)
	if args.OpcodeCosts != nil {
		// opcode costs are sometimes not initialized at this point in certain tests
		err = executor.SetOpcodeConfig(args.OpcodeVersion, args.OpcodeCosts)
		if err != nil {
			return nil, err
		}
	}

	return executor, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (wef *Wasmer2ExecutorFactory) IsInterfaceNil() bool {
	return wef == nil
}

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
//
// The opcode cost table is required. It used to be applied only when present, so that tests
// which had not built one could still create an executor - but an executor without a cost
// table does not fail, it meters every opcode at zero, and the contracts it runs are free.
// A caller with no schedule of its own passes config.MakeGasMapForTests.
func (wef *Wasmer2ExecutorFactory) CreateExecutor(args executor.ExecutorFactoryArgs) (executor.Executor, error) {
	if args.OpcodeCosts == nil {
		return nil, executor.ErrNilOpcodeCosts
	}

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
			return nil, executor.ErrVMHooksCannotHandleErrors
		}
		legacyVMHooks = executor.NewVMHooksLegacyAdapter(args.VMHooks, errorHandler)
	}

	wasmerExecutor, err := CreateExecutor()
	if err != nil {
		return nil, err
	}
	wasmerExecutor.initVMHooks(legacyVMHooks)
	err = wasmerExecutor.SetOpcodeConfig(args.OpcodeVersion, args.OpcodeCosts)
	if err != nil {
		// the native executor is not garbage collected, release it explicitly
		wasmerExecutor.destroy()
		return nil, err
	}

	return wasmerExecutor, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (wef *Wasmer2ExecutorFactory) IsInterfaceNil() bool {
	return wef == nil
}

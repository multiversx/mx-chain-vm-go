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

	wasmerExecutor, err := CreateExecutor()
	if err != nil {
		return nil, err
	}
	wasmerExecutor.initVMHooks(args.VMHooks)
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

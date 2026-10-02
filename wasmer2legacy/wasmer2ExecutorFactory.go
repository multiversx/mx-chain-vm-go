package wasmer2legacy

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

	executor, err := CreateExecutor()
	if err != nil {
		return nil, err
	}
	// The instance-lifetime variant, because this backend's hooks pointer is set once per
	// executor: a hook has no way to say which instance called it, so the provider resolves
	// that through the host's instance stack.
	executor.initVMHooks(args.VMHooksProvider.VMHooksForInstanceLifetime())
	if args.OpcodeCosts != nil {
		// opcode costs are sometimes not initialized at this point in certain tests
		executor.SetOpcodeConfig(args.OpcodeVersion, args.OpcodeCosts)
	}

	return executor, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (wef *Wasmer2ExecutorFactory) IsInterfaceNil() bool {
	return wef == nil
}

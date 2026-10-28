package wasmer2new

import (
	"os/signal"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ executor.ExecutorNewAbstractFactory = (*Wasmer2ExecutorFactory)(nil)

// Wasmer2ExecutorFactory builds Wasmer 2.2 executors reached through the new C API.
type Wasmer2ExecutorFactory struct{}

// ExecutorFactory returns the wasmer2new executor factory.
func ExecutorFactory() *Wasmer2ExecutorFactory {
	return &Wasmer2ExecutorFactory{}
}

// CreateExecutorNew creates a new ExecutorNew instance.
func (wef *Wasmer2ExecutorFactory) CreateExecutorNew(args executor.ExecutorFactoryArgsNew) (executor.ExecutorNew, error) {
	signal.Reset()

	return NewExecutor(args.Runtime)
}

// IsInterfaceNil returns true if there is no value under the interface
func (wef *Wasmer2ExecutorFactory) IsInterfaceNil() bool {
	return wef == nil
}

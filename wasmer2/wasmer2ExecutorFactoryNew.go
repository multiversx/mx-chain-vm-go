package wasmer2

import (
	"os/signal"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ executor.ExecutorNewAbstractFactory = (*Wasmer2ExecutorFactoryNew)(nil)

// Wasmer2ExecutorFactoryNew builds Wasmer 2.2 executors that speak the new interfaces.
type Wasmer2ExecutorFactoryNew struct{}

// ExecutorFactoryNew returns the Wasmer executor factory for the new interfaces.
func ExecutorFactoryNew() *Wasmer2ExecutorFactoryNew {
	return &Wasmer2ExecutorFactoryNew{}
}

// CreateExecutorNew creates a new ExecutorNew instance.
//
// There is nothing to configure afterwards: everything that is constant for the executor's
// life is in the args, and everything else travels per instance in the compilation options.
func (wef *Wasmer2ExecutorFactoryNew) CreateExecutorNew(args executor.ExecutorFactoryArgsNew) (executor.ExecutorNew, error) {
	signal.Reset()

	return CreateExecutorNew(args.Runtime)
}

// IsInterfaceNil returns true if there is no value under the interface
func (wef *Wasmer2ExecutorFactoryNew) IsInterfaceNil() bool {
	return wef == nil
}

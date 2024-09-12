package wasmer

import (
	"github.com/multiversx/mx-chain-vm-go/executor"
)

// WasmerExecutorFactory builds Wasmer Executors.
type WasmerExecutorFactory struct{}

// ExecutorFactory returns the Wasmer executor factory.
func ExecutorFactory() *WasmerExecutorFactory {
	return &WasmerExecutorFactory{}
}

// CreateExecutor creates a new Executor instance.
func (wef *WasmerExecutorFactory) CreateExecutor(args executor.ExecutorFactoryArgs) (executor.Executor, error) {
	if args.WasmerSIGSEGVPassthrough {
		SetSIGSEGVPassthrough()
	}

	exec, err := CreateExecutor()
	if err != nil {
		return nil, err
	}
	exec.initVMHooks(args.VMHooks)
	if args.OpcodeCosts.WASMOpcodeCost != nil {
		// opcode costs are sometimes not initialized at this point in certain tests
		exec.SetOpcodeCosts(args.OpcodeCosts)
	}

	SetRkyvSerializationEnabled(args.RkyvSerializationEnabled)
	return exec, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (wef *WasmerExecutorFactory) IsInterfaceNil() bool {
	return wef == nil
}

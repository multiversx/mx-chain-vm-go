package wasmer

import "github.com/ElrondNetwork/wasm-vm/executor"

// WasmerExecutorFactory builds Wasmer Executors.
type WasmerExecutorFactory struct{}

// ExecutorFactory returns the Wasmer executor factory.
func ExecutorFactory() *WasmerExecutorFactory {
	return &WasmerExecutorFactory{}
}

// NewExecutor creates a new Executor instance.
func (wef *WasmerExecutorFactory) NewExecutor(args executor.ExecutorFactoryArgs) (executor.Executor, error) {
	executor, err := NewExecutor()
	if err != nil {
		return nil, err
	}
	executor.initVMHooks(args.VMHooks)
	if args.OpcodeCosts != nil {
		// opcode costs are sometimes not initialized at this point in certain tests
		executor.SetOpcodeCosts(args.OpcodeCosts)
	}
	executor.SetRkyvSerializationEnabled(args.RkyvSerializationEnabled)
	if args.WasmerSIGSEGVPassthrough {
		executor.SetSIGSEGVPassthrough()
	}

	return executor, nil
}

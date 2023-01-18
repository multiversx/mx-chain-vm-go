package wasmer

import "github.com/multiversx/wasm-vm/executor"

// WasmerExecutorFactory builds Wasmer Executors.
type WasmerExecutorFactory struct{}

// ExecutorFactory returns the Wasmer executor factory.
func ExecutorFactory() *WasmerExecutorFactory {
	return &WasmerExecutorFactory{}
}

// CreateExecutor creates a new Executor instance.
func (wef *WasmerExecutorFactory) CreateExecutor(args executor.ExecutorFactoryArgs) (executor.Executor, error) {
	exec, err := CreateExecutor()
	if err != nil {
		return nil, err
	}
	exec.initVMHooks(args.VMHooks)
	if args.OpcodeCosts != nil {
		// opcode costs are sometimes not initialized at this point in certain tests
		exec.SetOpcodeCosts(args.OpcodeCosts)
	}
	SetRkyvSerializationEnabled(args.RkyvSerializationEnabled)
	if args.WasmerSIGSEGVPassthrough {
		SetSIGSEGVPassthrough()
	}

	return exec, nil
}

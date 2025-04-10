package executor

import "github.com/multiversx/mx-chain-core-go/core/check"

type VMOpcodeCost struct {
	EVMOpcodeCost  *EVMOpcodeCost
	WASMOpcodeCost *WASMOpcodeCost
}

// ExecutorFactoryArgs define the Executor configurations that come from the VM, especially the hooks and the gas costs.
type ExecutorFactoryArgs struct {
	EvmHooks                 EVMHooks
	VMHooks                  VMHooks
	OpcodeCosts              VMOpcodeCost
	RkyvSerializationEnabled bool
	WasmerSIGSEGVPassthrough bool
}

// ExecutorAbstractFactory defines an object to be passed to the VM to configure the instantiation of the Executor.
// The VM needs to create the Executor, because the VM hooks, gas costs and other configurations come from it.
type ExecutorAbstractFactory interface {
	check.NilInterfaceChecker

	// CreateExecutor produces a new Executor instance.
	CreateExecutor(args ExecutorFactoryArgs) (Executor, error)
}

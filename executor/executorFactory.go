package executor

// ExecutorFactoryArgs define the Executor configurations that come from the VM, especially the hooks and the gas costs.
type ExecutorFactoryArgs struct {
	VMHooks                  VMHooks
	OpcodeCosts              *WASMOpcodeCost
	RkyvSerializationEnabled bool
	WasmerSIGSEGVPassthrough bool
}

// ExecutorFactory defines an object to be passed to the VM to configure the instantiation of the Executor.
// The VM needs to create the Executor, because the VM hooks, gas costs and other configurations come from it.
type ExecutorFactory interface {
	// NewExecutor creates a new Executor instance.
	NewExecutor(args ExecutorFactoryArgs) (Executor, error)
}
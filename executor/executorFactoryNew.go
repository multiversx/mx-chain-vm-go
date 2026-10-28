package executor

import "github.com/multiversx/mx-chain-core-go/core/check"

// ExecutorFactoryArgsNew is everything an [ExecutorNew] is constructed with.
//
// One struct rather than a parameter list, so that an executor is configured directly and at
// once: there is nothing to set afterwards, and therefore no window in which an executor
// exists but is not yet usable.
//
// It holds one field, and that is the honest answer rather than an oversight. Under the new
// interfaces, everything [ExecutorFactoryArgs] carries besides the hooks — the opcode
// version, the gas costs — is compilation configuration, and travels per instance in
// [CompilationOptionsNew]. What is left is the one thing that really is constant for an
// executor's whole life: the host behind it.
type ExecutorFactoryArgsNew struct {
	// Runtime is the host behind the executor: supplied once, never replaced.
	Runtime ExecutorRuntime
}

// ExecutorNewAbstractFactory is an object passed to the VM to configure the instantiation of
// the executor. The VM needs to create the executor, because the VM hooks come from it.
//
// Counterpart of the Rust `ExecutorFactory` trait, which lives in the separate
// executor-factory crate rather than in new_traits. Here it stays next to the interfaces, as
// [ExecutorAbstractFactory] always has: in the Go VM the factory is how a node selects its
// backend, not merely how the tests do.
type ExecutorNewAbstractFactory interface {
	check.NilInterfaceChecker

	// CreateExecutorNew produces a new executor.
	//
	// New suffix: ExecutorAbstractFactory.CreateExecutor returns the old [Executor].
	CreateExecutorNew(args ExecutorFactoryArgsNew) (ExecutorNew, error)
}

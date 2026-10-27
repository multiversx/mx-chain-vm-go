package executor

import "github.com/multiversx/mx-chain-core-go/core/check"

// ExecutorNew creates instances, and nothing else.
//
// It is the Go counterpart of the Rust `Executor` trait (new_traits/executor_new.rs), and
// replaces [Executor] as the migration progresses. The old name is kept on the old interface
// for now; renaming is deliberately postponed.
//
// An executor is a factory and an instance-type selector. It holds only what is constant
// across every instance it will ever make: the [ExecutorRuntime] behind it, and whatever
// plumbing the backend itself needs, such as logging. Everything that is execution logic —
// the opcode costs, the memory limits — arrives per instance, in [CompilationOptionsNew].
//
// There are no setters. A host that needs different compilation options passes different
// options; a host that needs a different runtime builds a different executor. This is what
// SetOpcodeConfig and FunctionNames on [Executor] turn into: the first moves into the
// compilation options, the second is the host's own business.
type ExecutorNew interface {
	check.NilInterfaceChecker

	// NewInstance creates a new executor instance.
	NewInstance(
		wasmBytes []byte,
		compilationOptions *CompilationOptionsNew) (InstanceNew, error)

	// NewInstanceFromCache restores an executor instance from cache.
	NewInstanceFromCache(
		cacheBytes []byte,
		compilationOptions *CompilationOptionsNew) (InstanceNew, error)
}

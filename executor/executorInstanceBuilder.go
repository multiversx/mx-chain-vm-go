package executor

// CompilationOptions contains configurations for instantiating an executor instance.
type CompilationOptions struct {
	GasLimit           uint64
	UnmeteredLocals    uint64
	MaxMemoryGrow      uint64
	MaxMemoryGrowDelta uint64
	OpcodeTrace        bool
	Metering           bool
	RuntimeBreakpoints bool
}

// InstanceBuilder defines the functionality needed to create any executor instance.
// TODO: rename to Executor or VMExecutor.
type InstanceBuilder interface {
	NewInstanceWithOptions(
		contractCode []byte,
		options CompilationOptions) (InstanceHandler, error)

	NewInstanceFromCompiledCodeWithOptions(
		compiledCode []byte,
		options CompilationOptions) (InstanceHandler, error)
}

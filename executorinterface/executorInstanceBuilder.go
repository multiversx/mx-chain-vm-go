package executorinterface

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
type InstanceBuilder interface {
	NewInstanceWithOptions(
		contractCode []byte,
		options CompilationOptions) (InstanceHandler, error)

	NewInstanceFromCompiledCodeWithOptions(
		compiledCode []byte,
		options CompilationOptions) (InstanceHandler, error)
}

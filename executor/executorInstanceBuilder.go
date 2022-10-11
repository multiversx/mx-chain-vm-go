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

// OpcodeCount is the number of opcodes that we account for when setting gas costs.
const OpcodeCount = 448

// InstanceBuilder defines the functionality needed to create any executor instance.
// TODO: rename to Executor or VMExecutor.
type InstanceBuilder interface {
	// SetOpcodeCosts sets gas costs globally inside an executor.
	SetOpcodeCosts(opcodeCosts *[OpcodeCount]uint32)
	// NewInstanceWithOptions creates a new executor instance.
	NewInstanceWithOptions(
		contractCode []byte,
		options CompilationOptions) (InstanceHandler, error)

	// NewInstanceWithOptions is used to restore an executor instance from cache.
	NewInstanceFromCompiledCodeWithOptions(
		compiledCode []byte,
		options CompilationOptions) (InstanceHandler, error)
}

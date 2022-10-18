package executor

import vmcommon "github.com/ElrondNetwork/elrond-vm-common"

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

// Executor defines the functionality needed to create any executor instance.
type Executor interface {
	// SetOpcodeCosts sets gas costs globally inside an executor.
	SetOpcodeCosts(opcodeCosts *[OpcodeCount]uint32)

	// SetRkyvSerializationEnabled controls a Wasmer flag.
	SetRkyvSerializationEnabled(enabled bool)

	// SetSIGSEGVPassthrough controls a Wasmer flag.
	SetSIGSEGVPassthrough()

	// FunctionNames return the low-level function names provided to contracts.
	FunctionNames() vmcommon.FunctionNames

	// NewInstanceWithOptions creates a new executor instance.
	NewInstanceWithOptions(
		contractCode []byte,
		options CompilationOptions) (InstanceHandler, error)

	// NewInstanceWithOptions is used to restore an executor instance from cache.
	NewInstanceFromCompiledCodeWithOptions(
		compiledCode []byte,
		options CompilationOptions) (InstanceHandler, error)
}

package executor

import (
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/config"
)

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

// Executor defines the functionality needed to create any executor instance.
type Executor interface {
	// SetOpcodeCosts sets gas costs globally inside an executor.
	SetOpcodeCosts(opcodeCosts *config.WASMOpcodeCost)

	// SetRkyvSerializationEnabled controls a Wasmer flag.
	SetRkyvSerializationEnabled(enabled bool)

	// SetSIGSEGVPassthrough controls a Wasmer flag.
	SetSIGSEGVPassthrough()

	// FunctionNames return the low-level function names provided to contracts.
	FunctionNames() vmcommon.FunctionNames

	// NewInstanceWithOptions creates a new executor instance.
	NewInstanceWithOptions(
		contractCode []byte,
		options CompilationOptions) (Instance, error)

	// NewInstanceFromCompiledCodeWithOptions is used to restore an executor instance from cache.
	NewInstanceFromCompiledCodeWithOptions(
		compiledCode []byte,
		options CompilationOptions) (Instance, error)

	// GetVMHooks returns the VM hooks.
	GetVMHooks() VMHooks

	// InitVMHooks inits the VM hooks.
	InitVMHooks(vmHooks VMHooks)
}

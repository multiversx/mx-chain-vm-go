// Package executor contains the interfaces and definitions for the VM Executor
package executor

import (
	"github.com/multiversx/mx-chain-core-go/core/check"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
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
	check.NilInterfaceChecker

	// SetOpcodeCosts sets gas costs globally inside an executor.
	SetOpcodeCosts(opcodeCosts *WASMOpcodeCost)

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
}

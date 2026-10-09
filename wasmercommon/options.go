package wasmercommon

import "github.com/multiversx/mx-chain-vm-go/executor"

// ExecutorConfig is everything an executor is configured with, in one struct.
//
// The legacy API configured an executor by a sequence of setters after construction
// (SetOpcodeConfig, SetRkyvSerializationEnabled, SetSIGSEGVPassthrough). The new one takes its
// configuration once, which is both why there is nothing to mutate afterwards and why there
// is no window in which an executor exists but is not yet usable.
type ExecutorConfig struct {
	OpcodeVersion  executor.OpcodeVersion
	WASMOpcodeCost *executor.WASMOpcodeCost
}

// CompilationOptions is the per-instance half of the configuration.
//
// Four fields against the legacy seven. GasLimit moved to the call, where it belongs — it
// changes on every invocation while the compiled module does not — and metering and runtime
// breakpoints are no longer optional, so there is nothing to switch off.
//
// How this is reached from executor.CompilationOptions is a question for the interface
// convergence, not for the boundary; see PLAN-go-ffi.md, Phase 4b.
type CompilationOptions struct {
	UnmeteredLocals    uint64
	MaxMemoryGrow      uint64
	MaxMemoryGrowDelta uint64
	OpcodeTrace        bool
}

// ToCompilationOptions keeps the four options these executors have and drops the three they do
// not, rather than reinterpreting the struct: GasLimit moved to the call, and metering and
// runtime breakpoints are no longer optional.
func ToCompilationOptions(options executor.CompilationOptions) CompilationOptions {
	return CompilationOptions{
		UnmeteredLocals:    options.UnmeteredLocals,
		MaxMemoryGrow:      options.MaxMemoryGrow,
		MaxMemoryGrowDelta: options.MaxMemoryGrowDelta,
		OpcodeTrace:        options.OpcodeTrace,
	}
}

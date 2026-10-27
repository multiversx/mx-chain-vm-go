package executor

// OpcodeConfig is the opcode set and the per-opcode gas costs a module is compiled with.
//
// Counterpart of the Rust `OpcodeConfig`. It used to sit on the executor, behind
// Executor.SetOpcodeConfig; it is compilation configuration, so it travels with the
// compilation options instead.
type OpcodeConfig struct {
	OpcodeVersion OpcodeVersion
	OpcodeCosts   *WASMOpcodeCost
}

// CompilationOptionsNew is everything that determines how a wasm module is compiled.
//
// Counterpart of the Rust `CompilationOptions` (new_traits/compilation_options_new.rs).
// Named with a New suffix only because [CompilationOptions] is still in use; the rename is
// postponed.
//
// Passed to [ExecutorNew.NewInstance] on every instance, rather than configured on the
// executor: these are essential execution logic, identical for every instance the host cares
// about, and a compiled module is only interchangeable with another compiled from the same
// options. Keeping them here means the pair (wasm, options) fully determines the result,
// which is what lets a module cache be keyed on it.
//
// That includes the opcode costs, which used to sit on the executor behind a mutable setter.
// A gas schedule change is therefore a new set of options — and, for a host that caches, a
// different cache key — rather than a reconfigured executor.
//
// Gone, relative to [CompilationOptions]: GasLimit, which is now an argument of
// InstanceNew.Call, and Metering/RuntimeBreakpoints, which were never anything but true in
// production.
type CompilationOptionsNew struct {
	UnmeteredLocals    uint64
	MaxMemoryGrow      uint64
	MaxMemoryGrowDelta uint64
	OpcodeTrace        bool

	// OpcodeConfig is shared rather than owned: the cost table is some hundreds of bytes and
	// every instance under one gas schedule uses the same one, so the host builds it once.
	OpcodeConfig *OpcodeConfig
}

// ToLegacy converts to the old compilation options, for a backend that is still reached
// through the old interface.
func (options *CompilationOptionsNew) ToLegacy() CompilationOptions {
	return CompilationOptions{
		GasLimit:           0,
		UnmeteredLocals:    options.UnmeteredLocals,
		MaxMemoryGrow:      options.MaxMemoryGrow,
		MaxMemoryGrowDelta: options.MaxMemoryGrowDelta,
		OpcodeTrace:        options.OpcodeTrace,
		Metering:           true,
		RuntimeBreakpoints: true,
	}
}

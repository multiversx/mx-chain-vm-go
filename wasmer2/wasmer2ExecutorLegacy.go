package wasmer2

// Transitional: executor.Executor as it is today, plus the factory that lets a test select
// this backend. See the note at the top of wasmer2InstanceLegacy.go.

import (
	"os/signal"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ = (executor.Executor)((*Wasmer2Executor)(nil))
var _ = (executor.ExecutorAbstractFactory)((*Wasmer2ExecutorFactory)(nil))

// Wasmer2ExecutorFactory builds Wasmer 2.2 executors.
type Wasmer2ExecutorFactory struct{}

// ExecutorFactory returns the Wasmer 2.2 executor factory.
func ExecutorFactory() *Wasmer2ExecutorFactory {
	return &Wasmer2ExecutorFactory{}
}

// CreateExecutor creates a new executor.
func (wef *Wasmer2ExecutorFactory) CreateExecutor(args executor.ExecutorFactoryArgs) (executor.Executor, error) {
	// Same as the legacy factory does, and for the same reason: the runtime installs its own
	// trap handlers, and Go's must not be left in front of them.
	signal.Reset()

	opcodeCosts := args.OpcodeCosts
	if opcodeCosts == nil {
		// The legacy factory tolerates this, because some tests reach it before the gas
		// schedule is in place. Metering is not optional in this executor, so an absent cost
		// table has to become an explicit one rather than a null pointer.
		opcodeCosts = &executor.WASMOpcodeCost{}
	}

	wasmerExecutor, err := NewExecutor(ExecutorConfig{
		OpcodeVersion:  args.OpcodeVersion,
		WASMOpcodeCost: opcodeCosts,
	})
	if err != nil {
		return nil, err
	}

	wasmerExecutor.provider = args.VMHooksProvider

	return wasmerExecutor, nil
}

// IsInterfaceNil returns true if there is no value under the interface.
func (wef *Wasmer2ExecutorFactory) IsInterfaceNil() bool {
	return wef == nil
}

// SetOpcodeConfig applies a new opcode version and cost table by rebuilding the executor.
//
// This API takes its configuration once, at creation, and offers no setter — so the only way
// to change it is to build another one. That is behaviour-preserving for the single caller,
// GasScheduleChange, which already clears the warm instance cache: no instance survives the
// change, so nothing observes the swap. See PLAN-go-ffi.md, gate G4.
func (wasmerExecutor *Wasmer2Executor) SetOpcodeConfig(
	opcodeVersion executor.OpcodeVersion,
	wasmOps *executor.WASMOpcodeCost,
) {
	rebuilt, err := NewExecutor(ExecutorConfig{
		OpcodeVersion:  opcodeVersion,
		WASMOpcodeCost: wasmOps,
	})
	if err != nil {
		logWasmer2.Error("could not apply new opcode config", "err", err)
		return
	}

	provider := wasmerExecutor.provider
	wasmerExecutor.Close()

	wasmerExecutor.cExecutor = rebuilt.cExecutor
	wasmerExecutor.config = rebuilt.config
	wasmerExecutor.provider = provider
}

// FunctionNames returns the low-level function names provided to contracts.
func (wasmerExecutor *Wasmer2Executor) FunctionNames() vmcommon.FunctionNames {
	return executor.VMHookNames
}

// NewInstanceWithOptions creates a new instance from WASM bytecode.
func (wasmerExecutor *Wasmer2Executor) NewInstanceWithOptions(
	contractCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	instance, err := wasmerExecutor.NewInstance(
		contractCode,
		toCompilationOptions(options),
		wasmerExecutor.provider,
	)
	if err != nil {
		return nil, err
	}

	// GasLimit is a compilation option in the legacy interface and a call argument in this
	// one, so it is held on the instance until the call. Metering and RuntimeBreakpoints are
	// dropped: they are always on here, and were never switched off in practice.
	instance.SetGasLimit(options.GasLimit)

	return instance, nil
}

// NewInstanceFromCompiledCodeWithOptions is not supported yet; see PLAN-go-ffi.md, Phase 6.
func (wasmerExecutor *Wasmer2Executor) NewInstanceFromCompiledCodeWithOptions(
	_ []byte,
	_ executor.CompilationOptions,
) (executor.Instance, error) {
	return nil, ErrCacheUnsupported
}

// IsInterfaceNil returns true if there is no value under the interface.
func (wasmerExecutor *Wasmer2Executor) IsInterfaceNil() bool {
	return wasmerExecutor == nil
}

// toCompilationOptions keeps the four options this executor has and drops the three it does
// not, rather than reinterpreting the struct: GasLimit moved to the call, and metering and
// runtime breakpoints are no longer optional.
func toCompilationOptions(options executor.CompilationOptions) CompilationOptions {
	return CompilationOptions{
		UnmeteredLocals:    options.UnmeteredLocals,
		MaxMemoryGrow:      options.MaxMemoryGrow,
		MaxMemoryGrowDelta: options.MaxMemoryGrowDelta,
		OpcodeTrace:        options.OpcodeTrace,
	}
}

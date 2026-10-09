package wasmer6

// Transitional: executor.Executor as it is today, plus the factory that lets a test select
// this backend. See the note at the top of wasmer6InstanceLegacy.go.

import (
	"os/signal"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/wasmercommon"
)

var _ = (executor.Executor)((*Wasmer6Executor)(nil))
var _ = (executor.ExecutorAbstractFactory)((*Wasmer6ExecutorFactory)(nil))

// Wasmer6ExecutorFactory builds Wasmer 6 executors.
type Wasmer6ExecutorFactory struct{}

// ExecutorFactory returns the Wasmer 6 executor factory.
func ExecutorFactory() *Wasmer6ExecutorFactory {
	return &Wasmer6ExecutorFactory{}
}

// CreateExecutor creates a new executor.
func (wef *Wasmer6ExecutorFactory) CreateExecutor(args executor.ExecutorFactoryArgs) (executor.Executor, error) {
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

	wasmerExecutor, err := NewExecutor(wasmercommon.ExecutorConfig{
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
func (wef *Wasmer6ExecutorFactory) IsInterfaceNil() bool {
	return wef == nil
}

// SetOpcodeConfig applies a new opcode version and cost table by rebuilding the executor.
//
// This API takes its configuration once, at creation, and offers no setter — so the only way
// to change it is to build another one. That is behaviour-preserving for the single caller,
// GasScheduleChange, which already clears the warm instance cache: no instance survives the
// change, so nothing observes the swap. See PLAN-go-ffi.md, gate G4.
func (wasmerExecutor *Wasmer6Executor) SetOpcodeConfig(
	opcodeVersion executor.OpcodeVersion,
	wasmOps *executor.WASMOpcodeCost,
) {
	rebuilt, err := NewExecutor(wasmercommon.ExecutorConfig{
		OpcodeVersion:  opcodeVersion,
		WASMOpcodeCost: wasmOps,
	})
	if err != nil {
		logWasmer6.Error("could not apply new opcode config", "err", err)
		return
	}

	provider := wasmerExecutor.provider
	wasmerExecutor.Close()

	wasmerExecutor.cExecutor = rebuilt.cExecutor
	wasmerExecutor.config = rebuilt.config
	wasmerExecutor.provider = provider
}

// FunctionNames returns the low-level function names provided to contracts.
func (wasmerExecutor *Wasmer6Executor) FunctionNames() vmcommon.FunctionNames {
	return executor.VMHookNames
}

// NewInstanceWithOptions creates a new instance from WASM bytecode.
func (wasmerExecutor *Wasmer6Executor) NewInstanceWithOptions(
	contractCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	instance, err := wasmerExecutor.NewInstance(
		contractCode,
		wasmercommon.ToCompilationOptions(options),
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
func (wasmerExecutor *Wasmer6Executor) NewInstanceFromCompiledCodeWithOptions(
	_ []byte,
	_ executor.CompilationOptions,
) (executor.Instance, error) {
	return nil, wasmercommon.ErrCacheUnsupported
}

// IsInterfaceNil returns true if there is no value under the interface.
func (wasmerExecutor *Wasmer6Executor) IsInterfaceNil() bool {
	return wasmerExecutor == nil
}

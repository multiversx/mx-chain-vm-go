package wasmer7

// Transitional: executor.Executor as it is today, plus the factory that lets a test select
// this backend. See the note at the top of wasmer7InstanceLegacy.go.

import (
	"os/signal"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/wasmercommon"
)

var _ = (executor.Executor)((*Wasmer7Executor)(nil))
var _ = (executor.ExecutorAbstractFactory)((*Wasmer7ExecutorFactory)(nil))

// Wasmer7ExecutorFactory builds Wasmer 7 executors.
type Wasmer7ExecutorFactory struct{}

// ExecutorFactory returns the Wasmer 7 executor factory.
func ExecutorFactory() *Wasmer7ExecutorFactory {
	return &Wasmer7ExecutorFactory{}
}

// CreateExecutor creates a new executor.
func (wef *Wasmer7ExecutorFactory) CreateExecutor(args executor.ExecutorFactoryArgs) (executor.Executor, error) {
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
func (wef *Wasmer7ExecutorFactory) IsInterfaceNil() bool {
	return wef == nil
}

// SetOpcodeConfig applies a new opcode version and cost table by rebuilding the executor.
//
// This API takes its configuration once, at creation, and offers no setter — so the only way
// to change it is to build another one. That is behaviour-preserving for the single caller,
// GasScheduleChange, which already clears the warm instance cache: no instance survives the
// change, so nothing observes the swap. See PLAN-go-ffi.md, gate G4.
func (wasmerExecutor *Wasmer7Executor) SetOpcodeConfig(
	opcodeVersion executor.OpcodeVersion,
	wasmOps *executor.WASMOpcodeCost,
) {
	rebuilt, err := NewExecutor(wasmercommon.ExecutorConfig{
		OpcodeVersion:  opcodeVersion,
		WASMOpcodeCost: wasmOps,
	})
	if err != nil {
		logWasmer7.Error("could not apply new opcode config", "err", err)
		return
	}

	provider := wasmerExecutor.provider
	wasmerExecutor.Close()

	wasmerExecutor.cExecutor = rebuilt.cExecutor
	wasmerExecutor.config = rebuilt.config
	wasmerExecutor.provider = provider
}

// FunctionNames returns the low-level function names provided to contracts.
func (wasmerExecutor *Wasmer7Executor) FunctionNames() vmcommon.FunctionNames {
	return executor.VMHookNames
}

// NewInstanceWithOptions creates a new instance from WASM bytecode.
func (wasmerExecutor *Wasmer7Executor) NewInstanceWithOptions(
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

// NewInstanceFromCompiledCodeWithOptions creates an instance from a previously cached
// compiled module, skipping compilation.
//
// compiledCode comes from the VM's own compiled-code store, where Instance.Cache put it. It
// is loaded as executable code rather than validated, so this must not be reached with bytes
// from anywhere else; the caller is runtimeContext.makeInstanceFromCompiledCode, which reads
// them back by code hash.
//
// A caller that gets an error here is expected to fall back to compiling the bytecode, which
// is what makes an artifact from an older Wasmer version — rejected on its ABI version — a
// recoverable miss rather than a failed execution.
func (wasmerExecutor *Wasmer7Executor) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	instance, err := wasmerExecutor.NewInstanceFromCache(
		compiledCode,
		wasmercommon.ToCompilationOptions(options),
		wasmerExecutor.provider,
	)
	if err != nil {
		return nil, err
	}

	// As in NewInstanceWithOptions: the gas limit is a compilation option in the legacy
	// interface and a call argument in this one.
	instance.SetGasLimit(options.GasLimit)

	return instance, nil
}

// IsInterfaceNil returns true if there is no value under the interface.
func (wasmerExecutor *Wasmer7Executor) IsInterfaceNil() bool {
	return wasmerExecutor == nil
}

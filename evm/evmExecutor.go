package evm

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	evm "github.com/multiversx/mx-chain-vm-go/evm/interpreter"
	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ executor.Executor = (*EVMExecutor)(nil)

// EVMExecutor oversees the creation of EVM instances and execution.
type EVMExecutor struct {
	evmHooks       executor.EVMHooks
	gasConfig      *evm.GasConfig
	instructionSet evm.JumpTable
}

// CreateExecutor creates a new EVM executor.
func CreateExecutor(args executor.ExecutorFactoryArgs) (*EVMExecutor, error) {
	evmExecutor := &EVMExecutor{evmHooks: args.EvmHooks}
	evmExecutor.SetOpcodeCosts(args.OpcodeCosts)
	return evmExecutor, nil
}

func (evmExecutor *EVMExecutor) SetOpcodeCosts(opcodeCost executor.VMOpcodeCost) {
	if opcodeCost.EVMOpcodeCost != nil {
		evmExecutor.gasConfig = extractOpcodeCost(opcodeCost.EVMOpcodeCost)
		evmExecutor.instructionSet = evm.NewCancunInstructionSet(evmExecutor.gasConfig)
	}
}

// HasFunctionNameChecks returns true if the instance requires function name checks.
func (evmExecutor *EVMExecutor) HasFunctionNameChecks() bool {
	return false
}

func (evmExecutor *EVMExecutor) FunctionNames() vmcommon.FunctionNames {
	return map[string]struct{}{}
}

// NewInstanceWithOptions creates a new EVM instance from EVM bytecode,
// respecting the provided options
func (evmExecutor *EVMExecutor) NewInstanceWithOptions(
	contractCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	return newInstance(evmExecutor, false, contractCode, options)
}

// NewInstanceFromCompiledCodeWithOptions creates a new EVM instance from compiled code,
// respecting the provided options
func (evmExecutor *EVMExecutor) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	return newInstance(evmExecutor, true, compiledCode, options)
}

// IsInterfaceNil returns true if underlying object is nil
func (evmExecutor *EVMExecutor) IsInterfaceNil() bool {
	return evmExecutor == nil
}

func extractOpcodeCost(opcodeCost *executor.EVMOpcodeCost) *evm.GasConfig {
	return &evm.GasConfig{
		QuickStep:             opcodeCost.QuickStep,
		FastestStep:           opcodeCost.FastestStep,
		FastStep:              opcodeCost.FastStep,
		MidStep:               opcodeCost.MidStep,
		SlowStep:              opcodeCost.SlowStep,
		ExtStep:               opcodeCost.ExtStep,
		Ecrecover:             opcodeCost.Ecrecover,
		Sha256PerWord:         opcodeCost.Sha256PerWord,
		Sha256Base:            opcodeCost.Sha256Base,
		Ripemd160PerWord:      opcodeCost.Ripemd160PerWord,
		Ripemd160Base:         opcodeCost.Ripemd160Base,
		IdentityPerWord:       opcodeCost.IdentityPerWord,
		IdentityBase:          opcodeCost.IdentityBase,
		Bn256Add:              opcodeCost.Bn256Add,
		Bn256ScalarMul:        opcodeCost.Bn256ScalarMul,
		Bn256PairingBase:      opcodeCost.Bn256PairingBase,
		Bn256PairingPerPoint:  opcodeCost.Bn256PairingPerPoint,
		BlobTxPointEvaluation: opcodeCost.BlobTxPointEvaluation,
		Keccak256:             opcodeCost.Keccak256,
		Balance:               opcodeCost.Balance,
		ExtcodeSize:           opcodeCost.ExtcodeSize,
		ExtcodeCopy:           opcodeCost.ExtcodeCopy,
		ExtcodeHash:           opcodeCost.ExtcodeHash,
		Sload:                 opcodeCost.Sload,
		Sstore:                opcodeCost.Sstore,
		Jumpdest:              opcodeCost.Jumpdest,
		Tload:                 opcodeCost.Tload,
		Tstore:                opcodeCost.Tstore,
		Create:                opcodeCost.Create,
		Call:                  opcodeCost.Call,
		Create2:               opcodeCost.Create2,
		Selfdestruct:          opcodeCost.Selfdestruct,
		Memory:                opcodeCost.Memory,
		Copy:                  opcodeCost.Copy,
		Log:                   opcodeCost.Log,
		LogTopic:              opcodeCost.LogTopic,
		LogData:               opcodeCost.LogData,
		Keccak256Word:         opcodeCost.Keccak256Word,
		InitCodeWord:          opcodeCost.InitCodeWord,
		ExpByte:               opcodeCost.ExpByte,
		Exp:                   opcodeCost.Exp,
	}
}

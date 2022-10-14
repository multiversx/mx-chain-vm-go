package mock

import (
	"testing"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/executor"
	worldmock "github.com/ElrondNetwork/wasm-vm/mock/world"
	"github.com/ElrondNetwork/wasm-vm/wasmer"
)

// ExecutorMock can be passed to RuntimeContext as an InstanceBuilder to
// create mocked Wasmer instances.
type ExecutorMock struct {
	InstanceMap map[string]InstanceMock
	World       *worldmock.MockWorld
}

// NewExecutorMock constructs a new InstanceBuilderMock
func NewExecutorMock(world *worldmock.MockWorld) *ExecutorMock {
	return &ExecutorMock{
		InstanceMap: make(map[string]InstanceMock),
		World:       world,
	}
}

// SetOpcodeCosts should set gas costs, but it does nothing in the case of this mock.
func (executorMock *ExecutorMock) SetOpcodeCosts(opcodeCosts *[executor.OpcodeCount]uint32) {
}

// SetRkyvSerializationEnabled controls a Wasmer flag, but it does nothing in the case of this mock.
func (executorMock *ExecutorMock) SetRkyvSerializationEnabled(enabled bool) {
}

// SetSIGSEGVPassthrough controls a Wasmer flag, but it does nothing in the case of this mock.
func (executorMock *ExecutorMock) SetSIGSEGVPassthrough() {
}

func (executorMock *ExecutorMock) FunctionNames() vmcommon.FunctionNames {
	return nil
}

// CreateAndStoreInstanceMock creates a new InstanceMock and registers it as a
// smart contract account in the World, using `code` as the address of the account
func (executorMock *ExecutorMock) CreateAndStoreInstanceMock(t testing.TB, host arwen.VMHost, code []byte, codeHash []byte, codeMetadata []byte, ownerAddress []byte, shardID uint32, balance int64, createAccount bool) *InstanceMock {
	instance := NewInstanceMock(code)
	instance.Address = code
	instance.T = t
	instance.Host = host
	executorMock.InstanceMap[string(code)] = *instance

	if createAccount {
		account := executorMock.World.AcctMap.CreateSmartContractAccountWithCodeHash(nil, code, code, codeHash, executorMock.World)
		account.SetBalance(balance)
		account.ShardID = shardID
		account.CodeMetadata = codeMetadata
		account.OwnerAddress = ownerAddress
	}

	return instance
}

// getNewCopyOfStoredInstance retrieves and initializes a stored Wasmer instance, or
// nil if it doesn't exist
func (executorMock *ExecutorMock) getNewCopyOfStoredInstance(code []byte, gasLimit uint64) (executor.Instance, bool) {
	// this is a map to InstanceMock(s), and copies of these instances will be returned (as the method name indicates)
	instance, ok := executorMock.InstanceMap[string(code)]
	if ok {
		instance.SetPointsUsed(0)
		instance.SetGasLimit(gasLimit)
		return &instance, true
	}
	return nil, false
}

// NewInstanceWithOptions attempts to load a prepared instance using
// GetStoredInstance; if it doesn't exist, it creates a true Wasmer
// instance with the provided contract code.
func (executorMock *ExecutorMock) NewInstanceWithOptions(
	contractCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {

	instance, ok := executorMock.getNewCopyOfStoredInstance(contractCode, options.GasLimit)
	if ok {
		return instance, nil
	}
	return wasmer.NewInstanceWithOptions(contractCode, options)
}

// FIXME
func (executorMock *ExecutorMock) SetVMHooks(instance executor.Instance, hooks executor.VMHooks) {

}

// FIXME
func (executorMock *ExecutorMock) GetVMHooks() executor.VMHooks {
	return nil
}

// NewInstanceFromCompiledCodeWithOptions attempts to load a prepared instance
// using GetStoredInstance; if it doesn't exist, it creates a true Wasmer
// instance with the provided precompiled code.
func (executorMock *ExecutorMock) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	instance, ok := executorMock.getNewCopyOfStoredInstance(compiledCode, options.GasLimit)
	if ok {
		return instance, nil
	}
	return wasmer.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
}

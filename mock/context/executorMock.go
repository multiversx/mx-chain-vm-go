package mock

import (
	"testing"

	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
)

// ExecutorMockFactory is the factory for the ExecutorRecorderMock.
type ExecutorMockFactory struct {
	World *worldmock.MockWorld

	// gives access to the created Executor in tests
	LastCreatedExecutor *ExecutorMock
}

// NewExecutorMockFactory returns the Wasmer executor factory.
func NewExecutorMockFactory(world *worldmock.MockWorld) *ExecutorMockFactory {
	return &ExecutorMockFactory{
		World: world,
	}
}

// CreateExecutor creates a new Executor instance.
func (emf *ExecutorMockFactory) CreateExecutor(_ executor.ExecutorFactoryArgs) (executor.Executor, error) {
	executorMock := NewExecutorMock(emf.World)
	emf.LastCreatedExecutor = executorMock
	return executorMock, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (emf *ExecutorMockFactory) IsInterfaceNil() bool {
	return emf == nil
}

// ExecutorMock can be passed to RuntimeContext as an InstanceBuilder to
// create mocked Wasmer instances.
type ExecutorMock struct {
	InstanceMap    map[string]InstanceMock
	World          *worldmock.MockWorld
	WasmerExecutor *wasmer2.Wasmer2Executor
}

// NewExecutorMock constructs a new InstanceBuilderMock
func NewExecutorMock(world *worldmock.MockWorld) *ExecutorMock {
	exec, _ := wasmer2.CreateExecutor()
	return &ExecutorMock{
		InstanceMap:    make(map[string]InstanceMock),
		World:          world,
		WasmerExecutor: exec,
	}
}

// SetOpcodeCosts should set gas costs, but it does nothing in the case of this mock.
func (executorMock *ExecutorMock) SetOpcodeCosts(_ *executor.WASMOpcodeCost) {
}

// FunctionNames mocked method
func (executorMock *ExecutorMock) FunctionNames() vmcommon.FunctionNames {
	return functionNames
}

// CreateAndStoreInstanceMock creates a new InstanceMock and registers it as a
// smart contract account in the World, using `code` as the address of the account
func (executorMock *ExecutorMock) CreateAndStoreInstanceMock(t testing.TB, host vmhost.VMHost, code []byte, codeHash []byte, codeMetadata []byte, ownerAddress []byte, shardID uint32, balance int64, createAccount bool) *InstanceMock {
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
	return executorMock.WasmerExecutor.NewInstanceWithOptions(contractCode, options)
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
	return executorMock.WasmerExecutor.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
}

// IsInterfaceNil returns true if there is no value under the interface
func (executorMock *ExecutorMock) IsInterfaceNil() bool {
	return executorMock == nil
}

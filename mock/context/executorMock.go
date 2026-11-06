package mock

import (
	"testing"

	"github.com/multiversx/mx-chain-scenario-go/worldmock"
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

// CreateExecutorNew creates a new ExecutorNew instance.
func (emf *ExecutorMockFactory) CreateExecutorNew(args executor.ExecutorFactoryArgsNew) (executor.ExecutorNew, error) {
	executorMock, err := NewExecutorMock(emf.World, args.Runtime)
	if err != nil {
		return nil, err
	}

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
	WasmerExecutor *wasmer2.Wasmer2ExecutorNew
}

// NewExecutorMock constructs a new InstanceBuilderMock
//
// It falls back to a real Wasmer 2 executor for any contract no mock was registered for,
// which is why it needs the runtime: that fallback is a real executor and services real
// hooks.
func NewExecutorMock(world *worldmock.MockWorld, runtime executor.ExecutorRuntime) (*ExecutorMock, error) {
	exec, err := wasmer2.CreateExecutorNew(runtime)
	if err != nil {
		return nil, err
	}

	return &ExecutorMock{
		InstanceMap:    make(map[string]InstanceMock),
		World:          world,
		WasmerExecutor: exec,
	}, nil
}

// SetOpcodeConfig should set gas costs, but it does nothing in the case of this mock.
func (executorMock *ExecutorMock) SetOpcodeConfig(_ executor.OpcodeVersion, _ *executor.WASMOpcodeCost) error {
	return nil
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
		if codeHash != nil {
			account.CodeHash = codeHash
		}
		account.OwnerAddress = ownerAddress
	}

	return instance
}

// getNewCopyOfStoredInstance retrieves and initializes a stored Wasmer instance, or
// nil if it doesn't exist
func (executorMock *ExecutorMock) getNewCopyOfStoredInstance(code []byte) (executor.InstanceNew, bool) {
	// this is a map to InstanceMock(s), and copies of these instances will be returned (as the method name indicates)
	instance, ok := executorMock.InstanceMap[string(code)]
	if ok {
		// Zeroed directly rather than through a setter: Call now carries the gas limit, and
		// a copy of a stored instance has to start where a fresh one would.
		instance.Points = 0
		return &instance, true
	}
	return nil, false
}

// NewInstance attempts to load a prepared instance using GetStoredInstance; if it doesn't
// exist, it creates a true Wasmer instance with the provided contract code.
func (executorMock *ExecutorMock) NewInstance(
	contractCode []byte,
	options *executor.CompilationOptionsNew,
) (executor.InstanceNew, error) {

	instance, ok := executorMock.getNewCopyOfStoredInstance(contractCode)
	if ok {
		return instance, nil
	}
	return executorMock.WasmerExecutor.NewInstance(contractCode, options)
}

// NewInstanceFromCache attempts to load a prepared instance using GetStoredInstance; if it
// doesn't exist, it creates a true Wasmer instance with the provided precompiled code.
func (executorMock *ExecutorMock) NewInstanceFromCache(
	compiledCode []byte,
	options *executor.CompilationOptionsNew,
) (executor.InstanceNew, error) {
	instance, ok := executorMock.getNewCopyOfStoredInstance(compiledCode)
	if ok {
		return instance, nil
	}
	return executorMock.WasmerExecutor.NewInstanceFromCache(compiledCode, options)
}

// IsInterfaceNil returns true if there is no value under the interface
func (executorMock *ExecutorMock) IsInterfaceNil() bool {
	return executorMock == nil
}

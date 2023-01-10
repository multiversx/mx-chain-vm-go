package mock

import (
	"testing"

	"github.com/multiversx/mx-chain-vm-v1_4-go/arwen"
	worldmock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/world"
	"github.com/multiversx/mx-chain-vm-v1_4-go/wasmer"
)

// InstanceBuilderMock can be passed to RuntimeContext as an InstanceBuilder to
// create mocked Wasmer instances.
type InstanceBuilderMock struct {
	InstanceMap map[string]InstanceMock
	World       *worldmock.MockWorld
}

// NewInstanceBuilderMock constructs a new InstanceBuilderMock
func NewInstanceBuilderMock(world *worldmock.MockWorld) *InstanceBuilderMock {
	return &InstanceBuilderMock{
		InstanceMap: make(map[string]InstanceMock),
		World:       world,
	}
}

// Id -
func (builder *InstanceBuilderMock) Id() string {
	return ""
}

// CreateAndStoreInstanceMock creates a new InstanceMock and registers it as a
// smart contract account in the World, using `code` as the address of the account
func (builder *InstanceBuilderMock) CreateAndStoreInstanceMock(t testing.TB, host arwen.VMHost, code []byte, codeHash []byte, codeMetadata []byte, ownerAddress []byte, shardID uint32, balance int64) *InstanceMock {
	instance := NewInstanceMock(code)
	instance.Address = code
	instance.T = t
	instance.Host = host
	builder.InstanceMap[string(code)] = *instance

	account := builder.World.AcctMap.CreateSmartContractAccountWithCodeHash(nil, code, code, codeHash, builder.World)
	account.SetBalance(balance)
	account.ShardID = shardID
	account.CodeMetadata = codeMetadata
	account.OwnerAddress = ownerAddress

	return instance
}

// getNewCopyOfStoredInstance retrieves and initializes a stored Wasmer instance, or
// nil if it doesn't exist
func (builder *InstanceBuilderMock) getNewCopyOfStoredInstance(code []byte, gasLimit uint64) (wasmer.InstanceHandler, bool) {
	// this is a map to InstanceMock(s), and copies of these instances will be returned (as the method name indicates)
	instance, ok := builder.InstanceMap[string(code)]
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
func (builder *InstanceBuilderMock) NewInstanceWithOptions(
	contractCode []byte,
	options wasmer.CompilationOptions,
) (wasmer.InstanceHandler, error) {

	instance, ok := builder.getNewCopyOfStoredInstance(contractCode, options.GasLimit)
	if ok {
		return instance, nil
	}
	return wasmer.NewInstanceWithOptions(contractCode, options)
}

// NewInstanceFromCompiledCodeWithOptions attempts to load a prepared instance
// using GetStoredInstance; if it doesn't exist, it creates a true Wasmer
// instance with the provided precompiled code.
func (builder *InstanceBuilderMock) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options wasmer.CompilationOptions,
) (wasmer.InstanceHandler, error) {
	instance, ok := builder.getNewCopyOfStoredInstance(compiledCode, options.GasLimit)
	if ok {
		return instance, nil
	}
	return wasmer.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
}

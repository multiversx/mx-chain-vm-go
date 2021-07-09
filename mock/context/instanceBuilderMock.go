package mock

import (
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_2/mock/world"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_2/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// InstanceBuilderMock can be passed to RuntimeContext as an InstanceBuilder to
// create mocked Wasmer instances.
type InstanceBuilderMock struct {
	InstanceMap map[string]wasmer.InstanceHandler
	World       *worldmock.MockWorld
}

// NewInstanceBuilderMock constructs a new InstanceBuilderMock
func NewInstanceBuilderMock(world *worldmock.MockWorld) *InstanceBuilderMock {
	return &InstanceBuilderMock{
		InstanceMap: make(map[string]wasmer.InstanceHandler),
		World:       world,
	}
}

// CreateAndStoreInstanceMock creates a new InstanceMock and registers it as a
// smart contract account in the World, using `code` as the address of the account
func (builder *InstanceBuilderMock) CreateAndStoreInstanceMock(code []byte, balance int64) *InstanceMock {
	instance := NewInstanceMock(code)
	builder.InstanceMap[string(code)] = instance

	account := builder.World.AcctMap.CreateAccount(code)
	account.IsSmartContract = true
	account.SetBalance(balance)
	account.Code = code
	account.CodeMetadata = []byte{0, vmcommon.MetadataPayable}

	return instance
}

// GetStoredInstance retrieves and initializes a stored Wasmer instance, or
// nil if it doesn't exist
func (builder *InstanceBuilderMock) GetStoredInstance(code []byte, gasLimit uint64) (wasmer.InstanceHandler, bool) {
	instance, ok := builder.InstanceMap[string(code)]
	if ok {
		instance.SetPointsUsed(0)
		instance.SetGasLimit(gasLimit)
		return instance, true
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

	instance, ok := builder.GetStoredInstance(contractCode, options.GasLimit)
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
	instance, ok := builder.GetStoredInstance(compiledCode, options.GasLimit)
	if ok {
		return instance, nil
	}
	return wasmer.NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
}

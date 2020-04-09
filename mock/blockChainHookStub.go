package mock

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ vmcommon.BlockchainHook = (*BlockchainHookStub)(nil)

type BlockchainHookStub struct {
	AccountExtistsCalled          func(address []byte) (bool, error)
	NewAddressCalled              func(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error)
	GetBalanceCalled              func(address []byte) (*big.Int, error)
	GetNonceCalled                func(address []byte) (uint64, error)
	GetStorageDataCalled          func(accountsAddress []byte, index []byte) ([]byte, error)
	IsCodeEmptyCalled             func(address []byte) (bool, error)
	GetCodeCalled                 func(address []byte) ([]byte, error)
	GetBlockHashCalled            func(nonce uint64) ([]byte, error)
	LastNonceCalled               func() uint64
	LastRoundCalled               func() uint64
	LastTimeStampCalled           func() uint64
	LastRandomSeedCalled          func() []byte
	LastEpochCalled               func() uint32
	GetStateRootHashCalled        func() []byte
	CurrentNonceCalled            func() uint64
	CurrentRoundCalled            func() uint64
	CurrentTimeStampCalled        func() uint64
	CurrentRandomSeedCalled       func() []byte
	CurrentEpochCalled            func() uint32
	ProcessBuiltInFunctionCalled  func(input *vmcommon.ContractCallInput) (*big.Int, uint64, error)
	GetBuiltinFunctionNamesCalled func() vmcommon.FunctionNames
}

func (b *BlockchainHookStub) AccountExists(address []byte) (bool, error) {
	if b.AccountExtistsCalled != nil {
		return b.AccountExtistsCalled(address)
	}
	return false, nil
}

func (b *BlockchainHookStub) NewAddress(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
	if b.NewAddressCalled != nil {
		return b.NewAddressCalled(creatorAddress, creatorNonce, vmType)
	}
	return []byte("newAddress"), nil
}

func (b *BlockchainHookStub) GetBalance(address []byte) (*big.Int, error) {
	if b.GetBalanceCalled != nil {
		return b.GetBalanceCalled(address)
	}
	return big.NewInt(0), nil
}

func (b *BlockchainHookStub) GetNonce(address []byte) (uint64, error) {
	if b.GetNonceCalled != nil {
		return b.GetNonceCalled(address)
	}
	return 0, nil
}

func (b *BlockchainHookStub) GetStorageData(accountAddress []byte, index []byte) ([]byte, error) {
	if b.GetStorageDataCalled != nil {
		return b.GetStorageDataCalled(accountAddress, index)
	}
	return nil, nil
}

func (b *BlockchainHookStub) IsCodeEmpty(address []byte) (bool, error) {
	if b.IsCodeEmptyCalled != nil {
		return b.IsCodeEmptyCalled(address)
	}
	return true, nil
}

func (b *BlockchainHookStub) GetCode(address []byte) ([]byte, error) {
	if b.GetCodeCalled != nil {
		return b.GetCodeCalled(address)
	}
	return nil, nil
}

func (b *BlockchainHookStub) GetBlockhash(nonce uint64) ([]byte, error) {
	if b.GetBlockHashCalled != nil {
		return b.GetBlockHashCalled(nonce)
	}
	return []byte("roothash"), nil
}

func (b *BlockchainHookStub) LastNonce() uint64 {
	if b.LastNonceCalled != nil {
		return b.LastNonceCalled()
	}
	return 0
}

func (b *BlockchainHookStub) LastRound() uint64 {
	if b.LastRoundCalled != nil {
		return b.LastRoundCalled()
	}
	return 0
}

func (b *BlockchainHookStub) LastTimeStamp() uint64 {
	if b.LastTimeStampCalled != nil {
		return b.LastTimeStampCalled()
	}
	return 0
}

func (b *BlockchainHookStub) LastRandomSeed() []byte {
	if b.LastRandomSeedCalled != nil {
		return b.LastRandomSeedCalled()
	}
	return []byte("seed")
}

func (b *BlockchainHookStub) LastEpoch() uint32 {
	if b.LastEpochCalled != nil {
		return b.LastEpochCalled()
	}
	return 0
}

func (b *BlockchainHookStub) GetStateRootHash() []byte {
	if b.GetStateRootHashCalled != nil {
		return b.GetStateRootHashCalled()
	}
	return []byte("roothash")
}

func (b *BlockchainHookStub) CurrentNonce() uint64 {
	if b.CurrentNonceCalled != nil {
		return b.CurrentNonceCalled()
	}
	return 0
}

func (b *BlockchainHookStub) CurrentRound() uint64 {
	if b.CurrentRoundCalled != nil {
		return b.CurrentRoundCalled()
	}
	return 0
}

func (b *BlockchainHookStub) CurrentTimeStamp() uint64 {
	if b.CurrentTimeStampCalled != nil {
		return b.CurrentTimeStampCalled()
	}
	return 0
}

func (b *BlockchainHookStub) CurrentRandomSeed() []byte {
	if b.CurrentRandomSeedCalled != nil {
		return b.CurrentRandomSeedCalled()
	}
	return []byte("seed")
}

func (b *BlockchainHookStub) CurrentEpoch() uint32 {
	if b.CurrentEpochCalled != nil {
		return b.CurrentEpochCalled()
	}
	return 0
}

func (b *BlockchainHookStub) ProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*big.Int, uint64, error) {
	if b.ProcessBuiltInFunctionCalled != nil {
		return b.ProcessBuiltInFunction(input)
	}
	return arwen.Zero, 0, nil
}

func (b *BlockchainHookStub) GetBuiltinFunctionNames() vmcommon.FunctionNames {
	if b.GetBuiltinFunctionNamesCalled != nil {
		return b.GetBuiltinFunctionNamesCalled()
	}
	return make(vmcommon.FunctionNames)
}

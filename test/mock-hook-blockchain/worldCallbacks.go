package callbackblockchain

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

var _ vmcommon.BlockchainHook = (*BlockchainHookMock)(nil)

var zero = big.NewInt(0)

// NewAddress adapts between K model and elrond function
func (b *BlockchainHookMock) NewAddress(creatorAddress []byte, creatorNonce uint64, _ []byte) ([]byte, error) {
	// explicit new address mocks
	for _, newAddressMock := range b.NewAddressMocks {
		if bytes.Equal(creatorAddress, newAddressMock.CreatorAddress) && creatorNonce == newAddressMock.CreatorNonce {
			return newAddressMock.NewAddress, nil
		}
	}

	// a simple mock algorithm
	if b.mockAddressGenerationEnabled {
		result := make([]byte, 32)
		result[10] = 0x11
		result[11] = 0x11
		result[12] = 0x11
		result[13] = 0x11
		copy(result[14:29], creatorAddress)

		result[29] = byte(creatorNonce)

		copy(result[30:], creatorAddress[30:])

		return result, nil
	}
	// empty byte array signals not implemented, fallback to default
	return []byte{}, nil
}

// GetStorageData yields the storage value for a certain account and index.
// Should return an empty byte array if the key is missing from the account storage
func (b *BlockchainHookMock) GetStorageData(accountAddress []byte, index []byte) ([]byte, error) {
	acct := b.AcctMap.GetAccount(accountAddress)
	if acct == nil {
		return []byte{}, nil
	}
	return acct.StorageValue(string(index)), nil
}

// GetBlockhash should return the hash of the nth previous blockchain.
// Offset specifies how many blocks we need to look back.
func (b *BlockchainHookMock) GetBlockhash(nonce uint64) ([]byte, error) {
	offsetInt32 := int(nonce)
	if offsetInt32 >= len(b.Blockhashes) {
		return nil, errors.New("blockhash offset exceeds the blockhashes slice")
	}
	return b.Blockhashes[offsetInt32], nil
}

// LastNonce returns the nonce from from the last committed block
func (b *BlockchainHookMock) LastNonce() uint64 {
	if b.PreviousBlockInfo == nil {
		return 0
	}
	return b.PreviousBlockInfo.BlockNonce
}

// LastRound returns the round from the last committed block
func (b *BlockchainHookMock) LastRound() uint64 {
	if b.PreviousBlockInfo == nil {
		return 0
	}
	return b.PreviousBlockInfo.BlockRound
}

// LastTimeStamp returns the timeStamp from the last committed block
func (b *BlockchainHookMock) LastTimeStamp() uint64 {
	if b.PreviousBlockInfo == nil {
		return 0
	}
	return b.PreviousBlockInfo.BlockTimestamp
}

// LastRandomSeed returns the random seed from the last committed block
func (b *BlockchainHookMock) LastRandomSeed() []byte {
	return nil
}

// LastEpoch returns the epoch from the last committed block
func (b *BlockchainHookMock) LastEpoch() uint32 {
	if b.PreviousBlockInfo == nil {
		return 0
	}
	return b.PreviousBlockInfo.BlockEpoch
}

// GetStateRootHash returns the state root hash from the last committed block
func (b *BlockchainHookMock) GetStateRootHash() []byte {
	return nil
}

// CurrentNonce returns the nonce from the current block
func (b *BlockchainHookMock) CurrentNonce() uint64 {
	if b.CurrentBlockInfo == nil {
		return 0
	}
	return b.CurrentBlockInfo.BlockNonce
}

// CurrentRound returns the round from the current block
func (b *BlockchainHookMock) CurrentRound() uint64 {
	if b.CurrentBlockInfo == nil {
		return 0
	}
	return b.CurrentBlockInfo.BlockRound
}

// CurrentTimeStamp return the timestamp from the current block
func (b *BlockchainHookMock) CurrentTimeStamp() uint64 {
	if b.CurrentBlockInfo == nil {
		return 0
	}
	return b.CurrentBlockInfo.BlockTimestamp
}

// CurrentRandomSeed returns the random seed from the current header
func (b *BlockchainHookMock) CurrentRandomSeed() []byte {
	return nil
}

// CurrentEpoch returns the current epoch
func (b *BlockchainHookMock) CurrentEpoch() uint32 {
	if b.CurrentBlockInfo == nil {
		return 0
	}
	return b.CurrentBlockInfo.BlockEpoch
}

// ProcessBuiltInFunction -
func (b *BlockchainHookMock) ProcessBuiltInFunction(_ *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	return &vmcommon.VMOutput{}, nil
}

// GetBuiltinFunctionNames -
func (b *BlockchainHookMock) GetBuiltinFunctionNames() vmcommon.FunctionNames {
	return make(vmcommon.FunctionNames)
}

// GetAllState -
func (b *BlockchainHookMock) GetAllState(_ []byte) (map[string][]byte, error) {
	return make(map[string][]byte), nil
}

// GetUserAccount retrieves account info from map, or error if not found.
func (b *BlockchainHookMock) GetUserAccount(address []byte) (vmcommon.UserAccountHandler, error) {
	account := b.AcctMap.GetAccount(address)
	if account == nil {
		return nil, fmt.Errorf("account not found: %s", hex.EncodeToString(address))
	}

	return account, nil
}

// GetShardOfAddress -
func (b *BlockchainHookMock) GetShardOfAddress(address []byte) uint32 {
	account := b.AcctMap.GetAccount(address)
	if account == nil {
		return 0
	}

	return account.ShardID
}

// IsSmartContract -
func (b *BlockchainHookMock) IsSmartContract(address []byte) bool {
	account := b.AcctMap.GetAccount(address)
	if account == nil {
		return false
	}

	return account.IsSmartContract
}

func (b *BlockchainHookMock) IsPayable(address []byte) (bool, error) {
	account := b.AcctMap.GetAccount(address)
	if account == nil {
		return true, nil
	}

	if !account.IsSmartContract {
		return true, nil
	}

	metadata := vmcommon.CodeMetadataFromBytes(account.CodeMetadata)
	return metadata.Payable, nil
}

func (b *BlockchainHookMock) SaveCompiledCode(_ []byte, _ []byte) {
}

func (b *BlockchainHookMock) GetCompiledCode(_ []byte) (bool, []byte) {
	return false, nil
}

func (b *BlockchainHookMock) ClearCompiledCodes() {
}

// IsInterfaceNil returns true if underlying implementation is nil
func (b *BlockchainHookMock) IsInterfaceNil() bool {
	return b == nil
}

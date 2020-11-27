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
	// custom error
	if b.Err != nil {
		return nil, b.Err
	}

	// explicit new address mocks
	for _, newAddressMock := range b.NewAddressMocks {
		if bytes.Equal(creatorAddress, newAddressMock.CreatorAddress) && creatorNonce == newAddressMock.CreatorNonce {
			b.LastCreatedContractAddress = newAddressMock.NewAddress
			return newAddressMock.NewAddress, nil
		}
	}

	// a simple mock algorithm
	if b.mockAddressGenerationEnabled {
		result := GenerateMockAddress(creatorAddress, creatorNonce)
		b.LastCreatedContractAddress = result
		return result, nil
	}
	// empty byte array signals not implemented, fallback to default
	return []byte{}, nil
}

// GetStorageData yields the storage value for a certain account and index.
// Should return an empty byte array if the key is missing from the account storage
func (b *BlockchainHookMock) GetStorageData(accountAddress []byte, index []byte) ([]byte, error) {
	// custom error
	if b.Err != nil {
		return nil, b.Err
	}

	acct := b.AcctMap.GetAccount(accountAddress)
	if acct == nil {
		return []byte{}, nil
	}
	return acct.StorageValue(string(index)), nil
}

// GetBlockhash should return the hash of the nth previous blockchain.
// Offset specifies how many blocks we need to look back.
func (b *BlockchainHookMock) GetBlockhash(nonce uint64) ([]byte, error) {
	if b.Err != nil {
		return nil, b.Err
	}
	currentNonce := b.CurrentNonce()
	if nonce > currentNonce {
		return nil, errors.New("blockhash nonce exceeds current nonce")
	}
	offsetInt32 := int(currentNonce - nonce)
	if offsetInt32 >= len(b.Blockhashes) {
		return nil, errors.New("blockhash nonce is older than what is available")
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
	if b.PreviousBlockInfo == nil {
		return nil
	}
	return b.PreviousBlockInfo.RandomSeed
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
	return b.StateRootHash
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
	if b.CurrentBlockInfo == nil {
		return nil
	}
	return b.CurrentBlockInfo.RandomSeed
}

// CurrentEpoch returns the current epoch
func (b *BlockchainHookMock) CurrentEpoch() uint32 {
	if b.CurrentBlockInfo == nil {
		return 0
	}
	return b.CurrentBlockInfo.BlockEpoch
}

// ProcessBuiltInFunction -
func (b *BlockchainHookMock) ProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	// custom error
	if b.Err != nil {
		return nil, b.Err
	}

	// outPutAccounts := make(map[string]*vmcommon.OutputAccount)
	// outPutAccounts[string(input.CallerAddr)] = &vmcommon.OutputAccount{BalanceDelta: b.Value}

	// return &vmcommon.VMOutput{
	// 	GasRemaining:   b.Gas,
	// 	OutputAccounts: outPutAccounts,
	// }, nil
	return &vmcommon.VMOutput{}, nil
}

// GetBuiltinFunctionNames -
func (b *BlockchainHookMock) GetBuiltinFunctionNames() vmcommon.FunctionNames {
	return make(vmcommon.FunctionNames)
}

// GetAllState simply returns the storage as-is.
func (b *BlockchainHookMock) GetAllState(accountAddress []byte) (map[string][]byte, error) {
	account := b.AcctMap.GetAccount(accountAddress)
	if account == nil {
		return nil, fmt.Errorf("account not found: %s", hex.EncodeToString(accountAddress))
	}
	return account.Storage, nil
}

// GetUserAccount retrieves account info from map, or error if not found.
func (b *BlockchainHookMock) GetUserAccount(address []byte) (vmcommon.UserAccountHandler, error) {
	// custom error
	if b.Err != nil {
		return nil, b.Err
	}

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

func (b *BlockchainHookMock) SaveCompiledCode(codeHash []byte, code []byte) {
	b.CompiledCode[string(codeHash)] = code
}

func (b *BlockchainHookMock) GetCompiledCode(codeHash []byte) (bool, []byte) {
	code, found := b.CompiledCode[string(codeHash)]
	return found, code
}

func (b *BlockchainHookMock) ClearCompiledCodes() {
	b.CompiledCode = make(map[string][]byte)
}

// IsInterfaceNil returns true if underlying implementation is nil
func (b *BlockchainHookMock) IsInterfaceNil() bool {
	return b == nil
}

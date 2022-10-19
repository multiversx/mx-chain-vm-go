package worldmock

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data/esdt"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ vmcommon.BlockchainHook = (*MockWorld)(nil)

// ErrBuiltinFuncWrapperNotInitialized means that the builtin function wrapper was used before initialization.
var ErrBuiltinFuncWrapperNotInitialized = errors.New("builtin function not found or container not initialized")

var zero = big.NewInt(0)

// NewAddress provides the address for a new account.
// It looks up the explicit new address mocks, if none found generates one using a fake but realistic algorithm.
func (b *MockWorld) NewAddress(creatorAddress []byte, creatorNonce uint64, _ []byte) ([]byte, error) {
	// custom error
	if b.Err != nil {
		return nil, b.Err
	}

	// explicit new address mocks
	// matched by creator address and nonce
	for _, newAddressMock := range b.NewAddressMocks {
		if bytes.Equal(creatorAddress, newAddressMock.CreatorAddress) && creatorNonce == newAddressMock.CreatorNonce {
			b.LastCreatedContractAddress = newAddressMock.NewAddress
			return newAddressMock.NewAddress, nil
		}
	}

	// If a mock address wasn't registered for the specified creatorAddress, generate one automatically.
	// This is not the real algorithm but it's simple and close enough.
	result := GenerateMockAddress(creatorAddress, creatorNonce)
	b.LastCreatedContractAddress = result
	return result, nil
}

// GetStorageData yields the storage value for a certain account and storage key.
// Should return an empty byte array if the key is missing from the account storage
func (b *MockWorld) GetStorageData(accountAddress []byte, key []byte) ([]byte, uint32, error) {
	// custom error
	if b.Err != nil {
		return nil, 0, b.Err
	}

	acct := b.AcctMap.GetAccount(accountAddress)
	if acct == nil {
		return []byte{}, 0, nil
	}

	foundValue := acct.StorageValue(string(key))
	if len(foundValue) == 0 &&
		b.ProvidedBlockchainHook != nil {
		bhStoredValue, _, err := b.ProvidedBlockchainHook.GetStorageData(accountAddress, key)
		if err == nil {
			foundValue = bhStoredValue
		}
	}

	return foundValue, 0, nil
}

// GetBlockhash should return the hash of the nth previous blockchain.
// Offset specifies how many blocks we need to look back.
func (b *MockWorld) GetBlockhash(nonce uint64) ([]byte, error) {
	if b.Err != nil {
		return nil, b.Err
	}
	currentNonce := b.CurrentNonce()
	if nonce > currentNonce {
		return nil, errors.New("requested nonce is greater than current nonce")
	}
	offsetInt32 := int(currentNonce - nonce)
	if offsetInt32 >= len(b.Blockhashes) {
		return nil, errors.New("requested nonce is older than the oldest available block nonce")
	}
	return b.Blockhashes[offsetInt32], nil
}

// LastNonce returns the nonce from from the last committed block
func (b *MockWorld) LastNonce() uint64 {
	if b.PreviousBlockInfo == nil {
		return 0
	}
	return b.PreviousBlockInfo.BlockNonce
}

// LastRound returns the round from the last committed block
func (b *MockWorld) LastRound() uint64 {
	if b.PreviousBlockInfo == nil {
		return 0
	}
	return b.PreviousBlockInfo.BlockRound
}

// LastTimeStamp returns the timeStamp from the last committed block
func (b *MockWorld) LastTimeStamp() uint64 {
	if b.PreviousBlockInfo == nil {
		return 0
	}
	return b.PreviousBlockInfo.BlockTimestamp
}

// LastRandomSeed returns the random seed from the last committed block
func (b *MockWorld) LastRandomSeed() []byte {
	if b.PreviousBlockInfo == nil {
		return nil
	}
	return b.PreviousBlockInfo.GetRandomSeedSlice()
}

// LastEpoch returns the epoch from the last committed block
func (b *MockWorld) LastEpoch() uint32 {
	if b.PreviousBlockInfo == nil {
		return 0
	}
	return b.PreviousBlockInfo.BlockEpoch
}

// GetStateRootHash returns the state root hash from the last committed block
func (b *MockWorld) GetStateRootHash() []byte {
	return b.StateRootHash
}

// CurrentNonce returns the nonce from the current block
func (b *MockWorld) CurrentNonce() uint64 {
	if b.CurrentBlockInfo == nil {
		return 0
	}
	return b.CurrentBlockInfo.BlockNonce
}

// CurrentRound returns the round from the current block
func (b *MockWorld) CurrentRound() uint64 {
	if b.CurrentBlockInfo == nil {
		return 0
	}
	return b.CurrentBlockInfo.BlockRound
}

// CurrentTimeStamp return the timestamp from the current block
func (b *MockWorld) CurrentTimeStamp() uint64 {
	if b.CurrentBlockInfo == nil {
		return 0
	}
	return b.CurrentBlockInfo.BlockTimestamp
}

// CurrentRandomSeed returns the random seed from the current header
func (b *MockWorld) CurrentRandomSeed() []byte {
	if b.CurrentBlockInfo == nil {
		return nil
	}
	return b.CurrentBlockInfo.GetRandomSeedSlice()
}

// CurrentEpoch returns the current epoch
func (b *MockWorld) CurrentEpoch() uint32 {
	if b.CurrentBlockInfo == nil {
		return 0
	}
	return b.CurrentBlockInfo.BlockEpoch
}

// ProcessBuiltInFunction -
func (b *MockWorld) ProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	// custom error
	if b.Err != nil {
		return nil, b.Err
	}

	if b.BuiltinFuncs == nil {
		return nil, ErrBuiltinFuncWrapperNotInitialized
	}

	return b.BuiltinFuncs.ProcessBuiltInFunction(input)
}

// GetESDTToken -
func (b *MockWorld) GetESDTToken(address []byte, tokenIdentifier []byte, nonce uint64) (*esdt.ESDigitalToken, error) {
	// custom error
	if b.Err != nil {
		return nil, b.Err
	}

	if b.BuiltinFuncs == nil {
		return nil, ErrBuiltinFuncWrapperNotInitialized
	}

	return b.BuiltinFuncs.GetTokenData(address, tokenIdentifier, nonce)
}

// GetBuiltinFunctionNames -
func (b *MockWorld) GetBuiltinFunctionNames() vmcommon.FunctionNames {
	return b.BuiltinFuncs.GetBuiltinFunctionNames()
}

// GetAllState simply returns the storage as-is.
func (b *MockWorld) GetAllState(accountAddress []byte) (map[string][]byte, error) {
	account := b.AcctMap.GetAccount(accountAddress)
	if account == nil {
		return nil, fmt.Errorf("account not found: %s", hex.EncodeToString(accountAddress))
	}
	return account.Storage, nil
}

// GetUserAccount retrieves account info from map, or error if not found.
func (b *MockWorld) GetUserAccount(address []byte) (vmcommon.UserAccountHandler, error) {
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

// GetCode retrieves the code from the given account, or nil if not found
func (b *MockWorld) GetCode(acc vmcommon.UserAccountHandler) []byte {
	account := b.AcctMap.GetAccount(acc.AddressBytes())
	if account == nil {
		return nil
	}

	return account.Code
}

// GetShardOfAddress -
func (b *MockWorld) GetShardOfAddress(address []byte) uint32 {
	account := b.AcctMap.GetAccount(address)
	if account == nil {
		return 0
	}

	return account.ShardID
}

// IsSmartContract -
func (b *MockWorld) IsSmartContract(address []byte) bool {
	account := b.AcctMap.GetAccount(address)
	if account == nil {
		return vmcommon.IsSmartContractAddress(address)
	}

	return account.IsSmartContract
}

// IsPayable -
func (b *MockWorld) IsPayable(sndAddress []byte, rcvAddress []byte) (bool, error) {
	account := b.AcctMap.GetAccount(rcvAddress)
	if account == nil {
		return true, nil
	}

	if !account.IsSmartContract {
		return true, nil
	}

	metadata := vmcommon.CodeMetadataFromBytes(account.CodeMetadata)
	if core.IsSmartContractAddress(sndAddress) {
		return metadata.PayableBySC || metadata.Payable, nil
	}

	return metadata.Payable, nil
}

// SaveCompiledCode -
func (b *MockWorld) SaveCompiledCode(codeHash []byte, code []byte) {
	b.CompiledCode[string(codeHash)] = code
}

// GetCompiledCode -
func (b *MockWorld) GetCompiledCode(codeHash []byte) (bool, []byte) {
	code, found := b.CompiledCode[string(codeHash)]
	return found, code
}

// ClearCompiledCodes -
func (b *MockWorld) ClearCompiledCodes() {
	b.CompiledCode = make(map[string][]byte)
}

// IsPaused -
func (b *MockWorld) IsPaused(_ []byte) bool {
	return b.IsPausedValue
}

// IsLimitedTransfer -
func (b *MockWorld) IsLimitedTransfer(_ []byte) bool {
	return b.IsLimitedTransferValue
}

// IsInterfaceNil returns true if underlying implementation is nil
func (b *MockWorld) IsInterfaceNil() bool {
	return b == nil
}

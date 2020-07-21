package arwendebug

import (
	"math/big"
	"strconv"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ vmcommon.BlockchainHook = (*BlockchainHookMock)(nil)

// BlockchainHookMock -
type BlockchainHookMock struct {
	Accounts      AccountsMap
	BlockHash     []byte
	LNonce        uint64
	LRound        uint64
	CNonce        uint64
	CRound        uint64
	LTimeStamp    uint64
	CTimeStamp    uint64
	LRandomSeed   []byte
	CRandomSeed   []byte
	LEpoch        uint32
	CEpoch        uint32
	StateRootHash []byte
	Value         *big.Int
	Gas           uint64

	LastCreatedContractAddress []byte
}

// NewBlockchainHookMock -
func NewBlockchainHookMock() *BlockchainHookMock {
	return &BlockchainHookMock{
		Accounts: make(AccountsMap),
	}
}

// AddAccount -
func (b *BlockchainHookMock) AddAccount(account *Account) {
	b.Accounts[account.AddressHex] = account
}

// AddAccounts -
func (b *BlockchainHookMock) AddAccounts(accounts []*Account) {
	for _, account := range accounts {
		b.AddAccount(account)
	}
}

// NewAddress -
func (b *BlockchainHookMock) NewAddress(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
	if len(creatorAddress) != arwen.AddressLen {
		panic("mock: bad creator address")
	}

	address := make([]byte, arwen.AddressLen)
	copy(address, creatorAddress)
	copy(address, []byte("contract"))
	copy(address[len("contract"):], strconv.Itoa(int(creatorNonce)))
	b.LastCreatedContractAddress = address
	return address, nil
}

// GetStorageData -
func (b *BlockchainHookMock) GetStorageData(address []byte, index []byte) ([]byte, error) {
	account, ok := b.Accounts[toHex(address)]
	if !ok {
		return []byte{}, ErrAccountDoesntExist
	}

	return fromHex(account.Storage[toHex(index)])
}

// GetBlockhash -
func (b *BlockchainHookMock) GetBlockhash(nonce uint64) ([]byte, error) {
	return b.BlockHash, nil
}

// LastNonce -
func (b *BlockchainHookMock) LastNonce() uint64 {
	return b.LNonce
}

// LastRound -
func (b *BlockchainHookMock) LastRound() uint64 {
	return b.LRound
}

// LastTimeStamp -
func (b *BlockchainHookMock) LastTimeStamp() uint64 {
	return b.LTimeStamp
}

// LastRandomSeed -
func (b *BlockchainHookMock) LastRandomSeed() []byte {
	return b.LRandomSeed
}

// LastEpoch -
func (b *BlockchainHookMock) LastEpoch() uint32 {
	return b.LEpoch
}

// GetStateRootHash -
func (b *BlockchainHookMock) GetStateRootHash() []byte {
	return b.StateRootHash
}

// CurrentNonce -
func (b *BlockchainHookMock) CurrentNonce() uint64 {
	return b.CNonce
}

// CurrentRound -
func (b *BlockchainHookMock) CurrentRound() uint64 {
	return b.CRound
}

// CurrentTimeStamp -
func (b *BlockchainHookMock) CurrentTimeStamp() uint64 {
	return b.CTimeStamp
}

// CurrentRandomSeed -
func (b *BlockchainHookMock) CurrentRandomSeed() []byte {
	return b.CRandomSeed
}

// CurrentEpoch -
func (b *BlockchainHookMock) CurrentEpoch() uint32 {
	return b.CEpoch
}

// ProcessBuiltInFunction -
func (b *BlockchainHookMock) ProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	return &vmcommon.VMOutput{}, nil
}

// GetBuiltinFunctionNames -
func (b *BlockchainHookMock) GetBuiltinFunctionNames() vmcommon.FunctionNames {
	return make(vmcommon.FunctionNames)
}

// GetAllState -
func (b *BlockchainHookMock) GetAllState(address []byte) (map[string][]byte, error) {
	account, ok := b.Accounts[toHex(address)]
	if !ok {
		return nil, ErrAccountDoesntExist
	}

	allState := make(map[string][]byte)
	for key, value := range account.Storage {
		keyAsBytes, err := fromHex(key)
		if err != nil {
			return nil, err
		}

		valueAsBytes, err := fromHex(value)
		if err != nil {
			return nil, err
		}

		allState[string(keyAsBytes)] = valueAsBytes
	}

	return allState, nil
}

// GetUserAccount -
func (b *BlockchainHookMock) GetUserAccount(address []byte) (vmcommon.UserAccountHandler, error) {
	account, ok := b.Accounts[toHex(address)]
	if !ok {
		return nil, ErrAccountDoesntExist
	}

	return account, nil
}

// GetShardOfAddress -
func (b *BlockchainHookMock) GetShardOfAddress(address []byte) uint32 {
	account, ok := b.Accounts[toHex(address)]
	if !ok {
		return 0
	}

	return account.ShardID
}

// IsSmartContract -
func (b *BlockchainHookMock) IsSmartContract(address []byte) bool {
	account, ok := b.Accounts[toHex(address)]
	if !ok {
		return false
	}

	return len(account.CodeHex) > 0
}

// IsSmartContract -
func (b *BlockchainHookMock) IsPayable(address []byte) (bool, error) {
	account, ok := b.Accounts[toHex(address)]
	if !ok {
		return true, nil
	}

	if !b.IsSmartContract(address) {
		return true, nil
	}

	metadata := vmcommon.CodeMetadataFromBytes(account.GetCodeMetadata())
	return metadata.Payable, nil
}

// UpdateAccounts -
func (b *BlockchainHookMock) UpdateAccounts(outputAccounts map[string]*vmcommon.OutputAccount) {
	for address, outputAccount := range outputAccounts {
		addressHex := toHex([]byte(address))
		account, exists := b.Accounts[addressHex]
		if !exists {
			account = NewAccount(outputAccount.Address, 0, nil)
		}

		account.Balance.Add(account.Balance, outputAccount.BalanceDelta)

		if outputAccount.Nonce > account.Nonce {
			account.Nonce = outputAccount.Nonce
		}
		if len(outputAccount.Code) > 0 {
			account.CodeHex = toHex(outputAccount.Code)

		}
		if len(outputAccount.CodeMetadata) > 0 {
			account.CodeMetadataHex = toHex(outputAccount.CodeMetadata)
		}

		mergeStorageUpdates(account, outputAccount)
		b.Accounts[addressHex] = account
	}
}

func mergeStorageUpdates(leftAccount *Account, rightAccount *vmcommon.OutputAccount) {
	for key, update := range rightAccount.StorageUpdates {
		leftAccount.Storage[toHex([]byte(key))] = toHex(update.Data)
	}
}

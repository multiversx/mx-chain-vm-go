package arwendebug

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var ErrAccountDoesntExist = errors.New("account does not exist")

//var ErrCantDetermineAccountExists = errors.New("can't determine whether account exists")

var zero = big.NewInt(0)

var _ vmcommon.BlockchainHook = (*BlockchainHookMock)(nil)

// Account holds the account info
type Account struct {
	Address []byte
	Nonce   uint64
	Balance *big.Int
	Storage map[string][]byte
	Code    []byte
}

// AccountMap is a map from address to account
type AccountsMap map[string]*Account

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

func NewBlockchainHookMock() *BlockchainHookMock {
	return &BlockchainHookMock{
		Accounts: make(AccountsMap),
	}
}

func (b *BlockchainHookMock) AddAccount(account *Account) {
	if account.Storage == nil {
		account.Storage = make(map[string][]byte)
	}
	if account.Balance == nil {
		account.Balance = big.NewInt(0)
	}
	b.Accounts[string(account.Address)] = account
}

func (b *BlockchainHookMock) AddAccounts(accounts []*Account) {
	for _, account := range accounts {
		b.AddAccount(account)
	}
}

func (b *BlockchainHookMock) AccountExists(address []byte) (bool, error) {
	_, ok := b.Accounts[string(address)]
	if !ok {
		return false, nil
	}

	return true, nil
}

func (b *BlockchainHookMock) NewAddress(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
	if len(creatorAddress) == 0 {
		panic("mock: bad creator address")
	}

	address := make([]byte, arwen.AddressLen)
	copy(address, creatorAddress)
	copy(address, []byte("contract"))
	copy(address[len("contract"):], strconv.Itoa(int(creatorNonce)))
	b.LastCreatedContractAddress = address
	return address, nil
}

func (b *BlockchainHookMock) GetBalance(address []byte) (*big.Int, error) {
	account, ok := b.Accounts[string(address)]
	if !ok {
		return nil, ErrAccountDoesntExist
	}

	return account.Balance, nil
}

func (b *BlockchainHookMock) GetNonce(address []byte) (uint64, error) {
	account, ok := b.Accounts[string(address)]
	if !ok {
		return 0, ErrAccountDoesntExist
	}

	return account.Nonce, nil
}

func (b *BlockchainHookMock) GetStorageData(address []byte, index []byte) ([]byte, error) {
	account, ok := b.Accounts[string(address)]
	if !ok {
		return []byte{}, ErrAccountDoesntExist
	}

	return account.Storage[string(index)], nil
}

func (b *BlockchainHookMock) IsCodeEmpty(address []byte) (bool, error) {
	account, ok := b.Accounts[string(address)]
	if !ok {
		return false, ErrAccountDoesntExist
	}

	empty := len(account.Code) == 0
	return empty, nil
}

func (b *BlockchainHookMock) GetCode(address []byte) ([]byte, error) {
	account, ok := b.Accounts[string(address)]
	if !ok {
		return []byte{}, ErrAccountDoesntExist
	}

	return account.Code, nil
}

func (b *BlockchainHookMock) GetBlockhash(nonce uint64) ([]byte, error) {
	return b.BlockHash, nil
}

func (b *BlockchainHookMock) LastNonce() uint64 {
	return b.LNonce
}

func (b *BlockchainHookMock) LastRound() uint64 {
	return b.LRound
}

func (b *BlockchainHookMock) LastTimeStamp() uint64 {
	return b.LTimeStamp
}

func (b *BlockchainHookMock) LastRandomSeed() []byte {
	return b.LRandomSeed
}

func (b *BlockchainHookMock) LastEpoch() uint32 {
	return b.LEpoch
}

func (b *BlockchainHookMock) GetStateRootHash() []byte {
	return b.StateRootHash
}

func (b *BlockchainHookMock) CurrentNonce() uint64 {
	return b.CNonce
}

func (b *BlockchainHookMock) CurrentRound() uint64 {
	return b.CRound
}

func (b *BlockchainHookMock) CurrentTimeStamp() uint64 {
	return b.CTimeStamp
}

func (b *BlockchainHookMock) CurrentRandomSeed() []byte {
	return b.CRandomSeed
}

func (b *BlockchainHookMock) CurrentEpoch() uint32 {
	return b.CEpoch
}

func (b *BlockchainHookMock) ProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*big.Int, uint64, error) {
	return b.Value, b.Gas, nil
}

func (b *BlockchainHookMock) GetBuiltinFunctionNames() vmcommon.FunctionNames {
	return make(vmcommon.FunctionNames)
}

func (b *BlockchainHookMock) UpdateAccounts(outputAccounts map[string]*vmcommon.OutputAccount) {
	for strAddress, outputAccount := range outputAccounts {
		account, exists := b.Accounts[strAddress]
		if !exists {
			account = &Account{
				Address: outputAccount.Address,
				Balance: big.NewInt(0),
				Code:    nil,
				Storage: make(map[string][]byte),
				Nonce:   0,
			}
		}

		if outputAccount.Nonce > account.Nonce {
			account.Nonce = outputAccount.Nonce
		}
		account.Balance.Add(account.Balance, outputAccount.BalanceDelta)
		if len(outputAccount.Code) > 0 {
			account.Code = outputAccount.Code
		}

		mergeStorageUpdates(account, outputAccount)
		b.Accounts[strAddress] = account
	}
}

func mergeStorageUpdates(
	leftAccount *Account,
	rightAccount *vmcommon.OutputAccount,
) {
	if leftAccount.Storage == nil {
		leftAccount.Storage = make(map[string][]byte)
	}
	for key, update := range rightAccount.StorageUpdates {
		leftAccount.Storage[key] = update.Data
	}
}

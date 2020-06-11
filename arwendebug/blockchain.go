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
func (blockchain *BlockchainHookMock) AddAccount(account *Account) {
	blockchain.Accounts[account.AddressHex] = account
}

// AddAccounts -
func (blockchain *BlockchainHookMock) AddAccounts(accounts []*Account) {
	for _, account := range accounts {
		blockchain.AddAccount(account)
	}
}

// AccountExists -
func (blockchain *BlockchainHookMock) AccountExists(address []byte) (bool, error) {
	_, ok := blockchain.Accounts[toHex(address)]
	if !ok {
		return false, nil
	}

	return true, nil
}

// NewAddress -
func (blockchain *BlockchainHookMock) NewAddress(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
	if len(creatorAddress) != arwen.AddressLen {
		panic("mock: bad creator address")
	}

	address := make([]byte, arwen.AddressLen)
	copy(address, creatorAddress)
	copy(address, []byte("contract"))
	copy(address[len("contract"):], strconv.Itoa(int(creatorNonce)))
	blockchain.LastCreatedContractAddress = address
	return address, nil
}

// GetBalance -
func (blockchain *BlockchainHookMock) GetBalance(address []byte) (*big.Int, error) {
	account, ok := blockchain.Accounts[toHex(address)]
	if !ok {
		return nil, errAccountDoesntExist
	}

	return account.Balance, nil
}

// GetNonce -
func (blockchain *BlockchainHookMock) GetNonce(address []byte) (uint64, error) {
	account, ok := blockchain.Accounts[toHex(address)]
	if !ok {
		return 0, errAccountDoesntExist
	}

	return account.Nonce, nil
}

// GetStorageData -
func (blockchain *BlockchainHookMock) GetStorageData(address []byte, index []byte) ([]byte, error) {
	account, ok := blockchain.Accounts[toHex(address)]
	if !ok {
		return []byte{}, errAccountDoesntExist
	}

	return fromHex(account.Storage[toHex(index)])
}

// IsCodeEmpty -
func (blockchain *BlockchainHookMock) IsCodeEmpty(address []byte) (bool, error) {
	account, ok := blockchain.Accounts[toHex(address)]
	if !ok {
		return false, errAccountDoesntExist
	}

	empty := len(account.CodeHex) == 0
	return empty, nil
}

// GetCode -
func (blockchain *BlockchainHookMock) GetCode(address []byte) ([]byte, error) {
	account, ok := blockchain.Accounts[toHex(address)]
	if !ok {
		return []byte{}, errAccountDoesntExist
	}

	return fromHex(account.CodeHex)
}

// GetBlockhash -
func (blockchain *BlockchainHookMock) GetBlockhash(nonce uint64) ([]byte, error) {
	return blockchain.BlockHash, nil
}

// LastNonce -
func (blockchain *BlockchainHookMock) LastNonce() uint64 {
	return blockchain.LNonce
}

// LastRound -
func (blockchain *BlockchainHookMock) LastRound() uint64 {
	return blockchain.LRound
}

// LastTimeStamp -
func (blockchain *BlockchainHookMock) LastTimeStamp() uint64 {
	return blockchain.LTimeStamp
}

// LastRandomSeed -
func (blockchain *BlockchainHookMock) LastRandomSeed() []byte {
	return blockchain.LRandomSeed
}

// LastEpoch -
func (blockchain *BlockchainHookMock) LastEpoch() uint32 {
	return blockchain.LEpoch
}

// GetStateRootHash -
func (blockchain *BlockchainHookMock) GetStateRootHash() []byte {
	return blockchain.StateRootHash
}

// CurrentNonce -
func (blockchain *BlockchainHookMock) CurrentNonce() uint64 {
	return blockchain.CNonce
}

// CurrentRound -
func (blockchain *BlockchainHookMock) CurrentRound() uint64 {
	return blockchain.CRound
}

// CurrentTimeStamp -
func (blockchain *BlockchainHookMock) CurrentTimeStamp() uint64 {
	return blockchain.CTimeStamp
}

// CurrentRandomSeed -
func (blockchain *BlockchainHookMock) CurrentRandomSeed() []byte {
	return blockchain.CRandomSeed
}

// CurrentEpoch -
func (blockchain *BlockchainHookMock) CurrentEpoch() uint32 {
	return blockchain.CEpoch
}

// ProcessBuiltInFunction -
func (blockchain *BlockchainHookMock) ProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	return &vmcommon.VMOutput{}, nil
}

// GetBuiltinFunctionNames -
func (blockchain *BlockchainHookMock) GetBuiltinFunctionNames() vmcommon.FunctionNames {
	return make(vmcommon.FunctionNames)
}

// GetAllState -
func (blockchain *BlockchainHookMock) GetAllState(address []byte) (map[string][]byte, error) {
	account, ok := blockchain.Accounts[toHex(address)]
	if !ok {
		return nil, errAccountDoesntExist
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

// UpdateAccounts -
func (blockchain *BlockchainHookMock) UpdateAccounts(outputAccounts map[string]*vmcommon.OutputAccount) {
	for address, outputAccount := range outputAccounts {
		addressHex := toHex([]byte(address))
		account, exists := blockchain.Accounts[addressHex]
		if !exists {
			account = NewAccount(outputAccount.Address, 0, nil)
		}

		if outputAccount.Nonce > account.Nonce {
			account.Nonce = outputAccount.Nonce
		}
		account.Balance.Add(account.Balance, outputAccount.BalanceDelta)
		if len(outputAccount.Code) > 0 {
			account.CodeHex = toHex(outputAccount.Code)
		}

		mergeStorageUpdates(account, outputAccount)
		blockchain.Accounts[addressHex] = account
	}
}

func mergeStorageUpdates(leftAccount *Account, rightAccount *vmcommon.OutputAccount) {
	for key, update := range rightAccount.StorageUpdates {
		leftAccount.Storage[toHex([]byte(key))] = toHex(update.Data)
	}
}

package worldmock

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/crypto/hashing"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/ElrondNetwork/elrond-go/data"
	"github.com/ElrondNetwork/elrond-go/data/state"
)

var ErrOperationNotPermitted = errors.New("operation not permitted")
var ErrInvalidAddressLength = errors.New("invalid address length")

var _ state.AccountHandler = (*Account)(nil)
var _ state.DataTrieTracker = (*Account)(nil)

// Account holds the account info
type Account struct {
	Exists          bool
	Address         []byte
	Nonce           uint64
	Balance         *big.Int
	BalanceDelta    *big.Int
	Storage         map[string][]byte
	RootHash        []byte
	Code            []byte
	CodeHash        []byte
	CodeMetadata    []byte
	AsyncCallData   string
	OwnerAddress    []byte
	Username        []byte
	DeveloperReward *big.Int
	ShardID         uint32
	IsSmartContract bool
	ESDTData        map[string]*ESDTData
	TrieTracker     state.DataTrieTracker
}

// ESDTData models an account holding an ESDT token
type ESDTData struct {
	Balance      *big.Int
	BalanceDelta *big.Int
	Frozen       bool
}

var storageDefaultValue = []byte{}

// StorageValue yields the storage value for key, default 0
func (a *Account) StorageValue(key string) []byte {
	value, found := a.Storage[key]
	if !found {
		return storageDefaultValue
	}
	return value
}

// SetCodeAndMetadata changes the account code, as well as all fields depending on it:
// CodeHash, IsSmartContract, CodeMetadata.
// The code metadata must be given explicitly.
func (a *Account) SetCodeAndMetadata(code []byte, codeMetadata *vmcommon.CodeMetadata) {
	a.Code = code
	hasher := hashing.NewHasher()
	hash, err := hasher.Sha256(code)
	if err != nil {
		logger.GetOrCreate("worldAccount").Trace("Account.SetCodeAndMetadata", "error", err)
	}

	a.CodeHash = hash
	a.IsSmartContract = true
	a.CodeMetadata = codeMetadata.ToBytes()
}

func (a *Account) SetPayable(payable bool) {
	if a.CodeMetadata == nil {
		a.CodeMetadata = []byte{0, 0}
	}

	if payable {
		a.CodeMetadata[1] |= vmcommon.MetadataPayable
	} else {
		a.CodeMetadata[1] &^= vmcommon.MetadataPayable
	}
}

// AddressBytes -
func (a *Account) AddressBytes() []byte {
	return a.Address
}

// GetNonce -
func (a *Account) GetNonce() uint64 {
	return a.Nonce
}

// GetCode -
func (a *Account) GetCode() []byte {
	return a.Code
}

// GetCodeMetadata -
func (a *Account) GetCodeMetadata() []byte {
	return a.CodeMetadata
}

// GetCodeHash -
func (a *Account) GetCodeHash() []byte {
	return a.CodeHash
}

// GetRootHash -
func (a *Account) GetRootHash() []byte {
	return a.RootHash
}

// GetBalance -
func (a *Account) GetBalance() *big.Int {
	return a.Balance
}

// SetBalance -
func (a *Account) SetBalance(balance int64) {
	a.Balance = big.NewInt(balance)
}

// GetDeveloperReward -
func (a *Account) GetDeveloperReward() *big.Int {
	return big.NewInt(0)
}

// GetOwnerAddress -
func (a *Account) GetOwnerAddress() []byte {
	return a.OwnerAddress
}

// GetUserName -
func (a *Account) GetUserName() []byte {
	return a.Username
}

// IsInterfaceNil -
func (a *Account) IsInterfaceNil() bool {
	return a == nil
}

func (a *Account) SetCode(code []byte) {
	a.Code = code
	hasher := hashing.NewHasher()
	a.CodeHash, _ = hasher.Sha256(code)
	a.IsSmartContract = true
}

func (a *Account) SetCodeMetadata(codeMetadata []byte) {
	a.CodeMetadata = codeMetadata
}

func (a *Account) SetCodeHash(hash []byte) {
	a.CodeHash = hash
}

func (a *Account) SetRootHash(hash []byte) {
	a.RootHash = hash
}

func (a *Account) DataTrieTracker() state.DataTrieTracker {
	return a
}

func (a *Account) AddToBalance(value *big.Int) error {
	newBalance := big.NewInt(0).Add(a.Balance, value)
	if newBalance.Cmp(zero) < 0 {
		return ErrInsufficientFunds
	}

	a.Balance = newBalance
	return nil
}

func (a *Account) SubFromBalance(value *big.Int) error {
	newBalance := big.NewInt(0).Sub(a.Balance, value)
	if newBalance.Cmp(zero) < 0 {
		return ErrInsufficientFunds
	}

	a.Balance = newBalance
	return nil
}

func (a *Account) ClaimDeveloperRewards(sender []byte) (*big.Int, error) {
	if !bytes.Equal(sender, a.OwnerAddress) {
		return nil, ErrOperationNotPermitted
	}

	oldValue := big.NewInt(0).Set(a.DeveloperReward)
	a.DeveloperReward = big.NewInt(0)

	return oldValue, nil
}

func (a *Account) AddToDeveloperReward(value *big.Int) {
	a.DeveloperReward = big.NewInt(0).Add(a.DeveloperReward, value)
}

func (a *Account) ChangeOwnerAddress(sender []byte, newAddress []byte) error {
	if !bytes.Equal(sender, a.OwnerAddress) {
		return ErrOperationNotPermitted
	}
	if len(newAddress) != len(a.Address) {
		return ErrInvalidAddressLength
	}

	a.OwnerAddress = newAddress

	return nil
}

func (a *Account) SetOwnerAddress(address []byte) {
	a.OwnerAddress = address
}

func (a *Account) SetUserName(userName []byte) {
	a.Username = make([]byte, len(userName))
	copy(a.Username, userName)
}

func (a *Account) IncreaseNonce(nonce uint64) {
	a.Nonce = a.Nonce + nonce
}

func (a *Account) RetrieveValue(key []byte) ([]byte, error) {
	return a.Storage[string(key)], nil
}

func (a *Account) SaveKeyValue(key []byte, value []byte) error {
	a.Storage[string(key)] = value
	return nil
}

func (a *Account) SetDataTrie(tr data.Trie) {
}

func (a *Account) DataTrie() data.Trie {
	return nil
}

func (a *Account) ClearDataCaches() {
}

func (a *Account) DirtyData() map[string][]byte {
	return a.Storage
}

package callbackblockchain

import (
	"math/big"
)

// AccountMap is a map from address to account
type AccountMap map[string]*Account

// Account holds the account info
type Account struct {
	Exists          bool
	Address         []byte
	Nonce           uint64
	Balance         *big.Int
	BalanceDelta    *big.Int
	Storage         map[string][]byte
	Code            []byte
	CodeMetadata    []byte
	AsyncCallData   string
	OwnerAddress    []byte
	Username        []byte
	ShardID         uint32
	IsSmartContract bool
}

var storageDefaultValue = []byte{}

// NewAccountMap creates a new AccountMap instance
func NewAccountMap() AccountMap {
	return AccountMap(make(map[string]*Account))
}

// PutAccount inserts account based on address
func (am AccountMap) PutAccount(acct *Account) {
	mp := (map[string]*Account)(am)
	mp[addressKey(acct.Address)] = acct
}

// GetAccount retrieves account based on address
func (am AccountMap) GetAccount(address []byte) *Account {
	mp := (map[string]*Account)(am)
	return mp[addressKey(address)]
}

// DeleteAccount removes account based on address
func (am AccountMap) DeleteAccount(address []byte) {
	mp := (map[string]*Account)(am)
	delete(mp, addressKey(address))
}

func addressKey(address []byte) string {
	return string(address)
}

// StorageKey builds a key for the mock
func StorageKey(address []byte) string {
	//return string(big.NewInt(0).SetBytes(address).Bytes())
	return string(address)
}

// StorageValue yields the storage value for key, default 0
func (a *Account) StorageValue(key string) []byte {
	value, found := a.Storage[key]
	if !found {
		return storageDefaultValue
	}
	return value
}

// AccountAddress converts to account address bytes from big.Int
func AccountAddress(i *big.Int) []byte {
	if i.Sign() < 0 {
		panic("address cannot be negative")
	}
	return i.Bytes()
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
	return []byte{}
}

// GetRootHash -
func (a *Account) GetRootHash() []byte {
	return []byte{}
}

// GetBalance -
func (a *Account) GetBalance() *big.Int {
	return a.Balance
}

// GetDeveloperReward -
func (a *Account) GetDeveloperReward() *big.Int {
	return big.NewInt(0)
}

// GetOwnerAddress -
func (a *Account) GetOwnerAddress() []byte {
	return a.OwnerAddress
}

// GetOwnerAddress -
func (a *Account) GetUserName() []byte {
	return a.Username
}

// IsInterfaceNil -
func (a *Account) IsInterfaceNil() bool {
	return a == nil
}

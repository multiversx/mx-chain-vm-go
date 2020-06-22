package mock

import (
	"math/big"
)

// Account holds the account info
type Account struct {
	Address      []byte
	Nonce        uint64
	Balance      *big.Int
	Storage      map[string][]byte
	Code         []byte
	CodeMetadata []byte
	OwnerAddress []byte
	UserName     []byte
	ShardID      uint32
	Err          error
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
	if a.Balance == nil {
		return big.NewInt(0)
	}
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
	return a.UserName
}

// IsInterfaceNil -
func (a *Account) IsInterfaceNil() bool {
	return a == nil
}

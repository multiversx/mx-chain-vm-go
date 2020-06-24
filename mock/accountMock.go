package mock

import (
	"math/big"
)

// AccountMock holds the account info
type AccountMock struct {
	Address      []byte
	Nonce        uint64
	Balance      *big.Int
	Code         []byte
	CodeMetadata []byte
	OwnerAddress []byte
	UserName     []byte

	ShardID uint32
	Storage map[string][]byte
	Err     error
}

// AddressBytes -
func (a *AccountMock) AddressBytes() []byte {
	return a.Address
}

// GetNonce -
func (a *AccountMock) GetNonce() uint64 {
	return a.Nonce
}

// GetCode -
func (a *AccountMock) GetCode() []byte {
	return a.Code
}

// GetCodeMetadata -
func (a *AccountMock) GetCodeMetadata() []byte {
	return a.CodeMetadata
}

// GetCodeHash -
func (a *AccountMock) GetCodeHash() []byte {
	return []byte{}
}

// GetRootHash -
func (a *AccountMock) GetRootHash() []byte {
	return []byte{}
}

// GetBalance -
func (a *AccountMock) GetBalance() *big.Int {
	if a.Balance == nil {
		return big.NewInt(0)
	}
	return a.Balance
}

// GetDeveloperReward -
func (a *AccountMock) GetDeveloperReward() *big.Int {
	return big.NewInt(0)
}

// GetOwnerAddress -
func (a *AccountMock) GetOwnerAddress() []byte {
	return a.OwnerAddress
}

// GetOwnerAddress -
func (a *AccountMock) GetUserName() []byte {
	return a.UserName
}

// IsInterfaceNil -
func (a *AccountMock) IsInterfaceNil() bool {
	return a == nil
}

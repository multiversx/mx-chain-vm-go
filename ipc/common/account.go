package common

import (
	"math/big"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// Account holds the account info (is a substructure of an IPC message)
type Account struct {
	Nonce           uint64
	Balance         *big.Int
	CodeHash        []byte
	RootHash        []byte
	Address         []byte
	DeveloperReward *big.Int
	OwnerAddress    []byte
	UserName        []byte
	CodeMetadata    []byte
}

// AddressBytes gets the address
func (a *Account) AddressBytes() []byte {
	return a.Address
}

// GetNonce gets the nonce
func (a *Account) GetNonce() uint64 {
	return a.Nonce
}

// GetCodeMetadata gets the code metadata
func (a *Account) GetCodeMetadata() []byte {
	return a.CodeMetadata
}

// GetCodeHash gets the code hash
func (a *Account) GetCodeHash() []byte {
	return a.CodeHash
}

// GetRootHash gets the root hash
func (a *Account) GetRootHash() []byte {
	return a.RootHash
}

// GetBalance gets the balance
func (a *Account) GetBalance() *big.Int {
	if a.Balance == nil {
		return big.NewInt(0)
	}
	return a.Balance
}

// GetDeveloperReward gets the developer reward
func (a *Account) GetDeveloperReward() *big.Int {
	if a.DeveloperReward == nil {
		return big.NewInt(0)
	}
	return a.DeveloperReward
}

// GetOwnerAddress gets the owner's address
func (a *Account) GetOwnerAddress() []byte {
	return a.OwnerAddress
}

// GetOwnerAddress gets the username
func (a *Account) GetUserName() []byte {
	return a.UserName
}

// DataTrieTracker -
func (a *Account) AccountDataHandler() vmcommon.AccountDataHandler {
	return nil
}

// AddToBalance -
func (a *Account) AddToBalance(_ *big.Int) error {
	return nil
}

// ClaimDeveloperRewards -
func (a *Account) ClaimDeveloperRewards(_ []byte) (*big.Int, error) {
	return big.NewInt(0), nil
}

// ChangeOwnerAddress -
func (a *Account) ChangeOwnerAddress(_ []byte, _ []byte) error {
	return nil
}

// SetOwnerAddress -
func (a *Account) SetOwnerAddress(_ []byte) {
}

// SetUserName -
func (a *Account) SetUserName(_ []byte) {
}

// IncreaseNonce -
func (a *Account) IncreaseNonce(_ uint64) {
}

// IsInterfaceNil returns true if there is no value under the interface
func (a *Account) IsInterfaceNil() bool {
	return a == nil
}

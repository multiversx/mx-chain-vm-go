package mock

import (
	"math/big"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// StubAccount is used with the blockchain stub in vm context tests.
// It has minimal functionality.
type StubAccount struct {
	Address      []byte
	Nonce        uint64
	Balance      *big.Int
	CodeHash     []byte
	CodeMetadata []byte
	OwnerAddress []byte
	UserName     []byte

	ShardID uint32
	Storage map[string][]byte
	Err     error
}

// AccountDataHandler -
func (a *StubAccount) AccountDataHandler() vmcommon.AccountDataHandler {
	panic("implement me")
}

// AddToBalance -
func (a *StubAccount) AddToBalance(_ *big.Int) error {
	panic("implement me")
}

// SubFromBalance -
func (a *StubAccount) SubFromBalance(_ *big.Int) error {
	panic("implement me")
}

// ClaimDeveloperRewards -
func (a *StubAccount) ClaimDeveloperRewards(_ []byte) (*big.Int, error) {
	panic("implement me")
}

// ChangeOwnerAddress -
func (a *StubAccount) ChangeOwnerAddress(_ []byte, _ []byte) error {
	panic("implement me")
}

// SetOwnerAddress -
func (a *StubAccount) SetOwnerAddress(_ []byte) {
	panic("implement me")
}

// SetUserName -
func (a *StubAccount) SetUserName(_ []byte) {
	panic("implement me")
}

// IncreaseNonce -
func (a *StubAccount) IncreaseNonce(_ uint64) {
	panic("implement me")
}

// AddressBytes -
func (a *StubAccount) AddressBytes() []byte {
	return a.Address
}

// GetNonce -
func (a *StubAccount) GetNonce() uint64 {
	return a.Nonce
}

// GetCodeMetadata -
func (a *StubAccount) GetCodeMetadata() []byte {
	return a.CodeMetadata
}

// GetCodeHash -
func (a *StubAccount) GetCodeHash() []byte {
	return a.CodeHash
}

// GetRootHash -
func (a *StubAccount) GetRootHash() []byte {
	return []byte{}
}

// GetBalance -
func (a *StubAccount) GetBalance() *big.Int {
	if a.Balance == nil {
		return big.NewInt(0)
	}
	return a.Balance
}

// GetDeveloperReward -
func (a *StubAccount) GetDeveloperReward() *big.Int {
	return big.NewInt(0)
}

// GetOwnerAddress -
func (a *StubAccount) GetOwnerAddress() []byte {
	return a.OwnerAddress
}

// GetUserName -
func (a *StubAccount) GetUserName() []byte {
	return a.UserName
}

// SetCodeMetadata -
func (a *StubAccount) SetCodeMetadata(_ []byte) {
}

// IsInterfaceNil -
func (a *StubAccount) IsInterfaceNil() bool {
	return a == nil
}

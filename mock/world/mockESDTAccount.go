package worldmock

import (
	"math/big"

	"github.com/multiversx/mx-chain-core-go/data/esdt"
	"github.com/multiversx/mx-chain-vm-v1_4-go/mandos-go/esdtconvert"
)

// GetTokenBalance returns the ESDT balance of the account, specified by the
// token key.
func (a *Account) GetTokenBalance(tokenIdentifier []byte, nonce uint64) (*big.Int, error) {
	return esdtconvert.GetTokenBalance(tokenIdentifier, nonce, a.Storage)
}

// GetTokenBalanceUint64 returns the ESDT balance of the account, specified by the
// token identifier.
func (a *Account) GetTokenBalanceUint64(tokenIdentifier []byte, nonce uint64) (uint64, error) {
	balance, err := a.GetTokenBalance(tokenIdentifier, nonce)
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}

// SetTokenBalance sets the ESDT balance of the account, specified by the token
// key.
func (a *Account) SetTokenBalance(tokenIdentifier []byte, nonce uint64, balance *big.Int) error {
	return esdtconvert.SetTokenBalance(tokenIdentifier, nonce, balance, a.Storage)
}

// SetTokenBalanceUint64 sets the ESDT balance of the account, specified by the
// token key.
func (a *Account) SetTokenBalanceUint64(tokenIdentifier []byte, nonce uint64, balance uint64) error {
	return esdtconvert.SetTokenBalance(tokenIdentifier, nonce, big.NewInt(0).SetUint64(balance), a.Storage)
}

// GetTokenData gets the ESDT information related to a token from the storage of the account.
func (a *Account) GetTokenData(tokenIdentifier []byte, nonce uint64, systemAccStorage map[string][]byte) (*esdt.ESDigitalToken, error) {
	return esdtconvert.GetTokenData(tokenIdentifier, nonce, a.Storage, systemAccStorage)
}

// SetTokenData sets the ESDT information related to a token into the storage of the account.
func (a *Account) SetTokenData(tokenIdentifier []byte, nonce uint64, tokenData *esdt.ESDigitalToken) error {
	return esdtconvert.SetTokenData(tokenIdentifier, nonce, tokenData, a.Storage)
}

// SetTokenRolesAsStrings sets the specified roles to the account, corresponding to the given tokenName.
func (a *Account) SetTokenRolesAsStrings(tokenIdentifier []byte, rolesAsStrings []string) error {
	return esdtconvert.SetTokenRolesAsStrings(tokenIdentifier, rolesAsStrings, a.Storage)
}

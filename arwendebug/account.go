package arwendebug

import (
	"math/big"
)

// AccountsMap is a map from address to account
type AccountsMap map[string]*Account

// Account holds the account info
type Account struct {
	AddressHex string
	Nonce      uint64
	Balance    *big.Int
	Code       []byte
	Storage    map[string][]byte
}

// NewAccount creates a new debug account
func NewAccount(address []byte, nonce uint64, balance *big.Int) *Account {
	if balance == nil {
		balance = big.NewInt(0)
	}

	return &Account{
		AddressHex: toHex(address),
		Nonce:      nonce,
		Balance:    balance,
		Code:       nil,
		Storage:    make(map[string][]byte),
	}
}

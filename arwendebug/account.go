package arwendebug

import (
	"math/big"
)

// AccountsMap is a map from address to account
type AccountsMap map[string]*Account

// Account is a debug account
type Account struct {
	AddressHex string
	Nonce      uint64
	Balance    *big.Int
	CodeHex    string
	Storage    map[string]string
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
		CodeHex:    "",
		Storage:    make(map[string]string),
	}
}

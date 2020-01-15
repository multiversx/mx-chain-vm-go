package contexts

import (
	"math/big"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type outputAccount struct {
	Address        []byte
	Nonce          uint64
	BalanceDelta   *big.Int
	StorageUpdates map[string]*vmcommon.StorageUpdate
	Code           []byte
	Data           []byte
	GasLimit       uint64
}

func (account *outputAccount) update(otherAccount *outputAccount) {
	account.Address = otherAccount.Address
	account.updateBalanceDelta(otherAccount.BalanceDelta)
	account.updateNonce(otherAccount.Nonce)
	account.updateStorage(otherAccount.StorageUpdates)
	account.Code = otherAccount.Code
	account.Data = otherAccount.Data
	account.GasLimit = otherAccount.GasLimit
}

func (account *outputAccount) updateBalanceDelta(delta *big.Int) {
	if account.BalanceDelta == nil {
		account.BalanceDelta = big.NewInt(0)
	}
	if delta != nil {
		account.BalanceDelta = big.NewInt(0).Add(account.BalanceDelta, delta)
	}
}

func (account *outputAccount) updateNonce(nonce uint64) {
	if nonce > account.Nonce {
		account.Nonce = nonce
	}
}

func (account *outputAccount) updateStorage(updates map[string]*vmcommon.StorageUpdate) {
	if account.StorageUpdates == nil {
		account.StorageUpdates = make(map[string]*vmcommon.StorageUpdate)
	}
	for key, update := range updates {
		account.StorageUpdates[key] = update
	}
}

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

func newOutputAccount(address []byte) *outputAccount {
  return &outputAccount{
    Address: address,
    Nonce: 0,
    BalanceDelta: big.NewInt(0),
    StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
  }
}

func (account *outputAccount) update(otherAccount *outputAccount) {
	account.Address = otherAccount.Address
	account.updateBalanceDelta(otherAccount.BalanceDelta)
	account.updateStorage(otherAccount.StorageUpdates)

	if len(otherAccount.Code) > 0 {
		account.Code = otherAccount.Code
	}
	if len(otherAccount.Data) > 0 {
		account.Data = otherAccount.Data
	}
	if otherAccount.Nonce > account.Nonce {
		account.Nonce = otherAccount.Nonce
	}
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

func (account *outputAccount) updateStorage(updates map[string]*vmcommon.StorageUpdate) {
	if account.StorageUpdates == nil {
		account.StorageUpdates = make(map[string]*vmcommon.StorageUpdate)
	}
	for key, update := range updates {
		account.StorageUpdates[key] = update
	}
}

func (account *outputAccount) AccountForVMOutput() *vmcommon.OutputAccount {
	vmOutputAccount := &vmcommon.OutputAccount{
		Address:        account.Address,
		Nonce:          account.Nonce,
		BalanceDelta:   account.BalanceDelta,
		StorageUpdates: nil,
		Code:           account.Code,
		Data:           account.Data,
		GasLimit:       account.GasLimit,
	}

	vmOutputAccount.StorageUpdates = make([]*vmcommon.StorageUpdate, len(account.StorageUpdates))
	i := 0
	for _, update := range account.StorageUpdates {
		vmOutputAccount.StorageUpdates[i] = update
		i++
	}

	return vmOutputAccount
}

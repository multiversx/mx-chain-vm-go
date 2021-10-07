package elrondgo_exporter

import "math/big"

type testAccount struct {
	nonce   uint64
	address []byte
	balance *big.Int
	storage map[string][]byte
}

func NewTestAccount() *testAccount {
	return &testAccount{
		address: make([]byte, 0),
		balance: big.NewInt(0),
		storage: make(map[string][]byte),
	}
}

func SetNewAccount(nonce uint64, address []byte, balance *big.Int, storage map[string][]byte) *testAccount {
	return NewTestAccount().WithNonce(nonce).WithAddress(address).WithBalance(balance).WithStorage(storage)
}

func (tAcc *testAccount) WithNonce(nonce uint64) *testAccount {
	tAcc.nonce = nonce
	return tAcc
}

func (tAcc *testAccount) WithAddress(address []byte) *testAccount {
	tAcc.address = append(tAcc.address, address...)
	return tAcc
}

func (tAcc *testAccount) WithBalance(balance *big.Int) *testAccount {
	tAcc.balance.Set(balance)
	return tAcc
}

func (tAcc *testAccount) WithStorage(storage map[string][]byte) *testAccount {
	tAcc.storage = storage
	return tAcc
}

func (tAcc *testAccount) GetNonce() uint64 {
	return tAcc.nonce
}

func (tAcc *testAccount) GetAddress() []byte {
	return tAcc.address
}

func (tAcc *testAccount) GetBalance() *big.Int {
	return tAcc.balance
}

func (tAcc *testAccount) GetStorage() map[string][]byte {
	return tAcc.storage
}

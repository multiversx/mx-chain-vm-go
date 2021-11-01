package elrondgo_exporter

import "math/big"

type TestAccount struct {
	nonce        uint64
	address      []byte
	balance      *big.Int
	storage      map[string][]byte
	code         []byte
	ownerAddress []byte
}

func NewTestAccount() *TestAccount {
	return &TestAccount{
		address:      make([]byte, 0),
		balance:      big.NewInt(0),
		storage:      make(map[string][]byte),
		code:         make([]byte, 0),
		ownerAddress: make([]byte, 0),
	}
}

func SetNewAccount(nonce uint64, address []byte, balance *big.Int, storage map[string][]byte, code []byte, ownerAddress []byte) *TestAccount {
	return NewTestAccount().WithNonce(nonce).WithAddress(address).WithBalance(balance).WithStorage(storage).WithCode(code).WithOwner(ownerAddress)
}

func (tAcc *TestAccount) WithNonce(nonce uint64) *TestAccount {
	tAcc.nonce = nonce
	return tAcc
}

func (tAcc *TestAccount) WithAddress(address []byte) *TestAccount {
	tAcc.address = make([]byte, 0)
	tAcc.address = address
	return tAcc
}

func (tAcc *TestAccount) WithBalance(balance *big.Int) *TestAccount {
	tAcc.balance.Set(balance)
	return tAcc
}

func (tAcc *TestAccount) WithStorage(storage map[string][]byte) *TestAccount {
	tAcc.storage = storage
	return tAcc
}

func (tAcc *TestAccount) WithCode(code []byte) *TestAccount {
	tAcc.code = append(tAcc.code, code...)
	return tAcc
}

func (tAcc *TestAccount) WithOwner(owner []byte) *TestAccount {
	tAcc.ownerAddress = append(tAcc.ownerAddress, owner...)
	return tAcc
}

func (tAcc *TestAccount) GetNonce() uint64 {
	return tAcc.nonce
}

func (tAcc *TestAccount) GetAddress() []byte {
	return tAcc.address
}

func (tAcc *TestAccount) GetBalance() *big.Int {
	return tAcc.balance
}

func (tAcc *TestAccount) GetStorage() map[string][]byte {
	return tAcc.storage
}

func (tAcc *TestAccount) GetCode() []byte {
	return tAcc.code
}

func (tAcc *TestAccount) GetOwner() []byte {
	return tAcc.ownerAddress
}

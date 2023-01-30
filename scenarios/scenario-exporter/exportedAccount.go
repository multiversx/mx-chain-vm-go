package scenario_exporter

import "math/big"

// TestAccount defines the test account structure
type TestAccount struct {
	nonce        uint64
	address      []byte
	balance      *big.Int
	storage      map[string][]byte
	code         []byte
	ownerAddress []byte
}

// NewTestAccount will create a new instance of type TestAccount
func NewTestAccount() *TestAccount {
	return &TestAccount{
		address:      make([]byte, 0),
		balance:      big.NewInt(0),
		storage:      make(map[string][]byte),
		code:         make([]byte, 0),
		ownerAddress: make([]byte, 0),
	}
}

// SetNewAccount will create a new TestAccount
func SetNewAccount(nonce uint64, address []byte, balance *big.Int, storage map[string][]byte, code []byte, ownerAddress []byte) *TestAccount {
	return NewTestAccount().WithNonce(nonce).WithAddress(address).WithBalance(balance).WithStorage(storage).WithCode(code).WithOwner(ownerAddress)
}

// WithNonce sets the nonce
func (tAcc *TestAccount) WithNonce(nonce uint64) *TestAccount {
	tAcc.nonce = nonce
	return tAcc
}

// WithAddress sets the address
func (tAcc *TestAccount) WithAddress(address []byte) *TestAccount {
	tAcc.address = make([]byte, 0)
	tAcc.address = address
	return tAcc
}

// WithBalance sets the balance
func (tAcc *TestAccount) WithBalance(balance *big.Int) *TestAccount {
	tAcc.balance.Set(balance)
	return tAcc
}

// WithStorage sets the account's storage
func (tAcc *TestAccount) WithStorage(storage map[string][]byte) *TestAccount {
	tAcc.storage = storage
	return tAcc
}

// WithCode sets the account's code
func (tAcc *TestAccount) WithCode(code []byte) *TestAccount {
	tAcc.code = append(tAcc.code, code...)
	return tAcc
}

// WithOwner sets the owner
func (tAcc *TestAccount) WithOwner(owner []byte) *TestAccount {
	tAcc.ownerAddress = append(tAcc.ownerAddress, owner...)
	return tAcc
}

// GetNonce gets the nonce
func (tAcc *TestAccount) GetNonce() uint64 {
	return tAcc.nonce
}

// GetAddress gets the address
func (tAcc *TestAccount) GetAddress() []byte {
	return tAcc.address
}

// GetBalance gets the balance
func (tAcc *TestAccount) GetBalance() *big.Int {
	return tAcc.balance
}

// GetStorage gets the storage
func (tAcc *TestAccount) GetStorage() map[string][]byte {
	return tAcc.storage
}

// GetCode gets the code
func (tAcc *TestAccount) GetCode() []byte {
	return tAcc.code
}

// GetOwner gets the owner
func (tAcc *TestAccount) GetOwner() []byte {
	return tAcc.ownerAddress
}

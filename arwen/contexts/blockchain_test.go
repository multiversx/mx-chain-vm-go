package contexts

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

var ErrAccountFault = errors.New("account fault")
var ErrTestError = errors.New("some test error")

var testAccounts = []*mock.Account{
	{Exists: true, Address: []byte("account_old"), Nonce: 12, Balance: big.NewInt(0)},
	{Exists: true, Address: []byte("account_newer"), Nonce: 8, Balance: big.NewInt(0)},
	{Exists: true, Address: []byte("account_new"), Nonce: 0, Balance: big.NewInt(0)},
	{Exists: true, Address: []byte("account_new_with_money"), Nonce: 0, Balance: big.NewInt(1000)},
	{Exists: true, Address: []byte("account_with_code"), Nonce: 4, Balance: big.NewInt(512), Code: []byte("somecode")},
	{Exists: true, Address: []byte("account_old_with_money"), Nonce: 56, Balance: big.NewInt(1024)},
	{Exists: false, Address: []byte("account_missing")},
	{Exists: true, Address: []byte("account_faulty"), Err: ErrAccountFault},
}

func TestNewBlockchainContext(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostStub{}
	bcHook := mock.NewBlockchainHookMock()

	bcc, err := NewBlockchainContext(host, bcHook)
	require.Nil(t, err)
	require.NotNil(t, bcc)
}

func TestBlockchainContext_AccountExists(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostStub{}
	mockBlockchain := mock.NewBlockchainHookMock()
	mockBlockchain.AddAccounts(testAccounts)

	bcc, _ := NewBlockchainContext(host, mockBlockchain)

	require.False(t, bcc.AccountExists([]byte("account_missing")))
	require.False(t, bcc.AccountExists([]byte("account_faulty")))
	require.True(t, bcc.AccountExists([]byte("account_old")))

	mockBlockchain.Err = ErrTestError
	require.False(t, bcc.AccountExists([]byte("account_something")))
}

func TestBlockchainContext_GetBalance(t *testing.T) {
	t.Parallel()

	mockBlockchain := mock.NewBlockchainHookMock()
	mockBlockchain.AddAccounts(testAccounts)
	mockOutput := &mock.OutputContextMock{}
	host := &mock.VmHostMock{}
	host.OutputContext = mockOutput
	bcc, _ := NewBlockchainContext(host, mockBlockchain)

	// Act as if the OutputContext has no OutputAccounts cached
	// (mockOutput.GetOutputAccount() always returns "is new")
	account := &vmcommon.OutputAccount{}
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = true

	// Test if error is propagated from BlockchainHook
	mockBlockchain.Err = ErrTestError
	balanceBytes := bcc.GetBalance([]byte("account_new_with_money"))
	value := big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, big.NewInt(0), value)
	mockBlockchain.Err = nil

	// Test that an account that doesn't exist will not be updated with any kind
	// of balance in the OutputAccounts, but GetBalance() must return 0
	account.Balance = big.NewInt(15)
	balanceBytes = bcc.GetBalance([]byte("account_missing"))
	value = big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, big.NewInt(0), value)
	require.Equal(t, big.NewInt(15), account.Balance)

	// Test that an account newly cached by the OutputAccounts will have its
	// Balance updated, if BlockchainHook.GetBalance() is successful
	account.Balance = big.NewInt(300)
	balanceBytes = bcc.GetBalance([]byte("account_new_with_money"))
	value = big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, big.NewInt(1000), value)
	require.Equal(t, big.NewInt(1000), account.Balance)

	// Act as if the OutputContext has the requested OutputAccount cached
	account.Balance = big.NewInt(42)
	mockOutput.OutputAccountIsNew = false
	balanceBytes = bcc.GetBalance([]byte("any account"))
	value = big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, big.NewInt(42), value)

}

func TestBlockchainContext_GetNonceAndIncrease(t *testing.T) {
	t.Parallel()

	account := &vmcommon.OutputAccount{
		Nonce: 3,
	}

	host := &mock.VmHostMock{}

	mockOutput := &mock.OutputContextMock{}
	host.OutputContext = mockOutput

	mockBlockchain := mock.NewBlockchainHookMock()
	mockBlockchain.AddAccounts(testAccounts)
	bcc, _ := NewBlockchainContext(host, mockBlockchain)

	// GetNonce: Test if error is propagated from BlockchainHook, and that the
	// cached OutputAccount doesn't lose its Nonce due to the error.
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = true
	mockBlockchain.Err = ErrTestError
	nonce, err := bcc.GetNonce([]byte("any account"))
	require.Equal(t, ErrTestError, err)
	require.Equal(t, uint64(0), nonce)
	require.Equal(t, uint64(3), account.Nonce)
	mockBlockchain.Err = nil

	// GetNonce: Test requesting the nonce of an Account not yet cached by
	// OutputAccounts
	mockOutput.OutputAccountIsNew = true
	nonce, err = bcc.GetNonce([]byte("account_old"))
	require.Equal(t, nil, err)
	require.Equal(t, uint64(12), nonce)
	require.Equal(t, uint64(12), account.Nonce)

	// GetNonce: Test requesting the nonce of an Account already cached by
	// OutputAccounts
	account.Nonce = 88
	mockOutput.OutputAccountIsNew = false
	nonce, err = bcc.GetNonce([]byte("any account"))
	require.Nil(t, err)
	require.Equal(t, uint64(88), nonce)

	// IncreaseNonce: Test increasing the nonce of an Account already cached by
	// OutputAccounts
	account.Nonce = 88
	mockOutput.OutputAccountIsNew = false
	bcc.IncreaseNonce([]byte("any account"))
	bcc.IncreaseNonce([]byte("any account"))
	bcc.IncreaseNonce([]byte("any account"))
	nonce, err = bcc.GetNonce([]byte("any account"))
	require.Nil(t, err)
	require.Equal(t, uint64(91), nonce)
}

func TestBlockchainContext_GetCodeHashAndSize(t *testing.T) {
	t.Parallel()

	mockCrypto := &mock.CryptoHookMock{}

	mockBlockchain := mock.NewBlockchainHookMock()
	mockBlockchain.AddAccounts(testAccounts)

	host := &mock.VmHostMock{}
	host.CryptoHook = mockCrypto

	bcc, _ := NewBlockchainContext(host, mockBlockchain)

	address := []byte("account_with_code")
	expectedCode := []byte("somecode")
	expectedCodeHash := []byte("code hash")

	// GetCode: Test if error is propagated from blockchain hook
	mockBlockchain.Err = ErrTestError
	code, err := bcc.GetCode(address)
	require.NotNil(t, err)
	require.Nil(t, code)
	mockBlockchain.Err = nil

	// GetCode: Test success
	code, err = bcc.GetCode(address)
	require.Nil(t, err)
	require.Equal(t, expectedCode, code)

	// GetCodeHash: Test if error is propagated from blockchain hook
	mockBlockchain.Err = ErrTestError
	codeHash, err := bcc.GetCodeHash(address)
	require.Equal(t, ErrTestError, err)
	require.Nil(t, codeHash)
	mockBlockchain.Err = nil

	// GetCodeHash: Test if error is propagated from crypto hook
	mockCrypto.Result = nil
	mockCrypto.Err = ErrTestError
	codeHash, err = bcc.GetCodeHash(address)
	require.Equal(t, ErrTestError, err)
	require.Nil(t, codeHash)

	// GetCodeHash: Test success
	mockCrypto.Result = expectedCodeHash
	mockCrypto.Err = nil
	codeHash, err = bcc.GetCodeHash(address)
	require.Equal(t, expectedCodeHash, codeHash)
	require.Nil(t, err)

	// GetCodeSize: Test if error is propagated from blockchain hook
	mockBlockchain.Err = ErrTestError
	size, err := bcc.GetCodeSize(address)
	require.Equal(t, ErrTestError, err)
	require.Equal(t, int32(0), size)
	mockBlockchain.Err = nil

	// GetCodeSize: Test success
	size, err = bcc.GetCodeSize(address)
	require.Nil(t, err)
	require.Equal(t, int32(len(expectedCode)), size)
}

func TestBlockchainContext_NewAddress(t *testing.T) {
	t.Parallel()

	mockOutput := &mock.OutputContextMock{}

	mockBlockchain := mock.NewBlockchainHookMock()
	mockBlockchain.AddAccounts(testAccounts)

	mockRuntime := &mock.RuntimeContextMock{}
	mockRuntime.VmType = []byte{0xF, 0xF}

	host := &mock.VmHostMock{
		OutputContext:  mockOutput,
		RuntimeContext: mockRuntime,
	}

	// Test error propagation from GetNonce()
	bcc, _ := NewBlockchainContext(host, mockBlockchain)
	creatorAddress := []byte("account_new")
	creatorAccount := mockBlockchain.Accounts[string(creatorAddress)]
	creatorOutputAccount := mockOutput.NewVMOutputAccountFromMockAccount(creatorAccount)
	mockOutput.OutputAccountMock = creatorOutputAccount
	mockOutput.OutputAccountIsNew = true
	mockBlockchain.Err = ErrTestError

	address, err := bcc.NewAddress(creatorAddress)
	require.Equal(t, ErrTestError, err)
	require.Nil(t, address)

	mockBlockchain.Err = nil

	// Test if nonce is not deducted if 0, before calling BlockchainHook.NewAddres()
	creatorAddress = []byte("account_new")
	creatorAccount = mockBlockchain.Accounts[string(creatorAddress)]
	creatorOutputAccount = mockOutput.NewVMOutputAccountFromMockAccount(creatorAccount)
	mockOutput.OutputAccountMock = creatorOutputAccount
	mockOutput.OutputAccountIsNew = true

	expectedCreatorAddres := creatorAddress
	stubBlockchain := &mock.BlockChainHookStub{
		GetNonceCalled: mockBlockchain.GetNonce,
		NewAddressCalled: func(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
			require.Equal(t, expectedCreatorAddres, creatorAddress)
			require.Equal(t, uint64(0), creatorNonce)
			require.Equal(t, mockRuntime.VmType, vmType)
			return []byte("new_address"), nil
		},
	}
	bcc, _ = NewBlockchainContext(host, stubBlockchain)

	address, err = bcc.NewAddress(creatorAddress)
	require.Nil(t, err)
	require.Equal(t, []byte("new_address"), address)

	// Test if nonce is correctly deducted if greater than 0, before calling BlockchainHook.NewAddres()
	creatorAddress = []byte("account_old_with_money")
	creatorAccount = mockBlockchain.Accounts[string(creatorAddress)]
	creatorOutputAccount = mockOutput.NewVMOutputAccountFromMockAccount(creatorAccount)
	mockOutput.OutputAccountMock = creatorOutputAccount
	mockOutput.OutputAccountIsNew = false

	expectedCreatorAddres = creatorAddress
	stubBlockchain = &mock.BlockChainHookStub{
		GetNonceCalled: mockBlockchain.GetNonce,
		NewAddressCalled: func(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
			require.Equal(t, expectedCreatorAddres, creatorAddress)
			require.Equal(t, uint64(55), creatorNonce)
			require.Equal(t, mockRuntime.VmType, vmType)
			return []byte("new_address"), nil
		},
	}
	bcc, _ = NewBlockchainContext(host, stubBlockchain)

	address, err = bcc.NewAddress(creatorAddress)
	require.Nil(t, err)
	require.Equal(t, []byte("new_address"), address)

	// Test if error is propagated from Blockchain.NewAddress
	creatorAddress = []byte("account_with_code")
	creatorAccount = mockBlockchain.Accounts[string(creatorAddress)]
	creatorOutputAccount = mockOutput.NewVMOutputAccountFromMockAccount(creatorAccount)
	mockOutput.OutputAccountMock = creatorOutputAccount
	mockOutput.OutputAccountIsNew = false

	expectedCreatorAddres = creatorAddress
	stubBlockchain = &mock.BlockChainHookStub{
		GetNonceCalled: mockBlockchain.GetNonce,
		NewAddressCalled: func(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
			require.Equal(t, expectedCreatorAddres, creatorAddress)
			require.Equal(t, uint64(3), creatorNonce)
			require.Equal(t, mockRuntime.VmType, vmType)
			return nil, ErrTestError
		},
	}
	bcc, _ = NewBlockchainContext(host, stubBlockchain)

	address, err = bcc.NewAddress(creatorAddress)
	require.Equal(t, ErrTestError, err)
	require.Nil(t, address)
}

func TestBlockchainContext_BlockHash(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostMock{}
	mockBlockchain := mock.NewBlockchainHookMock()
	bcc, _ := NewBlockchainContext(host, mockBlockchain)

	mockBlockchain.Err = ErrTestError
	hash := bcc.BlockHash(42)
	require.Nil(t, hash)
	mockBlockchain.Err = nil

	mockBlockchain.BlockHash = []byte("1234fa")
	hash = bcc.BlockHash(-5)
	require.Nil(t, hash)

	mockBlockchain.BlockHash = []byte("1234fb")
	hash = bcc.BlockHash(0)
	require.Equal(t, []byte("1234fb"), hash)

	mockBlockchain.BlockHash = []byte("1234fc")
	hash = bcc.BlockHash(42)
	require.Equal(t, []byte("1234fc"), hash)
}

func TestBlockchainContext_Getters(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostMock{}
	mockBlockchain := &mock.BlockchainHookMock{
		LEpoch: 3,
		CEpoch: 4,

		LNonce: 90,
		CNonce: 98,

		LRound: 96,
		CRound: 99,

		LTimeStamp: 6749,
		CTimeStamp: 6800,

		StateRootHash: []byte("root hash"),
		LRandomSeed:   []byte("last random seed"),
		CRandomSeed:   []byte("current random seed"),
	}

	bcc, _ := NewBlockchainContext(host, mockBlockchain)

	require.Equal(t, uint32(3), bcc.LastEpoch())
	require.Equal(t, uint32(4), bcc.CurrentEpoch())

	require.Equal(t, uint64(90), bcc.LastNonce())
	require.Equal(t, uint64(98), bcc.CurrentNonce())

	require.Equal(t, uint64(96), bcc.LastRound())
	require.Equal(t, uint64(99), bcc.CurrentRound())

	require.Equal(t, uint64(6749), bcc.LastTimeStamp())
	require.Equal(t, uint64(6800), bcc.CurrentTimeStamp())

	require.Equal(t, []byte("root hash"), bcc.GetStateRootHash())
	require.Equal(t, []byte("last random seed"), bcc.LastRandomSeed())
	require.Equal(t, []byte("current random seed"), bcc.CurrentRandomSeed())
}

package contexts

import (
	"bytes"
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
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
	{Exists: false, Address: []byte("account_missing"), Balance: big.NewInt(0)},
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
	require.True(t, bcc.AccountExists([]byte("account_old")))
	require.False(t, bcc.AccountExists([]byte("account_faulty")))

	mockBlockchain.Err = ErrTestError
	require.False(t, bcc.AccountExists([]byte("account_something")))
}

func TestBlockchainContext_GetBalance(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostMock{}
	host.OutputContext = &mock.OutputContextStub{
		GetOutputAccountCalled: func(address []byte) (*vmcommon.OutputAccount, bool) {
			return &vmcommon.OutputAccount{}, true
		},
	}

	mockBlockchain := mock.NewBlockchainHookMock()
	mockBlockchain.AddAccounts(testAccounts)

	bcc, _ := NewBlockchainContext(host, mockBlockchain)

	mockBlockchain.Err = ErrTestError
	balanceBytes := bcc.GetBalance([]byte("account_new_with_money"))
	value := big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, big.NewInt(0), value)
	mockBlockchain.Err = nil

	balanceBytes = bcc.GetBalance([]byte("account_faulty"))
	value = big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, big.NewInt(0), value)

	balanceBytes = bcc.GetBalance([]byte("account_missing"))
	value = big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, big.NewInt(0), value)

	balanceBytes = bcc.GetBalance([]byte("account_new"))
	value = big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, big.NewInt(0), value)

	balanceBytes = bcc.GetBalance([]byte("account_new_with_money"))
	value = big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, big.NewInt(1000), value)

	balanceBytes = bcc.GetBalance([]byte("account_old_with_money"))
	value = big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, big.NewInt(1024), value)
}

func TestBlockchainContext_GetNonceAndIncrease(t *testing.T) {
	t.Parallel()

	account := &vmcommon.OutputAccount{
		Address: []byte("account"),
		Nonce:   12,
	}

	accountIsNew := true
	accountIsNewPointer := &accountIsNew

	host := &mock.VmHostMock{}
	host.OutputContext = &mock.OutputContextStub{
		GetOutputAccountCalled: func(address []byte) (*vmcommon.OutputAccount, bool) {
			return account, *accountIsNewPointer
		},
	}

	mockBlockchain := mock.NewBlockchainHookMock()
	bcc, _ := NewBlockchainContext(host, mockBlockchain)

	mockBlockchain.Err = ErrTestError
	nonce, err := bcc.GetNonce([]byte("account"))
	require.Equal(t, ErrTestError, err)
	require.Equal(t, uint64(0), nonce)
	mockBlockchain.Err = nil

	account.Nonce = 12
	accountIsNew = false
	nonce, err = bcc.GetNonce([]byte("account"))
	require.Nil(t, err)
	require.Equal(t, uint64(12), nonce)

	bcc.IncreaseNonce([]byte("account"))
	bcc.IncreaseNonce([]byte("account"))
	bcc.IncreaseNonce([]byte("account"))
	nonce, err = bcc.GetNonce([]byte("account"))
	require.Nil(t, err)
	require.Equal(t, uint64(15), nonce)
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
	require.Equal(t, len(expectedCode), size)
}

func TestBlockchainContext_NewAddress(t *testing.T) {
	t.Parallel()

	returnedAddres := []byte("addr")
	localErr := errors.New("localErr")
	creatorAddr1 := []byte("creatorAddr")
	outCon := mock.NewOutputContextMock()
	host := &mock.VmHostStub{
		OutputCalled: func() arwen.OutputContext {
			return outCon
		},
		RuntimeCalled: func() arwen.RuntimeContext {
			r, _ := NewRuntimeContext(&mock.VmHostStub{}, nil, []byte("type1"))
			return r
		},
	}
	bcHook := &mock.BlockChainHookStub{
		GetNonceCalled: func(address []byte) (u uint64, err error) {
			if bytes.Equal(address, creatorAddr1) {
				return 0, localErr
			}
			return 1, nil
		},
		NewAddressCalled: func(creatorAddress []byte, creatorNonce uint64, vmType []byte) (i []byte, err error) {
			return returnedAddres, nil
		},
	}
	bcc, _ := NewBlockchainContext(host, bcHook)

	address, err := bcc.NewAddress(creatorAddr1)
	require.Equal(t, localErr, err)
	require.Nil(t, address)

	address, err = bcc.NewAddress([]byte("addr2"))
	require.Nil(t, err)
	require.Equal(t, returnedAddres, address)
}

func TestBlockchainContext_BlockHash(t *testing.T) {
	t.Parallel()

	blockHash := []byte("blockHash")
	host := &mock.VmHostStub{}
	bcHook := &mock.BlockChainHookStub{
		GetBlockHashCalled: func(nonce uint64) (i []byte, err error) {
			if nonce == 0 {
				return nil, errors.New("err")
			}
			return blockHash, nil
		},
	}
	bcc, _ := NewBlockchainContext(host, bcHook)

	hash := bcc.BlockHash(-1)
	require.Nil(t, hash)

	hash = bcc.BlockHash(0)
	require.Nil(t, hash)

	hash = bcc.BlockHash(1)
	require.Equal(t, blockHash, hash)
}

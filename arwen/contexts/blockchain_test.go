package contexts

import (
	"bytes"
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	mockhookcrypto "github.com/ElrondNetwork/elrond-vm-util/mock-hook-crypto"
	"github.com/stretchr/testify/require"
)

func TestNewBlockchainContext(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostStub{}
	bcHook := &mock.BlockChainHookStub{}

	bcc, err := NewBlockchainContext(host, bcHook)
	require.Nil(t, err)
	require.NotNil(t, bcc)
}

func TestBlockchainContext_AccountExists(t *testing.T) {
	t.Parallel()

	type testData struct {
		testName string
		testErr  error
		accExits bool
	}

	tData := []testData{
		{testName: "AccountExits", testErr: nil, accExits: true},
		{testName: "AccountNotExits", testErr: errors.New("test err"), accExits: false},
	}
	for _, data := range tData {
		t.Run(data.testName, func(t *testing.T) {
			host := &mock.VmHostStub{}
			bcHook := &mock.BlockChainHookStub{
				AccountExtistsCalled: func(address []byte) (bool, error) {
					if data.testErr != nil {
						return false, data.testErr
					}
					return true, nil
				},
			}
			bcc, _ := NewBlockchainContext(host, bcHook)

			exits := bcc.AccountExists([]byte("addr"))
			require.Equal(t, data.accExits, exits)
		})
	}
}

func TestBlockchainContext_GetBalance(t *testing.T) {
	t.Parallel()

	balance := big.NewInt(1000)
	addr1 := []byte("addr1")
	localErr := errors.New("localErr")
	outCon := mock.NewOutputContextMock()
	host := &mock.VmHostStub{
		OutputCalled: func() arwen.OutputContext {
			return outCon
		},
	}
	bcHook := &mock.BlockChainHookStub{
		GetBalanceCalled: func(address []byte) (b *big.Int, err error) {
			if bytes.Equal(address, addr1) {
				return nil, localErr
			}
			return balance, nil
		},
	}
	bcc, _ := NewBlockchainContext(host, bcHook)

	balanceBytes := bcc.GetBalance(addr1)
	value := big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, big.NewInt(0), value)

	balanceBytes = bcc.GetBalance(addr1)
	value = big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, big.NewInt(0), value)

	balanceBytes = bcc.GetBalance([]byte("add2"))
	value = big.NewInt(0).SetBytes(balanceBytes)
	require.Equal(t, balance, value)
}

func TestBlockchainContext_GetNonceAndIncrease(t *testing.T) {
	t.Parallel()

	addr1 := []byte("addr1")
	localErr := errors.New("localErr")
	outCon := mock.NewOutputContextMock()
	host := &mock.VmHostStub{
		OutputCalled: func() arwen.OutputContext {
			return outCon
		},
	}
	bcHook := &mock.BlockChainHookStub{
		GetNonceCalled: func(address []byte) (u uint64, err error) {
			return 0, localErr
		},
	}
	bcc, _ := NewBlockchainContext(host, bcHook)

	nonce, err := bcc.GetNonce(addr1)
	require.Equal(t, localErr, err)
	require.Equal(t, uint64(0), nonce)

	bcc.IncreaseNonce(addr1)
	nonce, err = bcc.GetNonce(addr1)
	require.Nil(t, err)
	require.Equal(t, uint64(1), nonce)
}

func TestBlockchainContext_GetCodeHash(t *testing.T) {
	t.Parallel()

	code := []byte("code")
	addr := []byte("addr")
	addr2 := []byte("addr2")
	localErr := errors.New("localErr")
	crypto := mockhookcrypto.KryptoHookMockInstance
	host := &mock.VmHostStub{
		CryptoCalled: func() vmcommon.CryptoHook {
			return crypto
		},
	}
	bcHook := &mock.BlockChainHookStub{
		GetCodeCalled: func(address []byte) (i []byte, err error) {
			if bytes.Equal(address, addr) {
				return nil, localErr
			}
			return code, nil
		},
	}
	bcc, _ := NewBlockchainContext(host, bcHook)

	codeHash, err := bcc.GetCodeHash(addr)
	require.Equal(t, localErr, err)
	require.Nil(t, codeHash)

	size, err := bcc.GetCodeSize(addr)
	require.Equal(t, int32(0), size)
	require.Equal(t, localErr, err)

	codeHash, err = bcc.GetCodeHash(addr2)
	require.Nil(t, err)

	c, _ := bcc.GetCode(addr2)
	require.Nil(t, err)

	size, err = bcc.GetCodeSize(addr2)
	require.Equal(t, int32(len(code)), size)
	require.Nil(t, err)

	expectedCodeHash, _ := crypto.Keccak256(c)
	require.Equal(t, expectedCodeHash, codeHash)
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

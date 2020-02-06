package contexts

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

func TestNewOutputContext(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostStub{}

	outputContext, err := NewOutputContext(host)
	require.Nil(t, err)
	require.NotNil(t, outputContext)
}

func TestOutputContext_PushPopState(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostStub{}
	outputContext, _ := NewOutputContext(host)

	//create a new account
	address := []byte("address")
	_, isNew := outputContext.GetOutputAccount(address)
	require.True(t, isNew)

	outputContext.PushState()
	require.Equal(t, 1, len(outputContext.stateStack))

	err := outputContext.PopState()
	require.Nil(t, err)

	err = outputContext.PopState()
	require.Equal(t, arwen.StateStackUnderflow, err)
}

func TestCheckMergeAccountsWorkCorrectly(t *testing.T) {
	t.Parallel()

	left := &vmcommon.OutputAccount{
		Address:        []byte("addr1"),
		Nonce:          1,
		Balance:        big.NewInt(1000),
		BalanceDelta:   big.NewInt(10000),
		StorageUpdates: nil,
		Code:           []byte("code1"),
		Data:           []byte("data2"),
		GasLimit:       100000,
	}
	right := &vmcommon.OutputAccount{
		Address:        []byte("addr2"),
		Nonce:          2,
		Balance:        big.NewInt(2000),
		BalanceDelta:   big.NewInt(20000),
		StorageUpdates: map[string]*vmcommon.StorageUpdate{"key": {Data: []byte("data"), Offset: []byte("offset")}},
		Code:           []byte("code2"),
		Data:           []byte("data2"),
		GasLimit:       100000,
	}

	expected := &vmcommon.OutputAccount{
		Address:        []byte("addr2"),
		Nonce:          2,
		Balance:        big.NewInt(1000),
		BalanceDelta:   big.NewInt(30000),
		StorageUpdates: map[string]*vmcommon.StorageUpdate{"key": {Data: []byte("data"), Offset: []byte("offset")}},
		Code:           []byte("code2"),
		Data:           []byte("data2"),
		GasLimit:       100000,
	}

	mergeOutputAccounts(left, right)
	require.Equal(t, expected, left)
}

func TestOutputContext_Transfer(t *testing.T) {
	t.Parallel()

	accAddr1 := []byte("accAddr1")
	accAddr2 := []byte("accAddr2")
	balance1 := big.NewInt(10000)
	valueToTransfer := big.NewInt(1000)

	host := &mock.VmHostStub{}
	outputContext, _ := NewOutputContext(host)
	outputContext.AddTxValueToAccount(accAddr1, balance1)

	outputContext.Transfer(accAddr2, accAddr1, 1000, valueToTransfer, []byte("input"))

	senderAccount, isNew := outputContext.GetOutputAccount(accAddr1)
	require.False(t, isNew)
	require.Equal(t, big.NewInt(9000), senderAccount.BalanceDelta)

	destAccount, isNew := outputContext.GetOutputAccount(accAddr2)
	require.False(t, isNew)
	require.Equal(t, valueToTransfer, destAccount.BalanceDelta)
}

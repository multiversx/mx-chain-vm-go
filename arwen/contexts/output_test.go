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

	output, err := NewOutputContext(host)
	require.Nil(t, err)
	require.NotNil(t, output)

	require.Equal(t, vmcommon.Ok, output.ReturnCode())
	require.NotNil(t, output.ReturnData())
	require.Equal(t, 0, len(output.ReturnData()))
	require.Equal(t, "", output.ReturnMessage())

	require.Equal(t, uint64(0), output.GetRefund())

	require.NotNil(t, output.outputState.OutputAccounts)
	require.Equal(t, 0, len(output.outputState.OutputAccounts))

	vmOutput := output.outputState
	require.NotNil(t, vmOutput.DeletedAccounts)
	require.Equal(t, 0, len(vmOutput.DeletedAccounts))
	require.NotNil(t, vmOutput.TouchedAccounts)
	require.Equal(t, 0, len(vmOutput.TouchedAccounts))
	require.NotNil(t, vmOutput.Logs)
	require.Equal(t, 0, len(vmOutput.Logs))
}

func TestOutputContext_PushPopState(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostStub{}
	outputContext, _ := NewOutputContext(host)

	address := []byte("address")
	account, isNew := outputContext.GetOutputAccount(address)
	require.True(t, isNew)
	require.Equal(t, 1, len(outputContext.outputState.OutputAccounts))

	account.Nonce = 99
	outputContext.PushState()
	outputContext.InitState()
	require.Equal(t, 1, len(outputContext.stateStack))
	require.Equal(t, 0, len(outputContext.outputState.OutputAccounts))

	err := outputContext.PopState()
	account, isNew = outputContext.GetOutputAccount(address)
	require.False(t, isNew)
	require.Nil(t, err)
	require.Equal(t, uint64(99), account.Nonce)
	require.Equal(t, 1, len(outputContext.outputState.OutputAccounts))
	require.Equal(t, 0, len(outputContext.stateStack))

	err = outputContext.PopState()
	require.Equal(t, arwen.StateStackUnderflow, err)
}

func TestOutputContext_GetOutputAccount(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostStub{}
	output, _ := NewOutputContext(host)

	require.Zero(t, len(output.outputState.OutputAccounts))

	account, isNew := output.GetOutputAccount([]byte("account"))
	require.Equal(t, 1, len(output.outputState.OutputAccounts))
	require.True(t, isNew)
	require.Equal(t, []byte("account"), account.Address)
	require.Zero(t, account.Nonce)
	require.Equal(t, big.NewInt(0), account.BalanceDelta)
	require.Equal(t, big.NewInt(0), account.Balance)
	require.Zero(t, len(account.StorageUpdates))

	account.Address = []byte("changed address")
	account.Nonce = 88
	account.Balance = big.NewInt(94)

	cachedAccount, isNew := output.GetOutputAccount([]byte("account"))
	require.False(t, isNew)
	require.Equal(t, []byte("changed address"), cachedAccount.Address)
	require.Equal(t, uint64(88), cachedAccount.Nonce)
	require.Equal(t, big.NewInt(94), cachedAccount.Balance)
	require.Zero(t, len(cachedAccount.StorageUpdates))
}

func TestOutputContext_GettersAndSetters(t *testing.T) {
	host := &mock.VmHostStub{}
	output, _ := NewOutputContext(host)

	output.SetRefund(24)
	require.Equal(t, uint64(24), output.GetRefund())

}

func TestOutputContext_FinishReturnData(t *testing.T) {
	host := &mock.VmHostStub{}
	output, _ := NewOutputContext(host)

	require.Equal(t, 0, len(output.ReturnData()))

	output.Finish([]byte("something"))
	expectedData := [][]byte{
		[]byte("something"),
	}
	require.Equal(t, expectedData, output.ReturnData())

	output.Finish([]byte("something else"))
	expectedData = [][]byte{
		[]byte("something"),
		[]byte("something else"),
	}
	require.Equal(t, expectedData, output.ReturnData())

	output.Finish(big.NewInt(1234567).Bytes())
	expectedData = [][]byte{
		[]byte("something"),
		[]byte("something else"),
		big.NewInt(1234567).Bytes(),
	}
	require.Equal(t, expectedData, output.ReturnData())

	// TODO remove this section after modifying Finish to accept empty []byte
	// slices
	output.Finish([]byte{})
	expectedData = [][]byte{
		[]byte("something"),
		[]byte("something else"),
		big.NewInt(1234567).Bytes(),
	}
	require.Equal(t, expectedData, output.ReturnData())
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

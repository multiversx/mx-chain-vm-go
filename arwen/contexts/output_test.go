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

	require.Equal(t, vmcommon.Ok, outputContext.ReturnCode())
	require.NotNil(t, outputContext.ReturnData())
	require.Equal(t, 0, len(outputContext.ReturnData()))
	require.Equal(t, "", outputContext.ReturnMessage())
	require.Equal(t, 0, len(outputContext.stateStack))

	require.Equal(t, uint64(0), outputContext.GetRefund())

	require.NotNil(t, outputContext.outputState.OutputAccounts)
	require.Equal(t, 0, len(outputContext.outputState.OutputAccounts))

	vmOutput := outputContext.outputState
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

	address1 := []byte("address1")
	address2 := []byte("address2")

	// Create an account with nonce 99 on the active state.
	account, isNew := outputContext.GetOutputAccount(address1)
	account.Nonce = 99
	require.True(t, isNew)
	require.Equal(t, 1, len(outputContext.outputState.OutputAccounts))

	// Copy active state onto the stack.
	outputContext.PushState()
	require.Equal(t, 1, len(outputContext.stateStack))

	// Clear the active state and create a new account with the same address as
	// the previous; the new account must not have nonce 99.
	outputContext.InitState()
	account, isNew = outputContext.GetOutputAccount(address1)
	require.True(t, isNew)
	require.Equal(t, 1, len(outputContext.outputState.OutputAccounts))
	require.Equal(t, uint64(0), account.Nonce)

	account.Nonce = 84

	// Copy active state onto the stack, then create a new account with nonce 42.
	outputContext.PushState()
	require.Equal(t, 2, len(outputContext.stateStack))

	account, isNew = outputContext.GetOutputAccount(address2)
	account.Nonce = 42
	require.True(t, isNew)
	require.Equal(t, 2, len(outputContext.outputState.OutputAccounts))

	// Revert to the previous state: account with nonce 42 is lost, and the
	// account with "address1" has nonce 84.
	outputContext.PopSetActiveState()
	account, isNew = outputContext.GetOutputAccount(address1)
	require.False(t, isNew)
	require.Equal(t, uint64(84), account.Nonce)
	require.Equal(t, 1, len(outputContext.outputState.OutputAccounts))
	require.Equal(t, 1, len(outputContext.stateStack))

	outputContext.PushState()
	require.Equal(t, 2, len(outputContext.stateStack))

	outputContext.PopDiscard()
	require.Equal(t, 1, len(outputContext.stateStack))

	account, isNew = outputContext.GetOutputAccount(address1)
	require.False(t, isNew)
	require.Equal(t, uint64(84), account.Nonce)
	require.Equal(t, 1, len(outputContext.outputState.OutputAccounts))
	require.Equal(t, 1, len(outputContext.stateStack))

	outputContext.ClearStateStack()
	require.Equal(t, 0, len(outputContext.stateStack))
}

func TestOutputContext_GetOutputAccount(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostStub{}
	outputContext, _ := NewOutputContext(host)
	require.Zero(t, len(outputContext.outputState.OutputAccounts))

	// Request an account that is missing from OutputAccounts
	account, isNew := outputContext.GetOutputAccount([]byte("account"))
	require.Equal(t, 1, len(outputContext.outputState.OutputAccounts))
	require.True(t, isNew)
	require.Equal(t, []byte("account"), account.Address)
	require.Zero(t, account.Nonce)
	require.Equal(t, arwen.Zero, account.BalanceDelta)
	require.Nil(t, account.Balance)
	require.Zero(t, len(account.StorageUpdates))

	// Change fields of the OutputAccount to ensure it will be returned on the
	// next call to GetOutputAccount(), from the OutputAccounts cache
	account.Address = []byte("changed address")
	account.Nonce = 88
	account.Balance = big.NewInt(94)
	cachedAccount, isNew := outputContext.GetOutputAccount([]byte("account"))
	require.False(t, isNew)
	require.Equal(t, []byte("changed address"), cachedAccount.Address)
	require.Equal(t, uint64(88), cachedAccount.Nonce)
	require.Equal(t, big.NewInt(94), cachedAccount.Balance)
	require.Zero(t, len(cachedAccount.StorageUpdates))
}

func TestOutputContext_GettersAndSetters(t *testing.T) {
	host := &mock.VmHostStub{}
	outputContext, _ := NewOutputContext(host)

	outputContext.SetRefund(24)
	require.Equal(t, uint64(24), outputContext.GetRefund())

	outputContext.SetReturnCode(vmcommon.ExecutionFailed)
	require.Equal(t, vmcommon.ExecutionFailed, outputContext.ReturnCode())

	outputContext.SetReturnMessage("rockets")
	require.Equal(t, "rockets", outputContext.ReturnMessage())
}

func TestOutputContext_FinishReturnData(t *testing.T) {
	host := &mock.VmHostStub{}
	outputContext, _ := NewOutputContext(host)

	require.Zero(t, len(outputContext.ReturnData()))

	outputContext.Finish([]byte("something"))
	expectedData := [][]byte{
		[]byte("something"),
	}
	require.Equal(t, expectedData, outputContext.ReturnData())

	outputContext.Finish([]byte("something else"))
	expectedData = append(expectedData, []byte("something else"))
	require.Equal(t, expectedData, outputContext.ReturnData())

	outputContext.Finish(big.NewInt(1234567).Bytes())
	expectedData = append(expectedData, big.NewInt(1234567).Bytes())
	require.Equal(t, expectedData, outputContext.ReturnData())

	outputContext.Finish([]byte{})
	expectedData = append(expectedData, []byte{})
	require.Equal(t, expectedData, outputContext.ReturnData())

	outputContext.ClearReturnData()
	require.Zero(t, len(outputContext.ReturnData()))
}

func TestOutputContext_MergeCompleteAccounts(t *testing.T) {
	t.Parallel()

	left := &vmcommon.OutputAccount{
		Address:        []byte("addr1"),
		Nonce:          1,
		Balance:        big.NewInt(1000),
		BalanceDelta:   big.NewInt(10000),
		StorageUpdates: nil,
		Code:           []byte("code1"),
		Data:           []byte("data1"),
		GasLimit:       99999,
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
		Balance:        big.NewInt(2000),
		BalanceDelta:   big.NewInt(20000),
		StorageUpdates: map[string]*vmcommon.StorageUpdate{"key": {Data: []byte("data"), Offset: []byte("offset")}},
		Code:           []byte("code2"),
		Data:           []byte("data2"),
		GasLimit:       100000,
	}

	mergeOutputAccounts(left, right)
	require.Equal(t, expected, left)
}

func TestOutputContext_MergeIncompleteAccounts(t *testing.T) {
	t.Parallel()

	left := &vmcommon.OutputAccount{}
	right := &vmcommon.OutputAccount{}
	expected := &vmcommon.OutputAccount{
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
		BalanceDelta:   big.NewInt(0),
	}
	mergeOutputAccounts(left, right)
	require.Equal(t, expected, left)

	left = &vmcommon.OutputAccount{
		GasLimit: 98,
	}
	right = &vmcommon.OutputAccount{
		BalanceDelta: big.NewInt(42),
	}
	expected = &vmcommon.OutputAccount{
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
		BalanceDelta:   big.NewInt(42),
	}
	mergeOutputAccounts(left, right)
	require.Equal(t, expected, left)

	left = &vmcommon.OutputAccount{
		BalanceDelta: big.NewInt(48),
	}
	right = &vmcommon.OutputAccount{}
	expected = &vmcommon.OutputAccount{
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
		BalanceDelta:   big.NewInt(48),
	}
	mergeOutputAccounts(left, right)
	require.Equal(t, expected, left)

	left = &vmcommon.OutputAccount{
		Address: []byte("left address"),
		Code:    []byte("left code"),
	}
	right = &vmcommon.OutputAccount{
		BalanceDelta: big.NewInt(42),
	}
	expected = &vmcommon.OutputAccount{
		Code:           []byte("left code"),
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
		BalanceDelta:   big.NewInt(42),
	}
	mergeOutputAccounts(left, right)
	require.Equal(t, expected, left)

	left = &vmcommon.OutputAccount{
		Data: []byte("left data"),
	}
	right = &vmcommon.OutputAccount{}
	expected = &vmcommon.OutputAccount{
		Data:           []byte("left data"),
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
		BalanceDelta:   big.NewInt(0),
	}
	mergeOutputAccounts(left, right)
	require.Equal(t, expected, left)

	left = &vmcommon.OutputAccount{
		Nonce: 44,
	}
	right = &vmcommon.OutputAccount{
		Nonce: 42,
	}
	expected = &vmcommon.OutputAccount{
		Nonce:          44,
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
		BalanceDelta:   big.NewInt(0),
	}
	mergeOutputAccounts(left, right)
	require.Equal(t, expected, left)
}

func TestOutputContext_MergeVMOutputs(t *testing.T) {
	t.Parallel()

	left := newVMOutput()
	right := newVMOutput()
	expected := newVMOutput()
	mergeVMOutputs(left, right)
	require.Equal(t, expected, left)

	left = newVMOutput()
	right = newVMOutput()
	right.OutputAccounts["address"] = newVMOutputAccount([]byte("address"))
	right.OutputAccounts["address"].Nonce = 84
	expected = newVMOutput()
	expected.OutputAccounts["address"] = newVMOutputAccount([]byte("address"))
	expected.OutputAccounts["address"].Nonce = 84
	mergeVMOutputs(left, right)
	require.Equal(t, expected, left)

	left = newVMOutput()
	left.OutputAccounts["address"] = newVMOutputAccount([]byte("address"))
	left.OutputAccounts["address"].Nonce = 84
	right = newVMOutput()
	right.OutputAccounts["address"] = newVMOutputAccount([]byte("address"))
	right.OutputAccounts["address"].Nonce = 92
	expected = newVMOutput()
	expected.OutputAccounts["address"] = newVMOutputAccount([]byte("address"))
	expected.OutputAccounts["address"].Nonce = 92
	mergeVMOutputs(left, right)
	require.Equal(t, expected, left)

	left = newVMOutput()
	left.OutputAccounts["left address"] = newVMOutputAccount([]byte("left address"))
	right = newVMOutput()
	right.OutputAccounts["right address"] = newVMOutputAccount([]byte("right address"))
	expected = newVMOutput()
	expected.OutputAccounts["left address"] = newVMOutputAccount([]byte("left address"))
	expected.OutputAccounts["right address"] = newVMOutputAccount([]byte("right address"))
	mergeVMOutputs(left, right)
	require.Equal(t, expected, left)

	left = newVMOutput()
	left.GasRemaining = 99
	left.GasRefund = big.NewInt(42)
	left.ReturnCode = vmcommon.ContractNotFound
	left.ReturnMessage = "fireflies"
	right = newVMOutput()
	right.ReturnData = [][]byte{[]byte("rockets"), []byte("albatross")}
	right.GasRemaining = 100
	right.GasRefund = big.NewInt(84)
	right.ReturnCode = vmcommon.ExecutionFailed
	right.ReturnMessage = "turtles all the way down"
	expected = newVMOutput()
	expected.ReturnData = [][]byte{[]byte("rockets"), []byte("albatross")}
	expected.GasRemaining = 100
	expected.GasRefund = big.NewInt(84)
	expected.ReturnCode = vmcommon.ExecutionFailed
	expected.ReturnMessage = "turtles all the way down"
	mergeVMOutputs(left, right)
	require.Equal(t, expected, left)
}

func TestOutputContext_VMOutputError(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostStub{}
	outputContext, _ := NewOutputContext(host)

	returnCode := vmcommon.ContractNotFound
	returnMessage := "computer not found"

	expected := &vmcommon.VMOutput{
		GasRemaining:  0,
		GasRefund:     big.NewInt(0),
		ReturnCode:    returnCode,
		ReturnMessage: returnMessage,
	}
	vmOutput := outputContext.CreateVMOutputInCaseOfError(returnCode, returnMessage)
	require.Equal(t, expected, vmOutput)
}

func TestOutputContext_Transfer(t *testing.T) {
	t.Parallel()

	sender := []byte("sender")
	receiver := []byte("receiver")
	balance := big.NewInt(10000)
	valueToTransfer := big.NewInt(1000)

	host := &mock.VmHostMock{}
	mockBlockchainHook := mock.NewBlockchainHookMock()
	mockBlockchainHook.AddAccount(&mock.Account{
		Exists:  true,
		Address: sender,
		Nonce:   42,
		Balance: balance,
	})

	blockchainContext, _ := NewBlockchainContext(host, mockBlockchainHook)
	outputContext, _ := NewOutputContext(host)

	host.OutputContext = outputContext
	host.BlockchainContext = blockchainContext

	err := outputContext.Transfer(receiver, sender, 54, valueToTransfer, []byte("txdata"))
	require.Nil(t, err)

	senderAccount, isNew := outputContext.GetOutputAccount(sender)
	require.False(t, isNew)
	require.Equal(t, big.NewInt(-1000), senderAccount.BalanceDelta)

	destAccount, isNew := outputContext.GetOutputAccount(receiver)
	require.False(t, isNew)
	require.Equal(t, valueToTransfer, destAccount.BalanceDelta)
	require.Equal(t, uint64(54), destAccount.GasLimit)
	require.Equal(t, []byte("txdata"), destAccount.Data)
}

func TestOutputContext_Transfer_Errors_And_Checks(t *testing.T) {
	t.Parallel()

	sender := []byte("sender")
	receiver := []byte("receiver")

	mockBlockchainHook := mock.NewBlockchainHookMock()
	mockBlockchainHook.AddAccount(&mock.Account{
		Exists:  true,
		Address: sender,
		Nonce:   88,
		Balance: big.NewInt(2000),
	})

	host := &mock.VmHostMock{}
	outputContext, _ := NewOutputContext(host)
	blockchainContext, _ := NewBlockchainContext(host, mockBlockchainHook)

	host.OutputContext = outputContext
	host.BlockchainContext = blockchainContext

	senderOutputAccount, _ := outputContext.GetOutputAccount(sender)
	require.Nil(t, senderOutputAccount.Balance)
	require.Equal(t, arwen.Zero, senderOutputAccount.BalanceDelta)

	// negative transfers are disallowed
	valueToTransfer := big.NewInt(-1000)
	err := outputContext.Transfer(receiver, sender, 54, valueToTransfer, []byte("txdata"))
	require.Equal(t, arwen.ErrTransferNegativeValue, err)
	require.Nil(t, senderOutputAccount.Balance)
	require.Equal(t, arwen.Zero, senderOutputAccount.BalanceDelta)

	// account must have enough money to transfer
	valueToTransfer = big.NewInt(5000)
	err = outputContext.Transfer(receiver, sender, 54, valueToTransfer, []byte("txdata"))
	require.Equal(t, arwen.ErrTransferInsufficientFunds, err)
	require.Equal(t, big.NewInt(2000), senderOutputAccount.Balance)
	require.Equal(t, arwen.Zero, senderOutputAccount.BalanceDelta)

	senderOutputAccount.BalanceDelta = big.NewInt(4000)
	valueToTransfer = big.NewInt(5000)
	err = outputContext.Transfer(receiver, sender, 54, valueToTransfer, []byte("txdata"))
	require.Nil(t, err)
	require.Equal(t, big.NewInt(-1000), senderOutputAccount.BalanceDelta)

	require.Equal(t, big.NewInt(1000), blockchainContext.GetBalanceBigInt(sender))
}

func TestOutputContext_WriteLog(t *testing.T) {
	// TODO first discuss how Logs should be implemented
}

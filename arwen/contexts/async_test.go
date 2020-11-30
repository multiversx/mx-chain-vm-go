package contexts

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

var OriginalCaller = []byte("address_original_caller")
var Alice = []byte("address_alice")
var Bob = []byte("address_bob")

func InitializeArwenAndWasmer_AsyncContext() (*mock.VmHostMock, *mock.BlockchainHookMock) {
	imports := MakeAPIImports()
	_ = wasmer.SetImports(imports)

	gasSchedule := config.MakeGasMapForTests()
	gasCostConfig, _ := config.CreateGasConfig(gasSchedule)
	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)

	host := &mock.VmHostMock{}
	host.SCAPIMethods = imports

	mockMetering := &mock.MeteringContextMock{}
	mockMetering.SetGasSchedule(gasSchedule)
	host.MeteringContext = mockMetering

	bhm := mock.NewBlockchainHookMock()
	host.BlockchainContext, _ = NewBlockchainContext(host, bhm)
	host.RuntimeContext, _ = NewRuntimeContext(host, []byte("vm"), false)
	host.OutputContext, _ = NewOutputContext(host)
	host.CryptoHook = crypto.NewVMCrypto()
	return host, bhm
}

func InitializeArwenAndWasmer_AsyncContext_AliceAndBob() (
	*mock.VmHostMock,
	*mock.BlockchainHookMock,
	*vmcommon.ContractCallInput,
) {
	host, bhm := InitializeArwenAndWasmer_AsyncContext()
	bhm.AddAccount(&mock.AccountMock{
		Address: Alice,
		Balance: big.NewInt(88),
	})
	bhm.AddAccount(&mock.AccountMock{
		Address: Bob,
		Balance: big.NewInt(12),
	})

	originalVMInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     OriginalCaller,
			Arguments:      nil,
			CallType:       vmcommon.DirectCall,
			GasPrice:       1,
			CurrentTxHash:  []byte("txhash"),
			PrevTxHash:     []byte("txhash"),
			OriginalTxHash: []byte("txhash"),
		},
		RecipientAddr: Alice,
		Function:      "alice_function",
	}

	return host, bhm, originalVMInput
}

func TestNewAsyncContext(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)
	require.NotNil(t, async)

	require.NotNil(t, async.host)
	require.Nil(t, async.stateStack)
	require.Nil(t, async.callerAddr)
	require.Nil(t, async.returnData)
	require.NotNil(t, async.asyncCallGroups)
	require.Empty(t, async.asyncCallGroups)
}

func TestAsyncContext_InitState(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)
	require.NotNil(t, async)

	async.callerAddr = []byte("some address")
	async.gasPrice = 1000
	async.returnData = []byte("some return data")
	async.asyncCallGroups = nil

	async.InitState()

	require.NotNil(t, async.callerAddr)
	require.Empty(t, async.callerAddr)
	require.Zero(t, async.gasPrice)
	require.NotNil(t, async.returnData)
	require.Empty(t, async.returnData)
	require.NotNil(t, async.asyncCallGroups)
	require.Empty(t, async.asyncCallGroups)
}

func TestAsyncContext_InitStateFromInput(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)
	require.NotNil(t, async)

	async.callerAddr = []byte("some address")
	async.gasPrice = 1000
	async.returnData = []byte("some return data")
	async.asyncCallGroups = nil

	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: []byte("real caller addr"),
			GasPrice:   42,
		},
	}

	async.InitStateFromInput(input)

	require.Equal(t, input.CallerAddr, async.callerAddr)
	require.Equal(t, uint64(42), async.gasPrice)
	require.NotNil(t, async.returnData)
	require.Empty(t, async.returnData)
	require.NotNil(t, async.asyncCallGroups)
	require.Empty(t, async.asyncCallGroups)
}

func TestAsyncContext_GettersAndSetters(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)
	require.NotNil(t, async)

	async.callerAddr = []byte("some address")
	async.gasPrice = 1000
	async.returnData = []byte("some return data")

	require.Equal(t, []byte("some address"), async.GetCallerAddress())
	require.Equal(t, uint64(1000), async.GetGasPrice())
	require.Equal(t, []byte("some return data"), async.GetReturnData())

	async.SetReturnData([]byte("rockets"))
	require.Equal(t, []byte("rockets"), async.GetReturnData())
}

func TestAsyncContext_AddCall_NewGroup_DeleteGroup(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()

	async := NewAsyncContext(host)
	require.NotNil(t, async)

	require.True(t, async.IsComplete())
	require.False(t, async.HasPendingCallGroups())

	group, exists := async.GetCallGroup("testGroup")
	require.Nil(t, group)
	require.False(t, exists)

	err := async.AddCall("testGroup", &arwen.AsyncCall{
		Destination: []byte("somewhere"),
		Data:        []byte("something"),
	})
	require.Nil(t, err)
	require.False(t, async.IsComplete())
	require.True(t, async.HasPendingCallGroups())
	require.Equal(t, 1, len(async.asyncCallGroups))

	group, exists = async.GetCallGroup("testGroup")
	require.NotNil(t, group)
	require.True(t, exists)
	require.True(t, group.HasPendingCalls())
	require.False(t, group.IsComplete())
	require.False(t, group.HasCallback())

	async.deleteCallGroupByID("testGroup")
	group, exists = async.GetCallGroup("testGroup")
	require.Nil(t, group)
	require.False(t, exists)
}

func TestAsyncContext_AddCall_ExistingGroup(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()

	async := NewAsyncContext(host)
	require.NotNil(t, async)
	require.Equal(t, 0, len(async.asyncCallGroups))

	async.addCallGroup(arwen.NewAsyncCallGroup("testGroup"))
	require.Equal(t, 1, len(async.asyncCallGroups))
	require.False(t, async.IsComplete())

	group, exists := async.GetCallGroup("testGroup")
	require.NotNil(t, group)
	require.True(t, exists)

	err := async.AddCall("testGroup", &arwen.AsyncCall{
		Destination: []byte("somewhere"),
		Data:        []byte("something"),
	})
	require.Nil(t, err)
	require.Equal(t, 1, len(async.asyncCallGroups))
	require.False(t, async.IsComplete())
}

func TestAsyncContext_AddCall_ValidationAndFields(t *testing.T) {
	// TODO execution mode
	// TODO non-nil destination
	// TODO locked gas
}

func TestAsyncContext_SetGroupCallback_GroupDoesntExist(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)

	err := async.SetGroupCallback("testGroup", "callbackFunction", []byte{}, 0)
	require.True(t, errors.Is(err, arwen.ErrAsyncCallGroupDoesNotExist))
}

func TestAsyncContext_SetGroupCallback_OutOfGas(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)

	mockMetering := host.Metering().(*mock.MeteringContextMock)
	mockMetering.Err = arwen.ErrNotEnoughGas

	err := async.AddCall("testGroup", &arwen.AsyncCall{
		Destination: []byte("somewhere"),
		Data:        []byte("something"),
	})

	err = async.SetGroupCallback("testGroup", "callbackFunction", []byte{}, 0)
	require.True(t, errors.Is(err, arwen.ErrNotEnoughGas))
}

func TestAsyncContext_SetGroupCallback_Success(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)

	mockMetering := host.Metering().(*mock.MeteringContextMock)
	mockMetering.GasComputedToLock = 42

	err := async.AddCall("testGroup", &arwen.AsyncCall{
		Destination: []byte("somewhere"),
		Data:        []byte("something"),
	})
	require.Nil(t, err)

	err = async.SetGroupCallback("testGroup", "callbackFunction", []byte{}, 0)
	require.Nil(t, err)

	group, exists := async.GetCallGroup("testGroup")
	require.NotNil(t, group)
	require.True(t, exists)
	require.Equal(t, "callbackFunction", group.Callback)
	require.Equal(t, []byte{}, group.CallbackData)
	require.Equal(t, uint64(42), group.GasLocked)
}

func TestAsyncContext_DetermineExecutionMode(t *testing.T) {
	leftAddress := []byte("left")
	leftAccount := &mock.AccountMock{
		Address: leftAddress,
		ShardID: 0,
	}

	rightAddress := []byte("right")
	rightAccount := &mock.AccountMock{
		Address: rightAddress,
		ShardID: 0,
	}

	host, bhm := InitializeArwenAndWasmer_AsyncContext()
	bhm.AddAccount(leftAccount)
	bhm.AddAccount(rightAccount)
	runtime := host.Runtime()

	async := NewAsyncContext(host)

	runtime.SetSCAddress(leftAddress)
	execMode, err := async.determineExecutionMode(rightAddress, []byte("func"))
	require.Nil(t, err)
	require.Equal(t, arwen.SyncExecution, execMode)

	execMode, err = async.determineExecutionMode(rightAddress, []byte(""))
	require.NotNil(t, err)
	require.Equal(t, arwen.AsyncUnknown, execMode)

	execMode, err = async.determineExecutionMode(rightAddress, []byte(""))
	require.NotNil(t, err)
	require.Equal(t, arwen.AsyncUnknown, execMode)

	host.IsBuiltinFunc = true
	runtime.SetSCAddress(leftAddress)
	execMode, err = async.determineExecutionMode(rightAddress, []byte("func"))
	require.Nil(t, err)
	require.Equal(t, arwen.SyncExecution, execMode)

	host.IsBuiltinFunc = false
	rightAccount.ShardID = 1
	runtime.SetSCAddress(leftAddress)
	execMode, err = async.determineExecutionMode(rightAddress, []byte("func"))
	require.Nil(t, err)
	require.Equal(t, arwen.AsyncUnknown, execMode)

	host.IsBuiltinFunc = true
	rightAccount.ShardID = 1
	runtime.SetSCAddress(leftAddress)
	execMode, err = async.determineExecutionMode(rightAddress, []byte("func"))
	require.Nil(t, err)
	require.Equal(t, arwen.AsyncBuiltinFunc, execMode)
}

func TestAsyncContext_IsValidCallbackName(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)

	require.True(t, async.isValidCallbackName("a"))
	require.True(t, async.isValidCallbackName("my_contract_method_22"))

	require.True(t, async.isValidCallbackName("not_builtin"))
	host.IsBuiltinFunc = true
	require.False(t, async.isValidCallbackName("builtin"))
	host.IsBuiltinFunc = false

	require.True(t, async.isValidCallbackName("callBack"))
	require.True(t, async.isValidCallbackName("callback"))
	require.True(t, async.isValidCallbackName("function_do"))

	require.False(t, async.isValidCallbackName("function-do"))
	require.False(t, async.isValidCallbackName("3_my_contract_method_22"))
	require.False(t, async.isValidCallbackName("init"))
	require.False(t, async.isValidCallbackName("München"))
	require.False(t, async.isValidCallbackName("Göteborg"))
	require.False(t, async.isValidCallbackName("東京"))
	require.False(t, async.isValidCallbackName("function.org"))
	require.False(t, async.isValidCallbackName("Ainulindalë"))
}

func TestAsyncContext_FindCall(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)

	groupID, index, err := async.findCall([]byte("somewhere"))
	require.Equal(t, "", groupID)
	require.Equal(t, -1, index)
	require.True(t, errors.Is(err, arwen.ErrAsyncCallNotFound))

	err = async.AddCall("testGroup", &arwen.AsyncCall{
		Destination: []byte("somewhere"),
		Data:        []byte("something"),
	})
	require.Nil(t, err)

	groupID, index, err = async.findCall([]byte("somewhere"))
	require.Nil(t, err)
	require.Equal(t, "testGroup", groupID)
	require.Equal(t, 0, index)

	err = async.AddCall("testGroup", &arwen.AsyncCall{
		Destination: []byte("somewhere_else"),
		Data:        []byte("something"),
	})
	require.Nil(t, err)

	groupID, index, err = async.findCall([]byte("somewhere_else"))
	require.Nil(t, err)
	require.Equal(t, "testGroup", groupID)
	require.Equal(t, 1, index)

	err = async.AddCall("another_testGroup", &arwen.AsyncCall{
		Destination: []byte("somewhere_else_entirely"),
		Data:        []byte("something"),
	})
	require.Nil(t, err)

	groupID, index, err = async.findCall([]byte("somewhere_else_entirely"))
	require.Nil(t, err)
	require.Equal(t, "another_testGroup", groupID)
	require.Equal(t, 0, index)
}

func TestAsyncContext_UpdateCurrentCallStatus(t *testing.T) {
	vmInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: []byte("caller"),
			Arguments:  [][]byte{{0}},
			CallType:   vmcommon.DirectCall,
		},
	}

	host, _ := InitializeArwenAndWasmer_AsyncContext()
	async := NewAsyncContext(host)

	// CallType == DirectCall, async.UpdateCurrentCallStatus() does nothing
	host.Runtime().InitStateFromInput(vmInput)
	call, err := async.UpdateCurrentCallStatus()
	require.Nil(t, call)
	require.Nil(t, err)

	// CallType == AsynchronousCall, async.UpdateCurrentCallStatus() does nothing
	vmInput.CallType = vmcommon.AsynchronousCall
	host.Runtime().InitStateFromInput(vmInput)
	call, err = async.UpdateCurrentCallStatus()
	require.Nil(t, call)
	require.Nil(t, err)

	// CallType == AsynchronousCallback, but no AsyncCalls registered in the
	// AsyncContext, so async.UpdateCurrentCallStatus() returns an error
	vmInput.CallType = vmcommon.AsynchronousCallBack
	vmInput.Arguments = nil
	host.Runtime().InitStateFromInput(vmInput)
	call, err = async.UpdateCurrentCallStatus()
	require.Nil(t, call)
	require.True(t, errors.Is(err, arwen.ErrCannotInterpretCallbackArgs))

	// CallType == AsynchronousCallback, but no AsyncCalls registered in the
	// AsyncContext, so async.UpdateCurrentCallStatus() returns an error
	vmInput.CallType = vmcommon.AsynchronousCallBack
	vmInput.Arguments = [][]byte{{0}}
	host.Runtime().InitStateFromInput(vmInput)
	call, err = async.UpdateCurrentCallStatus()
	require.Nil(t, call)
	require.True(t, errors.Is(err, arwen.ErrAsyncCallNotFound))

	// CallType == AsynchronousCallback, and there is an AsyncCall registered,
	// but it's not the expected one.
	err = async.AddCall("testGroup", &arwen.AsyncCall{
		Destination: []byte("some_address"),
		Data:        []byte("function"),
	})
	require.Nil(t, err)

	vmInput.CallType = vmcommon.AsynchronousCallBack
	vmInput.Arguments = [][]byte{{0}}
	host.Runtime().InitStateFromInput(vmInput)
	call, err = async.UpdateCurrentCallStatus()
	require.Nil(t, call)
	require.True(t, errors.Is(err, arwen.ErrAsyncCallNotFound))

	// CallType == AsynchronousCallback, but this time there is a corresponding AsyncCall
	// registered, causing async.UpdateCurrentCallStatus() to find and update the AsyncCall
	err = async.AddCall("testGroup", &arwen.AsyncCall{
		Destination: vmInput.CallerAddr,
		Data:        []byte("function"),
	})
	require.Nil(t, err)

	vmInput.CallType = vmcommon.AsynchronousCallBack
	vmInput.Arguments = [][]byte{{0}}
	host.Runtime().InitStateFromInput(vmInput)
	call, err = async.UpdateCurrentCallStatus()
	require.Nil(t, err)
	require.NotNil(t, call)
	require.Equal(t, arwen.AsyncCallResolved, call.Status)

	// CallType == AsynchronousCallback, there is a corresponding AsyncCall
	// registered, causing async.UpdateCurrentCallStatus() to find and update the
	// AsyncCall, but with AsyncCallRejected
	vmInput.CallType = vmcommon.AsynchronousCallBack
	vmInput.Arguments = [][]byte{{1}}
	host.Runtime().InitStateFromInput(vmInput)
	call, err = async.UpdateCurrentCallStatus()
	require.Nil(t, err)
	require.NotNil(t, call)
	require.Equal(t, arwen.AsyncCallRejected, call.Status)
}

func TestAsyncContext_SendAsyncCallCrossShard(t *testing.T) {
	host, bhm := InitializeArwenAndWasmer_AsyncContext()
	bhm.AddAccount(&mock.AccountMock{
		Address: []byte("smartcontract"),
		Balance: big.NewInt(88),
	})

	host.Runtime().SetSCAddress([]byte("smartcontract"))
	async := NewAsyncContext(host)

	call := &arwen.AsyncCall{
		Destination: []byte("destination"),
		GasLimit:    42,
		GasLocked:   98,
		ValueBytes:  big.NewInt(88).Bytes(),
		Data:        []byte("some_data"),
	}

	err := async.sendAsyncCallCrossShard(call)
	require.Nil(t, err)

	vmOutput := host.Output().GetVMOutput()
	require.NotNil(t, vmOutput)

	smartcontract := vmOutput.OutputAccounts["smartcontract"]
	require.Equal(t, big.NewInt(-88), smartcontract.BalanceDelta)
	require.Empty(t, smartcontract.OutputTransfers)

	destination := vmOutput.OutputAccounts["destination"]
	require.Equal(t, big.NewInt(88), destination.BalanceDelta)
	require.Len(t, destination.OutputTransfers, 1)

	asyncTransfer := destination.OutputTransfers[0]
	require.Equal(t, big.NewInt(88), asyncTransfer.Value)
	require.Equal(t, uint64(42), asyncTransfer.GasLimit)
	require.Equal(t, uint64(98), asyncTransfer.GasLocked)
	require.Equal(t, []byte("some_data"), asyncTransfer.Data)
	require.Equal(t, vmcommon.AsynchronousCall, asyncTransfer.CallType)
}

func TestAsyncContext_ExecuteSyncCall(t *testing.T) {
	host, _, originalVMInput := InitializeArwenAndWasmer_AsyncContext_AliceAndBob()
	host.Runtime().InitStateFromInput(originalVMInput)
	async := NewAsyncContext(host)
	call := &arwen.AsyncCall{
		Destination: Bob,
		ValueBytes:  big.NewInt(88).Bytes(),
	}

	// Test error propagation from async.createContractCallInput()
	call.Data = []byte("function")
	call.GasLimit = 1
	err := async.executeSyncCall(call)
	require.True(t, errors.Is(err, arwen.ErrNotEnoughGas))

	// Successful destination execution, but not enough gas for callback execution
	host.EnqueueVMOutput(&vmcommon.VMOutput{
		ReturnCode: vmcommon.Ok,
	})
	call.Data = []byte("function")
	call.GasLimit = 1
	call.GasLocked = 1
	_ = async.executeSyncCall(call)

	expectedDestinationInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     originalVMInput.RecipientAddr,
			Arguments:      [][]byte{[]byte("\x0A\x0B\x0C"), []byte("\x03")},
			CallValue:      big.NewInt(88),
			CallType:       vmcommon.AsynchronousCall,
			GasPrice:       originalVMInput.GasPrice,
			GasProvided:    1,
			CurrentTxHash:  originalVMInput.CurrentTxHash,
			PrevTxHash:     originalVMInput.PrevTxHash,
			OriginalTxHash: originalVMInput.OriginalTxHash,
		},
		RecipientAddr: Bob,
		Function:      "function",
	}

	require.Len(t, host.StoredInputs, 1)
	require.Equal(t, expectedDestinationInput, host.StoredInputs[0])
}

func TestAsyncContext_CreateContractCallInput(t *testing.T) {
	host, _, originalVMInput := InitializeArwenAndWasmer_AsyncContext_AliceAndBob()
	host.Runtime().InitStateFromInput(originalVMInput)
	async := NewAsyncContext(host)
	call := &arwen.AsyncCall{
		Destination: Bob,
		ValueBytes:  big.NewInt(88).Bytes(),
	}

	call.Data = []byte{}
	input, err := async.createContractCallInput(call)
	require.Nil(t, input)
	require.Error(t, err)

	call.Data = []byte("function")
	call.GasLimit = 1
	input, err = async.createContractCallInput(call)
	require.Nil(t, input)
	require.True(t, errors.Is(err, arwen.ErrNotEnoughGas))

	call.Data = []byte("function@0A0B0C@03")
	call.GasLimit = 2
	input, err = async.createContractCallInput(call)
	require.Nil(t, err)
	require.NotNil(t, input)

	expectedInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     Alice,
			Arguments:      [][]byte{[]byte("\x0A\x0B\x0C"), []byte("\x03")},
			CallValue:      big.NewInt(88),
			CallType:       vmcommon.AsynchronousCall,
			GasPrice:       originalVMInput.GasPrice,
			GasProvided:    1,
			CurrentTxHash:  originalVMInput.CurrentTxHash,
			PrevTxHash:     originalVMInput.PrevTxHash,
			OriginalTxHash: originalVMInput.OriginalTxHash,
		},
		RecipientAddr: Bob,
		Function:      "function",
	}
	require.Equal(t, expectedInput, input)
}

func TestAsyncContext_CreateCallbackInput_DestinationCallSuccessful(t *testing.T) {
	host, _, originalVMInput := InitializeArwenAndWasmer_AsyncContext_AliceAndBob()
	host.Runtime().InitStateFromInput(originalVMInput)
	async := NewAsyncContext(host)
	call := &arwen.AsyncCall{
		Destination:     Bob,
		GasLocked:       82,
		ValueBytes:      big.NewInt(88).Bytes(),
		SuccessCallback: "successCallback",
		ErrorCallback:   "errorCallback",
		Status:          arwen.AsyncCallResolved,
	}
	vmOutput := &vmcommon.VMOutput{
		ReturnCode: vmcommon.Ok,
		ReturnData: [][]byte{
			[]byte("first"),
			[]byte("second"),
			{},
			[]byte("third"),
		},
		ReturnMessage: "a message",
		GasRemaining:  12,
	}
	destinationErr := error(nil)
	callbackInput, err := async.createCallbackInput(call, vmOutput, destinationErr)
	require.Nil(t, err)

	dataLength := len("successCallback") + 1 + len("first") + len("second") + len("third")
	separators := 5
	expectedGasProvided := call.GasLocked + vmOutput.GasRemaining
	expectedGasProvided -= uint64(dataLength + separators)
	expectedGasProvided -= host.Metering().GasSchedule().ElrondAPICost.AsyncCallStep

	expectedInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: Bob,
			Arguments: [][]byte{
				{byte(vmcommon.Ok)},
				[]byte("first"),
				[]byte("second"),
				{},
				[]byte("third"),
			},
			CallValue:      big.NewInt(0),
			CallType:       vmcommon.AsynchronousCallBack,
			GasPrice:       originalVMInput.GasPrice,
			GasProvided:    expectedGasProvided,
			CurrentTxHash:  originalVMInput.CurrentTxHash,
			PrevTxHash:     originalVMInput.PrevTxHash,
			OriginalTxHash: originalVMInput.OriginalTxHash,
		},
		RecipientAddr: Alice,
		Function:      "successCallback",
	}
	require.Equal(t, expectedInput, callbackInput)
}

func TestAsyncContext_CreateCallbackInput_DestinationCallFailed(t *testing.T) {
	host, _, originalVMInput := InitializeArwenAndWasmer_AsyncContext_AliceAndBob()
	host.Runtime().InitStateFromInput(originalVMInput)
	async := NewAsyncContext(host)
	call := &arwen.AsyncCall{
		Destination:     Bob,
		GasLocked:       82,
		ValueBytes:      big.NewInt(88).Bytes(),
		SuccessCallback: "successCallback",
		ErrorCallback:   "errorCallback",
		Status:          arwen.AsyncCallRejected,
	}
	vmOutput := &vmcommon.VMOutput{
		ReturnCode:    vmcommon.UserError,
		ReturnData:    [][]byte{},
		ReturnMessage: "there was a user error",
		GasRemaining:  0,
	}
	destinationErr := arwen.ErrSignalError
	callbackInput, err := async.createCallbackInput(call, vmOutput, destinationErr)
	require.Nil(t, err)

	dataLength := len("errorCallback") + 1 + len(vmOutput.ReturnMessage)
	separators := 2
	expectedGasProvided := call.GasLocked + vmOutput.GasRemaining
	expectedGasProvided -= uint64(dataLength + separators)
	expectedGasProvided -= host.Metering().GasSchedule().ElrondAPICost.AsyncCallStep

	expectedInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: Bob,
			Arguments: [][]byte{
				{byte(vmcommon.UserError)},
				[]byte(vmOutput.ReturnMessage),
			},
			CallValue:      big.NewInt(0),
			CallType:       vmcommon.AsynchronousCallBack,
			GasPrice:       originalVMInput.GasPrice,
			GasProvided:    expectedGasProvided,
			CurrentTxHash:  originalVMInput.CurrentTxHash,
			PrevTxHash:     originalVMInput.PrevTxHash,
			OriginalTxHash: originalVMInput.OriginalTxHash,
		},
		RecipientAddr: Alice,
		Function:      "errorCallback",
	}
	require.Equal(t, expectedInput, callbackInput)
}

func TestAsyncContext_CreateCallbackInput_NotEnoughGas(t *testing.T) {
	// Due to dynamic gas locking, this situation should never happen
	host, _, originalVMInput := InitializeArwenAndWasmer_AsyncContext_AliceAndBob()
	host.Runtime().InitStateFromInput(originalVMInput)
	async := NewAsyncContext(host)
	call := &arwen.AsyncCall{
		Destination:     Bob,
		GasLocked:       0,
		ValueBytes:      big.NewInt(88).Bytes(),
		SuccessCallback: "successCallback",
		ErrorCallback:   "errorCallback",
		Status:          arwen.AsyncCallRejected,
	}
	vmOutput := &vmcommon.VMOutput{
		ReturnCode:    vmcommon.UserError,
		ReturnData:    [][]byte{},
		ReturnMessage: "there was a user error",
		GasRemaining:  0,
	}
	destinationErr := arwen.ErrSignalError
	callbackInput, err := async.createCallbackInput(call, vmOutput, destinationErr)
	require.Nil(t, callbackInput)
	require.True(t, errors.Is(err, arwen.ErrNotEnoughGas))
}

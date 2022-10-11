package contexts

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/builtInFunctions"
	"github.com/ElrondNetwork/elrond-vm-common/parsers"
	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/config"
	"github.com/ElrondNetwork/wasm-vm/crypto/factory"
	contextmock "github.com/ElrondNetwork/wasm-vm/mock/context"
	worldmock "github.com/ElrondNetwork/wasm-vm/mock/world"
	"github.com/ElrondNetwork/wasm-vm/wasmer"
	"github.com/stretchr/testify/require"
)

var mockWasmerInstance *wasmer.Instance
var OriginalCaller = []byte("address_original_caller")
var Alice = []byte("address_alice")
var Bob = []byte("address_bob")

const GasForAsyncStep = config.GasValueForTests

var marshalizer = &marshal.GogoProtoMarshalizer{}

func makeAsyncContext(t testing.TB, host arwen.VMHost, address []byte) *asyncContext {
	callParser := parsers.NewCallArgsParser()
	esdtParser, _ := parsers.NewESDTTransferParser(worldmock.WorldMarshalizer)
	async, err := NewAsyncContext(
		host,
		callParser,
		esdtParser,
		marshalizer,
	)
	require.Nil(t, err)
	require.NotNil(t, async)

	async.address = address

	return async
}

func initializeArwenAndWasmer_AsyncContext() (*contextmock.VMHostMock, *worldmock.MockWorld) {
	imports := MakeAPIImports()
	_ = wasmer.SetImports(imports)

	vmType := []byte("type")

	gasSchedule := config.MakeGasMapForTests()
	gasCostConfig, _ := config.CreateGasConfig(gasSchedule)
	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)

	host := &contextmock.VMHostMock{}
	host.SCAPIMethods = imports

	mockMetering := &contextmock.MeteringContextMock{GasLeftMock: 10000}
	mockMetering.SetGasSchedule(gasSchedule)
	host.MeteringContext = mockMetering

	world := worldmock.NewMockWorld()
	host.BlockchainContext, _ = NewBlockchainContext(host, world)

	mockWasmerInstance = &wasmer.Instance{
		Exports: make(wasmer.ExportsMap),
	}
	runtimeContext, _ := NewRuntimeContext(
		host,
		vmType,
		builtInFunctions.NewBuiltInFunctionContainer(),
		wasmer.NewExecutor(),
	)
	runtimeContext.instance = mockWasmerInstance
	host.RuntimeContext = runtimeContext

	storageContext, _ := NewStorageContext(host, world, elrondReservedTestPrefix)
	host.StorageContext = storageContext

	host.OutputContext, _ = NewOutputContext(host)
	host.CryptoHook = factory.NewVMCrypto()
	host.StorageContext, _ = NewStorageContext(host, world, elrondReservedTestPrefix)
	host.EnableEpochsHandlerField = worldmock.EnableEpochsHandlerStubNoFlags()

	return host, world
}

func initializeArwenAndWasmer_AsyncContextWithAliceAndBob() (
	*contextmock.VMHostMock,
	*worldmock.MockWorld,
	*vmcommon.ContractCallInput,
) {
	host, world := initializeArwenAndWasmer_AsyncContext()
	world.AcctMap.PutAccount(&worldmock.Account{
		Address: Alice,
		Balance: big.NewInt(88),
	})
	world.AcctMap.PutAccount(&worldmock.Account{
		Address: Bob,
		Balance: big.NewInt(12),
	})

	originalVMInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:     OriginalCaller,
			Arguments:      nil,
			CallType:       vm.DirectCall,
			GasPrice:       1,
			CurrentTxHash:  []byte("txhash"),
			PrevTxHash:     []byte("txhash"),
			OriginalTxHash: []byte("txhash"),
		},
		RecipientAddr: Alice,
		Function:      "alice_function",
	}

	return host, world, originalVMInput
}

func TestNewAsyncContext(t *testing.T) {
	host, _ := initializeArwenAndWasmer_AsyncContext()
	async := makeAsyncContext(t, host, nil)

	require.NotNil(t, async.host)
	require.Nil(t, async.stateStack)
	require.Nil(t, async.callerAddr)
	require.Nil(t, async.returnData)
	require.NotNil(t, async.asyncCallGroups)
	require.Empty(t, async.asyncCallGroups)
}

func TestAsyncContext_InitState(t *testing.T) {
	host, _ := initializeArwenAndWasmer_AsyncContext()
	async := makeAsyncContext(t, host, nil)

	async.callerAddr = []byte("some address")
	async.returnData = []byte("some return data")
	async.asyncCallGroups = nil

	async.InitState()

	require.NotNil(t, async.callerAddr)
	require.Empty(t, async.callerAddr)
	require.NotNil(t, async.returnData)
	require.Empty(t, async.returnData)
	require.NotNil(t, async.asyncCallGroups)
	require.Empty(t, async.asyncCallGroups)
}

func TestAsyncContext_InitStateFromContractCallInput(t *testing.T) {
	contract := []byte("contract")
	host, _ := initializeArwenAndWasmer_AsyncContext()
	async := makeAsyncContext(t, host, nil)

	async.callerAddr = []byte("some address")
	async.returnData = []byte("some return data")
	async.asyncCallGroups = nil

	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: []byte("real caller addr"),
			GasPrice:   42,
		},
		RecipientAddr: contract,
	}

	host.Runtime().InitStateFromContractCallInput(input)
	async.InitStateFromInput(input)

	require.Equal(t, input.CallerAddr, async.callerAddr)
	require.NotNil(t, async.returnData)
	require.Empty(t, async.returnData)
	require.NotNil(t, async.asyncCallGroups)
	require.Empty(t, async.asyncCallGroups)
	require.Equal(t, contract, async.address)
}

func TestAsyncContext_GettersAndSetters(t *testing.T) {
	host, _ := initializeArwenAndWasmer_AsyncContext()
	async := makeAsyncContext(t, host, nil)

	async.callerAddr = []byte("some address")
	async.returnData = []byte("some return data")

	require.Equal(t, []byte("some address"), async.GetCallerAddress())
	require.Equal(t, []byte("some return data"), async.GetReturnData())

	async.SetReturnData([]byte("rockets"))
	require.Equal(t, []byte("rockets"), async.GetReturnData())
}

func TestAsyncContext_RegisterAsyncCall_NewGroup_DeleteGroup(t *testing.T) {
	host, _ := initializeArwenAndWasmer_AsyncContext()

	async := makeAsyncContext(t, host, nil)

	require.False(t, async.HasPendingCallGroups())

	group, exists := async.GetCallGroup("testGroup")
	require.Nil(t, group)
	require.False(t, exists)

	err := async.RegisterAsyncCall("testGroup", &arwen.AsyncCall{
		Destination: []byte("somewhere"),
		Data:        []byte("something"),
	})
	require.Nil(t, err)
	require.True(t, async.HasPendingCallGroups())
	require.Equal(t, 1, len(async.asyncCallGroups))

	group, exists = async.GetCallGroup("testGroup")
	require.NotNil(t, group)
	require.True(t, exists)
	require.True(t, group.HasPendingCalls())
	require.False(t, group.HasCallback())

	async.deleteCallGroupByID("testGroup")
	group, exists = async.GetCallGroup("testGroup")
	require.Nil(t, group)
	require.False(t, exists)
}

func TestAsyncContext_RegisterAsyncCall_ExistingGroup(t *testing.T) {
	host, _ := initializeArwenAndWasmer_AsyncContext()

	async := makeAsyncContext(t, host, nil)
	require.Equal(t, 0, len(async.asyncCallGroups))

	err := async.AddCallGroup(arwen.NewAsyncCallGroup("testGroup"))
	require.Nil(t, err)
	require.Equal(t, 1, len(async.asyncCallGroups))
	require.True(t, async.HasPendingCallGroups())

	group, exists := async.GetCallGroup("testGroup")
	require.NotNil(t, group)
	require.True(t, exists)

	err = async.RegisterAsyncCall("testGroup", &arwen.AsyncCall{
		Destination: []byte("somewhere"),
		Data:        []byte("something"),
	})
	require.Nil(t, err)
	require.Equal(t, 1, len(async.asyncCallGroups))
	require.True(t, async.HasPendingCallGroups())
}

func TestAsyncContext_RegisterAsyncCall_ValidationAndFields(t *testing.T) {
	// TODO execution mode
	// TODO non-nil destination
	// TODO locked gas
}

func TestAsyncContext_DetermineExecutionMode(t *testing.T) {
	leftAddress := []byte("left")
	leftAccount := &worldmock.Account{
		Address: leftAddress,
		Code:    []byte("left code"),
		ShardID: 0,
	}

	rightAddress := []byte("right")
	rightAccount := &worldmock.Account{
		Address: rightAddress,
		Code:    []byte("right code"),
		ShardID: 0,
	}

	host, world := initializeArwenAndWasmer_AsyncContext()
	world.AcctMap.PutAccount(leftAccount)
	world.AcctMap.PutAccount(rightAccount)
	runtime := host.Runtime()

	async := makeAsyncContext(t, host, nil)

	initRuntime(runtime, leftAddress)
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
	initRuntime(runtime, leftAddress)
	execMode, err = async.determineExecutionMode(rightAddress, []byte("func"))
	require.Nil(t, err)
	require.Equal(t, arwen.AsyncBuiltinFuncIntraShard, execMode)

	host.IsBuiltinFunc = false
	rightAccount.Code = []byte{}
	rightAccount.ShardID = 1

	// Erase the code of the rightAccount from the Output context cache
	outputAccount, _ := host.Output().GetOutputAccount(rightAddress)
	outputAccount.Code = []byte{}

	initRuntime(runtime, leftAddress)
	execMode, err = async.determineExecutionMode(rightAddress, []byte("func"))
	require.Nil(t, err)
	require.Equal(t, arwen.AsyncUnknown, execMode)

	host.IsBuiltinFunc = true
	rightAccount.Code = []byte{}
	rightAccount.ShardID = 1

	// Erase the code of the rightAccount from the Output context cache
	outputAccount, _ = host.Output().GetOutputAccount(rightAddress)
	outputAccount.Code = []byte{}

	initRuntime(runtime, leftAddress)
	execMode, err = async.determineExecutionMode(rightAddress, []byte("func"))
	require.Nil(t, err)
	require.Equal(t, arwen.AsyncBuiltinFuncCrossShard, execMode)
}

func initRuntime(runtime arwen.RuntimeContext, address []byte) {
	runtime.InitStateFromContractCallInput(&vmcommon.ContractCallInput{
		RecipientAddr: address,
	})
}

func TestAsyncContext_IsValidCallbackName(t *testing.T) {
	host, _ := initializeArwenAndWasmer_AsyncContext()
	async := makeAsyncContext(t, host, nil)

	mockWasmerInstance.Exports["a"] = nil
	mockWasmerInstance.Exports["my_contract_method_22"] = nil
	mockWasmerInstance.Exports["not_builtin"] = nil
	mockWasmerInstance.Exports["callBack"] = nil
	mockWasmerInstance.Exports["callback"] = nil
	mockWasmerInstance.Exports["function_do"] = nil

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

func TestAsyncContext_UpdateCurrentCallStatus(t *testing.T) {
	contract := []byte("contract")

	vmInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: []byte("caller"),
			Arguments:  [][]byte{{0}},
			CallType:   vm.DirectCall,
		},
		RecipientAddr: contract,
	}

	host, _ := initializeArwenAndWasmer_AsyncContext()
	async := makeAsyncContext(t, host, contract)

	storedAsync := &asyncContext{}
	storedAsync.host = host
	storedAsync.Save()

	// CallType == DirectCall, async.UpdateCurrentCallStatus() does nothing
	host.Runtime().InitStateFromContractCallInput(vmInput)
	asyncCall, isLegacy, err := async.UpdateCurrentAsyncCallStatus(contract, []byte{}, &vmInput.VMInput)
	require.Nil(t, asyncCall)
	require.False(t, isLegacy)
	require.Nil(t, err)

	// CallType == AsynchronousCall, async.UpdateCurrentCallStatus() does nothing
	vmInput.CallType = vm.AsynchronousCall
	host.Runtime().InitStateFromContractCallInput(vmInput)
	asyncCall, isLegacy, err = async.UpdateCurrentAsyncCallStatus(contract, []byte{}, &vmInput.VMInput)
	require.Nil(t, asyncCall)
	require.False(t, isLegacy)
	require.Nil(t, err)

	// CallType == AsynchronousCallback, but no AsyncCalls registered in the
	// AsyncContext, so async.UpdateCurrentCallStatus() returns an error
	vmInput.CallType = vm.AsynchronousCallBack
	vmInput.Arguments = nil
	host.Runtime().InitStateFromContractCallInput(vmInput)
	asyncCall, isLegacy, err = async.UpdateCurrentAsyncCallStatus(contract, []byte{}, &vmInput.VMInput)
	require.Nil(t, asyncCall)
	require.False(t, isLegacy)
	require.True(t, errors.Is(err, arwen.ErrCannotInterpretCallbackArgs))

	// CallType == AsynchronousCallback, but no AsyncCalls registered in the
	// AsyncContext, so async.UpdateCurrentCallStatus() returns an error
	vmInput.CallType = vm.AsynchronousCallBack
	vmInput.Arguments = [][]byte{{0}}
	host.Runtime().InitStateFromContractCallInput(vmInput)
	asyncCall, isLegacy, err = async.UpdateCurrentAsyncCallStatus(contract, []byte{}, &vmInput.VMInput)
	require.Equal(t, asyncCall, &arwen.AsyncCall{
		Status:          arwen.AsyncCallResolved,
		Destination:     contract,
		SuccessCallback: arwen.CallbackFunctionName,
		ErrorCallback:   arwen.CallbackFunctionName,
		GasLimit:        vmInput.GasProvided,
		GasLocked:       vmInput.GasLocked,
	})
	require.True(t, isLegacy)
	require.Nil(t, err)

	// CallType == AsynchronousCallback, and there is an AsyncCall registered,
	// but it's not the expected one.
	err = async.RegisterAsyncCall("testGroup", &arwen.AsyncCall{
		CallID:      []byte("callID_1"),
		Destination: []byte("some_address"),
		Data:        []byte("function"),
	})
	require.Nil(t, err)
	async.Save()

	vmInput.CallType = vm.AsynchronousCallBack
	vmInput.Arguments = [][]byte{{0}}
	host.Runtime().InitStateFromContractCallInput(vmInput)
	asyncCall, isLegacy, err = async.UpdateCurrentAsyncCallStatus(contract, []byte("callID_2"), &vmInput.VMInput)
	require.Equal(t, asyncCall, &arwen.AsyncCall{
		Status:          arwen.AsyncCallResolved,
		Destination:     contract,
		SuccessCallback: arwen.CallbackFunctionName,
		ErrorCallback:   arwen.CallbackFunctionName,
		GasLimit:        vmInput.GasProvided,
		GasLocked:       vmInput.GasLocked,
	})
	require.True(t, isLegacy)
	require.Nil(t, err)

	// CallType == AsynchronousCallback, but this time there is a corresponding AsyncCall
	// registered, causing async.UpdateCurrentCallStatus() to find and update the AsyncCall
	err = async.RegisterAsyncCall("testGroup", &arwen.AsyncCall{
		Destination: vmInput.CallerAddr,
		Data:        []byte("function"),
	})
	require.Nil(t, err)

	asyncCtx := &asyncContext{
		asyncCallGroups: []*arwen.AsyncCallGroup{
			{
				Identifier: "",
				AsyncCalls: []*arwen.AsyncCall{
					{
						Destination: []byte("caller"),
					},
				},
			},
		},
	}
	asyncCtx.host = host
	asyncCtx.Save()

	vmInput.CallType = vm.AsynchronousCallBack
	vmInput.Arguments = [][]byte{{0}}
	host.Runtime().InitStateFromContractCallInput(vmInput)
	asyncCall, isLegacy, err = async.UpdateCurrentAsyncCallStatus(contract, []byte{}, &vmInput.VMInput)
	require.Nil(t, err)
	require.False(t, isLegacy)
	require.NotNil(t, asyncCall)
	require.Equal(t, arwen.AsyncCallResolved, asyncCall.Status)

	// CallType == AsynchronousCallback, there is a corresponding AsyncCall
	// registered, causing async.UpdateCurrentCallStatus() to find and update the
	// AsyncCall, but with AsyncCallRejected
	vmInput.CallType = vm.AsynchronousCallBack
	vmInput.Arguments = [][]byte{{1}}
	host.Runtime().InitStateFromContractCallInput(vmInput)
	asyncCall, isLegacy, err = async.UpdateCurrentAsyncCallStatus(contract, []byte{}, &vmInput.VMInput)
	require.Nil(t, err)
	require.False(t, isLegacy)
	require.NotNil(t, asyncCall)
	require.Equal(t, arwen.AsyncCallRejected, asyncCall.Status)
}

func TestAsyncContext_SendAsyncCallCrossShard(t *testing.T) {
	host, world := initializeArwenAndWasmer_AsyncContext()
	world.AcctMap.PutAccount(&worldmock.Account{
		Address: []byte("smartcontract"),
		Balance: big.NewInt(88),
	})

	initRuntime(host.Runtime(), []byte("smartcontract"))
	async := makeAsyncContext(t, host, nil)

	asyncCall := &arwen.AsyncCall{
		Destination: []byte("destination"),
		GasLimit:    42,
		GasLocked:   98,
		ValueBytes:  big.NewInt(88).Bytes(),
		Data:        []byte("some_data"),
	}

	host.Runtime().GetVMInput().GasProvided = 200

	err := async.sendAsyncCallCrossShard(asyncCall)
	require.Nil(t, err)

	mockMetering := host.Metering().(*contextmock.MeteringContextMock)
	mockMetering.GasProvidedMock = 200
	mockMetering.GasLeftMock = 60

	vmOutput := host.Output().GetVMOutput()
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	smartcontract, ok := vmOutput.OutputAccounts["smartcontract"]
	require.True(t, ok)
	require.Equal(t, big.NewInt(-88), smartcontract.BalanceDelta)
	require.Empty(t, smartcontract.OutputTransfers)

	destination, ok := vmOutput.OutputAccounts["destination"]
	require.True(t, ok)
	require.Equal(t, big.NewInt(88), destination.BalanceDelta)
	require.Len(t, destination.OutputTransfers, 1)

	asyncTransfer := destination.OutputTransfers[0]
	require.Equal(t, big.NewInt(88), asyncTransfer.Value)
	require.Equal(t, uint64(42), asyncTransfer.GasLimit)
	require.Equal(t, uint64(98), asyncTransfer.GasLocked)

	callParser := parsers.NewCallArgsParser()
	function, _, _ := callParser.ParseData(string(asyncTransfer.Data))
	require.Equal(t, "some_data", function)
	require.Equal(t, vm.AsynchronousCall, asyncTransfer.CallType)
}

func TestAsyncContext_ExecuteSyncCall_EarlyOutOfGas(t *testing.T) {
	// Scenario 1
	// Assert error propagation in async.executeSyncCall() from
	// async.createContractCallInput()
	host, _, originalVMInput := initializeArwenAndWasmer_AsyncContextWithAliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := makeAsyncContext(t, host, nil)

	asyncCall := defaultAsyncCall_AliceToBob()
	asyncCall.Data = []byte("function")
	asyncCall.GasLimit = 1
	err := async.executeAsyncLocalCall(asyncCall)
	require.True(t, errors.Is(err, arwen.ErrNotEnoughGas))
}

func TestAsyncContext_ExecuteSyncCall_Successful(t *testing.T) {
	// Scenario 3
	// Successful execution at destination, and successful callback execution;
	// the AsyncCall contains sufficient gas this time.
	host, _, originalVMInput := initializeArwenAndWasmer_AsyncContextWithAliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)

	mockWasmerInstance.Exports["successCallback"] = nil
	mockWasmerInstance.Exports["errorCallback"] = nil

	async := makeAsyncContext(t, host, Alice)

	asyncCall := defaultAsyncCall_AliceToBob()
	asyncCall.GasLimit = 200
	asyncCall.GasLocked = 90
	gasConsumedByDestination := uint64(23)
	gasConsumedByCallback := uint64(22)

	// The expected input passed to host.ExecuteOnDestContext() to call Bob as destination
	destInput := defaultCallInput_AliceToBob(originalVMInput)
	destInput.GasProvided = asyncCall.GasLimit
	destInput.GasLocked = asyncCall.GasLocked

	// Prepare the output of Bob (the destination call)
	destOutput := defaultDestOutput_Ok()
	destOutput.GasRemaining = destInput.GasProvided - gasConsumedByDestination

	// Prepare the input to Alice's callback
	callbackInput := defaultCallbackInput_BobToAlice(originalVMInput)
	callbackInput.GasProvided = destOutput.GasRemaining + asyncCall.GasLocked
	callbackInput.GasProvided -= defaultOutputDataLengthAsArgs(asyncCall, destOutput)
	callbackInput.GasProvided -= GasForAsyncStep
	callbackInput.GasLocked = 0

	// Prepare the output of Alice's callback
	callbackOutput := defaultCallbackOutput_Ok()
	callbackOutput.GasRemaining = callbackInput.GasProvided - gasConsumedByCallback

	// Enqueue the prepared VMOutputs
	host.EnqueueVMOutput(destOutput)
	host.EnqueueVMOutput(callbackOutput)

	host.Metering().RestoreGas(10000)

	err := async.RegisterAsyncCall("test", asyncCall)
	require.Nil(t, err)
	err = async.executeAsyncLocalCall(asyncCall)
	require.Nil(t, err)
	require.Equal(t, arwen.AsyncCallResolved, asyncCall.Status)

	// There were two calls to host.ExecuteOnDestContext()
	require.Len(t, host.StoredInputs, 2)

	host.StoredInputs[0].AsyncArguments = nil
	require.Equal(t, destInput, host.StoredInputs[0])
	host.StoredInputs[1].AsyncArguments = nil
	require.Equal(t, callbackInput, host.StoredInputs[1])

	// Verify the final output of the execution; GasRemaining is set to 0 because
	// the test uses a mocked host.ExecuteOnDestContext(), which does not know to
	// manipulate the state stack of the OutputContext, therefore VMOutputs are
	// not merged between executions.
	expectedOutput := arwen.MakeEmptyVMOutput()
	expectedOutput.ReturnCode = vmcommon.Ok
	expectedOutput.GasRemaining = host.Metering().GasLeft()

	// The expectedOutput must also contain an OutputAccount corresponding to
	// Alice, because of a call to host.Output().GetOutputAccount() in
	// host.Output().GetVMOutput(), which creates and caches an empty account for
	// her.
	arwen.AddNewOutputAccount(expectedOutput, Alice, 0, nil)

	host.Output().GetOutputAccount(Alice)
	actualOutput := host.Output().GetVMOutput()
	// Bob entry is retuned empty by the MockWorld.GetStorageData()
	delete(actualOutput.OutputAccounts, string(Bob))
	require.Equal(t, expectedOutput, actualOutput)
}

func TestAsyncContext_CreateContractCallInput(t *testing.T) {
	host, _, originalVMInput := initializeArwenAndWasmer_AsyncContextWithAliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := makeAsyncContext(t, host, nil)
	asyncCall := &arwen.AsyncCall{
		Destination: Bob,
		ValueBytes:  big.NewInt(88).Bytes(),
	}

	asyncCall.Data = []byte{}
	input, err := async.createContractCallInput(asyncCall)
	require.Nil(t, input)
	require.Error(t, err)

	asyncCall.Data = []byte("function")
	asyncCall.GasLimit = 1
	input, err = async.createContractCallInput(asyncCall)
	require.Nil(t, input)
	require.True(t, errors.Is(err, arwen.ErrNotEnoughGas))

	asyncCall.Data = []byte("function@0A0B0C@03")
	asyncCall.GasLimit = 2
	input, err = async.createContractCallInput(asyncCall)
	require.Nil(t, err)
	require.NotNil(t, input)

	expectedInput := defaultCallInput_AliceToBob(originalVMInput)
	expectedInput.GasProvided = 2
	input.AsyncArguments = nil
	require.Equal(t, expectedInput, input)
}

func TestAsyncContext_CreateCallbackInput_DestinationCallSuccessful(t *testing.T) {
	host, _, originalVMInput := initializeArwenAndWasmer_AsyncContextWithAliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := makeAsyncContext(t, host, Alice)

	asyncCall := defaultAsyncCall_AliceToBob()
	asyncCall.Status = arwen.AsyncCallResolved
	asyncCall.GasLocked = 300

	vmOutput := defaultDestOutput_Ok()
	vmOutput.GasRemaining = 12

	destinationErr := error(nil)
	callbackInput, err := async.createCallbackInput(asyncCall, vmOutput, 0, destinationErr)
	require.Nil(t, err)

	expectedGasProvided := asyncCall.GasLocked + vmOutput.GasRemaining
	expectedGasProvided -= defaultOutputDataLengthAsArgs(asyncCall, vmOutput)
	expectedGasProvided -= host.Metering().GasSchedule().ElrondAPICost.AsyncCallStep

	expectedInput := defaultCallbackInput_BobToAlice(originalVMInput)
	expectedInput.GasProvided = expectedGasProvided
	callbackInput.AsyncArguments = nil
	require.Equal(t, expectedInput, callbackInput)
}

func TestAsyncContext_CreateCallbackInput_DestinationCallFailed(t *testing.T) {
	host, _, originalVMInput := initializeArwenAndWasmer_AsyncContextWithAliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := makeAsyncContext(t, host, Alice)

	asyncCall := defaultAsyncCall_AliceToBob()
	asyncCall.Status = arwen.AsyncCallRejected
	asyncCall.GasLocked = 200

	vmOutput := defaultDestOutput_UserError()
	destinationErr := arwen.ErrSignalError
	callbackInput, err := async.createCallbackInput(asyncCall, vmOutput, 0, destinationErr)
	require.Nil(t, err)

	expectedGasProvided := asyncCall.GasLocked + vmOutput.GasRemaining
	expectedGasProvided -= defaultOutputDataLengthAsArgs(asyncCall, vmOutput)
	expectedGasProvided -= host.Metering().GasSchedule().ElrondAPICost.AsyncCallStep

	expectedInput := arwen.MakeContractCallInput(Bob, Alice, "errorCallback", 0)
	arwen.AddArgument(expectedInput, []byte{byte(vmcommon.UserError)})
	arwen.AddArgument(expectedInput, []byte(vmOutput.ReturnMessage))
	arwen.CopyTxHashes(expectedInput, originalVMInput)
	expectedInput.GasProvided = expectedGasProvided
	expectedInput.CallType = vm.AsynchronousCallBack
	expectedInput.ReturnCallAfterError = true
	callbackInput.AsyncArguments = nil
	require.Equal(t, expectedInput, callbackInput)
}

func TestAsyncContext_CreateCallbackInput_NotEnoughGas(t *testing.T) {
	// Due to dynamic gas locking, this situation should never happen
	host, _, originalVMInput := initializeArwenAndWasmer_AsyncContextWithAliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := makeAsyncContext(t, host, nil)

	asyncCall := defaultAsyncCall_AliceToBob()
	asyncCall.Status = arwen.AsyncCallRejected

	vmOutput := &vmcommon.VMOutput{
		ReturnCode:    vmcommon.UserError,
		ReturnData:    [][]byte{},
		ReturnMessage: "there was a user error",
		GasRemaining:  0,
	}

	destinationErr := arwen.ErrSignalError
	callbackInput, err := async.createCallbackInput(asyncCall, vmOutput, 0, destinationErr)
	require.Nil(t, callbackInput)
	require.True(t, errors.Is(err, arwen.ErrNotEnoughGas))
}

func TestAsyncContext_FinishSyncExecution_NilError_NilVMOutput(t *testing.T) {
	host, _, originalVMInput := initializeArwenAndWasmer_AsyncContextWithAliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := makeAsyncContext(t, host, nil)
	async.finishAsyncLocalCallbackExecution(nil, nil, 0)

	// The expectedOutput must also contain an OutputAccount corresponding to
	// Alice, because of a call to host.Output().GetOutputAccount() in
	// host.Output().GetVMOutput(), which creates and caches an empty account for
	// her.
	expectedOutput := arwen.MakeEmptyVMOutput()
	expectedOutput.GasRemaining = host.Metering().GasLeft()
	arwen.AddNewOutputAccount(expectedOutput, Alice, 0, nil)

	host.Output().GetOutputAccount(Alice)
	require.Equal(t, expectedOutput, host.Output().GetVMOutput())
}

func TestAsyncContext_FinishSyncExecution_Error_NilVMOutput(t *testing.T) {
	host, _, originalVMInput := initializeArwenAndWasmer_AsyncContextWithAliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := makeAsyncContext(t, host, nil)

	syncExecErr := arwen.ErrNotEnoughGas
	async.finishAsyncLocalCallbackExecution(nil, syncExecErr, 0)

	expectedOutput := arwen.MakeEmptyVMOutput()
	expectedOutput.GasRemaining = host.Metering().GasLeft()
	// expectedOutput.ReturnCode = vmcommon.OutOfGas
	// expectedOutput.ReturnMessage = syncExecErr.Error()
	// arwen.AddFinishData(expectedOutput, []byte(vmcommon.OutOfGas.String()))
	// arwen.AddFinishData(expectedOutput, originalVMInput.CurrentTxHash)

	// The expectedOutput must also contain an OutputAccount corresponding to
	// Alice, because of a call to host.Output().GetOutputAccount() in
	// host.Output().GetVMOutput(), which creates and caches an empty account for
	// her.
	arwen.AddNewOutputAccount(expectedOutput, Alice, 0, nil)

	host.Output().GetOutputAccount(Alice)
	require.Equal(t, expectedOutput, host.Output().GetVMOutput())
}

func TestAsyncContext_FinishSyncExecution_ErrorAndVMOutput(t *testing.T) {
	host, _, originalVMInput := initializeArwenAndWasmer_AsyncContextWithAliceAndBob()
	host.Runtime().InitStateFromContractCallInput(originalVMInput)
	async := makeAsyncContext(t, host, nil)

	syncExecOutput := arwen.MakeEmptyVMOutput()
	syncExecOutput.ReturnCode = vmcommon.UserError
	syncExecOutput.ReturnMessage = "user made an error"
	syncExecErr := arwen.ErrSignalError
	async.finishAsyncLocalCallbackExecution(syncExecOutput, syncExecErr, 0)

	expectedOutput := arwen.MakeEmptyVMOutput()
	expectedOutput.GasRemaining = host.Metering().GasLeft()
	// expectedOutput.ReturnCode = vmcommon.UserError
	// expectedOutput.ReturnMessage = "user made an error"
	// arwen.AddFinishData(expectedOutput, []byte(vmcommon.UserError.String()))
	// arwen.AddFinishData(expectedOutput, originalVMInput.CurrentTxHash)

	// The expectedOutput must also contain an OutputAccount corresponding to
	// Alice, because of a call to host.Output().GetOutputAccount() in
	// host.Output().GetVMOutput(), which creates and caches an empty account for
	// her.
	arwen.AddNewOutputAccount(expectedOutput, Alice, 0, nil)

	host.Output().GetOutputAccount(Alice)
	require.Equal(t, expectedOutput, host.Output().GetVMOutput())
}

func defaultAsyncCall_AliceToBob() *arwen.AsyncCall {
	return &arwen.AsyncCall{
		Destination:     Bob,
		Data:            []byte("function@0A0B0C@03"),
		GasLimit:        0,
		GasLocked:       0,
		ValueBytes:      big.NewInt(88).Bytes(),
		SuccessCallback: "successCallback",
		ErrorCallback:   "errorCallback",
		Status:          arwen.AsyncCallPending,
	}
}

func defaultCallInput_AliceToBob(originalVMInput *vmcommon.ContractCallInput) *vmcommon.ContractCallInput {
	destInput := arwen.MakeContractCallInput(Alice, Bob, "function", 88)
	arwen.CopyTxHashes(destInput, originalVMInput)
	arwen.AddArgument(destInput, []byte{10, 11, 12})
	arwen.AddArgument(destInput, []byte{3})
	destInput.CallType = vm.AsynchronousCall

	return destInput
}

func defaultDestOutput_UserError() *vmcommon.VMOutput {
	return &vmcommon.VMOutput{
		ReturnCode:    vmcommon.UserError,
		ReturnData:    [][]byte{},
		ReturnMessage: "user error occurred",
		GasRemaining:  0,
	}
}

func defaultDestOutput_Ok() *vmcommon.VMOutput {
	return &vmcommon.VMOutput{
		ReturnCode: vmcommon.Ok,
		ReturnData: [][]byte{
			[]byte("first"),
			[]byte("second"),
			{},
			[]byte("third"),
		},
		ReturnMessage: "a message",
		GasRemaining:  0,
	}
}

func defaultCallbackInput_BobToAlice(originalVMInput *vmcommon.ContractCallInput) *vmcommon.ContractCallInput {
	input := arwen.MakeContractCallInput(Bob, Alice, "successCallback", 0)
	arwen.AddArgument(input, []byte{0}) // vmcommon.Ok
	arwen.AddArgument(input, []byte("first"))
	arwen.AddArgument(input, []byte("second"))
	arwen.AddArgument(input, []byte{})
	arwen.AddArgument(input, []byte("third"))
	arwen.CopyTxHashes(input, originalVMInput)
	input.CallType = vm.AsynchronousCallBack
	return input
}

func defaultCallbackOutput_Ok() *vmcommon.VMOutput {
	vmOutput := arwen.MakeEmptyVMOutput()
	arwen.AddFinishData(vmOutput, []byte("cbFirst"))
	arwen.AddFinishData(vmOutput, []byte("cbSecond"))

	return vmOutput
}

func defaultOutputDataLengthAsArgs(asyncCall *arwen.AsyncCall, vmOutput *vmcommon.VMOutput) uint64 {
	separator := 1
	hexSize := 2
	returnCode := 1 * hexSize

	dataLength := 0
	if vmOutput.ReturnCode == vmcommon.Ok {
		dataLength += returnCode + len(asyncCall.SuccessCallback) + separator
		for _, data := range vmOutput.ReturnData {
			dataLength += separator
			dataLength += len(data) * hexSize
		}
	} else {
		dataLength += len(asyncCall.ErrorCallback) + separator + returnCode
		dataLength += separator + len(vmOutput.ReturnMessage)*hexSize
	}

	return uint64(dataLength)
}

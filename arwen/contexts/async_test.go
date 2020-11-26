package contexts

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

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

	blockchainHookMock := mock.NewBlockchainHookMock()
	host.BlockchainContext, _ = NewBlockchainContext(host, blockchainHookMock)
	host.RuntimeContext, _ = NewRuntimeContext(host, []byte("vm"), false)
	host.OutputContext, _ = NewOutputContext(host)
	host.CryptoHook = crypto.NewVMCrypto()
	return host, blockchainHookMock
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

func TestAsyncContext_AddCall_NewGroup(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	vmType := []byte("type")
	runtime, err := NewRuntimeContext(host, vmType, false)
	require.Nil(t, err)
	host.RuntimeContext = runtime

	async := NewAsyncContext(host)
	require.NotNil(t, async)

	require.True(t, async.IsComplete())
	require.False(t, async.HasPendingCallGroups())

	group, exists := async.GetCallGroup("testGroup")
	require.Nil(t, group)
	require.False(t, exists)

	err = async.AddCall("testGroup", &arwen.AsyncCall{
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
}

func TestAsyncContext_AddCall_ExistingGroup(t *testing.T) {
	host, _ := InitializeArwenAndWasmer_AsyncContext()
	vmType := []byte("type")
	runtime, err := NewRuntimeContext(host, vmType, false)
	require.Nil(t, err)
	host.RuntimeContext = runtime

	async := NewAsyncContext(host)
	require.NotNil(t, async)
	require.Equal(t, 0, len(async.asyncCallGroups))

	async.addCallGroup(arwen.NewAsyncCallGroup("testGroup"))
	require.Equal(t, 1, len(async.asyncCallGroups))
	require.False(t, async.IsComplete())

	group, exists := async.GetCallGroup("testGroup")
	require.NotNil(t, group)
	require.True(t, exists)

	err = async.AddCall("testGroup", &arwen.AsyncCall{
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

func TestAsyncContext_SetGroupCallback_Errors(t *testing.T) {
}

func TestAsyncContext_SetGroupCallback_Success(t *testing.T) {
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

func TestAsyncContext_DeleteCallGroup(t *testing.T) {
}

func TestAsyncContext_IsValidCallbackName(t *testing.T) {
}

func TestAsyncContext_FindCall(t *testing.T) {
}

func TestAsyncContext_UpdateCurrentCallStatus(t *testing.T) {
}

func TestAsyncContext_PrepareLegacyAsyncCall(t *testing.T) {
}

func TestAsyncContext_SendAsyncCallCrossShard(t *testing.T) {
}

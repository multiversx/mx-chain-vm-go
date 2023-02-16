package contexts

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/hashing/blake2b"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
	"github.com/multiversx/mx-chain-vm-v1_4-go/config"
	"github.com/multiversx/mx-chain-vm-v1_4-go/crypto/factory"
	contextmock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/context"
	worldmock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/world"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/cryptoapi"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/mock"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/vmhooks"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/vmhooksmeta"
	"github.com/multiversx/mx-chain-vm-v1_4-go/wasmer"
	"github.com/stretchr/testify/require"
)

var defaultHasher = blake2b.NewBlake2b()

const counterWasmCode = "./../../test/contracts/counter/output/counter.wasm"

var vmType = []byte("type")

func MakeAPIImports() *wasmer.Imports {
	imports := vmhooksmeta.NewEIFunctions()
	_ = vmhooks.BaseOpsAPIImports(imports)
	_ = vmhooks.BigIntImports(imports)
	_ = vmhooks.BigFloatImports(imports)
	_ = vmhooks.ManagedBufferImports(imports)
	_ = vmhooks.SmallIntImports(imports)
	_ = cryptoapi.CryptoImports(imports)
	return wasmer.ConvertImports(imports)
}

func InitializeVMAndWasmer() *contextmock.VMHostMock {
	imports := MakeAPIImports()
	_ = wasmer.SetImports(imports)

	gasSchedule := config.MakeGasMapForTests()
	gasCostConfig, _ := config.CreateGasConfig(gasSchedule)
	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)

	host := &contextmock.VMHostMock{}
	host.SCAPIMethods = imports
	host.EnableEpochsHandlerField = &mock.EnableEpochsHandlerStub{}

	mockMetering := &contextmock.MeteringContextMock{}
	mockMetering.SetGasSchedule(gasSchedule)
	host.MeteringContext = mockMetering
	host.BlockchainContext, _ = NewBlockchainContext(host, worldmock.NewMockWorld())
	host.OutputContext, _ = NewOutputContext(host)
	host.CryptoHook = factory.NewVMCrypto()
	return host
}

func makeDefaultRuntimeContext(t *testing.T, host vmhost.VMHost) *runtimeContext {
	runtimeContext, err := NewRuntimeContext(
		host,
		vmType,
		builtInFunctions.NewBuiltInFunctionContainer(),
		defaultHasher,
	)
	require.Nil(t, err)
	require.NotNil(t, runtimeContext)

	return runtimeContext
}

func TestNewRuntimeContextErrors(t *testing.T) {
	host := InitializeVMAndWasmer()

	bfc := builtInFunctions.NewBuiltInFunctionContainer()
	hasher := defaultHasher

	t.Run("NilHost", func(t *testing.T) {
		runtimeContext, err := NewRuntimeContext(nil, vmType, bfc, hasher)
		require.Nil(t, runtimeContext)
		require.ErrorIs(t, err, vmhost.ErrNilHost)
	})
	t.Run("NilVMType", func(t *testing.T) {
		runtimeContext, err := NewRuntimeContext(host, nil, bfc, hasher)
		require.Nil(t, runtimeContext)
		require.ErrorIs(t, err, vmhost.ErrNilVMType)
	})
	t.Run("NilBuiltinFuncContainer", func(t *testing.T) {
		runtimeContext, err := NewRuntimeContext(host, vmType, nil, hasher)
		require.Nil(t, runtimeContext)
		require.ErrorIs(t, err, vmhost.ErrNilBuiltInFunctionsContainer)
	})
	t.Run("NilHasher", func(t *testing.T) {
		runtimeContext, err := NewRuntimeContext(host, vmType, bfc, nil)
		require.Nil(t, runtimeContext)
		require.ErrorIs(t, err, vmhost.ErrNilHasher)
	})
}

func TestNewRuntimeContext(t *testing.T) {
	host := InitializeVMAndWasmer()
	host.EnableEpochsHandlerField = &mock.EnableEpochsHandlerStub{}
	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()

	require.Equal(t, &vmcommon.ContractCallInput{}, runtimeContext.vmInput)
	require.Equal(t, []byte{}, runtimeContext.codeAddress)
	require.Equal(t, "", runtimeContext.callFunction)
	require.Equal(t, false, runtimeContext.readOnly)
	require.Equal(t, uint64(0), runtimeContext.GetSCCodeSize())
	require.Nil(t, runtimeContext.asyncCallInfo)
}

func TestRuntimeContext_InitState(t *testing.T) {
	host := InitializeVMAndWasmer()
	host.EpochsStub().IsRuntimeCodeSizeFixEnabledField = true

	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()

	runtimeContext.vmInput = nil
	runtimeContext.codeAddress = []byte("some address")
	runtimeContext.callFunction = "a function"
	runtimeContext.readOnly = true
	runtimeContext.iTracker.codeSize = 1024
	runtimeContext.asyncCallInfo = &vmhost.AsyncCallInfo{}

	runtimeContext.InitState()

	require.Equal(t, &vmcommon.ContractCallInput{}, runtimeContext.vmInput)
	require.Equal(t, []byte{}, runtimeContext.codeAddress)
	require.Equal(t, "", runtimeContext.callFunction)
	require.Equal(t, false, runtimeContext.readOnly)
	require.Equal(t, uint64(0), runtimeContext.iTracker.codeSize)
	require.Nil(t, runtimeContext.asyncCallInfo)
}

func TestRuntimeContext_CodeSizeFix(t *testing.T) {
	host := InitializeVMAndWasmer()
	epochs := host.EpochsStub()
	host.EnableEpochsHandlerField = epochs

	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()

	runtimeContext.codeSize = 1024
	runtimeContext.iTracker.codeSize = 1024

	epochs.IsRuntimeCodeSizeFixEnabledField = false
	runtimeContext.InitState()
	require.Equal(t, uint64(1024), runtimeContext.GetSCCodeSize())

	epochs.IsRuntimeCodeSizeFixEnabledField = true
	runtimeContext.InitState()
	require.Equal(t, uint64(0), runtimeContext.GetSCCodeSize())
}

func TestRuntimeContext_NewWasmerInstance(t *testing.T) {
	host := InitializeVMAndWasmer()
	epochs := host.EpochsStub()
	host.EnableEpochsHandlerField = epochs

	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()

	runtimeContext.SetMaxInstanceCount(1)

	gasLimit := uint64(100000000)
	var dummy []byte
	err := runtimeContext.StartWasmerInstance(dummy, gasLimit, false)
	require.NotNil(t, err)
	require.True(t, errors.Is(err, wasmer.ErrInvalidBytecode))

	gasLimit = uint64(100000000)
	dummy = []byte("contract")
	err = runtimeContext.StartWasmerInstance(dummy, gasLimit, false)
	require.NotNil(t, err)

	path := counterWasmCode
	contractCode := vmhost.GetSCCode(path)

	t.Run("fix code size disabled", func(t *testing.T) {
		epochs.IsRuntimeCodeSizeFixEnabledField = false
		runtimeContext := makeDefaultRuntimeContext(t, host)
		runtimeContext.SetMaxInstanceCount(1)
		defer runtimeContext.ClearWarmInstanceCache()
		err = runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
		require.Nil(t, err)
		require.Equal(t, vmhost.BreakpointNone, runtimeContext.GetRuntimeBreakpointValue())
		require.Equal(t, uint64(0), runtimeContext.GetSCCodeSize())
	})
	t.Run("fix code size enabled", func(t *testing.T) {
		epochs.IsRuntimeCodeSizeFixEnabledField = true
		runtimeContext := makeDefaultRuntimeContext(t, host)
		runtimeContext.SetMaxInstanceCount(1)
		defer runtimeContext.ClearWarmInstanceCache()
		err = runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
		require.Nil(t, err)
		require.Equal(t, vmhost.BreakpointNone, runtimeContext.GetRuntimeBreakpointValue())
		require.Equal(t, uint64(len(contractCode)), runtimeContext.GetSCCodeSize())
	})
}

func TestRuntimeContext_IsFunctionImported(t *testing.T) {
	host := InitializeVMAndWasmer()
	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()

	runtimeContext.SetMaxInstanceCount(1)

	gasLimit := uint64(100000000)
	path := counterWasmCode
	contractCode := vmhost.GetSCCode(path)
	err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)
	require.Equal(t, vmhost.BreakpointNone, runtimeContext.GetRuntimeBreakpointValue())

	// These API functions exist, and are imported by 'counter'
	require.True(t, runtimeContext.IsFunctionImported("int64storageLoad"))
	require.True(t, runtimeContext.IsFunctionImported("int64storageStore"))
	require.True(t, runtimeContext.IsFunctionImported("int64finish"))

	// These API functions exist, but are not imported by 'counter'
	require.False(t, runtimeContext.IsFunctionImported("transferValue"))
	require.False(t, runtimeContext.IsFunctionImported("executeOnSameContext"))
	require.False(t, runtimeContext.IsFunctionImported("asyncCall"))

	// These API functions don't even exist
	require.False(t, runtimeContext.IsFunctionImported(""))
	require.False(t, runtimeContext.IsFunctionImported("*"))
	require.False(t, runtimeContext.IsFunctionImported("$@%"))
	require.False(t, runtimeContext.IsFunctionImported("doesNotExist"))
}

func TestRuntimeContext_StateSettersAndGetters(t *testing.T) {
	imports := MakeAPIImports()
	host := &contextmock.VMHostMock{}
	host.EnableEpochsHandlerField = &mock.EnableEpochsHandlerStub{}
	host.SCAPIMethods = imports

	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()

	arguments := [][]byte{[]byte("argument 1"), []byte("argument 2")}
	esdtTransfer := &vmcommon.ESDTTransfer{
		ESDTValue:      big.NewInt(4242),
		ESDTTokenName:  []byte("random_token"),
		ESDTTokenType:  uint32(core.NonFungible),
		ESDTTokenNonce: 94,
	}

	vmInput := vmcommon.VMInput{
		CallerAddr:    []byte("caller"),
		Arguments:     arguments,
		CallValue:     big.NewInt(0),
		ESDTTransfers: []*vmcommon.ESDTTransfer{esdtTransfer},
	}
	callInput := &vmcommon.ContractCallInput{
		VMInput:       vmInput,
		RecipientAddr: []byte("recipient"),
		Function:      "test function",
	}

	runtimeContext.InitStateFromContractCallInput(callInput)
	require.Equal(t, []byte("caller"), runtimeContext.GetVMInput().CallerAddr)
	require.Equal(t, []byte("recipient"), runtimeContext.GetContextAddress())
	require.Equal(t, "test function", runtimeContext.Function())
	require.Equal(t, vmType, runtimeContext.GetVMType())
	require.Equal(t, arguments, runtimeContext.Arguments())

	runtimeInput := runtimeContext.GetVMInput()
	require.Zero(t, big.NewInt(4242).Cmp(runtimeInput.ESDTTransfers[0].ESDTValue))
	require.True(t, bytes.Equal([]byte("random_token"), runtimeInput.ESDTTransfers[0].ESDTTokenName))
	require.Equal(t, uint32(core.NonFungible), runtimeInput.ESDTTransfers[0].ESDTTokenType)
	require.Equal(t, uint64(94), runtimeInput.ESDTTransfers[0].ESDTTokenNonce)

	vmInput2 := vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: []byte("caller2"),
			Arguments:  arguments,
			CallValue:  big.NewInt(0),
		},
	}
	runtimeContext.SetVMInput(&vmInput2)
	require.Equal(t, []byte("caller2"), runtimeContext.GetVMInput().CallerAddr)

	runtimeContext.SetCodeAddress([]byte("smartcontract"))
	require.Equal(t, []byte("smartcontract"), runtimeContext.codeAddress)
}

func TestRuntimeContext_PushPopInstance(t *testing.T) {
	host := InitializeVMAndWasmer()
	epochs := host.EpochsStub()
	gasLimit := uint64(100000000)
	path := counterWasmCode
	contractCode := vmhost.GetSCCode(path)
	oldCodeSize := uint64(len(contractCode))
	newCodeSize := oldCodeSize + 84

	prepare := func() *runtimeContext {
		runtimeContext := makeDefaultRuntimeContext(t, host)
		runtimeContext.SetMaxInstanceCount(1)

		err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
		require.Nil(t, err)
		return runtimeContext
	}

	t.Run("fix code size disabled", func(t *testing.T) {
		epochs.IsRuntimeCodeSizeFixEnabledField = false
		runtimeContext := prepare()
		defer runtimeContext.ClearWarmInstanceCache()
		require.Equal(t, uint64(0), runtimeContext.GetSCCodeSize())
		instance := runtimeContext.iTracker.instance
		runtimeContext.pushInstance()
		runtimeContext.codeSize = newCodeSize
		runtimeContext.iTracker.codeSize = newCodeSize
		require.Equal(t, newCodeSize, runtimeContext.GetSCCodeSize())
		runtimeContext.iTracker.instance = &wasmer.Instance{}
		require.Equal(t, 1, len(runtimeContext.iTracker.instanceStack))

		runtimeContext.popInstance()
		require.Equal(t, newCodeSize, runtimeContext.GetSCCodeSize())
		require.Equal(t, oldCodeSize, runtimeContext.iTracker.codeSize)
		require.NotNil(t, runtimeContext.iTracker.instance)
		require.Equal(t, instance, runtimeContext.iTracker.instance)
		require.Equal(t, 0, len(runtimeContext.iTracker.instanceStack))

		runtimeContext.pushInstance()
		require.Equal(t, 1, len(runtimeContext.iTracker.instanceStack))
	})
	t.Run("fix code size enabled", func(t *testing.T) {
		epochs.IsRuntimeCodeSizeFixEnabledField = true
		runtimeContext := prepare()
		defer runtimeContext.ClearWarmInstanceCache()
		require.Equal(t, oldCodeSize, runtimeContext.GetSCCodeSize())
		instance := runtimeContext.iTracker.instance
		runtimeContext.pushInstance()
		runtimeContext.iTracker.codeSize = newCodeSize
		require.Equal(t, newCodeSize, runtimeContext.GetSCCodeSize())
		runtimeContext.iTracker.instance = &wasmer.Instance{}
		require.Equal(t, 1, len(runtimeContext.iTracker.instanceStack))

		runtimeContext.popInstance()
		require.Equal(t, oldCodeSize, runtimeContext.GetSCCodeSize())
		require.NotNil(t, runtimeContext.iTracker.instance)
		require.Equal(t, instance, runtimeContext.iTracker.instance)
		require.Equal(t, 0, len(runtimeContext.iTracker.instanceStack))

		runtimeContext.pushInstance()
		require.Equal(t, 1, len(runtimeContext.iTracker.instanceStack))
	})
}

func TestRuntimeContext_PushPopState(t *testing.T) {
	host := &contextmock.VMHostMock{}
	host.SCAPIMethods = MakeAPIImports()
	host.EnableEpochsHandlerField = &mock.EnableEpochsHandlerStub{}
	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()

	runtimeContext.SetMaxInstanceCount(1)

	vmInput := vmcommon.VMInput{
		CallerAddr:    []byte("caller"),
		GasProvided:   1000,
		CallValue:     big.NewInt(0),
		ESDTTransfers: make([]*vmcommon.ESDTTransfer, 0),
	}

	funcName := "test_func"
	scAddress := []byte("smartcontract")
	input := &vmcommon.ContractCallInput{
		VMInput:       vmInput,
		RecipientAddr: scAddress,
		Function:      funcName,
	}
	runtimeContext.InitStateFromContractCallInput(input)
	runtimeContext.iTracker.instance = &wasmer.Instance{}
	runtimeContext.PushState()
	require.Equal(t, 1, len(runtimeContext.stateStack))

	// change state
	runtimeContext.SetCodeAddress([]byte("dummy"))
	runtimeContext.SetVMInput(nil)
	runtimeContext.SetReadOnly(true)

	require.Equal(t, []byte("dummy"), runtimeContext.codeAddress)
	require.Nil(t, runtimeContext.GetVMInput())
	require.True(t, runtimeContext.ReadOnly())

	runtimeContext.PopSetActiveState()

	// check state was restored correctly
	require.Equal(t, scAddress, runtimeContext.GetContextAddress())
	require.Equal(t, funcName, runtimeContext.Function())
	require.Equal(t, input, runtimeContext.GetVMInput())
	require.False(t, runtimeContext.ReadOnly())
	require.Nil(t, runtimeContext.Arguments())

	runtimeContext.iTracker.instance = &wasmer.Instance{}
	runtimeContext.PushState()
	require.Equal(t, 1, len(runtimeContext.stateStack))

	runtimeContext.iTracker.instance = &wasmer.Instance{}
	runtimeContext.PushState()
	require.Equal(t, 2, len(runtimeContext.stateStack))

	runtimeContext.PopDiscard()
	require.Equal(t, 1, len(runtimeContext.stateStack))

	runtimeContext.ClearStateStack()
	require.Equal(t, 0, len(runtimeContext.stateStack))
}

func TestRuntimeContext_Instance(t *testing.T) {
	host := InitializeVMAndWasmer()
	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()

	runtimeContext.SetMaxInstanceCount(1)

	gasLimit := uint64(100000000)
	path := counterWasmCode
	contractCode := vmhost.GetSCCode(path)
	err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	gasPoints := uint64(100)
	runtimeContext.SetPointsUsed(gasPoints)
	require.Equal(t, gasPoints, runtimeContext.GetPointsUsed())

	funcName := "increment"
	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallValue: big.NewInt(0),
		},
		RecipientAddr: []byte("addr"),
		Function:      funcName,
	}
	runtimeContext.InitStateFromContractCallInput(input)

	f, err := runtimeContext.GetFunctionToCall()
	require.Nil(t, err)
	require.NotNil(t, f)

	input.Function = "func"
	runtimeContext.InitStateFromContractCallInput(input)
	f, err = runtimeContext.GetFunctionToCall()
	require.Equal(t, vmhost.ErrFuncNotFound, err)
	require.Equal(t, "", f)

	initFunc := runtimeContext.GetInitFunction()
	require.NotNil(t, initFunc)

	runtimeContext.ClearWarmInstanceCache()
	require.Nil(t, runtimeContext.iTracker.instance)
}

func TestRuntimeContext_Breakpoints(t *testing.T) {
	host := InitializeVMAndWasmer()
	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()

	mockOutput := &contextmock.OutputContextMock{
		OutputAccountMock: NewVMOutputAccount([]byte("address")),
	}
	mockOutput.OutputAccountMock.Code = []byte("code")
	mockOutput.SetReturnMessage("")

	host.OutputContext = mockOutput

	runtimeContext.SetMaxInstanceCount(1)

	gasLimit := uint64(100000000)
	path := counterWasmCode
	contractCode := vmhost.GetSCCode(path)
	err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	// Set and get curent breakpoint value
	require.Equal(t, vmhost.BreakpointNone, runtimeContext.GetRuntimeBreakpointValue())
	runtimeContext.SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
	require.Equal(t, vmhost.BreakpointOutOfGas, runtimeContext.GetRuntimeBreakpointValue())

	runtimeContext.SetRuntimeBreakpointValue(vmhost.BreakpointNone)
	require.Equal(t, vmhost.BreakpointNone, runtimeContext.GetRuntimeBreakpointValue())

	// Signal user error
	mockOutput.SetReturnCode(vmcommon.Ok)
	mockOutput.SetReturnMessage("")
	runtimeContext.SetRuntimeBreakpointValue(vmhost.BreakpointNone)

	runtimeContext.SignalUserError("something happened")
	require.Equal(t, vmhost.BreakpointSignalError, runtimeContext.GetRuntimeBreakpointValue())
	require.Equal(t, vmcommon.UserError, mockOutput.ReturnCode())
	require.Equal(t, "something happened", mockOutput.ReturnMessage())

	// Fail execution
	mockOutput.SetReturnCode(vmcommon.Ok)
	mockOutput.SetReturnMessage("")
	runtimeContext.SetRuntimeBreakpointValue(vmhost.BreakpointNone)

	runtimeContext.FailExecution(nil)
	require.Equal(t, vmhost.BreakpointExecutionFailed, runtimeContext.GetRuntimeBreakpointValue())
	require.Equal(t, vmcommon.ExecutionFailed, mockOutput.ReturnCode())
	require.Equal(t, "execution failed", mockOutput.ReturnMessage())

	mockOutput.SetReturnCode(vmcommon.Ok)
	mockOutput.SetReturnMessage("")
	runtimeContext.SetRuntimeBreakpointValue(vmhost.BreakpointNone)
	require.Equal(t, vmhost.BreakpointNone, runtimeContext.GetRuntimeBreakpointValue())

	runtimeError := errors.New("runtime error")
	runtimeContext.FailExecution(runtimeError)
	require.Equal(t, vmhost.BreakpointExecutionFailed, runtimeContext.GetRuntimeBreakpointValue())
	require.Equal(t, vmcommon.ExecutionFailed, mockOutput.ReturnCode())
	require.Equal(t, runtimeError.Error(), mockOutput.ReturnMessage())
}

func TestRuntimeContext_MemLoadStoreOk(t *testing.T) {
	host := InitializeVMAndWasmer()
	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()

	runtimeContext.SetMaxInstanceCount(1)

	gasLimit := uint64(100000000)
	path := counterWasmCode
	contractCode := vmhost.GetSCCode(path)
	err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	memory := runtimeContext.iTracker.instance.GetMemory()

	memContents, err := runtimeContext.MemLoad(10, 10)
	require.Nil(t, err)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, memContents)

	pageSize := uint32(65536)
	require.Equal(t, 2*pageSize, memory.Length())

	memContents = []byte("test data")
	err = runtimeContext.MemStore(10, memContents)
	require.Nil(t, err)
	require.Equal(t, 2*pageSize, memory.Length())

	memContents, err = runtimeContext.MemLoad(10, 10)
	require.Nil(t, err)
	require.Equal(t, []byte{'t', 'e', 's', 't', ' ', 'd', 'a', 't', 'a', 0}, memContents)
}

func TestRuntimeContext_MemoryIsBlank(t *testing.T) {
	host := InitializeVMAndWasmer()
	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()

	runtimeContext.SetMaxInstanceCount(1)

	gasLimit := uint64(100000000)
	path := "./../../test/contracts/answer/output/answer.wasm"
	contractCode := vmhost.GetSCCode(path)
	err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	memory := runtimeContext.iTracker.instance.GetMemory()
	totalPages := 2
	memoryContents := memory.Data()
	require.Equal(t, memory.Length(), uint32(len(memoryContents)))
	require.Equal(t, totalPages*vmhost.WASMPageSize, len(memoryContents))

	for i, value := range memoryContents {
		if value != byte(0) {
			msg := fmt.Sprintf("Non-zero value found at %d in Wasmer memory: 0x%X", i, value)
			require.Fail(t, msg)
		}
	}
}

func TestRuntimeContext_MemLoadCases(t *testing.T) {
	host := InitializeVMAndWasmer()
	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()

	runtimeContext.SetMaxInstanceCount(1)

	gasLimit := uint64(100000000)
	path := counterWasmCode
	contractCode := vmhost.GetSCCode(path)
	err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	memory := runtimeContext.iTracker.instance.GetMemory()

	var offset int32
	var length int32
	// Offset too small
	offset = -3
	length = 10
	memContents, err := runtimeContext.MemLoad(offset, length)
	require.True(t, errors.Is(err, vmhost.ErrBadBounds))
	require.Nil(t, memContents)

	// Offset too larget
	offset = int32(memory.Length() + 1)
	length = 10
	memContents, err = runtimeContext.MemLoad(offset, length)
	require.True(t, errors.Is(err, vmhost.ErrBadBounds))
	require.Nil(t, memContents)

	// Negative length
	offset = 10
	length = -2
	memContents, err = runtimeContext.MemLoad(offset, length)
	require.True(t, errors.Is(err, vmhost.ErrNegativeLength))
	require.Nil(t, memContents)

	// Requested end too large
	memContents = []byte("test data")
	offset = int32(memory.Length() - 9)
	err = runtimeContext.MemStore(offset, memContents)
	require.Nil(t, err)

	offset = int32(memory.Length() - 9)
	length = 9
	memContents, err = runtimeContext.MemLoad(offset, length)
	require.Nil(t, err)
	require.Equal(t, []byte("test data"), memContents)

	offset = int32(memory.Length() - 8)
	length = 9
	memContents, err = runtimeContext.MemLoad(offset, length)
	require.Nil(t, err)
	require.Equal(t, []byte{'e', 's', 't', ' ', 'd', 'a', 't', 'a', 0}, memContents)

	// Zero length
	offset = int32(memory.Length() - 8)
	length = 0
	memContents, err = runtimeContext.MemLoad(offset, length)
	require.Nil(t, err)
	require.Equal(t, []byte{}, memContents)
}

func TestRuntimeContext_MemStoreCases(t *testing.T) {
	host := InitializeVMAndWasmer()
	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()

	runtimeContext.SetMaxInstanceCount(1)

	gasLimit := uint64(100000000)
	path := counterWasmCode
	contractCode := vmhost.GetSCCode(path)
	err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	memory := runtimeContext.iTracker.instance.GetMemory()
	require.Equal(t, 2*vmhost.WASMPageSize, int(memory.Length()))

	// Bad lower bounds
	memContents := []byte("test data")
	offset := int32(-2)
	err = runtimeContext.MemStore(offset, memContents)
	require.True(t, errors.Is(err, vmhost.ErrBadLowerBounds))

	// Write something, then overwrite, then overwrite with empty byte slice
	memContents = []byte("this is a message")
	offset = int32(memory.Length() - 100)
	err = runtimeContext.MemStore(offset, memContents)
	require.Nil(t, err)

	memContents, err = runtimeContext.MemLoad(offset, 17)
	require.Nil(t, err)
	require.Equal(t, []byte("this is a message"), memContents)

	memContents = []byte("this is something")
	err = runtimeContext.MemStore(offset, memContents)
	require.Nil(t, err)

	memContents, err = runtimeContext.MemLoad(offset, 17)
	require.Nil(t, err)
	require.Equal(t, []byte("this is something"), memContents)

	memContents = []byte{}
	err = runtimeContext.MemStore(offset, memContents)
	require.Nil(t, err)

	memContents, err = runtimeContext.MemLoad(offset, 17)
	require.Nil(t, err)
	require.Equal(t, []byte("this is something"), memContents)
}

func TestRuntimeContext_MemStoreForbiddenGrowth(t *testing.T) {
	host := InitializeVMAndWasmer()
	host.EpochsStub().IsRuntimeMemStoreLimitEnabledField = true

	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()

	runtimeContext.SetMaxInstanceCount(1)

	gasLimit := uint64(100000000)
	path := counterWasmCode
	contractCode := vmhost.GetSCCode(path)
	err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	memory := runtimeContext.iTracker.instance.GetMemory()
	require.Equal(t, 2*vmhost.WASMPageSize, int(memory.Length()))

	memContents := []byte("test data")

	// Memory growth via MemStore forbidden
	offset := int32(memory.Length() - 4)
	err = runtimeContext.MemStore(offset, memContents)
	require.True(t, errors.Is(err, vmhost.ErrBadUpperBounds))
	require.Equal(t, 2*vmhost.WASMPageSize, int(memory.Length()))

	// Memory growth via MemStore forbidden
	memContents = make([]byte, vmhost.WASMPageSize+100)
	offset = int32(memory.Length() - 50)
	err = runtimeContext.MemStore(offset, memContents)
	require.True(t, errors.Is(err, vmhost.ErrBadUpperBounds))
	require.Equal(t, 2*vmhost.WASMPageSize, int(memory.Length()))

}

func TestRuntimeContext_MemLoadStoreVsInstanceStack(t *testing.T) {
	host := InitializeVMAndWasmer()
	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()

	runtimeContext.SetMaxInstanceCount(2)

	gasLimit := uint64(100000000)
	path := counterWasmCode
	contractCode := vmhost.GetSCCode(path)
	err := runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	// Write "test data1" to the WASM memory of the current instance
	memContents := []byte("test data1")
	err = runtimeContext.MemStore(10, memContents)
	require.Nil(t, err)

	memContents, err = runtimeContext.MemLoad(10, 10)
	require.Nil(t, err)
	require.Equal(t, []byte("test data1"), memContents)

	// Push the current instance down the instance stack
	runtimeContext.pushInstance()
	require.Equal(t, 1, len(runtimeContext.iTracker.instanceStack))

	// Create a new Wasmer instance
	contractCode = vmhost.GetSCCode(path)
	err = runtimeContext.StartWasmerInstance(contractCode, gasLimit, false)
	require.Nil(t, err)

	// Write "test data2" to the WASM memory of the new instance
	memContents = []byte("test data2")
	err = runtimeContext.MemStore(10, memContents)
	require.Nil(t, err)

	memContents, err = runtimeContext.MemLoad(10, 10)
	require.Nil(t, err)
	require.Equal(t, []byte("test data2"), memContents)

	// Pop the initial instance from the stack, making it the 'current instance'
	runtimeContext.popInstance()
	require.Equal(t, 0, len(runtimeContext.iTracker.instanceStack))

	// Check whether the previously-written string "test data1" is still in the
	// memory of the initial instance
	memContents, err = runtimeContext.MemLoad(10, 10)
	require.Nil(t, err)
	require.Equal(t, []byte("test data1"), memContents)

	// Write "test data3" to the WASM memory of the initial instance (now current)
	memContents = []byte("test data3")
	err = runtimeContext.MemStore(10, memContents)
	require.Nil(t, err)

	memContents, err = runtimeContext.MemLoad(10, 10)
	require.Nil(t, err)
	require.Equal(t, []byte("test data3"), memContents)
}

func TestRuntimeContext_PopSetActiveStateIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()

	host := InitializeVMAndWasmer()
	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()

	runtimeContext.PopSetActiveState()

	require.Equal(t, 0, len(runtimeContext.stateStack))
}

func TestRuntimeContext_PopDiscardIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()

	host := InitializeVMAndWasmer()
	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()

	runtimeContext.PopDiscard()

	require.Equal(t, 0, len(runtimeContext.stateStack))
}

func TestRuntimeContext_PopInstanceIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()

	host := InitializeVMAndWasmer()

	runtimeContext := makeDefaultRuntimeContext(t, host)
	defer runtimeContext.ClearWarmInstanceCache()
	runtimeContext.popInstance()

	require.Equal(t, 0, len(runtimeContext.stateStack))
}

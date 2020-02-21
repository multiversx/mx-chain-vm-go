package contexts

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/elrondapi"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/ethapi"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

func InitializeWasmer() {
	imports, _ := elrondapi.ElrondEImports()
	imports, _ = elrondapi.BigIntImports(imports)
	imports, _ = ethapi.EthereumImports(imports)
	imports, _ = crypto.CryptoImports(imports)
	_ = wasmer.SetImports(imports)

	gasSchedule := config.MakeGasMap(1)
	gasCostConfig, _ := config.CreateGasConfig(gasSchedule)
	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)
}

func getSCCode(fileName string) []byte {
	code, _ := ioutil.ReadFile(filepath.Clean(fileName))

	return code
}

func TestNewRuntimeContext(t *testing.T) {
	InitializeWasmer()

	host := &mock.VmHostMock{}
	vmType := []byte("type")

	runtimeContext, err := NewRuntimeContext(host, vmType)
	require.Nil(t, err)
	require.NotNil(t, runtimeContext)

	require.Equal(t, &vmcommon.VMInput{}, runtimeContext.vmInput)
	require.Equal(t, []byte{}, runtimeContext.scAddress)
	require.Equal(t, "", runtimeContext.callFunction)
	require.Equal(t, false, runtimeContext.readOnly)
	require.NotNil(t, runtimeContext.argParser)
	require.Nil(t, runtimeContext.asyncCallInfo)
}

func TestRuntimeContext_InitState(t *testing.T) {
	InitializeWasmer()

	host := &mock.VmHostMock{}
	vmType := []byte("type")

	runtimeContext, err := NewRuntimeContext(host, vmType)
	require.Nil(t, err)
	require.NotNil(t, runtimeContext)

	runtimeContext.vmInput = nil
	runtimeContext.scAddress = []byte("some address")
	runtimeContext.callFunction = "a function"
	runtimeContext.readOnly = true
	runtimeContext.argParser = nil
	runtimeContext.asyncCallInfo = &arwen.AsyncCallInfo{}

	runtimeContext.InitState()

	require.Equal(t, &vmcommon.VMInput{}, runtimeContext.vmInput)
	require.Equal(t, []byte{}, runtimeContext.scAddress)
	require.Equal(t, "", runtimeContext.callFunction)
	require.Equal(t, false, runtimeContext.readOnly)
	require.NotNil(t, runtimeContext.argParser)
	require.Nil(t, runtimeContext.asyncCallInfo)
}

func TestRuntimeContext_NewWasmerInstance(t *testing.T) {
	InitializeWasmer()

	host := &mock.VmHostMock{}
	vmType := []byte("type")

	runtimeContext, err := NewRuntimeContext(host, vmType)
	gasLimit := uint64(100000000)
	dummy := []byte("contract")
	err = runtimeContext.CreateWasmerInstance(dummy, gasLimit)
	require.NotNil(t, err)

	path := "./../../test/contracts/counter.wasm"
	contractCode := getSCCode(path)
	err = runtimeContext.CreateWasmerInstance(contractCode, gasLimit)
	require.Nil(t, err)
	require.Equal(t, arwen.BreakpointNone, runtimeContext.GetRuntimeBreakpointValue())
}

func TestRuntimeContext_StateSettersAndGetters(t *testing.T) {
	host := &mock.VmHostMock{}
	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType)

	arguments := [][]byte{[]byte("argument 1"), []byte("argument 2")}
	vmInput := vmcommon.VMInput{
		CallerAddr: []byte("caller"),
		Arguments:  arguments,
	}
	callInput := &vmcommon.ContractCallInput{
		VMInput:       vmInput,
		RecipientAddr: []byte("recipient"),
		Function:      "test function",
	}

	runtimeContext.InitStateFromContractCallInput(callInput)
	require.Equal(t, []byte("caller"), runtimeContext.GetVMInput().CallerAddr)
	require.Equal(t, []byte("recipient"), runtimeContext.GetSCAddress())
	require.Equal(t, "test function", runtimeContext.Function())
	require.Equal(t, vmType, runtimeContext.GetVMType())
	require.NotNil(t, runtimeContext.ArgParser())
	require.Equal(t, arguments, runtimeContext.Arguments())

	vmInput2 := vmcommon.VMInput{
		CallerAddr: []byte("caller2"),
		Arguments:  arguments,
	}
	runtimeContext.SetVMInput(&vmInput2)
	require.Equal(t, []byte("caller2"), runtimeContext.GetVMInput().CallerAddr)

	runtimeContext.SetSCAddress([]byte("smartcontract"))
	require.Equal(t, []byte("smartcontract"), runtimeContext.GetSCAddress())
}

func TestRuntimeContext_PushPopInstance(t *testing.T) {
	InitializeWasmer()
	host := &mock.VmHostMock{}
	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType)

	gasLimit := uint64(100000000)
	path := "./../../test/contracts/counter.wasm"
	contractCode := getSCCode(path)
	err := runtimeContext.CreateWasmerInstance(contractCode, gasLimit)
	require.Nil(t, err)

	instance := runtimeContext.instance

	runtimeContext.PushInstance()
	runtimeContext.instance = nil
	require.Equal(t, 1, len(runtimeContext.instanceStack))

	runtimeContext.PopInstance()
	require.NotNil(t, runtimeContext.instance)
	require.Equal(t, instance, runtimeContext.instance)
	require.Equal(t, 0, len(runtimeContext.instanceStack))

	runtimeContext.PushInstance()
	require.Equal(t, 1, len(runtimeContext.instanceStack))
	runtimeContext.ClearInstanceStack()
	require.Equal(t, 0, len(runtimeContext.instanceStack))
}

func TestRuntimeContext_PushPopState(t *testing.T) {
	host := &mock.VmHostMock{}
	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType)

	vmInput := vmcommon.VMInput{
		CallerAddr:  []byte("caller"),
		GasProvided: 1000,
	}

	funcName := "test_func"
	scAddress := []byte("smartcontract")
	input := &vmcommon.ContractCallInput{
		VMInput:       vmInput,
		RecipientAddr: scAddress,
		Function:      funcName,
	}
	runtimeContext.InitStateFromContractCallInput(input)

	runtimeContext.PushState()
	require.Equal(t, 1, len(runtimeContext.stateStack))

	// change state
	runtimeContext.SetSCAddress([]byte("dummy"))
	runtimeContext.SetVMInput(nil)
	runtimeContext.SetReadOnly(true)

	runtimeContext.PopState()

	//check state was restored correctly
	require.Equal(t, scAddress, runtimeContext.GetSCAddress())
	require.Equal(t, funcName, runtimeContext.Function())
	require.Equal(t, &vmInput, runtimeContext.GetVMInput())
	require.False(t, runtimeContext.ReadOnly())
	require.Nil(t, runtimeContext.Arguments())

	runtimeContext.PushState()
	require.Equal(t, 1, len(runtimeContext.stateStack))
	runtimeContext.ClearStateStack()
	require.Equal(t, 0, len(runtimeContext.stateStack))
}

func TestRuntimeContext_Instance(t *testing.T) {
	InitializeWasmer()
	host := &mock.VmHostMock{}
	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType)

	gasLimit := uint64(100000000)
	path := "./../../test/contracts/counter.wasm"
	contractCode := getSCCode(path)
	err := runtimeContext.CreateWasmerInstance(contractCode, gasLimit)
	require.Nil(t, err)

	gasPoints := uint64(100)
	runtimeContext.SetPointsUsed(gasPoints)
	require.Equal(t, gasPoints, runtimeContext.GetPointsUsed())

	funcName := "increment"
	input := &vmcommon.ContractCallInput{
		VMInput:       vmcommon.VMInput{},
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
	require.Equal(t, arwen.ErrFuncNotFound, err)
	require.Nil(t, f)

	initFunc := runtimeContext.GetInitFunction()
	require.Nil(t, initFunc)

	runtimeContext.CleanInstance()
	require.Nil(t, runtimeContext.instance)
}

func TestRuntimeContext_Breakpoints(t *testing.T) {
	mockOutput := &mock.OutputContextMock{}
	mockOutput.SetReturnMessage("")
	host := &mock.VmHostMock{
		OutputContext: mockOutput,
	}
	InitializeWasmer()
	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType)
	gasLimit := uint64(100000000)
	path := "./../../test/contracts/counter.wasm"
	contractCode := getSCCode(path)
	err := runtimeContext.CreateWasmerInstance(contractCode, gasLimit)
	require.Nil(t, err)

	// Set and get curent breakpoint value
	require.Equal(t, arwen.BreakpointNone, runtimeContext.GetRuntimeBreakpointValue())
	runtimeContext.SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
	require.Equal(t, arwen.BreakpointOutOfGas, runtimeContext.GetRuntimeBreakpointValue())

	runtimeContext.SetRuntimeBreakpointValue(arwen.BreakpointNone)
	require.Equal(t, arwen.BreakpointNone, runtimeContext.GetRuntimeBreakpointValue())

	// Signal user error
	mockOutput.SetReturnCode(vmcommon.Ok)
	mockOutput.SetReturnMessage("")
	runtimeContext.SetRuntimeBreakpointValue(arwen.BreakpointNone)

	runtimeContext.SignalUserError("something happened")
	require.Equal(t, arwen.BreakpointSignalError, runtimeContext.GetRuntimeBreakpointValue())
	require.Equal(t, vmcommon.UserError, mockOutput.ReturnCode())
	require.Equal(t, "something happened", mockOutput.ReturnMessage())

	// Fail execution
	mockOutput.SetReturnCode(vmcommon.Ok)
	mockOutput.SetReturnMessage("")
	runtimeContext.SetRuntimeBreakpointValue(arwen.BreakpointNone)

	runtimeContext.FailExecution(nil)
	require.Equal(t, arwen.BreakpointExecutionFailed, runtimeContext.GetRuntimeBreakpointValue())
	require.Equal(t, vmcommon.ExecutionFailed, mockOutput.ReturnCode())
	require.Equal(t, "execution failed", mockOutput.ReturnMessage())

	mockOutput.SetReturnCode(vmcommon.Ok)
	mockOutput.SetReturnMessage("")
	runtimeContext.SetRuntimeBreakpointValue(arwen.BreakpointNone)
	require.Equal(t, arwen.BreakpointNone, runtimeContext.GetRuntimeBreakpointValue())

	runtimeError := errors.New("runtime error")
	runtimeContext.FailExecution(runtimeError)
	require.Equal(t, arwen.BreakpointExecutionFailed, runtimeContext.GetRuntimeBreakpointValue())
	require.Equal(t, vmcommon.ExecutionFailed, mockOutput.ReturnCode())
	require.Equal(t, runtimeError.Error(), mockOutput.ReturnMessage())
}

func TestRuntimeContext_MemLoadStoreOk(t *testing.T) {
	InitializeWasmer()
	host := &mock.VmHostMock{}
	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType)

	gasLimit := uint64(100000000)
	path := "./../../test/contracts/counter.wasm"
	contractCode := getSCCode(path)
	err := runtimeContext.CreateWasmerInstance(contractCode, gasLimit)
	require.Nil(t, err)

	memory := runtimeContext.instance.Memory
	runtimeContext.instanceContext = wasmer.NewInstanceContext(nil, *memory)

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

func TestRuntimeContext_MemLoadCases(t *testing.T) {
	InitializeWasmer()
	host := &mock.VmHostMock{}
	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType)

	gasLimit := uint64(100000000)
	path := "./../../test/contracts/counter.wasm"
	contractCode := getSCCode(path)
	err := runtimeContext.CreateWasmerInstance(contractCode, gasLimit)
	require.Nil(t, err)

	memory := runtimeContext.instance.Memory
	runtimeContext.instanceContext = wasmer.NewInstanceContext(nil, *memory)

	var offset int32
	var length int32
	// Offset too small
	offset = -3
	length = 10
	memContents, err := runtimeContext.MemLoad(offset, length)
	require.True(t, errors.Is(err, arwen.ErrBadBounds))
	require.Nil(t, memContents)

	// Offset too larget
	offset = int32(memory.Length() + 1)
	length = 10
	memContents, err = runtimeContext.MemLoad(offset, length)
	require.True(t, errors.Is(err, arwen.ErrBadBounds))
	require.Nil(t, memContents)

	// Negative length
	offset = 10
	length = -2
	memContents, err = runtimeContext.MemLoad(offset, length)
	require.True(t, errors.Is(err, arwen.ErrNegativeLength))
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
}

func TestRuntimeContext_MemStoreCases(t *testing.T) {
	InitializeWasmer()
	host := &mock.VmHostMock{}
	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, vmType)

	gasLimit := uint64(100000000)
	path := "./../../test/contracts/counter.wasm"
	contractCode := getSCCode(path)
	err := runtimeContext.CreateWasmerInstance(contractCode, gasLimit)
	require.Nil(t, err)

	pageSize := uint32(65536)
	memory := runtimeContext.instance.Memory
	require.Equal(t, 2*pageSize, memory.Length())
	runtimeContext.instanceContext = wasmer.NewInstanceContext(nil, *memory)

	// Bad lower bounds
	memContents := []byte("test data")
	offset := int32(-2)
	err = runtimeContext.MemStore(offset, memContents)
	require.True(t, errors.Is(err, arwen.ErrBadLowerBounds))

	// Memory growth
	require.Equal(t, 2*pageSize, memory.Length())
	offset = int32(memory.Length() - 4)
	err = runtimeContext.MemStore(offset, memContents)
	require.Nil(t, err)
	require.Equal(t, 3*pageSize, memory.Length())

	// Bad upper bounds - forcing the Wasmer memory to grow more than a page at a
	// time is not allowed
	memContents = make([]byte, pageSize+100)
	offset = int32(memory.Length() - 50)
	err = runtimeContext.MemStore(offset, memContents)
	require.True(t, errors.Is(err, arwen.ErrBadUpperBounds))
	require.Equal(t, 4*pageSize, memory.Length())
}

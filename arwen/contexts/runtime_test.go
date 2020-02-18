package contexts

import (
	"fmt"
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

func TestRuntimeContext_InstanceMemory(t *testing.T) {
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
	fmt.Printf("memory size: %d\n", memory.Length())
	fmt.Println(memory.Data())
}

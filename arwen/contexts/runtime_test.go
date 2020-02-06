package contexts

import (
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
	t.Parallel()

	InitializeWasmer()

	host := &mock.VmHostStub{}
	vmType := []byte("type")

	runtimeContext, err := NewRuntimeContext(host, nil, vmType)
	require.Nil(t, err)
	require.NotNil(t, runtimeContext)

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

func TestRuntimeContext_InitStateFromContractCallInput(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostStub{}
	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, nil, vmType)

	vmInput := &vmcommon.VMInput{
		CallerAddr:  []byte("caller"),
		CallType:    0,
		GasPrice:    0,
		GasProvided: 1000,
	}
	funcName := "test_func"
	scAddress := []byte("addr")
	input := &vmcommon.ContractCallInput{
		VMInput:       vmcommon.VMInput{},
		RecipientAddr: nil,
		Function:      funcName,
	}

	runtimeContext.InitStateFromContractCallInput(input)
	runtimeContext.SetVMInput(vmInput)
	result := runtimeContext.GetVMInput()
	require.Equal(t, vmInput, result)

	runtimeContext.SetSCAddress(scAddress)
	resAddr := runtimeContext.GetSCAddress()
	require.Equal(t, scAddress, resAddr)

	require.Equal(t, funcName, runtimeContext.Function())
	require.Equal(t, vmType, runtimeContext.GetVMType())
}

func TestRuntimeContext_PushPopInstance(t *testing.T) {
	t.Parallel()

	InitializeWasmer()
	host := &mock.VmHostStub{}
	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, nil, vmType)

	gasLimit := uint64(100000000)
	path := "./../../test/contracts/counter.wasm"
	contractCode := getSCCode(path)
	err := runtimeContext.CreateWasmerInstance(contractCode, gasLimit)
	require.Nil(t, err)

	runtimeContext.PushInstance()
	require.Equal(t, 1, len(runtimeContext.instanceStack))

	err = runtimeContext.PopInstance()
	require.Nil(t, err)

	err = runtimeContext.PopInstance()
	require.Equal(t, arwen.InstanceStackUnderflow, err)
}

func TestRuntimeContext_PushPopState(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostStub{}
	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, nil, vmType)

	funcName := "test_func"
	scAddress := []byte("addr")
	vmInput := &vmcommon.VMInput{
		CallerAddr:  []byte("caller"),
		CallType:    0,
		GasPrice:    0,
		GasProvided: 1000,
	}
	input := &vmcommon.ContractCallInput{
		VMInput:       *vmInput,
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

	err := runtimeContext.PopState()
	require.Nil(t, err)

	//check state was restored correctly
	require.Equal(t, scAddress, runtimeContext.GetSCAddress())
	require.Equal(t, funcName, runtimeContext.Function())
	require.Equal(t, vmInput, runtimeContext.GetVMInput())
	require.False(t, runtimeContext.ReadOnly())
	require.Nil(t, runtimeContext.Arguments())

	err = runtimeContext.PopState()
	require.Equal(t, arwen.StateStackUnderflow, err)
}

func TestRuntimeContext_Instance(t *testing.T) {
	t.Parallel()

	InitializeWasmer()
	host := &mock.VmHostStub{}
	vmType := []byte("type")
	runtimeContext, _ := NewRuntimeContext(host, nil, vmType)

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

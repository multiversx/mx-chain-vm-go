package host

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

var defaultVmType = []byte{0xF, 0xF}
var counterKey = []byte{'m', 'y', 'c', 'o', 'u', 'n', 't', 'e', 'r', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func TestNewArwen(t *testing.T) {
	host, err := defaultArwen(t, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, host)
}

func TestExecution_DeployNewAddressErr(t *testing.T) {
	mockCryptoHook := &mock.CryptoHookMock{}
	stubBlockchainHook := &mock.BlockchainHookStub{}

	errNewAddress := errors.New("new address error")

	host, _ := defaultArwen(t, stubBlockchainHook, mockCryptoHook)
	input := defaultContractCreateInput()
	stubBlockchainHook.GetNonceCalled = func(address []byte) (uint64, error) {
		require.Equal(t, input.CallerAddr, address)
		return 0, nil
	}
	stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
		require.Equal(t, input.CallerAddr, creatorAddress)
		require.Equal(t, uint64(0), nonce)
		require.Equal(t, defaultVmType, vmType)
		return nil, errNewAddress
	}

	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
	require.Equal(t, errNewAddress.Error(), vmOutput.ReturnMessage)
}

func TestExecution_DeployOutOfGas(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := defaultArwenForDeployment(t, 24, newAddress)
	input := defaultContractCreateInput()
	input.GasProvided = 8 // default deployment requires 9 units of Gas
	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.OutOfGas, vmOutput.ReturnCode)
	require.Equal(t, arwen.ErrNotEnoughGas.Error(), vmOutput.ReturnMessage)
}

func TestExecution_DeployNotWASM(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := defaultArwenForDeployment(t, 24, newAddress)
	input := defaultContractCreateInput()
	input.GasProvided = 9
	input.ContractCode = []byte("not WASM")
	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractInvalid, vmOutput.ReturnCode)
}

func TestExecution_DeployWASM_WithoutMemory(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := defaultArwenForDeployment(t, 24, newAddress)
	input := defaultContractCreateInput()
	input.GasProvided = 1000
	input.ContractCode = arwen.GetTestSCCode("memoryless", "../../")
	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractInvalid, vmOutput.ReturnCode)
}

func TestExecution_DeployWASM_WrongInit(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := defaultArwenForDeployment(t, 24, newAddress)
	input := defaultContractCreateInput()
	input.GasProvided = 1000
	input.ContractCode = arwen.GetTestSCCode("init-wrong", "../../")
	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.FunctionWrongSignature, vmOutput.ReturnCode)
}

func TestExecution_DeployWASM_Successful(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := defaultArwenForDeployment(t, 24, newAddress)
	input := defaultContractCreateInput()
	input.CallValue = big.NewInt(88)
	input.GasProvided = 1000
	input.ContractCode = arwen.GetTestSCCode("init-correct", "../../")

	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Len(t, vmOutput.ReturnData, 1)
	require.Equal(t, []byte("init successful"), vmOutput.ReturnData[0])
	require.Equal(t, uint64(783), vmOutput.GasRemaining)
	require.Len(t, vmOutput.OutputAccounts, 2)
	require.Equal(t, uint64(24), vmOutput.OutputAccounts["caller"].Nonce)
	require.Equal(t, input.ContractCode, vmOutput.OutputAccounts["new smartcontract"].Code)
	require.Equal(t, big.NewInt(88), vmOutput.OutputAccounts["new smartcontract"].BalanceDelta)
}

func TestExecution_Deploy_DisallowFloatingPoint(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := defaultArwenForDeployment(t, 24, newAddress)
	input := defaultContractCreateInput()
	input.CallValue = big.NewInt(88)
	input.GasProvided = 1000
	input.ContractCode = arwen.GetTestSCCode("num-with-fp", "../../")

	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractInvalid, vmOutput.ReturnCode)
}

func TestExecution_CallGetCodeErr(t *testing.T) {
	mockCryptoHook := &mock.CryptoHookMock{}
	stubBlockchainHook := &mock.BlockchainHookStub{}

	errGetCode := errors.New("get code error")

	host, _ := defaultArwen(t, stubBlockchainHook, mockCryptoHook)
	input := defaultContractCallInput()
	stubBlockchainHook.GetCodeCalled = func(address []byte) ([]byte, error) {
		return nil, errGetCode
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractInvalid, vmOutput.ReturnCode)
	require.Equal(t, errGetCode.Error(), vmOutput.ReturnMessage)
}

func TestExecution_CallOutOfGas(t *testing.T) {
	code := arwen.GetTestSCCode("counter", "../../")
	host, _ := defaultArwenForCall(t, code)
	input := defaultContractCallInput()
	input.Function = "increment"

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.OutOfGas, vmOutput.ReturnCode)
	require.Equal(t, arwen.ErrNotEnoughGas.Error(), vmOutput.ReturnMessage)
}

func TestExecution_CallWasmerError(t *testing.T) {
	code := []byte("not WASM")
	host, _ := defaultArwenForCall(t, code)
	input := defaultContractCallInput()
	input.GasProvided = 100000
	input.Function = "increment"

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractInvalid, vmOutput.ReturnCode)
}

func TestExecution_CallSCMethod(t *testing.T) {
	code := arwen.GetTestSCCode("counter", "../../")
	host, _ := defaultArwenForCall(t, code)
	input := defaultContractCallInput()
	input.GasProvided = 100000

	// Calling init() is forbidden
	input.Function = "init"
	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.UserError, vmOutput.ReturnCode)
	require.Equal(t, arwen.ErrInitFuncCalledInRun.Error(), vmOutput.ReturnMessage)

	// Handle calling a missing function
	input.Function = "wrong"
	vmOutput, err = host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.FunctionNotFound, vmOutput.ReturnCode)
}

func TestExecution_Call_Successful(t *testing.T) {
	code := arwen.GetTestSCCode("counter", "../../")
	host, stubBlockchainHook := defaultArwenForCall(t, code)
	stubBlockchainHook.GetStorageDataCalled = func(scAddress []byte, key []byte) ([]byte, error) {
		return big.NewInt(1001).Bytes(), nil
	}
	input := defaultContractCallInput()
	input.GasProvided = 100000
	input.Function = "increment"

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Len(t, vmOutput.OutputAccounts, 1)
	require.Len(t, vmOutput.OutputAccounts["smartcontract"].StorageUpdates, 1)

	storedBytes := vmOutput.OutputAccounts["smartcontract"].StorageUpdates[string(counterKey)].Data
	require.Equal(t, big.NewInt(1002).Bytes(), storedBytes)
}

func defaultArwenForDeployment(t *testing.T, ownerNonce uint64, newAddress []byte) *vmHost {
	mockCryptoHook := &mock.CryptoHookMock{}
	stubBlockchainHook := &mock.BlockchainHookStub{}
	stubBlockchainHook.GetNonceCalled = func(address []byte) (uint64, error) {
		return 24, nil
	}
	stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
		return newAddress, nil
	}

	host, _ := defaultArwen(t, stubBlockchainHook, mockCryptoHook)
	return host
}

func defaultArwenForCall(t *testing.T, code []byte) (*vmHost, *mock.BlockchainHookStub) {
	mockCryptoHook := &mock.CryptoHookMock{}
	stubBlockchainHook := &mock.BlockchainHookStub{}
	stubBlockchainHook.GetCodeCalled = func(address []byte) ([]byte, error) {
		return code, nil
	}
	host, _ := defaultArwen(t, stubBlockchainHook, mockCryptoHook)
	return host, stubBlockchainHook
}

func defaultArwen(t *testing.T, blockchain vmcommon.BlockchainHook, crypto vmcommon.CryptoHook) (*vmHost, error) {
	host, err := NewArwenVM(blockchain, crypto, defaultVmType, uint64(1000), config.MakeGasMap(1))
	require.Nil(t, err)
	require.NotNil(t, host)
	return host, err
}

func defaultContractCreateInput() *vmcommon.ContractCreateInput {
	return &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr: []byte("caller"),
			Arguments: [][]byte{
				[]byte("argument 1"),
				[]byte("argument 2"),
			},
			CallValue:   big.NewInt(0),
			CallType:    vmcommon.DirectCall,
			GasPrice:    0,
			GasProvided: 0,
		},
		ContractCode: []byte("contract"),
	}
}

func defaultContractCallInput() *vmcommon.ContractCallInput {
	return &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  []byte("caller"),
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    vmcommon.DirectCall,
			GasPrice:    0,
			GasProvided: 0,
		},
		RecipientAddr: []byte("smartcontract"),
		Function:      "function",
	}
}

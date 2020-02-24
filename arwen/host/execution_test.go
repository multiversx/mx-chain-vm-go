package host

import (
	"bytes"
	"errors"
	"fmt"
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
var ErrCodeNotFound = errors.New("code not found")
var firstAddress = []byte("firstSC.........................")
var secondAddress = []byte("secondSC........................")

func TestNewArwen(t *testing.T) {
	t.Parallel()
	host, err := defaultArwen(t, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, host)
}

func TestExecution_DeployNewAddressErr(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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

func TestExecution_ExecuteOnSameContext(t *testing.T) {
	parentCode := arwen.GetTestSCCode("exec-same-ctx-parent", "../../")

	// Execute the parent SC method "parentFunctionPrepare", which sets storage,
	// finish data and performs a transfer. This step validates the test to the
	// actual call to ExecuteOnSameContext().
	host, _ := defaultArwenForCall(t, parentCode)
	input := defaultContractCallInput()
	input.CallerAddr = []byte("user")
	input.RecipientAddr = firstAddress
	input.Function = "parentFunctionPrepare"
	input.GasProvided = 1000000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	fmt.Println(vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	expectedVMOutput := expectedVMOutputs("ExecuteOnSameContext_Prepare")
	require.Equal(t, expectedVMOutput, vmOutput)

	// Call parentFunctionWrongCall() of the parent SC, which will try to call a
	// non-existing SC.
	host, _ = defaultArwenForCall(t, parentCode)
	input = defaultContractCallInput()
	input.CallerAddr = []byte("user")
	input.RecipientAddr = firstAddress
	input.Function = "parentFunctionWrongCall"
	input.GasProvided = 1000000

	vmOutput, err = host.RunSmartContractCall(input)
	require.Nil(t, err)
	expectedVMOutput = expectedVMOutputs("ExecuteOnSameContext_WrongCall")
	require.Equal(t, expectedVMOutput, vmOutput)

	// TODO verify whether the child can access bigInts of the parent?
	childCode := arwen.GetTestSCCode("exec-same-ctx-child", "../../")
	host, _ = defaultArwenForTwoSCs(t, parentCode, childCode)
	input = defaultContractCallInput()
	input.CallerAddr = []byte("user")
	input.RecipientAddr = firstAddress
	input.Function = "parentFunctionChildCall"
	input.GasProvided = 1000000

	vmOutput, err = host.RunSmartContractCall(input)
	require.Nil(t, err)
	expectedVMOutput = expectedVMOutputs("ExecuteOnSameContext_ChildCall")
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_Call_Breakpoints(t *testing.T) {
	t.Parallel()

	code := arwen.GetTestSCCode("breakpoint", "../../")
	host, _ := defaultArwenForCall(t, code)
	input := defaultContractCallInput()
	input.GasProvided = 100000
	input.Function = "testFunc"

	// Send the number 15 to the SC, causing it to finish with the number 100
	input.Arguments = [][]byte{[]byte{15}}
	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Equal(t, [][]byte{[]byte{100}}, vmOutput.ReturnData)

	// Send the number 1 to the SC, causing it to exit with ReturnMessage "exit
	// here" if the breakpoint mechanism works properly, or with the
	// ReturnMessage "exit later" if the breakpoint mechanism fails to stop the
	// execution.
	input.Arguments = [][]byte{[]byte{1}}
	vmOutput, err = host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.UserError, vmOutput.ReturnCode)
	require.Len(t, vmOutput.ReturnData, 0)
	require.Equal(t, "exit here", vmOutput.ReturnMessage)
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

func defaultArwenForTwoSCs(t *testing.T, firstCode []byte, secondCode []byte) (*vmHost, *mock.BlockchainHookStub) {
	mockCryptoHook := &mock.CryptoHookMock{}
	stubBlockchainHook := &mock.BlockchainHookStub{}
	stubBlockchainHook.GetCodeCalled = func(scAddress []byte) ([]byte, error) {
		if bytes.Equal(scAddress, firstAddress) {
			return firstCode, nil
		}
		if bytes.Equal(scAddress, secondAddress) {
			return secondCode, nil
		}
		return nil, ErrCodeNotFound
	}
	host, _ := defaultArwen(t, stubBlockchainHook, mockCryptoHook)
	return host, stubBlockchainHook
}

func expectedVMOutputs(id string) *vmcommon.VMOutput {
	parentKeyA := []byte("parentKeyA......................")
	parentKeyB := []byte("parentKeyB......................")
	childKey := []byte("childKey........................")
	parentDataA := []byte("parentDataA")
	parentDataB := []byte("parentDataB")
	childData := []byte("childData")
	parentFinishA := []byte("parentFinishA")
	parentFinishB := []byte("parentFinishB")
	childFinish := []byte("childFinish")
	parentTransferReceiver := []byte("parentTransferReceiver..........")
	childTransferReceiver := []byte("asdfoottxxwlllllllllllwrraattttt")
	parentTransferValue := int64(42)
	parentTransferData := []byte("parentTransferData")

	parentAddress := firstAddress
	childAddress := secondAddress
	wrongAddress := []byte("wrongSC.........................")

	if id == "ExecuteOnSameContext_Prepare" {
		expectedVMOutput := mock.MakeVMOutput()
		expectedVMOutput.ReturnCode = vmcommon.Ok
		expectedVMOutput.GasRemaining = 998255
		mock.AddFinishData(expectedVMOutput, parentFinishA)
		mock.AddFinishData(expectedVMOutput, parentFinishB)
		parentAccount := mock.AddNewOutputAccount(
			expectedVMOutput,
			parentAddress,
			-parentTransferValue,
			nil,
		)
		mock.SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
		mock.SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
		_ = mock.AddNewOutputAccount(
			expectedVMOutput,
			parentTransferReceiver,
			parentTransferValue,
			parentTransferData,
		)

		return expectedVMOutput
	}
	if id == "ExecuteOnSameContext_WrongCall" {
		expectedVMOutput := expectedVMOutputs("ExecuteOnSameContext_Prepare")
		mock.AddFinishData(expectedVMOutput, []byte("failed"))
		expectedVMOutput.GasRemaining = 988131
		parentAccount := expectedVMOutput.OutputAccounts[string(parentAddress)]
		parentAccount.BalanceDelta = big.NewInt(-141)
		_ = mock.AddNewOutputAccount(
			expectedVMOutput,
			wrongAddress,
			99, // TODO this is not supposed to happen! this should be 0.
			nil,
		)
		return expectedVMOutput
	}
	if id == "ExecuteOnSameContext_ChildCall" {
		expectedVMOutput := expectedVMOutputs("ExecuteOnSameContext_Prepare")
		mock.AddFinishData(expectedVMOutput, childFinish)
		mock.AddFinishData(expectedVMOutput, []byte("success"))
		expectedVMOutput.GasRemaining = 998206
		parentAccount := expectedVMOutput.OutputAccounts[string(parentAddress)]
		parentAccount.BalanceDelta = big.NewInt(-141)
		childAccount := mock.AddNewOutputAccount(
			expectedVMOutput,
			childAddress,
			3,
			nil,
		)
		mock.SetStorageUpdate(childAccount, childKey, childData)
		_ = mock.AddNewOutputAccount(
			expectedVMOutput,
			childTransferReceiver,
			96,
			[]byte("qwerty"),
		)

		return expectedVMOutput
	}
	if id == "Nil" {
		expectedVMOutput := mock.MakeVMOutput()
		expectedVMOutput.GasRemaining = 0
		expectedVMOutput.ReturnData = nil
		expectedVMOutput.OutputAccounts = nil
		expectedVMOutput.TouchedAccounts = nil
		expectedVMOutput.DeletedAccounts = nil
		expectedVMOutput.Logs = nil
		return expectedVMOutput
	}

	return nil
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

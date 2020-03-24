package host

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

var counterKey = []byte{'m', 'y', 'c', 'o', 'u', 'n', 't', 'e', 'r', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func TestNewArwen(t *testing.T) {
	t.Parallel()
	host, err := DefaultTestArwen(t, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, host)
}

func TestSCMem(t *testing.T) {
	code := GetTestSCCode("misc", "../../")
	host, _ := DefaultTestArwenForCall(t, code)
	input := DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = "iterate_over_byte_array"
	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	testString := "this is some random string of bytes"
	expectedData := [][]byte{
		[]byte(testString),
		[]byte{35},
	}
	for _, c := range testString {
		expectedData = append(expectedData, []byte{byte(c)})
	}
	require.Equal(t, expectedData, vmOutput.ReturnData)
}

func TestExecution_DeployNewAddressErr(t *testing.T) {
	t.Parallel()

	mockCryptoHook := &mock.CryptoHookMock{}
	stubBlockchainHook := &mock.BlockchainHookStub{}

	errNewAddress := errors.New("new address error")

	host, _ := DefaultTestArwen(t, stubBlockchainHook, mockCryptoHook)
	input := DefaultTestContractCreateInput()
	stubBlockchainHook.GetNonceCalled = func(address []byte) (uint64, error) {
		require.Equal(t, input.CallerAddr, address)
		return 0, nil
	}
	stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
		require.Equal(t, input.CallerAddr, creatorAddress)
		require.Equal(t, uint64(0), nonce)
		require.Equal(t, defaultVMType, vmType)
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
	host := DefaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
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
	host := DefaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
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
	host := DefaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.GasProvided = 1000
	input.ContractCode = GetTestSCCode("memoryless", "../../")
	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractInvalid, vmOutput.ReturnCode)
}

func TestExecution_DeployWASM_WrongInit(t *testing.T) {
	t.Parallel()

	newAddress := []byte("new smartcontract")
	host := DefaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.GasProvided = 1000
	input.ContractCode = GetTestSCCode("init-wrong", "../../")
	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.FunctionWrongSignature, vmOutput.ReturnCode)
}

func TestExecution_DeployWASM_Successful(t *testing.T) {
	t.Parallel()

	newAddress := []byte("new smartcontract")
	host := DefaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.CallValue = big.NewInt(88)
	input.GasProvided = 1000
	input.ContractCode = GetTestSCCode("init-correct", "../../")

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

func TestExecution_ManyDeployments(t *testing.T) {
	ownerNonce := uint64(23)
	newAddress := "new smartcontract"
	mockCryptoHook := &mock.CryptoHookMock{}
	stubBlockchainHook := &mock.BlockchainHookStub{}
	stubBlockchainHook.GetNonceCalled = func(address []byte) (uint64, error) {
		return ownerNonce, nil
	}
	stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
		ownerNonce++
		return []byte(newAddress + " " + string(ownerNonce)), nil
	}

	host, _ := DefaultTestArwen(t, stubBlockchainHook, mockCryptoHook)
	input := DefaultTestContractCreateInput()
	input.CallerAddr = []byte("owner")
	input.Arguments = make([][]byte, 0)
	input.CallValue = big.NewInt(88)
	input.ContractCode = GetTestSCCode("init-correct", "../../")

	numDeployments := 100000
	for i := 0; i < numDeployments; i++ {
		input.GasProvided = 100000
		vmOutput, err := host.RunSmartContractCreate(input)
		require.Nil(t, err)
		require.NotNil(t, vmOutput)
		if vmOutput.ReturnCode != vmcommon.Ok {
			fmt.Printf("Deployed %d SCs\n", i)
			fmt.Printf(vmOutput.ReturnMessage)
		}
		require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	}
}

func TestExecution_Deploy_DisallowFloatingPoint(t *testing.T) {
	t.Parallel()

	newAddress := []byte("new smartcontract")
	host := DefaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.CallValue = big.NewInt(88)
	input.GasProvided = 1000
	input.ContractCode = GetTestSCCode("num-with-fp", "../../")

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

	host, _ := DefaultTestArwen(t, stubBlockchainHook, mockCryptoHook)
	input := DefaultTestContractCallInput()
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

	code := GetTestSCCode("counter", "../../")
	host, _ := DefaultTestArwenForCall(t, code)
	input := DefaultTestContractCallInput()
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
	host, _ := DefaultTestArwenForCall(t, code)
	input := DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = "increment"

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractInvalid, vmOutput.ReturnCode)
}

func TestExecution_CallSCMethod(t *testing.T) {
	t.Parallel()

	code := GetTestSCCode("counter", "../../")
	host, _ := DefaultTestArwenForCall(t, code)
	input := DefaultTestContractCallInput()
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

	code := GetTestSCCode("counter", "../../")
	host, stubBlockchainHook := DefaultTestArwenForCall(t, code)
	stubBlockchainHook.GetStorageDataCalled = func(scAddress []byte, key []byte) ([]byte, error) {
		return big.NewInt(1001).Bytes(), nil
	}
	input := DefaultTestContractCallInput()
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

func TestExecution_ExecuteOnSameContext_Simple(t *testing.T) {
	parentCode := GetTestSCCode("exec-same-ctx-simple-parent", "../../")
	childCode := GetTestSCCode("exec-same-ctx-simple-child", "../../")
	host, _ := DefaultTestArwenForTwoSCs(t, parentCode, childCode)
	input := DefaultTestContractCallInput()
	input.CallerAddr = []byte("user")
	input.RecipientAddr = firstAddress
	input.Function = "parentFunctionChildCall"
	input.GasProvided = 1000000

	vmOutput, err := host.RunSmartContractCall(input)
	fmt.Println(vmOutput.ReturnMessage)
	require.Nil(t, err)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Equal(t, "", vmOutput.ReturnMessage)
}

func TestExecution_ExecuteOnSameContext(t *testing.T) {
	parentCode := GetTestSCCode("exec-same-ctx-parent", "../../")
	firstSC := []byte("firstSC.........................")
	firstSCBalance := big.NewInt(1000)

	getBalanceCalled := func(address []byte) (*big.Int, error) {
		if bytes.Equal(firstSC, address) {
			return firstSCBalance, nil
		}
		return big.NewInt(0), nil
	}

	// Execute the parent SC method "parentFunctionPrepare", which sets storage,
	// finish data and performs a transfer. This step validates the test to the
	// actual call to ExecuteOnSameContext().
	host, stubBlockchainHook := DefaultTestArwenForCall(t, parentCode)
	stubBlockchainHook.GetBalanceCalled = getBalanceCalled
	input := DefaultTestContractCallInput()
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
	host, stubBlockchainHook = DefaultTestArwenForCall(t, parentCode)
	stubBlockchainHook.GetBalanceCalled = getBalanceCalled
	host.Output().AddTxValueToAccount(firstSC, big.NewInt(1000))
	input = DefaultTestContractCallInput()
	input.CallerAddr = []byte("user")
	input.RecipientAddr = firstAddress
	input.Function = "parentFunctionWrongCall"
	input.GasProvided = 1000000

	vmOutput, err = host.RunSmartContractCall(input)
	require.Nil(t, err)
	expectedVMOutput = expectedVMOutputs("ExecuteOnSameContext_WrongCall")
	require.Equal(t, expectedVMOutput, vmOutput)

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnSameContext().
	childCode := GetTestSCCode("exec-same-ctx-child", "../../")
	host, stubBlockchainHook = DefaultTestArwenForTwoSCs(t, parentCode, childCode)
	stubBlockchainHook.GetBalanceCalled = getBalanceCalled
	host.Output().AddTxValueToAccount(firstSC, big.NewInt(1000))
	input = DefaultTestContractCallInput()
	input.CallerAddr = []byte("user")
	input.RecipientAddr = firstAddress
	input.Function = "parentFunctionChildCall"
	input.GasProvided = 1000000

	vmOutput, err = host.RunSmartContractCall(input)
	require.Nil(t, err)
	expectedVMOutput = expectedVMOutputs("ExecuteOnSameContext_ChildCall")
	require.Equal(t, expectedVMOutput, vmOutput)

	// Call parentFunctionChildCall_BigInts() of the parent SC, which will call a
	// method of the child SC that takes some big Int references as arguments and
	// produce a new big Int out of the arguments.
	childCode = GetTestSCCode("exec-same-ctx-child", "../../")
	host, stubBlockchainHook = DefaultTestArwenForTwoSCs(t, parentCode, childCode)
	stubBlockchainHook.GetBalanceCalled = getBalanceCalled
	host.Output().AddTxValueToAccount(firstSC, big.NewInt(1000))
	input = DefaultTestContractCallInput()
	input.CallerAddr = []byte("user")
	input.RecipientAddr = firstAddress
	input.Function = "parentFunctionChildCall_BigInts"
	input.GasProvided = 1000000

	vmOutput, err = host.RunSmartContractCall(input)
	require.Nil(t, err)
	expectedVMOutput = expectedVMOutputs("ExecuteOnSameContext_ChildCall_BigInts")
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_Call_Breakpoints(t *testing.T) {
	t.Parallel()

	code := GetTestSCCode("breakpoint", "../../")
	host, _ := DefaultTestArwenForCall(t, code)
	input := DefaultTestContractCallInput()
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

	baseGasRemaining := uint64(980098)

	if id == "ExecuteOnSameContext_Prepare" {
		expectedVMOutput := MakeVMOutput()
		expectedVMOutput.ReturnCode = vmcommon.Ok
		expectedVMOutput.GasRemaining = baseGasRemaining + 17658
		AddFinishData(expectedVMOutput, parentFinishA)
		AddFinishData(expectedVMOutput, parentFinishB)
		AddFinishData(expectedVMOutput, []byte("success"))
		parentAccount := AddNewOutputAccount(
			expectedVMOutput,
			parentAddress,
			-parentTransferValue,
			nil,
		)
		parentAccount.Balance = big.NewInt(1000)
		SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
		SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
		_ = AddNewOutputAccount(
			expectedVMOutput,
			parentTransferReceiver,
			parentTransferValue,
			parentTransferData,
		)

		return expectedVMOutput
	}
	if id == "ExecuteOnSameContext_WrongCall" {
		expectedVMOutput := expectedVMOutputs("ExecuteOnSameContext_Prepare")
		AddFinishData(expectedVMOutput, []byte("failed"))
		expectedVMOutput.GasRemaining = baseGasRemaining + 7532
		parentAccount := expectedVMOutput.OutputAccounts[string(parentAddress)]
		parentAccount.BalanceDelta = big.NewInt(-141)
		_ = AddNewOutputAccount(
			expectedVMOutput,
			wrongAddress,
			99, // TODO this is not supposed to happen! this should be 0.
			nil,
		)
		return expectedVMOutput
	}
	if id == "ExecuteOnSameContext_ChildCall" {
		expectedVMOutput := expectedVMOutputs("ExecuteOnSameContext_Prepare")
		AddFinishData(expectedVMOutput, childFinish)
		AddFinishData(expectedVMOutput, parentDataA)
		for _, c := range parentDataA {
			AddFinishData(expectedVMOutput, []byte{c})
		}
		AddFinishData(expectedVMOutput, parentDataB)
		for _, c := range parentDataB {
			AddFinishData(expectedVMOutput, []byte{c})
		}
		AddFinishData(expectedVMOutput, []byte("child ok"))
		AddFinishData(expectedVMOutput, []byte("success"))
		expectedVMOutput.GasRemaining = baseGasRemaining + 14075
		parentAccount := expectedVMOutput.OutputAccounts[string(parentAddress)]
		parentAccount.BalanceDelta = big.NewInt(-141)
		childAccount := AddNewOutputAccount(
			expectedVMOutput,
			childAddress,
			3,
			nil,
		)
		childAccount.Balance = big.NewInt(0)
		SetStorageUpdate(parentAccount, childKey, childData)
		_ = AddNewOutputAccount(
			expectedVMOutput,
			childTransferReceiver,
			96,
			[]byte("qwerty"),
		)

		return expectedVMOutput
	}
	if id == "ExecuteOnSameContext_ChildCall_BigInts" {
		expectedVMOutput := MakeVMOutput()
		expectedVMOutput.ReturnCode = vmcommon.Ok
		AddFinishData(expectedVMOutput, []byte("child ok"))
		AddFinishData(expectedVMOutput, []byte("success"))
		expectedVMOutput.GasRemaining = baseGasRemaining + 14547
		parentAccount := AddNewOutputAccount(
			expectedVMOutput,
			parentAddress,
			-parentTransferValue,
			nil,
		)
		parentAccount.Balance = big.NewInt(1000)
		parentAccount.BalanceDelta = big.NewInt(-99)
		_ = AddNewOutputAccount(
			expectedVMOutput,
			childAddress,
			99,
			nil,
		)
		return expectedVMOutput
	}
	if id == "Nil" {
		expectedVMOutput := MakeVMOutput()
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

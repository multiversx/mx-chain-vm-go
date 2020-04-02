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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var counterKey = []byte{'m', 'y', 'c', 'o', 'u', 'n', 't', 'e', 'r', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func TestNewArwen(t *testing.T) {
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
	newAddress := []byte("new smartcontract")
	host := DefaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.GasProvided = 1000
	input.ContractCode = GetTestSCCode("init-wrong", "../../")
	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractInvalid, vmOutput.ReturnCode)
}

func TestExecution_DeployWASM_WrongMethods(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := DefaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.GasProvided = 1000
	input.ContractCode = GetTestSCCode("signatures", "../../")
	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractInvalid, vmOutput.ReturnCode)
}

func TestExecution_DeployWASM_Successful(t *testing.T) {
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
	code := GetTestSCCode("counter", "../../")
	host, _ := DefaultTestArwenForCall(t, code)
	input := DefaultTestContractCallInput()
	input.GasProvided = 100000

	// Calling init() directly is forbidden
	input.Function = "init"
	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.UserError, vmOutput.ReturnCode)
	require.Equal(t, arwen.ErrInitFuncCalledInRun.Error(), vmOutput.ReturnMessage)

	// Calling callBack() directly is forbidden
	input.Function = "callBack"
	vmOutput, err = host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.UserError, vmOutput.ReturnCode)
	require.Equal(t, arwen.ErrCallBackFuncCalledInRun.Error(), vmOutput.ReturnMessage)

	// Handle calling a missing function
	input.Function = "wrong"
	vmOutput, err = host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.FunctionNotFound, vmOutput.ReturnCode)
}

func TestExecution_Call_Successful(t *testing.T) {
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
	require.Len(t, vmOutput.OutputAccounts[string(parentAddress)].StorageUpdates, 1)

	storedBytes := vmOutput.OutputAccounts[string(parentAddress)].StorageUpdates[string(counterKey)].Data
	require.Equal(t, big.NewInt(1002).Bytes(), storedBytes)
}

func TestExecution_ExecuteOnSameContext_Simple(t *testing.T) {
	parentCode := GetTestSCCode("exec-same-ctx-simple-parent", "../../")
	childCode := GetTestSCCode("exec-same-ctx-simple-child", "../../")

	host, _ := DefaultTestArwenForTwoSCs(t, parentCode, childCode)
	input := DefaultTestContractCallInput()
	input.CallerAddr = []byte("user")
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall"
	input.GasProvided = 1000000

	vmOutput, err := host.RunSmartContractCall(input)
	fmt.Println(vmOutput.ReturnMessage)
	require.Nil(t, err)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Equal(t, "", vmOutput.ReturnMessage)
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

func TestExecution_ExecuteOnSameContext_Prepare(t *testing.T) {
	parentCode := GetTestSCCode("exec-same-ctx-parent", "../../")
	parentSCBalance := big.NewInt(1000)

	getBalanceCalled := func(address []byte) (*big.Int, error) {
		if bytes.Equal(parentAddress, address) {
			return parentSCBalance, nil
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
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionPrepare"
	input.GasProvided = 1000000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	fmt.Println(vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	expectedVMOutput := expectedVMOutput_SameCtx_Prepare()
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnSameContext_Wrong(t *testing.T) {
	parentCode := GetTestSCCode("exec-same-ctx-parent", "../../")
	parentSCBalance := big.NewInt(1000)

	getBalanceCalled := func(address []byte) (*big.Int, error) {
		if bytes.Equal(parentAddress, address) {
			return parentSCBalance, nil
		}
		return big.NewInt(0), nil
	}

	// Call parentFunctionWrongCall() of the parent SC, which will try to call a
	// non-existing SC.
	host, stubBlockchainHook := DefaultTestArwenForCall(t, parentCode)
	stubBlockchainHook.GetBalanceCalled = getBalanceCalled
	input := DefaultTestContractCallInput()
	input.CallerAddr = []byte("user")
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionWrongCall"
	input.GasProvided = 1000000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	expectedVMOutput := expectedVMOutput_SameCtx_WrongContractCalled()
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnSameContext_OutOfGas(t *testing.T) {
	// Scenario:
	// Parent sets data into the storage, finishes data and creates a bigint
	// Parent calls executeOnSameContext, sending some value as well
	// Parent provides insufficient gas to executeOnSameContext (enoguh to start the SC though)
	// Child SC starts executing: sets data into the storage, finishes data and changes the bigint
	// Child starts an infinite loop, which must surely end with OutOfGas
	// Execution returns to parent, which finishes with the result of executeOnSameContext
	// Assertions: modifications made by the child are did not take effect
	// Assertions: the value sent by the parent to the child was returned to the parent
	// Assertions: the parent lost all the gas provided to executeOnSameContext
	parentCode := GetTestSCCode("exec-same-ctx-parent", "../../")
	childCode := GetTestSCCode("exec-same-ctx-child", "../../")
	parentSCBalance := big.NewInt(1000)

	getBalanceCalled := func(address []byte) (*big.Int, error) {
		if bytes.Equal(parentAddress, address) {
			return parentSCBalance, nil
		}

		return big.NewInt(0), nil
	}

	// Call parentFunctionChildCall_OutOfGas() of the parent SC, which will call
	// the child SC using executeOnSameContext() with sufficient gas for
	// compilation and starting, but the child starts an infinite loop which will
	// end in OutOfGas.
	host, stubBlockchainHook := DefaultTestArwenForTwoSCs(t, parentCode, childCode)
	stubBlockchainHook.GetBalanceCalled = getBalanceCalled
	input := DefaultTestContractCallInput()
	input.CallerAddr = []byte("user")
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall_OutOfGas"
	input.GasProvided = 1000000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	expectedVMOutput := expectedVMOutput_SameCtx_OutOfGas()
	assert.Equal(t, int64(42), host.BigInt().GetOne(0).Int64())
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnSameContext_Recursive_Direct(t *testing.T) {
	// Scenario:
	// SC has a method "callRecursive" which takes a byte as argument (number of recursive calls)
	// callRecursive() saves to storage "keyNNN" → "valueNNN", where NNN is the argument
	// callRecursive() saves to storage a counter starting at 1, increased by every recursive call
	// callRecursive() creates a bigInt and increments it with every iteration
	// callRecursive() calls itself using executeOnSameContext(), with the argument decremented
	// callRecursive() handles argument == 0 as follows: saves to storage the
	//		value of the bigInt counter, then exits without recursive call
	// Assertions: the VMOutput must contain as many StorageUpdates as the argument requires
	// Assertions: the VMOutput must contain as many finished values as the argument requires
	// Assertions: there must be a StorageUpdate with the value of the bigInt counter
}

func TestExecution_ExecuteOnSameContext_Recursive_Mutual_Methods(t *testing.T) {
}

func TestExecution_ExecuteOnSameContext_Recursive_Mutual_SCs(t *testing.T) {
	// Scenario:
	// Parent has method parentCallChild()
	// Child has method childCallParent()
	// The two methods are identical, just named differently
	// The methods do the following:
	//		parent: save to storage "keyParentNNN" → "valueParentNNN"
	//		parent:	finish "ParentNNN"
	//		child:	save to storage "keyChildNNN" → "valueChildNNN"
	//		child:	finish "ChildNNN"
	//		both:		increment a shared bigInt counter
	//		both:		whoever exits must save the shared bigInt counter to storage
}

func TestExecution_ExecuteOnSameContext_Successful(t *testing.T) {
	parentCode := GetTestSCCode("exec-same-ctx-parent", "../../")
	childCode := GetTestSCCode("exec-same-ctx-child", "../../")
	parentSCBalance := big.NewInt(1000)

	getBalanceCalled := func(address []byte) (*big.Int, error) {
		if bytes.Equal(parentAddress, address) {
			return parentSCBalance, nil
		}

		return big.NewInt(0), nil
	}

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnSameContext().
	host, stubBlockchainHook := DefaultTestArwenForTwoSCs(t, parentCode, childCode)
	stubBlockchainHook.GetBalanceCalled = getBalanceCalled
	input := DefaultTestContractCallInput()
	input.CallerAddr = []byte("user")
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall"
	input.GasProvided = 1000000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	expectedVMOutput := expectedVMOutput_SameCtx_SuccessfulChildCall()
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnSameContext_Successful_BigInts(t *testing.T) {
	parentCode := GetTestSCCode("exec-same-ctx-parent", "../../")
	childCode := GetTestSCCode("exec-same-ctx-child", "../../")
	parentSCBalance := big.NewInt(1000)

	getBalanceCalled := func(address []byte) (*big.Int, error) {
		if bytes.Equal(parentAddress, address) {
			return parentSCBalance, nil
		}
		return big.NewInt(0), nil
	}

	// Call parentFunctionChildCall_BigInts() of the parent SC, which will call a
	// method of the child SC that takes some big Int references as arguments and
	// produce a new big Int out of the arguments.
	host, stubBlockchainHook := DefaultTestArwenForTwoSCs(t, parentCode, childCode)
	stubBlockchainHook.GetBalanceCalled = getBalanceCalled
	input := DefaultTestContractCallInput()
	input.CallerAddr = []byte("user")
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall_BigInts"
	input.GasProvided = 1000000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	expectedVMOutput := expectedVMOutput_SameCtx_SuccessfulChildCall_BigInts()
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnDestContext_Prepare(t *testing.T) {
	parentCode := GetTestSCCode("exec-dest-ctx-parent", "../../")
	parentSC := []byte("parentSC.........................")
	parentSCBalance := big.NewInt(1000)

	getBalanceCalled := func(address []byte) (*big.Int, error) {
		if bytes.Equal(parentSC, address) {
			return parentSCBalance, nil
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
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionPrepare"
	input.GasProvided = 1000000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	fmt.Println(vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	expectedVMOutput := expectedVMOutput_DestCtx_Prepare()
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnDestContext_Wrong(t *testing.T) {
	parentCode := GetTestSCCode("exec-dest-ctx-parent", "../../")
	parentSCBalance := big.NewInt(1000)

	getBalanceCalled := func(address []byte) (*big.Int, error) {
		if bytes.Equal(parentAddress, address) {
			return parentSCBalance, nil
		}
		return big.NewInt(0), nil
	}

	// Call parentFunctionWrongCall() of the parent SC, which will try to call a
	// non-existing SC.
	host, stubBlockchainHook := DefaultTestArwenForCall(t, parentCode)
	stubBlockchainHook.GetBalanceCalled = getBalanceCalled
	input := DefaultTestContractCallInput()
	input.CallerAddr = []byte("user")
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionWrongCall"
	input.GasProvided = 1000000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	expectedVMOutput := expectedVMOutput_DestCtx_WrongContractCalled()
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnDestContext_OutOfGas(t *testing.T) {
	// Scenario:
	// Parent sets data into the storage, finishes data and creates a bigint
	// Parent calls executeOnDestContext, sending some value as well
	// Parent provides insufficient gas to executeOnDestContext (enoguh to start the SC though)
	// Child SC starts executing: sets data into the storage, finishes data and changes the bigint
	// Child starts an infinite loop, which must surely end with OutOfGas
	// Execution returns to parent, which finishes with the result of executeOnDestContext
	// Assertions: modifications made by the child are did not take effect (no OutputAccount is created)
	// Assertions: the value sent by the parent to the child was returned to the parent
	// Assertions: the parent lost all the gas provided to executeOnDestContext
	parentCode := GetTestSCCode("exec-dest-ctx-parent", "../../")
	childCode := GetTestSCCode("exec-dest-ctx-child", "../../")
	parentSCBalance := big.NewInt(1000)

	getBalanceCalled := func(address []byte) (*big.Int, error) {
		if bytes.Equal(parentAddress, address) {
			return parentSCBalance, nil
		}

		return big.NewInt(0), nil
	}

	// Call parentFunctionChildCall_OutOfGas() of the parent SC, which will call
	// the child SC using executeOnDestContext() with sufficient gas for
	// compilation and starting, but the child starts an infinite loop which will
	// end in OutOfGas.
	host, stubBlockchainHook := DefaultTestArwenForTwoSCs(t, parentCode, childCode)
	stubBlockchainHook.GetBalanceCalled = getBalanceCalled
	input := DefaultTestContractCallInput()
	input.CallerAddr = []byte("user")
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall_OutOfGas"
	input.GasProvided = 1000000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	expectedVMOutput := expectedVMOutput_DestCtx_OutOfGas()
	assert.Equal(t, int64(42), host.BigInt().GetOne(12).Int64())
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnDestContext_Successful(t *testing.T) {
	parentCode := GetTestSCCode("exec-dest-ctx-parent", "../../")
	childCode := GetTestSCCode("exec-dest-ctx-child", "../../")
	parentSCBalance := big.NewInt(1000)

	getBalanceCalled := func(address []byte) (*big.Int, error) {
		if bytes.Equal(parentAddress, address) {
			return parentSCBalance, nil
		}

		return big.NewInt(0), nil
	}

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().
	host, stubBlockchainHook := DefaultTestArwenForTwoSCs(t, parentCode, childCode)
	stubBlockchainHook.GetBalanceCalled = getBalanceCalled
	input := DefaultTestContractCallInput()
	input.CallerAddr = []byte("user")
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall"
	input.GasProvided = 1000000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	expectedVMOutput := expectedVMOutput_DestCtx_SuccessfulChildCall()
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnDestContext_Successful_BigInts(t *testing.T) {
	parentCode := GetTestSCCode("exec-dest-ctx-parent", "../../")
	childCode := GetTestSCCode("exec-dest-ctx-child", "../../")
	parentSCBalance := big.NewInt(1000)

	getBalanceCalled := func(address []byte) (*big.Int, error) {
		if bytes.Equal(parentAddress, address) {
			return parentSCBalance, nil
		}
		return big.NewInt(0), nil
	}

	// Call parentFunctionChildCall_BigInts() of the parent SC, which will call a
	// method of the child SC that takes some big Int references as arguments and
	// produce a new big Int out of the arguments.
	host, stubBlockchainHook := DefaultTestArwenForTwoSCs(t, parentCode, childCode)
	stubBlockchainHook.GetBalanceCalled = getBalanceCalled
	input := DefaultTestContractCallInput()
	input.CallerAddr = []byte("user")
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall_BigInts"
	input.GasProvided = 1000000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	expectedVMOutput := expectedVMOutput_DestCtx_SuccessfulChildCall_BigInts()
	require.Equal(t, expectedVMOutput, vmOutput)
}

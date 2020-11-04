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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var counterKey = []byte("COUNTER")
var WASMLocalsLimit = uint64(4000)

func TestNewArwen(t *testing.T) {
	host, err := DefaultTestArwen(t, &mock.BlockchainHookStub{})
	require.Nil(t, err)
	require.NotNil(t, host)
}

func TestSCMem(t *testing.T) {
	code := GetTestSCCode("misc", "../../")
	host, _ := DefaultTestArwenForCall(t, code, nil)
	input := DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = "iterate_over_byte_array"
	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	testString := "this is some random string of bytes"
	expectedData := [][]byte{
		[]byte(testString),
		{35},
	}
	for _, c := range testString {
		expectedData = append(expectedData, []byte{byte(c)})
	}
	require.Equal(t, expectedData, vmOutput.ReturnData)
}

func TestExecution_DeployNewAddressErr(t *testing.T) {
	stubBlockchainHook := &mock.BlockchainHookStub{}

	errNewAddress := errors.New("new address error")

	host, _ := DefaultTestArwen(t, stubBlockchainHook)
	input := DefaultTestContractCreateInput()
	stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		require.Equal(t, input.CallerAddr, address)
		return &mock.AccountMock{}, nil
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
	input.Arguments = [][]byte{{0}}

	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Len(t, vmOutput.ReturnData, 1)
	require.Equal(t, []byte("init successful"), vmOutput.ReturnData[0])
	require.Equal(t, uint64(528), vmOutput.GasRemaining)
	require.Len(t, vmOutput.OutputAccounts, 2)
	require.Equal(t, uint64(24), vmOutput.OutputAccounts["caller"].Nonce)
	require.Equal(t, input.ContractCode, vmOutput.OutputAccounts[string(newAddress)].Code)
	require.Equal(t, big.NewInt(88), vmOutput.OutputAccounts[string(newAddress)].BalanceDelta)
}

func TestExecution_DeployWASM_Popcnt(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := DefaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.CallValue = big.NewInt(88)
	input.GasProvided = 1000
	input.ContractCode = GetTestSCCode("init-simple-popcnt", "../../")
	input.Arguments = [][]byte{}

	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Len(t, vmOutput.ReturnData, 1)
	require.Equal(t, []byte{3}, vmOutput.ReturnData[0])
}

func TestExecution_DeployWASM_AtMaximumLocals(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := DefaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.CallValue = big.NewInt(88)
	input.GasProvided = 1000
	input.ContractCode = makeBytecodeWithLocals(WASMLocalsLimit)
	input.Arguments = [][]byte{{0}}

	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
}

func TestExecution_DeployWASM_MoreThanMaximumLocals(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := DefaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.CallValue = big.NewInt(88)
	input.GasProvided = 1000
	input.ContractCode = makeBytecodeWithLocals(WASMLocalsLimit + 1)
	input.Arguments = [][]byte{{0}}

	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractInvalid, vmOutput.ReturnCode)
}

func TestExecution_DeployWASM_Init_Errors(t *testing.T) {
	newAddress := []byte("new smartcontract")
	host := DefaultTestArwenForDeployment(t, 24, newAddress)
	input := DefaultTestContractCreateInput()
	input.CallValue = big.NewInt(88)
	input.GasProvided = 1000
	input.ContractCode = GetTestSCCode("init-correct", "../../")

	// init() calls signalError()
	input.Arguments = [][]byte{{1}}
	vmOutput, err := host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.UserError, vmOutput.ReturnCode)

	// init() starts an infinite loop
	input.Arguments = [][]byte{{2}}
	vmOutput, err = host.RunSmartContractCreate(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.OutOfGas, vmOutput.ReturnCode)
}

func TestExecution_ManyDeployments(t *testing.T) {
	ownerNonce := uint64(23)
	newAddress := "new smartcontract"
	stubBlockchainHook := &mock.BlockchainHookStub{}
	stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		return &mock.AccountMock{Nonce: ownerNonce}, nil
	}
	stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
		ownerNonce++
		return []byte(newAddress + " " + fmt.Sprint(ownerNonce)), nil
	}

	host, _ := DefaultTestArwen(t, stubBlockchainHook)
	input := DefaultTestContractCreateInput()
	input.CallerAddr = []byte("owner")
	input.Arguments = make([][]byte, 0)
	input.CallValue = big.NewInt(88)
	input.ContractCode = GetTestSCCode("init-simple", "../../")

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

func TestExecution_CallGetUserAccountErr(t *testing.T) {
	stubBlockchainHook := &mock.BlockchainHookStub{}

	errGetAccount := errors.New("get code error")

	host, _ := DefaultTestArwen(t, stubBlockchainHook)
	input := DefaultTestContractCallInput()
	stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		return nil, errGetAccount
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.ContractNotFound, vmOutput.ReturnCode)
	require.Equal(t, arwen.ErrContractNotFound.Error(), vmOutput.ReturnMessage)
}

func TestExecution_CallOutOfGas(t *testing.T) {
	code := GetTestSCCode("counter", "../../")
	host, _ := DefaultTestArwenForCall(t, code, nil)
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
	host, _ := DefaultTestArwenForCall(t, code, nil)
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
	host, _ := DefaultTestArwenForCall(t, code, nil)
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
	host, stubBlockchainHook := DefaultTestArwenForCall(t, code, nil)
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

func TestExecution_Call_GasConsumptionOnLocals(t *testing.T) {
	gasWithZeroLocals, gasSchedule := callCustomSCAndGetGasUsed(t, 0)
	costPerLocal := uint64(gasSchedule.WASMOpcodeCost.LocalAllocate)

	UnmeteredLocals := uint64(gasSchedule.WASMOpcodeCost.LocalsUnmetered)

	// Any number of local variables below `UnmeteredLocals` must be instantiated
	// without metering, i.e. gas-free.
	for _, locals := range []uint64{1, UnmeteredLocals / 2, UnmeteredLocals} {
		gasUsed, _ := callCustomSCAndGetGasUsed(t, uint64(locals))
		require.Equal(t, gasWithZeroLocals, gasUsed)
	}

	// Any number of local variables above `UnmeteredLocals` must be instantiated
	// with metering, i.e. will cost gas.
	for _, locals := range []uint64{UnmeteredLocals + 1, UnmeteredLocals * 2, UnmeteredLocals * 4} {
		gasUsed, _ := callCustomSCAndGetGasUsed(t, uint64(locals))
		metered_locals := locals - UnmeteredLocals
		costOfLocals := costPerLocal * uint64(metered_locals)
		expectedGasUsed := gasWithZeroLocals + costOfLocals
		require.Equal(t, expectedGasUsed, gasUsed)
	}
}

func callCustomSCAndGetGasUsed(t *testing.T, locals uint64) (uint64, *config.GasCost) {
	code := makeBytecodeWithLocals(uint64(locals))
	host, _ := DefaultTestArwenForCall(t, code, nil)
	gasSchedule := host.Metering().GasSchedule()

	gasLimit := uint64(100000)
	input := DefaultTestContractCallInput()
	input.GasProvided = gasLimit
	input.Function = "answer"

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	compilationCost := uint64(len(code)) * gasSchedule.BaseOperationCost.CompilePerByte
	return gasLimit - vmOutput.GasRemaining - compilationCost, gasSchedule
}

func TestExecution_ExecuteOnSameContext_Simple(t *testing.T) {
	parentCode := GetTestSCCode("exec-same-ctx-simple-parent", "../../")
	childCode := GetTestSCCode("exec-same-ctx-simple-child", "../../")

	host, _ := DefaultTestArwenForTwoSCs(t, parentCode, childCode, big.NewInt(1000))
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)

	expectedVMOutput := expectedVMOutput_SameCtx_Simple(parentCode, childCode)
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_Call_Breakpoints(t *testing.T) {
	t.Parallel()

	code := GetTestSCCode("breakpoint", "../../")
	host, _ := DefaultTestArwenForCall(t, code, nil)
	input := DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = "testFunc"

	// Send the number 15 to the SC, causing it to finish with the number 100
	input.Arguments = [][]byte{{15}}
	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	require.Equal(t, [][]byte{{100}}, vmOutput.ReturnData)

	// Send the number 1 to the SC, causing it to exit with ReturnMessage "exit
	// here" if the breakpoint mechanism works properly, or with the
	// ReturnMessage "exit later" if the breakpoint mechanism fails to stop the
	// execution.
	input.Arguments = [][]byte{{1}}
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

	// Execute the parent SC method "parentFunctionPrepare", which sets storage,
	// finish data and performs a transfer. This step validates the test to the
	// actual call to ExecuteOnSameContext().
	host, _ := DefaultTestArwenForCall(t, parentCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionPrepare"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	expectedVMOutput := expectedVMOutput_SameCtx_Prepare(parentCode)
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnSameContext_Wrong(t *testing.T) {
	parentCode := GetTestSCCode("exec-same-ctx-parent", "../../")
	parentSCBalance := big.NewInt(1000)

	// Call parentFunctionWrongCall() of the parent SC, which will try to call a
	// non-existing SC.
	host, _ := DefaultTestArwenForCall(t, parentCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionWrongCall"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
		expectedVMOutput := expectedVMOutput_SameCtx_WrongContractCalled(parentCode)
		require.Equal(t, expectedVMOutput, vmOutput)
	} else {
		require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
		require.Equal(t, "account not found", vmOutput.ReturnMessage)
		require.Zero(t, vmOutput.GasRemaining)
	}
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

	// Call parentFunctionChildCall_OutOfGas() of the parent SC, which will call
	// the child SC using executeOnSameContext() with sufficient gas for
	// compilation and starting, but the child starts an infinite loop which will
	// end in OutOfGas.
	host, _ := DefaultTestArwenForTwoSCs(t, parentCode, childCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall_OutOfGas"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
		expectedVMOutput := expectedVMOutput_SameCtx_OutOfGas(parentCode, childCode)
		assert.Equal(t, int64(42), host.BigInt().GetOne(0).Int64())
		require.Equal(t, expectedVMOutput, vmOutput)
	} else {
		require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
		require.Equal(t, arwen.ErrNotEnoughGas.Error(), vmOutput.ReturnMessage)
		require.Zero(t, vmOutput.GasRemaining)
	}
}

func TestExecution_ExecuteOnSameContext_Successful(t *testing.T) {
	parentCode := GetTestSCCode("exec-same-ctx-parent", "../../")
	childCode := GetTestSCCode("exec-same-ctx-child", "../../")
	parentSCBalance := big.NewInt(1000)

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnSameContext().
	host, _ := DefaultTestArwenForTwoSCs(t, parentCode, childCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	expectedVMOutput := expectedVMOutput_SameCtx_SuccessfulChildCall(parentCode, childCode)
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnSameContext_Successful_BigInts(t *testing.T) {
	parentCode := GetTestSCCode("exec-same-ctx-parent", "../../")
	childCode := GetTestSCCode("exec-same-ctx-child", "../../")
	parentSCBalance := big.NewInt(1000)

	// Call parentFunctionChildCall_BigInts() of the parent SC, which will call a
	// method of the child SC that takes some big Int references as arguments and
	// produce a new big Int out of the arguments.
	host, _ := DefaultTestArwenForTwoSCs(t, parentCode, childCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall_BigInts"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	expectedVMOutput := expectedVMOutput_SameCtx_SuccessfulChildCall_BigInts(parentCode, childCode)
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnSameContext_Recursive_Direct(t *testing.T) {
	// Scenario:
	// SC has a method "callRecursive" which takes a byte as argument (number of recursive calls)
	// callRecursive() saves to storage "keyNNN" → "valueNNN", where NNN is the argument
	// callRecursive() saves to storage a counter starting at 1, increased by every recursive call
	// callRecursive() creates a bigInt and increments it with every iteration
	// callRecursive() finishes "finishNNN" in each iteration
	// callRecursive() calls itself using executeOnSameContext(), with the argument decremented
	// callRecursive() handles argument == 0 as follows: saves to storage the
	//		value of the bigInt counter, then exits without recursive call
	// Assertions: the VMOutput must contain as many StorageUpdates as the argument requires
	// Assertions: the VMOutput must contain as many finished values as the argument requires
	// Assertions: there must be a StorageUpdate with the value of the bigInt counter
	code := GetTestSCCode("exec-same-ctx-recursive", "../../")
	scBalance := big.NewInt(1000)

	host, _ := DefaultTestArwenForCall(t, code, scBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "callRecursive"
	input.GasProvided = gasProvided

	recursiveCalls := byte(5)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)

	// TODO set proper gas calculation in the expectedVMOutput, like the other
	// tests
	expectedVMOutput := expectedVMOutput_SameCtx_Recursive_Direct(code, int(recursiveCalls))
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
	require.Equal(t, int64(recursiveCalls+1), host.BigInt().GetOne(16).Int64())
}

func TestExecution_ExecuteOnSameContext_Recursive_Direct_ErrMaxInstances(t *testing.T) {
	code := GetTestSCCode("exec-same-ctx-recursive", "../../")
	scBalance := big.NewInt(1000)

	host, _ := DefaultTestArwenForCall(t, code, scBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "callRecursive"
	input.GasProvided = gasProvided

	recursiveCalls := byte(11)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
		expectedVMOutput := expectedVMOutput_SameCtx_Recursive_Direct_ErrMaxInstances(code, int(recursiveCalls))
		expectedVMOutput.GasRemaining = vmOutput.GasRemaining
		require.Equal(t, expectedVMOutput, vmOutput)
		require.Equal(t, int64(1), host.BigInt().GetOne(16).Int64())
	} else {
		require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
		require.Equal(t, arwen.ErrExecutionFailed.Error(), vmOutput.ReturnMessage)
		require.Zero(t, vmOutput.GasRemaining)
	}
}

func TestExecution_ExecuteOnSameContext_Recursive_Mutual_Methods(t *testing.T) {
	// Scenario:
	// SC has a method "callRecursiveMutualMethods" which takes a byte as
	//		argument (number of recursive calls)
	// callRecursiveMutualMethods() sets the finish value "start recursive mutual calls"
	// callRecursiveMutualMethods() calls recursiveMethodA() on the same context,
	//		passing the argument

	// recursiveMethodA() saves to storage "AkeyNNN" → "AvalueNNN", where NNN is the argument
	// recursiveMethodA() saves to storage a counter starting at 1, increased by every recursive call
	// recursiveMethodA() creates a bigInt and increments it with every iteration
	// recursiveMethodA() finishes "AfinishNNN" in each iteration
	// recursiveMethodA() calls recursiveMethodB() with the argument decremented
	// recursiveMethodB() is a copy of recursiveMethodA()
	// when argument == 0, either of them will save to storage the
	//		value of the bigInt counter, then exits without recursive call
	// callRecursiveMutualMethods() sets the finish value "end recursive mutual calls" and exits
	// Assertions: the VMOutput must contain as many StorageUpdates as the argument requires
	// Assertions: the VMOutput must contain as many finished values as the argument requires
	// Assertions: there must be a StorageUpdate with the value of the bigInt counter
	code := GetTestSCCode("exec-same-ctx-recursive", "../../")
	scBalance := big.NewInt(1000)

	host, _ := DefaultTestArwenForCall(t, code, scBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "callRecursiveMutualMethods"
	input.GasProvided = gasProvided

	recursiveCalls := byte(5)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	// TODO set proper gas calculation in the expectedVMOutput, like the other
	// tests
	expectedVMOutput := expectedVMOutput_SameCtx_Recursive_MutualMethods(code, int(recursiveCalls))
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
	require.Equal(t, int64(recursiveCalls+1), host.BigInt().GetOne(16).Int64())
}

func TestExecution_ExecuteOnSameContext_Recursive_Mutual_SCs(t *testing.T) {
	// Scenario:
	// Parent has method parentCallChild()
	// Child has method childCallParent()
	// The two methods are identical, just named differently
	// The methods do the following:
	//		parent: save to storage "PkeyNNN" → "PvalueNNN"
	//		parent:	finish "PfinishNNN"
	//		child:	save to storage "CkeyNNN" → "CvalueNNN"
	//		child:	finish "CfinishNNN"
	//		both:		increment a shared bigInt counter
	//		both:		whoever exits must save the shared bigInt counter to storage
	parentCode := GetTestSCCode("exec-same-ctx-recursive-parent", "../../")
	childCode := GetTestSCCode("exec-same-ctx-recursive-child", "../../")
	parentSCBalance := big.NewInt(1000)

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().
	host, _ := DefaultTestArwenForTwoSCs(t, parentCode, childCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentCallsChild"
	input.GasProvided = gasProvided

	recursiveCalls := byte(5)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)

	// TODO set proper gas calculation in the expectedVMOutput, like the other
	// tests
	expectedVMOutput := expectedVMOutput_SameCtx_Recursive_MutualSCs(parentCode, childCode, int(recursiveCalls))
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
	require.Equal(t, int64(recursiveCalls+1), host.BigInt().GetOne(88).Int64())
}

func TestExecution_ExecuteOnSameContext_Recursive_Mutual_SCs_OutOfGas(t *testing.T) {
	parentCode := GetTestSCCode("exec-same-ctx-recursive-parent", "../../")
	childCode := GetTestSCCode("exec-same-ctx-recursive-child", "../../")
	parentSCBalance := big.NewInt(1000)

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().
	host, _ := DefaultTestArwenForTwoSCs(t, parentCode, childCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentCallsChild"
	input.GasProvided = 10000

	recursiveCalls := byte(5)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
		require.Equal(t, vmcommon.OutOfGas, vmOutput.ReturnCode)
		require.Equal(t, arwen.ErrNotEnoughGas.Error(), vmOutput.ReturnMessage)
	} else {
		require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
		require.Equal(t, arwen.ErrExecutionFailed.Error(), vmOutput.ReturnMessage)
	}
}

func TestExecution_ExecuteOnDestContext_Prepare(t *testing.T) {
	parentCode := GetTestSCCode("exec-dest-ctx-parent", "../../")
	parentSCBalance := big.NewInt(1000)

	// Execute the parent SC method "parentFunctionPrepare", which sets storage,
	// finish data and performs a transfer. This step validates the test to the
	// actual call to ExecuteOnSameContext().
	host, _ := DefaultTestArwenForCall(t, parentCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionPrepare"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	fmt.Println(vmOutput.ReturnMessage)
	require.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)

	expectedVMOutput := expectedVMOutput_DestCtx_Prepare(parentCode)
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnDestContext_Wrong(t *testing.T) {
	parentCode := GetTestSCCode("exec-dest-ctx-parent", "../../")
	parentSCBalance := big.NewInt(1000)

	// Call parentFunctionWrongCall() of the parent SC, which will try to call a
	// non-existing SC.
	host, _ := DefaultTestArwenForCall(t, parentCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionWrongCall"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
		expectedVMOutput := expectedVMOutput_DestCtx_WrongContractCalled(parentCode)
		require.Equal(t, expectedVMOutput, vmOutput)
	} else {
		require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
		require.Equal(t, "account not found", vmOutput.ReturnMessage)
		require.Zero(t, vmOutput.GasRemaining)
	}
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

	// Call parentFunctionChildCall_OutOfGas() of the parent SC, which will call
	// the child SC using executeOnDestContext() with sufficient gas for
	// compilation and starting, but the child starts an infinite loop which will
	// end in OutOfGas.
	host, _ := DefaultTestArwenForTwoSCs(t, parentCode, childCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall_OutOfGas"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
		expectedVMOutput := expectedVMOutput_DestCtx_OutOfGas(parentCode)
		require.Equal(t, expectedVMOutput, vmOutput)
		require.Equal(t, int64(42), host.BigInt().GetOne(12).Int64())
	} else {
		require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
		require.Equal(t, arwen.ErrNotEnoughGas.Error(), vmOutput.ReturnMessage)
		require.Zero(t, vmOutput.GasRemaining)
	}
}

func TestExecution_ExecuteOnDestContext_Successful(t *testing.T) {
	parentCode := GetTestSCCode("exec-dest-ctx-parent", "../../")
	childCode := GetTestSCCode("exec-dest-ctx-child", "../../")
	parentSCBalance := big.NewInt(1000)

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().
	host, _ := DefaultTestArwenForTwoSCs(t, parentCode, childCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	expectedVMOutput := expectedVMOutput_DestCtx_SuccessfulChildCall(parentCode, childCode)
	expectedVMOutput.OutputAccounts[string(parentAddress)].StorageUpdates[string(childKey)] = &vmcommon.StorageUpdate{Offset: childKey, Data: nil}
	assert.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnDestContext_Successful_BigInts(t *testing.T) {
	parentCode := GetTestSCCode("exec-dest-ctx-parent", "../../")
	childCode := GetTestSCCode("exec-dest-ctx-child", "../../")
	parentSCBalance := big.NewInt(1000)

	// Call parentFunctionChildCall_BigInts() of the parent SC, which will call a
	// method of the child SC that takes some big Int references as arguments and
	// produce a new big Int out of the arguments.
	host, _ := DefaultTestArwenForTwoSCs(t, parentCode, childCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall_BigInts"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	expectedVMOutput := expectedVMOutput_DestCtx_SuccessfulChildCall_BigInts(parentCode, childCode)
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_ExecuteOnDestContext_Recursive_Direct(t *testing.T) {
	code := GetTestSCCode("exec-dest-ctx-recursive", "../../")
	scBalance := big.NewInt(1000)

	host, _ := DefaultTestArwenForCall(t, code, scBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "callRecursive"
	input.GasProvided = gasProvided

	recursiveCalls := byte(6)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)

	// TODO set proper gas calculation in the expectedVMOutput, like the other
	// tests
	expectedVMOutput := expectedVMOutput_DestCtx_Recursive_Direct(code, int(recursiveCalls))
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
	require.Equal(t, int64(1), host.BigInt().GetOne(16).Int64())
}

func TestExecution_ExecuteOnDestContext_Recursive_Mutual_Methods(t *testing.T) {
	code := GetTestSCCode("exec-dest-ctx-recursive", "../../")
	scBalance := big.NewInt(1000)

	host, _ := DefaultTestArwenForCall(t, code, scBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "callRecursiveMutualMethods"
	input.GasProvided = gasProvided

	recursiveCalls := byte(7)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)

	// TODO set proper gas calculation in the expectedVMOutput, like the other
	// tests
	expectedVMOutput := expectedVMOutput_DestCtx_Recursive_MutualMethods(code, int(recursiveCalls))
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
	require.Equal(t, int64(0), host.BigInt().GetOne(16).Int64())
}

func TestExecution_ExecuteOnDestContext_Recursive_Mutual_SCs(t *testing.T) {
	parentCode := GetTestSCCode("exec-dest-ctx-recursive-parent", "../../")
	childCode := GetTestSCCode("exec-dest-ctx-recursive-child", "../../")
	parentSCBalance := big.NewInt(1000)

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().
	host, _ := DefaultTestArwenForTwoSCs(t, parentCode, childCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentCallsChild"
	input.GasProvided = gasProvided

	recursiveCalls := byte(6)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)

	// TODO set proper gas calculation in the expectedVMOutput, like the other
	// tests
	expectedVMOutput := expectedVMOutput_DestCtx_Recursive_MutualSCs(parentCode, childCode, int(recursiveCalls))
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
	require.Equal(t, int64(1), host.BigInt().GetOne(88).Int64())
}

func TestExecution_ExecuteOnDestContext_Recursive_Mutual_SCs_OutOfGas(t *testing.T) {
	parentCode := GetTestSCCode("exec-dest-ctx-recursive-parent", "../../")
	childCode := GetTestSCCode("exec-dest-ctx-recursive-child", "../../")
	parentSCBalance := big.NewInt(1000)

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().
	host, _ := DefaultTestArwenForTwoSCs(t, parentCode, childCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentCallsChild"
	input.GasProvided = 10000

	recursiveCalls := byte(5)
	input.Arguments = [][]byte{
		{recursiveCalls},
	}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
		require.Equal(t, vmcommon.OutOfGas, vmOutput.ReturnCode)
		require.Equal(t, arwen.ErrNotEnoughGas.Error(), vmOutput.ReturnMessage)
	} else {
		require.Equal(t, vmcommon.ExecutionFailed, vmOutput.ReturnCode)
		require.Equal(t, arwen.ErrExecutionFailed.Error(), vmOutput.ReturnMessage)
		require.Zero(t, vmOutput.GasRemaining)
	}
}

func TestExecution_AsyncCall(t *testing.T) {
	// Scenario
	// Parent SC calls Child SC
	// Before asyncCall, Parent sets storage, makes a value transfer to ThirdParty and finishes some data
	// Parent performs asyncCall to Child with a sufficient amount of ERD, with arguments:
	//	* the address of ThirdParty
	//	* number of ERD the Child should send to ThirdParty
	//  * a string, to be set as the data on the transfer to ThirdParty
	// Child stores the received arguments to storage
	// Child performs two transfers:
	//	* to ThirdParty, sending the amount of ERD specified as argument in asyncCall
	//	* to the Vault, a fixed address known by the Child, sending exactly 4 ERD with the data provided by Parent
	// Child finishes with "thirdparty" if the transfer to ThirdParty was successful
	// Child finishes with "vault" if the transfer to Vault was successful
	// Parent callBack() verifies its arguments and expects both "thirdparty" and "vault"
	// Assertions: OutputAccounts for
	//		* Parent: negative balance delta (payment for child + thirdparty + vault => 2), storage
	//		* Child: zero balance delta, storage
	//		* ThirdParty: positive balance delta
	//		* Vault
	parentCode := GetTestSCCode("async-call-parent", "../../")
	childCode := GetTestSCCode("async-call-child", "../../")
	parentSCBalance := big.NewInt(1000)

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().
	host, _ := DefaultTestArwenForTwoSCs(t, parentCode, childCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentPerformAsyncCall"
	input.GasProvided = 1_000_000
	input.Arguments = [][]byte{{0}}

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)

	// TODO calculate expected remaining gas properly, instead of copying it from
	// the actual vmOutput.
	expectedVMOutput := expectedVMOutput_AsyncCall(parentCode, childCode)
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_AsyncCall_ChildFails(t *testing.T) {
	// Scenario
	// Identical to TestExecution_AsyncCall(), except that the child is
	// instructed to call signalError().
	// Because "vault" was not received by the callBack(), the Parent sends 4 ERD
	// to the Vault directly.
	parentCode := GetTestSCCode("async-call-parent", "../../")
	childCode := GetTestSCCode("async-call-child", "../../")
	parentSCBalance := big.NewInt(1000)

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().
	host, _ := DefaultTestArwenForTwoSCs(t, parentCode, childCode, parentSCBalance)
	host.Metering().GasSchedule().ElrondAPICost.AsyncCallbackGasLock = 3000

	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentPerformAsyncCall"
	input.GasProvided = 1_000_000
	input.Arguments = [][]byte{{1}}
	input.CurrentTxHash = []byte("txhash")

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)

	// TODO calculate expected remaining gas properly, instead of copying it from
	// the actual vmOutput.
	expectedVMOutput := expectedVMOutput_AsyncCall_ChildFails(parentCode, childCode)
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_AsyncCall_CallBackFails(t *testing.T) {
	// Scenario
	// Identical to TestExecution_AsyncCall(), except that the child is
	// instructed to call signalError().
	// Because "vault" was not received by the callBack(), the Parent sends 4 ERD
	// to the Vault directly.
	parentCode := GetTestSCCode("async-call-parent", "../../")
	childCode := GetTestSCCode("async-call-child", "../../")
	parentSCBalance := big.NewInt(1000)

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().
	host, _ := DefaultTestArwenForTwoSCs(t, parentCode, childCode, parentSCBalance)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentPerformAsyncCall"
	input.GasProvided = 1_000_000
	input.Arguments = [][]byte{{0, 3}}
	input.CurrentTxHash = []byte("txhash")

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)

	// TODO calculate expected remaining gas properly, instead of copying it from
	// the actual vmOutput.
	expectedVMOutput := expectedVMOutput_AsyncCall_CallBackFails(parentCode, childCode)
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_CreateNewContract_Success(t *testing.T) {
	parentCode := GetTestSCCode("deployer", "../../")
	childCode := GetTestSCCode("init-correct", "../../")
	parentBalance := big.NewInt(1000)

	host, stubBlockchainHook := DefaultTestArwenForCall(t, parentCode, parentBalance)
	stubBlockchainHook.GetStorageDataCalled = func(address []byte, key []byte) ([]byte, error) {
		if bytes.Equal(address, parentAddress) {
			if bytes.Equal(key, []byte{'A'}) {
				return childCode, nil
			}
			return nil, nil
		}
		return nil, arwen.ErrInvalidAccount
	}

	input := DefaultTestContractCallInput()
	input.Function = "deployChildContract"
	input.Arguments = [][]byte{{'A'}, {0}}
	input.GasProvided = 1_000_000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)

	expectedVMOutput := expectedVMOutput_CreateNewContract_Success(parentCode, childCode)

	// TODO calculate expected remaining gas properly, instead of copying it from
	// the actual vmOutput.
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
}

func TestExecution_CreateNewContract_Fail(t *testing.T) {
	parentCode := GetTestSCCode("deployer", "../../")
	childCode := GetTestSCCode("init-correct", "../../")
	parentBalance := big.NewInt(1000)

	host, stubBlockchainHook := DefaultTestArwenForCall(t, parentCode, parentBalance)
	stubBlockchainHook.GetStorageDataCalled = func(address []byte, key []byte) ([]byte, error) {
		if bytes.Equal(address, parentAddress) {
			if bytes.Equal(key, []byte{'A'}) {
				return childCode, nil
			}
			return nil, nil
		}
		return nil, arwen.ErrInvalidAccount
	}

	input := DefaultTestContractCallInput()
	input.Function = "deployChildContract"
	input.Arguments = [][]byte{{'A'}, {1}}
	input.GasProvided = 1_000_000

	vmOutput, err := host.RunSmartContractCall(input)
	require.Nil(t, err)

	expectedVMOutput := expectedVMOutput_CreateNewContract_Fail(parentCode, childCode)

	// TODO calculate expected remaining gas properly, instead of copying it from
	// the actual vmOutput.
	expectedVMOutput.GasRemaining = vmOutput.GasRemaining
	require.Equal(t, expectedVMOutput, vmOutput)
}

// makeBytecodeWithLocals rewrites the bytecode of "answer" to change the
// number of i64 locals it instantiates
func makeBytecodeWithLocals(numLocals uint64) []byte {
	originalCode := GetTestSCCode("answer", "../../")
	firstSlice := originalCode[:0x5B]
	secondSlice := originalCode[0x5C:]

	encodedNumLocals := arwen.U64ToLEB128(numLocals)
	extraBytes := len(encodedNumLocals) - 1

	result := make([]byte, 0)
	result = append(result, firstSlice...)
	result = append(result, encodedNumLocals...)
	result = append(result, secondSlice...)

	result[0x57] = byte(int(result[0x57]) + extraBytes)
	result[0x59] = byte(int(result[0x59]) + extraBytes)

	return result
}

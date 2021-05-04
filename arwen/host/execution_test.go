package host

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

var counterKey = []byte("COUNTER")
var WASMLocalsLimit = uint64(4000)
var maxUint8AsInt = int(math.MaxUint8)
var newAddress = []byte("new smartcontract")

const (
	get                     = "get"
	increment               = "increment"
	callRecursive           = "callRecursive"
	parentCallsChild        = "parentCallsChild"
	parentPerformAsyncCall  = "parentPerformAsyncCall"
	parentFunctionChildCall = "parentFunctionChildCall"
)

func TestSCMem(t *testing.T) {

	testString := "this is some random string of bytes"
	returnData := [][]byte{
		[]byte(testString),
		{35},
	}
	for _, c := range testString {
		returnData = append(returnData, []byte{byte(c)})
	}

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("misc", "../../"))).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withGasProvided(100000).
			withFunction("iterate_over_byte_array").
			build()).
		andAssertResults(func(host *vmHost, blockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				ReturnData(returnData...)
		})
}

func TestExecution_DeployNewAddressErr(t *testing.T) {

	errNewAddress := errors.New("new address error")

	input := createTestContractCreateInputBuilder().
		withGasProvided(1000).
		withContractCode([]byte("contract")).
		build()

	runInstanceCreatorTestBuilder(t).
		withInput(input).
		withSetup(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
				require.Equal(t, input.CallerAddr, address)
				return &contextmock.StubAccount{}, nil
			}
			stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
				require.Equal(t, input.CallerAddr, creatorAddress)
				require.Equal(t, uint64(0), nonce)
				require.Equal(t, defaultVMType, vmType)
				return nil, errNewAddress
			}
		}).
		andAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ExecutionFailed).
				ReturnMessage(errNewAddress.Error())
		})
}

func TestExecution_DeployOutOfGas(t *testing.T) {
	runInstanceCreatorTestBuilder(t).
		withInput(createTestContractCreateInputBuilder().
			withGasProvided(8).
			build()).
		withAddress(newAddress).
		andAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.OutOfGas).
				ReturnMessage(arwen.ErrNotEnoughGas.Error())
		})
}

func TestExecution_DeployNotWASM(t *testing.T) {
	runInstanceCreatorTestBuilder(t).
		withInput(createTestContractCreateInputBuilder().
			withGasProvided(9).
			withContractCode([]byte("not WASM")).
			build()).
		withAddress(newAddress).
		andAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ContractInvalid)
		})
}

func TestExecution_DeployWASM_WithoutMemory(t *testing.T) {
	runInstanceCreatorTestBuilder(t).
		withInput(createTestContractCreateInputBuilder().
			withGasProvided(1000).
			withContractCode(GetTestSCCode("memoryless", "../../")).
			build()).
		withAddress(newAddress).
		andAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ContractInvalid)
		})
}

func TestExecution_DeployWASM_WrongInit(t *testing.T) {
	runInstanceCreatorTestBuilder(t).
		withInput(createTestContractCreateInputBuilder().
			withGasProvided(1000).
			withContractCode(GetTestSCCode("init-wrong", "../../")).
			build()).
		withAddress(newAddress).
		andAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ContractInvalid)
		})
}

func TestExecution_DeployWASM_WrongMethods(t *testing.T) {
	runInstanceCreatorTestBuilder(t).
		withInput(createTestContractCreateInputBuilder().
			withGasProvided(1000).
			withContractCode(GetTestSCCode("signatures", "../../")).
			build()).
		withAddress(newAddress).
		andAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ContractInvalid)
		})
}

func TestExecution_DeployWASM_Successful(t *testing.T) {
	input := createTestContractCreateInputBuilder().
		withGasProvided(1000).
		withContractCode(GetTestSCCode("init-correct", "../../")).
		withCallValue(88).
		withArguments([]byte{0}).
		build()
	runInstanceCreatorTestBuilder(t).
		withInput(input).
		withAddress(newAddress).
		andAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				ReturnData([]byte("init successful")).
				GasRemaining(528).
				Nonce([]byte("caller"), 24).
				Code(newAddress, input.ContractCode).
				BalanceDelta(newAddress, 88)
		})
}

func TestExecution_DeployWASM_Popcnt(t *testing.T) {
	runInstanceCreatorTestBuilder(t).
		withInput(createTestContractCreateInputBuilder().
			withGasProvided(1000).
			withCallValue(88).
			withArguments().
			withContractCode(GetTestSCCode("init-simple-popcnt", "../../")).
			build()).
		withAddress(newAddress).
		andAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				ReturnData([]byte{3})
		})
}

func TestExecution_DeployWASM_AtMaximumLocals(t *testing.T) {
	runInstanceCreatorTestBuilder(t).
		withInput(createTestContractCreateInputBuilder().
			withGasProvided(1000).
			withCallValue(88).
			withContractCode(makeBytecodeWithLocals(WASMLocalsLimit)).
			build()).
		withAddress(newAddress).
		andAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok()
		})
}

func TestExecution_DeployWASM_MoreThanMaximumLocals(t *testing.T) {
	runInstanceCreatorTestBuilder(t).
		withInput(createTestContractCreateInputBuilder().
			withGasProvided(1000).
			withCallValue(88).
			withContractCode(makeBytecodeWithLocals(WASMLocalsLimit + 1)).
			build()).
		withAddress(newAddress).
		andAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ContractInvalid)
		})
}

func TestExecution_DeployWASM_Init_Errors(t *testing.T) {
	runInstanceCreatorTestBuilder(t).
		withInput(createTestContractCreateInputBuilder().
			withGasProvided(1000).
			withCallValue(88).
			withArguments([]byte{1}).
			withContractCode(GetTestSCCode("init-correct", "../../")).
			build()).
		withAddress(newAddress).
		andAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.UserError)
		})
}

func TestExecution_DeployWASM_Init_InfiniteLoop_Errors(t *testing.T) {
	runInstanceCreatorTestBuilder(t).
		withInput(createTestContractCreateInputBuilder().
			withGasProvided(1000).
			withCallValue(88).
			withArguments([]byte{2}).
			withContractCode(GetTestSCCode("init-correct", "../../")).
			build()).
		withAddress(newAddress).
		andAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.OutOfGas)
		})
}

func TestExecution_ManyDeployments(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ownerNonce := uint64(23)
	numDeployments := 1000

	for i := 0; i < numDeployments; i++ {
		runInstanceCreatorTestBuilder(t).
			withInput(createTestContractCreateInputBuilder().
				withGasProvided(100000).
				withCallValue(88).
				withCallerAddr([]byte("owner")).
				withContractCode(GetTestSCCode("init-simple", "../../")).
				build()).
			withAddress(newAddress).
			withSetup(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
				stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
					return &contextmock.StubAccount{Nonce: ownerNonce}, nil
				}
				stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
					ownerNonce++
					return []byte(string(newAddress) + " " + fmt.Sprint(ownerNonce)), nil
				}
			}).
			andAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
				verify.
					Ok()
			})
	}
}

func TestExecution_MultipleArwens_OverlappingContractInstanceData(t *testing.T) {
	code := GetTestSCCode("counter", "../../")

	input := DefaultTestContractCallInput()
	input.GasProvided = 1000000
	input.Function = get

	host1, instanceRecorder1 := defaultTestArwenForCallWithInstanceRecorderMock(t, code, nil)
	runtimeContextMock := contextmock.NewRuntimeContextWrapper(&host1.runtimeContext)
	runtimeContextMock.CleanWasmerInstanceFunc = func() {}
	host1.runtimeContext = runtimeContextMock

	for i := 0; i < 5; i++ {
		vmOutput, err := host1.RunSmartContractCall(input)
		verify := NewVMOutputVerifier(t, vmOutput, err)
		verify.Ok()
	}

	var host1InstancesData = make(map[interface{}]bool)
	for _, instance := range instanceRecorder1.GetContractInstances(code) {
		host1InstancesData[instance.GetData()] = true
	}

	host2, instanceRecorder2 := defaultTestArwenForCallWithInstanceRecorderMock(t, code, nil)
	runtimeContextMock = contextmock.NewRuntimeContextWrapper(&host2.runtimeContext)
	runtimeContextMock.CleanWasmerInstanceFunc = func() {}
	runtimeContextMock.GetSCCodeFunc = func() ([]byte, error) {
		return code, nil
	}
	host2.runtimeContext = runtimeContextMock

	for i := 0; i < maxUint8AsInt+1; i++ {
		vmOutput, err := host2.RunSmartContractCall(input)
		verify := NewVMOutputVerifier(t, vmOutput, err)
		verify.Ok()
	}

	for _, instance := range instanceRecorder2.GetContractInstances(code) {
		_, found := host1InstancesData[instance.GetData()]
		require.False(t, found)
	}
}

func TestExecution_MultipleArwens_CleanInstanceWhileOthersAreRunning(t *testing.T) {
	code := GetTestSCCode("counter", "../../")

	input := DefaultTestContractCallInput()
	input.GasProvided = 1000000
	input.Function = get

	interHostsChan := make(chan string)
	host1Chan := make(chan string)

	host1, _ := defaultTestArwenForCall(t, code, nil)
	runtimeContextMock := contextmock.NewRuntimeContextWrapper(&host1.runtimeContext)
	runtimeContextMock.FunctionFunc = func() string {
		interHostsChan <- "waitForHost2"
		return runtimeContextMock.GetWrappedRuntimeContext().Function()
	}
	host1.runtimeContext = runtimeContextMock

	var vmOutput1 *vmcommon.VMOutput
	var err1 error
	go func() {
		vmOutput1, err1 = host1.RunSmartContractCall(input)
		interHostsChan <- "finish"
		host1Chan <- "finish"
	}()

	host2, _ := defaultTestArwenForCall(t, code, nil)
	runtimeContextMock = contextmock.NewRuntimeContextWrapper(&host2.runtimeContext)
	runtimeContextMock.FunctionFunc = func() string {
		// wait to make sure host1 is running also
		<-interHostsChan
		// wait for host1 to finish
		<-interHostsChan
		return runtimeContextMock.GetWrappedRuntimeContext().Function()
	}
	host2.runtimeContext = runtimeContextMock

	vmOutput2, err2 := host2.RunSmartContractCall(input)

	<-host1Chan

	verify1 := NewVMOutputVerifier(t, vmOutput1, err1)
	verify1.Ok()

	verify2 := NewVMOutputVerifier(t, vmOutput2, err2)
	verify2.Ok()
}

func TestExecution_Deploy_DisallowFloatingPoint(t *testing.T) {
	runInstanceCreatorTestBuilder(t).
		withInput(createTestContractCreateInputBuilder().
			withGasProvided(1000).
			withCallValue(88).
			withArguments([]byte{2}).
			withContractCode(GetTestSCCode("num-with-fp", "../../")).
			build()).
		withAddress(newAddress).
		andAssertResults(func(blockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ContractInvalid)
		})
}

func TestExecution_CallGetUserAccountErr(t *testing.T) {
	errGetAccount := errors.New("get code error")
	runInstanceCallerTestBuilder(t).
		withInput(createTestContractCallInputBuilder().
			withGasProvided(100).
			build()).
		withSetup(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
				return nil, errGetAccount
			}
		}).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ContractNotFound).
				ReturnMessage(arwen.ErrContractNotFound.Error())
		})
}

func TestExecution_NotEnoughGasForGetCode(t *testing.T) {
	runInstanceCallerTestBuilder(t).
		withInput(createTestContractCallInputBuilder().
			withGasProvided(0).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.OutOfGas).
				ReturnMessage(arwen.ErrNotEnoughGas.Error())
		})
}

func TestExecution_CallOutOfGas(t *testing.T) {
	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("counter", "../../"))).
		withInput(createTestContractCallInputBuilder().
			withGasProvided(0).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.OutOfGas).
				ReturnMessage(arwen.ErrNotEnoughGas.Error())
		})
}

func TestExecution_CallWasmerError(t *testing.T) {
	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode([]byte("not WASM"))).
		withInput(createTestContractCallInputBuilder().
			withGasProvided(100000).
			withFunction(increment).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ContractInvalid)
		})
}

func TestExecution_CallSCMethod_Init(t *testing.T) {
	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("counter", "../../"))).
		withInput(createTestContractCallInputBuilder().
			withGasProvided(100000).
			withFunction("init").
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.UserError).
				ReturnMessage(arwen.ErrInitFuncCalledInRun.Error())
		})
}

func TestExecution_CallSCMethod_Callback(t *testing.T) {
	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("counter", "../../"))).
		withInput(createTestContractCallInputBuilder().
			withGasProvided(100000).
			withFunction("callBack").
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.UserError).
				ReturnMessage(arwen.ErrCallBackFuncCalledInRun.Error())
		})
}

func TestExecution_CallSCMethod_MissingFunction(t *testing.T) {
	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("counter", "../../"))).
		withInput(createTestContractCallInputBuilder().
			withGasProvided(100000).
			withFunction("wrong").
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.FunctionNotFound)
		})
}

func TestExecution_Call_Successful(t *testing.T) {
	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("counter", "../../"))).
		withInput(createTestContractCallInputBuilder().
			withGasProvided(100000).
			withFunction(increment).
			build()).
		withSetup(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.GetStorageDataCalled = func(scAddress []byte, key []byte) ([]byte, error) {
				return big.NewInt(1001).Bytes(), nil
			}
		}).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				Storage(
					createStoreEntry(parentAddress).withKey(counterKey).withValue(big.NewInt(1002).Bytes()),
				)
		})
}

func TestExecution_Call_GasConsumptionOnLocals(t *testing.T) {
	gasWithZeroLocals, gasSchedule := callCustomSCAndGetGasUsed(t, 0)
	costPerLocal := uint64(gasSchedule.WASMOpcodeCost.LocalAllocate)

	UnmeteredLocals := uint64(gasSchedule.WASMOpcodeCost.LocalsUnmetered)

	// Any number of local variables below `UnmeteredLocals` must be instantiated
	// without metering, i.e. gas-free.
	for _, locals := range []uint64{1, UnmeteredLocals / 2, UnmeteredLocals} {
		gasUsed, _ := callCustomSCAndGetGasUsed(t, locals)
		require.Equal(t, gasWithZeroLocals, gasUsed)
	}

	// Any number of local variables above `UnmeteredLocals` must be instantiated
	// with metering, i.e. will cost gas.
	for _, locals := range []uint64{UnmeteredLocals + 1, UnmeteredLocals * 2, UnmeteredLocals * 4} {
		gasUsed, _ := callCustomSCAndGetGasUsed(t, locals)
		meteredLocals := locals - UnmeteredLocals
		costOfLocals := costPerLocal * meteredLocals
		expectedGasUsed := gasWithZeroLocals + costOfLocals
		require.Equal(t, expectedGasUsed, gasUsed)
	}
}

func callCustomSCAndGetGasUsed(t *testing.T, locals uint64) (uint64, *config.GasCost) {

	var gasSchedule *config.GasCost
	var gasUsed uint64

	gasLimit := uint64(100000)
	code := makeBytecodeWithLocals(locals)

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(code)).
		withInput(createTestContractCallInputBuilder().
			withGasProvided(gasLimit).
			withFunction("answer").
			build()).
		withSetup(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			gasSchedule = host.Metering().GasSchedule()
		}).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			compilationCost := uint64(len(code)) * gasSchedule.BaseOperationCost.CompilePerByte
			gasUsed = gasLimit - verify.vmOutput.GasRemaining - compilationCost
			verify.
				Ok()
		})

	return gasUsed, gasSchedule
}

func TestExecution_ExecuteOnSameContext_Simple(t *testing.T) {

	parentGasUsed := uint64(521)
	childGasUsed := uint64(6870)
	executionCost := parentGasUsed + childGasUsed

	var returnData [][]byte

	returnData = append(returnData, []byte("child"))
	returnData = append(returnData, []byte{})
	for i := 1; i < 100; i++ {
		returnData = append(returnData, []byte{byte(i)})
	}
	returnData = append(returnData, []byte{})
	returnData = append(returnData, []byte("child"))
	returnData = append(returnData, []byte{})
	for i := 1; i < 100; i++ {
		returnData = append(returnData, []byte{byte(i)})
	}
	returnData = append(returnData, []byte{})
	returnData = append(returnData, []byte("parent"))

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-same-ctx-simple-parent", "../../")).
				withBalance(1000),
			createInstanceContract(childAddress).
				withCode(GetTestSCCode("exec-same-ctx-simple-child", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction(parentFunctionChildCall).
			withGasProvided(gasProvided).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				// parentAddress
				BalanceDelta(parentAddress, -198).
				GasUsed(parentAddress, parentGasUsed).
				// childAddress
				BalanceDelta(childAddress, 198).
				GasUsed(childAddress, childGasUsed).
				// other
				GasRemaining(gasProvided - executionCost).
				ReturnData(returnData...)
		})
}

func TestExecution_Call_Breakpoints(t *testing.T) {
	t.Parallel()

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("breakpoint", "../../"))).
		withInput(createTestContractCallInputBuilder().
			withGasProvided(100000).
			withFunction("testFunc").
			withArguments([]byte{15}).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				ReturnData([]byte{100})
		})
}

func TestExecution_Call_Breakpoints_UserError(t *testing.T) {
	t.Parallel()
	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("breakpoint", "../../"))).
		withInput(createTestContractCallInputBuilder().
			withGasProvided(100000).
			withFunction("testFunc").
			withArguments([]byte{1}).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				ReturnData().
				ReturnCode(vmcommon.UserError).
				ReturnMessage("exit here")
		})
}

func TestExecution_ExecuteOnSameContext_Prepare(t *testing.T) {
	expectedExecutionCost := uint64(138)
	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-same-ctx-parent", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction("parentFunctionPrepare").
			withGasProvided(gasProvided).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(parentAddress, 3405).
				Balance(parentAddress, 1000).
				BalanceDelta(parentAddress, -parentTransferValue).
				BalanceDelta(parentTransferReceiver, parentTransferValue).
				GasRemaining(gasProvided-
					parentCompilationCostSameCtx-
					expectedExecutionCost).
				ReturnData(parentFinishA, parentFinishB, []byte("succ")).
				Storage(
					createStoreEntry(parentAddress).withKey(parentKeyA).withValue(parentDataA),
					createStoreEntry(parentAddress).withKey(parentKeyB).withValue(parentDataB),
				).
				Transfers(
					createTransferEntry(parentTransferReceiver).
						withData(parentTransferData).
						withValue(big.NewInt(parentTransferValue)).
						withSenderAddress(parentAddress),
				)
		})
}

func TestExecution_ExecuteOnSameContext_Wrong(t *testing.T) {
	executionCostBeforeExecuteAPI := uint64(156)
	executeAPICost := uint64(39)
	gasLostOnFailure := uint64(50000)
	finalCost := uint64(44)

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-same-ctx-parent", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction("parentFunctionWrongCall").
			withGasProvided(gasProvided).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			if !host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() {
				verify.
					Ok().
					GasUsed(parentAddress, 3405).
					Balance(parentAddress, 1000).
					BalanceDelta(parentAddress, -parentTransferValue).
					BalanceDelta(parentTransferReceiver, parentTransferValue).
					GasRemaining(gasProvided-
						parentCompilationCostSameCtx-
						executionCostBeforeExecuteAPI-
						executeAPICost-
						gasLostOnFailure-
						finalCost).
					ReturnData(parentFinishA, parentFinishB, []byte("succ"), []byte("fail")).
					Storage(
						createStoreEntry(parentAddress).withKey(parentKeyA).withValue(parentDataA),
						createStoreEntry(parentAddress).withKey(parentKeyB).withValue(parentDataB),
						createStoreEntry(childAddress).withKey(childKey).withValue(childData),
					).
					Transfers(
						createTransferEntry(parentTransferReceiver).
							withData(parentTransferData).
							withValue(big.NewInt(parentTransferValue)).
							withSenderAddress(parentAddress),
					)
			} else {
				verify.
					ReturnCode(vmcommon.ExecutionFailed).
					ReturnMessage("account not found").
					GasRemaining(0)
			}
		})
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

	// Call parentFunctionChildCall_OutOfGas() of the parent SC, which will call
	// the child SC using executeOnSameContext() with sufficient gas for
	// compilation and starting, but the child starts an infinite loop which will
	// end in OutOfGas.

	executionCostBeforeExecuteAPI := uint64(90)
	executeAPICost := uint64(1)
	gasLostOnFailure := uint64(3500)
	finalCost := uint64(54)

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-same-ctx-parent", "../../")).
				withBalance(1000),
			createInstanceContract(childAddress).
				withCode(GetTestSCCode("exec-same-ctx-child", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction("parentFunctionChildCall_OutOfGas").
			withGasProvided(gasProvided).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			if !host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() {
				verify.
					Ok().
					Balance(parentAddress, 1000).
					BalanceDelta(parentAddress, 0).
					GasRemaining(gasProvided-
						parentCompilationCostSameCtx-
						executionCostBeforeExecuteAPI-
						executeAPICost-
						gasLostOnFailure-
						finalCost).
					ReturnData(parentFinishA, []byte("fail")).
					Storage(
						createStoreEntry(parentAddress).withKey(parentKeyA).withValue(parentDataA),
						createStoreEntry(parentAddress).withKey(parentKeyB).withValue(parentDataB),
					)
			} else {
				verify.
					ReturnCode(vmcommon.ExecutionFailed).
					ReturnMessage(arwen.ErrNotEnoughGas.Error()).
					GasRemaining(0)
			}
		})
}

func TestExecution_ExecuteOnSameContext_Successful(t *testing.T) {

	executeAPICost := uint64(39)
	childExecutionCost := uint64(437)
	parentGasBeforeExecuteAPI := uint64(172)
	finalCost := uint64(134)

	parentAccountBalance := int64(1000)

	returnData := [][]byte{parentFinishA, parentFinishB, []byte("succ")}
	returnData = append(returnData, childFinish, parentDataA)
	for _, c := range parentDataA {
		returnData = append(returnData, []byte{c})
	}
	returnData = append(returnData, parentDataB)
	for _, c := range parentDataB {
		returnData = append(returnData, []byte{c})
	}
	returnData = append(returnData, []byte("child ok"), []byte("succ"), []byte("succ"))

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnSameContext().

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-same-ctx-parent", "../../")).
				withBalance(parentAccountBalance),
			createInstanceContract(childAddress).
				withCode(GetTestSCCode("exec-same-ctx-child", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction(parentFunctionChildCall).
			withGasProvided(gasProvided).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				// parentAddress
				Balance(parentAddress, parentAccountBalance).
				BalanceDelta(parentAddress, -141).
				GasUsed(parentAddress, 3612).
				// childAddress
				Balance(childAddress, 1000).
				BalanceDelta(childAddress, 3).
				GasUsed(childAddress, childCompilationCostSameCtx+childExecutionCost).
				// others
				BalanceDelta(childTransferReceiver, 96).
				BalanceDelta(parentTransferReceiver, parentTransferValue).
				GasRemaining(gasProvided-
					parentCompilationCostSameCtx-
					parentGasBeforeExecuteAPI-
					executeAPICost-
					childCompilationCostSameCtx-
					childExecutionCost-
					finalCost).
				ReturnData(returnData...).
				Storage(
					createStoreEntry(parentAddress).withKey(parentKeyA).withValue(parentDataA),
					createStoreEntry(parentAddress).withKey(parentKeyB).withValue(parentDataB),
					createStoreEntry(parentAddress).withKey(childKey).withValue(childData),
				).
				Transfers(
					createTransferEntry(parentTransferReceiver).
						withData(parentTransferData).
						withValue(big.NewInt(parentTransferValue)).
						withSenderAddress(parentAddress),
					createTransferEntry(childTransferReceiver).
						withData([]byte("qwerty")).
						withValue(big.NewInt(96)).
						withSenderAddress(childAddress),
				)
		})
}

func TestExecution_ExecuteOnSameContext_Successful_BigInts(t *testing.T) {

	// Call parentFunctionChildCall_BigInts() of the parent SC, which will call a
	// method of the child SC that takes some big Int references as arguments and
	// produce a new big Int out of the arguments.

	childExecutionCost := uint64(108)
	parentGasBeforeExecuteAPI := uint64(114)
	executeAPICost := uint64(13)
	finalCost := uint64(67)

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-same-ctx-parent", "../../")).
				withBalance(1000),
			createInstanceContract(childAddress).
				withCode(GetTestSCCode("exec-same-ctx-child", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction("parentFunctionChildCall_BigInts").
			withGasProvided(gasProvided).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				// parentAddress
				Balance(parentAddress, 1000).
				BalanceDelta(parentAddress, -99).
				GasUsed(parentAddress, 3461).
				// childAddress
				BalanceDelta(childAddress, 99).
				GasUsed(childAddress, childCompilationCostSameCtx+childExecutionCost).
				// others
				GasRemaining(gasProvided-
					parentCompilationCostSameCtx-
					parentGasBeforeExecuteAPI-
					executeAPICost-
					childCompilationCostSameCtx-
					childExecutionCost-
					finalCost).
				ReturnData([]byte("child ok"), []byte("succ"), []byte("succ"))
		})
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

	recursiveCalls := 5
	var returnData [][]byte

	for i := recursiveCalls; i >= 0; i-- {
		finishString := fmt.Sprintf("Rfinish%03d", i)
		returnData = append(returnData, []byte(finishString))
	}
	for i := recursiveCalls - 1; i >= 0; i-- {
		returnData = append(returnData, []byte("succ"))
	}

	var storeEntries []storeEntry

	for i := 0; i <= recursiveCalls; i++ {
		key := fmt.Sprintf("Rkey%03d.........................", i)
		value := fmt.Sprintf("Rvalue%03d", i)
		storeEntries = append(storeEntries, createStoreEntry(parentAddress).withKey([]byte(key)).withValue([]byte(value)))
	}

	storeEntries = append(storeEntries, createStoreEntry(parentAddress).withKey(recursiveIterationCounterKey).withValue([]byte{byte(recursiveCalls + 1)}))
	storeEntries = append(storeEntries, createStoreEntry(parentAddress).withKey(recursiveIterationBigCounterKey).withValue(big.NewInt(int64(recursiveCalls+1)).Bytes()))

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-same-ctx-recursive", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction(callRecursive).
			withGasProvided(gasProvided).
			withArguments([]byte{byte(recursiveCalls)}).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				Balance(parentAddress, 1000).
				BalanceDelta(parentAddress, 0).
				GasUsed(parentAddress, 25424).
				ReturnData(returnData...).
				Storage(storeEntries...)

			require.Equal(t, int64(recursiveCalls+1), host.BigInt().GetOne(16).Int64())
		})
}

func TestExecution_ExecuteOnSameContext_Recursive_Direct_ErrMaxInstances(t *testing.T) {
	recursiveCalls := byte(11)
	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-same-ctx-recursive", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction(callRecursive).
			withGasProvided(gasProvided).
			withArguments([]byte{byte(recursiveCalls)}).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
				verify.
					Ok().
					Balance(parentAddress, 1000).
					BalanceDelta(parentAddress, 0).
					ReturnData(
						[]byte(fmt.Sprintf("Rfinish%03d", recursiveCalls)),
						[]byte("fail"),
					).
					Storage(
						createStoreEntry(parentAddress).
							withKey([]byte(fmt.Sprintf("Rkey%03d.........................", recursiveCalls))).
							withValue([]byte(fmt.Sprintf("Rvalue%03d", recursiveCalls))),
					)
				require.Equal(t, int64(1), host.BigInt().GetOne(16).Int64())
			} else {
				verify.
					ReturnCode(vmcommon.ExecutionFailed).
					ReturnMessage(arwen.ErrExecutionFailed.Error()).
					GasRemaining(0)
			}
		})
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

	recursiveCalls := 5

	var returnData [][]byte
	var storeEntries []storeEntry

	storeEntries = append(storeEntries, createStoreEntry(parentAddress).withKey(recursiveIterationCounterKey).withValue([]byte{byte(recursiveCalls + 1)}))
	storeEntries = append(storeEntries, createStoreEntry(parentAddress).withKey(recursiveIterationBigCounterKey).withValue(big.NewInt(int64(recursiveCalls+1)).Bytes()))

	returnData = append(returnData, []byte("start recursive mutual calls"))

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Afinish%03d", iteration)
			key = fmt.Sprintf("Akey%03d.........................", iteration)
			value = fmt.Sprintf("Avalue%03d", iteration)
		} else {
			finishData = fmt.Sprintf("Bfinish%03d", iteration)
			key = fmt.Sprintf("Bkey%03d.........................", iteration)
			value = fmt.Sprintf("Bvalue%03d", iteration)
		}
		storeEntries = append(storeEntries, createStoreEntry(parentAddress).withKey([]byte(key)).withValue([]byte(value)))
		returnData = append(returnData, []byte(finishData))
	}

	for i := recursiveCalls; i >= 0; i-- {
		returnData = append(returnData, []byte("succ"))
	}

	returnData = append(returnData, []byte("end recursive mutual calls"))

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-same-ctx-recursive", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction("callRecursiveMutualMethods").
			withGasProvided(gasProvided).
			withArguments([]byte{byte(recursiveCalls)}).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				Balance(parentAddress, 1000).
				BalanceDelta(parentAddress, (big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))).Int64()).
				GasUsed(parentAddress, 29593).
				ReturnData(returnData...).
				Storage(storeEntries...)

			require.Equal(t, int64(recursiveCalls+1), host.BigInt().GetOne(16).Int64())
		})
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

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().

	recursiveCalls := 4

	var expectedParentBalanceDelta, expectedChildBalanceDelta int64
	if recursiveCalls%2 == 1 {
		expectedParentBalanceDelta = -5
		expectedChildBalanceDelta = 5
	} else {
		expectedParentBalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1)).Int64()
		expectedChildBalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1)).Int64()
	}

	var returnData [][]byte
	var storeEntries []storeEntry

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Pfinish%03d", iteration)
			key = fmt.Sprintf("Pkey%03d.........................", iteration)
			value = fmt.Sprintf("Pvalue%03d", iteration)
		} else {
			finishData = fmt.Sprintf("Cfinish%03d", iteration)
			key = fmt.Sprintf("Ckey%03d.........................", iteration)
			value = fmt.Sprintf("Cvalue%03d", iteration)
		}
		storeEntries = append(storeEntries, createStoreEntry(parentAddress).withKey([]byte(key)).withValue([]byte(value)))
		returnData = append(returnData, []byte(finishData))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		returnData = append(returnData, []byte("succ"))
	}

	storeEntries = append(storeEntries, createStoreEntry(parentAddress).withKey(recursiveIterationCounterKey).withValue([]byte{byte(recursiveCalls + 1)}))
	storeEntries = append(storeEntries, createStoreEntry(parentAddress).withKey(recursiveIterationBigCounterKey).withValue(big.NewInt(int64(recursiveCalls+1)).Bytes()))

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-same-ctx-recursive-parent", "../../")).
				withBalance(1000),
			createInstanceContract(childAddress).
				withCode(GetTestSCCode("exec-same-ctx-recursive-child", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction(parentCallsChild).
			withGasProvided(gasProvided).
			withArguments([]byte{byte(recursiveCalls)}).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				// parentAddress
				Balance(parentAddress, 1000).
				BalanceDelta(parentAddress, expectedParentBalanceDelta).
				GasUsed(parentAddress, 5426).
				// childAddress
				Balance(childAddress, 1000).
				BalanceDelta(childAddress, expectedChildBalanceDelta).
				GasUsed(childAddress, 3652).
				// other
				ReturnData(returnData...).
				Storage(storeEntries...)

			require.Equal(t, int64(recursiveCalls+1), host.BigInt().GetOne(88).Int64())
		})
}

func TestExecution_ExecuteOnSameContext_Recursive_Mutual_SCs_OutOfGas(t *testing.T) {

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().
	recursiveCalls := byte(5)

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-same-ctx-recursive-parent", "../../")).
				withBalance(1000),
			createInstanceContract(childAddress).
				withCode(GetTestSCCode("exec-same-ctx-recursive-child", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction(parentCallsChild).
			withGasProvided(10000).
			withArguments([]byte{byte(recursiveCalls)}).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
				verify.
					ReturnCode(vmcommon.OutOfGas).
					ReturnMessage(arwen.ErrNotEnoughGas.Error()).
					GasRemaining(0)
			} else {
				verify.
					ReturnCode(vmcommon.ExecutionFailed).
					ReturnMessage(arwen.ErrExecutionFailed.Error()).
					GasRemaining(0)
			}
		})
}

func TestExecution_ExecuteOnDestContext_Prepare(t *testing.T) {

	// Execute the parent SC method "parentFunctionPrepare", which sets storage,
	// finish data and performs a transfer. This step validates the test to the
	// actual call to ExecuteOnSameContext().

	expectedExecutionCost := uint64(138)

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-dest-ctx-parent", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction("parentFunctionPrepare").
			withGasProvided(gasProvided).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				Balance(parentAddress, 1000).
				BalanceDelta(parentAddress, -parentTransferValue).
				GasUsed(parentAddress, 4309).
				BalanceDelta(parentTransferReceiver, parentTransferValue).
				GasRemaining(gasProvided-
					parentCompilationCostDestCtx-
					expectedExecutionCost).
				ReturnData(parentFinishA, parentFinishB, []byte("succ")).
				Storage(
					createStoreEntry(parentAddress).withKey(parentKeyA).withValue(parentDataA),
					createStoreEntry(parentAddress).withKey(parentKeyB).withValue(parentDataB),
				).
				Transfers(
					createTransferEntry(parentTransferReceiver).
						withData(parentTransferData).
						withValue(big.NewInt(parentTransferValue)).
						withSenderAddress(parentAddress),
				)
		})
}

func TestExecution_ExecuteOnDestContext_Wrong(t *testing.T) {
	// Call parentFunctionWrongCall() of the parent SC, which will try to call a
	// non-existing SC.

	executionCostBeforeExecuteAPI := uint64(156)
	executeAPICost := uint64(42)
	gasLostOnFailure := uint64(10000)
	finalCost := uint64(44)

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-dest-ctx-parent", "../../")).
				withBalance(1000)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction("parentFunctionWrongCall").
			withGasProvided(gasProvided).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
				verify.
					Ok().
					Balance(parentAddress, 1000).
					BalanceDelta(parentAddress, -42).
					GasUsed(parentAddress, 3612).
					BalanceDelta(childTransferReceiver, 96).
					BalanceDelta(parentTransferReceiver, parentTransferValue).
					GasRemaining(gasProvided-
						parentCompilationCostDestCtx-
						executionCostBeforeExecuteAPI-
						executeAPICost-
						gasLostOnFailure-
						finalCost).
					ReturnData(parentFinishA, parentFinishB, []byte("succ"), []byte("fail")).
					Storage(
						createStoreEntry(parentAddress).withKey(parentKeyA).withValue(parentDataA),
						createStoreEntry(parentAddress).withKey(parentKeyB).withValue(parentDataB),
						createStoreEntry(parentAddress).withKey(childKey).withValue(childData),
					).
					Transfers(
						createTransferEntry(childTransferReceiver).
							withData([]byte("qwerty")).
							withValue(big.NewInt(96)).
							withSenderAddress(childAddress),
						createTransferEntry(parentTransferReceiver).
							withData(parentTransferData).
							withValue(big.NewInt(parentTransferValue)).
							withSenderAddress(parentAddress),
					)
			} else {
				verify.
					ReturnCode(vmcommon.ExecutionFailed).
					ReturnMessage("account not found").
					GasRemaining(0)
			}
		})
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

	// Call parentFunctionChildCall_OutOfGas() of the parent SC, which will call
	// the child SC using executeOnDestContext() with sufficient gas for
	// compilation and starting, but the child starts an infinite loop which will
	// end in OutOfGas.

	executionCostBeforeExecuteAPI := uint64(90)
	executeAPICost := uint64(1)
	gasLostOnFailure := uint64(3500)
	finalCost := uint64(54)

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-dest-ctx-parent", "../../")).
				withBalance(1000),
			createInstanceContract(childAddress).
				withCode(GetTestSCCode("exec-dest-ctx-child", "../../")).
				withBalance(1000),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction("parentFunctionChildCall_OutOfGas").
			withGasProvided(gasProvided).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
				verify.
					Ok().
					Balance(parentAddress, 1000).
					GasRemaining(gasProvided-
						parentCompilationCostDestCtx-
						executionCostBeforeExecuteAPI-
						executeAPICost-
						gasLostOnFailure-
						finalCost).
					ReturnData(parentFinishA, []byte("fail")).
					Storage(
						createStoreEntry(parentAddress).withKey(parentKeyA).withValue(parentDataA),
						createStoreEntry(parentAddress).withKey(parentKeyB).withValue(parentDataB),
					)
				require.Equal(t, int64(42), host.BigInt().GetOne(12).Int64())
			} else {
				verify.
					ReturnCode(vmcommon.ExecutionFailed).
					ReturnMessage(arwen.ErrNotEnoughGas.Error()).
					GasRemaining(0)
			}
		})
}

func TestExecution_ExecuteOnDestContext_Successful(t *testing.T) {

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().

	parentGasBeforeExecuteAPI := uint64(168)
	executeAPICost := uint64(42)
	childExecutionCost := uint64(91)
	finalCost := uint64(65)
	childTransferValue := int64(12)

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-dest-ctx-parent", "../../")).
				withBalance(1000),
			createInstanceContract(childAddress).
				withCode(GetTestSCCode("exec-dest-ctx-child", "../../")).
				withBalance(1000),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction(parentFunctionChildCall).
			withGasProvided(gasProvided).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				// parentAddress
				Balance(parentAddress, 1000).
				BalanceDelta(parentAddress, -141).
				GasUsed(parentAddress, 4444).
				/// childAddress
				Balance(childAddress, 1000).
				BalanceDelta(childAddress, 99-childTransferValue).
				GasUsed(childAddress, 2256).
				// other
				BalanceDelta(childTransferReceiver, childTransferValue).
				GasRemaining(gasProvided-
					parentCompilationCostDestCtx-
					parentGasBeforeExecuteAPI-
					executeAPICost-
					childCompilationCostDestCtx-
					childExecutionCost-
					finalCost).
				ReturnData(parentFinishA, parentFinishB, []byte("succ"), childFinish, []byte("succ"), []byte("succ")).
				Storage(
					createStoreEntry(parentAddress).withKey(parentKeyA).withValue(parentDataA),
					createStoreEntry(parentAddress).withKey(parentKeyB).withValue(parentDataB),
					createStoreEntry(parentAddress).withKey(childKey).withValue(nil),
					createStoreEntry(childAddress).withKey(childKey).withValue(childData),
				).
				Transfers(
					createTransferEntry(childTransferReceiver).
						withData([]byte("Second sentence.")).
						withValue(big.NewInt(childTransferValue)).
						withSenderAddress(childAddress),
					createTransferEntry(parentTransferReceiver).
						withData(parentTransferData).
						withValue(big.NewInt(parentTransferValue)).
						withSenderAddress(parentAddress),
				)
		})
}

func TestExecution_ExecuteOnDestContext_Successful_ChildReturns(t *testing.T) {

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().

	parentGasBeforeExecuteAPI := uint64(168)
	executeAPICost := uint64(42)
	childExecutionCost := uint64(91)
	parentGasAfterExecuteAPI := uint64(273)
	childTransferValue := int64(12)

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-dest-ctx-parent", "../../")).
				withBalance(1000),
			createInstanceContract(childAddress).
				withCode(GetTestSCCode("exec-dest-ctx-child", "../../")).
				withBalance(1000),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction("parentFunctionChildCall_ReturnedData").
			withGasProvided(gasProvided).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				// parentAddress
				Balance(parentAddress, 1000).
				BalanceDelta(parentAddress, -141).
				GasUsed(parentAddress, 4652).
				/// childAddress
				Balance(childAddress, 1000).
				BalanceDelta(childAddress, 99-childTransferValue).
				GasUsed(childAddress, 2256).
				// other
				BalanceDelta(childTransferReceiver, childTransferValue).
				GasRemaining(gasProvided-
					parentCompilationCostDestCtx-
					parentGasBeforeExecuteAPI-
					executeAPICost-
					childCompilationCostDestCtx-
					childExecutionCost-
					parentGasAfterExecuteAPI).
				ReturnData(parentFinishA, parentFinishB, []byte("succ"), childFinish, []byte("succ")).
				Storage(
					createStoreEntry(parentAddress).withKey(parentKeyA).withValue(parentDataA),
					createStoreEntry(parentAddress).withKey(parentKeyB).withValue(parentDataB),
					createStoreEntry(childAddress).withKey(childKey).withValue(childData),
				).
				Transfers(
					createTransferEntry(childTransferReceiver).
						withData([]byte("Second sentence.")).
						withValue(big.NewInt(childTransferValue)).
						withSenderAddress(childAddress),
					createTransferEntry(parentTransferReceiver).
						withData(parentTransferData).
						withValue(big.NewInt(parentTransferValue)).
						withSenderAddress(parentAddress),
				)
		})
}

func TestExecution_ExecuteOnDestContext_GasRemaining(t *testing.T) {
	// This test ensures that host.ExecuteOnDestContext() calls
	// metering.GasLeft() on the Wasmer instance of the child, and not of the
	// parent.

	parentCode := GetTestSCCode("exec-dest-ctx-parent", "../../")
	childCode := GetTestSCCode("exec-dest-ctx-child", "../../")

	// Pretend that the execution of the parent SC was requested, with the
	// following ContractCallInput:
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "parentFunctionChildCall"
	input.GasProvided = gasProvided

	// Initialize the VM with the parent SC and child SC, but without really
	// executing the parent. The initialization emulates the behavior of
	// host.doRunSmartContractCall(). Gas cost for compilation is skipped.
	host, _ := defaultTestArwenForTwoSCs(t, parentCode, childCode, nil, nil)
	host.InitState()

	_, _, metering, output, runtime, storage := host.GetContexts()
	runtime.InitStateFromContractCallInput(input)
	output.AddTxValueToAccount(input.RecipientAddr, input.CallValue)
	storage.SetAddress(runtime.GetSCAddress())
	_ = metering.DeductInitialGasForExecution([]byte{})

	contract, err := runtime.GetSCCode()
	require.Nil(t, err)

	vmInput := runtime.GetVMInput()
	err = runtime.StartWasmerInstance(contract, vmInput.GasProvided, false)
	require.Nil(t, err)

	// Use a lot of gas on the parent contract
	metering.UseGas(500000)
	require.Equal(t, input.GasProvided-500001, metering.GasLeft())

	// Create a second ContractCallInput, used to call the child SC using
	// host.ExecuteOnDestContext().
	childInput := DefaultTestContractCallInput()
	childInput.CallerAddr = parentAddress
	childInput.CallValue = big.NewInt(99)
	childInput.Function = "childFunction"
	childInput.RecipientAddr = childAddress
	childInput.Arguments = [][]byte{
		[]byte("some data"),
		[]byte("argument"),
		[]byte("another argument"),
	}
	childInput.GasProvided = 10000

	childOutput, _, err := host.ExecuteOnDestContext(childInput)
	verify := NewVMOutputVerifier(t, childOutput, err)
	verify.
		Ok().
		GasRemaining(7752)

	host.Clean()
}

func TestExecution_ExecuteOnDestContext_Successful_BigInts(t *testing.T) {

	// Call parentFunctionChildCall_BigInts() of the parent SC, which will call a
	// method of the child SC that takes some big Int references as arguments and
	// produce a new big Int out of the arguments.

	parentGasBeforeExecuteAPI := uint64(115)
	executeAPICost := uint64(13)
	childExecutionCost := uint64(101)
	finalCost := uint64(68)

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-dest-ctx-parent", "../../")).
				withBalance(1000),
			createInstanceContract(childAddress).
				withCode(GetTestSCCode("exec-dest-ctx-child", "../../")).
				withBalance(1000),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction("parentFunctionChildCall_BigInts").
			withGasProvided(gasProvided).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				// parentAddress
				Balance(parentAddress, 1000).
				BalanceDelta(parentAddress, -99).
				GasUsed(parentAddress, 4366).
				/// childAddress
				BalanceDelta(childAddress, 99).
				GasUsed(childAddress, 2265).
				// other
				GasRemaining(gasProvided-
					parentCompilationCostDestCtx-
					parentGasBeforeExecuteAPI-
					executeAPICost-
					childCompilationCostDestCtx-
					childExecutionCost-
					finalCost).
				ReturnData([]byte("child ok"), []byte("succ"), []byte("succ"))
		})
}

func TestExecution_ExecuteOnDestContext_Recursive_Direct(t *testing.T) {

	recursiveCalls := 6

	var returnData [][]byte
	var storeEntries []storeEntry

	for i := recursiveCalls; i >= 0; i-- {
		finishString := fmt.Sprintf("Rfinish%03d", i)
		returnData = append(returnData, []byte(finishString))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		returnData = append(returnData, []byte("succ"))
	}

	for i := 0; i <= recursiveCalls; i++ {
		key := fmt.Sprintf("Rkey%03d.........................", i)
		value := fmt.Sprintf("Rvalue%03d", i)
		storeEntries = append(storeEntries, createStoreEntry(parentAddress).withKey([]byte(key)).withValue([]byte(value)))
	}

	storeEntries = append(storeEntries,
		createStoreEntry(parentAddress).withKey(recursiveIterationCounterKey).withValue([]byte{byte(recursiveCalls + 1)}),
		createStoreEntry(parentAddress).withKey(recursiveIterationBigCounterKey).withValue(big.NewInt(int64(1)).Bytes()))

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-dest-ctx-recursive", "../../")).
				withBalance(1000),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction(callRecursive).
			withGasProvided(gasProvided).
			withArguments([]byte{byte(recursiveCalls)}).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				Balance(parentAddress, 1000).
				BalanceDelta(parentAddress, big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1)).Int64()).
				GasUsed(parentAddress, 29670).
				ReturnData(returnData...).
				Storage(storeEntries...)

			require.Equal(t, int64(1), host.BigInt().GetOne(16).Int64())
		})
}

func TestExecution_ExecuteOnDestContext_Recursive_Mutual_Methods(t *testing.T) {

	recursiveCalls := 7

	var returnData [][]byte
	var storeEntries []storeEntry

	storeEntries = append(storeEntries, createStoreEntry(parentAddress).withKey(recursiveIterationCounterKey).withValue([]byte{byte(recursiveCalls + 1)}))
	storeEntries = append(storeEntries, createStoreEntry(parentAddress).withKey(recursiveIterationBigCounterKey).withValue(big.NewInt(int64(1)).Bytes()))

	returnData = append(returnData, []byte("start recursive mutual calls"))

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Afinish%03d", iteration)
			key = fmt.Sprintf("Akey%03d.........................", iteration)
			value = fmt.Sprintf("Avalue%03d", iteration)
		} else {
			finishData = fmt.Sprintf("Bfinish%03d", iteration)
			key = fmt.Sprintf("Bkey%03d.........................", iteration)
			value = fmt.Sprintf("Bvalue%03d", iteration)
		}
		storeEntries = append(storeEntries, createStoreEntry(parentAddress).withKey([]byte(key)).withValue([]byte(value)))
		returnData = append(returnData, []byte(finishData))
	}

	for i := recursiveCalls; i >= 0; i-- {
		returnData = append(returnData, []byte("succ"))
	}

	returnData = append(returnData, []byte("end recursive mutual calls"))

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-dest-ctx-recursive", "../../")).
				withBalance(1000),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction("callRecursiveMutualMethods").
			withGasProvided(gasProvided).
			withArguments([]byte{byte(recursiveCalls)}).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				Balance(parentAddress, 1000).
				BalanceDelta(parentAddress, big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1)).Int64()).
				GasUsed(parentAddress, 38083).
				ReturnData(returnData...).
				Storage(storeEntries...)

			require.Equal(t, int64(0), host.BigInt().GetOne(16).Int64())
		})
}

func TestExecution_ExecuteOnDestContext_Recursive_Mutual_SCs(t *testing.T) {
	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().

	recursiveCalls := 6

	parentIterations := (recursiveCalls / 2) + (recursiveCalls % 2)
	childIterations := recursiveCalls - parentIterations
	balanceDelta := int64(5*parentIterations - 3*childIterations)

	var returnData [][]byte
	var storeEntries []storeEntry

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Pfinish%03d", iteration)
			key = fmt.Sprintf("Pkey%03d.........................", iteration)
			value = fmt.Sprintf("Pvalue%03d", iteration)
			storeEntries = append(storeEntries, createStoreEntry(parentAddress).withKey([]byte(key)).withValue([]byte(value)))
		} else {
			finishData = fmt.Sprintf("Cfinish%03d", iteration)
			key = fmt.Sprintf("Ckey%03d.........................", iteration)
			value = fmt.Sprintf("Cvalue%03d", iteration)
			storeEntries = append(storeEntries, createStoreEntry(childAddress).withKey([]byte(key)).withValue([]byte(value)))
		}
		returnData = append(returnData, []byte(finishData))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		returnData = append(returnData, []byte("succ"))
	}

	counterValue := (recursiveCalls + recursiveCalls%2) / 2

	storeEntries = append(storeEntries, createStoreEntry(parentAddress).withKey(recursiveIterationCounterKey).withValue([]byte{byte(counterValue + 1)}))
	storeEntries = append(storeEntries, createStoreEntry(childAddress).withKey(recursiveIterationCounterKey).withValue(big.NewInt(int64(counterValue)).Bytes()))

	if recursiveCalls%2 == 0 {
		storeEntries = append(storeEntries, createStoreEntry(parentAddress).withKey(recursiveIterationBigCounterKey).withValue(big.NewInt(int64(1)).Bytes()))
	} else {
		storeEntries = append(storeEntries, createStoreEntry(childAddress).withKey(recursiveIterationBigCounterKey).withValue(big.NewInt(int64(1)).Bytes()))
	}

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-dest-ctx-recursive-parent", "../../")).
				withBalance(1000),
			createInstanceContract(childAddress).
				withCode(GetTestSCCode("exec-dest-ctx-recursive-child", "../../")).
				withBalance(1000),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction(parentCallsChild).
			withGasProvided(gasProvided).
			withArguments([]byte{byte(recursiveCalls)}).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				// parentAddress
				Balance(parentAddress, 1000).
				BalanceDelta(parentAddress, -balanceDelta).
				GasUsed(parentAddress, 7252).
				// childAddress
				Balance(childAddress, 1000).
				BalanceDelta(childAddress, balanceDelta).
				GasUsed(childAddress, 5464).
				// others
				ReturnData(returnData...).
				Storage(storeEntries...)

			require.Equal(t, int64(1), host.BigInt().GetOne(88).Int64())
		})
}

func TestExecution_ExecuteOnDestContext_Recursive_Mutual_SCs_OutOfGas(t *testing.T) {
	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using executeOnDestContext().

	recursiveCalls := byte(5)

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("exec-dest-ctx-recursive-parent", "../../")).
				withBalance(1000),
			createInstanceContract(childAddress).
				withCode(GetTestSCCode("exec-dest-ctx-recursive-child", "../../")).
				withBalance(1000),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction(parentCallsChild).
			withGasProvided(10000).
			withArguments([]byte{byte(recursiveCalls)}).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			if host.Runtime().ElrondSyncExecAPIErrorShouldFailExecution() == false {
				verify.
					ReturnCode(vmcommon.OutOfGas).
					ReturnMessage(arwen.ErrNotEnoughGas.Error())
			} else {
				verify.
					ReturnCode(vmcommon.ExecutionFailed).
					ReturnMessage(arwen.ErrExecutionFailed.Error()).
					GasRemaining(0)
			}
		})
}

func TestExecution_ExecuteOnSameContext_MultipleChildren(t *testing.T) {
	world := worldmock.NewMockWorld()
	host := defaultTestArwen(t, world)

	alphaCode := GetTestSCCodeModule("exec-sync-ctx-multiple/alpha", "alpha", "../../")
	alpha := AddTestSmartContractToWorld(world, "alphaSC", alphaCode)
	alpha.Balance = big.NewInt(100)

	betaCode := GetTestSCCodeModule("exec-sync-ctx-multiple/beta", "beta", "../../")
	gammaCode := GetTestSCCodeModule("exec-sync-ctx-multiple/gamma", "gamma", "../../")
	deltaCode := GetTestSCCodeModule("exec-sync-ctx-multiple/delta", "delta", "../../")

	_ = AddTestSmartContractToWorld(world, "betaSC", betaCode)
	_ = AddTestSmartContractToWorld(world, "gammaSC", gammaCode)
	_ = AddTestSmartContractToWorld(world, "deltaSC", deltaCode)

	expectedReturnData := [][]byte{
		[]byte("arg1"),
		[]byte("succ"),
		[]byte("arg2"),
		[]byte("succ"),
		[]byte("arg3"),
		[]byte("succ"),
	}

	// Alpha uses executeOnSameContext() to call beta, gamma and delta one after
	// the other, in the same transaction.
	input := DefaultTestContractCallInput()
	input.Function = "callChildrenDirectly_SameCtx"
	input.GasProvided = 1000000
	input.RecipientAddr = alpha.Address

	vmOutput, err := host.RunSmartContractCall(input)

	verify := NewVMOutputVerifier(t, vmOutput, err)
	verify.
		Ok().
		ReturnData(expectedReturnData...)
}

func TestExecution_ExecuteOnDestContext_MultipleChildren(t *testing.T) {
	world := worldmock.NewMockWorld()
	host := defaultTestArwen(t, world)

	alphaCode := GetTestSCCodeModule("exec-sync-ctx-multiple/alpha", "alpha", "../../")
	alpha := AddTestSmartContractToWorld(world, "alphaSC", alphaCode)
	alpha.Balance = big.NewInt(100)

	betaCode := GetTestSCCodeModule("exec-sync-ctx-multiple/beta", "beta", "../../")
	gammaCode := GetTestSCCodeModule("exec-sync-ctx-multiple/gamma", "gamma", "../../")
	deltaCode := GetTestSCCodeModule("exec-sync-ctx-multiple/delta", "delta", "../../")

	_ = AddTestSmartContractToWorld(world, "betaSC", betaCode)
	_ = AddTestSmartContractToWorld(world, "gammaSC", gammaCode)
	_ = AddTestSmartContractToWorld(world, "deltaSC", deltaCode)

	expectedReturnData := [][]byte{
		[]byte("arg1"),
		[]byte("succ"),
		[]byte("arg2"),
		[]byte("succ"),
		[]byte("arg3"),
		[]byte("succ"),
	}

	// Alpha uses executeOnDestContext() to call beta, gamma and delta one after
	// the other, in the same transaction.
	input := DefaultTestContractCallInput()
	input.Function = "callChildrenDirectly_DestCtx"
	input.GasProvided = 1000000
	input.RecipientAddr = alpha.Address

	vmOutput, err := host.RunSmartContractCall(input)

	verify := NewVMOutputVerifier(t, vmOutput, err)
	verify.
		Ok().
		ReturnData(expectedReturnData...)
}

func TestExecution_ExecuteOnDestContextByCaller_SimpleTransfer(t *testing.T) {
	// The child contract is designed to send some tokens back to its caller, as
	// many as requested. The parent calls the child using
	// executeOnDestContextByCaller(), which means that the child will not see
	// the parent as its caller, but the original caller of the transaction
	// instead. Thus the original caller (the user address) will receive 42
	// tokens, and not the parent, even if the parent is the one making the call
	// to the child.

	transferValue := int64(42)

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCodeModule("exec-dest-ctx-by-caller/parent", "parent", "../../")).
				withBalance(1000),
			createInstanceContract(childAddress).
				withCode(GetTestSCCodeModule("exec-dest-ctx-by-caller/child", "child", "../../")).
				withBalance(1000),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction("call_child").
			withGasProvided(2000).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				// parentAddress
				GasUsed(parentAddress, 762).
				/// childAddress
				Balance(childAddress, 1000).
				BalanceDelta(childAddress, -transferValue).
				GasUsed(childAddress, 667).
				// userAddress
				BalanceDelta(userAddress, transferValue).
				// other
				ReturnData([]byte("sent"), []byte("child called")).
				Transfers(
					createTransferEntry(userAddress).
						withData([]byte{}).
						withValue(big.NewInt(transferValue)).
						withSenderAddress(childAddress),
				)
		})
}

func TestExecution_AsyncCall_GasLimitConsumed(t *testing.T) {

	parentCode := GetTestSCCode("async-call-parent", "../../")
	childCode := GetTestSCCode("async-call-child", "../../")

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(parentCode).
				withBalance(1000),
			createInstanceContract(childAddress).
				withCode(childCode).
				withBalance(1000),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction(parentPerformAsyncCall).
			withGasProvided(1000000).
			withArguments([]byte{0}).
			build()).
		withSetup(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.GetUserAccountCalled = func(scAddress []byte) (vmcommon.UserAccountHandler, error) {
				if bytes.Equal(scAddress, parentAddress) {
					return &contextmock.StubAccount{
						Address: parentAddress,
						Balance: big.NewInt(1000),
					}, nil
				}
				return nil, errAccountNotFound
			}
			stubBlockchainHook.GetCodeCalled = func(account vmcommon.UserAccountHandler) []byte {
				if bytes.Equal(parentAddress, account.AddressBytes()) {
					return parentCode
				}
				return nil
			}
			stubBlockchainHook.GetShardOfAddressCalled = func(address []byte) uint32 {
				if bytes.Equal(address, parentAddress) {
					return 0
				}
				return 1
			}
		}).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				GasRemaining(0)
		})
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

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using asyncCall().

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("async-call-parent", "../../")).
				withBalance(1000),
			createInstanceContract(childAddress).
				withCode(GetTestSCCode("async-call-child", "../../")).
				withBalance(1000),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction(parentPerformAsyncCall).
			withGasProvided(116000).
			withArguments([]byte{0}).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(parentAddress, 9114).
				GasUsed(childAddress, 2534).
				GasRemaining(104352).
				Balance(parentAddress, 1000).
				Balance(childAddress, 1000).
				BalanceDelta(thirdPartyAddress, 6).
				ReturnData(parentFinishA, parentFinishB, []byte{0}, []byte("thirdparty"), []byte("vault"), []byte{0}, []byte("succ")).
				Storage(
					createStoreEntry(parentAddress).withKey(parentKeyA).withValue(parentDataA),
					createStoreEntry(parentAddress).withKey(parentKeyB).withValue(parentDataB),
					createStoreEntry(childAddress).withKey(childKey).withValue(childData),
				).
				Transfers(
					createTransferEntry(thirdPartyAddress).
						withData([]byte("hello")).
						withValue(big.NewInt(3)).
						withSenderAddress(parentAddress),
					createTransferEntry(thirdPartyAddress).
						withData([]byte(" there")).
						withValue(big.NewInt(3)).
						withSenderAddress(childAddress),
					createTransferEntry(vaultAddress).
						withData([]byte{}).
						withValue(big.NewInt(4)).
						withSenderAddress(childAddress),
				)
		})
}

func TestExecution_AsyncCall_ChildFails(t *testing.T) {
	// Scenario
	// Identical to TestExecution_AsyncCall(), except that the child is
	// instructed to call signalError().
	// Because "vault" was not received by the callBack(), the Parent sends 4 ERD
	// to the Vault directly.

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using asyncCall().

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("async-call-parent", "../../")).
				withBalance(1000),
			createInstanceContract(childAddress).
				withCode(GetTestSCCode("async-call-child", "../../")).
				withBalance(1000),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction(parentPerformAsyncCall).
			withGasProvided(1000000).
			withArguments([]byte{1}).
			withCurrentTxHash([]byte("txhash")).
			build()).
		withSetup(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			host.Metering().GasSchedule().ElrondAPICost.AsyncCallbackGasLock = 3000
		}).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(parentAddress, 998352).
				GasRemaining(1648).
				ReturnData(parentFinishA, parentFinishB, []byte("succ")).
				Storage(
					createStoreEntry(parentAddress).withKey(parentKeyA).withValue(parentDataA),
					createStoreEntry(parentAddress).withKey(parentKeyB).withValue(parentDataB),
				)
		})
}

func TestExecution_AsyncCall_CallBackFails(t *testing.T) {
	// Scenario
	// Identical to TestExecution_AsyncCall(), except that the callback is
	// instructed to call signalError().

	// Call parentFunctionChildCall() of the parent SC, which will call the child
	// SC and pass some arguments using asyncCall().

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("async-call-parent", "../../")).
				withBalance(1000),
			createInstanceContract(childAddress).
				withCode(GetTestSCCode("async-call-child", "../../")).
				withBalance(1000),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction(parentPerformAsyncCall).
			withGasProvided(200000).
			withArguments([]byte{0, 3}).
			withCurrentTxHash([]byte("txhash")).
			build()).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				ReturnMessage("callBack error").
				GasUsed(parentAddress, 197437).
				GasUsed(childAddress, 2534).
				// TODO Why is there a minuscule amount of gas remaining after the callback
				// fails? This is supposed to be 0.
				GasRemaining(29).
				BalanceDelta(thirdPartyAddress, 6).
				BalanceDelta(childAddress, big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1)).Int64()).
				ReturnData(parentFinishA, parentFinishB, []byte{3}, []byte("thirdparty"), []byte("vault"), []byte("user error"), []byte("txhash")).
				Storage(
					createStoreEntry(parentAddress).withKey(parentKeyA).withValue(parentDataA),
					createStoreEntry(parentAddress).withKey(parentKeyB).withValue(parentDataB),
					createStoreEntry(childAddress).withKey(childKey).withValue(childData),
				).
				Transfers(
					createTransferEntry(thirdPartyAddress).
						withData([]byte("hello")).
						withValue(big.NewInt(3)).
						withSenderAddress(parentAddress),
					createTransferEntry(thirdPartyAddress).
						withData([]byte(" there")).
						withValue(big.NewInt(3)).
						withSenderAddress(childAddress),
					createTransferEntry(vaultAddress).
						withData([]byte{}).
						withValue(big.NewInt(4)).
						withSenderAddress(childAddress),
				)
		})
}

func TestExecution_CreateNewContract_Success(t *testing.T) {

	childCode := GetTestSCCode("init-correct", "../../")
	childAddress := []byte("newAddress")
	l := len(childCode)

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("deployer", "../../")).
				withBalance(1000),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction("deployChildContract").
			withGasProvided(1_000_000).
			withArguments([]byte{'A'}, []byte{0}).
			withCurrentTxHash([]byte("txhash")).
			build()).
		withSetup(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.GetStorageDataCalled = func(address []byte, key []byte) ([]byte, error) {
				if bytes.Equal(address, parentAddress) {
					if bytes.Equal(key, []byte{'A'}) {
						return childCode, nil
					}
					return nil, nil
				}
				return nil, arwen.ErrInvalidAccount
			}
		}).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				// parentAddress
				Balance(parentAddress, 1000).
				GasUsed(parentAddress, 885).
				/// childAddress
				BalanceDelta(childAddress, 42).
				Code(childAddress, childCode).
				CodeMetadata(childAddress, []byte{1, 0}).
				CodeDeployerAddress(childAddress, parentAddress).
				GasUsed(childAddress, 472).
				// other
				ReturnData([]byte{byte(l / 256), byte(l % 256)}, []byte("init successful"), []byte("succ")).
				Storage(
					createStoreEntry(parentAddress).withKey([]byte{'A'}).withValue(childCode))
		})
}

func TestExecution_CreateNewContract_Fail(t *testing.T) {

	childCode := GetTestSCCode("init-correct", "../../")
	l := len(childCode)

	runInstanceCallerTestBuilder(t).
		withContracts(
			createInstanceContract(parentAddress).
				withCode(GetTestSCCode("deployer", "../../")).
				withBalance(1000),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withFunction("deployChildContract").
			withGasProvided(1_000_000).
			withArguments([]byte{'A'}, []byte{1}).
			build()).
		withSetup(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub) {
			stubBlockchainHook.GetStorageDataCalled = func(address []byte, key []byte) ([]byte, error) {
				if bytes.Equal(address, parentAddress) {
					if bytes.Equal(key, []byte{'A'}) {
						return childCode, nil
					}
					return nil, nil
				}
				return nil, arwen.ErrInvalidAccount
			}
		}).
		andAssertResults(func(host *vmHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(parentAddress, 2885).
				ReturnData([]byte{byte(l / 256), byte(l % 256)}, []byte("fail")).
				Storage(createStoreEntry(parentAddress).withKey([]byte{'A'}).withValue(childCode))
		})
}

func TestExecution_Mocked_Wasmer_Instances(t *testing.T) {

	runMockInstanceCallerTestBuilder(t).
		withContracts(
			createMockContract(parentAddress).
				withBalance(1000).
				withMethods(func(parentInstance *mock.InstanceMock, config interface{}) {
					parentInstance.AddMockMethod("callChild", func() *mock.InstanceMock {
						host := parentInstance.Host
						host.Output().Finish([]byte("parent returns this"))
						host.Metering().UseGas(500)
						_, err := host.Storage().SetStorage([]byte("parent"), []byte("parent storage"))
						require.Nil(t, err)
						childInput := DefaultTestContractCallInput()
						childInput.CallerAddr = parentAddress
						childInput.RecipientAddr = childAddress
						childInput.CallValue = big.NewInt(4)
						childInput.Function = "doSomething"
						childInput.GasProvided = 1000
						_, _, err = host.ExecuteOnDestContext(childInput)
						require.Nil(t, err)
						return parentInstance
					})
				}),
			createMockContract(childAddress).
				withBalance(0).
				withMethods(func(childInstance *mock.InstanceMock, config interface{}) {
					childInstance.AddMockMethod("doSomething", func() *mock.InstanceMock {
						host := childInstance.Host
						host.Output().Finish([]byte("child returns this"))
						host.Metering().UseGas(100)
						_, err := host.Storage().SetStorage([]byte("child"), []byte("child storage"))
						require.Nil(t, err)
						return childInstance
					})
				}),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withGasProvided(1000).
			withFunction("callChild").
			build()).
		andAssertResults(func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			verify.
				Ok().
				// parentAddress
				Balance(parentAddress, 1000).
				BalanceDelta(parentAddress, -4).
				GasUsed(parentAddress, 547).
				/// childAddress
				BalanceDelta(childAddress, 4).
				GasUsed(childAddress, 146).
				// other
				GasRemaining(307).
				ReturnData([]byte("parent returns this"), []byte("child returns this")).
				Storage(
					createStoreEntry(parentAddress).withKey([]byte("parent")).withValue([]byte("parent storage")),
					createStoreEntry(childAddress).withKey([]byte("child")).withValue([]byte("child storage")),
				)
		})
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

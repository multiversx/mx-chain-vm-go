package host

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

var gasUsedByParent = uint64(400)
var gasUsedByChild = uint64(200)
var gasProvidedToChild = uint64(300)
var gasUsedByBuiltinClaim = uint64(120)

func TestGasUsed_SingleContract(t *testing.T) {
	host, _, imb := defaultTestArwenForCallWithInstanceMocks(t)
	createTestParentContract(t, host, imb)
	setZeroCodeCosts(host)

	gasProvided := uint64(1000)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.GasProvided = gasProvided
	input.Function = "wasteGas"

	vmOutput, err := host.RunSmartContractCall(input)

	verify := NewVMOutputVerifier(t, vmOutput, err)
	verify.Ok().GasRemaining(gasProvided - gasUsedByParent)
	verify.GasUsed(parentAddress, gasUsedByParent)
}

func TestGasUsed_SingleContract_BuiltinCall(t *testing.T) {
	host, world, imb := defaultTestArwenForCallWithInstanceMocks(t)
	createTestParentContract(t, host, imb)
	createMockBuiltinFunctions(t, host, world)
	setZeroCodeCosts(host)

	gasProvided := uint64(1000)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.GasProvided = gasProvided
	input.Function = "execOnDestCtx"
	input.Arguments = [][]byte{
		parentAddress,
		[]byte("builtinClaim"),
		arwen.One.Bytes(),
	}

	vmOutput, err := host.RunSmartContractCall(input)

	verify := NewVMOutputVerifier(t, vmOutput, err)
	verify.Ok().GasRemaining(gasProvided - gasUsedByParent - gasUsedByBuiltinClaim)
	verify.GasUsed(parentAddress, gasUsedByParent+gasUsedByBuiltinClaim)
}

func TestGasUsed_TwoContracts_ExecuteOnSameCtx(t *testing.T) {
	host, _, imb := defaultTestArwenForCallWithInstanceMocks(t)
	createTestParentContract(t, host, imb)
	createTestChildContract(t, host, imb)
	setZeroCodeCosts(host)

	gasProvided := uint64(1000)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.GasProvided = gasProvided
	input.Function = "execOnSameCtx"

	for numCalls := uint64(0); numCalls < 3; numCalls++ {
		expectedGasRemaining := gasProvided - gasUsedByParent - gasUsedByChild*numCalls

		numCallsBytes := big.NewInt(0).SetUint64(numCalls).Bytes()
		input.Arguments = [][]byte{
			childAddress,
			[]byte("wasteGas"),
			numCallsBytes,
		}

		vmOutput, err := host.RunSmartContractCall(input)

		verify := NewVMOutputVerifier(t, vmOutput, err)
		verify.Ok().GasRemaining(expectedGasRemaining)
		verify.GasUsed(parentAddress, gasUsedByParent)
		if numCalls > 0 {
			verify.GasUsed(childAddress, gasUsedByChild*numCalls)
		}
	}
}

func TestGasUsed_TwoContracts_ExecuteOnDestCtx(t *testing.T) {
	host, _, imb := defaultTestArwenForCallWithInstanceMocks(t)
	createTestParentContract(t, host, imb)
	createTestChildContract(t, host, imb)
	setZeroCodeCosts(host)

	gasProvided := uint64(1000)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.GasProvided = gasProvided
	input.Function = "execOnDestCtx"

	for numCalls := uint64(0); numCalls < 3; numCalls++ {
		expectedGasRemaining := gasProvided - gasUsedByParent - gasUsedByChild*numCalls

		numCallsBytes := big.NewInt(0).SetUint64(numCalls).Bytes()
		input.Arguments = [][]byte{
			childAddress,
			[]byte("wasteGas"),
			numCallsBytes,
		}

		vmOutput, err := host.RunSmartContractCall(input)

		verify := NewVMOutputVerifier(t, vmOutput, err)
		verify.Ok().GasRemaining(expectedGasRemaining)
		verify.GasUsed(parentAddress, gasUsedByParent)
		if numCalls > 0 {
			verify.GasUsed(childAddress, gasUsedByChild*numCalls)
		}
	}
}

func TestGasUsed_ThreeContracts_ExecuteOnDestCtx(t *testing.T) {
	host, _, imb := defaultTestArwenForCallWithInstanceMocks(t)
	setZeroCodeCosts(host)

	alphaAddress := MakeTestSCAddress("alpha")
	betaAddress := MakeTestSCAddress("beta")
	gammaAddress := MakeTestSCAddress("gamma")

	gasProvided := uint64(1000)
	alphaCallGas := uint64(400)
	alphaGasToForwardToReceivers := uint64(300)
	receiverCallGas := uint64(200)

	expectedGasRemaining := gasProvided - alphaCallGas - 2*receiverCallGas

	alpha := imb.CreateAndStoreInstanceMock(t, host, alphaAddress, 0)
	addForwarderMethodsToInstanceMock(alpha, alphaCallGas, alphaGasToForwardToReceivers)

	beta := imb.CreateAndStoreInstanceMock(t, host, betaAddress, 0)
	gamma := imb.CreateAndStoreInstanceMock(t, host, gammaAddress, 0)
	addDummyMethodsToInstanceMock(beta, receiverCallGas)
	addDummyMethodsToInstanceMock(gamma, receiverCallGas)

	input := DefaultTestContractCallInput()
	input.RecipientAddr = alphaAddress
	input.GasProvided = gasProvided
	input.Function = "execOnDestCtx"
	input.Arguments = [][]byte{
		betaAddress,
		[]byte("wasteGas"),
		arwen.One.Bytes(),
		gammaAddress,
		[]byte("wasteGas"),
		arwen.One.Bytes(),
	}

	vmOutput, err := host.RunSmartContractCall(input)

	verify := NewVMOutputVerifier(t, vmOutput, err)
	verify.GasUsed(alphaAddress, alphaCallGas)
	verify.GasUsed(betaAddress, receiverCallGas)
	verify.GasUsed(gammaAddress, receiverCallGas)
	verify.Ok().GasRemaining(expectedGasRemaining)
}

type asyncCallBaseTestConfig struct {
	gasProvided        uint64
	gasUsedByParent    uint64
	gasProvidedToChild uint64
	gasUsedByChild     uint64
	gasUsedByCallback  uint64
	gasLockCost        uint64

	transferFromParentToChild int64

	parentBalance int64
	childBalance  int64
}

var asyncBaseTestConfig = asyncCallBaseTestConfig{
	gasProvided:        116000,
	gasUsedByParent:    400,
	gasProvidedToChild: 300,
	gasUsedByChild:     200,
	gasUsedByCallback:  100,
	gasLockCost:        150,

	transferFromParentToChild: 7,

	parentBalance: 1000,
	childBalance:  1000,
}

type asyncCallTestConfig struct {
	asyncCallBaseTestConfig
	transferToThirdParty int64
	transferToVault      int64
}

var asyncTestConfig = &asyncCallTestConfig{
	asyncCallBaseTestConfig: asyncBaseTestConfig,
	transferToThirdParty:    3,
	transferToVault:         4,
}

func TestGasUsed_AsyncCall(t *testing.T) {

	testConfig := asyncTestConfig

	parentMockContract := &mockSmartContract{
		address:     parentAddress,
		balance:     testConfig.parentBalance,
		config:      testConfig,
		initMethods: &initFunctions{addAsyncParentMethodsToInstanceMock},
	}

	childMockContract := &mockSmartContract{
		address:     childAddress,
		balance:     testConfig.childBalance,
		config:      testConfig,
		initMethods: &initFunctions{addAsyncChildMethodsToInstanceMock},
	}

	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "performAsyncCall"
	input.GasProvided = testConfig.gasProvided
	input.Arguments = [][]byte{{0}}

	runMockInstanceCallerTest(&mockInstancesTestTemplate{
		t:         t,
		contracts: &mockContracts{parentMockContract, childMockContract},
		input:     input,
		setup: func(host *vmHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.gasLockCost)
		},
		assertResults: func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			// TODO update verifier to allow other stuff in expectedVMOutputAsyncCallWithConfig(testConfig)
			// + delete function + refactor old one
			verify.RetCode(vmcommon.Ok)
			gasUsedByParent := testConfig.gasUsedByParent + testConfig.gasUsedByCallback
			verify.GasUsed(parentAddress, gasUsedByParent)
			gasUsedByChild := testConfig.gasUsedByChild
			verify.GasUsed(childAddress, gasUsedByChild)
			verify.Ok().GasRemaining(testConfig.gasProvided - gasUsedByParent - gasUsedByChild)
		},
	})
}

func TestGasUsed_AsyncCall_BuiltinCall(t *testing.T) {

	testConfig := asyncBaseTestConfig
	testConfig.gasProvided = 1000

	parentMockContract := &mockSmartContract{
		address:     parentAddress,
		balance:     asyncBaseTestConfig.parentBalance,
		config:      &testConfig,
		initMethods: &initFunctions{addAsyncBuiltinParentMethodsToInstanceMock},
	}

	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.GasProvided = testConfig.gasProvided
	input.Function = "forwardAsyncCall"
	input.Arguments = [][]byte{
		userAddress,
		[]byte("builtinClaim"),
	}

	runMockInstanceCallerTest(&mockInstancesTestTemplate{
		t:         t,
		contracts: &mockContracts{parentMockContract},
		input:     input,
		setup: func(host *vmHost, world *worldmock.MockWorld) {
			world.AcctMap.CreateAccount(userAddress)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, asyncBaseTestConfig.gasLockCost)
		},
		assertResults: func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			verify.RetCode(vmcommon.Ok)
			expectedGasUsedByParent := asyncBaseTestConfig.gasUsedByParent + asyncBaseTestConfig.gasUsedByCallback + gasUsedByBuiltinClaim
			verify.GasUsed(parentAddress, expectedGasUsedByParent)
			expectedGasUsedByChild := uint64(0) // all gas for builtin call is consummed on caller
			verify.GasUsed(userAddress, asyncBaseTestConfig.gasUsedByChild)
			verify.Ok().GasRemaining(asyncBaseTestConfig.gasProvided - expectedGasUsedByParent - expectedGasUsedByChild)
		},
	})
}

type asyncBuiltInCallTestConfig struct {
	asyncCallBaseTestConfig
	transferFromChildToParent int64
}

func TestGasUsed_AsyncCall_BuiltinMultiContractCall(t *testing.T) {

	// TODO no possible yet, reactivate when new async context is on
	t.Skip()

	testConfig := &asyncBuiltInCallTestConfig{
		asyncCallBaseTestConfig:   asyncBaseTestConfig,
		transferFromChildToParent: 5,
	}

	parentMockContract := &mockSmartContract{
		address:     parentAddress,
		balance:     testConfig.parentBalance,
		config:      testConfig,
		initMethods: &initFunctions{addAsyncBuiltinMultiContractParentMethodsToInstanceMock},
	}

	childMockContract := &mockSmartContract{
		address:     childAddress,
		balance:     testConfig.childBalance,
		config:      testConfig,
		initMethods: &initFunctions{addAsyncBuiltinMultiContractChildMethodsToInstanceMock},
	}

	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.GasProvided = testConfig.gasProvided
	input.Function = "forwardAsyncCall"
	input.Arguments = [][]byte{
		userAddress,
		childAddress,
		[]byte("childFunction"),
		[]byte("builtinClaim"),
	}

	runMockInstanceCallerTest(&mockInstancesTestTemplate{
		t:         t,
		contracts: &mockContracts{parentMockContract, childMockContract},
		input:     input,
		setup: func(host *vmHost, world *worldmock.MockWorld) {
			world.AcctMap.CreateAccount(userAddress)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.gasLockCost)
			createMockBuiltinFunctions(t, host, world)
		},
		assertResults: func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			verify.RetCode(vmcommon.Ok)
			expectedGasUsedByParent := testConfig.gasUsedByParent + testConfig.gasUsedByCallback
			verify.GasUsed(parentAddress, expectedGasUsedByParent)
			expectedGasUsedByChild := testConfig.gasUsedByChild + gasUsedByBuiltinClaim
			verify.GasUsed(childAddress, testConfig.gasUsedByChild)
			verify.Ok().GasRemaining(testConfig.gasProvided - expectedGasUsedByParent - expectedGasUsedByChild)
		},
	})
}

func TestGasUsed_AsyncCall_ChildFails(t *testing.T) {

	testConfig := asyncTestConfig

	parentMockContract := &mockSmartContract{
		address:     parentAddress,
		balance:     testConfig.parentBalance,
		config:      testConfig,
		initMethods: &initFunctions{addAsyncParentMethodsToInstanceMock},
	}

	childMockContract := &mockSmartContract{
		address:     childAddress,
		balance:     testConfig.childBalance,
		config:      testConfig,
		initMethods: &initFunctions{addAsyncChildMethodsToInstanceMock},
	}

	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "performAsyncCall"
	input.GasProvided = testConfig.gasProvided
	input.Arguments = [][]byte{{1}}

	runMockInstanceCallerTest(&mockInstancesTestTemplate{
		t:         t,
		contracts: &mockContracts{parentMockContract, childMockContract},
		input:     input,
		setup: func(host *vmHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.gasLockCost)
		},
		assertResults: func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			// TODO complete with expectedVMOutputAsyncCallChildFailsWithConfig()
			// + delete function + refactor old one
			verify.RetCode(vmcommon.Ok)
			expectedGasUsedByParent := testConfig.gasProvided - testConfig.gasLockCost + testConfig.gasUsedByCallback
			verify.GasUsed(parentAddress, expectedGasUsedByParent)
			verify.Ok().GasRemaining(testConfig.gasProvided - expectedGasUsedByParent)
		},
	})
}

func TestGasUsed_AsyncCall_CallBackFails(t *testing.T) {

	testConfig := asyncTestConfig

	parentMockContract := &mockSmartContract{
		address:     parentAddress,
		balance:     testConfig.parentBalance,
		config:      testConfig,
		initMethods: &initFunctions{addAsyncParentMethodsToInstanceMock},
	}

	childMockContract := &mockSmartContract{
		address:     childAddress,
		balance:     testConfig.childBalance,
		config:      testConfig,
		initMethods: &initFunctions{addAsyncChildMethodsToInstanceMock},
	}

	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "performAsyncCall"
	input.GasProvided = testConfig.gasProvided
	input.Arguments = [][]byte{{0, 3}}

	runMockInstanceCallerTest(&mockInstancesTestTemplate{
		t:         t,
		contracts: &mockContracts{parentMockContract, childMockContract},
		input:     input,
		setup: func(host *vmHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.gasLockCost)
		},
		assertResults: func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			// TODO complete with expectedVMOutputAsyncCallCallBackFailsWithConfig()
			// + delete function + refactor old one
			verify.RetCode(vmcommon.Ok)
			expectedGasUsedByParent := testConfig.gasProvided - testConfig.gasUsedByChild
			verify.GasUsed(parentAddress, expectedGasUsedByParent)
			expectedGasUsedByChild := testConfig.gasUsedByChild
			verify.GasUsed(childAddress, expectedGasUsedByChild)
			verify.Ok().GasRemaining(0)
		},
	})
}

type asyncCallRecursiveTestConfig struct {
	asyncCallBaseTestConfig
	recursiveChildCalls int
}

func TestGasUsed_AsyncCall_Recursive(t *testing.T) {

	// // TODO no possible yet, reactivate when new async context is on
	t.Skip()

	testConfig := &asyncCallRecursiveTestConfig{
		asyncCallBaseTestConfig: *&asyncBaseTestConfig,
		recursiveChildCalls:     2,
	}

	parentMockContract := &mockSmartContract{
		address:     parentAddress,
		balance:     testConfig.parentBalance,
		config:      testConfig,
		initMethods: &initFunctions{addAsyncRecursiveParentMethodsToInstanceMock},
	}

	childMockContract := &mockSmartContract{
		address:     childAddress,
		balance:     testConfig.childBalance,
		config:      &testConfig.asyncCallBaseTestConfig,
		initMethods: &initFunctions{addAsyncRecursiveChildMethodsToInstanceMock},
	}

	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "forwardAsyncCall"
	input.GasProvided = testConfig.gasProvided
	input.Arguments = [][]byte{childAddress, []byte("recursiveAsyncCall"), big.NewInt(int64(testConfig.recursiveChildCalls)).Bytes()}

	runMockInstanceCallerTest(&mockInstancesTestTemplate{
		t:         t,
		contracts: &mockContracts{parentMockContract, childMockContract},
		input:     input,
		setup: func(host *vmHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.gasLockCost)
		},
		assertResults: func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			// TODO complete with expectedVMOutputAsyncRecursiveCallWithConfig()
			// + delete function + refactor old one
			verify.RetCode(vmcommon.Ok)
			expectedGasUsedByParent := testConfig.gasUsedByParent + testConfig.gasUsedByCallback
			verify.GasUsed(parentAddress, expectedGasUsedByParent)
			expectedGasUsedByChild := uint64(testConfig.recursiveChildCalls)*testConfig.gasUsedByChild +
				uint64(testConfig.recursiveChildCalls-1)*testConfig.gasUsedByCallback
			verify.GasUsed(childAddress, expectedGasUsedByChild)
			verify.Ok().GasRemaining(testConfig.gasProvided - expectedGasUsedByParent - expectedGasUsedByChild)
		},
	})
}

type asyncCallMultiChildTestConfig struct {
	asyncCallBaseTestConfig
	childCalls int
}

func TestGasUsed_AsyncCall_MultiChild(t *testing.T) {

	// TODO no possible yet, reactivate when new async context is on
	t.Skip()

	testConfig := &asyncCallMultiChildTestConfig{
		asyncCallBaseTestConfig: *&asyncBaseTestConfig,
		childCalls:              2,
	}

	parentMockContract := &mockSmartContract{
		address:     parentAddress,
		balance:     testConfig.parentBalance,
		config:      testConfig,
		initMethods: &initFunctions{addAsyncMultiChildParentMethodsToInstanceMock},
	}

	childMockContract := &mockSmartContract{
		address:     childAddress,
		balance:     testConfig.childBalance,
		config:      &testConfig.asyncCallBaseTestConfig,
		initMethods: &initFunctions{addAsyncRecursiveChildMethodsToInstanceMock},
	}

	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.Function = "forwardAsyncCall"
	input.GasProvided = testConfig.gasProvided
	input.Arguments = [][]byte{childAddress, []byte("recursiveAsyncCall")}

	runMockInstanceCallerTest(&mockInstancesTestTemplate{
		t:         t,
		contracts: &mockContracts{parentMockContract, childMockContract},
		input:     input,
		setup: func(host *vmHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.gasLockCost)
		},
		assertResults: func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			// TODO complete with expectedVMOutputAsyncMultiChildCallWithConfig()
			// + delete function + refactor old one
			verify.RetCode(vmcommon.Ok)
			expectedGasUsedByParent := testConfig.gasUsedByParent + testConfig.gasUsedByCallback
			verify.GasUsed(parentAddress, expectedGasUsedByParent)
			expectedGasUsedByChild := uint64(testConfig.childCalls) * testConfig.gasUsedByChild
			verify.GasUsed(childAddress, expectedGasUsedByChild)
			verify.Ok().GasRemaining(testConfig.gasProvided - expectedGasUsedByParent - expectedGasUsedByChild)
		},
	})
}

func createTestParentContract(t testing.TB, host *vmHost, imb *mock.InstanceBuilderMock) {
	gasUsedByParent := uint64(400)
	gasProvidedToChild := uint64(300)

	parentInstance := imb.CreateAndStoreInstanceMock(t, host, parentAddress, 1000)
	addDummyMethodsToInstanceMock(parentInstance, gasUsedByParent)
	addForwarderMethodsToInstanceMock(parentInstance, gasUsedByParent, gasProvidedToChild)
}

func createTestChildContract(t testing.TB, host *vmHost, imb *mock.InstanceBuilderMock) {
	gasUsedByChild := uint64(200)

	childInstance := imb.CreateAndStoreInstanceMock(t, host, childAddress, 0)
	addDummyMethodsToInstanceMock(childInstance, gasUsedByChild)
}

func createMockBuiltinFunctions(tb testing.TB, host *vmHost, world *worldmock.MockWorld) {
	err := world.InitBuiltinFunctions(host.GetGasScheduleMap())
	require.Nil(tb, err)

	world.BuiltinFuncs.Container.Add("builtinClaim", &MockClaimBuiltin{
		AmountToGive: 42,
		GasCost:      gasUsedByBuiltinClaim,
	})

	host.protocolBuiltinFunctions = world.BuiltinFuncs.GetBuiltinFunctionNames()
}

func addDummyMethodsToInstanceMock(instanceMock *mock.InstanceMock, gasPerCall uint64) {
	instanceMock.AddMockMethod("wasteGas", func() *mock.InstanceMock {
		host := instanceMock.Host
		host.Metering().UseGas(gasPerCall)
		instance := mock.GetMockInstance(host)
		return instance
	})
}

func addForwarderMethodsToInstanceMock(instanceMock *mock.InstanceMock, gasPerCall uint64, gasToForward uint64) {
	input := DefaultTestContractCallInput()
	input.GasProvided = gasToForward

	instanceMock.AddMockMethod("execOnSameCtx", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		host.Metering().UseGas(gasPerCall)

		arguments := host.Runtime().Arguments()
		input.CallerAddr = instance.Address
		input.RecipientAddr = arguments[0]
		input.Function = string(arguments[1])
		numCalls := big.NewInt(0).SetBytes(arguments[2]).Uint64()

		for i := uint64(0); i < numCalls; i++ {
			_, err := host.ExecuteOnSameContext(input)
			require.Nil(t, err)
		}

		return instance
	})

	instanceMock.AddMockMethod("execOnDestCtx", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		host.Metering().UseGas(gasPerCall)

		argsPerCall := 3
		arguments := host.Runtime().Arguments()
		if len(arguments)%argsPerCall != 0 {
			host.Runtime().SignalUserError("need 3 arguments per individual call")
			return instance
		}

		input.CallerAddr = instance.Address

		for callIndex := 0; callIndex < len(arguments); callIndex += argsPerCall {
			input.RecipientAddr = arguments[callIndex+0]
			input.Function = string(arguments[callIndex+1])
			numCalls := big.NewInt(0).SetBytes(arguments[callIndex+2]).Uint64()

			for i := uint64(0); i < numCalls; i++ {
				_, _, err := host.ExecuteOnDestContext(input)
				require.Nil(t, err)
			}
		}

		return instance
	})
}

func setZeroCodeCosts(host *vmHost) {
	host.Metering().GasSchedule().BaseOperationCost.CompilePerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.AoTPreparePerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.GetCode = 0
	host.Metering().GasSchedule().BaseOperationCost.StorePerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.DataCopyPerByte = 0
	host.Metering().GasSchedule().ElrondAPICost.SignalError = 0
}

func setAsyncCosts(host *vmHost, gasLock uint64) {
	host.Metering().GasSchedule().ElrondAPICost.AsyncCallStep = 0
	host.Metering().GasSchedule().ElrondAPICost.AsyncCallbackGasLock = gasLock
}

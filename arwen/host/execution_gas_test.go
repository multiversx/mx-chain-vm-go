package host

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

var gasUsedByBuiltinClaim = uint64(120)

type directCallGasTestConfig struct {
	gasUsedByParent    uint64
	gasUsedByChild     uint64
	gasProvidedToChild uint64
	gasProvided        uint64
	parentBalance      int64
	childBalance       int64
}

var simpleGasTestConfig = directCallGasTestConfig{
	gasUsedByParent:    uint64(400),
	gasUsedByChild:     uint64(200),
	gasProvidedToChild: uint64(300),
	gasProvided:        uint64(1000),
	parentBalance:      int64(1000),
	childBalance:       int64(1000),
}

func TestGasUsed_SingleContract(t *testing.T) {

	parentMockContract := &mockSmartContract{
		address:     parentAddress,
		balance:     simpleGasTestConfig.parentBalance,
		config:      simpleGasTestConfig,
		initMethods: &initFunctions{addMethodsToParentInstanceMock},
	}

	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.GasProvided = simpleGasTestConfig.gasProvided
	input.Function = "wasteGas"

	runMockInstanceCallerTest(&mockInstancesTestTemplate{
		t:         t,
		contracts: &mockContracts{parentMockContract},
		input:     input,
		setup: func(host *vmHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
		},
		assertResults: func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			verify.
				Ok().
				GasRemaining(simpleGasTestConfig.gasProvided-simpleGasTestConfig.gasUsedByParent).
				GasUsed(parentAddress, simpleGasTestConfig.gasUsedByParent)
		},
	})
}

func TestGasUsed_SingleContract_BuiltinCall(t *testing.T) {

	parentMockContract := &mockSmartContract{
		address:     parentAddress,
		balance:     simpleGasTestConfig.parentBalance,
		config:      simpleGasTestConfig,
		initMethods: &initFunctions{addMethodsToParentInstanceMock},
	}

	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.GasProvided = simpleGasTestConfig.gasProvided
	input.Function = "execOnDestCtx"
	input.Arguments = [][]byte{
		parentAddress,
		[]byte("builtinClaim"),
		arwen.One.Bytes(),
	}

	runMockInstanceCallerTest(&mockInstancesTestTemplate{
		t:         t,
		contracts: &mockContracts{parentMockContract},
		input:     input,
		setup: func(host *vmHost, world *worldmock.MockWorld) {
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		},
		assertResults: func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			verify.
				Ok().
				GasRemaining(simpleGasTestConfig.gasProvided-simpleGasTestConfig.gasUsedByParent-gasUsedByBuiltinClaim).
				GasUsed(parentAddress, simpleGasTestConfig.gasUsedByParent+gasUsedByBuiltinClaim)
		},
	})
}

func TestGasUsed_TwoContracts_ExecuteOnSameCtx(t *testing.T) {

	parentMockContract := &mockSmartContract{
		address:     parentAddress,
		balance:     simpleGasTestConfig.parentBalance,
		config:      simpleGasTestConfig,
		initMethods: &initFunctions{addMethodsToParentInstanceMock},
	}

	childMockContract := &mockSmartContract{
		address:     childAddress,
		balance:     simpleGasTestConfig.childBalance,
		config:      simpleGasTestConfig,
		initMethods: &initFunctions{addMethodsToChildInstanceMock},
	}

	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.GasProvided = simpleGasTestConfig.gasProvided
	input.Function = "execOnSameCtx"

	for numCalls := uint64(0); numCalls < 3; numCalls++ {
		expectedGasRemaining := simpleGasTestConfig.gasProvided - simpleGasTestConfig.gasUsedByParent - simpleGasTestConfig.gasUsedByChild*numCalls

		numCallsBytes := big.NewInt(0).SetUint64(numCalls).Bytes()
		input.Arguments = [][]byte{
			childAddress,
			[]byte("wasteGas"),
			numCallsBytes,
		}

		runMockInstanceCallerTest(&mockInstancesTestTemplate{
			t:         t,
			contracts: &mockContracts{parentMockContract, childMockContract},
			input:     input,
			setup: func(host *vmHost, world *worldmock.MockWorld) {
				setZeroCodeCosts(host)
			},
			assertResults: func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
				verify.
					Ok().
					GasRemaining(expectedGasRemaining).
					GasUsed(parentAddress, simpleGasTestConfig.gasUsedByParent)
				if numCalls > 0 {
					verify.GasUsed(childAddress, simpleGasTestConfig.gasUsedByChild*numCalls)
				}
			},
		})
	}
}

func TestGasUsed_TwoContracts_ExecuteOnDestCtx(t *testing.T) {

	parentMockContract := &mockSmartContract{
		address:     parentAddress,
		balance:     simpleGasTestConfig.parentBalance,
		config:      simpleGasTestConfig,
		initMethods: &initFunctions{addMethodsToParentInstanceMock},
	}

	childMockContract := &mockSmartContract{
		address:     childAddress,
		balance:     simpleGasTestConfig.childBalance,
		config:      simpleGasTestConfig,
		initMethods: &initFunctions{addMethodsToChildInstanceMock},
	}

	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.GasProvided = simpleGasTestConfig.gasProvided
	input.Function = "execOnDestCtx"

	for numCalls := uint64(0); numCalls < 3; numCalls++ {
		expectedGasRemaining := simpleGasTestConfig.gasProvided - simpleGasTestConfig.gasUsedByParent - simpleGasTestConfig.gasUsedByChild*numCalls

		numCallsBytes := big.NewInt(0).SetUint64(numCalls).Bytes()
		input.Arguments = [][]byte{
			childAddress,
			[]byte("wasteGas"),
			numCallsBytes,
		}

		runMockInstanceCallerTest(&mockInstancesTestTemplate{
			t:         t,
			contracts: &mockContracts{parentMockContract, childMockContract},
			input:     input,
			setup: func(host *vmHost, world *worldmock.MockWorld) {
				setZeroCodeCosts(host)
			},
			assertResults: func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
				verify.
					Ok().
					GasRemaining(expectedGasRemaining).
					GasUsed(parentAddress, simpleGasTestConfig.gasUsedByParent)
				if numCalls > 0 {
					verify.GasUsed(childAddress, simpleGasTestConfig.gasUsedByChild*numCalls)
				}
			},
		})
	}
}

func TestGasUsed_ThreeContracts_ExecuteOnDestCtx(t *testing.T) {

	alphaAddress := MakeTestSCAddress("alpha")
	betaAddress := MakeTestSCAddress("beta")
	gammaAddress := MakeTestSCAddress("gamma")

	gasProvided := uint64(1000)
	alphaCallGas := uint64(400)
	alphaGasToForwardToReceivers := uint64(300)
	receiverCallGas := uint64(200)

	alphaMockContract := &mockSmartContract{
		address: alphaAddress,
		balance: 0,
		config: directCallGasTestConfig{
			gasUsedByParent:    alphaCallGas,
			gasProvidedToChild: alphaGasToForwardToReceivers,
			gasProvided:        gasProvided,
		},
		initMethods: &initFunctions{addMethodsToParentInstanceMock},
	}

	betaMockContract := &mockSmartContract{
		address: betaAddress,
		balance: 0,
		config: directCallGasTestConfig{
			gasUsedByChild: receiverCallGas,
		},
		initMethods: &initFunctions{addMethodsToChildInstanceMock},
	}

	gammaMockContract := &mockSmartContract{
		address: gammaAddress,
		balance: 0,
		config: directCallGasTestConfig{
			gasUsedByChild: receiverCallGas,
		},
		initMethods: &initFunctions{addMethodsToChildInstanceMock},
	}

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

	runMockInstanceCallerTest(&mockInstancesTestTemplate{
		t:         t,
		contracts: &mockContracts{alphaMockContract, betaMockContract, gammaMockContract},
		input:     input,
		setup: func(host *vmHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
		},
		assertResults: func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(alphaAddress, alphaCallGas).
				GasUsed(betaAddress, receiverCallGas).
				GasUsed(gammaAddress, receiverCallGas).
				GasRemaining(gasProvided - alphaCallGas - 2*receiverCallGas)
		},
	})
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
			gasUsedByParent := testConfig.gasUsedByParent + testConfig.gasUsedByCallback
			gasUsedByChild := testConfig.gasUsedByChild

			verify.
				Ok().
				GasUsed(parentAddress, gasUsedByParent).
				GasUsed(childAddress, gasUsedByChild).
				GasRemaining(testConfig.gasProvided-gasUsedByParent-gasUsedByChild).
				BalanceDelta(thirdPartyAddress, 2*testConfig.transferToThirdParty).
				ReturnData(parentFinishA, parentFinishB, []byte{0}, []byte("thirdparty"), []byte("vault"), []byte{0}, []byte("succ")).
				Storage(
					storeEntry{parentAddress, parentKeyA, parentDataA},
					storeEntry{parentAddress, parentKeyB, parentDataB},
					storeEntry{childAddress, childKey, childData},
				).
				Transfers(
					transferEntry{thirdPartyAddress,
						vmcommon.OutputTransfer{Data: []byte("hello"), Value: big.NewInt(testConfig.transferToThirdParty), SenderAddress: parentAddress}},
					transferEntry{thirdPartyAddress,
						vmcommon.OutputTransfer{Data: []byte(" there"), Value: big.NewInt(testConfig.transferToThirdParty), SenderAddress: childAddress}},
					transferEntry{vaultAddress,
						vmcommon.OutputTransfer{Data: []byte{}, Value: big.NewInt(testConfig.transferToVault), SenderAddress: childAddress}},
				)
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
			expectedGasUsedByParent := asyncBaseTestConfig.gasUsedByParent + asyncBaseTestConfig.gasUsedByCallback + gasUsedByBuiltinClaim
			expectedGasUsedByChild := uint64(0) // all gas for builtin call is consummed on caller
			verify.
				Ok().
				GasUsed(parentAddress, expectedGasUsedByParent).
				GasUsed(userAddress, asyncBaseTestConfig.gasUsedByChild).
				GasRemaining(asyncBaseTestConfig.gasProvided - expectedGasUsedByParent - expectedGasUsedByChild)
		},
	})
}

type asyncBuiltInCallTestConfig struct {
	asyncCallBaseTestConfig
	transferFromChildToParent int64
}

func TestGasUsed_AsyncCall_BuiltinMultiContractCall(t *testing.T) {

	// TODO no possible yet, reactivate when new async context is merged
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
			expectedGasUsedByParent := testConfig.gasUsedByParent + testConfig.gasUsedByCallback
			expectedGasUsedByChild := testConfig.gasUsedByChild + gasUsedByBuiltinClaim
			verify.
				Ok().
				GasUsed(parentAddress, expectedGasUsedByParent).
				GasUsed(childAddress, testConfig.gasUsedByChild).
				GasRemaining(testConfig.gasProvided - expectedGasUsedByParent - expectedGasUsedByChild)
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
	input.CurrentTxHash = []byte("txhash")

	runMockInstanceCallerTest(&mockInstancesTestTemplate{
		t:         t,
		contracts: &mockContracts{parentMockContract, childMockContract},
		input:     input,
		setup: func(host *vmHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.gasLockCost)
		},
		assertResults: func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			expectedGasUsedByParent := testConfig.gasProvided - testConfig.gasLockCost + testConfig.gasUsedByCallback

			verify.
				Ok().
				GasUsed(parentAddress, expectedGasUsedByParent).
				GasRemaining(testConfig.gasProvided-expectedGasUsedByParent).
				ReturnData(parentFinishA, parentFinishB, []byte("succ")).
				Storage(
					storeEntry{parentAddress, parentKeyA, parentDataA},
					storeEntry{parentAddress, parentKeyB, parentDataB},
				)
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
	input.CurrentTxHash = []byte("txhash")

	runMockInstanceCallerTest(&mockInstancesTestTemplate{
		t:         t,
		contracts: &mockContracts{parentMockContract, childMockContract},
		input:     input,
		setup: func(host *vmHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.gasLockCost)
		},
		assertResults: func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			expectedGasUsedByParent := testConfig.gasProvided - testConfig.gasUsedByChild
			expectedGasUsedByChild := testConfig.gasUsedByChild

			verify.
				Ok().
				ReturnMessage("callBack error").
				GasUsed(parentAddress, expectedGasUsedByParent).
				GasUsed(childAddress, expectedGasUsedByChild).
				GasRemaining(0).
				ReturnData(parentFinishA, parentFinishB, []byte{3}, []byte("thirdparty"), []byte("vault"), []byte("user error"), []byte("txhash")).
				Storage(
					storeEntry{parentAddress, parentKeyA, parentDataA},
					storeEntry{parentAddress, parentKeyB, parentDataB},
					storeEntry{childAddress, childKey, childData},
				).
				Transfers(
					transferEntry{thirdPartyAddress,
						vmcommon.OutputTransfer{Data: []byte("hello"), Value: big.NewInt(testConfig.transferToThirdParty), SenderAddress: parentAddress}},
					transferEntry{thirdPartyAddress,
						vmcommon.OutputTransfer{Data: []byte(" there"), Value: big.NewInt(testConfig.transferToThirdParty), SenderAddress: childAddress}},
					transferEntry{vaultAddress,
						vmcommon.OutputTransfer{Data: []byte{}, Value: big.NewInt(testConfig.transferToVault), SenderAddress: childAddress}},
				)
		},
	})
}

type asyncCallRecursiveTestConfig struct {
	asyncCallBaseTestConfig
	recursiveChildCalls int
}

func TestGasUsed_AsyncCall_Recursive(t *testing.T) {

	// TODO no possible yet, reactivate when new async context is merged
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
			expectedGasUsedByParent := testConfig.gasUsedByParent + testConfig.gasUsedByCallback
			expectedGasUsedByChild := uint64(testConfig.recursiveChildCalls)*testConfig.gasUsedByChild +
				uint64(testConfig.recursiveChildCalls-1)*testConfig.gasUsedByCallback
			verify.
				Ok().
				GasUsed(parentAddress, expectedGasUsedByParent).
				GasUsed(childAddress, expectedGasUsedByChild).
				GasRemaining(testConfig.gasProvided-expectedGasUsedByParent-expectedGasUsedByChild).
				BalanceDelta(childAddress, testConfig.transferFromParentToChild)
		},
	})
}

type asyncCallMultiChildTestConfig struct {
	asyncCallBaseTestConfig
	childCalls int
}

func TestGasUsed_AsyncCall_MultiChild(t *testing.T) {

	// TODO no possible yet, reactivate when new async context is merged
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
			expectedGasUsedByParent := testConfig.gasUsedByParent + testConfig.gasUsedByCallback
			expectedGasUsedByChild := uint64(testConfig.childCalls) * testConfig.gasUsedByChild
			verify.
				Ok().
				GasUsed(parentAddress, expectedGasUsedByParent).
				GasUsed(childAddress, expectedGasUsedByChild).
				GasRemaining(testConfig.gasProvided-expectedGasUsedByParent-expectedGasUsedByChild).
				BalanceDelta(childAddress, testConfig.transferFromParentToChild)
		},
	})
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

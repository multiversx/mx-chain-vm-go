package host

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
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

	runMockInstanceCallerTestBuilder(t).
		withContracts(
			createMockContract(parentAddress).
				withBalance(simpleGasTestConfig.parentBalance).
				withConfig(simpleGasTestConfig).
				withMethods(addMethodsToParentInstanceMock)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withGasProvided(simpleGasTestConfig.gasProvided).
			withFunction("wasteGas").
			build()).
		withSetup(func(host *vmHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
		}).
		andAssertResults(func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			verify.
				Ok().
				GasRemaining(simpleGasTestConfig.gasProvided-simpleGasTestConfig.gasUsedByParent).
				GasUsed(parentAddress, simpleGasTestConfig.gasUsedByParent)
		})
}

func TestGasUsed_SingleContract_BuiltinCall(t *testing.T) {
	runMockInstanceCallerTestBuilder(t).
		withContracts(
			createMockContract(parentAddress).
				withBalance(simpleGasTestConfig.parentBalance).
				withConfig(simpleGasTestConfig).
				withMethods(addMethodsToParentInstanceMock)).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withGasProvided(simpleGasTestConfig.gasProvided).
			withFunction("execOnDestCtx").
			withArguments(parentAddress, []byte("builtinClaim"), arwen.One.Bytes()).
			build()).
		withSetup(func(host *vmHost, world *worldmock.MockWorld) {
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		}).
		andAssertResults(func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			verify.
				Ok().
				GasRemaining(simpleGasTestConfig.gasProvided-simpleGasTestConfig.gasUsedByParent-gasUsedByBuiltinClaim).
				GasUsed(parentAddress, simpleGasTestConfig.gasUsedByParent+gasUsedByBuiltinClaim)
		})
}

func TestGasUsed_TwoContracts_ExecuteOnSameCtx(t *testing.T) {

	for numCalls := uint64(0); numCalls < 3; numCalls++ {
		expectedGasRemaining := simpleGasTestConfig.gasProvided - simpleGasTestConfig.gasUsedByParent - simpleGasTestConfig.gasUsedByChild*numCalls
		numCallsBytes := big.NewInt(0).SetUint64(numCalls).Bytes()

		runMockInstanceCallerTestBuilder(t).
			withContracts(
				createMockContract(parentAddress).
					withBalance(simpleGasTestConfig.parentBalance).
					withConfig(simpleGasTestConfig).
					withMethods(addMethodsToParentInstanceMock),
				createMockContract(childAddress).
					withBalance(simpleGasTestConfig.childBalance).
					withConfig(simpleGasTestConfig).
					withMethods(addMethodsToChildInstanceMock),
			).
			withInput(createTestContractCallInputBuilder().
				withRecipientAddr(parentAddress).
				withGasProvided(simpleGasTestConfig.gasProvided).
				withFunction("execOnSameCtx").
				withArguments(childAddress, []byte("wasteGas"), numCallsBytes).
				build()).
			withSetup(func(host *vmHost, world *worldmock.MockWorld) {
				setZeroCodeCosts(host)
			}).
			andAssertResults(func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
				verify.
					Ok().
					GasRemaining(expectedGasRemaining).
					GasUsed(parentAddress, simpleGasTestConfig.gasUsedByParent)
				if numCalls > 0 {
					verify.GasUsed(childAddress, simpleGasTestConfig.gasUsedByChild*numCalls)
				}
			})
	}
}

func TestGasUsed_TwoContracts_ExecuteOnDestCtx(t *testing.T) {

	for numCalls := uint64(0); numCalls < 3; numCalls++ {
		expectedGasRemaining := simpleGasTestConfig.gasProvided - simpleGasTestConfig.gasUsedByParent - simpleGasTestConfig.gasUsedByChild*numCalls
		numCallsBytes := big.NewInt(0).SetUint64(numCalls).Bytes()

		runMockInstanceCallerTestBuilder(t).
			withContracts(
				createMockContract(parentAddress).
					withBalance(simpleGasTestConfig.parentBalance).
					withConfig(simpleGasTestConfig).
					withMethods(addMethodsToParentInstanceMock),
				createMockContract(childAddress).
					withBalance(simpleGasTestConfig.childBalance).
					withConfig(simpleGasTestConfig).
					withMethods(addMethodsToChildInstanceMock),
			).
			withInput(createTestContractCallInputBuilder().
				withRecipientAddr(parentAddress).
				withGasProvided(simpleGasTestConfig.gasProvided).
				withFunction("execOnDestCtx").
				withArguments(childAddress, []byte("wasteGas"), numCallsBytes).
				build()).
			withSetup(func(host *vmHost, world *worldmock.MockWorld) {
				setZeroCodeCosts(host)
			}).
			andAssertResults(func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
				verify.
					Ok().
					GasRemaining(expectedGasRemaining).
					GasUsed(parentAddress, simpleGasTestConfig.gasUsedByParent)
				if numCalls > 0 {
					verify.GasUsed(childAddress, simpleGasTestConfig.gasUsedByChild*numCalls)
				}
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

	runMockInstanceCallerTestBuilder(t).
		withContracts(
			createMockContract(alphaAddress).
				withBalance(0).
				withConfig(directCallGasTestConfig{
					gasUsedByParent:    alphaCallGas,
					gasProvidedToChild: alphaGasToForwardToReceivers,
					gasProvided:        gasProvided,
				}).
				withMethods(addMethodsToParentInstanceMock),
			createMockContract(betaAddress).
				withBalance(0).
				withConfig(directCallGasTestConfig{
					gasUsedByChild: receiverCallGas,
				}).
				withMethods(addMethodsToChildInstanceMock),
			createMockContract(gammaAddress).
				withBalance(0).
				withConfig(directCallGasTestConfig{
					gasUsedByChild: receiverCallGas,
				}).
				withMethods(addMethodsToChildInstanceMock),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(alphaAddress).
			withGasProvided(gasProvided).
			withFunction("execOnDestCtx").
			withArguments(betaAddress, []byte("wasteGas"), arwen.One.Bytes(),
				gammaAddress, []byte("wasteGas"), arwen.One.Bytes()).
			build()).
		withSetup(func(host *vmHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
		}).
		andAssertResults(func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(alphaAddress, alphaCallGas).
				GasUsed(betaAddress, receiverCallGas).
				GasUsed(gammaAddress, receiverCallGas).
				GasRemaining(gasProvided - alphaCallGas - 2*receiverCallGas)
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

	gasUsedByParent := testConfig.gasUsedByParent + testConfig.gasUsedByCallback
	gasUsedByChild := testConfig.gasUsedByChild

	runMockInstanceCallerTestBuilder(t).
		withContracts(
			createMockContract(parentAddress).
				withBalance(testConfig.parentBalance).
				withConfig(testConfig).
				withMethods(addAsyncParentMethodsToInstanceMock),
			createMockContract(childAddress).
				withBalance(testConfig.childBalance).
				withConfig(testConfig).
				withMethods(addAsyncChildMethodsToInstanceMock),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withGasProvided(testConfig.gasProvided).
			withFunction("performAsyncCall").
			withArguments([]byte{0}).
			build()).
		withSetup(func(host *vmHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.gasLockCost)
		}).
		andAssertResults(func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(parentAddress, gasUsedByParent).
				GasUsed(childAddress, gasUsedByChild).
				GasRemaining(testConfig.gasProvided-gasUsedByParent-gasUsedByChild).
				BalanceDelta(thirdPartyAddress, 2*testConfig.transferToThirdParty).
				ReturnData(parentFinishA, parentFinishB, []byte{0}, []byte("thirdparty"), []byte("vault"), []byte{0}, []byte("succ")).
				Storage(
					createStoreEntry(parentAddress).withKey(parentKeyA).withValue(parentDataA),
					createStoreEntry(parentAddress).withKey(parentKeyB).withValue(parentDataB),
					createStoreEntry(childAddress).withKey(childKey).withValue(childData),
				).
				Transfers(
					createTransferEntry(thirdPartyAddress).
						withData([]byte("hello")).
						withValue(big.NewInt(testConfig.transferToThirdParty)).
						withSenderAddress(parentAddress),
					createTransferEntry(thirdPartyAddress).
						withData([]byte(" there")).
						withValue(big.NewInt(testConfig.transferToThirdParty)).
						withSenderAddress(childAddress),
					createTransferEntry(vaultAddress).withData([]byte{}).
						withValue(big.NewInt(testConfig.transferToVault)).
						withSenderAddress(childAddress),
				)
		})
}

func TestGasUsed_AsyncCall_BuiltinCall(t *testing.T) {

	testConfig := asyncBaseTestConfig
	testConfig.gasProvided = 1000

	expectedGasUsedByParent := asyncBaseTestConfig.gasUsedByParent + asyncBaseTestConfig.gasUsedByCallback + gasUsedByBuiltinClaim
	expectedGasUsedByChild := uint64(0) // all gas for builtin call is consummed on caller

	runMockInstanceCallerTestBuilder(t).
		withContracts(
			createMockContract(parentAddress).
				withBalance(asyncBaseTestConfig.parentBalance).
				withConfig(&testConfig).
				withMethods(addAsyncBuiltinParentMethodsToInstanceMock),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withGasProvided(testConfig.gasProvided).
			withFunction("forwardAsyncCall").
			withArguments(userAddress, []byte("builtinClaim")).
			build()).
		withSetup(func(host *vmHost, world *worldmock.MockWorld) {
			world.AcctMap.CreateAccount(userAddress)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, asyncBaseTestConfig.gasLockCost)
		}).
		andAssertResults(func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(parentAddress, expectedGasUsedByParent).
				GasUsed(userAddress, asyncBaseTestConfig.gasUsedByChild).
				GasRemaining(asyncBaseTestConfig.gasProvided - expectedGasUsedByParent - expectedGasUsedByChild)
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

	expectedGasUsedByParent := testConfig.gasUsedByParent + testConfig.gasUsedByCallback
	expectedGasUsedByChild := testConfig.gasUsedByChild + gasUsedByBuiltinClaim

	runMockInstanceCallerTestBuilder(t).
		withContracts(
			createMockContract(parentAddress).
				withBalance(testConfig.parentBalance).
				withConfig(testConfig).
				withMethods(addAsyncBuiltinMultiContractParentMethodsToInstanceMock),
			createMockContract(childAddress).
				withBalance(testConfig.childBalance).
				withConfig(testConfig).
				withMethods(addAsyncBuiltinMultiContractChildMethodsToInstanceMock),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withGasProvided(testConfig.gasProvided).
			withFunction("forwardAsyncCall").
			withArguments(userAddress, childAddress, []byte("childFunction"), []byte("builtinClaim")).
			build()).
		withSetup(func(host *vmHost, world *worldmock.MockWorld) {
			world.AcctMap.CreateAccount(userAddress)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.gasLockCost)
			createMockBuiltinFunctions(t, host, world)
		}).
		andAssertResults(func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(parentAddress, expectedGasUsedByParent).
				GasUsed(childAddress, testConfig.gasUsedByChild).
				GasRemaining(testConfig.gasProvided - expectedGasUsedByParent - expectedGasUsedByChild)
		})
}

func TestGasUsed_AsyncCall_ChildFails(t *testing.T) {

	testConfig := asyncTestConfig

	expectedGasUsedByParent := testConfig.gasProvided - testConfig.gasLockCost + testConfig.gasUsedByCallback

	runMockInstanceCallerTestBuilder(t).
		withContracts(
			createMockContract(parentAddress).
				withBalance(testConfig.parentBalance).
				withConfig(testConfig).
				withMethods(addAsyncParentMethodsToInstanceMock),
			createMockContract(childAddress).
				withBalance(testConfig.childBalance).
				withConfig(testConfig).
				withMethods(addAsyncChildMethodsToInstanceMock),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withGasProvided(testConfig.gasProvided).
			withFunction("performAsyncCall").
			withArguments([]byte{1}).
			withCurrentTxHash([]byte("txhash")).
			build()).
		withSetup(func(host *vmHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.gasLockCost)
		}).
		andAssertResults(func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			verify.
				Ok().
				BalanceDelta(parentAddress, -(testConfig.transferToThirdParty+testConfig.transferToVault)).
				BalanceDelta(thirdPartyAddress, testConfig.transferToThirdParty).
				GasUsed(parentAddress, expectedGasUsedByParent).
				GasRemaining(testConfig.gasProvided-expectedGasUsedByParent).
				ReturnData(parentFinishA, parentFinishB, []byte("succ")).
				Storage(
					createStoreEntry(parentAddress).withKey(parentKeyA).withValue(parentDataA),
					createStoreEntry(parentAddress).withKey(parentKeyB).withValue(parentDataB),
				).
				Transfers(
					createTransferEntry(vaultAddress).
						withData([]byte("child error")).
						withValue(big.NewInt(testConfig.transferToVault)).
						withSenderAddress(parentAddress),
					createTransferEntry(thirdPartyAddress).
						withData([]byte("hello")).
						withValue(big.NewInt(testConfig.transferToThirdParty)).
						withSenderAddress(parentAddress),
				)
		})
}

func TestGasUsed_AsyncCall_CallBackFails(t *testing.T) {

	testConfig := asyncTestConfig

	expectedGasUsedByParent := testConfig.gasProvided - testConfig.gasUsedByChild
	expectedGasUsedByChild := testConfig.gasUsedByChild

	runMockInstanceCallerTestBuilder(t).
		withContracts(
			createMockContract(parentAddress).
				withBalance(testConfig.parentBalance).
				withConfig(testConfig).
				withMethods(addAsyncParentMethodsToInstanceMock),
			createMockContract(childAddress).
				withBalance(testConfig.childBalance).
				withConfig(testConfig).
				withMethods(addAsyncChildMethodsToInstanceMock),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withGasProvided(testConfig.gasProvided).
			withFunction("performAsyncCall").
			withArguments([]byte{0, 3}).
			withCurrentTxHash([]byte("txhash")).
			build()).
		withSetup(func(host *vmHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.gasLockCost)
		}).
		andAssertResults(func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			verify.
				Ok().
				ReturnMessage("callBack error").
				BalanceDelta(parentAddress, -(2*testConfig.transferToThirdParty+testConfig.transferToVault)).
				BalanceDelta(thirdPartyAddress, 2*testConfig.transferToThirdParty).
				GasUsed(parentAddress, expectedGasUsedByParent).
				GasUsed(childAddress, expectedGasUsedByChild).
				GasRemaining(0).
				ReturnData(parentFinishA, parentFinishB, []byte{3}, []byte("thirdparty"), []byte("vault"), []byte("user error"), []byte("txhash")).
				Storage(
					createStoreEntry(parentAddress).withKey(parentKeyA).withValue(parentDataA),
					createStoreEntry(parentAddress).withKey(parentKeyB).withValue(parentDataB),
					createStoreEntry(childAddress).withKey(childKey).withValue(childData),
				).
				Transfers(
					createTransferEntry(thirdPartyAddress).
						withData([]byte("hello")).
						withValue(big.NewInt(testConfig.transferToThirdParty)).
						withSenderAddress(parentAddress),
					createTransferEntry(thirdPartyAddress).
						withData([]byte(" there")).
						withValue(big.NewInt(testConfig.transferToThirdParty)).
						withSenderAddress(childAddress),
					createTransferEntry(vaultAddress).
						withData([]byte{}).
						withValue(big.NewInt(testConfig.transferToVault)).
						withSenderAddress(childAddress),
				)
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

	expectedGasUsedByParent := testConfig.gasUsedByParent + testConfig.gasUsedByCallback
	expectedGasUsedByChild := uint64(testConfig.recursiveChildCalls)*testConfig.gasUsedByChild +
		uint64(testConfig.recursiveChildCalls-1)*testConfig.gasUsedByCallback

	runMockInstanceCallerTestBuilder(t).
		withContracts(
			createMockContract(parentAddress).
				withBalance(testConfig.parentBalance).
				withConfig(testConfig).
				withMethods(addAsyncRecursiveParentMethodsToInstanceMock),
			createMockContract(childAddress).
				withBalance(testConfig.childBalance).
				withConfig(&testConfig.asyncCallBaseTestConfig).
				withMethods(addAsyncRecursiveChildMethodsToInstanceMock),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withGasProvided(testConfig.gasProvided).
			withFunction("forwardAsyncCall").
			withArguments(childAddress, []byte("recursiveAsyncCall"), big.NewInt(int64(testConfig.recursiveChildCalls)).Bytes()).
			build()).
		withSetup(func(host *vmHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.gasLockCost)
		}).
		andAssertResults(func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			verify.
				Ok().
				BalanceDelta(parentAddress, -testConfig.transferFromParentToChild).
				Transfers(
					createTransferEntry(childAddress).
						withData([]byte("hello")).
						withValue(big.NewInt(testConfig.transferFromParentToChild)).
						withSenderAddress(parentAddress),
				).
				GasUsed(parentAddress, expectedGasUsedByParent).
				GasUsed(childAddress, expectedGasUsedByChild).
				GasRemaining(testConfig.gasProvided-expectedGasUsedByParent-expectedGasUsedByChild).
				BalanceDelta(childAddress, testConfig.transferFromParentToChild)
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

	expectedGasUsedByParent := testConfig.gasUsedByParent + testConfig.gasUsedByCallback
	expectedGasUsedByChild := uint64(testConfig.childCalls) * testConfig.gasUsedByChild

	runMockInstanceCallerTestBuilder(t).
		withContracts(
			createMockContract(parentAddress).
				withBalance(testConfig.parentBalance).
				withConfig(testConfig).
				withMethods(addAsyncMultiChildParentMethodsToInstanceMock),
			createMockContract(childAddress).
				withBalance(testConfig.childBalance).
				withConfig(&testConfig.asyncCallBaseTestConfig).
				withMethods(addAsyncRecursiveChildMethodsToInstanceMock),
		).
		withInput(createTestContractCallInputBuilder().
			withRecipientAddr(parentAddress).
			withGasProvided(testConfig.gasProvided).
			withFunction("forwardAsyncCall").
			withArguments(childAddress, []byte("recursiveAsyncCall")).
			build()).
		withSetup(func(host *vmHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.gasLockCost)
		}).
		andAssertResults(func(world *worldmock.MockWorld, verify *VMOutputVerifier) {
			verify.
				Ok().
				BalanceDelta(parentAddress, -testConfig.transferFromParentToChild).
				BalanceDelta(childAddress, testConfig.transferFromParentToChild).
				Transfers(
					createTransferEntry(childAddress).
						withData([]byte("hello")).
						withValue(big.NewInt(testConfig.transferFromParentToChild)).
						withSenderAddress(parentAddress),
				).
				GasUsed(parentAddress, expectedGasUsedByParent).
				GasUsed(childAddress, expectedGasUsedByChild).
				GasRemaining(testConfig.gasProvided - expectedGasUsedByParent - expectedGasUsedByChild)
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

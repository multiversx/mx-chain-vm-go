package hostCoretest

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/txDataBuilder"
	contextmock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/context"
	"github.com/multiversx/mx-chain-vm-v1_4-go/mock/contracts"
	worldmock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/world"
	gasSchedules "github.com/multiversx/mx-chain-vm-v1_4-go/scenarioexec/gasSchedules"
	test "github.com/multiversx/mx-chain-vm-v1_4-go/testcommon"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	"github.com/stretchr/testify/require"
)

var gasUsedByBuiltinClaim = uint64(120)

var simpleGasTestConfig = contracts.DirectCallGasTestConfig{
	GasUsedByParent:    uint64(400),
	GasUsedByChild:     uint64(200),
	GasProvidedToChild: uint64(300),
	GasProvided:        uint64(1000),
	ParentBalance:      int64(1000),
	ChildBalance:       int64(1000),
}

func TestGasUsed_SingleContract(t *testing.T) {
	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(simpleGasTestConfig.ParentBalance).
				WithConfig(simpleGasTestConfig).
				WithMethods(contracts.WasteGasParentMock)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(simpleGasTestConfig.GasProvided).
			WithFunction("wasteGas").
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasRemaining(simpleGasTestConfig.GasProvided-simpleGasTestConfig.GasUsedByParent).
				GasUsed(test.ParentAddress, simpleGasTestConfig.GasUsedByParent)
		})
}

func TestGasUsed_SingleContract_BuiltinCall(t *testing.T) {
	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(simpleGasTestConfig.ParentBalance).
				WithConfig(simpleGasTestConfig).
				WithMethods(contracts.ExecOnDestCtxParentMock)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(simpleGasTestConfig.GasProvided).
			WithFunction("execOnDestCtx").
			WithArguments(test.ParentAddress, []byte("builtinClaim"), vmhost.One.Bytes()).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasRemaining(simpleGasTestConfig.GasProvided-simpleGasTestConfig.GasUsedByParent-gasUsedByBuiltinClaim).
				GasUsed(test.ParentAddress, simpleGasTestConfig.GasUsedByParent+gasUsedByBuiltinClaim)
		})
}

func TestGasUsed_SingleContract_BuiltinCallFail(t *testing.T) {
	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(simpleGasTestConfig.ParentBalance).
				WithConfig(simpleGasTestConfig).
				WithMethods(contracts.ExecOnDestCtxSingleCallParentMock)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(simpleGasTestConfig.GasProvided).
			WithFunction("execOnDestCtxSingleCall").
			WithArguments(test.ParentAddress, []byte("builtinFail")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ExecutionFailed).
				ReturnMessage("return value 1").
				HasRuntimeErrors("whatdidyoudo").
				GasRemaining(0)
		})
}

func TestGasUsed_TwoContracts_ExecuteOnSameCtx(t *testing.T) {
	for numCalls := uint64(0); numCalls < 3; numCalls++ {
		expectedGasRemaining := simpleGasTestConfig.GasProvided - simpleGasTestConfig.GasUsedByParent - simpleGasTestConfig.GasUsedByChild*numCalls
		numCallsBytes := big.NewInt(0).SetUint64(numCalls).Bytes()

		test.BuildMockInstanceCallTest(t).
			WithContracts(
				test.CreateMockContract(test.ParentAddress).
					WithBalance(simpleGasTestConfig.ParentBalance).
					WithConfig(simpleGasTestConfig).
					WithMethods(contracts.ExecOnSameCtxParentMock, contracts.ExecOnDestCtxParentMock, contracts.WasteGasParentMock),
				test.CreateMockContract(test.ChildAddress).
					WithBalance(simpleGasTestConfig.ChildBalance).
					WithConfig(simpleGasTestConfig).
					WithMethods(contracts.WasteGasChildMock),
			).
			WithInput(test.CreateTestContractCallInputBuilder().
				WithRecipientAddr(test.ParentAddress).
				WithGasProvided(simpleGasTestConfig.GasProvided).
				WithFunction("execOnSameCtx").
				WithArguments(test.ChildAddress, []byte("wasteGas"), numCallsBytes).
				Build()).
			WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
				setZeroCodeCosts(host)
			}).
			AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
				verify.
					Ok().
					GasRemaining(expectedGasRemaining).
					GasUsed(test.ParentAddress, simpleGasTestConfig.GasUsedByParent+simpleGasTestConfig.GasUsedByChild*numCalls)
				if numCalls > 0 {
					verify.GasUsed(test.ChildAddress, 0)
				}
			})
	}
}

func TestGasUsed_TwoContracts_ExecuteOnDestCtx(t *testing.T) {
	for numCalls := uint64(0); numCalls < 3; numCalls++ {
		expectedGasRemaining := simpleGasTestConfig.GasProvided - simpleGasTestConfig.GasUsedByParent - simpleGasTestConfig.GasUsedByChild*numCalls
		numCallsBytes := big.NewInt(0).SetUint64(numCalls).Bytes()

		test.BuildMockInstanceCallTest(t).
			WithContracts(
				test.CreateMockContract(test.ParentAddress).
					WithBalance(simpleGasTestConfig.ParentBalance).
					WithConfig(simpleGasTestConfig).
					WithMethods(contracts.ExecOnSameCtxParentMock, contracts.ExecOnDestCtxParentMock, contracts.WasteGasParentMock),
				test.CreateMockContract(test.ChildAddress).
					WithBalance(simpleGasTestConfig.ChildBalance).
					WithConfig(simpleGasTestConfig).
					WithMethods(contracts.WasteGasChildMock),
			).
			WithInput(test.CreateTestContractCallInputBuilder().
				WithRecipientAddr(test.ParentAddress).
				WithGasProvided(simpleGasTestConfig.GasProvided).
				WithFunction("execOnDestCtx").
				WithArguments(test.ChildAddress, []byte("wasteGas"), numCallsBytes).
				Build()).
			WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
				setZeroCodeCosts(host)
			}).
			AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
				verify.
					Ok().
					GasRemaining(expectedGasRemaining).
					GasUsed(test.ParentAddress, simpleGasTestConfig.GasUsedByParent)
				if numCalls > 0 {
					verify.GasUsed(test.ChildAddress, simpleGasTestConfig.GasUsedByChild*numCalls)
				}
			})
	}
}

func TestGasUsed_ThreeContracts_ExecuteOnDestCtx(t *testing.T) {
	alphaAddress := test.MakeTestSCAddress("alpha")
	betaAddress := test.MakeTestSCAddress("beta")
	gammaAddress := test.MakeTestSCAddress("gamma")

	gasProvided := uint64(1000)
	alphaCallGas := uint64(400)
	alphaGasToForwardToReceivers := uint64(300)
	receiverCallGas := uint64(200)

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(alphaAddress).
				WithBalance(0).
				WithConfig(contracts.DirectCallGasTestConfig{
					GasUsedByParent:    alphaCallGas,
					GasProvidedToChild: alphaGasToForwardToReceivers,
					GasProvided:        gasProvided,
				}).
				WithMethods(contracts.ExecOnSameCtxParentMock, contracts.ExecOnDestCtxParentMock, contracts.WasteGasParentMock),
			test.CreateMockContract(betaAddress).
				WithBalance(0).
				WithConfig(contracts.DirectCallGasTestConfig{
					GasUsedByChild: receiverCallGas,
				}).
				WithMethods(contracts.WasteGasChildMock),
			test.CreateMockContract(gammaAddress).
				WithBalance(0).
				WithConfig(contracts.DirectCallGasTestConfig{
					GasUsedByChild: receiverCallGas,
				}).
				WithMethods(contracts.WasteGasChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(alphaAddress).
			WithGasProvided(gasProvided).
			WithFunction("execOnDestCtx").
			WithArguments(betaAddress, []byte("wasteGas"), vmhost.One.Bytes(),
				gammaAddress, []byte("wasteGas"), vmhost.One.Bytes()).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(alphaAddress, alphaCallGas).
				GasUsed(betaAddress, receiverCallGas).
				GasUsed(gammaAddress, receiverCallGas).
				GasRemaining(gasProvided - alphaCallGas - 2*receiverCallGas)
		})
}

func TestGasUsed_ESDTTransfer_ThenExecuteCall_Success(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)
	esdtTransferGasCost := uint64(1)

	testConfig := simpleGasTestConfig
	testConfig.ESDTTokensToTransfer = 5

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ExecESDTTransferAndCallChild),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.WasteGasChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execESDTTransferAndCall").
			WithArguments(test.ChildAddress, []byte("ESDTTransfer"), []byte("wasteGas")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(test.ParentAddress, testConfig.GasUsedByParent+esdtTransferGasCost).
				GasUsed(test.ChildAddress, testConfig.GasUsedByChild).
				GasRemaining(testConfig.GasProvided - esdtTransferGasCost - testConfig.GasUsedByParent - testConfig.GasUsedByChild)

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, testConfig.ESDTTokensToTransfer, childESDTBalance)
		})
}

func TestGasUsed_ESDTTransfer_ThenExecuteCall_Fail(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := simpleGasTestConfig
	testConfig.ESDTTokensToTransfer = 5

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ExecESDTTransferAndCallChild),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.FailChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execESDTTransferAndCall").
			WithArguments(test.ChildAddress, []byte("ESDTTransfer"), []byte("fail")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ExecutionFailed).
				HasRuntimeErrors("forced fail").
				GasRemaining(0)

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, uint64(0), childESDTBalance)
		})
}

func TestGasUsed_ESDTTransferFailed(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := simpleGasTestConfig
	testConfig.ESDTTokensToTransfer = 2 * initialESDTTokenBalance

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ExecESDTTransferAndCallChild),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.FailChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execESDTTransferAndCall").
			WithArguments(test.ChildAddress, []byte("ESDTTransfer"), []byte("fail")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				ReturnCode(vmcommon.ExecutionFailed).
				HasRuntimeErrors("insufficient funds").
				GasRemaining(0)

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, uint64(0), childESDTBalance)
		})
}

func TestMultipleTimes(t *testing.T) {
	for i := 0; i < 20; i++ {
		TestGasUsed_ESDTTransferFromParent_ChildBurnsAndThenFails(t)
	}
}

func TestGasUsed_ESDTTransferFromParent_ChildBurnsAndThenFails(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)
	testConfig := simpleGasTestConfig
	testConfig.ESDTTokensToTransfer = 10

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ExecESDTTransferWithAPICall),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.FailChildAndBurnESDTMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execESDTTransferWithAPICall").
			WithArguments(test.ChildAddress, []byte("failAndBurn"), big.NewInt(int64(testConfig.ESDTTokensToTransfer)).Bytes()).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			_ = childAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, 0)
			_ = childAccount.SetTokenRolesAsStrings(test.ESDTTestTokenName, []string{core.ESDTRoleLocalBurn})
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.ReturnCode(vmcommon.ExecutionFailed)

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, uint64(0), childESDTBalance)
		})
}

var asyncBaseTestConfig = contracts.AsyncCallBaseTestConfig{
	GasProvided:       1000,
	GasUsedByParent:   400,
	GasUsedByChild:    200,
	GasUsedByCallback: 100,
	GasLockCost:       150,

	TransferFromParentToChild: 7,

	ParentBalance: 1000,
	ChildBalance:  1000,
}

var asyncTestConfig = &contracts.AsyncCallTestConfig{
	AsyncCallBaseTestConfig: asyncBaseTestConfig,
	TransferToThirdParty:    3,
	TransferToVault:         4,
	ESDTTokensToTransfer:    0,
}

var transferAndExecuteTestConfig = contracts.TransferAndExecuteTestConfig{
	DirectCallGasTestConfig: contracts.DirectCallGasTestConfig{
		GasProvided:     1000,
		GasUsedByParent: 200,
		ParentBalance:   1000,
	},
	TransferFromParentToChild: 5,
	GasTransferToChild:        100,
}

func TestGasUsed_AsyncCall(t *testing.T) {
	testConfig := asyncTestConfig
	testConfig.GasProvided = 1000

	gasUsedByParent := testConfig.GasUsedByParent + testConfig.GasUsedByCallback
	gasUsedByChild := testConfig.GasUsedByChild

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.PerformAsyncCallParentMock, contracts.CallBackParentMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.TransferToThirdPartyAsyncChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("performAsyncCall").
			WithArguments([]byte{0}).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(test.ParentAddress, gasUsedByParent).
				GasUsed(test.ChildAddress, gasUsedByChild).
				GasRemaining(testConfig.GasProvided-gasUsedByParent-gasUsedByChild).
				BalanceDelta(test.ThirdPartyAddress, 2*testConfig.TransferToThirdParty).
				ReturnData(test.ParentFinishA, test.ParentFinishB, []byte{0}, []byte("thirdparty"), []byte("vault"), []byte{0}, []byte("succ")).
				Storage(
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
					test.CreateStoreEntry(test.ChildAddress).WithKey(test.ChildKey).WithValue(test.ChildData),
				).
				Transfers(
					test.CreateTransferEntry(test.ParentAddress, test.ThirdPartyAddress).
						WithData([]byte("hello")).
						WithValue(big.NewInt(testConfig.TransferToThirdParty)),
					test.CreateTransferEntry(test.ChildAddress, test.ThirdPartyAddress).
						WithData([]byte(" there")).
						WithValue(big.NewInt(testConfig.TransferToThirdParty)),
					test.CreateTransferEntry(test.ChildAddress, test.VaultAddress).
						WithData([]byte{}).
						WithValue(big.NewInt(testConfig.TransferToVault)),
				)
		})
}

func TestGasUsed_AsyncCall_CrossShard_InitCall(t *testing.T) {
	testConfig := asyncTestConfig
	testConfig.GasProvided = 1000

	gasUsedByParent := testConfig.GasUsedByParent

	asyncCallData := txDataBuilder.NewBuilder()
	asyncCallData.Func(contracts.AsyncChildFunction)
	asyncCallData.Int64(testConfig.TransferToThirdParty)
	asyncCallData.Str(contracts.AsyncChildData)
	// behavior param for child
	asyncCallData.Bytes([]byte{0})
	asyncChildArgs := asyncCallData.ToBytes()

	gasForAsyncCall := testConfig.GasProvided - gasUsedByParent - testConfig.GasLockCost

	parentContract := test.CreateMockContractOnShard(test.ParentAddress, 0).
		WithBalance(testConfig.ParentBalance).
		WithConfig(testConfig).
		WithMethods(contracts.PerformAsyncCallParentMock, contracts.CallBackParentMock)

	// direct parent call
	test.BuildMockInstanceCallTest(t).
		WithContracts(parentContract).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithCallerAddr(test.UserAddress).
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("performAsyncCall").
			WithArguments([]byte{0}).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			world.SelfShardID = 0
			if world.CurrentBlockInfo == nil {
				world.CurrentBlockInfo = &worldmock.BlockInfo{}
			}
			world.CurrentBlockInfo.BlockRound = 0
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(test.ParentAddress, gasUsedByParent).
				GasRemaining(0).
				ReturnData(test.ParentFinishA, test.ParentFinishB).
				Storage(
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
				).
				Transfers(
					test.CreateTransferEntry(test.ParentAddress, test.ThirdPartyAddress).
						WithData([]byte("hello")).
						WithValue(big.NewInt(testConfig.TransferToThirdParty)),
					test.CreateTransferEntry(test.ParentAddress, test.ChildAddress).
						WithData(asyncChildArgs).
						WithGasLimit(gasForAsyncCall).
						WithGasLocked(testConfig.GasLockCost).
						WithCallType(vm.AsynchronousCall).
						WithValue(big.NewInt(testConfig.TransferFromParentToChild)),
				)
		})
}

func TestGasUsed_AsyncCall_CrossShard_ExecuteCall(t *testing.T) {
	testConfig := asyncTestConfig
	gasUsedByChild := testConfig.GasUsedByChild
	gasUsedByParent := testConfig.GasUsedByParent
	gasForAsyncCall := testConfig.GasProvided - gasUsedByParent - testConfig.GasLockCost

	childAsyncReturnData := [][]byte{{0}, []byte("thirdparty"), []byte("vault")}

	// async cross-shard parent -> child
	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContractOnShard(test.ChildAddress, 1).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.TransferToThirdPartyAsyncChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithCallerAddr(test.ParentAddress).
			WithRecipientAddr(test.ChildAddress).
			WithCallValue(testConfig.TransferFromParentToChild).
			WithGasProvided(gasForAsyncCall).
			WithFunction(contracts.AsyncChildFunction).
			WithArguments(
				big.NewInt(testConfig.TransferToThirdParty).Bytes(),
				[]byte(contracts.AsyncChildData),
				[]byte{0}).
			WithCallType(vm.AsynchronousCall).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			world.SelfShardID = 1
			if world.CurrentBlockInfo == nil {
				world.CurrentBlockInfo = &worldmock.BlockInfo{}
			}
			world.CurrentBlockInfo.BlockRound = 1
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(test.ChildAddress, gasUsedByChild).
				GasRemaining(0).
				ReturnData(childAsyncReturnData...).
				Transfers(
					test.CreateTransferEntry(test.ChildAddress, test.ThirdPartyAddress).
						WithData([]byte(contracts.AsyncChildData)).
						WithValue(big.NewInt(testConfig.TransferToThirdParty)),
					test.CreateTransferEntry(test.ChildAddress, test.VaultAddress).
						WithData([]byte{}).
						WithValue(big.NewInt(testConfig.TransferToVault)),
					test.CreateTransferEntry(test.ChildAddress, test.ParentAddress).
						WithData(computeReturnDataForCallback(vmcommon.Ok, childAsyncReturnData)).
						WithGasLimit(gasForAsyncCall-gasUsedByChild).
						WithCallType(vm.AsynchronousCallBack).
						WithValue(big.NewInt(0)),
				)
		})
}

func TestGasUsed_AsyncCall_CrossShard_ExecuteCall_WithTransfer(t *testing.T) {
	testConfig := asyncTestConfig
	gasUsedByChild := testConfig.GasUsedByChild
	gasUsedByParent := testConfig.GasUsedByParent
	gasForAsyncCall := testConfig.GasProvided - gasUsedByParent - testConfig.GasLockCost

	// async cross-shard parent -> child
	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContractOnShard(test.ChildAddress, 1).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.TransferToAsyncParentOnCallbackChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithCallerAddr(test.ParentAddress).
			WithRecipientAddr(test.ChildAddress).
			WithCallValue(testConfig.TransferFromParentToChild).
			WithGasProvided(gasForAsyncCall).
			WithFunction(contracts.AsyncChildFunction).
			WithArguments(
				big.NewInt(testConfig.TransferToThirdParty).Bytes()).
			WithCallType(vm.AsynchronousCall).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			world.SelfShardID = 1
			if world.CurrentBlockInfo == nil {
				world.CurrentBlockInfo = &worldmock.BlockInfo{}
			}
			world.CurrentBlockInfo.BlockRound = 1
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(test.ChildAddress, gasUsedByChild).
				GasRemaining(0).
				ReturnData().
				Transfers(
					test.CreateTransferEntry(test.ChildAddress, test.ParentAddress).
						WithGasLimit(0).
						WithCallType(vm.DirectCall).
						WithValue(big.NewInt(testConfig.TransferToThirdParty)),
					test.CreateTransferEntry(test.ChildAddress, test.ParentAddress).
						WithData(computeReturnDataForCallback(vmcommon.Ok, nil)).
						WithGasLimit(gasForAsyncCall-gasUsedByChild).
						WithCallType(vm.AsynchronousCallBack).
						WithValue(big.NewInt(0)),
				)
		})
}

func TestGasUsed_AsyncCall_CrossShard_CallBack(t *testing.T) {
	testConfig := asyncTestConfig
	testConfig.GasProvided = 1000

	gasUsedByParent := testConfig.GasUsedByParent
	gasUsedByChild := testConfig.GasUsedByChild
	gasForAsyncCall := testConfig.GasProvided - gasUsedByParent - testConfig.GasLockCost

	parentContract := test.CreateMockContractOnShard(test.ParentAddress, 0).
		WithBalance(testConfig.ParentBalance).
		WithConfig(testConfig).
		WithMethods(contracts.PerformAsyncCallParentMock, contracts.CallBackParentMock)

	// async cross shard callback child -> parent
	test.BuildMockInstanceCallTest(t).
		WithContracts(parentContract).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithCallerAddr(test.ChildAddress).
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasForAsyncCall-gasUsedByChild+asyncBaseTestConfig.GasLockCost).
			WithFunction("callBack").
			WithArguments([]byte{}, []byte{0}, []byte("thirdparty"), []byte("vault")).
			WithCallType(vm.AsynchronousCallBack).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			world.SelfShardID = 0
			if world.CurrentBlockInfo == nil {
				world.CurrentBlockInfo = &worldmock.BlockInfo{}
			}
			world.CurrentBlockInfo.BlockRound = 2
			// Mock the storage as if the parent was already executed
			accountHandler, _ := world.GetUserAccount(test.ParentAddress)
			(accountHandler.(*worldmock.Account)).Storage[string(test.ParentKeyA)] = test.ParentDataA
			(accountHandler.(*worldmock.Account)).Storage[string(test.ParentKeyB)] = test.ParentDataB
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasRemaining(testConfig.GasProvided-gasUsedByParent-gasUsedByChild-asyncTestConfig.GasUsedByCallback).
				ReturnData([]byte{0}, []byte("succ"))
		})
}

func TestGasUsed_AsyncCall_BuiltinCall(t *testing.T) {
	testConfig := asyncBaseTestConfig
	testConfig.GasProvided = 1000

	expectedGasUsedByParent := testConfig.GasUsedByParent + testConfig.GasUsedByCallback + gasUsedByBuiltinClaim
	expectedGasUsedByChild := uint64(0) // all gas for builtin call is consummed on caller

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(&testConfig).
				WithMethods(contracts.ForwardAsyncCallParentBuiltinMock, contracts.CallBackParentBuiltinMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("forwardAsyncCall").
			WithArguments(test.UserAddress, []byte("builtinClaim")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			world.AcctMap.CreateAccount(test.UserAddress, world)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(test.ParentAddress, expectedGasUsedByParent).
				GasUsed(test.UserAddress, 0).
				GasRemaining(testConfig.GasProvided - expectedGasUsedByParent - expectedGasUsedByChild)
		})
}

func TestGasUsed_AsyncCall_BuiltinCallFail(t *testing.T) {
	testConfig := asyncBaseTestConfig
	testConfig.GasProvided = 1000

	// all will be spent in case of failure
	gasProvidedForBuiltinCall := testConfig.GasProvided - testConfig.GasUsedByParent - testConfig.GasLockCost

	expectedGasUsedByParent := testConfig.GasUsedByParent + gasProvidedForBuiltinCall + testConfig.GasUsedByCallback
	expectedGasUsedByChild := uint64(0) // all gas for builtin call is consummed on caller

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(&testConfig).
				WithMethods(contracts.ForwardAsyncCallParentBuiltinMock, contracts.CallBackParentBuiltinMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("forwardAsyncCall").
			WithArguments(test.UserAddress, []byte("builtinFail")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			world.AcctMap.CreateAccount(test.UserAddress, world)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				HasRuntimeErrors("whatdidyoudo").
				GasUsed(test.ParentAddress, expectedGasUsedByParent).
				GasUsed(test.UserAddress, 0).
				GasRemaining(testConfig.GasProvided - expectedGasUsedByParent - expectedGasUsedByChild)
		})
}

func TestGasUsed_AsyncCall_BuiltinMultiContractCall(t *testing.T) {
	// TODO no possible yet, reactivate when new async context is merged
	t.Skip()

	testConfig := &contracts.AsyncBuiltInCallTestConfig{
		AsyncCallBaseTestConfig:   asyncBaseTestConfig,
		TransferFromChildToParent: 5,
	}

	expectedGasUsedByParent := testConfig.GasUsedByParent + testConfig.GasUsedByCallback
	expectedGasUsedByChild := testConfig.GasUsedByChild + gasUsedByBuiltinClaim

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ForwardAsyncCallMultiContractParentMock, contracts.CallBackMultiContractParentMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.RecursiveAsyncCallRecursiveChildMock, contracts.CallBackRecursiveChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("forwardAsyncCall").
			WithArguments(test.UserAddress, test.ChildAddress, []byte("childFunction"), []byte("builtinClaim")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			world.AcctMap.CreateAccount(test.UserAddress, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
			createMockBuiltinFunctions(t, host, world)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(test.ParentAddress, expectedGasUsedByParent).
				GasUsed(test.ChildAddress, testConfig.GasUsedByChild).
				GasRemaining(testConfig.GasProvided - expectedGasUsedByParent - expectedGasUsedByChild)
		})
}

func TestGasUsed_AsyncCall_ChildFails(t *testing.T) {
	testConfig := asyncTestConfig
	testConfig.GasProvided = 1000

	expectedGasUsedByParent := testConfig.GasProvided - testConfig.GasLockCost + testConfig.GasUsedByCallback

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.PerformAsyncCallParentMock, contracts.CallBackParentMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.TransferToThirdPartyAsyncChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("performAsyncCall").
			WithArguments([]byte{1}).
			WithCurrentTxHash([]byte("txhash")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				HasRuntimeErrors("child error").
				BalanceDelta(test.ParentAddress, -(testConfig.TransferToThirdParty+testConfig.TransferToVault)).
				BalanceDelta(test.ThirdPartyAddress, testConfig.TransferToThirdParty).
				GasUsed(test.ParentAddress, expectedGasUsedByParent).
				GasUsed(test.ChildAddress, 0).
				GasRemaining(testConfig.GasProvided-expectedGasUsedByParent).
				ReturnData(test.ParentFinishA, test.ParentFinishB, []byte("succ")).
				Storage(
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
				).
				Transfers(
					test.CreateTransferEntry(test.ParentAddress, test.VaultAddress).
						WithData([]byte("child error")).
						WithValue(big.NewInt(testConfig.TransferToVault)),
					test.CreateTransferEntry(test.ParentAddress, test.ThirdPartyAddress).
						WithData([]byte("hello")).
						WithValue(big.NewInt(testConfig.TransferToThirdParty)),
				)
		})
}

func TestGasUsed_AsyncCall_CallBackFails(t *testing.T) {
	testConfig := asyncTestConfig

	expectedGasUsedByParent := testConfig.GasProvided - testConfig.GasUsedByChild
	expectedGasUsedByChild := testConfig.GasUsedByChild

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.PerformAsyncCallParentMock, contracts.CallBackParentMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.TransferToThirdPartyAsyncChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("performAsyncCall").
			WithArguments([]byte{3}).
			WithCurrentTxHash([]byte("txhash")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				ReturnMessage("callBack error").
				HasRuntimeErrors("callBack error").
				BalanceDelta(test.ParentAddress, -(2*testConfig.TransferToThirdParty+testConfig.TransferToVault)).
				BalanceDelta(test.ThirdPartyAddress, 2*testConfig.TransferToThirdParty).
				GasUsed(test.ParentAddress, expectedGasUsedByParent).
				GasUsed(test.ChildAddress, expectedGasUsedByChild).
				GasRemaining(0).
				ReturnData(test.ParentFinishA, test.ParentFinishB, []byte{3}, []byte("thirdparty"), []byte("vault"), []byte("user error"), []byte("txhash")).
				Storage(
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
					test.CreateStoreEntry(test.ChildAddress).WithKey(test.ChildKey).WithValue(test.ChildData),
				).
				Transfers(
					test.CreateTransferEntry(test.ParentAddress, test.ThirdPartyAddress).
						WithData([]byte("hello")).
						WithValue(big.NewInt(testConfig.TransferToThirdParty)),
					test.CreateTransferEntry(test.ChildAddress, test.ThirdPartyAddress).
						WithData([]byte(" there")).
						WithValue(big.NewInt(testConfig.TransferToThirdParty)),
					test.CreateTransferEntry(test.ChildAddress, test.VaultAddress).
						WithData([]byte{}).
						WithValue(big.NewInt(testConfig.TransferToVault)),
				)
		})
}

func TestGasUsed_AsyncCall_Recursive(t *testing.T) {
	// TODO no possible yet, reactivate when new async context is merged
	t.Skip()

	testConfig := &contracts.AsyncCallRecursiveTestConfig{
		AsyncCallBaseTestConfig: asyncBaseTestConfig,
		RecursiveChildCalls:     2,
	}

	expectedGasUsedByParent := testConfig.GasUsedByParent + testConfig.GasUsedByCallback
	expectedGasUsedByChild := uint64(testConfig.RecursiveChildCalls)*testConfig.GasUsedByChild +
		uint64(testConfig.RecursiveChildCalls-1)*testConfig.GasUsedByCallback

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ForwardAsyncCallRecursiveParentMock, contracts.CallBackRecursiveParentMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(&testConfig.AsyncCallBaseTestConfig).
				WithMethods(contracts.RecursiveAsyncCallRecursiveChildMock, contracts.CallBackRecursiveChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("forwardAsyncCall").
			WithArguments(test.ChildAddress, []byte("recursiveAsyncCall"), big.NewInt(int64(testConfig.RecursiveChildCalls)).Bytes()).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				BalanceDelta(test.ParentAddress, -testConfig.TransferFromParentToChild).
				Transfers(
					test.CreateTransferEntry(test.ParentAddress, test.ChildAddress).
						WithData([]byte("hello")).
						WithValue(big.NewInt(testConfig.TransferFromParentToChild)),
				).
				GasUsed(test.ParentAddress, expectedGasUsedByParent).
				GasUsed(test.ChildAddress, expectedGasUsedByChild).
				GasRemaining(testConfig.GasProvided-expectedGasUsedByParent-expectedGasUsedByChild).
				BalanceDelta(test.ChildAddress, testConfig.TransferFromParentToChild)
		})
}

func TestGasUsed_AsyncCall_MultiChild(t *testing.T) {
	// TODO no possible yet, reactivate when new async context is merged
	t.Skip()

	testConfig := &contracts.AsyncCallMultiChildTestConfig{
		AsyncCallBaseTestConfig: asyncBaseTestConfig,
		ChildCalls:              2,
	}

	expectedGasUsedByParent := testConfig.GasUsedByParent + testConfig.GasUsedByCallback
	expectedGasUsedByChild := uint64(testConfig.ChildCalls) * testConfig.GasUsedByChild

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ForwardAsyncCallMultiChildMock, contracts.CallBackMultiChildMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(&testConfig.AsyncCallBaseTestConfig).
				WithMethods(contracts.RecursiveAsyncCallRecursiveChildMock, contracts.CallBackRecursiveChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("forwardAsyncCall").
			WithArguments(test.ChildAddress, []byte("recursiveAsyncCall")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				BalanceDelta(test.ParentAddress, -testConfig.TransferFromParentToChild).
				BalanceDelta(test.ChildAddress, testConfig.TransferFromParentToChild).
				Transfers(
					test.CreateTransferEntry(test.ParentAddress, test.ChildAddress).
						WithData([]byte("hello")).
						WithValue(big.NewInt(testConfig.TransferFromParentToChild)),
				).
				GasUsed(test.ParentAddress, expectedGasUsedByParent).
				GasUsed(test.ChildAddress, expectedGasUsedByChild).
				GasRemaining(testConfig.GasProvided - expectedGasUsedByParent - expectedGasUsedByChild)
		})
}

func TestGasUsed_ESDTTransfer_ThenExecuteAsyncCall_Success(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := asyncTestConfig
	testConfig.ESDTTokensToTransfer = 5

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ExecESDTTransferAndAsyncCallChild),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.WasteGasChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execESDTTransferAndAsyncCall").
			WithArguments(test.ChildAddress, []byte("ESDTTransfer"), []byte("wasteGas")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok()

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, testConfig.ESDTTokensToTransfer, childESDTBalance)
		})
}

func TestGasUsed_ESDTTransfer_ThenExecuteAsyncCall_ChildFails(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := asyncTestConfig
	testConfig.ESDTTokensToTransfer = 5

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ExecESDTTransferAndAsyncCallChild),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.FailChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execESDTTransferAndAsyncCall").
			WithArguments(test.ChildAddress, []byte("ESDTTransfer"), []byte("fail")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				HasRuntimeErrors(vmhost.ErrNotEnoughGas.Error())

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, uint64(0), childESDTBalance)
		})
}

func TestGasUsed_ESDTTransfer_ThenExecuteAsyncCall_CallbackFails(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := asyncTestConfig
	testConfig.ESDTTokensToTransfer = 5

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ExecESDTTransferAndAsyncCallChild, contracts.CallBackParentMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.WasteGasChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execESDTTransferAndAsyncCall").
			WithArguments(test.ChildAddress, []byte("ESDTTransfer"), []byte("wasteGas")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				ReturnMessage("wrong num of arguments")

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, testConfig.ESDTTokensToTransfer, childESDTBalance)
		})
}

func TestGasUsed_ESDTTransferInCallback(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := asyncTestConfig
	testConfig.ESDTTokensToTransfer = 5
	testConfig.CallbackESDTTokensToTransfer = 2

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ExecESDTTransferAndAsyncCallChild, contracts.SimpleCallbackMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ESDTTransferToParentMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execESDTTransferAndAsyncCall").
			WithArguments(test.ChildAddress, []byte("ESDTTransfer"), []byte("transferESDTToParent")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok()

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer+testConfig.CallbackESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, testConfig.ESDTTokensToTransfer-testConfig.CallbackESDTTokensToTransfer, childESDTBalance)
		})
}

func TestGasUsed_ESDTTransferWrongArgNumberForCallback(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := asyncTestConfig
	testConfig.ESDTTokensToTransfer = 5
	testConfig.CallbackESDTTokensToTransfer = 2

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ExecESDTTransferAndAsyncCallChild, contracts.SimpleCallbackMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ESDTTransferToParentWrongESDTArgsNumberMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execESDTTransferAndAsyncCall").
			WithArguments(test.ChildAddress, []byte("ESDTTransfer"), []byte("transferESDTToParent")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok()

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance-(testConfig.ESDTTokensToTransfer-testConfig.CallbackESDTTokensToTransfer), parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, testConfig.ESDTTokensToTransfer-testConfig.CallbackESDTTokensToTransfer, childESDTBalance)
		})
}

func TestGasUsed_ESDTTransfer_CallbackFail(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := asyncTestConfig
	testConfig.ESDTTokensToTransfer = 5
	testConfig.CallbackESDTTokensToTransfer = 2

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ExecESDTTransferAndAsyncCallChild, contracts.SimpleCallbackMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ESDTTransferToParentCallbackWillFail),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execESDTTransferAndAsyncCall").
			WithArguments(test.ChildAddress, []byte("ESDTTransfer"), []byte("transferESDTToParent")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				HasRuntimeErrors("wrong num of arguments")

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, testConfig.ESDTTokensToTransfer, childESDTBalance)
		})
}

func TestGasUsed_TransferAndExecute_CrossShard(t *testing.T) {
	testConfig := transferAndExecuteTestConfig

	noOfTransfers := 3

	childContracts := []test.MockTestSmartContract{
		test.CreateMockContractOnShard(test.ParentAddress, 0).
			WithBalance(testConfig.ParentBalance).
			WithConfig(testConfig).
			WithMethods(contracts.TransferAndExecute),
	}

	startShard := 1
	for transfer := 0; transfer < noOfTransfers; transfer++ {
		childContracts = append(childContracts,
			test.CreateMockContractOnShard(contracts.GetChildAddressForTransfer(transfer), uint32(startShard+transfer)).
				WithBalance(0).
				WithConfig(testConfig).
				WithCodeMetadata([]byte{0, vmcommon.MetadataPayable}).
				WithMethods(contracts.WasteGasChildMock))
	}

	expectedTransfers := make([]test.TransferEntry, 0)
	expectedLogs := make([]vmcommon.LogEntry, 0)
	for transfer := 0; transfer < noOfTransfers; transfer++ {
		expectedTransfer := test.CreateTransferEntry(test.ParentAddress, contracts.GetChildAddressForTransfer(transfer)).
			WithData(big.NewInt(int64(transfer)).Bytes()).
			WithGasLimit(testConfig.GasTransferToChild).
			WithValue(big.NewInt(testConfig.TransferFromParentToChild))
		expectedTransfers = append(expectedTransfers, expectedTransfer)
		expectedLogs = append(expectedLogs, vmcommon.LogEntry{
			Address: test.ParentAddress,
			Topics: [][]byte{
				test.ParentAddress,
				contracts.GetChildAddressForTransfer(transfer),
				big.NewInt(testConfig.TransferFromParentToChild).Bytes()},
			Data:       []byte{},
			Identifier: []byte("transferValueOnly"),
		})
	}

	gasRemaining := testConfig.GasProvided - testConfig.GasUsedByParent - uint64(noOfTransfers)*testConfig.GasTransferToChild

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			childContracts...,
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithCallerAddr(test.UserAddress).
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction(contracts.TransferAndExecuteFuncName).
			WithArguments(big.NewInt(int64(noOfTransfers)).Bytes()).
			WithCallType(vm.DirectCall).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(test.ParentAddress, testConfig.GasUsedByParent).
				GasRemaining(gasRemaining).
				ReturnData(contracts.TransferAndExecuteReturnData).
				Transfers(expectedTransfers...).
				Logs(expectedLogs...)
		})
}

func TestGasUsed_AsyncCallManaged_Mocks(t *testing.T) {
	testConfig := asyncTestConfig
	startValue := uint64(3000)
	outOfGasValue := uint64(150)
	stopValue := uint64(100)
	decrement := uint64(1)

	tester := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.GasMismatchAsyncCallParentMock, contracts.GasMismatchCallBackParentMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.GasMismatchAsyncCallChildMock),
		).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		})

	for gasLimit := startValue; gasLimit >= stopValue; gasLimit -= decrement {

		tester.WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasLimit).
			WithFunction("gasMismatchParent").
			Build()).
			AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
				if gasLimit > outOfGasValue {
					verify.
						Ok()
				} else {
					verify.
						OutOfGas()
				}
			})
	}
}

func TestGasUsed_AsyncCallManaged(t *testing.T) {
	startValue := uint64(10000000)
	outOfGasValue := uint64(5400000)
	stopValue := uint64(5000000)
	decrement := uint64(1000)

	gasSchedule, err := gasSchedules.LoadGasScheduleConfig(gasSchedules.GetV4())
	require.Nil(t, err)

	tester := test.BuildInstanceCallTest(t).
		WithContracts(
			test.CreateInstanceContract(test.ParentAddress).
				WithCode(test.GetTestSCCode("async-call-parent-managed", "../../")).
				WithBalance(1000),
			test.CreateInstanceContract(test.ChildAddress).
				WithCode(test.GetTestSCCode("async-call-child-managed", "../../")).
				WithBalance(1000),
		).
		WithGasSchedule(gasSchedule)

	for gasLimit := startValue; gasLimit >= stopValue; gasLimit -= decrement {
		tester.WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithFunction("foo").
			WithGasProvided(gasLimit).
			WithArguments(test.ChildAddress).
			Build()).
			AndAssertResults(func(host vmhost.VMHost, stubBlockchainHook *contextmock.BlockchainHookStub, verify *test.VMOutputVerifier) {
				if gasLimit > outOfGasValue {
					verify.
						Ok()
				} else {
					verify.
						OutOfGas()
				}
			})
	}
}

func TestGasUsed_AsyncESDTTransfer(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := asyncTestConfig
	testConfig.ESDTTokensToTransfer = 5

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithOwnerAddress(test.UserAddress).
				WithMethods(contracts.ExecESDTTransferInAsyncCall, contracts.EvilCallback),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithOwnerAddress(test.UserAddress2).
				WithMethods(contracts.WasteGasChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("esdtTransferInAsyncCall").
			WithArguments(test.UserAddress2).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			world.AcctMap.CreateAccount(test.UserAddress, world)
			world.AcctMap.CreateAccount(test.UserAddress2, world)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				ReturnMessage("execution by caller failed")

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.UserAddress2)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, testConfig.ESDTTokensToTransfer, childESDTBalance)
		})
}

func TestGasUsed_Async_CallbackWithOnSameContext(t *testing.T) {
	testConfig := asyncTestConfig
	testConfig.GasProvided = 1000

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithOwnerAddress(test.UserAddress).
				WithMethods(contracts.PerformAsyncCallParentMock, contracts.CallbackWithOnSameContext),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithOwnerAddress(test.UserAddress2).
				WithMethods(contracts.TransferToThirdPartyAsyncChildMock, contracts.ExecutedOnSameContextByCallback),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("performAsyncCall").
			WithArguments([]byte{0}).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
			world.AcctMap.CreateAccount(test.UserAddress, world)
			world.AcctMap.CreateAccount(test.UserAddress2, world)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				Storage(
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
					// overriden by ExecutedOnSameContextByCallback called from CallbackWithOnSameContext
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataA),
					test.CreateStoreEntry(test.ChildAddress).WithKey(test.ChildKey).WithValue(test.ChildData),
				)
		})
}

type MockClaimBuiltin struct {
	test.MockBuiltin
	AmountToGive int64
	GasCost      uint64
}

func createMockBuiltinFunctions(tb testing.TB, host vmhost.VMHost, world *worldmock.MockWorld) {
	err := world.InitBuiltinFunctions(host.GetGasScheduleMap())
	require.Nil(tb, err)

	mockClaimBuiltin := &MockClaimBuiltin{
		AmountToGive: 42,
		GasCost:      gasUsedByBuiltinClaim,
	}

	_ = world.BuiltinFuncs.Container.Add("builtinClaim", &test.MockBuiltin{
		ProcessBuiltinFunctionCall: func(acntSnd, _ vmcommon.UserAccountHandler, vmInput *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
			vmOutput := test.MakeVMOutput()
			test.AddNewOutputAccount(
				vmOutput,
				nil,
				acntSnd.AddressBytes(),
				mockClaimBuiltin.AmountToGive,
				nil)
			vmOutput.GasRemaining = vmInput.GasProvided - mockClaimBuiltin.GasCost
			return vmOutput, nil
		},
	})

	_ = world.BuiltinFuncs.Container.Add("builtinFail", &test.MockBuiltin{
		ProcessBuiltinFunctionCall: func(acntSnd, _ vmcommon.UserAccountHandler, vmInput *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
			return nil, errors.New("whatdidyoudo")
		},
	})

	host.SetBuiltInFunctionsContainer(world.BuiltinFuncs.Container)
}

func setZeroCodeCosts(host vmhost.VMHost) {
	host.Metering().GasSchedule().BaseOperationCost.CompilePerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.AoTPreparePerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.GetCode = 0
	host.Metering().GasSchedule().BaseOperationCost.StorePerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.DataCopyPerByte = 0
	host.Metering().GasSchedule().BaseOpsAPICost.SignalError = 0
	host.Metering().GasSchedule().BaseOpsAPICost.ExecuteOnSameContext = 0
	host.Metering().GasSchedule().BaseOpsAPICost.ExecuteOnDestContext = 0
	host.Metering().GasSchedule().BaseOpsAPICost.TransferValue = 0
}

func setAsyncCosts(host vmhost.VMHost, gasLock uint64) {
	host.Metering().GasSchedule().BaseOpsAPICost.AsyncCallStep = 0
	host.Metering().GasSchedule().BaseOpsAPICost.AsyncCallbackGasLock = gasLock
}

func computeReturnDataForCallback(returnCode vmcommon.ReturnCode, returnData [][]byte) []byte {
	retData := []byte("@" + core.ConvertToEvenHex(int(returnCode)))
	for _, data := range returnData {
		retData = append(retData, []byte("@"+hex.EncodeToString(data))...)
	}
	return retData
}

package hosttest

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/contracts"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/world"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/testcommon"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	"github.com/stretchr/testify/require"
)

var gasUsedByBuiltinClaim = uint64(120)

func makeTestConfig() *test.TestConfig {
	return &test.TestConfig{
		GasProvided:           2000,
		GasProvidedToChild:    300,
		GasProvidedToCallback: 50,
		GasUsedByParent:       400,
		GasUsedByChild:        200,
		GasUsedByCallback:     100,
		GasLockCost:           150,

		TransferFromParentToChild: 7,

		ParentBalance:        1000,
		ChildBalance:         1000,
		TransferToThirdParty: 3,
		TransferToVault:      4,
		ESDTTokensToTransfer: 0,
	}
}

func TestGasUsed_SingleContract(t *testing.T) {
	testConfig := makeTestConfig()

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.WasteGasParentMock)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("wasteGas").
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasRemaining(testConfig.GasProvided-testConfig.GasUsedByParent).
				GasUsed(test.ParentAddress, testConfig.GasUsedByParent)
		})
}

func TestGasUsed_SingleContract_BuiltinCall(t *testing.T) {
	testConfig := makeTestConfig()

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ExecOnDestCtxParentMock)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execOnDestCtx").
			WithArguments(test.ParentAddress, []byte("builtinClaim"), arwen.One.Bytes()).
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasRemaining(testConfig.GasProvided-testConfig.GasUsedByParent-gasUsedByBuiltinClaim).
				GasUsed(test.ParentAddress, testConfig.GasUsedByParent+gasUsedByBuiltinClaim).
				BalanceDelta(test.ParentAddress, amountToGiveByBuiltinClaim)
		})
}

func TestGasUsed_SingleContract_BuiltinCallFail(t *testing.T) {
	testConfig := makeTestConfig()

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ExecOnDestCtxSingleCallParentMock)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execOnDestCtxSingleCall").
			WithArguments(test.ParentAddress, []byte("builtinFail")).
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.ExecutionFailed().
				ReturnMessage("Return value 1").
				HasRuntimeErrors("whatdidyoudo").
				GasRemaining(0)
		})
}

func TestGasUsed_TwoContracts_ExecuteOnSameCtx(t *testing.T) {
	testConfig := makeTestConfig()

	for numCalls := uint64(0); numCalls < 3; numCalls++ {
		expectedGasRemaining := testConfig.GasProvided - testConfig.GasUsedByParent - testConfig.GasUsedByChild*numCalls
		numCallsBytes := big.NewInt(0).SetUint64(numCalls).Bytes()

		test.BuildMockInstanceCallTest(t).
			WithContracts(
				test.CreateMockContract(test.ParentAddress).
					WithBalance(testConfig.ParentBalance).
					WithConfig(testConfig).
					WithMethods(contracts.ExecOnSameCtxParentMock, contracts.ExecOnDestCtxParentMock, contracts.WasteGasParentMock),
				test.CreateMockContract(test.ChildAddress).
					WithBalance(testConfig.ChildBalance).
					WithConfig(testConfig).
					WithMethods(contracts.WasteGasChildMock),
			).
			WithInput(test.CreateTestContractCallInputBuilder().
				WithRecipientAddr(test.ParentAddress).
				WithGasProvided(testConfig.GasProvided).
				WithFunction("execOnSameCtx").
				WithArguments(test.ChildAddress, []byte("wasteGas"), numCallsBytes).
				Build()).
			WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
				setZeroCodeCosts(host)
			}).
			AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
				verify.Ok().
					GasRemaining(expectedGasRemaining).
					GasUsed(test.ParentAddress, testConfig.GasUsedByParent)
				if numCalls > 0 {
					verify.GasUsed(test.ChildAddress, testConfig.GasUsedByChild*numCalls)
				}
			})
	}
}

func TestGasUsed_TwoContracts_ExecuteOnDestCtx(t *testing.T) {
	testConfig := makeTestConfig()

	for numCalls := uint64(0); numCalls < 3; numCalls++ {
		expectedGasRemaining := testConfig.GasProvided - testConfig.GasUsedByParent - testConfig.GasUsedByChild*numCalls
		numCallsBytes := big.NewInt(0).SetUint64(numCalls).Bytes()

		test.BuildMockInstanceCallTest(t).
			WithContracts(
				test.CreateMockContract(test.ParentAddress).
					WithBalance(testConfig.ParentBalance).
					WithConfig(testConfig).
					WithMethods(contracts.ExecOnSameCtxParentMock, contracts.ExecOnDestCtxParentMock, contracts.WasteGasParentMock),
				test.CreateMockContract(test.ChildAddress).
					WithBalance(testConfig.ChildBalance).
					WithConfig(testConfig).
					WithMethods(contracts.WasteGasChildMock),
			).
			WithInput(test.CreateTestContractCallInputBuilder().
				WithRecipientAddr(test.ParentAddress).
				WithGasProvided(testConfig.GasProvided).
				WithFunction("execOnDestCtx").
				WithArguments(test.ChildAddress, []byte("wasteGas"), numCallsBytes).
				Build()).
			WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
				setZeroCodeCosts(host)
			}).
			AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
				verify.Ok().
					GasRemaining(expectedGasRemaining).
					GasUsed(test.ParentAddress, testConfig.GasUsedByParent)
				if numCalls > 0 {
					verify.GasUsed(test.ChildAddress, testConfig.GasUsedByChild*numCalls)
				}
			})
	}
}

func TestGasUsed_ThreeContracts_ExecuteOnDestCtx(t *testing.T) {
	alphaAddress := test.MakeTestSCAddress("alpha")
	betaAddress := test.MakeTestSCAddress("beta")
	gammaAddress := test.MakeTestSCAddress("gamma")

	testConfig := &test.TestConfig{
		GasUsedByParent:    uint64(400),
		GasProvidedToChild: uint64(300),
		GasProvided:        uint64(1000),
		GasUsedByChild:     uint64(200),
	}

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(alphaAddress).
				WithBalance(0).
				WithConfig(testConfig).
				WithMethods(contracts.ExecOnSameCtxParentMock, contracts.ExecOnDestCtxParentMock, contracts.WasteGasParentMock),
			test.CreateMockContract(betaAddress).
				WithBalance(0).
				WithConfig(testConfig).
				WithMethods(contracts.WasteGasChildMock),
			test.CreateMockContract(gammaAddress).
				WithBalance(0).
				WithConfig(testConfig).
				WithMethods(contracts.WasteGasChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(alphaAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execOnDestCtx").
			WithArguments(betaAddress, []byte("wasteGas"), arwen.One.Bytes(),
				gammaAddress, []byte("wasteGas"), arwen.One.Bytes()).
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasUsed(alphaAddress, testConfig.GasUsedByParent).
				GasUsed(betaAddress, testConfig.GasUsedByChild).
				GasUsed(gammaAddress, testConfig.GasUsedByChild).
				GasRemaining(testConfig.GasProvided - testConfig.GasUsedByParent - 2*testConfig.GasUsedByChild)
		})
}

func TestGasUsed_ESDTTransfer_ThenExecuteCall_Success(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)
	esdtTransferGasCost := uint64(1)

	testConfig := makeTestConfig()
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenKey, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasUsed(test.ParentAddress, testConfig.GasUsedByParent+esdtTransferGasCost).
				GasUsed(test.ChildAddress, testConfig.GasUsedByChild).
				GasRemaining(testConfig.GasProvided - esdtTransferGasCost - testConfig.GasUsedByParent - testConfig.GasUsedByChild)

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, testConfig.ESDTTokensToTransfer, childESDTBalance)
		})
}

func TestGasUsed_ESDTTransfer_ThenExecuteCall_Fail(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenKey, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.ExecutionFailed().
				HasRuntimeErrors("forced fail").
				GasRemaining(0)

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, initialESDTTokenBalance, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, uint64(0), childESDTBalance)
		})
}

func TestGasUsed_ESDTTransferFailed(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenKey, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.ExecutionFailed().
				HasRuntimeErrors("insufficient funds").
				GasRemaining(0)

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, initialESDTTokenBalance, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, uint64(0), childESDTBalance)
		})
}

func TestGasUsed_ESDTTransferFromParent_ChildBurnsAndThenFails(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenKey, initialESDTTokenBalance)
			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			_ = childAccount.SetTokenBalanceUint64(test.ESDTTestTokenKey, 0)
			_ = childAccount.SetTokenRolesAsStrings(test.ESDTTestTokenName, []string{core.ESDTRoleLocalBurn})
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				HasRuntimeErrors("forced fail")

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, initialESDTTokenBalance, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, uint64(0), childESDTBalance)
		})
}

func TestGasUsed_AsyncCall(t *testing.T) {
	testConfig := makeTestConfig()
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
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
	testConfig := makeTestConfig()
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			world.SelfShardID = 0
			world.CurrentBlockInfo.BlockRound = 0
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
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
	testConfig := makeTestConfig()
	gasForAsyncCall := testConfig.GasProvided - testConfig.GasUsedByParent - testConfig.GasLockCost

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
			WithGasProvided(gasForAsyncCall).
			WithFunction(contracts.AsyncChildFunction).
			WithArguments(
				big.NewInt(testConfig.TransferToThirdParty).Bytes(),
				[]byte(contracts.AsyncChildData),
				[]byte{0}).
			WithCallType(vm.AsynchronousCall).
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			world.SelfShardID = 1
			world.CurrentBlockInfo.BlockRound = 1
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasUsed(test.ChildAddress, testConfig.GasUsedByChild).
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
						WithGasLimit(gasForAsyncCall-testConfig.GasUsedByChild).
						WithCallType(vm.AsynchronousCallBack).
						WithValue(big.NewInt(0)),
				)
		})
}

func TestGasUsed_AsyncCall_CrossShard_CallBack(t *testing.T) {
	testConfig := makeTestConfig()
	testConfig.GasProvided = 1000

	gasUsedByParent := testConfig.GasUsedByParent
	gasUsedByChild := testConfig.GasUsedByChild
	gasForAsyncCall := testConfig.GasProvided - gasUsedByParent - testConfig.GasLockCost

	parentContract := test.CreateMockContractOnShard(test.ParentAddress, 0).
		WithBalance(testConfig.ParentBalance).
		WithConfig(testConfig).
		WithMethods(contracts.PerformAsyncCallParentMock, contracts.CallBackParentMock)

	arguments := [][]byte{{}, {0}, []byte("thirdparty"), []byte("vault")}

	// async cross shard callback child -> parent
	test.BuildMockInstanceCallTest(t).
		WithContracts(parentContract).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithCallerAddr(test.ChildAddress).
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasForAsyncCall - gasUsedByChild + testConfig.GasLockCost).
			WithFunction("callBack").
			WithArguments(arguments...).
			WithCallType(vm.AsynchronousCallBack).
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			world.SelfShardID = 0
			world.CurrentBlockInfo.BlockRound = 2

			// Mock the storage as if the parent was already executed
			accountHandler, _ := world.GetUserAccount(test.ParentAddress)
			(accountHandler.(*worldmock.Account)).Storage[string(test.ParentKeyA)] = test.ParentDataA
			(accountHandler.(*worldmock.Account)).Storage[string(test.ParentKeyB)] = test.ParentDataB

			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)

			// TODO factor this setup out if necessary for other tests

			// The instance started below will be cached on the InstanceMockBuilder and reused by doRunSmartContractCall().
			// This is necessary for gas usage metering during Save() below.
			// Note that the InstanceMockBuilder uses the address of the contract as
			// if it were its bytecode, hence StartWasmerInstance() receives an
			// address as its first argument.
			host.Runtime().StartWasmerInstance(test.ParentAddress, testConfig.GasUsedByParent, false)

			fakeInput := host.Runtime().GetVMInput()
			fakeInput.GasProvided = 1000
			host.Metering().InitStateFromContractCallInput(fakeInput)

			contracts.RegisterAsyncCallToChild(host, testConfig, arguments)
			host.Async().Save()

			for _, account := range host.Output().GetVMOutput().OutputAccounts {
				for _, storageUpdate := range account.StorageUpdates {
					(accountHandler.(*worldmock.Account)).Storage[string(storageUpdate.Offset)] = storageUpdate.Data
				}
			}
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasRemaining(testConfig.GasProvided-gasUsedByParent-gasUsedByChild-testConfig.GasUsedByCallback).
				ReturnData([]byte{0}, []byte("succ"))
		})
}

func TestGasUsed_LegacyAsyncCall_InShard_BuiltinCall(t *testing.T) {
	testConfig := makeTestConfig()
	testConfig.GasProvided = 1000

	expectedGasUsedByParent := testConfig.GasUsedByParent + testConfig.GasUsedByCallback + gasUsedByBuiltinClaim
	expectedGasUsedByChild := uint64(0) // all gas for builtin call is consummed on caller

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ForwardAsyncCallParentBuiltinMock, contracts.CallBackParentBuiltinMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("forwardAsyncCall").
			WithArguments(test.UserAddress, []byte("builtinClaim")).
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			world.AcctMap.CreateAccount(test.UserAddress, world)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				BalanceDelta(test.ParentAddress, amountToGiveByBuiltinClaim).
				GasUsed(test.ParentAddress, expectedGasUsedByParent).
				GasUsed(test.UserAddress, 0).
				GasRemaining(testConfig.GasProvided - expectedGasUsedByParent - expectedGasUsedByChild)
		})
}

func TestGasUsed_LegacyAsyncCall_BuiltinCallFail(t *testing.T) {
	testConfig := makeTestConfig()
	testConfig.GasProvided = 1000

	// all will be spent in case of failure
	gasProvidedForBuiltinCall := testConfig.GasProvided - testConfig.GasUsedByParent - testConfig.GasLockCost

	expectedGasUsedByParent := testConfig.GasUsedByParent + gasProvidedForBuiltinCall + testConfig.GasUsedByCallback

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ForwardAsyncCallParentBuiltinMock, contracts.CallBackParentBuiltinMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("forwardAsyncCall").
			WithArguments(test.UserAddress, []byte("builtinFail")).
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			world.AcctMap.CreateAccount(test.UserAddress, world)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				HasRuntimeErrors("whatdidyoudo").
				GasUsed(test.ParentAddress, expectedGasUsedByParent).
				GasUsed(test.UserAddress, 0).
				GasRemaining(testConfig.GasProvided - expectedGasUsedByParent)
		})
}

func TestGasUsed_LegacyAsyncCall_CrossShard_BuiltinCall(t *testing.T) {
	testConfig := makeTestConfig()
	testConfig.GasProvided = 1000

	expectedGasUsedByParent := testConfig.GasUsedByParent + gasUsedByBuiltinClaim

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithShardID(1).
				WithMethods(contracts.ForwardAsyncCallParentBuiltinMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("forwardAsyncCall").
			WithArguments(test.UserAddress, []byte("sendMessage")).
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			world.SelfShardID = 1
			world.AcctMap.CreateAccount(test.UserAddress, world)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasRemaining(0).
				GasUsed(test.ParentAddress, expectedGasUsedByParent)
		})
}

func TestGasUsed_AsyncCall_BuiltinMultiContractChainCall(t *testing.T) {
	testConfig := makeTestConfig()
	testConfig.TransferFromChildToParent = 5

	// expectedGasUsedByParent := testConfig.GasUsedByParent + testConfig.GasUsedByCallback
	// expectedGasUsedByChild := testConfig.GasUsedByChild + gasUsedByBuiltinClaim

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ForwardAsyncCallMultiContractParentMock, contracts.CallBackMultiContractParentMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ForwardAsyncCallParentBuiltinMock, contracts.CallBackParentBuiltinMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("forwardAsyncCall").
			WithArguments(test.ChildAddress, []byte("forwardAsyncCall"), []byte("builtinClaim") /*, arwen.One.Bytes()*/).
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			world.AcctMap.CreateAccount(test.UserAddress, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
			createMockBuiltinFunctions(t, host, world)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok()
			// TODO matei-p add asserts when in-shard builtin works
			// TODO matei-p enable gas assertions
			// GasUsed(test.ParentAddress, expectedGasUsedByParent).
			// GasUsed(test.ChildAddress, testConfig.GasUsedByChild).
			// GasRemaining(testConfig.GasProvided - expectedGasUsedByParent - expectedGasUsedByChild)
		})
}

func TestGasUsed_AsyncCall_ChildFails(t *testing.T) {
	testConfig := makeTestConfig()
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
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
	testConfig := makeTestConfig()

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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.UserError().
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
	//TODO reenable test after contracts are allowed to call themselves
	// repeatedly with async calls (see restriction in asyncContext.addAsyncCall())
	t.Skip("recursive async self-call currently disabled")

	testConfig := makeTestConfig()
	testConfig.RecursiveChildCalls = 3

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
				WithConfig(testConfig).
				WithMethods(contracts.RecursiveAsyncCallRecursiveChildMock, contracts.CallBackRecursiveChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("forwardAsyncCall").
			WithArguments(test.ChildAddress, []byte("recursiveAsyncCall"), big.NewInt(int64(testConfig.RecursiveChildCalls)).Bytes()).
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				BalanceDelta(test.ParentAddress, -testConfig.TransferFromParentToChild).
				GasUsed(test.ParentAddress, expectedGasUsedByParent).
				GasUsed(test.ChildAddress, expectedGasUsedByChild).
				GasRemaining(testConfig.GasProvided-expectedGasUsedByParent-expectedGasUsedByChild).
				BalanceDelta(test.ChildAddress, testConfig.TransferFromParentToChild).
				ReturnData(big.NewInt(2).Bytes(), big.NewInt(1).Bytes(), big.NewInt(0).Bytes())
		})
}

func TestGasUsed_AsyncCall_MultiChild(t *testing.T) {
	testConfig := makeTestConfig()
	testConfig.ChildCalls = 2

	expectedGasUsedByParent := testConfig.GasUsedByParent + 2*testConfig.GasUsedByCallback
	expectedGasUsedByChild := uint64(testConfig.ChildCalls) * testConfig.GasUsedByChild

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ForwardAsyncCallMultiChildMock, contracts.CallBackMultiChildMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.RecursiveAsyncCallRecursiveChildMock, contracts.CallBackRecursiveChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("forwardAsyncCall").
			WithArguments(test.ChildAddress, []byte("recursiveAsyncCall")).
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				BalanceDelta(test.ParentAddress, -2*testConfig.TransferFromParentToChild).
				BalanceDelta(test.ChildAddress, 2*testConfig.TransferFromParentToChild).
				GasUsed(test.ParentAddress, expectedGasUsedByParent).
				GasUsed(test.ChildAddress, expectedGasUsedByChild).
				GasRemaining(testConfig.GasProvided-expectedGasUsedByParent-expectedGasUsedByChild).
				ReturnData(big.NewInt(0).Bytes(), big.NewInt(1).Bytes())
		})
}

func TestGasUsed_ESDTTransfer_ThenExecuteAsyncCall_Success(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenKey, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok()

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, testConfig.ESDTTokensToTransfer, childESDTBalance)
		})
}

func TestGasUsed_ESDTTransfer_ThenExecuteAsyncCall_ChildFails(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
	testConfig.ESDTTokensToTransfer = 5

	expectedGasRemaining := uint64(50)
	gasUsedByParent := testConfig.GasProvided - expectedGasRemaining

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ExecESDTTransferAndAsyncCallChild, contracts.SimpleCallbackMock),
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenKey, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasRemaining(50).
				GasUsed(test.ParentAddress, gasUsedByParent).
				GasUsed(test.ChildAddress, 0)

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, initialESDTTokenBalance, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, uint64(0), childESDTBalance)
		})
}

func TestGasUsed_ESDTTransfer_ThenExecuteAsyncCall_CallbackFails(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenKey, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.UserError().
				ReturnMessage("wrong num of arguments")

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, testConfig.ESDTTokensToTransfer, childESDTBalance)
		})
}

func TestGasUsed_ESDTTransferInCallback(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenKey, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok()

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer+testConfig.CallbackESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, testConfig.ESDTTokensToTransfer-testConfig.CallbackESDTTokensToTransfer, childESDTBalance)
		})
}

func TestGasUsed_ESDTTransferWrongArgNumberForCallback(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenKey, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.ExecutionFailed().
				HasRuntimeErrors("tokenize failed")

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, testConfig.ESDTTokensToTransfer, childESDTBalance)
		})
}

func TestGasUsed_ESDTTransfer_CallbackFail(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenKey, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.UserError().
				HasRuntimeErrors("callback failed intentionally").
				Print()

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenKey)
			require.Equal(t, testConfig.ESDTTokensToTransfer, childESDTBalance)
		})
}

func TestGasUsed_AsyncCall_Groups(t *testing.T) {
	testConfig := makeTestConfig()
	testConfig.GasProvided = 10_000
	testConfig.GasLockCost = 10
	testConfig.GasProvidedToCallback = 60

	asyncGroupCallbackEnabled := false
	asyncContextCallbackEnabled := false
	expectedReturnData := make([][]byte, 0)
	for _, groupConfig := range contracts.AsyncGroupsConfig {
		groupName := groupConfig[0]
		for g := 1; g < len(groupConfig); g++ {
			functionReturnData := groupConfig[g] + test.TestReturnDataSuffix
			expectedReturnData = append(expectedReturnData, []byte(functionReturnData))
			expectedReturnData = append(expectedReturnData, []byte(test.TestCallbackPrefix+functionReturnData))
		}
		if asyncGroupCallbackEnabled {
			expectedReturnData = append(expectedReturnData, []byte(test.TestCallbackPrefix+groupName+test.TestReturnDataSuffix))
		}
	}
	if asyncContextCallbackEnabled {
		expectedReturnData = append(expectedReturnData, []byte(test.TestContextCallbackFunction+test.TestReturnDataSuffix))
	}

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ForwardAsyncCallMultiGroupsMock, contracts.CallBackMultiGroupsMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ChildAsyncMultiGroupsMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("forwardMultiGroupAsyncCall").
			WithArguments(test.ChildAddress).
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				Print().
				ReturnDataDoesNotContain([]byte("out of gas")).
				ReturnData(expectedReturnData...)
		})
}

func TestGasUsed_AsyncCall_CallGraph(t *testing.T) {
	testConfig := makeTestConfig()
	testConfig.GasProvided = 100_000
	testConfig.GasProvidedToChild = 30_000

	// callGraph := test.CreateGraphTestSimple1()
	callGraph := test.CreateGraphTest2()

	runGraphCallTestTemplate(t, testConfig, callGraph)
}

func TestGasUsed_AsyncCall_CallGraph_ContextCallback(t *testing.T) {
	t.Skip("context callbacks are disabled")

	testConfig := makeTestConfig()
	testConfig.GasProvided = 100_000
	testConfig.GasProvidedToChild = 60_000

	callGraph := test.CreateTestCallGraph()
	sc1f1 := callGraph.AddStartNode("sc1", "f1", 0, 0)
	sc2f2 := callGraph.AddNode("sc2", "f2")
	sc3f3 := callGraph.AddNode("sc3", "f3")

	callGraph.AddAsyncEdge(sc1f1, sc2f2, "", "")
	callGraph.AddAsyncEdge(sc2f2, sc3f3, "", "")

	sc1ctxcb := callGraph.AddNode("sc1", "ctxcb1")
	callGraph.SetContextCallback(sc1f1, sc1ctxcb)

	sc2ctxcb := callGraph.AddNode("sc2", "ctxcb2")
	callGraph.SetContextCallback(sc2f2, sc2ctxcb)

	runGraphCallTestTemplate(t, testConfig, callGraph)
}

func runGraphCallTestTemplate(t *testing.T, testConfig *test.TestConfig, callGraph *test.TestCallGraph) {
	expectedReturnData := make([][]byte, 0)
	executionGraph := callGraph.CreateExecutionGraphFromCallGraph()
	startNode := executionGraph.GetStartNode()

	executionOrder := test.CreateRunExpectationOrder(executionGraph)
	for _, testCall := range executionOrder {
		expectedReturnData = append(expectedReturnData, []byte(string(testCall.ContractAddress)+"_"+testCall.FunctionName+test.TestReturnDataSuffix))
	}

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContractsFromAsyncTestCallGraph(callGraph, testConfig)...,
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr([]byte(startNode.Call.ContractAddress)).
			WithGasProvided(testConfig.GasProvided).
			WithFunction(startNode.Call.FunctionName).
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				ReturnData(expectedReturnData...)
		})
}

type MockClaimBuiltin struct {
	test.MockBuiltin
	AmountToGive int64
	GasCost      uint64
}

var amountToGiveByBuiltinClaim = int64(42)

func createMockBuiltinFunctions(tb testing.TB, host arwen.VMHost, world *worldmock.MockWorld) {
	err := world.InitBuiltinFunctions(host.GetGasScheduleMap())
	require.Nil(tb, err)

	mockClaimBuiltin := &MockClaimBuiltin{
		AmountToGive: amountToGiveByBuiltinClaim,
		GasCost:      gasUsedByBuiltinClaim,
	}

	_ = world.BuiltinFuncs.Container.Add("builtinClaim", &test.MockBuiltin{
		ProcessBuiltinFunctionCall: func(acntSnd, _ vmcommon.UserAccountHandler, vmInput *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
			vmOutput := test.MakeEmptyVMOutput()
			test.AddNewOutputTransfer(
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

	world.BuiltinFuncs.Container.Add("sendMessage", &test.MockBuiltin{
		ProcessBuiltinFunctionCall: func(acntSnd, acntRecv vmcommon.UserAccountHandler, vmInput *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
			vmOutput := test.MakeEmptyVMOutput()
			if acntRecv != nil {
				// acntSnd and acntRecv are in the same shard
				test.AddFinishData(vmOutput, []byte("ok"))
				vmOutput.GasRemaining = vmInput.GasProvided - 120
				return vmOutput, nil
			}

			// acntSnd and acntRecv are in different shards
			account := test.AddNewOutputTransfer(
				vmOutput,
				acntSnd.AddressBytes(),
				vmInput.RecipientAddr,
				0,
				[]byte("message"),
			)
			account.OutputTransfers[0].GasLimit = vmInput.GasProvided - 120
			vmOutput.GasRemaining = 0
			return vmOutput, nil
		},
	})

	host.SetBuiltInFunctionsContainer(world.BuiltinFuncs.Container)
}

func setZeroCodeCosts(host arwen.VMHost) {
	host.Metering().GasSchedule().BaseOperationCost.CompilePerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.AoTPreparePerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.GetCode = 0
	host.Metering().GasSchedule().BaseOperationCost.StorePerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.DataCopyPerByte = 0
	host.Metering().GasSchedule().ElrondAPICost.SignalError = 0
	host.Metering().GasSchedule().ElrondAPICost.ExecuteOnSameContext = 0
	host.Metering().GasSchedule().ElrondAPICost.ExecuteOnDestContext = 0
	host.Metering().GasSchedule().ElrondAPICost.StorageLoad = 0
	host.Metering().GasSchedule().ElrondAPICost.StorageStore = 0
}

func setAsyncCosts(host arwen.VMHost, gasLock uint64) {
	host.Metering().GasSchedule().ElrondAPICost.AsyncCallStep = 0
	host.Metering().GasSchedule().ElrondAPICost.AsyncCallbackGasLock = gasLock
}

func computeReturnDataForCallback(returnCode vmcommon.ReturnCode, returnData [][]byte) []byte {
	retCode := string(big.NewInt(int64(returnCode)).Bytes())
	retData := []byte("@" + retCode)
	for _, data := range returnData {
		retData = append(retData, []byte("@"+hex.EncodeToString(data))...)
	}
	return retData
}

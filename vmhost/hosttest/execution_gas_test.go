//nolint:all
package hostCoretest

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/txDataBuilder"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/mock/contracts"
	gasSchedules "github.com/multiversx/mx-chain-vm-go/scenario/gasSchedules"
	"github.com/multiversx/mx-chain-vm-go/testcommon"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var gasUsedByBuiltinClaim = uint64(120)

var LegacyAsyncCallType = []byte{0}
var NewAsyncCallType = []byte{1}

func makeTestConfig() *test.TestConfig {
	return &test.TestConfig{
		GasProvided:           2000,
		GasProvidedToChild:    300,
		GasProvidedToCallback: 50,
		GasUsedByParent:       400,
		GasUsedByChild:        200,
		GasUsedByCallback:     100,
		GasLockCost:           150,
		GasToLock:             150,

		TransferFromParentToChild: 7,

		ParentBalance:        1000,
		ChildBalance:         1000,
		TransferToThirdParty: 3,
		TransferToVault:      4,
		ESDTTokensToTransfer: 0,

		SuccessCallback: "myCallBack",
		ErrorCallback:   "myCallBack",
	}
}

func TestGasUsed_SingleContract(t *testing.T) {
	testConfig := makeTestConfig()

	_, err := test.BuildMockInstanceCallTest(t).
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
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasRemaining(testConfig.GasProvided-testConfig.GasUsedByParent).
				GasUsed(test.ParentAddress, testConfig.GasUsedByParent)
		})
	assert.Nil(t, err)
}

func TestGasUsed_SingleContract_BuiltinCall(t *testing.T) {
	testConfig := makeTestConfig()

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ExecOnDestCtxParentMock)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execOnDestCtx").
			WithArguments(test.ParentAddress, []byte("builtinClaim"), vmhost.One.Bytes()).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasRemaining(testConfig.GasProvided-testConfig.GasUsedByParent-gasUsedByBuiltinClaim).
				GasUsed(test.ParentAddress, testConfig.GasUsedByParent+gasUsedByBuiltinClaim).
				BalanceDelta(test.ParentAddress, amountToGiveByBuiltinClaim)
		})
	assert.Nil(t, err)
}

func TestGasUsed_SingleContract_BuiltinCallFail(t *testing.T) {
	testConfig := makeTestConfig()

	_, err := test.BuildMockInstanceCallTest(t).
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
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.ExecutionFailed().
				ReturnMessage("return value 1").
				HasRuntimeErrors("whatdidyoudo").
				GasRemaining(0)
		})
	assert.Nil(t, err)
}

func TestGasUsed_SingleContract_TransferFromChild(t *testing.T) {
	testConfig := makeTestConfig()

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithCodeMetadata([]byte{0, vmcommon.MetadataPayable}).
				WithConfig(testConfig).
				WithMethods(contracts.ExecOnDestCtxSingleCallParentMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.TransferEGLDToParent)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execOnDestCtxSingleCall").
			WithArguments(test.ChildAddress, []byte("transferEGLDToParent")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				Logs(vmcommon.LogEntry{
					Identifier: []byte("transferValueOnly"),
					Address:    test.ParentAddress,
					Topics:     [][]byte{big.NewInt(0).Bytes(), test.ChildAddress},
					Data:       vmcommon.FormatLogDataForCall("ExecuteOnDestContext", "transferEGLDToParent", [][]byte{}),
				},
					vmcommon.LogEntry{
						Identifier: []byte("transferValueOnly"),
						Address:    test.ChildAddress,
						Topics:     [][]byte{big.NewInt(testConfig.ChildBalance / 2).Bytes(), test.ParentAddress},
						Data:       vmcommon.FormatLogDataForCall("BackTransfer", "", [][]byte{}),
					})
		})
	assert.Nil(t, err)
}

func TestGasUsed_ExecuteOnDestChain(t *testing.T) {
	alphaAddress := test.MakeTestSCAddress("alpha")
	betaAddress := test.MakeTestSCAddress("beta")
	gammaAddress := test.MakeTestSCAddress("gamma")

	testConfig := &test.TestConfig{
		GasUsedByParent:    uint64(400),
		GasProvidedToChild: uint64(1000),
		GasProvided:        uint64(2000),
		GasUsedByChild:     uint64(200),
	}

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(alphaAddress).
				WithBalance(10).
				WithConfig(testConfig).
				WithMethods(contracts.ExecOnDestCtxSingleCallParentMock),
			test.CreateMockContract(betaAddress).
				WithBalance(0).
				WithConfig(testConfig).
				WithMethods(contracts.ExecOnDestCtxSingleCallParentMock),
			test.CreateMockContract(gammaAddress).
				WithBalance(0).
				WithConfig(testConfig).
				WithMethods(contracts.ReportOriginalCaller),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(alphaAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execOnDestCtxSingleCall").
			WithArguments(betaAddress, []byte("execOnDestCtxSingleCall"),
				gammaAddress, []byte("reportOriginalCaller")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				ReturnData(test.UserAddress, test.UserAddress, test.UserAddress)
		})

	assert.Nil(t, err)
}

func TestGasUsed_TwoContracts_ExecuteOnSameCtx(t *testing.T) {
	testConfig := makeTestConfig()

	for numCalls := uint64(0); numCalls < 3; numCalls++ {
		expectedGasRemaining := testConfig.GasProvided - testConfig.GasUsedByParent - testConfig.GasUsedByChild*numCalls
		numCallsBytes := big.NewInt(0).SetUint64(numCalls).Bytes()

		_, err := test.BuildMockInstanceCallTest(t).
			WithContracts(
				test.CreateMockContract(test.ParentAddress).
					WithBalance(testConfig.ParentBalance).
					WithConfig(testConfig).
					WithMethods(contracts.ExecOnSameCtxParentMock, contracts.ExecOnDestCtxParentMock, contracts.WasteGasParentMock),
				test.CreateMockContract(test.ChildAddress).
					WithBalance(testConfig.ChildBalance).
					WithConfig(testConfig).
					WithMethods(contracts.TransferEGLDToParent),
			).
			WithInput(test.CreateTestContractCallInputBuilder().
				WithRecipientAddr(test.ParentAddress).
				WithGasProvided(testConfig.GasProvided).
				WithFunction("execOnSameCtx").
				WithArguments(test.ChildAddress, []byte("transferEGLDToParent"), numCallsBytes).
				Build()).
			WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
				setZeroCodeCosts(host)
			}).
			AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
				verify.Ok().
					GasRemaining(expectedGasRemaining).
					GasUsed(test.ParentAddress, testConfig.GasUsedByParent+testConfig.GasUsedByChild*numCalls)
				if numCalls > 0 {
					verify.GasUsed(test.ChildAddress, 0)
				}
				if numCalls == 1 {
					verify.
						Logs(vmcommon.LogEntry{
							Identifier: []byte("transferValueOnly"),
							Address:    test.ParentAddress,
							Topics:     [][]byte{big.NewInt(0).Bytes(), test.ParentAddress},
							Data:       vmcommon.FormatLogDataForCall("ExecuteOnSameContext", "transferEGLDToParent", [][]byte{}),
						}, vmcommon.LogEntry{
							Identifier: []byte("transferValueOnly"),
							Address:    test.ParentAddress,
							Topics:     [][]byte{big.NewInt(testConfig.ChildBalance / 2).Bytes(), test.ParentAddress},
							Data:       vmcommon.FormatLogDataForCall("BackTransfer", "", [][]byte{}),
						})
				}
			})
		assert.Nil(t, err)
	}
}

func TestGasUsed_TwoContracts_ExecuteOnDestCtx(t *testing.T) {
	testConfig := makeTestConfig()

	for numCalls := uint64(0); numCalls < 3; numCalls++ {
		expectedGasRemaining := testConfig.GasProvided - testConfig.GasUsedByParent - testConfig.GasUsedByChild*numCalls
		numCallsBytes := big.NewInt(0).SetUint64(numCalls).Bytes()

		_, err := test.BuildMockInstanceCallTest(t).
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
			WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
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
		assert.Nil(t, err)
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

	_, err := test.BuildMockInstanceCallTest(t).
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
			WithArguments(betaAddress, []byte("wasteGas"), vmhost.One.Bytes(),
				gammaAddress, []byte("wasteGas"), vmhost.One.Bytes()).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasUsed(alphaAddress, testConfig.GasUsedByParent).
				GasUsed(betaAddress, testConfig.GasUsedByChild).
				GasUsed(gammaAddress, testConfig.GasUsedByChild).
				GasRemaining(testConfig.GasProvided - testConfig.GasUsedByParent - 2*testConfig.GasUsedByChild)
		})
	assert.Nil(t, err)
}

func TestGasUsed_ESDTTransfer_ThenExecuteCall_Success(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)
	esdtTransferGasCost := uint64(1)

	testConfig := makeTestConfig()
	testConfig.ESDTTokensToTransfer = 5

	_, err := test.BuildMockInstanceCallTest(t).
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
			verify.Ok().
				GasUsed(test.ParentAddress, testConfig.GasUsedByParent+esdtTransferGasCost).
				GasUsed(test.ChildAddress, testConfig.GasUsedByChild).
				GasRemaining(testConfig.GasProvided - esdtTransferGasCost - testConfig.GasUsedByParent - testConfig.GasUsedByChild).
				Logs(vmcommon.LogEntry{
					Identifier: []byte("ESDTTransfer"),
					Address:    test.ParentAddress,
					Topics:     [][]byte{test.ESDTTestTokenName, {}, big.NewInt(int64(testConfig.ESDTTokensToTransfer)).Bytes(), test.ChildAddress},
					Data: vmcommon.FormatLogDataForCall(
						"ExecuteOnDestContext",
						"ESDTTransfer",
						[][]byte{
							test.ESDTTestTokenName,
							big.NewInt(int64(testConfig.ESDTTokensToTransfer)).Bytes(),
							[]byte("wasteGas")}),
				})

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, testConfig.ESDTTokensToTransfer, childESDTBalance)
		})
	assert.Nil(t, err)
}

func TestGasUsed_ESDTTransfer_TwoChaninedCalls(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
	testConfig.GasProvided = 20000
	testConfig.GasProvidedToChild = 5000
	testConfig.ESDTTokensToTransfer = 5

	firstExpectedTransferFromParentToChild := txDataBuilder.NewBuilder()
	firstExpectedTransferFromParentToChild.
		TransferESDT(
			string(test.ESDTTestTokenName),
			int64(testConfig.ESDTTokensToTransfer)).
		Bytes([]byte("execESDTTransferAndCall")).
		Bytes(test.NephewAddress).
		Bytes([]byte("ESDTTransfer")).
		Bytes([]byte("wasteGas"))

	secondExpectedTransferFromParentToChild := txDataBuilder.NewBuilder()
	secondExpectedTransferFromParentToChild.
		TransferESDT(
			string(test.ESDTTestTokenName),
			int64(testConfig.ESDTTokensToTransfer)).
		Bytes([]byte("wasteGas"))

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ExecESDTTransferAndCallChild),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ExecESDTTransferAndCallChild),
			test.CreateMockContract(test.NephewAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.WasteGasChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execESDTTransferAndCall").
			WithArguments(test.ChildAddress, []byte("ESDTTransfer"), []byte("execESDTTransferAndCall"),
				test.NephewAddress, []byte("ESDTTransfer"), []byte("wasteGas")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				Transfers(
					test.CreateTransferEntry(test.ParentAddress, test.ChildAddress, 1).
						WithData(firstExpectedTransferFromParentToChild.ToBytes()).
						WithValue(big.NewInt(0)),
					test.CreateTransferEntry(test.ChildAddress, test.NephewAddress, 2).
						WithData(secondExpectedTransferFromParentToChild.ToBytes()).
						WithValue(big.NewInt(0)),
				)

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, uint64(0), childESDTBalance)

			nephewAccount := world.AcctMap.GetAccount(test.NephewAddress)
			nephewESDTBalance, _ := nephewAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, testConfig.ESDTTokensToTransfer, nephewESDTBalance)
		})
	assert.Nil(t, err)
}

func TestGasUsed_ESDTTransfer_ThenExecuteCall_Fail(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
	testConfig.ESDTTokensToTransfer = 5

	_, err := test.BuildMockInstanceCallTest(t).
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
			verify.ExecutionFailed().
				HasRuntimeErrors("forced fail").
				GasRemaining(0)

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, uint64(0), childESDTBalance)
		})
	assert.Nil(t, err)
}

func TestGasUsed_ESDTTransferFailed(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
	testConfig.ESDTTokensToTransfer = 2 * initialESDTTokenBalance

	_, err := test.BuildMockInstanceCallTest(t).
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
			verify.ExecutionFailed().
				HasRuntimeErrors("insufficient funds").
				GasRemaining(0)

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, uint64(0), childESDTBalance)
		})
	assert.Nil(t, err)
}

func TestMultipleTimes(t *testing.T) {
	for i := 0; i < 20; i++ {
		TestGasUsed_ESDTTransferFromParent_ChildBurnsAndThenFails(t)
	}
}

func TestGasUsed_ESDTTransferFromParent_ChildBurnsAndThenFails(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
	testConfig.ESDTTokensToTransfer = 10

	_, err := test.BuildMockInstanceCallTest(t).
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
			verify.ExecutionFailed().
				HasRuntimeErrors("forced fail")

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, uint64(0), childESDTBalance)
		})
	assert.Nil(t, err)
}

func TestGasUsed_AsyncCall_CrossShard_InitCall(t *testing.T) {
	testGasUsedAsyncCallCrossShardInitCall(t, false)
}

func TestGasUsed_LegacyAsyncCall_CrossShard_InitCall(t *testing.T) {
	testGasUsedAsyncCallCrossShardInitCall(t, true)
}

func testGasUsedAsyncCallCrossShardInitCall(t *testing.T, isLegacy bool) {
	testConfig := makeTestConfig()
	testConfig.GasProvided = 1000

	gasUsedByParent := testConfig.GasUsedByParent

	asyncCallData := txDataBuilder.NewBuilder()
	asyncCallData.Func(contracts.AsyncChildFunction)
	asyncCallData.Int64(testConfig.TransferToThirdParty)
	asyncCallData.Str(contracts.AsyncChildData)
	asyncCallData.Bytes([]byte{0})
	asyncChildArgs := asyncCallData.ToBytes()

	gasForAsyncCall := testConfig.GasProvided - gasUsedByParent - testConfig.GasToLock
	gasLocked := testConfig.GasToLock

	testConfig.IsLegacyAsync = isLegacy
	if !isLegacy {
		gasForAsyncCall -= testConfig.GasLockCost
		gasLocked += testConfig.GasLockCost
	}

	parentContract := test.CreateMockContractOnShard(test.ParentAddress, 0).
		WithBalance(testConfig.ParentBalance).
		WithConfig(testConfig).
		WithMethods(contracts.PerformAsyncCallParentMock, contracts.CallBackParentMock)

	expectedStorages := make([]testcommon.StoreEntry, 0)
	expectedStorages = append(expectedStorages,
		test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
		test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
		test.CreateStoreEntry(test.ParentAddress).WithKey(test.OriginalCallerParent).WithValue(test.UserAddress))

	expectedTransfers := make([]testcommon.TransferEntry, 0)
	expectedTransfers = append(expectedTransfers,
		test.CreateTransferEntry(test.ParentAddress, test.ThirdPartyAddress, 1).
			WithData([]byte("hello")).
			WithValue(big.NewInt(testConfig.TransferToThirdParty)),
		test.CreateTransferEntry(test.ParentAddress, test.ChildAddress, 2).
			WithData(asyncChildArgs).
			WithGasLimit(gasForAsyncCall).
			WithGasLocked(gasLocked).
			WithCallType(vm.AsynchronousCall).
			WithValue(big.NewInt(testConfig.TransferFromParentToChild)))

	// direct parent call
	_, err := test.BuildMockInstanceCallTest(t).
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
			expectedStorages = append(expectedStorages,
				test.CreateStoreEntry(test.ParentAddress).WithKey(
					host.Storage().GetVmProtectedPrefix(vmhost.AsyncDataPrefix)).IgnoreValue())
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasUsed(test.ParentAddress, gasUsedByParent).
				GasRemaining(0).
				ReturnData(test.ParentFinishA, test.ParentFinishB).
				Storage(expectedStorages...).
				Transfers(expectedTransfers...)
		})
	assert.Nil(t, err)
}

func TestGasUsed_AsyncCall_CrossShard_ExecuteCall(t *testing.T) {
	testConfig := makeTestConfig()
	gasForAsyncCall := testConfig.GasProvided - testConfig.GasUsedByParent - testConfig.GasLockCost

	childAsyncReturnData := [][]byte{{0}, []byte("thirdparty"), []byte("vault")}

	// async cross-shard parent -> child
	_, err := test.BuildMockInstanceCallTest(t).
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
			WithAsyncArguments(
				&vmcommon.AsyncArguments{CallID: []byte{0}, CallerCallID: []byte{0}},
			).
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
			verify.Ok().
				GasUsed(test.ChildAddress, testConfig.GasUsedByChild).
				GasRemaining(gasForAsyncCall-testConfig.GasUsedByChild).
				ReturnData(childAsyncReturnData...).
				Transfers(
					test.CreateTransferEntry(test.ChildAddress, test.ThirdPartyAddress, 1).
						WithData([]byte(contracts.AsyncChildData)).
						WithValue(big.NewInt(testConfig.TransferToThirdParty)),
					test.CreateTransferEntry(test.ChildAddress, test.VaultAddress, 2).
						WithData([]byte{}).
						WithValue(big.NewInt(testConfig.TransferToVault)),
				)
		})
	assert.Nil(t, err)
}

func TestGasUsed_AsyncCall_CrossShard_ExecuteCall_WithTransfer(t *testing.T) {
	testConfig := makeTestConfig()
	gasUsedByParent := testConfig.GasUsedByParent
	gasForAsyncCall := testConfig.GasProvided - gasUsedByParent - testConfig.GasLockCost

	// async cross-shard parent -> child
	_, err := test.BuildMockInstanceCallTest(t).
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
			WithAsyncArguments(
				&vmcommon.AsyncArguments{CallID: []byte{0}, CallerCallID: []byte{0}},
			).
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
				GasUsed(test.ChildAddress, testConfig.GasUsedByChild).
				GasRemaining(gasForAsyncCall - testConfig.GasUsedByChild).
				ReturnData().
				Transfers(
					test.CreateTransferEntry(test.ChildAddress, test.ParentAddress, 1).
						WithGasLimit(0).
						WithCallType(vm.DirectCall).
						WithValue(big.NewInt(testConfig.TransferToThirdParty)),
				)
		})
	assert.Nil(t, err)
}

func TestGasUsed_AsyncCall_CrossShard_CallBack_LegacyAsyncCall(t *testing.T) {
	testGasUsedAsyncCallCrossShardCallBack(t, true)
}

func TestGasUsed_AsyncCall_CrossShard_CallBack_AsyncCall(t *testing.T) {
	testGasUsedAsyncCallCrossShardCallBack(t, false)
}

func testGasUsedAsyncCallCrossShardCallBack(t *testing.T, isLegacy bool) {
	testConfig := makeTestConfig()
	testConfig.GasProvided = 1000

	gasUsedByParent := testConfig.GasUsedByParent
	gasUsedByChild := testConfig.GasUsedByChild
	gasForAsyncCall := testConfig.GasProvided - gasUsedByParent - testConfig.GasLockCost

	testConfig.IsLegacyAsync = isLegacy

	parentContract := test.CreateMockContractOnShard(test.ParentAddress, 0).
		WithBalance(testConfig.ParentBalance).
		WithConfig(testConfig).
		WithMethods(contracts.PerformAsyncCallParentMock, contracts.CallBackParentMock)

	asyncArguments := &vmcommon.AsyncArguments{
		CallID:                       []byte{1, 2, 3},
		CallerCallID:                 []byte{3, 2, 1},
		CallbackAsyncInitiatorCallID: []byte{4, 5, 6},
		GasAccumulated:               1,
	}
	arguments := [][]byte{[]byte("thirdparty"), []byte("vault"), {0}}

	// async cross shard callback child -> parent
	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(parentContract).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithCallerAddr(test.ChildAddress).
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(gasForAsyncCall - gasUsedByChild + testConfig.GasLockCost).
			WithFunction("callBack").
			WithAsyncArguments(asyncArguments).
			WithArguments(arguments...).
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

			// TODO factor this setup out if necessary for other tests

			// The instance started below will be cached on the InstanceMockBuilder and reused by doRunSmartContractCall().
			// This is necessary for gas usage metering during Save() below.
			// Note that the InstanceMockBuilder uses the address of the contract as
			// if it were its bytecode, hence StartWasmerInstance() receives an
			// address as its first argument.
			err := host.Runtime().StartWasmerInstance(test.ParentAddress, testConfig.GasUsedByParent, false)
			assert.Nil(t, err)

			fakeInput := &host.Runtime().GetVMInput().VMInput
			fakeInput.GasProvided = 1000
			host.Metering().InitStateFromContractCallInput(fakeInput)

			err = contracts.RegisterAsyncCallToChild(host, testConfig, arguments)
			assert.Nil(t, err)

			host.Async().SetCallID(asyncArguments.CallbackAsyncInitiatorCallID)
			host.Async().SetCallIDForCallInGroup(0, 0, asyncArguments.CallerCallID)
			err = host.Async().Save()
			assert.Nil(t, err)

			for _, account := range host.Output().GetVMOutput().OutputAccounts {
				for _, storageUpdate := range account.StorageUpdates {
					(accountHandler.(*worldmock.Account)).Storage[string(storageUpdate.Offset)] = storageUpdate.Data
				}
			}
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasRemaining(testConfig.GasProvided - gasUsedByParent - gasUsedByChild - testConfig.GasUsedByCallback).
				ReturnData([]byte("succ"))
		})
	assert.Nil(t, err)
}

func TestGasUsed_LegacyAsyncCall_InShard_BuiltinCall(t *testing.T) {
	// all gas for builtin call is consummed on caller
	inShardBuiltinCall(t, true)
}

func TestGasUsed_AsyncCall_InShard_BuiltinCall(t *testing.T) {
	// all gas for builtin call is consummed on caller
	inShardBuiltinCall(t, false)
}

func inShardBuiltinCall(t *testing.T, legacy bool) {
	testConfig := makeTestConfig()
	testConfig.GasProvided = 1000

	expectedGasUsedByParent := testConfig.GasUsedByParent + testConfig.GasUsedByCallback + gasUsedByBuiltinClaim
	expectedGasUsedByChild := uint64(0)

	testConfig.IsLegacyAsync = legacy

	_, err := test.BuildMockInstanceCallTest(t).
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
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
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
	assert.Nil(t, err)
}

func TestGasUsed_BuiltinCallFail_LegacyAsyncCall(t *testing.T) {
	testGasUsedBuiltinCallFail(t, true)
}

func TestGasUsed_BuiltinCallFail_AsyncCall(t *testing.T) {
	testGasUsedBuiltinCallFail(t, false)
}

func testGasUsedBuiltinCallFail(t *testing.T, isLegacy bool) {
	testConfig := makeTestConfig()
	testConfig.GasProvided = 1000

	gasProvidedForBuiltinCall := testConfig.GasProvided - testConfig.GasUsedByParent - testConfig.GasLockCost
	expectedGasUsedByParent := testConfig.GasUsedByParent + gasProvidedForBuiltinCall + testConfig.GasUsedByCallback

	testConfig.IsLegacyAsync = isLegacy
	if !isLegacy {
		expectedGasUsedByParent -= testConfig.GasLockCost
	}

	_, err := test.BuildMockInstanceCallTest(t).
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
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
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
	assert.Nil(t, err)
}

func TestGasUsed_LegacyAsyncCall_CrossShard_BuiltinCall(t *testing.T) {
	testConfig := makeTestConfig()
	testConfig.GasProvided = 1000

	expectedGasUsedByParent := testConfig.GasUsedByParent + gasUsedByBuiltinClaim

	testConfig.IsLegacyAsync = true

	_, err := test.BuildMockInstanceCallTest(t).
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
			WithArguments(test.UserAddress, []byte("sendMessage"), []byte{2}, vmhost.One.Bytes()).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			world.SelfShardID = 1
			world.AcctMap.CreateAccount(test.UserAddress, world)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasRemaining(0).
				GasUsed(test.ParentAddress, expectedGasUsedByParent).
				Transfers(
					test.CreateTransferEntry(test.ParentAddress, test.UserAddress, 1).
						WithData([]byte("message")).
						WithGasLimit(480).
						WithValue(big.NewInt(0)),
					test.CreateTransferEntry(test.ParentAddress, test.UserAddress, 2).
						WithData([]byte("message")).
						WithGasLimit(0).
						WithValue(big.NewInt(0)),
				)
		})
	assert.Nil(t, err)
}

func TestGasUsed_AsyncCall_BuiltinMultiContractChainCall(t *testing.T) {
	// TODO matei-p enable this test for R2
	t.Skip()

	testConfig := makeTestConfig()
	testConfig.TransferFromChildToParent = 5

	expectedGasUsedByParent := testConfig.GasUsedByParent + testConfig.GasUsedByCallback
	expectedGasUsedByChild :=
		testConfig.GasUsedByParent /* due to ForwardAsyncCallParentBuiltinMock */ +
			gasUsedByBuiltinClaim +
			testConfig.GasUsedByCallback /* due to CallBackParentBuiltinMock */

	_, err := test.BuildMockInstanceCallTest(t).
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
			WithArguments(test.ChildAddress, []byte("forwardAsyncCall"), []byte("builtinClaim"), vmhost.One.Bytes()).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			world.AcctMap.CreateAccount(test.UserAddress, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
			createMockBuiltinFunctions(t, host, world)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasUsed(test.ParentAddress, expectedGasUsedByParent).
				GasUsed(test.ChildAddress, expectedGasUsedByChild).
				GasRemaining(testConfig.GasProvided - expectedGasUsedByParent - expectedGasUsedByChild)
		})
	assert.Nil(t, err)
}

func TestGasUsed_AsyncCall_ChildFails(t *testing.T) {
	testGasUsedAsyncCallChildFails(t, false)
}

func TestGasUsed_LegacyAsyncCall_ChildFails(t *testing.T) {
	testGasUsedAsyncCallChildFails(t, true)
}

func testGasUsedAsyncCallChildFails(t *testing.T, isLegacy bool) {
	testConfig := makeTestConfig()
	testConfig.GasProvided = 1000

	testConfig.IsLegacyAsync = isLegacy
	expectedGasUsedByParent := testConfig.GasProvided - testConfig.GasToLock + testConfig.GasUsedByCallback

	if !isLegacy {
		expectedGasUsedByParent -= testConfig.GasLockCost
	}

	_, err := test.BuildMockInstanceCallTest(t).
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
			WithArguments(vmhost.One.Bytes()).
			WithCurrentTxHash([]byte("txhash")).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
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
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.CallbackKey).WithValue(test.CallbackData),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.OriginalCallerParent).WithValue(test.UserAddress),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.OriginalCallerCallback).WithValue(test.UserAddress),
				).
				Transfers(
					test.CreateTransferEntry(test.ParentAddress, test.VaultAddress, 2).
						WithData([]byte("child error")).
						WithValue(big.NewInt(testConfig.TransferToVault)),
					test.CreateTransferEntry(test.ParentAddress, test.ThirdPartyAddress, 1).
						WithData([]byte("hello")).
						WithValue(big.NewInt(testConfig.TransferToThirdParty)),
				)
		})
	assert.Nil(t, err)
}

func TestGasUsed_AsyncCall_CallBackFails(t *testing.T) {
	testGasUsedAsyncCallCallBackFails(t, false)
}

func TestGasUsed_LegacyAsyncCall_CallBackFails(t *testing.T) {
	testGasUsedAsyncCallCallBackFails(t, true)
}

func testGasUsedAsyncCallCallBackFails(t *testing.T, isLegacy bool) {
	testConfig := makeTestConfig()

	var expectedGasUsedByParent uint64
	var expectedRemainingGas uint64
	expectedGasUsedByChild := testConfig.GasUsedByChild

	testConfig.IsLegacyAsync = isLegacy
	if !isLegacy {
		expectedGasUsedByParent =
			testConfig.GasUsedByParent +
				testConfig.GasProvidedToChild +
				testConfig.GasLockCost +
				testConfig.GasToLock -
				testConfig.GasUsedByChild
		expectedRemainingGas =
			testConfig.GasProvided -
				(expectedGasUsedByParent + testConfig.GasUsedByChild)
	} else {
		expectedGasUsedByParent =
			testConfig.GasProvided - testConfig.GasUsedByChild
		expectedRemainingGas = 0
	}

	_, err := test.BuildMockInstanceCallTest(t).
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
			verify.Ok().
				HasRuntimeErrors("callBack error").
				BalanceDelta(test.ParentAddress, -(2*testConfig.TransferToThirdParty+testConfig.TransferToVault)).
				BalanceDelta(test.ThirdPartyAddress, 2*testConfig.TransferToThirdParty).
				GasUsed(test.ParentAddress, expectedGasUsedByParent).
				GasUsed(test.ChildAddress, expectedGasUsedByChild).
				GasRemaining(expectedRemainingGas).
				ReturnData(test.ParentFinishA, test.ParentFinishB, []byte{3}, []byte("thirdparty"), []byte("vault")).
				Storage(
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyA).WithValue(test.ParentDataA),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.ParentKeyB).WithValue(test.ParentDataB),
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.OriginalCallerParent).WithValue(test.UserAddress),
					test.CreateStoreEntry(test.ChildAddress).WithKey(test.ChildKey).WithValue(test.ChildData),
					test.CreateStoreEntry(test.ChildAddress).WithKey(test.OriginalCallerChild).WithValue(test.UserAddress),
				).
				Transfers(
					test.CreateTransferEntry(test.ParentAddress, test.ThirdPartyAddress, 1).
						WithData([]byte("hello")).
						WithValue(big.NewInt(testConfig.TransferToThirdParty)),
					test.CreateTransferEntry(test.ChildAddress, test.ThirdPartyAddress, 2).
						WithData([]byte(" there")).
						WithValue(big.NewInt(testConfig.TransferToThirdParty)),
					test.CreateTransferEntry(test.ChildAddress, test.VaultAddress, 3).
						WithData([]byte{}).
						WithValue(big.NewInt(testConfig.TransferToVault)),
				)
		})
	assert.Nil(t, err)
}

func TestGasUsed_AsyncCall_Recursive(t *testing.T) {
	// TODO reenable test correct assertions after contracts are allowed to call themselves
	// repeatedly with async calls (see restriction in asyncContext.addAsyncCall())

	testConfig := makeTestConfig()
	testConfig.RecursiveChildCalls = 3

	// expectedGasUsedByParent := testConfig.GasUsedByParent + testConfig.GasUsedByCallback
	// expectedGasUsedByChild := uint64(testConfig.RecursiveChildCalls)*testConfig.GasProvidedToChild +
	// 	uint64(testConfig.RecursiveChildCalls-1)*testConfig.GasUsedByCallback

	_, err := test.BuildMockInstanceCallTest(t).
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
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				HasRuntimeErrors(vmhost.ErrExecutionFailed.Error())
			// BalanceDelta(test.ParentAddress, -testConfig.TransferFromParentToChild).
			// GasUsed(test.ParentAddress, expectedGasUsedByParent).
			// GasUsed(test.ChildAddress, expectedGasUsedByChild).
			// GasRemaining(testConfig.GasProvided-expectedGasUsedByParent-expectedGasUsedByChild).
			// BalanceDelta(test.ChildAddress, testConfig.TransferFromParentToChild).
			// ReturnData(big.NewInt(2).Bytes(), big.NewInt(1).Bytes(), big.NewInt(0).Bytes())
		})
	assert.Nil(t, err)
}

func TestGasUsed_AsyncCall_MultiChild(t *testing.T) {
	testConfig := makeTestConfig()
	testConfig.ChildCalls = 2

	expectedGasUsedByParent := testConfig.GasUsedByParent + 2*testConfig.GasUsedByCallback
	expectedGasUsedByChild := uint64(testConfig.ChildCalls) * testConfig.GasUsedByChild

	_, err := test.BuildMockInstanceCallTest(t).
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
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
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
				ReturnData(big.NewInt(0).Bytes(), big.NewInt(1).Bytes()).
				Logs(vmcommon.LogEntry{
					Identifier: []byte("transferValueOnly"),
					Address:    test.ParentAddress,
					Topics:     [][]byte{big.NewInt(testConfig.TransferFromParentToChild).Bytes(), test.ChildAddress},
					Data:       vmcommon.FormatLogDataForCall("AsyncCall", "recursiveAsyncCall", [][]byte{{1}, {}}),
				}, vmcommon.LogEntry{
					Identifier: []byte("transferValueOnly"),
					Address:    test.ChildAddress,
					Topics:     [][]byte{big.NewInt(0).Bytes(), test.ParentAddress},
					Data:       vmcommon.FormatLogDataForCall("AsyncCallback", "callBack", [][]byte{{0}, {}}),
				}, vmcommon.LogEntry{
					Identifier: []byte("transferValueOnly"),
					Address:    test.ParentAddress,
					Topics:     [][]byte{big.NewInt(testConfig.TransferFromParentToChild).Bytes(), test.ChildAddress},
					Data:       vmcommon.FormatLogDataForCall("AsyncCall", "recursiveAsyncCall", [][]byte{{1}, {1}}),
				}, vmcommon.LogEntry{
					Identifier: []byte("transferValueOnly"),
					Address:    test.ChildAddress,
					Topics:     [][]byte{big.NewInt(0).Bytes(), test.ParentAddress},
					Data:       vmcommon.FormatLogDataForCall("AsyncCallback", "callBack", [][]byte{{0}, {1}}),
				})
		})
	assert.Nil(t, err)
}

func TestGasUsed_ESDTTransfer_ThenExecuteAsyncCall_Success(t *testing.T) {
	testGasUsedESDTTransferThenExecuteAsyncCallSuccess(t, false)
}

func TestGasUsed_Legacy_ESDTTransfer_ThenExecuteAsyncCall_Success(t *testing.T) {
	testGasUsedESDTTransferThenExecuteAsyncCallSuccess(t, true)
}

func testGasUsedESDTTransferThenExecuteAsyncCallSuccess(t *testing.T, isLegacy bool) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
	testConfig.ESDTTokensToTransfer = 5

	asyncCallType := NewAsyncCallType
	if isLegacy {
		asyncCallType = LegacyAsyncCallType
	}

	expectedTransferFromParentToChild := txDataBuilder.NewBuilder()
	expectedTransferFromParentToChild.TransferESDT(string(test.ESDTTestTokenName), int64(testConfig.ESDTTokensToTransfer))

	expectedTransfers := make([]testcommon.TransferEntry, 0)
	expectedTransfers = append(expectedTransfers,
		test.CreateTransferEntry(test.ParentAddress, test.ChildAddress, 1).
			WithData(expectedTransferFromParentToChild.ToBytes()).
			WithGasLimit(0).
			WithGasLocked(0).
			WithCallType(vm.AsynchronousCall).
			WithValue(big.NewInt(0)))

	_, err := test.BuildMockInstanceCallTest(t).
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
			WithArguments(test.ChildAddress, []byte("ESDTTransfer"), []byte("wasteGas"), asyncCallType).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				Transfers(expectedTransfers...)

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, testConfig.ESDTTokensToTransfer, childESDTBalance)
		})
	assert.Nil(t, err)
}

/*
ParentAddress.execESDTTransferAndAsyncCall -> ChildAddress.wasteGasOnNephew (with async with ESDTTransfer)
ChildAddress.wasteGasOnNephew -> NephewAddress.wasteGas
-> ParentAddress.callBack
ParentAddress.callBack -> ChildAddress.wasteGas
*/
func TestGasUsed_ESDTTransfer_ThenExecuteAsyncCall_ThenExecuteOnDest(t *testing.T) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
	testConfig.ESDTTokensToTransfer = 5

	expectedTransferFromParentToChild := txDataBuilder.NewBuilder()
	expectedTransferFromParentToChild.TransferESDT(string(test.ESDTTestTokenName), int64(testConfig.ESDTTokensToTransfer))

	expectedTransfers := make([]testcommon.TransferEntry, 0)
	expectedTransfers = append(expectedTransfers,
		test.CreateTransferEntry(test.ParentAddress, test.ChildAddress, 1).
			WithData(expectedTransferFromParentToChild.ToBytes()).
			WithGasLimit(0).
			WithGasLocked(0).
			WithCallType(vm.AsynchronousCall).
			WithValue(big.NewInt(0)))

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(
					contracts.ExecESDTTransferAndAsyncCallChild,
					contracts.WasteGasParentMock,
					contracts.LocalCallAnotherContract("callBack", test.ChildAddress, "wasteGas")),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(
					contracts.LocalCallAnotherContract("wasteGasOnNewphew", test.NephewAddress, "wasteGas"),
					contracts.WasteGasChildMock),
			test.CreateMockContract(test.NephewAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.WasteGasChildMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execESDTTransferAndAsyncCall").
			WithArguments(test.ChildAddress, []byte("ESDTTransfer"), []byte("wasteGasOnNewphew"), NewAsyncCallType).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				Transfers(expectedTransfers...).
				Logs(
					vmcommon.LogEntry{
						Identifier: []byte("ESDTTransfer"),
						Address:    test.ParentAddress,
						Topics:     [][]byte{test.ESDTTestTokenName, {}, big.NewInt(int64(testConfig.ESDTTokensToTransfer)).Bytes(), test.ChildAddress},
						Data: vmcommon.FormatLogDataForCall(
							"AsyncCall",
							"ESDTTransfer",
							[][]byte{
								test.ESDTTestTokenName,
								big.NewInt(int64(testConfig.ESDTTokensToTransfer)).Bytes(),
								[]byte("wasteGasOnNewphew"),
								{1}}),
					},
					vmcommon.LogEntry{
						Identifier: []byte("transferValueOnly"),
						Address:    test.ChildAddress,
						Topics:     [][]byte{{}, test.NephewAddress},
						Data:       vmcommon.FormatLogDataForCall("ExecuteOnDestContext", "wasteGas", [][]byte{}),
					},
					vmcommon.LogEntry{
						Identifier: []byte("transferValueOnly"),
						Address:    test.ChildAddress,
						Topics:     [][]byte{{}, test.ParentAddress},
						Data:       vmcommon.FormatLogDataForCall("AsyncCallback", "callBack", [][]byte{{0}, test.UserAddress}),
					},
					vmcommon.LogEntry{
						Identifier: []byte("transferValueOnly"),
						Address:    test.ParentAddress,
						Topics:     [][]byte{{}, test.ChildAddress},
						Data:       vmcommon.FormatLogDataForCall("ExecuteOnDestContext", "wasteGas", [][]byte{}),
					})

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, testConfig.ESDTTokensToTransfer, childESDTBalance)
		})
	assert.Nil(t, err)
}

func TestGasUsed_ESDTTransfer_ThenExecuteAsyncCall_ChildFails(t *testing.T) {
	testGasUsedESDTTransferThenExecuteAsyncCallChildFails(t, false)
}

func TestGasUsed_Legacy_ESDTTransfer_ThenExecuteAsyncCall_ChildFails(t *testing.T) {
	testGasUsedESDTTransferThenExecuteAsyncCallChildFails(t, true)
}

func testGasUsedESDTTransferThenExecuteAsyncCallChildFails(t *testing.T, isLegacy bool) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
	testConfig.ESDTTokensToTransfer = 5

	var asyncCallType []byte
	var expectedGasUsedByParent uint64
	var expectedGasRemaining uint64
	if !isLegacy {
		asyncCallType = NewAsyncCallType
		expectedGasUsedByParent =
			testConfig.GasUsedByParent +
				testConfig.GasProvidedToChild +
				testConfig.GasLockCost +
				testConfig.GasToLock -
				testConfig.GasUsedByChild
		expectedGasRemaining =
			testConfig.GasProvided -
				expectedGasUsedByParent
	} else {
		asyncCallType = LegacyAsyncCallType
		expectedGasRemaining = uint64(50)
		expectedGasUsedByParent = testConfig.GasProvided - expectedGasRemaining
	}

	expectedTransferFromParentToChild := txDataBuilder.NewBuilder()
	expectedTransferFromParentToChild.TransferESDT(string(test.ESDTTestTokenName), int64(testConfig.ESDTTokensToTransfer))

	expectedTransfers := make([]testcommon.TransferEntry, 0)
	expectedTransfers = append(expectedTransfers,
		test.CreateTransferEntry(test.ParentAddress, test.ChildAddress, 1).
			WithData(expectedTransferFromParentToChild.ToBytes()).
			WithGasLimit(0).
			WithGasLocked(0).
			WithCallType(vm.AsynchronousCall).
			WithValue(big.NewInt(0)))

	_, err := test.BuildMockInstanceCallTest(t).
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
			WithArguments(test.ChildAddress, []byte("ESDTTransfer"), []byte("fail"), asyncCallType).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasRemaining(expectedGasRemaining).
				GasUsed(test.ParentAddress, expectedGasUsedByParent).
				GasUsed(test.ChildAddress, 0).
				Transfers(expectedTransfers...)

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, uint64(0), childESDTBalance)
		})
	assert.Nil(t, err)
}

func TestGasUsed_ESDTTransfer_ThenExecuteAsyncCall_CallbackFails(t *testing.T) {
	testGasUsedESDTTransferThenExecuteAsyncCallCallbackFails(t, false)
}

func TestGasUsed_Legacy_ESDTTransfer_ThenExecuteAsyncCall_CallbackFails(t *testing.T) {
	testGasUsedESDTTransferThenExecuteAsyncCallCallbackFails(t, true)
}

func testGasUsedESDTTransferThenExecuteAsyncCallCallbackFails(t *testing.T, isLegacy bool) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
	testConfig.ESDTTokensToTransfer = 5

	var expectedGasUsedByParent uint64
	var expectedRemainingGas uint64

	asyncCallType := LegacyAsyncCallType
	if !isLegacy {
		asyncCallType = NewAsyncCallType
		expectedGasUsedByParent =
			testConfig.GasUsedByParent +
				testConfig.GasProvidedToChild +
				testConfig.GasLockCost +
				testConfig.GasToLock -
				testConfig.GasUsedByChild
		expectedRemainingGas =
			testConfig.GasProvided -
				(expectedGasUsedByParent + testConfig.GasUsedByChild)
	} else {
		expectedGasUsedByParent =
			testConfig.GasProvided - testConfig.GasUsedByChild
		expectedRemainingGas = 0
	}

	expectedTransferFromParentToChild := txDataBuilder.NewBuilder()
	expectedTransferFromParentToChild.TransferESDT(string(test.ESDTTestTokenName), int64(testConfig.ESDTTokensToTransfer))

	expectedTransfers := make([]testcommon.TransferEntry, 0)
	expectedTransfers = append(expectedTransfers,
		test.CreateTransferEntry(test.ParentAddress, test.ChildAddress, 1).
			WithData(expectedTransferFromParentToChild.ToBytes()).
			WithGasLimit(0).
			WithGasLocked(0).
			WithCallType(vm.AsynchronousCall).
			WithValue(big.NewInt(0)))

	_, err := test.BuildMockInstanceCallTest(t).
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
			WithArguments(test.ChildAddress, []byte("ESDTTransfer"), []byte("wasteGas"), asyncCallType).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasUsed(test.ParentAddress, expectedGasUsedByParent).
				GasUsed(test.ChildAddress, testConfig.GasUsedByChild).
				GasRemaining(expectedRemainingGas).
				Transfers(expectedTransfers...)

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, testConfig.ESDTTokensToTransfer, childESDTBalance)
		})
	assert.Nil(t, err)
}

func TestGasUsed_ESDTTransferInCallback(t *testing.T) {
	testGasUsedESDTTransferInCallback(t, false, 2)
}

func TestGasUsed_Legacy_ESDTTransferInCallback(t *testing.T) {
	testGasUsedESDTTransferInCallback(t, true, 2)
}

func testGasUsedESDTTransferInCallback(t *testing.T, isLegacy bool, numOfTransfersInChild int) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
	testConfig.GasProvidedToChild += 20000
	testConfig.GasProvided += 40000
	testConfig.ESDTTokensToTransfer = 5
	testConfig.CallbackESDTTokensToTransfer = 2

	asyncCallType := LegacyAsyncCallType
	if !isLegacy {
		asyncCallType = NewAsyncCallType
	}

	expectedTransferFromParentToChild := txDataBuilder.NewBuilder()
	expectedTransferFromParentToChild.TransferESDT(string(test.ESDTTestTokenName), int64(testConfig.ESDTTokensToTransfer))

	expectedTransferFromChildToParent := txDataBuilder.NewBuilder()
	expectedTransferFromChildToParent.TransferESDT(string(test.ESDTTestTokenName), int64(testConfig.CallbackESDTTokensToTransfer))

	expectedTransfers := make([]testcommon.TransferEntry, 0)
	expectedLogs := make([]vmcommon.LogEntry, 0)

	expectedTransfers = append(expectedTransfers,
		test.CreateTransferEntry(test.ParentAddress, test.ChildAddress, 1).
			WithData(expectedTransferFromParentToChild.ToBytes()).
			WithGasLimit(0).
			WithGasLocked(0).
			WithCallType(vm.AsynchronousCall).
			WithValue(big.NewInt(0)))
	expectedLogs = append(expectedLogs,
		vmcommon.LogEntry{
			Identifier: []byte("ESDTTransfer"),
			Address:    test.ParentAddress,
			Topics:     [][]byte{test.ESDTTestTokenName, {}, big.NewInt(int64(testConfig.ESDTTokensToTransfer)).Bytes(), test.ChildAddress},
			Data: vmcommon.FormatLogDataForCall(
				"AsyncCall",
				"ESDTTransfer",
				[][]byte{
					test.ESDTTestTokenName,
					big.NewInt(int64(testConfig.ESDTTokensToTransfer)).Bytes(),
					[]byte("transferESDTToParent"),
					big.NewInt(int64(testConfig.CallbackESDTTokensToTransfer)).Bytes()}),
		})

	for transfer := 0; transfer < numOfTransfersInChild; transfer++ {
		expectedTransfers = append(expectedTransfers,
			test.CreateTransferEntry(test.ChildAddress, test.ParentAddress, uint32(2+transfer)).
				WithData(expectedTransferFromChildToParent.ToBytes()).
				WithGasLimit(0).
				WithGasLocked(0).
				WithCallType(vm.DirectCall).
				WithValue(big.NewInt(0)))
		expectedLogs = append(expectedLogs,
			vmcommon.LogEntry{
				Identifier: []byte("ESDTTransfer"),
				Address:    test.ChildAddress,
				Topics:     [][]byte{test.ESDTTestTokenName, {}, big.NewInt(int64(testConfig.CallbackESDTTokensToTransfer)).Bytes(), test.ParentAddress},
				Data: vmcommon.FormatLogDataForCall(
					"BackTransfer",
					"ESDTTransfer",
					[][]byte{
						test.ESDTTestTokenName,
						big.NewInt(int64(testConfig.CallbackESDTTokensToTransfer)).Bytes()}),
			})
	}

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithCodeMetadata([]byte{0, vmcommon.MetadataPayable}).
				WithMethods(contracts.ExecESDTTransferAndAsyncCallChild, contracts.SimpleCallbackMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ESDTTransferToParentMockNoReturnData),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execESDTTransferAndAsyncCall").
			WithArguments(test.ChildAddress, []byte("ESDTTransfer"), []byte("transferESDTToParent"), asyncCallType, []byte{byte(numOfTransfersInChild)}).
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
				Transfers(expectedTransfers...).
				ReturnData(
					[]byte(test.ESDTTestTokenName),
					big.NewInt(int64(testConfig.CallbackESDTTokensToTransfer)).Bytes()).
				Logs(expectedLogs...)

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer+uint64(numOfTransfersInChild)*testConfig.CallbackESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, testConfig.ESDTTokensToTransfer-uint64(numOfTransfersInChild)*testConfig.CallbackESDTTokensToTransfer, childESDTBalance)
		})
	assert.Nil(t, err)
}

func TestGasUsed_ESDTTransferInCallbackAndTryNewAsync(t *testing.T) {
	testGasUsedESDTTransferInCallbackAndTryNewAsync(t, false)
}

func TestGasUsed_Legacy_ESDTTransferInCallbackAndTryNewAsync(t *testing.T) {
	testGasUsedESDTTransferInCallbackAndTryNewAsync(t, true)
}

func testGasUsedESDTTransferInCallbackAndTryNewAsync(t *testing.T, isLegacy bool) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
	testConfig.GasProvided += 4000
	testConfig.GasProvidedToChild += 2000

	testConfig.ESDTTokensToTransfer = 5
	testConfig.CallbackESDTTokensToTransfer = 1

	var expectedGasUsedByParent uint64
	var expectedRemainingGas uint64

	asyncCallType := LegacyAsyncCallType
	if !isLegacy {
		asyncCallType = NewAsyncCallType
		expectedGasUsedByParent =
			testConfig.GasUsedByParent +
				testConfig.GasProvidedToChild +
				testConfig.GasLockCost +
				testConfig.GasToLock -
				(testConfig.GasUsedByChild + 1)
		expectedRemainingGas =
			testConfig.GasProvided -
				(expectedGasUsedByParent + testConfig.GasUsedByChild + 1)
	} else {
		expectedGasUsedByParent =
			testConfig.GasProvided - (testConfig.GasUsedByChild + 1)
		expectedRemainingGas = 0
	}

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithCodeMetadata([]byte{0, vmcommon.MetadataPayable}).
				WithMethods(contracts.ExecESDTTransferAndAsyncCallChild, contracts.SimpleCallbackMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(testConfig.ChildBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ESDTTransferToParentAndNewAsyncFromCallbackMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execESDTTransferAndAsyncCall").
			WithArguments(test.ChildAddress, []byte("ESDTTransfer"), []byte("transferESDTToParent"), asyncCallType).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				GasUsed(test.ParentAddress, expectedGasUsedByParent).
				GasUsed(test.ChildAddress, testConfig.GasUsedByChild+1 /* ESDTTransfer */).
				GasRemaining(expectedRemainingGas)

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer+testConfig.CallbackESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, testConfig.ESDTTokensToTransfer-testConfig.CallbackESDTTokensToTransfer, childESDTBalance)
		})
	assert.Nil(t, err)
}

func TestGasUsed_ESDTTransfer_CallbackFail(t *testing.T) {
	testGasUsedESDTTransferCallbackFail(t, false)
}

func TestGasUsed_Legacy_ESDTTransfer_CallbackFail(t *testing.T) {
	testGasUsedESDTTransferCallbackFail(t, true)
}

func testGasUsedESDTTransferCallbackFail(t *testing.T, isLegacy bool) {
	var parentAccount *worldmock.Account
	initialESDTTokenBalance := uint64(100)

	testConfig := makeTestConfig()
	testConfig.ESDTTokensToTransfer = 5
	testConfig.CallbackESDTTokensToTransfer = 2

	asyncCallType := LegacyAsyncCallType
	if !isLegacy {
		asyncCallType = NewAsyncCallType
		testConfig.GasProvided += 4000
		testConfig.GasProvidedToChild += 2000
	}

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithCodeMetadata([]byte{0, vmcommon.MetadataPayable}).
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
			WithArguments(test.ChildAddress, []byte("ESDTTransfer"), []byte("transferESDTToParent"), asyncCallType).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			parentAccount = world.AcctMap.GetAccount(test.ParentAddress)
			_ = parentAccount.SetTokenBalanceUint64(test.ESDTTestTokenName, 0, initialESDTTokenBalance)
			createMockBuiltinFunctions(t, host, world)
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				HasRuntimeErrors("callback failed intentionally")

			parentESDTBalance, _ := parentAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, initialESDTTokenBalance-testConfig.ESDTTokensToTransfer+testConfig.CallbackESDTTokensToTransfer, parentESDTBalance)

			childAccount := world.AcctMap.GetAccount(test.ChildAddress)
			childESDTBalance, _ := childAccount.GetTokenBalanceUint64(test.ESDTTestTokenName, 0)
			require.Equal(t, testConfig.ESDTTokensToTransfer-testConfig.CallbackESDTTokensToTransfer, childESDTBalance)
		})
	assert.Nil(t, err)
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

	_, err := test.BuildMockInstanceCallTest(t).
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
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			setAsyncCosts(host, testConfig.GasLockCost)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				Print().
				ReturnDataDoesNotContain([]byte("out of gas")).
				ReturnData(expectedReturnData...)
		})
	assert.Nil(t, err)
}

func TestGasUsed_TransferAndExecute_CrossShard(t *testing.T) {
	testConfig := makeTestConfig()

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
		expectedTransfer := test.CreateTransferEntry(test.ParentAddress, contracts.GetChildAddressForTransfer(transfer), uint32(transfer+1)).
			WithData(big.NewInt(int64(transfer)).Bytes()).
			WithGasLimit(testConfig.GasProvidedToChild).
			WithCallType(vm.DirectCall).
			WithValue(big.NewInt(testConfig.TransferFromParentToChild))
		expectedTransfers = append(expectedTransfers, expectedTransfer)
		if transfer == 0 {
			expectedLog := vmcommon.LogEntry{
				Address: test.ParentAddress,
				Topics: [][]byte{
					big.NewInt(testConfig.TransferFromParentToChild).Bytes(),
					contracts.GetChildAddressForTransfer(transfer)},
				Data:       [][]byte{[]byte("DirectCall"), []byte("")},
				Identifier: []byte("transferValueOnly"),
			}
			expectedLogs = append(expectedLogs, expectedLog)
		} else {
			expectedLog := vmcommon.LogEntry{
				Address: test.ParentAddress,
				Topics: [][]byte{
					big.NewInt(testConfig.TransferFromParentToChild).Bytes(),
					contracts.GetChildAddressForTransfer(transfer)},
				Data:       [][]byte{[]byte("TransferAndExecute"), big.NewInt(int64(transfer)).Bytes()},
				Identifier: []byte("transferValueOnly"),
			}
			expectedLogs = append(expectedLogs, expectedLog)
		}
	}

	gasRemaining := testConfig.GasProvided - testConfig.GasUsedByParent - uint64(noOfTransfers)*testConfig.GasProvidedToChild

	_, err := test.BuildMockInstanceCallTest(t).
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
	assert.Nil(t, err)
}

func TestGasUsed_AsyncCallManaged_Mocks(t *testing.T) {
	testConfig := makeTestConfig()
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

		_, err := tester.WithInput(test.CreateTestContractCallInputBuilder().
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
		assert.Nil(t, err)
	}
}

func TestGasUsed_AsyncCallManaged(t *testing.T) {
	startValue := uint64(5000000)
	outOfGasValue := uint64(5300000)
	stopValue := uint64(3000000)
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
					verify.Ok()
				} else {
					verify.OutOfGas()
				}
			})
	}
}

func TestGasUsed_Async_CallbackWithOnSameContext(t *testing.T) {
	testConfig := makeTestConfig()
	testConfig.SuccessCallback = "callBack"
	testConfig.ErrorCallback = "callBack"
	testConfig.GasProvided = 1000

	_, err := test.BuildMockInstanceCallTest(t).
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
					test.CreateStoreEntry(test.ParentAddress).WithKey(test.OriginalCallerParent).WithValue(test.UserAddress),
					test.CreateStoreEntry(test.ChildAddress).WithKey(test.ChildKey).WithValue(test.ChildData),
					test.CreateStoreEntry(test.ChildAddress).WithKey(test.OriginalCallerChild).WithValue(test.UserAddress),
				)
		})
	assert.Nil(t, err)
}

func Test_DifferentVM_ExecuteOnDestCtx(t *testing.T) {
	testConfig := makeTestConfig()

	fakeVMType, _ := hex.DecodeString("beaf")
	childAddress, _ := hex.DecodeString("0000000000000000beaf00000000000022cd8429ce92f8973bba2a9fb51e0eb3a1")

	world := worldmock.NewMockWorld()
	vmOutput := vmhost.MakeEmptyVMOutput()
	vmhost.AddNewOutputAccountWithSender(vmOutput,
		test.ThirdPartyAddress,
		test.ParentAddress,
		testConfig.TransferToThirdParty,
		[]byte("test"))
	world.OtherVMOutputMap[string(fakeVMType)] = vmOutput

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.ExecOnDestCtxParentMock, contracts.WasteGasParentMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("execOnDestCtx").
			WithArguments(childAddress, []byte("wasteGas"), big.NewInt(2).Bytes()).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
		}).
		AndAssertResultsWithWorld(world, true, nil, nil, func(startNode *testcommon.TestCallNode, world *worldmock.MockWorld, verify *testcommon.VMOutputVerifier, expectedErrorsForRound []string) {
			verify.Ok().
				Transfers(
					test.CreateTransferEntry(test.ParentAddress, test.ThirdPartyAddress, 1).
						WithData([]byte("test")).
						WithValue(big.NewInt(testConfig.TransferToThirdParty)),
					test.CreateTransferEntry(test.ParentAddress, test.ThirdPartyAddress, 2).
						WithData([]byte("test")).
						WithValue(big.NewInt(testConfig.TransferToThirdParty)),
				)
		})

	if err != nil {
		t.Error(err)
	}
}

func TestGasUsed_MockCreateContract(t *testing.T) {
	testConfig := makeTestConfig()

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithMethods(contracts.InitFunctionMock, contracts.WasteGasParentMock)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("wasteGas").
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
		}).
		AndCreateAndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				ReturnData([]byte(vmhost.InitFunctionName))
		})
	assert.Nil(t, err)
}

func TestGasUsed_MockUpgradeContract(t *testing.T) {
	testConfig := makeTestConfig()

	codeMetadata := []byte{vmcommon.MetadataUpgradeable, 0}

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(testConfig.ParentBalance).
				WithConfig(testConfig).
				WithOwnerAddress(test.UserAddress).
				WithCodeMetadata(codeMetadata).
				WithMethods(contracts.UpgradeFunctionMock, contracts.WasteGasParentMock)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction(vmhost.UpgradeFunctionName).
			WithArguments(test.ParentAddress, codeMetadata).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok().
				ReturnData([]byte(vmhost.UpgradeFunctionName))
		})
	assert.Nil(t, err)
}

type MockClaimBuiltin struct {
	test.MockBuiltin
	AmountToGive int64
	GasCost      uint64
}

var amountToGiveByBuiltinClaim = int64(42)

func createMockBuiltinFunctions(tb testing.TB, host vmhost.VMHost, world *worldmock.MockWorld) {
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
				1,
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

	err = world.BuiltinFuncs.Container.Add("sendMessage", &test.MockBuiltin{
		ProcessBuiltinFunctionCall: func(acntSnd, acntRecv vmcommon.UserAccountHandler, vmInput *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
			vmOutput := test.MakeEmptyVMOutput()
			if acntRecv != nil {
				test.AddFinishData(vmOutput, []byte("ok"))
				vmOutput.GasRemaining = vmInput.GasProvided - mockClaimBuiltin.GasCost
				return vmOutput, nil
			}

			numberOfTransfers := 1
			args := vmInput.Arguments
			if len(args) > 0 {
				numberOfTransfers = int(big.NewInt(0).SetBytes(args[0]).Int64())
			}

			for transfer := 1; transfer <= numberOfTransfers; transfer++ {
				account := test.AddNewOutputTransfer(
					vmOutput,
					uint32(transfer),
					acntSnd.AddressBytes(),
					vmInput.RecipientAddr,
					0,
					[]byte("message"),
				)
				account.OutputTransfers[0].GasLimit = vmInput.GasProvided - mockClaimBuiltin.GasCost
			}
			vmOutput.GasRemaining = 0
			return vmOutput, nil
		},
	})
	assert.Nil(tb, err)

	host.SetBuiltInFunctionsContainer(world.BuiltinFuncs.Container)
}

func setZeroCodeCosts(host vmhost.VMHost) {
	host.Metering().GasSchedule().BaseOperationCost.CompilePerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.AoTPreparePerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.GetCode = 0
	host.Metering().GasSchedule().BaseOperationCost.StorePerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.DataCopyPerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.PersistPerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.ReleasePerByte = 0
	host.Metering().GasSchedule().BaseOpsAPICost.SignalError = 0
	host.Metering().GasSchedule().BaseOpsAPICost.ExecuteOnSameContext = 0
	host.Metering().GasSchedule().BaseOpsAPICost.ExecuteOnDestContext = 0
	host.Metering().GasSchedule().BaseOpsAPICost.StorageLoad = 0
	host.Metering().GasSchedule().BaseOpsAPICost.StorageStore = 0
	host.Metering().GasSchedule().BaseOpsAPICost.TransferValue = 0
	host.Metering().GasSchedule().BaseOpsAPICost.CreateContract = 0
	host.Metering().GasSchedule().DynamicStorageLoad.MinGasCost = 0
	host.Metering().GasSchedule().DynamicStorageLoad.Linear = 0
	host.Metering().GasSchedule().DynamicStorageLoad.Constant = 0
	host.Metering().GasSchedule().DynamicStorageLoad.Quadratic = 0
}

func setAsyncCosts(host vmhost.VMHost, gasLockCost uint64) {
	host.Metering().GasSchedule().BaseOpsAPICost.CreateAsyncCall = 0
	host.Metering().GasSchedule().BaseOpsAPICost.SetAsyncCallback = 0
	host.Metering().GasSchedule().BaseOpsAPICost.AsyncCallStep = 0
	host.Metering().GasSchedule().BaseOpsAPICost.GetCallbackClosure = 0
	host.Metering().GasSchedule().BaseOpsAPICost.AsyncCallbackGasLock = gasLockCost
}

func computeReturnDataForCallback(returnCode vmcommon.ReturnCode, returnData [][]byte) []byte {
	builtReturnData := txDataBuilder.NewBuilder()
	builtReturnData.Func("<callback>")
	builtReturnData.Bytes([]byte{})
	builtReturnData.Bytes([]byte{})
	builtReturnData.Bytes([]byte{})
	builtReturnData.Bytes([]byte{})
	builtReturnData.Int(int(returnCode))
	for _, data := range returnData {
		builtReturnData.Bytes(data)
	}
	return builtReturnData.ToBytes()
	// TODO(check) commented code

	// retCode := string(big.NewInt(int64(returnCode)).Bytes())
	// retData := []byte("@" + hex.EncodeToString(prevTxHash))
	// retData = append(retData, []byte("@"+retCode)...)
	// for _, data := range returnData {
	// 	retData = append(retData, []byte("@"+hex.EncodeToString(data))...)
	// }
	// return retData
}

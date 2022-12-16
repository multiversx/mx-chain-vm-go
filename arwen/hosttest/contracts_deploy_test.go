package hosttest

import (
	"math/big"
	"testing"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen"
	arwenMock "github.com/ElrondNetwork/wasm-vm-v1_4/arwen/mock"
	mock "github.com/ElrondNetwork/wasm-vm-v1_4/mock/context"
	"github.com/ElrondNetwork/wasm-vm-v1_4/mock/contracts"
	worldmock "github.com/ElrondNetwork/wasm-vm-v1_4/mock/world"
	"github.com/ElrondNetwork/wasm-vm-v1_4/testcommon"
	test "github.com/ElrondNetwork/wasm-vm-v1_4/testcommon"
	"github.com/stretchr/testify/require"
)

var sc1Address = testcommon.MakeTestSCAddress("sc1")
var sc2Address = testcommon.MakeTestSCAddress("sc2")

type deployFromSourceTestConfig struct {
	deployedContractAddress []byte
	gasUsedByInit           uint64
	gasProvided             uint64
	gasProvidedForInit      uint64
	asyncCallStepCost       uint64
	aoTPreparePerByteCost   uint64
	compilePerByteCost      uint64
}

type updateFromSourceTestConfig struct {
	deployFromSourceTestConfig
	contractToBeUpdatedAddress []byte
	owner                      []byte
	isFlagEnabled              bool
	hasCallback                bool
	callbackFails              bool
}

func getDeployFromSourceTestConfig() deployFromSourceTestConfig {
	return deployFromSourceTestConfig{
		deployedContractAddress: sc1Address,
		gasUsedByInit:           uint64(200),
		gasProvidedForInit:      uint64(300),
		gasProvided:             uint64(1000),
		aoTPreparePerByteCost:   uint64(1),
		compilePerByteCost:      uint64(2),
	}
}

func getUpdateFromSourceTestConfig() updateFromSourceTestConfig {
	return updateFromSourceTestConfig{
		deployFromSourceTestConfig: getDeployFromSourceTestConfig(),
		contractToBeUpdatedAddress: sc2Address,
		owner:                      test.ParentAddress,
		isFlagEnabled:              true,
		hasCallback:                true,
		callbackFails:              false,
	}
}

// GetGasUsedByChild
func (config deployFromSourceTestConfig) GetGasUsedByChild() uint64 {
	return config.gasUsedByInit
}

// GetGasUsedByChild
func (config updateFromSourceTestConfig) GetGasUsedByChild() uint64 {
	return config.gasUsedByInit
}

// CallbackFails
func (config updateFromSourceTestConfig) CallbackFails() bool {
	return config.callbackFails
}

func TestDeployFromSource_Success(t *testing.T) {
	testConfig := getDeployFromSourceTestConfig()
	deployedCode := testConfig.deployedContractAddress /* this is the actual mock code of the deployed contract */
	deployedCodeLen := uint64(len(deployedCode))
	runDeployFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		newContractAddress := verify.VmOutput.ReturnData[0]
		verify.
			Ok().
			Code(newContractAddress, deployedCode).
			GasRemaining(testConfig.gasProvided -
				testConfig.gasUsedByInit -
				deployedCodeLen*testConfig.compilePerByteCost -
				deployedCodeLen*testConfig.aoTPreparePerByteCost)
	})
}

func TestDeployFromSource_NoGasForInit(t *testing.T) {
	testConfig := getDeployFromSourceTestConfig()
	testConfig.gasProvidedForInit = uint64(100)
	// TODO investigate why the ReturnCode is ExecutionFailed instead of OutOfGas
	runDeployFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			ExecutionFailed().
			HasRuntimeErrors(arwen.ErrInputAndOutputGasDoesNotMatch.Error())
	})
}

func TestDeployFromSource_NoGasForAoTPrepare(t *testing.T) {
	testConfig := getDeployFromSourceTestConfig()
	testConfig.aoTPreparePerByteCost = uint64(10)
	runDeployFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas()
	})
}

func TestDeployFromSource_NoGasForCompile(t *testing.T) {
	testConfig := getDeployFromSourceTestConfig()
	testConfig.compilePerByteCost = uint64(100)
	runDeployFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas()
	})
}

func TestDeployFromSource_NoContract(t *testing.T) {
	testConfig := getDeployFromSourceTestConfig()
	testConfig.deployedContractAddress = nil
	runDeployFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			ExecutionFailed().
			HasRuntimeErrors(arwen.ErrContractInvalid.Error())
	})
}

func runDeployFromSourceTest(t *testing.T, testConfig *deployFromSourceTestConfig, asserts func(world *worldmock.MockWorld, verify *test.VMOutputVerifier)) {
	var deployedContract test.MockTestSmartContract
	if testConfig.deployedContractAddress != nil {
		deployedContract = test.CreateMockContract(testConfig.deployedContractAddress).
			WithConfig(testConfig).
			WithMethods(contracts.InitMockMethod)
	}
	test.BuildMockInstanceCallTest(t).
		WithContracts(
			deployedContract,
			test.CreateMockContract(test.ParentAddress).
				WithConfig(testConfig).
				WithMethods(contracts.DeployContractFromSourceMock)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.gasProvided).
			WithFunction("deployContractFromSource").
			WithArguments(testConfig.deployedContractAddress, []byte{0, 0}, big.NewInt(int64(testConfig.gasProvidedForInit)).Bytes()).
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			host.Metering().GasSchedule().BaseOperationCost.AoTPreparePerByte = testConfig.aoTPreparePerByteCost
			host.Metering().GasSchedule().BaseOperationCost.CompilePerByte = testConfig.compilePerByteCost
		}).
		AndAssertResults(asserts)
}

func TestUpdateFromSource_Success_EpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	updatedCode := testConfig.deployedContractAddress /* this is the actual mock code of the deployed contract */
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			Code(testConfig.contractToBeUpdatedAddress, updatedCode)
	})
}

func TestUpdateFromSource_Success_EpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	updatedCode := testConfig.deployedContractAddress /* this is the actual mock code of the deployed contract */
	testConfig.hasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			Code(testConfig.contractToBeUpdatedAddress, updatedCode)
	})
}

func TestUpdateFromSource_CallbackFails_EpochFlag(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.callbackFails = true
	updatedCode := testConfig.deployedContractAddress /* this is the actual mock code of the deployed contract */
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			UserError().
			Code(testConfig.contractToBeUpdatedAddress, updatedCode)
	})
}

func TestUpdateFromSource_CallbackFails_NoEpochFlag(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.isFlagEnabled = false
	testConfig.callbackFails = true
	updatedCode := testConfig.deployedContractAddress /* this is the actual mock code of the deployed contract */
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			HasRuntimeErrors(arwen.ErrSignalError.Error()).
			Code(testConfig.contractToBeUpdatedAddress, updatedCode)
	})
}

func TestUpdateFromSource_Success_NoEpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.isFlagEnabled = false
	updatedCode := testConfig.deployedContractAddress /* this is the actual mock code of the deployed contract */
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			Code(testConfig.contractToBeUpdatedAddress, updatedCode)
	})
}

func TestUpdateFromSource_Success_NoEpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.isFlagEnabled = false
	testConfig.hasCallback = false
	updatedCode := testConfig.deployedContractAddress /* this is the actual mock code of the deployed contract */
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			Code(testConfig.contractToBeUpdatedAddress, updatedCode)
	})
}

func TestUpdateFromSource_NoPermission_EpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.owner = nil
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			ExecutionFailed().
			HasRuntimeErrors(arwen.ErrUpgradeNotAllowed.Error()).
			Code(testConfig.contractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoPermission_EpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.owner = nil
	testConfig.hasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas(). // not enough gas to provide for callback execution
			HasRuntimeErrors(arwen.ErrUpgradeNotAllowed.Error()).
			Code(testConfig.contractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoPermission_NoEpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.isFlagEnabled = false
	testConfig.owner = nil
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			HasRuntimeErrors(arwen.ErrUpgradeNotAllowed.Error()).
			Code(testConfig.contractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoPermission_NoEpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.isFlagEnabled = false
	testConfig.owner = nil
	testConfig.hasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			HasRuntimeErrors(arwen.ErrUpgradeNotAllowed.Error()).
			Code(testConfig.contractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoGasForAsyncCall_EpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.gasProvidedForInit = uint64(100)
	testConfig.asyncCallStepCost = uint64(300)
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas()
		require.Nil(t, verify.VmOutput.OutputAccounts[string(testConfig.contractToBeUpdatedAddress)])
	})
}

func TestUpdateFromSource_NoGasForAsyncCall_EpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.gasProvidedForInit = uint64(100)
	testConfig.asyncCallStepCost = uint64(300)
	testConfig.hasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas()
		require.Nil(t, verify.VmOutput.OutputAccounts[string(testConfig.contractToBeUpdatedAddress)])
	})
}

func TestUpdateFromSource_NoGasForAsyncCall_NoEpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.isFlagEnabled = false
	testConfig.gasProvidedForInit = uint64(100)
	testConfig.asyncCallStepCost = uint64(300)
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas()
		require.Nil(t, verify.VmOutput.OutputAccounts[string(testConfig.contractToBeUpdatedAddress)])
	})
}

func TestUpdateFromSource_NoGasForAsyncCall_NoEpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.isFlagEnabled = false
	testConfig.gasProvidedForInit = uint64(100)
	testConfig.asyncCallStepCost = uint64(300)
	testConfig.hasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas()
		require.Nil(t, verify.VmOutput.OutputAccounts[string(testConfig.contractToBeUpdatedAddress)])
	})
}

func TestUpdateFromSource_NoGasForInit_EpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.gasUsedByInit = 1000
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas(). // not enough gas to provide for callback execution
			HasRuntimeErrors(arwen.ErrNotEnoughGas.Error()).
			Code(testConfig.contractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoGasForInit_EpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.gasUsedByInit = 1000
	testConfig.hasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas(). // not enough gas to provide for callback execution
			HasRuntimeErrors(arwen.ErrNotEnoughGas.Error()).
			Code(testConfig.contractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoGasForInit_NoEpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.isFlagEnabled = false
	testConfig.gasUsedByInit = 1000
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			HasRuntimeErrors(arwen.ErrNotEnoughGas.Error()).
			Code(testConfig.contractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoGasForInit_NoEpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.isFlagEnabled = false
	testConfig.gasUsedByInit = 1000
	testConfig.hasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			HasRuntimeErrors(arwen.ErrNotEnoughGas.Error()).
			Code(testConfig.contractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoGasForCompile_EpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.compilePerByteCost = uint64(50)
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas(). // not enough gas to provide for callback execution
			HasRuntimeErrors(arwen.ErrNotEnoughGas.Error()).
			Code(testConfig.contractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoGasForCompile_EpochFlag_NoCallBack(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.compilePerByteCost = uint64(50)
	testConfig.hasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas(). // not enough gas to provide for callback execution
			HasRuntimeErrors(arwen.ErrNotEnoughGas.Error()).
			Code(testConfig.contractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoGasForCompile_NoEpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.isFlagEnabled = false
	testConfig.compilePerByteCost = uint64(50)
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			HasRuntimeErrors(arwen.ErrNotEnoughGas.Error()).
			Code(testConfig.contractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoGasForCompile_NoEpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.isFlagEnabled = false
	testConfig.compilePerByteCost = uint64(50)
	testConfig.hasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			HasRuntimeErrors(arwen.ErrNotEnoughGas.Error()).
			Code(testConfig.contractToBeUpdatedAddress, nil)
	})
}

func runUpdateFromSourceTest(t *testing.T, testConfig *updateFromSourceTestConfig, asserts func(world *worldmock.MockWorld, verify *test.VMOutputVerifier)) {
	var deployedContract test.MockTestSmartContract
	var contractToUpdate test.MockTestSmartContract
	if testConfig.deployedContractAddress != nil {
		deployedContract = test.CreateMockContract(testConfig.deployedContractAddress).
			WithConfig(testConfig).
			WithMethods(contracts.InitMockMethod)
	}
	if testConfig.contractToBeUpdatedAddress != nil {
		contractToUpdate = test.CreateMockContract(testConfig.contractToBeUpdatedAddress).
			WithConfig(testConfig).
			WithCodeMetadata([]byte{vmcommon.MetadataUpgradeable, 0}).
			WithOwnerAddress(testConfig.owner).
			WithMethods(contracts.InitMockMethod)
	}

	methods := []func(*mock.InstanceMock, interface{}){contracts.UpdateContractFromSourceMock}
	if testConfig.hasCallback {
		methods = append(methods, contracts.CallbackMockMethodThatCouldFail)
	}

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			deployedContract,
			contractToUpdate,
			test.CreateMockContract(test.ParentAddress).
				WithConfig(testConfig).
				WithMethods(methods...)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.gasProvided).
			WithFunction("updateContractFromSource").
			WithArguments(testConfig.deployedContractAddress, testConfig.contractToBeUpdatedAddress,
				nil,
				big.NewInt(int64(testConfig.gasProvidedForInit)).Bytes()).
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			gasSchedule := host.Metering().GasSchedule()
			gasSchedule.ElrondAPICost.AsyncCallStep = testConfig.asyncCallStepCost
			gasSchedule.BaseOperationCost.AoTPreparePerByte = testConfig.aoTPreparePerByteCost
			gasSchedule.BaseOperationCost.CompilePerByte = testConfig.compilePerByteCost
			gasSchedule.ElrondAPICost.AsyncCallbackGasLock = 0

			if !testConfig.isFlagEnabled {
				enableEpochsHandler, _ := host.EnableEpochsHandler().(*arwenMock.EnableEpochsHandlerStub)
				enableEpochsHandler.IsStorageAPICostOptimizationFlagEnabledField = false
			}
		}).
		AndAssertResults(asserts)
}

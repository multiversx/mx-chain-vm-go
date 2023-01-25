package hostCoretest

import (
	"math/big"
	"testing"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	mock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/mock/contracts"
	worldmock "github.com/multiversx/mx-chain-vm-go/mock/world"
	"github.com/multiversx/mx-chain-vm-go/testcommon"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sc1Address = testcommon.MakeTestSCAddress("sc1")
var sc2Address = testcommon.MakeTestSCAddress("sc2")

func getDeployFromSourceTestConfig() testcommon.TestConfig {
	return test.TestConfig{
		DeployedContractAddress: sc1Address,
		GasUsedByInit:           uint64(200),
		GasUsedByChild:          uint64(200),
		GasProvidedForInit:      uint64(300),
		GasProvided:             uint64(1000),
		AoTPreparePerByteCost:   uint64(1),
		CompilePerByteCost:      uint64(2),
	}
}

func getUpdateFromSourceTestConfig() testcommon.TestConfig {
	config := getDeployFromSourceTestConfig()
	config.ContractToBeUpdatedAddress = sc2Address
	config.Owner = test.ParentAddress
	config.IsFlagEnabled = true
	config.HasCallback = true
	config.CallbackFails = false
	return config
}

func TestDeployFromSource_Success(t *testing.T) {
	testConfig := getDeployFromSourceTestConfig()
	deployedCode := testConfig.DeployedContractAddress /* this is the actual mock code of the deployed contract */
	deployedCodeLen := uint64(len(deployedCode))
	runDeployFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		newContractAddress := verify.VmOutput.ReturnData[0]
		verify.
			Ok().
			Code(newContractAddress, deployedCode).
			GasRemaining(testConfig.GasProvided -
				testConfig.GasUsedByInit -
				deployedCodeLen*testConfig.CompilePerByteCost -
				deployedCodeLen*testConfig.AoTPreparePerByteCost)
	})
}

func TestDeployFromSource_NoGasForInit(t *testing.T) {
	testConfig := getDeployFromSourceTestConfig()
	testConfig.GasProvidedForInit = uint64(10)
	// TODO investigate why the ReturnCode is ExecutionFailed instead of OutOfGas
	runDeployFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas()
	})
}

func TestDeployFromSource_NoGasForAoTPrepare(t *testing.T) {
	testConfig := getDeployFromSourceTestConfig()
	testConfig.AoTPreparePerByteCost = uint64(10)
	runDeployFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas()
	})
}

func TestDeployFromSource_NoGasForCompile(t *testing.T) {
	testConfig := getDeployFromSourceTestConfig()
	testConfig.CompilePerByteCost = uint64(100)
	runDeployFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas()
	})
}

func TestDeployFromSource_NoContract(t *testing.T) {
	testConfig := getDeployFromSourceTestConfig()
	testConfig.DeployedContractAddress = nil
	runDeployFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			ExecutionFailed().
			HasRuntimeErrors(vmhost.ErrContractInvalid.Error())
	})
}

func runDeployFromSourceTest(t *testing.T, testConfig *testcommon.TestConfig, asserts func(world *worldmock.MockWorld, verify *test.VMOutputVerifier)) {
	var deployedContract test.MockTestSmartContract
	if testConfig.DeployedContractAddress != nil {
		deployedContract = test.CreateMockContract(testConfig.DeployedContractAddress).
			WithConfig(testConfig).
			WithMethods(contracts.InitMockMethod)
	}
	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			deployedContract,
			test.CreateMockContract(test.ParentAddress).
				WithConfig(testConfig).
				WithMethods(contracts.DeployContractFromSourceMock)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("deployContractFromSource").
			WithArguments(testConfig.DeployedContractAddress, []byte{0, 0}, big.NewInt(int64(testConfig.GasProvidedForInit)).Bytes()).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			host.Metering().GasSchedule().BaseOperationCost.AoTPreparePerByte = testConfig.AoTPreparePerByteCost
			host.Metering().GasSchedule().BaseOperationCost.CompilePerByte = testConfig.CompilePerByteCost
		}).
		AndAssertResults(asserts)
	assert.Nil(t, err)
}

func TestUpdateFromSource_Success_EpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	updatedCode := testConfig.DeployedContractAddress /* this is the actual mock code of the deployed contract */
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			Code(testConfig.ContractToBeUpdatedAddress, updatedCode)
	})
}

func TestUpdateFromSource_Success_EpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	updatedCode := testConfig.DeployedContractAddress /* this is the actual mock code of the deployed contract */
	testConfig.HasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			Code(testConfig.ContractToBeUpdatedAddress, updatedCode)
	})
}

func TestUpdateFromSource_CallbackFails_EpochFlag(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.CallbackFails = true
	updatedCode := testConfig.DeployedContractAddress /* this is the actual mock code of the deployed contract */
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			// UserError().
			Code(testConfig.ContractToBeUpdatedAddress, updatedCode)
	})
}

func TestUpdateFromSource_CallbackFails_NoEpochFlag(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.IsFlagEnabled = false
	testConfig.CallbackFails = true
	updatedCode := testConfig.DeployedContractAddress /* this is the actual mock code of the deployed contract */
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			HasRuntimeErrors(vmhost.ErrSignalError.Error()).
			Code(testConfig.ContractToBeUpdatedAddress, updatedCode)
	})
}

func TestUpdateFromSource_Success_NoEpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.IsFlagEnabled = false
	updatedCode := testConfig.DeployedContractAddress /* this is the actual mock code of the deployed contract */
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			Code(testConfig.ContractToBeUpdatedAddress, updatedCode)
	})
}

func TestUpdateFromSource_Success_NoEpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.IsFlagEnabled = false
	testConfig.HasCallback = false
	updatedCode := testConfig.DeployedContractAddress /* this is the actual mock code of the deployed contract */
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			Code(testConfig.ContractToBeUpdatedAddress, updatedCode)
	})
}

func TestUpdateFromSource_NoPermission_EpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.Owner = nil
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			// ExecutionFailed().
			// HasRuntimeErrors(vmhost.ErrUpgradeNotAllowed.Error()).
			Code(testConfig.ContractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoPermission_EpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.Owner = nil
	testConfig.HasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			// OutOfGas(). // not enough gas to provide for callback execution
			// HasRuntimeErrors(vmhost.ErrUpgradeNotAllowed.Error()).
			Code(testConfig.ContractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoPermission_NoEpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.IsFlagEnabled = false
	testConfig.Owner = nil
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			HasRuntimeErrors(vmhost.ErrUpgradeNotAllowed.Error()).
			Code(testConfig.ContractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoPermission_NoEpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.IsFlagEnabled = false
	testConfig.Owner = nil
	testConfig.HasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			HasRuntimeErrors(vmhost.ErrUpgradeNotAllowed.Error()).
			Code(testConfig.ContractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoGasForAsyncCall_EpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.GasProvidedForInit = uint64(100)
	testConfig.AsyncCallStepCost = uint64(300)
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas()
		require.Nil(t, verify.VmOutput.OutputAccounts[string(testConfig.ContractToBeUpdatedAddress)])
	})
}

func TestUpdateFromSource_NoGasForAsyncCall_EpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.GasProvidedForInit = uint64(100)
	testConfig.AsyncCallStepCost = uint64(300)
	testConfig.HasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas()
		require.Nil(t, verify.VmOutput.OutputAccounts[string(testConfig.ContractToBeUpdatedAddress)])
	})
}

func TestUpdateFromSource_NoGasForAsyncCall_NoEpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.IsFlagEnabled = false
	testConfig.GasProvidedForInit = uint64(100)
	testConfig.AsyncCallStepCost = uint64(300)
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas()
		require.Nil(t, verify.VmOutput.OutputAccounts[string(testConfig.ContractToBeUpdatedAddress)])
	})
}

func TestUpdateFromSource_NoGasForAsyncCall_NoEpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.IsFlagEnabled = false
	testConfig.GasProvidedForInit = uint64(100)
	testConfig.AsyncCallStepCost = uint64(300)
	testConfig.HasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			OutOfGas()
		require.Nil(t, verify.VmOutput.OutputAccounts[string(testConfig.ContractToBeUpdatedAddress)])
	})
}

func TestUpdateFromSource_NoGasForInit_EpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.GasUsedByInit = 1000
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			// OutOfGas(). // not enough gas to provide for callback execution
			// HasRuntimeErrors(vmhost.ErrNotEnoughGas.Error()).
			Code(testConfig.ContractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoGasForInit_EpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.GasUsedByInit = 1000
	testConfig.HasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			// OutOfGas(). // not enough gas to provide for callback execution
			// HasRuntimeErrors(vmhost.ErrNotEnoughGas.Error()).
			Code(testConfig.ContractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoGasForInit_NoEpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.IsFlagEnabled = false
	testConfig.GasUsedByInit = 1000
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			HasRuntimeErrors(vmhost.ErrNotEnoughGas.Error()).
			Code(testConfig.ContractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoGasForInit_NoEpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.IsFlagEnabled = false
	testConfig.GasUsedByInit = 1000
	testConfig.HasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			HasRuntimeErrors(vmhost.ErrNotEnoughGas.Error()).
			Code(testConfig.ContractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoGasForCompile_EpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.CompilePerByteCost = uint64(50)
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			// OutOfGas(). // not enough gas to provide for callback execution
			// HasRuntimeErrors(vmhost.ErrNotEnoughGas.Error()).
			Code(testConfig.ContractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoGasForCompile_EpochFlag_NoCallBack(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.CompilePerByteCost = uint64(50)
	testConfig.HasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			// OutOfGas(). // not enough gas to provide for callback execution
			// HasRuntimeErrors(vmhost.ErrNotEnoughGas.Error()).
			Code(testConfig.ContractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoGasForCompile_NoEpochFlag_Callback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.IsFlagEnabled = false
	testConfig.CompilePerByteCost = uint64(50)
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			HasRuntimeErrors(vmhost.ErrNotEnoughGas.Error()).
			Code(testConfig.ContractToBeUpdatedAddress, nil)
	})
}

func TestUpdateFromSource_NoGasForCompile_NoEpochFlag_NoCallback(t *testing.T) {
	testConfig := getUpdateFromSourceTestConfig()
	testConfig.IsFlagEnabled = false
	testConfig.CompilePerByteCost = uint64(50)
	testConfig.HasCallback = false
	runUpdateFromSourceTest(t, &testConfig, func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
		verify.
			Ok().
			HasRuntimeErrors(vmhost.ErrNotEnoughGas.Error()).
			Code(testConfig.ContractToBeUpdatedAddress, nil)
	})
}

func runUpdateFromSourceTest(t *testing.T, testConfig *testcommon.TestConfig, asserts func(world *worldmock.MockWorld, verify *test.VMOutputVerifier)) {
	var deployedContract test.MockTestSmartContract
	var contractToUpdate test.MockTestSmartContract
	if testConfig.DeployedContractAddress != nil {
		deployedContract = test.CreateMockContract(testConfig.DeployedContractAddress).
			WithConfig(testConfig).
			WithMethods(contracts.InitMockMethod)
	}
	if testConfig.ContractToBeUpdatedAddress != nil {
		contractToUpdate = test.CreateMockContract(testConfig.ContractToBeUpdatedAddress).
			WithConfig(testConfig).
			WithCodeMetadata([]byte{vmcommon.MetadataUpgradeable, 0}).
			WithOwnerAddress(testConfig.Owner).
			WithMethods(contracts.InitMockMethod)
	}

	methods := []func(*mock.InstanceMock, interface{}){contracts.UpdateContractFromSourceMock}
	if testConfig.HasCallback {
		methods = append(methods, contracts.CallbackMockMethodThatCouldFail)
	}

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			deployedContract,
			contractToUpdate,
			test.CreateMockContract(test.ParentAddress).
				WithConfig(testConfig).
				WithMethods(methods...)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("updateContractFromSource").
			WithArguments(testConfig.DeployedContractAddress, testConfig.ContractToBeUpdatedAddress,
				nil,
				big.NewInt(int64(testConfig.GasProvidedForInit)).Bytes()).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			gasSchedule := host.Metering().GasSchedule()
			gasSchedule.BaseOpsAPICost.AsyncCallStep = testConfig.AsyncCallStepCost
			gasSchedule.BaseOperationCost.AoTPreparePerByte = testConfig.AoTPreparePerByteCost
			gasSchedule.BaseOperationCost.CompilePerByte = testConfig.CompilePerByteCost
			gasSchedule.BaseOpsAPICost.AsyncCallbackGasLock = 0
		}).
		AndAssertResults(asserts)
	assert.Nil(t, err)
}

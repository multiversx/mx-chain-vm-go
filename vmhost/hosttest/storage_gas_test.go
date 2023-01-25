package hostCoretest

import (
	"testing"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/mock/contracts"
	worldmock "github.com/multiversx/mx-chain-vm-go/mock/world"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/stretchr/testify/assert"
)

var smallKey = []byte("testKey")
var bigKey = make([]byte, 50)

const storageLoadGas = uint64(10)
const cachedStorageLoadGas = uint64(5)
const dataCopyGas = uint64(1)

func TestGasUsed_LoadStorage_SmallKey_FlagEnabled(t *testing.T) {
	loadStorage(t, smallKey)
}

func TestGasUsed_LoadStorage_BigKey_FlagEnabled(t *testing.T) {
	loadStorage(t, bigKey)
}

func loadStorage(t *testing.T, key []byte) {
	testConfig := makeTestConfig()
	value := []byte("testValue")

	expectedUsedGas := computeExpectedGasForGetStorage(key, value)

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(0).
				WithConfig(nil).
				WithMethods(contracts.LoadStore)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("loadStore").
			WithArguments(key).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			host.Metering().GasSchedule().BaseOpsAPICost.StorageLoad = storageLoadGas
			host.Metering().GasSchedule().BaseOpsAPICost.CachedStorageLoad = cachedStorageLoadGas
			host.Metering().GasSchedule().BaseOperationCost.DataCopyPerByte = dataCopyGas
			host.Metering().GasSchedule().BaseOperationCost.PersistPerByte = 0

			accountHandler, _ := world.GetUserAccount(test.ParentAddress)
			(accountHandler.(*worldmock.Account)).Storage[string(key)] = value
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasRemaining(testConfig.GasProvided-expectedUsedGas).
				GasUsed(test.ParentAddress, expectedUsedGas).
				ReturnData(value)
		})
	assert.Nil(t, err)
}

func TestGasUsed_LoadStorageFromAddress_SmallKey_FlagEnabled(t *testing.T) {
	loadStorageFromAddress(t, smallKey)
}

func TestGasUsed_LoadStorageFromAddress_BigKey_FlagEnabled(t *testing.T) {
	loadStorageFromAddress(t, bigKey)
}

func loadStorageFromAddress(t *testing.T, key []byte) {
	testConfig := makeTestConfig()
	value := []byte("testValue")

	expectedUsedGas := computeExpectedGasForGetStorage(key, value)

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.UserAddress).
				WithBalance(0).
				WithConfig(nil).
				WithMethods(),
			test.CreateMockContract(test.ParentAddress).
				WithBalance(0).
				WithConfig(nil).
				WithMethods(contracts.LoadStoreFromAddress)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("loadStoreFromAddress").
			WithArguments(test.UserAddress, key).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			host.Metering().GasSchedule().BaseOpsAPICost.StorageLoad = storageLoadGas
			host.Metering().GasSchedule().BaseOpsAPICost.CachedStorageLoad = cachedStorageLoadGas
			host.Metering().GasSchedule().BaseOperationCost.DataCopyPerByte = dataCopyGas
			host.Metering().GasSchedule().BaseOperationCost.PersistPerByte = 0

			account := world.AcctMap[string(test.UserAddress)]
			account.Storage[string(key)] = value
			account.CodeMetadata = []byte{vmcommon.MetadataReadable, 0}
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasRemaining(testConfig.GasProvided-expectedUsedGas).
				GasUsed(test.ParentAddress, expectedUsedGas).
				ReturnData(value)
		})
	assert.Nil(t, err)
}

func computeExpectedGasForGetStorage(key []byte, value []byte) uint64 {
	extraBytesForKey := len(key) - vmhost.AddressLen
	if extraBytesForKey < 0 {
		extraBytesForKey = 0
	}

	expectedUsedGas := storageLoadGas + uint64(len(value))*dataCopyGas + cachedStorageLoadGas + uint64(extraBytesForKey)*dataCopyGas
	return expectedUsedGas
}

func TestGasUsed_SetStorage_FlagEnabled(t *testing.T) {
	setStorage(t, smallKey)
}

func TestGasUsed_SetStorage_BigKey_FlagEnabled(t *testing.T) {
	setStorage(t, bigKey)
}

func setStorage(t *testing.T, key []byte) {
	testConfig := makeTestConfig()
	value := []byte("testValue")

	storageStoreGas := uint64(10)
	dataCopyGas := uint64(1)

	expectedUsedGas := 2 * storageStoreGas
	if len(key) > vmhost.AddressLen {
		expectedUsedGas += uint64(len(key) - vmhost.AddressLen)
	}

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.UserAddress).
				WithBalance(0).
				WithConfig(nil).
				WithMethods(),
			test.CreateMockContract(test.ParentAddress).
				WithBalance(0).
				WithConfig(nil).
				WithMethods(contracts.SetStore)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(testConfig.GasProvided).
			WithFunction("setStore").
			WithArguments(key, value).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			host.Metering().GasSchedule().BaseOpsAPICost.StorageStore = storageStoreGas
			host.Metering().GasSchedule().BaseOperationCost.DataCopyPerByte = dataCopyGas
			host.Metering().GasSchedule().BaseOperationCost.PersistPerByte = 0

			account := world.AcctMap[string(test.UserAddress)]
			account.Storage[string(key)] = value
			account.CodeMetadata = []byte{vmcommon.MetadataReadable, 0}
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(test.ParentAddress, expectedUsedGas).
				GasRemaining(testConfig.GasProvided - expectedUsedGas)
		})
	assert.Nil(t, err)
}

var expectedAddByParent = len(contracts.TestStorageValue1) +
	len(contracts.TestStorageValue2)

var expectedDeletedByParent = (len(contracts.TestStorageValue1) - len(contracts.TestStorageValue2)) +
	(len(contracts.TestStorageValue2) - len(contracts.TestStorageValue3))

var expectedAddByChild = len(contracts.TestStorageValue2) +
	(len(contracts.TestStorageValue1) - len(contracts.TestStorageValue2)) +
	len(contracts.TestStorageValue1)

var expectedDeletedByChild = len(contracts.TestStorageValue1) - len(contracts.TestStorageValue4)

func TestBytesCount_SetStorage_ExecuteOnSameCtx(t *testing.T) {
	simpleGasTestConfig := makeTestConfig()
	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(simpleGasTestConfig.ParentBalance).
				WithConfig(simpleGasTestConfig).
				WithMethods(contracts.ParentSetStorageMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(simpleGasTestConfig.ChildBalance).
				WithConfig(simpleGasTestConfig).
				WithMethods(contracts.ChildSetStorageMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(simpleGasTestConfig.GasProvided).
			WithFunction("parentSetStorage").
			WithArguments([]byte{0}).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				BytesAddedToStorage(test.ParentAddress,
					expectedAddByParent+expectedAddByChild).
				BytesAddedToStorage(test.ChildAddress,
					0).
				BytesDeletedFromStorage(test.ParentAddress,
					expectedDeletedByParent+expectedDeletedByChild).
				BytesDeletedFromStorage(test.ChildAddress,
					0)
		})
	assert.Nil(t, err)
}

func TestBytesCount_SetStorage_ExecuteOnDestCtx(t *testing.T) {
	simpleGasTestConfig := makeTestConfig()
	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(simpleGasTestConfig.ParentBalance).
				WithConfig(simpleGasTestConfig).
				WithMethods(contracts.ParentSetStorageMock),
			test.CreateMockContract(test.ChildAddress).
				WithBalance(simpleGasTestConfig.ChildBalance).
				WithConfig(simpleGasTestConfig).
				WithMethods(contracts.ChildSetStorageMock),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(simpleGasTestConfig.GasProvided).
			WithFunction("parentSetStorage").
			WithArguments([]byte{1}).
			Build()).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				BytesAddedToStorage(test.ParentAddress,
					expectedAddByParent).
				BytesAddedToStorage(test.ChildAddress,
					expectedAddByChild).
				BytesDeletedFromStorage(test.ParentAddress,
					expectedDeletedByParent).
				BytesDeletedFromStorage(test.ChildAddress,
					expectedDeletedByChild)
		})
	assert.Nil(t, err)
}

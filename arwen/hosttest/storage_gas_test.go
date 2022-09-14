package hosttest

import (
	"testing"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/mock/contracts"
	worldmock "github.com/ElrondNetwork/wasm-vm/mock/world"
	test "github.com/ElrondNetwork/wasm-vm/testcommon"
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

	test.BuildMockInstanceCallTest(t).
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			host.Metering().GasSchedule().ElrondAPICost.StorageLoad = storageLoadGas
			host.Metering().GasSchedule().ElrondAPICost.CachedStorageLoad = cachedStorageLoadGas
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

	test.BuildMockInstanceCallTest(t).
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			host.Metering().GasSchedule().ElrondAPICost.StorageLoad = storageLoadGas
			host.Metering().GasSchedule().ElrondAPICost.CachedStorageLoad = cachedStorageLoadGas
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
}

func computeExpectedGasForGetStorage(key []byte, value []byte) uint64 {
	extraBytesForKey := len(key) - arwen.AddressLen
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
	if len(key) > arwen.AddressLen {
		expectedUsedGas += uint64(len(key) - arwen.AddressLen)
	}

	test.BuildMockInstanceCallTest(t).
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			host.Metering().GasSchedule().ElrondAPICost.StorageStore = storageStoreGas
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
	test.BuildMockInstanceCallTest(t).
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
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
}

func TestBytesCount_SetStorage_ExecuteOnDestCtx(t *testing.T) {
	simpleGasTestConfig := makeTestConfig()
	test.BuildMockInstanceCallTest(t).
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
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
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
}

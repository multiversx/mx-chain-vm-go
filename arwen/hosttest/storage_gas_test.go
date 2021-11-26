package hosttest

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/contracts"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/world"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/testcommon"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var smallKey = []byte("testKey")
var bigKey = make([]byte, 50)

func TestGasUsed_LoadStorage_SmallKey_FlagEnabled(t *testing.T) {
	loadStorage(t, smallKey, true)
}

func TestGasUsed_LoadStorage_SmallKey_FlagDisabled(t *testing.T) {
	loadStorage(t, smallKey, false)
}

func TestGasUsed_LoadStorage_BigKey_FlagEnabled(t *testing.T) {
	loadStorage(t, bigKey, true)
}

func TestGasUsed_LoadStorage_BigKey_FlagDisabled(t *testing.T) {
	loadStorage(t, bigKey, false)
}

func loadStorage(t *testing.T, key []byte, flagEnabled bool) {
	value := []byte("testValue")

	storageLoadGas := uint64(10)
	cachedStorageLoadGas := uint64(5)
	dataCopyGas := uint64(1)

	extraBytesForKey := len(key) - arwen.AddressLen
	if extraBytesForKey < 0 {
		extraBytesForKey = 0
	}
	var expectedUsedGas uint64
	if flagEnabled {
		expectedUsedGas = storageLoadGas + uint64(len(value))*dataCopyGas + cachedStorageLoadGas + uint64(extraBytesForKey)*dataCopyGas
	} else {
		expectedUsedGas = 2 * (storageLoadGas + uint64(len(value))*dataCopyGas + uint64(extraBytesForKey)*dataCopyGas)
	}

	test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(0).
				WithConfig(nil).
				WithMethods(contracts.LoadStore)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(simpleGasTestConfig.GasProvided).
			WithFunction("loadStore").
			WithArguments(key).
			Build()).
		WithSetup(func(host arwen.VMHost, world *worldmock.MockWorld) {
			setZeroCodeCosts(host)
			host.Metering().GasSchedule().ElrondAPICost.StorageLoad = storageLoadGas
			host.Metering().GasSchedule().ElrondAPICost.CachedStorageLoad = cachedStorageLoadGas
			host.Metering().GasSchedule().BaseOperationCost.DataCopyPerByte = dataCopyGas
			host.Metering().GasSchedule().BaseOperationCost.PersistPerByte = 0

			if !flagEnabled {
				host.Storage().DisableUseDifferentGasCostFlag()
			}

			accountHandler, _ := world.GetUserAccount(test.ParentAddress)
			(accountHandler.(*worldmock.Account)).Storage[string(key)] = value
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasRemaining(simpleGasTestConfig.GasProvided-expectedUsedGas).
				GasUsed(test.ParentAddress, expectedUsedGas).
				ReturnData(value)
		})
}

func TestGasUsed_LoadStorageFromAddress_FlagEnabled(t *testing.T) {
	loadStorageFromAddress(t, true)
}

func TestGasUsed_LoadStorageFromAddress_FlagDisabled(t *testing.T) {
	loadStorageFromAddress(t, false)
}

func loadStorageFromAddress(t *testing.T, flagEnabled bool) {
	key := []byte("testKey")
	value := []byte("testValue")

	storageLoadGas := uint64(10)
	cachedStorageLoadGas := uint64(5)
	dataCopyGas := uint64(1)

	var expectedUsedGas uint64
	if flagEnabled {
		expectedUsedGas = storageLoadGas + uint64(len(value))*dataCopyGas + cachedStorageLoadGas
	} else {
		expectedUsedGas = 2 * (storageLoadGas + uint64(len(value))*dataCopyGas)
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
				WithMethods(contracts.LoadStoreFromAddress)).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(simpleGasTestConfig.GasProvided).
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

			if !flagEnabled {
				host.Storage().DisableUseDifferentGasCostFlag()
			}
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasRemaining(simpleGasTestConfig.GasProvided-expectedUsedGas).
				GasUsed(test.ParentAddress, expectedUsedGas).
				ReturnData(value)
		})
}

func TestGasUsed_SetStorage_FlagEnabled(t *testing.T) {
	setStorage(t, true)
}

func TestGasUsed_SetStorage_FlagDisabled(t *testing.T) {
	setStorage(t, false)
}

func setStorage(t *testing.T, flagEnabled bool) {
	key := []byte("testKey")
	value := []byte("testValue")

	storageStoreGas := uint64(10)
	dataCopyGas := uint64(1)

	var expectedUsedGas uint64
	if flagEnabled {
		expectedUsedGas = 2 * storageStoreGas
	} else {
		expectedUsedGas = 2*storageStoreGas + uint64(len(value))*dataCopyGas
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
			WithGasProvided(simpleGasTestConfig.GasProvided).
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

			if !flagEnabled {
				host.Storage().DisableUseDifferentGasCostFlag()
			}
		}).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.
				Ok().
				GasUsed(test.ParentAddress, expectedUsedGas).
				GasRemaining(simpleGasTestConfig.GasProvided - expectedUsedGas)
		})
}

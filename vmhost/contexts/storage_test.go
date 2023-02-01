package contexts

import (
	"bytes"
	"errors"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/config"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	worldmock "github.com/multiversx/mx-chain-vm-go/mock/world"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/require"
)

var reservedTestPrefix = []byte("RESERVED")

func TestNewStorageContext(t *testing.T) {
	t.Parallel()

	t.Run("empty protected key prefix should error", func(t *testing.T) {
		t.Parallel()

		host := &contextmock.VMHostMock{}
		mockBlockchain := worldmock.NewMockWorld()

		storageCtx, err := NewStorageContext(host, mockBlockchain, make([]byte, 0))
		require.Equal(t, vmhost.ErrEmptyProtectedKeyPrefix, err)
		require.True(t, check.IfNil(storageCtx))
	})
	t.Run("nil VM host should error", func(t *testing.T) {
		t.Parallel()

		mockBlockchain := worldmock.NewMockWorld()

		storageCtx, err := NewStorageContext(nil, mockBlockchain, reservedTestPrefix)
		require.Equal(t, vmhost.ErrNilVMHost, err)
		require.True(t, check.IfNil(storageCtx))
	})
	t.Run("nil blockchain hook should error", func(t *testing.T) {
		t.Parallel()

		host := &contextmock.VMHostMock{}

		storageCtx, err := NewStorageContext(host, nil, reservedTestPrefix)
		require.Equal(t, vmhost.ErrNilBlockChainHook, err)
		require.True(t, check.IfNil(storageCtx))
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		enableEpochsHandler := &worldmock.EnableEpochsHandlerStub{
			IsStorageAPICostOptimizationFlagEnabledField: true,
		}
		host := &contextmock.VMHostMock{
			EnableEpochsHandlerField: enableEpochsHandler,
		}
		mockBlockchain := worldmock.NewMockWorld()

		storageCtx, err := NewStorageContext(host, mockBlockchain, reservedTestPrefix)
		require.Nil(t, err)
		require.False(t, check.IfNil(storageCtx))
	})
}

func TestStorageContext_SetAddress(t *testing.T) {
	t.Parallel()

	addressA := []byte("accountA")
	addressB := []byte("accountB")
	stubOutput := &contextmock.OutputContextStub{}
	accountA := &vmcommon.OutputAccount{
		Address:        addressA,
		Nonce:          0,
		BalanceDelta:   big.NewInt(0),
		Balance:        big.NewInt(0),
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}
	accountB := &vmcommon.OutputAccount{
		Address:        addressB,
		Nonce:          0,
		BalanceDelta:   big.NewInt(0),
		Balance:        big.NewInt(0),
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}
	stubOutput.GetOutputAccountCalled = func(address []byte) (*vmcommon.OutputAccount, bool) {
		if bytes.Equal(address, addressA) {
			return accountA, false
		}
		if bytes.Equal(address, addressB) {
			return accountB, false
		}
		return nil, false
	}

	mockRuntime := &contextmock.RuntimeContextMock{}
	mockMetering := &contextmock.MeteringContextMock{}
	mockMetering.SetGasSchedule(config.MakeGasMapForTests())
	mockMetering.BlockGasLimitMock = uint64(15000)
	mockMetering.GasLeftMock = 20000
	enableEpochsHandler := &worldmock.EnableEpochsHandlerStub{
		IsStorageAPICostOptimizationFlagEnabledField: true,
	}

	host := &contextmock.VMHostMock{
		OutputContext:            stubOutput,
		MeteringContext:          mockMetering,
		RuntimeContext:           mockRuntime,
		EnableEpochsHandlerField: enableEpochsHandler,
	}
	bcHook := &contextmock.BlockchainHookStub{}

	storageCtx, _ := NewStorageContext(host, bcHook, reservedTestPrefix)

	keyA := []byte("keyA")
	valueA := []byte("valueA")

	storageCtx.SetAddress(addressA)
	storageStatus, err := storageCtx.SetStorage(keyA, valueA)
	require.Nil(t, err)
	require.Equal(t, vmhost.StorageAdded, storageStatus)
	require.Equal(t, uint64(len(valueA)), accountA.BytesAddedToStorage)
	require.Equal(t, uint64(0), accountA.BytesDeletedFromStorage)
	foundValueA, _, err := storageCtx.GetStorage(keyA)
	require.Nil(t, err)
	require.Equal(t, valueA, foundValueA)
	require.Len(t, storageCtx.GetStorageUpdates(addressA), 1)
	require.Len(t, storageCtx.GetStorageUpdates(addressB), 0)

	keyB := []byte("keyB")
	valueB := []byte("valueB")
	storageCtx.SetAddress(addressB)
	storageStatus, err = storageCtx.SetStorage(keyB, valueB)
	require.Equal(t, uint64(len(valueB)), accountB.BytesAddedToStorage)
	require.Equal(t, uint64(0), accountB.BytesDeletedFromStorage)
	require.Nil(t, err)
	require.Equal(t, vmhost.StorageAdded, storageStatus)
	foundValueB, _, err := storageCtx.GetStorage(keyB)
	require.Nil(t, err)
	require.Equal(t, valueB, foundValueB)
	require.Len(t, storageCtx.GetStorageUpdates(addressA), 1)
	require.Len(t, storageCtx.GetStorageUpdates(addressB), 1)
	foundValueA, _, err = storageCtx.GetStorage(keyA)
	require.Nil(t, err)
	require.Equal(t, []byte(nil), foundValueA)
}

func TestStorageContext_GetStorageUpdates(t *testing.T) {
	t.Parallel()

	mockOutput := &contextmock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount([]byte("account"))
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	account.StorageUpdates["update"] = &vmcommon.StorageUpdate{
		Offset: []byte("update"),
		Data:   []byte("some data"),
	}

	enableEpochsHandler := &worldmock.EnableEpochsHandlerStub{
		IsStorageAPICostOptimizationFlagEnabledField: true,
	}

	host := &contextmock.VMHostMock{
		OutputContext:            mockOutput,
		EnableEpochsHandlerField: enableEpochsHandler,
	}

	mockBlockchainHook := worldmock.NewMockWorld()
	storageCtx, _ := NewStorageContext(host, mockBlockchainHook, reservedTestPrefix)

	storageUpdates := storageCtx.GetStorageUpdates([]byte("account"))
	require.Equal(t, 1, len(storageUpdates))
	require.Equal(t, []byte("update"), storageUpdates["update"].Offset)
	require.Equal(t, []byte("some data"), storageUpdates["update"].Data)
}

func TestStorageContext_SetStorage(t *testing.T) {
	t.Parallel()

	address := []byte("account")
	mockOutput := &contextmock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount(address)
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	mockRuntime := &contextmock.RuntimeContextMock{}
	mockMetering := &contextmock.MeteringContextMock{}
	mockMetering.SetGasSchedule(config.MakeGasMapForTests())
	mockMetering.BlockGasLimitMock = uint64(15000)
	mockMetering.GasLeftMock = 20000

	enableEpochsHandler := &worldmock.EnableEpochsHandlerStub{
		IsStorageAPICostOptimizationFlagEnabledField: true,
	}

	host := &contextmock.VMHostMock{
		OutputContext:            mockOutput,
		MeteringContext:          mockMetering,
		RuntimeContext:           mockRuntime,
		EnableEpochsHandlerField: enableEpochsHandler,
	}
	bcHook := &contextmock.BlockchainHookStub{}
	storageCtx, _ := NewStorageContext(host, bcHook, reservedTestPrefix)
	storageCtx.SetAddress(address)

	val1 := []byte("value")
	val2 := []byte("newValue")
	val3 := []byte("v")

	key := []byte("key")
	value := val1
	addedBytes := len(value)

	storageStatus, err := storageCtx.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, vmhost.StorageAdded, storageStatus)
	require.Equal(t, uint64(addedBytes), account.BytesAddedToStorage)
	require.Equal(t, uint64(0), account.BytesDeletedFromStorage)
	foundValue, _, err := storageCtx.GetStorage(key)
	require.Nil(t, err)
	require.Equal(t, value, foundValue)
	require.Len(t, storageCtx.GetStorageUpdates(address), 1)

	value = val2
	addedBytes += len(value) - len(val1)

	storageStatus, err = storageCtx.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, vmhost.StorageModified, storageStatus)
	require.Equal(t, uint64(addedBytes), account.BytesAddedToStorage)
	require.Equal(t, uint64(0), account.BytesDeletedFromStorage)
	foundValue, _, err = storageCtx.GetStorage(key)
	require.Nil(t, err)
	require.Equal(t, value, foundValue)
	require.Len(t, storageCtx.GetStorageUpdates(address), 1)

	value = val2

	storageStatus, err = storageCtx.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, vmhost.StorageUnchanged, storageStatus)
	require.Equal(t, uint64(addedBytes), account.BytesAddedToStorage)
	require.Equal(t, uint64(0), account.BytesDeletedFromStorage)
	foundValue, _, err = storageCtx.GetStorage(key)
	require.Nil(t, err)
	require.Equal(t, value, foundValue)
	require.Len(t, storageCtx.GetStorageUpdates(address), 1)

	value = val1
	deletedBytes := len(val2) - len(val1)

	storageStatus, err = storageCtx.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, vmhost.StorageModified, storageStatus)
	require.Equal(t, uint64(addedBytes), account.BytesAddedToStorage)
	require.Equal(t, uint64(deletedBytes), account.BytesDeletedFromStorage)
	foundValue, _, err = storageCtx.GetStorage(key)
	require.Nil(t, err)
	require.Equal(t, value, foundValue)
	require.Len(t, storageCtx.GetStorageUpdates(address), 1)

	value = val3
	deletedBytes += len(val1) - len(val3)

	storageStatus, err = storageCtx.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, vmhost.StorageModified, storageStatus)
	require.Equal(t, uint64(addedBytes), account.BytesAddedToStorage)
	require.Equal(t, uint64(deletedBytes), account.BytesDeletedFromStorage)
	foundValue, _, err = storageCtx.GetStorage(key)
	require.Nil(t, err)
	require.Equal(t, value, foundValue)
	require.Len(t, storageCtx.GetStorageUpdates(address), 1)

	value = nil
	deletedBytes += len(val3)

	storageStatus, err = storageCtx.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, vmhost.StorageDeleted, storageStatus)
	require.Equal(t, uint64(addedBytes), account.BytesAddedToStorage)
	require.Equal(t, uint64(deletedBytes), account.BytesDeletedFromStorage)
	foundValue, _, err = storageCtx.GetStorage(key)
	require.Nil(t, err)
	require.Equal(t, []byte{}, foundValue)
	require.Len(t, storageCtx.GetStorageUpdates(address), 1)

	mockRuntime.SetReadOnly(true)
	value = val2
	storageStatus, err = storageCtx.SetStorage(key, value)
	require.Equal(t, err, vmhost.ErrCannotWriteOnReadOnly)
	require.Equal(t, vmhost.StorageUnchanged, storageStatus)
	foundValue, _, err = storageCtx.GetStorage(key)
	require.Nil(t, err)
	require.Equal(t, []byte{}, foundValue)
	require.Len(t, storageCtx.GetStorageUpdates(address), 1)

	mockRuntime.SetReadOnly(false)
	key = []byte("other_key")
	value = []byte("other_value")
	storageStatus, err = storageCtx.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, vmhost.StorageAdded, storageStatus)
	foundValue, _, err = storageCtx.GetStorage(key)
	require.Nil(t, err)
	require.Equal(t, value, foundValue)
	require.Len(t, storageCtx.GetStorageUpdates(address), 2)

	key = []byte("RESERVEDkey")
	value = []byte("doesn't matter")
	storageStatus, err = storageCtx.SetStorage(key, value)
	require.Equal(t, vmhost.StorageUnchanged, storageStatus)
	require.Equal(t, vmhost.ErrStoreReservedKey, err)

	key = []byte("RESERVED")
	value = []byte("doesn't matter")
	storageStatus, err = storageCtx.SetStorage(key, value)
	require.Equal(t, vmhost.StorageUnchanged, storageStatus)
	require.Equal(t, vmhost.ErrStoreReservedKey, err)
}

func TestStorageContext_SetStorage_GasUsage(t *testing.T) {
	address := []byte("account")
	mockOutput := &contextmock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount(address)
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	storeCost := 11
	persistCost := 7
	releaseCost := 5

	gasMap := config.MakeGasMapForTests()
	gasMap["BaseOperationCost"]["StorePerByte"] = uint64(storeCost)
	gasMap["BaseOperationCost"]["PersistPerByte"] = uint64(persistCost)
	gasMap["BaseOperationCost"]["ReleasePerByte"] = uint64(releaseCost)

	mockRuntime := &contextmock.RuntimeContextMock{}
	mockMetering := &contextmock.MeteringContextMock{}
	mockMetering.SetGasSchedule(gasMap)
	mockMetering.BlockGasLimitMock = uint64(15000)

	enableEpochsHandler := &worldmock.EnableEpochsHandlerStub{
		IsStorageAPICostOptimizationFlagEnabledField: true,
	}
	host := &contextmock.VMHostMock{
		OutputContext:            mockOutput,
		MeteringContext:          mockMetering,
		RuntimeContext:           mockRuntime,
		EnableEpochsHandlerField: enableEpochsHandler,
	}
	bcHook := &contextmock.BlockchainHookStub{}

	storageCtx, _ := NewStorageContext(host, bcHook, reservedTestPrefix)
	storageCtx.SetAddress(address)

	gasProvided := 100
	mockMetering.GasLeftMock = uint64(gasProvided)
	key := []byte("key")

	// Store new value
	value := []byte("value")
	storageStatus, err := storageCtx.SetStorage(key, value)
	gasLeft := gasProvided - storeCost*len(value)
	storedValue, _, err := storageCtx.GetStorage(key)
	require.Nil(t, err)
	require.Equal(t, vmhost.StorageAdded, storageStatus)
	require.Equal(t, gasLeft, int(mockMetering.GasLeft()))
	require.Equal(t, value, storedValue)

	// Update with longer value
	value2 := []byte("value2")
	mockMetering.GasLeftMock = uint64(gasProvided)
	storageStatus, err = storageCtx.SetStorage(key, value2)
	require.Nil(t, err)
	storedValue, _, err = storageCtx.GetStorage(key)
	require.Nil(t, err)
	gasLeft = gasProvided - persistCost*len(value) - storeCost*(len(value2)-len(value))
	require.Nil(t, err)
	require.Equal(t, vmhost.StorageModified, storageStatus)
	require.Equal(t, gasLeft, int(mockMetering.GasLeft()))
	require.Equal(t, value2, storedValue)

	// Revert to initial value
	mockMetering.GasLeftMock = uint64(gasProvided)
	storageStatus, err = storageCtx.SetStorage(key, value)
	require.Nil(t, err)
	gasLeft = gasProvided - persistCost*len(value)
	gasFreed := releaseCost * (len(value2) - len(value))
	storedValue, _, err = storageCtx.GetStorage(key)
	require.Nil(t, err)
	require.Equal(t, vmhost.StorageModified, storageStatus)
	require.Equal(t, gasLeft, int(mockMetering.GasLeft()))
	require.Equal(t, gasFreed, int(mockMetering.GasFreedMock))
	require.Equal(t, value, storedValue)
}

func TestStorageContext_StorageProtection(t *testing.T) {
	address := []byte("account")
	mockOutput := &contextmock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount(address)
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	mockRuntime := &contextmock.RuntimeContextMock{}
	mockMetering := &contextmock.MeteringContextMock{}
	mockMetering.SetGasSchedule(config.MakeGasMapForTests())
	mockMetering.BlockGasLimitMock = uint64(15000)
	mockMetering.GasLeftMock = 20000

	enableEpochsHandler := &worldmock.EnableEpochsHandlerStub{
		IsStorageAPICostOptimizationFlagEnabledField: true,
	}

	host := &contextmock.VMHostMock{
		OutputContext:            mockOutput,
		MeteringContext:          mockMetering,
		RuntimeContext:           mockRuntime,
		EnableEpochsHandlerField: enableEpochsHandler,
	}
	bcHook := &contextmock.BlockchainHookStub{}

	storageCtx, _ := NewStorageContext(host, bcHook, reservedTestPrefix)
	storageCtx.SetAddress(address)

	key := storageCtx.GetVmProtectedPrefix("something")
	value := []byte("data")

	storageStatus, err := storageCtx.SetStorage(key, value)
	require.Equal(t, vmhost.StorageUnchanged, storageStatus)
	require.True(t, errors.Is(err, vmhost.ErrCannotWriteProtectedKey))
	require.Len(t, storageCtx.GetStorageUpdates(address), 0)

	storageCtx.disableStorageProtection()
	protectedKey := append(reservedTestPrefix, []byte("ABC")...)
	storageStatus, err = storageCtx.SetStorage(protectedKey, value)
	require.Equal(t, vmhost.StorageUnchanged, storageStatus)
	require.True(t, errors.Is(err, vmhost.ErrStoreReservedKey))
	require.Len(t, storageCtx.GetStorageUpdates(address), 0)

	storageCtx.disableStorageProtection()
	storageStatus, err = storageCtx.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, vmhost.StorageAdded, storageStatus)
	require.Len(t, storageCtx.GetStorageUpdates(address), 1)

	storageCtx.enableStorageProtection()
	storageStatus, err = storageCtx.SetStorage(key, value)
	require.Equal(t, vmhost.StorageUnchanged, storageStatus)
	require.True(t, errors.Is(err, vmhost.ErrCannotWriteProtectedKey))
	require.Len(t, storageCtx.GetStorageUpdates(address), 1)
}

func TestStorageContext_GetStorageFromAddress(t *testing.T) {
	t.Parallel()

	scAddress := []byte("account")
	enableEpochsHandler := &worldmock.EnableEpochsHandlerStub{
		IsStorageAPICostOptimizationFlagEnabledField: true,
	}

	t.Run("blockchain hook errors", func(t *testing.T) {
		t.Parallel()

		mockOutput := &contextmock.OutputContextMock{}
		account := mockOutput.NewVMOutputAccount(scAddress)
		mockOutput.OutputAccountMock = account
		mockOutput.OutputAccountIsNew = false

		mockRuntime := &contextmock.RuntimeContextMock{}
		mockMetering := &contextmock.MeteringContextMock{}
		mockMetering.SetGasSchedule(config.MakeGasMapForTests())
		mockMetering.BlockGasLimitMock = uint64(15000)

		host := &contextmock.VMHostMock{
			OutputContext:            mockOutput,
			MeteringContext:          mockMetering,
			RuntimeContext:           mockRuntime,
			EnableEpochsHandlerField: enableEpochsHandler,
		}

		readable := []byte("readable")
		nonreadable := []byte("nonreadable")

		errTooManyRequests := errors.New("too many requests")
		bcHook := &contextmock.BlockchainHookStub{
			GetUserAccountCalled: func(address []byte) (vmcommon.UserAccountHandler, error) {
				if bytes.Equal(readable, address) {
					return &worldmock.Account{CodeMetadata: []byte{4, 0}}, nil
				}
				if bytes.Equal(nonreadable, address) || bytes.Equal(scAddress, address) {
					return &worldmock.Account{CodeMetadata: []byte{0, 0}}, nil
				}
				return nil, nil
			},
			GetStorageDataCalled: func(accountsAddress []byte, index []byte) ([]byte, uint32, error) {
				return nil, 0, errTooManyRequests
			},
		}

		storageCtx, _ := NewStorageContext(host, bcHook, reservedTestPrefix)
		storageCtx.SetAddress(scAddress)

		key := []byte("key")
		data, _, err := storageCtx.GetStorageFromAddress(scAddress, key)
		require.Nil(t, data)
		require.Equal(t, errTooManyRequests, err)

		data, _, _ = storageCtx.GetStorageFromAddress(readable, key)
		require.Nil(t, data)
		require.Equal(t, errTooManyRequests, err)

		data, _, _ = storageCtx.GetStorageFromAddress(nonreadable, key)
		require.Nil(t, data)
		require.Equal(t, errTooManyRequests, err)
	})
	t.Run("should work when blockchain hook does not error", func(t *testing.T) {
		t.Parallel()

		mockOutput := &contextmock.OutputContextMock{}
		account := mockOutput.NewVMOutputAccount(scAddress)
		mockOutput.OutputAccountMock = account
		mockOutput.OutputAccountIsNew = false

		mockRuntime := &contextmock.RuntimeContextMock{}
		mockMetering := &contextmock.MeteringContextMock{}
		mockMetering.SetGasSchedule(config.MakeGasMapForTests())
		mockMetering.BlockGasLimitMock = uint64(15000)

		host := &contextmock.VMHostMock{
			OutputContext:            mockOutput,
			MeteringContext:          mockMetering,
			RuntimeContext:           mockRuntime,
			EnableEpochsHandlerField: enableEpochsHandler,
		}

		readable := []byte("readable")
		nonreadable := []byte("nonreadable")
		internalData := []byte("internalData")

		bcHook := &contextmock.BlockchainHookStub{
			GetUserAccountCalled: func(address []byte) (vmcommon.UserAccountHandler, error) {
				if bytes.Equal(readable, address) {
					return &worldmock.Account{CodeMetadata: []byte{4, 0}}, nil
				}
				if bytes.Equal(nonreadable, address) || bytes.Equal(scAddress, address) {
					return &worldmock.Account{CodeMetadata: []byte{0, 0}}, nil
				}
				return nil, nil
			},
			GetStorageDataCalled: func(accountsAddress []byte, index []byte) ([]byte, uint32, error) {
				return internalData, 0, nil
			},
		}

		storageCtx, _ := NewStorageContext(host, bcHook, reservedTestPrefix)
		storageCtx.SetAddress(scAddress)

		key := []byte("key")
		data, _, err := storageCtx.GetStorageFromAddress(scAddress, key)
		require.Nil(t, err)
		require.Equal(t, data, internalData)

		data, _, err = storageCtx.GetStorageFromAddress(readable, key)
		require.Nil(t, err)
		require.Equal(t, data, internalData)

		data, _, err = storageCtx.GetStorageFromAddress(nonreadable, key)
		require.Nil(t, err)
		require.Nil(t, data)
	})
}

func TestStorageContext_LoadGasStoreGasPerKey(t *testing.T) {
	// TODO
}

func TestStorageContext_StoreGasPerKey(t *testing.T) {
	// TODO
}

func TestStorageContext_PopSetActiveStateIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()
	enableEpochsHandler := &worldmock.EnableEpochsHandlerStub{
		IsStorageAPICostOptimizationFlagEnabledField: true,
	}

	host := &contextmock.VMHostMock{
		EnableEpochsHandlerField: enableEpochsHandler,
	}

	storageCtx, _ := NewStorageContext(host, &contextmock.BlockchainHookStub{}, reservedTestPrefix)
	storageCtx.PopSetActiveState()

	require.Equal(t, 0, len(storageCtx.stateStack))
}

func TestStorageContext_PopDiscardIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()
	enableEpochsHandler := &worldmock.EnableEpochsHandlerStub{
		IsStorageAPICostOptimizationFlagEnabledField: true,
	}

	host := &contextmock.VMHostMock{
		EnableEpochsHandlerField: enableEpochsHandler,
	}

	storageCtx, _ := NewStorageContext(host, &contextmock.BlockchainHookStub{}, reservedTestPrefix)
	storageCtx.PopDiscard()

	require.Equal(t, 0, len(storageCtx.stateStack))
}

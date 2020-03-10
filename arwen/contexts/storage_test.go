package contexts

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	"github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

func TestNewStorageContext(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostMock{}
	mockBlockchain := &mock.BlockchainHookMock{}

	storageContext, err := NewStorageContext(host, mockBlockchain)
	require.Nil(t, err)
	require.NotNil(t, storageContext)
}

func TestStorageContext_GetStorageUpdates(t *testing.T) {
	t.Parallel()

	mockOutput := &mock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount([]byte("account"))
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	account.StorageUpdates["update"] = &vmcommon.StorageUpdate{
		Offset: []byte("update"),
		Data:   []byte("some data"),
	}

	host := &mock.VmHostMock{
		OutputContext: mockOutput,
	}

	mockBlockchainHook := &mock.BlockchainHookMock{}
	storageContext, _ := NewStorageContext(host, mockBlockchainHook)

	storageUpdates := storageContext.GetStorageUpdates([]byte("account"))
	require.Equal(t, 1, len(storageUpdates))
	require.Equal(t, []byte("update"), storageUpdates["update"].Offset)
	require.Equal(t, []byte("some data"), storageUpdates["update"].Data)
}

func TestStorageContext_SetStorage(t *testing.T) {
	t.Parallel()

	address := []byte("account")
	mockOutput := &mock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount(address)
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	mockRuntime := &mock.RuntimeContextMock{}
	mockMetering := &mock.MeteringContextMock{}
	mockMetering.SetGasSchedule(config.MakeGasMap(1))
	mockMetering.BlockGasLimitMock = uint64(15000)

	host := &mock.VmHostMock{
		OutputContext:   mockOutput,
		MeteringContext: mockMetering,
		RuntimeContext:  mockRuntime,
	}
	bcHook := &mock.BlockChainHookStub{}

	storageContext, _ := NewStorageContext(host, bcHook)

	key := []byte("key")
	value := []byte("value")

	storageStatus := storageContext.SetStorage(address, key, value)
	require.Equal(t, int32(arwen.StorageAdded), storageStatus)
	require.Equal(t, value, storageContext.GetStorage(address, key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	value = []byte("newValue")
	storageStatus = storageContext.SetStorage(address, key, value)
	require.Equal(t, int32(arwen.StorageModified), storageStatus)
	require.Equal(t, value, storageContext.GetStorage(address, key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	value = []byte("newValue")
	storageStatus = storageContext.SetStorage(address, key, value)
	require.Equal(t, int32(arwen.StorageUnchanged), storageStatus)
	require.Equal(t, value, storageContext.GetStorage(address, key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	value = nil
	storageStatus = storageContext.SetStorage(address, key, value)
	require.Equal(t, int32(arwen.StorageDeleted), storageStatus)
	require.Equal(t, []byte{}, storageContext.GetStorage(address, key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	mockRuntime.SetReadOnly(true)
	value = []byte("newValue")
	storageStatus = storageContext.SetStorage(address, key, value)
	require.Equal(t, int32(arwen.StorageUnchanged), storageStatus)
	require.Equal(t, []byte{}, storageContext.GetStorage(address, key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	mockRuntime.SetReadOnly(false)
	key = []byte("other_key")
	value = []byte("other_value")
	storageStatus = storageContext.SetStorage(address, key, value)
	require.Equal(t, int32(arwen.StorageAdded), storageStatus)
	require.Equal(t, value, storageContext.GetStorage(address, key))
	require.Len(t, storageContext.GetStorageUpdates(address), 2)
}

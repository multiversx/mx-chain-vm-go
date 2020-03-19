package contexts

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
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

func TestStorageContext_StateStack(t *testing.T) {
	t.Parallel()

	addressA := []byte("accountA")
	addressB := []byte("accountB")
	stubOutput := &mock.OutputContextStub{}
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

	mockRuntime := &mock.RuntimeContextMock{}
	mockMetering := &mock.MeteringContextMock{}
	mockMetering.SetGasSchedule(config.MakeGasMap(1))
	mockMetering.BlockGasLimitMock = uint64(15000)

	host := &mock.VmHostMock{
		OutputContext:   stubOutput,
		MeteringContext: mockMetering,
		RuntimeContext:  mockRuntime,
	}
	bcHook := &mock.BlockchainHookStub{}

	storageContext, _ := NewStorageContext(host, bcHook)

	keyA := []byte("keyA")
	valueA := []byte("valueA")

	storageContext.SetAddress(addressA)
	storageStatus := storageContext.SetStorage(keyA, valueA)
	require.Equal(t, int32(arwen.StorageAdded), storageStatus)
	require.Equal(t, valueA, storageContext.GetStorage(keyA))
	require.Len(t, storageContext.GetStorageUpdates(addressA), 1)
	require.Len(t, storageContext.GetStorageUpdates(addressB), 0)

	keyB := []byte("keyB")
	valueB := []byte("valueB")
	storageContext.SetAddress(addressB)
	storageStatus = storageContext.SetStorage(keyB, valueB)
	require.Equal(t, int32(arwen.StorageAdded), storageStatus)
	require.Equal(t, valueB, storageContext.GetStorage(keyB))
	require.Len(t, storageContext.GetStorageUpdates(addressA), 1)
	require.Len(t, storageContext.GetStorageUpdates(addressB), 1)
	require.Equal(t, []byte(nil), storageContext.GetStorage(keyA))
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
	bcHook := &mock.BlockchainHookStub{}

	storageContext, _ := NewStorageContext(host, bcHook)
	storageContext.SetAddress(address)

	key := []byte("key")
	value := []byte("value")

	storageStatus := storageContext.SetStorage(key, value)
	require.Equal(t, int32(arwen.StorageAdded), storageStatus)
	require.Equal(t, value, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	value = []byte("newValue")
	storageStatus = storageContext.SetStorage(key, value)
	require.Equal(t, int32(arwen.StorageModified), storageStatus)
	require.Equal(t, value, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	value = []byte("newValue")
	storageStatus = storageContext.SetStorage(key, value)
	require.Equal(t, int32(arwen.StorageUnchanged), storageStatus)
	require.Equal(t, value, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	value = nil
	storageStatus = storageContext.SetStorage(key, value)
	require.Equal(t, int32(arwen.StorageDeleted), storageStatus)
	require.Equal(t, []byte{}, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	mockRuntime.SetReadOnly(true)
	value = []byte("newValue")
	storageStatus = storageContext.SetStorage(key, value)
	require.Equal(t, int32(arwen.StorageUnchanged), storageStatus)
	require.Equal(t, []byte{}, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	mockRuntime.SetReadOnly(false)
	key = []byte("other_key")
	value = []byte("other_value")
	storageStatus = storageContext.SetStorage(key, value)
	require.Equal(t, int32(arwen.StorageAdded), storageStatus)
	require.Equal(t, value, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 2)
}

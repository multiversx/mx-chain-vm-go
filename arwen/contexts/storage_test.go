package contexts

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	"github.com/stretchr/testify/require"
)

func createMeteringContext() *meteringContext {
	runtimeMock := mock.NewRuntimeContextMock()
	outputContext := mock.NewOutputContextMock()
	host := &mock.VmHostStub{
		RuntimeCalled: func() arwen.RuntimeContext {
			return runtimeMock
		},
		OutputCalled: func() arwen.OutputContext {
			return outputContext
		},
	}
	blockGasLimit := uint64(15000)
	gasSchedule := config.MakeGasMap(1)

	metContext, _ := NewMeteringContext(host, gasSchedule, blockGasLimit)
	return metContext
}

func TestNewStorageContext(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostStub{}
	bcHook := &mock.BlockChainHookStub{}

	storageContext, err := NewStorageContext(host, bcHook)
	require.Nil(t, err)
	require.NotNil(t, storageContext)
}

func TestStorageContext_SetStorage(t *testing.T) {
	t.Parallel()

	outputContext := mock.NewOutputContextMock()
	meteringContext := createMeteringContext()
	runtimeMock := mock.NewRuntimeContextMock()

	host := &mock.VmHostStub{
		RuntimeCalled: func() arwen.RuntimeContext {
			return runtimeMock
		},
		MeteringCalled: func() arwen.MeteringContext {
			return meteringContext
		},
		OutputCalled: func() arwen.OutputContext {
			return outputContext
		},
	}
	bcHook := &mock.BlockChainHookStub{}

	storageContext, _ := NewStorageContext(host, bcHook)

	address := []byte("addr")
	key := []byte("key")
	//secondKey := []byte("sKey")
	value := []byte("value")

	storageStatus := storageContext.SetStorage(address, key, value)
	require.Equal(t, int32(arwen.StorageAdded), storageStatus)

	value = []byte("newValue")
	storageStatus = storageContext.SetStorage(address, key, value)
	require.Equal(t, int32(arwen.StorageModified), storageStatus)

	value = []byte("newValue")
	storageStatus = storageContext.SetStorage(address, key, value)
	require.Equal(t, int32(arwen.StorageUnchanged), storageStatus)

	value = nil
	storageStatus = storageContext.SetStorage(address, key, value)
	require.Equal(t, int32(arwen.StorageDeleted), storageStatus)

	runtimeMock.SetReadOnly(true)
	value = nil
	storageStatus = storageContext.SetStorage(address, key, value)
	require.Equal(t, int32(arwen.StorageUnchanged), storageStatus)
}

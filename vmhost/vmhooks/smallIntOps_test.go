package vmhooks

import (
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/mock/mockery"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	twos "github.com/multiversx/mx-components-big-int/twos-complement"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createTestVMHooks() (*VMHooksImpl, *mockery.MockVMHost, *mockery.MockRuntimeContext, *mockery.MockMeteringContext, *mockery.MockOutputContext, *mockery.MockStorageContext) {
	host := &mockery.MockVMHost{}
	runtime := &mockery.MockRuntimeContext{}
	metering := &mockery.MockMeteringContext{}
	output := &mockery.MockOutputContext{}
	storage := &mockery.MockStorageContext{}
	instance := &mockery.MockInstance{}

	host.On("Runtime").Return(runtime)
	host.On("Metering").Return(metering)
	host.On("Output").Return(output)
	host.On("Storage").Return(storage)
	runtime.On("GetInstance").Return(instance)
	instance.On("MemLoad", mock.Anything, mock.Anything).Return(nil, nil)
	instance.On("MemStore", mock.Anything, mock.Anything).Return(nil)

	gasSchedule, _ := config.CreateGasConfig(config.MakeGasMapForTests())
	metering.On("GasSchedule").Return(gasSchedule)
	metering.On("StartGasTracing", mock.Anything)
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	hooks := NewVMHooksImpl(host)
	return hooks, host, runtime, metering, output, storage
}

func TestVMHooksImpl_SmallIntGetUnsignedArgument(t *testing.T) {
	t.Parallel()

	hooks, _, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	t.Run("should work", func(t *testing.T) {
		runtime.On("Arguments").Return([][]byte{big.NewInt(42).Bytes()})
		val := hooks.SmallIntGetUnsignedArgument(0)
		require.Equal(t, int64(42), val)
	})

	t.Run("out of range", func(t *testing.T) {
		runtime.On("Arguments").Return([][]byte{})
		hooks.SmallIntGetUnsignedArgument(0)
		// expect fail execution
	})
}

func TestVMHooksImpl_SmallIntGetSignedArgument(t *testing.T) {
	t.Parallel()

	hooks, _, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	t.Run("should work", func(t *testing.T) {
		runtime.On("Arguments").Return([][]byte{twos.ToBytes(big.NewInt(-42))})
		val := hooks.SmallIntGetSignedArgument(0)
		require.Equal(t, int64(-42), val)
	})
}

func TestVMHooksImpl_SmallIntFinishUnsigned(t *testing.T) {
	t.Parallel()

	hooks, _, _, metering, output, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)
	output.On("Finish", mock.Anything).Return()

	hooks.SmallIntFinishUnsigned(42)
	output.AssertCalled(t, "Finish", big.NewInt(42).Bytes())
}

func TestVMHooksImpl_SmallIntFinishSigned(t *testing.T) {
	t.Parallel()

	hooks, _, _, metering, output, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)
	output.On("Finish", mock.Anything).Return()

	hooks.SmallIntFinishSigned(-42)
}

func TestVMHooksImpl_SmallIntStorageStoreUnsigned(t *testing.T) {
	t.Parallel()

	hooks, _, _, metering, _, storage := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)
	storage.On("SetStorage", mock.Anything, mock.Anything).Return(vmhost.StorageAdded, nil)
	hooks.SmallIntStorageStoreUnsigned(0, 0, 42)
}

func TestVMHooksImpl_SmallIntStorageStoreSigned(t *testing.T) {
	t.Parallel()

	hooks, _, _, metering, _, storage := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)
	storage.On("SetStorage", mock.Anything, mock.Anything).Return(vmhost.StorageAdded, nil)
	hooks.SmallIntStorageStoreSigned(0, 0, -42)
}

func TestVMHooksImpl_SmallIntStorageLoadUnsigned(t *testing.T) {
	t.Parallel()

	hooks, _, _, metering, _, storage := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)
	storage.On("GetStorage", mock.Anything).Return(big.NewInt(42).Bytes(), uint32(0), false, nil)
	storage.On("UseGasForStorageLoad", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	val := hooks.SmallIntStorageLoadUnsigned(0, 0)
	require.Equal(t, int64(42), val)
}

func TestVMHooksImpl_SmallIntStorageLoadSigned(t *testing.T) {
	t.Parallel()

	hooks, _, _, metering, _, storage := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)
	storage.On("GetStorage", mock.Anything).Return(twos.ToBytes(big.NewInt(-42)), uint32(0), false, nil)
	storage.On("UseGasForStorageLoad", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	val := hooks.SmallIntStorageLoadSigned(0, 0)
	require.Equal(t, int64(-42), val)
}

package vmhooks

import (
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	"github.com/multiversx/mx-chain-vm-go/mock/mockery"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestVMHooksImpl_MBufferNew(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("NewManagedBuffer").Return(int32(1))

	ret := hooks.MBufferNew()
	require.Equal(t, int32(1), ret)
}

func TestVMHooksImpl_MBufferNewFromBytes(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("NewManagedBufferFromBytes", mock.Anything).Return(int32(1))

	ret := hooks.MBufferNewFromBytes(0, 0)
	require.Equal(t, int32(1), ret)
}

func TestVMHooksImpl_MBufferGetLength(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("GetLength", mock.Anything).Return(int32(10))

	ret := hooks.MBufferGetLength(0)
	require.Equal(t, int32(10), ret)
}

func TestVMHooksImpl_MBufferGetBytes(t *testing.T) {
	t.Parallel()
	hooks, host, _, _, _, _ := createTestVMHooks()

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)

	ret := hooks.MBufferGetBytes(0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MBufferGetByteSlice(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)

	ret := hooks.MBufferGetByteSlice(0, 0, 4, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MBufferCopyByteSlice(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	ret := hooks.MBufferCopyByteSlice(0, 0, 4, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MBufferEq(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)

	ret := hooks.MBufferEq(0, 0)
	require.Equal(t, int32(1), ret)
}

func TestVMHooksImpl_MBufferSetBytes(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	ret := hooks.MBufferSetBytes(0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MBufferSetByteSlice(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("GetBytes", mock.Anything).Return(make([]byte, 10), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	ret := hooks.MBufferSetByteSlice(0, 0, 4, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MBufferAppend(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	managedType.On("AppendBytes", mock.Anything, mock.Anything).Return(true)

	ret := hooks.MBufferAppend(0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MBufferAppendBytes(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("AppendBytes", mock.Anything, mock.Anything).Return(true)

	ret := hooks.MBufferAppendBytes(0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MBufferToBigIntUnsigned(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	enableEpochs := &worldmock.EnableEpochsHandlerStub{}
	host.On("EnableEpochsHandler").Return(enableEpochs)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("GetBigIntOrCreate", mock.Anything).Return(big.NewInt(0))

	ret := hooks.MBufferToBigIntUnsigned(0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MBufferToBigIntSigned(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	enableEpochs := &worldmock.EnableEpochsHandlerStub{}
	host.On("EnableEpochsHandler").Return(enableEpochs)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("GetBigIntOrCreate", mock.Anything).Return(big.NewInt(0))

	ret := hooks.MBufferToBigIntSigned(0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MBufferFromBigIntUnsigned(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("GetBigInt", mock.Anything).Return(big.NewInt(100), nil)
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	ret := hooks.MBufferFromBigIntUnsigned(0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MBufferFromBigIntSigned(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("GetBigInt", mock.Anything).Return(big.NewInt(100), nil)
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	ret := hooks.MBufferFromBigIntSigned(0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MBufferToSmallIntUnsigned(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("GetBytes", mock.Anything).Return(big.NewInt(100).Bytes(), nil)

	ret := hooks.MBufferToSmallIntUnsigned(0)
	require.Equal(t, int64(100), ret)
}

func TestVMHooksImpl_MBufferToSmallIntSigned(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("GetBytes", mock.Anything).Return(big.NewInt(-100).Bytes(), nil)

	ret := hooks.MBufferToSmallIntSigned(0)
	require.Equal(t, int64(-100), ret)
}

func TestVMHooksImpl_MBufferFromSmallIntUnsigned(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	hooks.MBufferFromSmallIntUnsigned(0, 100)
}

func TestVMHooksImpl_MBufferFromSmallIntSigned(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	hooks.MBufferFromSmallIntSigned(0, -100)
}

func TestVMHooksImpl_MBufferStorageStore(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, storage := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	storage.On("SetStorage", mock.Anything, mock.Anything).Return(uint32(0), nil)

	ret := hooks.MBufferStorageStore(0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MBufferStorageLoad(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, storage := createTestVMHooks()

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("GetBytes", mock.Anything).Return([]byte("key"), nil)
	storage.On("GetStorage", mock.Anything).Return([]byte("data"), uint32(0), false, nil)
	storage.On("UseGasForStorageLoad", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	ret := hooks.MBufferStorageLoad(0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MBufferGetArgument(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	runtime.On("Arguments").Return([][]byte{[]byte("arg1")})
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	ret := hooks.MBufferGetArgument(0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MBufferFinish(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, output, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	output.On("Finish", mock.Anything).Return()

	ret := hooks.MBufferFinish(0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MBufferSetRandom(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	managedType.On("GetRandReader").Return(nil)
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	ret := hooks.MBufferSetRandom(0, 10)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MBufferToBigFloat(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	enableEpochs := &worldmock.EnableEpochsHandlerStub{}
	host.On("EnableEpochsHandler").Return(enableEpochs)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	managedType.On("EncodedBigFloatIsNotValid", mock.Anything).Return(false)
	managedType.On("GetBigFloatOrCreate", mock.Anything).Return(big.NewFloat(0), nil)

	ret := hooks.MBufferToBigFloat(0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MBufferFromBigFloat(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	enableEpochs := &worldmock.EnableEpochsHandlerStub{}
	host.On("EnableEpochsHandler").Return(enableEpochs)
	managedType.On("GetBigFloat", mock.Anything).Return(big.NewFloat(0), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	managedType.On("SetBytes", mock.Anything, mock.Anything).Return()

	ret := hooks.MBufferFromBigFloat(0, 0)
	require.Equal(t, int32(0), ret)
}

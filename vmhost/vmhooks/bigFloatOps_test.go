package vmhooks

import (
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	"github.com/multiversx/mx-chain-vm-go/mock/mockery"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestVMHooksImpl_BigFloatAdd(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)

	dest := big.NewFloat(0)
	managedType.On("GetBigFloatOrCreate", int32(0)).Return(dest, nil)
	managedType.On("GetTwoBigFloats", int32(1), int32(2)).Return(big.NewFloat(10.5), big.NewFloat(20.5), nil)
	managedType.On("BigFloatExpIsNotValid", mock.Anything).Return(false)

	hooks.BigFloatAdd(0, 1, 2)
	f, _ := dest.Float64()
	require.Equal(t, float64(31), f)
}

func TestVMHooksImpl_BigFloatSub(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)

	dest := big.NewFloat(0)
	managedType.On("GetBigFloatOrCreate", int32(0)).Return(dest, nil)
	managedType.On("GetTwoBigFloats", int32(1), int32(2)).Return(big.NewFloat(20.5), big.NewFloat(10.5), nil)
	managedType.On("BigFloatExpIsNotValid", mock.Anything).Return(false)

	hooks.BigFloatSub(0, 1, 2)
	f, _ := dest.Float64()
	require.Equal(t, float64(10), f)
}

func TestVMHooksImpl_BigFloatMul(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)

	dest := big.NewFloat(0)
	managedType.On("GetBigFloatOrCreate", int32(0)).Return(dest, nil)
	managedType.On("GetTwoBigFloats", int32(1), int32(2)).Return(big.NewFloat(10.5), big.NewFloat(2.0), nil)
	managedType.On("BigFloatExpIsNotValid", mock.Anything).Return(false)

	hooks.BigFloatMul(0, 1, 2)
	f, _ := dest.Float64()
	require.Equal(t, float64(21), f)
}

func TestVMHooksImpl_BigFloatDiv(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	hooks.GetVMHost().(*mockery.MockVMHost).On("ManagedTypes").Return(managedType)

	dest := big.NewFloat(0)
	managedType.On("GetBigFloatOrCreate", int32(0)).Return(dest, nil)
	managedType.On("GetTwoBigFloats", int32(1), int32(2)).Return(big.NewFloat(21.0), big.NewFloat(2.0), nil)
	managedType.On("BigFloatExpIsNotValid", mock.Anything).Return(false)

	hooks.BigFloatDiv(0, 1, 2)
	f, _ := dest.Float64()
	require.Equal(t, float64(10.5), f)
}

func TestVMHooksImpl_BigFloatAbs(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	hooks.GetVMHost().(*mockery.MockVMHost).On("ManagedTypes").Return(managedType)

	dest := big.NewFloat(0)
	managedType.On("GetBigFloatOrCreate", int32(0)).Return(dest, nil)
	managedType.On("GetBigFloat", int32(1)).Return(big.NewFloat(-10.5), nil)

	hooks.BigFloatAbs(0, 1)
	f, _ := dest.Float64()
	require.Equal(t, float64(10.5), f)
}

func TestVMHooksImpl_BigFloatNeg(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	hooks.GetVMHost().(*mockery.MockVMHost).On("ManagedTypes").Return(managedType)

	dest := big.NewFloat(0)
	managedType.On("GetBigFloatOrCreate", int32(0)).Return(dest, nil)
	managedType.On("GetBigFloat", int32(1)).Return(big.NewFloat(10.5), nil)

	hooks.BigFloatNeg(0, 1)
	f, _ := dest.Float64()
	require.Equal(t, float64(-10.5), f)
}

func TestVMHooksImpl_BigFloatCmp1(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	hooks.GetVMHost().(*mockery.MockVMHost).On("ManagedTypes").Return(managedType)

	managedType.On("GetTwoBigFloats", int32(1), int32(2)).Return(big.NewFloat(10.5), big.NewFloat(20.5), nil)
	res := hooks.BigFloatCmp(1, 2)
	require.Equal(t, int32(-1), res)
}

func TestVMHooksImpl_BigFloatCmp2(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	hooks.GetVMHost().(*mockery.MockVMHost).On("ManagedTypes").Return(managedType)

	managedType.On("GetTwoBigFloats", int32(1), int32(2)).Return(big.NewFloat(20.5), big.NewFloat(10.5), nil)
	res := hooks.BigFloatCmp(1, 2)
	require.Equal(t, int32(1), res)
}

func TestVMHooksImpl_BigFloatCmp3(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	hooks.GetVMHost().(*mockery.MockVMHost).On("ManagedTypes").Return(managedType)

	managedType.On("GetTwoBigFloats", int32(1), int32(2)).Return(big.NewFloat(10.5), big.NewFloat(10.5), nil)
	res := hooks.BigFloatCmp(1, 2)
	require.Equal(t, int32(0), res)
}

func TestVMHooksImpl_BigFloatSqrt(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	hooks.GetVMHost().(*mockery.MockVMHost).On("ManagedTypes").Return(managedType)

	dest := big.NewFloat(0)
	managedType.On("GetBigFloatOrCreate", int32(0)).Return(dest, nil)
	managedType.On("GetBigFloat", int32(1)).Return(big.NewFloat(9.0), nil)
	managedType.On("BigFloatExpIsNotValid", mock.Anything).Return(false)

	hooks.BigFloatSqrt(0, 1)
	f, _ := dest.Float64()
	require.Equal(t, float64(3), f)
}

func TestVMHooksImpl_BigFloatPow(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	hooks.GetVMHost().(*mockery.MockVMHost).On("ManagedTypes").Return(managedType)
	hooks.GetVMHost().(*mockery.MockVMHost).On("EnableEpochsHandler").Return(&worldmock.EnableEpochsHandlerStub{})

	dest := big.NewFloat(0)
	managedType.On("GetBigFloatOrCreate", int32(0)).Return(dest, nil)
	managedType.On("GetBigFloat", int32(1)).Return(big.NewFloat(2.0), nil)
	managedType.On("BigFloatExpIsNotValid", mock.Anything).Return(false)
	managedType.On("ConsumeGasForThisBigIntNumberOfBytes", mock.Anything).Return(nil)

	hooks.BigFloatPow(0, 1, 10)
	f, _ := dest.Float64()
	require.Equal(t, float64(1024), f)
}

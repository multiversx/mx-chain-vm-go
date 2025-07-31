package vmhooks

import (
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestVMHooksImpl_BigIntAdd(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &context.MockManagedTypesContext{}
	hooks.host.(*context.MockVMHost).On("ManagedTypes").Return(managedType)

	dest := big.NewInt(0)
	managedType.On("GetBigIntOrCreate", int32(0)).Return(dest)
	managedType.On("GetTwoBigInt", int32(1), int32(2)).Return(big.NewInt(10), big.NewInt(20), nil)
	managedType.On("ConsumeGasForBigIntCopy", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	hooks.BigIntAdd(0, 1, 2)
	require.Equal(t, big.NewInt(30), dest)
}

func TestVMHooksImpl_BigIntSub(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &context.MockManagedTypesContext{}
	hooks.host.(*context.MockVMHost).On("ManagedTypes").Return(managedType)

	dest := big.NewInt(0)
	managedType.On("GetBigIntOrCreate", int32(0)).Return(dest)
	managedType.On("GetTwoBigInt", int32(1), int32(2)).Return(big.NewInt(20), big.NewInt(10), nil)
	managedType.On("ConsumeGasForBigIntCopy", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	hooks.BigIntSub(0, 1, 2)
	require.Equal(t, big.NewInt(10), dest)
}

func TestVMHooksImpl_BigIntMul(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &context.MockManagedTypesContext{}
	hooks.host.(*context.MockVMHost).On("ManagedTypes").Return(managedType)

	dest := big.NewInt(0)
	managedType.On("GetBigIntOrCreate", int32(0)).Return(dest)
	managedType.On("GetTwoBigInt", int32(1), int32(2)).Return(big.NewInt(10), big.NewInt(20), nil)
	managedType.On("ConsumeGasForBigIntCopy", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	hooks.BigIntMul(0, 1, 2)
	require.Equal(t, big.NewInt(200), dest)
}

func TestVMHooksImpl_BigIntTDiv(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &context.MockManagedTypesContext{}
	hooks.host.(*context.MockVMHost).On("ManagedTypes").Return(managedType)

	dest := big.NewInt(0)
	managedType.On("GetBigIntOrCreate", int32(0)).Return(dest)
	managedType.On("GetTwoBigInt", int32(1), int32(2)).Return(big.NewInt(20), big.NewInt(10), nil)
	managedType.On("ConsumeGasForBigIntCopy", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	hooks.BigIntTDiv(0, 1, 2)
	require.Equal(t, big.NewInt(2), dest)
}

func TestVMHooksImpl_BigIntTMod(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &context.MockManagedTypesContext{}
	hooks.host.(*context.MockVMHost).On("ManagedTypes").Return(managedType)

	dest := big.NewInt(0)
	managedType.On("GetBigIntOrCreate", int32(0)).Return(dest)
	managedType.On("GetTwoBigInt", int32(1), int32(2)).Return(big.NewInt(23), big.NewInt(10), nil)
	managedType.On("ConsumeGasForBigIntCopy", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	hooks.BigIntTMod(0, 1, 2)
	require.Equal(t, big.NewInt(3), dest)
}

func TestVMHooksImpl_BigIntEDiv(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &context.MockManagedTypesContext{}
	hooks.host.(*context.MockVMHost).On("ManagedTypes").Return(managedType)

	dest := big.NewInt(0)
	managedType.On("GetBigIntOrCreate", int32(0)).Return(dest)
	managedType.On("GetTwoBigInt", int32(1), int32(2)).Return(big.NewInt(20), big.NewInt(10), nil)
	managedType.On("ConsumeGasForBigIntCopy", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	hooks.BigIntEDiv(0, 1, 2)
	require.Equal(t, big.NewInt(2), dest)
}

func TestVMHooksImpl_BigIntEMod(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &context.MockManagedTypesContext{}
	hooks.host.(*context.MockVMHost).On("ManagedTypes").Return(managedType)

	dest := big.NewInt(0)
	managedType.On("GetBigIntOrCreate", int32(0)).Return(dest)
	managedType.On("GetTwoBigInt", int32(1), int32(2)).Return(big.NewInt(23), big.NewInt(10), nil)
	managedType.On("ConsumeGasForBigIntCopy", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	hooks.BigIntEMod(0, 1, 2)
	require.Equal(t, big.NewInt(3), dest)
}

func TestVMHooksImpl_BigIntPow(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &context.MockManagedTypesContext{}
	hooks.host.(*context.MockVMHost).On("ManagedTypes").Return(managedType)

	dest := big.NewInt(0)
	managedType.On("GetBigIntOrCreate", int32(0)).Return(dest)
	managedType.On("GetTwoBigInt", int32(1), int32(2)).Return(big.NewInt(2), big.NewInt(10), nil)
	managedType.On("ConsumeGasForThisBigIntNumberOfBytes", mock.Anything).Return(nil)
	managedType.On("ConsumeGasForBigIntCopy", mock.Anything, mock.Anything).Return(nil)

	hooks.BigIntPow(0, 1, 2)
	require.Equal(t, big.NewInt(1024), dest)
}

func TestVMHooksImpl_BigIntSqrt(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &context.MockManagedTypesContext{}
	hooks.host.(*context.MockVMHost).On("ManagedTypes").Return(managedType)

	dest := big.NewInt(0)
	managedType.On("GetBigIntOrCreate", int32(0)).Return(dest)
	managedType.On("GetBigInt", int32(1)).Return(big.NewInt(1024), nil)
	managedType.On("ConsumeGasForBigIntCopy", mock.Anything, mock.Anything).Return(nil)

	hooks.BigIntSqrt(0, 1)
	require.Equal(t, big.NewInt(32), dest)
}

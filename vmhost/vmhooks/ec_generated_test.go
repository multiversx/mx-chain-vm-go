package vmhooks

import (
	"testing"

	"github.com/multiversx/mx-chain-crypto-go/zk/lowLevelFeatures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestManagedAddEC_Fail_GetPoint1Bytes(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("GetBytes", int32(1)).Return(nil, assert.AnError)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedAddEC(int32(lowLevelFeatures.BN254), int32(lowLevelFeatures.G1), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedPairingChecksEC_Fail_ReadVectors(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("ReadManagedVecOfManagedBuffers", int32(1)).Return(nil, uint64(0), assert.AnError)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedPairingChecksEC(int32(lowLevelFeatures.BN254), 1, 2)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedPairingChecksEC_Fail_InvalidCurve(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("ReadManagedVecOfManagedBuffers", int32(1)).Return([][]byte{}, uint64(0), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", int32(2)).Return([][]byte{}, uint64(0), nil)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedPairingChecksEC(int32(42), 1, 2)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedPairingChecksEC_Fail_PairingCheck(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("ReadManagedVecOfManagedBuffers", int32(1)).Return([][]byte{[]byte("dummy point")}, uint64(1), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", int32(2)).Return([][]byte{[]byte("dummy point")}, uint64(1), nil)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedPairingChecksEC(int32(lowLevelFeatures.BN254), 1, 2)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedMultiExpEC_Fail_ReadVectors(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("ReadManagedVecOfManagedBuffers", int32(1)).Return(nil, uint64(0), assert.AnError)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedMultiExpEC(int32(lowLevelFeatures.BN254), int32(lowLevelFeatures.G1), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedMultiExpEC_Fail_InvalidCurve(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("ReadManagedVecOfManagedBuffers", int32(1)).Return([][]byte{}, uint64(0), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", int32(2)).Return([][]byte{}, uint64(0), nil)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedMultiExpEC(int32(42), int32(lowLevelFeatures.G1), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedMultiExpEC_Fail_InvalidGroup(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("ReadManagedVecOfManagedBuffers", int32(1)).Return([][]byte{}, uint64(0), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", int32(2)).Return([][]byte{}, uint64(0), nil)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedMultiExpEC(int32(lowLevelFeatures.BN254), int32(42), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedMultiExpEC_Fail_MultiExp(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("ReadManagedVecOfManagedBuffers", int32(1)).Return([][]byte{[]byte("dummy point")}, uint64(1), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", int32(2)).Return([][]byte{[]byte("dummy scalar")}, uint64(1), nil)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedMultiExpEC(int32(lowLevelFeatures.BN254), int32(lowLevelFeatures.G1), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedMapToCurveEC_Fail_GetElementBytes(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("GetBytes", int32(1)).Return(nil, assert.AnError)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedMapToCurveEC(int32(lowLevelFeatures.BN254), int32(lowLevelFeatures.G1), 1, 2)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedMapToCurveEC_Fail_InvalidCurve(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("GetBytes", int32(1)).Return([]byte("dummy element"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedMapToCurveEC(int32(42), int32(lowLevelFeatures.G1), 1, 2)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedMapToCurveEC_Fail_InvalidGroup(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("GetBytes", int32(1)).Return([]byte("dummy element"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedMapToCurveEC(int32(lowLevelFeatures.BN254), int32(42), 1, 2)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedMapToCurveEC_Fail_Map(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("GetBytes", int32(1)).Return([]byte("dummy element"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedMapToCurveEC(int32(lowLevelFeatures.BN254), int32(lowLevelFeatures.G1), 1, 2)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedMulEC_Fail_GetPointBytes(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("GetBytes", int32(1)).Return(nil, assert.AnError)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedMulEC(int32(lowLevelFeatures.BN254), int32(lowLevelFeatures.G1), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedMulEC_Fail_GetScalarBytes(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("GetBytes", int32(1)).Return([]byte("dummy point"), nil)
	managedType.On("GetBytes", int32(2)).Return(nil, assert.AnError)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedMulEC(int32(lowLevelFeatures.BN254), int32(lowLevelFeatures.G1), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedMulEC_Fail_InvalidCurve(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("GetBytes", int32(1)).Return([]byte("dummy point"), nil)
	managedType.On("GetBytes", int32(2)).Return([]byte("dummy scalar"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil).Times(2)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedMulEC(int32(42), int32(lowLevelFeatures.G1), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedMulEC_Fail_InvalidGroup(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("GetBytes", int32(1)).Return([]byte("dummy point"), nil)
	managedType.On("GetBytes", int32(2)).Return([]byte("dummy scalar"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil).Times(2)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedMulEC(int32(lowLevelFeatures.BN254), int32(42), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedMulEC_Fail_Mul(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("GetBytes", int32(1)).Return([]byte("dummy point"), nil)
	managedType.On("GetBytes", int32(2)).Return([]byte("dummy scalar"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil).Times(2)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedMulEC(int32(lowLevelFeatures.BN254), int32(lowLevelFeatures.G1), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedAddEC_Fail_GetPoint2Bytes(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("GetBytes", int32(1)).Return([]byte("dummy point"), nil)
	managedType.On("GetBytes", int32(2)).Return(nil, assert.AnError)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedAddEC(int32(lowLevelFeatures.BN254), int32(lowLevelFeatures.G1), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedAddEC_Fail_InvalidCurve(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("GetBytes", int32(1)).Return([]byte("dummy point 1"), nil)
	managedType.On("GetBytes", int32(2)).Return([]byte("dummy point 2"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil).Times(2)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedAddEC(int32(42), int32(lowLevelFeatures.G1), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedAddEC_Fail_InvalidGroup(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("GetBytes", int32(1)).Return([]byte("dummy point 1"), nil)
	managedType.On("GetBytes", int32(2)).Return([]byte("dummy point 2"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil).Times(2)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedAddEC(int32(lowLevelFeatures.BN254), int32(42), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedAddEC_Fail_Add(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("GetBytes", int32(1)).Return([]byte("dummy point 1"), nil)
	managedType.On("GetBytes", int32(2)).Return([]byte("dummy point 2"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil).Times(2)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedAddEC(int32(lowLevelFeatures.BN254), int32(lowLevelFeatures.G1), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

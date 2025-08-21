package vmhooks

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestManagedVerifyGroth16_Fail_GetProofBytes(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime
	managedType.On("GetBytes", int32(1)).Return(nil, assert.AnError)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedVerifyGroth16(int32(ecc.BLS12_381), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedVerifyGroth16_Fail_GetVKBytes(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime
	managedType.On("GetBytes", int32(1)).Return([]byte("dummy proof"), nil)
	managedType.On("GetBytes", int32(2)).Return(nil, assert.AnError)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedVerifyGroth16(int32(ecc.BLS12_381), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedVerifyGroth16_Fail_GetPubWitnessBytes(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime
	managedType.On("GetBytes", int32(1)).Return([]byte("dummy proof"), nil)
	managedType.On("GetBytes", int32(2)).Return([]byte("dummy vk"), nil)
	managedType.On("GetBytes", int32(3)).Return(nil, assert.AnError)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil).Times(2)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedVerifyGroth16(int32(ecc.BLS12_381), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedVerifyGroth16_Fail_Verification(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime
	managedType.On("GetBytes", int32(1)).Return([]byte("dummy proof"), nil)
	managedType.On("GetBytes", int32(2)).Return([]byte("dummy vk"), nil)
	managedType.On("GetBytes", int32(3)).Return([]byte("dummy witness"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil).Times(3)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedVerifyGroth16(int32(ecc.BLS12_381), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedVerifyGroth16_Fail_InvalidCurve(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime
	managedType.On("GetBytes", int32(1)).Return([]byte("dummy proof"), nil)
	managedType.On("GetBytes", int32(2)).Return([]byte("dummy vk"), nil)
	managedType.On("GetBytes", int32(3)).Return([]byte("dummy witness"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil).Times(3)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedVerifyGroth16(int32(42), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedVerifyPlonk_Fail_GetProofBytes(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime
	managedType.On("GetBytes", int32(1)).Return(nil, assert.AnError)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedVerifyPlonk(int32(ecc.BLS12_381), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedVerifyPlonk_Fail_GetVKBytes(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime
	managedType.On("GetBytes", int32(1)).Return([]byte("dummy proof"), nil)
	managedType.On("GetBytes", int32(2)).Return(nil, assert.AnError)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedVerifyPlonk(int32(ecc.BLS12_381), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedVerifyPlonk_Fail_GetPubWitnessBytes(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime
	managedType.On("GetBytes", int32(1)).Return([]byte("dummy proof"), nil)
	managedType.On("GetBytes", int32(2)).Return([]byte("dummy vk"), nil)
	managedType.On("GetBytes", int32(3)).Return(nil, assert.AnError)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil).Times(2)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedVerifyPlonk(int32(ecc.BLS12_381), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedVerifyPlonk_Fail_Verification(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime
	managedType.On("GetBytes", int32(1)).Return([]byte("dummy proof"), nil)
	managedType.On("GetBytes", int32(2)).Return([]byte("dummy vk"), nil)
	managedType.On("GetBytes", int32(3)).Return([]byte("dummy witness"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil).Times(3)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedVerifyPlonk(int32(ecc.BLS12_381), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

func TestManagedVerifyPlonk_Fail_InvalidCurve(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime
	managedType.On("GetBytes", int32(1)).Return([]byte("dummy proof"), nil)
	managedType.On("GetBytes", int32(2)).Return([]byte("dummy vk"), nil)
	managedType.On("GetBytes", int32(3)).Return([]byte("dummy witness"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil).Times(3)
	runtime.On("IsUnsafeMode").Return(true)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := vmHooks.hooks.ManagedVerifyPlonk(int32(42), 1, 2, 3)
	assert.Equal(t, int32(-1), ret)
	runtime.AssertCalled(t, "FailExecution", mock.Anything)
}

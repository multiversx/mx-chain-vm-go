package vmhooks

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestVMHooksImpl_ManagedMapNew(t *testing.T) {
	t.Parallel()
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	hooks := vmHooks.hooks

	managedType.On("NewManagedMap").Return(int32(1))

	ret := hooks.ManagedMapNew()
	require.Equal(t, int32(1), ret)
}

func TestVMHooksImpl_ManagedMapPut(t *testing.T) {
	t.Parallel()
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	hooks := vmHooks.hooks

	managedType.On("ManagedMapPut", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ret := hooks.ManagedMapPut(0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_ManagedMapGet(t *testing.T) {
	t.Parallel()
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	hooks := vmHooks.hooks

	managedType.On("ManagedMapGet", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ret := hooks.ManagedMapGet(0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_ManagedMapRemove(t *testing.T) {
	t.Parallel()
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	hooks := vmHooks.hooks

	managedType.On("ManagedMapRemove", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ret := hooks.ManagedMapRemove(0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_ManagedMapContains(t *testing.T) {
	t.Parallel()
	vmHooks := createHooksWithBaseSetup()
	managedType := vmHooks.managedType
	hooks := vmHooks.hooks

	managedType.On("ManagedMapContains", mock.Anything, mock.Anything).Return(true, nil)

	ret := hooks.ManagedMapContains(0, 0)
	require.Equal(t, int32(1), ret)
}

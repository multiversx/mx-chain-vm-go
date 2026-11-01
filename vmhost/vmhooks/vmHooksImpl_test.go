package vmhooks

import (
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	mockcontext "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/require"
)

// Hooks built with no instance state must refuse to touch memory.
//
// The three cases this replaces - nil host, nil runtime, nil instance - were three ways of
// asking the same question through the runtime, which is how the hooks used to find the
// instance. The state now arrives directly, so there is one condition, and it is checked
// without consulting the host at all.
func TestVMHooksImpl_MemLoad_NoInstanceState(t *testing.T) {
	t.Parallel()

	t.Run("nil state, nil host", func(t *testing.T) {
		hooks := NewVMHooksImpl(nil, nil)
		data, err := hooks.MemLoad(0, 10)
		require.Nil(t, data)
		require.ErrorIs(t, err, vmhost.ErrNilInstance)
	})

	t.Run("nil state, live host", func(t *testing.T) {
		host := &mockcontext.VMHostMock{
			RuntimeContext: &mockcontext.RuntimeContextMock{},
		}
		hooks := NewVMHooksImpl(host, nil)
		data, err := hooks.MemLoad(0, 10)
		require.Nil(t, data)
		require.ErrorIs(t, err, vmhost.ErrNilInstance)
	})

	t.Run("typed nil state", func(t *testing.T) {
		var state *mockcontext.InstanceMock
		hooks := NewVMHooksImpl(nil, executor.InstanceState(state))
		data, err := hooks.MemLoad(0, 10)
		require.Nil(t, data)
		require.ErrorIs(t, err, vmhost.ErrNilInstance)
	})
}

func TestVMHooksImpl_MemStore_NoInstanceState(t *testing.T) {
	t.Parallel()

	t.Run("nil state, nil host", func(t *testing.T) {
		hooks := NewVMHooksImpl(nil, nil)
		err := hooks.MemStore(0, []byte{1, 2, 3})
		require.ErrorIs(t, err, vmhost.ErrNilInstance)
	})

	t.Run("nil state, live host", func(t *testing.T) {
		host := &mockcontext.VMHostMock{
			RuntimeContext: &mockcontext.RuntimeContextMock{},
		}
		hooks := NewVMHooksImpl(host, nil)
		err := hooks.MemStore(0, []byte{1, 2, 3})
		require.ErrorIs(t, err, vmhost.ErrNilInstance)
	})

	t.Run("typed nil state", func(t *testing.T) {
		var state *mockcontext.InstanceMock
		hooks := NewVMHooksImpl(nil, executor.InstanceState(state))
		err := hooks.MemStore(0, []byte{1, 2, 3})
		require.ErrorIs(t, err, vmhost.ErrNilInstance)
	})
}

// Hooks built with a live instance state reach its memory.
func TestVMHooksImpl_MemoryRoundTrip(t *testing.T) {
	t.Parallel()

	instance := mockcontext.NewInstanceMock(nil)
	hooks := NewVMHooksImpl(nil, instance)

	require.Nil(t, hooks.MemStore(0, []byte{1, 2, 3}))

	data, err := hooks.MemLoad(0, 3)
	require.Nil(t, err)
	require.Equal(t, []byte{1, 2, 3}, data)
}

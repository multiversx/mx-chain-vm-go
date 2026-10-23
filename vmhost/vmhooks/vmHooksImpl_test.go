package vmhooks

import (
	"testing"

	mockcontext "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/require"
)

func TestVMHooksImpl_MemLoad_NilInstance(t *testing.T) {
	t.Parallel()

	t.Run("nil host", func(t *testing.T) {
		hooks := NewVMHooksImpl(nil)
		data, err := hooks.MemLoad(0, 10)
		require.Nil(t, data)
		require.ErrorIs(t, err, vmhost.ErrNilInstance)
	})

	t.Run("nil runtime", func(t *testing.T) {
		host := &mockcontext.VMHostMock{
			RuntimeContext: nil,
		}
		hooks := NewVMHooksImpl(host)
		data, err := hooks.MemLoad(0, 10)
		require.Nil(t, data)
		require.ErrorIs(t, err, vmhost.ErrNilInstance)
	})

	t.Run("nil instance", func(t *testing.T) {
		host := &mockcontext.VMHostMock{
			RuntimeContext: &mockcontext.RuntimeContextMock{
				Instance: nil,
			},
		}
		hooks := NewVMHooksImpl(host)
		data, err := hooks.MemLoad(0, 10)
		require.Nil(t, data)
		require.ErrorIs(t, err, vmhost.ErrNilInstance)
	})
}

func TestVMHooksImpl_MemStore_NilInstance(t *testing.T) {
	t.Parallel()

	t.Run("nil host", func(t *testing.T) {
		hooks := NewVMHooksImpl(nil)
		err := hooks.MemStore(0, []byte{1, 2, 3})
		require.ErrorIs(t, err, vmhost.ErrNilInstance)
	})

	t.Run("nil runtime", func(t *testing.T) {
		host := &mockcontext.VMHostMock{
			RuntimeContext: nil,
		}
		hooks := NewVMHooksImpl(host)
		err := hooks.MemStore(0, []byte{1, 2, 3})
		require.ErrorIs(t, err, vmhost.ErrNilInstance)
	})

	t.Run("nil instance", func(t *testing.T) {
		host := &mockcontext.VMHostMock{
			RuntimeContext: &mockcontext.RuntimeContextMock{
				Instance: nil,
			},
		}
		hooks := NewVMHooksImpl(host)
		err := hooks.MemStore(0, []byte{1, 2, 3})
		require.ErrorIs(t, err, vmhost.ErrNilInstance)
	})
}

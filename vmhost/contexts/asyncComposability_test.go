package contexts

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/stretchr/testify/mock"











	"github.com/multiversx/mx-chain-vm-go/vmhost"

	"github.com/stretchr/testify/require"
)

func TestAsyncContext_NotifyChildIsComplete(t *testing.T) {
	t.Parallel()

	t.Run("child completes, no other pending calls", func(t *testing.T) {
		t.Parallel()

		host := &context.MockVMHost{}
		storage := &context.MockStorageContext{}
		host.On("Storage").Return(storage)
		storage.On("SetProtectedStorageToAddressUnmetered", mock.Anything, mock.Anything).Return(nil)).Return(nil)).Return(nil)).Return(nil)).Return(nil)).Return(nil), mock.Anything).Return(vmhost.StorageDeleted, nil)

		ac := &asyncContext{
			callsCounter: 1,
			host:         host,
		}
		ac.callID = []byte("test")

		err := ac.NotifyChildIsComplete([]byte("child1"), 100)
		require.Nil(t, err)
	})

	t.Run("child completes, other pending calls", func(t *testing.T) {
		t.Parallel()
		host := &context.MockVMHost{}
		storage := &context.MockStorageContext{}
		host.On("Storage").Return(storage)
		storage.On("SetProtectedStorageToAddressUnmetered", mock.Anything, mock.Anything).Return(nil)).Return(nil)).Return(nil)).Return(nil)).Return(nil)).Return(nil), mock.Anything).Return(vmhost.StorageModified, nil)

		ac := &asyncContext{
			callsCounter: 2,
			host:         host,
		}

		err := ac.NotifyChildIsComplete([]byte("child1"), 100)
		require.Nil(t, err)
		require.Equal(t, uint64(100), ac.gasAccumulated)
		require.Equal(t, uint64(1), ac.callsCounter)
	})

	t.Run("delete async call fails", func(t *testing.T) {
		t.Parallel()

		ac := &asyncContext{
			callsCounter: 1,
		}
		ac.asyncCallGroups = []*vmhost.AsyncCallGroup{
			{
				AsyncCalls: []*vmhost.AsyncCall{
					{CallID: []byte("child2")},
				},
			},
		}
		err := ac.NotifyChildIsComplete([]byte("child1"), 100)
		require.NotNil(t, err)
		require.Equal(t, vmhost.ErrAsyncCallNotFound, err)
	})
}

func TestAsyncContext_complete(t *testing.T) {
	t.Parallel()

	t.Run("first call", func(t *testing.T) {
		t.Parallel()

		host := &context.MockVMHost{}
		storage := &context.MockStorageContext{}
		host.On("Storage").Return(storage)
		storage.On("SetProtectedStorageToAddressUnmetered", mock.Anything, mock.Anything).Return(nil)).Return(nil)).Return(nil)).Return(nil)).Return(nil)).Return(nil), mock.Anything).Return(vmhost.StorageDeleted, nil)
		ac := &asyncContext{
			host:       host,
			parentAddr: nil,
		}

		err := ac.complete()
		require.Nil(t, err)
	})
}

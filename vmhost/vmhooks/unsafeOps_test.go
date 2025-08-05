package vmhooks

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/mock"
	"github.com/stretchr/testify/assert"
)

func TestUnsafeOps_ActivateUnsafeMode(t *testing.T) {
	t.Parallel()

	t.Run("should set unsafe mode to true", func(t *testing.T) {
		t.Parallel()

		host := &mock.UnsafeVMHostMock{}
		hooks := NewVMHooksImpl(host)

		hooks.ActivateUnsafeMode()

		assert.True(t, host.IsUnsafeMode())
	})

	t.Run("should fail when out of gas", func(t *testing.T) {
		t.Parallel()

		host := &mock.UnsafeVMHostMock{}
		host.Metering().(*mock.MeteringContextMock).UseGasBoundedAndAddTracedGasCalled = func(name string, gas uint64) error {
			return errors.New("out of gas")
		}
		hooks := NewVMHooksImpl(host)

		hooks.ActivateUnsafeMode()

		assert.True(t, host.FailExecutionCalled)
	})
}

func TestUnsafeOps_DeactivateUnsafeMode(t *testing.T) {
	t.Parallel()

	t.Run("should set unsafe mode to false", func(t *testing.T) {
		t.Parallel()

		host := &mock.UnsafeVMHostMock{}
		host.SetUnsafeMode(true)
		hooks := NewVMHooksImpl(host)

		hooks.DeactivateUnsafeMode()

		assert.False(t, host.IsUnsafeMode())
	})

	t.Run("should fail when out of gas", func(t *testing.T) {
		t.Parallel()

		host := &mock.UnsafeVMHostMock{}
		host.Metering().(*mock.MeteringContextMock).UseGasBoundedAndAddTracedGasCalled = func(name string, gas uint64) error {
			return errors.New("out of gas")
		}
		hooks := NewVMHooksImpl(host)

		hooks.DeactivateUnsafeMode()

		assert.True(t, host.FailExecutionCalled)
	})
}

package hostCore

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/stretchr/testify/assert"
)

func TestUnsafeMode(t *testing.T) {
	t.Run("safe mode", func(t *testing.T) {
		host := &vmHost{}
		runtimeContext := &context.RuntimeContextMock{}
		host.runtimeContext = runtimeContext

		host.SetUnsafeMode(false)

		host.FailExecutionConditionally(errors.New("test error"))

		assert.True(t, runtimeContext.FailExecutionCalled)
	})

	t.Run("unsafe mode", func(t *testing.T) {
		host := &vmHost{}
		runtimeContext := &context.RuntimeContextMock{}
		host.runtimeContext = runtimeContext

		host.SetUnsafeMode(true)

		host.FailExecutionConditionally(errors.New("test error"))

		assert.False(t, runtimeContext.FailExecutionCalled)
	})
}

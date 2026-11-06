package executorwrapper

import (
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/stretchr/testify/require"
)

func TestWrapperExecutorFactory_CreateExecutorNew_MissingDependencies(t *testing.T) {
	t.Run("nil wrapped factory", func(t *testing.T) {
		factory := SimpleWrappedExecutorFactory(nil)

		exec, err := factory.CreateExecutorNew(executor.ExecutorFactoryArgsNew{})
		require.ErrorIs(t, err, ErrNilWrappedExecutorFactory)
		require.Nil(t, exec)
	})

	t.Run("nil logger", func(t *testing.T) {
		factory := NewWrappedExecutorFactory(nil, SimpleWrappedExecutorFactory(nil))

		exec, err := factory.CreateExecutorNew(executor.ExecutorFactoryArgsNew{})
		require.ErrorIs(t, err, ErrNilExecutorLogger)
		require.Nil(t, exec)
	})
}

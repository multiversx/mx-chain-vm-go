package factory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func createComponents() (RunTypeComponentsHandler, error) {
	rcf := NewRunTypeComponentsFactory()
	return NewManagedRunTypeComponents(rcf)
}

func TestNewManagedRunTypeComponents(t *testing.T) {
	t.Parallel()

	t.Run("should error", func(t *testing.T) {
		managedRunTypeComponents, err := NewManagedRunTypeComponents(nil)
		require.ErrorIs(t, err, errNilRunTypeComponentsFactory)
		require.True(t, managedRunTypeComponents.IsInterfaceNil())
	})
	t.Run("should work", func(t *testing.T) {
		rcf := NewRunTypeComponentsFactory()
		managedRunTypeComponents, err := NewManagedRunTypeComponents(rcf)
		require.NoError(t, err)
		require.False(t, managedRunTypeComponents.IsInterfaceNil())
	})
}

func TestNewManagedRunTypeComponents_Create(t *testing.T) {
	t.Parallel()

	managedRunTypeComponents, err := createComponents()
	require.NoError(t, err)

	require.Nil(t, managedRunTypeComponents.BlockchainContextCreator())
	require.Nil(t, managedRunTypeComponents.ExecutorCreator())
	require.Nil(t, managedRunTypeComponents.RuntimeContextCreator())
	require.Nil(t, managedRunTypeComponents.MeteringContextCreator())
	require.Nil(t, managedRunTypeComponents.OutputContextCreator())
	require.Nil(t, managedRunTypeComponents.ExecuteOnSameContextHandler())
	require.Nil(t, managedRunTypeComponents.CreateNewContractHandler())
	require.Nil(t, managedRunTypeComponents.GasScheduleFactory())

	err = managedRunTypeComponents.Create()
	require.NoError(t, err)

	require.NotNil(t, managedRunTypeComponents.BlockchainContextCreator())
	require.NotNil(t, managedRunTypeComponents.ExecutorCreator())
	require.NotNil(t, managedRunTypeComponents.RuntimeContextCreator())
	require.NotNil(t, managedRunTypeComponents.MeteringContextCreator())
	require.NotNil(t, managedRunTypeComponents.OutputContextCreator())
	require.NotNil(t, managedRunTypeComponents.ExecuteOnSameContextHandler())
	require.NotNil(t, managedRunTypeComponents.CreateNewContractHandler())
	require.NotNil(t, managedRunTypeComponents.GasScheduleFactory())

	require.Equal(t, runTypeComponentsName, managedRunTypeComponents.String())
	require.NoError(t, managedRunTypeComponents.Close())
}

func TestNewManagedRunTypeComponents_Close(t *testing.T) {
	t.Parallel()

	managedRunTypeComponents, _ := createComponents()
	require.NoError(t, managedRunTypeComponents.Close())

	err := managedRunTypeComponents.Create()
	require.NoError(t, err)

	require.NoError(t, managedRunTypeComponents.Close())
	require.Nil(t, managedRunTypeComponents.BlockchainContextCreator())
	require.Nil(t, managedRunTypeComponents.ExecutorCreator())
	require.Nil(t, managedRunTypeComponents.RuntimeContextCreator())
	require.Nil(t, managedRunTypeComponents.MeteringContextCreator())
	require.Nil(t, managedRunTypeComponents.OutputContextCreator())
	require.Nil(t, managedRunTypeComponents.ExecuteOnSameContextHandler())
	require.Nil(t, managedRunTypeComponents.CreateNewContractHandler())
	require.Nil(t, managedRunTypeComponents.GasScheduleFactory())
}

func TestManagedRunTypeCoreComponents_CheckSubcomponents(t *testing.T) {
	t.Parallel()

	managedRunTypeComponents, _ := createComponents()
	err := managedRunTypeComponents.CheckSubcomponents()
	require.Equal(t, ErrNilRunTypeComponents, err)

	err = managedRunTypeComponents.Create()
	require.NoError(t, err)

	err = managedRunTypeComponents.CheckSubcomponents()
	require.NoError(t, err)

	require.NoError(t, managedRunTypeComponents.Close())
}

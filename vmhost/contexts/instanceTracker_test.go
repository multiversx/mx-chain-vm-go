package contexts

import (
	"testing"

	"github.com/multiversx/mx-chain-vm-v1_4-go/wasmer"
	"github.com/stretchr/testify/require"
)

// TODO test tracking precompiled & bytecode instances
func TestInstanceTracker_TrackInstance(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	newInstance := &wasmer.Instance{
		AlreadyClean: false,
	}

	iTracker.SetNewInstance(newInstance, Warm)
	iTracker.codeHash = []byte("testinst")

	require.Equal(t, newInstance, iTracker.instance)
	require.Equal(t, Warm, iTracker.cacheLevel)

	iTracker.SaveAsWarmInstance()

	warm, cold := iTracker.NumRunningInstances()
	require.Equal(t, 1, warm)
	require.Equal(t, 0, cold)
}

func TestInstanceTracker_UnsetInstance_AlreadyNil_Ok(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	iTracker.instance = nil
	iTracker.UnsetInstance()
	require.Nil(t, iTracker.instance)
}

func TestInstanceTracker_UnsetInstance_Ok(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	iTracker.instance = &wasmer.Instance{
		AlreadyClean: true,
	}
	iTracker.UnsetInstance()
	require.Nil(t, iTracker.instance)
}

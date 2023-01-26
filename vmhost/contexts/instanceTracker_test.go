package contexts

import (
	"strings"
	"testing"

	mock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/context"
	"github.com/multiversx/mx-chain-vm-v1_4-go/wasmer"
	"github.com/stretchr/testify/require"
)

func TestInstanceTracker_TrackInstance(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	newInstance := &wasmer.Instance{
		AlreadyClean: false,
	}

	iTracker.SetNewInstance(newInstance, Bytecode)
	iTracker.codeHash = []byte("testinst")

	require.False(t, iTracker.IsCodeHashOnTheStack(iTracker.codeHash))

	require.Equal(t, newInstance, iTracker.instance)
	require.Equal(t, Bytecode, iTracker.cacheLevel)

	iTracker.SaveAsWarmInstance()

	warm, cold := iTracker.NumRunningInstances()
	require.Equal(t, 1, warm)
	require.Equal(t, 0, cold)
}

func TestInstanceTracker_InitState(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	iTracker.SetNewInstance(mock.NewInstanceMock(nil), Bytecode)
	iTracker.SetNewInstance(mock.NewInstanceMock(nil), Bytecode)
	iTracker.SetNewInstance(mock.NewInstanceMock(nil), Bytecode)
	iTracker.SetNewInstance(mock.NewInstanceMock(nil), Bytecode)
	iTracker.SetNewInstance(mock.NewInstanceMock(nil), Bytecode)

	require.Equal(t, 5, iTracker.numRunningInstances)
	require.Len(t, iTracker.instances, 5)

	iTracker.InitState()

	require.Nil(t, iTracker.instance)
	require.Equal(t, 0, iTracker.numRunningInstances)
	require.Len(t, iTracker.codeHash, 0)
	require.Len(t, iTracker.instances, 0)
}

func TestInstanceTracker_GetWarmInstance(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	testData := []string{"warm1", "bytecode1", "bytecode2", "warm2"}

	for _, codeHash := range testData {
		iTracker.SetNewInstance(mock.NewInstanceMock([]byte(codeHash)), Bytecode)
		iTracker.codeHash = []byte(codeHash)
		if strings.Contains(codeHash, "warm") {
			iTracker.SaveAsWarmInstance()
		}
	}

	require.Equal(t, 4, iTracker.numRunningInstances)
	require.Len(t, iTracker.instances, 4)

	for _, codeHash := range testData {
		instance, ok := iTracker.GetWarmInstance([]byte(codeHash))

		if strings.Contains(codeHash, "warm") {
			require.NotNil(t, instance)
			require.True(t, ok)
			continue
		}

		require.Nil(t, instance)
		require.False(t, ok)
	}

}

func TestInstanceTracker_UserWarmInstance(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	testData := []string{"warm1", "bytecode1", "warm2", "bytecode2"}

	for _, codeHash := range testData {
		iTracker.SetNewInstance(mock.NewInstanceMock([]byte(codeHash)), Bytecode)
		iTracker.codeHash = []byte(codeHash)

		if strings.Contains(codeHash, "warm") {
			iTracker.SaveAsWarmInstance()
		}
	}

	require.Equal(t, []byte("bytecode2"), iTracker.CodeHash())

	for _, codeHash := range testData {
		ok := iTracker.UseWarmInstance([]byte(codeHash), false)

		if strings.Contains(codeHash, "warm") {
			require.True(t, ok)
			continue
		}

		require.False(t, ok)
	}
}

// a->b->a(cold)->b(cold)
// a->a(cold)->a(cold)
// a->b->c->b(cold)->c(cold)->d->a(cold)

func TestInstanceTracker_PopSetActiveWarmChain(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

		testData := []string{"first", "second", "first", "second"}

	for _, codeHash := range testData {
		iTracker.SetNewInstance(mock.NewInstanceMock([]byte(codeHash)), Bytecode)
		iTracker.codeHash = []byte(codeHash)
		iTracker.SaveAsWarmInstance()
		iTracker.PushState()
	}

	for _, codeHash := range testData {
		iTracker.SetNewInstance(mock.NewInstanceMock([]byte(codeHash)), Bytecode)
		iTracker.codeHash = []byte(codeHash)
		iTracker.SaveAsWarmInstance()
		iTracker.PushState()
	}

	require.Len(t, iTracker.codeHashStack, 4)
	require.Len(t, iTracker.instanceStack, 4)

	iTracker.SetNewInstance(mock.NewInstanceMock([]byte("active")), Bytecode)
	iTracker.codeHash = []byte("active")
	iTracker.SaveAsWarmInstance()
	require.Equal(t, []byte("active"), iTracker.codeHash)

	iTracker.PopSetActiveState()
	require.Equal(t, []byte("last"), iTracker.codeHash)

	iTracker.PopSetActiveState()
	iTracker.PopSetActiveState()
	iTracker.PopSetActiveState()
	require.Equal(t, []byte("first"), iTracker.codeHash)

	require.Equal(t, 5, iTracker.numRunningInstances)

	iTracker.ClearWarmInstanceCache()

	require.Len(t, iTracker.instanceStack, 0)
	require.Len(t, iTracker.codeHashStack, 0)
	require.Equal(t, 0, iTracker.numRunningInstances)
	require.Nil(t, iTracker.CheckInstances())
}

func TestInstanceTracker_ForceCleanInstanceWithBypass(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	testData := []string{"warm1", "bytecode1"}

	for _, codeHash := range testData {
		iTracker.SetNewInstance(mock.NewInstanceMock([]byte(codeHash)), Bytecode)
		iTracker.codeHash = []byte(codeHash)

		if strings.Contains(codeHash, "warm") {
			iTracker.SaveAsWarmInstance()
		}
	}

	warm, cold := iTracker.NumRunningInstances()
	require.Equal(t, 1, warm)
	require.Equal(t, 1, cold)

	iTracker.ForceCleanInstance(true)
	require.Nil(t, iTracker.instance)

	iTracker.UseWarmInstance([]byte("warm1"), false)
	require.NotNil(t, iTracker.instance)

	iTracker.ForceCleanInstance(true)
	require.Nil(t, iTracker.instance)

	require.Equal(t, 0, iTracker.numRunningInstances)
	require.Nil(t, iTracker.CheckInstances())
}

func TestInstanceTracker_DoubleForceClean(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	iTracker.SetNewInstance(mock.NewInstanceMock(nil), Bytecode)
	require.NotNil(t, iTracker.instance)
	require.Equal(t, 1, iTracker.numRunningInstances)

	iTracker.ForceCleanInstance(true)
	require.Equal(t, 0, iTracker.numRunningInstances)
	require.Nil(t, iTracker.CheckInstances())

	iTracker.ForceCleanInstance(true)
	require.Equal(t, 0, iTracker.numRunningInstances)
	require.Nil(t, iTracker.CheckInstances())
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

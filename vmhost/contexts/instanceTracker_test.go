package contexts

import (
	"fmt"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	mock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
	"github.com/stretchr/testify/require"
)

func TestInstanceTracker_TrackInstance(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	newInstance := &wasmer2.Wasmer2Instance{
		AlreadyClean: false}

	_ = iTracker.SetNewInstance(newInstance, Bytecode)
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
	require.Equal(t, 0, iTracker.numRunningInstances)

	for i := 0; i < 5; i++ {
		_ = iTracker.SetNewInstance(mock.NewInstanceMock(nil), Bytecode)
	}

	require.Equal(t, 5, iTracker.numRunningInstances)
	require.Len(t, iTracker.instances, 5)

	iTracker.codeSize = 12
	iTracker.InitState()

	require.Nil(t, iTracker.instance)
	require.Len(t, iTracker.codeHash, 0)
	require.Len(t, iTracker.instances, 0)
	require.Zero(t, iTracker.codeSize)

	// InitState() must not reset numRunningInstances
	require.Equal(t, 5, iTracker.numRunningInstances)
}

func TestInstanceTracker_GetWarmInstance(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	testData := []string{"warm1", "bytecode1", "bytecode2", "warm2"}

	for _, codeHash := range testData {
		_ = iTracker.SetNewInstance(mock.NewInstanceMock([]byte(codeHash)), Bytecode)
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

func TestInstanceTracker_UseWarmInstance(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	testData := []string{"warm1", "bytecode1", "warm2", "bytecode2"}

	for _, codeHash := range testData {
		_ = iTracker.SetNewInstance(mock.NewInstanceMock([]byte(codeHash)), Bytecode)
		iTracker.codeHash = []byte(codeHash)

		if strings.Contains(codeHash, "warm") {
			iTracker.SaveAsWarmInstance()
		}
	}

	require.Equal(t, []byte("bytecode2"), iTracker.CodeHash())

	for _, codeHash := range testData {
		ok, _ := iTracker.UseWarmInstance([]byte(codeHash), false)

		if strings.Contains(codeHash, "warm") {
			require.True(t, ok)
			continue
		}

		require.False(t, ok)
	}
}

func TestInstanceTracker_IsCodeHashOnStack_Ok(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	testData := []string{"alpha", "beta", "alpha", "active"}

	for i, codeHash := range testData {
		_ = iTracker.SetNewInstance(mock.NewInstanceMock([]byte(codeHash)), Bytecode)
		iTracker.codeHash = []byte(codeHash)
		if i < 2 || codeHash == "active" {
			iTracker.SaveAsWarmInstance()
		}
		if codeHash != "active" {
			iTracker.PushState()
		}
	}
	require.Len(t, iTracker.codeHashStack, 3)
	require.Len(t, iTracker.instanceStack, 3)

	warm, cold := iTracker.NumRunningInstances()
	require.Equal(t, 3, warm)
	require.Equal(t, 1, cold)

	iTracker.PopSetActiveState()
	require.Equal(t, []byte("alpha"), iTracker.CodeHash())
	require.True(t, iTracker.IsCodeHashOnTheStack(iTracker.codeHash))

	iTracker.PopSetActiveState()
	require.Equal(t, []byte("beta"), iTracker.CodeHash())
	require.False(t, iTracker.IsCodeHashOnTheStack(iTracker.codeHash))
}

// stack: alpha<-alpha(cold)<-alpha(cold)<-alpha(cold)
func TestInstanceTracker_PopSetActiveSelfScenario(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	testData := []string{"alpha", "alpha", "alpha", "alpha", "active"}

	for i, codeHash := range testData {
		_ = iTracker.SetNewInstance(mock.NewInstanceMock([]byte(codeHash)), Bytecode)
		iTracker.codeHash = []byte(codeHash)
		if i == 0 || codeHash == "active" {
			iTracker.SaveAsWarmInstance()
		}
		if codeHash != "active" {
			iTracker.PushState()
		}
	}
	require.Len(t, iTracker.codeHashStack, 4)
	require.Len(t, iTracker.instanceStack, 4)

	warm, cold := iTracker.NumRunningInstances()
	require.Equal(t, 2, warm)
	require.Equal(t, 3, cold)

	checkColdInstancesAfterEmptyingStack(t, iTracker)

	iTracker.ClearWarmInstanceCache()
	checkInstances(t, iTracker)
}

// stack: alpha<-beta<-alpha(cold)<-beta(cold)
func TestInstanceTracker_PopSetActiveSimpleScenario(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	testData := []string{"alpha", "beta", "alpha", "beta", "active"}

	for i, codeHash := range testData {
		_ = iTracker.SetNewInstance(mock.NewInstanceMock([]byte(codeHash)), Bytecode)
		iTracker.codeHash = []byte(codeHash)
		if i < 2 || codeHash == "active" {
			iTracker.SaveAsWarmInstance()
		}
		if codeHash != "active" {
			iTracker.PushState()
		}
	}
	require.Len(t, iTracker.codeHashStack, 4)
	require.Len(t, iTracker.instanceStack, 4)

	warm, cold := iTracker.NumRunningInstances()
	require.Equal(t, 3, warm)
	require.Equal(t, 2, cold)

	emptyInstanceStack(iTracker)

	warm, cold = iTracker.NumRunningInstances()
	require.Equal(t, 3, warm)
	require.Equal(t, 0, cold)

	require.Equal(t, 3, iTracker.numRunningInstances)
	iTracker.InitState()
	require.Equal(t, 3, iTracker.numRunningInstances)

	iTracker.ClearWarmInstanceCache()
	require.Equal(t, 0, iTracker.numRunningInstances)
	checkInstances(t, iTracker)
}

// stack: alpha<-beta<-gamma<-beta(cold)<-gamma(cold)<-delta<-alpha(cold)
func TestInstanceTracker_PopSetActiveComplexScenario(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	testData := []string{"alpha", "beta", "gamma", "beta", "gamma", "delta", "alpha", "active"}

	for i, codeHash := range testData {
		_ = iTracker.SetNewInstance(mock.NewInstanceMock([]byte(codeHash)), Bytecode)
		iTracker.codeHash = []byte(codeHash)
		if i < 3 || codeHash == "delta" || codeHash == "active" {
			iTracker.SaveAsWarmInstance()
		}
		if codeHash != "active" {
			iTracker.PushState()
		}
	}
	require.Len(t, iTracker.codeHashStack, 7)
	require.Len(t, iTracker.instanceStack, 7)

	warm, cold := iTracker.NumRunningInstances()
	require.Equal(t, 5, warm)
	require.Equal(t, 3, cold)

	checkColdInstancesAfterEmptyingStack(t, iTracker)

	iTracker.ClearWarmInstanceCache()
	checkInstances(t, iTracker)
}

func TestInstanceTracker_PopSetActiveWarmOnlyScenario(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	testData := []string{"alpha", "beta", "gamma", "delta", "active"}

	for _, codeHash := range testData {
		_ = iTracker.SetNewInstance(mock.NewInstanceMock([]byte(codeHash)), Bytecode)
		iTracker.codeHash = []byte(codeHash)
		iTracker.SaveAsWarmInstance()

		if codeHash != "active" {
			iTracker.PushState()
		}
	}
	require.Len(t, iTracker.codeHashStack, 4)
	require.Len(t, iTracker.instanceStack, 4)

	warm, cold := iTracker.NumRunningInstances()
	require.Equal(t, 5, warm)
	require.Equal(t, 0, cold)

	checkColdInstancesAfterEmptyingStack(t, iTracker)

	iTracker.ClearWarmInstanceCache()
	checkInstances(t, iTracker)
}

func TestInstanceTracker_ForceCleanInstanceWithBypass(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	testData := []string{"warm1", "bytecode1"}

	for _, codeHash := range testData {
		_ = iTracker.SetNewInstance(mock.NewInstanceMock([]byte(codeHash)), Bytecode)
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

	_, _ = iTracker.UseWarmInstance([]byte("warm1"), false)
	require.NotNil(t, iTracker.instance)

	iTracker.ForceCleanInstance(true)
	require.Nil(t, iTracker.instance)

	require.Equal(t, 0, iTracker.numRunningInstances)
	require.Nil(t, iTracker.CheckInstances())
}

func TestInstanceTracker_DoubleForceClean(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	_ = iTracker.SetNewInstance(mock.NewInstanceMock(nil), Bytecode)
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

	iTracker.instance = &wasmer2.Wasmer2Instance{
		AlreadyClean: true,
	}
	iTracker.UnsetInstance()
	require.Nil(t, iTracker.instance)
}

func checkColdInstancesAfterEmptyingStack(t *testing.T, iTracker *instanceTracker) {
	emptyInstanceStack(iTracker)
	_, cold := iTracker.NumRunningInstances()
	require.Equal(t, 0, cold)
}

func emptyInstanceStack(iTracker *instanceTracker) {
	n := len(iTracker.instanceStack)
	for i := 0; i < n; i++ {
		iTracker.PopSetActiveState()
	}
}

func checkInstances(t *testing.T, iTracker *instanceTracker) {
	require.Equal(t, 0, iTracker.numRunningInstances)
	require.Len(t, iTracker.instanceStack, 0)
	require.Len(t, iTracker.codeHashStack, 0)
	require.Nil(t, iTracker.CheckInstances())
}

func TestInstanceTracker_SetNewInstance_TooManyInstancesCleansColdInstanceAndLeavesStateIntact(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	var lastAcceptedInstance *mock.InstanceMock
	for i := 0; i < warmCacheSize-2; i++ {
		inst := mock.NewInstanceMock([]byte{byte(i)})
		err := iTracker.SetNewInstance(inst, Bytecode)
		require.Nil(t, err)
		lastAcceptedInstance = inst
	}

	require.Equal(t, warmCacheSize-2, len(iTracker.instances))
	require.Equal(t, warmCacheSize-2, iTracker.numRunningInstances)
	require.Equal(t, lastAcceptedInstance, iTracker.instance)

	// 99th cold instance must be rejected, cleaned, and must not mutate tracker state
	rejectedColdInstance := mock.NewInstanceMock([]byte("rejected"))
	err = iTracker.SetNewInstance(rejectedColdInstance, Bytecode)
	require.Equal(t, errTooManyInstances, err)
	require.True(t, rejectedColdInstance.IsAlreadyCleaned(), "rejected cold instance must be cleaned")
	require.Equal(t, warmCacheSize-2, len(iTracker.instances), "instances map must not be modified on error")
	require.Equal(t, warmCacheSize-2, iTracker.numRunningInstances, "numRunningInstances must not be incremented on error")
	require.Equal(t, lastAcceptedInstance, iTracker.instance, "active instance must remain unchanged")

	// Nil instance check
	require.Equal(t, errNilInstance, iTracker.SetNewInstance(nil, Bytecode))
}

func TestInstanceTracker_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	inst1 := mock.NewInstanceMock([]byte("inst1"))
	err = iTracker.SetNewInstance(inst1, Bytecode)
	require.Nil(t, err)

	done := make(chan struct{})
	go func() {
		for i := 0; i < 500; i++ {
			_ = iTracker.Instance()
		}
		close(done)
	}()

	for i := 0; i < 500; i++ {
		inst := mock.NewInstanceMock([]byte(fmt.Sprintf("loop%d", i)))
		iTracker.ReplaceInstance(inst)
		iTracker.UnsetInstance()
	}
	<-done
}

func TestInstanceTracker_AbandonInstance_Ok(t *testing.T) {
	iTracker, err := NewInstanceTracker()
	require.Nil(t, err)

	inst1 := mock.NewInstanceMock([]byte("inst1"))
	err = iTracker.SetNewInstance(inst1, Bytecode)
	require.Nil(t, err)
	iTracker.SetCodeHash([]byte("code1"))
	iTracker.SaveAsWarmInstance()

	warm, cold := iTracker.NumRunningInstances()
	require.Equal(t, 1, warm)
	require.Equal(t, 0, cold)

	iTracker.AbandonInstance()

	require.Nil(t, iTracker.Instance())
	require.True(t, inst1.IsAlreadyCleaned())
	warm, cold = iTracker.NumRunningInstances()
	require.Equal(t, 0, warm)
	require.Equal(t, 0, cold)
}

type mockEpochsHandler struct {
	mock.DefaultEnableEpochsHandlerStub
	isFlagEnabled bool
}

func (m *mockEpochsHandler) IsFlagEnabled(_ core.EnableEpochFlag) bool {
	return m.isFlagEnabled
}

type mockTrackedInstance struct {
	*mock.InstanceMock
	id string
}

func (m *mockTrackedInstance) ID() string {
	return m.id
}

func TestInstanceTracker_AttributeExtraGasUsageFlag_Gating(t *testing.T) {
	t.Parallel()

	// 1. When AttributeExtraGasUsageFlag is disabled: ID prefix (%p) reuses slot,
	// allowing 150 deployments with the same pointer without hitting errTooManyInstances.
	iTrackerDisabled, err := NewInstanceTracker()
	require.Nil(t, err)
	iTrackerDisabled.SetEnableEpochsHandler(&mockEpochsHandler{isFlagEnabled: false})

	baseMock := mock.NewInstanceMock(nil)
	for i := 0; i < 150; i++ {
		// Mock ID format "%p-%d" like Wasmer2Instance
		inst := &mockTrackedInstance{
			InstanceMock: baseMock,
			id:           fmt.Sprintf("0x7fa000000000-%d", i),
		}
		err = iTrackerDisabled.SetNewInstance(inst, Bytecode)
		require.Nil(t, err)
	}
	require.Equal(t, 1, len(iTrackerDisabled.instances), "Pointer prefix must be reused when flag is disabled")

	// 2. When AttributeExtraGasUsageFlag is enabled: strictly bounds total instances.
	iTrackerEnabled, err := NewInstanceTracker()
	require.Nil(t, err)
	iTrackerEnabled.SetEnableEpochsHandler(&mockEpochsHandler{isFlagEnabled: true})

	for i := 0; i < warmCacheSize-2; i++ {
		inst := &mockTrackedInstance{
			InstanceMock: baseMock,
			id:           fmt.Sprintf("0x7fa000000000-%d", i),
		}
		err = iTrackerEnabled.SetNewInstance(inst, Bytecode)
		require.Nil(t, err)
	}
	require.Equal(t, warmCacheSize-2, len(iTrackerEnabled.instances))

	// Next instance exceeds warmCacheSize-1
	rejectedInst := &mockTrackedInstance{
		InstanceMock: baseMock,
		id:           "0x7fa000000000-rejected",
	}
	err = iTrackerEnabled.SetNewInstance(rejectedInst, Bytecode)
	require.Equal(t, errTooManyInstances, err)
}


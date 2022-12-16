package contexts

import (
	"bytes"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-storage/lrucache"
	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen"
	"github.com/ElrondNetwork/wasm-vm-v1_4/wasmer"
)

type instanceCacheLevel int

const (
	// Warm indicates that the instance to track is a warm instance
	Warm instanceCacheLevel = iota

	// Precompiled indicates that the instance to track is cold and has been created from precompiled code
	Precompiled

	// Bytecode indicates that the instance to track is cold and has been created from raw bytecode
	Bytecode
)

var _ arwen.StateStack = (*instanceTracker)(nil)

type instanceTracker struct {
	codeHash            []byte
	numRunningInstances int
	warmInstanceCache   Cacher
	instance            wasmer.InstanceHandler
	cacheLevel          instanceCacheLevel
	instanceStack       []wasmer.InstanceHandler
	codeHashStack       [][]byte
}

// NewInstanceTracker creates a new instanceTracker instance
func NewInstanceTracker() (*instanceTracker, error) {
	tracker := &instanceTracker{
		instanceStack:       make([]wasmer.InstanceHandler, 0),
		numRunningInstances: 0,
	}

	var err error
	instanceEvictedCallback := tracker.makeInstanceEvictionCallback()
	if WarmInstancesEnabled {
		tracker.warmInstanceCache, err = lrucache.NewCacheWithEviction(warmCacheSize, instanceEvictedCallback)
	} else {
		tracker.warmInstanceCache = nil
	}
	if err != nil {
		return nil, err
	}

	return tracker, nil
}

func (tracker *instanceTracker) InitState() {
	tracker.codeHash = make([]byte, 0)
}

func (tracker *instanceTracker) PushState() {
	tracker.instanceStack = append(tracker.instanceStack, tracker.instance)
	tracker.codeHashStack = append(tracker.codeHashStack, tracker.codeHash)
	logRuntime.Trace("pushing instance", "id", tracker.instance.Id(), "codeHash", tracker.codeHash)
}

func (tracker *instanceTracker) PopSetActiveState() {
	instanceStackLen := len(tracker.instanceStack)
	if instanceStackLen == 0 {
		return
	}

	prevInstance := tracker.instanceStack[instanceStackLen-1]
	tracker.instanceStack = tracker.instanceStack[:instanceStackLen-1]

	if prevInstance == tracker.instance {
		// The current Wasmer instance was previously pushed on the instance stack,
		// but a new Wasmer instance has not been created in the meantime. This
		// means that the instance at the top of the stack is the same as the
		// current instance, so it cannot be cleaned, because the execution will
		// resume on it. Popping will therefore only remove the top of the stack,
		// without cleaning anything.
		return
	}

	if !check.IfNil(tracker.instance) {
		onStack := tracker.IsCodeHashOnTheStack(tracker.codeHash)
		coldOnlyEnabled := !WarmInstancesEnabled
		if onStack || coldOnlyEnabled {
			if tracker.instance.Clean() {
				tracker.updateNumRunningInstances(-1)
			}
		}

		logRuntime.Trace("pop instance", "id", tracker.instance.Id(), "codeHash", tracker.codeHash)
	}

	tracker.ReplaceInstance(prevInstance)

	prevCodeHash := tracker.codeHashStack[instanceStackLen-1]
	tracker.codeHashStack = tracker.codeHashStack[:instanceStackLen-1]
	tracker.codeHash = prevCodeHash
}

func (tracker *instanceTracker) PopDiscard() {
}

// ClearStateStack reinitializes the state stack.
func (tracker *instanceTracker) ClearStateStack() {
	tracker.codeHashStack = make([][]byte, 0)
	tracker.instanceStack = make([]wasmer.InstanceHandler, 0)
}

func (tracker *instanceTracker) StackSize() uint64 {
	return uint64(len(tracker.instanceStack))
}

func (tracker *instanceTracker) Instance() wasmer.InstanceHandler {
	return tracker.instance
}

func (tracker *instanceTracker) CodeHash() []byte {
	return tracker.codeHash
}

func (tracker *instanceTracker) ClearWarmInstanceCache() {
	if WarmInstancesEnabled {
		tracker.warmInstanceCache.Clear()
	}
}

func (tracker *instanceTracker) UseWarmInstance(codeHash []byte) bool {
	cachedObject, ok := tracker.warmInstanceCache.Get(codeHash)
	if !ok {
		return false
	}

	instance, ok := cachedObject.(wasmer.InstanceHandler)
	if !ok {
		return false
	}

	ok = instance.Reset()
	if !ok {
		// we must remove instance, which cleans it to free the memory
		tracker.warmInstanceCache.Remove(codeHash)
		return false
	}

	tracker.SetNewInstance(instance, Warm)
	return true
}

func (tracker *instanceTracker) ForceCleanInstance() {
	if check.IfNil(tracker.instance) {
		logRuntime.Trace("cannot clean, instance already nil")
		return
	}

	onStack := tracker.IsCodeHashOnTheStack(tracker.codeHash)
	coldOnlyEnabled := !WarmInstancesEnabled
	if onStack || coldOnlyEnabled {
		if tracker.instance.Clean() {
			tracker.updateNumRunningInstances(-1)
		}
	} else {
		tracker.warmInstanceCache.Remove(tracker.codeHash)
	}
	tracker.UnsetInstance()

	numWarmInstances, numColdInstances := tracker.NumRunningInstances()
	logRuntime.Trace("instance cleaned; num instances", "warm", numWarmInstances, "cold", numColdInstances)

}

func (tracker *instanceTracker) SaveAsWarmInstance() {
	lenCacheBeforeSaving := tracker.warmInstanceCache.Len()

	codeHashInWarmCache := tracker.warmInstanceCache.Has(tracker.codeHash)
	if !codeHashInWarmCache {
		logRuntime.Trace("warm instance not found, saving",
			"id", tracker.instance.Id(),
			"codeHash", tracker.codeHash,
		)
		tracker.warmInstanceCache.Put(
			tracker.codeHash,
			tracker.instance,
			1,
		)
	} else {
		tracker.updateNumRunningInstances(-1)
		logRuntime.Trace("warm instance already in cache", "id", tracker.instance.Id())
	}

	lenCacheAfterSaving := tracker.warmInstanceCache.Len()
	logRuntime.Trace("save warm instance length",
		"before", lenCacheBeforeSaving,
		"after", lenCacheAfterSaving,
	)

	logRuntime.Trace("save warm instance",
		"id", tracker.instance.Id(),
		"codeHash", tracker.codeHash,
	)
}

func (tracker *instanceTracker) SetCodeHash(codeHash []byte) {
	tracker.codeHash = codeHash
}

func (tracker *instanceTracker) SetNewInstance(instance wasmer.InstanceHandler, cacheLevel instanceCacheLevel) {
	tracker.ReplaceInstance(instance)
	tracker.cacheLevel = cacheLevel
	if cacheLevel != Warm {
		tracker.updateNumRunningInstances(+1)
	}
}

func (tracker *instanceTracker) ReplaceInstance(instance wasmer.InstanceHandler) {
	tracker.instance = instance

	if check.IfNil(tracker.instance) {
		logRuntime.Trace("ReplaceInstance: current instance is already nil")
		return
	}

	if !tracker.instance.AlreadyCleaned() {
		logRuntime.Trace("running instance about to be replaced without cleaning",
			"id", tracker.instance.Id(),
			"stacked", tracker.isInstanceOnTheStack(instance),
		)
	}
}

func (tracker *instanceTracker) UnsetInstance() {
	if check.IfNil(tracker.instance) {
		logRuntime.Trace("UnsetInstance: current instance is already nil")
		return
	}

	if !tracker.instance.AlreadyCleaned() {
		logRuntime.Trace("running instance about to be unset without cleaning",
			"id", tracker.instance.Id(),
			"stacked", tracker.isInstanceOnTheStack(tracker.instance),
		)
	}
	tracker.instance = nil
}

func (tracker *instanceTracker) LogCounts() {
	warm, cold := tracker.NumRunningInstances()
	logRuntime.Trace("num instances after starting new one",
		"warm", warm,
		"cold", cold,
		"total", tracker.numRunningInstances,
	)
}

// NumRunningInstances returns the number of currently running instances (cold and warm)
func (tracker *instanceTracker) NumRunningInstances() (int, int) {
	numWarmInstances := 0
	if WarmInstancesEnabled {
		numWarmInstances = tracker.warmInstanceCache.Len()
	}

	numColdInstances := tracker.numRunningInstances - numWarmInstances
	return numWarmInstances, numColdInstances
}

func (tracker *instanceTracker) IsCodeHashOnTheStack(codeHash []byte) bool {
	for _, stackedCodeHash := range tracker.codeHashStack {
		if bytes.Equal(codeHash, stackedCodeHash) {
			return true
		}
	}
	return false
}

func (tracker *instanceTracker) isInstanceOnTheStack(instance wasmer.InstanceHandler) bool {
	for _, stackedInstance := range tracker.instanceStack {
		if stackedInstance.Id() == instance.Id() {
			return true
		}
	}

	return false
}

func (tracker *instanceTracker) makeInstanceEvictionCallback() func(interface{}, interface{}) {
	return func(_ interface{}, value interface{}) {
		instance, ok := value.(wasmer.InstanceHandler)
		if !ok {
			return
		}

		logRuntime.Trace("evicted instance", "id", instance.Id())
		if instance.Clean() {
			tracker.updateNumRunningInstances(-1)
		}
	}
}

func (tracker *instanceTracker) updateNumRunningInstances(delta int) {
	tracker.numRunningInstances += delta
	logRuntime.Trace("num running instances updated", "delta", delta)
}

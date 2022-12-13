package contexts

import (
	"bytes"
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/storage"
	"github.com/ElrondNetwork/elrond-go-core/storage/lrucache"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen"
	"github.com/ElrondNetwork/wasm-vm-v1_4/wasmer"
)

type instanceCacheLevel int

const (
	// Warm indicates that an instance is warm
	Warm instanceCacheLevel = iota

	// Precompiled indicates that an instance has precompiled code
	Precompiled

	// Bytecode indicates that an instance must be compiled from bytecode
	Bytecode
)

var _ arwen.StateStack = (*instanceTracker)(nil)

var logTracker = logger.GetOrCreate("arwen/tracker")

type instanceTracker struct {
	codeHash            []byte
	numRunningInstances int
	warmInstanceCache   storage.Cacher
	instance            wasmer.InstanceHandler
	cacheLevel          instanceCacheLevel
	instanceStack       []wasmer.InstanceHandler
	codeHashStack       [][]byte

	instances map[string]wasmer.InstanceHandler
}

// NewInstanceTracker creates a new instanceTracker
func NewInstanceTracker() (*instanceTracker, error) {
	tracker := &instanceTracker{
		instances:           make(map[string]wasmer.InstanceHandler),
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
	logTracker.Trace("pushing instance", "id", tracker.instance.Id(), "codeHash", tracker.codeHash)
}

func (tracker *instanceTracker) PopSetActiveState() {
	instanceStackLen := len(tracker.instanceStack)
	if instanceStackLen == 0 {
		return
	}

	prevInstance := tracker.instanceStack[instanceStackLen-1]

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
		cold := !WarmInstancesEnabled
		if onStack || cold {
			if tracker.instance.Clean() {
				tracker.updateNumRunningInstances(-1)
			}
		}

		logTracker.Trace("pop instance", "id", tracker.instance.Id(), "codeHash", tracker.codeHash)
	}

	tracker.ReplaceInstance(prevInstance)
	tracker.instanceStack = tracker.instanceStack[:instanceStackLen-1]
	tracker.codeHash = tracker.codeHashStack[instanceStackLen-1]
	tracker.codeHashStack = tracker.codeHashStack[:instanceStackLen-1]
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

func (tracker *instanceTracker) UseWarmInstance(codeHash []byte, newCode bool) bool {
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
		tracker.warmInstanceCache.Remove(codeHash)
		return false
	}

	tracker.SetNewInstance(instance, Warm)
	return true
}

func (tracker *instanceTracker) ForceCleanInstance() {
	if check.IfNil(tracker.instance) {
		logTracker.Trace("cannot clean, instance already nil")
		return
	}

	onStack := tracker.IsCodeHashOnTheStack(tracker.codeHash)
	cold := !WarmInstancesEnabled
	if onStack || cold {
		if tracker.instance.Clean() {
			tracker.updateNumRunningInstances(-1)
		}
	} else {
		tracker.warmInstanceCache.Remove(tracker.codeHash)
	}
	tracker.UnsetInstance()

	numWarmInstances, numColdInstances := tracker.NumRunningInstances()
	logTracker.Trace("instance cleaned; num instances", "warm", numWarmInstances, "cold", numColdInstances)

}

func (tracker *instanceTracker) SaveAsWarmInstance() {
	lenCacheBeforeSaving := tracker.warmInstanceCache.Len()

	codeHashInWarmCache := tracker.warmInstanceCache.Has(tracker.codeHash)

	if codeHashInWarmCache {
		// Finding an instance in the warm cache at this point means that
		// context.instance is a new instance which must replace the one in the
		// warm cache, because they have the same bytecode. The old one is removed
		// and cleaned before the new one is added to the cache.
		logTracker.Trace("warm instance already in cache, evicting",
			"id", tracker.instance.Id(),
			"codeHash", tracker.codeHash)
		tracker.warmInstanceCache.Remove(tracker.codeHash)
	}

	logTracker.Trace("warm instance not found, saving",
		"id", tracker.instance.Id(),
		"codeHash", tracker.codeHash,
	)
	tracker.warmInstanceCache.Put(
		tracker.codeHash,
		tracker.instance,
		1,
	)

	lenCacheAfterSaving := tracker.warmInstanceCache.Len()
	logTracker.Trace("after saving, warm instance size",
		"before", lenCacheBeforeSaving,
		"after", lenCacheAfterSaving,
	)

	logTracker.Trace("save warm instance",
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
	tracker.instances[instance.Id()] = instance
}

func (tracker *instanceTracker) ReplaceInstance(instance wasmer.InstanceHandler) {
	if check.IfNil(tracker.instance) {
		logTracker.Trace("ReplaceInstance: previous instance was nil")
	}

	tracker.instance = instance
}

func (tracker *instanceTracker) UnsetInstance() {
	if check.IfNil(tracker.instance) {
		logTracker.Trace("UnsetInstance: current instance is already nil")
		return
	}

	logTracker.Trace("UnsetInstance",
		"id", tracker.instance.Id(),
		"codeHash", tracker.codeHash)
	tracker.instance = nil
	tracker.codeHash = nil

}

func (tracker *instanceTracker) LogCounts() {
	warm, cold := tracker.NumRunningInstances()
	logTracker.Trace("num instances after starting new one",
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

func (tracker *instanceTracker) CheckInstances() error {
	unclosedWarm := 0
	unclosedCold := 0

	warmInstanceCacheByID := make(map[string]wasmer.InstanceHandler)
	for _, key := range tracker.warmInstanceCache.Keys() {
		cachedObject, exists := tracker.warmInstanceCache.Get(key)
		if !exists {
			return fmt.Errorf("degenerate cache")
		}
		instance, ok := cachedObject.(wasmer.InstanceHandler)
		if !ok {
			return fmt.Errorf("degenerate cache")
		}
		warmInstanceCacheByID[instance.Id()] = instance
	}

	for id, instance := range tracker.instances {
		if instance.AlreadyCleaned() {
			continue
		}
		_, isWarm := warmInstanceCacheByID[id]
		if isWarm {
			unclosedWarm++
		} else {
			unclosedCold++
			logTracker.Trace("unclosed cold instance", "id", id)
		}
	}

	if unclosedWarm != tracker.warmInstanceCache.Len() {
		return fmt.Errorf(
			"there are %d closed warm instances in the cache",
			unclosedWarm)
	}

	if unclosedCold != 0 {
		return fmt.Errorf(
			"unclosed cold instances remaining: cold %d, warm %d",
			unclosedCold,
			unclosedWarm)
	}

	return nil
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

		logTracker.Trace("evicted instance", "id", instance.Id())
		if instance.Clean() {
			tracker.updateNumRunningInstances(-1)
		}
	}
}

func (tracker *instanceTracker) updateNumRunningInstances(delta int) {
	tracker.numRunningInstances += delta
	logTracker.Trace("num running instances updated", "delta", delta)
}

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
	logTracker.Trace("pushing instance", "id", tracker.instance.ID(), "codeHash", tracker.codeHash)
}

func (tracker *instanceTracker) PopSetActiveState() {
	instanceStackLen := len(tracker.instanceStack)
	if instanceStackLen == 0 {
		return
	}

	activeInstance := tracker.instance
	activeCodeHash := tracker.codeHash
	stackedPrevInstance := tracker.instanceStack[instanceStackLen-1]

	onStack := tracker.IsCodeHashOnTheStack(activeCodeHash)
	activeInstanceIsTopOfStack := stackedPrevInstance == activeInstance
	cold := !WarmInstancesEnabled

	if !activeInstanceIsTopOfStack && (onStack || cold) {
		tracker.cleanPoppedInstance(activeInstance, activeCodeHash)
	}

	tracker.ReplaceInstance(stackedPrevInstance)

	tracker.instanceStack = tracker.instanceStack[:instanceStackLen-1]
	tracker.codeHash = tracker.codeHashStack[instanceStackLen-1]
	tracker.codeHashStack = tracker.codeHashStack[:instanceStackLen-1]
}

func (tracker *instanceTracker) cleanPoppedInstance(instance wasmer.InstanceHandler, codeHash []byte) {
	if !check.IfNil(instance) {
		if instance.Clean() {
			tracker.updateNumRunningInstances(-1)
		}

		logTracker.Trace("clean popped instance", "id", instance.ID(), "codeHash", codeHash)
	}
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

func (tracker *instanceTracker) GetWarmInstance(codeHash []byte) (wasmer.InstanceHandler, bool) {
	cachedObject, ok := tracker.warmInstanceCache.Get(codeHash)
	if !ok {
		return nil, false
	}

	instance, ok := cachedObject.(wasmer.InstanceHandler)
	if !ok {
		return nil, false
	}

	return instance, true
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
	instance, ok := tracker.GetWarmInstance(codeHash)
	if !ok {
		return false
	}

	ok = instance.Reset()
	if !ok {
		tracker.warmInstanceCache.Remove(codeHash)
		return false
	}

	if newCode {
		// A warm instance was found, but newCode == true, meaning this is an
		// upgrade; the old warm instance must be cleaned
		tracker.ForceCleanInstance(false)
		return false
	}

	tracker.SetNewInstance(instance, Warm)
	return true
}

func (tracker *instanceTracker) ForceCleanInstance(bypassWarmAndStackChecks bool) {
	if check.IfNil(tracker.instance) {
		logTracker.Trace("cannot clean, instance already nil")
		return
	}

	defer func() {
		tracker.UnsetInstance()
		numWarmInstances, numColdInstances := tracker.NumRunningInstances()
		logTracker.Trace("instance cleaned; num instances", "warm", numWarmInstances, "cold", numColdInstances)
	}()

	if bypassWarmAndStackChecks {
		if tracker.instance.Clean() {
			tracker.updateNumRunningInstances(-1)
		}

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
			"id", tracker.instance.ID(),
			"codeHash", tracker.codeHash)
		tracker.warmInstanceCache.Remove(tracker.codeHash)
	}

	logTracker.Trace("warm instance not found, saving",
		"id", tracker.instance.ID(),
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
		"id", tracker.instance.ID(),
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
	tracker.instances[instance.ID()] = instance
}

func (tracker *instanceTracker) ReplaceInstance(instance wasmer.InstanceHandler) {
	var previousInstanceID string
	if check.IfNil(tracker.instance) {
		logTracker.Trace("ReplaceInstance: previous instance was nil")
		previousInstanceID = "nil"
	} else {
		previousInstanceID = tracker.instance.ID()
	}

	logTracker.Trace("replaced instance",
		"prev id", previousInstanceID,
		"new id", instance.ID())
	tracker.instance = instance
}

func (tracker *instanceTracker) UnsetInstance() {
	if check.IfNil(tracker.instance) {
		logTracker.Trace("UnsetInstance: current instance is already nil")
		return
	}

	logTracker.Trace("UnsetInstance",
		"id", tracker.instance.ID(),
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
		instance, ok := tracker.GetWarmInstance(key)
		if !ok {
			return fmt.Errorf("degenerate cache")
		}
		warmInstanceCacheByID[instance.ID()] = instance
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
		if stackedInstance.ID() == instance.ID() {
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

		logTracker.Trace("evicted instance", "id", instance.ID())
		if instance.Clean() {
			tracker.updateNumRunningInstances(-1)
		}
	}
}

func (tracker *instanceTracker) updateNumRunningInstances(delta int) {
	tracker.numRunningInstances += delta
	logTracker.Trace("num running instances updated", "delta", delta)
}

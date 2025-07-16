package contexts

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-storage-go/lrucache"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

var _ vmhost.InstanceTracker = (*instanceTracker)(nil)

type instanceCacheLevel int

var errTooManyInstances = errors.New("too many instances")

const (
	// Warm indicates that the instance to track is a warm instance
	Warm instanceCacheLevel = iota

	// Precompiled indicates that the instance to track is cold and has been created from precompiled code
	Precompiled

	// Bytecode indicates that the instance to track is cold and has been created from raw bytecode
	Bytecode
)

var _ vmhost.StateStack = (*instanceTracker)(nil)

var logTracker = logger.GetOrCreate("vm/tracker")

type instanceTracker struct {
	codeHash            []byte
	codeSize            uint64
	numRunningInstances int
	warmInstanceCache   Cacher
	instance            executor.Instance
	cacheLevel          instanceCacheLevel
	instanceStack       []executor.Instance
	codeHashStack       [][]byte
	codeSizeStack       []uint64

	instances map[string]executor.Instance
}

// NewInstanceTracker creates a new instanceTracker instance
func NewInstanceTracker() (*instanceTracker, error) {
	tracker := &instanceTracker{
		instances:           make(map[string]executor.Instance),
		instanceStack:       make([]executor.Instance, 0),
		codeHashStack:       make([][]byte, 0),
		codeSizeStack:       make([]uint64, 0),
		numRunningInstances: 0,
	}

	var err error
	instanceEvictedCallback := tracker.makeInstanceEvictionCallback()
	tracker.warmInstanceCache, err = lrucache.NewCacheWithEviction(warmCacheSize, instanceEvictedCallback)
	if err != nil {
		return nil, err
	}

	return tracker, nil
}

// InitState initializes the internal instanceTracker state
func (tracker *instanceTracker) InitState() {
	tracker.instance = nil
	tracker.codeHash = make([]byte, 0)
	tracker.instances = make(map[string]executor.Instance)
	tracker.codeSize = 0
}

// PushState pushes the active instance and codeHash on the state stacks
func (tracker *instanceTracker) PushState() {
	tracker.instanceStack = append(tracker.instanceStack, tracker.instance)
	tracker.codeHashStack = append(tracker.codeHashStack, tracker.codeHash)
	tracker.codeSizeStack = append(tracker.codeSizeStack, tracker.codeSize)
	logTracker.Trace("pushing instance", "id", tracker.instance.ID(), "codeHash", tracker.codeHash)
}

// PopSetActiveState pops the instance and codeHash from the state stacks and sets them as active
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
	if !activeInstanceIsTopOfStack && onStack {
		tracker.cleanPoppedInstance(activeInstance, activeCodeHash)
	}

	tracker.ReplaceInstance(stackedPrevInstance)

	tracker.instanceStack = tracker.instanceStack[:instanceStackLen-1]
	tracker.codeHash = tracker.codeHashStack[instanceStackLen-1]
	tracker.codeHashStack = tracker.codeHashStack[:instanceStackLen-1]

	tracker.codeSize = tracker.codeSizeStack[instanceStackLen-1]
	tracker.codeSizeStack = tracker.codeSizeStack[:instanceStackLen-1]
}

func (tracker *instanceTracker) cleanPoppedInstance(instance executor.Instance, codeHash []byte) {
	if !check.IfNil(instance) {
		if instance.Clean() {
			tracker.updateNumRunningInstances(-1)
		}

		logTracker.Trace("clean popped instance", "id", instance.ID(), "codeHash", codeHash)
	}
}

// PopDiscard does nothing for the instanceTracker
func (tracker *instanceTracker) PopDiscard() {
}

// ClearStateStack reinitializes the state stack.
func (tracker *instanceTracker) ClearStateStack() {
	tracker.codeHashStack = make([][]byte, 0)
	tracker.instanceStack = make([]executor.Instance, 0)
	tracker.codeSizeStack = make([]uint64, 0)
}

// StackSize returns the size of the instance stack
func (tracker *instanceTracker) StackSize() uint64 {
	return uint64(len(tracker.instanceStack))
}

// Instance returns the active instance
func (tracker *instanceTracker) Instance() executor.Instance {
	return tracker.instance
}

// GetWarmInstance retrieves a warm instance from the internal cache
func (tracker *instanceTracker) GetWarmInstance(codeHash []byte) (executor.Instance, bool) {
	cachedObject, ok := tracker.warmInstanceCache.Get(codeHash)
	if !ok {
		return nil, false
	}

	instance, ok := cachedObject.(executor.Instance)
	if !ok {
		return nil, false
	}

	return instance, true
}

// CodeHash returns the codeHash of the active instance
func (tracker *instanceTracker) CodeHash() []byte {
	return tracker.codeHash
}

// ClearWarmInstanceCache clears the internal warm instance cache
func (tracker *instanceTracker) ClearWarmInstanceCache() {
	tracker.warmInstanceCache.Clear()
}

// TrackedInstances returns the internal map of tracked instances
func (tracker *instanceTracker) TrackedInstances() map[string]executor.Instance {
	return tracker.instances
}

// UseWarmInstance attempts to retrieve a warm instance for the given codeHash
// and to set it as active; returns false if not possible
func (tracker *instanceTracker) UseWarmInstance(codeHash []byte, newCode bool) (bool, error) {
	instance, ok := tracker.GetWarmInstance(codeHash)
	if !ok {
		return false, nil
	}

	ok = instance.Reset()
	if !ok {
		tracker.warmInstanceCache.Remove(codeHash)
		return false, nil
	}

	if newCode {
		// A warm instance was found, but newCode == true, meaning this is an
		// upgrade; the old warm instance must be cleaned
		tracker.ForceCleanInstance(false)
		return false, nil
	}

	err := tracker.SetNewInstance(instance, Warm)
	return true, err
}

// ForceCleanInstance cleans the active instance and evicts it from the
// internal warm instance cache if possible
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
	if onStack {
		if tracker.instance.Clean() {
			tracker.updateNumRunningInstances(-1)
		}
	} else {
		tracker.warmInstanceCache.Remove(tracker.codeHash)
	}
}

// SaveAsWarmInstance saves the active instance into the internal warm instance cache
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

// SetCodeHash sets the active codeHash; it must correspond with the active instance
func (tracker *instanceTracker) SetCodeHash(codeHash []byte) {
	tracker.codeHash = codeHash
}

// SetCodeSize sets the size of the active code
func (tracker *instanceTracker) SetCodeSize(codeSize uint64) {
	tracker.codeSize = codeSize
}

// GetCodeSize returns the size of the active code
func (tracker *instanceTracker) GetCodeSize() uint64 {
	return tracker.codeSize
}

// SetNewInstance sets the given instance as active and tracks its creation
func (tracker *instanceTracker) SetNewInstance(instance executor.Instance, cacheLevel instanceCacheLevel) error {
	tracker.ReplaceInstance(instance)
	tracker.cacheLevel = cacheLevel
	if cacheLevel != Warm {
		tracker.updateNumRunningInstances(+1)
	}
	tracker.instances[instance.ID()] = instance

	if len(tracker.instances) >= warmCacheSize-1 {
		return errTooManyInstances
	}
	return nil
}

// ReplaceInstance replaces the currently active instance with the given one
func (tracker *instanceTracker) ReplaceInstance(instance executor.Instance) {
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

// UnsetInstance replaces the currently active instance with nil
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

// LogCounts prints the instance counter to the log
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
	numWarmInstances = tracker.warmInstanceCache.Len()
	numColdInstances := tracker.numRunningInstances - numWarmInstances
	return numWarmInstances, numColdInstances
}

// IsCodeHashOnTheStack returns true if the given codeHash is found on the codeHash stack
func (tracker *instanceTracker) IsCodeHashOnTheStack(codeHash []byte) bool {
	for _, stackedCodeHash := range tracker.codeHashStack {
		if bytes.Equal(codeHash, stackedCodeHash) {
			return true
		}
	}
	return false
}

// CheckInstances returns an error if there are tracked cold instances which
// have not been cleaned (leak detection)
func (tracker *instanceTracker) CheckInstances() error {
	unclosedWarm := 0
	unclosedCold := 0

	warmInstanceCacheByID := make(map[string]executor.Instance)
	for _, key := range tracker.warmInstanceCache.Keys() {
		instance, ok := tracker.GetWarmInstance(key)
		if !ok {
			return fmt.Errorf("degenerate cache")
		}
		warmInstanceCacheByID[instance.ID()] = instance
	}

	for id, instance := range tracker.instances {
		if instance.IsAlreadyCleaned() {
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

func (tracker *instanceTracker) makeInstanceEvictionCallback() func(interface{}, interface{}) {
	return func(_ interface{}, value interface{}) {
		instance, ok := value.(executor.Instance)
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

package wasmer2

import (
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/stretchr/testify/require"
)

func newTestInstanceState(t *testing.T) (*Wasmer2InstanceNew, executor.InstanceState, uint64) {
	wasmerExecutor, _, _ := newTestExecutor(t)
	instance := newCounterInstance(t, wasmerExecutor)

	// Reached outside a call, which only Wasmer 2.2 allows. Legitimate in a test of the
	// state itself; the VM must not do it, which is why executor.InstanceState promises the
	// shorter lifetime.
	state := instance.instanceState()

	memoryLength, err := state.MemoryLength()
	require.Nil(t, err)
	require.Greater(t, memoryLength, uint64(0))

	return instance, state, memoryLength
}

func TestWasmer2InstanceStateNew_MemoryStoreAndStrictLoad(t *testing.T) {
	instance, state, _ := newTestInstanceState(t)
	defer instance.Close()

	require.Nil(t, state.MemoryStore(64, []byte("hello")))

	loaded, err := state.MemoryLoad(64, 5)
	require.Nil(t, err)
	require.Equal(t, []byte("hello"), loaded)
}

func TestWasmer2InstanceStateNew_StrictLoadRejectsWhatClampedTruncates(t *testing.T) {
	instance, state, memoryLength := newTestInstanceState(t)
	defer instance.Close()

	lastByte := executor.MemPtr(memoryLength - 1)

	_, err := state.MemoryLoad(lastByte, 2)
	require.ErrorIs(t, err, executor.ErrMemoryBadBoundsUpper)

	// The truncating load is the behaviour contracts have already observed, so it stays
	// available, and distinct.
	clamped, err := state.MemoryLoadClamped(lastByte, 2)
	require.Nil(t, err)
	require.Len(t, clamped, 1)

	// An offset exactly at the end is a valid offset that nothing can be read from.
	_, err = state.MemoryLoad(executor.MemPtr(memoryLength), 1)
	require.ErrorIs(t, err, executor.ErrMemoryBadBoundsUpper)
}

func TestWasmer2InstanceStateNew_LoadBoundsAsymmetries(t *testing.T) {
	instance, state, memoryLength := newTestInstanceState(t)
	defer instance.Close()

	// A zero length returns empty before the offset is looked at.
	for _, badOffset := range []executor.MemPtr{-1, executor.MemPtr(memoryLength) + 1000} {
		loaded, err := state.MemoryLoad(badOffset, 0)
		require.Nil(t, err)
		require.Empty(t, loaded)
	}

	_, err := state.MemoryLoad(-1, 10)
	require.ErrorIs(t, err, executor.ErrMemoryBadBounds)

	_, err = state.MemoryLoad(0, -1)
	require.ErrorIs(t, err, executor.ErrMemoryNegativeLength)

	// An out-of-range offset takes precedence over a negative length.
	_, err = state.MemoryLoad(-1, -1)
	require.ErrorIs(t, err, executor.ErrMemoryBadBounds)

	// A length a contract is free to choose must not wrap the range into one that looks
	// valid.
	_, err = state.MemoryLoad(1, 0x7fffffff)
	require.ErrorIs(t, err, executor.ErrMemoryBadBoundsUpper)
}

func TestWasmer2InstanceStateNew_StoreIsNeverTruncated(t *testing.T) {
	instance, state, memoryLength := newTestInstanceState(t)
	defer instance.Close()

	err := state.MemoryStore(executor.MemPtr(memoryLength-1), []byte("ab"))
	require.ErrorIs(t, err, executor.ErrMemoryBadBoundsUpper)

	err = state.MemoryStore(-1, []byte("a"))
	require.ErrorIs(t, err, executor.ErrMemoryBadBoundsLower)

	// An empty write is accepted without validating the offset, as on the load side.
	require.Nil(t, state.MemoryStore(-1, nil))
}

func TestWasmer2InstanceStateNew_MemoryGrowReturnsThePagesBefore(t *testing.T) {
	instance, state, memoryLength := newTestInstanceState(t)
	defer instance.Close()

	pagesBefore, err := state.MemoryGrow(1)
	require.Nil(t, err)
	require.Equal(t, uint32(memoryLength/wasmPageSize), pagesBefore)

	grownLength, err := state.MemoryLength()
	require.Nil(t, err)
	require.Equal(t, memoryLength+wasmPageSize, grownLength)
}

func TestWasmer2InstanceStateNew_PointsUsed(t *testing.T) {
	instance, state, _ := newTestInstanceState(t)
	defer instance.Close()

	require.Nil(t, state.SetPointsUsedNew(4242))

	pointsUsed, err := state.GetPointsUsedNew()
	require.Nil(t, err)
	require.Equal(t, uint64(4242), pointsUsed)
}

func TestWasmer2InstanceStateNew_DetachedOnceTheInstanceIsClosed(t *testing.T) {
	instance, state, _ := newTestInstanceState(t)
	instance.Close()

	_, err := state.MemoryLength()
	require.ErrorIs(t, err, ErrDetachedInstanceState)

	_, err = state.MemoryLoad(0, 1)
	require.ErrorIs(t, err, ErrDetachedInstanceState)

	_, err = state.MemoryLoadClamped(0, 1)
	require.ErrorIs(t, err, ErrDetachedInstanceState)

	require.ErrorIs(t, state.MemoryStore(0, []byte("a")), ErrDetachedInstanceState)

	_, err = state.MemoryGrow(1)
	require.ErrorIs(t, err, ErrDetachedInstanceState)

	_, err = state.GetPointsUsedNew()
	require.ErrorIs(t, err, ErrDetachedInstanceState)

	require.ErrorIs(t, state.SetPointsUsedNew(1), ErrDetachedInstanceState)
}

// A state with no instance behind it is what hooks would get with no call in progress. It
// cannot arise, the executor announcing the instance around every C call, but it must not
// dereference nothing if it ever does.
func TestWasmer2InstanceStateNew_DetachedWithNoInstance(t *testing.T) {
	var instance *Wasmer2InstanceNew
	state := instance.instanceState()

	require.False(t, state.IsInterfaceNil())

	_, err := state.MemoryLength()
	require.ErrorIs(t, err, ErrDetachedInstanceState)
}

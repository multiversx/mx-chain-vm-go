package contexts

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/stretchr/testify/require"
)

func TestNewBigInt(t *testing.T) {
	t.Parallel()

	bigIntContext, err := NewBigIntContext()

	require.Nil(t, err)
	require.False(t, bigIntContext.IsInterfaceNil())
	require.NotNil(t, bigIntContext.values)
	require.NotNil(t, bigIntContext.stateStack)
	require.Equal(t, 0, len(bigIntContext.values))
	require.Equal(t, 0, len(bigIntContext.stateStack))
}

func TestBigIntContext_InitPushPopState(t *testing.T) {
	t.Parallel()

	value1, value2, value3 := int64(100), int64(200), int64(-42)
	bigIntContext, _ := NewBigIntContext()
	bigIntContext.InitState()

	// Create 2 values on the active state
	index1 := bigIntContext.Put(value1)
	require.Equal(t, int32(0), index1)

	index2 := bigIntContext.Put(value2)
	require.Equal(t, int32(1), index2)

	bigValue1, bigValue2 := bigIntContext.GetTwo(index1, index2)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Equal(t, big.NewInt(value2), bigValue2)

	// Copy active state to stack, then clean it. The previous 2 values should not
	// be accessible.
	bigIntContext.PushState()
	require.Equal(t, 1, len(bigIntContext.stateStack))
	bigIntContext.InitState()

	bigValue1, bigValue2 = bigIntContext.GetTwo(index1, index2)
	require.Equal(t, arwen.Zero, bigValue1)
	require.Equal(t, arwen.Zero, bigValue2)

	// Add a value on the current active state
	index3 := bigIntContext.Put(value3)
	require.Equal(t, int32(2), index3)
	bigValue3 := bigIntContext.GetOne(index3)
	require.Equal(t, big.NewInt(value3), bigValue3)

	// Copy active state to stack, then clean it. The previous 3 values should not
	// be accessible.
	bigIntContext.PushState()
	require.Equal(t, 2, len(bigIntContext.stateStack))
	bigIntContext.InitState()

	bigValue1, bigValue2, bigValue3 = bigIntContext.GetThree(index1, index2, index3)
	require.Equal(t, arwen.Zero, bigValue1)
	require.Equal(t, arwen.Zero, bigValue2)
	require.Equal(t, arwen.Zero, bigValue3)

	value4 := int64(84)
	index4 := bigIntContext.Put(value4)
	require.Equal(t, int32(3), index4)
	bigValue4 := bigIntContext.GetOne(index4)
	require.Equal(t, big.NewInt(value4), bigValue4)

	// Discard the top of the stack, losing value3; value4 should still be
	// accessible, since its in the active state.
	bigIntContext.PopDiscard()
	require.Equal(t, 1, len(bigIntContext.stateStack))
	bigValue4 = bigIntContext.GetOne(index4)
	require.Equal(t, big.NewInt(value4), bigValue4)

	// Restore the first active state by popping to the active state (which is
	// lost).
	bigIntContext.PopSetActiveState()
	require.Equal(t, 0, len(bigIntContext.stateStack))

	bigValue1, bigValue2 = bigIntContext.GetTwo(index1, index2)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Equal(t, big.NewInt(value2), bigValue2)
}

func TestBigIntContext_PutGet(t *testing.T) {
	t.Parallel()

	value1, value2, value3 := int64(100), int64(200), int64(-42)
	bigIntContext, _ := NewBigIntContext()

	index1 := bigIntContext.Put(value1)
	require.Equal(t, int32(0), index1)

	index2 := bigIntContext.Put(value2)
	require.Equal(t, int32(1), index2)

	index3 := bigIntContext.Put(value3)
	require.Equal(t, int32(2), index3)

	bigValue1, bigValue2 := bigIntContext.GetOne(index1), bigIntContext.GetOne(index2)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Equal(t, big.NewInt(value2), bigValue2)

	zeroRes := bigIntContext.GetOne(123)
	require.Equal(t, arwen.Zero, zeroRes)

	bigValue1, bigValue2 = bigIntContext.GetTwo(index1, index2)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Equal(t, big.NewInt(value2), bigValue2)

	bigValue1, bigValue2, zeroRes = bigIntContext.GetThree(index1, index2, 123)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Equal(t, arwen.Zero, zeroRes)

	bigValue1, bigValue2, bigValue3 := bigIntContext.GetThree(index1, index2, index3)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Equal(t, big.NewInt(value3), bigValue3)
}

func TestBigIntContext_PopSetActiveStateIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()

	bigIntContext, _ := NewBigIntContext()
	bigIntContext.PopSetActiveState()

	require.Equal(t, 0, len(bigIntContext.stateStack))
}

func TestBigIntContext_PopDiscardIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()

	bigIntContext, _ := NewBigIntContext()
	bigIntContext.PopDiscard()

	require.Equal(t, 0, len(bigIntContext.stateStack))
}

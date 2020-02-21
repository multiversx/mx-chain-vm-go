package contexts

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBigInt(t *testing.T) {
	t.Parallel()

	bigIntContext, err := NewBigIntContext()

	require.Nil(t, err)
	require.False(t, bigIntContext.IsInterfaceNil())
	require.NotNil(t, bigIntContext.mappedValues)
	require.NotNil(t, bigIntContext.stateStack)
	require.Equal(t, 0, len(bigIntContext.mappedValues))
	require.Equal(t, 0, len(bigIntContext.stateStack))
}

func TestBigIntContext_InitPushPopState(t *testing.T) {
	t.Parallel()

	bigIntContext, _ := NewBigIntContext()
	bigIntContext.InitState()

	bigIntContext.PushState()
	require.Equal(t, 1, len(bigIntContext.stateStack))

	bigIntContext.PopState()
	require.Equal(t, 0, len(bigIntContext.stateStack))

	bigIntContext.PushState()
	require.Equal(t, 1, len(bigIntContext.stateStack))

	bigIntContext.ClearStateStack()
	require.Equal(t, 0, len(bigIntContext.stateStack))
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

	bigRes1, bigRes2 := bigIntContext.GetOne(index1), bigIntContext.GetOne(index2)
	require.Equal(t, big.NewInt(value1), bigRes1)
	require.Equal(t, big.NewInt(value2), bigRes2)

	zeroRes := bigIntContext.GetOne(123)
	require.Equal(t, big.NewInt(0), zeroRes)

	bigRes1, bigRes2 = bigIntContext.GetTwo(index1, index2)
	require.Equal(t, big.NewInt(value1), bigRes1)
	require.Equal(t, big.NewInt(value2), bigRes2)

	bigRes1, bigRes2, zeroRes = bigIntContext.GetThree(index1, index2, 123)
	require.Equal(t, big.NewInt(value1), bigRes1)
	require.Equal(t, big.NewInt(value2), bigRes2)
	require.Equal(t, big.NewInt(0), zeroRes)

	bigRes1, bigRes2, bigRes3 := bigIntContext.GetThree(index1, index2, index3)
	require.Equal(t, big.NewInt(value1), bigRes1)
	require.Equal(t, big.NewInt(value2), bigRes2)
	require.Equal(t, big.NewInt(value3), bigRes3)
}

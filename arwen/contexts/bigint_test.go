package contexts

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/stretchr/testify/require"
)

func TestNewBigInt(t *testing.T) {
	t.Parallel()

	bic, err := NewBigIntContext()

	require.Nil(t, err)
	require.False(t, bic.IsInterfaceNil())
	require.NotNil(t, bic.mappedValues)
	require.NotNil(t, bic.stateStack)
	require.Equal(t, 0, len(bic.mappedValues))
	require.Equal(t, 0, len(bic.stateStack))
}

func TestBigIntContext_InitPushPopState(t *testing.T) {
	t.Parallel()

	bic, _ := NewBigIntContext()
	bic.InitState()

	err := bic.PopState()
	require.Equal(t, arwen.StateStackUnderflow, err)

	bic.PushState()
	require.Equal(t, 1, len(bic.stateStack))

	err = bic.PopState()
	require.Nil(t, err)
	require.Equal(t, 0, len(bic.stateStack))
}

func TestBigIntContext_PutGet(t *testing.T) {
	t.Parallel()

	value1, value2, value3 := int64(100), int64(200), int64(-42)
	bic, _ := NewBigIntContext()

	index1 := bic.Put(value1)
	require.Equal(t, int32(0), index1)

	index2 := bic.Put(value2)
	require.Equal(t, int32(1), index2)

	index3 := bic.Put(value3)
	require.Equal(t, int32(2), index3)

	bigRes1, bigRes2 := bic.GetOne(index1), bic.GetOne(index2)
	require.Equal(t, big.NewInt(value1), bigRes1)
	require.Equal(t, big.NewInt(value2), bigRes2)

	zeroRes := bic.GetOne(123)
	require.Equal(t, big.NewInt(0), zeroRes)

	bigRes1, bigRes2 = bic.GetTwo(index1, index2)
	require.Equal(t, big.NewInt(value1), bigRes1)
	require.Equal(t, big.NewInt(value2), bigRes2)

	bigRes1, bigRes2, zeroRes = bic.GetThree(index1, index2, 123)
	require.Equal(t, big.NewInt(value1), bigRes1)
	require.Equal(t, big.NewInt(value2), bigRes2)
	require.Equal(t, big.NewInt(0), zeroRes)

	bigRes1, bigRes2, bigRes3 := bic.GetThree(index1, index2, index3)
	require.Equal(t, big.NewInt(value1), bigRes1)
	require.Equal(t, big.NewInt(value2), bigRes2)
	require.Equal(t, big.NewInt(value3), bigRes3)
}

package contexts

import (
	"crypto/elliptic"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/stretchr/testify/require"
)

func TestNewManagedType(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostStub{}

	managedTypeContext, err := NewManagedTypeContext(host)

	require.Nil(t, err)
	require.False(t, managedTypeContext.IsInterfaceNil())
	require.NotNil(t, managedTypeContext.bigIntValues)
	require.NotNil(t, managedTypeContext.bigIntStateStack)
	require.NotNil(t, managedTypeContext.ecValues)
	require.NotNil(t, managedTypeContext.ecStateStack)
	require.Equal(t, 0, len(managedTypeContext.bigIntValues))
	require.Equal(t, 0, len(managedTypeContext.bigIntStateStack))
	require.Equal(t, 0, len(managedTypeContext.ecValues))
	require.Equal(t, 0, len(managedTypeContext.ecStateStack))
}

func TestManagedTypeContext_InitPushPopState(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}
	value1, value2, value3 := int64(100), int64(200), int64(-42)
	p224ec, p256ec, p384ec := elliptic.P224().Params(), elliptic.P256().Params(), elliptic.P384().Params()
	managedTypeContext, _ := NewManagedTypeContext(host)
	managedTypeContext.InitState()

	// Create 2 bigInt and 2 EC values on the active state
	index1 := managedTypeContext.PutBigInt(value1)
	require.Equal(t, int32(0), index1)
	index2 := managedTypeContext.PutBigInt(value2)
	require.Equal(t, int32(1), index2)

	bigValue1, err := managedTypeContext.GetBigInt(index1)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Equal(t, nil, err)
	bigValue2, err := managedTypeContext.GetBigInt(index2)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Equal(t, nil, err)

	ecIndex1 := managedTypeContext.PutEllipticCurve(p224ec)
	require.Equal(t, int32(0), ecIndex1)
	ecIndex2 := managedTypeContext.PutEllipticCurve(p256ec)
	require.Equal(t, int32(1), ecIndex2)

	ec1, err := managedTypeContext.GetEllipticCurve(ecIndex1)
	require.Nil(t, err)
	require.Equal(t, p224ec, ec1)
	ec2, err := managedTypeContext.GetEllipticCurve(ecIndex2)
	require.Nil(t, err)
	require.Equal(t, p256ec, ec2)

	// Copy active state to stack, then clean it. The previous 2 values should not
	// be accessible.
	managedTypeContext.PushState()
	require.Equal(t, 1, len(managedTypeContext.bigIntStateStack))
	managedTypeContext.InitState()

	bigValue1, err = managedTypeContext.GetBigInt(index1)
	require.Nil(t, bigValue1)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigValue2, err = managedTypeContext.GetBigInt(index2)
	require.Nil(t, bigValue2)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)

	ec1, err = managedTypeContext.GetEllipticCurve(ecIndex1)
	require.Nil(t, ec1)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)
	ec2, err = managedTypeContext.GetEllipticCurve(ecIndex2)
	require.Nil(t, ec2)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)

	// Add a value on the current active state
	index3 := managedTypeContext.PutBigInt(value3)
	require.Equal(t, int32(0), index3)
	bigValue3, err := managedTypeContext.GetBigInt(index3)
	require.Equal(t, big.NewInt(value3), bigValue3)
	require.Equal(t, nil, err)

	ecIndex3 := managedTypeContext.PutEllipticCurve(p384ec)
	require.Equal(t, int32(0), ecIndex3)
	ec3, err := managedTypeContext.GetEllipticCurve(ecIndex3)
	require.Nil(t, err)
	require.Equal(t, p384ec, ec3)

	// Copy active state to stack, then clean it. The previous 3 values should not
	// be accessible.
	managedTypeContext.PushState()
	require.Equal(t, 2, len(managedTypeContext.bigIntStateStack))
	managedTypeContext.InitState()

	bigValue1, err = managedTypeContext.GetBigInt(index1)
	require.Nil(t, bigValue1)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigValue2, err = managedTypeContext.GetBigInt(index2)
	require.Nil(t, bigValue2)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigValue3, err = managedTypeContext.GetBigInt(index3)
	require.Nil(t, bigValue3)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)

	ec1, err = managedTypeContext.GetEllipticCurve(ecIndex1)
	require.Nil(t, ec1)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)
	ec2, err = managedTypeContext.GetEllipticCurve(ecIndex2)
	require.Nil(t, ec2)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)
	ec3, err = managedTypeContext.GetEllipticCurve(ecIndex3)
	require.Nil(t, ec3)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)

	value4 := int64(84)
	index4 := managedTypeContext.PutBigInt(value4)
	require.Equal(t, int32(0), index4)
	bigValue4, err := managedTypeContext.GetBigInt(index4)
	require.Equal(t, big.NewInt(value4), bigValue4)
	require.Equal(t, nil, err)

	p521ec := elliptic.P521().Params()
	ecIndex4 := managedTypeContext.PutEllipticCurve(p521ec)
	require.Equal(t, int32(0), ecIndex4)
	ec4, err := managedTypeContext.GetEllipticCurve(ecIndex4)
	require.Equal(t, p521ec, ec4)
	require.Nil(t, err)

	// Discard the top of the stack, losing value3; value4 should still be
	// accessible, since its in the active state.
	managedTypeContext.PopDiscard()
	require.Equal(t, 1, len(managedTypeContext.bigIntStateStack))
	bigValue4, err = managedTypeContext.GetBigInt(index4)
	require.Equal(t, big.NewInt(value4), bigValue4)
	require.Equal(t, nil, err)

	ec4, err = managedTypeContext.GetEllipticCurve(ecIndex4)
	require.Equal(t, p521ec, ec4)
	require.Nil(t, err)
	// Restore the first active state by popping to the active state (which is
	// lost).
	managedTypeContext.PopSetActiveState()
	require.Equal(t, 0, len(managedTypeContext.bigIntStateStack))

	bigValue1, err = managedTypeContext.GetBigInt(index1)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Equal(t, nil, err)
	bigValue2, err = managedTypeContext.GetBigInt(index2)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Equal(t, nil, err)

	ec1, err = managedTypeContext.GetEllipticCurve(ecIndex1)
	require.Nil(t, err)
	require.Equal(t, p224ec, ec1)
	ec2, err = managedTypeContext.GetEllipticCurve(ecIndex2)
	require.Nil(t, err)
	require.Equal(t, p256ec, ec2)
}

func TestManagedTypeContext_PutGet(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}

	value1, value2, value3, value4 := int64(100), int64(200), int64(-42), int64(-80)
	p224ec, p256ec, p384ec, p521ec := elliptic.P224().Params(), elliptic.P256().Params(), elliptic.P384().Params(), elliptic.P521().Params()
	managedTypeContext, _ := NewManagedTypeContext(host)

	index1 := managedTypeContext.PutBigInt(value1)
	require.Equal(t, int32(0), index1)
	index2 := managedTypeContext.PutBigInt(value2)
	require.Equal(t, int32(1), index2)
	index3 := managedTypeContext.PutBigInt(value3)
	require.Equal(t, int32(2), index3)

	ecIndex1 := managedTypeContext.PutEllipticCurve(p224ec)
	require.Equal(t, int32(0), ecIndex1)
	ecIndex2 := managedTypeContext.PutEllipticCurve(p256ec)
	require.Equal(t, int32(1), ecIndex2)
	ecIndex3 := managedTypeContext.PutEllipticCurve(p384ec)
	require.Equal(t, int32(2), ecIndex3)

	bigValue1, err := managedTypeContext.GetBigInt(index1)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Equal(t, nil, err)
	bigValue2, err := managedTypeContext.GetBigInt(index2)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Equal(t, nil, err)
	bigValue4, err := managedTypeContext.GetBigInt(int32(3))
	require.Nil(t, bigValue4)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigValue4 = managedTypeContext.GetBigIntOrCreate(3)
	require.Equal(t, big.NewInt(0), bigValue4)

	ec1, err := managedTypeContext.GetEllipticCurve(ecIndex1)
	require.Nil(t, err)
	require.Equal(t, p224ec, ec1)
	ec2, err := managedTypeContext.GetEllipticCurve(ecIndex2)
	require.Nil(t, err)
	require.Equal(t, p256ec, ec2)
	ec4, err := managedTypeContext.GetEllipticCurve(int32(3))
	require.Nil(t, ec4)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)

	index4 := managedTypeContext.PutBigInt(value4)
	require.Equal(t, int32(4), index4)
	bigValue4 = managedTypeContext.GetBigIntOrCreate(4)
	require.Equal(t, big.NewInt(value4), bigValue4)

	ecIndex4 := managedTypeContext.PutEllipticCurve(p521ec)
	require.Equal(t, int32(3), ecIndex4)
	ec4, err = managedTypeContext.GetEllipticCurve(ecIndex4)
	require.Nil(t, err)
	require.Equal(t, p521ec, ec4)

	bigValue, err := managedTypeContext.GetBigInt(123)
	require.Nil(t, bigValue)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigValue1, err = managedTypeContext.GetBigInt(index1)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Equal(t, nil, err)
	bigValue2, err = managedTypeContext.GetBigInt(index2)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Equal(t, nil, err)

	bigValue1, err = managedTypeContext.GetBigInt(index1)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Equal(t, nil, err)
	bigValue2, err = managedTypeContext.GetBigInt(index2)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Equal(t, nil, err)
	bigValue3, err := managedTypeContext.GetBigInt(index3)
	require.Equal(t, big.NewInt(value3), bigValue3)
	require.Equal(t, nil, err)

}

func TestManagedTypeContext_PopSetActiveStateIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}

	managedTypeContext, _ := NewManagedTypeContext(host)
	managedTypeContext.PopSetActiveState()

	require.Equal(t, 0, len(managedTypeContext.bigIntStateStack))
	require.Equal(t, 0, len(managedTypeContext.ecStateStack))

}

func TestManagedTypeContext_PopDiscardIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}

	managedTypeContext, _ := NewManagedTypeContext(host)
	managedTypeContext.PopDiscard()

	require.Equal(t, 0, len(managedTypeContext.bigIntStateStack))
	require.Equal(t, 0, len(managedTypeContext.ecStateStack))

}

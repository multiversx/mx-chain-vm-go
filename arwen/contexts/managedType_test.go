package contexts

import (
	"crypto/elliptic"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/context"
	"github.com/stretchr/testify/require"
)

func TestNewManagedTypes(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostStub{}

	managedTypesContext, err := NewManagedTypesContext(host)

	require.Nil(t, err)
	require.False(t, managedTypesContext.IsInterfaceNil())
	require.NotNil(t, managedTypesContext.managedTypesValues.bigIntValues)
	require.NotNil(t, managedTypesContext.managedTypesValues.ecValues)
	require.NotNil(t, managedTypesContext.managedTypesValues.mBufferValues)
	require.NotNil(t, managedTypesContext.managedTypesStack)
	require.Equal(t, 0, len(managedTypesContext.managedTypesValues.bigIntValues))
	require.Equal(t, 0, len(managedTypesContext.managedTypesValues.ecValues))
	require.Equal(t, 0, len(managedTypesContext.managedTypesValues.mBufferValues))
	require.Equal(t, 0, len(managedTypesContext.managedTypesStack))
}

func TestManagedTypesContext_Randomness(t *testing.T) {
	t.Parallel()

	mockRuntime := &contextmock.RuntimeContextMock{
		CurrentTxHash: []byte{0xf, 0xf, 0xf, 0xf, 0xf, 0xf},
	}
	host := &contextmock.VMHostMock{
		RuntimeContext: mockRuntime,
	}
	mockBlockchain := &contextmock.BlockchainHookStub{
		CurrentRandomSeedCalled: func() []byte {
			return []byte{0xf, 0xf, 0xf, 0xf, 0xa, 0xb}
		},
	}
	blockchainContext, _ := NewBlockchainContext(host, mockBlockchain)
	host.BlockchainContext = blockchainContext
	copyHost := host

	managedTypesContext, _ := NewManagedTypesContext(host)
	require.Nil(t, managedTypesContext.srr)
	managedTypesContext.initRandomizer()
	firstRandomizer := managedTypesContext.srr

	managedTypesContextCopy, _ := NewManagedTypesContext(copyHost)
	require.Nil(t, managedTypesContextCopy.srr)
	managedTypesContextCopy.initRandomizer()
	secondRandomizer := managedTypesContextCopy.srr

	require.Equal(t, firstRandomizer, secondRandomizer)

	prg := managedTypesContext.GetRandReader()
	a := make([]byte, 100)
	prg.Read(a)
	b := make([]byte, 100)
	for i := 0; i < 1000; i++ {
		prg.Read(b)
		require.NotEqual(t, a, b)
		copy(a, b)
	}
}

func TestManagedTypesContext_ClearStateStack(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}
	value1, value2 := int64(100), int64(200)
	p224ec, p256ec := elliptic.P224().Params(), elliptic.P256().Params()
	managedTypesContext, _ := NewManagedTypesContext(host)
	managedTypesContext.InitState()

	index1 := managedTypesContext.PutBigInt(value1)
	index2 := managedTypesContext.PutBigInt(value2)
	ecIndex1 := managedTypesContext.PutEllipticCurve(p224ec)
	ecIndex2 := managedTypesContext.PutEllipticCurve(p256ec)

	managedTypesContext.PushState()
	require.Equal(t, 1, len(managedTypesContext.managedTypesStack))
	managedTypesContext.ClearStateStack()
	require.Equal(t, 0, len(managedTypesContext.managedTypesStack))

	bigValue1, err := managedTypesContext.GetBigInt(index1)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Nil(t, err)
	bigValue2, err := managedTypesContext.GetBigInt(index2)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Nil(t, err)
	ec1, err := managedTypesContext.GetEllipticCurve(ecIndex1)
	require.Nil(t, err)
	require.Equal(t, p224ec, ec1)
	ec2, err := managedTypesContext.GetEllipticCurve(ecIndex2)
	require.Nil(t, err)
	require.Equal(t, p256ec, ec2)

	managedTypesContext.InitState()
	bigValue1, err = managedTypesContext.GetBigInt(index1)
	require.Nil(t, bigValue1)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigValue2, err = managedTypesContext.GetBigInt(index2)
	require.Nil(t, bigValue2)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	ec1, err = managedTypesContext.GetEllipticCurve(ecIndex1)
	require.Nil(t, ec1)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)
	ec2, err = managedTypesContext.GetEllipticCurve(ecIndex2)
	require.Nil(t, ec2)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)
}

func TestManagedTypesContext_InitPushPopState(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}
	value1, value2, value3 := int64(100), int64(200), int64(-42)
	p224ec, p256ec, p384ec, p521ec := elliptic.P224().Params(), elliptic.P256().Params(), elliptic.P384().Params(), elliptic.P521().Params()
	bytes := []byte{2, 234, 64, 255}
	managedTypesContext, _ := NewManagedTypesContext(host)
	managedTypesContext.InitState()

	// Create 2 bigInt,2 EC, 2 managedBuffers on the active state
	index1 := managedTypesContext.PutBigInt(value1)
	require.Equal(t, int32(0), index1)
	index2 := managedTypesContext.PutBigInt(value2)
	require.Equal(t, int32(1), index2)

	bigValue1, err := managedTypesContext.GetBigInt(index1)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Nil(t, err)
	bigValue2, err := managedTypesContext.GetBigInt(index2)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Nil(t, err)

	ecIndex1 := managedTypesContext.PutEllipticCurve(p224ec)
	require.Equal(t, int32(0), ecIndex1)
	ecIndex2 := managedTypesContext.PutEllipticCurve(p256ec)
	require.Equal(t, int32(1), ecIndex2)

	ec1, err := managedTypesContext.GetEllipticCurve(ecIndex1)
	require.Nil(t, err)
	require.Equal(t, p224ec, ec1)
	ec2, err := managedTypesContext.GetEllipticCurve(ecIndex2)
	require.Nil(t, err)
	require.Equal(t, p256ec, ec2)

	mBufferHandle1 := managedTypesContext.NewManagedBufferFromBytes(bytes)
	mBuffer, _ := managedTypesContext.GetBytes(mBufferHandle1)
	require.Equal(t, bytes, mBuffer)

	p224NormalGasCostMultiplier := managedTypesContext.Get100xCurveGasCostMultiplier(ecIndex1)
	require.Equal(t, int32(100), p224NormalGasCostMultiplier)
	p256NormalGasCostMultiplier := managedTypesContext.Get100xCurveGasCostMultiplier(ecIndex2)
	require.Equal(t, int32(135), p256NormalGasCostMultiplier)
	p224ScalarMultGasCostMultiplier := managedTypesContext.GetScalarMult100xCurveGasCostMultiplier(ecIndex1)
	require.Equal(t, int32(100), p224ScalarMultGasCostMultiplier)
	p256ScalarMultGasCostMultiplier := managedTypesContext.GetScalarMult100xCurveGasCostMultiplier(ecIndex2)
	require.Equal(t, int32(110), p256ScalarMultGasCostMultiplier)
	p224UCompressedGasCostMultiplier := managedTypesContext.GetUCompressed100xCurveGasCostMultiplier(ecIndex1)
	require.Equal(t, int32(2000), p224UCompressedGasCostMultiplier)
	p256UCompressedGasCostMultiplier := managedTypesContext.GetUCompressed100xCurveGasCostMultiplier(ecIndex2)
	require.Equal(t, int32(100), p256UCompressedGasCostMultiplier)

	// Copy active state to stack, then clean it. The previous 2 values should not
	// be accessible.
	managedTypesContext.PushState()
	require.Equal(t, 1, len(managedTypesContext.managedTypesStack))
	managedTypesContext.InitState()

	bigValue1, err = managedTypesContext.GetBigInt(index1)
	require.Nil(t, bigValue1)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigValue2, err = managedTypesContext.GetBigInt(index2)
	require.Nil(t, bigValue2)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigValue1, bigValue2, err = managedTypesContext.GetTwoBigInt(index1, index2)
	require.Nil(t, bigValue1, bigValue2)
	require.Equal(t, err, arwen.ErrNoBigIntUnderThisHandle)

	ec1, err = managedTypesContext.GetEllipticCurve(ecIndex1)
	require.Nil(t, ec1)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)
	ec2, err = managedTypesContext.GetEllipticCurve(ecIndex2)
	require.Nil(t, ec2)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)

	mBuffer, err = managedTypesContext.GetBytes(mBufferHandle1)
	require.Nil(t, mBuffer)
	require.Equal(t, arwen.ErrNoManagedBufferUnderThisHandle, err)

	p224NormalGasCostMultiplier = managedTypesContext.Get100xCurveGasCostMultiplier(ecIndex1)
	require.Equal(t, int32(-1), p224NormalGasCostMultiplier)
	p256NormalGasCostMultiplier = managedTypesContext.Get100xCurveGasCostMultiplier(ecIndex2)
	require.Equal(t, int32(-1), p256NormalGasCostMultiplier)
	p224ScalarMultGasCostMultiplier = managedTypesContext.GetScalarMult100xCurveGasCostMultiplier(ecIndex1)
	require.Equal(t, int32(-1), p224ScalarMultGasCostMultiplier)
	p256ScalarMultGasCostMultiplier = managedTypesContext.GetScalarMult100xCurveGasCostMultiplier(ecIndex2)
	require.Equal(t, int32(-1), p256ScalarMultGasCostMultiplier)
	p224UCompressedGasCostMultiplier = managedTypesContext.GetUCompressed100xCurveGasCostMultiplier(ecIndex1)
	require.Equal(t, int32(-1), p224UCompressedGasCostMultiplier)
	p256UCompressedGasCostMultiplier = managedTypesContext.GetUCompressed100xCurveGasCostMultiplier(ecIndex2)
	require.Equal(t, int32(-1), p256UCompressedGasCostMultiplier)

	// Add a value on the current active state
	index3 := managedTypesContext.PutBigInt(value3)
	require.Equal(t, int32(0), index3)
	bigValue3, err := managedTypesContext.GetBigInt(index3)
	require.Equal(t, big.NewInt(value3), bigValue3)
	require.Nil(t, err)

	ecIndex3 := managedTypesContext.PutEllipticCurve(p384ec)
	require.Equal(t, int32(0), ecIndex3)
	ec3, err := managedTypesContext.GetEllipticCurve(ecIndex3)
	require.Nil(t, err)
	require.Equal(t, p384ec, ec3)

	p384NormalGasCostMultiplier := managedTypesContext.Get100xCurveGasCostMultiplier(ecIndex3)
	require.Equal(t, int32(200), p384NormalGasCostMultiplier)
	p384ScalarMultGasCostMultiplier := managedTypesContext.GetScalarMult100xCurveGasCostMultiplier(ecIndex3)
	require.Equal(t, int32(150), p384ScalarMultGasCostMultiplier)
	p384UCompressedGasCostMultiplier := managedTypesContext.GetUCompressed100xCurveGasCostMultiplier(ecIndex3)
	require.Equal(t, int32(200), p384UCompressedGasCostMultiplier)

	// Copy active state to stack, then clean it. The previous 3 values should not
	// be accessible.
	managedTypesContext.PushState()
	require.Equal(t, 2, len(managedTypesContext.managedTypesStack))
	managedTypesContext.InitState()

	bigValue1, err = managedTypesContext.GetBigInt(index1)
	require.Nil(t, bigValue1)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigValue2, err = managedTypesContext.GetBigInt(index2)
	require.Nil(t, bigValue2)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigValue3, err = managedTypesContext.GetBigInt(index3)
	require.Nil(t, bigValue3)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)

	ec1, err = managedTypesContext.GetEllipticCurve(ecIndex1)
	require.Nil(t, ec1)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)
	ec2, err = managedTypesContext.GetEllipticCurve(ecIndex2)
	require.Nil(t, ec2)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)
	ec3, err = managedTypesContext.GetEllipticCurve(ecIndex3)
	require.Nil(t, ec3)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)

	value4 := int64(84)
	index4 := managedTypesContext.PutBigInt(value4)
	require.Equal(t, int32(0), index4)
	bigValue4, err := managedTypesContext.GetBigInt(index4)
	require.Equal(t, big.NewInt(value4), bigValue4)
	require.Nil(t, err)
	bigValue4, bigValue3, err = managedTypesContext.GetTwoBigInt(index4, int32(1))
	require.Nil(t, bigValue3)
	require.Nil(t, bigValue4)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)

	ecIndex4 := managedTypesContext.PutEllipticCurve(p521ec)
	require.Equal(t, int32(0), ecIndex4)
	ec4, err := managedTypesContext.GetEllipticCurve(ecIndex4)
	require.Equal(t, p521ec, ec4)
	require.Nil(t, err)

	p521NormalGasCostMultiplier := managedTypesContext.Get100xCurveGasCostMultiplier(ecIndex4)
	require.Equal(t, int32(250), p521NormalGasCostMultiplier)
	p521ScalarMultGasCostMultiplier := managedTypesContext.GetScalarMult100xCurveGasCostMultiplier(ecIndex4)
	require.Equal(t, int32(190), p521ScalarMultGasCostMultiplier)
	p521UCompressedGasCostMultiplier := managedTypesContext.GetUCompressed100xCurveGasCostMultiplier(ecIndex4)
	require.Equal(t, int32(400), p521UCompressedGasCostMultiplier)

	// Discard the top of the stack, losing value3; value4 should still be
	// accessible, since its in the active state.
	managedTypesContext.PopDiscard()
	require.Equal(t, 1, len(managedTypesContext.managedTypesStack))
	bigValue4, err = managedTypesContext.GetBigInt(index4)
	require.Equal(t, big.NewInt(value4), bigValue4)
	require.Nil(t, err)

	ec4, err = managedTypesContext.GetEllipticCurve(ecIndex4)
	require.Equal(t, p521ec, ec4)
	require.Nil(t, err)
	// Restore the first active state by popping to the active state (which is
	// lost).
	managedTypesContext.PopSetActiveState()
	require.Equal(t, 0, len(managedTypesContext.managedTypesStack))

	bigValue1, err = managedTypesContext.GetBigInt(index1)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Nil(t, err)
	bigValue2, err = managedTypesContext.GetBigInt(index2)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Nil(t, err)

	bigValue1, bigValue2, err = managedTypesContext.GetTwoBigInt(index1, index2)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Nil(t, err)

	ec1, err = managedTypesContext.GetEllipticCurve(ecIndex1)
	require.Nil(t, err)
	require.Equal(t, p224ec, ec1)
	ec2, err = managedTypesContext.GetEllipticCurve(ecIndex2)
	require.Nil(t, err)
	require.Equal(t, p256ec, ec2)
}

func TestManagedTypesContext_PutGetBigInt(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}

	value1, value2, value3, value4 := int64(100), int64(200), int64(-42), int64(-80)
	managedTypesContext, _ := NewManagedTypesContext(host)

	index1 := managedTypesContext.PutBigInt(value1)
	require.Equal(t, int32(0), index1)
	index2 := managedTypesContext.PutBigInt(value2)
	require.Equal(t, int32(1), index2)
	index3 := managedTypesContext.PutBigInt(value3)
	require.Equal(t, int32(2), index3)

	bigValue1, err := managedTypesContext.GetBigInt(index1)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Nil(t, err)
	bigValue2, err := managedTypesContext.GetBigInt(index2)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Nil(t, err)
	bigValue4, err := managedTypesContext.GetBigInt(int32(3))
	require.Nil(t, bigValue4)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigValue4 = managedTypesContext.GetBigIntOrCreate(3)
	require.Equal(t, big.NewInt(0), bigValue4)

	index4 := managedTypesContext.PutBigInt(value4)
	require.Equal(t, int32(4), index4)
	bigValue4 = managedTypesContext.GetBigIntOrCreate(4)
	require.Equal(t, big.NewInt(value4), bigValue4)

	bigValue, err := managedTypesContext.GetBigInt(123)
	require.Nil(t, bigValue)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigValue1, err = managedTypesContext.GetBigInt(index1)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Nil(t, err)
	bigValue2, err = managedTypesContext.GetBigInt(index2)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Nil(t, err)

	bigValue1, err = managedTypesContext.GetBigInt(index1)
	require.Equal(t, big.NewInt(value1), bigValue1)
	require.Nil(t, err)
	bigValue2, err = managedTypesContext.GetBigInt(index2)
	require.Equal(t, big.NewInt(value2), bigValue2)
	require.Nil(t, err)
	bigValue3, err := managedTypesContext.GetBigInt(index3)
	require.Equal(t, big.NewInt(value3), bigValue3)
	require.Nil(t, err)
}

func TestManagedTypesContext_PutGetEllipticCurves(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}

	p224ec, p256ec, p384ec, p521ec := elliptic.P224().Params(), elliptic.P256().Params(), elliptic.P384().Params(), elliptic.P521().Params()
	managedTypesContext, _ := NewManagedTypesContext(host)

	ecIndex1 := managedTypesContext.PutEllipticCurve(p224ec)
	require.Equal(t, int32(0), ecIndex1)
	ecIndex2 := managedTypesContext.PutEllipticCurve(p256ec)
	require.Equal(t, int32(1), ecIndex2)
	ecIndex3 := managedTypesContext.PutEllipticCurve(p384ec)
	require.Equal(t, int32(2), ecIndex3)

	p224PrivKeyByteLength := managedTypesContext.GetPrivateKeyByteLengthEC(ecIndex1)
	require.Equal(t, int32(28), p224PrivKeyByteLength)
	p256PrivKeyByteLength := managedTypesContext.GetPrivateKeyByteLengthEC(ecIndex2)
	require.Equal(t, int32(32), p256PrivKeyByteLength)
	p384PrivKeyByteLength := managedTypesContext.GetPrivateKeyByteLengthEC(ecIndex3)
	require.Equal(t, int32(48), p384PrivKeyByteLength)
	nonExistentCurvePrivKeyByteLength := managedTypesContext.GetPrivateKeyByteLengthEC(int32(3))
	require.Equal(t, int32(-1), nonExistentCurvePrivKeyByteLength)

	ec1, err := managedTypesContext.GetEllipticCurve(ecIndex1)
	require.Nil(t, err)
	require.Equal(t, p224ec, ec1)
	ec2, err := managedTypesContext.GetEllipticCurve(ecIndex2)
	require.Nil(t, err)
	require.Equal(t, p256ec, ec2)
	ec4, err := managedTypesContext.GetEllipticCurve(int32(3))
	require.Nil(t, ec4)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)

	ecIndex4 := managedTypesContext.PutEllipticCurve(p521ec)
	require.Equal(t, int32(3), ecIndex4)
	ec4, err = managedTypesContext.GetEllipticCurve(ecIndex4)
	require.Nil(t, err)
	require.Equal(t, p521ec, ec4)
}

func TestManagedTypesContext_ManagedBuffersFunctionalities(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}
	managedTypesContext, _ := NewManagedTypesContext(host)
	bytes := []byte{2, 234, 64, 255}
	emptyBuffer := make([]byte, 0)

	// Calls for non-existent buffers
	noBufHandle := int32(379)
	byteArray, err := managedTypesContext.GetBytes(noBufHandle)
	require.Nil(t, byteArray)
	require.Equal(t, arwen.ErrNoManagedBufferUnderThisHandle, err)
	newBuf, err := managedTypesContext.DeleteSlice(noBufHandle, 0, 3)
	require.Nil(t, newBuf)
	require.Equal(t, arwen.ErrNoManagedBufferUnderThisHandle, err)
	newBuf, err = managedTypesContext.GetSlice(noBufHandle, -3, 2)
	require.Nil(t, newBuf)
	require.Equal(t, arwen.ErrNoManagedBufferUnderThisHandle, err)
	lengthOfmBuffer := managedTypesContext.GetLength(noBufHandle)
	require.Equal(t, int32(-1), lengthOfmBuffer)
	isSuccess := managedTypesContext.AppendBytes(noBufHandle, bytes)
	require.False(t, isSuccess)
	newBuf, err = managedTypesContext.InsertSlice(noBufHandle, 0, bytes)
	require.Nil(t, newBuf)
	require.Equal(t, arwen.ErrNoManagedBufferUnderThisHandle, err)

	// New/Get/Set Buffer
	mBufferHandle1 := managedTypesContext.NewManagedBuffer()
	require.Equal(t, int32(0), mBufferHandle1)
	byteArray, err = managedTypesContext.GetBytes(mBufferHandle1)
	require.Nil(t, err)
	require.Equal(t, emptyBuffer, byteArray)
	managedTypesContext.SetBytes(mBufferHandle1, bytes)
	mBufferBytes, _ := managedTypesContext.GetBytes(mBufferHandle1)
	require.Equal(t, bytes, mBufferBytes)
	mBufferHandle2 := managedTypesContext.NewManagedBufferFromBytes(bytes)
	require.Equal(t, int32(1), mBufferHandle2)
	mBufferBytes, _ = managedTypesContext.GetBytes(mBufferHandle2)
	require.Equal(t, bytes, mBufferBytes)

	// Get Slice
	bufSlice, err := managedTypesContext.GetSlice(noBufHandle, int32(3), int32(0))
	require.Nil(t, bufSlice)
	require.Equal(t, arwen.ErrNoManagedBufferUnderThisHandle, err)
	bufSlice, err = managedTypesContext.GetSlice(mBufferHandle1, int32(1), int32(10))
	require.Nil(t, bufSlice)
	require.Equal(t, arwen.ErrBadBounds, err)
	bufSlice, err = managedTypesContext.GetSlice(mBufferHandle1, int32(4), int32(-1))
	require.Nil(t, bufSlice)
	require.Equal(t, arwen.ErrBadBounds, err)
	bufSlice, err = managedTypesContext.GetSlice(mBufferHandle1, int32(3), int32(0))
	require.Nil(t, err)
	require.Equal(t, emptyBuffer, bufSlice)

	// Delete Slice
	newBuf, err = managedTypesContext.DeleteSlice(mBufferHandle1, 3, 1)
	require.Nil(t, err)
	require.Equal(t, bytes[:3], newBuf)
	newBuf, err = managedTypesContext.DeleteSlice(mBufferHandle1, 3, 0)
	require.Nil(t, err)
	require.Equal(t, bytes[:3], newBuf)
	newBuf, err = managedTypesContext.DeleteSlice(mBufferHandle1, -1, 0)
	require.Nil(t, newBuf)
	require.Equal(t, arwen.ErrBadBounds, err)
	newBuf, err = managedTypesContext.DeleteSlice(mBufferHandle1, 0, -1)
	require.Nil(t, newBuf)
	require.Equal(t, arwen.ErrBadBounds, err)
	newBuf, err = managedTypesContext.DeleteSlice(mBufferHandle1, 0, 10)
	require.Nil(t, err)
	require.Equal(t, emptyBuffer, newBuf)

	// Append, GetLength
	isSuccess = managedTypesContext.AppendBytes(mBufferHandle1, bytes)
	require.True(t, isSuccess)
	lengthOfmBuffer = managedTypesContext.GetLength(mBufferHandle1)
	require.Equal(t, int32(4), lengthOfmBuffer)
	isSuccess = managedTypesContext.AppendBytes(mBufferHandle1, bytes)
	require.True(t, isSuccess)
	mBufferBytes, _ = managedTypesContext.GetBytes(mBufferHandle1)
	require.Equal(t, append(bytes, bytes...), mBufferBytes)
	isSuccess = managedTypesContext.AppendBytes(mBufferHandle1, emptyBuffer)
	require.True(t, isSuccess)
	mBufferBytes, _ = managedTypesContext.GetBytes(mBufferHandle1)
	require.Equal(t, append(bytes, bytes...), mBufferBytes)

	managedTypesContext.SetBytes(mBufferHandle1, bytes)
	mBufferBytes, _ = managedTypesContext.GetBytes(mBufferHandle1)
	require.Equal(t, bytes, mBufferBytes)

	// Insert Slice
	newBuf, err = managedTypesContext.InsertSlice(mBufferHandle1, -1, bytes)
	require.Nil(t, newBuf)
	require.Equal(t, arwen.ErrBadBounds, err)
	newBuf, err = managedTypesContext.InsertSlice(mBufferHandle1, 4, bytes)
	require.Nil(t, newBuf)
	require.Equal(t, arwen.ErrBadBounds, err)
	bytesWithNewSlice := []byte{2, 234, 64, 2, 234, 64, 255, 255}
	newBuf, err = managedTypesContext.InsertSlice(mBufferHandle1, 3, bytes)
	require.Nil(t, err)
	require.Equal(t, bytesWithNewSlice, newBuf)
	bytesWithNewSlice = []byte{2, 234, 64, 255, 2, 234, 64, 2, 234, 64, 255, 255}
	newBuf, err = managedTypesContext.InsertSlice(mBufferHandle1, 0, bytes)
	require.Nil(t, err)
	require.Equal(t, bytesWithNewSlice, newBuf)

	mBufferBytes, _ = managedTypesContext.GetBytes(mBufferHandle1)
	require.Equal(t, bytesWithNewSlice, mBufferBytes)
}

func TestManagedTypesContext_PopSetActiveStateIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}

	managedTypesContext, _ := NewManagedTypesContext(host)
	managedTypesContext.PopSetActiveState()

	require.Equal(t, 0, len(managedTypesContext.managedTypesStack))
}

func TestManagedTypesContext_PopDiscardIfStackIsEmptyShouldNotPanic(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}

	managedTypesContext, _ := NewManagedTypesContext(host)
	managedTypesContext.PopDiscard()

	require.Equal(t, 0, len(managedTypesContext.managedTypesStack))
}

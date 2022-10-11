package contexts

import (
	"bytes"
	"crypto/elliptic"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen"
	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen/mock"
	contextmock "github.com/ElrondNetwork/wasm-vm-v1_4/mock/context"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManagedTypes(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostStub{}

	managedTypesContext, err := NewManagedTypesContext(host)
	currentStateValues := managedTypesContext.managedTypesValues

	require.Nil(t, err)
	require.False(t, managedTypesContext.IsInterfaceNil())
	require.NotNil(t, currentStateValues.bigIntValues)
	require.NotNil(t, currentStateValues.bigFloatValues)
	require.NotNil(t, currentStateValues.ecValues)
	require.NotNil(t, currentStateValues.mBufferValues)
	require.NotNil(t, managedTypesContext.managedTypesStack)
	require.Equal(t, 0, len(currentStateValues.bigIntValues))
	require.Equal(t, 0, len(currentStateValues.bigFloatValues))
	require.Equal(t, 0, len(currentStateValues.ecValues))
	require.Equal(t, 0, len(currentStateValues.mBufferValues))
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
	require.Nil(t, managedTypesContext.randomnessGenerator)
	managedTypesContext.initRandomizer()
	firstRandomizer := managedTypesContext.randomnessGenerator

	managedTypesContextCopy, _ := NewManagedTypesContext(copyHost)
	require.Nil(t, managedTypesContextCopy.randomnessGenerator)
	managedTypesContextCopy.initRandomizer()
	secondRandomizer := managedTypesContextCopy.randomnessGenerator

	require.Equal(t, firstRandomizer, secondRandomizer)

	prg := managedTypesContext.GetRandReader()
	a := make([]byte, 100)
	_, _ = prg.Read(a)
	b := make([]byte, 100)
	for i := 0; i < 1000; i++ {
		_, _ = prg.Read(b)
		require.NotEqual(t, a, b)
		copy(a, b)
	}
}

func TestManagedTypesContext_ClearStateStack(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{
		BlockchainCalled: func() arwen.BlockchainContext {
			return &mock.BlockchainContextMock{}
		},
		RuntimeCalled: func() arwen.RuntimeContext {
			return &contextmock.RuntimeContextMock{CurrentTxHash: bytes.Repeat([]byte{1}, 32)}
		},
	}
	intValue1, intValue2 := int64(100), int64(200)
	floatValue1, floatValue2 := float64(307.72), float64(78.008)
	p224ec, p256ec := elliptic.P224().Params(), elliptic.P256().Params()
	managedTypesContext, _ := NewManagedTypesContext(host)
	managedTypesContext.InitState()

	bigIntHandle1 := managedTypesContext.NewBigIntFromInt64(intValue1)
	bigIntHandle2 := managedTypesContext.NewBigIntFromInt64(intValue2)
	bigFloatHandle1, _ := managedTypesContext.PutBigFloat(new(big.Float).SetFloat64(floatValue1))
	bigFloatHandle2, _ := managedTypesContext.PutBigFloat(new(big.Float).SetFloat64(floatValue2))
	ecHandle1 := managedTypesContext.PutEllipticCurve(p224ec)
	ecHandle2 := managedTypesContext.PutEllipticCurve(p256ec)

	bigFloatHandle3, _ := managedTypesContext.PutBigFloat(nil)
	bigFloat3, _ := managedTypesContext.GetBigFloat(bigFloatHandle3)
	require.Equal(t, big.NewFloat(0), bigFloat3)

	_ = managedTypesContext.GetRandReader()
	assert.False(t, check.IfNil(managedTypesContext.randomnessGenerator))
	managedTypesContext.PushState()
	require.Equal(t, 1, len(managedTypesContext.managedTypesStack))
	managedTypesContext.ClearStateStack()
	require.Equal(t, 0, len(managedTypesContext.managedTypesStack))
	assert.True(t, check.IfNil(managedTypesContext.randomnessGenerator))

	bigInt1, err := managedTypesContext.GetBigInt(bigIntHandle1)
	require.Equal(t, big.NewInt(intValue1), bigInt1)
	require.Nil(t, err)
	bigInt2, err := managedTypesContext.GetBigInt(bigIntHandle2)
	require.Equal(t, big.NewInt(intValue2), bigInt2)
	require.Nil(t, err)
	bigFloat1, err := managedTypesContext.GetBigFloat(bigFloatHandle1)
	require.Equal(t, big.NewFloat(floatValue1), bigFloat1)
	require.Nil(t, err)
	bigFloat2, err := managedTypesContext.GetBigFloat(bigFloatHandle2)
	require.Equal(t, big.NewFloat(floatValue2), bigFloat2)
	require.Nil(t, err)
	ec1, err := managedTypesContext.GetEllipticCurve(ecHandle1)
	require.Nil(t, err)
	require.Equal(t, p224ec, ec1)
	ec2, err := managedTypesContext.GetEllipticCurve(ecHandle2)
	require.Nil(t, err)
	require.Equal(t, p256ec, ec2)

	managedTypesContext.InitState()
	bigInt1, err = managedTypesContext.GetBigInt(bigIntHandle1)
	require.Nil(t, bigInt1)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigInt2, err = managedTypesContext.GetBigInt(bigIntHandle2)
	require.Nil(t, bigInt2)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigFloat1, err = managedTypesContext.GetBigFloat(bigFloatHandle1)
	require.Nil(t, bigFloat1)
	require.Equal(t, arwen.ErrNoBigFloatUnderThisHandle, err)
	bigFloat2, err = managedTypesContext.GetBigFloat(bigFloatHandle2)
	require.Nil(t, bigFloat2)
	require.Equal(t, arwen.ErrNoBigFloatUnderThisHandle, err)
	ec1, err = managedTypesContext.GetEllipticCurve(ecHandle1)
	require.Nil(t, ec1)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)
	ec2, err = managedTypesContext.GetEllipticCurve(ecHandle2)
	require.Nil(t, ec2)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)
}

func TestManagedTypesContext_InitPushPopState(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}
	intValue1, intValue2, intValue3 := int64(100), int64(200), int64(-42)
	floatValue1, floatValue2, floatValue3 := float64(307.72), float64(78.008), float64(-37.84732)
	p224ec, p256ec, p384ec, p521ec := elliptic.P224().Params(), elliptic.P256().Params(), elliptic.P384().Params(), elliptic.P521().Params()
	mBytes := []byte{2, 234, 64, 255}
	managedTypesContext, _ := NewManagedTypesContext(host)
	managedTypesContext.InitState()

	// Create 2 bigInt, 2 bigFloat, 2 EC, 2 managedBuffers on the active state
	bigIntHandle1 := managedTypesContext.NewBigIntFromInt64(intValue1)
	require.Equal(t, int32(0), bigIntHandle1)
	bigIntHandle2 := managedTypesContext.NewBigIntFromInt64(intValue2)
	require.Equal(t, int32(1), bigIntHandle2)

	bigInt1, err := managedTypesContext.GetBigInt(bigIntHandle1)
	require.Equal(t, big.NewInt(intValue1), bigInt1)
	require.Nil(t, err)
	bigInt2, err := managedTypesContext.GetBigInt(bigIntHandle2)
	require.Equal(t, big.NewInt(intValue2), bigInt2)
	require.Nil(t, err)

	bigFloatHandle1, _ := managedTypesContext.PutBigFloat(new(big.Float).SetFloat64(floatValue1))
	require.Equal(t, int32(0), bigFloatHandle1)
	bigFloatHandle2, _ := managedTypesContext.PutBigFloat(new(big.Float).SetFloat64(floatValue2))
	require.Equal(t, int32(1), bigFloatHandle2)

	bigFloat1, err := managedTypesContext.GetBigFloat(bigFloatHandle1)
	require.Equal(t, big.NewFloat(floatValue1), bigFloat1)
	require.Nil(t, err)
	bigFloat2, err := managedTypesContext.GetBigFloat(bigFloatHandle2)
	require.Equal(t, big.NewFloat(floatValue2), bigFloat2)
	require.Nil(t, err)

	ecHandle1 := managedTypesContext.PutEllipticCurve(p224ec)
	require.Equal(t, int32(0), ecHandle1)
	ecHandle2 := managedTypesContext.PutEllipticCurve(p256ec)
	require.Equal(t, int32(1), ecHandle2)

	ec1, err := managedTypesContext.GetEllipticCurve(ecHandle1)
	require.Nil(t, err)
	require.Equal(t, p224ec, ec1)
	ec2, err := managedTypesContext.GetEllipticCurve(ecHandle2)
	require.Nil(t, err)
	require.Equal(t, p256ec, ec2)

	mBufferHandle1 := managedTypesContext.NewManagedBufferFromBytes(mBytes)
	mBuffer, _ := managedTypesContext.GetBytes(mBufferHandle1)
	require.Equal(t, mBytes, mBuffer)

	p224NormalGasCostMultiplier := managedTypesContext.Get100xCurveGasCostMultiplier(ecHandle1)
	require.Equal(t, int32(100), p224NormalGasCostMultiplier)
	p256NormalGasCostMultiplier := managedTypesContext.Get100xCurveGasCostMultiplier(ecHandle2)
	require.Equal(t, int32(135), p256NormalGasCostMultiplier)
	p224ScalarMultGasCostMultiplier := managedTypesContext.GetScalarMult100xCurveGasCostMultiplier(ecHandle1)
	require.Equal(t, int32(100), p224ScalarMultGasCostMultiplier)
	p256ScalarMultGasCostMultiplier := managedTypesContext.GetScalarMult100xCurveGasCostMultiplier(ecHandle2)
	require.Equal(t, int32(110), p256ScalarMultGasCostMultiplier)
	p224UCompressedGasCostMultiplier := managedTypesContext.GetUCompressed100xCurveGasCostMultiplier(ecHandle1)
	require.Equal(t, int32(2000), p224UCompressedGasCostMultiplier)
	p256UCompressedGasCostMultiplier := managedTypesContext.GetUCompressed100xCurveGasCostMultiplier(ecHandle2)
	require.Equal(t, int32(100), p256UCompressedGasCostMultiplier)

	// Copy active state to stack, then clean it. The previous 2 values should not
	// be accessible.
	managedTypesContext.PushState()
	require.Equal(t, 1, len(managedTypesContext.managedTypesStack))
	managedTypesContext.InitState()

	bigInt1, err = managedTypesContext.GetBigInt(bigIntHandle1)
	require.Nil(t, bigInt1)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigInt2, err = managedTypesContext.GetBigInt(bigIntHandle2)
	require.Nil(t, bigInt2)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigInt1, bigInt2, err = managedTypesContext.GetTwoBigInt(bigIntHandle1, bigIntHandle2)
	require.Nil(t, bigInt1, bigInt2)
	require.Equal(t, err, arwen.ErrNoBigIntUnderThisHandle)

	bigFloat1, err = managedTypesContext.GetBigFloat(bigFloatHandle1)
	require.Nil(t, bigFloat1)
	require.Equal(t, arwen.ErrNoBigFloatUnderThisHandle, err)
	bigFloat2, err = managedTypesContext.GetBigFloat(bigFloatHandle2)
	require.Nil(t, bigFloat2)
	require.Equal(t, arwen.ErrNoBigFloatUnderThisHandle, err)
	bigFloat1, bigFloat2, err = managedTypesContext.GetTwoBigFloats(bigFloatHandle1, bigFloatHandle2)
	require.Nil(t, bigFloat1, bigFloat2)
	require.Equal(t, arwen.ErrNoBigFloatUnderThisHandle, err)

	ec1, err = managedTypesContext.GetEllipticCurve(ecHandle1)
	require.Nil(t, ec1)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)
	ec2, err = managedTypesContext.GetEllipticCurve(ecHandle2)
	require.Nil(t, ec2)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)

	mBuffer, err = managedTypesContext.GetBytes(mBufferHandle1)
	require.Nil(t, mBuffer)
	require.Equal(t, arwen.ErrNoManagedBufferUnderThisHandle, err)

	p224NormalGasCostMultiplier = managedTypesContext.Get100xCurveGasCostMultiplier(ecHandle1)
	require.Equal(t, int32(-1), p224NormalGasCostMultiplier)
	p256NormalGasCostMultiplier = managedTypesContext.Get100xCurveGasCostMultiplier(ecHandle2)
	require.Equal(t, int32(-1), p256NormalGasCostMultiplier)
	p224ScalarMultGasCostMultiplier = managedTypesContext.GetScalarMult100xCurveGasCostMultiplier(ecHandle1)
	require.Equal(t, int32(-1), p224ScalarMultGasCostMultiplier)
	p256ScalarMultGasCostMultiplier = managedTypesContext.GetScalarMult100xCurveGasCostMultiplier(ecHandle2)
	require.Equal(t, int32(-1), p256ScalarMultGasCostMultiplier)
	p224UCompressedGasCostMultiplier = managedTypesContext.GetUCompressed100xCurveGasCostMultiplier(ecHandle1)
	require.Equal(t, int32(-1), p224UCompressedGasCostMultiplier)
	p256UCompressedGasCostMultiplier = managedTypesContext.GetUCompressed100xCurveGasCostMultiplier(ecHandle2)
	require.Equal(t, int32(-1), p256UCompressedGasCostMultiplier)

	// Add a value on the current active state
	bigIntHandle3 := managedTypesContext.NewBigIntFromInt64(intValue3)
	require.Equal(t, int32(0), bigIntHandle3)
	bigInt3, err := managedTypesContext.GetBigInt(bigIntHandle3)
	require.Equal(t, big.NewInt(intValue3), bigInt3)
	require.Nil(t, err)

	bigFloatHandle3, _ := managedTypesContext.PutBigFloat(new(big.Float).SetFloat64(floatValue3))
	require.Equal(t, int32(0), bigFloatHandle3)
	bigFloat3, err := managedTypesContext.GetBigFloat(bigFloatHandle3)
	require.Equal(t, big.NewFloat(floatValue3), bigFloat3)
	require.Nil(t, err)

	ecHandle3 := managedTypesContext.PutEllipticCurve(p384ec)
	require.Equal(t, int32(0), ecHandle3)
	ec3, err := managedTypesContext.GetEllipticCurve(ecHandle3)
	require.Nil(t, err)
	require.Equal(t, p384ec, ec3)

	p384NormalGasCostMultiplier := managedTypesContext.Get100xCurveGasCostMultiplier(ecHandle3)
	require.Equal(t, int32(200), p384NormalGasCostMultiplier)
	p384ScalarMultGasCostMultiplier := managedTypesContext.GetScalarMult100xCurveGasCostMultiplier(ecHandle3)
	require.Equal(t, int32(150), p384ScalarMultGasCostMultiplier)
	p384UCompressedGasCostMultiplier := managedTypesContext.GetUCompressed100xCurveGasCostMultiplier(ecHandle3)
	require.Equal(t, int32(200), p384UCompressedGasCostMultiplier)

	// Copy active state to stack, then clean it. The previous 3 values should not
	// be accessible.
	managedTypesContext.PushState()
	require.Equal(t, 2, len(managedTypesContext.managedTypesStack))
	managedTypesContext.InitState()

	bigInt1, err = managedTypesContext.GetBigInt(bigIntHandle1)
	require.Nil(t, bigInt1)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigInt2, err = managedTypesContext.GetBigInt(bigIntHandle2)
	require.Nil(t, bigInt2)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigInt3, err = managedTypesContext.GetBigInt(bigIntHandle3)
	require.Nil(t, bigInt3)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)

	bigFloat1, err = managedTypesContext.GetBigFloat(bigFloatHandle1)
	require.Nil(t, bigFloat1)
	require.Equal(t, arwen.ErrNoBigFloatUnderThisHandle, err)
	bigFloat2, err = managedTypesContext.GetBigFloat(bigFloatHandle2)
	require.Nil(t, bigFloat2)
	require.Equal(t, arwen.ErrNoBigFloatUnderThisHandle, err)
	bigFloat3, err = managedTypesContext.GetBigFloat(bigFloatHandle3)
	require.Nil(t, bigFloat3)
	require.Equal(t, arwen.ErrNoBigFloatUnderThisHandle, err)

	ec1, err = managedTypesContext.GetEllipticCurve(ecHandle1)
	require.Nil(t, ec1)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)
	ec2, err = managedTypesContext.GetEllipticCurve(ecHandle2)
	require.Nil(t, ec2)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)
	ec3, err = managedTypesContext.GetEllipticCurve(ecHandle3)
	require.Nil(t, ec3)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)

	intValue4 := int64(84)
	bigIntHandle4 := managedTypesContext.NewBigIntFromInt64(intValue4)
	require.Equal(t, int32(0), bigIntHandle4)
	bigInt4, err := managedTypesContext.GetBigInt(bigIntHandle4)
	require.Equal(t, big.NewInt(intValue4), bigInt4)
	require.Nil(t, err)
	bigInt4, bigInt3, err = managedTypesContext.GetTwoBigInt(bigIntHandle4, int32(1))
	require.Nil(t, bigInt3)
	require.Nil(t, bigInt4)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)

	floatValue4 := float64(89.3823)
	bigFloatHandle4, _ := managedTypesContext.PutBigFloat(new(big.Float).SetFloat64(floatValue4))
	require.Equal(t, int32(0), bigFloatHandle4)
	bigFloat4, err := managedTypesContext.GetBigFloat(bigFloatHandle4)
	require.Equal(t, big.NewFloat(floatValue4), bigFloat4)
	require.Nil(t, err)
	bigFloat4, bigFloat3, err = managedTypesContext.GetTwoBigFloats(bigFloatHandle4, int32(1))
	require.Nil(t, bigFloat3)
	require.Nil(t, bigFloat4)
	require.Equal(t, arwen.ErrNoBigFloatUnderThisHandle, err)

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
	bigInt4, err = managedTypesContext.GetBigInt(bigIntHandle4)
	require.Equal(t, big.NewInt(intValue4), bigInt4)
	require.Nil(t, err)

	bigFloat4, err = managedTypesContext.GetBigFloat(bigFloatHandle4)
	require.Equal(t, big.NewFloat(floatValue4), bigFloat4)
	require.Nil(t, err)

	ec4, err = managedTypesContext.GetEllipticCurve(ecIndex4)
	require.Equal(t, p521ec, ec4)
	require.Nil(t, err)
	// Restore the first active state by popping to the active state (which is
	// lost).
	managedTypesContext.PopSetActiveState()
	require.Equal(t, 0, len(managedTypesContext.managedTypesStack))

	bigInt1, err = managedTypesContext.GetBigInt(bigIntHandle1)
	require.Equal(t, big.NewInt(intValue1), bigInt1)
	require.Nil(t, err)
	bigInt2, err = managedTypesContext.GetBigInt(bigIntHandle2)
	require.Equal(t, big.NewInt(intValue2), bigInt2)
	require.Nil(t, err)
	bigInt1, bigInt2, err = managedTypesContext.GetTwoBigInt(bigIntHandle1, bigIntHandle2)
	require.Equal(t, big.NewInt(intValue1), bigInt1)
	require.Equal(t, big.NewInt(intValue2), bigInt2)
	require.Nil(t, err)

	bigFloat1, err = managedTypesContext.GetBigFloat(bigFloatHandle1)
	require.Equal(t, big.NewFloat(floatValue1), bigFloat1)
	require.Nil(t, err)
	bigFloat2, err = managedTypesContext.GetBigFloat(bigFloatHandle2)
	require.Equal(t, big.NewFloat(floatValue2), bigFloat2)
	require.Nil(t, err)
	bigFloat1, bigFloat2, err = managedTypesContext.GetTwoBigFloats(bigFloatHandle1, bigFloatHandle2)
	require.Equal(t, big.NewFloat(floatValue1), bigFloat1)
	require.Equal(t, big.NewFloat(floatValue2), bigFloat2)
	require.Nil(t, err)

	ec1, err = managedTypesContext.GetEllipticCurve(ecHandle1)
	require.Nil(t, err)
	require.Equal(t, p224ec, ec1)
	ec2, err = managedTypesContext.GetEllipticCurve(ecHandle2)
	require.Nil(t, err)
	require.Equal(t, p256ec, ec2)
}

func TestManagedTypesContext_PutGetBigInt(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}

	intValue1, intValue2, intValue3, intValue4 := int64(100), int64(200), int64(-42), int64(-80)
	managedTypesContext, _ := NewManagedTypesContext(host)

	bigIntHandle1 := managedTypesContext.NewBigIntFromInt64(intValue1)
	require.Equal(t, int32(0), bigIntHandle1)
	bigIntHandle2 := managedTypesContext.NewBigIntFromInt64(intValue2)
	require.Equal(t, int32(1), bigIntHandle2)
	bigIntHandle3 := managedTypesContext.NewBigIntFromInt64(intValue3)
	require.Equal(t, int32(2), bigIntHandle3)

	bigInt1, err := managedTypesContext.GetBigInt(bigIntHandle1)
	require.Equal(t, big.NewInt(intValue1), bigInt1)
	require.Nil(t, err)
	bigInt2, err := managedTypesContext.GetBigInt(bigIntHandle2)
	require.Equal(t, big.NewInt(intValue2), bigInt2)
	require.Nil(t, err)
	bigInt4, err := managedTypesContext.GetBigInt(int32(3))
	require.Nil(t, bigInt4)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigInt4 = managedTypesContext.GetBigIntOrCreate(3)
	require.Equal(t, big.NewInt(0), bigInt4)

	index4 := managedTypesContext.NewBigIntFromInt64(intValue4)
	require.Equal(t, int32(4), index4)
	bigInt4 = managedTypesContext.GetBigIntOrCreate(4)
	require.Equal(t, big.NewInt(intValue4), bigInt4)

	bigValue, err := managedTypesContext.GetBigInt(123)
	require.Nil(t, bigValue)
	require.Equal(t, arwen.ErrNoBigIntUnderThisHandle, err)
	bigInt1, err = managedTypesContext.GetBigInt(bigIntHandle1)
	require.Equal(t, big.NewInt(intValue1), bigInt1)
	require.Nil(t, err)
	bigInt2, err = managedTypesContext.GetBigInt(bigIntHandle2)
	require.Equal(t, big.NewInt(intValue2), bigInt2)
	require.Nil(t, err)

	bigInt1, err = managedTypesContext.GetBigInt(bigIntHandle1)
	require.Equal(t, big.NewInt(intValue1), bigInt1)
	require.Nil(t, err)
	bigInt2, err = managedTypesContext.GetBigInt(bigIntHandle2)
	require.Equal(t, big.NewInt(intValue2), bigInt2)
	require.Nil(t, err)
	bigInt3, err := managedTypesContext.GetBigInt(bigIntHandle3)
	require.Equal(t, big.NewInt(intValue3), bigInt3)
	require.Nil(t, err)
}

func TestManagedTypesContext_PutGetBigFloat(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}

	floatValue1, floatValue2, floatValue3, floatValue4 := float64(23.56), float64(62.8453), float64(-8234.6512), float64(-0.0001)
	managedTypesContext, _ := NewManagedTypesContext(host)

	bigFloatHandle1, _ := managedTypesContext.PutBigFloat(new(big.Float).SetFloat64(floatValue1))
	require.Equal(t, int32(0), bigFloatHandle1)
	bigFloatHandle2, _ := managedTypesContext.PutBigFloat(new(big.Float).SetFloat64(floatValue2))
	require.Equal(t, int32(1), bigFloatHandle2)
	bigFloatHandle3, _ := managedTypesContext.PutBigFloat(new(big.Float).SetFloat64(floatValue3))
	require.Equal(t, int32(2), bigFloatHandle3)

	bigFloat1, err := managedTypesContext.GetBigFloat(bigFloatHandle1)
	require.Equal(t, big.NewFloat(floatValue1), bigFloat1)
	require.Nil(t, err)
	bigFloat2, err := managedTypesContext.GetBigFloat(bigFloatHandle2)
	require.Equal(t, big.NewFloat(floatValue2), bigFloat2)
	require.Nil(t, err)
	bigFloat4, err := managedTypesContext.GetBigFloat(int32(3))
	require.Nil(t, bigFloat4)
	require.Equal(t, arwen.ErrNoBigFloatUnderThisHandle, err)
	bigFloat4, err = managedTypesContext.GetBigFloatOrCreate(3)
	require.Equal(t, big.NewFloat(0), bigFloat4)
	require.Nil(t, err)

	bigFloatHandle4, _ := managedTypesContext.PutBigFloat(new(big.Float).SetFloat64(floatValue4))
	require.Equal(t, int32(4), bigFloatHandle4)
	bigFloat4, err = managedTypesContext.GetBigFloatOrCreate(4)
	require.Equal(t, big.NewFloat(floatValue4), bigFloat4)
	require.Nil(t, err)

	bigValue, err := managedTypesContext.GetBigFloat(123)
	require.Nil(t, bigValue)
	require.Equal(t, arwen.ErrNoBigFloatUnderThisHandle, err)
	bigFloat1, err = managedTypesContext.GetBigFloat(bigFloatHandle1)
	require.Equal(t, big.NewFloat(floatValue1), bigFloat1)
	require.Nil(t, err)
	bigFloat2, err = managedTypesContext.GetBigFloat(bigFloatHandle2)
	require.Equal(t, big.NewFloat(floatValue2), bigFloat2)
	require.Nil(t, err)

	bigFloat1, err = managedTypesContext.GetBigFloat(bigFloatHandle1)
	require.Equal(t, big.NewFloat(floatValue1), bigFloat1)
	require.Nil(t, err)
	bigFloat2, err = managedTypesContext.GetBigFloat(bigFloatHandle2)
	require.Equal(t, big.NewFloat(floatValue2), bigFloat2)
	require.Nil(t, err)
	bigFloat3, err := managedTypesContext.GetBigFloat(bigFloatHandle3)
	require.Equal(t, big.NewFloat(floatValue3), bigFloat3)
	require.Nil(t, err)

	bigFloat1.SetInf(true)
	bigFloat2.SetInf(false)

	infFloat1, err := managedTypesContext.GetBigFloat(bigFloatHandle1)
	require.Nil(t, infFloat1)
	require.Equal(t, arwen.ErrInfinityFloatOperation, err)

	infFloat2, err := managedTypesContext.GetBigFloat(bigFloatHandle2)
	require.Nil(t, infFloat2)
	require.Equal(t, arwen.ErrInfinityFloatOperation, err)

	infFloat1, err = managedTypesContext.GetBigFloatOrCreate(bigFloatHandle1)
	require.Nil(t, infFloat1)
	require.Equal(t, arwen.ErrInfinityFloatOperation, err)

	infFloat2, err = managedTypesContext.GetBigFloatOrCreate(bigFloatHandle2)
	require.Nil(t, infFloat2)
	require.Equal(t, arwen.ErrInfinityFloatOperation, err)

	infFloat1, infFloat2, err = managedTypesContext.GetTwoBigFloats(bigFloatHandle1, bigFloatHandle2)
	require.Nil(t, infFloat1)
	require.Nil(t, infFloat2)
	require.Equal(t, arwen.ErrInfinityFloatOperation, err)

	infFloat1, nonInfFloat, err := managedTypesContext.GetTwoBigFloats(bigFloatHandle1, bigFloatHandle3)
	require.Nil(t, infFloat1)
	require.Nil(t, nonInfFloat)
	require.Equal(t, arwen.ErrInfinityFloatOperation, err)
}
func TestManagedTypesContext_NewBigIntCopied(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}
	managedTypesContext, _ := NewManagedTypesContext(host)

	originalBigInt := big.NewInt(3)
	index1 := managedTypesContext.NewBigInt(originalBigInt)

	retrievedValue, err := managedTypesContext.GetBigInt(index1)
	require.Nil(t, err)
	retrievedValue.Add(retrievedValue, big.NewInt(100)) // simulate a change of the value in the contract

	require.Equal(t, big.NewInt(3), originalBigInt)
}

func TestManagedTypesContext_PutGetEllipticCurves(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}

	p224ec, p256ec, p384ec, p521ec := elliptic.P224().Params(), elliptic.P256().Params(), elliptic.P384().Params(), elliptic.P521().Params()
	managedTypesContext, _ := NewManagedTypesContext(host)

	ecHandle1 := managedTypesContext.PutEllipticCurve(p224ec)
	require.Equal(t, int32(0), ecHandle1)
	ecHandle2 := managedTypesContext.PutEllipticCurve(p256ec)
	require.Equal(t, int32(1), ecHandle2)
	ecHandle3 := managedTypesContext.PutEllipticCurve(p384ec)
	require.Equal(t, int32(2), ecHandle3)

	p224PrivKeyByteLength := managedTypesContext.GetPrivateKeyByteLengthEC(ecHandle1)
	require.Equal(t, int32(28), p224PrivKeyByteLength)
	p256PrivKeyByteLength := managedTypesContext.GetPrivateKeyByteLengthEC(ecHandle2)
	require.Equal(t, int32(32), p256PrivKeyByteLength)
	p384PrivKeyByteLength := managedTypesContext.GetPrivateKeyByteLengthEC(ecHandle3)
	require.Equal(t, int32(48), p384PrivKeyByteLength)
	nonExistentCurvePrivKeyByteLength := managedTypesContext.GetPrivateKeyByteLengthEC(int32(3))
	require.Equal(t, int32(-1), nonExistentCurvePrivKeyByteLength)

	ec1, err := managedTypesContext.GetEllipticCurve(ecHandle1)
	require.Nil(t, err)
	require.Equal(t, p224ec, ec1)
	ec2, err := managedTypesContext.GetEllipticCurve(ecHandle2)
	require.Nil(t, err)
	require.Equal(t, p256ec, ec2)
	ec4, err := managedTypesContext.GetEllipticCurve(int32(3))
	require.Nil(t, ec4)
	require.Equal(t, arwen.ErrNoEllipticCurveUnderThisHandle, err)

	ecHandle4 := managedTypesContext.PutEllipticCurve(p521ec)
	require.Equal(t, int32(3), ecHandle4)
	ec4, err = managedTypesContext.GetEllipticCurve(ecHandle4)
	require.Nil(t, err)
	require.Equal(t, p521ec, ec4)
}

func TestManagedTypesContext_ManagedBuffersFunctionalities(t *testing.T) {
	t.Parallel()
	host := &contextmock.VMHostStub{}
	managedTypesContext, _ := NewManagedTypesContext(host)
	mBytes := []byte{2, 234, 64, 255}
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
	isSuccess := managedTypesContext.AppendBytes(noBufHandle, mBytes)
	require.False(t, isSuccess)
	newBuf, err = managedTypesContext.InsertSlice(noBufHandle, 0, mBytes)
	require.Nil(t, newBuf)
	require.Equal(t, arwen.ErrNoManagedBufferUnderThisHandle, err)

	// New/Get/Set Buffer
	mBufferHandle1 := managedTypesContext.NewManagedBuffer()
	require.Equal(t, int32(0), mBufferHandle1)
	byteArray, err = managedTypesContext.GetBytes(mBufferHandle1)
	require.Nil(t, err)
	require.Equal(t, emptyBuffer, byteArray)
	managedTypesContext.SetBytes(mBufferHandle1, mBytes)
	mBufferBytes, _ := managedTypesContext.GetBytes(mBufferHandle1)
	require.Equal(t, mBytes, mBufferBytes)
	mBufferHandle2 := managedTypesContext.NewManagedBufferFromBytes(mBytes)
	require.Equal(t, int32(1), mBufferHandle2)
	mBufferBytes, _ = managedTypesContext.GetBytes(mBufferHandle2)
	require.Equal(t, mBytes, mBufferBytes)

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
	require.Equal(t, mBytes[:3], newBuf)
	newBuf, err = managedTypesContext.DeleteSlice(mBufferHandle1, 3, 0)
	require.Nil(t, err)
	require.Equal(t, mBytes[:3], newBuf)
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
	isSuccess = managedTypesContext.AppendBytes(mBufferHandle1, mBytes)
	require.True(t, isSuccess)
	lengthOfmBuffer = managedTypesContext.GetLength(mBufferHandle1)
	require.Equal(t, int32(4), lengthOfmBuffer)
	isSuccess = managedTypesContext.AppendBytes(mBufferHandle1, mBytes)
	require.True(t, isSuccess)
	mBufferBytes, _ = managedTypesContext.GetBytes(mBufferHandle1)
	require.Equal(t, append(mBytes, mBytes...), mBufferBytes)
	isSuccess = managedTypesContext.AppendBytes(mBufferHandle1, emptyBuffer)
	require.True(t, isSuccess)
	mBufferBytes, _ = managedTypesContext.GetBytes(mBufferHandle1)
	require.Equal(t, append(mBytes, mBytes...), mBufferBytes)

	managedTypesContext.SetBytes(mBufferHandle1, mBytes)
	mBufferBytes, _ = managedTypesContext.GetBytes(mBufferHandle1)
	require.Equal(t, mBytes, mBufferBytes)

	// Insert Slice
	newBuf, err = managedTypesContext.InsertSlice(mBufferHandle1, -1, mBytes)
	require.Nil(t, newBuf)
	require.Equal(t, arwen.ErrBadBounds, err)
	newBuf, err = managedTypesContext.InsertSlice(mBufferHandle1, 4, mBytes)
	require.Nil(t, newBuf)
	require.Equal(t, arwen.ErrBadBounds, err)
	bytesWithNewSlice := []byte{2, 234, 64, 2, 234, 64, 255, 255}
	newBuf, err = managedTypesContext.InsertSlice(mBufferHandle1, 3, mBytes)
	require.Nil(t, err)
	require.Equal(t, bytesWithNewSlice, newBuf)
	bytesWithNewSlice = []byte{2, 234, 64, 255, 2, 234, 64, 2, 234, 64, 255, 255}
	newBuf, err = managedTypesContext.InsertSlice(mBufferHandle1, 0, mBytes)
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

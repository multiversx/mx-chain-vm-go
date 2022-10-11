package contexts

import (
	"bytes"
	"crypto/elliptic"
	"encoding/binary"
	"errors"
	"io"
	basicMath "math"
	"math/big"

	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen"
	"github.com/ElrondNetwork/wasm-vm-v1_4/math"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
)

const bigFloatPrecision = 53
const encodedBigFloatMaxByteLen = 18
const bigFloatMaxExponent = 65025
const bigFloatMinExponent = -65025

var positiveEncodedBigFloatPrefix = [...]byte{1, 10, 0, 0, 0, 53}
var negativeEncodedBigFloatPrefix = [...]byte{1, 11, 0, 0, 0, 53}
var encodedZeroBigFloat = [...]byte{1, 8, 0, 0, 0, 53}

const maxBigIntByteLenForNormalCost = 32
const p224CurveMultiplier = 100
const p256CurveMultiplier = 135
const p384CurveMultiplier = 200
const p521CurveMultiplier = 250

const p224CurveScalarMultMultiplier = 100
const p256CurveScalarMultMultiplier = 110
const p384CurveScalarMultMultiplier = 150
const p521CurveScalarMultMultiplier = 190

const p224CurveUnmarshalCompressedMultiplier = 2000
const p256CurveUnmarshalCompressedMultiplier = 100
const p384CurveUnmarshalCompressedMultiplier = 200
const p521CurveUnmarshalCompressedMultiplier = 400

const minEncodedBigFloatLength = 6
const handleLen = 4

type managedBufferMap map[int32][]byte
type bigIntMap map[int32]*big.Int
type bigFloatMap map[int32]*big.Float
type ellipticCurveMap map[int32]*elliptic.CurveParams

type managedTypesContext struct {
	host                arwen.VMHost
	managedTypesValues  managedTypesState
	managedTypesStack   []managedTypesState
	randomnessGenerator math.RandomnessGenerator
}

type managedTypesState struct {
	bigIntValues   bigIntMap
	bigFloatValues bigFloatMap
	ecValues       ellipticCurveMap
	mBufferValues  managedBufferMap
}

// NewManagedTypesContext creates a new managedTypesContext
func NewManagedTypesContext(host arwen.VMHost) (*managedTypesContext, error) {
	context := &managedTypesContext{
		host: host,
		managedTypesValues: managedTypesState{
			bigIntValues:   make(bigIntMap),
			bigFloatValues: make(bigFloatMap),
			ecValues:       make(ellipticCurveMap),
			mBufferValues:  make(managedBufferMap),
		},
		managedTypesStack:   make([]managedTypesState, 0),
		randomnessGenerator: nil,
	}

	return context, nil
}

func (context *managedTypesContext) initRandomizer() {
	blockchainContext := context.host.Blockchain()
	previousRandomSeed := blockchainContext.LastRandomSeed()
	currentRandomSeed := blockchainContext.CurrentRandomSeed()
	txHash := context.host.Runtime().GetCurrentTxHash()

	blocksRandomSeed := append(previousRandomSeed, currentRandomSeed...)
	randomSeed := append(blocksRandomSeed, txHash...)

	randomizer := math.NewSeedRandReader(randomSeed)
	context.randomnessGenerator = randomizer
}

// GetRandReader returns pseudo-randomness generator that implements io.Reader interface
func (context *managedTypesContext) GetRandReader() io.Reader {
	if check.IfNil(context.randomnessGenerator) {
		context.initRandomizer()
	}
	return context.randomnessGenerator
}

// InitState initializes the underlying values map
func (context *managedTypesContext) InitState() {
	context.managedTypesValues = managedTypesState{
		bigIntValues:   make(bigIntMap),
		bigFloatValues: make(bigFloatMap),
		ecValues:       make(ellipticCurveMap),
		mBufferValues:  make(managedBufferMap)}
}

// PushState appends the values map to the state stack
func (context *managedTypesContext) PushState() {
	newBigIntState, newBigFloatState, newEcState, newmBufferState := context.clone()
	context.managedTypesStack = append(context.managedTypesStack, managedTypesState{
		bigIntValues:   newBigIntState,
		bigFloatValues: newBigFloatState,
		ecValues:       newEcState,
		mBufferValues:  newmBufferState,
	})
}

// PopSetActiveState removes the latest entry from the state stack and sets it as the current values map
func (context *managedTypesContext) PopSetActiveState() {
	managedTypesStackLen := len(context.managedTypesStack)
	if managedTypesStackLen == 0 {
		return
	}

	prevState := context.managedTypesStack[managedTypesStackLen-1]
	prevBigIntValues := prevState.bigIntValues
	prevBigFloatValues := prevState.bigFloatValues
	prevEcValues := prevState.ecValues
	prevmBufferValues := prevState.mBufferValues
	context.managedTypesValues.bigIntValues = prevBigIntValues
	context.managedTypesValues.bigFloatValues = prevBigFloatValues
	context.managedTypesValues.ecValues = prevEcValues
	context.managedTypesValues.mBufferValues = prevmBufferValues
	context.managedTypesStack = context.managedTypesStack[:managedTypesStackLen-1]
}

// PopDiscard removes the latest entry from the state stack
func (context *managedTypesContext) PopDiscard() {
	managedTypesStackLen := len(context.managedTypesStack)
	if managedTypesStackLen == 0 {
		return
	}
	context.managedTypesStack = context.managedTypesStack[:managedTypesStackLen-1]
}

// ClearStateStack initializes the state stack
func (context *managedTypesContext) ClearStateStack() {
	context.managedTypesStack = make([]managedTypesState, 0)
	context.randomnessGenerator = nil
}

func (context *managedTypesContext) clone() (bigIntMap, bigFloatMap, ellipticCurveMap, managedBufferMap) {
	newBigIntState := make(bigIntMap, len(context.managedTypesValues.bigIntValues))
	newBigFloatState := make(bigFloatMap, len(context.managedTypesValues.bigFloatValues))
	newEcState := make(ellipticCurveMap, len(context.managedTypesValues.ecValues))
	newmBufferState := make(managedBufferMap, len(context.managedTypesValues.mBufferValues))
	for bigIntHandle, bigInt := range context.managedTypesValues.bigIntValues {
		newBigIntState[bigIntHandle] = big.NewInt(0).Set(bigInt)
	}
	for bigFloatHandle, bigFloat := range context.managedTypesValues.bigFloatValues {
		newBigFloatState[bigFloatHandle] = big.NewFloat(0).Set(bigFloat)
	}
	for ecHandle, ec := range context.managedTypesValues.ecValues {
		newEcState[ecHandle] = ec
	}
	for mBufferHandle, mBuffer := range context.managedTypesValues.mBufferValues {
		newmBufferState[mBufferHandle] = mBuffer
	}
	return newBigIntState, newBigFloatState, newEcState, newmBufferState
}

// IsInterfaceNil returns true if there is no value under the interface
func (context *managedTypesContext) IsInterfaceNil() bool {
	return context == nil
}

// ConsumeGasForBigIntCopy uses gas for Copy operations
func (context *managedTypesContext) ConsumeGasForBigIntCopy(values ...*big.Int) {
	for _, val := range values {
		byteLen := val.BitLen() / 8
		context.ConsumeGasForThisIntNumberOfBytes(byteLen)
	}
}

// ConsumeGasForThisIntNumberOfBytes uses gas for the number of bytes given
func (context *managedTypesContext) ConsumeGasForThisIntNumberOfBytes(byteLen int) {
	gasToUse := uint64(0)
	metering := context.host.Metering()
	if byteLen > maxBigIntByteLenForNormalCost {
		gasToUse = math.MulUint64(uint64(byteLen), metering.GasSchedule().BigIntAPICost.CopyPerByteForTooBig)
		metering.UseAndTraceGas(gasToUse)
	}
}

// ConsumeGasForBytes uses gas for the given bytes
func (context *managedTypesContext) ConsumeGasForBytes(bytes []byte) {
	metering := context.host.Metering()
	gasToUse := math.MulUint64(uint64(len(bytes)), metering.GasSchedule().BaseOperationCost.DataCopyPerByte)
	metering.UseAndTraceGas(gasToUse)
}

// ConsumeGasForThisBigIntNumberOfBytes uses gas for the number of bytes given that are being copied
func (context *managedTypesContext) ConsumeGasForThisBigIntNumberOfBytes(byteLen *big.Int) {
	metering := context.host.Metering()

	gasToUseBigInt := big.NewInt(0).Mul(byteLen, big.NewInt(int64(metering.GasSchedule().BigIntAPICost.CopyPerByteForTooBig)))
	maxGasBigInt := big.NewInt(0).SetUint64(basicMath.MaxUint64)
	gasToUse := uint64(basicMath.MaxUint64)
	if gasToUseBigInt.Cmp(maxGasBigInt) < 0 {
		gasToUse = gasToUseBigInt.Uint64()
	}
	metering.UseAndTraceGas(gasToUse)
}

// ConsumeGasForBigFloatCopy uses gas for the given big float values
func (context *managedTypesContext) ConsumeGasForBigFloatCopy(values ...*big.Float) {
	context.ConsumeGasForThisIntNumberOfBytes(encodedBigFloatMaxByteLen * len(values))
}

// BIGINT

// GetBigIntOrCreate returns the value at the given handle. If there is no value under that value, it will set a new one with value 0
func (context *managedTypesContext) GetBigIntOrCreate(handle int32) *big.Int {
	value, ok := context.managedTypesValues.bigIntValues[handle]
	if !ok {
		value = big.NewInt(0)
		context.managedTypesValues.bigIntValues[handle] = value
	}
	return value
}

// GetBigInt returns the value at the given handle. If there is no value under that handle, it will return error
func (context *managedTypesContext) GetBigInt(handle int32) (*big.Int, error) {
	value, ok := context.managedTypesValues.bigIntValues[handle]
	if !ok {
		return nil, arwen.ErrNoBigIntUnderThisHandle
	}
	return value, nil
}

// GetTwoBigInt returns the values at the two given handles. If there is at least one missing value, it will return error
func (context *managedTypesContext) GetTwoBigInt(handle1 int32, handle2 int32) (*big.Int, *big.Int, error) {
	bigIntValues := context.managedTypesValues.bigIntValues
	value1, ok := bigIntValues[handle1]
	if !ok {
		return nil, nil, arwen.ErrNoBigIntUnderThisHandle
	}
	value2, ok := bigIntValues[handle2]
	if !ok {
		return nil, nil, arwen.ErrNoBigIntUnderThisHandle
	}
	return value1, value2, nil
}

func (context *managedTypesContext) newBigIntNoCopy(value *big.Int) int32 {
	newHandle := int32(len(context.managedTypesValues.bigIntValues))
	for {
		if _, ok := context.managedTypesValues.bigIntValues[newHandle]; !ok {
			break
		}
		newHandle++
	}
	context.managedTypesValues.bigIntValues[newHandle] = value
	return newHandle
}

// BIG FLOAT

// BigFloatPrecIsNotValid checks if the precision of a big float is not valid (not equal to 53)
func (context *managedTypesContext) BigFloatPrecIsNotValid(precision uint) bool {
	return precision != bigFloatPrecision
}

// BigFloatExpIsNotValid checks if the exponent of a big float is not valid (smaller than -65025 or bigger than 65025)
func (context *managedTypesContext) BigFloatExpIsNotValid(exponent int) bool {
	return exponent < bigFloatMinExponent || exponent > bigFloatMaxExponent
}

// EncodedBigFloatIsNotValid checks if an encoded big float is not valid
func (context *managedTypesContext) EncodedBigFloatIsNotValid(encodedBigFloat []byte) bool {
	length := len(encodedBigFloat)
	if length < minEncodedBigFloatLength {
		return true
	} else if length == minEncodedBigFloatLength && !bytes.Equal(encodedBigFloat, encodedZeroBigFloat[:]) {
		return true
	}

	return !bytes.Equal(encodedBigFloat[:minEncodedBigFloatLength], positiveEncodedBigFloatPrefix[:]) &&
		!bytes.Equal(encodedBigFloat[:minEncodedBigFloatLength], negativeEncodedBigFloatPrefix[:])
}

// GetBigFloatOrCreate returns the value at the given handle. If there is no value under that value, it will set a new one with value 0
func (context *managedTypesContext) GetBigFloatOrCreate(handle int32) (*big.Float, error) {
	value, ok := context.managedTypesValues.bigFloatValues[handle]
	if !ok {
		value = big.NewFloat(0)
		context.managedTypesValues.bigFloatValues[handle] = value
	}
	if value.IsInf() {
		return nil, arwen.ErrInfinityFloatOperation
	} else {
		exponent := value.MantExp(nil)
		if exponent > bigFloatMaxExponent || exponent < bigFloatMinExponent {
			return nil, arwen.ErrExponentTooBigOrTooSmall
		}
	}
	return value, nil
}

// GetBigFloat returns the value at the given handle. If there is no value under that handle, it will return error
func (context *managedTypesContext) GetBigFloat(handle int32) (*big.Float, error) {
	value, ok := context.managedTypesValues.bigFloatValues[handle]
	if !ok {
		return nil, arwen.ErrNoBigFloatUnderThisHandle
	}
	if value.IsInf() {
		return nil, arwen.ErrInfinityFloatOperation
	} else {
		exponent := value.MantExp(nil)
		if exponent > bigFloatMaxExponent || exponent < bigFloatMinExponent {
			return nil, arwen.ErrExponentTooBigOrTooSmall
		}
	}
	return value, nil
}

// GetTwoBigFloats returns the values at the two given handles. If there is at least one missing value, it will return error
func (context *managedTypesContext) GetTwoBigFloats(handle1 int32, handle2 int32) (*big.Float, *big.Float, error) {
	bigFloatValues := context.managedTypesValues.bigFloatValues
	value1, ok := bigFloatValues[handle1]
	if !ok {
		return nil, nil, arwen.ErrNoBigFloatUnderThisHandle
	}
	value2, ok := bigFloatValues[handle2]
	if !ok {
		return nil, nil, arwen.ErrNoBigFloatUnderThisHandle
	}
	if value1.IsInf() || value2.IsInf() {
		return nil, nil, arwen.ErrInfinityFloatOperation
	} else {
		exponent1 := value1.MantExp(nil)
		exponent2 := value2.MantExp(nil)
		if context.BigFloatExpIsNotValid(exponent1) || context.BigFloatExpIsNotValid(exponent2) {
			return nil, nil, arwen.ErrExponentTooBigOrTooSmall

		}
	}
	return value1, value2, nil
}

// PutBigFloat adds the given value to the current values map and returns the handle. Returns error if exponent is incorrect
func (context *managedTypesContext) PutBigFloat(value *big.Float) (int32, error) {
	if value == nil {
		value = big.NewFloat(0)
	}
	exponent := value.MantExp(nil)
	if exponent > 65025 || exponent < -65025 {
		return 0, arwen.ErrExponentTooBigOrTooSmall
	}
	newHandle := int32(len(context.managedTypesValues.bigFloatValues))
	for {
		if _, ok := context.managedTypesValues.bigFloatValues[newHandle]; !ok {
			break
		}
		newHandle++
	}

	context.managedTypesValues.bigFloatValues[newHandle] = new(big.Float).Set(value)
	return newHandle, nil
}

// NewBigInt adds the given value to the current values map and returns the handle
func (context *managedTypesContext) NewBigInt(value *big.Int) int32 {
	return context.newBigIntNoCopy(big.NewInt(0).Set(value))
}

// NewBigIntFromInt64 adds the given value to the current values map and returns the handle
func (context *managedTypesContext) NewBigIntFromInt64(int64Value int64) int32 {
	return context.newBigIntNoCopy(big.NewInt(int64Value))
}

// ELLIPTIC CURVES

// GetEllipticCurve returns the elliptic curve under the given handle. If there is no value under that handle, it will return error
func (context *managedTypesContext) GetEllipticCurve(handle int32) (*elliptic.CurveParams, error) {
	curve, ok := context.managedTypesValues.ecValues[handle]
	if !ok {
		return nil, arwen.ErrNoEllipticCurveUnderThisHandle
	}
	return curve, nil
}

// PutEllipticCurve adds the given elliptic curve to the current ecValues map and returns the handle
func (context *managedTypesContext) PutEllipticCurve(curve *elliptic.CurveParams) int32 {
	newHandle := int32(len(context.managedTypesValues.ecValues))
	for {
		if _, ok := context.managedTypesValues.ecValues[newHandle]; !ok {
			break
		}
		newHandle++
	}
	context.managedTypesValues.ecValues[newHandle] = &elliptic.CurveParams{P: curve.P, N: curve.N, B: curve.B, Gx: curve.Gx, Gy: curve.Gy, BitSize: curve.BitSize, Name: curve.Name}
	return newHandle
}

// GetEllipticCurveSizeOfField returns the size of field of the curve under the given handle
func (context *managedTypesContext) GetEllipticCurveSizeOfField(ecHandle int32) int32 {
	curve, ok := context.managedTypesValues.ecValues[ecHandle]
	if !ok {
		return -1
	}
	return int32(curve.BitSize)
}

// Get100xCurveGasCostMultiplier returns (100*multiplier) to be used with the basic gasCost depending on which curve is used
func (context *managedTypesContext) Get100xCurveGasCostMultiplier(ecHandle int32) int32 {
	sizeOfField := context.GetEllipticCurveSizeOfField(ecHandle)
	if sizeOfField < 0 {
		return -1
	}
	switch sizeOfField {
	case 224:
		return p224CurveMultiplier
	case 256:
		return p256CurveMultiplier
	case 384:
		return p384CurveMultiplier
	case 521:
		return p521CurveMultiplier
	}
	return -1
}

// GetScalarMult100xCurveGasCostMultiplier returns (100*multiplier) to be used with the basic gasCost within ScalarMult/ScalarBaseMult function depending on which curve is used
func (context *managedTypesContext) GetScalarMult100xCurveGasCostMultiplier(ecHandle int32) int32 {
	sizeOfField := context.GetEllipticCurveSizeOfField(ecHandle)
	if sizeOfField < 0 {
		return -1
	}
	switch sizeOfField {
	case 224:
		return p224CurveScalarMultMultiplier
	case 256:
		return p256CurveScalarMultMultiplier
	case 384:
		return p384CurveScalarMultMultiplier
	case 521:
		return p521CurveScalarMultMultiplier
	}
	return -1
}

// GetUCompressed100xCurveGasCostMultiplier returns (100*multiplier) to be used with the basic gasCost within UnmarshalCompressed function depending on which curve is used
func (context *managedTypesContext) GetUCompressed100xCurveGasCostMultiplier(ecHandle int32) int32 {
	sizeOfField := context.GetEllipticCurveSizeOfField(ecHandle)
	if sizeOfField < 0 {
		return -1
	}
	switch sizeOfField {
	case 224:
		return p224CurveUnmarshalCompressedMultiplier
	case 256:
		return p256CurveUnmarshalCompressedMultiplier
	case 384:
		return p384CurveUnmarshalCompressedMultiplier
	case 521:
		return p521CurveUnmarshalCompressedMultiplier
	}
	return -1
}

// GetPrivateKeyByteLengthEC returns the length in bytes of the private key that will be generated
func (context *managedTypesContext) GetPrivateKeyByteLengthEC(ecHandle int32) int32 {
	curve, ok := context.managedTypesValues.ecValues[ecHandle]
	if !ok {
		return -1
	}
	return int32((curve.N.BitLen() + 7) / 8)
}

// MANAGED BUFFERS

// NewManagedBuffer creates a new empty buffer in the managed buffers map and returns the handle
func (context *managedTypesContext) NewManagedBuffer() int32 {
	newHandle := int32(len(context.managedTypesValues.mBufferValues))
	for {
		if _, ok := context.managedTypesValues.mBufferValues[newHandle]; !ok {
			break
		}
		newHandle++
	}
	newmBuffer := make([]byte, 0)
	context.managedTypesValues.mBufferValues[newHandle] = newmBuffer
	return newHandle
}

// NewManagedBufferFromBytes creates a new buffer in the managed buffers map, sets the bytes provided, and returns the handle
func (context *managedTypesContext) NewManagedBufferFromBytes(bytes []byte) int32 {
	mBufferHandle := context.NewManagedBuffer()
	context.SetBytes(mBufferHandle, bytes)
	return mBufferHandle
}

// SetBytes sets the bytes given as value for the managed buffer
func (context *managedTypesContext) SetBytes(mBufferHandle int32, bytes []byte) {
	_, ok := context.managedTypesValues.mBufferValues[mBufferHandle]
	if !ok {
		context.managedTypesValues.mBufferValues[mBufferHandle] = make([]byte, 0)
	}

	// always performing a copy,
	// so that changes to the byte buffer in the contract can never leak back into the blockchain
	bytesCopy := make([]byte, len(bytes))
	copy(bytesCopy, bytes)

	context.managedTypesValues.mBufferValues[mBufferHandle] = bytesCopy
}

// GetBytes returns the bytes for the managed buffer. Returns nil as value and error if buffer is non-existent
func (context *managedTypesContext) GetBytes(mBufferHandle int32) ([]byte, error) {
	mBuffer, ok := context.managedTypesValues.mBufferValues[mBufferHandle]
	if !ok {
		return nil, arwen.ErrNoManagedBufferUnderThisHandle
	}
	return mBuffer, nil
}

// AppendBytes appends the given bytes to the buffer at the end
func (context *managedTypesContext) AppendBytes(mBufferHandle int32, bytes []byte) bool {
	_, ok := context.managedTypesValues.mBufferValues[mBufferHandle]
	if !ok {
		return false
	}
	context.managedTypesValues.mBufferValues[mBufferHandle] = append(context.managedTypesValues.mBufferValues[mBufferHandle], bytes...)
	return true
}

// GetLength returns the length of the managed buffer
func (context *managedTypesContext) GetLength(mBufferHandle int32) int32 {
	mBuffer, ok := context.managedTypesValues.mBufferValues[mBufferHandle]
	if !ok {
		return -1
	}
	return int32(len(mBuffer))
}

// GetSlice returns a slice of given length beginning at given start position from the managed buffer
func (context *managedTypesContext) GetSlice(mBufferHandle int32, startPosition int32, lengthOfSlice int32) ([]byte, error) {
	mBuffer, ok := context.managedTypesValues.mBufferValues[mBufferHandle]
	if !ok {
		return nil, arwen.ErrNoManagedBufferUnderThisHandle
	}
	if int(lengthOfSlice) > len(mBuffer)-int(startPosition) || lengthOfSlice < 0 || startPosition < 0 {
		return nil, arwen.ErrBadBounds
	}
	return mBuffer[startPosition:(startPosition + lengthOfSlice)], nil
}

// DeleteSlice deletes a slice from the managed buffer. Returns (new buffer, nil) if success, (nil, error) otherwise
func (context *managedTypesContext) DeleteSlice(mBufferHandle int32, startPosition int32, lengthOfSlice int32) ([]byte, error) {
	mBuffer, ok := context.managedTypesValues.mBufferValues[mBufferHandle]
	if !ok {
		return nil, arwen.ErrNoManagedBufferUnderThisHandle
	}
	if lengthOfSlice < 0 || startPosition < 0 {
		return nil, arwen.ErrBadBounds
	}
	if int(lengthOfSlice) > len(mBuffer)-int(startPosition) {
		mBuffer = mBuffer[:startPosition]
	} else {
		mBuffer = append(mBuffer[:startPosition], mBuffer[startPosition+lengthOfSlice:]...)
	}
	context.managedTypesValues.mBufferValues[mBufferHandle] = mBuffer
	return context.managedTypesValues.mBufferValues[mBufferHandle], nil
}

// InsertSlice inserts a slice in the managed buffer at the given startPosition. Returns (new buffer, nil) if success, (nil, error) otherwise
func (context *managedTypesContext) InsertSlice(mBufferHandle int32, startPosition int32, slice []byte) ([]byte, error) {
	mBuffer, ok := context.managedTypesValues.mBufferValues[mBufferHandle]
	if !ok {
		return nil, arwen.ErrNoManagedBufferUnderThisHandle
	}
	if startPosition < 0 || startPosition > int32(len(mBuffer))-1 {
		return nil, arwen.ErrBadBounds
	}
	mBuffer = append(mBuffer[:startPosition], append(slice, mBuffer[startPosition:]...)...)
	context.managedTypesValues.mBufferValues[mBufferHandle] = mBuffer
	return context.managedTypesValues.mBufferValues[mBufferHandle], nil
}

// ReadManagedVecOfManagedBuffers converts a managed buffer of managed buffers to a slice of byte slices.
func (context *managedTypesContext) ReadManagedVecOfManagedBuffers(
	managedVecHandle int32,
) ([][]byte, uint64, error) {
	managedVecBytes, err := context.GetBytes(managedVecHandle)
	if err != nil {
		return nil, 0, err
	}
	context.ConsumeGasForBytes(managedVecBytes)

	if len(managedVecBytes)%handleLen != 0 {
		return nil, 0, errors.New("invalid managed vector of managed buffer handles")
	}

	numBuffers := len(managedVecBytes) / handleLen
	result := make([][]byte, 0, numBuffers)
	sumOfItemByteLengths := uint64(0)
	for i := 0; i < len(managedVecBytes); i += handleLen {
		itemHandle := int32(binary.BigEndian.Uint32(managedVecBytes[i : i+handleLen]))

		itemBytes, err := context.GetBytes(itemHandle)
		if err != nil {
			return nil, 0, err
		}
		context.ConsumeGasForBytes(itemBytes)

		sumOfItemByteLengths += uint64(len(itemBytes))
		result = append(result, itemBytes)
	}

	return result, sumOfItemByteLengths, nil
}

// WriteManagedVecOfManagedBuffers converts a slice of byte slices to a managed buffer of managed buffers.
func (context *managedTypesContext) WriteManagedVecOfManagedBuffers(
	data [][]byte,
	destinationHandle int32,
) {
	sumOfItemByteLengths := uint64(0)
	destinationBytes := make([]byte, handleLen*len(data))
	dataIndex := 0
	for _, itemBytes := range data {
		sumOfItemByteLengths += uint64(len(itemBytes))
		itemHandle := context.NewManagedBufferFromBytes(itemBytes)
		binary.BigEndian.PutUint32(destinationBytes[dataIndex:dataIndex+handleLen], uint32(itemHandle))
		dataIndex += handleLen
	}

	context.SetBytes(destinationHandle, destinationBytes)
	metering := context.host.Metering()
	metering.UseAndTraceGas(sumOfItemByteLengths * metering.GasSchedule().BaseOperationCost.DataCopyPerByte)
}

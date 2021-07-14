package contexts

import (
	"crypto/elliptic"
	basicMath "math"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/math"
)

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

type managedBufferMap map[int32][]byte
type bigIntMap map[int32]*big.Int
type ellipticCurveMap map[int32]*elliptic.CurveParams

type managedTypesContext struct {
	host               arwen.VMHost
	managedTypesValues managedTypesState
	managedTypesStack  []managedTypesState
}

type managedTypesState struct {
	bigIntValues bigIntMap
	ecValues     ellipticCurveMap
	manBufValues managedBufferMap
}

// NewBigIntContext creates a new bigIntContext
func NewManagedTypesContext(host arwen.VMHost) (*managedTypesContext, error) {
	context := &managedTypesContext{
		host: host,
		managedTypesValues: managedTypesState{
			bigIntValues: make(bigIntMap),
			ecValues:     make(ellipticCurveMap),
			manBufValues: make(managedBufferMap),
		},
		managedTypesStack: make([]managedTypesState, 0),
	}

	return context, nil
}

// InitState initializes the underlying values map
func (context *managedTypesContext) InitState() {
	context.managedTypesValues = managedTypesState{
		bigIntValues: make(bigIntMap),
		ecValues:     make(ellipticCurveMap),
		manBufValues: make(managedBufferMap)}
}

// PushState appends the values map to the state stack
func (context *managedTypesContext) PushState() {
	newBigIntState, newEcState, newManBufState := context.clone()
	context.managedTypesStack = append(context.managedTypesStack, managedTypesState{
		bigIntValues: newBigIntState,
		ecValues:     newEcState,
		manBufValues: newManBufState,
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
	prevEcValues := prevState.ecValues
	prevManBufValues := prevState.manBufValues
	context.managedTypesValues.bigIntValues = prevBigIntValues
	context.managedTypesValues.ecValues = prevEcValues
	context.managedTypesValues.manBufValues = prevManBufValues
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
}

func (context *managedTypesContext) clone() (bigIntMap, ellipticCurveMap, managedBufferMap) {
	newBigIntState := make(bigIntMap, len(context.managedTypesValues.bigIntValues))
	newEcState := make(ellipticCurveMap, len(context.managedTypesValues.ecValues))
	newManBufState := make(managedBufferMap, len(context.managedTypesValues.manBufValues))
	for bigIntHandle, bigInt := range context.managedTypesValues.bigIntValues {
		newBigIntState[bigIntHandle] = big.NewInt(0).Set(bigInt)
	}
	for ecHandle, ec := range context.managedTypesValues.ecValues {
		newEcState[ecHandle] = ec
	}
	for manBufHandle, manBuf := range context.managedTypesValues.manBufValues {
		newManBufState[manBufHandle] = manBuf
	}
	return newBigIntState, newEcState, newManBufState
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
	metering := context.host.Metering()
	if byteLen > maxBigIntByteLenForNormalCost {
		metering.UseGas(math.MulUint64(uint64(byteLen), metering.GasSchedule().BaseOperationCost.DataCopyPerByte))
	}
}

// ConsumeGasForThisBigIntNumberOfBytes uses gas for the number of bytes given that are being copied
func (context *managedTypesContext) ConsumeGasForThisBigIntNumberOfBytes(byteLen *big.Int) {
	metering := context.host.Metering()
	DataCopyPerByte := metering.GasSchedule().BaseOperationCost.DataCopyPerByte

	gasToUseBigInt := big.NewInt(0).Mul(byteLen, big.NewInt(int64(DataCopyPerByte)))
	maxGasBigInt := big.NewInt(0).SetUint64(basicMath.MaxUint64)
	gasToUse := uint64(basicMath.MaxUint64)
	if gasToUseBigInt.Cmp(maxGasBigInt) < 0 {
		gasToUse = gasToUseBigInt.Uint64()
	}
	metering.UseGas(gasToUse)
}

// BIGINT

// GetOneOrCreate returns the value at the given handle. If there is no value under that value, it will set a new one with value 0
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

// PutBigInt adds the given value to the current values map and returns the handle
func (context *managedTypesContext) PutBigInt(value int64) int32 {
	newHandle := int32(len(context.managedTypesValues.bigIntValues))
	for {
		if _, ok := context.managedTypesValues.bigIntValues[newHandle]; !ok {
			break
		}
		newHandle++
	}
	context.managedTypesValues.bigIntValues[newHandle] = big.NewInt(value)
	return newHandle
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
	newHandle := int32(len(context.managedTypesValues.manBufValues))
	for {
		if _, ok := context.managedTypesValues.manBufValues[newHandle]; !ok {
			break
		}
		newHandle++
	}
	newManBuf := make([]byte, 0)
	context.managedTypesValues.manBufValues[newHandle] = newManBuf
	return newHandle
}

// NewManagedBufferFromBytes creates a new buffer in the managed buffers map, sets the bytes provided, and returns the handle
func (context *managedTypesContext) NewManagedBufferFromBytes(bytes []byte) int32 {
	manBufHandle := context.NewManagedBuffer()
	context.SetBytesForThisManagedBuffer(manBufHandle, bytes)
	return manBufHandle
}

// SetBytesForThisManagedBuffer sets the bytes given as value for the managed buffer. Returns 0 if success, 1 otherwise
func (context *managedTypesContext) SetBytesForThisManagedBuffer(manBufHandle int32, bytes []byte) bool {
	_, ok := context.managedTypesValues.manBufValues[manBufHandle]
	if !ok {
		return false
	}
	context.managedTypesValues.manBufValues[manBufHandle] = bytes
	return true
}

// GetBytesForThisManagedBuffer returns the bytes for the managed buffer. Returns nil as value and error if buffer is non-existent
func (context *managedTypesContext) GetBytesForThisManagedBuffer(manBufHandle int32) ([]byte, error) {
	manBuf, ok := context.managedTypesValues.manBufValues[manBufHandle]
	if !ok {
		return nil, arwen.ErrNoManagedBufferUnderThisHandle
	}
	return manBuf, nil
}

// AppendBytesToThisManagedBuffer appends the given bytes to the buffer at the end
func (context *managedTypesContext) AppendBytesToThisManagedBuffer(manBufHandle int32, bytes []byte) bool {
	_, ok := context.managedTypesValues.manBufValues[manBufHandle]
	if !ok {
		return false
	}
	context.managedTypesValues.manBufValues[manBufHandle] = append(context.managedTypesValues.manBufValues[manBufHandle], bytes...)
	return true
}

// GetLengthForThisManagedBuffer returns the length of the managed buffer
func (context *managedTypesContext) GetLengthForThisManagedBuffer(manBufHandle int32) int32 {
	manBuf, ok := context.managedTypesValues.manBufValues[manBufHandle]
	if !ok {
		return -1
	}
	return int32(len(manBuf))
}

// GetSliceFromManagedBuffer returns a slice of given length beginning at given start position from the managed buffer
func (context *managedTypesContext) GetSliceFromManagedBuffer(manBufHandle int32, startPosition int32, lengthOfSlice int32) ([]byte, error) {
	manBuf, ok := context.managedTypesValues.manBufValues[manBufHandle]
	if !ok {
		return nil, arwen.ErrNoManagedBufferUnderThisHandle
	}
	if int(lengthOfSlice) > len(manBuf)-int(startPosition) || lengthOfSlice < 0 || startPosition < 0 {
		return nil, arwen.ErrBadBounds
	}
	return manBuf[startPosition:(startPosition + lengthOfSlice)], nil
}

// DeleteSliceFromManagedBuffer deletes a slice from the managed buffer. Returns (new buffer, nil) if success, (nil, error) otherwise
func (context *managedTypesContext) DeleteSliceFromManagedBuffer(manBufHandle int32, startPosition int32, lengthOfSlice int32) ([]byte, error) {
	manBuf, ok := context.managedTypesValues.manBufValues[manBufHandle]
	if !ok {
		return nil, arwen.ErrNoManagedBufferUnderThisHandle
	}
	if lengthOfSlice < 0 || startPosition < 0 {
		return nil, arwen.ErrBadBounds
	}
	if int(lengthOfSlice) > len(manBuf)-int(startPosition) {
		manBuf = manBuf[:startPosition]
	} else {
		manBuf = append(manBuf[:startPosition], manBuf[startPosition+lengthOfSlice:]...)
	}
	context.managedTypesValues.manBufValues[manBufHandle] = manBuf
	return context.managedTypesValues.manBufValues[manBufHandle], nil
}

// InsertSliceInManagedBuffer inserts a slice in the managed buffer at the given startPosition. Returns (new buffer, nil) if success, (nil, error) otherwise
func (context *managedTypesContext) InsertSliceInManagedBuffer(manBufHandle int32, startPosition int32, slice []byte) ([]byte, error) {
	manBuf, ok := context.managedTypesValues.manBufValues[manBufHandle]
	if !ok {
		return nil, arwen.ErrNoManagedBufferUnderThisHandle
	}
	if startPosition < 0 || startPosition > int32(len(manBuf))-1 {
		return nil, arwen.ErrBadBounds
	}
	manBuf = append(manBuf[:startPosition], append(slice, manBuf[startPosition:]...)...)
	context.managedTypesValues.manBufValues[manBufHandle] = manBuf
	return context.managedTypesValues.manBufValues[manBufHandle], nil
}

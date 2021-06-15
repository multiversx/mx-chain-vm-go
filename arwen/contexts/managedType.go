package contexts

import (
	"crypto/elliptic"
	basicMath "math"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/math"
)

const maxBigIntByteLenForNormalCost = 32

type bigIntMap map[int32]*big.Int
type ellipticCurveMap map[int32]*elliptic.CurveParams

type managedTypesContext struct {
	host             arwen.VMHost
	bigIntValues     bigIntMap
	ecValues         ellipticCurveMap
	ecStateStack     []ellipticCurveMap
	bigIntStateStack []bigIntMap
}

// NewBigIntContext creates a new bigIntContext
func NewManagedTypesContext(host arwen.VMHost) (*managedTypesContext, error) {
	context := &managedTypesContext{
		host:             host,
		bigIntValues:     make(bigIntMap),
		ecValues:         make(ellipticCurveMap),
		ecStateStack:     make([]ellipticCurveMap, 0),
		bigIntStateStack: make([]bigIntMap, 0),
	}

	return context, nil
}

// InitState initializes the underlying values map
func (context *managedTypesContext) InitState() {
	context.bigIntValues = make(bigIntMap)
	context.ecValues = make(ellipticCurveMap)
}

// PushState appends the values map to the state stack
func (context *managedTypesContext) PushState() {
	newBigIntState, newEcState := context.clone()
	context.bigIntStateStack = append(context.bigIntStateStack, newBigIntState)
	context.ecStateStack = append(context.ecStateStack, newEcState)
}

// PopSetActiveState removes the latest entry from the state stack and sets it as the current values map
func (context *managedTypesContext) PopSetActiveState() {
	bigIntStateStackLen := len(context.bigIntStateStack)
	ecStateStackLen := len(context.ecStateStack)
	if bigIntStateStackLen == 0 && ecStateStackLen == 0 {
		return
	}
	prevBigIntValues := context.bigIntStateStack[bigIntStateStackLen-1]
	context.bigIntStateStack = context.bigIntStateStack[:bigIntStateStackLen-1]
	context.bigIntValues = prevBigIntValues

	prevEcValues := context.ecStateStack[ecStateStackLen-1]
	context.ecStateStack = context.ecStateStack[:ecStateStackLen-1]
	context.ecValues = prevEcValues
}

// PopDiscard removes the latest entry from the state stack
func (context *managedTypesContext) PopDiscard() {
	bigIntStateStackLen := len(context.bigIntStateStack)
	ecStateStackLen := len(context.ecStateStack)
	if bigIntStateStackLen == 0 && ecStateStackLen == 0 {
		return
	}

	context.ecStateStack = context.ecStateStack[:ecStateStackLen-1]
	context.bigIntStateStack = context.bigIntStateStack[:bigIntStateStackLen-1]
}

// ClearStateStack initializes the state stack
func (context *managedTypesContext) ClearStateStack() {
	context.bigIntStateStack = make([]bigIntMap, 0)
	context.ecStateStack = make([]ellipticCurveMap, 0)
}

func (context *managedTypesContext) clone() (bigIntMap, ellipticCurveMap) {
	newBigIntState := make(bigIntMap, len(context.bigIntValues))
	newEcState := make(ellipticCurveMap, len(context.ecValues))
	for bigIntHandle, bigInt := range context.bigIntValues {
		newBigIntState[bigIntHandle] = big.NewInt(0).Set(bigInt)
	}
	for ecHandle, ec := range context.ecValues {
		newEcState[ecHandle] = ec
	}
	return newBigIntState, newEcState
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

// ConsumeGasForThisIntNumberOfBytes uses gas for the number of bytes given.
func (context *managedTypesContext) ConsumeGasForThisIntNumberOfBytes(byteLen int) {
	metering := context.host.Metering()
	if byteLen > maxBigIntByteLenForNormalCost {
		metering.UseGas(math.MulUint64(uint64(byteLen), metering.GasSchedule().BaseOperationCost.DataCopyPerByte))
	}
}

// ConsumeGasForThisBigIntNumberOfBytes uses gas for the number of bytes given that are being copied.
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

// GetOneOrCreate returns the value at the given handle. If there is no value under that value, it will set a new on with value 0.
func (context *managedTypesContext) GetBigIntOrCreate(handle int32) *big.Int {
	value, ok := context.bigIntValues[handle]
	if !ok {
		value = big.NewInt(0)
		context.bigIntValues[handle] = value
	}
	return value
}

// GetBigInt returns the value at the given handle. If there is no value under that handle, it will return error
func (context *managedTypesContext) GetBigInt(handle int32) (*big.Int, error) {
	value, ok := context.bigIntValues[handle]
	if !ok {
		return nil, arwen.ErrNoBigIntUnderThisHandle
	}
	return value, nil
}

// GetTwoBigInt returns the values at the two given handles. If there is at least one missing value, it will return error
func (context *managedTypesContext) GetTwoBigInt(handle1 int32, handle2 int32) (*big.Int, *big.Int, error) {
	value1, ok := context.bigIntValues[handle1]
	if !ok {
		return nil, nil, arwen.ErrNoBigIntUnderThisHandle
	}
	value2, ok := context.bigIntValues[handle2]
	if !ok {
		return nil, nil, arwen.ErrNoBigIntUnderThisHandle
	}
	return value1, value2, nil
}

// PutBigInt adds the given value to the current values map and returns the handle
func (context *managedTypesContext) PutBigInt(value int64) int32 {
	newHandle := int32(len(context.bigIntValues))
	for {
		if _, ok := context.bigIntValues[newHandle]; !ok {
			break
		}
		newHandle++
	}
	context.bigIntValues[newHandle] = big.NewInt(value)
	return newHandle
}

// ELLIPTIC CURVES

// GetOneEllipticCurve returns the elliptic curve under the given handle. If there is no value under that handle, it will return error
func (context *managedTypesContext) GetEllipticCurve(handle int32) (*elliptic.CurveParams, error) {
	curve, ok := context.ecValues[handle]
	if !ok {
		return nil, arwen.ErrNoEllipticCurveUnderThisHandle
	}
	return curve, nil
}

// PutEllipticCurve adds the given elliptic curve to the current ecValues map and returns the handle
func (context *managedTypesContext) PutEllipticCurve(curve *elliptic.CurveParams) int32 {
	newHandle := int32(len(context.ecValues))
	for {
		if _, ok := context.ecValues[newHandle]; !ok {
			break
		}
		newHandle++
	}
	context.ecValues[newHandle] = &elliptic.CurveParams{P: curve.P, N: curve.N, B: curve.B, Gx: curve.Gx, Gy: curve.Gy, BitSize: curve.BitSize, Name: curve.Name}
	return newHandle
}

// GetEllipticCurveLength returns the size of field of the curve under the given handle.
func (context *managedTypesContext) GetEllipticCurveSizeOfField(ecHandle int32) int32 {
	curve, ok := context.ecValues[ecHandle]
	if !ok {
		return -1
	}
	return int32(curve.BitSize)
}

// GetPrivateKeyByteLengthEC returns the length in bytes of the private key that will be generated.
func (context *managedTypesContext) GetPrivateKeyByteLengthEC(ecHandle int32) int32 {
	curve, ok := context.ecValues[ecHandle]
	if !ok {
		return -1
	}
	return int32((curve.N.BitLen() + 7) / 8)
}

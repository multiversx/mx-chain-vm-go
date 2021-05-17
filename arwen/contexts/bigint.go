package contexts

import (
	basicMath "math"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/math"
)

const maxBigIntByteLenForNormalCost = 32

type bigIntMap map[int32]*big.Int

type bigIntContext struct {
	host       arwen.VMHost
	values     bigIntMap
	stateStack []bigIntMap
}

// NewBigIntContext creates a new bigIntContext
func NewBigIntContext(host arwen.VMHost) (*bigIntContext, error) {
	context := &bigIntContext{
		host:       host,
		values:     make(bigIntMap),
		stateStack: make([]bigIntMap, 0),
	}

	return context, nil
}

// InitState initializes the underlying values map
func (context *bigIntContext) InitState() {
	context.values = make(bigIntMap)
}

// PushState appends the values map to the state stack
func (context *bigIntContext) PushState() {
	newState := context.clone()
	context.stateStack = append(context.stateStack, newState)
}

// PopSetActiveState removes the latest entry from the state stack and sets it as the current values map
func (context *bigIntContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	prevValues := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.values = prevValues
}

// PopDiscard removes the latest entry from the state stack
func (context *bigIntContext) PopDiscard() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	context.stateStack = context.stateStack[:stateStackLen-1]
}

// ClearStateStack initializes the state stack
func (context *bigIntContext) ClearStateStack() {
	context.stateStack = make([]bigIntMap, 0)
}

func (context *bigIntContext) clone() bigIntMap {
	newState := make(bigIntMap, len(context.values))
	for handle, bigInt := range context.values {
		newState[handle] = big.NewInt(0).Set(bigInt)
	}
	return newState
}

// Put adds the given value to the current values map and returns the handle
func (context *bigIntContext) Put(value int64) int32 {
	newHandle := int32(len(context.values))
	for {
		if _, ok := context.values[newHandle]; !ok {
			break
		}
		newHandle++
	}

	context.values[newHandle] = big.NewInt(value)

	return newHandle
}

// GetOne returns the value at the given handle. If there is no value under that handle, it will return error
func (context *bigIntContext) GetOne(handle int32) (*big.Int, error) {
	value, ok := context.values[handle]

	if !ok {
		return nil, arwen.ErrNoBigIntUnderThisHandle
	}

	return value, nil
}

// GetTwo returns the values at the given handles.
func (context *bigIntContext) GetTwo(handle1 int32, handle2 int32) (*big.Int, *big.Int, error) {
	firstBigInt, err := context.GetOne(handle1)
	if err != nil {
		return nil, nil, err
	}
	secondBigInt, err := context.GetOne(handle2)
	if err != nil {
		return nil, nil, err
	}
	return firstBigInt, secondBigInt, nil
}

// GetThree returns the values at the given handles.
func (context *bigIntContext) GetThree(handle1 int32, handle2 int32, handle3 int32) (*big.Int, *big.Int, *big.Int, error) {
	firstBigInt, secondBigInt, err := context.GetTwo(handle1, handle2)
	if err != nil {
		return nil, nil, nil, err
	}
	thirdBigInt, err := context.GetOne(handle3)
	if err != nil {
		return nil, nil, nil, err
	}
	return firstBigInt, secondBigInt, thirdBigInt, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (context *bigIntContext) IsInterfaceNil() bool {
	return context == nil
}

// ConsumeGasForBigIntCopy uses gas for Copy operations
func (context *bigIntContext) ConsumeGasForBigIntCopy(values ...*big.Int) {
	for _, val := range values {
		byteLen := val.BitLen() / 8
		context.ConsumeGasForThisIntNumberOfBytes(byteLen)
	}
}

// ConsumeGasForThisIntNumberOfBytes uses gas for the number of bytes given.
func (context *bigIntContext) ConsumeGasForThisIntNumberOfBytes(byteLen int) {
	metering := context.host.Metering()
	if byteLen > maxBigIntByteLenForNormalCost {
		metering.UseGas(math.MulUint64(uint64(byteLen), metering.GasSchedule().BaseOperationCost.DataCopyPerByte))
	}
}

// ConsumeGasForThisBigIntNumberOfBytes uses gas for the number of bytes given that are being copied.
func (context *bigIntContext) ConsumeGasForThisBigIntNumberOfBytes(byteLen *big.Int) {
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

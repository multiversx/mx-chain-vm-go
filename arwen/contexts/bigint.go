package contexts

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
)

type bigIntContext struct {
	mappedValues map[int32]*big.Int
	stateStack   []*bigIntContext
	nextIndex    int32
}

func NewBigIntContext() (*bigIntContext, error) {
	context := &bigIntContext{
		mappedValues: make(map[int32]*big.Int),
		stateStack:   make([]*bigIntContext, 0),
		nextIndex:    0,
	}

	return context, nil
}

func (context *bigIntContext) InitState() {
	context.mappedValues = make(map[int32]*big.Int)
}

func (context *bigIntContext) PushState() {
	newState := &bigIntContext{
		mappedValues: context.mappedValues,
	}

	context.stateStack = append(context.stateStack, newState)
}

func (context *bigIntContext) PopState() error {
	stateStackLen := len(context.stateStack)
	if stateStackLen < 1 {
		return arwen.StateStackUnderflow
	}

	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.mappedValues = prevState.mappedValues

	return nil
}

func (context *bigIntContext) putBigInt(bi *big.Int) int32 {
	for {
		if _, ok := context.mappedValues[context.nextIndex]; !ok {
			break
		}
		context.nextIndex++
	}

	newIndex := context.nextIndex
	context.mappedValues[newIndex] = bi
	context.nextIndex++
	return newIndex
}

func (context *bigIntContext) Put(value int64) int32 {
	newIndex := context.putBigInt(big.NewInt(value))
	return newIndex
}

func (context *bigIntContext) Clone(id int32) int32 {
	oldValue := context.GetOne(id)
	newIndex := context.putBigInt(big.NewInt(0).Set(oldValue))
	return newIndex
}

func (context *bigIntContext) Destruct(id int32) {
	delete(context.mappedValues, id)
}

func (context *bigIntContext) GetOne(id int32) *big.Int {
	if _, ok := context.mappedValues[id]; !ok {
		context.mappedValues[id] = big.NewInt(0)
	}

	return context.mappedValues[id]
}

func (context *bigIntContext) GetTwo(id1 int32, id2 int32) (*big.Int, *big.Int) {
	if _, ok := context.mappedValues[id1]; !ok {
		context.mappedValues[id1] = big.NewInt(0)
	}

	if _, ok := context.mappedValues[id2]; !ok {
		context.mappedValues[id2] = big.NewInt(0)
	}

	return context.mappedValues[id1], context.mappedValues[id2]
}

func (context *bigIntContext) GetThree(id1 int32, id2 int32, id3 int32) (*big.Int, *big.Int, *big.Int) {
	if _, ok := context.mappedValues[id1]; !ok {
		context.mappedValues[id1] = big.NewInt(0)
	}

	if _, ok := context.mappedValues[id2]; !ok {
		context.mappedValues[id2] = big.NewInt(0)
	}

	if _, ok := context.mappedValues[id3]; !ok {
		context.mappedValues[id3] = big.NewInt(0)
	}

	return context.mappedValues[id1], context.mappedValues[id2], context.mappedValues[id3]
}

func (context *bigIntContext) IsInterfaceNil() bool {
	return context == nil
}

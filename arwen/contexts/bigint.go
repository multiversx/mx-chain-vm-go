package contexts

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
)

type bigIntContext struct {
	mappedValues map[int32]*big.Int
	stateStack   []*bigIntContext
}

func NewBigIntContext() (*bigIntContext, error) {
	context := &bigIntContext{
		mappedValues: make(map[int32]*big.Int),
		stateStack:   make([]*bigIntContext, 0),
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

func (context *bigIntContext) Put(value int64) int32 {
	newIndex := int32(len(context.mappedValues))
	for {
		if _, ok := context.mappedValues[newIndex]; !ok {
			break
		}
		newIndex++
	}

	context.mappedValues[newIndex] = big.NewInt(value)

	return newIndex
}

func (context *bigIntContext) GetOne(id int32) *big.Int {
	if _, ok := context.mappedValues[id]; !ok {
		context.mappedValues[id] = big.NewInt(0)
	}

	return context.mappedValues[id]
}

func (context *bigIntContext) GetTwo(id1 int32, id2 int32) (*big.Int, *big.Int) {
	return context.GetOne(id1), context.GetOne(id2)
}

func (context *bigIntContext) GetThree(id1 int32, id2 int32, id3 int32) (*big.Int, *big.Int, *big.Int) {
	return context.GetOne(id1), context.GetOne(id2), context.GetOne(id3)
}

func (context *bigIntContext) IsInterfaceNil() bool {
	return context == nil
}

package contexts

import (
	"math/big"
)

type bigIntMap map[int32]*big.Int

type bigIntContext struct {
	values     bigIntMap
	stateStack []bigIntMap
}

// NewBigIntContext creates a new bigIntContext
func NewBigIntContext() (*bigIntContext, error) {
	context := &bigIntContext{
		values:     make(bigIntMap),
		stateStack: make([]bigIntMap, 0),
	}

	return context, nil
}

func (context *bigIntContext) InitState() {
	context.values = make(bigIntMap)
}

func (context *bigIntContext) PushState() {
	newState := context.clone()
	context.stateStack = append(context.stateStack, newState)
}

func (context *bigIntContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	prevValues := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.values = prevValues
}

func (context *bigIntContext) PopDiscard() {
	stateStackLen := len(context.stateStack)
	context.stateStack = context.stateStack[:stateStackLen-1]
}

func (context *bigIntContext) ClearStateStack() {
	context.stateStack = make([]bigIntMap, 0)
}

func (context *bigIntContext) clone() bigIntMap {
	newState := make(bigIntMap, len(context.values))
	for id, bigInt := range context.values {
		newState[id] = big.NewInt(0).Set(bigInt)
	}
	return newState
}

func (context *bigIntContext) Put(value int64) int32 {
	newIndex := int32(len(context.values))
	for {
		if _, ok := context.values[newIndex]; !ok {
			break
		}
		newIndex++
	}

	context.values[newIndex] = big.NewInt(value)

	return newIndex
}

func (context *bigIntContext) GetOne(id int32) *big.Int {
	if _, ok := context.values[id]; !ok {
		context.values[id] = big.NewInt(0)
	}

	return context.values[id]
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

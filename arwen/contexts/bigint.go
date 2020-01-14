package contexts

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
)

type BigInt struct {
	mappedValues map[int32]*big.Int
	stateStack   []*BigInt
}

func NewBigIntContext() (*BigInt, error) {
	bigInt := &BigInt{
		mappedValues: make(map[int32]*big.Int),
		stateStack:   make([]*BigInt, 0),
	}

	return bigInt, nil
}

func (bigInt *BigInt) InitState() {
	bigInt.mappedValues = make(map[int32]*big.Int)
}

func (bigInt *BigInt) PushState() {
	newState := &BigInt{
		mappedValues: bigInt.mappedValues,
	}

	bigInt.stateStack = append(bigInt.stateStack, newState)
}

func (bigInt *BigInt) PopState() error {
	stateStackLen := len(bigInt.stateStack)
	if stateStackLen < 1 {
		return arwen.StateStackUnderflow
	}

	prevState := bigInt.stateStack[stateStackLen-1]
	bigInt.stateStack = bigInt.stateStack[:stateStackLen-1]

	bigInt.mappedValues = prevState.mappedValues

	return nil
}

func (bigInt *BigInt) Put(value int64) int32 {
	newIndex := int32(len(bigInt.mappedValues))
	for {
		if _, ok := bigInt.mappedValues[newIndex]; !ok {
			break
		}
		newIndex++
	}

	bigInt.mappedValues[newIndex] = big.NewInt(value)

	return newIndex
}

func (bigInt *BigInt) GetOne(id int32) *big.Int {
	if _, ok := bigInt.mappedValues[id]; !ok {
		bigInt.mappedValues[id] = big.NewInt(0)
	}

	return bigInt.mappedValues[id]
}

func (bigInt *BigInt) GetTwo(id1 int32, id2 int32) (*big.Int, *big.Int) {
	if _, ok := bigInt.mappedValues[id1]; !ok {
		bigInt.mappedValues[id1] = big.NewInt(0)
	}

	if _, ok := bigInt.mappedValues[id2]; !ok {
		bigInt.mappedValues[id2] = big.NewInt(0)
	}

	return bigInt.mappedValues[id1], bigInt.mappedValues[id2]
}

func (bigInt *BigInt) GetThree(id1 int32, id2 int32, id3 int32) (*big.Int, *big.Int, *big.Int) {
	if _, ok := bigInt.mappedValues[id1]; !ok {
		bigInt.mappedValues[id1] = big.NewInt(0)
	}

	if _, ok := bigInt.mappedValues[id2]; !ok {
		bigInt.mappedValues[id2] = big.NewInt(0)
	}

	if _, ok := bigInt.mappedValues[id3]; !ok {
		bigInt.mappedValues[id3] = big.NewInt(0)
	}

	return bigInt.mappedValues[id1], bigInt.mappedValues[id2], bigInt.mappedValues[id3]
}

func (bigInt *BigInt) IsInterfaceNil() bool {
	return bigInt == nil
}

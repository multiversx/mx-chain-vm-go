package subcontexts

import (
	"math/big"
	"sync"
)

type BigInt struct {
	mutex        sync.Mutex
	mappedValues map[int32]*big.Int
}

func NewBigIntSubcontext() *BigInt {
	return &BigInt{
		mutex: sync.Mutex{},
		mappedValues: make(map[int32]*big.Int),
	}
}

func (bigInt *BigInt) CreateStateCopy() *BigInt {
	return &BigInt{
		mutex: sync.Mutex{},
		mappedValues: bigInt.mappedValues,
	}
}

func (bigInt *BigInt) LoadFromStateCopy(otherBigInt *BigInt) {
	bigInt.mappedValues = otherBigInt.mappedValues
}

func (bigInt *BigInt) Clean() {
	bigInt.mutex.Lock()
	bigInt.mappedValues = make(map[int32]*big.Int)
	bigInt.mutex.Unlock()
}

func (bigInt *BigInt) Put(value int64) int32 {
	bigInt.mutex.Lock()

	newIndex := int32(len(bigInt.mappedValues))
	for {
		if _, ok := bigInt.mappedValues[newIndex]; !ok {
			break
		}
		newIndex++
	}

	bigInt.mappedValues[newIndex] = big.NewInt(value)

	bigInt.mutex.Unlock()

	return newIndex
}

func (bigInt *BigInt) GetOne(id int32) *big.Int {
	bigInt.mutex.Lock()
	defer bigInt.mutex.Unlock()

	if _, ok := bigInt.mappedValues[id]; !ok {
		bigInt.mappedValues[id] = big.NewInt(0)
	}

	return bigInt.mappedValues[id]
}

func (bigInt *BigInt) GetTwo(id1 int32, id2 int32) (*big.Int, *big.Int) {
	bigInt.mutex.Lock()
	defer bigInt.mutex.Unlock()

	if _, ok := bigInt.mappedValues[id1]; !ok {
		bigInt.mappedValues[id1] = big.NewInt(0)
	}

	if _, ok := bigInt.mappedValues[id2]; !ok {
		bigInt.mappedValues[id2] = big.NewInt(0)
	}

	return bigInt.mappedValues[id1], bigInt.mappedValues[id2]
}

func (bigInt *BigInt) GetThree(id1 int32, id2 int32, id3 int32) (*big.Int, *big.Int, *big.Int) {
	bigInt.mutex.Lock()
	defer bigInt.mutex.Unlock()

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

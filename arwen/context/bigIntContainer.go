package context

import (
	"math/big"
	"sync"
)

type bigIntContainer struct {
	mutex        sync.Mutex
	mappedValues map[int32]*big.Int
}

func NewBigIntContainer() *bigIntContainer {
	return &bigIntContainer{
		mutex:        sync.Mutex{},
		mappedValues: make(map[int32]*big.Int),
	}
}

func (b *bigIntContainer) Clean() {
	b.mutex.Lock()
	b.mappedValues = make(map[int32]*big.Int)
	b.mutex.Unlock()
}

func (b *bigIntContainer) Put(value int64) int32 {
	b.mutex.Lock()

	newIndex := int32(len(b.mappedValues))
	for {
		if _, ok := b.mappedValues[newIndex]; !ok {
			break
		}
		newIndex++
	}

	b.mappedValues[newIndex] = big.NewInt(value)

	b.mutex.Unlock()

	return newIndex
}

func (b *bigIntContainer) GetOne(id int32) *big.Int {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, ok := b.mappedValues[id]; !ok {
		b.mappedValues[id] = big.NewInt(0)
	}

	return b.mappedValues[id]
}

func (b *bigIntContainer) GetTwo(id1, id2 int32) (*big.Int, *big.Int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, ok := b.mappedValues[id1]; !ok {
		b.mappedValues[id1] = big.NewInt(0)
	}

	if _, ok := b.mappedValues[id2]; !ok {
		b.mappedValues[id2] = big.NewInt(0)
	}

	return b.mappedValues[id1], b.mappedValues[id2]
}

func (b *bigIntContainer) GetThree(id1, id2, id3 int32) (*big.Int, *big.Int, *big.Int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, ok := b.mappedValues[id1]; !ok {
		b.mappedValues[id1] = big.NewInt(0)
	}

	if _, ok := b.mappedValues[id2]; !ok {
		b.mappedValues[id2] = big.NewInt(0)
	}

	if _, ok := b.mappedValues[id3]; !ok {
		b.mappedValues[id3] = big.NewInt(0)
	}

	return b.mappedValues[id1], b.mappedValues[id2], b.mappedValues[id3]
}

func (b *bigIntContainer) IsInterfaceNil() bool {
	return b == nil
}

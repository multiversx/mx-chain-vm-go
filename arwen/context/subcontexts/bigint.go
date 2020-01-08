package subcontexts

import (
	"math/big"
)

type BigInt struct {
}

func (b *BigInt) Put(value int64) int32 {
	panic("not implemented")
}

func (b *BigInt) GetOne(id int32) *big.Int {
	panic("not implemented")
}

func (b *BigInt) GetTwo(id1 int32, id2 int32) (*big.Int, *big.Int) {
	panic("not implemented")
}

func (b *BigInt) GetThree(id1 int32, id2 int32, id3 int32) (*big.Int, *big.Int, *big.Int) {
	panic("not implemented")
}

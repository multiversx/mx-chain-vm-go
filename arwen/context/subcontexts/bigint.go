package subcontexts

import (
	"math/big"
)

type BigInt struct {
}

func (bigInt *BigInt) Put(value int64) int32 {
	panic("not implemented")
}

func (bigInt *BigInt) GetOne(id int32) *big.Int {
	panic("not implemented")
}

func (bigInt *BigInt) GetTwo(id1 int32, id2 int32) (*big.Int, *big.Int) {
	panic("not implemented")
}

func (bigInt *BigInt) GetThree(id1 int32, id2 int32, id3 int32) (*big.Int, *big.Int, *big.Int) {
	panic("not implemented")
}

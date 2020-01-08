package subcontexts

import (
	"math/big"
)

type Output struct {
}

func (o *Output) WriteLog(addr []byte, topics [][]byte, data []byte) {
	panic("not implemented")
}

func (o *Output) Transfer(destination []byte, sender []byte, gasLimit uint64, value *big.Int, input []byte) {
	panic("not implemented")
}

func (o *Output) ReturnData() [][]byte {
	panic("not implemented")
}

func (o *Output) ClearReturnData() {
	panic("not implemented")
}

func (o *Output) Finish(data []byte) {
	panic("not implemented")
}


package subcontexts

import (
	"math/big"
)

type Output struct {
}

func (output *Output) WriteLog(addr []byte, topics [][]byte, data []byte) {
	panic("not implemented")
}

func (output *Output) Transfer(destination []byte, sender []byte, gasLimit uint64, value *big.Int, input []byte) {
	panic("not implemented")
}

func (output *Output) ReturnData() [][]byte {
	panic("not implemented")
}

func (output *Output) ClearReturnData() {
	panic("not implemented")
}

func (output *Output) Finish(data []byte) {
	panic("not implemented")
}


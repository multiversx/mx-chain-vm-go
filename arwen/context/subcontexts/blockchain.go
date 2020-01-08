package subcontexts

import (
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type Blockchain struct {
}

func (blockchain *Blockchain) AccountExists(addr []byte) bool {
	panic("not implemented")
}

func (blockchain *Blockchain) GetBalance(addr []byte) []byte {
	panic("not implemented")
}

func (blockchain *Blockchain) GetNonce(addr []byte) uint64 {
	panic("not implemented")
}

func (blockchain *Blockchain) GetCodeHash(addr []byte) ([]byte, error) {
	panic("not implemented")
}

func (blockchain *Blockchain) GetCode(addr []byte) ([]byte, error) {
	panic("not implemented")
}

func (blockchain *Blockchain) SelfDestruct(addr []byte, beneficiary []byte) {
	panic("not implemented")
}

func (blockchain *Blockchain) GetVMInput() vmcommon.VMInput {
	panic("not implemented")
}

func (blockchain *Blockchain) BlockHash(number int64) []byte {
	panic("not implemented")
}



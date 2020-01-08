package subcontexts

import (
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type Blockchain struct {
}

func (b *Blockchain) AccountExists(addr []byte) bool {
	panic("not implemented")
}

func (b *Blockchain) GetBalance(addr []byte) []byte {
	panic("not implemented")
}

func (b *Blockchain) GetNonce(addr []byte) uint64 {
	panic("not implemented")
}

func (b *Blockchain) GetCodeHash(addr []byte) ([]byte, error) {
	panic("not implemented")
}

func (b *Blockchain) GetCode(addr []byte) ([]byte, error) {
	panic("not implemented")
}

func (b *Blockchain) SelfDestruct(addr []byte, beneficiary []byte) {
	panic("not implemented")
}

func (b *Blockchain) GetVMInput() vmcommon.VMInput {
	panic("not implemented")
}

func (b *Blockchain) BlockHash(number int64) []byte {
	panic("not implemented")
}



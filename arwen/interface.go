package arwen

import (
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"math/big"
)

type HostContext interface {
	Arguments() []*big.Int
	Function() string
	AccountExists(addr []byte) bool
	GetStorage(addr []byte, key []byte) []byte
	SetStorage(addr []byte, key []byte, value []byte) int32
	GetBalance(addr []byte) []byte
	GetCodeSize(addr []byte) int
	BlockHash(nonce int64) []byte
	GetCodeHash(addr []byte) []byte
	GetCode(addr []byte) []byte
	SelfDestruct(addr []byte, beneficiary []byte)
	GetVMInput() vmcommon.VMInput
	GetSCAddress() []byte
	WriteLog(addr []byte, topics [][]byte, data []byte)
	Transfer(destination []byte, sender []byte, value *big.Int, input []byte, gas int64) (gasLeft int64, err error)
	Finish(data []byte)
	BlockChainHook() vmcommon.BlockchainHook
	BigIntContainer() BigIntContainer
	SignalUserError()
}

type EthContext interface {
	CallData() []byte
	UseGas(gas int64)
	GasLeft() int64
	BlockGasLimit() int64
}

type BigIntContainer interface {
	Clean()
	GetOne(id int32) *big.Int
	GetTwo(id1, id2 int32) (*big.Int, *big.Int)
	GetThree(id1, id2, id3 int32) (*big.Int, *big.Int, *big.Int)
	IsInterfaceNil() bool
}

package arwen

import (
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"math/big"
)

type EthContext interface {
	GetSCAddress() []byte
	CallData() []byte
	UseGas(gas int64)
	GasLeft() int64
	BlockGasLimit() int64
	GetBalance(addr []byte) []byte
	BlockHash(nonce int64) []byte
	GetVMInput() vmcommon.VMInput
	Transfer(destination []byte, sender []byte, value *big.Int, input []byte, gas int64) (gasLeft int64, err error)
	GetStorage(addr []byte, key []byte) []byte
	SetStorage(addr []byte, key []byte, value []byte) int32
	GetCode(addr []byte) []byte
	GetCodeSize(addr []byte) int32
	SignalUserError()
	Finish(data []byte)
	WriteLog(addr []byte, topics [][]byte, data []byte)
	SelfDestruct(addr []byte, beneficiary []byte)
}

type HostContext interface {
	Function() string
	Arguments() []*big.Int
	GetStorage(addr []byte, key []byte) []byte
	SetStorage(addr []byte, key []byte, value []byte) int32
	GetBalance(addr []byte) []byte
	BlockHash(nonce int64) []byte
	GetVMInput() vmcommon.VMInput
	GetSCAddress() []byte
	WriteLog(addr []byte, topics [][]byte, data []byte)
	Transfer(destination []byte, sender []byte, value *big.Int, input []byte, gas int64) (gasLeft int64, err error)
	Finish(data []byte)
	BlockChainHook() vmcommon.BlockchainHook
	SignalUserError()
}

type BigIntContext interface {
	Put(value int64) int32
	GetOne(id int32) *big.Int
	GetTwo(id1, id2 int32) (*big.Int, *big.Int)
	GetThree(id1, id2, id3 int32) (*big.Int, *big.Int, *big.Int)
	Arguments() []*big.Int
	GetStorage(addr []byte, key []byte) []byte
	SetStorage(addr []byte, key []byte, value []byte) int32
	GetVMInput() vmcommon.VMInput
	GetBalance(addr []byte) []byte
	GetSCAddress() []byte
	Finish(data []byte)
}

type VMContext interface {
	EthContext() EthContext
	CoreContext() HostContext
	BigInContext() BigIntContext
}

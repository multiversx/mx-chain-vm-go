package arwen

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"math/big"
)

type EthContext interface {
	GasSchedule() *config.GasCost
	GetSCAddress() []byte
	CallData() []byte
	UseGas(gas uint64)
	GasLeft() uint64
	BoundGasLimit(int64) uint64
	BlockGasLimit() uint64
	GetBalance(addr []byte) []byte
	BlockHash(nonce int64) []byte
	GetVMInput() vmcommon.VMInput
	GetStorage(addr []byte, key []byte) []byte
	SetStorage(addr []byte, key []byte, value []byte) int32
	GetCode(addr []byte) []byte
	GetCodeSize(addr []byte) int32
	SignalUserError()
	Finish(data []byte)
	WriteLog(addr []byte, topics [][]byte, data []byte)
	SelfDestruct(addr []byte, beneficiary []byte)
	BlockChainHook() vmcommon.BlockchainHook
	Transfer(destination []byte, sender []byte, value *big.Int, input []byte)
	ReturnData() [][]byte
	ClearReturnData()

	SetReadOnly(readOnly bool)
	CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error)
	ExecuteOnSameContext(input *vmcommon.ContractCallInput) error
	ExecuteOnDestContext(input *vmcommon.ContractCallInput) error
}

type HostContext interface {
	GasSchedule() *config.GasCost
	UseGas(gas uint64)
	GasLeft() uint64
	BoundGasLimit(int64) uint64
	Function() string
	Arguments() [][]byte
	GetStorage(addr []byte, key []byte) []byte
	SetStorage(addr []byte, key []byte, value []byte) int32
	GetBalance(addr []byte) []byte
	BlockHash(nonce int64) []byte
	GetVMInput() vmcommon.VMInput
	GetSCAddress() []byte
	WriteLog(addr []byte, topics [][]byte, data []byte)
	Transfer(destination []byte, sender []byte, value *big.Int, input []byte)
	Finish(data []byte)
	BlockChainHook() vmcommon.BlockchainHook
	SignalUserError()
	ReturnData() [][]byte

	SetReadOnly(readOnly bool)
	CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error)
	ExecuteOnSameContext(input *vmcommon.ContractCallInput) error
	ExecuteOnDestContext(input *vmcommon.ContractCallInput) error
}

type BigIntContext interface {
	GasSchedule() *config.GasCost
	UseGas(gas uint64)
	GasLeft() uint64
	Put(value int64) int32
	GetOne(id int32) *big.Int
	GetTwo(id1, id2 int32) (*big.Int, *big.Int)
	GetThree(id1, id2, id3 int32) (*big.Int, *big.Int, *big.Int)
	Arguments() [][]byte
	GetStorage(addr []byte, key []byte) []byte
	SetStorage(addr []byte, key []byte, value []byte) int32
	GetVMInput() vmcommon.VMInput
	GetBalance(addr []byte) []byte
	GetSCAddress() []byte
	Finish(data []byte)
}

type CryptoContext interface {
	GasSchedule() *config.GasCost
	UseGas(gas uint64)
	CryptoHooks() vmcommon.CryptoHook
}

type VMContext interface {
	EthContext() EthContext
	CoreContext() HostContext
	BigInContext() BigIntContext
	CryptoContext() CryptoContext
}

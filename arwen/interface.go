package arwen

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type StorageStatus int

const (
	StorageUnchanged StorageStatus = 0
	StorageModified  StorageStatus = 1
	StorageAdded     StorageStatus = 3
	StorageDeleted   StorageStatus = 4
)

type BreakpointValue uint64

const (
	BreakpointNone        BreakpointValue = 0
	BreakpointAbort       BreakpointValue = 1
	BreakpointAsyncCall   BreakpointValue = 2
	BreakpointSignalError BreakpointValue = 3
	BreakpointOutOfGas    BreakpointValue = 4
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
	GetCode(addr []byte) ([]byte, error)
	GetCodeSize(addr []byte) (int32, error)
	SignalUserError()
	Finish(data []byte)
	WriteLog(addr []byte, topics [][]byte, data []byte)
	SelfDestruct(addr []byte, beneficiary []byte)
	BlockChainHook() vmcommon.BlockchainHook
	Transfer(destination []byte, sender []byte, gasLimit uint64, value *big.Int, input []byte)
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
	Transfer(destination []byte, sender []byte, gasLimit uint64, value *big.Int, input []byte)
	Finish(data []byte)
	BlockChainHook() vmcommon.BlockchainHook
	SignalUserError()
	ReturnData() [][]byte

	SetRuntimeBreakpointValue(value BreakpointValue)
	GetRuntimeBreakpointValue() BreakpointValue

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
	GasLeft() uint64
	CryptoHooks() vmcommon.CryptoHook
	SignalUserError()
}

type VMContext interface {
	StateStack

	Crypto() vmcommon.CryptoHook
	Blockchain() BlockchainSubcontext
	Runtime() RuntimeSubcontext
	BigInt() BigIntSubcontext
	Output() OutputSubcontext
	Metering() MeteringSubcontext
	Storage() StorageSubcontext
}

type StateStack interface {
	PushState()
	PopState() error
}

type BlockchainSubcontext interface {
	InitState()
	NewAddress(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error)
	AccountExists(addr []byte) bool
	GetBalance(addr []byte) []byte
	GetNonce(addr []byte) uint64
	CurrentEpoch() uint32
	GetStateRootHash() []byte
	LastTimeStamp() uint64
	LastNonce() uint64
	LastRound() uint64
	LastEpoch() uint32
	CurrentRound() uint64
	CurrentNonce() uint64
	CurrentTimeStamp() uint64
	CurrentRandomSeed() []byte
	LastRandomSeed() []byte
	IncreaseNonce(addr []byte)
	GetCodeHash(addr []byte) ([]byte, error)
	GetCode(addr []byte) ([]byte, error)
	GetCodeSize(addr []byte) (int32, error)
	BlockHash(number int64) []byte
}

type RuntimeSubcontext interface {
	StateStack

	InitState()
	GetVMInput() *vmcommon.VMInput
	GetSCAddress() []byte
	Function() string
	Arguments() [][]byte
	SignalUserError()
	SetRuntimeBreakpointValue(value BreakpointValue)
	GetRuntimeBreakpointValue() BreakpointValue
	CallData() []byte
	ReadOnly() bool
	SetReadOnly(readOnly bool)
	CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error)
	ExecuteOnSameContext(input *vmcommon.ContractCallInput) error
	ExecuteOnDestContext(input *vmcommon.ContractCallInput) error
	SetInstanceContext(instCtx *wasmer.InstanceContext)
	GetInstanceContext() *wasmer.InstanceContext
	GetPointsUsed() uint64
	SetPointsUsed(gasPoints uint64)
	MemStore(offset int32, data []byte) error
	MemLoad(offset int32, length int32) ([]byte, error)
	Clean()
	SetInstanceContextId(id int)
}

type BigIntSubcontext interface {
	StateStack

	InitState()
	Put(value int64) int32
	GetOne(id int32) *big.Int
	GetTwo(id1, id2 int32) (*big.Int, *big.Int)
	GetThree(id1, id2, id3 int32) (*big.Int, *big.Int, *big.Int)
}

// TODO find better name
type OutputSubcontext interface {
	StateStack

	InitState()
	GetOutputAccounts() map[string]*vmcommon.OutputAccount
	GetStorageUpdates() map[string](map[string][]byte)
	WriteLog(addr []byte, topics [][]byte, data []byte)
	Transfer(destination []byte, sender []byte, gasLimit uint64, value *big.Int, input []byte)
	SelfDestruct(addr []byte, beneficiary []byte)
	GetRefund() uint64
	SetRefund(refund uint64)
	ReturnCode() vmcommon.ReturnCode
	SetReturnCode(returnCode vmcommon.ReturnCode)
	ReturnData() [][]byte
	ClearReturnData()
	Finish(data []byte)
  CreateVMOutput(result []byte) *vmcommon.VMOutput
}

type MeteringSubcontext interface {
	InitState()
	GasSchedule() *config.GasCost
	UseGas(gas uint64)
	FreeGas(gas uint64)
	GasLeft() uint64
	BoundGasLimit(value int64) uint64
	BlockGasLimit() uint64
}

type StorageSubcontext interface {
	InitState()
	GetStorage(addr []byte, key []byte) []byte
	SetStorage(addr []byte, key []byte, value []byte) int32
}

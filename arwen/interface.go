package arwen

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type StateStack interface {
	InitState()
	PushState()
	PopState() error
}

type VMHost interface {
	StateStack

	Crypto() vmcommon.CryptoHook
	Blockchain() BlockchainContext
	Runtime() RuntimeContext
	BigInt() BigIntContext
	Output() OutputContext
	Metering() MeteringContext
	Storage() StorageContext

	CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error)
	ExecuteOnSameContext(input *vmcommon.ContractCallInput) error
	ExecuteOnDestContext(input *vmcommon.ContractCallInput) error
	EthereumCallData() []byte
}

type BlockchainContext interface {
	NewAddress(creatorAddress []byte) ([]byte, error)
	AccountExists(addr []byte) bool
	GetBalance(addr []byte) []byte
	GetNonce(addr []byte) (uint64, error)
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

type RuntimeContext interface {
	StateStack

	InitStateFromContractCallInput(input *vmcommon.ContractCallInput)
	GetVMInput() *vmcommon.VMInput
	SetVMInput(vmInput *vmcommon.VMInput)
	GetSCAddress() []byte
	SetSCAddress(scAddress []byte)
	GetVMType() []byte
	Function() string
	Arguments() [][]byte
	SignalUserError(message string)
	SignalExit(exitCode int)
	SetRuntimeBreakpointValue(value BreakpointValue)
	GetRuntimeBreakpointValue() BreakpointValue
	PushInstance()
	PopInstance() error
	ReadOnly() bool
	SetReadOnly(readOnly bool)
	CreateWasmerInstance(contract []byte, gasLimit uint64) error
	SetInstanceContext(instCtx *wasmer.InstanceContext)
	GetInstanceContext() *wasmer.InstanceContext
	GetInstanceExports() wasmer.ExportsMap
	GetInitFunction() wasmer.ExportedFunctionCallback
	GetFunctionToCall() (wasmer.ExportedFunctionCallback, error)
	GetPointsUsed() uint64
	SetPointsUsed(gasPoints uint64)
	MemStore(offset int32, data []byte) error
	MemLoad(offset int32, length int32) ([]byte, error)
	CleanInstance()
	SetInstanceContextId(id int)
}

type BigIntContext interface {
	StateStack

	Put(value int64) int32
	GetOne(id int32) *big.Int
	GetTwo(id1, id2 int32) (*big.Int, *big.Int)
	GetThree(id1, id2, id3 int32) (*big.Int, *big.Int, *big.Int)
}

// TODO find a better name
type OutputContext interface {
	StateStack

	GetOutputAccount(address []byte) (*vmcommon.OutputAccount, bool)
	WriteLog(address []byte, topics [][]byte, data []byte)
	Transfer(destination []byte, sender []byte, gasLimit uint64, value *big.Int, input []byte)
	SelfDestruct(address []byte, beneficiary []byte)
	GetRefund() uint64
	SetRefund(refund uint64)
	ReturnCode() vmcommon.ReturnCode
	SetReturnCode(returnCode vmcommon.ReturnCode)
	ReturnMessage() string
	SetReturnMessage(message string)
	ReturnData() [][]byte
	ClearReturnData()
	Finish(data []byte)
	FinishValue(value wasmer.Value)
	GetVMOutput(result wasmer.Value) *vmcommon.VMOutput
	AddTxValueToAccount(address []byte, value *big.Int)
	DeployCode(address []byte, code []byte)
	CreateVMOutputInCaseOfError(errCode vmcommon.ReturnCode, message string) *vmcommon.VMOutput
}

type MeteringContext interface {
	GasSchedule() *config.GasCost
	UseGas(gas uint64)
	FreeGas(gas uint64)
	GasLeft() uint64
	BoundGasLimit(value int64) uint64
	BlockGasLimit() uint64
	DeductInitialGasForExecution(input *vmcommon.ContractCallInput, contract []byte) (uint64, error)
	DeductInitialGasForDirectDeployment(input *vmcommon.ContractCreateInput) (uint64, error)
	DeductInitialGasForIndirectDeployment(input *vmcommon.ContractCreateInput) (uint64, error)
}

type StorageContext interface {
	GetStorageUpdates(address []byte) map[string]*vmcommon.StorageUpdate
	GetStorage(address []byte, key []byte) []byte
	SetStorage(address []byte, key []byte, value []byte) int32
}

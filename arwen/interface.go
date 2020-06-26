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
	PopSetActiveState()
	PopDiscard()
	ClearStateStack()
}

// CallArgsParser defines the functionality to parse transaction data for a smart contract call
type CallArgsParser interface {
	ParseData(data string) (string, [][]byte, error)
	IsInterfaceNil() bool
}

type VMHost interface {
	Crypto() vmcommon.CryptoHook
	Blockchain() BlockchainContext
	Runtime() RuntimeContext
	BigInt() BigIntContext
	Output() OutputContext
	Metering() MeteringContext
	Storage() StorageContext

	CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error)
	ExecuteOnSameContext(input *vmcommon.ContractCallInput) (*AsyncContext, error)
	ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, *AsyncContext, error)
	EthereumCallData() []byte
	GetAPIMethods() *wasmer.Imports
	GetProtocolBuiltinFunctions() vmcommon.FunctionNames
}

type BlockchainContext interface {
	NewAddress(creatorAddress []byte) ([]byte, error)
	AccountExists(addr []byte) bool
	GetBalance(addr []byte) []byte
	GetBalanceBigInt(addr []byte) *big.Int
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
	GetOwnerAddress() ([]byte, error)
	GetShardOfAddress(addr []byte) uint32
	IsSmartContract(addr []byte) bool
}

type RuntimeContext interface {
	StateStack

	InitStateFromContractCallInput(input *vmcommon.ContractCallInput)
	SetCustomCallFunction(callFunction string)
	GetVMInput() *vmcommon.VMInput
	SetVMInput(vmInput *vmcommon.VMInput)
	GetSCAddress() []byte
	SetSCAddress(scAddress []byte)
	GetVMType() []byte
	Function() string
	Arguments() [][]byte
	GetCurrentTxHash() []byte
	GetOriginalTxHash() []byte
	ExtractCodeUpgradeFromArgs() ([]byte, []byte, error)
	SignalUserError(message string)
	FailExecution(err error)
	SetRuntimeBreakpointValue(value BreakpointValue)
	GetRuntimeBreakpointValue() BreakpointValue
	GetAsyncCallInfo() *AsyncCallInfo
	SetAsyncCallInfo(asyncCallInfo *AsyncCallInfo)
	AddAsyncCall(contextIdentifier []byte, asyncCall *AsyncCall) error
	GetAsyncContext() *AsyncContext
	GetAsyncCallGroup(groupID []byte) (*AsyncCallGroup, error)
	PushInstance()
	PopInstance()
	RunningInstancesCount() uint64
	ClearInstanceStack()
	ReadOnly() bool
	SetReadOnly(readOnly bool)
	StartWasmerInstance(contract []byte, gasLimit uint64) error
	SetMaxInstanceCount(uint64)
	VerifyContractCode() error
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
	SetInstanceContextID(id int)
	ElrondAPIErrorShouldFailExecution() bool
	CryptoAPIErrorShouldFailExecution() bool
	BigIntAPIErrorShouldFailExecution() bool
}

type BigIntContext interface {
	StateStack

	Put(value int64) int32
	GetOne(id int32) *big.Int
	GetTwo(id1, id2 int32) (*big.Int, *big.Int)
	GetThree(id1, id2, id3 int32) (*big.Int, *big.Int, *big.Int)
}

type OutputContext interface {
	StateStack
	PopMergeActiveState()
	CensorVMOutput()
	AddToActiveState(rightOutput *vmcommon.VMOutput)

	GetOutputAccount(address []byte) (*vmcommon.OutputAccount, bool)
	WriteLog(address []byte, topics [][]byte, data []byte)
	Transfer(destination []byte, sender []byte, gasLimit uint64, value *big.Int, input []byte) error
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
	GetVMOutput() *vmcommon.VMOutput
	AddTxValueToAccount(address []byte, value *big.Int)
	DeployCode(input CodeDeployInput)
	CreateVMOutputInCaseOfError(err error) *vmcommon.VMOutput
}

type MeteringContext interface {
	GasSchedule() *config.GasCost
	UseGas(gas uint64)
	FreeGas(gas uint64)
	RestoreGas(gas uint64)
	GasLeft() uint64
	BoundGasLimit(value int64) uint64
	BlockGasLimit() uint64
	DeductInitialGasForExecution(contract []byte) error
	DeductInitialGasForDirectDeployment(input CodeDeployInput) error
	DeductInitialGasForIndirectDeployment(input CodeDeployInput) error
	UnlockGasIfAsyncStep()
	GetGasLockedForAsyncStep() uint64
}

type StorageStatus int

const (
	StorageUnchanged StorageStatus = iota
	StorageModified
	StorageAdded
	StorageDeleted
)

type StorageContext interface {
	StateStack

	SetAddress(address []byte)
	GetStorageUpdates(address []byte) map[string]*vmcommon.StorageUpdate
	GetStorage(key []byte) []byte
	SetStorage(key []byte, value []byte) (StorageStatus, error)
}

type AsyncCallInfoHandler interface {
	GetDestination() []byte
	GetData() []byte
	GetGasLimit() uint64
	GetValueBytes() []byte
}

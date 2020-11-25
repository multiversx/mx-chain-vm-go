package arwen

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
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
	Crypto() crypto.VMCrypto
	Blockchain() BlockchainContext
	Runtime() RuntimeContext
	Async() AsyncContext
	BigInt() BigIntContext
	Output() OutputContext
	Metering() MeteringContext
	Storage() StorageContext
	IsArwenV2Enabled() bool
	IsAheadOfTimeCompileEnabled() bool
	IsDynamicGasLockingEnabled() bool

	CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error)
	ExecuteOnSameContext(input *vmcommon.ContractCallInput) error
	ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error)
	EthereumCallData() []byte
	GetAPIMethods() *wasmer.Imports
	GetProtocolBuiltinFunctions() vmcommon.FunctionNames
	IsBuiltinFunctionName(functionName string) bool
	CallArgsParser() CallArgsParser
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
	GetCodeHash(addr []byte) []byte
	GetCode(addr []byte) ([]byte, error)
	GetCodeSize(addr []byte) (int32, error)
	BlockHash(number uint64) []byte
	GetOwnerAddress() ([]byte, error)
	GetShardOfAddress(addr []byte) uint32
	IsSmartContract(addr []byte) bool
	IsPayable(address []byte) (bool, error)
	SaveCompiledCode(codeHash []byte, code []byte)
	GetCompiledCode(codeHash []byte) (bool, []byte)
}

type RuntimeContext interface {
	StateStack

	InitStateFromInput(input *vmcommon.ContractCallInput)
	SetCustomCallFunction(callFunction string)
	GetVMInput() *vmcommon.VMInput
	SetVMInput(vmInput *vmcommon.VMInput)
	GetSCAddress() []byte
	SetSCAddress(scAddress []byte)
	GetSCCode() ([]byte, error)
	GetSCCodeSize() uint64
	GetVMType() []byte
	Function() string
	Arguments() [][]byte
	HasCallbackMethod() bool
	GetCurrentTxHash() []byte
	GetOriginalTxHash() []byte
	GetPrevTxHash() []byte
	ExtractCodeUpgradeFromArgs() ([]byte, []byte, error)
	SignalUserError(message string)
	FailExecution(err error)
	MustVerifyNextContractCode()
	ValidateCallbackName(callbackName string) error
	SetRuntimeBreakpointValue(value BreakpointValue)
	GetRuntimeBreakpointValue() BreakpointValue
	PushInstance()
	PopInstance()
	RunningInstancesCount() uint64
	ClearInstanceStack()
	IsWarmInstance() bool
	ResetWarmInstance()
	ReadOnly() bool
	SetReadOnly(readOnly bool)
	StartWasmerInstance(contract []byte, gasLimit uint64, newCode bool) error
	CleanWasmerInstance()
	SetMaxInstanceCount(uint64)
	VerifyContractCode() error
	GetInstanceExports() wasmer.ExportsMap
	GetInitFunction() wasmer.ExportedFunctionCallback
	GetFunctionToCall() (wasmer.ExportedFunctionCallback, error)
	GetPointsUsed() uint64
	SetPointsUsed(gasPoints uint64)
	MemStore(offset int32, data []byte) error
	MemLoad(offset int32, length int32) ([]byte, error)
	ElrondAPIErrorShouldFailExecution() bool
	ElrondSyncExecAPIErrorShouldFailExecution() bool
	CryptoAPIErrorShouldFailExecution() bool
	BigIntAPIErrorShouldFailExecution() bool
}

type AsyncContext interface {
	StateStack

	InitStateFromInput(input *vmcommon.ContractCallInput)
	AddCall(groupID string, call *AsyncCall) error
	AddCallGroup(group *AsyncCallGroup) error
	HasPendingCallGroups() bool
	IsComplete() bool
	FindCall(destination []byte) (string, int, error)
	GetCallGroup(groupID string) (*AsyncCallGroup, bool)
	DeleteCallGroupByID(groupID string)
	DeleteCallGroup(index int)
	SetCaller(caller []byte)
	SetGasPrice(gasPrice uint64)
	SetGroupCallback(groupID string, callbackName string, data []byte, gas uint64) error
	PostprocessCrossShardCallback() error
	GetCallerAddress() []byte
	GetReturnData() []byte
	Load() error
	Save() error
	Delete() error
	Execute() error
	PrepareLegacyAsyncCall(address []byte, data []byte, value []byte) error
	UpdateCurrentCallStatus() (*AsyncCall, error)
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
	ResetGas()
	AddToActiveState(rightOutput *vmcommon.VMOutput)

	GetOutputAccount(address []byte) (*vmcommon.OutputAccount, bool)
	DeleteOutputAccount(address []byte)
	WriteLog(address []byte, topics [][]byte, data []byte)
	TransferValueOnly(destination []byte, sender []byte, value *big.Int) error
	Transfer(destination []byte, sender []byte, gasLimit uint64, gasLocked uint64, value *big.Int, input []byte, callType vmcommon.CallType) error
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
	SetGasSchedule(gasMap config.GasScheduleMap)
	GasSchedule() *config.GasCost
	UseGas(gas uint64)
	FreeGas(gas uint64)
	RestoreGas(gas uint64)
	GasLeft() uint64
	BoundGasLimit(limit uint64) uint64
	BlockGasLimit() uint64
	DeductInitialGasForExecution(contract []byte) error
	DeductInitialGasForDirectDeployment(input CodeDeployInput) error
	DeductInitialGasForIndirectDeployment(input CodeDeployInput) error
	ComputeGasLockedForAsync() uint64
	UseGasForAsyncStep() error
	UseGasBounded(gasToUse uint64) error
	UnlockGasIfAsyncCallback()
	GetGasLocked() uint64
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
	GetStorageFromAddress(address []byte, key []byte) []byte
	GetStorage(key []byte) []byte
	GetStorageUnmetered(key []byte) []byte
	SetStorage(key []byte, value []byte) (StorageStatus, error)
}

type AsyncCallHandler interface {
	GetDestination() []byte
	GetData() []byte
	GetGasLimit() uint64
	GetGasLocked() uint64
	GetValueBytes() []byte
}

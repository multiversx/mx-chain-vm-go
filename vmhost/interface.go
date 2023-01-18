package vmhost

import (
	"crypto/elliptic"
	"io"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/data/esdt"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-v1_4-go/config"
	"github.com/multiversx/mx-chain-vm-v1_4-go/crypto"
	"github.com/multiversx/mx-chain-vm-v1_4-go/wasmer"
)

// StateStack defines the functionality for working with a state stack
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

// VMHost defines the functionality for working with the VM
type VMHost interface {
	vmcommon.VMExecutionHandler
	Crypto() crypto.VMCrypto
	Blockchain() BlockchainContext
	Runtime() RuntimeContext
	ManagedTypes() ManagedTypesContext
	Output() OutputContext
	Metering() MeteringContext
	Storage() StorageContext
	EnableEpochsHandler() vmcommon.EnableEpochsHandler

	ExecuteESDTTransfer(destination []byte, sender []byte, esdtTransfers []*vmcommon.ESDTTransfer, callType vm.CallType) (*vmcommon.VMOutput, uint64, error)
	CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error)
	ExecuteOnSameContext(input *vmcommon.ContractCallInput) (*AsyncContextInfo, error)
	ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, *AsyncContextInfo, error)
	GetAPIMethods() *wasmer.Imports
	IsBuiltinFunctionName(functionName string) bool
	AreInSameShard(leftAddress []byte, rightAddress []byte) bool

	GetGasScheduleMap() config.GasScheduleMap
	GetContexts() (ManagedTypesContext, BlockchainContext, MeteringContext, OutputContext, RuntimeContext, StorageContext)
	SetRuntimeContext(runtime RuntimeContext)

	SetBuiltInFunctionsContainer(builtInFuncs vmcommon.BuiltInFunctionContainer)
	InitState()

	FixOOGReturnCodeEnabled() bool
	FixFailExecutionEnabled() bool
	CreateNFTOnExecByCallerEnabled() bool
	DisableExecByCaller() bool
	CheckExecuteReadOnly() bool
	Reset()
}

// BlockchainContext defines the functionality needed for interacting with the blockchain context
type BlockchainContext interface {
	StateStack

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
	BlockHash(number int64) []byte
	GetOwnerAddress() ([]byte, error)
	GetShardOfAddress(addr []byte) uint32
	IsSmartContract(addr []byte) bool
	IsPayable(sndAddress, rcvAddress []byte) (bool, error)
	SaveCompiledCode(codeHash []byte, code []byte)
	GetCompiledCode(codeHash []byte) (bool, []byte)
	GetESDTToken(address []byte, tokenID []byte, nonce uint64) (*esdt.ESDigitalToken, error)
	IsLimitedTransfer(tokenID []byte) bool
	IsPaused(tokenID []byte) bool
	GetUserAccount(address []byte) (vmcommon.UserAccountHandler, error)
	ProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error)
	GetSnapshot() int
	RevertToSnapshot(snapshot int)
	ClearCompiledCodes()
}

// RuntimeContext defines the functionality needed for interacting with the runtime context
type RuntimeContext interface {
	StateStack

	InitStateFromContractCallInput(input *vmcommon.ContractCallInput)
	SetCustomCallFunction(callFunction string)
	GetVMInput() *vmcommon.ContractCallInput
	SetVMInput(vmInput *vmcommon.ContractCallInput)
	GetContextAddress() []byte
	SetCodeAddress(scAddress []byte)
	GetSCCode() ([]byte, error)
	GetSCCodeSize() uint64
	GetVMType() []byte
	Function() string
	Arguments() [][]byte
	GetCurrentTxHash() []byte
	GetOriginalTxHash() []byte
	ExtractCodeUpgradeFromArgs() ([]byte, []byte, error)
	SignalUserError(message string)
	FailExecution(err error)
	MustVerifyNextContractCode()
	SetRuntimeBreakpointValue(value BreakpointValue)
	GetRuntimeBreakpointValue() BreakpointValue
	GetAsyncCallInfo() *AsyncCallInfo
	SetAsyncCallInfo(asyncCallInfo *AsyncCallInfo)
	AddAsyncContextCall(contextIdentifier []byte, asyncCall *AsyncGeneratedCall) error
	GetAsyncContextInfo() *AsyncContextInfo
	GetAsyncContext(contextIdentifier []byte) (*AsyncContext, error)
	RunningInstancesCount() uint64
	IsFunctionImported(name string) bool
	ReadOnly() bool
	SetReadOnly(readOnly bool)
	StartWasmerInstance(contract []byte, gasLimit uint64, newCode bool) error
	ClearWarmInstanceCache()
	SetMaxInstanceCount(uint64)
	VerifyContractCode() error
	GetInstance() wasmer.InstanceHandler
	GetWarmInstance(codeHash []byte) (wasmer.InstanceHandler, bool)
	GetInstanceExports() wasmer.ExportsMap
	GetInitFunction() wasmer.ExportedFunctionCallback
	GetFunctionToCall() (wasmer.ExportedFunctionCallback, error)
	GetPointsUsed() uint64
	SetPointsUsed(gasPoints uint64)
	MemStore(offset int32, data []byte) error
	MemLoad(offset int32, length int32) ([]byte, error)
	MemLoadMultiple(offset int32, lengths []int32) ([][]byte, error)
	ElrondAPIErrorShouldFailExecution() bool
	ElrondSyncExecAPIErrorShouldFailExecution() bool
	CryptoAPIErrorShouldFailExecution() bool
	BigIntAPIErrorShouldFailExecution() bool
	BigFloatAPIErrorShouldFailExecution() bool
	ManagedBufferAPIErrorShouldFailExecution() bool
	ExecuteAsyncCall(address []byte, data []byte, value []byte) error
	CleanInstance()
	AddError(err error, otherInfo ...string)
	GetAllErrors() error
	ReplaceInstanceBuilder(builder InstanceBuilder)
	EndExecution()
	ValidateInstances() error
}

// ManagedTypesContext defines the functionality needed for interacting with the big int context
type ManagedTypesContext interface {
	StateStack

	GetRandReader() io.Reader
	ConsumeGasForThisBigIntNumberOfBytes(byteLen *big.Int)
	ConsumeGasForThisIntNumberOfBytes(byteLen int)
	ConsumeGasForBytes(bytes []byte)
	ConsumeGasForBigIntCopy(values ...*big.Int)
	ConsumeGasForBigFloatCopy(values ...*big.Float)
	NewBigInt(value *big.Int) int32
	NewBigIntFromInt64(int64Value int64) int32
	GetBigIntOrCreate(handle int32) *big.Int
	GetBigInt(id int32) (*big.Int, error)
	GetTwoBigInt(handle1 int32, handle2 int32) (*big.Int, *big.Int, error)
	PutBigFloat(value *big.Float) (int32, error)
	BigFloatPrecIsNotValid(precision uint) bool
	BigFloatExpIsNotValid(exponent int) bool
	EncodedBigFloatIsNotValid(encodedBigFloat []byte) bool
	GetBigFloatOrCreate(handle int32) (*big.Float, error)
	GetBigFloat(handle int32) (*big.Float, error)
	GetTwoBigFloats(handle1 int32, handle2 int32) (*big.Float, *big.Float, error)
	PutEllipticCurve(ec *elliptic.CurveParams) int32
	GetEllipticCurve(handle int32) (*elliptic.CurveParams, error)
	GetEllipticCurveSizeOfField(ecHandle int32) int32
	Get100xCurveGasCostMultiplier(ecHandle int32) int32
	GetScalarMult100xCurveGasCostMultiplier(ecHandle int32) int32
	GetUCompressed100xCurveGasCostMultiplier(ecHandle int32) int32
	GetPrivateKeyByteLengthEC(ecHandle int32) int32
	NewManagedBuffer() int32
	NewManagedBufferFromBytes(bytes []byte) int32
	SetBytes(mBufferHandle int32, bytes []byte)
	GetBytes(mBufferHandle int32) ([]byte, error)
	AppendBytes(mBufferHandle int32, bytes []byte) bool
	GetLength(mBufferHandle int32) int32
	GetSlice(mBufferHandle int32, startPosition int32, lengthOfSlice int32) ([]byte, error)
	DeleteSlice(mBufferHandle int32, startPosition int32, lengthOfSlice int32) ([]byte, error)
	InsertSlice(mBufferHandle int32, startPosition int32, slice []byte) ([]byte, error)
	ReadManagedVecOfManagedBuffers(managedVecHandle int32) ([][]byte, uint64, error)
	WriteManagedVecOfManagedBuffers(data [][]byte, destinationHandle int32)
}

// OutputContext defines the functionality needed for interacting with the output context
type OutputContext interface {
	StateStack
	PopMergeActiveState()
	CensorVMOutput()
	AddToActiveState(rightOutput *vmcommon.VMOutput)

	GetOutputAccount(address []byte) (*vmcommon.OutputAccount, bool)
	GetOutputAccounts() map[string]*vmcommon.OutputAccount
	DeleteOutputAccount(address []byte)
	WriteLog(address []byte, topics [][]byte, data []byte)
	TransferValueOnly(destination []byte, sender []byte, value *big.Int, checkPayable bool) error
	Transfer(destination []byte, sender []byte, gasLimit uint64, gasLocked uint64, value *big.Int, input []byte, callType vm.CallType) error
	TransferESDT(destination []byte, sender []byte, transfers []*vmcommon.ESDTTransfer, callInput *vmcommon.ContractCallInput) (uint64, error)
	SelfDestruct(address []byte, beneficiary []byte)
	GetRefund() uint64
	SetRefund(refund uint64)
	ReturnCode() vmcommon.ReturnCode
	SetReturnCode(returnCode vmcommon.ReturnCode)
	ReturnMessage() string
	SetReturnMessage(message string)
	ReturnData() [][]byte
	ClearReturnData()
	RemoveReturnData(index uint32)
	Finish(data []byte)
	PrependFinish(data []byte)
	DeleteFirstReturnData()
	GetVMOutput() *vmcommon.VMOutput
	RemoveNonUpdatedStorage()
	AddTxValueToAccount(address []byte, value *big.Int)
	DeployCode(input CodeDeployInput)
	CreateVMOutputInCaseOfError(err error) *vmcommon.VMOutput
}

// MeteringContext defines the functionality needed for interacting with the metering context
type MeteringContext interface {
	StateStack
	PopMergeActiveState()

	InitStateFromContractCallInput(input *vmcommon.VMInput)
	SetGasSchedule(gasMap config.GasScheduleMap)
	GasSchedule() *config.GasCost
	UseGas(gas uint64)
	UseAndTraceGas(gas uint64)
	UseGasAndAddTracedGas(functionName string, gas uint64)
	FreeGas(gas uint64)
	RestoreGas(gas uint64)
	GasLeft() uint64
	GasUsedForExecution() uint64
	GasSpentByContract() uint64
	GetGasForExecution() uint64
	GetGasProvided() uint64
	GetSCPrepareInitialCost() uint64
	BoundGasLimit(value int64) uint64
	BlockGasLimit() uint64
	DeductInitialGasForExecution(contract []byte) error
	DeductInitialGasForDirectDeployment(input CodeDeployInput) error
	DeductInitialGasForIndirectDeployment(input CodeDeployInput) error
	ComputeGasLockedForAsync() uint64
	UseGasForAsyncStep() error
	UseGasBounded(gasToUse uint64) error
	GetGasLocked() uint64
	UpdateGasStateOnSuccess(vmOutput *vmcommon.VMOutput) error
	UpdateGasStateOnFailure(vmOutput *vmcommon.VMOutput)
	TrackGasUsedByBuiltinFunction(builtinInput *vmcommon.ContractCallInput, builtinOutput *vmcommon.VMOutput, postBuiltinInput *vmcommon.ContractCallInput)
	StartGasTracing(functionName string)
	SetGasTracing(enableGasTracing bool)
	GetGasTrace() map[string]map[string][]uint64
}

// StorageStatus defines the states the storage can be in
type StorageStatus int

const (
	// StorageUnchanged signals that the storage was not changed
	StorageUnchanged StorageStatus = iota

	// StorageModified signals that the storage has been modified
	StorageModified

	// StorageAdded signals that something was added to storage
	StorageAdded

	// StorageDeleted signals that something was removed from storage
	StorageDeleted
)

// StorageContext defines the functionality needed for interacting with the storage context
type StorageContext interface {
	StateStack

	SetAddress(address []byte)
	GetStorageUpdates(address []byte) map[string]*vmcommon.StorageUpdate
	GetStorageFromAddress(address []byte, key []byte) ([]byte, bool, error)
	GetStorage(key []byte) ([]byte, bool, error)
	GetStorageUnmetered(key []byte) ([]byte, bool, error)
	SetStorage(key []byte, value []byte) (StorageStatus, error)
	SetProtectedStorage(key []byte, value []byte) (StorageStatus, error)
	UseGasForStorageLoad(tracedFunctionName string, blockChainLoadCost uint64, usedCache bool)
	IsUseDifferentGasCostFlagSet() bool
}

// AsyncCallInfoHandler defines the functionality for working with AsyncCallInfo
type AsyncCallInfoHandler interface {
	GetDestination() []byte
	GetData() []byte
	GetGasLimit() uint64
	GetGasLocked() uint64
	GetValueBytes() []byte
}

// InstanceBuilder defines the functionality needed to create Wasmer instances
type InstanceBuilder interface {
	NewInstanceWithOptions(contractCode []byte, options wasmer.CompilationOptions) (wasmer.InstanceHandler, error)
	NewInstanceFromCompiledCodeWithOptions(compiledCode []byte, options wasmer.CompilationOptions) (wasmer.InstanceHandler, error)
}

// GasTracing defines the functionality needed for a gas tracing
type GasTracing interface {
	BeginTrace(scAddress string, functionName string)
	AddToCurrentTrace(usedGas uint64)
	AddTracedGas(scAddress string, functionName string, usedGas uint64)
	GetGasTrace() map[string]map[string][]uint64
	IsInterfaceNil() bool
}

// HashComputer provides hash computation
type HashComputer interface {
	Compute(string) []byte
	Size() int
	IsInterfaceNil() bool
}

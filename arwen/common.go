package arwen

import (
	"github.com/ElrondNetwork/wasm-vm/config"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// ArwenVersion returns the current arwen version
const ArwenVersion = "v1.4"

// BreakpointValue encodes Wasmer runtime breakpoint types
type BreakpointValue uint64

const (
	// BreakpointNone means the lack of a breakpoint
	BreakpointNone BreakpointValue = iota

	// BreakpointExecutionFailed means that Wasmer must stop immediately
	// due to failure indicated by Arwen
	BreakpointExecutionFailed

	// BreakpointAsyncCall means that Wasmer must stop immediately
	// so the VM can execute an AsyncCall
	BreakpointAsyncCall

	// BreakpointSignalError means that Wasmer must stop immediately
	// due to a contract-signalled error
	BreakpointSignalError

	// BreakpointOutOfGas means that Wasmer must stop immediately
	// due to gas being exhausted
	BreakpointOutOfGas

	// BreakpointMemoryLimit means that Wasmer must stop immediately
	// due to over-allocation of WASM memory
	BreakpointMemoryLimit
)

const (
	// BreakpointNoneString is the human-readable name of BreakpointNone
	BreakpointNoneString = "BreakpointNone"

	// BreakpointExecutionFailedString is the human-readable name of BreakpointExecutionFailed
	BreakpointExecutionFailedString = "BreakpointExecutionFailed"

	// BreakpointAsyncCallString is the human-readable name of BreakpointAsyncCall
	BreakpointAsyncCallString = "BreakpointAsyncCall"

	// BreakpointSignalErrorString is the human-readable name of BreakpointSignalError
	BreakpointSignalErrorString = "BreakpointSignalError"

	// BreakpointOutOfGasString is the human-readable name of BreakpointOutOfGas
	BreakpointOutOfGasString = "BreakpointOutOfGas"

	// UnknownBreakpointString is the human-readable label for an unknown breakpoint value
	UnknownBreakpointString = "unknown breakpoint"
)

// String returns the human-readable name of a BreakpointValue
func (b BreakpointValue) String() string {
	switch b {
	case BreakpointNone:
		return BreakpointNoneString
	case BreakpointExecutionFailed:
		return BreakpointExecutionFailedString
	case BreakpointAsyncCall:
		return BreakpointAsyncCallString
	case BreakpointSignalError:
		return BreakpointSignalErrorString
	case BreakpointOutOfGas:
		return BreakpointOutOfGasString
	default:
		return UnknownBreakpointString
	}
}

// AsyncCallExecutionMode encodes the execution modes of an AsyncCall
type AsyncCallExecutionMode uint

const (
	// SyncExecution indicates that the async call can be executed synchronously,
	// with its corresponding callback
	SyncExecution AsyncCallExecutionMode = iota

	// AsyncBuiltinFuncIntraShard indicates that the async call is an intra-shard built in function call
	AsyncBuiltinFuncIntraShard

	// AsyncBuiltinFuncCrossShard indicates that the async call is a cross-shard call to a
	// built-in function, which is executed half in-shard, half cross-shard
	AsyncBuiltinFuncCrossShard

	// ESDTTransferOnCallBack indicated that the async call is actually a callback with ESDT transfer
	ESDTTransferOnCallBack

	// AsyncUnknown indicates that the async call cannot be executed locally, and
	// must be forwarded to the destination account
	AsyncUnknown
)

// CallbackFunctionName is the name of the default asynchronous callback
// function of a smart contract
const CallbackFunctionName = "callBack"

// TimeLockKeyPrefix is the storage key prefix used for timelock-related storage.
const TimeLockKeyPrefix = "TIMELOCK"

// AsyncDataPrefix is the storage key prefix used for AsyncContext-related storage.
const AsyncDataPrefix = "ASYNC"

// AsyncCallStatus represents the different status an async call can have
type AsyncCallStatus uint8

// LegacyAsyncCallGroupID is the AsyncCallGroup identifier reserved for the
// implementation of the legacy asyncCall() EEI function
const LegacyAsyncCallGroupID = "LegacyAsync"

const (
	// AsyncCallPending is the status of an async call that awaits complete execution
	AsyncCallPending AsyncCallStatus = iota

	// AsyncCallResolved is the status of an async call that was executed completely and successfully
	AsyncCallResolved

	// AsyncCallRejected is the status of an async call that was executed completely but unsuccessfully
	AsyncCallRejected

	// AddressLen specifies the length of the address
	AddressLen = 32

	// HashLen specifies the lenghth of a hash
	HashLen = 32

	// BalanceLen specifies the number of bytes on which the balance is stored
	BalanceLen = 32

	// CodeMetadataLen specifies the length of the code metadata
	CodeMetadataLen = 2

	// InitFunctionName specifies the name for the init function
	InitFunctionName = "init"

	// InitFunctionNameEth specifies the name for the init function on Ethereum
	InitFunctionNameEth = "solidity.ctor"

	// UpgradeFunctionName specifies if the call is an upgradeContract call
	UpgradeFunctionName = "upgradeContract"

	// DeleteFunctionName specifies if the call is an deleteContract call
	DeleteFunctionName = "deleteContract"
)

// CodeDeployInput contains code deploy state, whether it comes from a ContractCreateInput or a ContractCallInput
type CodeDeployInput struct {
	ContractCode         []byte
	ContractCodeMetadata []byte
	ContractAddress      []byte
	CodeDeployerAddress  []byte
}

// VMHostParameters represents the parameters to be passed to VMHost
type VMHostParameters struct {
	VMType                              []byte
	BlockGasLimit                       uint64
	GasSchedule                         config.GasScheduleMap
	BuiltInFuncContainer                vmcommon.BuiltInFunctionContainer
	ESDTTransferParser                  vmcommon.ESDTTransferParser
	ElrondProtectedKeyPrefix            []byte
	WasmerSIGSEGVPassthrough            bool
	EpochNotifier                       vmcommon.EpochNotifier
	EnableEpochsHandler                 vmcommon.EnableEpochsHandler
	TimeOutForSCExecutionInMilliseconds uint32
}

// AsyncCallInfo contains the information required to handle the asynchronous call of another SmartContract
type AsyncCallInfo struct {
	Destination []byte
	Data        []byte
	GasLimit    uint64
	GasLocked   uint64
	ValueBytes  []byte
}

// GetDestination returns the destination of an async call
func (aci *AsyncCallInfo) GetDestination() []byte {
	return aci.Destination
}

// GetData returns the transaction data of the async call
func (aci *AsyncCallInfo) GetData() []byte {
	return aci.Data
}

// GetGasLimit returns the gas limit of the current async call
func (aci *AsyncCallInfo) GetGasLimit() uint64 {
	return aci.GasLimit
}

// GetGasLocked returns the gas locked for the async callback
func (aci *AsyncCallInfo) GetGasLocked() uint64 {
	return aci.GasLocked
}

// GetValueBytes returns the byte representation of the value of the async call
func (aci *AsyncCallInfo) GetValueBytes() []byte {
	return aci.ValueBytes
}

// AsyncGeneratedCall holds the information abount an async call
type AsyncGeneratedCall struct {
	Status          AsyncCallStatus
	Destination     []byte
	Data            []byte
	GasLimit        uint64
	ValueBytes      []byte
	SuccessCallback string
	ErrorCallback   string
	ProvidedGas     uint64
}

// OldAsyncContext is a structure containing a group of async calls and a callback
//  that should be called when all these async calls are resolved
type OldAsyncContext struct {
	Callback   string
	AsyncCalls []*AsyncGeneratedCall
}

// GetDestination returns the destination of an async call
func (ac *AsyncGeneratedCall) GetDestination() []byte {
	return ac.Destination
}

// GetData returns the transaction data of the async call
func (ac *AsyncGeneratedCall) GetData() []byte {
	return ac.Data
}

// GetGasLimit returns the gas limit of the current async call
func (ac *AsyncGeneratedCall) GetGasLimit() uint64 {
	return ac.GasLimit
}

// GetGasLocked returns the gas locked for the async callback
func (ac *AsyncGeneratedCall) GetGasLocked() uint64 {
	return 0
}

// GetValueBytes returns the byte representation of the value of the async call
func (ac *AsyncGeneratedCall) GetValueBytes() []byte {
	return ac.ValueBytes
}

// IsInterfaceNil returns true if there is no value under the interface
func (ac *AsyncGeneratedCall) IsInterfaceNil() bool {
	return ac == nil
}

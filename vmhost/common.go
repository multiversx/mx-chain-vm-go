// Package vm contains the top-level components and definitions of the VM
package vmhost

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-v1_4-go/config"
)

// VMVersion returns the current vm version
const VMVersion = "v1.4"

// WASMPageSize size in bytes of a WASM Linear Memory Page
const WASMPageSize = 65536

// BreakpointValue encodes Wasmer runtime breakpoint types
type BreakpointValue uint64

const (
	// BreakpointNone means the lack of a breakpoint
	BreakpointNone BreakpointValue = iota

	// BreakpointExecutionFailed means that Wasmer must stop immediately
	// due to failure indicated by VM
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

// AsyncCallExecutionMode encodes the execution modes of an AsyncCall
type AsyncCallExecutionMode uint

const (
	// SyncCall indicates that the async call can be executed synchronously,
	// with its corresponding callback
	SyncCall AsyncCallExecutionMode = iota

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

// ProtectedStoragePrefix is the storage key prefix that will be protected by
// VM explicitly, and implicitly by the Elrond node due to '@'; the
// protection can be disabled temporarily by the StorageContext
const ProtectedStoragePrefix = "A"+"R"+"W"+"E"+"N@"

// TimeLockKeyPrefix is the storage key prefix used for timelock-related storage.
const TimeLockKeyPrefix = ProtectedStoragePrefix + "TIMELOCK"

// AsyncDataPrefix is the storage key prefix used for AsyncContext-related storage.
const AsyncDataPrefix = ProtectedStoragePrefix + "ASYNC"

// AsyncCallStatus represents the different status an async call can have
type AsyncCallStatus uint8

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
	ProtectedKeyPrefix                  []byte
	WasmerSIGSEGVPassthrough            bool
	EpochNotifier                       vmcommon.EpochNotifier
	EnableEpochsHandler                 vmcommon.EnableEpochsHandler
	Hasher                              HashComputer
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

// AsyncContext is a structure containing a group of async calls and a callback
//  that should be called when all these async calls are resolved
type AsyncContext struct {
	Callback   string
	AsyncCalls []*AsyncGeneratedCall
}

// AsyncContextInfo is the structure resulting after a smart contract call that has initiated
// one or more async calls. It will
type AsyncContextInfo struct {
	CallerAddr      []byte
	ReturnData      []byte
	AsyncContextMap map[string]*AsyncContext
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

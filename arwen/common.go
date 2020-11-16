package arwen

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

// CallbackFunctionName is the name of the default asynchronous callback
// function of a smart contract
const CallbackFunctionName = "callBack"

// TimeLockKeyPrevix is the storage key prefix used for timelock-related storage;
// not protected by Arwen, nor by the Elrond node
const TimeLockKeyPrefix = "timelock"

// AsyncDataPrefix is the storage key prefix used for AsyncContext-related
// storage; protected by Arwen explicitly, and implicitly by the Elrond node due to '@'
const AsyncDataPrefix = "ARWEN@"

type BreakpointValue uint64

const (
	BreakpointNone BreakpointValue = iota
	BreakpointExecutionFailed
	BreakpointAsyncCall
	BreakpointSignalError
	BreakpointOutOfGas
)

type AsyncCallExecutionMode uint

const (
	// SyncExecution indicates that the async call can be executed synchronously,
	// with its corresponding callback
	SyncExecution AsyncCallExecutionMode = iota

	// AsyncBuiltinFunc indicates that the async call is a cross-shard call to a
	// built-in function, which is executed half in-shard, half cross-shard
	AsyncBuiltinFunc

	// AsyncUnknown indicates that the async call cannot be executed locally, and
	// must be forwarded to the destination account
	AsyncUnknown
)

// AsyncCallStatus represents the different status an async call can have
type AsyncCallStatus uint8

const (
	AsyncCallPending AsyncCallStatus = iota
	AsyncCallResolved
	AsyncCallRejected
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
	VMType                   []byte
	BlockGasLimit            uint64
	GasSchedule              config.GasScheduleMap
	ProtocolBuiltinFunctions vmcommon.FunctionNames
	ElrondProtectedKeyPrefix []byte
	ArwenV2EnableEpoch       uint32
	AheadOfTimeEnableEpoch   uint32
	DynGasLockEnableEpoch    uint32
	UseWarmInstance          bool
}

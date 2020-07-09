package arwen

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

const TimeLockKeyPrefix = "timelock"
const CallbackDefault = "callBack"
const AsyncDataPrefix = "asyncCalls"

type BreakpointValue uint64

const (
	BreakpointNone BreakpointValue = iota
	BreakpointExecutionFailed
	BreakpointAsyncCall
	BreakpointSignalError
	BreakpointSignalExit
	BreakpointOutOfGas
)

type AsyncCallExecutionMode uint

const (
	// SyncCall indicates that the async call can be executed synchronously, with
	// its corresponding callback
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
}

// VMHostParameters represents the parameters to be passed to VMHost
type VMHostParameters struct {
	VMType                   []byte
	BlockGasLimit            uint64
	GasSchedule              config.GasScheduleMap
	ProtocolBuiltinFunctions vmcommon.FunctionNames
	ElrondProtectedKeyPrefix []byte
}

package arwen

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

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
	SyncCall AsyncCallExecutionMode = iota
	AsyncBuiltinFunc
	AsyncUnknown
)

const CallbackFunctionName = "callBack"
const TimeLockKeyPrefix = "timelock"
const AsyncDataPrefix = "asyncCalls"

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
	UseWarmInstance          bool
}

// AsyncCallInfo contains the information required to handle the asynchronous call of another SmartContract
type AsyncCallInfo struct {
	Destination []byte
	Data        []byte
	GasLimit    uint64
	GasLocked   uint64
	ValueBytes  []byte
}

func (aci *AsyncCallInfo) GetDestination() []byte {
	return aci.Destination
}

func (aci *AsyncCallInfo) GetData() []byte {
	return aci.Data
}

func (aci *AsyncCallInfo) GetGasLimit() uint64 {
	return aci.GasLimit
}

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

// GetValueBytes returns the byte representation of the value of the async call
func (ac *AsyncGeneratedCall) GetValueBytes() []byte {
	return ac.ValueBytes
}

// IsInterfaceNil returns true if there is no value under the interface
func (ac *AsyncGeneratedCall) IsInterfaceNil() bool {
	return ac == nil
}

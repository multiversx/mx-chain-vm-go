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
	BreakpointSignalExit
	BreakpointOutOfGas
)

const CallbackDefault = "callBack"
const TimeLockKeyPrefix = "timelock"
const AsyncDataPrefix = "asyncCalls"

// AsyncCallInfo contains the information required to handle the asynchronous call of another SmartContract
type AsyncCallInfo struct {
	Destination []byte
	Data        []byte
	GasLimit    uint64
	ValueBytes  []byte
}

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



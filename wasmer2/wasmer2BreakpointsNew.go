package wasmer2

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// The breakpoint values this backend writes into an instance's breakpoint global.
//
// They are the wire values of the 2.2 C API, and are the same numbers as the
// vmhost.BreakpointValue constants - which is where they are documented, and which this
// package cannot import, vmhost/hostCore importing wasmer2. Repeated here rather than
// shared, because under the new interfaces the legacy breakpoint is the backend's private
// mechanism for stopping a contract: executor.BreakpointValueNew, which is what the VM sees,
// has only the three values a middleware can raise.
const (
	legacyBreakpointNone uint64 = iota
	legacyBreakpointExecutionFailed
	legacyBreakpointAsyncCall
	legacyBreakpointSignalError
	legacyBreakpointOutOfGas
	legacyBreakpointMemoryLimit
)

// legacyAsyncCallCode is the early exit of a legacy async call. Code 0 is reserved for it,
// which is why executor.EarlyExit never carries vmcommon.Ok.
const legacyAsyncCallCode = uint64(vmcommon.Ok)

// legacyBreakpointForEarlyExit says which breakpoint halts the instance for an early exit.
//
// Counterpart of the Rust VMHooksEarlyExit::breakpoint_legacy, and the same switch as
// applyEarlyExit in vmhost/vmhooks/earlyExit.go: a hook that fails has to stop the contract,
// and under 2.2 the only way to do that is to trip the global.
//
// Everything else is a failure, which is also what an unrecognised code has to be: the
// instance must stop either way, and stopping it as a failure is the safe reading.
func legacyBreakpointForEarlyExit(code uint64) uint64 {
	switch code {
	case legacyAsyncCallCode:
		return legacyBreakpointAsyncCall
	case uint64(vmcommon.UserError):
		return legacyBreakpointSignalError
	case uint64(vmcommon.OutOfGas):
		return legacyBreakpointOutOfGas
	default:
		return legacyBreakpointExecutionFailed
	}
}

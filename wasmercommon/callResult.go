package wasmercommon

// CallResultKind says how a call ended.
//
// This is the enum that replaces the legacy arrangement, where a call returned a bare
// ok/error and the VM then read a wasm global to find out what had actually happened. The
// reason now comes back with the result.
type CallResultKind int

const (
	// CallOk means the function returned normally.
	CallOk CallResultKind = iota

	// CallFunctionNotFound means there is no exported function by that name.
	CallFunctionNotFound

	// CallRuntimeError means the instance trapped for a reason that is neither a VM hook's
	// doing nor a middleware breakpoint.
	CallRuntimeError

	// CallEarlyExit means a VM hook asked for execution to stop.
	CallEarlyExit

	// CallBreakpoint means a middleware stopped execution: out of gas, or the memory limit.
	CallBreakpoint
)

// CallResult is the outcome of a call.
type CallResult struct {
	Kind CallResultKind

	// EarlyExitCode is meaningful for CallEarlyExit. It is the vmcommon.ReturnCode the hook
	// exited with.
	EarlyExitCode uint64

	// Breakpoint is meaningful for CallBreakpoint, and matches executor.BreakpointValue.
	Breakpoint uint64

	// Message is informational, and carries the trap text for CallRuntimeError.
	//
	// For CallEarlyExit it is empty when the hook was a Go one, by construction: the hook
	// recorded its own message in the VM's output state before returning, and only the code
	// came back through the boundary. That is why the text does not have to travel.
	Message string
}

// HasMessage reports whether the error buffer holds text for this outcome.
//
// For the other kinds it holds whatever the previous call left behind, so it must not be read.
func (kind CallResultKind) HasMessage() bool {
	return kind == CallRuntimeError || kind == CallEarlyExit
}

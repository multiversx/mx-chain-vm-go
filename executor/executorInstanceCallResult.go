package executor

// InstanceCallResultKind says why a call stopped.
//
// It replaces the older arrangement in which a call reported only an error and the VM then
// polled the instance for a breakpoint to find out what had actually happened. A backend that
// knows why it stopped says so here; see PLAN-go-breakpoint-removal.md.
type InstanceCallResultKind int

const (
	// CallOk means the function returned normally.
	CallOk InstanceCallResultKind = iota

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

// InstanceCallResult is the outcome of calling an exported function.
type InstanceCallResult struct {
	Kind InstanceCallResultKind

	// Err carries the original error for CallRuntimeError, rather than only its text, so that
	// errors.Is and errors.As keep working across the boundary. A backend that can only
	// report a trap string wraps it here.
	Err error

	// EarlyExitCode is meaningful for CallEarlyExit: the vmcommon.ReturnCode the hook exited
	// with.
	EarlyExitCode uint64

	// Breakpoint is meaningful for CallBreakpoint, and matches BreakpointValue.
	Breakpoint uint64
}

// CallResultOk reports a call that returned normally.
func CallResultOk() InstanceCallResult {
	return InstanceCallResult{Kind: CallOk}
}

// CallResultError reports a call that trapped, keeping the error itself.
//
// This is what a backend whose C API reports only success or failure can say: it cannot tell
// a trap from a breakpoint, so the VM keeps reading the breakpoint global for those.
func CallResultError(err error) InstanceCallResult {
	return InstanceCallResult{Kind: CallRuntimeError, Err: err}
}

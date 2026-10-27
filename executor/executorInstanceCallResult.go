package executor

// InstanceCallResultKind says why a call stopped.
//
// Counterpart of the Rust `InstanceCallResult` enum (new_traits/instance_new.rs). It
// replaces the older arrangement in which a call reported only an error and the VM then
// polled the instance for a breakpoint to find out what had actually happened: a backend
// that knows why it stopped says so here.
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
//
// A struct rather than a sum type, Go having none: Kind selects which of the other fields is
// meaningful, and the rest are zero.
type InstanceCallResult struct {
	Kind InstanceCallResultKind

	// Err carries the original error, rather than only its text, so that errors.Is and
	// errors.As keep working across the boundary.
	//
	// For CallRuntimeError it is the trap. For CallEarlyExit it is the hook's own error,
	// which the VM resolves into a *EarlyExit; a backend that transports the exit
	// across an FFI boundary carrying only a code leaves it nil and fills EarlyExitCode
	// instead.
	Err error

	// EarlyExitCode is meaningful for CallEarlyExit: the vmcommon.ReturnCode the hook
	// exited with.
	EarlyExitCode uint64

	// Breakpoint is meaningful for CallBreakpoint.
	Breakpoint BreakpointValueNew
}

// CallResultOk reports a call that returned normally.
func CallResultOk() InstanceCallResult {
	return InstanceCallResult{Kind: CallOk}
}

// CallResultFunctionNotFound reports a call to a function the module does not export.
func CallResultFunctionNotFound() InstanceCallResult {
	return InstanceCallResult{Kind: CallFunctionNotFound}
}

// CallResultRuntimeError reports a call that trapped, keeping the error itself.
func CallResultRuntimeError(err error) InstanceCallResult {
	return InstanceCallResult{Kind: CallRuntimeError, Err: err}
}

// CallResultEarlyExit reports a call stopped by a VM hook.
func CallResultEarlyExit(err error, code uint64) InstanceCallResult {
	return InstanceCallResult{Kind: CallEarlyExit, Err: err, EarlyExitCode: code}
}

// CallResultBreakpoint reports a call stopped by one of the executor's own middlewares.
func CallResultBreakpoint(breakpoint BreakpointValueNew) InstanceCallResult {
	return InstanceCallResult{Kind: CallBreakpoint, Breakpoint: breakpoint}
}

package hostCore

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

// handleCallResult turns the executor's account of a finished call into the VM's.
//
// This is what reading the breakpoint global after a call turns into. The question is the
// same one - the call ended, why? - but it is now answered by the side that knows, in the
// value the call returned, rather than by polling a global the executor happened to leave
// behind.
//
// Every error produced here is also recorded on the runtime's error list, exactly as
// handleBreakpointIfAny did, so that the VMOutput's error trace is unchanged.
func (host *vmHost) handleCallResult(result executor.InstanceCallResult) error {
	if result.Kind == executor.CallOk {
		return nil
	}

	runtime := host.Runtime()

	// A trap that is nobody's deliberate exit. The underlying error is what goes on the
	// trace - it is the only description of what happened - while the caller is told the
	// generic failure, which is what the breakpoint-less path always did.
	if result.Kind == executor.CallRuntimeError {
		log.Trace("wasmer execution error", "err", result.Err)
		runtime.AddError(result.Err, runtime.FunctionName())
		return vmhost.ErrExecutionFailed
	}

	err := host.resolveCallResult(result)
	runtime.AddError(err, runtime.FunctionName())

	return err
}

// resolveCallResult names the outcome of a call that did not simply return.
func (host *vmHost) resolveCallResult(result executor.InstanceCallResult) error {
	switch result.Kind {
	case executor.CallFunctionNotFound:
		return executor.ErrFuncNotFound

	case executor.CallEarlyExit:
		return host.resolveEarlyExit(result.EarlyExitCode)

	case executor.CallBreakpoint:
		return resolveBreakpoint(result.Breakpoint)

	case executor.CallInterrupted:
		// The host asked for this call to stop, which only the execution deadline does. The
		// watchdog replaces this with ErrExecutionFailedWithTimeout once it regains
		// control; what matters here is that the execution does not continue.
		return vmhost.ErrExecutionFailed

	default:
		return vmhost.ErrUnhandledRuntimeBreakpoint
	}
}

// resolveEarlyExit names the outcome of a call a VM hook stopped.
//
// The code is a vmcommon.ReturnCode as it crossed the executor boundary. vmcommon.Ok is not
// a return code here: no ordinary exit uses it, which is what leaves it free to mean that
// the contract suspended itself for a legacy async call.
func (host *vmHost) resolveEarlyExit(code uint64) error {
	switch vmcommon.ReturnCode(code) {
	case vmcommon.Ok:
		return host.handleAsyncCallBreakpoint()
	case vmcommon.UserError:
		return vmhost.ErrSignalError
	case vmcommon.OutOfGas:
		return vmhost.ErrNotEnoughGas
	case vmcommon.ExecutionFailed:
		return vmhost.ErrExecutionFailed
	default:
		return vmhost.ErrUnhandledRuntimeBreakpoint
	}
}

// resolveBreakpoint names the outcome of a call one of the executor's own middlewares
// stopped.
func resolveBreakpoint(breakpoint executor.BreakpointValueNew) error {
	switch breakpoint {
	case executor.BreakpointOutOfGasNew:
		return vmhost.ErrNotEnoughGas
	case executor.BreakpointMemoryLimitNew:
		return vmhost.ErrMemoryLimit
	default:
		return vmhost.ErrUnhandledRuntimeBreakpoint
	}
}

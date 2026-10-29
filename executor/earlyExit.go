package executor

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// defaultEarlyExitCode is what an error that does not describe its own early exit means.
//
// Named rather than inlined into ExitExecutionFailed, because an executor that learns to
// report early exits of its own will have to agree with it.
const defaultEarlyExitCode = vmcommon.ExecutionFailed

// EarlyExit is an error that describes how execution must be terminated from inside a VM
// hook.
//
// Counterpart of the Rust `VMHooksEarlyExit` (new_traits/vm_hooks_early_exit.rs). Hooks only
// need to build one for the cases that deviate from the default failure behaviour; for
// everything else a plain error is enough, and the VM supplies the default.
//
// It lives here rather than in the vmhost package because it is what crosses the boundary: an
// executor that reports early exits to the VM reads Code off it. The other fields are the
// VM's own business, and no executor touches them.
type EarlyExit struct {
	// Code selects which termination routine the VM applies, and is also the code reported
	// to the executor. Only vmcommon.ExecutionFailed, vmcommon.UserError and
	// vmcommon.OutOfGas are meaningful today.
	//
	// Never vmcommon.Ok, which is what leaves code 0 free to mean "legacy async call" on
	// the wire.
	Code vmcommon.ReturnCode

	// Message is the VMOutput return message. For a user error this is the contract's own
	// text, which is not the same as Cause.Error().
	Message string

	// ConsumeAllGas burns all remaining gas before the termination routine runs.
	ConsumeAllGas bool

	// Cause is the underlying error, used for the error trace and errors.Is.
	Cause error
}

// Error implements the error interface.
func (e *EarlyExit) Error() string { return e.Message }

// Unwrap keeps errors.Is/errors.As working through an EarlyExit.
func (e *EarlyExit) Unwrap() error { return e.Cause }

// ExitExecutionFailed is the early exit for an error that does not describe one itself.
//
// This is the default the overwhelming majority of hook failures take, and the only place it
// is defined. The constructors for the two exits that deviate from it, ExitUserError and
// ExitOutOfGas, stay in the vmhost package with the sentinel errors they carry.
func ExitExecutionFailed(err error) *EarlyExit {
	return &EarlyExit{
		Code:          defaultEarlyExitCode,
		Message:       err.Error(),
		ConsumeAllGas: true,
		Cause:         err,
	}
}

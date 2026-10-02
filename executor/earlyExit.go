package executor

import (
	"errors"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// defaultEarlyExitCode is what an error that does not describe its own early exit means.
//
// Named rather than repeated, because both the struct below and EarlyExitCodeOf have to agree
// on it, and they are reached by different paths: the VM builds the full exit to apply its
// side effects, while an executor only ever reads the code.
const defaultEarlyExitCode = vmcommon.ExecutionFailed

// EarlyExit is an error that describes how execution must be terminated from inside a VM
// hook.
//
// Hooks only need to build one for the cases that deviate from the default failure behaviour;
// for everything else a plain error is enough, and the VM supplies the default. It lives here
// rather than in the VM because it is what crosses the boundary: an executor that reports
// early exits through its own call result reads Code off it. The other fields are the VM's
// own business, and no executor touches them.
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
// is defined.
func ExitExecutionFailed(err error) *EarlyExit {
	return &EarlyExit{
		Code:          defaultEarlyExitCode,
		Message:       err.Error(),
		ConsumeAllGas: true,
		Cause:         err,
	}
}

// EarlyExitCodeOf maps a VM hook's error to the code that travels to the executor.
//
// errors.As rather than a type assertion, so the code survives the fmt.Errorf("%w") wrapping
// that hook helpers do on the way out.
func EarlyExitCodeOf(err error) uint64 {
	var earlyExit *EarlyExit
	if errors.As(err, &earlyExit) {
		return uint64(earlyExit.Code)
	}

	return uint64(defaultEarlyExitCode)
}

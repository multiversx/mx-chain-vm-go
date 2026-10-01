package vmhost

import vmcommon "github.com/multiversx/mx-chain-vm-common-go"

// EarlyExit is an error that describes how execution must be terminated from
// inside a VM hook. Hooks only need it for the cases that deviate from the
// default failure behaviour; for everything else a plain error is enough and
// resolveEarlyExit supplies the default.
type EarlyExit struct {
	// Code selects which termination routine is applied. Only
	// vmcommon.ExecutionFailed, vmcommon.UserError and vmcommon.OutOfGas are
	// meaningful today.
	Code vmcommon.ReturnCode

	// Message is the VMOutput return message. For a user error this is the
	// contract's own text, which is not the same as Cause.Error().
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

// ExitUserError terminates execution with a user error carrying the contract's
// own message. It does not consume gas, matching runtime.SignalUserError.
func ExitUserError(message string) *EarlyExit {
	return &EarlyExit{
		Code:          vmcommon.UserError,
		Message:       message,
		ConsumeAllGas: false,
		Cause:         ErrSignalError,
	}
}

// ExitOutOfGas terminates execution with an out-of-gas breakpoint only, without
// touching the output state and without consuming the remaining gas. This is the
// legacy async-call behaviour.
func ExitOutOfGas() *EarlyExit {
	return &EarlyExit{
		Code:          vmcommon.OutOfGas,
		Message:       ErrNotEnoughGas.Error(),
		ConsumeAllGas: false,
		Cause:         ErrNotEnoughGas,
	}
}

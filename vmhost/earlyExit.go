package vmhost

import (
	"errors"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
)

// ErrAsyncCallPending signals that a contract suspended itself to let the host run a legacy
// async call. Not a failure: it is the cause carried by ExitAsyncCall, so that the exit can
// be recognised by errors.Is wherever it surfaces.
var errAsyncCallPending = errors.New("async call pending")

// The EarlyExit type itself lives in the executor package, because it is what crosses the
// boundary. The two constructors that deviate from the default failure behaviour stay here,
// with the sentinel errors they carry.

// ExitUserError terminates execution with a user error carrying the contract's
// own message. It does not consume gas, matching runtime.SignalUserError.
func ExitUserError(message string) *executor.EarlyExit {
	return &executor.EarlyExit{
		Code:          vmcommon.UserError,
		Message:       message,
		ConsumeAllGas: false,
		Cause:         ErrSignalError,
	}
}

// ExitAsyncCall suspends the contract so that the host can run the legacy async call it
// just registered.
//
// Code vmcommon.Ok is what makes this one distinguishable: the three ordinary exits never
// use it, so 0 on the wire means exactly "legacy async call". It is not a failure, so it
// records nothing in the VMOutput and consumes no gas - the whole of its effect is that the
// call ends and the host takes over.
//
// This is what BreakpointAsyncCall turns into. The breakpoint said the same thing by
// halting the instance from the outside; the hook now says it by returning.
func ExitAsyncCall() *executor.EarlyExit {
	return &executor.EarlyExit{
		Code:          vmcommon.Ok,
		Message:       "async call",
		ConsumeAllGas: false,
		Cause:         errAsyncCallPending,
	}
}

// errNotEnoughGas is the plain error ErrNotEnoughGas is built on, so that the exit has a
// cause distinct from itself.
var errNotEnoughGas = errors.New("not enough gas")

// ErrNotEnoughGas signals that there is not enough gas for the operation.
//
// It is an executor.EarlyExit, not a plain error, and that is the whole point: a hook that
// returns it says "the contract ran out of gas", and every classifier agrees - the VM's,
// which applies the side effects, and each backend's, which reports the call result - without
// anyone having to recognise a particular sentinel belonging to another package.
//
// This is what runtime.FailExecution used to do with its errors.Is(err, ErrNotEnoughGas)
// check: it chose the out-of-gas breakpoint over the execution-failed one, and that choice is
// what put ReturnCode OutOfGas in the VMOutput. The breakpoint is gone, and the choice now
// travels on the error itself, which is the only thing both sides of the executor boundary
// see.
//
// Consumes all remaining gas, as the execution-failed default it replaces did. At most call
// sites this is already true: UseGasBounded burns whatever is left before reporting failure.
//
// Still works with errors.Is and %w wrapping exactly as a sentinel does, by identity.
var ErrNotEnoughGas = &executor.EarlyExit{
	Code:          vmcommon.OutOfGas,
	Message:       errNotEnoughGas.Error(),
	ConsumeAllGas: true,
	Cause:         errNotEnoughGas,
}

// ExitOutOfGas terminates execution with an out-of-gas breakpoint only, without
// touching the output state and without consuming the remaining gas. This is the
// legacy async-call behaviour.
func ExitOutOfGas() *executor.EarlyExit {
	return &executor.EarlyExit{
		Code:          vmcommon.OutOfGas,
		Message:       ErrNotEnoughGas.Error(),
		ConsumeAllGas: false,
		Cause:         ErrNotEnoughGas,
	}
}

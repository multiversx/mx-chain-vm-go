package vmhost

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
)

// The EarlyExit type itself lives in the executor package, because it is what crosses the
// boundary. The constructors stay here, with the sentinel errors they carry.

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

package executor

import (
	"errors"
	"fmt"
)

// ErrInvalidFunction signals that the function is invalid
var ErrInvalidFunction = errors.New("invalid function")

// ErrFunctionNonvoidSignature signals that the signature for the function is invalid
var ErrFunctionNonvoidSignature = fmt.Errorf("%w (nonvoid signature)", ErrInvalidFunction)

// ErrFuncNotFound signals that the the function does not exist
var ErrFuncNotFound = fmt.Errorf("%w (not found)", ErrInvalidFunction)

// ErrNilOpcodeCosts signals that an executor was configured without a WASM opcode cost table.
//
// Rejected rather than defaulted, because an absent table is indistinguishable from one in
// which every opcode costs zero, and that is a contract running for free. OpcodeVersion needs
// no such check: its zero value is OpcodeVersionV1, a real choice.
var ErrNilOpcodeCosts = errors.New("nil WASM opcode costs")

// ErrVMHooksCannotHandleErrors signals that the VM hooks given to a backend (or to a wrapper
// of one) cannot absorb an error returned by a hook. See VMHooksErrorHandler for why
// such backends need them to.
//
// Transitional code: it is pending deletion in a future PR.
var ErrVMHooksCannotHandleErrors = errors.New("vm hooks do not implement VMHooksErrorHandler")

// ErrMemoryBadBounds signals that a certain variable is out of bounds
var ErrMemoryBadBounds = errors.New("bad bounds")

// ErrMemoryBadBoundsLower signals that a certain variable is lower than allowed
var ErrMemoryBadBoundsLower = fmt.Errorf("%w (lower)", ErrMemoryBadBounds)

// ErrMemoryBadBoundsUpper signals that a certain variable is higher than allowed
var ErrMemoryBadBoundsUpper = fmt.Errorf("%w (upper)", ErrMemoryBadBounds)

// ErrMemoryNegativeLength signals that the given length is less than 0
var ErrMemoryNegativeLength = errors.New("negative length")

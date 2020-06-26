package arwen

import (
	"errors"
	"fmt"
)

var ErrFunctionRunError = errors.New("function run error")

var ErrReturnCodeNotOk = errors.New("return not is not ok")

var ErrInvalidCallOnReadOnlyMode = errors.New("operation not permitted in read only mode")

var ErrNotEnoughGas = errors.New("not enough gas")

var ErrUnhandledRuntimeBreakpoint = errors.New("unhandled runtime breakpoint")

var ErrSignalError = errors.New("error signalled by smartcontract")

var ErrExecutionFailed = errors.New("execution failed")

var ErrInvalidAPICall = errors.New("invalid API call")

var ErrBadBounds = errors.New("bad bounds")

var ErrBadLowerBounds = fmt.Errorf("%w (lower)", ErrBadBounds)

var ErrBadUpperBounds = fmt.Errorf("%w (upper)", ErrBadBounds)

var ErrNegativeLength = errors.New("negative length")

var ErrFailedTransfer = errors.New("failed transfer")

var ErrTransferInsufficientFunds = fmt.Errorf("%w (insufficient funds)", ErrFailedTransfer)

var ErrFailedTransferDuringAsyncCall = fmt.Errorf("%w (failed during async call)", ErrFailedTransfer)

var ErrTransferNegativeValue = fmt.Errorf("%w (negative value)", ErrFailedTransfer)

var ErrUpgradeFailed = errors.New("upgrade failed")

var ErrInvalidUpgradeArguments = fmt.Errorf("%w (invalid arguments)", ErrUpgradeFailed)

var ErrInvalidFunction = errors.New("invalid function")

var ErrInitFuncCalledInRun = fmt.Errorf("%w (calling init() directly is forbidden)", ErrInvalidFunction)

var ErrCallBackFuncCalledInRun = fmt.Errorf("%w (calling callBack() directly is forbidden)", ErrInvalidFunction)

var ErrCallBackFuncNotExpected = fmt.Errorf("%w (unexpected callback was received)", ErrInvalidFunction)

var ErrFuncNotFound = fmt.Errorf("%w (not found)", ErrInvalidFunction)

var ErrInvalidFunctionName = fmt.Errorf("%w (invalid name)", ErrInvalidFunction)

var ErrFunctionNonvoidSignature = fmt.Errorf("%w (nonvoid signature)", ErrInvalidFunction)

var ErrContractInvalid = fmt.Errorf("invalid contract code")

var ErrContractNotFound = fmt.Errorf("%w (not found)", ErrContractInvalid)

var ErrMemoryDeclarationMissing = fmt.Errorf("%w (missing memory declaration)", ErrContractInvalid)

var ErrMaxInstancesReached = fmt.Errorf("%w (max instances reached)", ErrExecutionFailed)

var ErrStoreElrondReservedKey = errors.New("cannot write to storage under Elrond reserved key")

var ErrArgIndexOutOfRange = errors.New("argument index out of range")

var ErrArgOutOfRange = errors.New("argument out of range")

var ErrDivZero = errors.New("division by 0")

var ErrBitwiseNegative = errors.New("bitwise operations only allowed on positive integers")

var ErrShiftNegative = errors.New("bitwise shift operations only allowed on positive integers and by a positive amount")

var ErrAsyncCallGroupDoesNotExist = errors.New("async call group was not created yet")

var ErrAsync = errors.New("invalid gas percentage provided for async call")

var ErrInvalidAccount = errors.New("account does not exist")

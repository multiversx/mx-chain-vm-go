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

var ErrMemoryDeclarationMissing = errors.New("wasm memory declaration missing")

var ErrFailedTransfer = errors.New("failed transfer")

var ErrTransferInsufficientFunds = fmt.Errorf("%w (insufficient funds)", ErrFailedTransfer)

var ErrFailedTransferDuringAsyncCall = fmt.Errorf("%w (failed during async call)", ErrFailedTransfer)

var ErrTransferNegativeValue = fmt.Errorf("%w (negative value)", ErrFailedTransfer)

var ErrUpgradeFailed = errors.New("upgrade failed")

var ErrInvalidUpgradeArguments = fmt.Errorf("%w (invalid arguments)", ErrUpgradeFailed)

var ErrInvalidFunction = errors.New("invalid function")

var ErrInitFuncCalledInRun = fmt.Errorf("%w (calling init() directly is forbidden)", ErrInvalidFunction)

var ErrCallBackFuncCalledInRun = fmt.Errorf("%w (calling callBack() directly is forbidden)", ErrInvalidFunction)

var ErrFuncNotFound = fmt.Errorf("%w (not found)", ErrInvalidFunction)

var ErrInvalidFunctionName = fmt.Errorf("%w (invalid name)", ErrInvalidFunction)

var ErrFunctionNonvoidSignature = fmt.Errorf("%w (nonvoid signature)", ErrInvalidFunction)

var ErrContractInvalid = fmt.Errorf("invalid contract code")

var ErrContractNotFound = fmt.Errorf("%w (not found)", ErrContractInvalid)

package arwen

import (
	"errors"
	"fmt"
)

var ErrInitFuncCalledInRun = errors.New("it is not allowed to call init in run")

var ErrCallBackFuncCalledInRun = errors.New("it is not allowed to call callBack in run")

var ErrFunctionRunError = errors.New("function run error")

var ErrFunctionNonvoidSignature = errors.New("function signature not void")

var ErrReturnCodeNotOk = errors.New("return not is not ok")

var ErrInvalidCallOnReadOnlyMode = errors.New("operation not permitted in read only mode")

var ErrNotEnoughGas = errors.New("not enough gas")

var ErrUnhandledRuntimeBreakpoint = errors.New("unhandled runtime breakpoint")

var ErrFuncNotFound = errors.New("function not found")

var ErrSignalError = errors.New("error signalled by smartcontract")

var ErrExecutionFailed = errors.New("execution failed")

var ErrInvalidAPICall = errors.New("invalid API call")

var ErrBadBounds = errors.New("bad bounds")

var ErrBadLowerBounds = fmt.Errorf("%w (lower)", ErrBadBounds)

var ErrBadUpperBounds = fmt.Errorf("%w (upper)", ErrBadBounds)

var ErrNegativeLength = errors.New("negative length")

var ErrMemoryDeclarationMissing = errors.New("wasm memory declaration missing")

var ErrInvalidFunctionName = errors.New("invalid function name")

var ErrFailedTransfer = errors.New("failed transfer")

var ErrFailedTransferDuringAsyncCall = errors.New("failed transfer during async call")

var ErrTransferInsufficientFunds = errors.New("insufficient funds for transfer")

var ErrTransferNegativeValue = errors.New("cannot transfer negative value")

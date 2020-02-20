package arwen

import "errors"

var ErrInitFuncCalledInRun = errors.New("it is not allowed to call init in run")

var ErrFunctionRunError = errors.New("function run error")

var ErrReturnCodeNotOk = errors.New("return not is not ok")

var ErrInvalidCallOnReadOnlyMode = errors.New("operation not permitted in read only mode")

var ErrNotEnoughGas = errors.New("not enough gas")

var ErrUnhandledRuntimeBreakpoint = errors.New("unhandled runtime breakpoint")

var StateStackUnderflow = errors.New("state stack underflow")

var InstanceStackUnderflow = errors.New("instance stack underflow")

var ErrFuncNotFound = errors.New("function not found")

var ErrSignalError = errors.New("error signalled by smartcontract")

var ErrExecutionFailed = errors.New("execution failed")

var ErrInvalidAPICall = errors.New("invalid API call")

var ErrMemLoadBadBounds = errors.New("LoadBytes: bad bounds")

var ErrMemLoadNegativeLength = errors.New("LoadBytes: negative length")

var ErrMemStoreBadLowerBounds = errors.New("StoreBytes: bad lower bounds")

var ErrMemStoreBadUpperBounds = errors.New("StoreBytes: bad upper bounds")

var ErrMemoryDeclarationMissing = errors.New("wasm memory declaration missing")

var ErrInvalidFunctionName = errors.New("invalid function name")

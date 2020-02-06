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

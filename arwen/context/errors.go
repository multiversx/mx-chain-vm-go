package context

import "errors"

var ErrInitFuncCalledInRun = errors.New("it is not allowed to call init in run")

var ErrInvalidCallOnReadOnlyMode = errors.New("operation not permitted in read only mode")

var ErrFunctionRunError = errors.New("function run error")

var ErrFuncNotFound = errors.New("function not found")

var ErrReturnCodeNotOk = errors.New("return not is not ok")

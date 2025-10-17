package evm

import "errors"

var ErrActionNotSupported = errors.New("action not supported on EVM")

var ErrCodeNotCompiled = errors.New("code is not compiled")

var ErrExecutionAborted = errors.New("execution aborted")

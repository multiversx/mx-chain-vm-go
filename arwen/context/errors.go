package context

import "errors"

var ErrInitFuncCalledInRun = errors.New("it is not allowed to call init in run")

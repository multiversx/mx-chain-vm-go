package executorFactory

import (
	"errors"
)

// ErrNilCreateExecutorFactory signals that a nil create executor factory has been provided
var ErrNilCreateExecutorFactory = errors.New("nil create executor factory has been provided")

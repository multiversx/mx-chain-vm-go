package contexts

import (
	"errors"
)

// ErrNilBlockchainContextFactory signals that a nil blockchain context factory has been provided
var ErrNilBlockchainContextFactory = errors.New("nil blockchain context factory has been provided")

// ErrNilRuntimeContextFactory signals that a nil runtime context factory has been provided
var ErrNilRuntimeContextFactory = errors.New("nil runtime context factory has been provided")

// ErrNilMeteringContextFactory signals that a nil metering context factory has been provided
var ErrNilMeteringContextFactory = errors.New("nil metering context factory has been provided")

// ErrNilOutputContextFactory signals that a nil output context factory has been provided
var ErrNilOutputContextFactory = errors.New("nil output context factory has been provided")

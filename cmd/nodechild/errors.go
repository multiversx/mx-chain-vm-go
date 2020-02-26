package main

import (
	"errors"
	"fmt"
)

// ErrCriticalError signals a critical error
var ErrCriticalError = errors.New("critical error")

// ErrBadCommandFromNode signals a bad command from node
var ErrBadCommandFromNode = fmt.Errorf("%w: bad command from node", ErrCriticalError)

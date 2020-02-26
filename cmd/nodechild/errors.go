package main

import (
	"errors"
	"fmt"
)

// ErrCriticalError signals a critical error
var ErrCriticalError = errors.New("critical error")

// ErrBadRequestFromNode signals a critical error
var ErrBadRequestFromNode = fmt.Errorf("%w: bad request from node", ErrCriticalError)

// ErrCannotSendContractRequest signals a critical error
var ErrCannotSendContractRequest = fmt.Errorf("%w: cannot send contract request", ErrCriticalError)

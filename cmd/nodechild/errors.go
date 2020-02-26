package main

import (
	"fmt"
)

// ErrCriticalError signals a critical error
var ErrCriticalError = fmt.Errorf("critical error")

// ErrBadResponseTag signals a critical error
var ErrBadResponseTag = fmt.Errorf("%w: bad response tag", ErrCriticalError)

// ErrBadRequestFromNode signals a critical error
var ErrBadRequestFromNode = fmt.Errorf("%w: bad request from node", ErrCriticalError)

// ErrCannotSendContractRequest signals a critical error
var ErrCannotSendContractRequest = fmt.Errorf("%w: cannot send contract request", ErrCriticalError)

// ErrCannotSendHookCallRequest signals a critical error
var ErrCannotSendHookCallRequest = fmt.Errorf("%w: cannot send hook call request", ErrCriticalError)

// ErrCannotReceiveHookCallResponse signals a critical error
var ErrCannotReceiveHookCallResponse = fmt.Errorf("%w: cannot receive hook call response", ErrCriticalError)

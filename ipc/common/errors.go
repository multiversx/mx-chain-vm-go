package common

import (
	"fmt"
)

// ErrCriticalError signals a critical error
var ErrCriticalError = fmt.Errorf("critical error")

// ErrInvalidMessageNonce signals a critical error
var ErrInvalidMessageNonce = fmt.Errorf("%w: invalid nonce in message", ErrCriticalError)

// ErrStopPerNodeRequest signals a critical error
var ErrStopPerNodeRequest = fmt.Errorf("%w: will stop, as node requested", ErrCriticalError)

// ErrBadRequestFromNode signals a critical error
var ErrBadRequestFromNode = fmt.Errorf("%w: bad request from node", ErrCriticalError)

// ErrBadMessageFromArwen signals a critical error
var ErrBadMessageFromArwen = fmt.Errorf("%w: bad message from Arwen", ErrCriticalError)

// ErrCannotSendContractRequest signals a critical error
var ErrCannotSendContractRequest = fmt.Errorf("%w: cannot send contract request", ErrCriticalError)

// ErrCannotSendHookCallResponse signals a critical error
var ErrCannotSendHookCallResponse = fmt.Errorf("%w: cannot hook call response", ErrCriticalError)

// ErrCannotSendHookCallRequest signals a critical error
var ErrCannotSendHookCallRequest = fmt.Errorf("%w: cannot send hook call request", ErrCriticalError)

// ErrCannotReceiveHookCallResponse signals a critical error
var ErrCannotReceiveHookCallResponse = fmt.Errorf("%w: cannot receive hook call response", ErrCriticalError)

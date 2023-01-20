package arwendebug

import (
	"errors"
	"fmt"
)

// RequestError signals an error
type RequestError struct {
	Message  string
	InnerErr error
}

// NewRequestError -
func NewRequestError(message string) *RequestError {
	return &RequestError{Message: message}
}

// NewRequestErrorMessageInner -
func NewRequestErrorMessageInner(message string, err error) *RequestError {
	return &RequestError{Message: message, InnerErr: err}
}

func (err *RequestError) Error() string {
	if err.InnerErr == nil {
		return fmt.Sprintf("request error: message=%s", err.Message)
	}

	return fmt.Sprintf("request error: message=%s; inner=%v", err.Message, err.InnerErr)
}

// Unwrap unwraps the inner error
func (err *RequestError) Unwrap() error {
	return err.InnerErr
}

// ErrInvalidArgumentEncoding signals an error
var ErrInvalidArgumentEncoding = errors.New("invalid contract argument encoding")

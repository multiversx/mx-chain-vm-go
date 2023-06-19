package math

import (
	"fmt"
	builtinMath "math"
)

type overflowHandler struct {
	err error
}

// NewOverflowHandler creates a new instance of type overflowHandler
// This struct is not concurrent safe and should be instantiated for each chain of operations required to be executed
func NewOverflowHandler() *overflowHandler {
	return &overflowHandler{}
}

// AddInt64 will add the two provided values. If an overflow occurs, the error will be stored internally
func (handler *overflowHandler) AddInt64(a, b int64) int64 {
	if handler.err != nil {
		return builtinMath.MaxInt64
	}

	s := a + b
	if (s > a) == (b > 0) {
		return s
	}

	if handler.err == nil {
		handler.err = fmt.Errorf("%w when adding %d with %d", ErrAdditionOverflow, a, b)
	}

	return builtinMath.MaxInt64
}

// MulInt64 will multiply the two provided values. If an overflow occurs, the error will be stored internally
func (handler *overflowHandler) MulInt64(a, b int64) int64 {
	if handler.err != nil {
		return builtinMath.MaxInt64
	}

	res := a * b
	if a == 0 || b == 0 || a == res/b {
		return res
	}

	if handler.err == nil {
		handler.err = fmt.Errorf("%w when multiplying %d with %d", ErrMultiplicationOverflow, a, b)
	}

	return builtinMath.MaxInt64
}

// Error returns the stored error
func (handler *overflowHandler) Error() error {
	return handler.err
}

package executor

import (
	"errors"
	"fmt"
)

// ErrInvalidFunction signals that the function is invalid
var ErrInvalidFunction = errors.New("invalid function")

// ErrFunctionNonvoidSignature signals that the signature for the function is invalid
var ErrFunctionNonvoidSignature = fmt.Errorf("%w (nonvoid signature)", ErrInvalidFunction)

// ErrFuncNotFound signals that the the function does not exist
var ErrFuncNotFound = fmt.Errorf("%w (not found)", ErrInvalidFunction)

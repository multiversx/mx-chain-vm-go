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

// ErrMemoryBadBounds signals that a certain variable is out of bounds
var ErrMemoryBadBounds = errors.New("bad bounds")

// ErrMemoryBadBoundsLower signals that a certain variable is lower than allowed
var ErrMemoryBadBoundsLower = fmt.Errorf("%w (lower)", ErrMemoryBadBounds)

// ErrMemoryBadBoundsUpper signals that a certain variable is higher than allowed
var ErrMemoryBadBoundsUpper = fmt.Errorf("%w (upper)", ErrMemoryBadBounds)

// ErrMemoryNegativeLength signals that the given length is less than 0
var ErrMemoryNegativeLength = errors.New("negative length")

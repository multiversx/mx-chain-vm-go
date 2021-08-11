package math

import (
	"errors"
)

// ErrAdditionOverflow is raised when there is an overflow because of the addition of two numbers
var ErrAdditionOverflow = errors.New("addition overflow")

// ErrSubtractionUnderflow is raised when there is an underflow because of the subtraction of two numbers
var ErrSubtractionUnderflow = errors.New("subtraction underflow")

// ErrMultiplicationOverflow is raised when there is an overflow because of the multiplication of two numbers
var ErrMultiplicationOverflow = errors.New("multiplication overflow")

package math

import (
	"errors"
)

// ErrAdditionOverflow is raised when there is an overflow because of the addition of two numbers
var ErrAdditionOverflow = errors.New("addition overflow")

// ErrSubtractionOverflow is raised when there is an overflow because of the subtraction of two numbers
var ErrSubtractionOverflow = errors.New("subtraction overflow")

// ErrMultiplicationOverflow is raised when there is an overflow because of the multiplication of two numbers
var ErrMultiplicationOverflow = errors.New("multiplication overflow")

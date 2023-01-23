package math

import (
	"errors"
)

// ErrAdditionOverflow is raised when there is an overflow because of the addition of two numbers
var ErrAdditionOverflow = errors.New("addition overflow")

// ErrMultiplicationOverflow is raised when there is an overflow because of the multiplication of two numbers
var ErrMultiplicationOverflow = errors.New("multiplication overflow")

// ErrBigFloatSub is raised when sub of floats produces a panic
var ErrBigFloatSub = errors.New("this big Float operation is not permitted while doing float.Sub")

// ErrBigFloatAdd is raised when add of floats produces a panic
var ErrBigFloatAdd = errors.New("this big Float operation is not permitted while doing float.Add")

// ErrBigFloatQuo is raised when quo of floats produces a panic
var ErrBigFloatQuo = errors.New("this big Float operation is not permitted while doing float.Quo")

// ErrBigFloatMul is raised when mul of floats produces a panic
var ErrBigFloatMul = errors.New("this big Float operation is not permitted while doing float.Mul")

// ErrBigFloatSqrt is raised when sqrt of floats produces a panic
var ErrBigFloatSqrt = errors.New("this big Float operation is not permitted while doing float.Sqrt")

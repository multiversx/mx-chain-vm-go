package math

import (
	"math/big"
)

// SubBigFloat subtraction implementation with error handling for big float
func SubBigFloat(op1, op2 *big.Float) (result *big.Float, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrBigFloatSub
			result = big.NewFloat(0)
		}
	}()
	result = new(big.Float)
	result.Sub(op1, op2)
	return
}

// AddBigFloat addition implementation with error handling for big float
func AddBigFloat(op1, op2 *big.Float) (result *big.Float, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrBigFloatAdd
			result = big.NewFloat(0)
		}
	}()
	result = new(big.Float)
	result.Add(op1, op2)
	return
}

// QuoBigFloat quotient implementation with error handling for big float
func QuoBigFloat(op1, op2 *big.Float) (result *big.Float, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrBigFloatQuo
			result = big.NewFloat(0)
		}
	}()
	result = new(big.Float)
	result.Quo(op1, op2)
	return
}

// MulBigFloat multiplication implementation with error handling for big float
func MulBigFloat(op1, op2 *big.Float) (result *big.Float, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrBigFloatMul
			result = big.NewFloat(0)
		}
	}()
	result = new(big.Float)
	result.Mul(op1, op2)
	return
}

// SqrtBigFloat sqrt implementation with error handling for big float
func SqrtBigFloat(op *big.Float) (result *big.Float, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrBigFloatSqrt
			result = big.NewFloat(0)
		}
	}()
	result = new(big.Float)
	result.Sqrt(op)
	return
}

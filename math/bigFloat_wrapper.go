package math

import (
	"fmt"
	"math/big"
)

func SubBigFloat(op1, op2 *big.Float) (result *big.Float, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w while doing float.Sub", ErrOperationCausingPanic)
			result = big.NewFloat(0)
		}
	}()
	result = new(big.Float)
	result.Sub(op1, op2)
	return
}

func AddBigFloat(op1, op2 *big.Float) (result *big.Float, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w while doing float.Add", ErrOperationCausingPanic)
			result = big.NewFloat(0)
		}
	}()
	result = new(big.Float)
	result.Add(op1, op2)
	return
}

func QuoBigFloat(op1, op2 *big.Float) (result *big.Float, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w while doing float.Quo", ErrOperationCausingPanic)
			result = big.NewFloat(0)
		}
	}()
	result = new(big.Float)
	result.Quo(op1, op2)
	return
}

func MulBigFloat(op1, op2 *big.Float) (result *big.Float, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w while doing float.Mul", ErrOperationCausingPanic)
			result = big.NewFloat(0)
		}
	}()
	result = new(big.Float)
	result.Mul(op1, op2)
	return
}

func SqrtBigFloat(op *big.Float) (result *big.Float, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w while doing float.Sqrt", ErrOperationCausingPanic)
			result = big.NewFloat(0)
		}
	}()
	result = new(big.Float)
	result.Sqrt(op)
	return
}

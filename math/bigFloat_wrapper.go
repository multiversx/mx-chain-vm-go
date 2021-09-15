package math

import (
	"fmt"
	"math/big"
)

func Sub(result, op1, op2 *big.Float) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w while doing float.Sub", ErrOperationCausingPanic)
			result = big.NewFloat(0)
		}
	}()
	result.Sub(op1, op2)
	return err
}

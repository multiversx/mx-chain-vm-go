package math

import (
	builtinMath "math"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

var log = logger.GetOrCreate("arwen/overflow")

// AddUint64 performs addition on uint64 and logs an error if the addition overflows
func AddUint64(a, b uint64) uint64 {
	res, err := AddUint64WithErr(a, b)
	if err != nil {
		log.Error("AddUint64 overflow", a, b)
		return builtinMath.MaxUint64
	}

	return res
}

// AddUint64WithErr performs addition on uint64 and returns an error if the addition overflows
func AddUint64WithErr(a, b uint64) (uint64, error) {
	s := a + b
	if s >= a && s >= b {
		return s, nil
	}

	return s, ErrAdditionOverflow
}

// AddInt64 performs addition on int64 and returns an error if the addition overflows
func AddInt64(a, b int64) int64 {
	s := a + b
	if s >= a && s >= b {
		return s
	}

	return builtinMath.MaxInt64
}

// MulUint64 performs multiplication on uint64 and logs an error if the multiplication overflows
func MulUint64(a, b uint64) uint64 {
	res, err := MulUint64WithErr(a, b)
	if err != nil {
		log.Error("MulUint64 overflow", a, b)
		return builtinMath.MaxUint64
	}

	return res
}

// MulUint64WithErr performs multiplication on uint64 and returns an error if the multiplication overflows
func MulUint64WithErr(a, b uint64) (uint64, error) {
	res := a * b
	if a == 0 || b == 0 || a == res/b {
		return res, nil
	}

	return 0, ErrMultiplicationOverflow
}

// AddInt32 performs addition on int32 and returns an error if the addition overflows
func AddInt32(a, b int32) (int32, error) {
	res := a + b
	if a < 0 && b < 0 && res > 0 {
		return 0, ErrAdditionOverflow
	}
	if a > 0 && b > 0 && res < 0 {
		return 0, ErrAdditionOverflow
	}

	return res, nil
}

// SubInt performs subtraction on int and returns an error if the subtraction overflows
func SubInt(a, b int) int {
	res := a - b
	if res <= a {
		return res
	}

	return builtinMath.MinInt64
}

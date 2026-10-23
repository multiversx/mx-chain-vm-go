package math

import (
	builtinMath "math"

	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("vm/overflow")

// AddUint64 performs addition on uint64 and logs an error if the addition overflows
func AddUint64(a, b uint64) uint64 {
	res, err := AddUint64WithErr(a, b)
	if err != nil {
		log.Trace("AddUint64 overflow", "a", a, "b", b)
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

	return builtinMath.MaxUint64, ErrAdditionOverflow
}

// SubUint64 performs subtraction on uint64, in case of underflow returns 0
func SubUint64(a, b uint64) uint64 {
	res, err := SubUint64WithErr(a, b)
	if err != nil {
		return 0
	}

	return res
}

// SubUint64WithErr performs subtraction on uint64 and returns an error if the subtraction underflows
func SubUint64WithErr(a, b uint64) (uint64, error) {
	if a < b {
		return 0, ErrSubtractionUnderflow
	}

	return a - b, nil
}

// MulUint64 performs multiplication on uint64 and logs an error if the multiplication overflows
func MulUint64(a, b uint64) uint64 {
	res, err := MulUint64WithErr(a, b)
	if err != nil {
		log.Trace("MulUint64 overflow", "a", a, "b", b)
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

	return builtinMath.MaxUint64, ErrMultiplicationOverflow
}

// AddInt32 performs addition on int32 and logs an error if the addition overflows
func AddInt32(a, b int32) int32 {
	res, err := AddInt32WithErr(a, b)
	if err != nil {
		log.Trace("AddInt32 overflow", "a", a, "b", b)
		return builtinMath.MaxInt32
	}

	return res
}

// AddInt32WithErr performs addition on int32 and returns an error if the addition overflows
func AddInt32WithErr(a, b int32) (int32, error) {
	s := a + b
	if (s > a) == (b > 0) {
		return s, nil
	}

	return builtinMath.MaxInt32, ErrAdditionOverflow
}

// SubInt performs subtraction on int and logs an error if the subtraction overflows
func SubInt(a, b int) int {
	res := a - b
	if (res < a) == (b > 0) {
		return res
	}

	log.Trace("SubInt underflow", "a", a, "b", b)
	return builtinMath.MinInt64
}

// MulInt32 performs multiplication on int32 and logs an error if the multiplication overflows
func MulInt32(a, b int32) int32 {
	res, err := MulInt32WithErr(a, b)
	if err != nil {
		log.Trace("MulInt32 overflow", "a", a, "b", b)
		return builtinMath.MaxInt32
	}

	return res
}

// MulInt32WithErr performs multiplication on int32 and returns an error if the multiplication overflows
func MulInt32WithErr(a, b int32) (int32, error) {
	wide := int64(a) * int64(b)
	if wide > builtinMath.MaxInt32 || wide < builtinMath.MinInt32 {
		return builtinMath.MaxInt32, ErrMultiplicationOverflow
	}
	return int32(wide), nil
}

// SubInt32 performs subtraction on int32 and logs an error if the subtraction underflows
func SubInt32(a, b int32) int32 {
	res, err := SubInt32WithErr(a, b)
	if err != nil {
		log.Trace("SubInt32 underflow", "a", a, "b", b)
	}

	return res
}

// SubInt32WithErr performs subtraction on int32 and returns an error if the subtraction underflows or overflows
func SubInt32WithErr(a, b int32) (int32, error) {
	diff := int64(a) - int64(b)
	if diff < builtinMath.MinInt32 {
		return builtinMath.MinInt32, ErrSubtractionUnderflow
	}
	if diff > builtinMath.MaxInt32 {
		return builtinMath.MaxInt32, ErrSubtractionUnderflow
	}

	return int32(diff), nil
}

// CheckBoundsInt32 validates that start and length are non-negative and int64(start)+int64(length) <= int64(maxLen)
func CheckBoundsInt32(start, length int32, maxLen int) error {
	if start < 0 || length < 0 {
		return ErrOutOfBounds
	}
	if int64(start)+int64(length) > int64(maxLen) {
		return ErrOutOfBounds
	}

	return nil
}

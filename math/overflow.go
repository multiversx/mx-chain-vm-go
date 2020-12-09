package math

// AddUint64 performs addition on uint64 and returns an error if the addition overflows
func AddUint64(a, b uint64) (uint64, error) {
	s := a + b
	if s >= a && s >= b {
		return s, nil
	}

	return s, ErrAdditionOverflow
}

// MulUint64 performs multiplication on uint64 and returns an error if the multiplication overflows
func MulUint64(a, b uint64) (uint64, error) {
	res := a * b
	if a == 0 || b == 0 || a == res/b {
		return res, nil
	}

	return 0, ErrMultiplicationOverflow
}

// SubUint64 performs subtraction on uint64 and returns an error if the subtraction overflows
func SubUint64(a, b uint64) (uint64, error) {
	res := a - b
	if res <= a {
		return res, nil
	}

	return 0, ErrSubtractionOverflow
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

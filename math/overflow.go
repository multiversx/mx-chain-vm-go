package math

// AddUint64 performs addition on uint64 and returns an error if the addition overflows
func AddUint64(a, b uint64) (uint64, error) {
	s := a + b
	if s >= a && s >= b {
		return s, nil
	}

	return s, ErrAdditionOverflow
}

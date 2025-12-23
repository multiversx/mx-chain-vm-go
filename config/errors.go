package config

import (
	"errors"
)

// ErrNilGasScheduleFactory signals that a nil gas schedule factory has been provided
var ErrNilGasScheduleFactory = errors.New("nil gas schedule factory has been provided")

package factory

import (
	"errors"
)

var errNilRunTypeComponentsFactory = errors.New("nil run type components factory")

// ErrNilRunTypeComponents signals that nil runType components has been provided
var ErrNilRunTypeComponents = errors.New("nil run type components")

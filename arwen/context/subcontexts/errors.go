package subcontexts

import "errors"

var StateStackUnderflow = errors.New("State stack underflow")

var ErrInvalidCallOnReadOnlyMode = errors.New("operation not permitted in read only mode")

var ErrNotEnoughGas = errors.New("not enough gas")

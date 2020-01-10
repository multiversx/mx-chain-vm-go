package subcontexts

import "errors"

var StateStackUnderflow = errors.New("State stack underflow")

var InstanceStackUnderflow = errors.New("Instance stack underflow")

var ErrNotEnoughGas = errors.New("not enough gas")

var ErrFuncNotFound = errors.New("function not found")

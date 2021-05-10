package host

import (
	"errors"

	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/process"
)

type mockBuiltin struct {
	processBuiltinFunction func(acntSnd, _ state.UserAccountHandler, vmInput *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error)
	setNewGasConfig        func(_ *process.GasCost)
	isInterfaceNil         func() bool
}

// ProcessBuiltInFunction - see BuiltinFunction.ProcessBuiltInFunction()
func (m *mockBuiltin) ProcessBuiltinFunction(acntSnd, acntRcv state.UserAccountHandler, vmInput *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	if m.processBuiltinFunction == nil {
		return nil, errors.New("Undefined processBuiltinFunction")
	}
	return m.processBuiltinFunction(acntSnd, acntRcv, vmInput)
}

// SetNewGasConfig - see BuiltinFunction.SetNewGasConfig()
func (m *mockBuiltin) SetNewGasConfig(gasCost *process.GasCost) {
	if m.setNewGasConfig != nil {
		m.setNewGasConfig(gasCost)
	}
}

// IsInterfaceNil - see BuiltinFunction.IsInterfaceNil()
func (m *mockBuiltin) IsInterfaceNil() bool {
	if m.isInterfaceNil == nil {
		return m == nil
	}
	return m.isInterfaceNil()
}

package testcommon

import (
	"errors"

	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/process"
)

// MockBuiltin defined the functions that can be replaced in order to mock a builtin
type MockBuiltin struct {
	ProcessBuiltinFunctionCall func(acntSnd, _ state.UserAccountHandler, vmInput *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error)
	setNewGasConfig            func(_ *process.GasCost)
	isInterfaceNil             func() bool
}

// ProcessBuiltinFunction - see BuiltinFunction.ProcessBuiltInFunction()
func (m *MockBuiltin) ProcessBuiltinFunction(acntSnd, acntRcv state.UserAccountHandler, vmInput *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	if m.ProcessBuiltinFunctionCall == nil {
		return nil, errors.New("Undefined processBuiltinFunction")
	}
	return m.ProcessBuiltinFunctionCall(acntSnd, acntRcv, vmInput)
}

// SetNewGasConfig - see BuiltinFunction.SetNewGasConfig()
func (m *MockBuiltin) SetNewGasConfig(gasCost *process.GasCost) {
	if m.setNewGasConfig != nil {
		m.setNewGasConfig(gasCost)
	}
}

// IsInterfaceNil - see BuiltinFunction.IsInterfaceNil()
func (m *MockBuiltin) IsInterfaceNil() bool {
	if m.isInterfaceNil == nil {
		return m == nil
	}
	return m.isInterfaceNil()
}

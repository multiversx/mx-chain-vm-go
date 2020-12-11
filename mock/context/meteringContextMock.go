package mock

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
)

var _ arwen.MeteringContext = (*MeteringContextMock)(nil)

// MeteringContextMock is used in tests to check the MeteringContext interface method calls
type MeteringContextMock struct {
	GasCost           *config.GasCost
	GasLeftMock       uint64
	GasLockedMock     uint64
	GasComputedToLock uint64
	BlockGasLimitMock uint64
	Err               error
}

// SetGasSchedule mocked method
func (m *MeteringContextMock) SetGasSchedule(gasSchedule config.GasScheduleMap) {
	gasCostConfig, _ := config.CreateGasConfig(gasSchedule)
	m.GasCost = gasCostConfig
}

// GasSchedule mocked method
func (m *MeteringContextMock) GasSchedule() *config.GasCost {
	return m.GasCost
}

// UseGas mocked method
func (m *MeteringContextMock) UseGas(gas uint64) {
}

// FreeGas mocked method
func (m *MeteringContextMock) FreeGas(gas uint64) {
}

// RestoreGas mocked method
func (m *MeteringContextMock) RestoreGas(gas uint64) {
}

// GasLeft mocked method
func (m *MeteringContextMock) GasLeft() uint64 {
	return m.GasLeftMock
}

// BoundGasLimit mocked method
func (m *MeteringContextMock) BoundGasLimit(limit uint64) uint64 {
	gasLeft := m.GasLeft()

	if gasLeft < limit {
		return gasLeft
	}
	return limit
}

// ComputeGasLockedForAsync mocked method
func (m *MeteringContextMock) ComputeGasLockedForAsync() uint64 {
	return m.GasComputedToLock
}

// DeductGasIfAsyncStep mocked method
func (m *MeteringContextMock) DeductGasIfAsyncStep() error {
	return m.Err
}

// UseGasBounded mocked method
func (m *MeteringContextMock) UseGasBounded(gas uint64) error {
	return m.Err
}

// UnlockGasIfAsyncCallback mocked method
func (m *MeteringContextMock) UnlockGasIfAsyncCallback() {
}

// UseGasForAsyncStep mocked method
func (m *MeteringContextMock) UseGasForAsyncStep() error {
	return m.Err
}

// UnlockGasIfAsyncStep mocked method
func (m *MeteringContextMock) UnlockGasIfAsyncStep() {
}

// GetGasLocked mocked method
func (m *MeteringContextMock) GetGasLocked() uint64 {
	return m.GasLockedMock
}

// BlockGasLimit mocked method
func (m *MeteringContextMock) BlockGasLimit() uint64 {
	return m.BlockGasLimitMock
}

// DeductInitialGasForExecution mocked method
func (m *MeteringContextMock) DeductInitialGasForExecution(contract []byte) error {
	return m.Err
}

// DeductInitialGasForDirectDeployment mocked method
func (m *MeteringContextMock) DeductInitialGasForDirectDeployment(input arwen.CodeDeployInput) error {
	return m.Err
}

// DeductInitialGasForIndirectDeployment mocked method
func (m *MeteringContextMock) DeductInitialGasForIndirectDeployment(input arwen.CodeDeployInput) error {
	return m.Err
}

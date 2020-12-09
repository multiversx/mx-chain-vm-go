package mock

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
)

var _ arwen.MeteringContext = (*MeteringContextMock)(nil)

// MeteringContextMock -
type MeteringContextMock struct {
	GasCost           *config.GasCost
	GasLeftMock       uint64
	GasLockedMock     uint64
	GasComputedToLock uint64
	BlockGasLimitMock uint64
	Err               error
}

// SetGasSchedule -
func (m *MeteringContextMock) SetGasSchedule(gasSchedule config.GasScheduleMap) {
	gasCostConfig, _ := config.CreateGasConfig(gasSchedule)
	m.GasCost = gasCostConfig
}

// GasSchedule -
func (m *MeteringContextMock) GasSchedule() *config.GasCost {
	return m.GasCost
}

// UseGas -
func (m *MeteringContextMock) UseGas(_ uint64) {
}

// FreeGas -
func (m *MeteringContextMock) FreeGas(_ uint64) {
}

// RestoreGas -
func (m *MeteringContextMock) RestoreGas(_ uint64) {
}

// GasLeft -
func (m *MeteringContextMock) GasLeft() uint64 {
	return m.GasLeftMock
}

// BoundGasLimit -
func (m *MeteringContextMock) BoundGasLimit(value int64) uint64 {
	gasLeft := m.GasLeft()
	limit := uint64(value)

	if gasLeft < limit {
		return gasLeft
	}
	return limit
}

// ComputeGasLockedForAsync -
func (m *MeteringContextMock) ComputeGasLockedForAsync() uint64 {
	return m.GasComputedToLock
}

// DeductGasIfAsyncStep -
func (m *MeteringContextMock) DeductGasIfAsyncStep() error {
	return m.Err
}

// UseGasBounded -
func (m *MeteringContextMock) UseGasBounded(_ uint64) error {
	return m.Err
}

// UnlockGasIfAsyncCallback -
func (m *MeteringContextMock) UnlockGasIfAsyncCallback() {
}

// UseGasForAsyncStep -
func (m *MeteringContextMock) UseGasForAsyncStep() error {
	return m.Err
}

// UnlockGasIfAsyncStep -
func (m *MeteringContextMock) UnlockGasIfAsyncStep() {
}

// GetGasLocked -
func (m *MeteringContextMock) GetGasLocked() uint64 {
	return m.GasLockedMock
}

// BlockGasLimit -
func (m *MeteringContextMock) BlockGasLimit() uint64 {
	return m.BlockGasLimitMock
}

// DeductInitialGasForExecution -
func (m *MeteringContextMock) DeductInitialGasForExecution(_ []byte) error {
	return m.Err
}

// DeductInitialGasForDirectDeployment -
func (m *MeteringContextMock) DeductInitialGasForDirectDeployment(_ arwen.CodeDeployInput) error {
	return m.Err
}

// DeductInitialGasForIndirectDeployment -
func (m *MeteringContextMock) DeductInitialGasForIndirectDeployment(_ arwen.CodeDeployInput) error {
	return m.Err
}

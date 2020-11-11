package mock

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
)

var _ arwen.MeteringContext = (*MeteringContextMock)(nil)

type MeteringContextMock struct {
	GasCost           *config.GasCost
	GasLeftMock       uint64
	GasLockedMock     uint64
	GasComputedToLock uint64
	BlockGasLimitMock uint64
	Err               error
}

func (m *MeteringContextMock) SetGasSchedule(gasSchedule config.GasScheduleMap) {
	gasCostConfig, _ := config.CreateGasConfig(gasSchedule)
	m.GasCost = gasCostConfig
}

func (m *MeteringContextMock) GasSchedule() *config.GasCost {
	return m.GasCost
}

func (m *MeteringContextMock) UseGas(gas uint64) {
}

func (m *MeteringContextMock) FreeGas(gas uint64) {
}

func (m *MeteringContextMock) RestoreGas(gas uint64) {
}

func (m *MeteringContextMock) GasLeft() uint64 {
	return m.GasLeftMock
}

func (m *MeteringContextMock) BoundGasLimit(limit uint64) uint64 {
	gasLeft := m.GasLeft()

	if gasLeft < limit {
		return gasLeft
	}
	return limit
}

func (m *MeteringContextMock) ComputeGasLockedForAsync() uint64 {
	return m.GasComputedToLock
}

func (m *MeteringContextMock) DeductGasIfAsyncStep() error {
	return m.Err
}

func (m *MeteringContextMock) UseGasBounded(gas uint64) error {
	return m.Err
}

func (m *MeteringContextMock) UnlockGasIfAsyncCallback() {
}

func (m *MeteringContextMock) UseGasForAsyncStep() error {
	return m.Err
}

func (m *MeteringContextMock) UnlockGasIfAsyncStep() {
}

func (m *MeteringContextMock) GetGasLocked() uint64 {
	return m.GasLockedMock
}

func (m *MeteringContextMock) BlockGasLimit() uint64 {
	return m.BlockGasLimitMock
}

func (m *MeteringContextMock) DeductInitialGasForExecution(contract []byte) error {
	return m.Err
}

func (m *MeteringContextMock) DeductInitialGasForDirectDeployment(input arwen.CodeDeployInput) error {
	return m.Err
}

func (m *MeteringContextMock) DeductInitialGasForIndirectDeployment(input arwen.CodeDeployInput) error {
	return m.Err
}

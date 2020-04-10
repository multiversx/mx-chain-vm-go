package mock

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
)

var _ arwen.MeteringContext = (*MeteringContextMock)(nil)

type MeteringContextMock struct {
	GasCost           *config.GasCost
	GasLeftMock       uint64
	BlockGasLimitMock uint64
	GasLocked         uint64
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

func (m *MeteringContextMock) BoundGasLimit(value int64) uint64 {
	gasLeft := m.GasLeft()
	limit := uint64(value)

	if gasLeft < limit {
		return gasLeft
	}
	return limit
}

func (m *MeteringContextMock) UnlockGasIfAsyncStep() {
}

func (m *MeteringContextMock) GetGasLockedForAsyncStep() uint64 {
	return m.GasLocked
}

func (m *MeteringContextMock) BlockGasLimit() uint64 {
	return m.BlockGasLimitMock
}

func (m *MeteringContextMock) DeductInitialGasForExecution(contract []byte) error {
	if m.Err != nil {
		return m.Err
	}
	return nil
}

func (m *MeteringContextMock) DeductInitialGasForDirectDeployment(input arwen.CodeDeployInput) error {
	if m.Err != nil {
		return m.Err
	}
	return nil
}

func (m *MeteringContextMock) DeductInitialGasForIndirectDeployment(input arwen.CodeDeployInput) error {
	if m.Err != nil {
		return m.Err
	}
	return nil
}

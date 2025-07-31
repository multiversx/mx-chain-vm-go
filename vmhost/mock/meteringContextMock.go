package mock

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

var _ vmhost.MeteringContext = (*MeteringContextMock)(nil)

// MeteringContextMock is used in tests to check the MeteringContext interface method calls
type MeteringContextMock struct {
	GasCost           *config.GasCost
	GasLeftMock       uint64
	GasFreedMock      uint64
	GasLockedMock     uint64
	GasProvidedMock   uint64
	GasComputedToLock uint64
	BlockGasLimitMock uint64
	Err               error
}

// InitState mocked method
func (m *MeteringContextMock) InitState() {}

// PushState mocked method
func (m *MeteringContextMock) PushState() {}

// PopSetActiveState mocked method
func (m *MeteringContextMock) PopSetActiveState() {}

// PopMergeActiveState mocked method
func (m *MeteringContextMock) PopMergeActiveState() {}

// PopDiscard mocked method
func (m *MeteringContextMock) PopDiscard() {}

// ClearStateStack mocked method
func (m *MeteringContextMock) ClearStateStack() {}

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
func (m *MeteringContextMock) useGas(gas uint64) {
	if gas > m.GasLeftMock {
		m.GasLeftMock = 0
		return
	}

	m.GasLeftMock -= gas
}

// UseGasBoundedAndAddTracedGas -
func (m *MeteringContextMock) UseGasBoundedAndAddTracedGas(_ string, _ uint64) error {
	return nil
}

// FreeGas mocked method
func (m *MeteringContextMock) FreeGas(gas uint64) {
	m.GasFreedMock += gas
}

// RestoreGas mocked method
func (m *MeteringContextMock) RestoreGas(gas uint64) {
	m.GasLeftMock += gas
}

// GasLeft mocked method
func (m *MeteringContextMock) GasLeft() uint64 {
	return m.GasLeftMock
}

// UpdateGasStateOnSuccess mocked method
func (m *MeteringContextMock) UpdateGasStateOnSuccess(_ *vmcommon.VMOutput) error {
	return nil
}

// UpdateGasStateOnFailure mocked method
func (m *MeteringContextMock) UpdateGasStateOnFailure(_ *vmcommon.VMOutput) {}

// InitStateFromContractCallInput mocked method
func (m *MeteringContextMock) InitStateFromContractCallInput(_ *vmcommon.VMInput) {}

// TrackGasUsedByOutOfVMFunction mocked method
func (m *MeteringContextMock) TrackGasUsedByOutOfVMFunction(_ *vmcommon.ContractCallInput, _ *vmcommon.VMOutput, _ *vmcommon.ContractCallInput) {
}

// GasUsedByContract mocked method
func (m *MeteringContextMock) GasUsedByContract() (uint64, uint64) {
	return 0, 0
}

// GasUsedForExecution mocked method
func (m *MeteringContextMock) GasUsedForExecution() uint64 {
	return 0
}

// GasSpentByContract mocked method
func (m *MeteringContextMock) GasSpentByContract() uint64 {
	return 0
}

// GetGasForExecution mocked method
func (m *MeteringContextMock) GetGasForExecution() uint64 {
	return 0
}

// GetGasProvided mocked method
func (m *MeteringContextMock) GetGasProvided() uint64 {
	return 0
}

// GetSCPrepareInitialCost mocked method
func (m *MeteringContextMock) GetSCPrepareInitialCost() uint64 {
	return 0
}

// BoundGasLimit mocked method
func (m *MeteringContextMock) BoundGasLimit(value int64) uint64 {
	gasLeft := m.GasLeft()
	limit := uint64(value)

	if gasLeft < limit {
		return gasLeft
	}
	return limit
}

// ComputeExtraGasLockedForAsync mocked method
func (m *MeteringContextMock) ComputeExtraGasLockedForAsync() uint64 {
	return m.GasComputedToLock
}

// DeductGasIfAsyncStep mocked method
func (m *MeteringContextMock) DeductGasIfAsyncStep() error {
	return m.Err
}

// UseGasForContractInit mocked method
func (m *MeteringContextMock) UseGasForContractInit(gasToUse uint64) {
	m.useGas(gasToUse)
}

// UseGasBounded mocked method
func (m *MeteringContextMock) UseGasBounded(gas uint64) error {
	if m.Err != nil {
		return m.Err
	}
	if m.GasLeft() <= gas {
		return vmhost.ErrNotEnoughGas
	}
	m.useGas(gas)
	return nil
}

// UnlockGasIfAsyncCallback mocked method
func (m *MeteringContextMock) UnlockGasIfAsyncCallback() {}

// UseGasForAsyncStep mocked method
func (m *MeteringContextMock) UseGasForAsyncStep() error {
	return m.Err
}

// UnlockGasIfAsyncStep mocked method
func (m *MeteringContextMock) UnlockGasIfAsyncStep() {}

// GetGasLocked mocked method
func (m *MeteringContextMock) GetGasLocked() uint64 {
	return m.GasLockedMock
}

// BlockGasLimit mocked method
func (m *MeteringContextMock) BlockGasLimit() uint64 {
	return m.BlockGasLimitMock
}

// DeductInitialGasForExecution mocked method
func (m *MeteringContextMock) DeductInitialGasForExecution(_ []byte) error {
	return m.Err
}

// DeductInitialGasForDirectDeployment mocked method
func (m *MeteringContextMock) DeductInitialGasForDirectDeployment(_ vmhost.CodeDeployInput) error {
	return m.Err
}

// DeductInitialGasForIndirectDeployment mocked method
func (m *MeteringContextMock) DeductInitialGasForIndirectDeployment(_ vmhost.CodeDeployInput) error {
	return m.Err
}

// EnableRestoreGas mocked method
func (m *MeteringContextMock) EnableRestoreGas() {}

// DisableRestoreGas mocked method
func (m *MeteringContextMock) DisableRestoreGas() {}

// StartGasTracing mocked method
func (m *MeteringContextMock) StartGasTracing(_ string) {}

// SetGasTracing mocked method
func (m *MeteringContextMock) SetGasTracing(_ bool) {}

// GetGasTrace returns nil
func (m *MeteringContextMock) GetGasTrace() map[string]map[string][]uint64 {
	return nil
}

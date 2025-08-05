package mock

import (
	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/crypto"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

// UnsafeVMHostMock is a mock for the VMHost interface that helps testing the unsafe mode
type UnsafeVMHostMock struct {
	unsafeMode            bool
	FailExecutionCalled   bool
	FailConditionallyCalled bool
	BreakpointValue       vmhost.BreakpointValue
	meteringContext       vmhost.MeteringContext
	runtimeContext        vmhost.RuntimeContext
}

// GetVersion -
func (h *UnsafeVMHostMock) GetVersion() string {
	return "unsafe mock"
}

// Crypto -
func (h *UnsafeVMHostMock) Crypto() crypto.VMCrypto {
	return nil
}

// Blockchain -
func (h *UnsafeVMHostMock) Blockchain() vmhost.BlockchainContext {
	return nil
}

// Runtime -
func (h *UnsafeVMHostMock) Runtime() vmhost.RuntimeContext {
	if h.runtimeContext == nil {
		h.runtimeContext = &RuntimeContextMock{
			VMHost: h,
		}
	}
	return h.runtimeContext
}

// Output -
func (h *UnsafeVMHostMock) Output() vmhost.OutputContext {
	return nil
}

// Metering -
func (h *UnsafeVMHostMock) Metering() vmhost.MeteringContext {
	if h.meteringContext == nil {
		h.meteringContext = &MeteringContextMock{}
	}
	return h.meteringContext
}

// Storage -
func (h *UnsafeVMHostMock) Storage() vmhost.StorageContext {
	return nil
}

// EnableEpochsHandler -
func (h *UnsafeVMHostMock) EnableEpochsHandler() vmhost.EnableEpochsHandler {
	return nil
}

// ManagedTypes -
func (h *UnsafeVMHostMock) ManagedTypes() vmhost.ManagedTypesContext {
	return nil
}

// AreInSameShard -
func (h *UnsafeVMHostMock) AreInSameShard(left []byte, right []byte) bool {
	return true
}

// IsAllowedToExecute -
func (h *UnsafeVMHostMock) IsAllowedToExecute(_ string) bool {
	return true
}

// ExecuteESDTTransfer -
func (h *UnsafeVMHostMock) ExecuteESDTTransfer(_ *vmhost.ESDTTransfersArgs, _ vm.CallType) (*vmcommon.VMOutput, uint64, error) {
	return nil, 0, nil
}

// CreateNewContract -
func (h *UnsafeVMHostMock) CreateNewContract(_ *vmcommon.ContractCreateInput, _ int) ([]byte, error) {
	return nil, nil
}

// ExecuteOnSameContext -
func (h *UnsafeVMHostMock) ExecuteOnSameContext(_ *vmcommon.ContractCallInput) error {
	return nil
}

// ExecuteOnDestContext -
func (h *UnsafeVMHostMock) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, bool, error) {
	return nil, true, nil
}

// InitState -
func (h *UnsafeVMHostMock) InitState() {
}

// IsBuiltinFunctionName -
func (h *UnsafeVMHostMock) IsBuiltinFunctionName(_ string) bool {
	return false
}

// IsBuiltinFunctionCall -
func (h *UnsafeVMHostMock) IsBuiltinFunctionCall(_ []byte) bool {
	return false
}

// GetGasScheduleMap -
func (h *UnsafeVMHostMock) GetGasScheduleMap() config.GasScheduleMap {
	return make(config.GasScheduleMap)
}

// RunSmartContractCall -
func (h *UnsafeVMHostMock) RunSmartContractCall(_ *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, err error) {
	return nil, nil
}

// RunSmartContractCreate -
func (h *UnsafeVMHostMock) RunSmartContractCreate(_ *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput, err error) {
	return nil, nil
}

// GasScheduleChange -
func (h *UnsafeVMHostMock) GasScheduleChange(_ config.GasScheduleMap) {
}

// SetBuiltInFunctionsContainer -
func (h *UnsafeVMHostMock) SetBuiltInFunctionsContainer(_ vmcommon.BuiltInFunctionContainer) {
}

// IsInterfaceNil -
func (h *UnsafeVMHostMock) IsInterfaceNil() bool {
	return false
}

// GetContexts -
func (h *UnsafeVMHostMock) GetContexts() (
	vmhost.ManagedTypesContext,
	vmhost.BlockchainContext,
	vmhost.MeteringContext,
	vmhost.OutputContext,
	vmhost.RuntimeContext,
	vmhost.AsyncContext,
	vmhost.StorageContext,
) {
	return nil, nil, h.Metering(), nil, h.Runtime(), nil, nil
}

// SetRuntimeContext -
func (h *UnsafeVMHostMock) SetRuntimeContext(runtime vmhost.RuntimeContext) {
}

// Async -
func (h *UnsafeVMHostMock) Async() vmhost.AsyncContext {
	return nil
}

// CompleteLogEntriesWithCallType -
func (h *UnsafeVMHostMock) CompleteLogEntriesWithCallType(vmOutput *vmcommon.VMOutput, callType string) {
}

// Close -
func (h *UnsafeVMHostMock) Close() error {
	return nil
}

// Reset -
func (h *UnsafeVMHostMock) Reset() {
}

// SetGasTracing -
func (h *UnsafeVMHostMock) SetGasTracing(enableGasTracing bool) {
}

// GetGasTrace -
func (h *UnsafeVMHostMock) GetGasTrace() map[string]map[string][]uint64 {
	return make(map[string]map[string][]uint64)
}

// SetUnsafeMode sets the unsafe mode flag
func (h *UnsafeVMHostMock) SetUnsafeMode(unsafeMode bool) {
	h.unsafeMode = unsafeMode
}

// IsUnsafeMode returns true if the unsafe mode is enabled
func (h *UnsafeVMHostMock) IsUnsafeMode() bool {
	return h.unsafeMode
}

// FailExecution marks that FailExecution was called
func (h *UnsafeVMHostMock) FailExecution(_ error) {
	h.FailExecutionCalled = true
	h.BreakpointValue = vmhost.BreakpointExecutionFailed
}

// FailExecutionConditionally marks that FailExecutionConditionally was called
func (h *UnsafeVMHostMock) FailExecutionConditionally(err error) {
	h.FailConditionallyCalled = true
	if !h.unsafeMode {
		h.FailExecution(err)
	}
}

// RuntimeContextMock is a mock for the RuntimeContext interface
type RuntimeContextMock struct {
	vmhost.RuntimeContext
	VMHost              vmhost.VMHost
	FailExecutionCalled bool
}

// FailExecution -
func (r *RuntimeContextMock) FailExecution(err error) {
	r.VMHost.(*UnsafeVMHostMock).FailExecution(err)
}

// FailExecutionConditionally -
func (r *RuntimeContextMock) FailExecutionConditionally(err error) {
	r.VMHost.FailExecutionConditionally(err)
}

// MeteringContextMock is a mock for the MeteringContext interface
type MeteringContextMock struct {
	vmhost.MeteringContext
	UseGasBoundedAndAddTracedGasCalled func(name string, gas uint64) error
	UseGasBoundedCalled                func(gas uint64) error
	GasLeftCalled                      func() uint64
}

// UseGasBoundedAndAddTracedGas -
func (m *MeteringContextMock) UseGasBoundedAndAddTracedGas(name string, gas uint64) error {
	if m.UseGasBoundedAndAddTracedGasCalled != nil {
		return m.UseGasBoundedAndAddTracedGasCalled(name, gas)
	}
	return nil
}

// GasLeft -
func (m *MeteringContextMock) GasLeft() uint64 {
	if m.GasLeftCalled != nil {
		return m.GasLeftCalled()
	}
	return 0
}

// UseGasBounded -
func (m *MeteringContextMock) UseGasBounded(gas uint64) error {
	if m.UseGasBoundedCalled != nil {
		return m.UseGasBoundedCalled(gas)
	}
	return nil
}

package mock

import (
	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	"github.com/multiversx/mx-chain-vm-v1_4-go/config"
	"github.com/multiversx/mx-chain-vm-v1_4-go/crypto"
	"github.com/multiversx/mx-chain-vm-v1_4-go/wasmer"
)

var _ vmhost.VMHost = (*VMHostStub)(nil)

// VMHostStub is used in tests to check the VMHost interface method calls
type VMHostStub struct {
	InitStateCalled       func()
	PushStateCalled       func()
	PopStateCalled        func()
	ClearStateStackCalled func()
	GetVersionCalled      func() string

	CryptoCalled                func() crypto.VMCrypto
	BlockchainCalled            func() vmhost.BlockchainContext
	RuntimeCalled               func() vmhost.RuntimeContext
	ManagedTypesCalled          func() vmhost.ManagedTypesContext
	OutputCalled                func() vmhost.OutputContext
	MeteringCalled              func() vmhost.MeteringContext
	StorageCalled               func() vmhost.StorageContext
	EnableEpochsHandlerCalled   func() vmcommon.EnableEpochsHandler
	ExecuteESDTTransferCalled   func(destination []byte, sender []byte, transfers []*vmcommon.ESDTTransfer, callType vm.CallType) (*vmcommon.VMOutput, uint64, error)
	CreateNewContractCalled     func(input *vmcommon.ContractCreateInput) ([]byte, error)
	ExecuteOnSameContextCalled  func(input *vmcommon.ContractCallInput) (*vmhost.AsyncContextInfo, error)
	ExecuteOnDestContextCalled  func(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, *vmhost.AsyncContextInfo, error)
	GetAPIMethodsCalled         func() *wasmer.Imports
	IsBuiltinFunctionNameCalled func(functionName string) bool
	AreInSameShardCalled        func(left []byte, right []byte) bool

	RunSmartContractCallCalled   func(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, err error)
	RunSmartContractCreateCalled func(input *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput, err error)
	GetGasScheduleMapCalled      func() config.GasScheduleMap
	GasScheduleChangeCalled      func(newGasSchedule config.GasScheduleMap)
	IsInterfaceNilCalled         func() bool

	SetRuntimeContextCalled func(runtime vmhost.RuntimeContext)
	GetContextsCalled       func() (vmhost.ManagedTypesContext, vmhost.BlockchainContext, vmhost.MeteringContext, vmhost.OutputContext, vmhost.RuntimeContext, vmhost.StorageContext)

	SetBuiltInFunctionsContainerCalled func(builtInFuncs vmcommon.BuiltInFunctionContainer)
}

// GetVersion mocked method
func (vhs *VMHostStub) GetVersion() string {
	if vhs.GetVersionCalled != nil {
		return vhs.GetVersionCalled()
	}

	return "stub"
}

// InitState mocked method
func (vhs *VMHostStub) InitState() {
	if vhs.InitStateCalled != nil {
		vhs.InitStateCalled()
	}
}

// PushState mocked method
func (vhs *VMHostStub) PushState() {
	if vhs.PushStateCalled != nil {
		vhs.PushStateCalled()
	}
}

// PopState mocked method
func (vhs *VMHostStub) PopState() {
	if vhs.PopStateCalled != nil {
		vhs.PopStateCalled()
	}
}

// ClearStateStack mocked method
func (vhs *VMHostStub) ClearStateStack() {
	if vhs.ClearStateStackCalled != nil {
		vhs.ClearStateStackCalled()
	}
}

// Crypto mocked method
func (vhs *VMHostStub) Crypto() crypto.VMCrypto {
	if vhs.CryptoCalled != nil {
		return vhs.CryptoCalled()
	}
	return nil
}

// Blockchain mocked method
func (vhs *VMHostStub) Blockchain() vmhost.BlockchainContext {
	if vhs.BlockchainCalled != nil {
		return vhs.BlockchainCalled()
	}
	return nil
}

// Runtime mocked method
func (vhs *VMHostStub) Runtime() vmhost.RuntimeContext {
	if vhs.RuntimeCalled != nil {
		return vhs.RuntimeCalled()
	}
	return nil
}

// ManagedTypes mocked method
func (vhs *VMHostStub) ManagedTypes() vmhost.ManagedTypesContext {
	if vhs.ManagedTypesCalled != nil {
		return vhs.ManagedTypesCalled()
	}
	return nil
}

// IsArwenV2Enabled mocked method
func (vhs *VMHostStub) IsArwenV2Enabled() bool {
	return true
}

// IsArwenV3Enabled mocked method
func (vhs *VMHostStub) IsArwenV3Enabled() bool {
	return true
}

// IsAheadOfTimeCompileEnabled mocked method
func (vhs *VMHostStub) IsAheadOfTimeCompileEnabled() bool {
	return true
}

// IsDynamicGasLockingEnabled mocked method
func (vhs *VMHostStub) IsDynamicGasLockingEnabled() bool {
	return true
}

// IsESDTFunctionsEnabled mocked method
func (vhs *VMHostStub) IsESDTFunctionsEnabled() bool {
	return true
}

// Output mocked method
func (vhs *VMHostStub) Output() vmhost.OutputContext {
	if vhs.OutputCalled != nil {
		return vhs.OutputCalled()
	}
	return nil
}

// Metering mocked method
func (vhs *VMHostStub) Metering() vmhost.MeteringContext {
	if vhs.MeteringCalled != nil {
		return vhs.MeteringCalled()
	}
	return nil
}

// Storage mocked method
func (vhs *VMHostStub) Storage() vmhost.StorageContext {
	if vhs.StorageCalled != nil {
		return vhs.StorageCalled()
	}
	return nil
}

// EnableEpochsHandler mocked method
func (vhs *VMHostStub) EnableEpochsHandler() vmcommon.EnableEpochsHandler {
	if vhs.EnableEpochsHandlerCalled != nil {
		return vhs.EnableEpochsHandlerCalled()
	}
	return nil
}

// ExecuteESDTTransfer mocked method
func (vhs *VMHostStub) ExecuteESDTTransfer(destination []byte, sender []byte, transfers []*vmcommon.ESDTTransfer, callType vm.CallType) (*vmcommon.VMOutput, uint64, error) {
	if vhs.ExecuteESDTTransferCalled != nil {
		return vhs.ExecuteESDTTransferCalled(destination, sender, transfers, callType)
	}
	return nil, 0, nil
}

// CreateNewContract mocked method
func (vhs *VMHostStub) CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error) {
	if vhs.CreateNewContractCalled != nil {
		return vhs.CreateNewContractCalled(input)
	}
	return nil, nil
}

// ExecuteOnSameContext mocked method
func (vhs *VMHostStub) ExecuteOnSameContext(input *vmcommon.ContractCallInput) (*vmhost.AsyncContextInfo, error) {
	if vhs.ExecuteOnSameContextCalled != nil {
		return vhs.ExecuteOnSameContextCalled(input)
	}
	return nil, nil
}

// ExecuteOnDestContext mocked method
func (vhs *VMHostStub) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, *vmhost.AsyncContextInfo, error) {
	if vhs.ExecuteOnDestContextCalled != nil {
		return vhs.ExecuteOnDestContextCalled(input)
	}
	return nil, nil, nil
}

// AreInSameShard mocked method
func (vhs *VMHostStub) AreInSameShard(left []byte, right []byte) bool {
	if vhs.AreInSameShardCalled != nil {
		return vhs.AreInSameShardCalled(left, right)
	}
	return true
}

// GetAPIMethods mocked method
func (vhs *VMHostStub) GetAPIMethods() *wasmer.Imports {
	if vhs.GetAPIMethodsCalled != nil {
		return vhs.GetAPIMethodsCalled()
	}
	return nil
}

// IsBuiltinFunctionName mocked method
func (vhs *VMHostStub) IsBuiltinFunctionName(functionName string) bool {
	if vhs.IsBuiltinFunctionNameCalled != nil {
		return vhs.IsBuiltinFunctionNameCalled(functionName)
	}
	return false
}

// GetGasScheduleMap returns the currently stored gas schedule
func (vhs *VMHostStub) GetGasScheduleMap() config.GasScheduleMap {
	if vhs.GetGasScheduleMapCalled != nil {
		return vhs.GetGasScheduleMapCalled()
	}
	return nil
}

// RunSmartContractCall mocked method
func (vhs *VMHostStub) RunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, err error) {
	if vhs.RunSmartContractCallCalled != nil {
		return vhs.RunSmartContractCallCalled(input)
	}
	return nil, nil
}

// RunSmartContractCreate mocked method
func (vhs *VMHostStub) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput, err error) {
	if vhs.RunSmartContractCreateCalled != nil {
		return vhs.RunSmartContractCreateCalled(input)
	}
	return nil, nil
}

// GasScheduleChange mocked method
func (vhs *VMHostStub) GasScheduleChange(newGasSchedule config.GasScheduleMap) {
	if vhs.GasScheduleChangeCalled != nil {
		vhs.GasScheduleChangeCalled(newGasSchedule)
	}
}

// SetBuiltInFunctionsContainer mocked method
func (vhs *VMHostStub) SetBuiltInFunctionsContainer(builtInFuncs vmcommon.BuiltInFunctionContainer) {
	if vhs.SetBuiltInFunctionsContainerCalled != nil {
		vhs.SetBuiltInFunctionsContainerCalled(builtInFuncs)
	}
}

// IsInterfaceNil mocked method
func (vhs *VMHostStub) IsInterfaceNil() bool {
	if vhs.IsInterfaceNilCalled != nil {
		return vhs.IsInterfaceNilCalled()
	}
	return false
}

// GetContexts mocked method
func (vhs *VMHostStub) GetContexts() (
	vmhost.ManagedTypesContext,
	vmhost.BlockchainContext,
	vmhost.MeteringContext,
	vmhost.OutputContext,
	vmhost.RuntimeContext,
	vmhost.StorageContext,
) {
	if vhs.GetContextsCalled != nil {
		return vhs.GetContextsCalled()
	}
	return nil, nil, nil, nil, nil, nil
}

// SetRuntimeContext mocked method
func (vhs *VMHostStub) SetRuntimeContext(runtime vmhost.RuntimeContext) {
	if vhs.SetRuntimeContextCalled != nil {
		vhs.SetRuntimeContextCalled(runtime)
	}
}

// FixOOGReturnCodeEnabled mocked method
func (vhs *VMHostStub) FixOOGReturnCodeEnabled() bool {
	return true
}

// FixFailExecutionEnabled mocked method
func (vhs *VMHostStub) FixFailExecutionEnabled() bool {
	return true
}

// CreateNFTOnExecByCallerEnabled mocked method
func (vhs *VMHostStub) CreateNFTOnExecByCallerEnabled() bool {
	return true
}

// DisableExecByCaller mocked method
func (vhs *VMHostStub) DisableExecByCaller() bool {
	return true
}

// CheckExecuteReadOnly mocked method
func (vhs *VMHostStub) CheckExecuteReadOnly() bool {
	return true
}

// Close -
func (vhs *VMHostStub) Close() error {
	return nil
}

// Reset -
func (vhs *VMHostStub) Reset() {
}

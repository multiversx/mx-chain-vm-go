package mock

import (
	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/crypto"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

var _ vmhost.VMHost = (*VMHostMock)(nil)

// VMHostMock is used in tests to check the VMHost interface method calls
type VMHostMock struct {
	BlockChainHook vmcommon.BlockchainHook
	CryptoHook     crypto.VMCrypto

	EthInput []byte

	BlockchainContext        vmhost.BlockchainContext
	RuntimeContext           vmhost.RuntimeContext
	AsyncContext             vmhost.AsyncContext
	OutputContext            vmhost.OutputContext
	MeteringContext          vmhost.MeteringContext
	StorageContext           vmhost.StorageContext
	EnableEpochsHandlerField vmhost.EnableEpochsHandler
	ManagedTypesContext      vmhost.ManagedTypesContext

	IsBuiltinFunc bool

	StoredInputs []*vmcommon.ContractCallInput

	VMOutputQueue    []*vmcommon.VMOutput
	VMOutputToReturn int
	Err              error
}

// GetVersion mocked method
func (host *VMHostMock) GetVersion() string {
	return "mock"
}

// Crypto mocked method
func (host *VMHostMock) Crypto() crypto.VMCrypto {
	return host.CryptoHook
}

// Blockchain mocked method
func (host *VMHostMock) Blockchain() vmhost.BlockchainContext {
	return host.BlockchainContext
}

// Runtime mocked method
func (host *VMHostMock) Runtime() vmhost.RuntimeContext {
	return host.RuntimeContext
}

// Output mocked method
func (host *VMHostMock) Output() vmhost.OutputContext {
	return host.OutputContext
}

// Metering mocked method
func (host *VMHostMock) Metering() vmhost.MeteringContext {
	return host.MeteringContext
}

// Storage mocked method
func (host *VMHostMock) Storage() vmhost.StorageContext {
	return host.StorageContext
}

// EnableEpochsHandler mocked method
func (host *VMHostMock) EnableEpochsHandler() vmhost.EnableEpochsHandler {
	return host.EnableEpochsHandlerField
}

// ManagedTypes mocked method
func (host *VMHostMock) ManagedTypes() vmhost.ManagedTypesContext {
	return host.ManagedTypesContext
}

// IsAheadOfTimeCompileEnabled mocked method
func (host *VMHostMock) IsAheadOfTimeCompileEnabled() bool {
	return true
}

// IsDynamicGasLockingEnabled mocked method
func (host *VMHostMock) IsDynamicGasLockingEnabled() bool {
	return true
}

// IsESDTFunctionsEnabled mocked method
func (host *VMHostMock) IsESDTFunctionsEnabled() bool {
	return true
}

// AreInSameShard mocked method
func (host *VMHostMock) AreInSameShard(left []byte, right []byte) bool {
	leftShard := host.BlockchainContext.GetShardOfAddress(left)
	rightShard := host.BlockchainContext.GetShardOfAddress(right)
	return leftShard == rightShard
}

// ExecuteESDTTransfer mocked method
func (host *VMHostMock) ExecuteESDTTransfer(_ *vmhost.ESDTTransfersArgs, _ vm.CallType) (*vmcommon.VMOutput, uint64, error) {
	return nil, 0, nil
}

// CreateNewContract mocked method
func (host *VMHostMock) CreateNewContract(_ *vmcommon.ContractCreateInput, _ int) ([]byte, error) {
	return nil, nil
}

// ExecuteOnSameContext mocked method
func (host *VMHostMock) ExecuteOnSameContext(_ *vmcommon.ContractCallInput) error {
	return nil
}

// ExecuteOnDestContext mocked method
func (host *VMHostMock) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, bool, error) {
	if host.Err != nil {
		return nil, true, host.Err
	}
	host.StoreInput(input)
	return host.GetNextVMOutput(), true, nil
}

// InitState mocked method
func (host *VMHostMock) InitState() {
}

// PushState mocked method
func (host *VMHostMock) PushState() {
}

// PopState mocked method
func (host *VMHostMock) PopState() {
}

// ClearStateStack mocked method
func (host *VMHostMock) ClearStateStack() {
}

// IsBuiltinFunctionName mocked method
func (host *VMHostMock) IsBuiltinFunctionName(_ string) bool {
	return host.IsBuiltinFunc
}

// IsBuiltinFunctionCall mocked method
func (host *VMHostMock) IsBuiltinFunctionCall(_ []byte) bool {
	return host.IsBuiltinFunc
}

// GetGasScheduleMap mocked method
func (host *VMHostMock) GetGasScheduleMap() config.GasScheduleMap {
	return make(config.GasScheduleMap)
}

// RunSmartContractCall mocked method
func (host *VMHostMock) RunSmartContractCall(_ *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, err error) {
	return nil, nil
}

// RunSmartContractCreate mocked method
func (host *VMHostMock) RunSmartContractCreate(_ *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput, err error) {
	return nil, nil
}

// GasScheduleChange mocked method
func (host *VMHostMock) GasScheduleChange(_ config.GasScheduleMap) {
}

// SetBuiltInFunctionsContainer mocked method
func (host *VMHostMock) SetBuiltInFunctionsContainer(_ vmcommon.BuiltInFunctionContainer) {
}

// IsInterfaceNil mocked method
func (host *VMHostMock) IsInterfaceNil() bool {
	return false
}

// GetContexts mocked method
func (host *VMHostMock) GetContexts() (
	vmhost.ManagedTypesContext,
	vmhost.BlockchainContext,
	vmhost.MeteringContext,
	vmhost.OutputContext,
	vmhost.RuntimeContext,
	vmhost.AsyncContext,
	vmhost.StorageContext,
) {
	return host.ManagedTypesContext, host.BlockchainContext, host.MeteringContext, host.OutputContext, host.RuntimeContext, host.AsyncContext, host.StorageContext
}

// SetRuntimeContext mocked method
func (host *VMHostMock) SetRuntimeContext(runtime vmhost.RuntimeContext) {
	host.RuntimeContext = runtime
}

// Async mocked method
func (host *VMHostMock) Async() vmhost.AsyncContext {
	return host.AsyncContext
}

// FixOOGReturnCodeEnabled mocked method
func (host *VMHostMock) FixOOGReturnCodeEnabled() bool {
	return true
}

// StoreInput enqueues the given ContractCallInput
func (host *VMHostMock) StoreInput(input *vmcommon.ContractCallInput) {
	if host.StoredInputs == nil {
		host.StoredInputs = make([]*vmcommon.ContractCallInput, 0)
	}
	host.StoredInputs = append(host.StoredInputs, input)
}

// EnqueueVMOutput enqueues the given VMOutput
func (host *VMHostMock) EnqueueVMOutput(vmOutput *vmcommon.VMOutput) {
	if host.VMOutputQueue == nil {
		host.VMOutputQueue = make([]*vmcommon.VMOutput, 1)
		host.VMOutputQueue[0] = vmOutput
		host.VMOutputToReturn = 0
		return
	}

	host.VMOutputQueue = append(host.VMOutputQueue, vmOutput)
}

// GetNextVMOutput returns the next VMOutput in the queue
func (host *VMHostMock) GetNextVMOutput() *vmcommon.VMOutput {
	if host.VMOutputToReturn >= len(host.VMOutputQueue) {
		return nil
	}

	vmOutput := host.VMOutputQueue[host.VMOutputToReturn]
	host.VMOutputToReturn += 1
	return vmOutput
}

// CompleteLogEntriesWithCallType mocked method
func (host *VMHostMock) CompleteLogEntriesWithCallType(vmOutput *vmcommon.VMOutput, callType string) {
}

// Close -
func (host *VMHostMock) Close() error {
	return nil
}

// Reset -
func (host *VMHostMock) Reset() {
}

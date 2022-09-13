package mock

import (
	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/config"
	"github.com/ElrondNetwork/wasm-vm/crypto"
	"github.com/ElrondNetwork/wasm-vm/wasmer"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ arwen.VMHost = (*VMHostMock)(nil)

// VMHostMock is used in tests to check the VMHost interface method calls
type VMHostMock struct {
	BlockChainHook vmcommon.BlockchainHook
	CryptoHook     crypto.VMCrypto

	EthInput []byte

	BlockchainContext   arwen.BlockchainContext
	RuntimeContext      arwen.RuntimeContext
	AsyncContext        arwen.AsyncContext
	OutputContext       arwen.OutputContext
	MeteringContext     arwen.MeteringContext
	StorageContext      arwen.StorageContext
	ManagedTypesContext arwen.ManagedTypesContext

	SCAPIMethods  *wasmer.Imports
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
func (host *VMHostMock) Blockchain() arwen.BlockchainContext {
	return host.BlockchainContext
}

// Runtime mocked method
func (host *VMHostMock) Runtime() arwen.RuntimeContext {
	return host.RuntimeContext
}

// Output mocked method
func (host *VMHostMock) Output() arwen.OutputContext {
	return host.OutputContext
}

// Metering mocked method
func (host *VMHostMock) Metering() arwen.MeteringContext {
	return host.MeteringContext
}

// Storage mocked method
func (host *VMHostMock) Storage() arwen.StorageContext {
	return host.StorageContext
}

// BigInt mocked method
func (host *VMHostMock) ManagedTypes() arwen.ManagedTypesContext {
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
func (host *VMHostMock) ExecuteESDTTransfer(_ []byte, _ []byte, _ []*vmcommon.ESDTTransfer, _ vm.CallType) (*vmcommon.VMOutput, uint64, error) {
	return nil, 0, nil
}

// CreateNewContract mocked method
func (host *VMHostMock) CreateNewContract(_ *vmcommon.ContractCreateInput) ([]byte, error) {
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

// GetAPIMethods mocked method
func (host *VMHostMock) GetAPIMethods() *wasmer.Imports {
	return host.SCAPIMethods
}

// IsBuiltinFunctionName mocked method
func (host *VMHostMock) IsBuiltinFunctionName(_ string) bool {
	return host.IsBuiltinFunc
}

// IsBuiltinFunctionName mocked method
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
	arwen.ManagedTypesContext,
	arwen.BlockchainContext,
	arwen.MeteringContext,
	arwen.OutputContext,
	arwen.RuntimeContext,
	arwen.AsyncContext,
	arwen.StorageContext,
) {
	return host.ManagedTypesContext, host.BlockchainContext, host.MeteringContext, host.OutputContext, host.RuntimeContext, host.AsyncContext, host.StorageContext
}

// SetRuntimeContext mocked method
func (host *VMHostMock) SetRuntimeContext(runtime arwen.RuntimeContext) {
	host.RuntimeContext = runtime
}

// Async mocked method
func (host *VMHostMock) Async() arwen.AsyncContext {
	return host.AsyncContext
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

// Close -
func (host *VMHostMock) Close() error {
	return nil
}

// Reset -
func (host *VMHostMock) Reset() {
}

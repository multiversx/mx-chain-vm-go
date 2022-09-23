package mock

import (
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/executor"
)

var _ arwen.RuntimeContext = (*RuntimeContextMock)(nil)

// RuntimeContextMock is used in tests to check the RuntimeContextMock interface method calls
type RuntimeContextMock struct {
	Err                      error
	VMInput                  *vmcommon.ContractCallInput
	SCAddress                []byte
	SCCode                   []byte
	SCCodeSize               uint64
	CallFunction             string
	VMType                   []byte
	ReadOnlyFlag             bool
	VerifyCode               bool
	CurrentBreakpointValue   arwen.BreakpointValue
	PointsUsed               uint64
	InstanceCtxID            int
	MemLoadResult            []byte
	MemLoadMultipleResult    [][]byte
	FailCryptoAPI            bool
	FailElrondAPI            bool
	FailElrondSyncExecAPI    bool
	FailBigIntAPI            bool
	FailBigFloatAPI          bool
	FailManagedBuffersAPI    bool
	AsyncCallInfo            *arwen.AsyncCallInfo
	RunningInstances         uint64
	CurrentTxHash            []byte
	OriginalTxHash           []byte
	TraceGasEnabled          bool
	GasTrace                 map[string]map[string][]uint64
	SameContractOnStackCount uint64
	HasFunctionResult        bool
}

// InitState mocked method
func (r *RuntimeContextMock) InitState() {
}

// ReplaceInstanceBuilder mocked method()
func (r *RuntimeContextMock) ReplaceInstanceBuilder(_ executor.InstanceBuilder) {
}

// StartWasmerInstance mocked method
func (r *RuntimeContextMock) StartWasmerInstance(_ []byte, _ uint64, _ bool) error {
	if r.Err != nil {
		return r.Err
	}
	return nil
}

// SetCaching mocked method
func (r *RuntimeContextMock) SetCaching(_ bool) {
}

// InitStateFromContractCallInput mocked method
func (r *RuntimeContextMock) InitStateFromContractCallInput(_ *vmcommon.ContractCallInput) {
}

// PushState mocked method
func (r *RuntimeContextMock) PushState() {
}

// PopSetActiveState mocked method
func (r *RuntimeContextMock) PopSetActiveState() {
}

// PopDiscard mocked method
func (r *RuntimeContextMock) PopDiscard() {
}

// MustVerifyNextContractCode mocked method
func (r *RuntimeContextMock) MustVerifyNextContractCode() {
}

// ClearStateStack mocked method
func (r *RuntimeContextMock) ClearStateStack() {
}

// PushInstance mocked method
func (r *RuntimeContextMock) PushInstance() {
}

// PopInstance mocked method
func (r *RuntimeContextMock) PopInstance() {
}

// IsWarmInstance mocked method
func (r *RuntimeContextMock) IsWarmInstance() bool {
	return false
}

// ResetWarmInstance mocked method
func (r *RuntimeContextMock) ResetWarmInstance() {
}

// RunningInstancesCount mocked method
func (r *RuntimeContextMock) RunningInstancesCount() uint64 {
	return r.RunningInstances
}

// SetMaxInstanceCount mocked method
func (r *RuntimeContextMock) SetMaxInstanceCount(uint64) {
}

// ClearInstanceStack mocked method
func (r *RuntimeContextMock) ClearInstanceStack() {
}

// GetVMType mocked method
func (r *RuntimeContextMock) GetVMType() []byte {
	return r.VMType
}

// GetVMInput mocked method
func (r *RuntimeContextMock) GetVMInput() *vmcommon.ContractCallInput {
	return r.VMInput
}

// SetVMInput mocked method
func (r *RuntimeContextMock) SetVMInput(vmInput *vmcommon.ContractCallInput) {
	r.VMInput = vmInput
}

// CountSameContractInstancesOnStack mocked method
func (r *RuntimeContextMock) CountSameContractInstancesOnStack(address []byte) uint64 {
	return r.SameContractOnStackCount
}

// GetSCAddress mocked method
func (r *RuntimeContextMock) GetContextAddress() []byte {
	return r.SCAddress
}

// SetCodeAddress mocked method
func (r *RuntimeContextMock) SetCodeAddress(scAddress []byte) {
	r.SCAddress = scAddress
}

// GetSCCode mocked method
func (r *RuntimeContextMock) GetSCCode() ([]byte, error) {
	return r.SCCode, r.Err
}

// GetSCCodeSize mocked method
func (r *RuntimeContextMock) GetSCCodeSize() uint64 {
	return r.SCCodeSize
}

// Function mocked method
func (r *RuntimeContextMock) Function() string {
	return r.CallFunction
}

// Arguments mocked method
func (r *RuntimeContextMock) Arguments() [][]byte {
	return r.VMInput.Arguments
}

// GetCurrentTxHash mocked method
func (r *RuntimeContextMock) GetCurrentTxHash() []byte {
	return r.CurrentTxHash
}

// GetOriginalTxHash mocked method
func (r *RuntimeContextMock) GetOriginalTxHash() []byte {
	return r.OriginalTxHash
}

// ExtractCodeUpgradeFromArgs mocked method
func (r *RuntimeContextMock) ExtractCodeUpgradeFromArgs() ([]byte, []byte, error) {
	arguments := r.VMInput.Arguments
	if len(arguments) < 2 {
		panic("ExtractCodeUpgradeFromArgs: bad test setup")
	}

	return r.VMInput.Arguments[0], r.VMInput.Arguments[1], nil
}

// SignalExit mocked method
func (r *RuntimeContextMock) SignalExit(_ int) {
}

// SignalUserError mocked method
func (r *RuntimeContextMock) SignalUserError(_ string) {
}

// SetRuntimeBreakpointValue mocked method
func (r *RuntimeContextMock) SetRuntimeBreakpointValue(_ arwen.BreakpointValue) {
}

// GetRuntimeBreakpointValue mocked method
func (r *RuntimeContextMock) GetRuntimeBreakpointValue() arwen.BreakpointValue {
	return r.CurrentBreakpointValue
}

// PrepareLegacyAsyncCall mocked method
func (r *RuntimeContextMock) PrepareLegacyAsyncCall(address []byte, data []byte, value []byte) error {
	return r.Err
}

// VerifyContractCode mocked method
func (r *RuntimeContextMock) VerifyContractCode() error {
	return r.Err
}

// GetPointsUsed mocked method
func (r *RuntimeContextMock) GetPointsUsed() uint64 {
	return r.PointsUsed
}

// SetPointsUsed mocked method
func (r *RuntimeContextMock) SetPointsUsed(gasPoints uint64) {
	r.PointsUsed = gasPoints
}

// ReadOnly mocked method
func (r *RuntimeContextMock) ReadOnly() bool {
	return r.ReadOnlyFlag
}

// SetReadOnly mocked method
func (r *RuntimeContextMock) SetReadOnly(readOnly bool) {
	r.ReadOnlyFlag = readOnly
}

// GetInstance mocked method()
func (r *RuntimeContextMock) GetInstance() executor.InstanceHandler {
	return nil
}

// ClearWarmInstanceCache mocked method
func (r *RuntimeContextMock) ClearWarmInstanceCache() {
}

// FunctionNameChecked mocked method
func (r *RuntimeContextMock) FunctionNameChecked() (string, error) {
	if r.Err != nil {
		return "", r.Err
	}
	return "", nil
}

// CallSCFunction mocked method
func (r *RuntimeContextMock) CallSCFunction(functionName string) error {
	return r.Err
}

// MemLoad mocked method
func (r *RuntimeContextMock) MemLoad(_ int32, _ int32) ([]byte, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	return r.MemLoadResult, nil
}

// MemLoadMultiple mocked method
func (r *RuntimeContextMock) MemLoadMultiple(_ int32, _ []int32) ([][]byte, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	return r.MemLoadMultipleResult, nil
}

// MemStore mocked method
func (r *RuntimeContextMock) MemStore(_ int32, _ []byte) error {
	return r.Err
}

// ElrondAPIErrorShouldFailExecution mocked method
func (r *RuntimeContextMock) ElrondAPIErrorShouldFailExecution() bool {
	return r.FailElrondAPI
}

// ElrondSyncExecAPIErrorShouldFailExecution mocked method
func (r *RuntimeContextMock) ElrondSyncExecAPIErrorShouldFailExecution() bool {
	return r.FailElrondSyncExecAPI
}

// CryptoAPIErrorShouldFailExecution mocked method
func (r *RuntimeContextMock) CryptoAPIErrorShouldFailExecution() bool {
	return r.FailCryptoAPI
}

// BigIntAPIErrorShouldFailExecution mocked method
func (r *RuntimeContextMock) BigIntAPIErrorShouldFailExecution() bool {
	return r.FailBigIntAPI
}

// BigFloatAPIErrorShouldFailExecution mocked method
func (r *RuntimeContextMock) BigFloatAPIErrorShouldFailExecution() bool {
	return r.FailBigFloatAPI
}

// ManagedBufferAPIErrorShouldFailExecution mocked method
func (r *RuntimeContextMock) ManagedBufferAPIErrorShouldFailExecution() bool {
	return r.FailManagedBuffersAPI
}

// FailExecution mocked method
func (r *RuntimeContextMock) FailExecution(_ error) {
}

// AddAsyncContextCall mocked method
func (r *RuntimeContextMock) AddAsyncContextCall(_ []byte, _ *arwen.AsyncGeneratedCall) error {
	return r.Err
}

// SetCustomCallFunction mocked method
func (r *RuntimeContextMock) SetCustomCallFunction(_ string) {
}

// IsFunctionImported mocked method
func (r *RuntimeContextMock) IsFunctionImported(_ string) bool {
	return true
}

// AddError mocked method
func (r *RuntimeContextMock) AddError(_ error, _ ...string) {
}

// GetAllErrors mocked method
func (r *RuntimeContextMock) GetAllErrors() error {
	return nil
}

// ValidateCallbackName mocked method
func (r *RuntimeContextMock) ValidateCallbackName(callbackName string) error {
	return nil
}

// HasFunction mocked method
func (r *RuntimeContextMock) HasFunction(functionName string) bool {
	return r.HasFunctionResult
}

// GetPrevTxHash mocked method
func (r *RuntimeContextMock) GetPrevTxHash() []byte {
	return nil
}

// CleanInstance mocked method
func (r *RuntimeContextMock) CleanInstance() {
}

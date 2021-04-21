package mock

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

var _ RuntimeContextWrapper = (*runtimeContextWrapper)(nil)

type RuntimeContextWrapper interface {
	arwen.RuntimeContext
}

type runtimeContextWrapper struct {
	runtimeContext arwen.RuntimeContext

	InitStateFromContractCallInputFunc            func(input *vmcommon.ContractCallInput)
	SetCustomCallFunctionFunc                     func(callFunction string)
	GetVMInputFunc                                func() *vmcommon.VMInput
	SetVMInputFunc                                func(vmInput *vmcommon.VMInput)
	GetSCAddressFunc                              func() []byte
	SetSCAddressFunc                              func(scAddress []byte)
	GetSCCodeFunc                                 func() ([]byte, error)
	GetSCCodeSizeFunc                             func() uint64
	GetVMTypeFunc                                 func() []byte
	FunctionFunc                                  func() string
	ArgumentsFunc                                 func() [][]byte
	GetCurrentTxHashFunc                          func() []byte
	GetOriginalTxHashFunc                         func() []byte
	ExtractCodeUpgradeFromArgsFunc                func() ([]byte, []byte, error)
	SignalUserErrorFunc                           func(message string)
	FailExecutionFunc                             func(err error)
	MustVerifyNextContractCodeFunc                func()
	SetRuntimeBreakpointValueFunc                 func(value arwen.BreakpointValue)
	GetRuntimeBreakpointValueFunc                 func() arwen.BreakpointValue
	IsContractOnTheStackFunc                      func(address []byte) bool
	GetAsyncCallInfoFunc                          func() *arwen.AsyncCallInfo
	SetAsyncCallInfoFunc                          func(asyncCallInfo *arwen.AsyncCallInfo)
	AddAsyncContextCallFunc                       func(contextIdentifier []byte, asyncCall *arwen.AsyncGeneratedCall) error
	GetAsyncContextInfoFunc                       func() *arwen.AsyncContextInfo
	GetAsyncContextFunc                           func(contextIdentifier []byte) (*arwen.AsyncContext, error)
	RunningInstancesCountFunc                     func() uint64
	IsFunctionImportedFunc                        func(name string) bool
	IsWarmInstanceFunc                            func() bool
	ResetWarmInstanceFunc                         func()
	ReadOnlyFunc                                  func() bool
	SetReadOnlyFunc                               func(readOnly bool)
	StartWasmerInstanceFunc                       func(contract []byte, gasLimit uint64, newCode bool) error
	CleanWasmerInstanceFunc                       func()
	SetMaxInstanceCountFunc                       func(maxInstances uint64)
	VerifyContractCodeFunc                        func() error
	GetInstanceExportsFunc                        func() wasmer.ExportsMap
	GetInitFunctionFunc                           func() wasmer.ExportedFunctionCallback
	GetFunctionToCallFunc                         func() (wasmer.ExportedFunctionCallback, error)
	GetPointsUsedFunc                             func() uint64
	SetPointsUsedFunc                             func(gasPoints uint64)
	MemStoreFunc                                  func(offset int32, data []byte) error
	MemLoadFunc                                   func(offset int32, length int32) ([]byte, error)
	MemLoadMultipleFunc                           func(offset int32, lengths []int32) ([][]byte, error)
	ElrondAPIErrorShouldFailExecutionFunc         func() bool
	ElrondSyncExecAPIErrorShouldFailExecutionFunc func() bool
	CryptoAPIErrorShouldFailExecutionFunc         func() bool
	BigIntAPIErrorShouldFailExecutionFunc         func() bool
	ExecuteAsyncCallFunc                          func(address []byte, data []byte, value []byte) error
	ReplaceInstanceBuilderFunc                    func(builder arwen.InstanceBuilder)

	InitStateFunc         func()
	PushStateFunc         func()
	PopSetActiveStateFunc func()
	PopDiscardFunc        func()
	ClearStateStackFunc   func()
}

func NewRuntimeContextWrapper(inputRuntimeContext *arwen.RuntimeContext) *runtimeContextWrapper {

	runtimeWrapper := &runtimeContextWrapper{runtimeContext: *inputRuntimeContext}

	/*
		default implementations delegate to wrapped context
	*/

	runtimeWrapper.InitStateFromContractCallInputFunc = func(input *vmcommon.ContractCallInput) {
		runtimeWrapper.runtimeContext.InitStateFromContractCallInput(input)
	}

	runtimeWrapper.SetCustomCallFunctionFunc = func(callFunction string) {
		runtimeWrapper.runtimeContext.SetCustomCallFunction(callFunction)
	}

	runtimeWrapper.GetVMInputFunc = func() *vmcommon.VMInput {
		return runtimeWrapper.runtimeContext.GetVMInput()
	}

	runtimeWrapper.SetVMInputFunc = func(vmInput *vmcommon.VMInput) {
		runtimeWrapper.runtimeContext.SetVMInput(vmInput)
	}

	runtimeWrapper.GetSCAddressFunc = func() []byte {
		return runtimeWrapper.runtimeContext.GetSCAddress()
	}

	runtimeWrapper.SetSCAddressFunc = func(scAddress []byte) {
		runtimeWrapper.runtimeContext.SetSCAddress(scAddress)
	}

	runtimeWrapper.GetSCCodeFunc = func() ([]byte, error) {
		return runtimeWrapper.runtimeContext.GetSCCode()
	}

	runtimeWrapper.GetSCCodeSizeFunc = func() uint64 {
		return runtimeWrapper.runtimeContext.GetSCCodeSize()
	}

	runtimeWrapper.GetVMTypeFunc = func() []byte {
		return runtimeWrapper.runtimeContext.GetVMType()
	}

	runtimeWrapper.FunctionFunc = func() string {
		return runtimeWrapper.runtimeContext.Function()
	}

	runtimeWrapper.ArgumentsFunc = func() [][]byte {
		return runtimeWrapper.runtimeContext.Arguments()
	}

	runtimeWrapper.GetCurrentTxHashFunc = func() []byte {
		return runtimeWrapper.runtimeContext.GetCurrentTxHash()
	}

	runtimeWrapper.GetOriginalTxHashFunc = func() []byte {
		return runtimeWrapper.runtimeContext.GetOriginalTxHash()
	}

	runtimeWrapper.ExtractCodeUpgradeFromArgsFunc = func() ([]byte, []byte, error) {
		return runtimeWrapper.runtimeContext.ExtractCodeUpgradeFromArgs()
	}

	runtimeWrapper.SignalUserErrorFunc = func(message string) {
		runtimeWrapper.runtimeContext.SignalUserError(message)
	}

	runtimeWrapper.FailExecutionFunc = func(err error) {
		runtimeWrapper.runtimeContext.FailExecution(err)
	}

	runtimeWrapper.MustVerifyNextContractCodeFunc = func() {
		runtimeWrapper.runtimeContext.MustVerifyNextContractCode()
	}

	runtimeWrapper.SetRuntimeBreakpointValueFunc = func(value arwen.BreakpointValue) {
		runtimeWrapper.runtimeContext.SetRuntimeBreakpointValue(value)
	}

	runtimeWrapper.GetRuntimeBreakpointValueFunc = func() arwen.BreakpointValue {
		return runtimeWrapper.runtimeContext.GetRuntimeBreakpointValue()
	}

	runtimeWrapper.IsContractOnTheStackFunc = func(address []byte) bool {
		return runtimeWrapper.runtimeContext.IsContractOnTheStack(address)
	}

	runtimeWrapper.GetAsyncCallInfoFunc = func() *arwen.AsyncCallInfo {
		return runtimeWrapper.runtimeContext.GetAsyncCallInfo()
	}

	runtimeWrapper.SetAsyncCallInfoFunc = func(asyncCallInfo *arwen.AsyncCallInfo) {
		runtimeWrapper.runtimeContext.SetAsyncCallInfo(asyncCallInfo)
	}

	runtimeWrapper.AddAsyncContextCallFunc = func(contextIdentifier []byte, asyncCall *arwen.AsyncGeneratedCall) error {
		return runtimeWrapper.runtimeContext.AddAsyncContextCall(contextIdentifier, asyncCall)
	}

	runtimeWrapper.GetAsyncContextInfoFunc = func() *arwen.AsyncContextInfo {
		return runtimeWrapper.runtimeContext.GetAsyncContextInfo()
	}

	runtimeWrapper.GetAsyncContextFunc = func(contextIdentifier []byte) (*arwen.AsyncContext, error) {
		return runtimeWrapper.runtimeContext.GetAsyncContext(contextIdentifier)
	}

	runtimeWrapper.RunningInstancesCountFunc = func() uint64 {
		return runtimeWrapper.runtimeContext.RunningInstancesCount()
	}

	runtimeWrapper.IsFunctionImportedFunc = func(name string) bool {
		return runtimeWrapper.runtimeContext.IsFunctionImported(name)
	}

	runtimeWrapper.IsWarmInstanceFunc = func() bool {
		return runtimeWrapper.runtimeContext.IsWarmInstance()
	}

	runtimeWrapper.ResetWarmInstanceFunc = func() {
		runtimeWrapper.runtimeContext.ResetWarmInstance()
	}

	runtimeWrapper.ReadOnlyFunc = func() bool {
		return runtimeWrapper.runtimeContext.ReadOnly()
	}

	runtimeWrapper.SetReadOnlyFunc = func(readOnly bool) {
		runtimeWrapper.runtimeContext.SetReadOnly(readOnly)
	}

	runtimeWrapper.StartWasmerInstanceFunc = func(contract []byte, gasLimit uint64, newCode bool) error {
		return runtimeWrapper.runtimeContext.StartWasmerInstance(contract, gasLimit, newCode)
	}

	runtimeWrapper.CleanWasmerInstanceFunc = func() {
		runtimeWrapper.runtimeContext.CleanWasmerInstance()
	}

	runtimeWrapper.SetMaxInstanceCountFunc = func(maxInstances uint64) {
		runtimeWrapper.runtimeContext.SetMaxInstanceCount(maxInstances)
	}

	runtimeWrapper.VerifyContractCodeFunc = func() error {
		return runtimeWrapper.runtimeContext.VerifyContractCode()
	}

	runtimeWrapper.GetInstanceExportsFunc = func() wasmer.ExportsMap {
		return runtimeWrapper.runtimeContext.GetInstanceExports()
	}

	runtimeWrapper.GetInitFunctionFunc = func() wasmer.ExportedFunctionCallback {
		return runtimeWrapper.runtimeContext.GetInitFunction()
	}

	runtimeWrapper.GetFunctionToCallFunc = func() (wasmer.ExportedFunctionCallback, error) {
		return runtimeWrapper.runtimeContext.GetFunctionToCall()
	}

	runtimeWrapper.GetPointsUsedFunc = func() uint64 {
		return runtimeWrapper.runtimeContext.GetPointsUsed()
	}

	runtimeWrapper.SetPointsUsedFunc = func(gasPoints uint64) {
		runtimeWrapper.runtimeContext.SetPointsUsed(gasPoints)
	}

	runtimeWrapper.MemStoreFunc = func(offset int32, data []byte) error {
		return runtimeWrapper.runtimeContext.MemStore(offset, data)
	}

	runtimeWrapper.MemLoadFunc = func(offset int32, length int32) ([]byte, error) {
		return runtimeWrapper.runtimeContext.MemLoad(offset, length)
	}

	runtimeWrapper.MemLoadMultipleFunc = func(offset int32, lengths []int32) ([][]byte, error) {
		return runtimeWrapper.runtimeContext.MemLoadMultiple(offset, lengths)
	}

	runtimeWrapper.ElrondAPIErrorShouldFailExecutionFunc = func() bool {
		return runtimeWrapper.runtimeContext.ElrondAPIErrorShouldFailExecution()
	}

	runtimeWrapper.ElrondSyncExecAPIErrorShouldFailExecutionFunc = func() bool {
		return runtimeWrapper.runtimeContext.ElrondSyncExecAPIErrorShouldFailExecution()
	}

	runtimeWrapper.CryptoAPIErrorShouldFailExecutionFunc = func() bool {
		return runtimeWrapper.runtimeContext.CryptoAPIErrorShouldFailExecution()
	}

	runtimeWrapper.BigIntAPIErrorShouldFailExecutionFunc = func() bool {
		return runtimeWrapper.runtimeContext.BigIntAPIErrorShouldFailExecution()
	}

	runtimeWrapper.ExecuteAsyncCallFunc = func(address []byte, data []byte, value []byte) error {
		return runtimeWrapper.runtimeContext.ExecuteAsyncCall(address, data, value)
	}

	runtimeWrapper.ReplaceInstanceBuilderFunc = func(builder arwen.InstanceBuilder) {
		runtimeWrapper.runtimeContext.ReplaceInstanceBuilder(builder)
	}

	runtimeWrapper.InitStateFunc = func() {
		runtimeWrapper.runtimeContext.InitState()
	}

	runtimeWrapper.PushStateFunc = func() {
		runtimeWrapper.runtimeContext.PushState()
	}

	runtimeWrapper.PopSetActiveStateFunc = func() {
		runtimeWrapper.runtimeContext.PopSetActiveState()
	}

	runtimeWrapper.PopDiscardFunc = func() {
		runtimeWrapper.runtimeContext.PopDiscard()
	}

	runtimeWrapper.ClearStateStackFunc = func() {
		runtimeWrapper.runtimeContext.ClearStateStack()
	}

	return runtimeWrapper
}

// gets the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) GetWrappedRuntimeContext() arwen.RuntimeContext {
	return contextWrapper.runtimeContext
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) InitStateFromContractCallInput(input *vmcommon.ContractCallInput) {
	contextWrapper.InitStateFromContractCallInputFunc(input)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) SetCustomCallFunction(callFunction string) {
	contextWrapper.SetCustomCallFunctionFunc(callFunction)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) GetVMInput() *vmcommon.VMInput {
	return contextWrapper.GetVMInputFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) SetVMInput(vmInput *vmcommon.VMInput) {
	contextWrapper.SetVMInputFunc(vmInput)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) GetSCAddress() []byte {
	return contextWrapper.GetSCAddressFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) SetSCAddress(scAddress []byte) {
	contextWrapper.SetSCAddressFunc(scAddress)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) GetSCCode() ([]byte, error) {
	return contextWrapper.GetSCCodeFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) GetSCCodeSize() uint64 {
	return contextWrapper.GetSCCodeSizeFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) GetVMType() []byte {
	return contextWrapper.GetVMTypeFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) Function() string {
	return contextWrapper.FunctionFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) Arguments() [][]byte {
	return contextWrapper.ArgumentsFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) GetCurrentTxHash() []byte {
	return contextWrapper.GetCurrentTxHashFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) GetOriginalTxHash() []byte {
	return contextWrapper.GetOriginalTxHashFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) ExtractCodeUpgradeFromArgs() ([]byte, []byte, error) {
	return contextWrapper.ExtractCodeUpgradeFromArgsFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) SignalUserError(message string) {
	contextWrapper.SignalUserErrorFunc(message)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) FailExecution(err error) {
	contextWrapper.FailExecutionFunc(err)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) MustVerifyNextContractCode() {
	contextWrapper.MustVerifyNextContractCodeFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) SetRuntimeBreakpointValue(value arwen.BreakpointValue) {
	contextWrapper.SetRuntimeBreakpointValueFunc(value)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) GetRuntimeBreakpointValue() arwen.BreakpointValue {
	return contextWrapper.GetRuntimeBreakpointValueFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) IsContractOnTheStack(address []byte) bool {
	return contextWrapper.IsContractOnTheStackFunc(address)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) GetAsyncCallInfo() *arwen.AsyncCallInfo {
	return contextWrapper.GetAsyncCallInfoFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) SetAsyncCallInfo(asyncCallInfo *arwen.AsyncCallInfo) {
	contextWrapper.SetAsyncCallInfoFunc(asyncCallInfo)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) AddAsyncContextCall(contextIdentifier []byte, asyncCall *arwen.AsyncGeneratedCall) error {
	return contextWrapper.AddAsyncContextCallFunc(contextIdentifier, asyncCall)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) GetAsyncContextInfo() *arwen.AsyncContextInfo {
	return contextWrapper.GetAsyncContextInfoFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) GetAsyncContext(contextIdentifier []byte) (*arwen.AsyncContext, error) {
	return contextWrapper.GetAsyncContextFunc(contextIdentifier)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) RunningInstancesCount() uint64 {
	return contextWrapper.RunningInstancesCountFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) IsFunctionImported(name string) bool {
	return contextWrapper.IsFunctionImportedFunc(name)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) IsWarmInstance() bool {
	return contextWrapper.IsWarmInstanceFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) ResetWarmInstance() {
	contextWrapper.ResetWarmInstanceFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) ReadOnly() bool {
	return contextWrapper.ReadOnlyFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) SetReadOnly(readOnly bool) {
	contextWrapper.SetReadOnlyFunc(readOnly)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) StartWasmerInstance(contract []byte, gasLimit uint64, newCode bool) error {
	return contextWrapper.StartWasmerInstanceFunc(contract, gasLimit, newCode)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) CleanWasmerInstance() {
	contextWrapper.CleanWasmerInstanceFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) SetMaxInstanceCount(maxInstances uint64) {
	contextWrapper.SetMaxInstanceCountFunc(maxInstances)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) VerifyContractCode() error {
	return contextWrapper.VerifyContractCodeFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) GetInstanceExports() wasmer.ExportsMap {
	return contextWrapper.GetInstanceExportsFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) GetInitFunction() wasmer.ExportedFunctionCallback {
	return contextWrapper.GetInitFunctionFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) GetFunctionToCall() (wasmer.ExportedFunctionCallback, error) {
	return contextWrapper.GetFunctionToCallFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) GetPointsUsed() uint64 {
	return contextWrapper.GetPointsUsedFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) SetPointsUsed(gasPoints uint64) {
	contextWrapper.SetPointsUsedFunc(gasPoints)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) MemStore(offset int32, data []byte) error {
	return contextWrapper.MemStoreFunc(offset, data)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) MemLoad(offset int32, length int32) ([]byte, error) {
	return contextWrapper.MemLoadFunc(offset, length)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) MemLoadMultiple(offset int32, lengths []int32) ([][]byte, error) {
	return contextWrapper.MemLoadMultipleFunc(offset, lengths)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) ElrondAPIErrorShouldFailExecution() bool {
	return contextWrapper.ElrondAPIErrorShouldFailExecutionFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) ElrondSyncExecAPIErrorShouldFailExecution() bool {
	return contextWrapper.ElrondSyncExecAPIErrorShouldFailExecutionFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) CryptoAPIErrorShouldFailExecution() bool {
	return contextWrapper.CryptoAPIErrorShouldFailExecutionFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) BigIntAPIErrorShouldFailExecution() bool {
	return contextWrapper.BigIntAPIErrorShouldFailExecutionFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) ExecuteAsyncCall(address []byte, data []byte, value []byte) error {
	return contextWrapper.ExecuteAsyncCallFunc(address, data, value)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) ReplaceInstanceBuilder(builder arwen.InstanceBuilder) {
	contextWrapper.ReplaceInstanceBuilderFunc(builder)
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) InitState() {
	contextWrapper.InitStateFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) PushState() {
	contextWrapper.PushStateFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) PopSetActiveState() {
	contextWrapper.PopSetActiveStateFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) PopDiscard() {
	contextWrapper.PopDiscardFunc()
}

// calls corresponding xxxFunc function, that by default in turn calls the original method of the wrapped RuntimeContext
func (contextWrapper *runtimeContextWrapper) ClearStateStack() {
	contextWrapper.ClearStateStackFunc()
}

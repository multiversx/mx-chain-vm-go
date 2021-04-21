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

func (contextWrapper *runtimeContextWrapper) GetWrappedRuntimeContext() arwen.RuntimeContext {
	return contextWrapper.runtimeContext
}

func (contextWrapper *runtimeContextWrapper) InitStateFromContractCallInput(input *vmcommon.ContractCallInput) {
	contextWrapper.InitStateFromContractCallInputFunc(input)
}

func (contextWrapper *runtimeContextWrapper) SetCustomCallFunction(callFunction string) {
	contextWrapper.SetCustomCallFunctionFunc(callFunction)
}

func (contextWrapper *runtimeContextWrapper) GetVMInput() *vmcommon.VMInput {
	return contextWrapper.GetVMInputFunc()
}

func (contextWrapper *runtimeContextWrapper) SetVMInput(vmInput *vmcommon.VMInput) {
	contextWrapper.SetVMInputFunc(vmInput)
}

func (contextWrapper *runtimeContextWrapper) GetSCAddress() []byte {
	return contextWrapper.GetSCAddressFunc()
}

func (contextWrapper *runtimeContextWrapper) SetSCAddress(scAddress []byte) {
	contextWrapper.SetSCAddressFunc(scAddress)
}

func (contextWrapper *runtimeContextWrapper) GetSCCode() ([]byte, error) {
	return contextWrapper.GetSCCodeFunc()
}

func (contextWrapper *runtimeContextWrapper) GetSCCodeSize() uint64 {
	return contextWrapper.GetSCCodeSizeFunc()
}

func (contextWrapper *runtimeContextWrapper) GetVMType() []byte {
	return contextWrapper.GetVMTypeFunc()
}

func (contextWrapper *runtimeContextWrapper) Function() string {
	return contextWrapper.FunctionFunc()
}

func (contextWrapper *runtimeContextWrapper) Arguments() [][]byte {
	return contextWrapper.ArgumentsFunc()
}

func (contextWrapper *runtimeContextWrapper) GetCurrentTxHash() []byte {
	return contextWrapper.GetCurrentTxHashFunc()
}

func (contextWrapper *runtimeContextWrapper) GetOriginalTxHash() []byte {
	return contextWrapper.GetOriginalTxHashFunc()
}

func (contextWrapper *runtimeContextWrapper) ExtractCodeUpgradeFromArgs() ([]byte, []byte, error) {
	return contextWrapper.ExtractCodeUpgradeFromArgsFunc()
}

func (contextWrapper *runtimeContextWrapper) SignalUserError(message string) {
	contextWrapper.SignalUserErrorFunc(message)
}

func (contextWrapper *runtimeContextWrapper) FailExecution(err error) {
	contextWrapper.FailExecutionFunc(err)
}

func (contextWrapper *runtimeContextWrapper) MustVerifyNextContractCode() {
	contextWrapper.MustVerifyNextContractCodeFunc()
}

func (contextWrapper *runtimeContextWrapper) SetRuntimeBreakpointValue(value arwen.BreakpointValue) {
	contextWrapper.SetRuntimeBreakpointValueFunc(value)
}

func (contextWrapper *runtimeContextWrapper) GetRuntimeBreakpointValue() arwen.BreakpointValue {
	return contextWrapper.GetRuntimeBreakpointValueFunc()
}

func (contextWrapper *runtimeContextWrapper) IsContractOnTheStack(address []byte) bool {
	return contextWrapper.IsContractOnTheStackFunc(address)
}

func (contextWrapper *runtimeContextWrapper) GetAsyncCallInfo() *arwen.AsyncCallInfo {
	return contextWrapper.GetAsyncCallInfoFunc()
}

func (contextWrapper *runtimeContextWrapper) SetAsyncCallInfo(asyncCallInfo *arwen.AsyncCallInfo) {
	contextWrapper.SetAsyncCallInfoFunc(asyncCallInfo)
}

func (contextWrapper *runtimeContextWrapper) AddAsyncContextCall(contextIdentifier []byte, asyncCall *arwen.AsyncGeneratedCall) error {
	return contextWrapper.AddAsyncContextCallFunc(contextIdentifier, asyncCall)
}

func (contextWrapper *runtimeContextWrapper) GetAsyncContextInfo() *arwen.AsyncContextInfo {
	return contextWrapper.GetAsyncContextInfoFunc()
}

func (contextWrapper *runtimeContextWrapper) GetAsyncContext(contextIdentifier []byte) (*arwen.AsyncContext, error) {
	return contextWrapper.GetAsyncContextFunc(contextIdentifier)
}

func (contextWrapper *runtimeContextWrapper) RunningInstancesCount() uint64 {
	return contextWrapper.RunningInstancesCountFunc()
}

func (contextWrapper *runtimeContextWrapper) IsFunctionImported(name string) bool {
	return contextWrapper.IsFunctionImportedFunc(name)
}

func (contextWrapper *runtimeContextWrapper) IsWarmInstance() bool {
	return contextWrapper.IsWarmInstanceFunc()
}

func (contextWrapper *runtimeContextWrapper) ResetWarmInstance() {
	contextWrapper.ResetWarmInstanceFunc()
}

func (contextWrapper *runtimeContextWrapper) ReadOnly() bool {
	return contextWrapper.ReadOnlyFunc()
}

func (contextWrapper *runtimeContextWrapper) SetReadOnly(readOnly bool) {
	contextWrapper.SetReadOnlyFunc(readOnly)
}

func (contextWrapper *runtimeContextWrapper) StartWasmerInstance(contract []byte, gasLimit uint64, newCode bool) error {
	return contextWrapper.StartWasmerInstanceFunc(contract, gasLimit, newCode)
}

func (contextWrapper *runtimeContextWrapper) CleanWasmerInstance() {
	contextWrapper.CleanWasmerInstanceFunc()
}

func (contextWrapper *runtimeContextWrapper) SetMaxInstanceCount(maxInstances uint64) {
	contextWrapper.SetMaxInstanceCountFunc(maxInstances)
}

func (contextWrapper *runtimeContextWrapper) VerifyContractCode() error {
	return contextWrapper.VerifyContractCodeFunc()
}

func (contextWrapper *runtimeContextWrapper) GetInstanceExports() wasmer.ExportsMap {
	return contextWrapper.GetInstanceExportsFunc()
}

func (contextWrapper *runtimeContextWrapper) GetInitFunction() wasmer.ExportedFunctionCallback {
	return contextWrapper.GetInitFunctionFunc()
}

func (contextWrapper *runtimeContextWrapper) GetFunctionToCall() (wasmer.ExportedFunctionCallback, error) {
	return contextWrapper.GetFunctionToCallFunc()
}

func (contextWrapper *runtimeContextWrapper) GetPointsUsed() uint64 {
	return contextWrapper.GetPointsUsedFunc()
}

func (contextWrapper *runtimeContextWrapper) SetPointsUsed(gasPoints uint64) {
	contextWrapper.SetPointsUsedFunc(gasPoints)
}

func (contextWrapper *runtimeContextWrapper) MemStore(offset int32, data []byte) error {
	return contextWrapper.MemStoreFunc(offset, data)
}

func (contextWrapper *runtimeContextWrapper) MemLoad(offset int32, length int32) ([]byte, error) {
	return contextWrapper.MemLoadFunc(offset, length)
}

func (contextWrapper *runtimeContextWrapper) MemLoadMultiple(offset int32, lengths []int32) ([][]byte, error) {
	return contextWrapper.MemLoadMultipleFunc(offset, lengths)
}

func (contextWrapper *runtimeContextWrapper) ElrondAPIErrorShouldFailExecution() bool {
	return contextWrapper.ElrondAPIErrorShouldFailExecutionFunc()
}

func (contextWrapper *runtimeContextWrapper) ElrondSyncExecAPIErrorShouldFailExecution() bool {
	return contextWrapper.ElrondSyncExecAPIErrorShouldFailExecutionFunc()
}

func (contextWrapper *runtimeContextWrapper) CryptoAPIErrorShouldFailExecution() bool {
	return contextWrapper.CryptoAPIErrorShouldFailExecutionFunc()
}

func (contextWrapper *runtimeContextWrapper) BigIntAPIErrorShouldFailExecution() bool {
	return contextWrapper.BigIntAPIErrorShouldFailExecutionFunc()
}

func (contextWrapper *runtimeContextWrapper) ExecuteAsyncCall(address []byte, data []byte, value []byte) error {
	return contextWrapper.ExecuteAsyncCallFunc(address, data, value)
}

func (contextWrapper *runtimeContextWrapper) ReplaceInstanceBuilder(builder arwen.InstanceBuilder) {
	contextWrapper.ReplaceInstanceBuilderFunc(builder)
}

func (contextWrapper *runtimeContextWrapper) InitState() {
	contextWrapper.InitStateFunc()
}

func (contextWrapper *runtimeContextWrapper) PushState() {
	contextWrapper.PushStateFunc()
}

func (contextWrapper *runtimeContextWrapper) PopSetActiveState() {
	contextWrapper.PopSetActiveStateFunc()
}

func (contextWrapper *runtimeContextWrapper) PopDiscard() {
	contextWrapper.PopDiscardFunc()
}

func (contextWrapper *runtimeContextWrapper) ClearStateStack() {
	contextWrapper.ClearStateStackFunc()
}

package vmhooks

import (
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

type mockRuntimeContext struct {
	failExecutionCalled           func(err error)
	failExecutionConditionallyCalled func(err error)
	addErrorCalled                func(err error, otherInfo ...string)
	isUnsafeModeCalled            func() bool
}

func (m *mockRuntimeContext) AddError(err error, otherInfo ...string) {
	if m.addErrorCalled != nil {
		m.addErrorCalled(err, otherInfo...)
	}
}

func (m *mockRuntimeContext) IsUnsafeMode() bool {
	if m.isUnsafeModeCalled != nil {
		return m.isUnsafeModeCalled()
	}
	return false
}

func (m *mockRuntimeContext) FailExecution(err error) {
	if m.failExecutionCalled != nil {
		m.failExecutionCalled(err)
	}
}

func (m *mockRuntimeContext) FailExecutionConditionally(err error) {
	if m.failExecutionConditionallyCalled != nil {
		m.failExecutionConditionallyCalled(err)
	}
}

// below are the rest of the methods of the interface, that are not used in these tests

func (m *mockRuntimeContext) PushState()                                         {}
func (m *mockRuntimeContext) PopSetActiveState()                                 {}
func (m *mockRuntimeContext) PopDiscard()                                        {}
func (m *mockRuntimeContext) ClearStateStack()                                   {}
func (m *mockRuntimeContext) GetVMExecutor() executor.Executor                   { return nil }
func (m *mockRuntimeContext) ReplaceVMExecutor(vmExecutor executor.Executor)     {}
func (m *mockRuntimeContext) InitStateFromContractCallInput(input *vmcommon.ContractCallInput) {}
func (m *mockRuntimeContext) SetCustomCallFunction(callFunction string)          {}
func (m *mockRuntimeContext) GetVMInput() *vmcommon.ContractCallInput              { return nil }
func (m *mockRuntimeContext) SetVMInput(vmInput *vmcommon.ContractCallInput)       {}
func (m *mockRuntimeContext) GetContextAddress() []byte                          { return nil }
func (m *mockRuntimeContext) GetOriginalCallerAddress() []byte                   { return nil }
func (m *mockRuntimeContext) SetCodeAddress(scAddress []byte)                    {}
func (m *mockRuntimeContext) GetSCCode() ([]byte, error)                         { return nil, nil }
func (m *mockRuntimeContext) GetSCCodeSize() uint64                              { return 0 }
func (m *mockRuntimeContext) GetVMType() []byte                                  { return nil }
func (m *mockRuntimeContext) FunctionName() string                               { return "" }
func (m *mockRuntimeContext) Arguments() [][]byte                                { return nil }
func (m *mockRuntimeContext) GetCurrentTxHash() []byte                           { return nil }
func (m *mockRuntimeContext) GetOriginalTxHash() []byte                          { return nil }
func (m *mockRuntimeContext) RemoveCodeUpgradeFromArgs()                         {}
func (m *mockRuntimeContext) SignalUserError(message string)                     {}
func (m *mockRuntimeContext) MustVerifyNextContractCode()                        {}
func (m *mockRuntimeContext) SetRuntimeBreakpointValue(value vmhost.BreakpointValue) {}
func (m *mockRuntimeContext) GetRuntimeBreakpointValue() vmhost.BreakpointValue    { return 0 }
func (m *mockRuntimeContext) GetInstanceStackSize() uint64                       { return 0 }
func (m *mockRuntimeContext) CountSameContractInstancesOnStack(address []byte) uint64 { return 0 }
func (m *mockRuntimeContext) IsFunctionImported(name string) bool                { return false }
func (m *mockRuntimeContext) ReadOnly() bool                                     { return false }
func (m *mockRuntimeContext) SetReadOnly(readOnly bool)                          {}
func (m *mockRuntimeContext) SetUnsafeMode(unsafeMode bool)                      {}
func (m *mockRuntimeContext) StartWasmerInstance(contract []byte, gasLimit uint64, newCode bool) error {
	return nil
}
func (m *mockRuntimeContext) ClearWarmInstanceCache()                            {}
func (m *mockRuntimeContext) SetMaxInstanceStackSize(uint64)                     {}
func (m *mockRuntimeContext) VerifyContractCode() error                          { return nil }
func (m *mockRuntimeContext) GetInstance() executor.Instance                     { return nil }
func (m *mockRuntimeContext) GetInstanceTracker() vmhost.InstanceTracker         { return nil }
func (m *mockRuntimeContext) FunctionNameChecked() (string, error)               { return "", nil }
func (m *mockRuntimeContext) CallSCFunction(functionName string) error           { return nil }
func (m *mockRuntimeContext) GetPointsUsed() uint64                              { return 0 }
func (m *mockRuntimeContext) SetPointsUsed(gasPoints uint64)                     {}
func (m *mockRuntimeContext) UseGasBoundedShouldFailExecution() bool             { return false }
func (m *mockRuntimeContext) CleanInstance()                                     {}
func (m *mockRuntimeContext) GetAllErrors() error                                { return nil }
func (m *mockRuntimeContext) ValidateCallbackName(callbackName string) error     { return nil }
func (m *mockRuntimeContext) IsReservedFunctionName(functionName string) bool    { return false }
func (m *mockRuntimeContext) HasFunction(functionName string) bool               { return false }
func (m *mockRuntimeContext) GetPrevTxHash() []byte                              { return nil }
func (m *mockRuntimeContext) EndExecution()                                      {}
func (m *mockRuntimeContext) ValidateInstances() error                           { return nil }
func (m *mockRuntimeContext) InitState()                                         {}
func (m *mockRuntimeContext) IsInterfaceNil() bool                               { return m == nil }

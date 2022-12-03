package executorwrappers

import (
	"fmt"

	"github.com/ElrondNetwork/wasm-vm/executor"
)

// WrapperInstance is a wrapper around an executor instance, which adds the possibility of logging operations.
type WrapperInstance struct {
	logger          ExecutorLogger
	wrappedInstance executor.Instance
}

// GetPointsUsed wraps the call to the underlying instance.
func (inst *WrapperInstance) GetPointsUsed() uint64 {
	result := inst.wrappedInstance.GetPointsUsed()
	inst.logger.LogExecutorEvent(fmt.Sprintf("GetPointsUsed: %d", result))
	return result
}

// SetPointsUsed wraps the call to the underlying instance.
func (inst *WrapperInstance) SetPointsUsed(points uint64) {
	inst.logger.LogExecutorEvent(fmt.Sprintf("SetPointsUsed: %d", points))
	inst.wrappedInstance.SetPointsUsed(points)

}

// SetGasLimit wraps the call to the underlying instance.
func (inst *WrapperInstance) SetGasLimit(gasLimit uint64) {
	inst.logger.LogExecutorEvent(fmt.Sprintf("SetGasLimit: %d", gasLimit))
	inst.wrappedInstance.SetGasLimit(gasLimit)
}

// SetBreakpointValue wraps the call to the underlying instance.
func (inst *WrapperInstance) SetBreakpointValue(value uint64) {
	inst.logger.LogExecutorEvent(fmt.Sprintf("SetBreakpointValue: %d", value))
	inst.wrappedInstance.SetBreakpointValue(value)
}

// GetBreakpointValue wraps the call to the underlying instance.
func (inst *WrapperInstance) GetBreakpointValue() uint64 {
	result := inst.wrappedInstance.GetBreakpointValue()
	inst.logger.LogExecutorEvent(fmt.Sprintf("GetBreakpointValue: %d", result))
	return result
}

// Cache wraps the call to the underlying instance.
func (inst *WrapperInstance) Cache() ([]byte, error) {
	return inst.wrappedInstance.Cache()
}

// Clean wraps the call to the underlying instance.
func (inst *WrapperInstance) Clean() {
	inst.wrappedInstance.Clean()
}

// CallFunction wraps the call to the underlying instance.
func (inst *WrapperInstance) CallFunction(functionName string) error {
	inst.logger.LogExecutorEvent(fmt.Sprintf("CallFunction:(%s)", functionName))
	return inst.wrappedInstance.CallFunction(functionName)
}

// HasFunction wraps the call to the underlying instance.
func (inst *WrapperInstance) HasFunction(functionName string) bool {
	result := inst.wrappedInstance.HasFunction(functionName)
	inst.logger.LogExecutorEvent(fmt.Sprintf("HasFunction(%s): %t", functionName, result))
	return result
}

// GetFunctionNames wraps the call to the underlying instance.
func (inst *WrapperInstance) GetFunctionNames() []string {
	result := inst.wrappedInstance.GetFunctionNames()
	inst.logger.LogExecutorEvent(fmt.Sprintf("GetFunctionNames: %s", result))
	return result
}

// ValidateVoidFunction wraps the call to the underlying instance.
func (inst *WrapperInstance) ValidateVoidFunction(functionName string) error {
	result := inst.wrappedInstance.ValidateVoidFunction(functionName)
	inst.logger.LogExecutorEvent(fmt.Sprintf("ValidateVoidFunction: %s", result))
	return result
}

// HasMemory wraps the call to the underlying instance.
func (inst *WrapperInstance) HasMemory() bool {
	return inst.wrappedInstance.HasMemory()
}

// GetMemory wraps the call to the underlying instance.
func (inst *WrapperInstance) GetMemory() executor.Memory {
	return inst.wrappedInstance.GetMemory()
}

// IsFunctionImported wraps the call to the underlying instance.
func (inst *WrapperInstance) IsFunctionImported(name string) bool {
	result := inst.wrappedInstance.IsFunctionImported(name)
	inst.logger.LogExecutorEvent(fmt.Sprintf("IsFunctionImported(%s): %t", name, result))
	return result
}

// IsInterfaceNil returns true if there is no value under the interface.
func (inst *WrapperInstance) IsInterfaceNil() bool {
	return inst == nil
}

// Reset wraps the call to the underlying instance.
func (inst *WrapperInstance) Reset() bool {
	result := inst.wrappedInstance.Reset()
	inst.logger.LogExecutorEvent(fmt.Sprintf("Reset: %t", result))
	return result
}

// SetVMHooksPtr wraps the call to the underlying instance.
func (inst *WrapperInstance) SetVMHooksPtr(vmHooksPtr uintptr) {
	inst.wrappedInstance.SetVMHooksPtr(vmHooksPtr)
}

// GetVMHooksPtr wraps the call to the underlying instance.
func (inst *WrapperInstance) GetVMHooksPtr() uintptr {
	return inst.wrappedInstance.GetVMHooksPtr()
}

// Id wraps the call to the underlying instance.
func (inst *WrapperInstance) Id() string {
	return inst.wrappedInstance.Id()
}

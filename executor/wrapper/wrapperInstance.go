package executorwrapper

import (
	"fmt"
	"sort"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ executor.InstanceNew = (*WrapperInstance)(nil)

// WrapperInstance is a wrapper around an executor instance, which adds the possibility of logging operations.
type WrapperInstance struct {
	logger          ExecutorLogger
	wrappedInstance executor.InstanceNew
}

// GetPointsUsedNew wraps the call to the underlying instance.
func (inst *WrapperInstance) GetPointsUsedNew() (uint64, error) {
	points, err := inst.wrappedInstance.GetPointsUsedNew()
	inst.logger.LogExecutorEvent(fmt.Sprintf("GetPointsUsed: %d", points))
	return points, err
}

// SetPointsUsedNew wraps the call to the underlying instance.
func (inst *WrapperInstance) SetPointsUsedNew(points uint64) error {
	err := inst.wrappedInstance.SetPointsUsedNew(points)
	inst.logger.LogExecutorEvent(fmt.Sprintf("SetPointsUsed: %d", points))
	return err
}

// Cache wraps the call to the underlying instance.
func (inst *WrapperInstance) Cache() ([]byte, error) {
	return inst.wrappedInstance.Cache()
}

// Close wraps the call to the underlying instance.
func (inst *WrapperInstance) Close() {
	inst.wrappedInstance.Close()
	inst.logger.LogExecutorEvent("Clean: true")
}

// IsAlreadyCleaned wraps the call to the underlying instance.
//
// Not part of executor.InstanceNew; required of anything the VM tracks.
func (inst *WrapperInstance) IsAlreadyCleaned() bool {
	result := instanceIsAlreadyCleaned(inst.wrappedInstance)
	inst.logger.LogExecutorEvent(fmt.Sprintf("IsAlreadyCleaned: %t", result))
	return result
}

// Abandon wraps the call to the underlying instance.
//
// Not part of executor.InstanceNew; required of anything the VM tracks.
func (inst *WrapperInstance) Abandon() {
	abandonInstance(inst.wrappedInstance)
	inst.logger.LogExecutorEvent("Abandon")
}

// Call wraps the call to the underlying instance.
func (inst *WrapperInstance) Call(functionName string, pointsLimit uint64) executor.InstanceCallResult {
	inst.logger.LogExecutorEvent(fmt.Sprintf("CallFunction(%s):", functionName))
	return inst.wrappedInstance.Call(functionName, pointsLimit)
}

// RequestStop wraps the call to the underlying instance.
func (inst *WrapperInstance) RequestStop() {
	inst.wrappedInstance.RequestStop()
}

// HasFunction wraps the call to the underlying instance.
func (inst *WrapperInstance) HasFunction(functionName string) bool {
	result := inst.wrappedInstance.HasFunction(functionName)
	inst.logger.LogExecutorEvent(fmt.Sprintf("HasFunction(%s): %t", functionName, result))
	return result
}

// GetExportedFunctionNames wraps the call to the underlying instance.
func (inst *WrapperInstance) GetExportedFunctionNames() []string {
	result := inst.wrappedInstance.GetExportedFunctionNames()
	sort.Strings(result) // to get consistent logs, function names must be sorted
	inst.logger.LogExecutorEvent(fmt.Sprintf("GetFunctionNames: %s", result))
	return result
}

// IsFunctionImported wraps the call to the underlying instance.
func (inst *WrapperInstance) IsFunctionImported(name string) bool {
	return inst.wrappedInstance.IsFunctionImported(name)
}

// CheckSignatures wraps the call to the underlying instance.
func (inst *WrapperInstance) CheckSignatures() bool {
	result := inst.wrappedInstance.CheckSignatures()
	inst.logger.LogExecutorEvent(fmt.Sprintf("ValidateFunctionArities: %t", result))
	return result
}

// IsInterfaceNil returns true if there is no value under the interface.
func (inst *WrapperInstance) IsInterfaceNil() bool {
	return inst == nil
}

// ResetNew wraps the call to the underlying instance.
func (inst *WrapperInstance) ResetNew() error {
	err := inst.wrappedInstance.ResetNew()
	inst.logger.LogExecutorEvent(fmt.Sprintf("Reset: %t", err == nil))
	return err
}

// ID wraps the call to the underlying instance.
func (inst *WrapperInstance) ID() string {
	return inst.wrappedInstance.ID()
}

// instanceCleanupOps is what the VM requires of an instance beyond executor.InstanceNew.
// Named here rather than imported from vmhost, to keep this package below it.
type instanceCleanupOps interface {
	Abandon()
	IsAlreadyCleaned() bool
}

// instanceIsAlreadyCleaned asks the wrapped instance, or assumes it is live if the backend
// does not offer the question - in which case the host would have rejected it anyway.
func instanceIsAlreadyCleaned(instance executor.InstanceNew) bool {
	cleanup, ok := instance.(instanceCleanupOps)
	if !ok {
		return false
	}

	return cleanup.IsAlreadyCleaned()
}

func abandonInstance(instance executor.InstanceNew) {
	cleanup, ok := instance.(instanceCleanupOps)
	if !ok {
		return
	}

	cleanup.Abandon()
}

package mock

import (
	"errors"
	"testing"

	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/executor"
	"github.com/ElrondNetwork/wasm-vm/wasmer"
)

// InstanceMock is a mock for Wasmer instances; it allows creating mock smart
// contracts within tests, without needing actual WASM smart contracts.
type InstanceMock struct {
	Code            []byte
	Exports         wasmer.ExportsMap
	Points          uint64
	Data            executor.VMHooks
	GasLimit        uint64
	BreakpointValue arwen.BreakpointValue
	Memory          executor.Memory
	Host            arwen.VMHost
	T               testing.TB
	Address         []byte
}

// NewInstanceMock creates a new InstanceMock
func NewInstanceMock(code []byte) *InstanceMock {
	return &InstanceMock{
		Code:            code,
		Exports:         make(wasmer.ExportsMap),
		Points:          0,
		Data:            nil,
		GasLimit:        0,
		BreakpointValue: 0,
		Memory:          NewMemoryMock(),
	}
}

// AddMockMethod adds the provided function as a mocked method to the instance under the specified name.
func (instance *InstanceMock) AddMockMethod(name string, method func() *InstanceMock) {
	wrappedMethod := func(...interface{}) (wasmer.Value, error) {
		instance := method()
		breakpoint := arwen.BreakpointValue(instance.GetBreakpointValue())
		var err error
		if breakpoint != arwen.BreakpointNone {
			err = errors.New(breakpoint.String())
		}
		return wasmer.Void(), err
	}

	instance.Exports[name] = wrappedMethod
}

// SetVMHooks mocked method
func (instance *InstanceMock) SetVMHooks(callbacks executor.VMHooks) {
	instance.Data = callbacks
}

// GetVMHooks mocked method
func (instance *InstanceMock) GetVMHooks() executor.VMHooks {
	return instance.Data
}

// HasMemory mocked method
func (instance *InstanceMock) HasMemory() bool {
	return true
}

// GetPointsUsed mocked method
func (instance *InstanceMock) GetPointsUsed() uint64 {
	return instance.Points
}

// SetPointsUsed mocked method
func (instance *InstanceMock) SetPointsUsed(points uint64) {
	instance.Points = points
}

// SetGasLimit mocked method
func (instance *InstanceMock) SetGasLimit(gasLimit uint64) {
	instance.GasLimit = gasLimit
}

// SetBreakpointValue mocked method
func (instance *InstanceMock) SetBreakpointValue(value uint64) {
	instance.BreakpointValue = arwen.BreakpointValue(value)
}

// GetBreakpointValue mocked method
func (instance *InstanceMock) GetBreakpointValue() uint64 {
	return uint64(instance.BreakpointValue)
}

// Cache mocked method
func (instance *InstanceMock) Cache() ([]byte, error) {
	return instance.Code, nil
}

// Clean mocked method
func (instance *InstanceMock) Clean() {
}

// Reset mocked method
func (instance *InstanceMock) Reset() bool {
	return true
}

// CallFunction mocked method
func (instance *InstanceMock) CallFunction(functionName string) error {
	if function, ok := instance.Exports[functionName]; ok {
		_, err := function()
		return err
	}

	return executor.ErrFuncNotFound
}

// HasFunction mocked method
func (instance *InstanceMock) HasFunction(functionName string) bool {
	_, ok := instance.Exports[functionName]
	return ok
}

// GetFunctionNames mocked method
func (instance *InstanceMock) GetFunctionNames() []string {
	var functionNames []string
	for functionName := range instance.Exports {
		functionNames = append(functionNames, functionName)
	}
	return functionNames
}

// ValidateVoidFunction mocked method
func (instance *InstanceMock) ValidateVoidFunction(functionName string) error {
	return nil
}

// GetInstanceCtxMemory mocked method
func (instance *InstanceMock) GetInstanceCtxMemory() executor.Memory {
	return instance.Memory
}

// GetMemory mocked method
func (instance *InstanceMock) GetMemory() executor.Memory {
	return instance.Memory
}

// IsFunctionImported mocked method
func (instance *InstanceMock) IsFunctionImported(name string) bool {
	_, ok := instance.Exports[name]
	return ok
}

// GetMockInstance gets the mock instance from the runtime of the provided host
func GetMockInstance(host arwen.VMHost) *InstanceMock {
	instance := host.Runtime().GetInstance().(*InstanceMock)
	return instance
}

// SetMemory mocked method
func (instance *InstanceMock) SetMemory(_ []byte) bool {
	return true
}

// IsInterfaceNil mocked method
func (instance *InstanceMock) IsInterfaceNil() bool {
	return instance == nil
}

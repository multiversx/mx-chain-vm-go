package mock

import (
	"errors"
	"testing"

	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen"
	"github.com/ElrondNetwork/wasm-vm-v1_4/wasmer"
)

// InstanceMock is a mock for Wasmer instances; it allows creating mock smart
// contracts within tests, without needing actual WASM smart contracts.
type InstanceMock struct {
	Code            []byte
	Exports         wasmer.ExportsMap
	Points          uint64
	Data            uintptr
	GasLimit        uint64
	BreakpointValue arwen.BreakpointValue
	Memory          wasmer.MemoryHandler
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
		Data:            0,
		GasLimit:        0,
		BreakpointValue: 0,
		Memory:          NewMemoryMock(),
	}
}

func (instance *InstanceMock) Id() string { return "" }

// AddMockMethod adds the provided function as a mocked method to the instance under the specified name.
func (instance *InstanceMock) AddMockMethod(name string, method func() *InstanceMock) {
	instance.AddMockMethodWithError(name, method, nil)
}

// AddMockMethodWithError adds the provided function as a mocked method to the instance under the specified name and returns an error
func (instance *InstanceMock) AddMockMethodWithError(name string, method func() *InstanceMock, err error) {
	wrappedMethod := func(...interface{}) (wasmer.Value, error) {
		instance := method()
		if arwen.BreakpointValue(instance.GetBreakpointValue()) != arwen.BreakpointNone {
			var errMsg string
			if arwen.BreakpointValue(instance.GetBreakpointValue()) == arwen.BreakpointAsyncCall {
				errMsg = "breakpoint"
			} else {
				errMsg = instance.Host.Output().GetVMOutput().ReturnMessage
			}
			err = errors.New(errMsg)
		}
		return wasmer.Void(), err
	}

	instance.Exports[name] = wrappedMethod
}

// HasMemory mocked method
func (instance *InstanceMock) HasMemory() bool {
	return true
}

// SetContextData mocked method
func (instance *InstanceMock) SetContextData(data uintptr) {
	instance.Data = data
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
func (instance *InstanceMock) Clean() bool {
	return true
}

// AlreadyCleaned mocked method
func (instance *InstanceMock) AlreadyCleaned() bool {
	return false
}

// Reset mocked method
func (instance *InstanceMock) Reset() bool {
	return true
}

// GetExports mocked method
func (instance *InstanceMock) GetExports() wasmer.ExportsMap {
	return instance.Exports
}

// GetSignature mocked method
func (instance *InstanceMock) GetSignature(functionName string) (*wasmer.ExportedFunctionSignature, bool) {
	_, ok := instance.Exports[functionName]

	if !ok {
		return nil, false
	}

	return &wasmer.ExportedFunctionSignature{
		InputArity:  0,
		OutputArity: 0,
	}, true
}

// GetData mocked method
func (instance *InstanceMock) GetData() uintptr {
	return instance.Data
}

// GetInstanceCtxMemory mocked method
func (instance *InstanceMock) GetInstanceCtxMemory() wasmer.MemoryHandler {
	return instance.Memory
}

// GetMemory mocked method
func (instance *InstanceMock) GetMemory() wasmer.MemoryHandler {
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

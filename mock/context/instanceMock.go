package mock

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
)

type InstanceMock struct {
	Code            []byte
	Exports         wasmer.ExportsMap
	Points          uint64
	Data            int
	GasLimit        uint64
	BreakpointValue uint64
	Memory          wasmer.MemoryHandler
}

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

func (instance *InstanceMock) AddMockMethod(name string, method func()) {
	wrappedMethod := func(...interface{}) (wasmer.Value, error) {
		method()
		return wasmer.Void(), nil
	}

	instance.Exports[name] = wrappedMethod
}

// HasMemory mocked method
func (instance *InstanceMock) HasMemory() bool {
	return true
}

// SetContextData mocked method
func (instance *InstanceMock) SetContextData(data int) {
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
	instance.BreakpointValue = value
}

// GetBreakpointValue mocked method
func (instance *InstanceMock) GetBreakpointValue() uint64 {
	return instance.BreakpointValue
}

// Cache mocked method
func (instance *InstanceMock) Cache() ([]byte, error) {
	return instance.Code, nil
}

// Clean mocked method
func (instance *InstanceMock) Clean() {
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
func (instance *InstanceMock) GetData() *int {
	return &instance.Data
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

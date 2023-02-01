package mock

import (
	"errors"
	"fmt"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/wasmer"
)

// InstanceMock is a mock for Wasmer instances; it allows creating mock smart
// contracts within tests, without needing actual WASM smart contracts.
type InstanceMock struct {
	Code            []byte
	Exports         wasmer.ExportsMap
	Points          uint64
	Data            executor.VMHooks
	GasLimit        uint64
	BreakpointValue vmhost.BreakpointValue
	Memory          executor.Memory
	Host            vmhost.VMHost
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
		breakpoint := vmhost.BreakpointValue(instance.GetBreakpointValue())
		var err error
		if breakpoint != vmhost.BreakpointNone {
			err = errors.New(breakpoint.String())
		}
		return wasmer.Void(), err
	}

	instance.Exports[name] = wrappedMethod
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
	instance.BreakpointValue = vmhost.BreakpointValue(value)
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
func (instance *InstanceMock) ValidateVoidFunction(_ string) error {
	return nil
}

// ValidateFunctionArities mocked method
func (instance *InstanceMock) ValidateFunctionArities() error {
	return nil
}

// HasMemory mocked method
func (instance *InstanceMock) HasMemory() bool {
	return true
}

// MemLoad returns the contents from the given offset of the WASM memory.
func (instance *InstanceMock) MemLoad(memPtr executor.MemPtr, length executor.MemLength) ([]byte, error) {
	return executor.MemLoadFromMemory(instance.Memory, memPtr, length)
}

// MemStore stores the given data in the WASM memory at the given offset.
func (instance *InstanceMock) MemStore(memPtr executor.MemPtr, data []byte) error {
	return executor.MemStoreToMemory(instance.Memory, memPtr, data)
}

// MemLength returns the length of the allocated memory. Only called directly in tests.
func (instance *InstanceMock) MemLength() uint32 {
	return instance.Memory.Length()
}

// MemGrow allocates more pages to the current memory. Only called directly in tests.
func (instance *InstanceMock) MemGrow(pages uint32) error {
	return instance.Memory.Grow(pages)
}

// MemDump yields the entire contents of the memory. Only used in tests.
func (instance *InstanceMock) MemDump() []byte {
	return instance.Memory.Data()
}

// IsFunctionImported mocked method
func (instance *InstanceMock) IsFunctionImported(name string) bool {
	_, ok := instance.Exports[name]
	return ok
}

// GetMockInstance gets the mock instance from the runtime of the provided host
func GetMockInstance(host vmhost.VMHost) *InstanceMock {
	instance := host.Runtime().GetInstance().(*InstanceMock)
	return instance
}

// Id returns an identifier for the instance, unique at runtime
func (instance *InstanceMock) Id() string {
	return fmt.Sprintf("%p", instance)
}

// IsInterfaceNil mocked method
func (instance *InstanceMock) IsInterfaceNil() bool {
	return instance == nil
}

// SetVMHooksPtr mocked method
func (instance *InstanceMock) SetVMHooksPtr(_ uintptr) {
}

// GetVMHooksPtr mocked method
func (instance *InstanceMock) GetVMHooksPtr() uintptr {
	return 0
}

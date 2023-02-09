package mock

import (
	"errors"
	"fmt"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/wasmer"
)

var _ executor.Instance = (*InstanceMock)(nil)

type mockMethod func() *InstanceMock

// InstanceMock is a mock for Wasmer instances; it allows creating mock smart
// contracts within tests, without needing actual WASM smart contracts.
type InstanceMock struct {
	Code            []byte
	Exports         wasmer.ExportsMap
	DefaultErrors   map[string]error
	Methods         map[string]mockMethod
	Points          uint64
	Data            executor.VMHooks
	GasLimit        uint64
	BreakpointValue vmhost.BreakpointValue
	Memory          executor.Memory
	Host            vmhost.VMHost
	T               testing.TB
	Address         []byte
	AlreadyClean    bool
}

// NewInstanceMock creates a new InstanceMock
func NewInstanceMock(code []byte) *InstanceMock {
	return &InstanceMock{
		Code:            code,
		Exports:         make(wasmer.ExportsMap),
		DefaultErrors:   make(map[string]error),
		Methods:         make(map[string]mockMethod),
		Points:          0,
		Data:            nil,
		GasLimit:        0,
		BreakpointValue: 0,
		Memory:          NewMemoryMock(),
		AlreadyClean:    false,
	}
}

// AddMockMethod adds the provided function as a mocked method to the instance under the specified name.
func (instance *InstanceMock) AddMockMethod(name string, method mockMethod) {
	instance.AddMockMethodWithError(name, method, nil)
}

// AddMockMethodWithError adds the provided function as a mocked method to the instance under the specified name and returns an error
func (instance *InstanceMock) AddMockMethodWithError(name string, method mockMethod, err error) {
	instance.Methods[name] = method
	instance.DefaultErrors[name] = err
	instance.Exports[name] = &wasmer.ExportedFunctionCallInfo{}
}

// CallFunction mocked method
func (instance *InstanceMock) CallFunction(funcName string) error {
	err := instance.DefaultErrors[funcName]
	method := instance.Methods[funcName]
	newInstance := method()
	if vmhost.BreakpointValue(instance.GetBreakpointValue()) != vmhost.BreakpointNone {
		var errMsg string
		if vmhost.BreakpointValue(instance.GetBreakpointValue()) == vmhost.BreakpointAsyncCall {
			errMsg = "breakpoint"
		} else {
			errMsg = newInstance.Host.Output().GetVMOutput().ReturnMessage
		}
		err = errors.New(errMsg)
	}
	return err
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
func (instance *InstanceMock) Clean() bool {
	instance.AlreadyClean = true
	return true
}

// IsAlreadyCleaned mocked method
func (instance *InstanceMock) IsAlreadyCleaned() bool {
	return instance.AlreadyClean
}

// Reset mocked method
func (instance *InstanceMock) Reset() bool {
	return true
}

// HasFunction mocked method
func (instance *InstanceMock) HasFunction(name string) bool {
	_, has := instance.Methods[name]
	return has
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

// ID returns an identifier for the instance, unique at runtime
func (instance *InstanceMock) ID() string {
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

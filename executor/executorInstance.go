package executor

// Instance defines the functionality of a Wasmer instance
type Instance interface {
	HasMemory() bool
	GetPointsUsed() uint64
	SetPointsUsed(points uint64)
	SetGasLimit(gasLimit uint64)
	SetBreakpointValue(value uint64)
	GetBreakpointValue() uint64
	Cache() ([]byte, error)
	Clean()
	CallFunction(functionName string) error
	HasFunction(functionName string) bool
	GetFunctionNames() []string
	ValidateVoidFunction(functionName string) error
	GetMemory() Memory
	IsFunctionImported(name string) bool
	IsInterfaceNil() bool
	Reset() bool
	SetVMHooksPtr(vmHooksPtr uintptr)
	GetVMHooksPtr() uintptr
}

// Memory defines the functionality of the memory of a Wasmer instance
type Memory interface {
	Length() uint32
	Data() []byte
	Grow(pages uint32) error
	Destroy()
	IsInterfaceNil() bool
}

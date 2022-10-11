package executor

// InstanceHandler defines the functionality of a Wasmer instance
type InstanceHandler interface {
	HasMemory() bool
	SetCallbacks(callbacks ImportsInterface)
	GetCallbacks() ImportsInterface
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
	GetInstanceCtxMemory() MemoryHandler
	GetMemory() MemoryHandler
	SetMemory(data []byte) bool
	IsFunctionImported(name string) bool
	IsInterfaceNil() bool
	Reset() bool
}

// MemoryHandler defines the functionality of the memory of a Wasmer instance
type MemoryHandler interface {
	Length() uint32
	Data() []byte
	Grow(pages uint32) error
	Destroy()
	IsInterfaceNil() bool
}

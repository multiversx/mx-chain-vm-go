package wasmer

// InstanceHandler defines the functionality of a Wasmer instance
type InstanceHandler interface {
	HasMemory() bool
	SetContextData(data uintptr)
	GetPointsUsed() uint64
	SetPointsUsed(points uint64)
	SetGasLimit(gasLimit uint64)
	SetBreakpointValue(value uint64)
	GetBreakpointValue() uint64
	Cache() ([]byte, error)
	Clean() bool
	AlreadyCleaned() bool
	GetExports() ExportsMap
	GetSignature(functionName string) (*ExportedFunctionSignature, bool)
	GetData() uintptr
	GetInstanceCtxMemory() MemoryHandler
	GetMemory() MemoryHandler
	SetMemory(data []byte) bool
	IsFunctionImported(name string) bool
	IsInterfaceNil() bool
	Reset() bool
	ID() string
}

// MemoryHandler defines the functionality of the memory of a Wasmer instance
type MemoryHandler interface {
	Length() uint32
	Data() []byte
	Grow(pages uint32) error
	Destroy()
	IsInterfaceNil() bool
}

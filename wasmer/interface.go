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
	Clean()
	GetExports() ExportsMap
	GetSignature(functionName string) (*ExportedFunctionSignature, bool)
	GetData() uintptr
	GetInstanceCtxMemory() MemoryHandler
	GetMemory() MemoryHandler
	IsFunctionImported(name string) bool
}

// MemoryHandler defines the functionality of the memory of a Wasmer instance
type MemoryHandler interface {
	Length() uint32
	Data() []byte
	Grow(pages uint32) error
	Destroy()
}

package executor

// Instance defines the functionality of a Wasmer instance
type Instance interface {
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
	HasMemory() bool
	MemLoad(offset int32, length int32) ([]byte, error)
	MemLoadMultiple(offset int32, lengths []int32) ([][]byte, error)
	MemStore(offset int32, data []byte) error
	MemLength() uint32
	MemGrow(pages uint32) error
	MemDump() []byte
	IsFunctionImported(name string) bool
	IsInterfaceNil() bool
	Reset() bool
	SetVMHooksPtr(vmHooksPtr uintptr)
	GetVMHooksPtr() uintptr
	Id() string
}

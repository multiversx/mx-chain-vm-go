package executor

// Instance defines the functionality of a Wasmer instance
type Instance interface {
	GetPointsUsed() uint64
	SetPointsUsed(points uint64)
	SetGasLimit(gasLimit uint64)
	SetBreakpointValue(value uint64)
	GetBreakpointValue() uint64
	Cache() ([]byte, error)
	Clean() bool
	CallFunction(functionName string) error
	HasFunction(functionName string) bool
	GetFunctionNames() []string
	ValidateFunctionArities() error
	HasMemory() bool
	MemLoad(memPtr MemPtr, length MemLength) ([]byte, error)
	MemStore(memPtr MemPtr, data []byte) error
	MemLength() uint32
	MemGrow(pages uint32) error
	MemDump() []byte
	IsFunctionImported(name string) bool
	IsInterfaceNil() bool
	Reset() bool
	SetVMHooksPtr(vmHooksPtr uintptr)
	GetVMHooksPtr() uintptr
	ID() string
	IsAlreadyCleaned() bool
}

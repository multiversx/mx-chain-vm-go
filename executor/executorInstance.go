package executor

// Instance is the view of a running contract from *outside* a host call: calling it, gas
// accounting, module introspection and lifecycle.
//
// What a VM hook needs while a call is running is not here — that is InstanceState, reached
// through the hook's own context. Keeping the two apart is what lets a backend whose instance
// state is only reachable for the duration of a call implement this interface honestly.
type Instance interface {
	GetPointsUsed() uint64
	SetPointsUsed(points uint64)
	SetGasLimit(gasLimit uint64)
	SetBreakpointValue(value uint64)
	GetBreakpointValue() uint64
	Cache() ([]byte, error)
	Clean() bool
	CallFunction(functionName string) InstanceCallResult
	HasFunction(functionName string) bool
	GetFunctionNames() []string
	ValidateFunctionArities() error
	MemLoad(memPtr MemPtr, length MemLength) ([]byte, error)
	MemStore(memPtr MemPtr, data []byte) error
	MemLength() uint32
	MemGrow(pages uint32) error
	IsFunctionImported(name string) bool
	IsInterfaceNil() bool
	Reset() bool
	ID() string
	IsAlreadyCleaned() bool
}

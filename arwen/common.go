package arwen

type StorageStatus int

const (
	StorageUnchanged StorageStatus = iota
	StorageModified  StorageStatus = iota
	StorageAdded     StorageStatus = iota
	StorageDeleted   StorageStatus = iota
)

type BreakpointValue uint64

const (
	BreakpointNone            BreakpointValue = iota
	BreakpointExecutionFailed BreakpointValue = iota
	BreakpointAsyncCall       BreakpointValue = iota
	BreakpointSignalError     BreakpointValue = iota
	BreakpointSignalExit      BreakpointValue = iota
	BreakpointOutOfGas        BreakpointValue = iota
)

// AsyncCallInfo contains the information required to handle the asynchronous call of another SmartContract
type AsyncCallInfo struct {
	Destination []byte
	Data        []byte
	GasLimit    uint64
	ValueBytes  []byte
}

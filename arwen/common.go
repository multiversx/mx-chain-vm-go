package arwen

type StorageStatus int

const (
	StorageUnchanged StorageStatus = 0
	StorageModified  StorageStatus = 1
	StorageAdded     StorageStatus = 3
	StorageDeleted   StorageStatus = 4
)

type BreakpointValue uint64

const (
	BreakpointNone            BreakpointValue = 0
	BreakpointExecutionFailed BreakpointValue = 1
	BreakpointAsyncCall       BreakpointValue = 2
	BreakpointSignalError     BreakpointValue = 3
	BreakpointSignalExit      BreakpointValue = 4
	BreakpointOutOfGas        BreakpointValue = 5
)

// AsyncCallInfo contains the information required to handle the asynchronous call of another SmartContract
type AsyncCallInfo struct {
	Destination []byte
	Data        []byte
	GasLimit    uint64
	ValueBytes  []byte
}

package async

const CallbackDefault = "callBack"
const AsyncDataPrefix = "asyncCalls"

type AsyncCallExecutionMode uint

const (
	// SyncCall indicates that the async call can be executed synchronously, with
	// its corresponding callback
	SyncCall AsyncCallExecutionMode = iota

	// AsyncBuiltinFunc indicates that the async call is a cross-shard call to a
	// built-in function, which is executed half in-shard, half cross-shard
	AsyncBuiltinFunc

	// AsyncUnknown indicates that the async call cannot be executed locally, and
	// must be forwarded to the destination account
	AsyncUnknown
)

// AsyncCallStatus represents the different status an async call can have
type AsyncCallStatus uint8

const (
	AsyncCallPending AsyncCallStatus = iota
	AsyncCallResolved
	AsyncCallRejected
)

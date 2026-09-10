package vmhost

type AsyncStatus int

const (
	AsyncSuccess        AsyncStatus = 0
	AsyncFailure        AsyncStatus = 1
	AsyncPartialFailure AsyncStatus = 2
)

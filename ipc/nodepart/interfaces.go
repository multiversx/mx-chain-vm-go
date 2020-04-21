package nodepart

import "io"

// ParentLogsPart defines the interface for the Node's part of the logging dialogue
type ParentLogsPart interface {
	StartLoop(childStdout io.Reader, childStderr io.Reader) error
	StopLoop()
}

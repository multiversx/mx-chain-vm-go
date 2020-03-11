package logger

import (
	"fmt"
	"os"
)

var _ Logger = (*PipeLogger)(nil)

// PipeLogger is a pipe-based logger
type PipeLogger struct {
	level LogLevel
}

// NewPipeLogger creates a new pipe logger
func NewPipeLogger(level LogLevel, pipeWrite *os.File) Logger {
	return &PipeLogger{
		level: level,
	}
}

func (pipeLogger *PipeLogger) Trace(message string, args ...interface{}) {
	if pipeLogger.level > LogTrace {
		return
	}

	//pipeLogger.messenger.Send()
}

func (pipeLogger *PipeLogger) Debug(message string, args ...interface{}) {
	if pipeLogger.level > LogDebug {
		return
	}

	fmt.Printf("DEBUG:"+message+"\n", args...)
}

func (pipeLogger *PipeLogger) Info(message string, args ...interface{}) {
	if pipeLogger.level > LogInfo {
		return
	}

	fmt.Printf("INFO:"+message+"\n", args...)
}

func (pipeLogger *PipeLogger) Warn(message string, args ...interface{}) {
	if pipeLogger.level > LogWarning {
		return
	}

	fmt.Printf("WARN:"+message+"\n", args...)
}

func (pipeLogger *PipeLogger) Error(message string, args ...interface{}) {
	if pipeLogger.level > LogError {
		return
	}

	fmt.Printf("ERROR:"+message+"\n", args...)
}

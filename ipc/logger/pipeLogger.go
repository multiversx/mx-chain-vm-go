package logger

import (
	"encoding/binary"
	"encoding/json"
	"os"
)

var _ Logger = (*PipeLogger)(nil)

// LogMessage is a log message
type LogMessage struct {
	Level   LogLevel
	Message string
	Args    interface{}
}

// PipeLogger is a pipe-based logger
type PipeLogger struct {
	pipe     *os.File
	level    LogLevel
	fallback Logger
}

// NewPipeLogger creates a new pipe logger
func NewPipeLogger(level LogLevel, pipe *os.File) Logger {
	return &PipeLogger{
		level:    level,
		pipe:     pipe,
		fallback: NewDefaultLogger(level),
	}
}

// Trace logs
func (pipeLogger *PipeLogger) Trace(message string, args ...interface{}) {
	if pipeLogger.level > LogTrace {
		return
	}

	pipeLogger.sendMessage(&LogMessage{
		Level:   LogTrace,
		Message: message,
		Args:    args,
	})
}

// Debug logs
func (pipeLogger *PipeLogger) Debug(message string, args ...interface{}) {
	if pipeLogger.level > LogDebug {
		return
	}

	pipeLogger.sendMessage(&LogMessage{
		Level:   LogDebug,
		Message: message,
		Args:    args,
	})
}

// Info logs
func (pipeLogger *PipeLogger) Info(message string, args ...interface{}) {
	if pipeLogger.level > LogInfo {
		return
	}

	pipeLogger.sendMessage(&LogMessage{
		Level:   LogInfo,
		Message: message,
		Args:    args,
	})
}

// Warn logs
func (pipeLogger *PipeLogger) Warn(message string, args ...interface{}) {
	if pipeLogger.level > LogWarning {
		return
	}

	pipeLogger.sendMessage(&LogMessage{
		Level:   LogWarning,
		Message: message,
		Args:    args,
	})

}

// Error logs
func (pipeLogger *PipeLogger) Error(message string, args ...interface{}) {
	if pipeLogger.level > LogError {
		return
	}

	pipeLogger.sendMessage(&LogMessage{
		Level:   LogError,
		Message: message,
		Args:    args,
	})
}

func (pipeLogger *PipeLogger) sendMessage(message *LogMessage) {
	payload, err := pipeLogger.marshal(message)
	if err != nil {
		pipeLogger.fallback.Error("sendMessage marshal error", err.Error())
		return
	}

	// Send length
	length := len(payload)
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, uint32(length))
	_, err = pipeLogger.pipe.Write(buffer)
	if err != nil {
		pipeLogger.fallback.Error("sendMessage send length error", err.Error())
		return
	}

	// Send payload
	_, err = pipeLogger.pipe.Write(payload)
	if err != nil {
		pipeLogger.fallback.Error("sendMessage send payload error", err.Error())
		return
	}
}

func (pipeLogger *PipeLogger) marshal(data interface{}) ([]byte, error) {
	return marshalJSON(data)
}

func (pipeLogger *PipeLogger) unmarshal(dataBytes []byte, data interface{}) error {
	return unmarshalJSON(dataBytes, data)
}

func marshalJSON(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func unmarshalJSON(dataBytes []byte, data interface{}) error {
	return json.Unmarshal(dataBytes, data)
}

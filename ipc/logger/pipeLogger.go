package logger

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/marshaling"
)

var _ Logger = (*PipeLogger)(nil)

// LogMessage is a log message
type LogMessage struct {
	Level   LogLevel
	Message string
	Args    []interface{}
}

// PipeLogger is a pipe-based logger
type PipeLogger struct {
	pipe        *os.File
	level       LogLevel
	fallback    Logger
	marshalizer marshaling.Marshalizer
}

// NewPipeLogger creates a new pipe logger
func NewPipeLogger(level LogLevel, pipe *os.File, marshalizer marshaling.Marshalizer) Logger {
	return &PipeLogger{
		level:       level,
		pipe:        pipe,
		fallback:    NewDefaultLogger(level),
		marshalizer: marshalizer,
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
	payload, err := pipeLogger.marshalizer.MarshalItem(message)
	if err != nil {
		pipeLogger.fallback.Error("pipeLogger.sendMessage() marshal error", err.Error())
		return
	}

	length := len(payload)
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, uint32(length))
	_, err = pipeLogger.pipe.Write(buffer)
	if err != nil {
		pipeLogger.fallback.Error("pipeLogger.sendMessage() send length error", err.Error())
		return
	}

	_, err = pipeLogger.pipe.Write(payload)
	if err != nil {
		pipeLogger.fallback.Error("pipeLogger.sendMessage() send payload error", err.Error())
		return
	}
}

// ReceiveLogThroughPipe reads a log message from the pipe and sends it to a regular Node logger
// TODO: refactor, create LogsReceiver component
func ReceiveLogThroughPipe(receivingLogger Logger, pipe *os.File, marshalizer marshaling.Marshalizer) error {
	buffer := make([]byte, 4)
	_, err := io.ReadFull(pipe, buffer)
	if err != nil {
		return err
	}

	length := binary.LittleEndian.Uint32(buffer)
	buffer = make([]byte, length)
	_, err = io.ReadFull(pipe, buffer)
	if err != nil {
		return err
	}

	logMessage := &LogMessage{}
	err = marshalizer.UnmarshalItem(buffer, logMessage)
	if err != nil {
		return err
	}

	switch logMessage.Level {
	case LogTrace:
		receivingLogger.Trace(logMessage.Message, logMessage.Args...)
	case LogDebug:
		receivingLogger.Debug(logMessage.Message, logMessage.Args...)
	case LogInfo:
		receivingLogger.Info(logMessage.Message, logMessage.Args...)
	case LogWarning:
		receivingLogger.Warn(logMessage.Message, logMessage.Args...)
	case LogError:
		receivingLogger.Error(logMessage.Message, logMessage.Args...)
	default:
		receivingLogger.Error("Unknown log message from Arwen")
	}

	return nil
}

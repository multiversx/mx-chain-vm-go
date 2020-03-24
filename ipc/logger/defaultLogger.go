package logger

import "fmt"

var _ Logger = (*defaultLogger)(nil)

type defaultLogger struct {
	level LogLevel
}

// NewDefaultLogger creates a logger
func NewDefaultLogger(level LogLevel) Logger {
	return &defaultLogger{
		level: level,
	}
}

// Trace logs
func (logger *defaultLogger) Trace(message string, args ...interface{}) {
	if logger.level > LogTrace {
		return
	}

	fmt.Println("TRACE:", message, args)
}

// Debug logs
func (logger *defaultLogger) Debug(message string, args ...interface{}) {
	if logger.level > LogDebug {
		return
	}

	fmt.Println("DEBUG:", message, args)
}

// Info logs
func (logger *defaultLogger) Info(message string, args ...interface{}) {
	if logger.level > LogInfo {
		return
	}

	fmt.Println("INFO:", message, args)
}

// Warn logs
func (logger *defaultLogger) Warn(message string, args ...interface{}) {
	if logger.level > LogWarning {
		return
	}

	fmt.Println("WARN:", message, args)
}

// Error logs
func (logger *defaultLogger) Error(message string, args ...interface{}) {
	if logger.level > LogError {
		return
	}

	fmt.Println("ERROR:", message, args)
}

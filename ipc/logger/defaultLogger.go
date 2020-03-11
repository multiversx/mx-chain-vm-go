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

func (logger *defaultLogger) Trace(message string, args ...interface{}) {
	if logger.level > LogTrace {
		return
	}

	fmt.Printf("TRACE:"+message+"\n", args...)
}

func (logger *defaultLogger) Debug(message string, args ...interface{}) {
	if logger.level > LogDebug {
		return
	}

	fmt.Printf("DEBUG:"+message+"\n", args...)
}

func (logger *defaultLogger) Info(message string, args ...interface{}) {
	if logger.level > LogInfo {
		return
	}

	fmt.Printf("INFO:"+message+"\n", args...)
}

func (logger *defaultLogger) Warn(message string, args ...interface{}) {
	if logger.level > LogWarning {
		return
	}

	fmt.Printf("WARN:"+message+"\n", args...)
}

func (logger *defaultLogger) Error(message string, args ...interface{}) {
	if logger.level > LogError {
		return
	}

	fmt.Printf("ERROR:"+message+"\n", args...)
}

package logger

// Logger is the logger interface
type Logger interface {
	Trace(message string, args ...interface{})
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
}

// LogLevel defines the priority level of a log line. Trace is the lowest priority level, Error is the highest
type LogLevel uint8

// These constants are the string representation of the package logging levels.
const (
	LogTrace   LogLevel = 0
	LogDebug   LogLevel = 1
	LogInfo    LogLevel = 2
	LogWarning LogLevel = 3
	LogError   LogLevel = 4
	LogNone    LogLevel = 5
)

package executorwrappers

import "strings"

type ExecutorLogger interface {
	LogExecutorEvent(description string)
}

// StringLogger is a simple ExecutorLogger that records data into a string builder.
type StringLogger struct {
	sb strings.Builder
}

func NewStringLogger() *StringLogger {
	return &StringLogger{}
}

func (sl *StringLogger) LogExecutorEvent(description string) {
	sl.sb.WriteString(description)
	sl.sb.WriteRune('\n')
}

func (sl *StringLogger) String() string {
	return sl.sb.String()
}

// NoLogger is an ExecutorLogger implementation that does nothing.
type NoLogger struct{}

func (*NoLogger) LogExecutorEvent(description string) {}

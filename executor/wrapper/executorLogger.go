package executorwrapper

import (
	"fmt"
	"strings"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
)

var log = logger.GetOrCreate("vm/executor")

// ExecutorLogger defines a logging interface for the WrapperExecutor.
type ExecutorLogger interface {
	LogExecutorEvent(description string)
	LogVMHookCallBefore(callInfo string)
	LogVMHookCallAfter(callInfo string)
}

// ConsoleLogger is a simple ExecutorLogger that records data into the console.
type ConsoleLogger struct {
}

// NewConsoleLogger creates a new ConsoleLogger, which records events into the console.
func NewConsoleLogger() *ConsoleLogger {
	cl := &ConsoleLogger{}
	log.Trace("Starting Console Logger:")
	return cl
}

// LogExecutorEvent logs a custom event from the executor.
func (cl *ConsoleLogger) LogExecutorEvent(description string) {
	log.Trace(description)
}

// LogVMHookCallBefore is called before processing a wrapped VM hook.
func (cl *ConsoleLogger) LogVMHookCallBefore(callInfo string) {
	log.Trace(fmt.Sprintf("VM hook begin: %s", callInfo))
}

// LogVMHookCallAfter is called after processing a wrapped VM hook.
func (cl *ConsoleLogger) LogVMHookCallAfter(callInfo string) {
	log.Trace(fmt.Sprintf("VM hook end: %s", callInfo))
}

// StringLogger is a simple ExecutorLogger that records data into a string builder.
type StringLogger struct {
	sb strings.Builder
}

// NewStringLogger creates a new StringLogger, which records events into a string builder.
func NewStringLogger() *StringLogger {
	sl := &StringLogger{}
	sl.sb.WriteString("starting log:\n")
	return sl
}

// LogExecutorEvent logs a custom event from the executor.
func (sl *StringLogger) LogExecutorEvent(description string) {
	sl.sb.WriteString(description)
	sl.sb.WriteRune('\n')
}

// LogVMHookCallBefore is called before processing a wrapped VM hook.
func (sl *StringLogger) LogVMHookCallBefore(callInfo string) {
	sl.sb.WriteString("VM hook begin: ")
	sl.sb.WriteString(callInfo)
	sl.sb.WriteRune('\n')
}

// LogVMHookCallAfter is called after processing a wrapped VM hook.
func (sl *StringLogger) LogVMHookCallAfter(callInfo string) {
	sl.sb.WriteString("VM hook end:   ")
	sl.sb.WriteString(callInfo)
	sl.sb.WriteRune('\n')
}

// String yields the logs accumulated up to this point.
func (sl *StringLogger) String() string {
	return sl.sb.String()
}

// NoLogger is an ExecutorLogger implementation that does nothing.
type NoLogger struct{}

// SetCurrentInstance does nothing.
func (*NoLogger) SetCurrentInstance(_ executor.Instance) {}

// LogExecutorEvent does nothing.
func (*NoLogger) LogExecutorEvent(_ string) {}

// LogVMHookCallBefore does nothing.
func (*NoLogger) LogVMHookCallBefore(_ string) {}

// LogVMHookCallAfter does nothing.
func (*NoLogger) LogVMHookCallAfter(_ string) {}

package executorwrapper

import (
	"fmt"
	"strings"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-vm-go/executor"
)

var consoleLogger = logger.GetOrCreate("vm/consoleLogger")

// ExecutorLogger defines a logging interface for the WrapperExecutor.
type ExecutorLogger interface {
	SetCurrentInstance(instance executor.Instance)
	LogExecutorEvent(description string)
	LogVMHookCallBefore(callInfo string)
	LogVMHookCallAfter(callInfo string)
}

// ConsoleLogger is a simple ExecutorLogger that records data into the console.
type ConsoleLogger struct {
	currentInstance executor.Instance
}

// NewConsoleLogger creates a new ConsoleLogger, which records events into the console.
func NewConsoleLogger() *ConsoleLogger {
	cl := &ConsoleLogger{}
	consoleLogger.Trace("Starting Console Logger:")
	return cl
}

// SetCurrentInstance adds context pertaining to the current instance, when running tests.
func (cl *ConsoleLogger) SetCurrentInstance(instance executor.Instance) {
	cl.currentInstance = instance
}

// LogExecutorEvent logs a custom event from the executor.
func (cl *ConsoleLogger) LogExecutorEvent(description string) {
	consoleLogger.Trace(description)
}

// LogVMHookCallBefore is called before processing a wrapped VM hook.
func (cl *ConsoleLogger) LogVMHookCallBefore(callInfo string) {
	consoleLogger.Trace(fmt.Sprintf("VM hook begin: %s", callInfo))
	if !cl.currentInstance.IsInterfaceNil() {
		consoleLogger.Trace((fmt.Sprintf("Points used: %d", cl.currentInstance.GetPointsUsed())))
	}
}

// LogVMHookCallAfter is called after processing a wrapped VM hook.
func (cl *ConsoleLogger) LogVMHookCallAfter(callInfo string) {
	consoleLogger.Trace(fmt.Sprintf("VM hook end: %s", callInfo))
	if !cl.currentInstance.IsInterfaceNil() {
		consoleLogger.Trace(fmt.Sprintf("Points used: %d", cl.currentInstance.GetPointsUsed()))
	}
}

// StringLogger is a simple ExecutorLogger that records data into a string builder.
type StringLogger struct {
	sb              strings.Builder
	currentInstance executor.Instance
}

// NewStringLogger creates a new StringLogger, which records events into a string builder.
func NewStringLogger() *StringLogger {
	sl := &StringLogger{}
	sl.sb.WriteString("starting log:\n")
	return sl
}

// SetCurrentInstance adds context pertaining to the current instance, when running tests.
func (sl *StringLogger) SetCurrentInstance(instance executor.Instance) {
	sl.currentInstance = instance
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
	if !sl.currentInstance.IsInterfaceNil() {
		sl.sb.WriteString(fmt.Sprintf(" points used: %d", sl.currentInstance.GetPointsUsed()))
	}
	sl.sb.WriteRune('\n')
}

// LogVMHookCallAfter is called after processing a wrapped VM hook.
func (sl *StringLogger) LogVMHookCallAfter(callInfo string) {
	sl.sb.WriteString("VM hook end:   ")
	sl.sb.WriteString(callInfo)
	if !sl.currentInstance.IsInterfaceNil() {
		sl.sb.WriteString(fmt.Sprintf(" points used: %d", sl.currentInstance.GetPointsUsed()))
	}
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

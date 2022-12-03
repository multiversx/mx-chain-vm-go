package executorwrapper

import (
	"fmt"
	"strings"

	"github.com/ElrondNetwork/wasm-vm/executor"
)

type ExecutorLogger interface {
	SetCurrentInstance(instance executor.Instance)
	LogExecutorEvent(description string)
	LogVMHookCallBefore(callInfo string)
	LogVMHookCallAfter(callInfo string)
}

// StringLogger is a simple ExecutorLogger that records data into a string builder.
type StringLogger struct {
	sb              strings.Builder
	currentInstance executor.Instance
}

func NewStringLogger() *StringLogger {
	sl := &StringLogger{}
	sl.sb.WriteString("starting log:\n")
	return sl
}

func (sl *StringLogger) SetCurrentInstance(instance executor.Instance) {
	sl.currentInstance = instance
}

func (sl *StringLogger) LogExecutorEvent(description string) {
	sl.sb.WriteString(description)
	sl.sb.WriteRune('\n')
}

func (sl *StringLogger) LogVMHookCallBefore(callInfo string) {
	sl.sb.WriteString("VM hook begin: ")
	sl.sb.WriteString(callInfo)
	if !sl.currentInstance.IsInterfaceNil() {
		sl.sb.WriteString(fmt.Sprintf(" points used: %d", sl.currentInstance.GetPointsUsed()))
	}
	sl.sb.WriteRune('\n')
}

func (sl *StringLogger) LogVMHookCallAfter(callInfo string) {
	sl.sb.WriteString("VM hook end:   ")
	sl.sb.WriteString(callInfo)
	if !sl.currentInstance.IsInterfaceNil() {
		sl.sb.WriteString(fmt.Sprintf(" points used: %d", sl.currentInstance.GetPointsUsed()))
	}
	sl.sb.WriteRune('\n')
}

func (sl *StringLogger) String() string {
	return sl.sb.String()
}

// NoLogger is an ExecutorLogger implementation that does nothing.
type NoLogger struct{}

func (*NoLogger) SetCurrentInstance(instance executor.Instance) {}
func (*NoLogger) LogExecutorEvent(description string)           {}
func (*NoLogger) LogVMHookCallBefore(callInfo string)           {}
func (*NoLogger) LogVMHookCallAfter(callInfo string)            {}

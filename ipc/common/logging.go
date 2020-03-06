package common

import (
	"fmt"
)

// LogDebug logs
func LogDebug(format string, values ...interface{}) {
	//fmt.Printf(format+"\n", values...)
}

// LogInfo logs
func LogInfo(format string, values ...interface{}) {
	fmt.Printf(format+"\n", values...)
}

// LogDebugJSON logs
func LogDebugJSON(message string, value interface{}) {
	// jsonValue, _ := json.MarshalIndent(value, "", "\t")
	// fmt.Println(message)
	// fmt.Println(string(jsonValue))
}

// LogError logs
func LogError(format string, values ...interface{}) {
	fmt.Printf("ERROR: "+format+"\n", values...)
}

// NodeLogger is
type NodeLogger interface {
	Trace(message string, args ...interface{})
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
}

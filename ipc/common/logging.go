package common

import (
	"fmt"
)

// LogDebug logs
func LogDebug(format string, values ...interface{}) {
	//fmt.Printf(format+"\n", values...)
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

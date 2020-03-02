package common

import "fmt"

// LogDebug logs
func LogDebug(format string, values ...interface{}) {
	//fmt.Printf(format+"\n", values...)
}

// LogError logs
func LogError(format string, values ...interface{}) {
	fmt.Printf("ERROR: "+format+"\n", values...)
}

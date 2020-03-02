package common

import (
	"fmt"
)

// LogDebug logs
func LogDebug(format string, values ...interface{}) {
	fmt.Printf(format+"\n", values...)
}

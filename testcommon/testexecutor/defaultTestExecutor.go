// Package testexecutor provides with a default executor for testing
package testexecutor

import (
	"os"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
)

// EnvVMEXECUTOR is the name of the environment variable that controls the default test executor
var EnvVMEXECUTOR = "VMEXECUTOR"

// ExecWasmer1 is the value of the EnvVMEXECUTOR variable which selects Wasmer 1
var ExecWasmer1 = "wasmer1"

// ExecWasmer2 is the value of the EnvVMEXECUTOR variable which selects Wasmer 2
var ExecWasmer2 = "wasmer2"

var defaultExecutorString = ExecWasmer2

// NewDefaultTestExecutorFactory instantiates an executor factory based on the $VMEXECUTOR environment variable
func NewDefaultTestExecutorFactory(_ testing.TB) executor.ExecutorAbstractFactory {
	return wasmer2.ExecutorFactory()
}

// IsWasmer1Allowed returns true if the default test executor is Wasmer 1.
// If the default test executor is Wasmer 2, it is not allowed to instantiate a
// Wasmer 1 executor due to low-level conflicts between Wasmer 1 and 2.
func IsWasmer1Allowed() bool {
	return getVMExecutorString() == ExecWasmer1
}

func getVMExecutorString() string {
	execStr := os.Getenv(EnvVMEXECUTOR)

	if len(execStr) == 0 {
		execStr = defaultExecutorString
	}

	return execStr
}

// Package testexecutor provides with a default executor for testing
package testexecutor

import (
	"fmt"
	"os"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
	"github.com/multiversx/mx-chain-vm-go/wasmer3"
)

// EnvVMEXECUTOR is the name of the environment variable that controls the default test executor
var EnvVMEXECUTOR = "VMEXECUTOR"

// ExecWasmer2 is the value of the EnvVMEXECUTOR variable which selects Wasmer 2
var ExecWasmer2 = "wasmer2"

// ExecWasmer3 is the value of the EnvVMEXECUTOR variable which selects Wasmer 6, through the
// second executor interface set. Not the default: see PLAN-go-ffi.md, Phase 4b, for what the
// Go side still has to converge before it can be.
var ExecWasmer3 = "wasmer3"

var defaultExecutorString = ExecWasmer2

// NewDefaultTestExecutorFactory instantiates an executor factory based on the $VMEXECUTOR environment variable
func NewDefaultTestExecutorFactory(tb testing.TB) executor.ExecutorAbstractFactory {
	execStr := getVMExecutorString()

	if execStr == ExecWasmer2 {
		return wasmer2.ExecutorFactory()
	}

	if execStr == ExecWasmer3 {
		return wasmer3.ExecutorFactory()
	}

	if tb == (testing.TB)(nil) {
		panic(fmt.Sprintf("executor %s not recognized", execStr))
	}
	tb.Fatalf("executor %s not recognized", execStr)

	return nil
}

func getVMExecutorString() string {
	execStr := os.Getenv(EnvVMEXECUTOR)

	if len(execStr) == 0 {
		execStr = defaultExecutorString
	}

	return execStr
}

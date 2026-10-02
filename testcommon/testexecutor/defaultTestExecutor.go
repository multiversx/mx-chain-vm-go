// Package testexecutor provides with a default executor for testing
package testexecutor

import (
	"fmt"
	"os"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/wasmer2legacy"
	"github.com/multiversx/mx-chain-vm-go/wasmer6"
)

// EnvVMEXECUTOR is the name of the environment variable that controls the default test executor
var EnvVMEXECUTOR = "VMEXECUTOR"

// ExecWasmer2Legacy is the value of the EnvVMEXECUTOR variable which selects Wasmer 2.2
// through the legacy C API and the legacy executor interfaces. The default, and what
// production runs today.
var ExecWasmer2Legacy = "wasmer2legacy"

// ExecWasmer2 is reserved for Wasmer 2.2 through the *new* C API and interfaces, which does
// not exist yet. Named now so that the two Wasmer 2.2 integrations are never confused: they
// will be kept side by side until the migration is shown to be correct, which is what the
// hook-sequence differential in integrationTests/json is for.
var ExecWasmer2 = "wasmer2"

// ExecWasmer6 is the value of the EnvVMEXECUTOR variable which selects Wasmer 6, through the
// second executor interface set. Not the default: see PLAN-go-ffi.md, Phase 4b, for what the
// Go side still has to converge before it can be.
var ExecWasmer6 = "wasmer6"

var defaultExecutorString = ExecWasmer2Legacy

// NewDefaultTestExecutorFactory instantiates an executor factory based on the $VMEXECUTOR environment variable
func NewDefaultTestExecutorFactory(tb testing.TB) executor.ExecutorAbstractFactory {
	execStr := getVMExecutorString()

	if execStr == ExecWasmer2Legacy {
		return wasmer2legacy.ExecutorFactory()
	}

	if execStr == ExecWasmer6 {
		return wasmer6.ExecutorFactory()
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

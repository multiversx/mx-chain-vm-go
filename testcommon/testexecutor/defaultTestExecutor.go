// Package testexecutor provides with a default executor for testing
package testexecutor

import (
	"fmt"
	"os"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
	"github.com/multiversx/mx-chain-vm-go/wasmer2legacy"
	"github.com/multiversx/mx-chain-vm-go/wasmer7"
)

// EnvVMEXECUTOR is the name of the environment variable that controls the default test executor
var EnvVMEXECUTOR = "VMEXECUTOR"

// ExecWasmer2Legacy is the value of the EnvVMEXECUTOR variable which selects Wasmer 2.2
// through the legacy C API and the legacy executor interfaces. The default, and what
// production runs today.
var ExecWasmer2Legacy = "wasmer2legacy"

// ExecWasmer2 is the value of the EnvVMEXECUTOR variable which selects Wasmer 2.2 through the
// *new* C API and interfaces.
//
// The same runtime as ExecWasmer2Legacy, reached through the new boundary, which is what
// makes the pair worth running against each other: a difference between them is attributable
// to the migration, with no runtime change to confound it. Comparing against ExecWasmer7
// cannot say that, because it moves both at once. The two are kept side by side until the
// migration is shown to be correct — see the hook-sequence differential in
// integrationTests/json.
var ExecWasmer2 = "wasmer2"

// ExecWasmer7 is the value of the EnvVMEXECUTOR variable which selects Wasmer 7, through the
// second executor interface set. Not the default: see PLAN-go-ffi.md, Phase 4b, for what the
// Go side still has to converge before it can be.
var ExecWasmer7 = "wasmer7"

var defaultExecutorString = ExecWasmer2Legacy

// NewDefaultTestExecutorFactory instantiates an executor factory based on the $VMEXECUTOR environment variable
func NewDefaultTestExecutorFactory(tb testing.TB) executor.ExecutorAbstractFactory {
	execStr := getVMExecutorString()

	if execStr == ExecWasmer2Legacy {
		return wasmer2legacy.ExecutorFactory()
	}

	if execStr == ExecWasmer2 {
		return wasmer2.ExecutorFactory()
	}

	if execStr == ExecWasmer7 {
		return wasmer7.ExecutorFactory()
	}

	// Named rather than just rejected, because the value that gets here is almost always a
	// stale one: the names have been renamed once already (wasmer2 -> wasmer2legacy), and a
	// script or Makefile target carrying the old spelling is the likely cause.
	message := fmt.Sprintf("executor %q not recognized; %s must be one of %s, %s or %s",
		execStr, EnvVMEXECUTOR, ExecWasmer2Legacy, ExecWasmer2, ExecWasmer7)

	if tb == (testing.TB)(nil) {
		panic(message)
	}
	tb.Fatal(message)

	return nil
}

func getVMExecutorString() string {
	execStr := os.Getenv(EnvVMEXECUTOR)

	if len(execStr) == 0 {
		execStr = defaultExecutorString
	}

	return execStr
}

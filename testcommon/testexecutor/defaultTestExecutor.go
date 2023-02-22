// Package testexecutor provides with a default executor for testing
package testexecutor

import (
	"fmt"
	"os"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/wasmer"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
)

var defaultExecutorString = "wasmer1"

// NewDefaultTestExecutorFactory instantiates an executor factory based on a CLI flag specified to `go test`
func NewDefaultTestExecutorFactory(tb testing.TB) executor.ExecutorAbstractFactory {
	execStr := getVMExecutorString()

	if execStr == "wasmer1" {
		return wasmer.ExecutorFactory()
	}
	if execStr == "wasmer2" {
		return wasmer2.ExecutorFactory()
	}

	if tb == (testing.TB)(nil) {
		panic(fmt.Sprintf("executor %s not recognized", execStr))
	}
	tb.Fatalf("executor %s not recognized", execStr)

	return nil
}

func getVMExecutorString() string {
	execStr := os.Getenv("VMEXECUTOR")

	if len(execStr) == 0 {
		execStr = defaultExecutorString
	}

	return execStr
}

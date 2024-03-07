package testcommon

import (
	"runtime"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/testcommon/testexecutor"
)

func SkipTestOnDarwin(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("skipping test on darwin")
	}
}

func SkipTestOnARM64(t *testing.T) {
	if runtime.GOARCH == "arm64" {
		t.Skip("skipping test on arm64")
	}
}

func SkipTestIfWasmer1NotAllowed(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}
}

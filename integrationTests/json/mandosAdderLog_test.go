package vmjsonintegrationtest

import (
	"testing"

	"github.com/ElrondNetwork/wasm-vm/wasmer2"
)

func TestRustAdderLog(t *testing.T) {
	expected := MandosTest(t).
		Folder("adder/mandos").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	MandosTest(t).
		Folder("adder/mandos").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

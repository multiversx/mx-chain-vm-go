package vmjsonintegrationtest

import (
	"github.com/multiversx/mx-chain-vm-go/evm"
	"testing"
)

func TestAdderEVM(t *testing.T) {
	ScenariosTest(t).
		Folder("evm/adder/scenarios").
		WithExecutorFactory(evm.ExecutorFactory()).
		Run().
		CheckNoError()
}

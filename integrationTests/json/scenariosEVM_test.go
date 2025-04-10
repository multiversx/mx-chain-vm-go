package vmjsonintegrationtest

import (
	"github.com/multiversx/mx-chain-vm-go/evm"
	"github.com/multiversx/mx-chain-vm-go/scenario"
	"testing"
)

func TestAdderEVM(t *testing.T) {
	ScenariosTest(t).
		Folder("evm/adder/scenarios").
		WithExecutorFactory(evm.ExecutorFactory()).
		WithVMType(scenario.EVMType).
		WithOmitFunctionNameChecks(true).
		WithOmitDefaultCodeChanges(true).
		Run().
		CheckNoError()
}

package vmjsonintegrationtest

import (
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
)

func TestCErc20Executors_TwiceW1ThenTwiceW2(t *testing.T) {
	testCERC20WithExecutorFactory(t, wasmer2.ExecutorFactoryNew())
	testCERC20WithExecutorFactory(t, wasmer2.ExecutorFactoryNew())
}

func testCERC20WithExecutorFactory(t *testing.T, factory executor.ExecutorNewAbstractFactory) {
	ScenariosTest(t).
		Folder("erc20-c").
		WithExecutorFactory(factory).
		WithExecutorLogs().
		Run().
		CheckNoError()
}

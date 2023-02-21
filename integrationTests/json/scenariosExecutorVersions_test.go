package vmjsonintegrationtest

import (
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/wasmer"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
)

func TestCErc20Executors_Works(t *testing.T) {
	testCERC20WithExecutorFactory(t, wasmer.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer2.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer2.ExecutorFactory())
}

func TestCErc20Executors_AlsoWorks(t *testing.T) {
	testCERC20WithExecutorFactory(t, wasmer.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer2.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer2.ExecutorFactory())
}

func TestCErc20Executors_Works3(t *testing.T) {
	testCERC20WithExecutorFactory(t, wasmer.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer2.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer2.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer2.ExecutorFactory())
}

func TestCErc20Executors_Fails1(t *testing.T) {
	testCERC20WithExecutorFactory(t, wasmer2.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer2.ExecutorFactory())
}

func TestCErc20Executors_Fails2(t *testing.T) {
	testCERC20WithExecutorFactory(t, wasmer2.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer2.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer2.ExecutorFactory())
}

func testCERC20WithExecutorFactory(t *testing.T, factory executor.ExecutorAbstractFactory) {
	ScenariosTest(t).
		Folder("erc20-c").
		WithExecutorFactory(factory).
		WithExecutorLogs().
		Run().
		CheckNoError()
}

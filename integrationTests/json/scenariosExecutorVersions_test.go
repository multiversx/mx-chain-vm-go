package vmjsonintegrationtest

import (
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/wasmer2legacy"
	"github.com/multiversx/mx-chain-vm-go/wasmer7"
)

// TestCErc20Executors runs the same scenario suite under each executor in turn, in one process.
//
// This test dates from the Wasmer 1 to Wasmer 2 migration, where it ran each version twice to
// catch state leaking between executors of the same kind and between kinds. It had decayed to
// running Wasmer 2 twice once Wasmer 1 was dropped; here it is pointed at the Wasmer 2 to
// Wasmer 7 migration, which is what it was built for.
//
// Running them in one process is the point: both libraries are loaded and both runtimes'
// trap handlers are installed, so this also exercises the coexistence that
// wasmer7/dualbackend/ tests directly.
func TestCErc20Executors(t *testing.T) {
	testCERC20WithExecutorFactory(t, wasmer2legacy.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer2legacy.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer7.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer7.ExecutorFactory())
	testCERC20WithExecutorFactory(t, wasmer2legacy.ExecutorFactory())
}

func testCERC20WithExecutorFactory(t *testing.T, factory executor.ExecutorAbstractFactory) {
	ScenariosTest(t).
		Folder("erc20-c").
		WithExecutorFactory(factory).
		WithExecutorLogs().
		Run().
		CheckNoError()
}

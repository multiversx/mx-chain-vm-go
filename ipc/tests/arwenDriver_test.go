package tests

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/nodepart"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	"github.com/stretchr/testify/require"
)

var arwenVirtualMachine = []byte{5, 0}

func TestArwenDriver_DiagnoseWait(t *testing.T) {
	blockchain := &mock.BlockChainHookStub{}
	driver, err := nodepart.NewArwenDriver(blockchain, arwenVirtualMachine, uint64(10000000), config.MakeGasMap(1))
	require.Nil(t, err)
	require.NotNil(t, driver)
	err = driver.DiagnoseWait(100)
	require.Nil(t, err)
}

func TestArwenDriver_DiagnoseWaitWithTimeout(t *testing.T) {
	blockchain := &mock.BlockChainHookStub{}
	driver, err := nodepart.NewArwenDriver(blockchain, arwenVirtualMachine, uint64(10000000), config.MakeGasMap(1))
	require.Nil(t, err)
	require.NotNil(t, driver)
	err = driver.DiagnoseWait(5000)
	require.True(t, common.IsCriticalError(err))
	require.Contains(t, err.Error(), "timeout")
}

// Test, restarts upon critical error

package tests

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/nodepart"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

var arwenVirtualMachine = []byte{5, 0}

func TestArwenDriver_StopsArwenOnTimeout(t *testing.T) {
	blockchain := &mock.BlockChainHookStub{}
	driver, err := nodepart.NewArwenDriver(blockchain, arwenVirtualMachine, uint64(10000000), config.MakeGasMap(1))
	require.Nil(t, err)
	require.NotNil(t, driver)
	vmOutput, err := driver.RunSmartContractCall(&vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  []byte("me"),
			Arguments:   [][]byte{},
			CallValue:   big.NewInt(0),
			GasPrice:    100000000,
			GasProvided: 2000000,
		},
		RecipientAddr: []byte("contract"),
		Function:      "foo",
	})

	require.Nil(t, vmOutput)
	require.NotNil(t, err)
}

// Test, restarts upon critical error

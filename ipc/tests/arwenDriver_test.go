package tests

import (
	"fmt"
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
	cryptoHook := &mock.CryptoHookMock{}

	// blockchain.GetCodeCalled = func(address []byte) ([]byte, error) {
	// 	fmt.Println(":::: I WAIT")
	// 	time.Sleep(5 * time.Second)
	// 	fmt.Println(":::: I WAIT DONE")
	// 	return nil, nil
	// }

	driver, err := nodepart.NewArwenDriver(blockchain, cryptoHook, arwenVirtualMachine, uint64(10000000), config.MakeGasMap(1))
	require.Nil(t, err)
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

	fmt.Println("VMoutput", vmOutput)
	fmt.Println("err", err)

	//wait for the...
}

// Test, restarts upon critical error

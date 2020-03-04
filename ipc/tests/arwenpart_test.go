package tests

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/arwenpart"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/nodepart"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testFiles struct {
	outputOfNode  *os.File
	inputOfArwen  *os.File
	outputOfArwen *os.File
	inputOfNode   *os.File
}

func TestArwenPart_SendBadRequest(t *testing.T) {
	blockchain := &mock.BlockChainHookStub{}
	cryptoHook := &mock.CryptoHookMock{}

	response, err := doContractRequest(t, "1", &common.ContractRequest{Action: "foobar"}, blockchain, cryptoHook)
	require.Nil(t, response)
	require.Error(t, err, common.ErrBadRequestFromNode)
}

func TestArwenPart_SendDeployRequest(t *testing.T) {
	blockchain := &mock.BlockChainHookStub{}
	cryptoHook := &mock.CryptoHookMock{}

	response, err := doContractRequest(t, "2", createDeployRequest(), blockchain, cryptoHook)
	require.NotNil(t, response)
	require.Nil(t, err)
}

func TestArwenPart_SendCallRequestWhenNoContract(t *testing.T) {
	blockchain := &mock.BlockChainHookStub{}
	cryptoHook := &mock.CryptoHookMock{}

	response, err := doContractRequest(t, "3", createCallRequest("increment"), blockchain, cryptoHook)
	require.NotNil(t, response)
	require.Nil(t, err)
}

func TestArwenPart_SendCallRequest(t *testing.T) {
	blockchain := &mock.BlockChainHookStub{}
	cryptoHook := &mock.CryptoHookMock{}

	blockchain.GetCodeCalled = func(address []byte) ([]byte, error) {
		return getSCCode("./../../test/contracts/counter.wasm"), nil
	}
	response, err := doContractRequest(t, "3", createCallRequest("increment"), blockchain, cryptoHook)
	require.NotNil(t, response)
	require.Nil(t, err)
}

func doContractRequest(
	t *testing.T,
	tag string,
	request *common.ContractRequest,
	blockchain vmcommon.BlockchainHook,
	cryptoHook vmcommon.CryptoHook,
) (*common.HookCallRequestOrContractResponse, error) {
	files := createTestFiles(t, tag)
	var response *common.HookCallRequestOrContractResponse
	var responseError error

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		part, err := arwenpart.NewArwenPart(files.inputOfArwen, files.outputOfArwen)
		assert.Nil(t, err)
		part.StartLoop()
		wg.Done()
	}()

	go func() {
		part, err := nodepart.NewNodePart(files.inputOfNode, files.outputOfNode, blockchain, cryptoHook)
		assert.Nil(t, err)
		response, responseError = part.StartLoop(request)
		part.SendStopSignal()
		wg.Done()
	}()

	wg.Wait()

	return response, responseError
}

func createTestFiles(t *testing.T, tag string) testFiles {
	files := testFiles{}

	var err error
	files.inputOfArwen, files.outputOfNode, err = os.Pipe()
	require.Nil(t, err)
	files.inputOfNode, files.outputOfArwen, err = os.Pipe()
	require.Nil(t, err)

	return files
}

func createDeployRequest() *common.ContractRequest {
	path := "./../../test/contracts/counter.wasm"
	code := getSCCode(path)

	return &common.ContractRequest{
		Action: "Deploy",
		CreateInput: &vmcommon.ContractCreateInput{
			VMInput: vmcommon.VMInput{
				CallerAddr:  []byte("me"),
				Arguments:   [][]byte{},
				CallValue:   big.NewInt(0),
				GasPrice:    100000000,
				GasProvided: 2000000,
			},
			ContractCode: code,
		},
	}
}

func getSCCode(fileName string) []byte {
	code, err := ioutil.ReadFile(filepath.Clean(fileName))
	if err != nil {
		panic(fmt.Sprintf("Cannot read file [%s].", fileName))
	}

	return code
}

func createCallRequest(function string) *common.ContractRequest {
	return &common.ContractRequest{
		Action: "Call",
		CallInput: &vmcommon.ContractCallInput{
			VMInput: vmcommon.VMInput{
				CallerAddr:  []byte("me"),
				Arguments:   [][]byte{},
				CallValue:   big.NewInt(0),
				GasPrice:    100000000,
				GasProvided: 2000000,
			},
			RecipientAddr: []byte("contract"),
			Function:      function,
		},
	}
}

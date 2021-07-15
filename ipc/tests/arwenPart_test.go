package tests

import (
	"os"
	"sync"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/ipc/arwenpart"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/ipc/common"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/ipc/marshaling"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/ipc/nodepart"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/world"
	"github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/builtInFunctions"
	"github.com/ElrondNetwork/elrond-vm-common/parsers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testFiles struct {
	outputOfNode  *os.File
	inputOfArwen  *os.File
	outputOfArwen *os.File
	inputOfNode   *os.File
}

func TestArwenPart_SendDeployRequest(t *testing.T) {
	blockchain := &contextmock.BlockchainHookStub{}

	response, err := doContractRequest(t, "2", createDeployRequest(bytecodeCounter), blockchain)
	require.NotNil(t, response)
	require.Nil(t, err)
}

func TestArwenPart_SendCallRequestWhenNoContract(t *testing.T) {
	blockchain := &contextmock.BlockchainHookStub{}

	response, err := doContractRequest(t, "3", createCallRequest("increment"), blockchain)
	require.NotNil(t, response)
	require.Nil(t, err)
}

func TestArwenPart_SendCallRequest(t *testing.T) {
	blockchain := &contextmock.BlockchainHookStub{}

	blockchain.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		return &worldmock.Account{Code: bytecodeCounter}, nil
	}

	response, err := doContractRequest(t, "3", createCallRequest("increment"), blockchain)
	require.NotNil(t, response)
	require.Nil(t, err)
}

func doContractRequest(
	t *testing.T,
	tag string,
	request common.MessageHandler,
	blockchain vmcommon.BlockchainHook,
) (common.MessageHandler, error) {
	files := createTestFiles(t, tag)
	var response common.MessageHandler
	var responseError error

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		esdtTransferParser, _ := parsers.NewESDTTransferParser(worldmock.WorldMarshalizer)
		vmHostParameters := &arwen.VMHostParameters{
			VMType:                   []byte{5, 0},
			BlockGasLimit:            uint64(10000000),
			GasSchedule:              config.MakeGasMapForTests(),
			ElrondProtectedKeyPrefix: []byte("ELROND"),
			BuiltInFuncContainer:     builtInFunctions.NewBuiltInFunctionContainer(),
			ESDTTransferParser:       esdtTransferParser,
		}

		part, err := arwenpart.NewArwenPart(
			"testversion",
			files.inputOfArwen,
			files.outputOfArwen,
			vmHostParameters,
			marshaling.CreateMarshalizer(marshaling.JSON),
		)
		assert.Nil(t, err)
		_ = part.StartLoop()
		wg.Done()
	}()

	go func() {
		part, err := nodepart.NewNodePart(
			files.inputOfNode,
			files.outputOfNode,
			blockchain,
			nodepart.Config{MaxLoopTime: 1000},
			marshaling.CreateMarshalizer(marshaling.JSON),
		)
		assert.Nil(t, err)
		response, responseError = part.StartLoop(request)
		_ = part.SendStopSignal()
		wg.Done()
	}()

	wg.Wait()

	return response, responseError
}

func createTestFiles(t *testing.T, _ string) testFiles {
	files := testFiles{}

	var err error
	files.inputOfArwen, files.outputOfNode, err = os.Pipe()
	require.Nil(t, err)
	files.inputOfNode, files.outputOfArwen, err = os.Pipe()
	require.Nil(t, err)

	return files
}

func createDeployRequest(contractCode []byte) common.MessageHandler {
	return common.NewMessageContractDeployRequest(createDeployInput(contractCode))
}

func createCallRequest(function string) common.MessageHandler {
	return common.NewMessageContractCallRequest(createCallInput(function))
}

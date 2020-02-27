package tests

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/arwenpart"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/nodepart"
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
	response, err := doContractRequest(t, "SendBadRequest", &common.ContractRequest{Tag: "foobar"})
	require.Nil(t, response)
	require.Error(t, err, common.ErrBadRequestFromNode)
}

func TestArwenPart_SendDeployRequest(t *testing.T) {
	response, err := doContractRequest(t, "SendDeployRequest", createDeployRequest())
	require.Nil(t, response)
	require.Error(t, err, common.ErrBadRequestFromNode)
}

func createDeployRequest() *common.ContractRequest {
	return &common.ContractRequest{
		Tag: "Deploy",
		CreateInput: &vmcommon.ContractCreateInput{
			VMInput: vmcommon.VMInput{
				CallerAddr:  []byte{},
				Arguments:   [][]byte{},
				CallValue:   big.NewInt(0),
				GasPrice:    100000000,
				GasProvided: 2000000,
			},
			ContractCode: []byte{},
		},
	}
}

func doContractRequest(t *testing.T, tag string, request *common.ContractRequest) (*common.HookCallRequestOrContractResponse, error) {
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
		part, err := nodepart.NewNodePart(files.inputOfNode, files.outputOfNode)
		assert.Nil(t, err)
		response, responseError = part.StartLoop(request)
		wg.Done()
	}()

	wg.Wait()

	return response, responseError
}

func createTestFiles(t *testing.T, tag string) testFiles {
	folder := filepath.Join(".", "testdata", "streams")
	os.MkdirAll(folder, os.ModePerm)

	nodeToArwen := filepath.Join(folder, fmt.Sprintf("node-to-arwen-%s.bin", tag))
	arwenToNode := filepath.Join(folder, fmt.Sprintf("arwen-to-node-%s.bin", tag))

	files := testFiles{}

	var err error
	files.outputOfNode, err = os.Create(nodeToArwen)
	require.Nil(t, err)
	files.outputOfArwen, err = os.Create(arwenToNode)
	require.Nil(t, err)
	files.inputOfNode, err = os.Open(arwenToNode)
	require.Nil(t, err)
	files.inputOfArwen, err = os.Open(nodeToArwen)
	require.Nil(t, err)

	return files
}

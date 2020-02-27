package tests

import (
	"bufio"
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

func TestChildServer_SendBadRequest(t *testing.T) {
	flow := func(node *nodepart.NodeMessenger) {
		response, err := node.SendContractRequest(&common.ContractRequest{Tag: "foobar"})
		assert.Nil(t, response)
		assert.Error(t, err, common.ErrBadRequestFromNode)
	}

	runChildServer(t, "foo", flow)
}

func TestChildServer_SendDeployRequest(t *testing.T) {
	flow := func(node *nodepart.NodeMessenger) {
		response, err := node.SendContractRequest(createDeployRequest())
		assert.Nil(t, response)
		assert.Error(t, err, common.ErrBadRequestFromNode)
		_, err = node.SendContractRequest(&common.ContractRequest{Tag: "Stop"})
		assert.Error(t, err, common.ErrStopPerNodeRequest)
	}

	runChildServer(t, "bar", flow)
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

func runChildServer(t *testing.T, tag string, nodeFlow func(node *nodepart.NodeMessenger)) {
	files := createTestFiles(t, tag)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		server, err := arwenpart.NewChildServer(files.inputOfArwen, files.outputOfArwen)
		assert.Nil(t, err)
		server.Start()
		wg.Done()
	}()

	go func() {
		node := nodepart.NewNodeMessenger(bufio.NewReader(files.inputOfNode), bufio.NewWriter(files.outputOfNode))
		nodeFlow(node)
		wg.Done()
	}()

	wg.Wait()
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

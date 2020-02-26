package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type testFiles struct {
	outputOfNode  *os.File
	inputOfArwen  *os.File
	outputOfArwen *os.File
	inputOfNode   *os.File
}

func Test_Loop(t *testing.T) {
	files := createTestFiles(t, "foo")

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		doMain(files.inputOfArwen, files.outputOfArwen)
		wg.Done()
	}()

	go func() {
		node := NewNodeMessenger(bufio.NewReader(files.inputOfNode), bufio.NewWriter(files.outputOfNode))
		response, err := node.SendContractRequest(&ContractRequest{Tag: "foobar"})
		require.Nil(t, response)
		require.Equal(t, ErrCannotSendContractRequest, err)
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

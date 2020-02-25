package main

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Loop(t *testing.T) {
	nodeToArwen := "testdata/node-to-arwen_42"
	arwenToNode := "testdata/arwen-to-node_42"

	outputOfNode, err := os.Create(nodeToArwen)
	require.Nil(t, err)
	outputOfArwen, err := os.Create(arwenToNode)
	require.Nil(t, err)
	_, err = os.Open(arwenToNode) // inputOfNode
	require.Nil(t, err)
	inputOfArwen, err := os.Open(nodeToArwen)
	require.Nil(t, err)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		beginMessageLoop(bufio.NewReader(inputOfArwen), bufio.NewWriter(outputOfArwen))
	}()

	go func() {
		writeUint32(outputOfNode, 3)
		outputOfNode.Write([]byte("foobar"))
	}()

	wg.Wait()
	//time.Sleep(3 * time.Second)
}

func writeUint32(file *os.File, value uint32) {
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, value)
	file.Write(buffer)
}

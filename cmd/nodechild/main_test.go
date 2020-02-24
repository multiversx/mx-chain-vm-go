package main

import (
	"bufio"
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

	outputOfNode.Write([]byte("foo"))

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		beginMessageLoop(bufio.NewReader(inputOfArwen), bufio.NewWriter(outputOfArwen))
	}()

	go func() {
		outputOfNode.Write([]byte("foo"))
	}()

	wg.Wait()
	//time.Sleep(3 * time.Second)
}

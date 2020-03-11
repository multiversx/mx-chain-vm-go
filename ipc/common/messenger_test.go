package common

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSetDeadline(t *testing.T) {
	readFile, writeFile, err := os.Pipe()
	require.Nil(t, err)

	future := time.Now().Add(500 * time.Millisecond)
	err = readFile.SetReadDeadline(future)
	require.Nil(t, err)

	go func() {
		time.Sleep(2 * time.Second)
		writeFile.WriteString("foo")
	}()

	buff := make([]byte, 100)
	n, err := io.ReadFull(readFile, buff)
	require.Nil(t, err)
	require.Equal(t, 3, n)
}

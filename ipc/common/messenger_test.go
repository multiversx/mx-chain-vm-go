package common

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSetDeadline(t *testing.T) {
	r, w, err := os.Pipe()
	require.Nil(t, err)

	future := time.Now().Add(500 * time.Millisecond)
	err = r.SetReadDeadline(future)
	require.Nil(t, err)

	go func() {
		time.Sleep(2 * time.Second)
		w.WriteString("foo")
	}()

	buff := make([]byte, 100)
	n, err := io.ReadFull(r, buff) // r.Read(buff)
	require.Nil(t, err)
	require.Equal(t, 3, n)
}

package wasmer2

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLastError_ThreadPinning(t *testing.T) {
	// Calling GetLastError returns nil Go error
	_, err := GetLastError()
	assert.NoError(t, err)

	// Ensure concurrent calls from multiple goroutines don't panic or corrupt thread affinity
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			runtime.LockOSThread()
			defer runtime.UnlockOSThread()
			for j := 0; j < 50; j++ {
				_, _ = GetLastError()
			}
			done <- struct{}{}
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

package hosttest

import (
	"testing"

	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
	"github.com/stretchr/testify/assert"
)

func TestElrondEI_validateToken(t *testing.T) {
	var result bool
	result = vmhooks.ValidateToken([]byte("EGLDRIDEFL-08d8eff"))
	assert.False(t, result)
	result = vmhooks.ValidateToken([]byte("EGLDRIDEFL-08d8e"))
	assert.False(t, result)
	result = vmhooks.ValidateToken([]byte("EGLDRIDEFL08d8ef"))
	assert.False(t, result)
	result = vmhooks.ValidateToken([]byte("EGLDRIDEFl-08d8ef"))
	assert.False(t, result)
	result = vmhooks.ValidateToken([]byte("EGLDRIDEF*-08d8ef"))
	assert.False(t, result)
	result = vmhooks.ValidateToken([]byte("EGLDRIDEFL-08d8eF"))
	assert.False(t, result)
	result = vmhooks.ValidateToken([]byte("EGLDRIDEFL-08d*ef"))
	assert.False(t, result)

	result = vmhooks.ValidateToken([]byte("ALC6258d2"))
	assert.False(t, result)
	result = vmhooks.ValidateToken([]byte("AL-C6258d2"))
	assert.False(t, result)
	result = vmhooks.ValidateToken([]byte("alc-6258d2"))
	assert.False(t, result)
	result = vmhooks.ValidateToken([]byte("ALC-6258D2"))
	assert.False(t, result)
	result = vmhooks.ValidateToken([]byte("ALC-6258d2ff"))
	assert.False(t, result)
	result = vmhooks.ValidateToken([]byte("AL-6258d2"))
	assert.False(t, result)
	result = vmhooks.ValidateToken([]byte("ALCCCCCCCCC-6258d2"))
	assert.False(t, result)

	result = vmhooks.ValidateToken([]byte("EGLDRIDEF2-08d8ef"))
	assert.True(t, result)
	result = vmhooks.ValidateToken([]byte("EGLDRIDEFL-08d8ef"))
	assert.True(t, result)
	result = vmhooks.ValidateToken([]byte("ALC-6258d2"))
	assert.True(t, result)
	result = vmhooks.ValidateToken([]byte("ALC123-6258d2"))
	assert.True(t, result)
	result = vmhooks.ValidateToken([]byte("12345-6258d2"))
	assert.True(t, result)
}

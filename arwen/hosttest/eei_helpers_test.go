package hosttest

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen/elrondapi"
	"github.com/stretchr/testify/assert"
)

func TestElrondEI_validateToken(t *testing.T) {
	var result int32
	result = elrondapi.ValidateToken([]byte("EGLDRIDEFL-08d8eff"))
	assert.Equal(t, result, int32(0))
	result = elrondapi.ValidateToken([]byte("EGLDRIDEFL-08d8e"))
	assert.Equal(t, result, int32(0))
	result = elrondapi.ValidateToken([]byte("EGLDRIDEFL08d8ef"))
	assert.Equal(t, result, int32(0))
	result = elrondapi.ValidateToken([]byte("EGLDRIDEFl-08d8ef"))
	assert.Equal(t, result, int32(0))
	result = elrondapi.ValidateToken([]byte("EGLDRIDEF*-08d8ef"))
	assert.Equal(t, result, int32(0))
	result = elrondapi.ValidateToken([]byte("EGLDRIDEFL-08d8eF"))
	assert.Equal(t, result, int32(0))
	result = elrondapi.ValidateToken([]byte("EGLDRIDEFL-08d*ef"))
	assert.Equal(t, result, int32(0))

	result = elrondapi.ValidateToken([]byte("ALC6258d2"))
	assert.Equal(t, result, int32(0))
	result = elrondapi.ValidateToken([]byte("AL-C6258d2"))
	assert.Equal(t, result, int32(0))
	result = elrondapi.ValidateToken([]byte("alc-6258d2"))
	assert.Equal(t, result, int32(0))
	result = elrondapi.ValidateToken([]byte("ALC-6258D2"))
	assert.Equal(t, result, int32(0))
	result = elrondapi.ValidateToken([]byte("ALC-6258d2ff"))
	assert.Equal(t, result, int32(0))
	result = elrondapi.ValidateToken([]byte("AL-6258d2"))
	assert.Equal(t, result, int32(0))
	result = elrondapi.ValidateToken([]byte("ALCCCCCCCCC-6258d2"))
	assert.Equal(t, result, int32(0))

	result = elrondapi.ValidateToken([]byte("EGLDRIDEF2-08d8ef"))
	assert.Equal(t, result, int32(1))
	result = elrondapi.ValidateToken([]byte("EGLDRIDEFL-08d8ef"))
	assert.Equal(t, result, int32(1))
	result = elrondapi.ValidateToken([]byte("ALC-6258d2"))
	assert.Equal(t, result, int32(1))
	result = elrondapi.ValidateToken([]byte("ALC123-6258d2"))
	assert.Equal(t, result, int32(1))
	result = elrondapi.ValidateToken([]byte("12345-6258d2"))
	assert.Equal(t, result, int32(1))
}

package arwen

import (
	"github.com/ElrondNetwork/go-ext-wasm/wasmer"
	"sync"
	"unsafe"
)

const addressLen = 32
const hashLen = 32
const balanceLen = 32

var (
	vmContextCounter uint8
	vmContextMap     map[uint8]*vmContext
	vmContextMapMu   sync.Mutex
)

func addHostContext(ctx *vmContext) int {
	vmContextMapMu.Lock()
	id := vmContextCounter
	vmContextCounter++
	vmContextMap[id] = ctx
	vmContextMapMu.Unlock()
	return int(id)
}

func removeHostContext(idx int) {
	vmContextMapMu.Lock()
	delete(vmContextMap, uint8(idx))
	vmContextMapMu.Unlock()
}

func GetHostContext(pointer unsafe.Pointer) HostContext {
	var idx = *(*int)(pointer)

	vmContextMapMu.Lock()
	ctx := vmContextMap[uint8(idx)]
	vmContextMapMu.Unlock()

	return ctx
}

func GetEthContext(pointer unsafe.Pointer) EthContext {
	var idx = *(*int)(pointer)

	vmContextMapMu.Lock()
	ctx := vmContextMap[uint8(idx)]
	vmContextMapMu.Unlock()

	return ctx
}

func LoadBytes(from *wasmer.Memory, offset int32, length int32) []byte {
	if from.Length() < uint32(offset+length) {
		return from.Data()[offset:]
	}

	return from.Data()[offset : offset+length]
}

func StoreBytes(to *wasmer.Memory, offset int32, data []byte) error {
	var memoryData = to.Data()
	length := int32(len(data))

	if to.Length() < uint32(offset+length) {
		err := to.Grow(1)
		if err != nil {
			return err
		}
	}

	copy(memoryData[offset:offset+length], data)

	return nil
}

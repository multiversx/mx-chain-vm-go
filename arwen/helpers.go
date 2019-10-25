package arwen

import (
	"math/big"

	"github.com/ElrondNetwork/go-ext-wasm/wasmer"
)

const addressLen = 32
const hashLen = 32
const balanceLen = 32

func loadBytes(from *wasmer.Memory, offset int32, length int32) []byte {
	if from.Length() < uint32(offset+length) {
		return from.Data()[offset:]
	}

	return from.Data()[offset : offset+length]
}

func storeBytes(to *wasmer.Memory, offset int32, data []byte) error {
	length := int32(len(data))

	if to.Length() < uint32(offset+length) {
		err := to.Grow(1)
		if err != nil {
			return err
		}
	}

	var memoryData = to.Data()
	copy(memoryData[offset:offset+length], data)

	return nil
}

func convertReturnValue(wasmValue wasmer.Value) *big.Int {
	switch wasmValue.GetType() {
	case wasmer.TypeVoid:
		return big.NewInt(0)
	case wasmer.TypeI32:
		return big.NewInt(wasmValue.ToI64())
	case wasmer.TypeI64:
		return big.NewInt(wasmValue.ToI64())
	}

	panic("unsupported return type")
}

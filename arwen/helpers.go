package arwen

import (
	"fmt"
	"math/big"

	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

const AddressLen = 32
const AddressLenEth = 20
const HashLen = 32
const ArgumentLenEth = 32
const BalanceLen = 32
const InitFunctionName = "init"
const InitFunctionNameEth = "solidity.ctor"

var (
	vmContextCounter uint8
	vmContextMap     map[uint8]VMContext
)

func AddHostContext(ctx VMContext) int {
	id := vmContextCounter
	vmContextCounter++
	if vmContextMap == nil {
		vmContextMap = make(map[uint8]VMContext)
	}
	vmContextMap[id] = ctx
	return int(id)
}

func RemoveHostContext(idx int) {
	delete(vmContextMap, uint8(idx))
}

func GetVmContext(context unsafe.Pointer) VMContext {
	instCtx := wasmer.IntoInstanceContext(context)
	var idx = *(*int)(instCtx.Data())

	ctx := vmContextMap[uint8(idx)]
	ctx.Runtime().SetInstanceContext(&instCtx)

	return ctx
}

func GetBlockchainSubcontext(context unsafe.Pointer) BlockchainSubcontext {
	return GetVmContext(context).Blockchain()
}

func GetRuntimeSubcontext(context unsafe.Pointer) RuntimeSubcontext {
	return GetVmContext(context).Runtime()
}

func GetCryptoSubcontext(context unsafe.Pointer) vmcommon.CryptoHook {
	return GetVmContext(context).Crypto()
}

func GetBigIntSubcontext(context unsafe.Pointer) BigIntSubcontext {
	return GetVmContext(context).BigInt()
}

func GetOutputSubcontext(context unsafe.Pointer) OutputSubcontext {
	return GetVmContext(context).Output()
}

func GetMeteringSubcontext(context unsafe.Pointer) MeteringSubcontext {
	return GetVmContext(context).Metering()
}

func GetStorageSubcontext(context unsafe.Pointer) StorageSubcontext {
	return GetVmContext(context).Storage()
}

func ConvertReturnValue(wasmValue wasmer.Value) *big.Int {
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

func GuardedMakeByteSlice2D(length int32) ([][]byte, error) {
	if length < 0 {
		return nil, fmt.Errorf("GuardedMakeByteSlice2D: negative length (%d)", length)
	}

	result := make([][]byte, length)
	return result, nil
}

func GuardedGetBytesSlice(data []byte, offset int32, length int32) ([]byte, error) {
	dataLength := uint32(len(data))
	isOffsetTooSmall := offset < 0
	isOffsetTooLarge := uint32(offset) > dataLength
	requestedEnd := uint32(offset + length)
	isRequestedEndTooLarge := requestedEnd > dataLength
	isLengthNegative := length < 0

	if isOffsetTooSmall || isOffsetTooLarge {
		return nil, fmt.Errorf("GuardedGetBytesSlice: bad bounds")
	}

	if isRequestedEndTooLarge {
		return nil, fmt.Errorf("GuardedGetBytesSlice: bad bounds")
	}

	if isLengthNegative {
		return nil, fmt.Errorf("GuardedGetBytesSlice: negative length")
	}

	result := data[offset : offset+length]
	return result, nil
}

func InverseBytes(data []byte) []byte {
	length := len(data)
	invBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		invBytes[length-i-1] = data[i]
	}
	return invBytes
}

package arwen

import (
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
)

const AddressLen = 32
const AddressLenEth = 20
const HashLen = 32
const ArgumentLenEth = 32
const BalanceLen = 32
const InitFunctionName = "init"
const InitFunctionNameEth = "solidity.ctor"
const CallBackFunctionName = "callBack"
const UpgradeFunctionName = "upgradeContract"

var (
	vmContextCounter uint8
	vmContextMap     map[uint8]VMHost
)

func AddHostContext(ctx VMHost) int {
	id := vmContextCounter
	vmContextCounter++
	if vmContextMap == nil {
		vmContextMap = make(map[uint8]VMHost)
	}
	vmContextMap[id] = ctx
	return int(id)
}

func RemoveAllHostContexts() {
	vmContextMap = make(map[uint8]VMHost)
}

func RemoveHostContext(idx int) {
	delete(vmContextMap, uint8(idx))
}

func GetVmContext(context unsafe.Pointer) VMHost {
	instCtx := wasmer.IntoInstanceContext(context)
	var idx = *(*int)(instCtx.Data())

	return vmContextMap[uint8(idx)]
}

func GetBlockchainContext(context unsafe.Pointer) BlockchainContext {
	return GetVmContext(context).Blockchain()
}

func GetRuntimeContext(context unsafe.Pointer) RuntimeContext {
	return GetVmContext(context).Runtime()
}

func GetCryptoContext(context unsafe.Pointer) crypto.VMCrypto {
	return GetVmContext(context).Crypto()
}

func GetBigIntContext(context unsafe.Pointer) BigIntContext {
	return GetVmContext(context).BigInt()
}

func GetOutputContext(context unsafe.Pointer) OutputContext {
	return GetVmContext(context).Output()
}

func GetMeteringContext(context unsafe.Pointer) MeteringContext {
	return GetVmContext(context).Metering()
}

func GetStorageContext(context unsafe.Pointer) StorageContext {
	return GetVmContext(context).Storage()
}

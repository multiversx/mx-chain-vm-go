package crypto

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t sha256(void* context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t keccak256(void *context, int32_t dataOffset, int32_t length, int32_t resultOffset);
import "C"

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	"unsafe"
)

func CryptoImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("env")
	imports, err := imports.Append("sha256", sha256, C.sha256)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("keccak256", keccak256, C.keccak256)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export sha256
func sha256(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)

	data, err := runtime.MemLoad(dataOffset, length)
	if withFault(err, context) {
		return 1
	}

	result, err := crypto.Sha256(data)
	if err != nil {
		return 1
	}

	err = runtime.MemStore(resultOffset, result)
	if withFault(err, context) {
		return 1
	}

	metering := arwen.GetMeteringContext(context)
	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	return 0
}

//export keccak256
func keccak256(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)

	data, err := runtime.MemLoad(dataOffset, length)
	if withFault(err, context) {
		return 1
	}

	result, err := crypto.Keccak256(data)
	if err != nil {
		return 1
	}

	err = runtime.MemStore(resultOffset, result)
	if withFault(err, context) {
		return 1
	}

	metering := arwen.GetMeteringContext(context)
	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	return 0
}

func withFault(err error, context unsafe.Pointer) bool {
	if err != nil {
		runtime := arwen.GetRuntimeContext(context)
		metering := arwen.GetMeteringContext(context)

		runtime.SignalUserError(err.Error())
		metering.UseGas(metering.GasLeft())

		return true
	}

	return false
}

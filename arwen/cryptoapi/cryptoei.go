package cryptoapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t sha256(void* context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t keccak256(void *context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t ripemd160(void *context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t blsVerify(void *context, int32_t keyOffset, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern int32_t ed25519Verify(void *context, int32_t keyOffset, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern int32_t secp256k1Verify(void *context, int32_t keyOffset, int32_t keyLength, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
import "C"

import (
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
)

const blsPublicKeyLength = 96
const blsSignatureLength = 48
const ed25519PublicKeyLength = 32
const ed25519SignatureLength = 64
const secp256k1CompressedPublicKeyLength = 33
const secp256k1UncompressedPublicKeyLength = 65
const secp256k1SignatureLength = 64

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
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	result, err := crypto.Sha256(data)
	if err != nil {
		return 1
	}

	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
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
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	result, err := crypto.Keccak256(data)
	if err != nil {
		return 1
	}

	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	metering := arwen.GetMeteringContext(context)
	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	return 0
}

//export ripemd160
func ripemd160(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	result, err := crypto.Ripemd160(data)
	if err != nil {
		return 1
	}

	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	metering := arwen.GetMeteringContext(context)
	gasToUse := metering.GasSchedule().CryptoAPICost.Ripemd160
	metering.UseGas(gasToUse)

	return 0
}

//export blsVerify
func blsVerify(
	context unsafe.Pointer,
	keyOffset int32,
	messageOffset int32,
	messageLength int32,
	sigOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)

	key, err := runtime.MemLoad(keyOffset, blsPublicKeyLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	message, err := runtime.MemLoad(messageOffset, messageLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	sig, err := runtime.MemLoad(sigOffset, blsSignatureLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	invalidSigErr := crypto.BLSVerify(key, message, sig)

	metering := arwen.GetMeteringContext(context)
	gasToUse := metering.GasSchedule().CryptoAPICost.BLSVerify
	metering.UseGas(gasToUse)

	if invalidSigErr != nil {
		return 1
	}

	return 0
}

//export ed25519Verify
func ed25519Verify(
	context unsafe.Pointer,
	keyOffset int32,
	messageOffset int32,
	messageLength int32,
	sigOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)

	key, err := runtime.MemLoad(keyOffset, ed25519PublicKeyLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	message, err := runtime.MemLoad(messageOffset, messageLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	sig, err := runtime.MemLoad(sigOffset, ed25519SignatureLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	invalidSigErr := crypto.Ed25519Verify(key, message, sig)

	metering := arwen.GetMeteringContext(context)
	gasToUse := metering.GasSchedule().CryptoAPICost.Ed25519Verify
	metering.UseGas(gasToUse)

	if invalidSigErr != nil {
		return 1
	}

	return 0
}

//export secp256k1Verify
func secp256k1Verify(
	context unsafe.Pointer,
	keyOffset int32,
	keyLength int32,
	messageOffset int32,
	messageLength int32,
	sigOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)

	if keyLength != secp256k1CompressedPublicKeyLength && keyLength != secp256k1UncompressedPublicKeyLength {
		arwen.WithFault(arwen.ErrInvalidPublicKeySize, context, runtime.ElrondAPIErrorShouldFailExecution())
		return 1
	}

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	message, err := runtime.MemLoad(messageOffset, messageLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	sig, err := runtime.MemLoad(sigOffset, secp256k1SignatureLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	invalidSigErr := crypto.Secp256k1Verify(key, message, sig)

	metering := arwen.GetMeteringContext(context)
	gasToUse := metering.GasSchedule().CryptoAPICost.Secp256k1Verify
	metering.UseGas(gasToUse)

	if invalidSigErr != nil {
		return 1
	}

	return 0
}

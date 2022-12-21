package cryptoapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t v1_4_sha256(void* context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t v1_4_managedSha256(void *context, int32_t inputHanle, int32_t outputHandle);
// extern int32_t v1_4_keccak256(void *context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t v1_4_managedKeccak256(void *context, int32_t inputHanle, int32_t outputHandle);
// extern int32_t v1_4_ripemd160(void *context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t v1_4_managedRipemd160(void *context, int32_t dataHandle, int32_t resultHandle);
// extern int32_t v1_4_verifyBLS(void *context, int32_t keyOffset, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern int32_t v1_4_managedVerifyBLS(void *context, int32_t keyHandle, int32_t messageHandle, int32_t sigHandle);
// extern int32_t v1_4_verifyEd25519(void *context, int32_t keyOffset, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern int32_t v1_4_managedVerifyEd25519(void *context, int32_t keyHandle, int32_t messageHandle, int32_t sigHandle);
// extern int32_t v1_4_verifySecp256k1(void *context, int32_t keyOffset, int32_t keyLength, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern int32_t v1_4_managedVerifySecp256k1(void *context, int32_t keyHandle, int32_t messageHandle, int32_t sigHandle);
// extern int32_t v1_4_verifyCustomSecp256k1(void *context, int32_t keyOffset, int32_t keyLength, int32_t messageOffset, int32_t messageLength, int32_t sigOffset, int32_t hashType);
// extern int32_t v1_4_managedVerifyCustomSecp256k1(void *context, int32_t keyHandle, int32_t messageHandle, int32_t sigHandle, int32_t hashType);
// extern int32_t v1_4_encodeSecp256k1DerSignature(void *context, int32_t rOffset, int32_t rLength, int32_t sOffset, int32_t sLength, int32_t sigOffset);
// extern int32_t v1_4_managedEncodeSecp256k1DerSignature(void *context, int32_t rHandle, int32_t sHandle, int32_t sigHandle);
// extern void v1_4_addEC(void *context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t fstPointXHandle, int32_t fstPointYHandle, int32_t sndPointXHandle, int32_t sndPointYHandle);
// extern void v1_4_doubleEC(void *context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t pointXHandle, int32_t pointYHandle);
// extern int32_t v1_4_isOnCurveEC(void *context, int32_t ecHandle, int32_t pointXHandle, int32_t pointYHandle);
// extern int32_t v1_4_scalarBaseMultEC(void *context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t dataOffset, int32_t length);
// extern int32_t v1_4_managedScalarBaseMultEC(void *context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t dataHandle);
// extern int32_t v1_4_scalarMultEC(void *context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t pointXHandle, int32_t pointYHandle, int32_t dataOffset, int32_t length);
// extern int32_t v1_4_managedScalarMultEC(void *context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t pointXHandle, int32_t pointYHandle, int32_t dataHandle);
// extern int32_t v1_4_marshalEC(void *context, int32_t xPairHandle, int32_t yPairHandle, int32_t ecHandle, int32_t resultOffset);
// extern int32_t v1_4_managedMarshalEC(void *context, int32_t xPairHandle, int32_t yPairHandle, int32_t ecHandle, int32_t resultHandle);
// extern int32_t v1_4_unmarshalEC(void *context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t dataOffset, int32_t length);
// extern int32_t v1_4_managedUnmarshalEC(void *context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t dataHandle);
// extern int32_t v1_4_marshalCompressedEC(void *context, int32_t xPairHandle, int32_t yPairHandle, int32_t ecHandle, int32_t resultOffset);
// extern int32_t v1_4_managedMarshalCompressedEC(void *context, int32_t xPairHandle, int32_t yPairHandle, int32_t ecHandle, int32_t Handle);
// extern int32_t v1_4_unmarshalCompressedEC(void *context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t dataOffset, int32_t length);
// extern int32_t v1_4_managedUnmarshalCompressedEC(void *context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t dataHandle);
// extern int32_t v1_4_generateKeyEC(void *context, int32_t xPubKeyHandle, int32_t yPubKeyHandle, int32_t ecHandle, int32_t resultOffset);
// extern int32_t v1_4_managedGenerateKeyEC(void *context, int32_t xPubKeyHandle, int32_t yPubKeyHandle, int32_t ecHandle, int32_t resultHandle);
// extern int32_t v1_4_createEC(void *context, int32_t dataOffset, int32_t dataLength);
// extern int32_t v1_4_managedCreateEC(void *context, int32_t dataHandle);
// extern int32_t v1_4_getCurveLengthEC(void *context, int32_t ecHandle);
// extern int32_t v1_4_getPrivKeyByteLengthEC(void *context, int32_t ecHandle);
// extern int32_t v1_4_ellipticCurveGetValues(void *context, int32_t ecHandle, int32_t fieldOrderHandle, int32_t basePointOrderHandle, int32_t eqConstantHandle, int32_t xBasePointHandle, int32_t yBasePointHandle);
import "C"

import (
	"crypto/elliptic"
	"unsafe"

	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen"
	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen/elrondapimeta"
	"github.com/ElrondNetwork/wasm-vm-v1_4/crypto/signing/secp256k1"
	"github.com/ElrondNetwork/wasm-vm-v1_4/math"
)

const blsPublicKeyLength = 96
const blsSignatureLength = 48
const ed25519PublicKeyLength = 32
const ed25519SignatureLength = 64
const secp256k1CompressedPublicKeyLength = 33
const secp256k1UncompressedPublicKeyLength = 65
const curveNameLength = 4

const (
	sha256Name                      = "sha256"
	keccak256Name                   = "keccak256"
	ripemd160Name                   = "ripemd160"
	verifyBLSName                   = "verifyBLS"
	verifyEd25519Name               = "verifyEd25519"
	verifyCustomSecp256k1Name       = "verifyCustomSecp256k1"
	encodeSecp256k1DerSignatureName = "encodeSecp256k1DerSignature"
	addECName                       = "addEC"
	doubleECName                    = "doubleEC"
	isOnCurveECName                 = "isOnCurveEC"
	scalarBaseMultECName            = "scalarBaseMultEC"
	scalarMultECName                = "scalarMultEC"
	marshalECName                   = "marshalEC"
	unmarshalECName                 = "unmarshalEC"
	marshalCompressedECName         = "marshalCompressedEC"
	unmarshalCompressedECName       = "unmarshalCompressedEC"
	generateKeyECName               = "generateKeyEC"
	createECName                    = "createEC"
	getCurveLengthECName            = "getCurveLengthEC"
	getPrivKeyByteLengthECName      = "getPrivKeyByteLengthEC"
	ellipticCurveGetValuesName      = "ellipticCurveGetValues"
)

// CryptoImports adds some crypto imports to the Wasmer Imports map
func CryptoImports(imports elrondapimeta.EIFunctionReceiver) error {
	imports.Namespace("env")
	err := imports.Append("sha256", v1_4_sha256, C.v1_4_sha256)
	if err != nil {
		return err
	}

	err = imports.Append("managedSha256", v1_4_managedSha256, C.v1_4_managedSha256)
	if err != nil {
		return err
	}

	err = imports.Append("keccak256", v1_4_keccak256, C.v1_4_keccak256)
	if err != nil {
		return err
	}

	err = imports.Append("managedKeccak256", v1_4_managedKeccak256, C.v1_4_managedKeccak256)
	if err != nil {
		return err
	}

	err = imports.Append("ripemd160", v1_4_ripemd160, C.v1_4_ripemd160)
	if err != nil {
		return err
	}

	err = imports.Append("managedRipemd160", v1_4_managedRipemd160, C.v1_4_managedRipemd160)
	if err != nil {
		return err
	}

	err = imports.Append("verifyBLS", v1_4_verifyBLS, C.v1_4_verifyBLS)
	if err != nil {
		return err
	}

	err = imports.Append("managedVerifyBLS", v1_4_managedVerifyBLS, C.v1_4_managedVerifyBLS)
	if err != nil {
		return err
	}

	err = imports.Append("verifyEd25519", v1_4_verifyEd25519, C.v1_4_verifyEd25519)
	if err != nil {
		return err
	}

	err = imports.Append("managedVerifyEd25519", v1_4_managedVerifyEd25519, C.v1_4_managedVerifyEd25519)
	if err != nil {
		return err
	}

	err = imports.Append("verifySecp256k1", v1_4_verifySecp256k1, C.v1_4_verifySecp256k1)
	if err != nil {
		return err
	}

	err = imports.Append("managedVerifySecp256k1", v1_4_managedVerifySecp256k1, C.v1_4_managedVerifySecp256k1)
	if err != nil {
		return err
	}

	err = imports.Append("verifyCustomSecp256k1", v1_4_verifyCustomSecp256k1, C.v1_4_verifyCustomSecp256k1)
	if err != nil {
		return err
	}

	err = imports.Append("managedVerifyCustomSecp256k1", v1_4_managedVerifyCustomSecp256k1, C.v1_4_managedVerifyCustomSecp256k1)
	if err != nil {
		return err
	}

	err = imports.Append("encodeSecp256k1DerSignature", v1_4_encodeSecp256k1DerSignature, C.v1_4_encodeSecp256k1DerSignature)
	if err != nil {
		return err
	}

	err = imports.Append("managedEncodeSecp256k1DerSignature", v1_4_managedEncodeSecp256k1DerSignature, C.v1_4_managedEncodeSecp256k1DerSignature)
	if err != nil {
		return err
	}

	err = imports.Append("addEC", v1_4_addEC, C.v1_4_addEC)
	if err != nil {
		return err
	}

	err = imports.Append("doubleEC", v1_4_doubleEC, C.v1_4_doubleEC)
	if err != nil {
		return err
	}

	err = imports.Append("isOnCurveEC", v1_4_isOnCurveEC, C.v1_4_isOnCurveEC)
	if err != nil {
		return err
	}

	err = imports.Append("scalarBaseMultEC", v1_4_scalarBaseMultEC, C.v1_4_scalarBaseMultEC)
	if err != nil {
		return err
	}

	err = imports.Append("managedScalarBaseMultEC", v1_4_managedScalarBaseMultEC, C.v1_4_managedScalarBaseMultEC)
	if err != nil {
		return err
	}

	err = imports.Append("scalarMultEC", v1_4_scalarMultEC, C.v1_4_scalarMultEC)
	if err != nil {
		return err
	}

	err = imports.Append("managedScalarMultEC", v1_4_managedScalarMultEC, C.v1_4_managedScalarMultEC)
	if err != nil {
		return err
	}

	err = imports.Append("marshalEC", v1_4_marshalEC, C.v1_4_marshalEC)
	if err != nil {
		return err
	}

	err = imports.Append("managedMarshalEC", v1_4_managedMarshalEC, C.v1_4_managedMarshalEC)
	if err != nil {
		return err
	}

	err = imports.Append("unmarshalEC", v1_4_unmarshalEC, C.v1_4_unmarshalEC)
	if err != nil {
		return err
	}

	err = imports.Append("managedUnmarshalEC", v1_4_managedUnmarshalEC, C.v1_4_managedUnmarshalEC)
	if err != nil {
		return err
	}

	err = imports.Append("marshalCompressedEC", v1_4_marshalCompressedEC, C.v1_4_marshalCompressedEC)
	if err != nil {
		return err
	}

	err = imports.Append("managedMarshalCompressedEC", v1_4_managedMarshalCompressedEC, C.v1_4_managedMarshalCompressedEC)
	if err != nil {
		return err
	}

	err = imports.Append("unmarshalCompressedEC", v1_4_unmarshalCompressedEC, C.v1_4_unmarshalCompressedEC)
	if err != nil {
		return err
	}

	err = imports.Append("managedUnmarshalCompressedEC", v1_4_managedUnmarshalCompressedEC, C.v1_4_managedUnmarshalCompressedEC)
	if err != nil {
		return err
	}

	err = imports.Append("generateKeyEC", v1_4_generateKeyEC, C.v1_4_generateKeyEC)
	if err != nil {
		return err
	}

	err = imports.Append("managedGenerateKeyEC", v1_4_managedGenerateKeyEC, C.v1_4_managedGenerateKeyEC)
	if err != nil {
		return err
	}

	err = imports.Append("createEC", v1_4_createEC, C.v1_4_createEC)
	if err != nil {
		return err
	}

	err = imports.Append("managedCreateEC", v1_4_managedCreateEC, C.v1_4_managedCreateEC)
	if err != nil {
		return err
	}

	err = imports.Append("getCurveLengthEC", v1_4_getCurveLengthEC, C.v1_4_getCurveLengthEC)
	if err != nil {
		return err
	}

	err = imports.Append("getPrivKeyByteLengthEC", v1_4_getPrivKeyByteLengthEC, C.v1_4_getPrivKeyByteLengthEC)
	if err != nil {
		return err
	}

	err = imports.Append("ellipticCurveGetValues", v1_4_ellipticCurveGetValues, C.v1_4_ellipticCurveGetValues)
	if err != nil {
		return err
	}

	return nil
}

//export v1_4_sha256
func v1_4_sha256(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	memLoadGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(metering.GasSchedule().CryptoAPICost.SHA256, memLoadGas)
	metering.UseGasAndAddTracedGas(sha256Name, gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	result, err := crypto.Sha256(data)
	if err != nil {
		arwen.WithFaultIfFailAlwaysActive(err, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export v1_4_managedSha256
func v1_4_managedSha256(context unsafe.Pointer, inputHandle, outputHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	metering.UseGasAndAddTracedGas(sha256Name, metering.GasSchedule().CryptoAPICost.SHA256)

	inputBytes, err := managedType.GetBytes(inputHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(inputBytes)

	resultBytes, err := crypto.Sha256(inputBytes)
	if err != nil {
		arwen.WithFaultIfFailAlwaysActive(err, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.SetBytes(outputHandle, resultBytes)

	return 0
}

//export v1_4_keccak256
func v1_4_keccak256(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	memLoadGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(metering.GasSchedule().CryptoAPICost.Keccak256, memLoadGas)
	metering.UseGasAndAddTracedGas(keccak256Name, gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	result, err := crypto.Keccak256(data)
	if err != nil {
		arwen.WithFaultIfFailAlwaysActive(err, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export v1_4_managedKeccak256
func v1_4_managedKeccak256(context unsafe.Pointer, inputHandle, outputHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	metering.UseGasAndAddTracedGas(keccak256Name, metering.GasSchedule().CryptoAPICost.Keccak256)

	inputBytes, err := managedType.GetBytes(inputHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(inputBytes)

	resultBytes, err := crypto.Keccak256(inputBytes)
	if err != nil {
		arwen.WithFaultIfFailAlwaysActive(err, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.SetBytes(outputHandle, resultBytes)

	return 0
}

//export v1_4_ripemd160
func v1_4_ripemd160(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	memLoadGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(metering.GasSchedule().CryptoAPICost.Ripemd160, memLoadGas)
	metering.UseGasAndAddTracedGas(ripemd160Name, gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	result, err := crypto.Ripemd160(data)
	if err != nil {
		arwen.WithFaultIfFailAlwaysActive(err, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export v1_4_managedRipemd160
func v1_4_managedRipemd160(context unsafe.Pointer, inputHandle int32, outputHandle int32) int32 {
	host := arwen.GetVMHost(context)
	return ManagedRipemd160WithHost(host, inputHandle, outputHandle)
}

func ManagedRipemd160WithHost(host arwen.VMHost, inputHandle int32, outputHandle int32) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	crypto := host.Crypto()

	metering.UseGasAndAddTracedGas(ripemd160Name, metering.GasSchedule().CryptoAPICost.Ripemd160)

	inputBytes, err := managedType.GetBytes(inputHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(inputBytes)

	result, err := crypto.Ripemd160(inputBytes)
	if err != nil {
		arwen.WithFaultAndHostIfFailAlwaysActive(err, host, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.SetBytes(outputHandle, result)

	return 0
}

//export v1_4_verifyBLS
func v1_4_verifyBLS(
	context unsafe.Pointer,
	keyOffset int32,
	messageOffset int32,
	messageLength int32,
	sigOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)
	metering.StartGasTracing(verifyBLSName)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifyBLS
	metering.UseAndTraceGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, blsPublicKeyLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(messageLength))
	metering.UseAndTraceGas(gasToUse)

	message, err := runtime.MemLoad(messageOffset, messageLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	sig, err := runtime.MemLoad(sigOffset, blsSignatureLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	invalidSigErr := crypto.VerifyBLS(key, message, sig)
	if invalidSigErr != nil {
		arwen.WithFaultIfFailAlwaysActive(invalidSigErr, context, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	return 0
}

//export v1_4_managedVerifyBLS
func v1_4_managedVerifyBLS(
	context unsafe.Pointer,
	keyHandle int32,
	messageHandle int32,
	sigHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	return ManagedVerifyBLSWithHost(host, keyHandle, messageHandle, sigHandle)
}

func ManagedVerifyBLSWithHost(
	host arwen.VMHost,
	keyHandle int32,
	messageHandle int32,
	sigHandle int32,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	crypto := host.Crypto()
	metering.StartGasTracing(verifyBLSName)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifyBLS
	metering.UseAndTraceGas(gasToUse)

	keyBytes, err := managedType.GetBytes(keyHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(keyBytes)

	msgBytes, err := managedType.GetBytes(messageHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(msgBytes)

	sigBytes, err := managedType.GetBytes(sigHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(sigBytes)

	invalidSigErr := crypto.VerifyBLS(keyBytes, msgBytes, sigBytes)
	if invalidSigErr != nil {
		arwen.WithFaultAndHostIfFailAlwaysActive(invalidSigErr, host, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	return 0
}

//export v1_4_verifyEd25519
func v1_4_verifyEd25519(
	context unsafe.Pointer,
	keyOffset int32,
	messageOffset int32,
	messageLength int32,
	sigOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)
	metering.StartGasTracing(verifyEd25519Name)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifyEd25519
	metering.UseAndTraceGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, ed25519PublicKeyLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(messageLength))
	metering.UseAndTraceGas(gasToUse)

	message, err := runtime.MemLoad(messageOffset, messageLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	sig, err := runtime.MemLoad(sigOffset, ed25519SignatureLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	invalidSigErr := crypto.VerifyEd25519(key, message, sig)
	if invalidSigErr != nil {
		arwen.WithFaultIfFailAlwaysActive(invalidSigErr, context, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	return 0
}

//export v1_4_managedVerifyEd25519
func v1_4_managedVerifyEd25519(
	context unsafe.Pointer,
	keyHandle, messageHandle, sigHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	return ManagedVerifyEd25519WithHost(host, keyHandle, messageHandle, sigHandle)
}

func ManagedVerifyEd25519WithHost(
	host arwen.VMHost,
	keyHandle, messageHandle, sigHandle int32,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	crypto := host.Crypto()
	metering.StartGasTracing(verifyEd25519Name)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifyEd25519
	metering.UseAndTraceGas(gasToUse)

	keyBytes, err := managedType.GetBytes(keyHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(keyBytes)

	msgBytes, err := managedType.GetBytes(messageHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(msgBytes)

	sigBytes, err := managedType.GetBytes(sigHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(sigBytes)

	invalidSigErr := crypto.VerifyEd25519(keyBytes, msgBytes, sigBytes)
	if invalidSigErr != nil {
		arwen.WithFaultAndHostIfFailAlwaysActive(invalidSigErr, host, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	return 0
}

//export v1_4_verifyCustomSecp256k1
func v1_4_verifyCustomSecp256k1(
	context unsafe.Pointer,
	keyOffset int32,
	keyLength int32,
	messageOffset int32,
	messageLength int32,
	sigOffset int32,
	hashType int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)
	metering.StartGasTracing(verifyCustomSecp256k1Name)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifySecp256k1
	metering.UseAndTraceGas(gasToUse)

	if keyLength != secp256k1CompressedPublicKeyLength && keyLength != secp256k1UncompressedPublicKeyLength {
		_ = arwen.WithFault(arwen.ErrInvalidPublicKeySize, context, runtime.ElrondAPIErrorShouldFailExecution())
		return 1
	}

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(messageLength))
	metering.UseAndTraceGas(gasToUse)

	message, err := runtime.MemLoad(messageOffset, messageLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	// read the 2 leading bytes first
	// byte1: 0x30, header
	// byte2: the remaining buffer length
	const sigHeaderLength = 2
	sigHeader, err := runtime.MemLoad(sigOffset, sigHeaderLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}
	sigLength := int32(sigHeader[1]) + sigHeaderLength
	sig, err := runtime.MemLoad(sigOffset, sigLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	invalidSigErr := crypto.VerifySecp256k1(key, message, sig, uint8(hashType))
	if invalidSigErr != nil {
		arwen.WithFaultIfFailAlwaysActive(invalidSigErr, context, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	return 0
}

//export v1_4_managedVerifyCustomSecp256k1
func v1_4_managedVerifyCustomSecp256k1(
	context unsafe.Pointer,
	keyHandle, messageHandle, sigHandle int32,
	hashType int32,
) int32 {
	host := arwen.GetVMHost(context)
	return ManagedVerifyCustomSecp256k1WithHost(
		host,
		keyHandle,
		messageHandle,
		sigHandle,
		hashType)
}

func ManagedVerifyCustomSecp256k1WithHost(
	host arwen.VMHost,
	keyHandle, messageHandle, sigHandle int32,
	hashType int32,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	crypto := host.Crypto()
	metering.StartGasTracing(verifyCustomSecp256k1Name)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifySecp256k1
	metering.UseAndTraceGas(gasToUse)

	keyBytes, err := managedType.GetBytes(keyHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(keyBytes)

	msgBytes, err := managedType.GetBytes(messageHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(msgBytes)

	sigBytes, err := managedType.GetBytes(sigHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(sigBytes)

	invalidSigErr := crypto.VerifySecp256k1(keyBytes, msgBytes, sigBytes, uint8(hashType))
	if invalidSigErr != nil {
		arwen.WithFaultAndHostIfFailAlwaysActive(invalidSigErr, host, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	return 0
}

//export v1_4_verifySecp256k1
func v1_4_verifySecp256k1(
	context unsafe.Pointer,
	keyOffset int32,
	keyLength int32,
	messageOffset int32,
	messageLength int32,
	sigOffset int32,
) int32 {
	return v1_4_verifyCustomSecp256k1(
		context,
		keyOffset,
		keyLength,
		messageOffset,
		messageLength,
		sigOffset,
		int32(secp256k1.ECDSADoubleSha256),
	)
}

//export v1_4_managedVerifySecp256k1
func v1_4_managedVerifySecp256k1(
	context unsafe.Pointer,
	keyHandle, messageHandle, sigHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	return ManagedVerifySecp256k1WithHost(host, keyHandle, messageHandle, sigHandle)
}

func ManagedVerifySecp256k1WithHost(
	host arwen.VMHost,
	keyHandle, messageHandle, sigHandle int32,
) int32 {
	return ManagedVerifyCustomSecp256k1WithHost(
		host,
		keyHandle,
		messageHandle,
		sigHandle,
		int32(secp256k1.ECDSADoubleSha256),
	)
}

//export v1_4_encodeSecp256k1DerSignature
func v1_4_encodeSecp256k1DerSignature(
	context unsafe.Pointer,
	rOffset int32,
	rLength int32,
	sOffset int32,
	sLength int32,
	sigOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.EncodeDERSig
	metering.UseGasAndAddTracedGas(encodeSecp256k1DerSignatureName, gasToUse)

	r, err := runtime.MemLoad(rOffset, rLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	s, err := runtime.MemLoad(sOffset, sLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	derSig := crypto.EncodeSecp256k1DERSignature(r, s)
	err = runtime.MemStore(sigOffset, derSig)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export v1_4_managedEncodeSecp256k1DerSignature
func v1_4_managedEncodeSecp256k1DerSignature(
	context unsafe.Pointer,
	rHandle, sHandle, sigHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	return ManagedEncodeSecp256k1DerSignatureWithHost(host, rHandle, sHandle, sigHandle)
}

func ManagedEncodeSecp256k1DerSignatureWithHost(
	host arwen.VMHost,
	rHandle, sHandle, sigHandle int32,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	crypto := host.Crypto()

	gasToUse := metering.GasSchedule().CryptoAPICost.EncodeDERSig
	metering.UseGasAndAddTracedGas(encodeSecp256k1DerSignatureName, gasToUse)

	r, err := managedType.GetBytes(rHandle)
	if arwen.WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	s, err := managedType.GetBytes(sHandle)
	if arwen.WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	derSig := crypto.EncodeSecp256k1DERSignature(r, s)
	managedType.SetBytes(sigHandle, derSig)

	return 0
}

//export v1_4_addEC
func v1_4_addEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	fstPointXHandle int32,
	fstPointYHandle int32,
	sndPointXHandle int32,
	sndPointYHandle int32,
) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(addECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.AddECC * uint64(curveMultiplier) / 100
	metering.UseAndTraceGas(gasToUse)

	ec, err1 := managedType.GetEllipticCurve(ecHandle)
	if arwen.WithFault(err1, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return
	}

	xResult, yResult, err := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	if err != nil {
		_ = arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	x1, y1, err := managedType.GetTwoBigInt(fstPointXHandle, fstPointYHandle)
	if err != nil {
		_ = arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	x2, y2, err := managedType.GetTwoBigInt(sndPointXHandle, sndPointYHandle)
	if err != nil {
		_ = arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}

	if !ec.IsOnCurve(x1, y1) || !ec.IsOnCurve(x2, y2) {
		_ = arwen.WithFault(arwen.ErrPointNotOnCurve, context, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}

	managedType.ConsumeGasForBigIntCopy(xResult, yResult, ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x1, y1, x2, y2)
	xResultAdd, yResultAdd := ec.Add(x1, y1, x2, y2)
	xResult.Set(xResultAdd)
	yResult.Set(yResultAdd)
}

//export v1_4_doubleEC
func v1_4_doubleEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	pointXHandle int32,
	pointYHandle int32,
) {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(doubleECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.DoubleECC * uint64(curveMultiplier) / 100
	metering.UseAndTraceGas(gasToUse)

	ec, err1 := managedType.GetEllipticCurve(ecHandle)
	if arwen.WithFault(err1, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return
	}

	xResult, yResult, err1 := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	x, y, err2 := managedType.GetTwoBigInt(pointXHandle, pointYHandle)
	if err1 != nil || err2 != nil {
		_ = arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}
	if !ec.IsOnCurve(x, y) {
		_ = arwen.WithFault(arwen.ErrPointNotOnCurve, context, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}

	managedType.ConsumeGasForBigIntCopy(xResult, yResult, ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)

	xResultDouble, yResultDouble := ec.Double(x, y)
	xResult.Set(xResultDouble)
	yResult.Set(yResultDouble)
}

//export v1_4_isOnCurveEC
func v1_4_isOnCurveEC(
	context unsafe.Pointer,
	ecHandle int32,
	pointXHandle int32,
	pointYHandle int32,
) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering.StartGasTracing(isOnCurveECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.IsOnCurveECC * uint64(curveMultiplier) / 100
	metering.UseAndTraceGas(gasToUse)

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return -1
	}

	x, y, err := managedType.GetTwoBigInt(pointXHandle, pointYHandle)
	if err != nil || x == nil || y == nil {
		_ = arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	if ec.IsOnCurve(x, y) {
		return 1
	}

	return 0
}

//export v1_4_scalarBaseMultEC
func v1_4_scalarBaseMultEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataOffset int32,
	length int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)
	metering.StartGasTracing(scalarBaseMultECName)

	if length < 0 {
		_ = arwen.WithFault(arwen.ErrNegativeLength, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	curveMultiplier := managedType.GetScalarMult100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	oneByteScalarGasCost := metering.GasSchedule().CryptoAPICost.ScalarMultECC * uint64(curveMultiplier) / 100
	gasToUse := oneByteScalarGasCost + uint64(length)*oneByteScalarGasCost
	metering.UseAndTraceGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	host := arwen.GetVMHost(context)
	return commonScalarBaseMultEC(host, xResultHandle, yResultHandle, ecHandle, data)
}

//export v1_4_managedScalarBaseMultEC
func v1_4_managedScalarBaseMultEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	return ManagedScalarBaseMultECWithHost(
		host,
		xResultHandle,
		yResultHandle,
		ecHandle,
		dataHandle,
	)
}

func ManagedScalarBaseMultECWithHost(
	host arwen.VMHost,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataHandle int32,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(scalarBaseMultECName)

	curveMultiplier := managedType.GetScalarMult100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = arwen.WithFaultAndHost(host, arwen.ErrNoEllipticCurveUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	data, err := managedType.GetBytes(dataHandle)
	if arwen.WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	oneByteScalarGasCost := metering.GasSchedule().CryptoAPICost.ScalarMultECC * uint64(curveMultiplier) / 100
	gasToUse := oneByteScalarGasCost + uint64(len(data))*oneByteScalarGasCost
	metering.UseAndTraceGas(gasToUse)

	return commonScalarBaseMultEC(host, xResultHandle, yResultHandle, ecHandle, data)
}

func commonScalarBaseMultEC(
	host arwen.VMHost,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	data []byte,
) int32 {
	runtime := host.Runtime()
	managedType := host.ManagedTypes()

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if arwen.WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	xResult, yResult, err := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	if arwen.WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xResult, yResult)

	xResultSBM, yResultSBM := ec.ScalarBaseMult(data)
	if !ec.IsOnCurve(xResultSBM, yResultSBM) {
		_ = arwen.WithFaultAndHost(host, arwen.ErrPointNotOnCurve, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	xResult.Set(xResultSBM)
	yResult.Set(yResultSBM)

	return 0
}

//export v1_4_scalarMultEC
func v1_4_scalarMultEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	pointXHandle int32,
	pointYHandle int32,
	dataOffset int32,
	length int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)
	metering.StartGasTracing(scalarMultECName)

	if length < 0 {
		_ = arwen.WithFault(arwen.ErrNegativeLength, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	curveMultiplier := managedType.GetScalarMult100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	oneByteScalarGasCost := metering.GasSchedule().CryptoAPICost.ScalarMultECC * uint64(curveMultiplier) / 100
	gasToUse := oneByteScalarGasCost + uint64(length)*oneByteScalarGasCost
	metering.UseAndTraceGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	host := arwen.GetVMHost(context)
	return commonScalarMultEC(host, xResultHandle, yResultHandle, ecHandle, pointXHandle, pointYHandle, data)
}

//export v1_4_managedScalarMultEC
func v1_4_managedScalarMultEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	pointXHandle int32,
	pointYHandle int32,
	dataHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	return ManagedScalarMultECWithHost(
		host,
		xResultHandle,
		yResultHandle,
		ecHandle,
		pointXHandle,
		pointYHandle,
		dataHandle,
	)
}

func ManagedScalarMultECWithHost(
	host arwen.VMHost,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	pointXHandle int32,
	pointYHandle int32,
	dataHandle int32,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(scalarMultECName)

	curveMultiplier := managedType.GetScalarMult100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = arwen.WithFaultAndHost(host, arwen.ErrNoEllipticCurveUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	data, err := managedType.GetBytes(dataHandle)
	if arwen.WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	oneByteScalarGasCost := metering.GasSchedule().CryptoAPICost.ScalarMultECC * uint64(curveMultiplier) / 100
	gasToUse := oneByteScalarGasCost + uint64(len(data))*oneByteScalarGasCost
	metering.UseAndTraceGas(gasToUse)

	return commonScalarMultEC(host, xResultHandle, yResultHandle, ecHandle, pointXHandle, pointYHandle, data)
}

func commonScalarMultEC(
	host arwen.VMHost,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	pointXHandle int32,
	pointYHandle int32,
	data []byte,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(scalarMultECName)

	ec, err1 := managedType.GetEllipticCurve(ecHandle)
	if arwen.WithFaultAndHost(host, err1, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	xResult, yResult, err1 := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	x, y, err2 := managedType.GetTwoBigInt(pointXHandle, pointYHandle)
	if err1 != nil || err2 != nil {
		_ = arwen.WithFaultAndHost(host, arwen.ErrNoBigIntUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	if !ec.IsOnCurve(x, y) {
		_ = arwen.WithFaultAndHost(host, arwen.ErrPointNotOnCurve, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(xResult, yResult, ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	xResultSM, yResultSM := ec.ScalarMult(x, y, data)
	if !ec.IsOnCurve(xResultSM, yResultSM) {
		_ = arwen.WithFaultAndHost(host, arwen.ErrPointNotOnCurve, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	xResult.Set(xResultSM)
	yResult.Set(yResultSM)

	return 0
}

//export v1_4_marshalEC
func v1_4_marshalEC(
	context unsafe.Pointer,
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	resultOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	host := arwen.GetVMHost(context)
	result, err := commonMarshalEC(host, xPairHandle, yPairHandle, ecHandle)
	if err != nil {
		_ = arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return -1
	}
	return int32(len(result))
}

//export v1_4_managedMarshalEC
func v1_4_managedMarshalEC(
	context unsafe.Pointer,
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	resultHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	return ManagedMarshalECWithHost(
		host,
		xPairHandle,
		yPairHandle,
		ecHandle,
		resultHandle,
	)
}

func ManagedMarshalECWithHost(
	host arwen.VMHost,
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	resultHandle int32,
) int32 {
	result, err := commonMarshalEC(host, xPairHandle, yPairHandle, ecHandle)
	if err != nil {
		_ = arwen.WithFaultAndHost(host, err, true)
		return -1
	}

	managedType := host.ManagedTypes()
	managedType.SetBytes(resultHandle, result)
	return int32(len(result))
}

func commonMarshalEC(
	host arwen.VMHost,
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
) ([]byte, error) {
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(marshalECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		return nil, arwen.ErrNoEllipticCurveUnderThisHandle
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.MarshalECC * uint64(curveMultiplier) / 100
	metering.UseAndTraceGas(gasToUse)

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		return nil, err
	}

	x, y, err := managedType.GetTwoBigInt(xPairHandle, yPairHandle)
	if err != nil {
		return nil, err
	}
	if !ec.IsOnCurve(x, y) {
		return nil, arwen.ErrPointNotOnCurve
	}
	if x.BitLen() > int(ec.BitSize) || y.BitLen() > int(ec.BitSize) {
		return nil, arwen.ErrLengthOfBufferNotCorrect
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)

	result := elliptic.Marshal(ec, x, y)
	return result, nil
}

//export v1_4_marshalCompressedEC
func v1_4_marshalCompressedEC(
	context unsafe.Pointer,
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	resultOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	host := arwen.GetVMHost(context)
	result, err := commonMarshalCompressedEC(host, xPairHandle, yPairHandle, ecHandle)
	if err != nil {
		_ = arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return -1
	}
	return int32(len(result))
}

//export v1_4_managedMarshalCompressedEC
func v1_4_managedMarshalCompressedEC(
	context unsafe.Pointer,
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	resultHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	return ManagedMarshalCompressedECWithHost(
		host,
		xPairHandle,
		yPairHandle,
		ecHandle,
		resultHandle,
	)
}

func ManagedMarshalCompressedECWithHost(
	host arwen.VMHost,
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	resultHandle int32,
) int32 {
	runtime := host.Runtime()
	managedType := host.ManagedTypes()
	result, err := commonMarshalCompressedEC(host, xPairHandle, yPairHandle, ecHandle)
	if err != nil {
		_ = arwen.WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	managedType.SetBytes(resultHandle, result)
	return int32(len(result))
}

func commonMarshalCompressedEC(host arwen.VMHost,
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
) ([]byte, error) {
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(marshalCompressedECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		return nil, arwen.ErrNoEllipticCurveUnderThisHandle
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.MarshalCompressedECC * uint64(curveMultiplier) / 100
	metering.UseAndTraceGas(gasToUse)

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		return nil, err
	}

	x, y, err := managedType.GetTwoBigInt(xPairHandle, yPairHandle)
	if err != nil || x == nil || y == nil {
		return nil, arwen.ErrNoBigIntUnderThisHandle
	}
	if !ec.IsOnCurve(x, y) {
		return nil, arwen.ErrPointNotOnCurve
	}
	if x.BitLen() > int(ec.BitSize) || y.BitLen() > int(ec.BitSize) {
		return nil, arwen.ErrLengthOfBufferNotCorrect
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)

	result := elliptic.MarshalCompressed(ec, x, y)
	return result, nil
}

//export v1_4_unmarshalEC
func v1_4_unmarshalEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataOffset int32,
	length int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)
	metering.StartGasTracing(unmarshalECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.UnmarshalECC * uint64(curveMultiplier) / 100
	metering.UseAndTraceGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	host := arwen.GetVMHost(context)
	return commonUnmarshalEC(host, xResultHandle, yResultHandle, ecHandle, data)
}

//export v1_4_managedUnmarshalEC
func v1_4_managedUnmarshalEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	return ManagedUnmarshalECWithHost(
		host,
		xResultHandle,
		yResultHandle,
		ecHandle,
		dataHandle,
	)
}

func ManagedUnmarshalECWithHost(
	host arwen.VMHost,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataHandle int32,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(unmarshalECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = arwen.WithFaultAndHost(host, arwen.ErrNoEllipticCurveUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.UnmarshalECC * uint64(curveMultiplier) / 100
	metering.UseAndTraceGas(gasToUse)

	data, err := managedType.GetBytes(dataHandle)
	if err != nil {
		return 1
	}

	return commonUnmarshalEC(host, xResultHandle, yResultHandle, ecHandle, data)
}

func commonUnmarshalEC(
	host arwen.VMHost,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	data []byte,
) int32 {
	runtime := host.Runtime()
	managedType := host.ManagedTypes()

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if arwen.WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}
	byteLen := (ec.BitSize + 7) / 8
	if len(data) != 1+2*byteLen {
		_ = arwen.WithFaultAndHost(host, arwen.ErrLengthOfBufferNotCorrect, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	xResult, yResult, err := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	if arwen.WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xResult, yResult)

	xResultU, yResultU := elliptic.Unmarshal(ec, data)
	if xResultU == nil || yResultU == nil || !ec.IsOnCurve(xResultU, yResultU) {
		_ = arwen.WithFaultAndHost(host, arwen.ErrPointNotOnCurve, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	xResult.Set(xResultU)
	yResult.Set(yResultU)

	return 0
}

//export v1_4_unmarshalCompressedEC
func v1_4_unmarshalCompressedEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataOffset int32,
	length int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)
	metering.StartGasTracing(unmarshalCompressedECName)

	curveMultiplier := managedType.GetUCompressed100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.UnmarshalCompressedECC * uint64(curveMultiplier) / 100
	metering.UseAndTraceGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(data))
	}

	host := arwen.GetVMHost(context)
	return commonUnmarshalCompressedEC(host, xResultHandle, yResultHandle, ecHandle, data)
}

//export v1_4_managedUnmarshalCompressedEC
func v1_4_managedUnmarshalCompressedEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	return ManagedUnmarshalCompressedECWithHost(
		host,
		xResultHandle,
		yResultHandle,
		ecHandle,
		dataHandle,
	)
}

func ManagedUnmarshalCompressedECWithHost(
	host arwen.VMHost,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataHandle int32,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(unmarshalCompressedECName)

	curveMultiplier := managedType.GetUCompressed100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = arwen.WithFaultAndHost(host, arwen.ErrNoEllipticCurveUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.UnmarshalCompressedECC * uint64(curveMultiplier) / 100
	metering.UseAndTraceGas(gasToUse)

	data, err := managedType.GetBytes(dataHandle)
	if arwen.WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(data))
	}

	return commonUnmarshalCompressedEC(host, xResultHandle, yResultHandle, ecHandle, data)
}

func commonUnmarshalCompressedEC(
	host arwen.VMHost,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	data []byte,
) int32 {
	runtime := host.Runtime()
	managedType := host.ManagedTypes()

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if arwen.WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}
	byteLen := (ec.BitSize+7)/8 + 1
	if len(data) != byteLen {
		_ = arwen.WithFaultAndHost(host, arwen.ErrLengthOfBufferNotCorrect, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	xResult, yResult, err := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	if arwen.WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xResult, yResult)

	xResultUC, yResultUC := elliptic.UnmarshalCompressed(ec, data)
	if xResultUC == nil || yResultUC == nil || !ec.IsOnCurve(xResultUC, yResultUC) {
		_ = arwen.WithFaultAndHost(host, arwen.ErrPointNotOnCurve, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	xResult.Set(xResultUC)
	yResult.Set(yResultUC)
	return 0
}

//export v1_4_generateKeyEC
func v1_4_generateKeyEC(
	context unsafe.Pointer,
	xPubKeyHandle int32,
	yPubKeyHandle int32,
	ecHandle int32,
	resultOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	host := arwen.GetVMHost(context)
	result, err := commonGenerateEC(host, xPubKeyHandle, yPubKeyHandle, ecHandle)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(result))
	}

	return 0
}

//export v1_4_managedGenerateKeyEC
func v1_4_managedGenerateKeyEC(
	context unsafe.Pointer,
	xPubKeyHandle int32,
	yPubKeyHandle int32,
	ecHandle int32,
	resultHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	return ManagedGenerateKeyECWithHost(
		host,
		xPubKeyHandle,
		yPubKeyHandle,
		ecHandle,
		resultHandle,
	)
}

func ManagedGenerateKeyECWithHost(
	host arwen.VMHost,
	xPubKeyHandle int32,
	yPubKeyHandle int32,
	ecHandle int32,
	resultHandle int32,
) int32 {
	runtime := host.Runtime()
	managedType := host.ManagedTypes()
	result, err := commonGenerateEC(host, xPubKeyHandle, yPubKeyHandle, ecHandle)
	if arwen.WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.SetBytes(resultHandle, result)
	return 0
}

func commonGenerateEC(
	host arwen.VMHost,
	xPubKeyHandle int32,
	yPubKeyHandle int32,
	ecHandle int32,
) ([]byte, error) {
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(generateKeyECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		return nil, arwen.ErrNoEllipticCurveUnderThisHandle
	}
	if curveMultiplier == 250 {
		curveMultiplier = 500
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.GenerateKeyECC * uint64(curveMultiplier) / 100
	metering.UseAndTraceGas(gasToUse)

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		return nil, err
	}

	xPubKey, yPubKey, err := managedType.GetTwoBigInt(xPubKeyHandle, yPubKeyHandle)
	if err != nil {
		return nil, err
	}
	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xPubKey, yPubKey)

	ioReader := managedType.GetRandReader()
	result, xPubKeyGK, yPubKeyGK, err := elliptic.GenerateKey(ec, ioReader)
	if err != nil {
		return nil, err
	}

	xPubKey.Set(xPubKeyGK)
	yPubKey.Set(yPubKeyGK)

	return result, nil
}

//export v1_4_createEC
func v1_4_createEC(context unsafe.Pointer, dataOffset int32, dataLength int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.EllipticCurveNew
	metering.UseGasAndAddTracedGas(createECName, gasToUse)

	if dataLength != curveNameLength {
		_ = arwen.WithFault(arwen.ErrBadBounds, context, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}
	data, err := runtime.MemLoad(dataOffset, dataLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return -1
	}
	curveChoice := string(data[:])
	switch curveChoice {
	case "p224":
		curveParams := elliptic.P224().Params()
		return managedType.PutEllipticCurve(curveParams)
	case "p256":
		curveParams := elliptic.P256().Params()
		return managedType.PutEllipticCurve(curveParams)
	case "p384":
		curveParams := elliptic.P384().Params()
		return managedType.PutEllipticCurve(curveParams)
	case "p521":
		curveParams := elliptic.P521().Params()
		return managedType.PutEllipticCurve(curveParams)
	}
	return -1
}

//export v1_4_managedCreateEC
func v1_4_managedCreateEC(context unsafe.Pointer, dataHandle int32) int32 {
	host := arwen.GetVMHost(context)
	return ManagedCreateECWithHost(host, dataHandle)
}

func ManagedCreateECWithHost(host arwen.VMHost, dataHandle int32) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().CryptoAPICost.EllipticCurveNew
	metering.UseGasAndAddTracedGas(createECName, gasToUse)

	data, err := managedType.GetBytes(dataHandle)
	if arwen.WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return -1
	}
	curveChoice := string(data[:])
	switch curveChoice {
	case "p224":
		curveParams := elliptic.P224().Params()
		return managedType.PutEllipticCurve(curveParams)
	case "p256":
		curveParams := elliptic.P256().Params()
		return managedType.PutEllipticCurve(curveParams)
	case "p384":
		curveParams := elliptic.P384().Params()
		return managedType.PutEllipticCurve(curveParams)
	case "p521":
		curveParams := elliptic.P521().Params()
		return managedType.PutEllipticCurve(curveParams)
	}

	_ = arwen.WithFaultAndHost(host, arwen.ErrBadBounds, runtime.CryptoAPIErrorShouldFailExecution())
	return -1
}

//export v1_4_getCurveLengthEC
func v1_4_getCurveLengthEC(context unsafe.Pointer, ecHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetInt64
	metering.UseGasAndAddTracedGas(getCurveLengthECName, gasToUse)

	ecLength := managedType.GetEllipticCurveSizeOfField(ecHandle)
	if ecLength == -1 {
		_ = arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.BigIntAPIErrorShouldFailExecution())
	}

	return ecLength
}

//export v1_4_getPrivKeyByteLengthEC
func v1_4_getPrivKeyByteLengthEC(context unsafe.Pointer, ecHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetInt64
	metering.UseGasAndAddTracedGas(getPrivKeyByteLengthECName, gasToUse)

	byteLength := managedType.GetPrivateKeyByteLengthEC(ecHandle)
	if byteLength == -1 {
		_ = arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.BigIntAPIErrorShouldFailExecution())
	}

	return byteLength
}

//export v1_4_ellipticCurveGetValues
func v1_4_ellipticCurveGetValues(context unsafe.Pointer, ecHandle int32, fieldOrderHandle int32, basePointOrderHandle int32, eqConstantHandle int32, xBasePointHandle int32, yBasePointHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetInt64 * 5
	metering.UseGasAndAddTracedGas(ellipticCurveGetValuesName, gasToUse)

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return -1
	}
	fieldOrder, basePointOrder, err := managedType.GetTwoBigInt(fieldOrderHandle, basePointOrderHandle)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return -1
	}
	eqConstant, err := managedType.GetBigInt(eqConstantHandle)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return -1
	}
	xBasePoint, yBasePoint, err := managedType.GetTwoBigInt(xBasePointHandle, yBasePointHandle)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return -1
	}
	fieldOrder.Set(ec.P)
	basePointOrder.Set(ec.N)
	eqConstant.Set(ec.B)
	xBasePoint.Set(ec.Gx)
	yBasePoint.Set(ec.Gy)
	return ecHandle
}

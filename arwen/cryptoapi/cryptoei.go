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
// extern int32_t verifyBLS(void *context, int32_t keyOffset, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern int32_t verifyEd25519(void *context, int32_t keyOffset, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern int32_t verifySecp256k1(void *context, int32_t keyOffset, int32_t keyLength, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern void addEC(void *context, int32_t destination1, int32_t destination2, int32_t fieldOrder, int32_t basePointOrder, int32_t eqConstant, int32_t xBasePoint, int32_t yBasePoint, int32_t sizeOfField, int32_t fstPointX, int32_t fstPointY, int32_t sndPointX, int32_t sndPointY);
// extern void doubleEC(void *context, int32_t destination1, int32_t destination2, int32_t fieldOrder, int32_t basePointOrder, int32_t eqConstant, int32_t xBasePoint, int32_t yBasePoint, int32_t sizeOfField, int32_t pointX, int32_t pointY);
// extern int32_t isOnCurveEC(void *context, int32_t fieldOrder, int32_t basePointOrder, int32_t eqConstant, int32_t xBasePoint, int32_t yBasePoint, int32_t sizeOfField, int32_t pointX, int32_t pointY);
// extern int32_t scalarBaseMultEC(void *context, int32_t destination1, int32_t destination2, int32_t fieldOrder, int32_t basePointOrder, int32_t eqConstant, int32_t xBasePoint, int32_t yBasePoint, int32_t sizeOfField, int32_t kOffset, int32_t length);
// extern int32_t scalarMultEC(void *context, int32_t destination1, int32_t destination2, int32_t fieldOrder, int32_t basePointOrder, int32_t eqConstant, int32_t xBasePoint, int32_t yBasePoint, int32_t sizeOfField, int32_t pointX, int32_t pointY, int32_t kOffset, int32_t length);
// extern int32_t marshalEC(void *context, int32_t fieldOrder, int32_t basePointOrder, int32_t eqConstant, int32_t xBasePoint, int32_t yBasePoint, int32_t sizeOfField, int32_t xPairHandle, int32_t yPairHandle, int32_t resultOffest);
// extern int32_t unmarshalEC(void *context, int32_t xPairHandle, int32_t yPairHandle, int32_t fieldOrder, int32_t basePointOrder, int32_t eqConstant, int32_t xBasePoint, int32_t yBasePoint, int32_t sizeOfField, int32_t dataOffest, int32_t length);
// extern int32_t marshalCompressedEC(void *context, int32_t fieldOrder, int32_t basePointOrder, int32_t eqConstant, int32_t xBasePoint, int32_t yBasePoint, int32_t sizeOfField, int32_t xPairHandle, int32_t yPairHandle, int32_t resultOffest);
// extern int32_t unmarshalCompressedEC(void *context, int32_t xPairHandle, int32_t yPairHandle, int32_t fieldOrder, int32_t basePointOrder, int32_t eqConstant, int32_t xBasePoint, int32_t yBasePoint, int32_t sizeOfField, int32_t dataOffest, int32_t length);
// extern int32_t generateKeyEC(void *context, int32_t xPubKeyHandle, int32_t yPubKeyHandle, int32_t fieldOrder, int32_t basePointOrder, int32_t eqConstant, int32_t xBasePoint, int32_t yBasePoint, int32_t sizeOfField, int32_t resultOffset);
import "C"

import (
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	elliptic_curve "github.com/ElrondNetwork/arwen-wasm-vm/crypto/elliptic_curves"
	"github.com/ElrondNetwork/arwen-wasm-vm/math"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
)

const blsPublicKeyLength = 96
const blsSignatureLength = 48
const ed25519PublicKeyLength = 32
const ed25519SignatureLength = 64
const secp256k1CompressedPublicKeyLength = 33
const secp256k1UncompressedPublicKeyLength = 65
const secp256k1SignatureLength = 64

// CryptoImports adds some crypto imports to the Wasmer Imports map
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

	imports, err = imports.Append("ripemd160", ripemd160, C.ripemd160)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("verifyBLS", verifyBLS, C.verifyBLS)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("verifyEd25519", verifyEd25519, C.verifyEd25519)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("verifySecp256k1", verifySecp256k1, C.verifySecp256k1)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("addEC", addEC, C.addEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("doubleEC", doubleEC, C.doubleEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("isOnCurveEC", isOnCurveEC, C.isOnCurveEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("scalarBaseMultEC", scalarBaseMultEC, C.scalarBaseMultEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("scalarMultEC", scalarMultEC, C.scalarMultEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("marshalEC", marshalEC, C.marshalEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("unmarshalEC", unmarshalEC, C.unmarshalEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("marshalCompressedEC", marshalCompressedEC, C.marshalCompressedEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("unmarshalCompressedEC", unmarshalCompressedEC, C.unmarshalCompressedEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("generateKeyEC", generateKeyEC, C.generateKeyEC)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export sha256
func sha256(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	memLoadGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(metering.GasSchedule().CryptoAPICost.SHA256, memLoadGas)
	metering.UseGas(gasToUse)

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

	return 0
}

//export keccak256
func keccak256(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	memLoadGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(metering.GasSchedule().CryptoAPICost.Keccak256, memLoadGas)
	metering.UseGas(gasToUse)

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

	return 0
}

//export ripemd160
func ripemd160(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	memLoadGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(metering.GasSchedule().CryptoAPICost.Ripemd160, memLoadGas)
	metering.UseGas(gasToUse)

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

	return 0
}

//export verifyBLS
func verifyBLS(
	context unsafe.Pointer,
	keyOffset int32,
	messageOffset int32,
	messageLength int32,
	sigOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifyBLS
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, blsPublicKeyLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(messageLength))
	metering.UseGas(gasToUse)

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
		return -1
	}

	return 0
}

//export verifyEd25519
func verifyEd25519(
	context unsafe.Pointer,
	keyOffset int32,
	messageOffset int32,
	messageLength int32,
	sigOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifyEd25519
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, ed25519PublicKeyLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(messageLength))
	metering.UseGas(gasToUse)

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
		return -1
	}

	return 0
}

//export verifySecp256k1
func verifySecp256k1(
	context unsafe.Pointer,
	keyOffset int32,
	keyLength int32,
	messageOffset int32,
	messageLength int32,
	sigOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	crypto := arwen.GetCryptoContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifySecp256k1
	metering.UseGas(gasToUse)

	if keyLength != secp256k1CompressedPublicKeyLength && keyLength != secp256k1UncompressedPublicKeyLength {
		arwen.WithFault(arwen.ErrInvalidPublicKeySize, context, runtime.ElrondAPIErrorShouldFailExecution())
		return 1
	}

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(messageLength))
	metering.UseGas(gasToUse)

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

	invalidSigErr := crypto.VerifySecp256k1(key, message, sig)
	if invalidSigErr != nil {
		return -1
	}

	return 0
}

//export addEC
func addEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	fieldOrderHandle int32,
	basePointOrderHandle int32,
	eqConstantHandle int32,
	xBasePointHandle int32,
	yBasePointHandle int32,
	sizeOfField int32,
	fstPointXHandle int32,
	fstPointYHandle int32,
	sndPointXHandle int32,
	sndPointYHandle int32,
) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	xResult, yResult, P := bigInt.GetThree(xResultHandle, yResultHandle, fieldOrderHandle)
	N, B, Gx := bigInt.GetThree(basePointOrderHandle, eqConstantHandle, xBasePointHandle)
	Gy, x1, y1 := bigInt.GetThree(yBasePointHandle, fstPointXHandle, fstPointYHandle)
	x2, y2 := bigInt.GetTwo(sndPointXHandle, sndPointYHandle)
	bigInt.ConsumeGasForBigIntCopy(xResult, yResult, P, N, B, Gx, Gy, x1, y1, x2, y2)

	xResultAdd, yResultAdd := elliptic_curve.Add(P, N, B, Gx, Gy, int(sizeOfField), x1, y1, x2, y2)
	xResult.Set(xResultAdd)
	yResult.Set(yResultAdd)
}

//export doubleEC
func doubleEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	fieldOrderHandle int32,
	basePointOrderHandle int32,
	eqConstantHandle int32,
	xBasePointHandle int32,
	yBasePointHandle int32,
	sizeOfField int32,
	pointXHandle int32,
	pointYHandle int32,
) {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	xResult, yResult, P := bigInt.GetThree(xResultHandle, yResultHandle, fieldOrderHandle)
	N, B, Gx := bigInt.GetThree(basePointOrderHandle, eqConstantHandle, xBasePointHandle)
	Gy, x, y := bigInt.GetThree(yBasePointHandle, pointXHandle, pointYHandle)
	bigInt.ConsumeGasForBigIntCopy(xResult, yResult, P, N, B, Gx, Gy, x, y)

	xResultDouble, yResultDouble := elliptic_curve.Double(P, N, B, Gx, Gy, int(sizeOfField), x, y)
	xResult.Set(xResultDouble)
	yResult.Set(yResultDouble)
}

//export isOnCurveEC
func isOnCurveEC(
	context unsafe.Pointer,
	fieldOrderHandle int32,
	basePointOrderHandle int32,
	eqConstantHandle int32,
	xBasePointHandle int32,
	yBasePointHandle int32,
	sizeOfField int32,
	pointXHandle int32,
	pointYHandle int32,
) int32 {
	bigInt := arwen.GetBigIntContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	x, y, P := bigInt.GetThree(pointXHandle, pointYHandle, fieldOrderHandle)
	N, B, Gx := bigInt.GetThree(basePointOrderHandle, eqConstantHandle, xBasePointHandle)
	Gy := bigInt.GetOne(yBasePointHandle)
	bigInt.ConsumeGasForBigIntCopy(P, N, B, Gx, Gy, x, y)

	if elliptic_curve.IsOnCurve(P, N, B, Gx, Gy, int(sizeOfField), x, y) {
		return 1
	}

	return 0
}

//export scalarBaseMultEC
func scalarBaseMultEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	fieldOrderHandle int32,
	basePointOrderHandle int32,
	eqConstantHandle int32,
	xBasePointHandle int32,
	yBasePointHandle int32,
	sizeOfField int32,
	kOffset int32,
	length int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	bigInt := arwen.GetBigIntContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	k, err := runtime.MemLoad(kOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	xResult, yResult, P := bigInt.GetThree(xResultHandle, yResultHandle, fieldOrderHandle)
	N, B, Gx := bigInt.GetThree(basePointOrderHandle, eqConstantHandle, xBasePointHandle)
	Gy := bigInt.GetOne(yBasePointHandle)
	bigInt.ConsumeGasForBigIntCopy(P, N, B, Gx, Gy, xResult, yResult)

	xResultSBM, yResultSBM := elliptic_curve.ScalarBaseMult(P, N, B, Gx, Gy, int(sizeOfField), k)
	xResult.Set(xResultSBM)
	yResult.Set(yResultSBM)

	return 0
}

//export scalarMultEC
func scalarMultEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	fieldOrderHandle int32,
	basePointOrderHandle int32,
	eqConstantHandle int32,
	xBasePointHandle int32,
	yBasePointHandle int32,
	sizeOfField int32,
	pointXHandle int32,
	pointYHandle int32,
	kOffset int32,
	length int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	bigInt := arwen.GetBigIntContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	k, err := runtime.MemLoad(kOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	xResult, yResult, P := bigInt.GetThree(xResultHandle, yResultHandle, fieldOrderHandle)
	N, B, Gx := bigInt.GetThree(basePointOrderHandle, eqConstantHandle, xBasePointHandle)
	Gy, x, y := bigInt.GetThree(yBasePointHandle, pointXHandle, pointYHandle)
	bigInt.ConsumeGasForBigIntCopy(xResult, yResult, P, N, B, Gx, Gy, x, y)

	xResultSM, yResultSM := elliptic_curve.ScalarMult(P, N, B, Gx, Gy, int(sizeOfField), x, y, k)
	xResult.Set(xResultSM)
	yResult.Set(yResultSM)

	return 0
}

//export marshalEC
func marshalEC(
	context unsafe.Pointer,
	fieldOrderHandle int32,
	basePointOrderHandle int32,
	eqConstantHandle int32,
	xBasePointHandle int32,
	yBasePointHandle int32,
	sizeOfField int32,
	pointXHandle int32,
	pointYHandle int32,
	resultOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	bigInt := arwen.GetBigIntContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	x, y, P := bigInt.GetThree(pointXHandle, pointYHandle, fieldOrderHandle)
	N, B, Gx := bigInt.GetThree(basePointOrderHandle, eqConstantHandle, xBasePointHandle)
	Gy := bigInt.GetOne(yBasePointHandle)
	bigInt.ConsumeGasForBigIntCopy(P, N, B, Gx, Gy, x, y)

	result := elliptic_curve.Marshal(P, N, B, Gx, Gy, int(sizeOfField), x, y)

	err := runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(result))
	}

	return 0
}

//export marshalCompressedEC
func marshalCompressedEC(
	context unsafe.Pointer,
	fieldOrderHandle int32,
	basePointOrderHandle int32,
	eqConstantHandle int32,
	xBasePointHandle int32,
	yBasePointHandle int32,
	sizeOfField int32,
	pointXHandle int32,
	pointYHandle int32,
	resultOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	bigInt := arwen.GetBigIntContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	x, y, P := bigInt.GetThree(pointXHandle, pointYHandle, fieldOrderHandle)
	N, B, Gx := bigInt.GetThree(basePointOrderHandle, eqConstantHandle, xBasePointHandle)
	Gy := bigInt.GetOne(yBasePointHandle)
	bigInt.ConsumeGasForBigIntCopy(P, N, B, Gx, Gy, x, y)

	result := elliptic_curve.MarshalCompressed(P, N, B, Gx, Gy, int(sizeOfField), x, y)

	err := runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(result))
	}

	return 0
}

//export unmarshalEC
func unmarshalEC(
	context unsafe.Pointer,
	xPairHandle int32,
	yPairHandle int32,
	fieldOrderHandle int32,
	basePointOrderHandle int32,
	eqConstantHandle int32,
	xBasePointHandle int32,
	yBasePointHandle int32,
	sizeOfField int32,
	dataOffset int32,
	length int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	bigInt := arwen.GetBigIntContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	xPair, yPair, P := bigInt.GetThree(xPairHandle, yPairHandle, fieldOrderHandle)
	N, B, Gx := bigInt.GetThree(basePointOrderHandle, eqConstantHandle, xBasePointHandle)
	Gy := bigInt.GetOne(yBasePointHandle)
	bigInt.ConsumeGasForBigIntCopy(P, N, B, Gx, Gy, xPair, yPair)

	xPairU, yPairU := elliptic_curve.Unmarshal(P, N, B, Gx, Gy, int(sizeOfField), data)
	xPair.Set(xPairU)
	yPair.Set(yPairU)

	return 0
}

//export unmarshalCompressedEC
func unmarshalCompressedEC(
	context unsafe.Pointer,
	xPairHandle int32,
	yPairHandle int32,
	fieldOrderHandle int32,
	basePointOrderHandle int32,
	eqConstantHandle int32,
	xBasePointHandle int32,
	yBasePointHandle int32,
	sizeOfField int32,
	dataOffset int32,
	length int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	bigInt := arwen.GetBigIntContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(data))
	}

	xPair, yPair, P := bigInt.GetThree(xPairHandle, yPairHandle, fieldOrderHandle)
	N, B, Gx := bigInt.GetThree(basePointOrderHandle, eqConstantHandle, xBasePointHandle)
	Gy := bigInt.GetOne(yBasePointHandle)
	bigInt.ConsumeGasForBigIntCopy(P, N, B, Gx, Gy, xPair, yPair)

	xPairUC, yPairUC := elliptic_curve.UnmarshalCompressed(P, N, B, Gx, Gy, int(sizeOfField), data)
	xPair.Set(xPairUC)
	yPair.Set(yPairUC)

	return 0
}

//export generateKeyEC
func generateKeyEC(
	context unsafe.Pointer,
	xPubKeyHandle int32,
	yPubKeyHandle int32,
	fieldOrderHandle int32,
	basePointOrderHandle int32,
	eqConstantHandle int32,
	xBasePointHandle int32,
	yBasePointHandle int32,
	sizeOfField int32,
	resultOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	bigInt := arwen.GetBigIntContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	xPubKey, yPubKey, P := bigInt.GetThree(xPubKeyHandle, yPubKeyHandle, fieldOrderHandle)
	N, B, Gx := bigInt.GetThree(basePointOrderHandle, eqConstantHandle, xBasePointHandle)
	Gy := bigInt.GetOne(yBasePointHandle)
	bigInt.ConsumeGasForBigIntCopy(P, N, B, Gx, Gy, xPubKey, yPubKey)

	result, xPubKeyGK, yPubKeyGK, err := elliptic_curve.GenerateKey(P, N, N, Gx, Gy, int(sizeOfField))
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(result))
	}

	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(result))
	}

	xPubKey.Set(xPubKeyGK)
	yPubKey.Set(yPubKeyGK)

	return 0
}

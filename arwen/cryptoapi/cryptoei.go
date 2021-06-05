package cryptoapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t v1_3_sha256(void* context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t v1_3_keccak256(void *context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t v1_3_ripemd160(void *context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t v1_3_verifyBLS(void *context, int32_t keyOffset, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern int32_t v1_3_verifyEd25519(void *context, int32_t keyOffset, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern int32_t v1_3_verifySecp256k1(void *context, int32_t keyOffset, int32_t keyLength, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern void v1_3_addEC(void *context, int32_t destination1, int32_t destination2, int32_t ecHandle, int32_t fstPointX, int32_t fstPointY, int32_t sndPointX, int32_t sndPointY);
// extern void v1_3_doubleEC(void *context, int32_t destination1, int32_t destination2, int32_t ecHandle, int32_t pointX, int32_t pointY);
// extern int32_t v1_3_isOnCurveEC(void *context, int32_t ecHandle, int32_t pointX, int32_t pointY);
// extern int32_t v1_3_scalarBaseMultEC(void *context, int32_t destination1, int32_t destination2, int32_t ecHandle, int32_t kOffset, int32_t length);
// extern int32_t v1_3_scalarMultEC(void *context, int32_t destination1, int32_t destination2, int32_t ecHandle, int32_t pointX, int32_t pointY, int32_t kOffset, int32_t length);
// extern int32_t v1_3_marshalEC(void *context, int32_t ecHandle, int32_t xPairHandle, int32_t yPairHandle, int32_t resultOffest);
// extern int32_t v1_3_unmarshalEC(void *context, int32_t xPairHandle, int32_t yPairHandle, int32_t ecHandle, int32_t dataOffest, int32_t length);
// extern int32_t v1_3_marshalCompressedEC(void *context, int32_t ecHandle, int32_t xPairHandle, int32_t yPairHandle, int32_t resultOffest);
// extern int32_t v1_3_unmarshalCompressedEC(void *context, int32_t xPairHandle, int32_t yPairHandle, int32_t ecHandle, int32_t dataOffest, int32_t length);
// extern int32_t v1_3_generateKeyEC(void *context, int32_t xPubKeyHandle, int32_t yPubKeyHandle, int32_t ecHandle, int32_t resultOffset);
// extern int32_t v1_3_ellipticCurveNew(void *context, int32_t fieldOrderHandle, int32_t basePointOrderHandle, int32_t eqConstantHandle, int32_t xBasePointHandle, int32_t yBasePointHandle, int32_t sizeOfField);
// extern int32_t v1_3_p224Ec(void *context);
// extern int32_t v1_3_p256Ec(void *context);
// extern int32_t v1_3_p384Ec(void *context);
// extern int32_t v1_3_p521Ec(void *context);
// extern int32_t v1_3_getCurveLengthEC(void *context, int32_t ecHandle);
// extern int32_t v1_3_getCurveByteLengthEC(void *context, int32_t ecHandle);
import "C"

import (
	"crypto/elliptic"
	"crypto/rand"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/math"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/wasmer"
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
	imports, err := imports.Append("sha256", v1_3_sha256, C.v1_3_sha256)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("keccak256", v1_3_keccak256, C.v1_3_keccak256)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("ripemd160", v1_3_ripemd160, C.v1_3_ripemd160)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("verifyBLS", v1_3_verifyBLS, C.v1_3_verifyBLS)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("verifyEd25519", v1_3_verifyEd25519, C.v1_3_verifyEd25519)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("verifySecp256k1", v1_3_verifySecp256k1, C.v1_3_verifySecp256k1)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("addEC", v1_3_addEC, C.v1_3_addEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("doubleEC", v1_3_doubleEC, C.v1_3_doubleEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("isOnCurveEC", v1_3_isOnCurveEC, C.v1_3_isOnCurveEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("scalarBaseMultEC", v1_3_scalarBaseMultEC, C.v1_3_scalarBaseMultEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("scalarMultEC", v1_3_scalarMultEC, C.v1_3_scalarMultEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("marshalEC", v1_3_marshalEC, C.v1_3_marshalEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("unmarshalEC", v1_3_unmarshalEC, C.v1_3_unmarshalEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("marshalCompressedEC", v1_3_marshalCompressedEC, C.v1_3_marshalCompressedEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("unmarshalCompressedEC", v1_3_unmarshalCompressedEC, C.v1_3_unmarshalCompressedEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("generateKeyEC", v1_3_generateKeyEC, C.v1_3_generateKeyEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("ellipticCurveNew", v1_3_ellipticCurveNew, C.v1_3_ellipticCurveNew)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("p224Ec", v1_3_p224Ec, C.v1_3_p224Ec)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCurveLengthEC", v1_3_getCurveLengthEC, C.v1_3_getCurveLengthEC)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrivKeyLengthEC", v1_3_getCurveByteLengthEC, C.v1_3_getCurveByteLengthEC)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export v1_3_sha256
func v1_3_sha256(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
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

//export v1_3_keccak256
func v1_3_keccak256(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
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

//export v1_3_ripemd160
func v1_3_ripemd160(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
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

//export v1_3_verifyBLS
func v1_3_verifyBLS(
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

//export v1_3_verifyEd25519
func v1_3_verifyEd25519(
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

//export v1_3_verifySecp256k1
func v1_3_verifySecp256k1(
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

//export v1_3_addEC
func v1_3_addEC(
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

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	ec, err1 := managedType.GetEllipticCurve(ecHandle)
	if err1 != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}

	xResult, err1 := managedType.GetBigInt(xResultHandle)
	yResult, err2 := managedType.GetBigInt(yResultHandle)
	x1, err3 := managedType.GetBigInt(fstPointXHandle)
	y1, err4 := managedType.GetBigInt(fstPointYHandle)
	x2, err5 := managedType.GetBigInt(sndPointXHandle)
	y2, err6 := managedType.GetBigInt(sndPointYHandle)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}

	managedType.ConsumeGasForBigIntCopy(xResult, yResult, ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x1, y1, x2, y2)
	xResultAdd, yResultAdd := ec.Add(x1, x2, y1, y2)
	xResult.Set(xResultAdd)
	yResult.Set(yResultAdd)
}

//export v1_3_doubleEC
func v1_3_doubleEC(
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

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	ec, err1 := managedType.GetEllipticCurve(ecHandle)
	if err1 != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}

	xResult, err1 := managedType.GetBigInt(xResultHandle)
	yResult, err2 := managedType.GetBigInt(yResultHandle)
	x, err3 := managedType.GetBigInt(pointXHandle)
	y, err4 := managedType.GetBigInt(pointYHandle)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}

	managedType.ConsumeGasForBigIntCopy(xResult, yResult, ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	xResultDouble, yResultDouble := ec.Double(x, y)
	xResult.Set(xResultDouble)
	yResult.Set(yResultDouble)
}

//export v1_3_isOnCurveEC
func v1_3_isOnCurveEC(
	context unsafe.Pointer,
	ecHandle int32,
	pointXHandle int32,
	pointYHandle int32,
) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	x, err1 := managedType.GetBigInt(pointXHandle)
	y, err2 := managedType.GetBigInt(pointYHandle)
	if err1 != nil || err2 != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	if ec.IsOnCurve(x, y) {
		return 1
	}

	return 0
}

//export v1_3_scalarBaseMultEC
func v1_3_scalarBaseMultEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	kOffset int32,
	length int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	k, err := runtime.MemLoad(kOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	xResult, err1 := managedType.GetBigInt(xResultHandle)
	yResult, err2 := managedType.GetBigInt(yResultHandle)
	if err1 != nil || err2 != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xResult, yResult)
	xResultSBM, yResultSBM := ec.ScalarBaseMult(k)
	xResult.Set(xResultSBM)
	yResult.Set(yResultSBM)

	return 0
}

//export v1_3_scalarMultEC
func v1_3_scalarMultEC(
	context unsafe.Pointer,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	pointXHandle int32,
	pointYHandle int32,
	kOffset int32,
	length int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	k, err := runtime.MemLoad(kOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	ec, err1 := managedType.GetEllipticCurve(ecHandle)
	if err1 != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	xResult, err1 := managedType.GetBigInt(xResultHandle)
	yResult, err2 := managedType.GetBigInt(yResultHandle)
	x, err3 := managedType.GetBigInt(pointXHandle)
	y, err4 := managedType.GetBigInt(pointYHandle)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(xResult, yResult, ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	xResultSM, yResultSM := ec.ScalarMult(x, y, k)
	xResult.Set(xResultSM)
	yResult.Set(yResultSM)

	return 0
}

//export v1_3_marshalEC
func v1_3_marshalEC(
	context unsafe.Pointer,
	pointXHandle int32,
	pointYHandle int32,
	ecHandle int32,
	resultOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	x, err1 := managedType.GetBigInt(pointXHandle)
	y, err2 := managedType.GetBigInt(pointYHandle)
	if err1 != nil || err2 != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	if x.BitLen() > int(ec.BitSize) || y.BitLen() > int(ec.BitSize) {
		arwen.WithFault(arwen.ErrBufNotBigEnough, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	result := elliptic.Marshal(ec, x, y)
	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(result))
	}
	return 0
}

//export v1_3_marshalCompressedEC
func v1_3_marshalCompressedEC(
	context unsafe.Pointer,
	pointXHandle int32,
	pointYHandle int32,
	ecHandle int32,
	resultOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	x, err1 := managedType.GetBigInt(pointXHandle)
	y, err2 := managedType.GetBigInt(pointYHandle)
	if err1 != nil || err2 != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	if x.BitLen() > int(ec.BitSize) || y.BitLen() > int(ec.BitSize) {
		arwen.WithFault(arwen.ErrBufNotBigEnough, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	result := elliptic.MarshalCompressed(ec, x, y)
	err = runtime.MemStore(resultOffset, result)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(result))
	}
	return 0
}

//export v1_3_unmarshalEC
func v1_3_unmarshalEC(
	context unsafe.Pointer,
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	dataOffset int32,
	length int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	xPair, err1 := managedType.GetBigInt(xPairHandle)
	yPair, err2 := managedType.GetBigInt(yPairHandle)
	if err1 != nil || err2 != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xPair, yPair)
	xPairU, yPairU := elliptic.Unmarshal(ec, data)
	xPair.Set(xPairU)
	yPair.Set(yPairU)

	return 0
}

//export v1_3_unmarshalCompressedEC
func v1_3_unmarshalCompressedEC(
	context unsafe.Pointer,
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	dataOffset int32,
	length int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(data))
	}

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	xPair, err1 := managedType.GetBigInt(xPairHandle)
	yPair, err2 := managedType.GetBigInt(yPairHandle)
	if err1 != nil || err2 != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xPair, yPair)
	xPairUC, yPairUC := elliptic.UnmarshalCompressed(ec, data)
	xPair.Set(xPairUC)
	yPair.Set(yPairUC)
	return 0
}

//export v1_3_generateKeyEC
func v1_3_generateKeyEC(
	context unsafe.Pointer,
	xPubKeyHandle int32,
	yPubKeyHandle int32,
	ecHandle int32,
	resultOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256
	metering.UseGas(gasToUse)

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	xPubKey, err1 := managedType.GetBigInt(xPubKeyHandle)
	yPubKey, err2 := managedType.GetBigInt(yPubKeyHandle)
	if err1 != nil || err2 != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xPubKey, yPubKey)

	result, xPubKeyGK, yPubKeyGK, err := elliptic.GenerateKey(ec, rand.Reader)
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

//export v1_3_ellipticCurveNew
func v1_3_ellipticCurveNew(context unsafe.Pointer, fieldOrderHandle int32, basePointOrderHandle int32, eqConstantHandle int32, xBasePointHandle int32, yBasePointHandle int32, sizeOfField int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.EllipticCurveNew
	metering.UseGas(gasToUse)

	P, err1 := managedType.GetBigInt(fieldOrderHandle)
	N, err2 := managedType.GetBigInt(basePointOrderHandle)
	B, err3 := managedType.GetBigInt(eqConstantHandle)
	Gx, err4 := managedType.GetBigInt(xBasePointHandle)
	Gy, err5 := managedType.GetBigInt(yBasePointHandle)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.BigIntAPIErrorShouldFailExecution())
		return -1
	}
	// TODO
	// should I verify? are the bigInt values in this case topencoded?
	// if P.BitLen() != int(sizeOfField) || N.BitLen() != int(sizeOfField) || B.BitLen() != int(sizeOfField) || Gx.BitLen() != int(sizeOfField) || Gy.BitLen() != int(sizeOfField) {
	// 	arwen.WithFault(arwen.ErrNoBigIntUnderThisHandle, context, runtime.BigIntAPIErrorShouldFailExecution())
	// 	return -1
	// }
	curve := elliptic.CurveParams{P: P, N: N, B: B, Gx: Gx, Gy: Gy, BitSize: int(sizeOfField), Name: "EC"}

	return managedType.PutEllipticCurve(&curve)
}

//export v1_3_p224Ec
func v1_3_p224Ec(context unsafe.Pointer) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.EllipticCurveNew
	metering.UseGas(gasToUse)

	curveParams := elliptic.P224().Params()
	return managedType.PutEllipticCurve(curveParams)
}

//export v1_3_p256Ec
func v1_3_p256Ec(context unsafe.Pointer) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.EllipticCurveNew
	metering.UseGas(gasToUse)

	curveParams := elliptic.P256().Params()
	return managedType.PutEllipticCurve(curveParams)
}

//export v1_3_p384Ec
func v1_3_p384Ec(context unsafe.Pointer) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.EllipticCurveNew
	metering.UseGas(gasToUse)

	curveParams := elliptic.P384().Params()
	return managedType.PutEllipticCurve(curveParams)
}

//export v1_3_p521Ec
func v1_3_p521Ec(context unsafe.Pointer) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.EllipticCurveNew
	metering.UseGas(gasToUse)

	curveParams := elliptic.P521().Params()
	return managedType.PutEllipticCurve(curveParams)
}

//export v1_3_getCurveLengthEC
func v1_3_getCurveLengthEC(context unsafe.Pointer, ecHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.EllipticCurveNew / 5
	metering.UseGas(gasToUse)

	ecLength := managedType.GetEllipticCurveSizeOfField(ecHandle)
	if ecLength == -1 {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.BigIntAPIErrorShouldFailExecution())
	}

	return ecLength
}

//export v1_3_getCurveByteLengthEC
func v1_3_getCurveByteLengthEC(context unsafe.Pointer, ecHandle int32) int32 {
	managedType := arwen.GetManagedTypesContext(context)
	metering := arwen.GetMeteringContext(context)
	runtime := arwen.GetRuntimeContext(context)

	gasToUse := metering.GasSchedule().BigIntAPICost.EllipticCurveNew / 5
	metering.UseGas(gasToUse)

	privKeyLength := managedType.GetEllipticCurveByteLength(ecHandle)
	if privKeyLength == -1 {
		arwen.WithFault(arwen.ErrNoEllipticCurveUnderThisHandle, context, runtime.BigIntAPIErrorShouldFailExecution())
	}

	return privKeyLength
}

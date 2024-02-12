package vmhooks

import (
	"crypto/elliptic"
	"github.com/multiversx/mx-chain-vm-go/crypto/signing/secp256"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/math"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
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
	verifyBLSSignatureShare         = "verifyBLSSignatureShare"
	verifyBLSAggregatedSignature    = "verifyBLSAggregatedSignature"
	verifySecp256R1Signature        = "verifySecp256R1Signature"
)

// Sha256 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) Sha256(
	dataOffset executor.MemPtr,
	length executor.MemLength,
	resultOffset executor.MemPtr) int32 {

	runtime := context.GetRuntimeContext()
	crypto := context.GetCryptoContext()
	metering := context.GetMeteringContext()

	memLoadGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(metering.GasSchedule().CryptoAPICost.SHA256, memLoadGas)
	metering.UseGasAndAddTracedGas(sha256Name, gasToUse)

	data, err := context.MemLoad(dataOffset, length)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	result, err := crypto.Sha256(data)
	if err != nil {
		context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	err = context.MemStore(resultOffset, result)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

// ManagedSha256 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedSha256(inputHandle, outputHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	crypto := context.GetCryptoContext()
	metering := context.GetMeteringContext()

	metering.UseGasAndAddTracedGas(sha256Name, metering.GasSchedule().CryptoAPICost.SHA256)

	inputBytes, err := managedType.GetBytes(inputHandle)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(inputBytes)

	resultBytes, err := crypto.Sha256(inputBytes)
	if err != nil {
		context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.SetBytes(outputHandle, resultBytes)

	return 0
}

// Keccak256 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) Keccak256(dataOffset executor.MemPtr, length executor.MemLength, resultOffset executor.MemPtr) int32 {
	runtime := context.GetRuntimeContext()
	crypto := context.GetCryptoContext()
	metering := context.GetMeteringContext()

	memLoadGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(metering.GasSchedule().CryptoAPICost.Keccak256, memLoadGas)
	metering.UseGasAndAddTracedGas(keccak256Name, gasToUse)

	data, err := context.MemLoad(dataOffset, length)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	result, err := crypto.Keccak256(data)
	if err != nil {
		context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	err = context.MemStore(resultOffset, result)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

// ManagedKeccak256 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedKeccak256(inputHandle, outputHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	crypto := context.GetCryptoContext()
	metering := context.GetMeteringContext()

	metering.UseGasAndAddTracedGas(keccak256Name, metering.GasSchedule().CryptoAPICost.Keccak256)

	inputBytes, err := managedType.GetBytes(inputHandle)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(inputBytes)

	resultBytes, err := crypto.Keccak256(inputBytes)
	if err != nil {
		context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.SetBytes(outputHandle, resultBytes)

	return 0
}

// Ripemd160 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) Ripemd160(dataOffset executor.MemPtr, length executor.MemLength, resultOffset executor.MemPtr) int32 {
	runtime := context.GetRuntimeContext()
	crypto := context.GetCryptoContext()
	metering := context.GetMeteringContext()

	memLoadGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(metering.GasSchedule().CryptoAPICost.Ripemd160, memLoadGas)
	metering.UseGasAndAddTracedGas(ripemd160Name, gasToUse)

	data, err := context.MemLoad(dataOffset, length)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	result, err := crypto.Ripemd160(data)
	if err != nil {
		context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	err = context.MemStore(resultOffset, result)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

// ManagedRipemd160 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedRipemd160(inputHandle int32, outputHandle int32) int32 {
	host := context.GetVMHost()
	return ManagedRipemd160WithHost(host, inputHandle, outputHandle)
}

// ManagedRipemd160WithHost VMHooks implementation.
func ManagedRipemd160WithHost(host vmhost.VMHost, inputHandle int32, outputHandle int32) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	crypto := host.Crypto()

	metering.UseGasAndAddTracedGas(ripemd160Name, metering.GasSchedule().CryptoAPICost.Ripemd160)

	inputBytes, err := managedType.GetBytes(inputHandle)
	if WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(inputBytes)

	result, err := crypto.Ripemd160(inputBytes)
	if err != nil {
		WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.SetBytes(outputHandle, result)

	return 0
}

// VerifyBLS VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) VerifyBLS(
	keyOffset executor.MemPtr,
	messageOffset executor.MemPtr,
	messageLength executor.MemLength,
	sigOffset executor.MemPtr,
) int32 {
	runtime := context.GetRuntimeContext()
	crypto := context.GetCryptoContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(verifyBLSName)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifyBLS
	metering.UseAndTraceGas(gasToUse)

	key, err := context.MemLoad(keyOffset, blsPublicKeyLength)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(messageLength))
	metering.UseAndTraceGas(gasToUse)

	message, err := context.MemLoad(messageOffset, messageLength)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	sig, err := context.MemLoad(sigOffset, blsSignatureLength)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	invalidSigErr := crypto.VerifyBLS(key, message, sig)
	if invalidSigErr != nil {
		context.WithFault(invalidSigErr, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	return 0
}

// ManagedVerifyBLS VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedVerifyBLS(
	keyHandle int32,
	messageHandle int32,
	sigHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedVerifyBLSWithHost(host, keyHandle, messageHandle, sigHandle, verifyBLSName)
}

func useGasForCryptoVerify(
	metering vmhost.MeteringContext,
	sigVerificationType string,
) {
	metering.StartGasTracing(sigVerificationType)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifyBLS
	switch sigVerificationType {
	case verifyCustomSecp256k1Name:
		gasToUse = metering.GasSchedule().CryptoAPICost.VerifySecp256k1
	case verifySecp256R1Signature:
		gasToUse = metering.GasSchedule().CryptoAPICost.VerifySecp256r1
	case verifyBLSName:
		gasToUse = metering.GasSchedule().CryptoAPICost.VerifyBLS
	case verifyBLSSignatureShare:
		gasToUse = metering.GasSchedule().CryptoAPICost.VerifyBLSSignatureShare
	case verifyBLSAggregatedSignature:
		gasToUse = metering.GasSchedule().CryptoAPICost.VerifyBLSMultiSig
	}
	metering.UseAndTraceGas(gasToUse)
}

// ManagedVerifyBLSWithHost VMHooks implementation.
func ManagedVerifyBLSWithHost(
	host vmhost.VMHost,
	keyHandle int32,
	messageHandle int32,
	sigHandle int32,
	sigVerificationType string,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	crypto := host.Crypto()
	useGasForCryptoVerify(metering, sigVerificationType)

	keyBytes, err := managedType.GetBytes(keyHandle)
	if WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(keyBytes)

	msgBytes, err := managedType.GetBytes(messageHandle)
	if WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(msgBytes)

	sigBytes, err := managedType.GetBytes(sigHandle)
	if WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(sigBytes)

	invalidSigErr := vmhost.ErrInvalidArgument
	switch sigVerificationType {
	case verifyBLSName:
		invalidSigErr = crypto.VerifyBLS(keyBytes, msgBytes, sigBytes)
	case verifyBLSSignatureShare:
		invalidSigErr = crypto.VerifySignatureShare(keyBytes, msgBytes, sigBytes)
	case verifyBLSAggregatedSignature:
		var pubKeyBytes [][]byte
		pubKeyBytes, _, invalidSigErr = managedType.ReadManagedVecOfManagedBuffers(keyHandle)
		if invalidSigErr != nil {
			break
		}
		invalidSigErr = crypto.VerifyAggregatedSig(pubKeyBytes, msgBytes, sigBytes)
	}

	if invalidSigErr != nil {
		WithFaultAndHost(host, invalidSigErr, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	return 0
}

// VerifyEd25519 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) VerifyEd25519(
	keyOffset executor.MemPtr,
	messageOffset executor.MemPtr,
	messageLength executor.MemLength,
	sigOffset executor.MemPtr,
) int32 {
	runtime := context.GetRuntimeContext()
	crypto := context.GetCryptoContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(verifyEd25519Name)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifyEd25519
	metering.UseAndTraceGas(gasToUse)

	key, err := context.MemLoad(keyOffset, ed25519PublicKeyLength)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(messageLength))
	metering.UseAndTraceGas(gasToUse)

	message, err := context.MemLoad(messageOffset, messageLength)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	sig, err := context.MemLoad(sigOffset, ed25519SignatureLength)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	invalidSigErr := crypto.VerifyEd25519(key, message, sig)
	if invalidSigErr != nil {
		context.WithFault(invalidSigErr, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	return 0
}

// ManagedVerifyEd25519 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedVerifyEd25519(
	keyHandle, messageHandle, sigHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedVerifyEd25519WithHost(host, keyHandle, messageHandle, sigHandle)
}

// ManagedVerifyEd25519WithHost VMHooks implementation.
func ManagedVerifyEd25519WithHost(
	host vmhost.VMHost,
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
	if WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(keyBytes)

	msgBytes, err := managedType.GetBytes(messageHandle)
	if WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(msgBytes)

	sigBytes, err := managedType.GetBytes(sigHandle)
	if WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(sigBytes)

	invalidSigErr := crypto.VerifyEd25519(keyBytes, msgBytes, sigBytes)
	if invalidSigErr != nil {
		WithFaultAndHost(host, invalidSigErr, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	return 0
}

// VerifyCustomSecp256k1 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) VerifyCustomSecp256k1(
	keyOffset executor.MemPtr,
	keyLength executor.MemLength,
	messageOffset executor.MemPtr,
	messageLength executor.MemLength,
	sigOffset executor.MemPtr,
	hashType int32,
) int32 {
	runtime := context.GetRuntimeContext()
	crypto := context.GetCryptoContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(verifyCustomSecp256k1Name)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifySecp256k1
	metering.UseAndTraceGas(gasToUse)

	if keyLength != secp256k1CompressedPublicKeyLength && keyLength != secp256k1UncompressedPublicKeyLength {
		_ = context.WithFault(vmhost.ErrInvalidPublicKeySize, runtime.BaseOpsErrorShouldFailExecution())
		return 1
	}

	key, err := context.MemLoad(keyOffset, keyLength)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(messageLength))
	metering.UseAndTraceGas(gasToUse)

	message, err := context.MemLoad(messageOffset, messageLength)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	// read the 2 leading bytes first
	// byte1: 0x30, header
	// byte2: the remaining buffer length
	const sigHeaderLength = 2
	sigHeader, err := context.MemLoad(sigOffset, sigHeaderLength)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}
	sigLength := int32(sigHeader[1]) + sigHeaderLength
	sig, err := context.MemLoad(sigOffset, sigLength)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	invalidSigErr := crypto.VerifySecp256k1(key, message, sig, uint8(hashType))
	if invalidSigErr != nil {
		context.WithFault(invalidSigErr, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	return 0
}

// ManagedVerifyCustomSecp256k1 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedVerifyCustomSecp256k1(
	keyHandle, messageHandle, sigHandle int32,
	hashType int32,
) int32 {
	host := context.GetVMHost()
	return ManagedVerifyCustomSecp256k1WithHost(
		host,
		keyHandle,
		messageHandle,
		sigHandle,
		hashType,
		verifyCustomSecp256k1Name)
}

// ManagedVerifyCustomSecp256k1WithHost VMHooks implementation.
func ManagedVerifyCustomSecp256k1WithHost(
	host vmhost.VMHost,
	keyHandle, messageHandle, sigHandle int32,
	hashType int32,
	verifyCryptoFunc string,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	crypto := host.Crypto()

	useGasForCryptoVerify(metering, verifyCryptoFunc)

	keyBytes, err := managedType.GetBytes(keyHandle)
	if WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(keyBytes)

	msgBytes, err := managedType.GetBytes(messageHandle)
	if WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(msgBytes)

	sigBytes, err := managedType.GetBytes(sigHandle)
	if WithFaultAndHost(host, err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return 1
	}
	managedType.ConsumeGasForBytes(sigBytes)

	invalidSigErr := vmhost.ErrInvalidArgument
	switch verifyCryptoFunc {
	case verifyCustomSecp256k1Name:
		invalidSigErr = crypto.VerifySecp256k1(keyBytes, msgBytes, sigBytes, uint8(hashType))
	case verifySecp256R1Signature:
		invalidSigErr = crypto.VerifySecp256r1(keyBytes, msgBytes, sigBytes)
	}

	if invalidSigErr != nil {
		WithFaultAndHost(host, invalidSigErr, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	return 0
}

// VerifySecp256k1 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) VerifySecp256k1(
	keyOffset executor.MemPtr,
	keyLength executor.MemLength,
	messageOffset executor.MemPtr,
	messageLength executor.MemLength,
	sigOffset executor.MemPtr,
) int32 {
	return context.VerifyCustomSecp256k1(
		keyOffset,
		keyLength,
		messageOffset,
		messageLength,
		sigOffset,
		int32(secp256.ECDSADoubleSha256),
	)
}

// ManagedVerifySecp256k1 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedVerifySecp256k1(
	keyHandle, messageHandle, sigHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedVerifySecp256k1WithHost(host, keyHandle, messageHandle, sigHandle)
}

// ManagedVerifySecp256k1WithHost VMHooks implementation.
func ManagedVerifySecp256k1WithHost(
	host vmhost.VMHost,
	keyHandle, messageHandle, sigHandle int32,
) int32 {
	return ManagedVerifyCustomSecp256k1WithHost(
		host,
		keyHandle,
		messageHandle,
		sigHandle,
		int32(secp256.ECDSADoubleSha256),
		verifyCustomSecp256k1Name,
	)
}

// EncodeSecp256k1DerSignature VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) EncodeSecp256k1DerSignature(
	rOffset executor.MemPtr,
	rLength executor.MemLength,
	sOffset executor.MemPtr,
	sLength executor.MemLength,
	sigOffset executor.MemPtr,
) int32 {
	runtime := context.GetRuntimeContext()
	crypto := context.GetCryptoContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().CryptoAPICost.EncodeDERSig
	metering.UseGasAndAddTracedGas(encodeSecp256k1DerSignatureName, gasToUse)

	r, err := context.MemLoad(rOffset, rLength)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	s, err := context.MemLoad(sOffset, sLength)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	derSig := crypto.EncodeSecp256k1DERSignature(r, s)
	err = context.MemStore(sigOffset, derSig)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

// ManagedEncodeSecp256k1DerSignature VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedEncodeSecp256k1DerSignature(
	rHandle, sHandle, sigHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedEncodeSecp256k1DerSignatureWithHost(host, rHandle, sHandle, sigHandle)
}

// ManagedEncodeSecp256k1DerSignatureWithHost VMHooks implementation.
func ManagedEncodeSecp256k1DerSignatureWithHost(
	host vmhost.VMHost,
	rHandle, sHandle, sigHandle int32,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	crypto := host.Crypto()

	gasToUse := metering.GasSchedule().CryptoAPICost.EncodeDERSig
	metering.UseGasAndAddTracedGas(encodeSecp256k1DerSignatureName, gasToUse)

	r, err := managedType.GetBytes(rHandle)
	if WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	s, err := managedType.GetBytes(sHandle)
	if WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	derSig := crypto.EncodeSecp256k1DERSignature(r, s)
	managedType.SetBytes(sigHandle, derSig)

	return 0
}

// AddEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) AddEC(
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	fstPointXHandle int32,
	fstPointYHandle int32,
	sndPointXHandle int32,
	sndPointYHandle int32,
) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(addECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = context.WithFault(vmhost.ErrNoEllipticCurveUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.AddECC * uint64(curveMultiplier) / 100
	metering.UseAndTraceGas(gasToUse)

	ec, err1 := managedType.GetEllipticCurve(ecHandle)
	if context.WithFault(err1, runtime.CryptoAPIErrorShouldFailExecution()) {
		return
	}

	xResult, yResult, err := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	if err != nil {
		_ = context.WithFault(vmhost.ErrNoBigIntUnderThisHandle, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	x1, y1, err := managedType.GetTwoBigInt(fstPointXHandle, fstPointYHandle)
	if err != nil {
		_ = context.WithFault(vmhost.ErrNoBigIntUnderThisHandle, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}
	x2, y2, err := managedType.GetTwoBigInt(sndPointXHandle, sndPointYHandle)
	if err != nil {
		_ = context.WithFault(vmhost.ErrNoBigIntUnderThisHandle, runtime.BigIntAPIErrorShouldFailExecution())
		return
	}

	if !ec.IsOnCurve(x1, y1) || !ec.IsOnCurve(x2, y2) {
		_ = context.WithFault(vmhost.ErrPointNotOnCurve, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}

	managedType.ConsumeGasForBigIntCopy(xResult, yResult, ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x1, y1, x2, y2)
	xResultAdd, yResultAdd := ec.Add(x1, y1, x2, y2)
	xResult.Set(xResultAdd)
	yResult.Set(yResultAdd)
}

// DoubleEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) DoubleEC(
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	pointXHandle int32,
	pointYHandle int32,
) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(doubleECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = context.WithFault(vmhost.ErrNoEllipticCurveUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.DoubleECC * uint64(curveMultiplier) / 100
	metering.UseAndTraceGas(gasToUse)

	ec, err1 := managedType.GetEllipticCurve(ecHandle)
	if context.WithFault(err1, runtime.CryptoAPIErrorShouldFailExecution()) {
		return
	}

	xResult, yResult, err1 := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	x, y, err2 := managedType.GetTwoBigInt(pointXHandle, pointYHandle)
	if err1 != nil || err2 != nil {
		_ = context.WithFault(vmhost.ErrNoBigIntUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}
	if !ec.IsOnCurve(x, y) {
		_ = context.WithFault(vmhost.ErrPointNotOnCurve, runtime.CryptoAPIErrorShouldFailExecution())
		return
	}

	managedType.ConsumeGasForBigIntCopy(xResult, yResult, ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)

	xResultDouble, yResultDouble := ec.Double(x, y)
	xResult.Set(xResultDouble)
	yResult.Set(yResultDouble)
}

// IsOnCurveEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) IsOnCurveEC(
	ecHandle int32,
	pointXHandle int32,
	pointYHandle int32,
) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()
	metering.StartGasTracing(isOnCurveECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = context.WithFault(vmhost.ErrNoEllipticCurveUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.IsOnCurveECC * uint64(curveMultiplier) / 100
	metering.UseAndTraceGas(gasToUse)

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return -1
	}

	x, y, err := managedType.GetTwoBigInt(pointXHandle, pointYHandle)
	if err != nil || x == nil || y == nil {
		_ = context.WithFault(vmhost.ErrNoBigIntUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	if ec.IsOnCurve(x, y) {
		return 1
	}

	return 0
}

// ScalarBaseMultEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ScalarBaseMultEC(
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataOffset executor.MemPtr,
	length executor.MemLength,
) int32 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()
	metering.StartGasTracing(scalarBaseMultECName)

	if length < 0 {
		_ = context.WithFault(vmhost.ErrNegativeLength, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	curveMultiplier := managedType.GetScalarMult100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = context.WithFault(vmhost.ErrNoEllipticCurveUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	oneByteScalarGasCost := metering.GasSchedule().CryptoAPICost.ScalarMultECC * uint64(curveMultiplier) / 100
	gasToUse := oneByteScalarGasCost + uint64(length)*oneByteScalarGasCost
	metering.UseAndTraceGas(gasToUse)

	data, err := context.MemLoad(dataOffset, length)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	host := context.GetVMHost()
	return commonScalarBaseMultEC(host, xResultHandle, yResultHandle, ecHandle, data)
}

// ManagedScalarBaseMultEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedScalarBaseMultEC(
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedScalarBaseMultECWithHost(
		host,
		xResultHandle,
		yResultHandle,
		ecHandle,
		dataHandle,
	)
}

// ManagedScalarBaseMultECWithHost VMHooks implementation.
func ManagedScalarBaseMultECWithHost(
	host vmhost.VMHost,
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
		_ = WithFaultAndHost(host, vmhost.ErrNoEllipticCurveUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	data, err := managedType.GetBytes(dataHandle)
	if WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	oneByteScalarGasCost := metering.GasSchedule().CryptoAPICost.ScalarMultECC * uint64(curveMultiplier) / 100
	gasToUse := oneByteScalarGasCost + uint64(len(data))*oneByteScalarGasCost
	metering.UseAndTraceGas(gasToUse)

	return commonScalarBaseMultEC(host, xResultHandle, yResultHandle, ecHandle, data)
}

func commonScalarBaseMultEC(
	host vmhost.VMHost,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	data []byte,
) int32 {
	runtime := host.Runtime()
	managedType := host.ManagedTypes()

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	xResult, yResult, err := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	if WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xResult, yResult)

	xResultSBM, yResultSBM := ec.ScalarBaseMult(data)
	if !ec.IsOnCurve(xResultSBM, yResultSBM) {
		_ = WithFaultAndHost(host, vmhost.ErrPointNotOnCurve, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	xResult.Set(xResultSBM)
	yResult.Set(yResultSBM)

	return 0
}

// ScalarMultEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ScalarMultEC(
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	pointXHandle int32,
	pointYHandle int32,
	dataOffset executor.MemPtr,
	length executor.MemLength,
) int32 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()
	metering.StartGasTracing(scalarMultECName)

	if length < 0 {
		_ = context.WithFault(vmhost.ErrNegativeLength, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	curveMultiplier := managedType.GetScalarMult100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = context.WithFault(vmhost.ErrNoEllipticCurveUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	oneByteScalarGasCost := metering.GasSchedule().CryptoAPICost.ScalarMultECC * uint64(curveMultiplier) / 100
	gasToUse := oneByteScalarGasCost + uint64(length)*oneByteScalarGasCost
	metering.UseAndTraceGas(gasToUse)

	data, err := context.MemLoad(dataOffset, length)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	host := context.GetVMHost()
	return commonScalarMultEC(host, xResultHandle, yResultHandle, ecHandle, pointXHandle, pointYHandle, data)
}

// ManagedScalarMultEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedScalarMultEC(
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	pointXHandle int32,
	pointYHandle int32,
	dataHandle int32,
) int32 {
	host := context.GetVMHost()
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

// ManagedScalarMultECWithHost VMHooks implementation.
func ManagedScalarMultECWithHost(
	host vmhost.VMHost,
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
		_ = WithFaultAndHost(host, vmhost.ErrNoEllipticCurveUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	data, err := managedType.GetBytes(dataHandle)
	if WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	oneByteScalarGasCost := metering.GasSchedule().CryptoAPICost.ScalarMultECC * uint64(curveMultiplier) / 100
	gasToUse := oneByteScalarGasCost + uint64(len(data))*oneByteScalarGasCost
	metering.UseAndTraceGas(gasToUse)

	return commonScalarMultEC(host, xResultHandle, yResultHandle, ecHandle, pointXHandle, pointYHandle, data)
}

func commonScalarMultEC(
	host vmhost.VMHost,
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
	if WithFaultAndHost(host, err1, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	xResult, yResult, err1 := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	x, y, err2 := managedType.GetTwoBigInt(pointXHandle, pointYHandle)
	if err1 != nil || err2 != nil {
		_ = WithFaultAndHost(host, vmhost.ErrNoBigIntUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	if !ec.IsOnCurve(x, y) {
		_ = WithFaultAndHost(host, vmhost.ErrPointNotOnCurve, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(xResult, yResult, ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	xResultSM, yResultSM := ec.ScalarMult(x, y, data)
	if !ec.IsOnCurve(xResultSM, yResultSM) {
		_ = WithFaultAndHost(host, vmhost.ErrPointNotOnCurve, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	xResult.Set(xResultSM)
	yResult.Set(yResultSM)

	return 0
}

// MarshalEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MarshalEC(
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	resultOffset executor.MemPtr,
) int32 {
	runtime := context.GetRuntimeContext()
	host := context.GetVMHost()
	result, err := commonMarshalEC(host, xPairHandle, yPairHandle, ecHandle)
	if err != nil {
		_ = context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	err = context.MemStore(resultOffset, result)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return -1
	}
	return int32(len(result))
}

// ManagedMarshalEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMarshalEC(
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedMarshalECWithHost(
		host,
		xPairHandle,
		yPairHandle,
		ecHandle,
		resultHandle,
	)
}

// ManagedMarshalECWithHost VMHooks implementation.
func ManagedMarshalECWithHost(
	host vmhost.VMHost,
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	resultHandle int32,
) int32 {
	result, err := commonMarshalEC(host, xPairHandle, yPairHandle, ecHandle)
	if err != nil {
		_ = WithFaultAndHost(host, err, true)
		return -1
	}

	managedType := host.ManagedTypes()
	managedType.SetBytes(resultHandle, result)
	return int32(len(result))
}

func commonMarshalEC(
	host vmhost.VMHost,
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
) ([]byte, error) {
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(marshalECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		return nil, vmhost.ErrNoEllipticCurveUnderThisHandle
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
		return nil, vmhost.ErrPointNotOnCurve
	}
	if x.BitLen() > ec.BitSize || y.BitLen() > ec.BitSize {
		return nil, vmhost.ErrLengthOfBufferNotCorrect
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)

	result := elliptic.Marshal(ec, x, y)
	return result, nil
}

// MarshalCompressedEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MarshalCompressedEC(
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	resultOffset executor.MemPtr,
) int32 {
	runtime := context.GetRuntimeContext()
	host := context.GetVMHost()
	result, err := commonMarshalCompressedEC(host, xPairHandle, yPairHandle, ecHandle)
	if err != nil {
		_ = context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	err = context.MemStore(resultOffset, result)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return -1
	}
	return int32(len(result))
}

// ManagedMarshalCompressedEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMarshalCompressedEC(
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedMarshalCompressedECWithHost(
		host,
		xPairHandle,
		yPairHandle,
		ecHandle,
		resultHandle,
	)
}

// ManagedMarshalCompressedECWithHost VMHooks implementation.
func ManagedMarshalCompressedECWithHost(
	host vmhost.VMHost,
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	resultHandle int32,
) int32 {
	runtime := host.Runtime()
	managedType := host.ManagedTypes()
	result, err := commonMarshalCompressedEC(host, xPairHandle, yPairHandle, ecHandle)
	if err != nil {
		_ = WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}

	managedType.SetBytes(resultHandle, result)
	return int32(len(result))
}

func commonMarshalCompressedEC(host vmhost.VMHost,
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
) ([]byte, error) {
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(marshalCompressedECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		return nil, vmhost.ErrNoEllipticCurveUnderThisHandle
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.MarshalCompressedECC * uint64(curveMultiplier) / 100
	metering.UseAndTraceGas(gasToUse)

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		return nil, err
	}

	x, y, err := managedType.GetTwoBigInt(xPairHandle, yPairHandle)
	if err != nil || x == nil || y == nil {
		return nil, vmhost.ErrNoBigIntUnderThisHandle
	}
	if !ec.IsOnCurve(x, y) {
		return nil, vmhost.ErrPointNotOnCurve
	}
	if x.BitLen() > ec.BitSize || y.BitLen() > ec.BitSize {
		return nil, vmhost.ErrLengthOfBufferNotCorrect
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)

	result := elliptic.MarshalCompressed(ec, x, y)
	return result, nil
}

// UnmarshalEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) UnmarshalEC(
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataOffset executor.MemPtr,
	length executor.MemLength,
) int32 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()
	metering.StartGasTracing(unmarshalECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = context.WithFault(vmhost.ErrNoEllipticCurveUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.UnmarshalECC * uint64(curveMultiplier) / 100
	metering.UseAndTraceGas(gasToUse)

	data, err := context.MemLoad(dataOffset, length)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	host := context.GetVMHost()
	return commonUnmarshalEC(host, xResultHandle, yResultHandle, ecHandle, data)
}

// ManagedUnmarshalEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedUnmarshalEC(
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedUnmarshalECWithHost(
		host,
		xResultHandle,
		yResultHandle,
		ecHandle,
		dataHandle,
	)
}

// ManagedUnmarshalECWithHost VMHooks implementation.
func ManagedUnmarshalECWithHost(
	host vmhost.VMHost,
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
		_ = WithFaultAndHost(host, vmhost.ErrNoEllipticCurveUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
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
	host vmhost.VMHost,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	data []byte,
) int32 {
	runtime := host.Runtime()
	managedType := host.ManagedTypes()

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}
	byteLen := (ec.BitSize + 7) / 8
	if len(data) != 1+2*byteLen {
		_ = WithFaultAndHost(host, vmhost.ErrLengthOfBufferNotCorrect, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	xResult, yResult, err := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	if WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xResult, yResult)

	xResultU, yResultU := elliptic.Unmarshal(ec, data)
	if xResultU == nil || yResultU == nil || !ec.IsOnCurve(xResultU, yResultU) {
		_ = WithFaultAndHost(host, vmhost.ErrPointNotOnCurve, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	xResult.Set(xResultU)
	yResult.Set(yResultU)

	return 0
}

// UnmarshalCompressedEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) UnmarshalCompressedEC(
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataOffset executor.MemPtr,
	length executor.MemLength,
) int32 {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()
	metering.StartGasTracing(unmarshalCompressedECName)

	curveMultiplier := managedType.GetUCompressed100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		_ = context.WithFault(vmhost.ErrNoEllipticCurveUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.UnmarshalCompressedECC * uint64(curveMultiplier) / 100
	metering.UseAndTraceGas(gasToUse)

	data, err := context.MemLoad(dataOffset, length)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(data))
	}

	host := context.GetVMHost()
	return commonUnmarshalCompressedEC(host, xResultHandle, yResultHandle, ecHandle, data)
}

// ManagedUnmarshalCompressedEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedUnmarshalCompressedEC(
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedUnmarshalCompressedECWithHost(
		host,
		xResultHandle,
		yResultHandle,
		ecHandle,
		dataHandle,
	)
}

// ManagedUnmarshalCompressedECWithHost VMHooks implementation.
func ManagedUnmarshalCompressedECWithHost(
	host vmhost.VMHost,
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
		_ = WithFaultAndHost(host, vmhost.ErrNoEllipticCurveUnderThisHandle, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.UnmarshalCompressedECC * uint64(curveMultiplier) / 100
	metering.UseAndTraceGas(gasToUse)

	data, err := managedType.GetBytes(dataHandle)
	if WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(data))
	}

	return commonUnmarshalCompressedEC(host, xResultHandle, yResultHandle, ecHandle, data)
}

func commonUnmarshalCompressedEC(
	host vmhost.VMHost,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	data []byte,
) int32 {
	runtime := host.Runtime()
	managedType := host.ManagedTypes()

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}
	byteLen := (ec.BitSize+7)/8 + 1
	if len(data) != byteLen {
		_ = WithFaultAndHost(host, vmhost.ErrLengthOfBufferNotCorrect, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}

	xResult, yResult, err := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	if WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xResult, yResult)

	xResultUC, yResultUC := elliptic.UnmarshalCompressed(ec, data)
	if xResultUC == nil || yResultUC == nil || !ec.IsOnCurve(xResultUC, yResultUC) {
		_ = WithFaultAndHost(host, vmhost.ErrPointNotOnCurve, runtime.CryptoAPIErrorShouldFailExecution())
		return 1
	}
	xResult.Set(xResultUC)
	yResult.Set(yResultUC)
	return 0
}

// GenerateKeyEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GenerateKeyEC(
	xPubKeyHandle int32,
	yPubKeyHandle int32,
	ecHandle int32,
	resultOffset executor.MemPtr,
) int32 {
	runtime := context.GetRuntimeContext()
	host := context.GetVMHost()
	result, err := commonGenerateEC(host, xPubKeyHandle, yPubKeyHandle, ecHandle)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	err = context.MemStore(resultOffset, result)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return int32(len(result))
	}

	return 0
}

// ManagedGenerateKeyEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGenerateKeyEC(
	xPubKeyHandle int32,
	yPubKeyHandle int32,
	ecHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedGenerateKeyECWithHost(
		host,
		xPubKeyHandle,
		yPubKeyHandle,
		ecHandle,
		resultHandle,
	)
}

// ManagedGenerateKeyECWithHost VMHooks implementation.
func ManagedGenerateKeyECWithHost(
	host vmhost.VMHost,
	xPubKeyHandle int32,
	yPubKeyHandle int32,
	ecHandle int32,
	resultHandle int32,
) int32 {
	runtime := host.Runtime()
	managedType := host.ManagedTypes()
	result, err := commonGenerateEC(host, xPubKeyHandle, yPubKeyHandle, ecHandle)
	if WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.SetBytes(resultHandle, result)
	return 0
}

func commonGenerateEC(
	host vmhost.VMHost,
	xPubKeyHandle int32,
	yPubKeyHandle int32,
	ecHandle int32,
) ([]byte, error) {
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(generateKeyECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		return nil, vmhost.ErrNoEllipticCurveUnderThisHandle
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

// CreateEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) CreateEC(dataOffset executor.MemPtr, dataLength executor.MemLength) int32 {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().CryptoAPICost.EllipticCurveNew
	metering.UseGasAndAddTracedGas(createECName, gasToUse)

	if dataLength != curveNameLength {
		_ = context.WithFault(vmhost.ErrBadBounds, runtime.CryptoAPIErrorShouldFailExecution())
		return -1
	}
	data, err := context.MemLoad(dataOffset, dataLength)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
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

// ManagedCreateEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedCreateEC(dataHandle int32) int32 {
	host := context.GetVMHost()
	return ManagedCreateECWithHost(host, dataHandle)
}

// ManagedCreateECWithHost VMHooks implementation.
func ManagedCreateECWithHost(host vmhost.VMHost, dataHandle int32) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().CryptoAPICost.EllipticCurveNew
	metering.UseGasAndAddTracedGas(createECName, gasToUse)

	data, err := managedType.GetBytes(dataHandle)
	if WithFaultAndHost(host, err, runtime.CryptoAPIErrorShouldFailExecution()) {
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

	_ = WithFaultAndHost(host, vmhost.ErrBadBounds, runtime.CryptoAPIErrorShouldFailExecution())
	return -1
}

// GetCurveLengthEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetCurveLengthEC(ecHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetInt64
	metering.UseGasAndAddTracedGas(getCurveLengthECName, gasToUse)

	ecLength := managedType.GetEllipticCurveSizeOfField(ecHandle)
	if ecLength == -1 {
		_ = context.WithFault(vmhost.ErrNoEllipticCurveUnderThisHandle, runtime.BigIntAPIErrorShouldFailExecution())
	}

	return ecLength
}

// GetPrivKeyByteLengthEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetPrivKeyByteLengthEC(ecHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetInt64
	metering.UseGasAndAddTracedGas(getPrivKeyByteLengthECName, gasToUse)

	byteLength := managedType.GetPrivateKeyByteLengthEC(ecHandle)
	if byteLength == -1 {
		_ = context.WithFault(vmhost.ErrNoEllipticCurveUnderThisHandle, runtime.BigIntAPIErrorShouldFailExecution())
	}

	return byteLength
}

// EllipticCurveGetValues VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) EllipticCurveGetValues(ecHandle int32, fieldOrderHandle int32, basePointOrderHandle int32, eqConstantHandle int32, xBasePointHandle int32, yBasePointHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetInt64 * 5
	metering.UseGasAndAddTracedGas(ellipticCurveGetValuesName, gasToUse)

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return -1
	}
	fieldOrder, basePointOrder, err := managedType.GetTwoBigInt(fieldOrderHandle, basePointOrderHandle)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return -1
	}
	eqConstant, err := managedType.GetBigInt(eqConstantHandle)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return -1
	}
	xBasePoint, yBasePoint, err := managedType.GetTwoBigInt(xBasePointHandle, yBasePointHandle)
	if context.WithFault(err, runtime.CryptoAPIErrorShouldFailExecution()) {
		return -1
	}
	fieldOrder.Set(ec.P)
	basePointOrder.Set(ec.N)
	eqConstant.Set(ec.B)
	xBasePoint.Set(ec.Gx)
	yBasePoint.Set(ec.Gy)
	return ecHandle
}

// ManagedVerifySecp256r1 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedVerifySecp256r1(
	keyHandle, messageHandle, sigHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedVerifyCustomSecp256k1WithHost(
		host,
		keyHandle,
		messageHandle,
		sigHandle,
		0,
		verifySecp256R1Signature)
}

// ManagedVerifyBLSSignatureShare VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedVerifyBLSSignatureShare(
	keyHandle int32,
	messageHandle int32,
	sigHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedVerifyBLSWithHost(host, keyHandle, messageHandle, sigHandle, verifyBLSSignatureShare)
}

// ManagedVerifyBLSAggregatedSignature VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedVerifyBLSAggregatedSignature(
	keyHandle int32,
	messageHandle int32,
	sigHandle int32,
) int32 {
	host := context.GetVMHost()
	return ManagedVerifyBLSWithHost(host, keyHandle, messageHandle, sigHandle, verifyBLSAggregatedSignature)
}

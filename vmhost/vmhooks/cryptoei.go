package vmhooks

import (
	"crypto/elliptic"
	"fmt"

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
const maxBLSAggregateSigners = 400

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
	resultOffset executor.MemPtr) (int32, error) {

	crypto := context.GetCryptoContext()
	enableEpochsHandler := context.host.EnableEpochsHandler()
	metering := context.GetMeteringContext()

	memLoadGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(metering.GasSchedule().CryptoAPICost.SHA256, memLoadGas)
	err := metering.UseGasBoundedAndAddTracedGas(sha256Name, gasToUse)
	if err != nil {
		return 1, err
	}

	data, err := context.MemLoad(dataOffset, length)
	if err != nil {
		return 1, err
	}

	result, err := crypto.Sha256(data)
	if err != nil {

		if enableEpochsHandler.IsFlagEnabled(vmhost.MaskInternalDependenciesErrorsFlag) {
			err = vmhost.ErrSha256Hash
		}

		return 1, err
	}

	err = context.MemStore(resultOffset, result)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

func getBytesAndConsumeGas(managedType vmhost.ManagedTypesContext, handle int32) ([]byte, error) {
	data, err := managedType.GetBytes(handle)
	if err != nil {
		return nil, err
	}
	err = managedType.ConsumeGasForBytes(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getVerifyBuffers(managedType vmhost.ManagedTypesContext, keyHandle, messageHandle, sigHandle int32) ([]byte, []byte, []byte, error) {
	key, err := getBytesAndConsumeGas(managedType, keyHandle)
	if err != nil {
		return nil, nil, nil, err
	}
	msg, err := getBytesAndConsumeGas(managedType, messageHandle)
	if err != nil {
		return nil, nil, nil, err
	}
	sig, err := getBytesAndConsumeGas(managedType, sigHandle)
	if err != nil {
		return nil, nil, nil, err
	}
	return key, msg, sig, nil
}

// ManagedSha256 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedSha256(inputHandle, outputHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	crypto := context.GetCryptoContext()
	enableEpochsHandler := context.host.EnableEpochsHandler()
	metering := context.GetMeteringContext()

	err := metering.UseGasBoundedAndAddTracedGas(sha256Name, metering.GasSchedule().CryptoAPICost.SHA256)
	if err != nil {
		return 1, err
	}

	inputBytes, err := getBytesAndConsumeGas(managedType, inputHandle)
	if err != nil {
		return 1, err
	}

	resultBytes, err := crypto.Sha256(inputBytes)
	if err != nil {
		if enableEpochsHandler.IsFlagEnabled(vmhost.MaskInternalDependenciesErrorsFlag) {
			err = vmhost.ErrSha256Hash
		}

		return 1, err
	}

	managedType.SetBytes(outputHandle, resultBytes)

	return 0, nil
}

// Keccak256 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) Keccak256(dataOffset executor.MemPtr, length executor.MemLength, resultOffset executor.MemPtr) (int32, error) {
	crypto := context.GetCryptoContext()
	enableEpochsHandler := context.host.EnableEpochsHandler()
	metering := context.GetMeteringContext()

	memLoadGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(metering.GasSchedule().CryptoAPICost.Keccak256, memLoadGas)
	err := metering.UseGasBoundedAndAddTracedGas(keccak256Name, gasToUse)
	if err != nil {
		return 1, err
	}

	data, err := context.MemLoad(dataOffset, length)
	if err != nil {
		return 1, err
	}

	result, err := crypto.Keccak256(data)
	if err != nil {
		if enableEpochsHandler.IsFlagEnabled(vmhost.MaskInternalDependenciesErrorsFlag) {
			err = vmhost.ErrKeccak256Hash
		}

		return 1, err
	}

	err = context.MemStore(resultOffset, result)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

// ManagedKeccak256 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedKeccak256(inputHandle, outputHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	crypto := context.GetCryptoContext()
	enableEpochsHandler := context.host.EnableEpochsHandler()
	metering := context.GetMeteringContext()

	err := metering.UseGasBoundedAndAddTracedGas(keccak256Name, metering.GasSchedule().CryptoAPICost.Keccak256)
	if err != nil {
		return 1, err
	}

	inputBytes, err := getBytesAndConsumeGas(managedType, inputHandle)
	if err != nil {
		return 1, err
	}

	resultBytes, err := crypto.Keccak256(inputBytes)
	if err != nil {
		if enableEpochsHandler.IsFlagEnabled(vmhost.MaskInternalDependenciesErrorsFlag) {
			err = vmhost.ErrKeccak256Hash
		}

		return 1, err
	}

	managedType.SetBytes(outputHandle, resultBytes)

	return 0, nil
}

// Ripemd160 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) Ripemd160(dataOffset executor.MemPtr, length executor.MemLength, resultOffset executor.MemPtr) (int32, error) {
	crypto := context.GetCryptoContext()
	enableEpochsHandler := context.host.EnableEpochsHandler()
	metering := context.GetMeteringContext()

	memLoadGas := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(length))
	gasToUse := math.AddUint64(metering.GasSchedule().CryptoAPICost.Ripemd160, memLoadGas)
	err := metering.UseGasBoundedAndAddTracedGas(ripemd160Name, gasToUse)
	if err != nil {
		return 1, err
	}

	data, err := context.MemLoad(dataOffset, length)
	if err != nil {
		return 1, err
	}

	result, err := crypto.Ripemd160(data)
	if err != nil {
		if enableEpochsHandler.IsFlagEnabled(vmhost.MaskInternalDependenciesErrorsFlag) {
			err = vmhost.ErrRipemd160Hash
		}

		return 1, err
	}

	err = context.MemStore(resultOffset, result)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

// ManagedRipemd160 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedRipemd160(inputHandle int32, outputHandle int32) (int32, error) {
	host := context.GetVMHost()
	return ManagedRipemd160WithHost(host, inputHandle, outputHandle)
}

// ManagedRipemd160WithHost VMHooks implementation.
func ManagedRipemd160WithHost(host vmhost.VMHost, inputHandle int32, outputHandle int32) (int32, error) {
	metering := host.Metering()
	managedType := host.ManagedTypes()
	crypto := host.Crypto()
	enableEpochsHandler := host.EnableEpochsHandler()

	err := metering.UseGasBoundedAndAddTracedGas(ripemd160Name, metering.GasSchedule().CryptoAPICost.Ripemd160)
	if err != nil {
		return 1, err
	}

	inputBytes, err := getBytesAndConsumeGas(managedType, inputHandle)
	if err != nil {
		return 1, err
	}

	result, err := crypto.Ripemd160(inputBytes)
	if err != nil {
		if enableEpochsHandler.IsFlagEnabled(vmhost.MaskInternalDependenciesErrorsFlag) {
			err = vmhost.ErrRipemd160Hash
		}

		return 1, err
	}

	managedType.SetBytes(outputHandle, result)

	return 0, nil
}

// VerifyBLS VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) VerifyBLS(
	keyOffset executor.MemPtr,
	messageOffset executor.MemPtr,
	messageLength executor.MemLength,
	sigOffset executor.MemPtr,
) (int32, error) {
	crypto := context.GetCryptoContext()
	enableEpochsHandler := context.host.EnableEpochsHandler()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(verifyBLSName)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifyBLS
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	key, err := context.MemLoad(keyOffset, blsPublicKeyLength)
	if err != nil {
		return 1, err
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(messageLength))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	message, err := context.MemLoad(messageOffset, messageLength)
	if err != nil {
		return 1, err
	}

	sig, err := context.MemLoad(sigOffset, blsSignatureLength)
	if err != nil {
		return 1, err
	}

	invalidSigErr := crypto.VerifyBLS(key, message, sig)
	if invalidSigErr != nil {
		if enableEpochsHandler.IsFlagEnabled(vmhost.MaskInternalDependenciesErrorsFlag) {
			invalidSigErr = vmhost.ErrBlsVerify
		}

		return -1, invalidSigErr
	}

	return 0, nil
}

// ManagedVerifyBLS VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedVerifyBLS(
	keyHandle int32,
	messageHandle int32,
	sigHandle int32,
) (int32, error) {
	host := context.GetVMHost()
	return ManagedVerifyBLSWithHost(host, keyHandle, messageHandle, sigHandle, verifyBLSName)
}

func useGasForCryptoVerify(
	metering vmhost.MeteringContext,
	sigVerificationType string,
) error {
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

	return metering.UseGasBounded(gasToUse)
}

// ManagedVerifyBLSWithHost VMHooks implementation.
func ManagedVerifyBLSWithHost(
	host vmhost.VMHost,
	keyHandle int32,
	messageHandle int32,
	sigHandle int32,
	sigVerificationType string,
) (int32, error) {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	crypto := host.Crypto()
	enableEpochsHandler := host.EnableEpochsHandler()
	err := useGasForCryptoVerify(metering, sigVerificationType)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return 1, err
	}

	keyBytes, msgBytes, sigBytes, err := getVerifyBuffers(managedType, keyHandle, messageHandle, sigHandle)
	if err != nil {
		return 1, err
	}

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

		if host.Runtime().FixCorrectSmallIntAndGasLimits() {
			if len(pubKeyBytes) == 0 {
				return 1, vmhost.ErrInvalidArgument
			}
			if len(pubKeyBytes) > maxBLSAggregateSigners {
				return 1, vmhost.ErrArgOutOfRange
			}
			for _, pubKey := range pubKeyBytes {
				if len(pubKey) != blsPublicKeyLength {
					return 1, vmhost.ErrLengthOfBufferNotCorrect
				}
			}

			if len(pubKeyBytes) > 1 {
				additionalSigners := uint64(len(pubKeyBytes) - 1)
				perSignerCost := metering.GasSchedule().CryptoAPICost.VerifyBLSSignatureShare
				if perSignerCost == 0 {
					perSignerCost = metering.GasSchedule().CryptoAPICost.VerifyBLS
				}
				cryptoGasToUse := math.MulUint64(perSignerCost, additionalSigners)
				err = metering.UseGasBounded(cryptoGasToUse)
				if err != nil {
					return 1, err
				}
			}
		}

		invalidSigErr = crypto.VerifyAggregatedSig(pubKeyBytes, msgBytes, sigBytes)
	}

	if invalidSigErr != nil {
		if enableEpochsHandler.IsFlagEnabled(vmhost.MaskInternalDependenciesErrorsFlag) {
			invalidSigErr = vmhost.ErrBlsVerify
		}

		return -1, invalidSigErr
	}

	return 0, nil
}

// VerifyEd25519 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) VerifyEd25519(
	keyOffset executor.MemPtr,
	messageOffset executor.MemPtr,
	messageLength executor.MemLength,
	sigOffset executor.MemPtr,
) (int32, error) {
	crypto := context.GetCryptoContext()
	enableEpochsHandler := context.host.EnableEpochsHandler()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(verifyEd25519Name)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifyEd25519
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	key, err := context.MemLoad(keyOffset, ed25519PublicKeyLength)
	if err != nil {
		return 1, err
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(messageLength))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	message, err := context.MemLoad(messageOffset, messageLength)
	if err != nil {
		return 1, err
	}

	sig, err := context.MemLoad(sigOffset, ed25519SignatureLength)
	if err != nil {
		return 1, err
	}

	invalidSigErr := crypto.VerifyEd25519(key, message, sig)
	if invalidSigErr != nil {
		if enableEpochsHandler.IsFlagEnabled(vmhost.MaskInternalDependenciesErrorsFlag) {
			invalidSigErr = vmhost.ErrEd25519Verify
		}

		return -1, invalidSigErr
	}

	return 0, nil
}

// ManagedVerifyEd25519 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedVerifyEd25519(
	keyHandle, messageHandle, sigHandle int32,
) (int32, error) {
	host := context.GetVMHost()
	return ManagedVerifyEd25519WithHost(host, keyHandle, messageHandle, sigHandle)
}

// ManagedVerifyEd25519WithHost VMHooks implementation.
func ManagedVerifyEd25519WithHost(
	host vmhost.VMHost,
	keyHandle, messageHandle, sigHandle int32,
) (int32, error) {
	metering := host.Metering()
	managedType := host.ManagedTypes()
	enableEpochsHandler := host.EnableEpochsHandler()
	crypto := host.Crypto()
	metering.StartGasTracing(verifyEd25519Name)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifyEd25519
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	keyBytes, msgBytes, sigBytes, err := getVerifyBuffers(managedType, keyHandle, messageHandle, sigHandle)
	if err != nil {
		return 1, err
	}

	invalidSigErr := crypto.VerifyEd25519(keyBytes, msgBytes, sigBytes)
	if invalidSigErr != nil {
		if enableEpochsHandler.IsFlagEnabled(vmhost.MaskInternalDependenciesErrorsFlag) {
			invalidSigErr = vmhost.ErrEd25519Verify
		}

		return -1, invalidSigErr
	}

	return 0, nil
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
) (int32, error) {
	crypto := context.GetCryptoContext()
	enableEpochsHandler := context.host.EnableEpochsHandler()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(verifyCustomSecp256k1Name)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifySecp256k1
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	if keyLength != secp256k1CompressedPublicKeyLength && keyLength != secp256k1UncompressedPublicKeyLength {
		return 1, vmhost.ErrInvalidPublicKeySize
	}

	key, err := context.MemLoad(keyOffset, keyLength)
	if err != nil {
		return 1, err
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(messageLength))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	message, err := context.MemLoad(messageOffset, messageLength)
	if err != nil {
		return 1, err
	}

	// read the 2 leading bytes first
	// byte1: 0x30, header
	// byte2: the remaining buffer length
	const sigHeaderLength = 2
	sigHeader, err := context.MemLoad(sigOffset, sigHeaderLength)
	if err != nil {
		return 1, err
	}
	if len(sigHeader) < sigHeaderLength {
		return 1, fmt.Errorf("sigHeader: %w", executor.ErrMemoryBadBounds)
	}
	sigLength := int32(sigHeader[1]) + sigHeaderLength
	sig, err := context.MemLoad(sigOffset, sigLength)
	if err != nil {
		return 1, err
	}

	invalidSigErr := crypto.VerifySecp256k1(key, message, sig, uint8(hashType))
	if invalidSigErr != nil {
		if enableEpochsHandler.IsFlagEnabled(vmhost.MaskInternalDependenciesErrorsFlag) {
			invalidSigErr = vmhost.ErrSecp256k1Verify
		}

		return -1, invalidSigErr
	}

	return 0, nil
}

// ManagedVerifyCustomSecp256k1 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedVerifyCustomSecp256k1(
	keyHandle, messageHandle, sigHandle int32,
	hashType int32,
) (int32, error) {
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
) (int32, error) {
	runtime := host.Runtime()
	enableEpochsHandler := host.EnableEpochsHandler()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	crypto := host.Crypto()

	err := useGasForCryptoVerify(metering, verifyCryptoFunc)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return 1, err
	}

	keyBytes, msgBytes, sigBytes, err := getVerifyBuffers(managedType, keyHandle, messageHandle, sigHandle)
	if err != nil {
		return 1, err
	}

	invalidSigErr := vmhost.ErrInvalidArgument
	switch verifyCryptoFunc {
	case verifyCustomSecp256k1Name:
		invalidSigErr = crypto.VerifySecp256k1(keyBytes, msgBytes, sigBytes, uint8(hashType))
	case verifySecp256R1Signature:
		invalidSigErr = crypto.VerifySecp256r1(keyBytes, msgBytes, sigBytes)
	}

	if invalidSigErr != nil {
		if enableEpochsHandler.IsFlagEnabled(vmhost.MaskInternalDependenciesErrorsFlag) {
			invalidSigErr = vmhost.ErrSecp256k1Verify
		}

		return -1, invalidSigErr
	}

	return 0, nil
}

// VerifySecp256k1 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) VerifySecp256k1(
	keyOffset executor.MemPtr,
	keyLength executor.MemLength,
	messageOffset executor.MemPtr,
	messageLength executor.MemLength,
	sigOffset executor.MemPtr,
) (int32, error) {
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
) (int32, error) {
	host := context.GetVMHost()
	return ManagedVerifySecp256k1WithHost(host, keyHandle, messageHandle, sigHandle)
}

// ManagedVerifySecp256k1WithHost VMHooks implementation.
func ManagedVerifySecp256k1WithHost(
	host vmhost.VMHost,
	keyHandle, messageHandle, sigHandle int32,
) (int32, error) {
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
) (int32, error) {
	crypto := context.GetCryptoContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().CryptoAPICost.EncodeDERSig
	err := metering.UseGasBoundedAndAddTracedGas(encodeSecp256k1DerSignatureName, gasToUse)
	if err != nil {
		return 1, err
	}

	r, err := context.MemLoad(rOffset, rLength)
	if err != nil {
		return 1, err
	}

	s, err := context.MemLoad(sOffset, sLength)
	if err != nil {
		return 1, err
	}

	derSig := crypto.EncodeSecp256k1DERSignature(r, s)
	err = context.MemStore(sigOffset, derSig)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

// ManagedEncodeSecp256k1DerSignature VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedEncodeSecp256k1DerSignature(
	rHandle, sHandle, sigHandle int32,
) (int32, error) {
	host := context.GetVMHost()
	return ManagedEncodeSecp256k1DerSignatureWithHost(host, rHandle, sHandle, sigHandle)
}

// ManagedEncodeSecp256k1DerSignatureWithHost VMHooks implementation.
func ManagedEncodeSecp256k1DerSignatureWithHost(
	host vmhost.VMHost,
	rHandle, sHandle, sigHandle int32,
) (int32, error) {
	metering := host.Metering()
	managedType := host.ManagedTypes()
	crypto := host.Crypto()

	gasToUse := metering.GasSchedule().CryptoAPICost.EncodeDERSig
	err := metering.UseGasBoundedAndAddTracedGas(encodeSecp256k1DerSignatureName, gasToUse)
	if err != nil {
		return 1, err
	}

	r, err := managedType.GetBytes(rHandle)
	if err != nil {
		return 1, err
	}

	s, err := managedType.GetBytes(sHandle)
	if err != nil {
		return 1, err
	}

	derSig := crypto.EncodeSecp256k1DERSignature(r, s)
	managedType.SetBytes(sigHandle, derSig)

	return 0, nil
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
) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(addECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		return vmhost.ErrNoEllipticCurveUnderThisHandle
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.AddECC * uint64(curveMultiplier) / 100
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	// NOTE on the original fall-through: the legacy code called
	// context.FailExecution(err1) here WITHOUT returning, and kept using a
	// possibly-nil `ec` afterwards. That branch is unreachable in practice:
	// Get100xCurveGasCostMultiplier above reads the very same ecValues map
	// under the very same ecHandle and already returned early when the
	// handle is missing, so GetEllipticCurve below can never fail at this
	// point. Returning here is therefore an exact behavioural no-op; see the
	// migration report for the full reasoning.
	ec, err1 := managedType.GetEllipticCurve(ecHandle)
	if err1 != nil {
		return err1
	}

	xResult, yResult, err := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	if err != nil {
		return vmhost.ErrNoBigIntUnderThisHandle
	}
	x1, y1, err := managedType.GetTwoBigInt(fstPointXHandle, fstPointYHandle)
	if err != nil {
		return vmhost.ErrNoBigIntUnderThisHandle
	}
	x2, y2, err := managedType.GetTwoBigInt(sndPointXHandle, sndPointYHandle)
	if err != nil {
		return vmhost.ErrNoBigIntUnderThisHandle
	}

	if !ec.IsOnCurve(x1, y1) || !ec.IsOnCurve(x2, y2) {
		return vmhost.ErrPointNotOnCurve
	}

	err = managedType.ConsumeGasForBigIntCopy(xResult, yResult, ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x1, y1, x2, y2)
	if err != nil {
		return err
	}
	xResultAdd, yResultAdd := ec.Add(x1, y1, x2, y2)
	xResult.Set(xResultAdd)
	yResult.Set(yResultAdd)
	return nil
}

// DoubleEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) DoubleEC(
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	pointXHandle int32,
	pointYHandle int32,
) error {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(doubleECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		return vmhost.ErrNoEllipticCurveUnderThisHandle
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.DoubleECC * uint64(curveMultiplier) / 100
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	ec, err1 := managedType.GetEllipticCurve(ecHandle)
	if err1 != nil {
		return err1
	}

	xResult, yResult, err1 := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	x, y, err2 := managedType.GetTwoBigInt(pointXHandle, pointYHandle)
	if err1 != nil || err2 != nil {
		return vmhost.ErrNoBigIntUnderThisHandle
	}
	if !ec.IsOnCurve(x, y) {
		return vmhost.ErrPointNotOnCurve
	}

	err = managedType.ConsumeGasForBigIntCopy(xResult, yResult, ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	if err != nil {
		return err
	}

	xResultDouble, yResultDouble := ec.Double(x, y)
	xResult.Set(xResultDouble)
	yResult.Set(yResultDouble)
	return nil
}

// IsOnCurveEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) IsOnCurveEC(
	ecHandle int32,
	pointXHandle int32,
	pointYHandle int32,
) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(isOnCurveECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		return 1, vmhost.ErrNoEllipticCurveUnderThisHandle
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.IsOnCurveECC * uint64(curveMultiplier) / 100
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return -1, err
	}

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		return -1, err
	}

	x, y, err := managedType.GetTwoBigInt(pointXHandle, pointYHandle)
	if err != nil || x == nil || y == nil {
		return -1, vmhost.ErrNoBigIntUnderThisHandle
	}

	err = managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	if err != nil {
		return -1, err
	}

	if ec.IsOnCurve(x, y) {
		return 1, nil
	}

	return 0, nil
}

// ScalarBaseMultEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ScalarBaseMultEC(
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataOffset executor.MemPtr,
	length executor.MemLength,
) (int32, error) {
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()
	metering.StartGasTracing(scalarBaseMultECName)

	if length < 0 {
		return 1, vmhost.ErrNegativeLength
	}

	curveMultiplier := managedType.GetScalarMult100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		return 1, vmhost.ErrNoEllipticCurveUnderThisHandle
	}
	oneByteScalarGasCost := metering.GasSchedule().CryptoAPICost.ScalarMultECC * uint64(curveMultiplier) / 100
	gasToUse := oneByteScalarGasCost + uint64(length)*oneByteScalarGasCost
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return -1, err
	}

	data, err := context.MemLoad(dataOffset, length)
	if err != nil {
		return 1, err
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
) (int32, error) {
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
) (int32, error) {
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(scalarBaseMultECName)

	curveMultiplier := managedType.GetScalarMult100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		return 1, vmhost.ErrNoEllipticCurveUnderThisHandle
	}

	data, err := managedType.GetBytes(dataHandle)
	if err != nil {
		return 1, err
	}

	oneByteScalarGasCost := metering.GasSchedule().CryptoAPICost.ScalarMultECC * uint64(curveMultiplier) / 100
	gasToUse := oneByteScalarGasCost + uint64(len(data))*oneByteScalarGasCost
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return -1, err
	}

	return commonScalarBaseMultEC(host, xResultHandle, yResultHandle, ecHandle, data)
}

func commonScalarBaseMultEC(
	host vmhost.VMHost,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	data []byte,
) (int32, error) {
	managedType := host.ManagedTypes()

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		return 1, err
	}

	xResult, yResult, err := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	if err != nil {
		return 1, err
	}

	err = managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xResult, yResult)
	if err != nil {
		return 1, err
	}

	xResultSBM, yResultSBM := ec.ScalarBaseMult(data)
	if !ec.IsOnCurve(xResultSBM, yResultSBM) {
		return 1, vmhost.ErrPointNotOnCurve
	}
	xResult.Set(xResultSBM)
	yResult.Set(yResultSBM)

	return 0, nil
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
) (int32, error) {
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()
	metering.StartGasTracing(scalarMultECName)

	if length < 0 {
		return 1, vmhost.ErrNegativeLength
	}

	curveMultiplier := managedType.GetScalarMult100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		return 1, vmhost.ErrNoEllipticCurveUnderThisHandle
	}
	oneByteScalarGasCost := metering.GasSchedule().CryptoAPICost.ScalarMultECC * uint64(curveMultiplier) / 100
	gasToUse := oneByteScalarGasCost + uint64(length)*oneByteScalarGasCost
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return -1, err
	}

	data, err := context.MemLoad(dataOffset, length)
	if err != nil {
		return 1, err
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
) (int32, error) {
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
) (int32, error) {
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(scalarMultECName)

	curveMultiplier := managedType.GetScalarMult100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		return 1, vmhost.ErrNoEllipticCurveUnderThisHandle
	}

	data, err := managedType.GetBytes(dataHandle)
	if err != nil {
		return 1, err
	}

	oneByteScalarGasCost := metering.GasSchedule().CryptoAPICost.ScalarMultECC * uint64(curveMultiplier) / 100
	gasToUse := oneByteScalarGasCost + uint64(len(data))*oneByteScalarGasCost
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return -1, err
	}

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
) (int32, error) {
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(scalarMultECName)

	ec, err1 := managedType.GetEllipticCurve(ecHandle)
	if err1 != nil {
		return 1, err1
	}

	xResult, yResult, err1 := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	x, y, err2 := managedType.GetTwoBigInt(pointXHandle, pointYHandle)
	if err1 != nil || err2 != nil {
		return 1, vmhost.ErrNoBigIntUnderThisHandle
	}
	if !ec.IsOnCurve(x, y) {
		return 1, vmhost.ErrPointNotOnCurve
	}

	err := managedType.ConsumeGasForBigIntCopy(xResult, yResult, ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	if err != nil {
		return 1, err
	}

	xResultSM, yResultSM := ec.ScalarMult(x, y, data)
	if !ec.IsOnCurve(xResultSM, yResultSM) {
		return 1, vmhost.ErrPointNotOnCurve
	}
	xResult.Set(xResultSM)
	yResult.Set(yResultSM)

	return 0, nil
}

// MarshalEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MarshalEC(
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	resultOffset executor.MemPtr,
) (int32, error) {
	host := context.GetVMHost()
	result, err := commonMarshalEC(host, xPairHandle, yPairHandle, ecHandle)
	if err != nil {
		return -1, err
	}

	err = context.MemStore(resultOffset, result)
	if err != nil {
		return -1, err
	}
	return int32(len(result)), nil
}

// ManagedMarshalEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMarshalEC(
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	resultHandle int32,
) (int32, error) {
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
) (int32, error) {
	result, err := commonMarshalEC(host, xPairHandle, yPairHandle, ecHandle)
	if err != nil {
		return -1, err
	}

	managedType := host.ManagedTypes()
	managedType.SetBytes(resultHandle, result)
	return int32(len(result)), nil
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
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return nil, err
	}

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

	err = managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	if err != nil {
		return nil, err
	}

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
) (int32, error) {
	host := context.GetVMHost()
	result, err := commonMarshalCompressedEC(host, xPairHandle, yPairHandle, ecHandle)
	if err != nil {
		return -1, err
	}

	err = context.MemStore(resultOffset, result)
	if err != nil {
		return -1, err
	}
	return int32(len(result)), nil
}

// ManagedMarshalCompressedEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMarshalCompressedEC(
	xPairHandle int32,
	yPairHandle int32,
	ecHandle int32,
	resultHandle int32,
) (int32, error) {
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
) (int32, error) {
	managedType := host.ManagedTypes()
	result, err := commonMarshalCompressedEC(host, xPairHandle, yPairHandle, ecHandle)
	if err != nil {
		return -1, err
	}

	managedType.SetBytes(resultHandle, result)
	return int32(len(result)), nil
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
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return nil, err
	}

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

	err = managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, x, y)
	if err != nil {
		return nil, err
	}

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
) (int32, error) {
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()
	metering.StartGasTracing(unmarshalECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		return 1, vmhost.ErrNoEllipticCurveUnderThisHandle
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.UnmarshalECC * uint64(curveMultiplier) / 100
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return -1, err
	}

	data, err := context.MemLoad(dataOffset, length)
	if err != nil {
		return 1, err
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
) (int32, error) {
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
) (int32, error) {
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(unmarshalECName)

	curveMultiplier := managedType.Get100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		return 1, vmhost.ErrNoEllipticCurveUnderThisHandle
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.UnmarshalECC * uint64(curveMultiplier) / 100
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return -1, err
	}

	data, err := managedType.GetBytes(dataHandle)
	if err != nil {
		return 1, err
	}

	return commonUnmarshalEC(host, xResultHandle, yResultHandle, ecHandle, data)
}

func commonUnmarshalEC(
	host vmhost.VMHost,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	data []byte,
) (int32, error) {
	managedType := host.ManagedTypes()

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		return 1, err
	}
	byteLen := (ec.BitSize + 7) / 8
	if len(data) != 1+2*byteLen {
		return 1, vmhost.ErrLengthOfBufferNotCorrect
	}

	xResult, yResult, err := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	if err != nil {
		return 1, err
	}

	err = managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xResult, yResult)
	if err != nil {
		return 1, err
	}

	xResultU, yResultU := elliptic.Unmarshal(ec, data)
	if xResultU == nil || yResultU == nil || !ec.IsOnCurve(xResultU, yResultU) {
		return 1, vmhost.ErrPointNotOnCurve
	}
	xResult.Set(xResultU)
	yResult.Set(yResultU)

	return 0, nil
}

// UnmarshalCompressedEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) UnmarshalCompressedEC(
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	dataOffset executor.MemPtr,
	length executor.MemLength,
) (int32, error) {
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()
	metering.StartGasTracing(unmarshalCompressedECName)

	curveMultiplier := managedType.GetUCompressed100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		return 1, vmhost.ErrNoEllipticCurveUnderThisHandle
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.UnmarshalCompressedECC * uint64(curveMultiplier) / 100
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return -1, err
	}

	data, err := context.MemLoad(dataOffset, length)
	if err != nil {
		return int32(len(data)), err
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
) (int32, error) {
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
) (int32, error) {
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(unmarshalCompressedECName)

	curveMultiplier := managedType.GetUCompressed100xCurveGasCostMultiplier(ecHandle)
	if curveMultiplier < 0 {
		return 1, vmhost.ErrNoEllipticCurveUnderThisHandle
	}
	gasToUse := metering.GasSchedule().CryptoAPICost.UnmarshalCompressedECC * uint64(curveMultiplier) / 100
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return -1, err
	}

	data, err := managedType.GetBytes(dataHandle)
	if err != nil {
		return int32(len(data)), err
	}

	return commonUnmarshalCompressedEC(host, xResultHandle, yResultHandle, ecHandle, data)
}

func commonUnmarshalCompressedEC(
	host vmhost.VMHost,
	xResultHandle int32,
	yResultHandle int32,
	ecHandle int32,
	data []byte,
) (int32, error) {
	managedType := host.ManagedTypes()

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		return 1, err
	}
	byteLen := (ec.BitSize+7)/8 + 1
	if len(data) != byteLen {
		return 1, vmhost.ErrLengthOfBufferNotCorrect
	}

	xResult, yResult, err := managedType.GetTwoBigInt(xResultHandle, yResultHandle)
	if err != nil {
		return 1, err
	}

	err = managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xResult, yResult)
	if err != nil {
		return 1, err
	}

	xResultUC, yResultUC := elliptic.UnmarshalCompressed(ec, data)
	if xResultUC == nil || yResultUC == nil || !ec.IsOnCurve(xResultUC, yResultUC) {
		return 1, vmhost.ErrPointNotOnCurve
	}
	xResult.Set(xResultUC)
	yResult.Set(yResultUC)
	return 0, nil
}

// GenerateKeyEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GenerateKeyEC(
	xPubKeyHandle int32,
	yPubKeyHandle int32,
	ecHandle int32,
	resultOffset executor.MemPtr,
) (int32, error) {
	host := context.GetVMHost()
	result, err := commonGenerateEC(host, xPubKeyHandle, yPubKeyHandle, ecHandle)
	if err != nil {
		return 1, err
	}

	err = context.MemStore(resultOffset, result)
	if err != nil {
		return int32(len(result)), err
	}

	return 0, nil
}

// ManagedGenerateKeyEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGenerateKeyEC(
	xPubKeyHandle int32,
	yPubKeyHandle int32,
	ecHandle int32,
	resultHandle int32,
) (int32, error) {
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
) (int32, error) {
	managedType := host.ManagedTypes()
	result, err := commonGenerateEC(host, xPubKeyHandle, yPubKeyHandle, ecHandle)
	if err != nil {
		return 1, err
	}

	managedType.SetBytes(resultHandle, result)
	return 0, nil
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
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return nil, err
	}

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		return nil, err
	}

	xPubKey, yPubKey, err := managedType.GetTwoBigInt(xPubKeyHandle, yPubKeyHandle)
	if err != nil {
		return nil, err
	}

	err = managedType.ConsumeGasForBigIntCopy(ec.P, ec.N, ec.B, ec.Gx, ec.Gy, xPubKey, yPubKey)
	if err != nil {
		return nil, err
	}

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
func (context *VMHooksImpl) CreateEC(dataOffset executor.MemPtr, dataLength executor.MemLength) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().CryptoAPICost.EllipticCurveNew
	err := metering.UseGasBoundedAndAddTracedGas(createECName, gasToUse)
	if err != nil {
		return -1, err
	}

	if dataLength != curveNameLength {
		return -1, vmhost.ErrBadBounds
	}
	data, err := context.MemLoad(dataOffset, dataLength)
	if err != nil {
		return -1, err
	}
	curveChoice := string(data[:])
	switch curveChoice {
	case "p224":
		curveParams := elliptic.P224().Params()
		return managedType.PutEllipticCurve(curveParams), nil
	case "p256":
		curveParams := elliptic.P256().Params()
		return managedType.PutEllipticCurve(curveParams), nil
	case "p384":
		curveParams := elliptic.P384().Params()
		return managedType.PutEllipticCurve(curveParams), nil
	case "p521":
		curveParams := elliptic.P521().Params()
		return managedType.PutEllipticCurve(curveParams), nil
	}
	if context.host.Runtime().FixCorrectSmallIntAndGasLimits() {
		return -1, vmhost.ErrBadBounds
	}
	return -1, nil
}

// ManagedCreateEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedCreateEC(dataHandle int32) (int32, error) {
	host := context.GetVMHost()
	return ManagedCreateECWithHost(host, dataHandle)
}

// ManagedCreateECWithHost VMHooks implementation.
func ManagedCreateECWithHost(host vmhost.VMHost, dataHandle int32) (int32, error) {
	metering := host.Metering()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().CryptoAPICost.EllipticCurveNew
	err := metering.UseGasBoundedAndAddTracedGas(createECName, gasToUse)
	if err != nil {
		return -1, err
	}

	data, err := managedType.GetBytes(dataHandle)
	if err != nil {
		return -1, err
	}
	curveChoice := string(data[:])
	switch curveChoice {
	case "p224":
		curveParams := elliptic.P224().Params()
		return managedType.PutEllipticCurve(curveParams), nil
	case "p256":
		curveParams := elliptic.P256().Params()
		return managedType.PutEllipticCurve(curveParams), nil
	case "p384":
		curveParams := elliptic.P384().Params()
		return managedType.PutEllipticCurve(curveParams), nil
	case "p521":
		curveParams := elliptic.P521().Params()
		return managedType.PutEllipticCurve(curveParams), nil
	}

	return -1, vmhost.ErrBadBounds
}

// GetCurveLengthEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetCurveLengthEC(ecHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetInt64
	err := metering.UseGasBoundedAndAddTracedGas(getCurveLengthECName, gasToUse)
	if err != nil {
		return -1, err
	}

	ecLength := managedType.GetEllipticCurveSizeOfField(ecHandle)
	if ecLength == -1 {
		return -1, vmhost.ErrNoEllipticCurveUnderThisHandle
	}

	return ecLength, nil
}

// GetPrivKeyByteLengthEC VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetPrivKeyByteLengthEC(ecHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetInt64
	err := metering.UseGasBoundedAndAddTracedGas(getPrivKeyByteLengthECName, gasToUse)
	if err != nil {
		return -1, err
	}

	byteLength := managedType.GetPrivateKeyByteLengthEC(ecHandle)
	if byteLength == -1 {
		return -1, vmhost.ErrNoEllipticCurveUnderThisHandle
	}

	return byteLength, nil
}

// EllipticCurveGetValues VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) EllipticCurveGetValues(ecHandle int32, fieldOrderHandle int32, basePointOrderHandle int32, eqConstantHandle int32, xBasePointHandle int32, yBasePointHandle int32) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BigIntAPICost.BigIntGetInt64 * 5
	err := metering.UseGasBoundedAndAddTracedGas(ellipticCurveGetValuesName, gasToUse)
	if err != nil {
		return -1, err
	}

	ec, err := managedType.GetEllipticCurve(ecHandle)
	if err != nil {
		return -1, err
	}
	fieldOrder, basePointOrder, err := managedType.GetTwoBigInt(fieldOrderHandle, basePointOrderHandle)
	if err != nil {
		return -1, err
	}
	eqConstant, err := managedType.GetBigInt(eqConstantHandle)
	if err != nil {
		return -1, err
	}
	xBasePoint, yBasePoint, err := managedType.GetTwoBigInt(xBasePointHandle, yBasePointHandle)
	if err != nil {
		return -1, err
	}
	fieldOrder.Set(ec.P)
	basePointOrder.Set(ec.N)
	eqConstant.Set(ec.B)
	xBasePoint.Set(ec.Gx)
	yBasePoint.Set(ec.Gy)
	return ecHandle, nil
}

// ManagedVerifySecp256r1 VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedVerifySecp256r1(
	keyHandle, messageHandle, sigHandle int32,
) (int32, error) {
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
) (int32, error) {
	host := context.GetVMHost()
	return ManagedVerifyBLSWithHost(host, keyHandle, messageHandle, sigHandle, verifyBLSSignatureShare)
}

// ManagedVerifyBLSAggregatedSignature VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedVerifyBLSAggregatedSignature(
	keyHandle int32,
	messageHandle int32,
	sigHandle int32,
) (int32, error) {
	host := context.GetVMHost()
	return ManagedVerifyBLSWithHost(host, keyHandle, messageHandle, sigHandle, verifyBLSAggregatedSignature)
}

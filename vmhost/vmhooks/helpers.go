package vmhooks

import (
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

import (
	"github.com/multiversx/mx-chain-core-go/data/esdt"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

// This file will contain helper functions to reduce boilerplate and duplication in the other files.

type hashFunc func(data []byte) ([]byte, error)

func (context *VMHooksImpl) managedHash(
	inputHandle int32,
	outputHandle int32,
	traceName string,
	gasCost uint64,
	hf hashFunc,
	failError error,
) int32 {
	host := context.GetVMHost()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	enableEpochsHandler := host.EnableEpochsHandler()

	err := metering.UseGasBoundedAndAddTracedGas(traceName, gasCost)
	if err != nil {
		FailExecution(host, err)
		return 1
	}

	inputBytes, err := managedType.GetBytes(inputHandle)
	if err != nil {
		FailExecution(host, err)
		return 1
	}

	err = managedType.ConsumeGasForBytes(inputBytes)
	if err != nil {
		FailExecution(host, err)
		return 1
	}

	resultBytes, err := hf(inputBytes)
	if err != nil {
		if enableEpochsHandler.IsFlagEnabled(vmhost.MaskInternalDependenciesErrorsFlag) {
			err = failError
		}

		FailExecution(host, err)
		return 1
	}

	managedType.SetBytes(outputHandle, resultBytes)

	return 0
}

type esdtDataHandler func(context *VMHooksImpl, esdtData *esdt.ESDigitalToken) int32

func (context *VMHooksImpl) withESDTData(
	addressOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	nonce int64,
	traceName string,
	handler esdtDataHandler,
) int32 {
	metering := context.GetMeteringContext()
	metering.StartGasTracing(traceName)

	esdtData, err := context.GetESDTDataFromBlockchainHook(addressOffset, tokenIDOffset, tokenIDLen, nonce)

	if err != nil {
		context.FailExecution(err)
		return -1
	}

	return handler(context, esdtData)
}

func (context *VMHooksImpl) GetESDTDataFromBlockchainHook(
	addressOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	nonce int64,
) (*esdt.ESDigitalToken, error) {
	metering := context.GetMeteringContext()
	blockchain := context.GetBlockchainContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return nil, err
	}

	address, err := context.MemLoad(addressOffset, vmhost.AddressLen)
	if err != nil {
		return nil, err
	}

	tokenID, err := context.MemLoad(tokenIDOffset, tokenIDLen)
	if err != nil {
		return nil, err
	}

	esdtToken, err := blockchain.GetESDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		return nil, err
	}

	return esdtToken, nil
}

func (context *VMHooksImpl) getSignatureOperands(keyHandle, messageHandle, sigHandle int32) ([]byte, []byte, []byte, error) {
	managedType := context.GetManagedTypesContext()

	keyBytes, err := managedType.GetBytes(keyHandle)
	if err != nil {
		return nil, nil, nil, err
	}
	if err = managedType.ConsumeGasForBytes(keyBytes); err != nil {
		return nil, nil, nil, err
	}

	msgBytes, err := managedType.GetBytes(messageHandle)
	if err != nil {
		return nil, nil, nil, err
	}
	if err = managedType.ConsumeGasForBytes(msgBytes); err != nil {
		return nil, nil, nil, err
	}

	sigBytes, err := managedType.GetBytes(sigHandle)
	if err != nil {
		return nil, nil, nil, err
	}
	if err = managedType.ConsumeGasForBytes(sigBytes); err != nil {
		return nil, nil, nil, err
	}

	return keyBytes, msgBytes, sigBytes, nil
}

func (context *VMHooksImpl) managedVerifyWithOperands(
	sigVerificationType string,
	failError error,
	verify func() error,
) int32 {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	enableEpochsHandler := host.EnableEpochsHandler()

	err := useGasForCryptoVerify(metering, sigVerificationType)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		FailExecution(host, err)
		return 1
	}

	invalidSigErr := verify()
	if invalidSigErr != nil {
		if enableEpochsHandler.IsFlagEnabled(vmhost.MaskInternalDependenciesErrorsFlag) {
			invalidSigErr = failError
		}

		FailExecution(host, invalidSigErr)
		return -1
	}

	return 0
}

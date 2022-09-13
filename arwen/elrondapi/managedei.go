package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
// typedef unsigned long long uint64_t;
//
// extern void	v1_5_managedSCAddress(void *context, int32_t addressHandle);
// extern void	v1_5_managedOwnerAddress(void *context, int32_t addressHandle);
// extern void	v1_5_managedCaller(void *context, int32_t addressHandle);
// extern void	v1_5_managedSignalError(void* context, int32_t errHandle1);
// extern void	v1_5_managedWriteLog(void* context, int32_t topicsHandle, int32_t dataHandle);
//
// extern int32_t	v1_5_managedMultiTransferESDTNFTExecute(void *context, int32_t dstHandle, int32_t tokenTransfersHandle, long long gasLimit, int32_t functionHandle, int32_t argumentsHandle);
// extern int32_t	v1_5_managedTransferValueExecute(void *context, int32_t dstHandle, int32_t valueHandle, long long gasLimit, int32_t functionHandle, int32_t argumentsHandle);
// extern int32_t	v1_5_managedExecuteOnDestContext(void *context, long long gas, int32_t addressHandle, int32_t valueHandle, int32_t functionHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern int32_t	v1_5_managedExecuteOnSameContext(void *context, long long gas, int32_t addressHandle, int32_t valueHandle, int32_t functionHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern int32_t	v1_5_managedExecuteReadOnly(void *context, long long gas, int32_t addressHandle, int32_t functionHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern int32_t	v1_5_managedCreateContract(void *context, long long gas, int32_t valueHandle, int32_t codeHandle, int32_t codeMetadataHandle, int32_t argumentsHandle, int32_t resultAddressHandle, int32_t resultHandle);
// extern int32_t	v1_5_managedDeployFromSourceContract(void *context, long long gas, int32_t valueHandle, int32_t addressHandle, int32_t codeMetadataHandle, int32_t argumentsHandle, int32_t resultAddressHandle, int32_t resultHandle);
// extern void		v1_5_managedUpgradeContract(void *context, int32_t dstHandle, long long gas, int32_t valueHandle, int32_t codeHandle, int32_t codeMetadataHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern void		v1_5_managedUpgradeFromSourceContract(void *context, int32_t dstHandle, long long gas, int32_t valueHandle, int32_t addressHandle, int32_t codeMetadataHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern void		v1_5_managedDeleteContract(void *context, int32_t dstHandle, long long gas, int32_t argsHandle);
// extern void		v1_5_managedAsyncCall(void *context, int32_t dstHandle, int32_t valueHandle, int32_t functionHandle, int32_t argumentsHandle);
// extern int32_t	v1_5_managedCreateAsyncCall(void *context, int32_t destHandle, int32_t valueHandle, int32_t functionHandle, int32_t argumentsHandle, int32_t successCallback, int32_t successLength, int32_t errorCallback, int32_t errorLength, long long gas, long long extraGasForCallback, int32_t callbackClosureHandle);
// extern void		v1_5_managedGetCallbackClosure(void *context, int32_t callbackClosureHandle);
//
// extern void		v1_5_managedGetMultiESDTCallValue(void *context, int32_t multiCallValueHandle);
// extern void		v1_5_managedGetESDTBalance(void *context, int32_t addressHandle, int32_t tokenIDHandle, long long nonce, int32_t valueHandle);
// extern void		v1_5_managedGetESDTTokenData(void *context, int32_t addressHandle, int32_t tokenIDHandle, long long nonce, int32_t valueHandle, int32_t propertiesHandle, int32_t hashHandle, int32_t nameHandle, int32_t attributesHandle, int32_t creatorHandle, int32_t royaltiesHandle, int32_t urisHandle);
//
// extern void		v1_5_managedGetReturnData(void *context, int32_t resultID, int32_t resultHandle);
// extern void		v1_5_managedGetPrevBlockRandomSeed(void *context, int32_t resultHandle);
// extern void		v1_5_managedGetBlockRandomSeed(void *context, int32_t resultHandle);
// extern void		v1_5_managedGetStateRootHash(void *context, int32_t resultHandle);
// extern void		v1_5_managedGetOriginalTxHash(void *context, int32_t resultHandle);
//
// extern int32_t   v1_5_managedIsESDTFrozen(void *context, int32_t addressHandle, int32_t tokenIDHandle, long long nonce);
// extern int32_t   v1_5_managedIsESDTPaused(void *context, int32_t tokenIDHandle);
// extern int32_t   v1_5_managedIsESDTLimitedTransfer(void *context, int32_t tokenIDHandle);
// extern void      v1_5_managedBufferToHex(void *context, int32_t sourceHandle, int32_t destHandle);
import "C"

import (
	"encoding/hex"
	"errors"
	"unsafe"

	"github.com/ElrondNetwork/elrond-vm-common/builtInFunctions"

	"github.com/ElrondNetwork/wasm-vm/arwen"
	"github.com/ElrondNetwork/wasm-vm/math"
	"github.com/ElrondNetwork/wasm-vm/wasmer"
)

const (
	managedSCAddressName                    = "managedSCAddress"
	managedOwnerAddressName                 = "managedOwnerAddress"
	managedCallerName                       = "managedCaller"
	managedSignalErrorName                  = "managedSignalError"
	managedWriteLogName                     = "managedWriteLog"
	managedMultiTransferESDTNFTExecuteName  = "managedMultiTransferESDTNFTExecute"
	managedTransferValueExecuteName         = "managedTransferValueExecute"
	managedExecuteOnDestContextName         = "managedExecuteOnDestContext"
	managedExecuteOnDestContextByCallerName = "managedExecuteOnDestContextByCaller"
	managedExecuteOnSameContextName         = "managedExecuteOnSameContext"
	managedExecuteReadOnlyName              = "managedExecuteReadOnly"
	managedCreateContractName               = "managedCreateContract"
	managedDeployFromSourceContractName     = "managedDeployFromSourceContract"
	managedUpgradeContractName              = "managedUpgradeContract"
	managedUpgradeFromSourceContractName    = "managedUpgradeFromSourceContract"
	managedAsyncCallName                    = "managedAsyncCall"
	managedCreateAsyncCallName              = "managedCreateAsyncCall"
	managedGetCallbackClosure               = "managedGetCallbackClosure"
	managedGetMultiESDTCallValueName        = "managedGetMultiESDTCallValue"
	managedGetESDTBalanceName               = "managedGetESDTBalance"
	managedGetESDTTokenDataName             = "managedGetESDTTokenData"
	managedGetReturnDataName                = "managedGetReturnData"
	managedGetPrevBlockRandomSeedName       = "managedGetPrevBlockRandomSeed"
	managedGetBlockRandomSeedName           = "managedGetBlockRandomSeed"
	managedGetStateRootHashName             = "managedGetStateRootHash"
	managedGetOriginalTxHashName            = "managedGetOriginalTxHash"
	managedIsESDTFrozenName                 = "managedIsESDTFrozen"
	managedIsESDTLimitedTransferName        = "managedIsESDTLimitedTransfer"
	managedIsESDTPausedName                 = "managedIsESDTPaused"
	managedBufferToHexName                  = "managedBufferToHex"
)

// ManagedEIImports creates a new wasmer.Imports populated with variants of the API methods that use managed types only.
func ManagedEIImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("env")

	imports, err := imports.Append("managedSCAddress", v1_5_managedSCAddress, C.v1_5_managedSCAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedOwnerAddress", v1_5_managedOwnerAddress, C.v1_5_managedOwnerAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedCaller", v1_5_managedCaller, C.v1_5_managedCaller)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedSignalError", v1_5_managedSignalError, C.v1_5_managedSignalError)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedWriteLog", v1_5_managedWriteLog, C.v1_5_managedWriteLog)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedMultiTransferESDTNFTExecute", v1_5_managedMultiTransferESDTNFTExecute, C.v1_5_managedMultiTransferESDTNFTExecute)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedTransferValueExecute", v1_5_managedTransferValueExecute, C.v1_5_managedTransferValueExecute)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedExecuteOnDestContext", v1_5_managedExecuteOnDestContext, C.v1_5_managedExecuteOnDestContext)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedExecuteOnSameContext", v1_5_managedExecuteOnSameContext, C.v1_5_managedExecuteOnSameContext)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedExecuteReadOnly", v1_5_managedExecuteReadOnly, C.v1_5_managedExecuteReadOnly)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedCreateContract", v1_5_managedCreateContract, C.v1_5_managedCreateContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedDeployFromSourceContract", v1_5_managedDeployFromSourceContract, C.v1_5_managedDeployFromSourceContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedUpgradeContract", v1_5_managedUpgradeContract, C.v1_5_managedUpgradeContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedUpgradeFromSourceContract", v1_5_managedUpgradeFromSourceContract, C.v1_5_managedUpgradeFromSourceContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedDeleteContract", v1_5_managedDeleteContract, C.v1_5_managedDeleteContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedAsyncCall", v1_5_managedAsyncCall, C.v1_5_managedAsyncCall)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedCreateAsyncCall", v1_5_managedCreateAsyncCall, C.v1_5_managedCreateAsyncCall)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedGetCallbackClosure", v1_5_managedGetCallbackClosure, C.v1_5_managedGetCallbackClosure)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedGetMultiESDTCallValue", v1_5_managedGetMultiESDTCallValue, C.v1_5_managedGetMultiESDTCallValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedGetESDTBalance", v1_5_managedGetESDTBalance, C.v1_5_managedGetESDTBalance)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedGetESDTTokenData", v1_5_managedGetESDTTokenData, C.v1_5_managedGetESDTTokenData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedGetReturnData", v1_5_managedGetReturnData, C.v1_5_managedGetReturnData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedGetPrevBlockRandomSeed", v1_5_managedGetPrevBlockRandomSeed, C.v1_5_managedGetPrevBlockRandomSeed)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedGetBlockRandomSeed", v1_5_managedGetBlockRandomSeed, C.v1_5_managedGetBlockRandomSeed)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedGetStateRootHash", v1_5_managedGetStateRootHash, C.v1_5_managedGetStateRootHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedGetOriginalTxHash", v1_5_managedGetOriginalTxHash, C.v1_5_managedGetOriginalTxHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedIsESDTFrozen", v1_5_managedIsESDTFrozen, C.v1_5_managedIsESDTFrozen)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedIsESDTPaused", v1_5_managedIsESDTPaused, C.v1_5_managedIsESDTPaused)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedIsESDTLimitedTransfer", v1_5_managedIsESDTLimitedTransfer, C.v1_5_managedIsESDTLimitedTransfer)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedBufferToHex", v1_5_managedBufferToHex, C.v1_5_managedBufferToHex)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export v1_5_managedSCAddress
func v1_5_managedSCAddress(context unsafe.Pointer, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetSCAddress
	metering.UseGasAndAddTracedGas(managedSCAddressName, gasToUse)

	scAddress := runtime.GetContextAddress()

	managedType.SetBytes(destinationHandle, scAddress)
}

//export v1_5_managedOwnerAddress
func v1_5_managedOwnerAddress(context unsafe.Pointer, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetOwnerAddress
	metering.UseGasAndAddTracedGas(managedOwnerAddressName, gasToUse)

	owner, err := blockchain.GetOwnerAddress()
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	managedType.SetBytes(destinationHandle, owner)
}

//export v1_5_managedCaller
func v1_5_managedCaller(context unsafe.Pointer, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCaller
	metering.UseGasAndAddTracedGas(managedCallerName, gasToUse)

	caller := runtime.GetVMInput().CallerAddr
	managedType.SetBytes(destinationHandle, caller)
}

//export v1_5_managedSignalError
func v1_5_managedSignalError(context unsafe.Pointer, errHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	metering.StartGasTracing(managedSignalErrorName)

	gasToUse := metering.GasSchedule().ElrondAPICost.SignalError
	metering.UseAndTraceGas(gasToUse)

	errBytes, err := managedType.GetBytes(errHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBytes(errBytes)

	gasToUse = metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(len(errBytes))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		_ = arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}

	runtime.SignalUserError(string(errBytes))
}

//export v1_5_managedWriteLog
func v1_5_managedWriteLog(
	context unsafe.Pointer,
	topicsHandle int32,
	dataHandle int32,
) {
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)
	metering.StartGasTracing(managedWriteLogName)

	topics, sumOfTopicByteLengths, err := managedType.ReadManagedVecOfManagedBuffers(topicsHandle)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	dataBytes, err := managedType.GetBytes(dataHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBytes(dataBytes)
	dataByteLen := uint64(len(dataBytes))

	gasToUse := metering.GasSchedule().ElrondAPICost.Log
	gasForData := math.MulUint64(
		metering.GasSchedule().BaseOperationCost.DataCopyPerByte,
		sumOfTopicByteLengths+dataByteLen)
	gasToUse = math.AddUint64(gasToUse, gasForData)
	metering.UseAndTraceGas(gasToUse)

	output.WriteLog(runtime.GetContextAddress(), topics, dataBytes)
}

//export v1_5_managedGetOriginalTxHash
func v1_5_managedGetOriginalTxHash(context unsafe.Pointer, resultHandle int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetOriginalTxHash
	metering.UseGasAndAddTracedGas(managedGetOriginalTxHashName, gasToUse)

	managedType.SetBytes(resultHandle, runtime.GetOriginalTxHash())
}

//export v1_5_managedGetStateRootHash
func v1_5_managedGetStateRootHash(context unsafe.Pointer, resultHandle int32) {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetStateRootHash
	metering.UseGasAndAddTracedGas(managedGetStateRootHashName, gasToUse)

	managedType.SetBytes(resultHandle, blockchain.GetStateRootHash())
}

//export v1_5_managedGetBlockRandomSeed
func v1_5_managedGetBlockRandomSeed(context unsafe.Pointer, resultHandle int32) {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockRandomSeed
	metering.UseGasAndAddTracedGas(managedGetBlockRandomSeedName, gasToUse)

	managedType.SetBytes(resultHandle, blockchain.CurrentRandomSeed())
}

//export v1_5_managedGetPrevBlockRandomSeed
func v1_5_managedGetPrevBlockRandomSeed(context unsafe.Pointer, resultHandle int32) {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockRandomSeed
	metering.UseGasAndAddTracedGas(managedGetPrevBlockRandomSeedName, gasToUse)

	managedType.SetBytes(resultHandle, blockchain.LastRandomSeed())
}

//export v1_5_managedGetReturnData
func v1_5_managedGetReturnData(context unsafe.Pointer, resultID int32, resultHandle int32) {
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetReturnData
	metering.UseGasAndAddTracedGas(managedGetReturnDataName, gasToUse)

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) || resultID < 0 {
		_ = arwen.WithFault(arwen.ErrArgOutOfRange, context, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}

	managedType.SetBytes(resultHandle, returnData[resultID])
}

//export v1_5_managedGetMultiESDTCallValue
func v1_5_managedGetMultiESDTCallValue(context unsafe.Pointer, multiCallValueHandle int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCallValue
	metering.UseGasAndAddTracedGas(managedGetMultiESDTCallValueName, gasToUse)

	esdtTransfers := runtime.GetVMInput().ESDTTransfers
	multiCallBytes := writeESDTTransfersToBytes(managedType, esdtTransfers)
	managedType.ConsumeGasForBytes(multiCallBytes)

	managedType.SetBytes(multiCallValueHandle, multiCallBytes)
}

//export v1_5_managedGetESDTBalance
func v1_5_managedGetESDTBalance(context unsafe.Pointer, addressHandle int32, tokenIDHandle int32, nonce int64, valueHandle int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetExternalBalance
	metering.UseGasAndAddTracedGas(managedGetESDTBalanceName, gasToUse)

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		_ = arwen.WithFault(arwen.ErrArgOutOfRange, context, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}
	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		_ = arwen.WithFault(arwen.ErrArgOutOfRange, context, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}

	esdtToken, err := blockchain.GetESDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		_ = arwen.WithFault(arwen.ErrArgOutOfRange, context, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}

	value := managedType.GetBigIntOrCreate(valueHandle)
	value.Set(esdtToken.Value)
}

//export v1_5_managedGetESDTTokenData
func v1_5_managedGetESDTTokenData(
	context unsafe.Pointer,
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64,
	valueHandle, propertiesHandle, hashHandle, nameHandle, attributesHandle, creatorHandle, royaltiesHandle, urisHandle int32) {
	host := arwen.GetVMHost(context)
	ManagedGetESDTTokenDataWithHost(
		host,
		addressHandle,
		tokenIDHandle,
		nonce,
		valueHandle, propertiesHandle, hashHandle, nameHandle, attributesHandle, creatorHandle, royaltiesHandle, urisHandle)

}

func ManagedGetESDTTokenDataWithHost(
	host arwen.VMHost,
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64,
	valueHandle, propertiesHandle, hashHandle, nameHandle, attributesHandle, creatorHandle, royaltiesHandle, urisHandle int32) {
	runtime := host.Runtime()
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedGetESDTTokenDataName)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetExternalBalance
	metering.UseAndTraceGas(gasToUse)

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		_ = arwen.WithFaultAndHost(host, arwen.ErrArgOutOfRange, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}
	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		_ = arwen.WithFaultAndHost(host, arwen.ErrArgOutOfRange, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}

	esdtToken, err := blockchain.GetESDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		_ = arwen.WithFaultAndHost(host, arwen.ErrArgOutOfRange, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}

	value := managedType.GetBigIntOrCreate(valueHandle)
	value.Set(esdtToken.Value)

	managedType.SetBytes(propertiesHandle, esdtToken.Properties)
	if esdtToken.TokenMetaData != nil {
		managedType.SetBytes(hashHandle, esdtToken.TokenMetaData.Hash)
		managedType.ConsumeGasForBytes(esdtToken.TokenMetaData.Hash)
		managedType.SetBytes(nameHandle, esdtToken.TokenMetaData.Name)
		managedType.ConsumeGasForBytes(esdtToken.TokenMetaData.Name)
		managedType.SetBytes(attributesHandle, esdtToken.TokenMetaData.Attributes)
		managedType.ConsumeGasForBytes(esdtToken.TokenMetaData.Attributes)
		managedType.SetBytes(creatorHandle, esdtToken.TokenMetaData.Creator)
		managedType.ConsumeGasForBytes(esdtToken.TokenMetaData.Creator)
		royalties := managedType.GetBigIntOrCreate(royaltiesHandle)
		royalties.SetUint64(uint64(esdtToken.TokenMetaData.Royalties))

		managedType.WriteManagedVecOfManagedBuffers(esdtToken.TokenMetaData.URIs, urisHandle)
	}

}

//export v1_5_managedAsyncCall
func v1_5_managedAsyncCall(
	context unsafe.Pointer,
	destHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32) {
	host := arwen.GetVMHost(context)
	ManagedAsyncCallWithHost(
		host,
		destHandle,
		valueHandle,
		functionHandle,
		argumentsHandle)
}

func ManagedAsyncCallWithHost(
	host arwen.VMHost,
	destHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32) {
	runtime := host.Runtime()
	async := host.Async()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedAsyncCallName)

	gasSchedule := metering.GasSchedule()
	gasToUse := gasSchedule.ElrondAPICost.AsyncCallStep
	metering.UseAndTraceGas(gasToUse)

	vmInput, err := readDestinationFunctionArguments(host, destHandle, functionHandle, argumentsHandle)
	if arwen.WithFaultAndHost(host, err, host.Runtime().ElrondAPIErrorShouldFailExecution()) {
		return
	}

	data := makeCrossShardCallFromInput(vmInput.function, vmInput.arguments)

	value, err := managedType.GetBigInt(valueHandle)
	if err != nil {
		_ = arwen.WithFaultAndHost(host, arwen.ErrArgOutOfRange, host.Runtime().ElrondAPIErrorShouldFailExecution())
		return
	}

	gasToUse = math.MulUint64(gasSchedule.BaseOperationCost.DataCopyPerByte, uint64(len(data)))
	metering.UseAndTraceGas(gasToUse)

	err = async.RegisterLegacyAsyncCall(vmInput.destination, []byte(data), value.Bytes())
	if errors.Is(err, arwen.ErrNotEnoughGas) {
		runtime.SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
		return
	}
	if arwen.WithFaultAndHost(host, err, host.Runtime().ElrondAPIErrorShouldFailExecution()) {
		return
	}
}

//export v1_5_managedCreateAsyncCall
func v1_5_managedCreateAsyncCall(
	context unsafe.Pointer,
	destHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	successOffset int32,
	successLength int32,
	errorOffset int32,
	errorLength int32,
	gas int64,
	extraGasForCallback int64,
	callbackClosureHandle int32,
) int32 {

	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	managedType := host.ManagedTypes()

	vmInput, err := readDestinationFunctionArguments(host, destHandle, functionHandle, argumentsHandle)
	if arwen.WithFaultAndHost(host, err, host.Runtime().ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	data := makeCrossShardCallFromInput(vmInput.function, vmInput.arguments)

	value, err := managedType.GetBigInt(valueHandle)
	if err != nil {
		_ = arwen.WithFault(arwen.ErrArgOutOfRange, context, runtime.ElrondAPIErrorShouldFailExecution())
		return 1
	}

	successFunc, err := runtime.MemLoad(successOffset, successLength)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	errorFunc, err := runtime.MemLoad(errorOffset, errorLength)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	callbackClosure, err := managedType.GetBytes(callbackClosureHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	return CreateAsyncCallWithTypedArgs(host,
		vmInput.destination,
		value.Bytes(),
		[]byte(data),
		successFunc,
		errorFunc,
		gas,
		extraGasForCallback,
		callbackClosure)
}

//export v1_5_managedGetCallbackClosure
func v1_5_managedGetCallbackClosure(
	context unsafe.Pointer,
	callbackClosureHandle int32,
) {
	host := arwen.GetVMHost(context)
	GetCallbackClosureWithHost(host, callbackClosureHandle)
}

func GetCallbackClosureWithHost(
	host arwen.VMHost,
	callbackClosureHandle int32,
) {
	runtime := host.Runtime()
	async := host.Async()
	metering := host.Metering()
	managedTypes := host.ManagedTypes()

	metering.StartGasTracing(managedGetCallbackClosure)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCallbackClosure
	metering.UseAndTraceGas(gasToUse)

	callbackClosure, err := async.GetCallbackClosure()
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	managedTypes.SetBytes(callbackClosureHandle, callbackClosure)
}

//export v1_5_managedUpgradeFromSourceContract
func v1_5_managedUpgradeFromSourceContract(
	context unsafe.Pointer,
	destHandle int32,
	gas int64,
	valueHandle int32,
	addressHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) {
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedUpgradeFromSourceContractName)

	gasToUse := metering.GasSchedule().ElrondAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	vmInput, err := readDestinationValueArguments(host, destHandle, valueHandle, argumentsHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	sourceContractAddress, err := managedType.GetBytes(addressHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	lenReturnData := len(host.Output().ReturnData())

	UpgradeFromSourceContractWithTypedArgs(
		host,
		sourceContractAddress,
		vmInput.destination,
		vmInput.value.Bytes(),
		vmInput.arguments,
		gas,
		codeMetadata,
	)
	setReturnDataIfExists(host, lenReturnData, resultHandle)
}

//export v1_5_managedUpgradeContract
func v1_5_managedUpgradeContract(
	context unsafe.Pointer,
	destHandle int32,
	gas int64,
	valueHandle int32,
	codeHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) {
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedUpgradeContractName)

	gasToUse := metering.GasSchedule().ElrondAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	vmInput, err := readDestinationValueArguments(host, destHandle, valueHandle, argumentsHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	code, err := managedType.GetBytes(codeHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	lenReturnData := len(host.Output().ReturnData())

	upgradeContract(host, vmInput.destination, code, codeMetadata, vmInput.value.Bytes(), vmInput.arguments, gas)
	setReturnDataIfExists(host, lenReturnData, resultHandle)
}

//export v1_5_managedDeleteContract
func v1_5_managedDeleteContract(
	context unsafe.Pointer,
	destHandle int32,
	gasLimit int64,
	argumentsHandle int32,
) {
	host := arwen.GetVMHost(context)
	ManagedDeleteContractWithHost(
		host,
		destHandle,
		gasLimit,
		argumentsHandle,
	)
}

func ManagedDeleteContractWithHost(
	host arwen.VMHost,
	destHandle int32,
	gasLimit int64,
	argumentsHandle int32,
) {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(deleteContractName)

	gasToUse := metering.GasSchedule().ElrondAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	calledSCAddress, err := managedType.GetBytes(destHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	data, _, err := managedType.ReadManagedVecOfManagedBuffers(argumentsHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	deleteContract(
		host,
		calledSCAddress,
		data,
		gasLimit,
	)
}

//export v1_5_managedDeployFromSourceContract
func v1_5_managedDeployFromSourceContract(
	context unsafe.Pointer,
	gas int64,
	valueHandle int32,
	addressHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultAddressHandle int32,
	resultHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedDeployFromSourceContractName)

	gasToUse := metering.GasSchedule().ElrondAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	vmInput, err := readDestinationValueArguments(host, addressHandle, valueHandle, argumentsHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	lenReturnData := len(host.Output().ReturnData())

	newAddress, err := DeployFromSourceContractWithTypedArgs(
		host,
		vmInput.destination,
		codeMetadata,
		vmInput.value,
		vmInput.arguments,
		gas,
	)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.SetBytes(resultAddressHandle, newAddress)
	setReturnDataIfExists(host, lenReturnData, resultHandle)

	return 0
}

//export v1_5_managedCreateContract
func v1_5_managedCreateContract(
	context unsafe.Pointer,
	gas int64,
	valueHandle int32,
	codeHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultAddressHandle int32,
	resultHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedCreateContractName)

	gasToUse := metering.GasSchedule().ElrondAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	sender := runtime.GetContextAddress()
	value, err := managedType.GetBigInt(valueHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	data, actualLen, err := managedType.ReadManagedVecOfManagedBuffers(argumentsHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, actualLen)
	metering.UseAndTraceGas(gasToUse)

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	code, err := managedType.GetBytes(codeHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	lenReturnData := len(host.Output().ReturnData())
	newAddress, err := createContract(sender, data, value, metering, gas, code, codeMetadata, host, runtime)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	managedType.SetBytes(resultAddressHandle, newAddress)
	setReturnDataIfExists(host, lenReturnData, resultHandle)

	return 0
}

func setReturnDataIfExists(
	host arwen.VMHost,
	oldLen int,
	resultHandle int32,
) {
	returnData := host.Output().ReturnData()
	if len(returnData) > oldLen {
		host.ManagedTypes().WriteManagedVecOfManagedBuffers(returnData[oldLen:], resultHandle)
	} else {
		host.ManagedTypes().SetBytes(resultHandle, make([]byte, 0))
	}
}

//export v1_5_managedExecuteReadOnly
func v1_5_managedExecuteReadOnly(
	context unsafe.Pointer,
	gas int64,
	addressHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteReadOnlyName)

	vmInput, err := readDestinationFunctionArguments(host, addressHandle, functionHandle, argumentsHandle)
	if arwen.WithFaultAndHost(host, err, host.Runtime().ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	lenReturnData := len(host.Output().ReturnData())
	returnVal := ExecuteReadOnlyWithTypedArguments(
		host,
		gas,
		[]byte(vmInput.function),
		vmInput.destination,
		vmInput.arguments,
	)
	setReturnDataIfExists(host, lenReturnData, resultHandle)
	return returnVal
}

//export v1_5_managedExecuteOnSameContext
func v1_5_managedExecuteOnSameContext(
	context unsafe.Pointer,
	gas int64,
	addressHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteOnSameContextName)

	vmInput, err := readDestinationValueFunctionArguments(host, addressHandle, valueHandle, functionHandle, argumentsHandle)
	if arwen.WithFaultAndHost(host, err, host.Runtime().ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	lenReturnData := len(host.Output().ReturnData())
	returnVal := ExecuteOnSameContextWithTypedArgs(
		host,
		gas,
		vmInput.value,
		[]byte(vmInput.function),
		vmInput.destination,
		vmInput.arguments,
	)
	setReturnDataIfExists(host, lenReturnData, resultHandle)
	return returnVal
}

//export v1_5_managedExecuteOnDestContext
func v1_5_managedExecuteOnDestContext(
	context unsafe.Pointer,
	gas int64,
	addressHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteOnDestContextName)

	vmInput, err := readDestinationValueFunctionArguments(host, addressHandle, valueHandle, functionHandle, argumentsHandle)
	if arwen.WithFaultAndHost(host, err, host.Runtime().ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	lenReturnData := len(host.Output().ReturnData())
	returnVal := ExecuteOnDestContextWithTypedArgs(
		host,
		gas,
		vmInput.value,
		[]byte(vmInput.function),
		vmInput.destination,
		vmInput.arguments,
	)
	setReturnDataIfExists(host, lenReturnData, resultHandle)
	return returnVal
}

//export v1_5_managedMultiTransferESDTNFTExecute
func v1_5_managedMultiTransferESDTNFTExecute(
	context unsafe.Pointer,
	dstHandle int32,
	tokenTransfersHandle int32,
	gasLimit int64,
	functionHandle int32,
	argumentsHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	managedType := host.ManagedTypes()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(managedMultiTransferESDTNFTExecuteName)

	vmInput, err := readDestinationFunctionArguments(host, dstHandle, functionHandle, argumentsHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	transfers, err := readESDTTransfers(managedType, tokenTransfersHandle)
	if arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return TransferESDTNFTExecuteWithTypedArgs(
		host,
		vmInput.destination,
		transfers,
		gasLimit,
		[]byte(vmInput.function),
		vmInput.arguments,
	)
}

//export v1_5_managedTransferValueExecute
func v1_5_managedTransferValueExecute(
	context unsafe.Pointer,
	dstHandle int32,
	valueHandle int32,
	gasLimit int64,
	functionHandle int32,
	argumentsHandle int32,
) int32 {
	host := arwen.GetVMHost(context)
	metering := host.Metering()
	metering.StartGasTracing(managedTransferValueExecuteName)

	vmInput, err := readDestinationValueFunctionArguments(host, dstHandle, valueHandle, functionHandle, argumentsHandle)
	if arwen.WithFaultAndHost(host, err, host.Runtime().ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return TransferValueExecuteWithTypedArgs(
		host,
		vmInput.destination,
		vmInput.value,
		gasLimit,
		[]byte(vmInput.function),
		vmInput.arguments,
	)
}

//export v1_5_managedIsESDTFrozen
func v1_5_managedIsESDTFrozen(
	context unsafe.Pointer,
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64) int32 {
	host := arwen.GetVMHost(context)
	return ManagedIsESDTFrozenWithHost(host, addressHandle, tokenIDHandle, nonce)
}

func ManagedIsESDTFrozenWithHost(
	host arwen.VMHost,
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().ElrondAPICost.GetExternalBalance
	metering.UseGasAndAddTracedGas(managedIsESDTFrozenName, gasToUse)

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		_ = arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution())
		return -1
	}
	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		_ = arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution())
		return -1
	}

	esdtToken, err := blockchain.GetESDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		_ = arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution())
		return -1
	}

	esdtUserData := builtInFunctions.ESDTUserMetadataFromBytes(esdtToken.Properties)
	if esdtUserData.Frozen {
		return 1
	}
	return 0
}

//export v1_5_managedIsESDTLimitedTransfer
func v1_5_managedIsESDTLimitedTransfer(context unsafe.Pointer, tokenIDHandle int32) int32 {
	host := arwen.GetVMHost(context)
	return ManagedIsESDTLimitedTransferWithHost(host, tokenIDHandle)
}

func ManagedIsESDTLimitedTransferWithHost(host arwen.VMHost, tokenIDHandle int32) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().ElrondAPICost.GetExternalBalance
	metering.UseGasAndAddTracedGas(managedIsESDTLimitedTransferName, gasToUse)

	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		_ = arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution())
		return -1
	}

	if blockchain.IsLimitedTransfer(tokenID) {
		return 1
	}

	return 0
}

//export v1_5_managedIsESDTPaused
func v1_5_managedIsESDTPaused(context unsafe.Pointer, tokenIDHandle int32) int32 {
	host := arwen.GetVMHost(context)
	return ManagedIsESDTPausedWithHost(host, tokenIDHandle)
}

func ManagedIsESDTPausedWithHost(host arwen.VMHost, tokenIDHandle int32) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().ElrondAPICost.GetExternalBalance
	metering.UseGasAndAddTracedGas(managedIsESDTPausedName, gasToUse)

	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		_ = arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution())
		return -1
	}

	if blockchain.IsPaused(tokenID) {
		return 1
	}

	return 0
}

//export v1_5_managedBufferToHex
func v1_5_managedBufferToHex(context unsafe.Pointer, sourceHandle int32, destHandle int32) {
	host := arwen.GetVMHost(context)
	ManagedBufferToHexWithHost(host, sourceHandle, destHandle)
}

func ManagedBufferToHexWithHost(host arwen.VMHost, sourceHandle int32, destHandle int32) {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferSetBytes
	metering.UseGasAndAddTracedGas(managedBufferToHexName, gasToUse)

	mBuff, err := managedType.GetBytes(sourceHandle)
	if err != nil {
		arwen.WithFaultAndHost(host, err, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}

	encoded := hex.EncodeToString(mBuff)
	managedType.SetBytes(destHandle, []byte(encoded))
}

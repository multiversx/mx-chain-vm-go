package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern void	v1_4_managedSCAddress(void *context, int32_t addressHandle);
// extern void	v1_4_managedOwnerAddress(void *context, int32_t addressHandle);
// extern void	v1_4_managedCaller(void *context, int32_t addressHandle);
// extern void	v1_4_managedSignalError(void* context, int32_t errHandle1);
// extern void	v1_4_managedWriteLog(void* context, int32_t topicsHandle, int32_t dataHandle);
//
// extern int32_t	v1_4_managedMultiTransferESDTNFTExecute(void *context, int32_t dstHandle, int32_t tokenTransfersHandle, long long gasLimit, int32_t functionHandle, int32_t argumentsHandle);
// extern int32_t	v1_4_managedTransferValueExecute(void *context, int32_t dstHandle, int32_t valueHandle, long long gasLimit, int32_t functionHandle, int32_t argumentsHandle);
// extern int32_t	v1_4_managedExecuteOnDestContext(void *context, long long gas, int32_t addressHandle, int32_t valueHandle, int32_t functionHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern int32_t	v1_4_managedExecuteOnDestContextByCaller(void *context, long long gas, int32_t addressHandle, int32_t valueHandle, int32_t functionHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern int32_t	v1_4_managedExecuteOnSameContext(void *context, long long gas, int32_t addressHandle, int32_t valueHandle, int32_t functionHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern int32_t	v1_4_managedDelegateExecution(void *context, long long gas, int32_t addressHandle, int32_t valueHandle, int32_t functionHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern int32_t	v1_4_managedExecuteReadOnly(void *context, long long gas, int32_t addressHandle, int32_t valueHandle, int32_t functionHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern int32_t	v1_4_managedCreateContract(void *context, long long gas, int32_t valueHandle, int32_t codeHandle, int32_t codeMetadataHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern int32_t	v1_4_managedDeployFromSourceContract(void *context, long long gas, int32_t valueHandle, int32_t addressHandle, int32_t codeMetadataHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern void		v1_4_managedUpgradeContract(void *context, int32_t dstHandle, long long gas, int32_t valueHandle, int32_t codeHandle, int32_t codeMetadataHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern void		v1_4_managedUpgradeFromSourceContract(void *context, int32_t dstHandle, long long gas, int32_t valueHandle, int32_t addressHandle, int32_t codeMetadataHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern void		v1_4_managedAsyncCall(void *context, int32_t dstHandle, int32_t valueHandle, int32_t dataHandle);
//
// extern int32_t	v1_4_managedGetESDTBalance(void *context, int32_t addressHandle, int32_t tokenIDHandle, long long nonce, int32_t valueHandle);
// extern int32_t	v1_4_managedGetESDTTokenData(void *context, int32_t addressHandle, int32_t tokenIDHandle, long long nonce, int32_t valueHandle, int32_t propertiesHandle, int32_t hashHandle, int32_t nameHandle, int32_t attributesHandle, int32_t creatorHandle, int32_t royaltiesHandle, int32_t urisHandle);
//
// extern int32_t	v1_4_managedGetReturnData(void *context, int32_t resultID, int32_t resultHandle);
// extern void		v1_4_managedGetPrevBlockRandomSeed(void *context, int32_t resultHandle);
// extern void		v1_4_managedGetBlockRandomSeed(void *context, int32_t resultHandle);
// extern void		v1_4_managedGetStateRootHash(void *context, int32_t resultHandle);
// extern void		v1_4_managedGetOriginalTxHash(void *context, int32_t resultHandle);
//
import "C"

import (
	"encoding/binary"
	"errors"
	"math/big"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/math"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/wasmer"
)

// ManagedEIImports creates a new wasmer.Imports populated with variants of the API methods that use managed types only.
func ManagedEIImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("env")

	imports, err := imports.Append("managedSCAddress", v1_4_managedSCAddress, C.v1_4_managedSCAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedOwnerAddress", v1_4_managedOwnerAddress, C.v1_4_managedOwnerAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedCaller", v1_4_managedCaller, C.v1_4_managedCaller)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedSignalError", v1_4_managedSignalError, C.v1_4_managedSignalError)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedWriteLog", v1_4_managedWriteLog, C.v1_4_managedWriteLog)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedMultiTransferESDTNFTExecute", v1_4_managedMultiTransferESDTNFTExecute, C.v1_4_managedMultiTransferESDTNFTExecute)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedTransferValueExecute", v1_4_managedTransferValueExecute, C.v1_4_managedTransferValueExecute)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedExecuteOnDestContext", v1_4_managedExecuteOnDestContext, C.v1_4_managedExecuteOnDestContext)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedExecuteOnDestContextByCaller", v1_4_managedExecuteOnDestContextByCaller, C.v1_4_managedExecuteOnDestContextByCaller)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedExecuteOnSameContext", v1_4_managedExecuteOnSameContext, C.v1_4_managedExecuteOnSameContext)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedDelegateExecution", v1_4_managedDelegateExecution, C.v1_4_managedDelegateExecution)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedExecuteReadOnly", v1_4_managedExecuteReadOnly, C.v1_4_managedExecuteReadOnly)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedCreateContract", v1_4_managedCreateContract, C.v1_4_managedCreateContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedDeployFromSourceContract", v1_4_managedDeployFromSourceContract, C.v1_4_managedDeployFromSourceContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedUpgradeContract", v1_4_managedUpgradeContract, C.v1_4_managedUpgradeContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedUpgradeFromSourceContract", v1_4_managedUpgradeFromSourceContract, C.v1_4_managedUpgradeFromSourceContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedAsyncCall", v1_4_managedAsyncCall, C.v1_4_managedAsyncCall)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedGetESDTBalance", v1_4_managedGetESDTBalance, C.v1_4_managedGetESDTBalance)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedGetESDTTokenData", v1_4_managedGetESDTTokenData, C.v1_4_managedGetESDTTokenData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedGetReturnData", v1_4_managedGetReturnData, C.v1_4_managedGetReturnData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedGetPrevBlockRandomSeed", v1_4_managedGetPrevBlockRandomSeed, C.v1_4_managedGetPrevBlockRandomSeed)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedGetBlockRandomSeed", v1_4_managedGetBlockRandomSeed, C.v1_4_managedGetBlockRandomSeed)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedGetStateRootHash", v1_4_managedGetStateRootHash, C.v1_4_managedGetStateRootHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("v1_4_managedGetOriginalTxHash", v1_4_managedGetOriginalTxHash, C.v1_4_managedGetOriginalTxHash)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export v1_4_managedSCAddress
func v1_4_managedSCAddress(context unsafe.Pointer, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetSCAddress
	metering.UseGas(gasToUse)

	scAddress := runtime.GetSCAddress()

	managedType.SetBytes(destinationHandle, scAddress)
}

//export v1_4_managedOwnerAddress
func v1_4_managedOwnerAddress(context unsafe.Pointer, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetOwnerAddress
	metering.UseGas(gasToUse)

	owner, err := blockchain.GetOwnerAddress()
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	managedType.SetBytes(destinationHandle, owner)
}

//export v1_4_managedCaller
func v1_4_managedCaller(context unsafe.Pointer, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCaller
	metering.UseGas(gasToUse)

	caller := runtime.GetVMInput().CallerAddr
	managedType.SetBytes(destinationHandle, caller)
}

//export v1_4_managedSignalError
func v1_4_managedSignalError(context unsafe.Pointer, errHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.SignalError
	metering.UseGas(gasToUse)

	errBytes, err := managedType.GetBytes(errHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForThisIntNumberOfBytes(len(errBytes))

	runtime.SignalUserError(string(errBytes))
}

//export v1_4_managedWriteLog
func v1_4_managedWriteLog(
	context unsafe.Pointer,
	topicsHandle int32,
	dataHandle int32,
) {
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	topics, sumOfTopicByteLengths, err := readManagedVecOfManagedBuffers(managedType, topicsHandle)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	dataBytes, err := managedType.GetBytes(dataHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForThisIntNumberOfBytes(len(dataBytes))
	dataByteLen := uint64(len(dataBytes))

	gasToUse := metering.GasSchedule().ElrondAPICost.Log
	gasForData := math.MulUint64(
		metering.GasSchedule().BaseOperationCost.DataCopyPerByte,
		sumOfTopicByteLengths+dataByteLen)
	gasToUse = math.AddUint64(gasToUse, gasForData)
	metering.UseGas(gasToUse)

	output.WriteLog(runtime.GetSCAddress(), topics, dataBytes)
}

func readManagedVecOfManagedBuffers(
	managedType arwen.ManagedTypesContext,
	managedVecHandle int32,
) ([][]byte, uint64, error) {
	managedVecBytes, err := managedType.GetBytes(managedVecHandle)
	if err != nil {
		return nil, 0, err
	}
	managedType.ConsumeGasForThisIntNumberOfBytes(len(managedVecBytes))

	if len(managedVecBytes)%4 != 0 {
		return nil, 0, errors.New("invalid managed vector of managed buffers")
	}

	numBuffers := len(managedVecBytes) / 4
	result := make([][]byte, 0, numBuffers)
	sumOfItemByteLengths := uint64(0)
	for i := 0; i < len(managedVecBytes); i += 4 {
		itemHandle := int32(binary.BigEndian.Uint32(managedVecBytes[i : i+4]))

		itemBytes, err := managedType.GetBytes(itemHandle)
		if err != nil {
			return nil, 0, err
		}
		managedType.ConsumeGasForThisIntNumberOfBytes(len(itemBytes))

		sumOfItemByteLengths += uint64(len(itemBytes))
		result = append(result, itemBytes)
	}

	return result, sumOfItemByteLengths, nil
}

func writeManagedVecOfManagedBuffers(
	managedType arwen.ManagedTypesContext,
	data [][]byte,
	destinationHandle int32,
) uint64 {
	sumOfItemByteLengths := uint64(0)
	destinationBytes := make([]byte, 4*len(data))
	dataIndex := 0
	for _, itemBytes := range data {
		sumOfItemByteLengths += uint64(len(itemBytes))
		itemHandle := managedType.NewManagedBufferFromBytes(itemBytes)
		binary.BigEndian.PutUint32(destinationBytes[dataIndex:dataIndex+4], uint32(itemHandle))
		dataIndex += 4
	}

	managedType.SetBytes(destinationHandle, destinationBytes)

	return sumOfItemByteLengths
}

//export v1_4_managedGetOriginalTxHash
func v1_4_managedGetOriginalTxHash(context unsafe.Pointer, resultHandle int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockHash
	metering.UseGas(gasToUse)

	managedType.SetBytes(resultHandle, runtime.GetOriginalTxHash())
}

//export v1_4_managedGetStateRootHash
func v1_4_managedGetStateRootHash(context unsafe.Pointer, resultHandle int32) {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetStateRootHash
	metering.UseGas(gasToUse)

	managedType.SetBytes(resultHandle, blockchain.GetStateRootHash())
}

//export v1_4_managedGetBlockRandomSeed
func v1_4_managedGetBlockRandomSeed(context unsafe.Pointer, resultHandle int32) {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockRandomSeed
	metering.UseGas(gasToUse)

	managedType.SetBytes(resultHandle, blockchain.CurrentRandomSeed())
}

//export v1_4_managedGetPrevBlockRandomSeed
func v1_4_managedGetPrevBlockRandomSeed(context unsafe.Pointer, resultHandle int32) {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockRandomSeed
	metering.UseGas(gasToUse)

	managedType.SetBytes(resultHandle, blockchain.LastRandomSeed())
}

//export v1_4_managedGetPrevBlockRandomSeed
func v1_4_managedGetReturnData(context unsafe.Pointer, resultID int32, resultHandle int32) {
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetReturnData
	metering.UseGas(gasToUse)

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) {
		_ = arwen.WithFault(arwen.ErrArgOutOfRange, context, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}

	managedType.SetBytes(resultHandle, returnData[resultID])
}

//export v1_4_managedGetESDTBalance
func v1_4_managedGetESDTBalance(context unsafe.Pointer, addressHandle int32, tokenIDHandle int32, nonce int64, valueHandle int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetExternalBalance
	metering.UseGas(gasToUse)

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

//export v1_4_managedGetESDTBalance
func v1_4_managedGetESDTTokenData(context unsafe.Pointer, addressHandle int32, tokenIDHandle int32, nonce int64,
	valueHandle, propertiesHandle, hashHandle, nameHandle, attributesHandle, creatorHandle, royaltiesHandle, urisHandle int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetExternalBalance
	metering.UseGas(gasToUse)

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

	managedType.SetBytes(propertiesHandle, esdtToken.Properties)
	if esdtToken.TokenMetaData != nil {
		managedType.SetBytes(hashHandle, esdtToken.TokenMetaData.Hash)
		managedType.ConsumeGasForThisIntNumberOfBytes(len(esdtToken.TokenMetaData.Hash))
		managedType.SetBytes(nameHandle, esdtToken.TokenMetaData.Name)
		managedType.ConsumeGasForThisIntNumberOfBytes(len(esdtToken.TokenMetaData.Name))
		managedType.SetBytes(attributesHandle, esdtToken.TokenMetaData.Attributes)
		managedType.ConsumeGasForThisIntNumberOfBytes(len(esdtToken.TokenMetaData.Attributes))
		managedType.SetBytes(creatorHandle, esdtToken.TokenMetaData.Creator)
		managedType.ConsumeGasForThisIntNumberOfBytes(len(esdtToken.TokenMetaData.Creator))
		royalties := managedType.GetBigIntOrCreate(royaltiesHandle)
		royalties.SetUint64(uint64(esdtToken.TokenMetaData.Royalties))

		totalBytes := writeManagedVecOfManagedBuffers(managedType, esdtToken.TokenMetaData.URIs, urisHandle)
		managedType.ConsumeGasForThisIntNumberOfBytes(int(totalBytes))
	}
}

//export v1_4_managedAsyncCall
func v1_4_managedAsyncCall(context unsafe.Pointer, destHandle int32, valueHandle int32, dataHandle int32) {
	vmHost := arwen.GetVMHost(context)
	runtime := vmHost.Runtime()
	metering := vmHost.Metering()
	managedType := vmHost.ManagedTypes()

	gasSchedule := metering.GasSchedule()
	gasToUse := gasSchedule.ElrondAPICost.AsyncCallStep
	metering.UseGas(gasToUse)

	address, err := managedType.GetBytes(destHandle)
	if err != nil {
		_ = arwen.WithFault(arwen.ErrArgOutOfRange, context, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}

	data, err := managedType.GetBytes(dataHandle)
	if err != nil {
		_ = arwen.WithFault(arwen.ErrArgOutOfRange, context, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}

	value, err := managedType.GetBigInt(valueHandle)
	if err != nil {
		_ = arwen.WithFault(arwen.ErrArgOutOfRange, context, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}

	gasToUse = math.MulUint64(gasSchedule.BaseOperationCost.DataCopyPerByte, uint64(len(data)))
	metering.UseGas(gasToUse)

	err = runtime.ExecuteAsyncCall(address, data, value.Bytes())
	if errors.Is(err, arwen.ErrNotEnoughGas) {
		runtime.SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
		return
	}
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}
}

type vmInputData struct {
	destination []byte
	function    string
	value       *big.Int
	arguments   [][]byte
}

func readDestinationValueFunctionArguments(
	vmHost arwen.VMHost,
	destHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
) (*vmInputData, error) {
	managedType := vmHost.ManagedTypes()
	runtime := vmHost.Runtime()

	vmInput, err := readDestinationValueArguments(vmHost, destHandle, valueHandle, argumentsHandle)
	if err != nil {
		return nil, err
	}

	function, err := managedType.GetBytes(functionHandle)
	if err != nil {
		_ = arwen.WithFaultAndHost(vmHost, arwen.ErrArgOutOfRange, runtime.ElrondAPIErrorShouldFailExecution())
		return nil, err
	}
	vmInput.function = string(function)

	return vmInput, err
}

func readDestinationValueArguments(
	vmHost arwen.VMHost,
	destHandle int32,
	valueHandle int32,
	argumentsHandle int32,
) (*vmInputData, error) {
	managedType := vmHost.ManagedTypes()
	runtime := vmHost.Runtime()
	metering := vmHost.Metering()

	var err error
	vmInput := &vmInputData{}

	vmInput.destination, err = managedType.GetBytes(destHandle)
	if err != nil {
		_ = arwen.WithFaultAndHost(vmHost, arwen.ErrArgOutOfRange, runtime.ElrondAPIErrorShouldFailExecution())
		return nil, err
	}

	vmInput.value, err = managedType.GetBigInt(valueHandle)
	if err != nil {
		_ = arwen.WithFaultAndHost(vmHost, arwen.ErrArgOutOfRange, runtime.ElrondAPIErrorShouldFailExecution())
		return nil, err
	}

	data, actualLen, err := readManagedVecOfManagedBuffers(managedType, argumentsHandle)
	if err != nil {
		_ = arwen.WithFaultAndHost(vmHost, arwen.ErrArgOutOfRange, runtime.ElrondAPIErrorShouldFailExecution())
		return nil, err
	}
	vmInput.arguments = data

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, actualLen)
	metering.UseGas(gasToUse)

	return vmInput, err
}

//export v1_4_managedUpgradeFromSourceContract
func v1_4_managedUpgradeFromSourceContract(
	context unsafe.Pointer,
	destHandle int32,
	gas int64,
	valueHandle int32,
	addressHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultsHandle int32,
) {
	vmHost := arwen.GetVMHost(context)
	runtime := vmHost.Runtime()
	metering := vmHost.Metering()
	managedType := vmHost.ManagedTypes()

	gasToUse := metering.GasSchedule().ElrondAPICost.CreateContract
	metering.UseGas(gasToUse)

	vmInput, err := readDestinationValueArguments(vmHost, destHandle, valueHandle, argumentsHandle)
	if err != nil {
		return
	}

	sourceContractAddress, err := managedType.GetBytes(addressHandle)
	if err != nil {
		_ = arwen.WithFaultAndHost(vmHost, arwen.ErrArgOutOfRange, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if err != nil {
		_ = arwen.WithFaultAndHost(vmHost, arwen.ErrArgOutOfRange, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}

	lenReturnData := len(vmHost.Output().ReturnData())

	UpgradeFromSourceContractWithTypedArgs(
		vmHost,
		sourceContractAddress,
		vmInput.destination,
		vmInput.value.Bytes(),
		vmInput.arguments,
		gas,
		codeMetadata,
	)

	returnData := vmHost.Output().ReturnData()
	if len(returnData) > lenReturnData {
		_ = writeManagedVecOfManagedBuffers(managedType, returnData[lenReturnData:], resultsHandle)
	} else {
		managedType.SetBytes(resultsHandle, make([]byte, 0))
	}
}

//export v1_4_upgradeContract
func v1_4_managedUpgradeContract(
	context unsafe.Pointer,
	destHandle int32,
	gas int64,
	codeHandle int32,
	valueHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultsHandle int32,
) {
	vmHost := arwen.GetVMHost(context)
	runtime := vmHost.Runtime()
	metering := vmHost.Metering()
	managedType := vmHost.ManagedTypes()

	gasToUse := metering.GasSchedule().ElrondAPICost.CreateContract
	metering.UseGas(gasToUse)

	vmInput, err := readDestinationValueArguments(vmHost, destHandle, valueHandle, argumentsHandle)
	if err != nil {
		return
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if err != nil {
		_ = arwen.WithFaultAndHost(vmHost, arwen.ErrArgOutOfRange, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}

	code, err := managedType.GetBytes(codeHandle)
	if err != nil {
		_ = arwen.WithFaultAndHost(vmHost, arwen.ErrArgOutOfRange, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}

	if arwen.WithFaultAndHost(vmHost, err, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	lenReturnData := len(vmHost.Output().ReturnData())

	upgradeContract(vmHost, vmInput.destination, code, codeMetadata, vmInput.value.Bytes(), vmInput.arguments, gas)

	returnData := vmHost.Output().ReturnData()
	if len(returnData) > lenReturnData {
		_ = writeManagedVecOfManagedBuffers(managedType, returnData[lenReturnData:], resultsHandle)
	} else {
		managedType.SetBytes(resultsHandle, make([]byte, 0))
	}
}

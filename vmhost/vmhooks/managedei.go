package vmhooks

import (
	"encoding/hex"
	"errors"

	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"

	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/math"
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

// ManagedSCAddress VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedSCAddress(destinationHandle int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetSCAddress
	metering.UseGasAndAddTracedGas(managedSCAddressName, gasToUse)

	scAddress := runtime.GetContextAddress()

	managedType.SetBytes(destinationHandle, scAddress)
}

// ManagedOwnerAddress VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedOwnerAddress(destinationHandle int32) {
	managedType := context.GetManagedTypesContext()
	blockchain := context.GetBlockchainContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetOwnerAddress
	metering.UseGasAndAddTracedGas(managedOwnerAddressName, gasToUse)

	owner, err := blockchain.GetOwnerAddress()
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	managedType.SetBytes(destinationHandle, owner)
}

// ManagedCaller VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedCaller(destinationHandle int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCaller
	metering.UseGasAndAddTracedGas(managedCallerName, gasToUse)

	caller := runtime.GetVMInput().CallerAddr
	managedType.SetBytes(destinationHandle, caller)
}

// ManagedSignalError VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedSignalError(errHandle int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(managedSignalErrorName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.SignalError
	metering.UseAndTraceGas(gasToUse)

	errBytes, err := managedType.GetBytes(errHandle)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBytes(errBytes)

	gasToUse = metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(len(errBytes))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		_ = context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	runtime.SignalUserError(string(errBytes))
}

// ManagedWriteLog VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedWriteLog(
	topicsHandle int32,
	dataHandle int32,
) {
	runtime := context.GetRuntimeContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()
	metering.StartGasTracing(managedWriteLogName)

	topics, sumOfTopicByteLengths, err := managedType.ReadManagedVecOfManagedBuffers(topicsHandle)
	if context.WithFault(err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	dataBytes, err := managedType.GetBytes(dataHandle)
	if context.WithFault(err, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBytes(dataBytes)
	dataByteLen := uint64(len(dataBytes))

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Log
	gasForData := math.MulUint64(
		metering.GasSchedule().BaseOperationCost.DataCopyPerByte,
		sumOfTopicByteLengths+dataByteLen)
	gasToUse = math.AddUint64(gasToUse, gasForData)
	metering.UseAndTraceGas(gasToUse)

	output.WriteLog(runtime.GetContextAddress(), topics, dataBytes)
}

// ManagedGetOriginalTxHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedGetOriginalTxHash(resultHandle int32) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetOriginalTxHash
	metering.UseGasAndAddTracedGas(managedGetOriginalTxHashName, gasToUse)

	managedType.SetBytes(resultHandle, runtime.GetOriginalTxHash())
}

// ManagedGetStateRootHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedGetStateRootHash(resultHandle int32) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetStateRootHash
	metering.UseGasAndAddTracedGas(managedGetStateRootHashName, gasToUse)

	managedType.SetBytes(resultHandle, blockchain.GetStateRootHash())
}

// ManagedGetBlockRandomSeed VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedGetBlockRandomSeed(resultHandle int32) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRandomSeed
	metering.UseGasAndAddTracedGas(managedGetBlockRandomSeedName, gasToUse)

	managedType.SetBytes(resultHandle, blockchain.CurrentRandomSeed())
}

// ManagedGetPrevBlockRandomSeed VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedGetPrevBlockRandomSeed(resultHandle int32) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRandomSeed
	metering.UseGasAndAddTracedGas(managedGetPrevBlockRandomSeedName, gasToUse)

	managedType.SetBytes(resultHandle, blockchain.LastRandomSeed())
}

// ManagedGetReturnData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedGetReturnData(resultID int32, resultHandle int32) {
	runtime := context.GetRuntimeContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetReturnData
	metering.UseGasAndAddTracedGas(managedGetReturnDataName, gasToUse)

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) || resultID < 0 {
		_ = context.WithFault(vmhost.ErrArgOutOfRange, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	managedType.SetBytes(resultHandle, returnData[resultID])
}

// ManagedGetMultiESDTCallValue VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedGetMultiESDTCallValue(multiCallValueHandle int32) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	metering.UseGasAndAddTracedGas(managedGetMultiESDTCallValueName, gasToUse)

	esdtTransfers := runtime.GetVMInput().ESDTTransfers
	multiCallBytes := writeESDTTransfersToBytes(managedType, esdtTransfers)
	managedType.ConsumeGasForBytes(multiCallBytes)

	managedType.SetBytes(multiCallValueHandle, multiCallBytes)
}

// ManagedGetESDTBalance VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedGetESDTBalance(addressHandle int32, tokenIDHandle int32, nonce int64, valueHandle int32) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	blockchain := context.GetBlockchainContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	metering.UseGasAndAddTracedGas(managedGetESDTBalanceName, gasToUse)

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		_ = context.WithFault(vmhost.ErrArgOutOfRange, runtime.BaseOpsErrorShouldFailExecution())
		return
	}
	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		_ = context.WithFault(vmhost.ErrArgOutOfRange, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	esdtToken, err := blockchain.GetESDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		_ = context.WithFault(vmhost.ErrArgOutOfRange, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	value := managedType.GetBigIntOrCreate(valueHandle)
	value.Set(esdtToken.Value)
}

// ManagedGetESDTTokenData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedGetESDTTokenData(
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64,
	valueHandle, propertiesHandle, hashHandle, nameHandle, attributesHandle, creatorHandle, royaltiesHandle, urisHandle int32) {
	host := context.GetVMHost()
	ManagedGetESDTTokenDataWithHost(
		host,
		addressHandle,
		tokenIDHandle,
		nonce,
		valueHandle, propertiesHandle, hashHandle, nameHandle, attributesHandle, creatorHandle, royaltiesHandle, urisHandle)

}

func ManagedGetESDTTokenDataWithHost(
	host vmhost.VMHost,
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64,
	valueHandle, propertiesHandle, hashHandle, nameHandle, attributesHandle, creatorHandle, royaltiesHandle, urisHandle int32) {
	runtime := host.Runtime()
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedGetESDTTokenDataName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	metering.UseAndTraceGas(gasToUse)

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		_ = WithFaultAndHost(host, vmhost.ErrArgOutOfRange, runtime.BaseOpsErrorShouldFailExecution())
		return
	}
	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		_ = WithFaultAndHost(host, vmhost.ErrArgOutOfRange, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	esdtToken, err := blockchain.GetESDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		_ = WithFaultAndHost(host, vmhost.ErrArgOutOfRange, runtime.BaseOpsErrorShouldFailExecution())
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

// ManagedAsyncCall VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedAsyncCall(
	destHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32) {
	host := context.GetVMHost()
	ManagedAsyncCallWithHost(
		host,
		destHandle,
		valueHandle,
		functionHandle,
		argumentsHandle)
}

func ManagedAsyncCallWithHost(
	host vmhost.VMHost,
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
	gasToUse := gasSchedule.BaseOpsAPICost.AsyncCallStep
	metering.UseAndTraceGas(gasToUse)

	vmInput, err := readDestinationFunctionArguments(host, destHandle, functionHandle, argumentsHandle)
	if WithFaultAndHost(host, err, host.Runtime().BaseOpsErrorShouldFailExecution()) {
		return
	}

	data := makeCrossShardCallFromInput(vmInput.function, vmInput.arguments)

	value, err := managedType.GetBigInt(valueHandle)
	if err != nil {
		_ = WithFaultAndHost(host, vmhost.ErrArgOutOfRange, host.Runtime().BaseOpsErrorShouldFailExecution())
		return
	}

	gasToUse = math.MulUint64(gasSchedule.BaseOperationCost.DataCopyPerByte, uint64(len(data)))
	metering.UseAndTraceGas(gasToUse)

	err = async.RegisterLegacyAsyncCall(vmInput.destination, []byte(data), value.Bytes())
	if errors.Is(err, vmhost.ErrNotEnoughGas) {
		runtime.SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
		return
	}
	if WithFaultAndHost(host, err, host.Runtime().BaseOpsErrorShouldFailExecution()) {
		return
	}
}

// ManagedCreateAsyncCall VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedCreateAsyncCall(
	destHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	successOffset executor.MemPtr,
	successLength executor.MemLength,
	errorOffset executor.MemPtr,
	errorLength executor.MemLength,
	gas int64,
	extraGasForCallback int64,
	callbackClosureHandle int32,
) int32 {

	host := context.GetVMHost()
	runtime := host.Runtime()
	managedType := host.ManagedTypes()

	vmInput, err := readDestinationFunctionArguments(host, destHandle, functionHandle, argumentsHandle)
	if WithFaultAndHost(host, err, host.Runtime().BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	data := makeCrossShardCallFromInput(vmInput.function, vmInput.arguments)

	value, err := managedType.GetBigInt(valueHandle)
	if err != nil {
		_ = context.WithFault(vmhost.ErrArgOutOfRange, runtime.BaseOpsErrorShouldFailExecution())
		return 1
	}

	successFunc, err := context.MemLoad(successOffset, successLength)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	errorFunc, err := context.MemLoad(errorOffset, errorLength)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	callbackClosure, err := managedType.GetBytes(callbackClosureHandle)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
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

// ManagedGetCallbackClosure VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedGetCallbackClosure(
	callbackClosureHandle int32,
) {
	host := context.GetVMHost()
	GetCallbackClosureWithHost(host, callbackClosureHandle)
}

func GetCallbackClosureWithHost(
	host vmhost.VMHost,
	callbackClosureHandle int32,
) {
	runtime := host.Runtime()
	async := host.Async()
	metering := host.Metering()
	managedTypes := host.ManagedTypes()

	metering.StartGasTracing(managedGetCallbackClosure)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallbackClosure
	metering.UseAndTraceGas(gasToUse)

	callbackClosure, err := async.GetCallbackClosure()
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	managedTypes.SetBytes(callbackClosureHandle, callbackClosure)
}

// ManagedUpgradeFromSourceContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedUpgradeFromSourceContract(
	destHandle int32,
	gas int64,
	valueHandle int32,
	addressHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedUpgradeFromSourceContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	vmInput, err := readDestinationValueArguments(host, destHandle, valueHandle, argumentsHandle)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	sourceContractAddress, err := managedType.GetBytes(addressHandle)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
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

// ManagedUpgradeContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedUpgradeContract(
	destHandle int32,
	gas int64,
	valueHandle int32,
	codeHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedUpgradeContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	vmInput, err := readDestinationValueArguments(host, destHandle, valueHandle, argumentsHandle)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	code, err := managedType.GetBytes(codeHandle)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	lenReturnData := len(host.Output().ReturnData())

	upgradeContract(host, vmInput.destination, code, codeMetadata, vmInput.value.Bytes(), vmInput.arguments, gas)
	setReturnDataIfExists(host, lenReturnData, resultHandle)
}

// ManagedDeleteContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedDeleteContract(
	destHandle int32,
	gasLimit int64,
	argumentsHandle int32,
) {
	host := context.GetVMHost()
	ManagedDeleteContractWithHost(
		host,
		destHandle,
		gasLimit,
		argumentsHandle,
	)
}

func ManagedDeleteContractWithHost(
	host vmhost.VMHost,
	destHandle int32,
	gasLimit int64,
	argumentsHandle int32,
) {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(deleteContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	calledSCAddress, err := managedType.GetBytes(destHandle)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	data, _, err := managedType.ReadManagedVecOfManagedBuffers(argumentsHandle)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	deleteContract(
		host,
		calledSCAddress,
		data,
		gasLimit,
	)
}

// ManagedDeployFromSourceContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedDeployFromSourceContract(
	gas int64,
	valueHandle int32,
	addressHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultAddressHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedDeployFromSourceContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	vmInput, err := readDestinationValueArguments(host, addressHandle, valueHandle, argumentsHandle)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
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
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	managedType.SetBytes(resultAddressHandle, newAddress)
	setReturnDataIfExists(host, lenReturnData, resultHandle)

	return 0
}

// ManagedCreateContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedCreateContract(
	gas int64,
	valueHandle int32,
	codeHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultAddressHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedCreateContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	sender := runtime.GetContextAddress()
	value, err := managedType.GetBigInt(valueHandle)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	data, actualLen, err := managedType.ReadManagedVecOfManagedBuffers(argumentsHandle)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, actualLen)
	metering.UseAndTraceGas(gasToUse)

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	code, err := managedType.GetBytes(codeHandle)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	lenReturnData := len(host.Output().ReturnData())
	newAddress, err := createContract(sender, data, value, metering, gas, code, codeMetadata, host, runtime)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	managedType.SetBytes(resultAddressHandle, newAddress)
	setReturnDataIfExists(host, lenReturnData, resultHandle)

	return 0
}

func setReturnDataIfExists(
	host vmhost.VMHost,
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

// ManagedExecuteReadOnly VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedExecuteReadOnly(
	gas int64,
	addressHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteReadOnlyName)

	vmInput, err := readDestinationFunctionArguments(host, addressHandle, functionHandle, argumentsHandle)
	if WithFaultAndHost(host, err, host.Runtime().BaseOpsErrorShouldFailExecution()) {
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

// ManagedExecuteOnSameContext VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedExecuteOnSameContext(
	gas int64,
	addressHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteOnSameContextName)

	vmInput, err := readDestinationValueFunctionArguments(host, addressHandle, valueHandle, functionHandle, argumentsHandle)
	if WithFaultAndHost(host, err, host.Runtime().BaseOpsErrorShouldFailExecution()) {
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

// ManagedExecuteOnDestContext VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedExecuteOnDestContext(
	gas int64,
	addressHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteOnDestContextName)

	vmInput, err := readDestinationValueFunctionArguments(host, addressHandle, valueHandle, functionHandle, argumentsHandle)
	if WithFaultAndHost(host, err, host.Runtime().BaseOpsErrorShouldFailExecution()) {
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

// ManagedMultiTransferESDTNFTExecute VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedMultiTransferESDTNFTExecute(
	dstHandle int32,
	tokenTransfersHandle int32,
	gasLimit int64,
	functionHandle int32,
	argumentsHandle int32,
) int32 {
	host := context.GetVMHost()
	managedType := host.ManagedTypes()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(managedMultiTransferESDTNFTExecuteName)

	vmInput, err := readDestinationFunctionArguments(host, dstHandle, functionHandle, argumentsHandle)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	transfers, err := readESDTTransfers(managedType, tokenTransfersHandle)
	if WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
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

// ManagedTransferValueExecute VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedTransferValueExecute(
	dstHandle int32,
	valueHandle int32,
	gasLimit int64,
	functionHandle int32,
	argumentsHandle int32,
) int32 {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(managedTransferValueExecuteName)

	vmInput, err := readDestinationValueFunctionArguments(host, dstHandle, valueHandle, functionHandle, argumentsHandle)
	if WithFaultAndHost(host, err, host.Runtime().BaseOpsErrorShouldFailExecution()) {
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

// ManagedIsESDTFrozen VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedIsESDTFrozen(
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64) int32 {
	host := context.GetVMHost()
	return ManagedIsESDTFrozenWithHost(host, addressHandle, tokenIDHandle, nonce)
}

func ManagedIsESDTFrozenWithHost(
	host vmhost.VMHost,
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	metering.UseGasAndAddTracedGas(managedIsESDTFrozenName, gasToUse)

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		_ = WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution())
		return -1
	}
	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		_ = WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution())
		return -1
	}

	esdtToken, err := blockchain.GetESDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		_ = WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution())
		return -1
	}

	esdtUserData := builtInFunctions.ESDTUserMetadataFromBytes(esdtToken.Properties)
	if esdtUserData.Frozen {
		return 1
	}
	return 0
}

// ManagedIsESDTLimitedTransfer VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedIsESDTLimitedTransfer(tokenIDHandle int32) int32 {
	host := context.GetVMHost()
	return ManagedIsESDTLimitedTransferWithHost(host, tokenIDHandle)
}

func ManagedIsESDTLimitedTransferWithHost(host vmhost.VMHost, tokenIDHandle int32) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	metering.UseGasAndAddTracedGas(managedIsESDTLimitedTransferName, gasToUse)

	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		_ = WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution())
		return -1
	}

	if blockchain.IsLimitedTransfer(tokenID) {
		return 1
	}

	return 0
}

// ManagedIsESDTPaused VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedIsESDTPaused(tokenIDHandle int32) int32 {
	host := context.GetVMHost()
	return ManagedIsESDTPausedWithHost(host, tokenIDHandle)
}

func ManagedIsESDTPausedWithHost(host vmhost.VMHost, tokenIDHandle int32) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	metering.UseGasAndAddTracedGas(managedIsESDTPausedName, gasToUse)

	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		_ = WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution())
		return -1
	}

	if blockchain.IsPaused(tokenID) {
		return 1
	}

	return 0
}

// ManagedBufferToHex VMHooks implementation.
// @autogenerate(VMHooks)
func (context *ElrondApi) ManagedBufferToHex(sourceHandle int32, destHandle int32) {
	host := context.GetVMHost()
	ManagedBufferToHexWithHost(host, sourceHandle, destHandle)
}

func ManagedBufferToHexWithHost(host vmhost.VMHost, sourceHandle int32, destHandle int32) {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferSetBytes
	metering.UseGasAndAddTracedGas(managedBufferToHexName, gasToUse)

	mBuff, err := managedType.GetBytes(sourceHandle)
	if err != nil {
		WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	encoded := hex.EncodeToString(mBuff)
	managedType.SetBytes(destHandle, []byte(encoded))
}

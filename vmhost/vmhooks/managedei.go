package vmhooks

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/math"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

const (
	managedSCAddressName                         = "managedSCAddress"
	managedOwnerAddressName                      = "managedOwnerAddress"
	managedCallerName                            = "managedCaller"
	managedSignalErrorName                       = "managedSignalError"
	managedWriteLogName                          = "managedWriteLog"
	managedMultiTransferESDTNFTExecuteName       = "managedMultiTransferESDTNFTExecute"
	managedTransferValueExecuteName              = "managedTransferValueExecute"
	managedExecuteOnDestContextName              = "managedExecuteOnDestContext"
	managedExecuteOnDestContextWithReturnName    = "managedExecuteOnDestContextWithReturn"
	managedExecuteOnDestContextByCallerName      = "managedExecuteOnDestContextByCaller"
	managedExecuteOnSameContextName              = "managedExecuteOnSameContext"
	managedExecuteReadOnlyName                   = "managedExecuteReadOnly"
	managedCreateContractName                    = "managedCreateContract"
	managedDeployFromSourceContractName          = "managedDeployFromSourceContract"
	managedUpgradeContractName                   = "managedUpgradeContract"
	managedUpgradeFromSourceContractName         = "managedUpgradeFromSourceContract"
	managedAsyncCallName                         = "managedAsyncCall"
	managedCreateAsyncCallName                   = "managedCreateAsyncCall"
	managedGetCallbackClosure                    = "managedGetCallbackClosure"
	managedGetMultiESDTCallValueName             = "managedGetMultiESDTCallValue"
	managedGetAllTransfersCallValue              = "managedGetAllTransfersCallValue"
	managedGetESDTBalanceName                    = "managedGetESDTBalance"
	managedGetESDTTokenDataName                  = "managedGetESDTTokenData"
	managedGetESDTTokenTypeName                  = "managedGetESDTTokenType"
	managedGetReturnDataName                     = "managedGetReturnData"
	managedGetPrevBlockRandomSeedName            = "managedGetPrevBlockRandomSeed"
	managedGetBlockRandomSeedName                = "managedGetBlockRandomSeed"
	managedGetStateRootHashName                  = "managedGetStateRootHash"
	managedGetOriginalTxHashName                 = "managedGetOriginalTxHash"
	managedIsESDTFrozenName                      = "managedIsESDTFrozen"
	managedIsESDTLimitedTransferName             = "managedIsESDTLimitedTransfer"
	managedIsESDTPausedName                      = "managedIsESDTPaused"
	managedBufferToHexName                       = "managedBufferToHex"
	managedGetCodeMetadataName                   = "managedGetCodeMetadata"
	managedGetCodeHashName                       = "managedGetCodeHash"
	managedIsBuiltinFunction                     = "managedIsBuiltinFunction"
	managedMultiTransferESDTNFTExecuteByUser     = "managedMultiTransferESDTNFTExecuteByUser"
	managedMultiTransferESDTNFTExecuteWithReturn = "managedMultiTransferESDTNFTExecuteByWithReturn"
)

const EGLDTokenName = "EGLD-000000" // TODO: maybe move to core?
const safeAttributeLength = 100

// ManagedSCAddress VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedSCAddress(destinationHandle int32) error {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetSCAddress
	err := metering.UseGasBoundedAndAddTracedGas(managedSCAddressName, gasToUse)
	if err != nil {
		return err
	}

	scAddress := runtime.GetContextAddress()

	managedType.SetBytes(destinationHandle, scAddress)
	return nil
}

// ManagedOwnerAddress VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedOwnerAddress(destinationHandle int32) error {
	managedType := context.GetManagedTypesContext()
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetOwnerAddress
	err := metering.UseGasBoundedAndAddTracedGas(managedOwnerAddressName, gasToUse)
	if err != nil {
		return err
	}

	owner, err := blockchain.GetOwnerAddress()
	if err != nil {
		return err
	}

	managedType.SetBytes(destinationHandle, owner)
	return nil
}

// ManagedCaller VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedCaller(destinationHandle int32) error {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCaller
	err := metering.UseGasBoundedAndAddTracedGas(managedCallerName, gasToUse)
	if err != nil {
		return err
	}

	caller := runtime.GetVMInput().CallerAddr
	managedType.SetBytes(destinationHandle, caller)
	return nil
}

// ManagedGetOriginalCallerAddr VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetOriginalCallerAddr(destinationHandle int32) error {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCaller
	err := metering.UseGasBoundedAndAddTracedGas(managedCallerName, gasToUse)
	if err != nil {
		return err
	}

	caller := runtime.GetVMInput().OriginalCallerAddr
	managedType.SetBytes(destinationHandle, caller)
	return nil
}

// ManagedGetRelayerAddr VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetRelayerAddr(destinationHandle int32) error {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCaller
	err := metering.UseGasBoundedAndAddTracedGas(managedCallerName, gasToUse)
	if err != nil {
		return err
	}

	caller := runtime.GetVMInput().RelayerAddr
	managedType.SetBytes(destinationHandle, caller)
	return nil
}

// ManagedSignalError VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedSignalError(errHandle int32) error {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(managedSignalErrorName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.SignalError
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	errBytes, err := managedType.GetBytes(errHandle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBytes(errBytes)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	gasToUse = metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(len(errBytes))
	err = metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	return vmhost.ExitUserError(string(errBytes))
}

// ManagedWriteLog VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedWriteLog(
	topicsHandle int32,
	dataHandle int32,
) error {
	runtime := context.GetRuntimeContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()
	metering.StartGasTracing(managedWriteLogName)

	topics, sumOfTopicByteLengths, err := managedType.ReadManagedVecOfManagedBuffers(topicsHandle)
	if err != nil {
		return err
	}

	dataBytes, err := managedType.GetBytes(dataHandle)
	if err != nil {
		return err
	}

	err = managedType.ConsumeGasForBytes(dataBytes)
	if err != nil {
		return err
	}

	dataByteLen := uint64(len(dataBytes))

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Log
	gasForData := math.MulUint64(
		metering.GasSchedule().BaseOperationCost.DataCopyPerByte,
		sumOfTopicByteLengths+dataByteLen)
	gasToUse = math.AddUint64(gasToUse, gasForData)
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	output.WriteLog(runtime.GetContextAddress(), topics, [][]byte{dataBytes})
	return nil
}

// ManagedGetOriginalTxHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetOriginalTxHash(resultHandle int32) error {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetOriginalTxHash
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	managedType.SetBytes(resultHandle, runtime.GetOriginalTxHash())
	return nil
}

// ManagedGetStateRootHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetStateRootHash(resultHandle int32) error {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetStateRootHash
	err := metering.UseGasBoundedAndAddTracedGas(managedGetStateRootHashName, gasToUse)
	if err != nil {
		return err
	}

	managedType.SetBytes(resultHandle, blockchain.GetStateRootHash())
	return nil
}

// ManagedGetBlockRandomSeed VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetBlockRandomSeed(resultHandle int32) error {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRandomSeed
	err := metering.UseGasBoundedAndAddTracedGas(managedGetBlockRandomSeedName, gasToUse)
	if err != nil {
		return err
	}

	managedType.SetBytes(resultHandle, blockchain.CurrentRandomSeed())
	return nil
}

// ManagedGetPrevBlockRandomSeed VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetPrevBlockRandomSeed(resultHandle int32) error {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRandomSeed
	err := metering.UseGasBoundedAndAddTracedGas(managedGetPrevBlockRandomSeedName, gasToUse)
	if err != nil {
		return err
	}

	managedType.SetBytes(resultHandle, blockchain.LastRandomSeed())
	return nil
}

// ManagedGetReturnData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetReturnData(resultID int32, resultHandle int32) error {
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetReturnData
	err := metering.UseGasBoundedAndAddTracedGas(managedGetReturnDataName, gasToUse)
	if err != nil {
		return err
	}

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) || resultID < 0 {
		return vmhost.ErrArgOutOfRange
	}

	managedType.SetBytes(resultHandle, returnData[resultID])
	return nil
}

// ManagedGetMultiESDTCallValue VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetMultiESDTCallValue(multiCallValueHandle int32) error {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(managedGetMultiESDTCallValueName, gasToUse)
	if err != nil {
		return err
	}

	esdtTransfers := runtime.GetVMInput().ESDTTransfers
	multiCallBytes := writeESDTTransfersToBytes(managedType, esdtTransfers)
	err = managedType.ConsumeGasForBytes(multiCallBytes)
	if err != nil {
		return err
	}

	managedType.SetBytes(multiCallValueHandle, multiCallBytes)
	return nil
}

// ManagedGetAllTransfersCallValue VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetAllTransfersCallValue(transferCallValuesListHandle int32) error {
	host := context.GetVMHost()
	managedType := host.ManagedTypes()

	allTransfers, err := ManagedGetAllTransfersCallValueTyped(host)
	if err != nil {
		return err
	}

	allTransfersBytes := writeESDTTransfersToBytes(managedType, allTransfers)
	err = managedType.ConsumeGasForBytes(allTransfersBytes)
	if err != nil {
		return err
	}

	managedType.SetBytes(transferCallValuesListHandle, allTransfersBytes)
	return nil
}

// ManagedGetAllTransfersCallValueTyped returns a combined list of all transfers (ESDT and EGLD)
func ManagedGetAllTransfersCallValueTyped(
	host vmhost.VMHost,
) ([]*vmcommon.ESDTTransfer, error) {
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(managedGetAllTransfersCallValue, gasToUse)
	if err != nil {
		return nil, err
	}

	input := runtime.GetVMInput()
	egldCallValue := input.CallValue
	hasCallValue := egldCallValue.Sign() > 0

	if hasCallValue {
		return []*vmcommon.ESDTTransfer{
			{
				ESDTValue:      egldCallValue,
				ESDTTokenName:  []byte(EGLDTokenName),
				ESDTTokenType:  uint32(core.Fungible),
				ESDTTokenNonce: 0,
			},
		}, nil
	}

	return input.ESDTTransfers, nil
}

// ManagedGetBackTransfers VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetBackTransfers(esdtTransfersValueHandle int32, egldValueHandle int32) error {
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(managedGetMultiESDTCallValueName, gasToUse)
	if err != nil {
		return err
	}

	esdtTransfers, transferValue := managedType.GetBackTransfers()
	multiCallBytes := writeESDTTransfersToBytes(managedType, esdtTransfers)
	err = managedType.ConsumeGasForBytes(multiCallBytes)
	if err != nil {
		return err
	}

	managedType.SetBytes(esdtTransfersValueHandle, multiCallBytes)
	egldValue := managedType.GetBigIntOrCreate(egldValueHandle)
	egldValue.SetBytes(transferValue.Bytes())
	return nil
}

// ManagedGetESDTBalance VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetESDTBalance(addressHandle int32, tokenIDHandle int32, nonce int64, valueHandle int32) error {
	metering := context.GetMeteringContext()
	blockchain := context.GetBlockchainContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	err := metering.UseGasBoundedAndAddTracedGas(managedGetESDTBalanceName, gasToUse)
	if err != nil {
		return err
	}

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		return vmhost.ErrArgOutOfRange
	}
	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		return vmhost.ErrArgOutOfRange
	}

	esdtToken, err := blockchain.GetESDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		return vmhost.ErrArgOutOfRange
	}

	value := managedType.GetBigIntOrCreate(valueHandle)
	value.Set(esdtToken.Value)
	return nil
}

// ManagedGetESDTTokenData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetESDTTokenData(
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64,
	valueHandle, propertiesHandle, hashHandle, nameHandle, attributesHandle, creatorHandle, royaltiesHandle, urisHandle int32) error {
	host := context.GetVMHost()
	return ManagedGetESDTTokenDataWithHost(
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
	valueHandle, propertiesHandle, hashHandle, nameHandle, attributesHandle, creatorHandle, royaltiesHandle, urisHandle int32) error {
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedGetESDTTokenDataName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		return vmhost.ErrArgOutOfRange
	}
	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		return vmhost.ErrArgOutOfRange
	}

	esdtToken, err := blockchain.GetESDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		return vmhost.ErrArgOutOfRange
	}

	value := managedType.GetBigIntOrCreate(valueHandle)
	value.Set(esdtToken.Value)

	managedType.SetBytes(propertiesHandle, esdtToken.Properties)
	if esdtToken.TokenMetaData != nil {
		managedType.SetBytes(hashHandle, esdtToken.TokenMetaData.Hash)
		err = managedType.ConsumeGasForBytes(esdtToken.TokenMetaData.Hash)
		if err != nil {
			return err
		}
		managedType.SetBytes(nameHandle, esdtToken.TokenMetaData.Name)
		err = managedType.ConsumeGasForBytes(esdtToken.TokenMetaData.Name)
		if err != nil {
			return err
		}
		managedType.SetBytes(attributesHandle, esdtToken.TokenMetaData.Attributes)
		err = managedType.ConsumeGasForBytes(esdtToken.TokenMetaData.Attributes)
		if err != nil {
			return err
		}
		managedType.SetBytes(creatorHandle, esdtToken.TokenMetaData.Creator)
		err = managedType.ConsumeGasForBytes(esdtToken.TokenMetaData.Creator)
		if err != nil {
			return err
		}
		royalties := managedType.GetBigIntOrCreate(royaltiesHandle)
		royalties.SetUint64(uint64(esdtToken.TokenMetaData.Royalties))

		err = managedType.WriteManagedVecOfManagedBuffers(esdtToken.TokenMetaData.URIs, urisHandle)
		if err != nil {
			return err
		}
	}

	return nil
}

// ManagedGetESDTTokenData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetESDTTokenType(
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64,
	typeHandle int32) error {
	host := context.GetVMHost()
	return ManagedGetESDTTokenTypeWithHost(
		host,
		addressHandle,
		tokenIDHandle,
		nonce,
		typeHandle)
}

// ManagedGetESDTTokenTypeWithHost implements the ManagedGetESDTTokenType VMHook logic.
func ManagedGetESDTTokenTypeWithHost(
	host vmhost.VMHost,
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64,
	typeHandle int32) error {
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedGetESDTTokenTypeName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		return vmhost.ErrArgOutOfRange
	}
	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		return vmhost.ErrArgOutOfRange
	}

	esdtToken, err := blockchain.GetESDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		return vmhost.ErrArgOutOfRange
	}

	esdtType := managedType.GetBigIntOrCreate(typeHandle)
	esdtType.SetUint64(uint64(esdtToken.Type))
	return nil
}

// ManagedAsyncCall VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedAsyncCall(
	destHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32) error {
	host := context.GetVMHost()
	return ManagedAsyncCallWithHost(
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
	argumentsHandle int32) error {
	runtime := host.Runtime()
	async := host.Async()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedAsyncCallName)

	gasSchedule := metering.GasSchedule()
	gasToUse := gasSchedule.BaseOpsAPICost.AsyncCallStep
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	vmInput, err := readDestinationFunctionArguments(host, destHandle, functionHandle, argumentsHandle)
	if err != nil {
		return err
	}

	data := makeCrossShardCallFromInput(vmInput.function, vmInput.arguments)

	value, err := managedType.GetBigInt(valueHandle)
	if err != nil {
		return vmhost.ErrArgOutOfRange
	}

	gasToUse = math.MulUint64(gasSchedule.BaseOperationCost.DataCopyPerByte, uint64(len(data)))
	err = metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	err = async.RegisterLegacyAsyncCall(vmInput.destination, []byte(data), value.Bytes())
	if errors.Is(err, vmhost.ErrNotEnoughGas) {
		return vmhost.ExitOutOfGas()
	}
	if err != nil {
		return err
	}
	return nil
}

// ManagedCreateAsyncCall VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedCreateAsyncCall(
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
) (int32, error) {

	host := context.GetVMHost()
	managedType := host.ManagedTypes()

	vmInput, err := readDestinationFunctionArguments(host, destHandle, functionHandle, argumentsHandle)
	if err != nil {
		return 1, err
	}

	data := makeCrossShardCallFromInput(vmInput.function, vmInput.arguments)

	value, err := managedType.GetBigInt(valueHandle)
	if err != nil {
		return 1, vmhost.ErrArgOutOfRange
	}

	successFunc, err := context.MemLoad(successOffset, successLength)
	if err != nil {
		return 1, err
	}

	errorFunc, err := context.MemLoad(errorOffset, errorLength)
	if err != nil {
		return 1, err
	}

	callbackClosure, err := managedType.GetBytes(callbackClosureHandle)
	if err != nil {
		return 1, err
	}

	return CreateAsyncCallWithTypedArgs(host,
		vmInput.destination,
		value.Bytes(),
		[]byte(data),
		successFunc,
		errorFunc,
		gas,
		extraGasForCallback,
		callbackClosure), nil
}

// ManagedGetCallbackClosure VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetCallbackClosure(
	callbackClosureHandle int32,
) error {
	host := context.GetVMHost()
	return GetCallbackClosureWithHost(host, callbackClosureHandle)
}

func GetCallbackClosureWithHost(
	host vmhost.VMHost,
	callbackClosureHandle int32,
) error {
	async := host.Async()
	metering := host.Metering()
	managedTypes := host.ManagedTypes()

	metering.StartGasTracing(managedGetCallbackClosure)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallbackClosure
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	callbackClosure, err := async.GetCallbackClosure()
	if err != nil {
		return err
	}

	managedTypes.SetBytes(callbackClosureHandle, callbackClosure)
	return nil
}

// ManagedUpgradeFromSourceContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedUpgradeFromSourceContract(
	destHandle int32,
	gas int64,
	valueHandle int32,
	addressHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) error {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedUpgradeFromSourceContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	vmInput, err := readDestinationValueArguments(host, destHandle, valueHandle, argumentsHandle)
	if err != nil {
		return err
	}

	sourceContractAddress, err := managedType.GetBytes(addressHandle)
	if err != nil {
		return err
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if err != nil {
		return err
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
	err = setReturnDataIfExists(host, lenReturnData, resultHandle)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}
	return nil
}

// ManagedUpgradeContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedUpgradeContract(
	destHandle int32,
	gas int64,
	valueHandle int32,
	codeHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) error {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedUpgradeContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	vmInput, err := readDestinationValueArguments(host, destHandle, valueHandle, argumentsHandle)
	if err != nil {
		return err
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if err != nil {
		return err
	}

	code, err := managedType.GetBytes(codeHandle)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	lenReturnData := len(host.Output().ReturnData())

	upgradeContract(host, vmInput.destination, code, codeMetadata, vmInput.value.Bytes(), vmInput.arguments, gas)
	err = setReturnDataIfExists(host, lenReturnData, resultHandle)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}
	return nil
}

// ManagedDeleteContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedDeleteContract(
	destHandle int32,
	gasLimit int64,
	argumentsHandle int32,
) error {
	host := context.GetVMHost()
	return ManagedDeleteContractWithHost(
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
) error {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(deleteContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	calledSCAddress, err := managedType.GetBytes(destHandle)
	if err != nil {
		return err
	}

	data, _, err := managedType.ReadManagedVecOfManagedBuffers(argumentsHandle)
	if err != nil {
		return err
	}

	deleteContract(
		host,
		calledSCAddress,
		data,
		gasLimit,
	)
	return nil
}

// ManagedDeployFromSourceContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedDeployFromSourceContract(
	gas int64,
	valueHandle int32,
	addressHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultAddressHandle int32,
	resultHandle int32,
) (int32, error) {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedDeployFromSourceContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return -1, err
	}

	vmInput, err := readDestinationValueArguments(host, addressHandle, valueHandle, argumentsHandle)
	if err != nil {
		return 1, err
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if err != nil {
		return 1, err
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
	if err != nil {
		return 1, err
	}

	managedType.SetBytes(resultAddressHandle, newAddress)
	err = setReturnDataIfExists(host, lenReturnData, resultHandle)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return 1, err
	}

	return 0, nil
}

// ManagedCreateContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedCreateContract(
	gas int64,
	valueHandle int32,
	codeHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultAddressHandle int32,
	resultHandle int32,
) (int32, error) {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedCreateContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return -1, err
	}

	sender := runtime.GetContextAddress()
	value, err := managedType.GetBigInt(valueHandle)
	if err != nil {
		return 1, err
	}

	data, actualLen, err := managedType.ReadManagedVecOfManagedBuffers(argumentsHandle)
	if err != nil {
		return 1, err
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, actualLen)
	err = metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return -1, err
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if err != nil {
		return 1, err
	}

	code, err := managedType.GetBytes(codeHandle)
	if err != nil {
		return 1, err
	}

	lenReturnData := len(host.Output().ReturnData())
	newAddress, err := createContract(sender, data, value, gas, code, codeMetadata, host, CreateContract)
	if err != nil {
		return 1, err
	}

	managedType.SetBytes(resultAddressHandle, newAddress)
	err = setReturnDataIfExists(host, lenReturnData, resultHandle)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return 1, err
	}

	return 0, nil
}

func setReturnDataIfExists(
	host vmhost.VMHost,
	oldLen int,
	resultHandle int32,
) error {
	returnData := host.Output().ReturnData()
	if len(returnData) > oldLen {
		return host.ManagedTypes().WriteManagedVecOfManagedBuffers(returnData[oldLen:], resultHandle)
	}

	host.ManagedTypes().SetBytes(resultHandle, make([]byte, 0))
	return nil
}

// ManagedExecuteReadOnly VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedExecuteReadOnly(
	gas int64,
	addressHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) (int32, error) {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteReadOnlyName)

	vmInput, err := readDestinationFunctionArguments(host, addressHandle, functionHandle, argumentsHandle)
	if err != nil {
		return -1, err
	}

	lenReturnData := len(host.Output().ReturnData())
	returnVal := ExecuteReadOnlyWithTypedArguments(
		host,
		gas,
		[]byte(vmInput.function),
		vmInput.destination,
		vmInput.arguments,
	)
	err = setReturnDataIfExists(host, lenReturnData, resultHandle)
	if err != nil && host.Runtime().UseGasBoundedShouldFailExecution() {
		return -1, err
	}

	return returnVal, nil
}

// ManagedExecuteOnSameContext VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedExecuteOnSameContext(
	gas int64,
	addressHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) (int32, error) {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteOnSameContextName)

	vmInput, err := readDestinationValueFunctionArguments(host, addressHandle, valueHandle, functionHandle, argumentsHandle)
	if err != nil {
		return -1, err
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
	err = setReturnDataIfExists(host, lenReturnData, resultHandle)
	if err != nil && host.Runtime().UseGasBoundedShouldFailExecution() {
		return -1, err
	}

	return returnVal, nil
}

// ManagedExecuteOnDestContext VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedExecuteOnDestContext(
	gas int64,
	addressHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) (int32, error) {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteOnDestContextName)

	vmInput, err := readDestinationValueFunctionArguments(host, addressHandle, valueHandle, functionHandle, argumentsHandle)
	if err != nil {
		return -1, err
	}

	lenReturnData := len(host.Output().ReturnData())
	returnVal := ExecuteOnDestContextWithTypedArgs(
		host,
		gas,
		vmInput.value,
		[]byte(vmInput.function),
		vmInput.destination,
		vmInput.arguments,
		true,
	)
	err = setReturnDataIfExists(host, lenReturnData, resultHandle)
	if err != nil && host.Runtime().UseGasBoundedShouldFailExecution() {
		return -1, err
	}

	return returnVal, nil
}

// ManagedExecuteOnDestContextWithErrorReturn VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedExecuteOnDestContextWithErrorReturn(
	gas int64,
	addressHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) (int32, error) {
	host := context.GetVMHost()

	vmInput, err := readDestinationValueFunctionArguments(host, addressHandle, valueHandle, functionHandle, argumentsHandle)
	if err != nil {
		return -1, err
	}

	return ManagedExecuteOnDestContextWithErrorReturnWithHost(host, gas, vmInput.value, vmInput.function, vmInput.destination, vmInput.arguments, resultHandle)
}

// ManagedExecuteOnDestContextWithErrorReturnWithHost - execute on dest context and return error instead of failing execution
func ManagedExecuteOnDestContextWithErrorReturnWithHost(
	host vmhost.VMHost,
	gas int64,
	value *big.Int,
	function string,
	destination []byte,
	arguments [][]byte,
	resultHandle int32,
) (int32, error) {
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteOnDestContextWithReturnName)

	lenReturnData := len(host.Output().ReturnData())
	returnVal := ExecuteOnDestContextWithTypedArgs(
		host,
		gas,
		value,
		[]byte(function),
		destination,
		arguments,
		false,
	)

	err := setReturnDataIfExists(host, lenReturnData, resultHandle)
	if err != nil {
		return -1, err
	}

	return returnVal, nil
}

// ManagedMultiTransferESDTNFTExecute VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMultiTransferESDTNFTExecute(
	dstHandle int32,
	tokenTransfersHandle int32,
	gasLimit int64,
	functionHandle int32,
	argumentsHandle int32,
) (int32, error) {
	host := context.GetVMHost()
	managedType := host.ManagedTypes()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(managedMultiTransferESDTNFTExecuteName)

	vmInput, err := readDestinationFunctionArguments(host, dstHandle, functionHandle, argumentsHandle)
	if err != nil {
		return -1, err
	}

	transfers, err := readESDTTransfers(managedType, runtime, tokenTransfersHandle)
	if err != nil {
		return -1, err
	}

	return TransferESDTNFTExecuteWithTypedArgs(
		host,
		vmInput.destination,
		transfers,
		gasLimit,
		[]byte(vmInput.function),
		vmInput.arguments,
	), nil
}

// ManagedMultiTransferESDTNFTExecuteWithReturn VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMultiTransferESDTNFTExecuteWithReturn(
	dstHandle int32,
	tokenTransfersHandle int32,
	gasLimit int64,
	functionHandle int32,
	argumentsHandle int32,
) (int32, error) {
	host := context.GetVMHost()
	managedType := host.ManagedTypes()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(managedMultiTransferESDTNFTExecuteWithReturn)

	vmInput, err := readDestinationFunctionArguments(host, dstHandle, functionHandle, argumentsHandle)
	if err != nil {
		return -1, err
	}

	transfers, err := readESDTTransfers(managedType, runtime, tokenTransfersHandle)
	if err != nil {
		return -1, err
	}

	return TransferESDTNFTExecuteWithTypedArgsWithFailure(
		host,
		vmInput.destination,
		transfers,
		gasLimit,
		[]byte(vmInput.function),
		vmInput.arguments,
		false,
	), nil
}

// ManagedMultiTransferESDTNFTExecuteByUser VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMultiTransferESDTNFTExecuteByUser(
	userHandle int32,
	dstHandle int32,
	tokenTransfersHandle int32,
	gasLimit int64,
	functionHandle int32,
	argumentsHandle int32,
) (int32, error) {
	host := context.GetVMHost()
	managedType := host.ManagedTypes()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(managedMultiTransferESDTNFTExecuteByUser)

	if !host.IsAllowedToExecute(managedMultiTransferESDTNFTExecuteByUser) {
		return -1, vmhost.ErrOpcodeIsNotAllowed
	}

	user, err := managedType.GetBytes(userHandle)
	if err != nil {
		return -1, err
	}

	vmInput, err := readDestinationFunctionArguments(host, dstHandle, functionHandle, argumentsHandle)
	if err != nil {
		return -1, err
	}

	transfers, err := readESDTTransfers(managedType, runtime, tokenTransfersHandle)
	if err != nil {
		return -1, err
	}

	return TransferESDTNFTExecuteByUserWithTypedArgs(
		host,
		user,
		vmInput.destination,
		transfers,
		gasLimit,
		[]byte(vmInput.function),
		vmInput.arguments,
	), nil
}

// ManagedTransferValueExecute VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedTransferValueExecute(
	dstHandle int32,
	valueHandle int32,
	gasLimit int64,
	functionHandle int32,
	argumentsHandle int32,
) (int32, error) {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(managedTransferValueExecuteName)

	vmInput, err := readDestinationValueFunctionArguments(host, dstHandle, valueHandle, functionHandle, argumentsHandle)
	if err != nil {
		return -1, err
	}

	return TransferValueExecuteWithTypedArgs(
		host,
		vmInput.destination,
		vmInput.value,
		gasLimit,
		[]byte(vmInput.function),
		vmInput.arguments,
	), nil
}

// ManagedIsESDTFrozen VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedIsESDTFrozen(
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64) (int32, error) {
	host := context.GetVMHost()
	return ManagedIsESDTFrozenWithHost(host, addressHandle, tokenIDHandle, nonce)
}

func ManagedIsESDTFrozenWithHost(
	host vmhost.VMHost,
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64) (int32, error) {
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	err := metering.UseGasBoundedAndAddTracedGas(managedIsESDTFrozenName, gasToUse)
	if err != nil {
		return -1, err
	}

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		return -1, err
	}
	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		return -1, err
	}

	esdtToken, err := blockchain.GetESDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		return -1, err
	}

	esdtUserData := builtInFunctions.ESDTUserMetadataFromBytes(esdtToken.Properties)
	if esdtUserData.Frozen {
		return 1, nil
	}
	return 0, nil
}

// ManagedIsESDTLimitedTransfer VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedIsESDTLimitedTransfer(tokenIDHandle int32) (int32, error) {
	host := context.GetVMHost()
	return ManagedIsESDTLimitedTransferWithHost(host, tokenIDHandle)
}

func ManagedIsESDTLimitedTransferWithHost(host vmhost.VMHost, tokenIDHandle int32) (int32, error) {
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	err := metering.UseGasBoundedAndAddTracedGas(managedIsESDTLimitedTransferName, gasToUse)
	if err != nil {
		return -1, err
	}

	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		return -1, err
	}

	if blockchain.IsLimitedTransfer(tokenID) {
		return 1, nil
	}

	return 0, nil
}

// ManagedIsESDTPaused VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedIsESDTPaused(tokenIDHandle int32) (int32, error) {
	host := context.GetVMHost()
	return ManagedIsESDTPausedWithHost(host, tokenIDHandle)
}

func ManagedIsESDTPausedWithHost(host vmhost.VMHost, tokenIDHandle int32) (int32, error) {
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	err := metering.UseGasBoundedAndAddTracedGas(managedIsESDTPausedName, gasToUse)
	if err != nil {
		return -1, err
	}

	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		return -1, err
	}

	err = chargeGasForExtraLength(len(tokenID), host)
	if err != nil {
		return -1, err
	}

	if blockchain.IsPaused(tokenID) {
		return 1, nil
	}

	return 0, nil
}

// ManagedBufferToHex VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedBufferToHex(sourceHandle int32, destHandle int32) error {
	host := context.GetVMHost()
	return ManagedBufferToHexWithHost(host, sourceHandle, destHandle)
}

func ManagedBufferToHexWithHost(host vmhost.VMHost, sourceHandle int32, destHandle int32) error {
	metering := host.Metering()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferSetBytes
	err := metering.UseGasBoundedAndAddTracedGas(managedBufferToHexName, gasToUse)
	if err != nil {
		return err
	}

	mBuff, err := managedType.GetBytes(sourceHandle)
	if err != nil {
		return err
	}

	encoded := hex.EncodeToString(mBuff)
	managedType.SetBytes(destHandle, []byte(encoded))
	return nil
}

// ManagedGetCodeMetadata VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetCodeMetadata(addressHandle int32, responseHandle int32) error {
	host := context.GetVMHost()
	return ManagedGetCodeMetadataWithHost(host, addressHandle, responseHandle)
}

func ManagedGetCodeMetadataWithHost(host vmhost.VMHost, addressHandle int32, responseHandle int32) error {
	metering := host.Metering()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCodeMetadata
	err := metering.UseGasBoundedAndAddTracedGas(managedGetCodeMetadataName, gasToUse)
	if err != nil {
		return err
	}

	gasToUse = metering.GasSchedule().ManagedBufferAPICost.MBufferSetBytes
	err = metering.UseGasBoundedAndAddTracedGas(managedGetCodeMetadataName, gasToUse)
	if err != nil {
		return err
	}

	mBuffAddress, err := managedType.GetBytes(addressHandle)
	if err != nil {
		return err
	}

	contract, err := host.Blockchain().GetUserAccount(mBuffAddress)
	if err != nil || check.IfNil(contract) {
		return err
	}

	codeMetadata := contract.GetCodeMetadata()

	managedType.SetBytes(responseHandle, codeMetadata)
	return nil
}

// ManagedGetCodeHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetCodeHash(addressHandle int32, codeHashHandle int32) error {
	host := context.GetVMHost()
	managedType := host.ManagedTypes()

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		return err
	}

	codeHash, err := ManagedGetCodeHashTyped(host, address)
	if err != nil {
		return err
	}

	managedType.SetBytes(codeHashHandle, codeHash)
	return nil
}

// ManagedGetCodeHashWithHost returns the code hash at some address
func ManagedGetCodeHashTyped(
	host vmhost.VMHost,
	address []byte,
) ([]byte, error) {
	metering := host.Metering()
	blockchain := host.Blockchain()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCodeHash
	err := metering.UseGasBoundedAndAddTracedGas(managedGetCodeHashName, gasToUse)
	if err != nil {
		return nil, err
	}

	codeHash := blockchain.GetCodeHash(address)

	return codeHash, nil
}

// ManagedIsBuiltinFunction VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedIsBuiltinFunction(functionNameHandle int32) (int32, error) {
	host := context.GetVMHost()
	return ManagedIsBuiltinFunctionWithHost(host, functionNameHandle)
}

func ManagedIsBuiltinFunctionWithHost(host vmhost.VMHost, functionNameHandle int32) (int32, error) {
	metering := host.Metering()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.IsBuiltinFunction
	err := metering.UseGasBoundedAndAddTracedGas(managedIsBuiltinFunction, gasToUse)
	if err != nil {
		return -1, err
	}

	mBuffFunctionName, err := managedType.GetBytes(functionNameHandle)
	if err != nil {
		return -1, err
	}

	lenFuncName := len(mBuffFunctionName)
	err = chargeGasForExtraLength(lenFuncName, host)
	if err != nil {
		return -1, err
	}

	isBuiltinFunction := host.IsBuiltinFunctionName(string(mBuffFunctionName))
	if isBuiltinFunction {
		return 1, nil
	}

	return 0, nil
}

func chargeGasForExtraLength(lenAttribute int, host vmhost.VMHost) error {
	if !host.Runtime().AttributeExtraGasUsage() || lenAttribute < safeAttributeLength {
		return nil
	}

	metering := host.Metering()
	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(lenAttribute))
	return metering.UseGasBounded(gasToUse)
}

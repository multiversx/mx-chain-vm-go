package vmhooks

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/esdt"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	logger "github.com/multiversx/mx-chain-logger-go"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/math"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

const (
	getSCAddressName                 = "getSCAddress"
	getOwnerAddressName              = "getOwnerAddress"
	getShardOfAddressName            = "getShardOfAddress"
	isSmartContractName              = "isSmartContract"
	getExternalBalanceName           = "getExternalBalance"
	blockHashName                    = "blockHash"
	transferValueName                = "transferValue"
	transferESDTExecuteName          = "transferESDTExecute"
	transferESDTNFTExecuteName       = "transferESDTNFTExecute"
	multiTransferESDTNFTExecuteName  = "multiTransferESDTNFTExecute"
	transferValueExecuteName         = "transferValueExecute"
	createAsyncCallName              = "createAsyncCall"
	setAsyncGroupCallbackName        = "setAsyncGroupCallback"
	setAsyncContextCallbackName      = "setAsyncContextCallback"
	getArgumentLengthName            = "getArgumentLength"
	getArgumentName                  = "getArgument"
	getFunctionName                  = "getFunction"
	getNumArgumentsName              = "getNumArguments"
	storageStoreName                 = "storageStore"
	storageLoadLengthName            = "storageLoadLength"
	storageLoadName                  = "storageLoad"
	storageLoadFromAddressName       = "storageLoadFromAddress"
	getCallerName                    = "getCaller"
	checkNoPaymentName               = "checkNoPayment"
	callValueName                    = "callValue"
	getESDTValueName                 = "getESDTValue"
	getESDTTokenNameName             = "getESDTTokenName"
	getESDTTokenNonceName            = "getESDTTokenNonce"
	getESDTTokenTypeName             = "getESDTTokenType"
	getCallValueTokenNameName        = "getCallValueTokenName"
	getESDTValueByIndexName          = "getESDTValueByIndex"
	getESDTTokenNameByIndexName      = "getESDTTokenNameByIndex"
	getESDTTokenNonceByIndexName     = "getESDTTokenNonceByIndex"
	getESDTTokenTypeByIndexName      = "getESDTTokenTypeByIndex"
	getCallValueTokenNameByIndexName = "getCallValueTokenNameByIndex"
	getNumESDTTransfersName          = "getNumESDTTransfers"
	getCurrentESDTNFTNonceName       = "getCurrentESDTNFTNonce"
	writeLogName                     = "writeLog"
	writeEventLogName                = "writeEventLog"
	returnDataName                   = "returnData"
	signalErrorName                  = "signalError"
	getGasLeftName                   = "getGasLeft"
	getESDTBalanceName               = "getESDTBalance"
	getESDTNFTNameLengthName         = "getESDTNFTNameLength"
	getESDTNFTAttributeLengthName    = "getESDTNFTAttributeLength"
	getESDTNFTURILengthName          = "getESDTNFTURILength"
	getESDTTokenDataName             = "getESDTTokenData"
	getESDTLocalRolesName            = "getESDTLocalRoles"
	validateTokenIdentifierName      = "validateTokenIdentifier"
	executeOnDestContextName         = "executeOnDestContext"
	executeOnSameContextName         = "executeOnSameContext"
	executeReadOnlyName              = "executeReadOnly"
	createContractName               = "createContract"
	deployFromSourceContractName     = "deployFromSourceContract"
	upgradeContractName              = "upgradeContract"
	upgradeFromSourceContractName    = "upgradeFromSourceContract"
	deleteContractName               = "deleteContract"
	asyncCallName                    = "asyncCall"
	getNumReturnDataName             = "getNumReturnData"
	getReturnDataSizeName            = "getReturnDataSize"
	getReturnDataName                = "getReturnData"
	cleanReturnDataName              = "cleanReturnData"
	deleteFromReturnDataName         = "deleteFromReturnData"
	setStorageLockName               = "setStorageLock"
	getStorageLockName               = "getStorageLock"
	isStorageLockedName              = "isStorageLocked"
	clearStorageLockName             = "clearStorageLock"
	getBlockTimestampName            = "getBlockTimestamp"
	getBlockTimestampMsName          = "getBlockTimestampMs"
	getBlockNonceName                = "getBlockNonce"
	getBlockRoundName                = "getBlockRound"
	getBlockEpochName                = "getBlockEpoch"
	getBlockRandomSeedName           = "getBlockRandomSeed"
	getStateRootHashName             = "getStateRootHash"
	getPrevBlockTimestampName        = "getPrevBlockTimestamp"
	getPrevBlockTimestampMsName      = "getPrevBlockTimestampMs"
	getPrevBlockNonceName            = "getPrevBlockNonce"
	getPrevBlockRoundName            = "getPrevBlockRound"
	getPrevBlockEpochName            = "getPrevBlockEpoch"
	getPrevBlockRandomSeedName       = "getPrevBlockRandomSeed"
	getOriginalTxHashName            = "getOriginalTxHash"
	getCurrentTxHashName             = "getCurrentTxHash"
	getPrevTxHashName                = "getPrevTxHash"
	getBlockRoundTimeMsName          = "getBlockRoundTimeMs"
	epochStartBlockTimestampMsName   = "epochStartBlockTimestampMs"
	epochStartBlockNonceName         = "epochStartBlockNonce"
	epochStartBlockRoundName         = "epochStartBlockRound"
)

type CreateContractCallType int

const (
	CreateContract = iota
	DeployContract
)

// one memory page 64KB / 4 as min bytes per arguments
const maxNumArgumentsFromMemory = 16000000
const maxNumArgumentsFromMemoryCapped = 1024

var logEEI = logger.GetOrCreate("vm/eei")

func getESDTTransferFromInputFailIfWrongIndex(host vmhost.VMHost, index int32) (*vmcommon.ESDTTransfer, error) {
	esdtTransfers := host.Runtime().GetVMInput().ESDTTransfers
	if int32(len(esdtTransfers))-1 < index || index < 0 {
		return nil, vmhost.ErrInvalidTokenIndex
	}
	return esdtTransfers[index], nil
}

func failIfMoreThanOneESDTTransfer(context *VMHooksImpl) error {
	runtime := context.GetRuntimeContext()
	if len(runtime.GetVMInput().ESDTTransfers) > 1 {
		return vmhost.ErrTooManyESDTTransfers
	}
	return nil
}

// GetGasLeft VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetGasLeft() (int64, error) {
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetGasLeft
	err := metering.UseGasBoundedAndAddTracedGas(getGasLeftName, gasToUse)
	if err != nil {
		return 0, err
	}

	return int64(metering.GasLeft()), nil
}

// GetSCAddress VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetSCAddress(resultOffset executor.MemPtr) error {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetSCAddress
	err := metering.UseGasBoundedAndAddTracedGas(getSCAddressName, gasToUse)
	if err != nil {
		return err
	}

	owner := runtime.GetContextAddress()
	err = context.MemStore(resultOffset, owner)
	if err != nil {
		return err
	}
	return nil
}

// GetOwnerAddress VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetOwnerAddress(resultOffset executor.MemPtr) error {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetOwnerAddress
	err := metering.UseGasBoundedAndAddTracedGas(getOwnerAddressName, gasToUse)
	if err != nil {
		return err
	}

	owner, err := blockchain.GetOwnerAddress()
	if err != nil {
		return err
	}

	err = context.MemStore(resultOffset, owner)
	if err != nil {
		return err
	}
	return nil
}

// GetShardOfAddress VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetShardOfAddress(addressOffset executor.MemPtr) (int32, error) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetShardOfAddress
	err := metering.UseGasBoundedAndAddTracedGas(getShardOfAddressName, gasToUse)
	if err != nil {
		return -1, err
	}

	address, err := context.MemLoad(addressOffset, vmhost.AddressLen)
	if err != nil {
		return -1, err
	}

	return int32(blockchain.GetShardOfAddress(address)), nil
}

// IsSmartContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) IsSmartContract(addressOffset executor.MemPtr) (int32, error) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.IsSmartContract
	err := metering.UseGasBoundedAndAddTracedGas(isSmartContractName, gasToUse)
	if err != nil {
		return -1, err
	}

	address, err := context.MemLoad(addressOffset, vmhost.AddressLen)
	if err != nil {
		return -1, err
	}

	isSmartContract := blockchain.IsSmartContract(address)

	return int32(vmhost.BooleanToInt(isSmartContract)), nil
}

// SignalError VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SignalError(messageOffset executor.MemPtr, messageLength executor.MemLength) error {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(signalErrorName)

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(messageLength))
	gasToUse = math.AddUint64(gasToUse, metering.GasSchedule().BaseOpsAPICost.SignalError)

	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	message, err := context.MemLoad(messageOffset, messageLength)
	if err != nil {
		return err
	}

	return vmhost.ExitUserError(string(message))
}

// GetExternalBalance VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetExternalBalance(addressOffset executor.MemPtr, resultOffset executor.MemPtr) error {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	err := metering.UseGasBoundedAndAddTracedGas(getExternalBalanceName, gasToUse)
	if err != nil {
		return err
	}

	address, err := context.MemLoad(addressOffset, vmhost.AddressLen)
	if err != nil {
		return err
	}

	balance := blockchain.GetBalance(address)

	err = context.MemStore(resultOffset, balance)
	if err != nil {
		return err
	}
	return nil
}

// GetBlockHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetBlockHash(nonce int64, resultOffset executor.MemPtr) (int32, error) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockHash
	err := metering.UseGasBoundedAndAddTracedGas(blockHashName, gasToUse)
	if err != nil {
		return 1, err
	}

	hash := blockchain.BlockHash(uint64(nonce))
	err = context.MemStore(resultOffset, hash)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

func getESDTDataFromBlockchainHook(
	context *VMHooksImpl,
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

// GetESDTBalance VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetESDTBalance(
	addressOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	nonce int64,
	resultOffset executor.MemPtr,
) (int32, error) {
	metering := context.GetMeteringContext()
	metering.StartGasTracing(getESDTBalanceName)

	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)
	if err != nil {
		return -1, err
	}

	bytes := esdtData.Value.Bytes()
	err = context.MemStore(resultOffset, bytes)
	if err != nil {
		return -1, err
	}

	return int32(len(bytes)), nil
}

// GetESDTNFTNameLength VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetESDTNFTNameLength(
	addressOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	nonce int64,
) (int32, error) {
	metering := context.GetMeteringContext()
	metering.StartGasTracing(getESDTNFTNameLengthName)

	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)
	if err != nil {
		return -1, err
	}
	if esdtData == nil || esdtData.TokenMetaData == nil {
		return 0, vmhost.ErrNilESDTData
	}

	return int32(len(esdtData.TokenMetaData.Name)), nil
}

// GetESDTNFTAttributeLength VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetESDTNFTAttributeLength(
	addressOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	nonce int64,
) (int32, error) {
	metering := context.GetMeteringContext()
	metering.StartGasTracing(getESDTNFTAttributeLengthName)

	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)
	if err != nil {
		return -1, err
	}
	if esdtData == nil || esdtData.TokenMetaData == nil {
		return 0, vmhost.ErrNilESDTData
	}

	return int32(len(esdtData.TokenMetaData.Attributes)), nil
}

// GetESDTNFTURILength VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetESDTNFTURILength(
	addressOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	nonce int64,
) (int32, error) {
	metering := context.GetMeteringContext()
	metering.StartGasTracing(getESDTNFTURILengthName)

	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)
	if err != nil {
		return -1, err
	}
	if esdtData == nil || esdtData.TokenMetaData == nil {
		return 0, vmhost.ErrNilESDTData
	}
	if len(esdtData.TokenMetaData.URIs) == 0 {
		return 0, nil
	}

	return int32(len(esdtData.TokenMetaData.URIs[0])), nil
}

// GetESDTTokenData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetESDTTokenData(
	addressOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	nonce int64,
	valueHandle int32,
	propertiesOffset executor.MemPtr,
	hashOffset executor.MemPtr,
	nameOffset executor.MemPtr,
	attributesOffset executor.MemPtr,
	creatorOffset executor.MemPtr,
	royaltiesHandle int32,
	urisOffset executor.MemPtr,
) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(getESDTTokenDataName)

	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)
	if err != nil {
		return -1, err
	}

	value := managedType.GetBigIntOrCreate(valueHandle)
	value.Set(esdtData.Value)

	err = context.MemStore(propertiesOffset, esdtData.Properties)
	if err != nil {
		return -1, err
	}

	if esdtData.TokenMetaData == nil {
		return int32(len(esdtData.Value.Bytes())), nil
	}

	err = context.MemStore(hashOffset, esdtData.TokenMetaData.Hash)
	if err != nil {
		return -1, err
	}
	err = context.MemStore(nameOffset, esdtData.TokenMetaData.Name)
	if err != nil {
		return -1, err
	}
	err = context.MemStore(attributesOffset, esdtData.TokenMetaData.Attributes)
	if err != nil {
		return -1, err
	}
	err = context.MemStore(creatorOffset, esdtData.TokenMetaData.Creator)
	if err != nil {
		return -1, err
	}

	royalties := managedType.GetBigIntOrCreate(royaltiesHandle)
	royalties.SetUint64(uint64(esdtData.TokenMetaData.Royalties))

	if len(esdtData.TokenMetaData.URIs) > 0 {
		err = context.MemStore(urisOffset, esdtData.TokenMetaData.URIs[0])
		if err != nil {
			return -1, err
		}
	}
	return int32(len(esdtData.Value.Bytes())), nil
}

// GetESDTLocalRoles VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetESDTLocalRoles(tokenIdHandle int32) (int64, error) {
	managedType := context.GetManagedTypesContext()
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	tokenID, err := managedType.GetBytes(tokenIdHandle)
	if err != nil {
		return -1, err
	}

	esdtRoleKeyPrefix := []byte(core.ProtectedKeyPrefix + core.ESDTRoleIdentifier + core.ESDTKeyIdentifier)
	key := []byte(string(esdtRoleKeyPrefix) + string(tokenID))

	data, trieDepth, usedCache, err := storage.GetStorage(key)
	if err != nil {
		return -1, err
	}

	err = storage.UseGasForStorageLoad(
		storageLoadName,
		int64(trieDepth),
		metering.GasSchedule().BaseOpsAPICost.StorageLoad,
		usedCache)
	if err != nil {
		return -1, err
	}

	enableEpochsHandler := context.host.EnableEpochsHandler()
	return getESDTRoles(data, enableEpochsHandler.IsFlagEnabled(vmhost.CryptoOpcodesV2Flag)), nil
}

// ValidateTokenIdentifier VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ValidateTokenIdentifier(
	tokenIdHandle int32,
) (int32, error) {
	managedType := context.GetManagedTypesContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetArgument
	err := metering.UseGasBoundedAndAddTracedGas(validateTokenIdentifierName, gasToUse)
	if err != nil {
		return -1, err
	}

	tokenID, err := managedType.GetBytes(tokenIdHandle)
	if err != nil {
		return -1, err
	}

	if ValidateToken(tokenID) {
		return 1, nil
	} else {
		return 0, nil
	}

}

// TransferValue VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) TransferValue(
	destOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	dataOffset executor.MemPtr,
	length executor.MemLength) (int32, error) {

	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()
	metering.StartGasTracing(transferValueName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.TransferValue
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return 1, err
	}

	sender := runtime.GetContextAddress()
	dest, err := context.MemLoad(destOffset, vmhost.AddressLen)
	if err != nil {
		return 1, err
	}

	valueBytes, err := context.MemLoad(valueOffset, vmhost.BalanceLen)
	if err != nil {
		return 1, err
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(length))
	err = metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return 1, err
	}

	data, err := context.MemLoad(dataOffset, length)
	if err != nil {
		return 1, err
	}

	if host.IsBuiltinFunctionCall(data) {
		return 1, vmhost.ErrTransferValueOnESDTCall
	}

	err = output.Transfer(dest, sender, 0, 0, big.NewInt(0).SetBytes(valueBytes), nil, data, vm.DirectCall)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

type indirectContractCallArguments struct {
	dest      []byte
	value     *big.Int
	function  []byte
	args      [][]byte
	actualLen int32
}

func (context *VMHooksImpl) extractIndirectContractCallArgumentsWithValue(
	host vmhost.VMHost,
	destOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (*indirectContractCallArguments, error) {
	return context.extractIndirectContractCallArguments(
		host,
		destOffset,
		valueOffset,
		true,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
}

func (context *VMHooksImpl) extractIndirectContractCallArgumentsWithoutValue(
	host vmhost.VMHost,
	destOffset executor.MemPtr,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (*indirectContractCallArguments, error) {
	return context.extractIndirectContractCallArguments(
		host,
		destOffset,
		0,
		false,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
}

func (context *VMHooksImpl) extractIndirectContractCallArguments(
	host vmhost.VMHost,
	destOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	hasValueOffset bool,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (*indirectContractCallArguments, error) {
	dest, err := context.MemLoad(destOffset, vmhost.AddressLen)
	if err != nil {
		return nil, err
	}

	var value *big.Int

	if hasValueOffset {
		valueBytes, err := context.MemLoad(valueOffset, vmhost.BalanceLen)
		if err != nil {
			return nil, err
		}
		value = big.NewInt(0).SetBytes(valueBytes)
	}

	function, err := context.MemLoad(functionOffset, functionLength)
	if err != nil {
		return nil, err
	}

	args, actualLen, err := context.getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if err != nil {
		return nil, err
	}

	return &indirectContractCallArguments{
		dest:      dest,
		value:     value,
		function:  function,
		args:      args,
		actualLen: actualLen,
	}, nil
}

// TransferValueExecute VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) TransferValueExecute(
	destOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	gasLimit int64,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (int32, error) {
	host := context.GetVMHost()
	return context.TransferValueExecuteWithHost(
		host,
		destOffset,
		valueOffset,
		gasLimit,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
}

// TransferValueExecuteWithHost - transferValueExecute with host instead of pointer context
func (context *VMHooksImpl) TransferValueExecuteWithHost(
	host vmhost.VMHost,
	destOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	gasLimit int64,
	functionOffset executor.MemPtr,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (int32, error) {
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(transferValueExecuteName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.TransferValue
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return 1, err
	}

	callArgs, err := context.extractIndirectContractCallArgumentsWithValue(
		host, destOffset, valueOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if err != nil {
		return 1, err
	}

	return TransferValueExecuteWithTypedArgs(
		host,
		callArgs.dest,
		callArgs.value,
		gasLimit,
		callArgs.function,
		callArgs.args,
	)
}

// TransferValueExecuteWithTypedArgs - transferValueExecute with args already read from memory
func TransferValueExecuteWithTypedArgs(
	host vmhost.VMHost,
	dest []byte,
	value *big.Int,
	gasLimit int64,
	function []byte,
	args [][]byte,
) (int32, error) {
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.TransferValue
	err := metering.UseGasBounded(gasToUse)

	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return 1, err
	}

	sender := runtime.GetContextAddress()

	var contractCallInput *vmcommon.ContractCallInput
	if len(function) > 0 {
		contractCallInput, err = prepareIndirectContractCallInput(
			host,
			sender,
			value,
			gasLimit,
			dest,
			function,
			args,
			gasToUse,
			false,
		)
		if err != nil {
			return 1, err
		}
	}

	if contractCallInput != nil {
		if host.IsBuiltinFunctionName(contractCallInput.Function) {
			return 1, vmhost.ErrNilESDTData
		}
	}

	data := ""
	if contractCallInput != nil {
		data = makeCrossShardCallFromInput(contractCallInput.Function, contractCallInput.Arguments)
	}

	lastRound := host.Blockchain().LastRound()
	if host.IsBuiltinFunctionCall([]byte(data)) &&
		lastRound >= uint64(host.EnableEpochsHandler().GetActivationEpoch(vmhost.CheckBuiltInCallOnTransferValueAndFailExecutionFlag)) {
		return 1, vmhost.ErrTransferValueOnESDTCall
	}

	if host.AreInSameShard(sender, dest) && contractCallInput != nil && host.Blockchain().IsSmartContract(dest) {
		logEEI.Trace("eGLD pre-transfer execution begin")
		vmOutput, err := executeOnDestContextFromAPI(host, contractCallInput)
		if err != nil {
			logEEI.Trace("eGLD pre-transfer execution failed", "error", err)
			return 1, err
		}
		host.CompleteLogEntriesWithCallType(vmOutput, vmhost.TransferAndExecuteString)

		return 0, nil
	}

	err = metering.UseGasBounded(uint64(gasLimit))
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return 1, err
	}

	err = output.Transfer(dest, sender, uint64(gasLimit), 0, value, nil, []byte(data), vm.DirectCall)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

func makeCrossShardCallFromInput(function string, arguments [][]byte) string {
	var builder strings.Builder
	builder.WriteString(function)
	for _, arg := range arguments {
		builder.WriteByte('@')
		builder.WriteString(hex.EncodeToString(arg))
	}

	return builder.String()
}

// TransferESDTExecute VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) TransferESDTExecute(
	destOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	valueOffset executor.MemPtr,
	gasLimit int64,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (int32, error) {

	return context.TransferESDTNFTExecute(destOffset, tokenIDOffset, tokenIDLen, valueOffset, 0,
		gasLimit, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
}

// TransferESDTNFTExecute VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) TransferESDTNFTExecute(
	destOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	valueOffset executor.MemPtr,
	nonce int64,
	gasLimit int64,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (int32, error) {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(transferESDTNFTExecuteName)
	return context.TransferESDTNFTExecuteWithHost(
		host,
		destOffset,
		tokenIDOffset,
		tokenIDLen,
		valueOffset,
		nonce,
		gasLimit,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset)
}

// MultiTransferESDTNFTExecute VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) MultiTransferESDTNFTExecute(
	destOffset executor.MemPtr,
	numTokenTransfers int32,
	tokenTransfersArgsLengthOffset executor.MemPtr,
	tokenTransferDataOffset executor.MemPtr,
	gasLimit int64,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (int32, error) {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(multiTransferESDTNFTExecuteName)

	if numTokenTransfers == 0 {
		return 1, vmhost.ErrFailedTransfer
	}

	callArgs, err := context.extractIndirectContractCallArgumentsWithoutValue(
		host, destOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if err != nil {
		return 1, err
	}

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(callArgs.actualLen))
	err = metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return 1, err
	}

	numArgsFromMemory := math.MulInt32(numTokenTransfers, parsers.ArgsPerTransfer)
	transferArgs, _, err := context.getArgumentsFromMemory(
		host,
		numArgsFromMemory,
		tokenTransfersArgsLengthOffset,
		tokenTransferDataOffset,
	)
	if err != nil {
		return 1, err
	}

	transfers := make([]*vmcommon.ESDTTransfer, numTokenTransfers)
	for i := int32(0); i < numTokenTransfers; i++ {
		tokenStartIndex := i * parsers.ArgsPerTransfer
		transfer := &vmcommon.ESDTTransfer{
			ESDTTokenName:  transferArgs[tokenStartIndex],
			ESDTTokenNonce: big.NewInt(0).SetBytes(transferArgs[tokenStartIndex+1]).Uint64(),
			ESDTValue:      big.NewInt(0).SetBytes(transferArgs[tokenStartIndex+2]),
			ESDTTokenType:  uint32(core.Fungible),
		}
		if transfer.ESDTTokenNonce > 0 {
			transfer.ESDTTokenType = uint32(core.NonFungible)
		}
		transfers[i] = transfer
	}

	return TransferESDTNFTExecuteWithTypedArgs(
		host,
		callArgs.dest,
		transfers,
		gasLimit,
		callArgs.function,
		callArgs.args,
	)
}

// TransferESDTNFTExecuteWithHost contains only memory reading of arguments
func (context *VMHooksImpl) TransferESDTNFTExecuteWithHost(
	host vmhost.VMHost,
	destOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength,
	valueOffset executor.MemPtr,
	nonce int64,
	gasLimit int64,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (int32, error) {
	runtime := host.Runtime()
	metering := host.Metering()

	tokenIdentifier, executeErr := context.MemLoad(tokenIDOffset, tokenIDLen)
	if executeErr != nil {
		return 1, executeErr
	}

	callArgs, err := context.extractIndirectContractCallArgumentsWithValue(
		host, destOffset, valueOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if err != nil {
		return 1, err
	}

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(callArgs.actualLen))
	err = metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return 1, err
	}

	transfer := &vmcommon.ESDTTransfer{
		ESDTValue:      callArgs.value,
		ESDTTokenName:  tokenIdentifier,
		ESDTTokenNonce: uint64(nonce),
		ESDTTokenType:  uint32(core.Fungible),
	}
	if nonce > 0 {
		transfer.ESDTTokenType = uint32(core.NonFungible)
	}
	return TransferESDTNFTExecuteWithTypedArgs(
		host,
		callArgs.dest,
		[]*vmcommon.ESDTTransfer{transfer},
		gasLimit,
		callArgs.function,
		callArgs.args,
	)
}

// TransferESDTNFTExecuteWithTypedArgs defines the actual transfer ESDT execute logic
func TransferESDTNFTExecuteWithTypedArgs(
	host vmhost.VMHost,
	dest []byte,
	transfers []*vmcommon.ESDTTransfer,
	gasLimit int64,
	function []byte,
	data [][]byte,
) (int32, error) {
	return TransferESDTNFTExecuteWithTypedArgsWithFailure(host, dest, transfers, gasLimit, function, data, true)
}

// TransferESDTNFTExecuteWithTypedArgsWithFailure defines the actual transfer ESDT execute logic
func TransferESDTNFTExecuteWithTypedArgsWithFailure(
	host vmhost.VMHost,
	dest []byte,
	transfers []*vmcommon.ESDTTransfer,
	gasLimit int64,
	function []byte,
	data [][]byte,
	withFailure bool,
) (int32, error) {
	var executeErr error

	runtime := host.Runtime()
	metering := host.Metering()

	output := host.Output()

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOpsAPICost.TransferValue, uint64(len(transfers)))
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return 1, err
	}

	sender := runtime.GetContextAddress()

	var contractCallInput *vmcommon.ContractCallInput
	if len(function) > 0 {
		contractCallInput, executeErr = prepareIndirectContractCallInput(
			host,
			sender,
			big.NewInt(0),
			gasLimit,
			dest,
			function,
			data,
			gasToUse,
			false,
		)
		if executeErr != nil {
			return 1, executeErr
		}

		contractCallInput.ESDTTransfers = transfers
	}

	snapshotBeforeTransfer := host.Blockchain().GetSnapshot()
	output.PushState()

	originalCaller := host.Runtime().GetOriginalCallerAddress()
	transfersArgs := &vmhost.ESDTTransfersArgs{
		Destination:    dest,
		OriginalCaller: originalCaller,
		Sender:         sender,
		Transfers:      transfers,
		SenderForExec:  sender,
	}
	gasLimitForExec, executeErr := output.TransferESDT(transfersArgs, contractCallInput)
	if executeErr != nil {
		output.PopSetActiveState()
		errRevert := host.Blockchain().RevertToSnapshot(snapshotBeforeTransfer)
		if withFailure {
			return 1, executeErr
		}
		if errRevert != nil {
			return 1, errRevert
		}
		return 1, nil
	}

	if host.AreInSameShard(sender, dest) && contractCallInput != nil && host.Blockchain().IsSmartContract(dest) {
		contractCallInput.GasProvided = gasLimitForExec
		contractCallInput.CallerAddr = sender
		logEEI.Trace("ESDT post-transfer execution begin")
		vmOutput, executeErr := executeOnDestContextFromAPI(host, contractCallInput)
		if executeErr != nil {
			logEEI.Trace("ESDT post-transfer execution failed", "error", executeErr)
			output.PopSetActiveState()
			errRevert := host.Blockchain().RevertToSnapshot(snapshotBeforeTransfer)
			if vmOutput == nil || withFailure {
				return 1, executeErr
			}
			if errRevert != nil {
				return 1, errRevert
			}
			return 1, nil
		}

		output.PopDiscard()
		return 0, nil
	}

	output.PopDiscard()
	return 0, nil
}

// TransferESDTNFTExecuteByUserWithTypedArgs defines the actual transfer ESDT execute logic and execution
func TransferESDTNFTExecuteByUserWithTypedArgs(
	host vmhost.VMHost,
	callerForExecution []byte,
	dest []byte,
	transfers []*vmcommon.ESDTTransfer,
	gasLimit int64,
	function []byte,
	data [][]byte,
) (int32, error) {
	var executeErr error

	runtime := host.Runtime()
	metering := host.Metering()

	output := host.Output()

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOpsAPICost.TransferValue, uint64(len(transfers)))
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return 1, err
	}

	sender := runtime.GetContextAddress()

	var contractCallInput *vmcommon.ContractCallInput
	if len(function) > 0 {
		contractCallInput, executeErr = prepareIndirectContractCallInput(
			host,
			sender,
			big.NewInt(0),
			gasLimit,
			dest,
			function,
			data,
			gasToUse,
			false,
		)
		if executeErr != nil {
			return 1, executeErr
		}

		contractCallInput.ESDTTransfers = transfers
	}
	snapshotBeforeTransfer := host.Blockchain().GetSnapshot()
	output.PushState()

	originalCaller := host.Runtime().GetOriginalCallerAddress()
	transfersArgs := &vmhost.ESDTTransfersArgs{
		Destination:    dest,
		OriginalCaller: originalCaller,
		Sender:         sender,
		Transfers:      transfers,
		SenderForExec:  callerForExecution,
	}
	gasLimitForExec, executeErr := output.TransferESDT(transfersArgs, contractCallInput)
	if executeErr != nil {
		output.PopSetActiveState()
		errRevert := host.Blockchain().RevertToSnapshot(snapshotBeforeTransfer)
		// no fail execution is needed here - transfer was not successful, returning error which can be treated at SC level
		return 1, errRevert
	}

	if host.AreInSameShard(sender, dest) && contractCallInput != nil && host.Blockchain().IsSmartContract(dest) {
		contractCallInput.GasProvided = gasLimitForExec
		contractCallInput.CallerAddr = callerForExecution
		logEEI.Trace("ESDT post-transfer execution begin")
		_, executeErr = executeOnDestContextFromAPI(host, contractCallInput)
		if executeErr != nil {
			logEEI.Trace("ESDT post-transfer execution failed, started transfer to user", "error", executeErr)
			// in case of failed execution, the funds have to be moved to the user
			returnTransferArgs := &vmhost.ESDTTransfersArgs{
				Destination:      callerForExecution,
				OriginalCaller:   originalCaller,
				Sender:           dest,
				Transfers:        transfers,
				SenderForExec:    dest,
				ReturnAfterError: true,
			}
			_, executeErr = output.TransferESDT(returnTransferArgs, nil)
			if executeErr != nil {
				output.PopSetActiveState()
				_ = host.Blockchain().RevertToSnapshot(snapshotBeforeTransfer)
				// fail execution is needed here - tokens are at destination contract, so fail is needed to revert everything
				return 1, executeErr
			}

			output.PopDiscard()
			return 1, nil
		}

		output.PopDiscard()
		return 0, nil
	}

	output.PopDiscard()
	return 0, nil
}

// CreateAsyncCall VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) CreateAsyncCall(
	destOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	dataOffset executor.MemPtr,
	dataLength executor.MemLength,
	successOffset executor.MemPtr,
	successLength executor.MemLength,
	errorOffset executor.MemPtr,
	errorLength executor.MemLength,
	gas int64,
	extraGasForCallback int64,
) (int32, error) {
	host := context.GetVMHost()
	return context.CreateAsyncCallWithHost(
		host,
		destOffset,
		valueOffset,
		dataOffset,
		dataLength,
		successOffset,
		successLength,
		errorOffset,
		errorLength,
		gas,
		extraGasForCallback)
}

// CreateAsyncCallWithHost - createAsyncCall with host instead of pointer
func (context *VMHooksImpl) CreateAsyncCallWithHost(host vmhost.VMHost,
	destOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	dataOffset executor.MemPtr,
	dataLength executor.MemLength,
	successOffset executor.MemPtr,
	successLength executor.MemLength,
	errorOffset executor.MemPtr,
	errorLength executor.MemLength,
	gas int64,
	extraGasForCallback int64,
) (int32, error) {

	calledSCAddress, err := context.MemLoad(destOffset, vmhost.AddressLen)
	if err != nil {
		return 1, err
	}

	value, err := context.MemLoad(valueOffset, vmhost.BalanceLen)
	if err != nil {
		return 1, err
	}

	data, err := context.MemLoad(dataOffset, dataLength)
	if err != nil {
		return 1, err
	}

	successFunc, err := context.MemLoad(successOffset, successLength)
	if err != nil {
		return 1, err
	}

	errorFunc, err := context.MemLoad(errorOffset, errorLength)
	if err != nil {
		return 1, err
	}

	return CreateAsyncCallWithTypedArgs(host,
		calledSCAddress,
		value,
		data,
		successFunc,
		errorFunc,
		gas,
		extraGasForCallback,
		nil)
}

// CreateAsyncCallWithTypedArgs - createAsyncCall with arguments already read from memory
func CreateAsyncCallWithTypedArgs(host vmhost.VMHost,
	calledSCAddress []byte,
	value []byte,
	data []byte,
	successFunc []byte,
	errorFunc []byte,
	gas int64,
	extraGasForCallback int64,
	callbackClosure []byte) (int32, error) {

	metering := host.Metering()
	runtime := host.Runtime()
	async := host.Async()

	metering.StartGasTracing(createAsyncCallName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateAsyncCall
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return 1, err
	}

	asyncCall := &vmhost.AsyncCall{
		Status:          vmhost.AsyncCallPending,
		Destination:     calledSCAddress,
		Data:            data,
		ValueBytes:      value,
		GasLimit:        uint64(gas),
		SuccessCallback: string(successFunc),
		ErrorCallback:   string(errorFunc),
		GasLocked:       uint64(extraGasForCallback),
		CallbackClosure: callbackClosure,
	}

	if asyncCall.HasDefinedAnyCallback() {
		gasToUse = metering.GasSchedule().BaseOpsAPICost.SetAsyncCallback
		err = metering.UseGasBounded(gasToUse)
		if err != nil && runtime.UseGasBoundedShouldFailExecution() {
			return 1, err
		}
	}

	err = async.RegisterAsyncCall("", asyncCall)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

// SetAsyncContextCallback VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SetAsyncContextCallback(
	callback executor.MemPtr,
	callbackLength executor.MemLength,
	data executor.MemPtr,
	dataLength executor.MemLength,
	gas int64,
) (int32, error) {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	async := host.Async()
	metering.StartGasTracing(setAsyncContextCallbackName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.SetAsyncContextCallback
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return 1, err
	}

	callbackNameBytes, err := context.MemLoad(callback, callbackLength)
	if err != nil {
		return 1, err
	}

	dataBytes, err := context.MemLoad(data, dataLength)
	if err != nil {
		return 1, err
	}

	err = async.SetContextCallback(
		string(callbackNameBytes),
		dataBytes,
		uint64(gas))
	if err != nil {
		return 1, err
	}

	return 0, nil
}

// UpgradeContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) UpgradeContract(
	destOffset executor.MemPtr,
	gasLimit int64,
	valueOffset executor.MemPtr,
	codeOffset executor.MemPtr,
	codeMetadataOffset executor.MemPtr,
	length executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) error {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(upgradeContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	value, err := context.MemLoad(valueOffset, vmhost.BalanceLen)
	if err != nil {
		return err
	}

	code, err := context.MemLoad(codeOffset, length)
	if err != nil {
		return err
	}

	codeMetadata, err := context.MemLoad(codeMetadataOffset, vmhost.CodeMetadataLen)
	if err != nil {
		return err
	}

	data, _, err := context.getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	calledSCAddress, err := context.MemLoad(destOffset, vmhost.AddressLen)
	if err != nil {
		return err
	}

	gasSchedule := metering.GasSchedule()
	gasToUse = math.MulUint64(gasSchedule.BaseOperationCost.DataCopyPerByte, uint64(length))
	err = metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	return upgradeContract(host, calledSCAddress, code, codeMetadata, value, data, gasLimit)
}

// UpgradeFromSourceContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) UpgradeFromSourceContract(
	destOffset executor.MemPtr,
	gasLimit int64,
	valueOffset executor.MemPtr,
	sourceContractAddressOffset executor.MemPtr,
	codeMetadataOffset executor.MemPtr,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) error {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(upgradeFromSourceContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	value, err := context.MemLoad(valueOffset, vmhost.BalanceLen)
	if err != nil {
		return err
	}

	sourceContractAddress, err := context.MemLoad(sourceContractAddressOffset, vmhost.AddressLen)
	if err != nil {
		return err
	}

	codeMetadata, err := context.MemLoad(codeMetadataOffset, vmhost.CodeMetadataLen)
	if err != nil {
		return err
	}

	data, _, err := context.getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	calledSCAddress, err := context.MemLoad(destOffset, vmhost.AddressLen)
	if err != nil {
		return err
	}

	return UpgradeFromSourceContractWithTypedArgs(
		host,
		sourceContractAddress,
		calledSCAddress,
		value,
		data,
		gasLimit,
		codeMetadata,
	)
}

// UpgradeFromSourceContractWithTypedArgs - upgradeFromSourceContract with args already read from memory
func UpgradeFromSourceContractWithTypedArgs(
	host vmhost.VMHost,
	sourceContractAddress []byte,
	destContractAddress []byte,
	value []byte,
	data [][]byte,
	gasLimit int64,
	codeMetadata []byte,
) error {
	blockchain := host.Blockchain()

	code, err := blockchain.GetCode(sourceContractAddress)
	if err != nil {
		return err
	}

	return upgradeContract(host, destContractAddress, code, codeMetadata, value, data, gasLimit)
}

func upgradeContract(
	host vmhost.VMHost,
	destContractAddress []byte,
	code []byte,
	codeMetadata []byte,
	value []byte,
	data [][]byte,
	gasLimit int64,
) error {
	metering := host.Metering()
	gasSchedule := metering.GasSchedule()
	minAsyncCallCost := math.AddUint64(
		math.MulUint64(2, gasSchedule.BaseOpsAPICost.AsyncCallStep),
		gasSchedule.BaseOpsAPICost.AsyncCallbackGasLock)
	if gasLimit < 0 || uint64(gasLimit) < minAsyncCallCost {
		return vmhost.ExitOutOfGas()
	}

	// Set up the async call as if it is not known whether the called SC
	// is in the same shard with the caller or not. This will be later resolved
	// by runtime.ExecuteAsyncCall().
	args := make([][]byte, 2+len(data))
	args[0] = code
	args[1] = codeMetadata
	copy(args[2:], data)
	callData := makeCrossShardCallFromInput(vmhost.UpgradeFunctionName, args)

	async := host.Async()
	err := async.RegisterLegacyAsyncCall(
		destContractAddress,
		[]byte(callData),
		value,
	)
	logEEI.Trace("upgradeContract", "error", err)

	if errors.Is(err, vmhost.ErrNotEnoughGas) {
		return vmhost.ExitOutOfGas()
	}
	if err != nil {
		return err
	}
	return nil
}

// DeleteContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) DeleteContract(
	destOffset executor.MemPtr,
	gasLimit int64,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) error {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(deleteContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	data, _, err := context.getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	calledSCAddress, err := context.MemLoad(destOffset, vmhost.AddressLen)
	if err != nil {
		return err
	}

	return deleteContract(
		host,
		calledSCAddress,
		data,
		gasLimit,
	)
}

func deleteContract(
	host vmhost.VMHost,
	dest []byte,
	data [][]byte,
	gasLimit int64,
) error {
	metering := host.Metering()
	gasSchedule := metering.GasSchedule()
	minAsyncCallCost := math.AddUint64(
		math.MulUint64(2, gasSchedule.BaseOpsAPICost.AsyncCallStep),
		gasSchedule.BaseOpsAPICost.AsyncCallbackGasLock)
	if gasLimit < 0 || uint64(gasLimit) < minAsyncCallCost {
		return vmhost.ExitOutOfGas()
	}

	callData := makeCrossShardCallFromInput(vmhost.DeleteFunctionName, data)

	async := host.Async()
	err := async.RegisterLegacyAsyncCall(
		dest,
		[]byte(callData),
		big.NewInt(0).Bytes(),
	)
	logEEI.Trace("deleteContract", "error", err)

	if errors.Is(err, vmhost.ErrNotEnoughGas) {
		return vmhost.ExitOutOfGas()
	}
	if err != nil {
		return err
	}
	return nil
}

// AsyncCall VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) AsyncCall(
	destOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	dataOffset executor.MemPtr,
	length executor.MemLength) error {

	host := context.GetVMHost()
	runtime := host.Runtime()
	async := host.Async()
	metering := host.Metering()
	metering.StartGasTracing(asyncCallName)

	gasSchedule := metering.GasSchedule()
	gasToUse := gasSchedule.BaseOpsAPICost.AsyncCallStep
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	calledSCAddress, err := context.MemLoad(destOffset, vmhost.AddressLen)
	if err != nil {
		return err
	}

	value, err := context.MemLoad(valueOffset, vmhost.BalanceLen)
	if err != nil {
		return err
	}

	gasToUse = math.MulUint64(gasSchedule.BaseOperationCost.DataCopyPerByte, uint64(length))
	err = metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return err
	}

	data, err := context.MemLoad(dataOffset, length)
	if err != nil {
		return err
	}

	err = async.RegisterLegacyAsyncCall(calledSCAddress, data, value)
	if errors.Is(err, vmhost.ErrNotEnoughGas) {
		return vmhost.ExitOutOfGas()
	}
	if err != nil {
		return err
	}
	return nil
}

// GetArgumentLength VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetArgumentLength(id int32) (int32, error) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetArgument
	err := metering.UseGasBoundedAndAddTracedGas(getArgumentLengthName, gasToUse)
	if err != nil {
		return -1, err
	}

	args := runtime.Arguments()
	if id < 0 || int32(len(args)) <= id {
		return -1, vmhost.ErrInvalidArgument
	}

	return int32(len(args[id])), nil
}

// GetArgument VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetArgument(id int32, argOffset executor.MemPtr) (int32, error) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetArgument
	err := metering.UseGasBoundedAndAddTracedGas(getArgumentName, gasToUse)
	if err != nil {
		return -1, err
	}

	args := runtime.Arguments()
	if id < 0 || int32(len(args)) <= id {
		return -1, vmhost.ErrInvalidArgument
	}

	err = context.MemStore(argOffset, args[id])
	if err != nil {
		return -1, err
	}

	return int32(len(args[id])), nil
}

// GetFunction VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetFunction(functionOffset executor.MemPtr) (int32, error) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetFunction
	err := metering.UseGasBoundedAndAddTracedGas(getFunctionName, gasToUse)
	if err != nil {
		return -1, err
	}

	function := runtime.FunctionName()
	err = context.MemStore(functionOffset, []byte(function))
	if err != nil {
		return -1, err
	}

	return int32(len(function)), nil
}

// GetNumArguments VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetNumArguments() (int32, error) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetNumArguments
	err := metering.UseGasBoundedAndAddTracedGas(getNumArgumentsName, gasToUse)
	if err != nil {
		return -1, err
	}

	args := runtime.Arguments()
	return int32(len(args)), nil
}

// StorageStore VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) StorageStore(
	keyOffset executor.MemPtr,
	keyLength executor.MemLength,
	dataOffset executor.MemPtr,
	dataLength executor.MemLength) (int32, error) {

	host := context.GetVMHost()
	return context.StorageStoreWithHost(
		host,
		keyOffset,
		keyLength,
		dataOffset,
		dataLength,
	)
}

// StorageStoreWithHost - storageStore with host instead of pointer context
func (context *VMHooksImpl) StorageStoreWithHost(
	host vmhost.VMHost,
	keyOffset executor.MemPtr,
	keyLength executor.MemLength,
	dataOffset executor.MemPtr,
	dataLength executor.MemLength) (int32, error) {

	key, err := context.MemLoad(keyOffset, keyLength)
	if err != nil {
		return -1, err
	}

	data, err := context.MemLoad(dataOffset, dataLength)
	if err != nil {
		return -1, err
	}

	return StorageStoreWithTypedArgs(host, key, data)
}

// StorageStoreWithTypedArgs - storageStore with args already read from memory
func StorageStoreWithTypedArgs(host vmhost.VMHost, key []byte, data []byte) (int32, error) {
	storage := host.Storage()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.StorageStore
	err := metering.UseGasBoundedAndAddTracedGas(storageStoreName, gasToUse)
	if err != nil {
		return -1, err
	}

	storageStatus, err := storage.SetStorage(key, data)
	if err != nil {
		return -1, err
	}

	return int32(storageStatus), nil
}

// StorageLoadLength VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) StorageLoadLength(keyOffset executor.MemPtr, keyLength executor.MemLength) (int32, error) {
	storage := context.GetStorageContext()
	metering := context.GetMeteringContext()

	key, err := context.MemLoad(keyOffset, keyLength)
	if err != nil {
		return -1, err
	}

	data, trieDepth, usedCache, err := storage.GetStorageUnmetered(key)
	if err != nil {
		return -1, err
	}

	err = storage.UseGasForStorageLoad(
		storageLoadLengthName,
		int64(trieDepth),
		metering.GasSchedule().BaseOpsAPICost.StorageLoad,
		usedCache)
	if err != nil {
		return -1, err
	}

	return int32(len(data)), nil
}

// StorageLoadFromAddress VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) StorageLoadFromAddress(
	addressOffset executor.MemPtr,
	keyOffset executor.MemPtr,
	keyLength executor.MemLength,
	dataOffset executor.MemPtr) (int32, error) {

	host := context.GetVMHost()
	return context.StorageLoadFromAddressWithHost(
		host,
		addressOffset,
		keyOffset,
		keyLength,
		dataOffset,
	)
}

// StorageLoadFromAddressWithHost - storageLoadFromAddress with host instead of pointer context
func (context *VMHooksImpl) StorageLoadFromAddressWithHost(
	host vmhost.VMHost,
	addressOffset executor.MemPtr,
	keyOffset executor.MemPtr,
	keyLength executor.MemLength,
	dataOffset executor.MemPtr) (int32, error) {

	key, err := context.MemLoad(keyOffset, keyLength)
	if err != nil {
		return -1, err
	}

	address, err := context.MemLoad(addressOffset, vmhost.AddressLen)
	if err != nil {
		return -1, err
	}

	data, err := StorageLoadFromAddressWithTypedArgs(host, address, key)
	if err != nil {
		return -1, err
	}

	err = context.MemStore(dataOffset, data)
	if err != nil {
		return -1, err
	}

	return int32(len(data)), nil
}

// StorageLoadFromAddressWithTypedArgs - storageLoadFromAddress with args already read from memory
func StorageLoadFromAddressWithTypedArgs(host vmhost.VMHost, address []byte, key []byte) ([]byte, error) {
	storage := host.Storage()
	metering := host.Metering()
	data, trieDepth, usedCache, err := storage.GetStorageFromAddress(address, key)
	if err != nil {
		return nil, err
	}
	err = storage.UseGasForStorageLoad(
		storageLoadFromAddressName,
		int64(trieDepth),
		metering.GasSchedule().BaseOpsAPICost.StorageLoad,
		usedCache)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// StorageLoad VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) StorageLoad(keyOffset executor.MemPtr, keyLength executor.MemLength, dataOffset executor.MemPtr) (int32, error) {
	host := context.GetVMHost()
	return context.StorageLoadWithHost(
		host,
		keyOffset,
		keyLength,
		dataOffset,
	)
}

// StorageLoadWithHost - storageLoad with host instead of pointer context
func (context *VMHooksImpl) StorageLoadWithHost(host vmhost.VMHost, keyOffset executor.MemPtr, keyLength executor.MemLength, dataOffset executor.MemPtr) (int32, error) {

	key, err := context.MemLoad(keyOffset, keyLength)
	if err != nil {
		return -1, err
	}

	data, err := StorageLoadWithWithTypedArgs(host, key)
	if err != nil {
		return -1, err
	}

	err = context.MemStore(dataOffset, data)
	if err != nil {
		return -1, err
	}

	return int32(len(data)), nil
}

// StorageLoadWithWithTypedArgs - storageLoad with args already read from memory
func StorageLoadWithWithTypedArgs(host vmhost.VMHost, key []byte) ([]byte, error) {
	storage := host.Storage()
	metering := host.Metering()
	data, trieDepth, usedCache, err := storage.GetStorage(key)
	if err != nil {
		return nil, err
	}

	err = storage.UseGasForStorageLoad(
		storageLoadName,
		int64(trieDepth),
		metering.GasSchedule().BaseOpsAPICost.StorageLoad,
		usedCache)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// SetStorageLock VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) SetStorageLock(keyOffset executor.MemPtr, keyLength executor.MemLength, lockTimestamp int64) (int32, error) {
	host := context.GetVMHost()
	return context.SetStorageLockWithHost(
		host,
		keyOffset,
		keyLength,
		lockTimestamp,
	)
}

// SetStorageLockWithHost - setStorageLock with host instead of pointer context
func (context *VMHooksImpl) SetStorageLockWithHost(host vmhost.VMHost, keyOffset executor.MemPtr, keyLength executor.MemLength, lockTimestamp int64) (int32, error) {
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Int64StorageStore
	err := metering.UseGasBoundedAndAddTracedGas(setStorageLockName, gasToUse)
	if err != nil {
		return -1, err
	}

	key, err := context.MemLoad(keyOffset, keyLength)
	if err != nil {
		return -1, err
	}

	return SetStorageLockWithTypedArgs(host, key, lockTimestamp)
}

// SetStorageLockWithTypedArgs - setStorageLock with args already read from memory
func SetStorageLockWithTypedArgs(host vmhost.VMHost, key []byte, lockTimestamp int64) (int32, error) {
	storage := host.Storage()
	timeLockKeyPrefix := string(storage.GetVmProtectedPrefix(vmhost.TimeLockKeyPrefix))
	timeLockKey := vmhost.CustomStorageKey(timeLockKeyPrefix, key)
	bigTimestamp := big.NewInt(0).SetInt64(lockTimestamp)
	storageStatus, err := storage.SetProtectedStorage(timeLockKey, bigTimestamp.Bytes())
	if err != nil {
		return -1, err
	}
	return int32(storageStatus), nil
}

// GetStorageLock VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetStorageLock(keyOffset executor.MemPtr, keyLength executor.MemLength) (int64, error) {
	metering := context.GetMeteringContext()
	storage := context.GetStorageContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.StorageLoad
	err := metering.UseGasBoundedAndAddTracedGas(getStorageLockName, gasToUse)
	if err != nil {
		return -1, err
	}

	key, err := context.MemLoad(keyOffset, keyLength)
	if err != nil {
		return -1, err
	}

	timeLockKeyPrefix := string(storage.GetVmProtectedPrefix(vmhost.TimeLockKeyPrefix))
	timeLockKey := vmhost.CustomStorageKey(timeLockKeyPrefix, key)

	data, trieDepth, usedCache, err := storage.GetStorage(timeLockKey)
	if err != nil {
		return -1, err
	}

	err = storage.UseGasForStorageLoad(
		getStorageLockName,
		int64(trieDepth),
		metering.GasSchedule().BaseOpsAPICost.StorageLoad,
		usedCache)
	if err != nil {
		return -1, err
	}

	timeLock := big.NewInt(0).SetBytes(data).Int64()

	// TODO if timelock <= currentTimeStamp { fail somehow }

	return timeLock, nil
}

// IsStorageLocked VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) IsStorageLocked(keyOffset executor.MemPtr, keyLength executor.MemLength) (int32, error) {
	timeLock, err := context.GetStorageLock(keyOffset, keyLength)
	if err != nil {
		return -1, err
	}
	if timeLock < 0 {
		return -1, nil
	}

	currentTimestamp, err := context.GetBlockTimestamp()
	if err != nil {
		return -1, err
	}
	if timeLock <= currentTimestamp {
		return 0, nil
	}

	return 1, nil
}

// ClearStorageLock VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ClearStorageLock(keyOffset executor.MemPtr, keyLength executor.MemLength) (int32, error) {
	return context.SetStorageLock(keyOffset, keyLength, 0)
}

// GetCaller VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetCaller(resultOffset executor.MemPtr) error {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCaller
	err := metering.UseGasBoundedAndAddTracedGas(getCallerName, gasToUse)
	if err != nil {
		return err
	}

	caller := runtime.GetVMInput().CallerAddr

	err = context.MemStore(resultOffset, caller)
	if err != nil {
		return err
	}
	return nil
}

// CheckNoPayment VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) CheckNoPayment() error {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(checkNoPaymentName, gasToUse)
	if err != nil {
		return err
	}

	vmInput := runtime.GetVMInput()
	if vmInput.CallValue.Sign() > 0 {
		return vmhost.ErrNonPayableFunctionEgld
	}
	if len(vmInput.ESDTTransfers) > 0 {
		return vmhost.ErrNonPayableFunctionEsdt
	}
	return nil
}

// GetCallValue VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetCallValue(resultOffset executor.MemPtr) (int32, error) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(callValueName, gasToUse)
	if err != nil {
		return -1, err
	}

	value := runtime.GetVMInput().CallValue.Bytes()
	value = vmhost.PadBytesLeft(value, vmhost.BalanceLen)

	err = context.MemStore(resultOffset, value)
	if err != nil {
		return -1, err
	}

	return int32(len(value)), nil
}

// GetESDTValue VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetESDTValue(resultOffset executor.MemPtr) (int32, error) {
	err := failIfMoreThanOneESDTTransfer(context)
	if err != nil {
		return -1, err
	}
	return context.GetESDTValueByIndex(resultOffset, 0)
}

// GetESDTValueByIndex VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetESDTValueByIndex(resultOffset executor.MemPtr, index int32) (int32, error) {
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(getESDTValueByIndexName, gasToUse)
	if err != nil {
		return -1, err
	}

	var value []byte

	esdtTransfer, err := getESDTTransferFromInputFailIfWrongIndex(context.GetVMHost(), index)
	if err != nil {
		return 0, err
	}
	if esdtTransfer != nil && esdtTransfer.ESDTValue.Cmp(vmhost.Zero) > 0 {
		value = esdtTransfer.ESDTValue.Bytes()
		value = vmhost.PadBytesLeft(value, vmhost.BalanceLen)
	}

	err = context.MemStore(resultOffset, value)
	if err != nil {
		return -1, err
	}

	return int32(len(value)), nil
}

// GetESDTTokenName VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetESDTTokenName(resultOffset executor.MemPtr) (int32, error) {
	err := failIfMoreThanOneESDTTransfer(context)
	if err != nil {
		return -1, err
	}
	return context.GetESDTTokenNameByIndex(resultOffset, 0)
}

// GetESDTTokenNameByIndex VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetESDTTokenNameByIndex(resultOffset executor.MemPtr, index int32) (int32, error) {
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(getESDTTokenNameByIndexName, gasToUse)
	if err != nil {
		return -1, err
	}

	esdtTransfer, err := getESDTTransferFromInputFailIfWrongIndex(context.GetVMHost(), index)
	if err != nil {
		return 0, err
	}
	var tokenName []byte
	if esdtTransfer != nil {
		tokenName = esdtTransfer.ESDTTokenName
	}

	err = context.MemStore(resultOffset, tokenName)
	if err != nil {
		return -1, err
	}

	return int32(len(tokenName)), nil
}

// GetESDTTokenNonce VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetESDTTokenNonce() (int64, error) {
	err := failIfMoreThanOneESDTTransfer(context)
	if err != nil {
		return -1, err
	}
	return context.GetESDTTokenNonceByIndex(0)
}

// GetESDTTokenNonceByIndex VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetESDTTokenNonceByIndex(index int32) (int64, error) {
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(getESDTTokenNonceByIndexName, gasToUse)
	if err != nil {
		return -1, err
	}

	esdtTransfer, err := getESDTTransferFromInputFailIfWrongIndex(context.GetVMHost(), index)
	if err != nil {
		return 0, err
	}
	nonce := uint64(0)
	if esdtTransfer != nil {
		nonce = esdtTransfer.ESDTTokenNonce
	}
	return int64(nonce), nil
}

// GetCurrentESDTNFTNonce VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetCurrentESDTNFTNonce(
	addressOffset executor.MemPtr,
	tokenIDOffset executor.MemPtr,
	tokenIDLen executor.MemLength) (int64, error) {

	metering := context.GetMeteringContext()
	storage := context.GetStorageContext()

	destination, err := context.MemLoad(addressOffset, vmhost.AddressLen)
	if err != nil {
		return 0, err
	}

	tokenID, err := context.MemLoad(tokenIDOffset, tokenIDLen)
	if err != nil {
		return 0, err
	}

	key := []byte(core.ProtectedKeyPrefix + core.ESDTNFTLatestNonceIdentifier + string(tokenID))
	data, trieDepth, _, err := storage.GetStorageFromAddress(destination, key)
	if err != nil {
		return 0, err
	}

	err = storage.UseGasForStorageLoad(
		getCurrentESDTNFTNonceName,
		int64(trieDepth),
		metering.GasSchedule().BaseOpsAPICost.StorageLoad,
		false)
	if err != nil {
		return 0, err
	}

	nonce := big.NewInt(0).SetBytes(data).Uint64()
	return int64(nonce), nil
}

// GetESDTTokenType VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetESDTTokenType() (int32, error) {
	err := failIfMoreThanOneESDTTransfer(context)
	if err != nil {
		return -1, err
	}
	return context.GetESDTTokenTypeByIndex(0)
}

// GetESDTTokenTypeByIndex VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetESDTTokenTypeByIndex(index int32) (int32, error) {
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(getESDTTokenTypeByIndexName, gasToUse)
	if err != nil {
		return -1, err
	}

	esdtTransfer, err := getESDTTransferFromInputFailIfWrongIndex(context.GetVMHost(), index)
	if err != nil {
		return 0, err
	}
	if esdtTransfer != nil {
		return int32(esdtTransfer.ESDTTokenType), nil
	}
	return 0, nil
}

// GetNumESDTTransfers VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetNumESDTTransfers() (int32, error) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(getNumESDTTransfersName, gasToUse)
	if err != nil {
		return -1, err
	}

	return int32(len(runtime.GetVMInput().ESDTTransfers)), nil
}

// GetCallValueTokenName VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetCallValueTokenName(callValueOffset executor.MemPtr, tokenNameOffset executor.MemPtr) (int32, error) {
	err := failIfMoreThanOneESDTTransfer(context)
	if err != nil {
		return -1, err
	}
	return context.GetCallValueTokenNameByIndex(callValueOffset, tokenNameOffset, 0)
}

// GetCallValueTokenNameByIndex VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetCallValueTokenNameByIndex(
	callValueOffset executor.MemPtr,
	tokenNameOffset executor.MemPtr,
	index int32) (int32, error) {

	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(getCallValueTokenNameByIndexName, gasToUse)
	if err != nil {
		return -1, err
	}

	callValue := runtime.GetVMInput().CallValue.Bytes()
	tokenName := make([]byte, 0)
	esdtTransfer, err := getESDTTransferFromInputFailIfWrongIndex(context.GetVMHost(), index)
	if err != nil {
		return 0, err
	}

	if esdtTransfer != nil {
		tokenName = make([]byte, len(esdtTransfer.ESDTTokenName))
		copy(tokenName, esdtTransfer.ESDTTokenName)
		callValue = esdtTransfer.ESDTValue.Bytes()
	}
	callValue = vmhost.PadBytesLeft(callValue, vmhost.BalanceLen)

	err = context.MemStore(tokenNameOffset, tokenName)
	if err != nil {
		return -1, err
	}

	err = context.MemStore(callValueOffset, callValue)
	if err != nil {
		return -1, err
	}

	return int32(len(tokenName)), nil
}

// IsReservedFunctionName VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) IsReservedFunctionName(nameHandle int32) (int32, error) {
	host := context.host
	managedTypes := context.GetManagedTypesContext()
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.IsReservedFunctionName
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return -1, err
	}

	name, err := managedTypes.GetBytes(nameHandle)
	if err != nil {
		return -1, err
	}

	err = chargeGasForExtraLengthUnderFlag(len(name), host)
	if err != nil {
		return -1, err
	}

	if runtime.IsReservedFunctionName(string(name)) {
		return 1, nil
	}

	return 0, nil
}

// WriteLog VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) WriteLog(
	dataPointer executor.MemPtr,
	dataLength executor.MemLength,
	topicPtr executor.MemPtr,
	numTopics int32) error {

	// note: deprecated
	runtime := context.GetRuntimeContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Log
	fullDataForGas := math.MulUint64(uint64(numTopics), vmhost.HashLen)
	fullDataForGas = math.AddUint64(fullDataForGas, uint64(dataLength))
	gas := math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, fullDataForGas)
	gasToUse = math.AddUint64(gasToUse, gas)

	if numTopics < 0 || dataLength < 0 {
		return vmhost.ErrNegativeLength
	}

	err := metering.UseGasBoundedAndAddTracedGas(writeLogName, gasToUse)
	if err != nil {
		return err
	}

	log, err := context.MemLoad(dataPointer, dataLength)
	if err != nil {
		return err
	}

	topics := make([][]byte, numTopics)
	for i := int32(0); i < numTopics; i++ {
		topics[i], err = context.MemLoad(topicPtr.Offset(i*vmhost.HashLen), vmhost.HashLen)
		if err != nil {
			return err
		}
	}

	output.WriteLog(runtime.GetContextAddress(), topics, [][]byte{log})
	return nil
}

// WriteEventLog VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) WriteEventLog(
	numTopics int32,
	topicLengthsOffset executor.MemPtr,
	topicOffset executor.MemPtr,
	dataOffset executor.MemPtr,
	dataLength executor.MemLength,
) error {

	host := context.GetVMHost()
	runtime := context.GetRuntimeContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(writeEventLogName)

	topics, _, err := context.getArgumentsFromMemory(
		host,
		numTopics,
		topicLengthsOffset,
		topicOffset,
	)
	if err != nil {
		return err
	}

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Log
	gasForData := math.MulUint64(
		metering.GasSchedule().BaseOperationCost.DataCopyPerByte,
		uint64(dataLength))
	gasToUse = math.AddUint64(gasToUse, gasForData)
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	data, err := context.MemLoad(dataOffset, dataLength)
	if err != nil {
		return err
	}

	output.WriteLog(runtime.GetContextAddress(), topics, [][]byte{data})
	return nil
}

// GetBlockTimestamp VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetBlockTimestamp() (int64, error) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockTimeStamp
	err := metering.UseGasBoundedAndAddTracedGas(getBlockTimestampName, gasToUse)
	if err != nil {
		return -1, err
	}

	return int64(blockchain.CurrentTimeStamp()), nil
}

// GetBlockTimestampMs VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetBlockTimestampMs() (int64, error) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockTimeStamp
	err := metering.UseGasBoundedAndAddTracedGas(getBlockTimestampMsName, gasToUse)
	if err != nil {
		return -1, err
	}

	return int64(blockchain.CurrentTimeStampMs()), nil
}

// GetBlockNonce VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetBlockNonce() (int64, error) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockNonce
	err := metering.UseGasBoundedAndAddTracedGas(getBlockNonceName, gasToUse)
	if err != nil {
		return -1, err
	}

	return int64(blockchain.CurrentNonce()), nil
}

// GetBlockRound VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetBlockRound() (int64, error) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRound
	err := metering.UseGasBoundedAndAddTracedGas(getBlockRoundName, gasToUse)
	if err != nil {
		return -1, err
	}

	return int64(blockchain.CurrentRound()), nil
}

// GetBlockEpoch VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetBlockEpoch() (int64, error) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockEpoch
	err := metering.UseGasBoundedAndAddTracedGas(getBlockEpochName, gasToUse)
	if err != nil {
		return -1, err
	}

	return int64(blockchain.CurrentEpoch()), nil
}

// GetBlockRandomSeed VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetBlockRandomSeed(pointer executor.MemPtr) error {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRandomSeed
	err := metering.UseGasBoundedAndAddTracedGas(getBlockRandomSeedName, gasToUse)
	if err != nil {
		return err
	}

	randomSeed := blockchain.CurrentRandomSeed()
	err = context.MemStore(pointer, randomSeed)
	if err != nil {
		return err
	}
	return nil
}

// GetStateRootHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetStateRootHash(pointer executor.MemPtr) error {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetStateRootHash
	err := metering.UseGasBoundedAndAddTracedGas(getStateRootHashName, gasToUse)
	if err != nil {
		return err
	}

	stateRootHash := blockchain.GetStateRootHash()
	err = context.MemStore(pointer, stateRootHash)
	if err != nil {
		return err
	}
	return nil
}

// GetPrevBlockTimestamp VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetPrevBlockTimestamp() (int64, error) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockTimeStamp
	err := metering.UseGasBoundedAndAddTracedGas(getPrevBlockTimestampName, gasToUse)
	if err != nil {
		return -1, err
	}

	return int64(blockchain.LastTimeStamp()), nil
}

// GetPrevBlockTimestampMs VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetPrevBlockTimestampMs() (int64, error) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockTimeStamp
	err := metering.UseGasBoundedAndAddTracedGas(getPrevBlockTimestampMsName, gasToUse)
	if err != nil {
		return -1, err
	}

	return int64(blockchain.LastTimeStampMs()), nil
}

// GetPrevBlockNonce VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetPrevBlockNonce() (int64, error) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockNonce
	err := metering.UseGasBoundedAndAddTracedGas(getPrevBlockNonceName, gasToUse)
	if err != nil {
		return -1, err
	}

	return int64(blockchain.LastNonce()), nil
}

// GetPrevBlockRound VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetPrevBlockRound() (int64, error) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRound
	err := metering.UseGasBoundedAndAddTracedGas(getPrevBlockRoundName, gasToUse)
	if err != nil {
		return -1, err
	}

	return int64(blockchain.LastRound()), nil
}

// GetPrevBlockEpoch VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetPrevBlockEpoch() (int64, error) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockEpoch
	err := metering.UseGasBoundedAndAddTracedGas(getPrevBlockEpochName, gasToUse)
	if err != nil {
		return -1, err
	}

	return int64(blockchain.LastEpoch()), nil
}

// GetPrevBlockRandomSeed VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetPrevBlockRandomSeed(pointer executor.MemPtr) error {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRandomSeed
	err := metering.UseGasBoundedAndAddTracedGas(getPrevBlockRandomSeedName, gasToUse)
	if err != nil {
		return err
	}

	randomSeed := blockchain.LastRandomSeed()
	err = context.MemStore(pointer, randomSeed)
	if err != nil {
		return err
	}
	return nil
}

// GetBlockRoundTimeMs VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetBlockRoundTimeMs() (int64, error) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetRoundTime
	err := metering.UseGasBoundedAndAddTracedGas(getBlockRoundTimeMsName, gasToUse)
	if err != nil {
		return -1, err
	}

	return int64(blockchain.RoundTime()), nil
}

// EpochStartBlockTimestampMs VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) EpochStartBlockTimestampMs() (int64, error) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.EpochStartBlockTimeStamp
	err := metering.UseGasBoundedAndAddTracedGas(epochStartBlockTimestampMsName, gasToUse)
	if err != nil {
		return -1, err
	}

	return int64(blockchain.EpochStartBlockTimeStampMs()), nil
}

// EpochStartBlockNonce VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) EpochStartBlockNonce() (int64, error) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.EpochStartBlockNonce
	err := metering.UseGasBoundedAndAddTracedGas(epochStartBlockNonceName, gasToUse)
	if err != nil {
		return -1, err
	}

	return int64(blockchain.EpochStartBlockNonce()), nil
}

// EpochStartBlockRound VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) EpochStartBlockRound() (int64, error) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.EpochStartBlockRound

	err := metering.UseGasBoundedAndAddTracedGas(epochStartBlockRoundName, gasToUse)
	if err != nil {
		return -1, err
	}

	return int64(blockchain.EpochStartBlockRound()), nil
}

// Finish VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) Finish(pointer executor.MemPtr, length executor.MemLength) error {
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(returnDataName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Finish
	gas := math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(length))
	gasToUse = math.AddUint64(gasToUse, gas)
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return err
	}

	data, err := context.MemLoad(pointer, length)
	if err != nil {
		return err
	}

	output.Finish(data)
	return nil
}

// ExecuteOnSameContext VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ExecuteOnSameContext(
	gasLimit int64,
	addressOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (int32, error) {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(executeOnSameContextName)

	return context.ExecuteOnSameContextWithHost(
		host,
		gasLimit,
		addressOffset,
		valueOffset,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
}

// ExecuteOnSameContextWithHost - executeOnSameContext with host instead of pointer context
func (context *VMHooksImpl) ExecuteOnSameContextWithHost(
	host vmhost.VMHost,
	gasLimit int64,
	addressOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (int32, error) {

	callArgs, err := context.extractIndirectContractCallArgumentsWithValue(
		host, addressOffset, valueOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if err != nil {
		return 1, err
	}

	return ExecuteOnSameContextWithTypedArgs(
		host,
		gasLimit,
		callArgs.value,
		callArgs.function,
		callArgs.dest,
		callArgs.args,
	)
}

// ExecuteOnSameContextWithTypedArgs - executeOnSameContext with args already read from memory
func ExecuteOnSameContextWithTypedArgs(
	host vmhost.VMHost,
	gasLimit int64,
	value *big.Int,
	function []byte,
	dest []byte,
	args [][]byte,
) (int32, error) {
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.ExecuteOnSameContext
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return -1, err
	}

	sender := runtime.GetContextAddress()

	contractCallInput, err := prepareIndirectContractCallInput(
		host,
		sender,
		value,
		gasLimit,
		dest,
		function,
		args,
		gasToUse,
		true,
	)
	if err != nil {
		return -1, err
	}

	if host.IsBuiltinFunctionName(contractCallInput.Function) {
		return 1, vmhost.ErrInvalidBuiltInFunctionCall
	}

	err = host.ExecuteOnSameContext(contractCallInput)
	if err != nil {
		return -1, err
	}

	return 0, nil
}

// ExecuteOnDestContext VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ExecuteOnDestContext(
	gasLimit int64,
	addressOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (int32, error) {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(executeOnDestContextName)

	return context.ExecuteOnDestContextWithHost(
		host,
		gasLimit,
		addressOffset,
		valueOffset,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
}

// ExecuteOnDestContextWithHost - executeOnDestContext with host instead of pointer context
func (context *VMHooksImpl) ExecuteOnDestContextWithHost(
	host vmhost.VMHost,
	gasLimit int64,
	addressOffset executor.MemPtr,
	valueOffset executor.MemPtr,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (int32, error) {

	callArgs, err := context.extractIndirectContractCallArgumentsWithValue(
		host, addressOffset, valueOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if err != nil {
		return 1, err
	}

	return ExecuteOnDestContextWithTypedArgs(
		host,
		gasLimit,
		callArgs.value,
		callArgs.function,
		callArgs.dest,
		callArgs.args,
		true,
	)
}

// ExecuteOnDestContextWithTypedArgs - executeOnDestContext with args already read from memory
func ExecuteOnDestContextWithTypedArgs(
	host vmhost.VMHost,
	gasLimit int64,
	value *big.Int,
	function []byte,
	dest []byte,
	args [][]byte,
	failExecution bool,
) (int32, error) {
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.ExecuteOnDestContext
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return -1, err
	}

	sender := runtime.GetContextAddress()

	contractCallInput, err := prepareIndirectContractCallInput(
		host,
		sender,
		value,
		gasLimit,
		dest,
		function,
		args,
		gasToUse,
		true,
	)
	if err != nil {
		return 1, err
	}

	vmOutput, err := executeOnDestContextFromAPI(host, contractCallInput)
	if err != nil {
		if vmOutput == nil || failExecution {
			return 1, err
		}
		return 1, nil
	}

	host.CompleteLogEntriesWithCallType(vmOutput, vmhost.ExecuteOnDestContextString)

	return 0, nil
}

// ExecuteReadOnly VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ExecuteReadOnly(
	gasLimit int64,
	addressOffset executor.MemPtr,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (int32, error) {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(executeReadOnlyName)

	return context.ExecuteReadOnlyWithHost(
		host,
		gasLimit,
		addressOffset,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
}

// ExecuteReadOnlyWithHost - executeReadOnly with host instead of pointer context
func (context *VMHooksImpl) ExecuteReadOnlyWithHost(
	host vmhost.VMHost,
	gasLimit int64,
	addressOffset executor.MemPtr,
	functionOffset executor.MemPtr,
	functionLength executor.MemLength,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (int32, error) {

	callArgs, err := context.extractIndirectContractCallArgumentsWithoutValue(
		host, addressOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if err != nil {
		return -1, err
	}

	return ExecuteReadOnlyWithTypedArguments(
		host,
		gasLimit,
		callArgs.function,
		callArgs.dest,
		callArgs.args,
	)
}

// ExecuteReadOnlyWithTypedArguments - executeReadOnly with args already read from memory
func ExecuteReadOnlyWithTypedArguments(
	host vmhost.VMHost,
	gasLimit int64,
	function []byte,
	dest []byte,
	args [][]byte,
) (int32, error) {
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.ExecuteReadOnly
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		return -1, err
	}

	sender := runtime.GetContextAddress()

	contractCallInput, err := prepareIndirectContractCallInput(
		host,
		sender,
		big.NewInt(0),
		gasLimit,
		dest,
		function,
		args,
		gasToUse,
		true,
	)
	if err != nil {
		return -1, err
	}

	if host.IsBuiltinFunctionName(contractCallInput.Function) {
		return 1, vmhost.ErrInvalidBuiltInFunctionCall
	}

	wasReadOnly := runtime.ReadOnly()
	runtime.SetReadOnly(true)
	_, err = executeOnDestContextFromAPI(host, contractCallInput)
	runtime.SetReadOnly(wasReadOnly)

	if err != nil {
		return -1, err
	}

	return 0, nil
}

// CreateContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) CreateContract(
	gasLimit int64,
	valueOffset executor.MemPtr,
	codeOffset executor.MemPtr,
	codeMetadataOffset executor.MemPtr,
	length executor.MemLength,
	resultOffset executor.MemPtr,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (int32, error) {
	host := context.GetVMHost()
	return context.createContractWithHost(
		host,
		gasLimit,
		valueOffset,
		codeOffset,
		codeMetadataOffset,
		length,
		resultOffset,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
}

func (context *VMHooksImpl) createContractWithHost(
	host vmhost.VMHost,
	gasLimit int64,
	valueOffset executor.MemPtr,
	codeOffset executor.MemPtr,
	codeMetadataOffset executor.MemPtr,
	length executor.MemLength,
	resultOffset executor.MemPtr,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (int32, error) {
	runtime := host.Runtime()

	metering := host.Metering()
	metering.StartGasTracing(createContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return -1, err
	}

	sender := runtime.GetContextAddress()
	value, err := context.MemLoad(valueOffset, vmhost.BalanceLen)
	if err != nil {
		return 1, err
	}

	code, err := context.MemLoad(codeOffset, length)
	if err != nil {
		return 1, err
	}

	codeMetadata, err := context.MemLoad(codeMetadataOffset, vmhost.CodeMetadataLen)
	if err != nil {
		return 1, err
	}

	data, _, err := context.getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if err != nil {
		return 1, err
	}

	valueAsInt := big.NewInt(0).SetBytes(value)
	newAddress, err := createContract(sender, data, valueAsInt, gasLimit, code, codeMetadata, host, CreateContract)

	if err != nil {
		return 1, err
	}

	err = context.MemStore(resultOffset, newAddress)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

// DeployFromSourceContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) DeployFromSourceContract(
	gasLimit int64,
	valueOffset executor.MemPtr,
	sourceContractAddressOffset executor.MemPtr,
	codeMetadataOffset executor.MemPtr,
	resultAddressOffset executor.MemPtr,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) (int32, error) {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(deployFromSourceContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return 1, err
	}

	value, err := context.MemLoad(valueOffset, vmhost.BalanceLen)
	if err != nil {
		return 1, err
	}

	sourceContractAddress, err := context.MemLoad(sourceContractAddressOffset, vmhost.AddressLen)
	if err != nil {
		return 1, err
	}

	codeMetadata, err := context.MemLoad(codeMetadataOffset, vmhost.CodeMetadataLen)
	if err != nil {
		return 1, err
	}

	data, _, err := context.getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if err != nil {
		return 1, err
	}

	newAddress, err := DeployFromSourceContractWithTypedArgs(
		host,
		sourceContractAddress,
		codeMetadata,
		big.NewInt(0).SetBytes(value),
		data,
		gasLimit,
	)

	if err != nil {
		return 1, err
	}

	err = context.MemStore(resultAddressOffset, newAddress)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

// DeployFromSourceContractWithTypedArgs - deployFromSourceContract with args already read from memory
func DeployFromSourceContractWithTypedArgs(
	host vmhost.VMHost,
	sourceContractAddress []byte,
	codeMetadata []byte,
	value *big.Int,
	data [][]byte,
	gasLimit int64,
) ([]byte, error) {
	runtime := host.Runtime()
	sender := runtime.GetContextAddress()

	blockchain := host.Blockchain()
	code, err := blockchain.GetCode(sourceContractAddress)
	if err != nil {
		return nil, err
	}

	return createContract(sender, data, value, gasLimit, code, codeMetadata, host, DeployContract)
}

func createContract(
	sender []byte,
	data [][]byte,
	value *big.Int,
	gasLimit int64,
	code []byte,
	codeMetadata []byte,
	host vmhost.VMHost,
	createContractCallType CreateContractCallType,
) ([]byte, error) {
	originalCaller := host.Runtime().GetOriginalCallerAddress()
	metering := host.Metering()
	contractCreate := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			OriginalCallerAddr: originalCaller,
			CallerAddr:         sender,
			Arguments:          data,
			CallValue:          value,
			GasPrice:           0,
			GasProvided:        metering.BoundGasLimit(gasLimit),
		},
		ContractCode:         code,
		ContractCodeMetadata: codeMetadata,
	}

	currentVMInput := host.Runtime().GetVMInput()
	if len(currentVMInput.RelayerAddr) > 0 {
		contractCreate.RelayerAddr = make([]byte, len(currentVMInput.RelayerAddr))
		copy(contractCreate.RelayerAddr, currentVMInput.RelayerAddr)
	}

	return host.CreateNewContract(contractCreate, int(createContractCallType))
}

// GetNumReturnData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetNumReturnData() (int32, error) {
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetNumReturnData
	err := metering.UseGasBoundedAndAddTracedGas(getNumReturnDataName, gasToUse)
	if err != nil {
		return -1, err
	}

	returnData := output.ReturnData()
	return int32(len(returnData)), nil
}

// GetReturnDataSize VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetReturnDataSize(resultID int32) (int32, error) {
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetReturnDataSize
	err := metering.UseGasBoundedAndAddTracedGas(getReturnDataSizeName, gasToUse)
	if err != nil {
		return -1, err
	}

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) || resultID < 0 {
		return 0, vmhost.ErrInvalidArgument
	}

	return int32(len(returnData[resultID])), nil
}

// GetReturnData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetReturnData(resultID int32, dataOffset executor.MemPtr) (int32, error) {
	host := context.GetVMHost()

	result, err := GetReturnDataWithHostAndTypedArgs(host, resultID)
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, nil
	}

	err = context.MemStore(dataOffset, result)
	if err != nil {
		return 0, err
	}

	return int32(len(result)), nil
}

func GetReturnDataWithHostAndTypedArgs(host vmhost.VMHost, resultID int32) ([]byte, error) {
	output := host.Output()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetReturnData
	err := metering.UseGasBoundedAndAddTracedGas(getReturnDataName, gasToUse)
	if err != nil {
		return nil, err
	}

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) || resultID < 0 {
		return nil, vmhost.ErrInvalidArgument
	}

	return returnData[resultID], nil
}

// CleanReturnData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) CleanReturnData() error {
	host := context.GetVMHost()
	return CleanReturnDataWithHost(host)
}

// CleanReturnDataWithHost - exposed version of v1_5_deleteFromReturnData for tests
func CleanReturnDataWithHost(host vmhost.VMHost) error {
	output := host.Output()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CleanReturnData
	err := metering.UseGasBoundedAndAddTracedGas(cleanReturnDataName, gasToUse)
	if err != nil {
		return err
	}

	output.ClearReturnData()
	return nil
}

// DeleteFromReturnData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) DeleteFromReturnData(resultID int32) error {
	host := context.GetVMHost()
	return DeleteFromReturnDataWithHost(host, resultID)
}

// DeleteFromReturnDataWithHost - exposed version of v1_5_deleteFromReturnData for tests
func DeleteFromReturnDataWithHost(host vmhost.VMHost, resultID int32) error {
	output := host.Output()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.DeleteFromReturnData
	err := metering.UseGasBoundedAndAddTracedGas(deleteFromReturnDataName, gasToUse)
	if err != nil {
		return err
	}

	returnData := output.ReturnData()
	if resultID >= 0 && resultID < int32(len(returnData)) {
		output.RemoveReturnData(uint32(resultID))
	}
	return nil
}

// GetOriginalTxHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetOriginalTxHash(dataOffset executor.MemPtr) error {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetOriginalTxHash
	err := metering.UseGasBoundedAndAddTracedGas(getOriginalTxHashName, gasToUse)
	if err != nil {
		return err
	}

	err = context.MemStore(dataOffset, runtime.GetOriginalTxHash())
	if err != nil {
		return err
	}
	return nil
}

// GetCurrentTxHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetCurrentTxHash(dataOffset executor.MemPtr) error {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCurrentTxHash
	err := metering.UseGasBoundedAndAddTracedGas(getCurrentTxHashName, gasToUse)
	if err != nil {
		return err
	}

	err = context.MemStore(dataOffset, runtime.GetCurrentTxHash())
	if err != nil {
		return err
	}
	return nil
}

// GetPrevTxHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) GetPrevTxHash(dataOffset executor.MemPtr) error {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetPrevTxHash
	err := metering.UseGasBoundedAndAddTracedGas(getPrevTxHashName, gasToUse)
	if err != nil {
		return err
	}

	err = context.MemStore(dataOffset, runtime.GetPrevTxHash())
	if err != nil {
		return err
	}
	return nil
}

func prepareIndirectContractCallInput(
	host vmhost.VMHost,
	sender []byte,
	value *big.Int,
	gasLimit int64,
	destination []byte,
	function []byte,
	data [][]byte,
	_ uint64,
	syncExecutionRequired bool,
) (*vmcommon.ContractCallInput, error) {
	runtime := host.Runtime()
	metering := host.Metering()

	if syncExecutionRequired && !host.AreInSameShard(runtime.GetContextAddress(), destination) {
		return nil, vmhost.ErrSyncExecutionNotInSameShard
	}

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			OriginalCallerAddr: host.Runtime().GetOriginalCallerAddress(),
			CallerAddr:         sender,
			Arguments:          data,
			CallValue:          value,
			GasPrice:           0,
			GasProvided:        metering.BoundGasLimit(gasLimit),
			CallType:           vm.DirectCall,
		},
		RecipientAddr: bytes.Clone(destination),
		Function:      string(function),
	}

	currentVMInput := runtime.GetVMInput()
	if len(currentVMInput.RelayerAddr) > 0 {
		contractCallInput.RelayerAddr = make([]byte, len(currentVMInput.RelayerAddr))
		copy(contractCallInput.RelayerAddr, currentVMInput.RelayerAddr)
	}

	return contractCallInput, nil
}

func (context *VMHooksImpl) getArgumentsFromMemory(
	_ vmhost.VMHost,
	numArguments int32,
	argumentsLengthOffset executor.MemPtr,
	dataOffset executor.MemPtr,
) ([][]byte, int32, error) {
	maxArgs := int32(maxNumArgumentsFromMemory)
	if context.GetRuntimeContext().FixCorrectSmallIntAndGasLimits() {
		maxArgs = maxNumArgumentsFromMemoryCapped
	}
	if numArguments < 0 || numArguments > maxArgs {
		return nil, 0, fmt.Errorf("numArguments out of range (%d)", numArguments)
	}

	metering := context.GetMeteringContext()
	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(numArguments))
	err := metering.UseGasBounded(gasToUse)
	if err != nil && context.GetRuntimeContext().UseGasBoundedShouldFailExecution() {
		return nil, 0, err
	}

	argumentsLengthData, err := context.MemLoad(argumentsLengthOffset, numArguments*4)
	if err != nil {
		return nil, 0, err
	}
	if context.GetRuntimeContext().FixCorrectSmallIntAndGasLimits() && len(argumentsLengthData) < int(numArguments)*4 {
		return nil, 0, fmt.Errorf("arguments length data truncated: %w", executor.ErrMemoryBadBounds)
	}

	argumentLengths := createInt32Array(argumentsLengthData, numArguments)
	data, err := context.MemLoadMultiple(dataOffset, argumentLengths)
	if err != nil {
		return nil, 0, err
	}

	totalArgumentBytes := int32(0)
	for _, length := range argumentLengths {
		totalArgumentBytes += length
	}

	restOfLenForGasCopy := math.SubUint64(uint64(totalArgumentBytes), uint64(numArguments))
	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, restOfLenForGasCopy)
	err = metering.UseGasBounded(gasToUse)
	if err != nil && context.GetRuntimeContext().UseGasBoundedShouldFailExecution() {
		return nil, 0, err
	}

	return data, totalArgumentBytes, nil
}

func createInt32Array(rawData []byte, numIntegers int32) []int32 {
	integers := make([]int32, numIntegers)
	index := 0
	for cursor := 0; cursor+4 <= len(rawData) && index < len(integers); cursor += 4 {
		rawInt := rawData[cursor : cursor+4]
		actualInt := binary.LittleEndian.Uint32(rawInt)
		integers[index] = int32(actualInt)
		index++
	}
	return integers
}

func executeOnDestContextFromAPI(host vmhost.VMHost, input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	host.Async().SetAsyncArgumentsForCall(input)
	vmOutput, isChildComplete, err := host.ExecuteOnDestContext(input)
	if err != nil {
		if host.Runtime().FixCorrectSmallIntAndGasLimits() {
			_ = host.Async().CompleteChildConditional(true, nil, 0)
		}
		return vmOutput, err
	}

	err = host.Async().CompleteChildConditional(isChildComplete, nil, 0)
	if err != nil {
		return nil, err
	}

	return vmOutput, err
}

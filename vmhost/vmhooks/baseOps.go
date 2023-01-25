package vmhooks

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern void		v1_4_getSCAddress(void *context, int32_t resultOffset);
// extern void		v1_4_getOwnerAddress(void *context, int32_t resultOffset);
// extern int32_t	v1_4_getShardOfAddress(void *context, int32_t addressOffset);
// extern int32_t	v1_4_isSmartContract(void *context, int32_t addressOffset);
// extern void		v1_4_getExternalBalance(void *context, int32_t addressOffset, int32_t resultOffset);
// extern int32_t	v1_4_blockHash(void *context, long long nonce, int32_t resultOffset);
// extern int32_t	v1_4_transferValue(void *context, int32_t dstOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
// extern int32_t	v1_4_transferESDTExecute(void *context, int32_t dstOffset, int32_t tokenIDOffset, int32_t tokenIdLen, int32_t valueOffset, long long gasLimit, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t	v1_4_transferESDTNFTExecute(void *context, int32_t dstOffset, int32_t tokenIDOffset, int32_t tokenIdLen, int32_t valueOffset, long long nonce, long long gasLimit, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t	v1_4_multiTransferESDTNFTExecute(void *context, int32_t dstOffset, int32_t numTokenTransfers, int32_t tokenTransfersArgsLengthOffset, int32_t tokenTransferDataOffset, long long gasLimit, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t	v1_4_transferValueExecute(void *context, int32_t dstOffset, int32_t valueOffset, long long gasLimit, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t	v1_4_getArgumentLength(void *context, int32_t id);
// extern int32_t	v1_4_getArgument(void *context, int32_t id, int32_t argOffset);
// extern int32_t	v1_4_getFunction(void *context, int32_t functionOffset);
// extern int32_t	v1_4_getNumArguments(void *context);
// extern int32_t	v1_4_storageStore(void *context, int32_t keyOffset, int32_t keyLength , int32_t dataOffset, int32_t dataLength);
// extern int32_t	v1_4_storageLoadLength(void *context, int32_t keyOffset, int32_t keyLength );
// extern int32_t	v1_4_storageLoad(void *context, int32_t keyOffset, int32_t keyLength , int32_t dataOffset);
// extern int32_t	v1_4_storageLoadFromAddress(void *context, int32_t addressOffset, int32_t keyOffset, int32_t keyLength , int32_t dataOffset);
// extern void		v1_4_getCaller(void *context, int32_t resultOffset);
// extern void		v1_4_checkNoPayment(void *context);
// extern int32_t	v1_4_callValue(void *context, int32_t resultOffset);
// extern int32_t	v1_4_getESDTValue(void *context, int32_t resultOffset);
// extern int32_t	v1_4_getESDTTokenName(void *context, int32_t resultOffset);
// extern long long	v1_4_getESDTTokenNonce(void *context);
// extern int32_t	v1_4_getESDTTokenType(void *context);
// extern int32_t	v1_4_getCallValueTokenName(void *context, int32_t callValueOffset, int32_t tokenNameOffset);
// extern int32_t	v1_4_getESDTValueByIndex(void *context, int32_t resultOffset, int32_t index);
// extern int32_t	v1_4_getESDTTokenNameByIndex(void *context, int32_t resultOffset, int32_t index);
// extern long long	v1_4_getESDTTokenNonceByIndex(void *context, int32_t index);
// extern int32_t	v1_4_getESDTTokenTypeByIndex(void *context, int32_t index);
// extern int32_t	v1_4_getCallValueTokenNameByIndex(void *context, int32_t callValueOffset, int32_t tokenNameOffset, int32_t index);
// extern int32_t	v1_4_getNumESDTTransfers(void *context);
// extern long long v1_4_getCurrentESDTNFTNonce(void *context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen);
// extern void		v1_4_writeLog(void *context, int32_t pointer, int32_t length, int32_t topicPtr, int32_t numTopics);
// extern void		v1_4_writeEventLog(void *context, int32_t numTopics, int32_t topicLengthsOffset, int32_t topicOffset, int32_t dataOffset, int32_t dataLength);
// extern void		v1_4_returnData(void* context, int32_t dataOffset, int32_t length);
// extern void		v1_4_signalError(void* context, int32_t messageOffset, int32_t messageLength);
// extern long long v1_4_getGasLeft(void *context);
// extern int32_t	v1_4_getESDTBalance(void *context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen, long long nonce, int32_t resultOffset);
// extern int32_t	v1_4_getESDTNFTNameLength(void *context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen, long long nonce);
// extern int32_t	v1_4_getESDTNFTAttributeLength(void *context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen, long long nonce);
// extern int32_t	v1_4_getESDTNFTURILength(void *context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen, long long nonce);
// extern int32_t	v1_4_getESDTTokenData(void *context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen, long long nonce, int32_t valueOffset, int32_t propertiesOffset, int32_t hashOffset, int32_t nameOffset, int32_t attributesOffset, int32_t creatorOffset, int32_t royaltiesOffset, int32_t urisOffset);
// extern long long	v1_4_getESDTLocalRoles(void *context, int32_t tokenIdHandle);
// extern int32_t	v1_4_validateTokenIdentifier(void *context, int32_t tokenIdHandle);
//
// extern int32_t	v1_4_executeOnDestContext(void *context, long long gas, int32_t addressOffset, int32_t valueOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t	v1_4_executeOnDestContextByCaller(void *context, long long gas, int32_t addressOffset, int32_t valueOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t	v1_4_executeOnSameContext(void *context, long long gas, int32_t addressOffset, int32_t valueOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t	v1_4_executeReadOnly(void *context, long long gas, int32_t addressOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t	v1_4_createContract(void *context, long long gas, int32_t valueOffset, int32_t codeOffset, int32_t codeMetadataOffset, int32_t length, int32_t resultOffset, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t	v1_4_deployFromSourceContract(void *context, long long gas, int32_t valueOffset, int32_t addressOffset, int32_t codeMetadataOffset, int32_t resultOffset, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern void		v1_4_upgradeContract(void *context, int32_t dstOffset, long long gas, int32_t valueOffset, int32_t codeOffset, int32_t codeMetadataOffset, int32_t length, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern void		v1_4_upgradeFromSourceContract(void *context, int32_t dstOffset, long long gas, int32_t valueOffset, int32_t addressOffset, int32_t codeMetadataOffset, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern void		v1_4_asyncCall(void *context, int32_t dstOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
//
// extern int32_t	v1_4_getNumReturnData(void *context);
// extern int32_t	v1_4_getReturnDataSize(void *context, int32_t resultID);
// extern int32_t	v1_4_getReturnData(void *context, int32_t resultID, int32_t dataOffset);
// extern void		v1_4_cleanReturnData(void *context);
// extern void		v1_4_deleteFromReturnData(void *context, int32_t resultID);
//
// extern int32_t	v1_4_setStorageLock(void *context, int32_t keyOffset, int32_t keyLength, long long lockTimestamp);
// extern long long v1_4_getStorageLock(void *context, int32_t keyOffset, int32_t keyLength);
// extern int32_t	v1_4_isStorageLocked(void *context, int32_t keyOffset, int32_t keyLength);
// extern int32_t	v1_4_clearStorageLock(void *context, int32_t keyOffset, int32_t keyLength);
//
// extern long long v1_4_getBlockTimestamp(void *context);
// extern long long v1_4_getBlockNonce(void *context);
// extern long long v1_4_getBlockRound(void *context);
// extern long long v1_4_getBlockEpoch(void *context);
// extern void		v1_4_getBlockRandomSeed(void *context, int32_t resultOffset);
// extern void		v1_4_getStateRootHash(void *context, int32_t resultOffset);
//
// extern long long v1_4_getPrevBlockTimestamp(void *context);
// extern long long v1_4_getPrevBlockNonce(void *context);
// extern long long v1_4_getPrevBlockRound(void *context);
// extern long long v1_4_getPrevBlockEpoch(void *context);
// extern void		v1_4_getPrevBlockRandomSeed(void *context, int32_t resultOffset);
// extern void		v1_4_getOriginalTxHash(void *context, int32_t resultOffset);
import "C"

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"unsafe"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/esdt"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	logger "github.com/multiversx/mx-chain-logger-go"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	"github.com/multiversx/mx-chain-vm-v1_4-go/math"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/vmhooksmeta"
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
	executeOnDestContextByCallerName = "executeOnDestContextByCaller"
	executeOnSameContextName         = "executeOnSameContext"
	executeReadOnlyName              = "executeReadOnly"
	createContractName               = "createContract"
	deployFromSourceContractName     = "deployFromSourceContract"
	upgradeContractName              = "upgradeContract"
	upgradeFromSourceContractName    = "upgradeFromSourceContract"
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
	getBlockNonceName                = "getBlockNonce"
	getBlockRoundName                = "getBlockRound"
	getBlockEpochName                = "getBlockEpoch"
	getBlockRandomSeedName           = "getBlockRandomSeed"
	getStateRootHashName             = "getStateRootHash"
	getPrevBlockTimestampName        = "getPrevBlockTimestamp"
	getPrevBlockNonceName            = "getPrevBlockNonce"
	getPrevBlockRoundName            = "getPrevBlockRound"
	getPrevBlockEpochName            = "getPrevBlockEpoch"
	getPrevBlockRandomSeedName       = "getPrevBlockRandomSeed"
	getOriginalTxHashName            = "getOriginalTxHash"
)

var logEEI = logger.GetOrCreate("vm/eei")

func getESDTTransferFromInputFailIfWrongIndex(host vmhost.VMHost, index int32) *vmcommon.ESDTTransfer {
	esdtTransfers := host.Runtime().GetVMInput().ESDTTransfers
	if int32(len(esdtTransfers))-1 < index || index < 0 {
		vmhost.WithFaultAndHostIfFailAlwaysActive(vmhost.ErrInvalidTokenIndex, host, host.Runtime().BaseOpsErrorShouldFailExecution())
		return nil
	}
	return esdtTransfers[index]
}

func failIfMoreThanOneESDTTransfer(context unsafe.Pointer) bool {
	runtime := vmhost.GetRuntimeContext(context)
	if len(runtime.GetVMInput().ESDTTransfers) > 1 {
		return vmhost.WithFault(vmhost.ErrTooManyESDTTransfers, context, true)
	}
	return false
}

// BaseOpsAPIImports creates a new wasmer.Imports populated with the BaseOpsAPI API methods
func BaseOpsAPIImports(imports vmhooksmeta.EIFunctionReceiver) error {
	imports.Namespace("env")

	err := imports.Append("getSCAddress", v1_4_getSCAddress, C.v1_4_getSCAddress)
	if err != nil {
		return err
	}

	err = imports.Append("getOwnerAddress", v1_4_getOwnerAddress, C.v1_4_getOwnerAddress)
	if err != nil {
		return err
	}

	err = imports.Append("getShardOfAddress", v1_4_getShardOfAddress, C.v1_4_getShardOfAddress)
	if err != nil {
		return err
	}

	err = imports.Append("isSmartContract", v1_4_isSmartContract, C.v1_4_isSmartContract)
	if err != nil {
		return err
	}

	err = imports.Append("getExternalBalance", v1_4_getExternalBalance, C.v1_4_getExternalBalance)
	if err != nil {
		return err
	}

	err = imports.Append("getBlockHash", v1_4_blockHash, C.v1_4_blockHash)
	if err != nil {
		return err
	}

	err = imports.Append("transferValue", v1_4_transferValue, C.v1_4_transferValue)
	if err != nil {
		return err
	}

	err = imports.Append("transferESDTExecute", v1_4_transferESDTExecute, C.v1_4_transferESDTExecute)
	if err != nil {
		return err
	}

	err = imports.Append("transferESDTNFTExecute", v1_4_transferESDTNFTExecute, C.v1_4_transferESDTNFTExecute)
	if err != nil {
		return err
	}

	err = imports.Append("multiTransferESDTNFTExecute", v1_4_multiTransferESDTNFTExecute, C.v1_4_multiTransferESDTNFTExecute)
	if err != nil {
		return err
	}

	err = imports.Append("transferValueExecute", v1_4_transferValueExecute, C.v1_4_transferValueExecute)
	if err != nil {
		return err
	}

	err = imports.Append("asyncCall", v1_4_asyncCall, C.v1_4_asyncCall)
	if err != nil {
		return err
	}

	err = imports.Append("getArgumentLength", v1_4_getArgumentLength, C.v1_4_getArgumentLength)
	if err != nil {
		return err
	}

	err = imports.Append("getArgument", v1_4_getArgument, C.v1_4_getArgument)
	if err != nil {
		return err
	}

	err = imports.Append("getFunction", v1_4_getFunction, C.v1_4_getFunction)
	if err != nil {
		return err
	}

	err = imports.Append("getNumArguments", v1_4_getNumArguments, C.v1_4_getNumArguments)
	if err != nil {
		return err
	}

	err = imports.Append("storageStore", v1_4_storageStore, C.v1_4_storageStore)
	if err != nil {
		return err
	}

	err = imports.Append("storageLoadLength", v1_4_storageLoadLength, C.v1_4_storageLoadLength)
	if err != nil {
		return err
	}

	err = imports.Append("storageLoad", v1_4_storageLoad, C.v1_4_storageLoad)
	if err != nil {
		return err
	}

	err = imports.Append("storageLoadFromAddress", v1_4_storageLoadFromAddress, C.v1_4_storageLoadFromAddress)
	if err != nil {
		return err
	}

	err = imports.Append("getStorageLock", v1_4_getStorageLock, C.v1_4_getStorageLock)
	if err != nil {
		return err
	}

	err = imports.Append("setStorageLock", v1_4_setStorageLock, C.v1_4_setStorageLock)
	if err != nil {
		return err
	}

	err = imports.Append("isStorageLocked", v1_4_isStorageLocked, C.v1_4_isStorageLocked)
	if err != nil {
		return err
	}

	err = imports.Append("clearStorageLock", v1_4_clearStorageLock, C.v1_4_clearStorageLock)
	if err != nil {
		return err
	}

	err = imports.Append("getCaller", v1_4_getCaller, C.v1_4_getCaller)
	if err != nil {
		return err
	}

	err = imports.Append("checkNoPayment", v1_4_checkNoPayment, C.v1_4_checkNoPayment)
	if err != nil {
		return err
	}

	err = imports.Append("getCallValue", v1_4_callValue, C.v1_4_callValue)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTValue", v1_4_getESDTValue, C.v1_4_getESDTValue)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTTokenName", v1_4_getESDTTokenName, C.v1_4_getESDTTokenName)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTTokenType", v1_4_getESDTTokenType, C.v1_4_getESDTTokenType)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTTokenNonce", v1_4_getESDTTokenNonce, C.v1_4_getESDTTokenNonce)
	if err != nil {
		return err
	}

	err = imports.Append("getCallValueTokenName", v1_4_getCallValueTokenName, C.v1_4_getCallValueTokenName)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTValueByIndex", v1_4_getESDTValueByIndex, C.v1_4_getESDTValueByIndex)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTTokenNameByIndex", v1_4_getESDTTokenNameByIndex, C.v1_4_getESDTTokenNameByIndex)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTTokenTypeByIndex", v1_4_getESDTTokenTypeByIndex, C.v1_4_getESDTTokenTypeByIndex)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTTokenNonceByIndex", v1_4_getESDTTokenNonceByIndex, C.v1_4_getESDTTokenNonceByIndex)
	if err != nil {
		return err
	}

	err = imports.Append("getCallValueTokenNameByIndex", v1_4_getCallValueTokenNameByIndex, C.v1_4_getCallValueTokenNameByIndex)
	if err != nil {
		return err
	}

	err = imports.Append("getNumESDTTransfers", v1_4_getNumESDTTransfers, C.v1_4_getNumESDTTransfers)
	if err != nil {
		return err
	}

	err = imports.Append("getCurrentESDTNFTNonce", v1_4_getCurrentESDTNFTNonce, C.v1_4_getCurrentESDTNFTNonce)
	if err != nil {
		return err
	}

	err = imports.Append("validateTokenIdentifier", v1_4_validateTokenIdentifier, C.v1_4_validateTokenIdentifier)
	if err != nil {
		return err
	}

	err = imports.Append("writeLog", v1_4_writeLog, C.v1_4_writeLog)
	if err != nil {
		return err
	}

	err = imports.Append("writeEventLog", v1_4_writeEventLog, C.v1_4_writeEventLog)
	if err != nil {
		return err
	}

	err = imports.Append("finish", v1_4_returnData, C.v1_4_returnData)
	if err != nil {
		return err
	}

	err = imports.Append("signalError", v1_4_signalError, C.v1_4_signalError)
	if err != nil {
		return err
	}

	err = imports.Append("getBlockTimestamp", v1_4_getBlockTimestamp, C.v1_4_getBlockTimestamp)
	if err != nil {
		return err
	}

	err = imports.Append("getBlockNonce", v1_4_getBlockNonce, C.v1_4_getBlockNonce)
	if err != nil {
		return err
	}

	err = imports.Append("getBlockRound", v1_4_getBlockRound, C.v1_4_getBlockRound)
	if err != nil {
		return err
	}

	err = imports.Append("getBlockEpoch", v1_4_getBlockEpoch, C.v1_4_getBlockEpoch)
	if err != nil {
		return err
	}

	err = imports.Append("getBlockRandomSeed", v1_4_getBlockRandomSeed, C.v1_4_getBlockRandomSeed)
	if err != nil {
		return err
	}

	err = imports.Append("getStateRootHash", v1_4_getStateRootHash, C.v1_4_getStateRootHash)
	if err != nil {
		return err
	}

	err = imports.Append("getPrevBlockTimestamp", v1_4_getPrevBlockTimestamp, C.v1_4_getPrevBlockTimestamp)
	if err != nil {
		return err
	}

	err = imports.Append("getPrevBlockNonce", v1_4_getPrevBlockNonce, C.v1_4_getPrevBlockNonce)
	if err != nil {
		return err
	}

	err = imports.Append("getPrevBlockRound", v1_4_getPrevBlockRound, C.v1_4_getPrevBlockRound)
	if err != nil {
		return err
	}

	err = imports.Append("getPrevBlockEpoch", v1_4_getPrevBlockEpoch, C.v1_4_getPrevBlockEpoch)
	if err != nil {
		return err
	}

	err = imports.Append("getPrevBlockRandomSeed", v1_4_getPrevBlockRandomSeed, C.v1_4_getPrevBlockRandomSeed)
	if err != nil {
		return err
	}

	err = imports.Append("getOriginalTxHash", v1_4_getOriginalTxHash, C.v1_4_getOriginalTxHash)
	if err != nil {
		return err
	}

	err = imports.Append("getGasLeft", v1_4_getGasLeft, C.v1_4_getGasLeft)
	if err != nil {
		return err
	}

	err = imports.Append("executeOnDestContext", v1_4_executeOnDestContext, C.v1_4_executeOnDestContext)
	if err != nil {
		return err
	}

	err = imports.Append("executeOnDestContextByCaller", v1_4_executeOnDestContextByCaller, C.v1_4_executeOnDestContextByCaller)
	if err != nil {
		return err
	}

	err = imports.Append("executeOnSameContext", v1_4_executeOnSameContext, C.v1_4_executeOnSameContext)
	if err != nil {
		return err
	}

	err = imports.Append("createContract", v1_4_createContract, C.v1_4_createContract)
	if err != nil {
		return err
	}

	err = imports.Append("deployFromSourceContract", v1_4_deployFromSourceContract, C.v1_4_deployFromSourceContract)
	if err != nil {
		return err
	}

	err = imports.Append("upgradeContract", v1_4_upgradeContract, C.v1_4_upgradeContract)
	if err != nil {
		return err
	}

	err = imports.Append("upgradeFromSourceContract", v1_4_upgradeFromSourceContract, C.v1_4_upgradeFromSourceContract)
	if err != nil {
		return err
	}

	err = imports.Append("executeReadOnly", v1_4_executeReadOnly, C.v1_4_executeReadOnly)
	if err != nil {
		return err
	}

	err = imports.Append("getNumReturnData", v1_4_getNumReturnData, C.v1_4_getNumReturnData)
	if err != nil {
		return err
	}

	err = imports.Append("getReturnDataSize", v1_4_getReturnDataSize, C.v1_4_getReturnDataSize)
	if err != nil {
		return err
	}

	err = imports.Append("getReturnData", v1_4_getReturnData, C.v1_4_getReturnData)
	if err != nil {
		return err
	}

	err = imports.Append("cleanReturnData", v1_4_cleanReturnData, C.v1_4_cleanReturnData)
	if err != nil {
		return err
	}

	err = imports.Append("deleteFromReturnData", v1_4_deleteFromReturnData, C.v1_4_deleteFromReturnData)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTBalance", v1_4_getESDTBalance, C.v1_4_getESDTBalance)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTTokenData", v1_4_getESDTTokenData, C.v1_4_getESDTTokenData)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTLocalRoles", v1_4_getESDTLocalRoles, C.v1_4_getESDTLocalRoles)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTNFTNameLength", v1_4_getESDTNFTNameLength, C.v1_4_getESDTNFTNameLength)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTNFTAttributeLength", v1_4_getESDTNFTAttributeLength, C.v1_4_getESDTNFTAttributeLength)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTNFTURILength", v1_4_getESDTNFTURILength, C.v1_4_getESDTNFTURILength)
	if err != nil {
		return err
	}

	return nil
}

//export v1_4_getGasLeft
func v1_4_getGasLeft(context unsafe.Pointer) int64 {
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetGasLeft
	metering.UseGasAndAddTracedGas(getGasLeftName, gasToUse)

	return int64(metering.GasLeft())
}

//export v1_4_getSCAddress
func v1_4_getSCAddress(context unsafe.Pointer, resultOffset int32) {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetSCAddress
	metering.UseGasAndAddTracedGas(getSCAddressName, gasToUse)

	owner := runtime.GetContextAddress()
	err := runtime.MemStore(resultOffset, owner)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}
}

//export v1_4_getOwnerAddress
func v1_4_getOwnerAddress(context unsafe.Pointer, resultOffset int32) {
	blockchain := vmhost.GetBlockchainContext(context)
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetOwnerAddress
	metering.UseGasAndAddTracedGas(getOwnerAddressName, gasToUse)

	owner, err := blockchain.GetOwnerAddress()
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	err = runtime.MemStore(resultOffset, owner)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}
}

//export v1_4_getShardOfAddress
func v1_4_getShardOfAddress(context unsafe.Pointer, addressOffset int32) int32 {
	blockchain := vmhost.GetBlockchainContext(context)
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetShardOfAddress
	metering.UseGasAndAddTracedGas(getShardOfAddressName, gasToUse)

	address, err := runtime.MemLoad(addressOffset, vmhost.AddressLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(blockchain.GetShardOfAddress(address))
}

//export v1_4_isSmartContract
func v1_4_isSmartContract(context unsafe.Pointer, addressOffset int32) int32 {
	blockchain := vmhost.GetBlockchainContext(context)
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.IsSmartContract
	metering.UseGasAndAddTracedGas(isSmartContractName, gasToUse)

	address, err := runtime.MemLoad(addressOffset, vmhost.AddressLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	isSmartContract := blockchain.IsSmartContract(address)

	return int32(vmhost.BooleanToInt(isSmartContract))
}

//export v1_4_signalError
func v1_4_signalError(context unsafe.Pointer, messageOffset int32, messageLength int32) {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)
	metering.StartGasTracing(signalErrorName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.SignalError
	gasToUse += metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(messageLength)

	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		_ = vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	message, err := runtime.MemLoad(messageOffset, messageLength)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}
	runtime.SignalUserError(string(message))
}

//export v1_4_getExternalBalance
func v1_4_getExternalBalance(context unsafe.Pointer, addressOffset int32, resultOffset int32) {
	blockchain := vmhost.GetBlockchainContext(context)
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	metering.UseGasAndAddTracedGas(getExternalBalanceName, gasToUse)

	address, err := runtime.MemLoad(addressOffset, vmhost.AddressLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	balance := blockchain.GetBalance(address)

	err = runtime.MemStore(resultOffset, balance)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}
}

//export v1_4_blockHash
func v1_4_blockHash(context unsafe.Pointer, nonce int64, resultOffset int32) int32 {
	blockchain := vmhost.GetBlockchainContext(context)
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockHash
	metering.UseGasAndAddTracedGas(blockHashName, gasToUse)

	hash := blockchain.BlockHash(nonce)
	err := runtime.MemStore(resultOffset, hash)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

func isBuiltInCall(data string, host vmhost.VMHost) bool {
	argParser := parsers.NewCallArgsParser()
	functionName, _, _ := argParser.ParseData(data)
	return host.IsBuiltinFunctionName(functionName)
}

func getESDTDataFromBlockchainHook(
	context unsafe.Pointer,
	addressOffset int32,
	tokenIDOffset int32,
	tokenIDLen int32,
	nonce int64,
) (*esdt.ESDigitalToken, error) {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)
	blockchain := vmhost.GetBlockchainContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	metering.UseAndTraceGas(gasToUse)

	address, err := runtime.MemLoad(addressOffset, vmhost.AddressLen)
	if err != nil {
		return nil, err
	}

	tokenID, err := runtime.MemLoad(tokenIDOffset, tokenIDLen)
	if err != nil {
		return nil, err
	}

	esdtToken, err := blockchain.GetESDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		return nil, err
	}

	return esdtToken, nil
}

//export v1_4_getESDTBalance
func v1_4_getESDTBalance(
	context unsafe.Pointer,
	addressOffset int32,
	tokenIDOffset int32,
	tokenIDLen int32,
	nonce int64,
	resultOffset int32,
) int32 {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)
	metering.StartGasTracing(getESDTBalanceName)

	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)

	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}
	err = runtime.MemStore(resultOffset, esdtData.Value.Bytes())
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(esdtData.Value.Bytes()))
}

//export v1_4_getESDTNFTNameLength
func v1_4_getESDTNFTNameLength(
	context unsafe.Pointer,
	addressOffset int32,
	tokenIDOffset int32,
	tokenIDLen int32,
	nonce int64,
) int32 {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)
	metering.StartGasTracing(getESDTNFTNameLengthName)

	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)

	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}
	if esdtData == nil || esdtData.TokenMetaData == nil {
		vmhost.WithFaultIfFailAlwaysActive(vmhost.ErrNilESDTData, context, runtime.BaseOpsErrorShouldFailExecution())
		return 0
	}

	return int32(len(esdtData.TokenMetaData.Name))
}

//export v1_4_getESDTNFTAttributeLength
func v1_4_getESDTNFTAttributeLength(
	context unsafe.Pointer,
	addressOffset int32,
	tokenIDOffset int32,
	tokenIDLen int32,
	nonce int64,
) int32 {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)
	metering.StartGasTracing(getESDTNFTAttributeLengthName)

	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)

	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}
	if esdtData == nil || esdtData.TokenMetaData == nil {
		vmhost.WithFaultIfFailAlwaysActive(vmhost.ErrNilESDTData, context, runtime.BaseOpsErrorShouldFailExecution())
		return 0
	}

	return int32(len(esdtData.TokenMetaData.Attributes))
}

//export v1_4_getESDTNFTURILength
func v1_4_getESDTNFTURILength(
	context unsafe.Pointer,
	addressOffset int32,
	tokenIDOffset int32,
	tokenIDLen int32,
	nonce int64,
) int32 {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)
	metering.StartGasTracing(getESDTNFTURILengthName)

	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)

	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}
	if esdtData == nil || esdtData.TokenMetaData == nil {
		vmhost.WithFaultIfFailAlwaysActive(vmhost.ErrNilESDTData, context, runtime.BaseOpsErrorShouldFailExecution())
		return 0
	}
	if len(esdtData.TokenMetaData.URIs) == 0 {
		return 0
	}

	return int32(len(esdtData.TokenMetaData.URIs[0]))
}

//export v1_4_getESDTTokenData
func v1_4_getESDTTokenData(
	context unsafe.Pointer,
	addressOffset int32,
	tokenIDOffset int32,
	tokenIDLen int32,
	nonce int64,
	valueHandle int32,
	propertiesOffset int32,
	hashOffset int32,
	nameOffset int32,
	attributesOffset int32,
	creatorOffset int32,
	royaltiesHandle int32,
	urisOffset int32,
) int32 {
	managedType := vmhost.GetManagedTypesContext(context)
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)
	metering.StartGasTracing(getESDTTokenDataName)

	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)

	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	value := managedType.GetBigIntOrCreate(valueHandle)
	value.Set(esdtData.Value)

	err = runtime.MemStore(propertiesOffset, esdtData.Properties)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	if esdtData.TokenMetaData != nil {
		err = runtime.MemStore(hashOffset, esdtData.TokenMetaData.Hash)
		if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
			return -1
		}
		err = runtime.MemStore(nameOffset, esdtData.TokenMetaData.Name)
		if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
			return -1
		}
		err = runtime.MemStore(attributesOffset, esdtData.TokenMetaData.Attributes)
		if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
			return -1
		}
		err = runtime.MemStore(creatorOffset, esdtData.TokenMetaData.Creator)
		if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
			return -1
		}

		royalties := managedType.GetBigIntOrCreate(royaltiesHandle)
		royalties.SetUint64(uint64(esdtData.TokenMetaData.Royalties))

		if len(esdtData.TokenMetaData.URIs) > 0 {
			err = runtime.MemStore(urisOffset, esdtData.TokenMetaData.URIs[0])
			if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
				return -1
			}
		}
	}
	return int32(len(esdtData.Value.Bytes()))
}

//export v1_4_getESDTLocalRoles
func v1_4_getESDTLocalRoles(context unsafe.Pointer, tokenIdHandle int32) int64 {
	managedType := vmhost.GetManagedTypesContext(context)
	runtime := vmhost.GetRuntimeContext(context)
	storage := vmhost.GetStorageContext(context)
	metering := vmhost.GetMeteringContext(context)

	tokenID, err := managedType.GetBytes(tokenIdHandle)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	esdtRoleKeyPrefix := []byte(core.ProtectedKeyPrefix + core.ESDTRoleIdentifier + core.ESDTKeyIdentifier)
	key := []byte(string(esdtRoleKeyPrefix) + string(tokenID))

	data, trieDepth, usedCache, err := storage.GetStorage(key)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	err = storage.UseGasForStorageLoad(storageLoadName, int64(trieDepth), metering.GasSchedule().BaseOpsAPICost.StorageLoad, usedCache)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return getESDTRoles(data)
}

//export v1_4_validateTokenIdentifier
func v1_4_validateTokenIdentifier(
	context unsafe.Pointer,
	tokenIdHandle int32,
) int32 {
	managedType := vmhost.GetManagedTypesContext(context)
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetArgument
	metering.UseGasAndAddTracedGas(validateTokenIdentifierName, gasToUse)

	tokenID, err := managedType.GetBytes(tokenIdHandle)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	if ValidateToken(tokenID) {
		return 1
	} else {
		return 0
	}

}

//export v1_4_transferValue
func v1_4_transferValue(context unsafe.Pointer, destOffset int32, valueOffset int32, dataOffset int32, length int32) int32 {
	host := vmhost.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()
	metering.StartGasTracing(transferValueName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.TransferValue
	metering.UseAndTraceGas(gasToUse)

	sender := runtime.GetContextAddress()
	dest, err := runtime.MemLoad(destOffset, vmhost.AddressLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	valueBytes, err := runtime.MemLoad(valueOffset, vmhost.BalanceLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(length))
	metering.UseAndTraceGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	if isBuiltInCall(string(data), host) {
		vmhost.WithFaultIfFailAlwaysActive(vmhost.ErrTransferValueOnESDTCall, context, runtime.BaseOpsErrorShouldFailExecution())
		return 1
	}

	err = output.Transfer(dest, sender, 0, 0, big.NewInt(0).SetBytes(valueBytes), data, vm.DirectCall)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

type indirectContractCallArguments struct {
	dest      []byte
	value     *big.Int
	function  []byte
	args      [][]byte
	actualLen int32
}

func extractIndirectContractCallArgumentsWithValue(
	host vmhost.VMHost,
	destOffset int32,
	valueOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) (*indirectContractCallArguments, error) {
	return extractIndirectContractCallArguments(
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

func extractIndirectContractCallArgumentsWithoutValue(
	host vmhost.VMHost,
	destOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) (*indirectContractCallArguments, error) {
	return extractIndirectContractCallArguments(
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

func extractIndirectContractCallArguments(
	host vmhost.VMHost,
	destOffset int32,
	valueOffset int32,
	hasValueOffset bool,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) (*indirectContractCallArguments, error) {
	runtime := host.Runtime()
	metering := host.Metering()

	dest, err := runtime.MemLoad(destOffset, vmhost.AddressLen)
	if err != nil {
		return nil, err
	}

	var value *big.Int

	if hasValueOffset {
		valueBytes, err := runtime.MemLoad(valueOffset, vmhost.BalanceLen)
		if err != nil {
			return nil, err
		}
		value = big.NewInt(0).SetBytes(valueBytes)
	}

	function, err := runtime.MemLoad(functionOffset, functionLength)
	if err != nil {
		return nil, err
	}

	args, actualLen, err := getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if err != nil {
		return nil, err
	}

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	metering.UseAndTraceGas(gasToUse)

	return &indirectContractCallArguments{
		dest:      dest,
		value:     value,
		function:  function,
		args:      args,
		actualLen: actualLen,
	}, nil
}

//export v1_4_transferValueExecute
func v1_4_transferValueExecute(
	context unsafe.Pointer,
	destOffset int32,
	valueOffset int32,
	gasLimit int64,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := vmhost.GetVMHost(context)
	return TransferValueExecuteWithHost(
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
func TransferValueExecuteWithHost(
	host vmhost.VMHost,
	destOffset int32,
	valueOffset int32,
	gasLimit int64,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(transferValueExecuteName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.TransferValue
	metering.UseAndTraceGas(gasToUse)

	callArgs, err := extractIndirectContractCallArgumentsWithValue(
		host, destOffset, valueOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
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
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.TransferValue
	metering.UseAndTraceGas(gasToUse)

	sender := runtime.GetContextAddress()

	var err error
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
		if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
			return 1
		}
	}

	if contractCallInput != nil {
		if host.IsBuiltinFunctionName(contractCallInput.Function) {
			vmhost.WithFaultAndHostIfFailAlwaysActive(vmhost.ErrNilESDTData, host, runtime.BaseOpsErrorShouldFailExecution())
			return 1
		}
	}

	if host.AreInSameShard(sender, dest) && contractCallInput != nil && host.Blockchain().IsSmartContract(dest) {
		logEEI.Trace("eGLD pre-transfer execution begin")
		_, _, err = host.ExecuteOnDestContext(contractCallInput)
		if err != nil {
			logEEI.Trace("eGLD pre-transfer execution failed", "error", err)
			vmhost.WithFaultAndHostIfFailAlwaysActive(err, host, runtime.BaseOpsErrorShouldFailExecution())
			return 1
		}

		return 0
	}

	data := ""
	if contractCallInput != nil {
		data = makeCrossShardCallFromInput(contractCallInput.Function, contractCallInput.Arguments)
	}

	metering.UseAndTraceGas(uint64(gasLimit))
	err = output.Transfer(dest, sender, uint64(gasLimit), 0, value, []byte(data), vm.DirectCall)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

func makeCrossShardCallFromInput(function string, arguments [][]byte) string {
	txData := function
	for _, arg := range arguments {
		txData += "@" + hex.EncodeToString(arg)
	}

	return txData
}

//export v1_4_transferESDTExecute
func v1_4_transferESDTExecute(
	context unsafe.Pointer,
	destOffset int32,
	tokenIDOffset int32,
	tokenIDLen int32,
	valueOffset int32,
	gasLimit int64,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {

	return v1_4_transferESDTNFTExecute(context, destOffset, tokenIDOffset, tokenIDLen, valueOffset, 0,
		gasLimit, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
}

//export v1_4_transferESDTNFTExecute
func v1_4_transferESDTNFTExecute(
	context unsafe.Pointer,
	destOffset int32,
	tokenIDOffset int32,
	tokenIDLen int32,
	valueOffset int32,
	nonce int64,
	gasLimit int64,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := vmhost.GetVMHost(context)
	metering := host.Metering()
	metering.StartGasTracing(transferESDTNFTExecuteName)
	return TransferESDTNFTExecuteWithHost(
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

//export v1_4_multiTransferESDTNFTExecute
func v1_4_multiTransferESDTNFTExecute(
	context unsafe.Pointer,
	destOffset int32,
	numTokenTransfers int32,
	tokenTransfersArgsLengthOffset int32,
	tokenTransferDataOffset int32,
	gasLimit int64,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := vmhost.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(multiTransferESDTNFTExecuteName)

	if numTokenTransfers == 0 {
		_ = vmhost.WithFaultAndHost(host, vmhost.ErrFailedTransfer, runtime.BaseOpsErrorShouldFailExecution())
		return 1
	}

	callArgs, err := extractIndirectContractCallArgumentsWithoutValue(
		host, destOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(callArgs.actualLen))
	metering.UseAndTraceGas(gasToUse)

	transferArgs, actualLen, err := getArgumentsFromMemory(
		host,
		numTokenTransfers*parsers.ArgsPerTransfer,
		tokenTransfersArgsLengthOffset,
		tokenTransferDataOffset,
	)

	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	metering.UseAndTraceGas(gasToUse)

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
func TransferESDTNFTExecuteWithHost(
	host vmhost.VMHost,
	destOffset int32,
	tokenIDOffset int32,
	tokenIDLen int32,
	valueOffset int32,
	nonce int64,
	gasLimit int64,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()

	tokenIdentifier, executeErr := runtime.MemLoad(tokenIDOffset, tokenIDLen)
	if vmhost.WithFaultAndHost(host, executeErr, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	callArgs, err := extractIndirectContractCallArgumentsWithValue(
		host, destOffset, valueOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(callArgs.actualLen))
	metering.UseAndTraceGas(gasToUse)

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
) int32 {
	var executeErr error

	runtime := host.Runtime()
	metering := host.Metering()

	output := host.Output()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.TransferValue * uint64(len(transfers))
	metering.UseAndTraceGas(gasToUse)

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
		if vmhost.WithFaultAndHost(host, executeErr, runtime.SyncExecAPIErrorShouldFailExecution()) {
			return 1
		}

		contractCallInput.ESDTTransfers = transfers
	}

	snapshotBeforeTransfer := host.Blockchain().GetSnapshot()

	gasLimitForExec, executeErr := output.TransferESDT(dest, sender, transfers, contractCallInput)
	if vmhost.WithFaultAndHost(host, executeErr, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	if host.AreInSameShard(sender, dest) && contractCallInput != nil && host.Blockchain().IsSmartContract(dest) {
		contractCallInput.GasProvided = gasLimitForExec
		logEEI.Trace("ESDT post-transfer execution begin")
		_, _, executeErr = host.ExecuteOnDestContext(contractCallInput)
		if executeErr != nil {
			logEEI.Trace("ESDT post-transfer execution failed", "error", executeErr)
			host.Blockchain().RevertToSnapshot(snapshotBeforeTransfer)
			vmhost.WithFaultAndHostIfFailAlwaysActive(executeErr, host, runtime.BaseOpsErrorShouldFailExecution())
			return 1
		}

		return 0
	}

	return 0
}

//export v1_4_createAsyncCall
func v1_4_createAsyncCall(context unsafe.Pointer,
	asyncContextIdentifier int32,
	identifierLength int32,
	destOffset int32,
	valueOffset int32,
	dataOffset int32,
	length int32,
	successOffset int32,
	successLength int32,
	errorOffset int32,
	errorLength int32,
	gas int64,
) {
	host := vmhost.GetVMHost(context)
	runtime := host.Runtime()

	// TODO consume gas

	acIdentifier, err := runtime.MemLoad(asyncContextIdentifier, identifierLength)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	calledSCAddress, err := runtime.MemLoad(destOffset, vmhost.AddressLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	value, err := runtime.MemLoad(valueOffset, vmhost.BalanceLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	data, err := runtime.MemLoad(dataOffset, length)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	successFunc, err := runtime.MemLoad(successOffset, successLength)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	errorFunc, err := runtime.MemLoad(errorOffset, errorLength)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	err = runtime.AddAsyncContextCall(acIdentifier, &vmhost.AsyncGeneratedCall{
		Destination:     calledSCAddress,
		Data:            data,
		ValueBytes:      value,
		SuccessCallback: string(successFunc),
		ErrorCallback:   string(errorFunc),
		ProvidedGas:     uint64(gas),
	})
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}
}

//export v1_4_setAsyncContextCallback
func v1_4_setAsyncContextCallback(context unsafe.Pointer,
	asyncContextIdentifier int32,
	identifierLength int32,
	callback int32,
	callbackLength int32,
) int32 {
	host := vmhost.GetVMHost(context)
	runtime := host.Runtime()

	// TODO consume gas

	acIdentifier, err := runtime.MemLoad(asyncContextIdentifier, identifierLength)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	asyncContext, err := runtime.GetAsyncContext(acIdentifier)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	callbackFunc, err := runtime.MemLoad(callback, callbackLength)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	asyncContext.Callback = string(callbackFunc)

	return 0
}

//export v1_4_upgradeContract
func v1_4_upgradeContract(
	context unsafe.Pointer,
	destOffset int32,
	gasLimit int64,
	valueOffset int32,
	codeOffset int32,
	codeMetadataOffset int32,
	length int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) {
	host := vmhost.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(upgradeContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	value, err := runtime.MemLoad(valueOffset, vmhost.BalanceLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	code, err := runtime.MemLoad(codeOffset, length)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	codeMetadata, err := runtime.MemLoad(codeMetadataOffset, vmhost.CodeMetadataLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	data, actualLen, err := getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	metering.UseAndTraceGas(gasToUse)

	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	calledSCAddress, err := runtime.MemLoad(destOffset, vmhost.AddressLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	gasSchedule := metering.GasSchedule()
	gasToUse = math.MulUint64(gasSchedule.BaseOperationCost.DataCopyPerByte, uint64(length))
	metering.UseAndTraceGas(gasToUse)

	upgradeContract(host, calledSCAddress, code, codeMetadata, value, data, gasLimit)
}

//export v1_4_upgradeFromSourceContract
func v1_4_upgradeFromSourceContract(
	context unsafe.Pointer,
	destOffset int32,
	gasLimit int64,
	valueOffset int32,
	sourceContractAddressOffset int32,
	codeMetadataOffset int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) {
	host := vmhost.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(upgradeFromSourceContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	value, err := runtime.MemLoad(valueOffset, vmhost.BalanceLen)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	sourceContractAddress, err := runtime.MemLoad(sourceContractAddressOffset, vmhost.AddressLen)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	codeMetadata, err := runtime.MemLoad(codeMetadataOffset, vmhost.CodeMetadataLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	data, actualLen, err := getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	metering.UseAndTraceGas(gasToUse)

	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	calledSCAddress, err := runtime.MemLoad(destOffset, vmhost.AddressLen)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	UpgradeFromSourceContractWithTypedArgs(
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
) {
	runtime := host.Runtime()
	blockchain := host.Blockchain()

	code, err := blockchain.GetCode(sourceContractAddress)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	upgradeContract(host, destContractAddress, code, codeMetadata, value, data, gasLimit)
}

func upgradeContract(
	host vmhost.VMHost,
	destContractAddress []byte,
	code []byte,
	codeMetadata []byte,
	value []byte,
	data [][]byte,
	gasLimit int64,
) {
	runtime := host.Runtime()
	metering := host.Metering()
	gasSchedule := metering.GasSchedule()
	minAsyncCallCost := math.AddUint64(
		math.MulUint64(2, gasSchedule.BaseOpsAPICost.AsyncCallStep),
		gasSchedule.BaseOpsAPICost.AsyncCallbackGasLock)
	if uint64(gasLimit) < minAsyncCallCost {
		runtime.SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
		return
	}

	// Set up the async call as if it is not known whether the called SC
	// is in the same shard with the caller or not. This will be later resolved
	// by runtime.ExecuteAsyncCall().
	callData := vmhost.UpgradeFunctionName + "@" + hex.EncodeToString(code) + "@" + hex.EncodeToString(codeMetadata)
	for _, arg := range data {
		callData += "@" + hex.EncodeToString(arg)
	}

	err := runtime.ExecuteAsyncCall(
		destContractAddress,
		[]byte(callData),
		value,
	)
	logEEI.Trace("upgradeContract", "error", err)

	storage := host.Storage()
	if storage.IsUseDifferentGasCostFlagSet() {
		if errors.Is(err, vmhost.ErrNotEnoughGas) {
			runtime.SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
			return
		}
		if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
			return
		}
	}
}

//export v1_4_asyncCall
func v1_4_asyncCall(context unsafe.Pointer, destOffset int32, valueOffset int32, dataOffset int32, length int32) {
	host := vmhost.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(asyncCallName)

	gasSchedule := metering.GasSchedule()
	gasToUse := gasSchedule.BaseOpsAPICost.AsyncCallStep
	metering.UseAndTraceGas(gasToUse)

	calledSCAddress, err := runtime.MemLoad(destOffset, vmhost.AddressLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	value, err := runtime.MemLoad(valueOffset, vmhost.BalanceLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	gasToUse = math.MulUint64(gasSchedule.BaseOperationCost.DataCopyPerByte, uint64(length))
	metering.UseAndTraceGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	err = runtime.ExecuteAsyncCall(calledSCAddress, data, value)
	if errors.Is(err, vmhost.ErrNotEnoughGas) {
		runtime.SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
		return
	}
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}
}

//export v1_4_getArgumentLength
func v1_4_getArgumentLength(context unsafe.Pointer, id int32) int32 {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetArgument
	metering.UseGasAndAddTracedGas(getArgumentLengthName, gasToUse)

	args := runtime.Arguments()
	if id < 0 || int32(len(args)) <= id {
		vmhost.WithFaultIfFailAlwaysActive(vmhost.ErrInvalidArgument, context, runtime.BaseOpsErrorShouldFailExecution())
		return -1
	}

	return int32(len(args[id]))
}

//export v1_4_getArgument
func v1_4_getArgument(context unsafe.Pointer, id int32, argOffset int32) int32 {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetArgument
	metering.UseGasAndAddTracedGas(getArgumentName, gasToUse)

	args := runtime.Arguments()
	if id < 0 || int32(len(args)) <= id {
		vmhost.WithFaultIfFailAlwaysActive(vmhost.ErrInvalidArgument, context, runtime.BaseOpsErrorShouldFailExecution())
		return -1
	}

	err := runtime.MemStore(argOffset, args[id])
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(args[id]))
}

//export v1_4_getFunction
func v1_4_getFunction(context unsafe.Pointer, functionOffset int32) int32 {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetFunction
	metering.UseGasAndAddTracedGas(getFunctionName, gasToUse)

	function := runtime.Function()
	err := runtime.MemStore(functionOffset, []byte(function))
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(function))
}

//export v1_4_getNumArguments
func v1_4_getNumArguments(context unsafe.Pointer) int32 {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetNumArguments
	metering.UseGasAndAddTracedGas(getNumArgumentsName, gasToUse)

	args := runtime.Arguments()
	return int32(len(args))
}

//export v1_4_storageStore
func v1_4_storageStore(context unsafe.Pointer, keyOffset int32, keyLength int32, dataOffset int32, dataLength int32) int32 {
	host := vmhost.GetVMHost(context)
	return StorageStoreWithHost(
		host,
		keyOffset,
		keyLength,
		dataOffset,
		dataLength,
	)
}

// StorageStoreWithHost - storageStore with host instead of pointer context
func StorageStoreWithHost(host vmhost.VMHost, keyOffset int32, keyLength int32, dataOffset int32, dataLength int32) int32 {
	runtime := host.Runtime()

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return StorageStoreWithTypedArgs(host, key, data)
}

// StorageStoreWithTypedArgs - storageStore with args already read from memory
func StorageStoreWithTypedArgs(host vmhost.VMHost, key []byte, data []byte) int32 {
	runtime := host.Runtime()
	storage := host.Storage()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.StorageStore
	metering.UseGasAndAddTracedGas(storageStoreName, gasToUse)

	storageStatus, err := storage.SetStorage(key, data)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(storageStatus)
}

//export v1_4_storageLoadLength
func v1_4_storageLoadLength(context unsafe.Pointer, keyOffset int32, keyLength int32) int32 {
	runtime := vmhost.GetRuntimeContext(context)
	storage := vmhost.GetStorageContext(context)
	metering := vmhost.GetMeteringContext(context)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	data, trieDepth, usedCache, err := storage.GetStorageUnmetered(key)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	err = storage.UseGasForStorageLoad(storageLoadLengthName, int64(trieDepth), metering.GasSchedule().BaseOpsAPICost.StorageLoad, usedCache)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(data))
}

//export v1_4_storageLoadFromAddress
func v1_4_storageLoadFromAddress(context unsafe.Pointer, addressOffset int32, keyOffset int32, keyLength int32, dataOffset int32) int32 {
	host := vmhost.GetVMHost(context)
	return StorageLoadFromAddressWithHost(
		host,
		addressOffset,
		keyOffset,
		keyLength,
		dataOffset,
	)
}

// StorageLoadFromAddressWithHost - storageLoadFromAddress with host instead of pointer context
func StorageLoadFromAddressWithHost(host vmhost.VMHost, addressOffset int32, keyOffset int32, keyLength int32, dataOffset int32) int32 {
	runtime := host.Runtime()

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	address, err := runtime.MemLoad(addressOffset, vmhost.AddressLen)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	data, err := StorageLoadFromAddressWithTypedArgs(host, address, key)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	err = runtime.MemStore(dataOffset, data)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(data))
}

// StorageLoadFromAddressWithTypedArgs - storageLoadFromAddress with args already read from memory
func StorageLoadFromAddressWithTypedArgs(host vmhost.VMHost, address []byte, key []byte) ([]byte, error) {
	storage := host.Storage()
	metering := host.Metering()
	data, trieDepth, usedCache, err := storage.GetStorageFromAddress(address, key)
	if err != nil {
		return nil, err
	}

	err = storage.UseGasForStorageLoad(storageLoadFromAddressName, int64(trieDepth), metering.GasSchedule().BaseOpsAPICost.StorageLoad, usedCache)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

//export v1_4_storageLoad
func v1_4_storageLoad(context unsafe.Pointer, keyOffset int32, keyLength int32, dataOffset int32) int32 {
	host := vmhost.GetVMHost(context)
	return StorageLoadWithHost(
		host,
		keyOffset,
		keyLength,
		dataOffset,
	)
}

// StorageLoadWithHost - storageLoad with host instead of pointer context
func StorageLoadWithHost(host vmhost.VMHost, keyOffset int32, keyLength int32, dataOffset int32) int32 {
	runtime := host.Runtime()

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	data, err := StorageLoadWithWithTypedArgs(host, key)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	err = runtime.MemStore(dataOffset, data)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(data))
}

// StorageLoadWithWithTypedArgs - storageLoad with args already read from memory
func StorageLoadWithWithTypedArgs(host vmhost.VMHost, key []byte) ([]byte, error) {
	storage := host.Storage()
	metering := host.Metering()
	data, trieDepth, usedCache, err := storage.GetStorage(key)
	if err != nil {
		return nil, err
	}

	err = storage.UseGasForStorageLoad(storageLoadName, int64(trieDepth), metering.GasSchedule().BaseOpsAPICost.StorageLoad, usedCache)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

//export v1_4_setStorageLock
func v1_4_setStorageLock(context unsafe.Pointer, keyOffset int32, keyLength int32, lockTimestamp int64) int32 {
	host := vmhost.GetVMHost(context)
	return SetStorageLockWithHost(
		host,
		keyOffset,
		keyLength,
		lockTimestamp,
	)
}

// SetStorageLockWithHost - setStorageLock with host instead of pointer context
func SetStorageLockWithHost(host vmhost.VMHost, keyOffset int32, keyLength int32, lockTimestamp int64) int32 {
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Int64StorageStore
	metering.UseGasAndAddTracedGas(setStorageLockName, gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return SetStorageLockWithTypedArgs(host, key, lockTimestamp)
}

// SetStorageLockWithTypedArgs - setStorageLock with args already read from memory
func SetStorageLockWithTypedArgs(host vmhost.VMHost, key []byte, lockTimestamp int64) int32 {
	runtime := host.Runtime()
	storage := host.Storage()
	timeLockKey := vmhost.CustomStorageKey(vmhost.TimeLockKeyPrefix, key)
	bigTimestamp := big.NewInt(0).SetInt64(lockTimestamp)
	storageStatus, err := storage.SetProtectedStorage(timeLockKey, bigTimestamp.Bytes())
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}
	return int32(storageStatus)
}

//export v1_4_getStorageLock
func v1_4_getStorageLock(context unsafe.Pointer, keyOffset int32, keyLength int32) int64 {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)
	storage := vmhost.GetStorageContext(context)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	timeLockKey := vmhost.CustomStorageKey(vmhost.TimeLockKeyPrefix, key)
	data, trieDepth, usedCache, err := storage.GetStorage(timeLockKey)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 0
	}

	err = storage.UseGasForStorageLoad(getStorageLockName, int64(trieDepth), metering.GasSchedule().BaseOpsAPICost.StorageLoad, usedCache)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	timeLock := big.NewInt(0).SetBytes(data).Int64()

	// TODO if timelock <= currentTimeStamp { fail somehow }

	return timeLock
}

//export v1_4_isStorageLocked
func v1_4_isStorageLocked(context unsafe.Pointer, keyOffset int32, keyLength int32) int32 {

	timeLock := v1_4_getStorageLock(context, keyOffset, keyLength)
	if timeLock < 0 {
		return -1
	}

	currentTimestamp := v1_4_getBlockTimestamp(context)
	if timeLock <= currentTimestamp {
		return 0
	}

	return 1
}

//export v1_4_clearStorageLock
func v1_4_clearStorageLock(context unsafe.Pointer, keyOffset int32, keyLength int32) int32 {
	return v1_4_setStorageLock(context, keyOffset, keyLength, 0)
}

//export v1_4_getCaller
func v1_4_getCaller(context unsafe.Pointer, resultOffset int32) {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCaller
	metering.UseGasAndAddTracedGas(getCallerName, gasToUse)

	caller := runtime.GetVMInput().CallerAddr

	err := runtime.MemStore(resultOffset, caller)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}
}

//export v1_4_checkNoPayment
func v1_4_checkNoPayment(context unsafe.Pointer) {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	metering.UseGasAndAddTracedGas(checkNoPaymentName, gasToUse)

	vmInput := runtime.GetVMInput()
	if vmInput.CallValue.Sign() > 0 {
		_ = vmhost.WithFault(vmhost.ErrNonPayableFunctionEgld, context, runtime.BaseOpsErrorShouldFailExecution())
		return
	}
	if len(vmInput.ESDTTransfers) > 0 {
		_ = vmhost.WithFault(vmhost.ErrNonPayableFunctionEsdt, context, runtime.BaseOpsErrorShouldFailExecution())
		return
	}
}

//export v1_4_callValue
func v1_4_callValue(context unsafe.Pointer, resultOffset int32) int32 {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	metering.UseGasAndAddTracedGas(callValueName, gasToUse)

	value := runtime.GetVMInput().CallValue.Bytes()
	value = vmhost.PadBytesLeft(value, vmhost.BalanceLen)

	err := runtime.MemStore(resultOffset, value)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(value))
}

//export v1_4_getESDTValue
func v1_4_getESDTValue(context unsafe.Pointer, resultOffset int32) int32 {
	isFail := failIfMoreThanOneESDTTransfer(context)
	if isFail {
		return -1
	}
	return v1_4_getESDTValueByIndex(context, resultOffset, 0)
}

//export v1_4_getESDTValueByIndex
func v1_4_getESDTValueByIndex(context unsafe.Pointer, resultOffset int32, index int32) int32 {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	metering.UseGasAndAddTracedGas(getESDTValueByIndexName, gasToUse)

	var value []byte

	esdtTransfer := getESDTTransferFromInputFailIfWrongIndex(vmhost.GetVMHost(context), index)
	if esdtTransfer != nil && esdtTransfer.ESDTValue.Cmp(vmhost.Zero) > 0 {
		value = esdtTransfer.ESDTValue.Bytes()
		value = vmhost.PadBytesLeft(value, vmhost.BalanceLen)
	}

	err := runtime.MemStore(resultOffset, value)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(value))
}

//export v1_4_getESDTTokenName
func v1_4_getESDTTokenName(context unsafe.Pointer, resultOffset int32) int32 {
	isFail := failIfMoreThanOneESDTTransfer(context)
	if isFail {
		return -1
	}
	return v1_4_getESDTTokenNameByIndex(context, resultOffset, 0)
}

//export v1_4_getESDTTokenNameByIndex
func v1_4_getESDTTokenNameByIndex(context unsafe.Pointer, resultOffset int32, index int32) int32 {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	metering.UseGasAndAddTracedGas(getESDTTokenNameByIndexName, gasToUse)

	esdtTransfer := getESDTTransferFromInputFailIfWrongIndex(vmhost.GetVMHost(context), index)
	var tokenName []byte
	if esdtTransfer != nil {
		tokenName = esdtTransfer.ESDTTokenName
	}

	err := runtime.MemStore(resultOffset, tokenName)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(tokenName))
}

//export v1_4_getESDTTokenNonce
func v1_4_getESDTTokenNonce(context unsafe.Pointer) int64 {
	isFail := failIfMoreThanOneESDTTransfer(context)
	if isFail {
		return -1
	}
	return v1_4_getESDTTokenNonceByIndex(context, 0)
}

//export v1_4_getESDTTokenNonceByIndex
func v1_4_getESDTTokenNonceByIndex(context unsafe.Pointer, index int32) int64 {
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	metering.UseGasAndAddTracedGas(getESDTTokenNonceByIndexName, gasToUse)

	esdtTransfer := getESDTTransferFromInputFailIfWrongIndex(vmhost.GetVMHost(context), index)
	nonce := uint64(0)
	if esdtTransfer != nil {
		nonce = esdtTransfer.ESDTTokenNonce
	}
	return int64(nonce)
}

//export v1_4_getCurrentESDTNFTNonce
func v1_4_getCurrentESDTNFTNonce(context unsafe.Pointer, addressOffset int32, tokenIDOffset int32, tokenIDLen int32) int64 {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)
	storage := vmhost.GetStorageContext(context)

	destination, err := runtime.MemLoad(addressOffset, vmhost.AddressLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 0
	}

	tokenID, err := runtime.MemLoad(tokenIDOffset, tokenIDLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 0
	}

	key := []byte(core.ProtectedKeyPrefix + core.ESDTNFTLatestNonceIdentifier + string(tokenID))
	data, trieDepth, _, err := storage.GetStorageFromAddress(destination, key)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 0
	}

	err = storage.UseGasForStorageLoad(getCurrentESDTNFTNonceName, int64(trieDepth), metering.GasSchedule().BaseOpsAPICost.StorageLoad, false)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 0
	}

	nonce := big.NewInt(0).SetBytes(data).Uint64()
	return int64(nonce)
}

//export v1_4_getESDTTokenType
func v1_4_getESDTTokenType(context unsafe.Pointer) int32 {
	isFail := failIfMoreThanOneESDTTransfer(context)
	if isFail {
		return -1
	}
	return v1_4_getESDTTokenTypeByIndex(context, 0)
}

//export v1_4_getESDTTokenTypeByIndex
func v1_4_getESDTTokenTypeByIndex(context unsafe.Pointer, index int32) int32 {
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	metering.UseGasAndAddTracedGas(getESDTTokenTypeByIndexName, gasToUse)

	esdtTransfer := getESDTTransferFromInputFailIfWrongIndex(vmhost.GetVMHost(context), index)
	if esdtTransfer != nil {
		return int32(esdtTransfer.ESDTTokenType)
	}
	return 0
}

//export v1_4_getNumESDTTransfers
func v1_4_getNumESDTTransfers(context unsafe.Pointer) int32 {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	metering.UseGasAndAddTracedGas(getNumESDTTransfersName, gasToUse)

	return int32(len(runtime.GetVMInput().ESDTTransfers))
}

//export v1_4_getCallValueTokenName
func v1_4_getCallValueTokenName(context unsafe.Pointer, callValueOffset int32, tokenNameOffset int32) int32 {
	isFail := failIfMoreThanOneESDTTransfer(context)
	if isFail {
		return -1
	}
	return v1_4_getCallValueTokenNameByIndex(context, callValueOffset, tokenNameOffset, 0)
}

//export v1_4_getCallValueTokenNameByIndex
func v1_4_getCallValueTokenNameByIndex(context unsafe.Pointer, callValueOffset int32, tokenNameOffset int32, index int32) int32 {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	metering.UseGasAndAddTracedGas(getCallValueTokenNameByIndexName, gasToUse)

	callValue := runtime.GetVMInput().CallValue.Bytes()
	tokenName := make([]byte, 0)
	esdtTransfer := getESDTTransferFromInputFailIfWrongIndex(vmhost.GetVMHost(context), index)

	if esdtTransfer != nil {
		tokenName = make([]byte, len(esdtTransfer.ESDTTokenName))
		copy(tokenName, esdtTransfer.ESDTTokenName)
		callValue = esdtTransfer.ESDTValue.Bytes()
	}
	callValue = vmhost.PadBytesLeft(callValue, vmhost.BalanceLen)

	err := runtime.MemStore(tokenNameOffset, tokenName)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	err = runtime.MemStore(callValueOffset, callValue)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(tokenName))
}

//export v1_4_writeLog
func v1_4_writeLog(context unsafe.Pointer, dataPointer int32, dataLength int32, topicPtr int32, numTopics int32) {
	// note: deprecated
	runtime := vmhost.GetRuntimeContext(context)
	output := vmhost.GetOutputContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Log
	gas := math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(numTopics*vmhost.HashLen+dataLength))
	gasToUse = math.AddUint64(gasToUse, gas)
	metering.UseGasAndAddTracedGas(writeLogName, gasToUse)

	log, err := runtime.MemLoad(dataPointer, dataLength)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	topics, err := vmhost.GuardedMakeByteSlice2D(numTopics)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	for i := int32(0); i < numTopics; i++ {
		topics[i], err = runtime.MemLoad(topicPtr+i*vmhost.HashLen, vmhost.HashLen)
		if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
			return
		}
	}

	output.WriteLog(runtime.GetContextAddress(), topics, log)
}

//export v1_4_writeEventLog
func v1_4_writeEventLog(
	context unsafe.Pointer,
	numTopics int32,
	topicLengthsOffset int32,
	topicOffset int32,
	dataOffset int32,
	dataLength int32,
) {

	host := vmhost.GetVMHost(context)
	runtime := vmhost.GetRuntimeContext(context)
	output := vmhost.GetOutputContext(context)
	metering := vmhost.GetMeteringContext(context)

	topics, topicDataTotalLen, err := getArgumentsFromMemory(
		host,
		numTopics,
		topicLengthsOffset,
		topicOffset,
	)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Log
	gasForData := math.MulUint64(
		metering.GasSchedule().BaseOperationCost.DataCopyPerByte,
		uint64(topicDataTotalLen+dataLength))
	gasToUse = math.AddUint64(gasToUse, gasForData)
	metering.UseGasAndAddTracedGas(writeEventLogName, gasToUse)

	output.WriteLog(runtime.GetContextAddress(), topics, data)
}

//export v1_4_getBlockTimestamp
func v1_4_getBlockTimestamp(context unsafe.Pointer) int64 {
	blockchain := vmhost.GetBlockchainContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockTimeStamp
	metering.UseGasAndAddTracedGas(getBlockTimestampName, gasToUse)

	return int64(blockchain.CurrentTimeStamp())
}

//export v1_4_getBlockNonce
func v1_4_getBlockNonce(context unsafe.Pointer) int64 {
	blockchain := vmhost.GetBlockchainContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockNonce
	metering.UseGasAndAddTracedGas(getBlockNonceName, gasToUse)

	return int64(blockchain.CurrentNonce())
}

//export v1_4_getBlockRound
func v1_4_getBlockRound(context unsafe.Pointer) int64 {
	blockchain := vmhost.GetBlockchainContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRound
	metering.UseGasAndAddTracedGas(getBlockRoundName, gasToUse)

	return int64(blockchain.CurrentRound())
}

//export v1_4_getBlockEpoch
func v1_4_getBlockEpoch(context unsafe.Pointer) int64 {
	blockchain := vmhost.GetBlockchainContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockEpoch
	metering.UseGasAndAddTracedGas(getBlockEpochName, gasToUse)

	return int64(blockchain.CurrentEpoch())
}

//export v1_4_getBlockRandomSeed
func v1_4_getBlockRandomSeed(context unsafe.Pointer, pointer int32) {
	runtime := vmhost.GetRuntimeContext(context)
	blockchain := vmhost.GetBlockchainContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRandomSeed
	metering.UseGasAndAddTracedGas(getBlockRandomSeedName, gasToUse)

	randomSeed := blockchain.CurrentRandomSeed()
	err := runtime.MemStore(pointer, randomSeed)
	vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution())
}

//export v1_4_getStateRootHash
func v1_4_getStateRootHash(context unsafe.Pointer, pointer int32) {
	runtime := vmhost.GetRuntimeContext(context)
	blockchain := vmhost.GetBlockchainContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetStateRootHash
	metering.UseGasAndAddTracedGas(getStateRootHashName, gasToUse)

	stateRootHash := blockchain.GetStateRootHash()
	err := runtime.MemStore(pointer, stateRootHash)
	vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution())
}

//export v1_4_getPrevBlockTimestamp
func v1_4_getPrevBlockTimestamp(context unsafe.Pointer) int64 {
	blockchain := vmhost.GetBlockchainContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockTimeStamp
	metering.UseGasAndAddTracedGas(getPrevBlockTimestampName, gasToUse)

	return int64(blockchain.LastTimeStamp())
}

//export v1_4_getPrevBlockNonce
func v1_4_getPrevBlockNonce(context unsafe.Pointer) int64 {
	blockchain := vmhost.GetBlockchainContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockNonce
	metering.UseGasAndAddTracedGas(getPrevBlockNonceName, gasToUse)

	return int64(blockchain.LastNonce())
}

//export v1_4_getPrevBlockRound
func v1_4_getPrevBlockRound(context unsafe.Pointer) int64 {
	blockchain := vmhost.GetBlockchainContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRound
	metering.UseGasAndAddTracedGas(getPrevBlockRoundName, gasToUse)

	return int64(blockchain.LastRound())
}

//export v1_4_getPrevBlockEpoch
func v1_4_getPrevBlockEpoch(context unsafe.Pointer) int64 {
	blockchain := vmhost.GetBlockchainContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockEpoch
	metering.UseGasAndAddTracedGas(getPrevBlockEpochName, gasToUse)

	return int64(blockchain.LastEpoch())
}

//export v1_4_getPrevBlockRandomSeed
func v1_4_getPrevBlockRandomSeed(context unsafe.Pointer, pointer int32) {
	runtime := vmhost.GetRuntimeContext(context)
	blockchain := vmhost.GetBlockchainContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRandomSeed
	metering.UseGasAndAddTracedGas(getPrevBlockRandomSeedName, gasToUse)

	randomSeed := blockchain.LastRandomSeed()
	err := runtime.MemStore(pointer, randomSeed)
	vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution())
}

//export v1_4_returnData
func v1_4_returnData(context unsafe.Pointer, pointer int32, length int32) {
	runtime := vmhost.GetRuntimeContext(context)
	output := vmhost.GetOutputContext(context)
	metering := vmhost.GetMeteringContext(context)
	metering.StartGasTracing(returnDataName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Finish
	gas := math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(length))
	gasToUse = math.AddUint64(gasToUse, gas)
	err := metering.UseGasBounded(gasToUse)

	if err != nil {
		_ = vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	data, err := runtime.MemLoad(pointer, length)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	output.Finish(data)
}

//export v1_4_executeOnSameContext
func v1_4_executeOnSameContext(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	valueOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := vmhost.GetVMHost(context)
	metering := host.Metering()
	metering.StartGasTracing(executeOnSameContextName)

	return ExecuteOnSameContextWithHost(
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
func ExecuteOnSameContextWithHost(
	host vmhost.VMHost,
	gasLimit int64,
	addressOffset int32,
	valueOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	runtime := host.Runtime()

	callArgs, err := extractIndirectContractCallArgumentsWithValue(
		host, addressOffset, valueOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
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
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.ExecuteOnSameContext
	metering.UseAndTraceGas(gasToUse)

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
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	if isBuiltInCall(contractCallInput.Function, host) {
		vmhost.WithFaultAndHostIfFailAlwaysActive(vmhost.ErrInvalidBuiltInFunctionCall, host, runtime.BaseOpsErrorShouldFailExecution())
		return 1
	}

	_, err = host.ExecuteOnSameContext(contractCallInput)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return 0
}

//export v1_4_executeOnDestContext
func v1_4_executeOnDestContext(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	valueOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := vmhost.GetVMHost(context)
	metering := host.Metering()
	metering.StartGasTracing(executeOnDestContextName)

	return ExecuteOnDestContextWithHost(
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
func ExecuteOnDestContextWithHost(
	host vmhost.VMHost,
	gasLimit int64,
	addressOffset int32,
	valueOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	runtime := host.Runtime()

	callArgs, err := extractIndirectContractCallArgumentsWithValue(
		host, addressOffset, valueOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return ExecuteOnDestContextWithTypedArgs(
		host,
		gasLimit,
		callArgs.value,
		callArgs.function,
		callArgs.dest,
		callArgs.args,
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
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.ExecuteOnDestContext
	metering.UseAndTraceGas(gasToUse)

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
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	_, _, err = host.ExecuteOnDestContext(contractCallInput)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export v1_4_executeOnDestContextByCaller
func v1_4_executeOnDestContextByCaller(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	valueOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := vmhost.GetVMHost(context)
	metering := host.Metering()
	metering.StartGasTracing(executeOnDestContextByCallerName)

	return ExecuteOnDestContextByCallerWithHost(
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

// ExecuteOnDestContextByCallerWithHost - executeOnDestContextByCaller with host instead of pointer context
func ExecuteOnDestContextByCallerWithHost(
	host vmhost.VMHost,
	gasLimit int64,
	addressOffset int32,
	valueOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	runtime := host.Runtime()

	callArgs, err := extractIndirectContractCallArgumentsWithValue(
		host, addressOffset, valueOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return ExecuteOnDestContextByCallerWithTypedArgs(
		host,
		gasLimit,
		callArgs.value,
		callArgs.function,
		callArgs.dest,
		callArgs.args,
	)
}

// ExecuteOnDestContextByCallerWithTypedArgs - executeOnDestContextByCaller with args already read from memory
func ExecuteOnDestContextByCallerWithTypedArgs(
	host vmhost.VMHost,
	gasLimit int64,
	value *big.Int,
	function []byte,
	dest []byte,
	args [][]byte,
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.ExecuteOnDestContext
	metering.UseAndTraceGas(gasToUse)

	if host.DisableExecByCaller() {
		_ = vmhost.WithFaultAndHost(host, core.ErrInvalidValue, runtime.BaseOpsErrorShouldFailExecution())
		return -1
	}

	send := runtime.GetVMInput().CallerAddr
	contractCallInput, err := prepareIndirectContractCallInput(
		host,
		send,
		value,
		gasLimit,
		dest,
		function,
		args,
		gasToUse,
		true,
	)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	if isBuiltInCall(contractCallInput.Function, host) {
		if !host.CreateNFTOnExecByCallerEnabled() {
			return 1
		}

		if contractCallInput.Function != core.BuiltInFunctionESDTNFTCreate {
			if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
				return -1
			}
			return -1
		}

		contractCallInput.CallType = vm.ExecOnDestByCaller
		contractCallInput.Arguments = append(contractCallInput.Arguments, runtime.GetContextAddress())
	}

	_, _, err = host.ExecuteOnDestContext(contractCallInput)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return 0
}

//export v1_4_executeReadOnly
func v1_4_executeReadOnly(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := vmhost.GetVMHost(context)
	metering := host.Metering()
	metering.StartGasTracing(executeReadOnlyName)

	return ExecuteReadOnlyWithHost(
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
func ExecuteReadOnlyWithHost(
	host vmhost.VMHost,
	gasLimit int64,
	addressOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	runtime := host.Runtime()

	callArgs, err := extractIndirectContractCallArgumentsWithoutValue(
		host, addressOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
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
) int32 {
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.ExecuteReadOnly
	metering.UseAndTraceGas(gasToUse)

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
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	if isBuiltInCall(contractCallInput.Function, host) {
		vmhost.WithFaultAndHostIfFailAlwaysActive(vmhost.ErrInvalidBuiltInFunctionCall, host, runtime.BaseOpsErrorShouldFailExecution())
		return 1
	}

	wasReadOnly := runtime.ReadOnly()
	runtime.SetReadOnly(true)
	_, _, err = host.ExecuteOnDestContext(contractCallInput)
	runtime.SetReadOnly(false)
	if host.CheckExecuteReadOnly() {
		runtime.SetReadOnly(wasReadOnly)
	}

	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return 0
}

//export v1_4_createContract
func v1_4_createContract(
	context unsafe.Pointer,
	gasLimit int64,
	valueOffset int32,
	codeOffset int32,
	codeMetadataOffset int32,
	length int32,
	resultOffset int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := vmhost.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(createContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	sender := runtime.GetContextAddress()
	value, err := runtime.MemLoad(valueOffset, vmhost.BalanceLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	code, err := runtime.MemLoad(codeOffset, length)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	codeMetadata, err := runtime.MemLoad(codeMetadataOffset, vmhost.CodeMetadataLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	data, actualLen, err := getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	metering.UseAndTraceGas(gasToUse)

	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	valueAsInt := big.NewInt(0).SetBytes(value)
	newAddress, err := createContract(sender, data, valueAsInt, metering, gasLimit, code, codeMetadata, host, runtime)

	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	err = runtime.MemStore(resultOffset, newAddress)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export v1_4_deployFromSourceContract
func v1_4_deployFromSourceContract(
	context unsafe.Pointer,
	gasLimit int64,
	valueOffset int32,
	sourceContractAddressOffset int32,
	codeMetadataOffset int32,
	resultAddressOffset int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := vmhost.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(deployFromSourceContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	value, err := runtime.MemLoad(valueOffset, vmhost.BalanceLen)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	sourceContractAddress, err := runtime.MemLoad(sourceContractAddressOffset, vmhost.AddressLen)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	codeMetadata, err := runtime.MemLoad(codeMetadataOffset, vmhost.CodeMetadataLen)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	data, actualLen, err := getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	metering.UseAndTraceGas(gasToUse)

	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	newAddress, err := DeployFromSourceContractWithTypedArgs(
		host,
		sourceContractAddress,
		codeMetadata,
		big.NewInt(0).SetBytes(value),
		data,
		gasLimit,
	)

	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	err = runtime.MemStore(resultAddressOffset, newAddress)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	return 0
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
	metering := host.Metering()
	sender := runtime.GetContextAddress()

	blockchain := host.Blockchain()
	code, err := blockchain.GetCode(sourceContractAddress)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return nil, err
	}

	return createContract(sender, data, value, metering, gasLimit, code, codeMetadata, host, runtime)
}

func createContract(
	sender []byte,
	data [][]byte,
	value *big.Int,
	metering vmhost.MeteringContext,
	gasLimit int64,
	code []byte,
	codeMetadata []byte,
	host vmhost.VMHost,
	_ vmhost.RuntimeContext,
) ([]byte, error) {
	contractCreate := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   data,
			CallValue:   value,
			GasPrice:    0,
			GasProvided: metering.BoundGasLimit(gasLimit),
		},
		ContractCode:         code,
		ContractCodeMetadata: codeMetadata,
	}

	return host.CreateNewContract(contractCreate)
}

//export v1_4_getNumReturnData
func v1_4_getNumReturnData(context unsafe.Pointer) int32 {
	output := vmhost.GetOutputContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetNumReturnData
	metering.UseGasAndAddTracedGas(getNumReturnDataName, gasToUse)

	returnData := output.ReturnData()
	return int32(len(returnData))
}

//export v1_4_getReturnDataSize
func v1_4_getReturnDataSize(context unsafe.Pointer, resultID int32) int32 {
	runtime := vmhost.GetRuntimeContext(context)
	output := vmhost.GetOutputContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetReturnDataSize
	metering.UseGasAndAddTracedGas(getReturnDataSizeName, gasToUse)

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) || resultID < 0 {
		vmhost.WithFaultAndHostIfFailAlwaysActive(vmhost.ErrInvalidArgument, vmhost.GetVMHost(context), runtime.BaseOpsErrorShouldFailExecution())
		return 0
	}

	return int32(len(returnData[resultID]))
}

//export v1_4_getReturnData
func v1_4_getReturnData(context unsafe.Pointer, resultID int32, dataOffset int32) int32 {
	host := vmhost.GetVMHost(context)

	result := GetReturnDataWithHostAndTypedArgs(host, resultID)
	if result == nil {
		return 0
	}

	runtime := vmhost.GetRuntimeContext(context)
	err := runtime.MemStore(dataOffset, result)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return 0
	}

	return int32(len(result))
}

func GetReturnDataWithHostAndTypedArgs(host vmhost.VMHost, resultID int32) []byte {
	output := host.Output()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetReturnData
	metering.UseGasAndAddTracedGas(getReturnDataName, gasToUse)

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) || resultID < 0 {
		vmhost.WithFaultAndHostIfFailAlwaysActive(vmhost.ErrInvalidArgument, host, host.Runtime().BaseOpsErrorShouldFailExecution())
		return nil
	}

	return returnData[resultID]
}

//export v1_4_cleanReturnData
func v1_4_cleanReturnData(context unsafe.Pointer) {
	host := vmhost.GetVMHost(context)
	CleanReturnDataWithHost(host)
}

// CleanReturnDataWithHost - exposed version of v1_4_deleteFromReturnData for tests
func CleanReturnDataWithHost(host vmhost.VMHost) {
	output := host.Output()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CleanReturnData
	metering.UseGasAndAddTracedGas(cleanReturnDataName, gasToUse)

	output.ClearReturnData()
}

//export v1_4_deleteFromReturnData
func v1_4_deleteFromReturnData(context unsafe.Pointer, resultID int32) {
	host := vmhost.GetVMHost(context)
	DeleteFromReturnDataWithHost(host, resultID)
}

// DeleteFromReturnDataWithHost - exposed version of v1_4_deleteFromReturnData for tests
func DeleteFromReturnDataWithHost(host vmhost.VMHost, resultID int32) {
	output := host.Output()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.DeleteFromReturnData
	metering.UseGasAndAddTracedGas(deleteFromReturnDataName, gasToUse)

	returnData := output.ReturnData()
	if resultID < int32(len(returnData)) {
		output.RemoveReturnData(uint32(resultID))
	}
}

//export v1_4_getOriginalTxHash
func v1_4_getOriginalTxHash(context unsafe.Pointer, dataOffset int32) {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetOriginalTxHash
	metering.UseGasAndAddTracedGas(getOriginalTxHashName, gasToUse)

	err := runtime.MemStore(dataOffset, runtime.GetOriginalTxHash())
	_ = vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution())
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
			CallerAddr:  sender,
			Arguments:   data,
			CallValue:   value,
			GasPrice:    0,
			GasProvided: metering.BoundGasLimit(gasLimit),
		},
		RecipientAddr: destination,
		Function:      string(function),
	}

	return contractCallInput, nil
}

func getArgumentsFromMemory(
	host vmhost.VMHost,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) ([][]byte, int32, error) {
	runtime := host.Runtime()

	if numArguments < 0 {
		return nil, 0, fmt.Errorf("negative numArguments (%d)", numArguments)
	}

	argumentsLengthData, err := runtime.MemLoad(argumentsLengthOffset, numArguments*4)
	if err != nil {
		return nil, 0, err
	}

	argumentLengths := createInt32Array(argumentsLengthData, numArguments)
	data, err := runtime.MemLoadMultiple(dataOffset, argumentLengths)
	if err != nil {
		return nil, 0, err
	}

	totalArgumentBytes := int32(0)
	for _, length := range argumentLengths {
		totalArgumentBytes += length
	}

	return data, totalArgumentBytes, nil
}

func createInt32Array(rawData []byte, numIntegers int32) []int32 {
	integers := make([]int32, numIntegers)
	index := 0
	for cursor := 0; cursor < len(rawData); cursor += 4 {
		rawInt := rawData[cursor : cursor+4]
		actualInt := binary.LittleEndian.Uint32(rawInt)
		integers[index] = int32(actualInt)
		index++
	}
	return integers
}

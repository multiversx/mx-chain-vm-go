package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern void			v1_2_getSCAddress(void *context, int32_t resultOffset);
// extern void			v1_2_getOwnerAddress(void *context, int32_t resultOffset);
// extern int32_t 	v1_2_getShardOfAddress(void *context, int32_t addressOffset);
// extern int32_t 	v1_2_isSmartContract(void *context, int32_t addressOffset);
// extern void			v1_2_getExternalBalance(void *context, int32_t addressOffset, int32_t resultOffset);
// extern int32_t		v1_2_blockHash(void *context, long long nonce, int32_t resultOffset);
// extern int32_t 	v1_2_transferValue(void *context, int32_t dstOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
// extern int32_t 	v1_2_transferESDT(void *context, int32_t dstOffset, int32_t tokenIDOffset, int32_t tokenIdLen, int32_t valueOffset, long long gasLimit, int32_t dataOffset, int32_t length);
// extern int32_t 	v1_2_transferESDTExecute(void *context, int32_t dstOffset, int32_t tokenIDOffset, int32_t tokenIdLen, int32_t valueOffset, long long gasLimit, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t 	v1_2_transferESDTNFTExecute(void *context, int32_t dstOffset, int32_t tokenIDOffset, int32_t tokenIdLen, int32_t valueOffset, long long nonce, long long gasLimit, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t 	v1_2_transferValueExecute(void *context, int32_t dstOffset, int32_t valueOffset, long long gasLimit, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t 	v1_2_getArgumentLength(void *context, int32_t id);
// extern int32_t 	v1_2_getArgument(void *context, int32_t id, int32_t argOffset);
// extern int32_t 	v1_2_getFunction(void *context, int32_t functionOffset);
// extern int32_t 	v1_2_getNumArguments(void *context);
// extern int32_t 	v1_2_storageStore(void *context, int32_t keyOffset, int32_t keyLength , int32_t dataOffset, int32_t dataLength);
// extern int32_t 	v1_2_storageLoadLength(void *context, int32_t keyOffset, int32_t keyLength );
// extern int32_t 	v1_2_storageLoad(void *context, int32_t keyOffset, int32_t keyLength , int32_t dataOffset);
// extern int32_t 	v1_2_storageLoadFromAddress(void *context, int32_t addressOffset, int32_t keyOffset, int32_t keyLength , int32_t dataOffset);
// extern void			v1_2_getCaller(void *context, int32_t resultOffset);
// extern void		 	v1_2_checkNoPayment(void *context);
// extern int32_t		v1_2_callValue(void *context, int32_t resultOffset);
// extern int32_t		v1_2_getESDTValue(void *context, int32_t resultOffset);
// extern int32_t		v1_2_getESDTTokenName(void *context, int32_t resultOffset);
// extern long long v1_2_getESDTTokenNonce(void *context);
// extern int32_t		v1_2_getESDTTokenType(void *context);
// extern long long v1_2_getCurrentESDTNFTNonce(void *context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen);
// extern int32_t		v1_2_getCallValueTokenName(void *context, int32_t callValueOffset, int32_t tokenNameOffset);
// extern void			v1_2_writeLog(void *context, int32_t pointer, int32_t length, int32_t topicPtr, int32_t numTopics);
// extern void 			v1_2_writeEventLog(void *context, int32_t numTopics, int32_t topicLengthsOffset, int32_t topicOffset, int32_t dataOffset, int32_t dataLength);
// extern void 			v1_2_returnData(void* context, int32_t dataOffset, int32_t length);
// extern void 			v1_2_signalError(void* context, int32_t messageOffset, int32_t messageLength);
// extern long long v1_2_getGasLeft(void *context);
// extern int32_t		v1_2_getESDTBalance(void *context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen, long long nonce, int32_t resultOffset);
// extern int32_t 	v1_2_getESDTNFTNameLength(void *context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen, long long nonce);
// extern int32_t 	v1_2_getESDTNFTAttributeLength(void *context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen, long long nonce);
// extern int32_t 	v1_2_getESDTNFTURILength(void *context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen, long long nonce);
// extern int32_t 	v1_2_getESDTTokenData(void *context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen, long long nonce, int32_t valueOffset, int32_t propertiesOffset, int32_t hashOffset, int32_t nameOffset, int32_t attributesOffset, int32_t creatorOffset, int32_t royaltiesOffset, int32_t urisOffset);
//
// extern int32_t		v1_2_executeOnDestContext(void *context, long long gas, int32_t addressOffset, int32_t valueOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t 	v1_2_executeOnDestContextByCaller(void *context, long long gas, int32_t addressOffset, int32_t valueOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t 	v1_2_executeOnSameContext(void *context, long long gas, int32_t addressOffset, int32_t valueOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t 	v1_2_delegateExecution(void *context, long long gas, int32_t addressOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t 	v1_2_executeReadOnly(void *context, long long gas, int32_t addressOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t 	v1_2_createContract(void *context, long long gas, int32_t valueOffset, int32_t codeOffset, int32_t codeMetadataOffset, int32_t length, int32_t resultOffset, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern void			v1_2_upgradeContract(void *context, int32_t dstOffset, long long gas, int32_t valueOffset, int32_t codeOffset, int32_t codeMetadataOffset, int32_t length, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern void 			v1_2_asyncCall(void *context, int32_t dstOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
// extern void 			v1_2_createAsyncCall(void *context, int32_t identifierOffset, int32_t identifierLength, int32_t dstOffset, int32_t valueOffset, int32_t dataOffset, int32_t length, int32_t successCallback, int32_t successLength, int32_t errorCallback, int32_t errorLength, long long gas);
// extern int32_t		v1_2_setAsyncContextCallback(void *context, int32_t identifierOffset, int32_t identifierLength, int32_t callback, int32_t callbackLength);
//
// extern int32_t		v1_2_getNumReturnData(void *context);
// extern int32_t 	v1_2_getReturnDataSize(void *context, int32_t resultID);
// extern int32_t 	v1_2_getReturnData(void *context, int32_t resultID, int32_t dataOffset);
//
// extern int32_t		v1_2_setStorageLock(void *context, int32_t keyOffset, int32_t keyLength, long long lockTimestamp);
// extern long long v1_2_getStorageLock(void *context, int32_t keyOffset, int32_t keyLength);
// extern int32_t		v1_2_isStorageLocked(void *context, int32_t keyOffset, int32_t keyLength);
// extern int32_t		v1_2_clearStorageLock(void *context, int32_t keyOffset, int32_t keyLength);
//
// extern long long v1_2_getBlockTimestamp(void *context);
// extern long long v1_2_getBlockNonce(void *context);
// extern long long v1_2_getBlockRound(void *context);
// extern long long v1_2_getBlockEpoch(void *context);
// extern void			v1_2_getBlockRandomSeed(void *context, int32_t resultOffset);
// extern void			v1_2_getStateRootHash(void *context, int32_t resultOffset);
//
// extern long long v1_2_getPrevBlockTimestamp(void *context);
// extern long long v1_2_getPrevBlockNonce(void *context);
// extern long long v1_2_getPrevBlockRound(void *context);
// extern long long v1_2_getPrevBlockEpoch(void *context);
// extern void			v1_2_getPrevBlockRandomSeed(void *context, int32_t resultOffset);
// extern void			v1_2_getOriginalTxHash(void *context, int32_t resultOffset);
import "C"

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_2/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_2/math"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_2/wasmer"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data/esdt"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/elrond-vm-common/parsers"
)

var logEEI = logger.GetOrCreate("arwen/eei")

func getFirstESDTTransferIfExist(vmInput *vmcommon.VMInput) *vmcommon.ESDTTransfer {
	esdtTransfers := vmInput.ESDTTransfers
	if len(esdtTransfers) > 0 {
		return esdtTransfers[0]
	}
	return &vmcommon.ESDTTransfer{
		ESDTValue: big.NewInt(0),
	}
}

// ElrondEIImports creates a new wasmer.Imports populated with the ElrondEI API methods
func ElrondEIImports() (*wasmer.Imports, error) {
	imports := wasmer.NewImports()
	imports = imports.Namespace("env")

	imports, err := imports.Append("getSCAddress", v1_2_getSCAddress, C.v1_2_getSCAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getOwnerAddress", v1_2_getOwnerAddress, C.v1_2_getOwnerAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getShardOfAddress", v1_2_getShardOfAddress, C.v1_2_getShardOfAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("isSmartContract", v1_2_isSmartContract, C.v1_2_isSmartContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getExternalBalance", v1_2_getExternalBalance, C.v1_2_getExternalBalance)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockHash", v1_2_blockHash, C.v1_2_blockHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("transferValue", v1_2_transferValue, C.v1_2_transferValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("transferESDTExecute", v1_2_transferESDTExecute, C.v1_2_transferESDTExecute)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("transferESDTNFTExecute", v1_2_transferESDTNFTExecute, C.v1_2_transferESDTNFTExecute)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("transferESDT", v1_2_transferESDT, C.v1_2_transferESDT)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("transferValueExecute", v1_2_transferValueExecute, C.v1_2_transferValueExecute)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("asyncCall", v1_2_asyncCall, C.v1_2_asyncCall)
	if err != nil {
		return nil, err
	}

	// imports, err = imports.Append("createAsyncCall", createAsyncCall, C.createAsyncCall)
	// if err != nil {
	// 	return nil, err
	// }

	// imports, err = imports.Append("setAsyncContextCallback", setAsyncContextCallback, C.setAsyncContextCallback)
	// if err != nil {
	// 	return nil, err
	// }

	imports, err = imports.Append("getArgumentLength", v1_2_getArgumentLength, C.v1_2_getArgumentLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getArgument", v1_2_getArgument, C.v1_2_getArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getFunction", v1_2_getFunction, C.v1_2_getFunction)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getNumArguments", v1_2_getNumArguments, C.v1_2_getNumArguments)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageStore", v1_2_storageStore, C.v1_2_storageStore)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageLoadLength", v1_2_storageLoadLength, C.v1_2_storageLoadLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageLoad", v1_2_storageLoad, C.v1_2_storageLoad)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageLoadFromAddress", v1_2_storageLoadFromAddress, C.v1_2_storageLoadFromAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getStorageLock", v1_2_getStorageLock, C.v1_2_getStorageLock)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("setStorageLock", v1_2_setStorageLock, C.v1_2_setStorageLock)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("isStorageLocked", v1_2_isStorageLocked, C.v1_2_isStorageLocked)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("clearStorageLock", v1_2_clearStorageLock, C.v1_2_clearStorageLock)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCaller", v1_2_getCaller, C.v1_2_getCaller)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("checkNoPayment", v1_2_checkNoPayment, C.v1_2_checkNoPayment)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCallValue", v1_2_callValue, C.v1_2_callValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getESDTValue", v1_2_getESDTValue, C.v1_2_getESDTValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getESDTTokenName", v1_2_getESDTTokenName, C.v1_2_getESDTTokenName)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getESDTTokenType", v1_2_getESDTTokenType, C.v1_2_getESDTTokenType)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getESDTTokenNonce", v1_2_getESDTTokenNonce, C.v1_2_getESDTTokenNonce)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCurrentESDTNFTNonce", v1_2_getCurrentESDTNFTNonce, C.v1_2_getCurrentESDTNFTNonce)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCallValueTokenName", v1_2_getCallValueTokenName, C.v1_2_getCallValueTokenName)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("writeLog", v1_2_writeLog, C.v1_2_writeLog)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("writeEventLog", v1_2_writeEventLog, C.v1_2_writeEventLog)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("finish", v1_2_returnData, C.v1_2_returnData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("signalError", v1_2_signalError, C.v1_2_signalError)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockTimestamp", v1_2_getBlockTimestamp, C.v1_2_getBlockTimestamp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockNonce", v1_2_getBlockNonce, C.v1_2_getBlockNonce)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockRound", v1_2_getBlockRound, C.v1_2_getBlockRound)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockEpoch", v1_2_getBlockEpoch, C.v1_2_getBlockEpoch)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockRandomSeed", v1_2_getBlockRandomSeed, C.v1_2_getBlockRandomSeed)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getStateRootHash", v1_2_getStateRootHash, C.v1_2_getStateRootHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockTimestamp", v1_2_getPrevBlockTimestamp, C.v1_2_getPrevBlockTimestamp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockNonce", v1_2_getPrevBlockNonce, C.v1_2_getPrevBlockNonce)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockRound", v1_2_getPrevBlockRound, C.v1_2_getPrevBlockRound)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockEpoch", v1_2_getPrevBlockEpoch, C.v1_2_getPrevBlockEpoch)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockRandomSeed", v1_2_getPrevBlockRandomSeed, C.v1_2_getPrevBlockRandomSeed)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getOriginalTxHash", v1_2_getOriginalTxHash, C.v1_2_getOriginalTxHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getGasLeft", v1_2_getGasLeft, C.v1_2_getGasLeft)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("executeOnDestContext", v1_2_executeOnDestContext, C.v1_2_executeOnDestContext)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("executeOnDestContextByCaller", v1_2_executeOnDestContextByCaller, C.v1_2_executeOnDestContextByCaller)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("executeOnSameContext", v1_2_executeOnSameContext, C.v1_2_executeOnSameContext)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("delegateExecution", v1_2_delegateExecution, C.v1_2_delegateExecution)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("createContract", v1_2_createContract, C.v1_2_createContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("upgradeContract", v1_2_upgradeContract, C.v1_2_upgradeContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("executeReadOnly", v1_2_executeReadOnly, C.v1_2_executeReadOnly)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getNumReturnData", v1_2_getNumReturnData, C.v1_2_getNumReturnData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getReturnDataSize", v1_2_getReturnDataSize, C.v1_2_getReturnDataSize)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getReturnData", v1_2_getReturnData, C.v1_2_getReturnData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getESDTBalance", v1_2_getESDTBalance, C.v1_2_getESDTBalance)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getESDTTokenData", v1_2_getESDTTokenData, C.v1_2_getESDTTokenData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getESDTNFTNameLength", v1_2_getESDTNFTNameLength, C.v1_2_getESDTNFTNameLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getESDTNFTAttributeLength", v1_2_getESDTNFTAttributeLength, C.v1_2_getESDTNFTAttributeLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getESDTNFTURILength", v1_2_getESDTNFTURILength, C.v1_2_getESDTNFTURILength)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export v1_2_getGasLeft
func v1_2_getGasLeft(context unsafe.Pointer) int64 {
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetGasLeft
	metering.UseGas(gasToUse)

	return int64(metering.GasLeft())
}

//export v1_2_getSCAddress
func v1_2_getSCAddress(context unsafe.Pointer, resultOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetSCAddress
	metering.UseGas(gasToUse)

	owner := runtime.GetSCAddress()
	err := runtime.MemStore(resultOffset, owner)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}
}

//export v1_2_getOwnerAddress
func v1_2_getOwnerAddress(context unsafe.Pointer, resultOffset int32) {
	blockchain := arwen.GetBlockchainContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetOwnerAddress
	metering.UseGas(gasToUse)

	owner, err := blockchain.GetOwnerAddress()
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	err = runtime.MemStore(resultOffset, owner)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}
}

//export v1_2_getShardOfAddress
func v1_2_getShardOfAddress(context unsafe.Pointer, addressOffset int32) int32 {
	blockchain := arwen.GetBlockchainContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetShardOfAddress
	metering.UseGas(gasToUse)

	address, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}

	return int32(blockchain.GetShardOfAddress(address))
}

//export v1_2_isSmartContract
func v1_2_isSmartContract(context unsafe.Pointer, addressOffset int32) int32 {
	blockchain := arwen.GetBlockchainContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.IsSmartContract
	metering.UseGas(gasToUse)

	address, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}

	isSmartContract := blockchain.IsSmartContract(address)
	return int32(arwen.BooleanToInt(isSmartContract))
}

//export v1_2_signalError
func v1_2_signalError(context unsafe.Pointer, messageOffset int32, messageLength int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.SignalError
	metering.UseGas(gasToUse)

	message, err := runtime.MemLoad(messageOffset, messageLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}
	runtime.SignalUserError(string(message))
}

//export v1_2_getExternalBalance
func v1_2_getExternalBalance(context unsafe.Pointer, addressOffset int32, resultOffset int32) {
	blockchain := arwen.GetBlockchainContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetExternalBalance
	metering.UseGas(gasToUse)

	address, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	balance := blockchain.GetBalance(address)

	err = runtime.MemStore(resultOffset, balance)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}
}

//export v1_2_blockHash
func v1_2_blockHash(context unsafe.Pointer, nonce int64, resultOffset int32) int32 {
	blockchain := arwen.GetBlockchainContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockHash
	metering.UseGas(gasToUse)

	hash := blockchain.BlockHash(nonce)
	err := runtime.MemStore(resultOffset, hash)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

func isBuiltInCall(data string, host arwen.VMHost) bool {
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
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	blockchain := arwen.GetBlockchainContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetExternalBalance
	metering.UseGas(gasToUse)

	address, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
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

//export v1_2_getESDTBalance
func v1_2_getESDTBalance(
	context unsafe.Pointer,
	addressOffset int32,
	tokenIDOffset int32,
	tokenIDLen int32,
	nonce int64,
	resultOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}
	err = runtime.MemStore(resultOffset, esdtData.Value.Bytes())
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}

	return int32(len(esdtData.Value.Bytes()))
}

//export v1_2_getESDTNFTNameLength
func v1_2_getESDTNFTNameLength(
	context unsafe.Pointer,
	addressOffset int32,
	tokenIDOffset int32,
	tokenIDLen int32,
	nonce int64,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}
	if esdtData == nil || esdtData.TokenMetaData == nil {
		return 0
	}

	return int32(len(esdtData.TokenMetaData.Name))
}

//export v1_2_getESDTNFTAttributeLength
func v1_2_getESDTNFTAttributeLength(
	context unsafe.Pointer,
	addressOffset int32,
	tokenIDOffset int32,
	tokenIDLen int32,
	nonce int64,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}
	if esdtData == nil || esdtData.TokenMetaData == nil {
		return 0
	}

	return int32(len(esdtData.TokenMetaData.Attributes))
}

//export v1_2_getESDTNFTURILength
func v1_2_getESDTNFTURILength(
	context unsafe.Pointer,
	addressOffset int32,
	tokenIDOffset int32,
	tokenIDLen int32,
	nonce int64,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}
	if esdtData == nil || esdtData.TokenMetaData == nil {
		return 0
	}
	if len(esdtData.TokenMetaData.URIs) == 0 {
		return 0
	}

	return int32(len(esdtData.TokenMetaData.URIs[0]))
}

//export v1_2_getESDTTokenData
func v1_2_getESDTTokenData(
	context unsafe.Pointer,
	addressOffset int32,
	tokenIDOffset int32,
	tokenIDLen int32,
	nonce int64,
	valueOffset int32,
	propertiesOffset int32,
	hashOffset int32,
	nameOffset int32,
	attributesOffset int32,
	creatorOffset int32,
	royaltiesOffset int32,
	urisOffset int32,
) int32 {
	runtime := arwen.GetRuntimeContext(context)
	esdtData, err := getESDTDataFromBlockchainHook(context, addressOffset, tokenIDOffset, tokenIDLen, nonce)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}

	err = runtime.MemStore(valueOffset, esdtData.Value.Bytes())
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}
	err = runtime.MemStore(propertiesOffset, esdtData.Properties)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}

	if esdtData.TokenMetaData != nil {
		err = runtime.MemStore(hashOffset, esdtData.TokenMetaData.Hash)
		if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
			return 0
		}
		err = runtime.MemStore(nameOffset, esdtData.TokenMetaData.Name)
		if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
			return 0
		}
		err = runtime.MemStore(attributesOffset, esdtData.TokenMetaData.Attributes)
		if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
			return 0
		}
		err = runtime.MemStore(creatorOffset, esdtData.TokenMetaData.Creator)
		if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
			return 0
		}
		err = runtime.MemStore(royaltiesOffset, big.NewInt(int64(esdtData.TokenMetaData.Royalties)).Bytes())
		if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
			return 0
		}
		if len(esdtData.TokenMetaData.URIs) > 0 {
			err = runtime.MemStore(urisOffset, esdtData.TokenMetaData.URIs[0])
			if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
				return 0
			}
		}
	}
	return int32(len(esdtData.Value.Bytes()))
}

//export v1_2_transferValue
func v1_2_transferValue(context unsafe.Pointer, destOffset int32, valueOffset int32, dataOffset int32, length int32) int32 {
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()

	gasToUse := metering.GasSchedule().ElrondAPICost.TransferValue
	metering.UseGas(gasToUse)

	send := runtime.GetSCAddress()
	dest, err := runtime.MemLoad(destOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	valueBytes, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(length))
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	if isBuiltInCall(string(data), host) {
		return 1
	}

	err = output.Transfer(dest, send, 0, 0, big.NewInt(0).SetBytes(valueBytes), data, vm.DirectCall)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export v1_2_transferValueExecute
func v1_2_transferValueExecute(
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
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()

	gasToUse := metering.GasSchedule().ElrondAPICost.TransferValue
	metering.UseGas(gasToUse)

	send := runtime.GetSCAddress()
	dest, err := runtime.MemLoad(destOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	valueBytes, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	var contractCallInput *vmcommon.ContractCallInput
	if functionLength > 0 {
		contractCallInput, err = prepareIndirectContractCallInput(
			host,
			send,
			big.NewInt(0).SetBytes(valueBytes),
			gasLimit,
			destOffset,
			functionOffset,
			functionLength,
			numArguments,
			argumentsLengthOffset,
			dataOffset,
			false,
		)
		if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
			return 1
		}
	}

	if contractCallInput != nil {
		if host.IsBuiltinFunctionName(contractCallInput.Function) {
			return 1
		}
	}

	if host.AreInSameShard(send, dest) && contractCallInput != nil && host.Blockchain().IsSmartContract(dest) {
		logEEI.Trace("eGLD pre-transfer execution begin")
		_, _, _, err = host.ExecuteOnDestContext(contractCallInput)
		if err != nil {
			logEEI.Trace("eGLD pre-transfer execution failed", "error", err)
			return 1
		}

		return 0
	}

	data := makeCrossShardCallFromInput(contractCallInput)
	err = output.Transfer(dest, send, uint64(gasLimit), 0, big.NewInt(0).SetBytes(valueBytes), []byte(data), vm.DirectCall)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

func makeCrossShardCallFromInput(vmInput *vmcommon.ContractCallInput) string {
	if vmInput == nil {
		return ""
	}

	txData := vmInput.Function
	for _, arg := range vmInput.Arguments {
		txData += "@" + hex.EncodeToString(arg)
	}

	return txData
}

//export v1_2_transferESDT
func v1_2_transferESDT(
	context unsafe.Pointer,
	destOffset int32,
	tokenIDOffset int32,
	tokenIDLen int32,
	valueOffset int32,
	gasLimit int64,
	dataOffset int32,
	length int32,
) int32 {
	host := arwen.GetVMHost(context)
	metering := host.Metering()

	gasToUse := metering.GasSchedule().ElrondAPICost.TransferValue
	metering.UseGas(gasToUse)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(length))
	metering.UseGas(gasToUse)
	logEEI.Warn("transferESDT() is deprecated")
	// this is only for backward compatibility - function deprecated
	return 1
}

//export v1_2_transferESDTExecute
func v1_2_transferESDTExecute(
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
	return v1_2_transferESDTNFTExecute(context, destOffset, tokenIDOffset, tokenIDLen, valueOffset, 0,
		gasLimit, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
}

//export v1_2_transferESDTNFTExecute
func v1_2_transferESDTNFTExecute(
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
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	output := host.Output()

	gasToUse := metering.GasSchedule().ElrondAPICost.TransferValue
	metering.UseGas(gasToUse)

	sender := runtime.GetSCAddress()
	dest, err := runtime.MemLoad(destOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	valueBytes, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	tokenIdentifier, err := runtime.MemLoad(tokenIDOffset, tokenIDLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	var contractCallInput *vmcommon.ContractCallInput
	if functionLength > 0 {
		contractCallInput, err = prepareIndirectContractCallInput(
			host,
			sender,
			big.NewInt(0),
			gasLimit,
			destOffset,
			functionOffset,
			functionLength,
			numArguments,
			argumentsLengthOffset,
			dataOffset,
			false,
		)
		if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
			return 1
		}

		esdtTokenType := core.Fungible
		if nonce > 0 {
			esdtTokenType = core.NonFungible
		}
		contractCallInput.ESDTTransfers = make([]*vmcommon.ESDTTransfer, 1)
		contractCallInput.ESDTTransfers[0] = &vmcommon.ESDTTransfer{
			ESDTValue:      big.NewInt(0).SetBytes(valueBytes),
			ESDTTokenName:  tokenIdentifier,
			ESDTTokenType:  uint32(esdtTokenType),
			ESDTTokenNonce: uint64(nonce),
		}
	}

	gasLimitForExec, err := output.TransferESDT(dest, sender, tokenIdentifier, uint64(nonce), big.NewInt(0).SetBytes(valueBytes), contractCallInput)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	if host.AreInSameShard(sender, dest) && contractCallInput != nil && host.Blockchain().IsSmartContract(dest) {
		contractCallInput.GasProvided = gasLimitForExec
		logEEI.Trace("ESDT post-transfer execution begin")
		_, _, _, err = host.ExecuteOnDestContext(contractCallInput)
		if err != nil {
			logEEI.Trace("ESDT post-transfer execution failed", "error", err)
			_, _, err = host.ExecuteESDTTransfer(sender, dest, tokenIdentifier, uint64(nonce), big.NewInt(0).SetBytes(valueBytes), vm.AsynchronousCallBack, true)
			if err != nil {
				logEEI.Warn("ESDT revert failed - forced fail execution for context", "error", err)
				_ = arwen.WithFault(err, context, true)
			}
			return 1
		}

		return 0
	}

	return 0
}

//export v1_2_createAsyncCall
func v1_2_createAsyncCall(context unsafe.Pointer,
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
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()

	// TODO consume gas

	acIdentifier, err := runtime.MemLoad(asyncContextIdentifier, identifierLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	calledSCAddress, err := runtime.MemLoad(destOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	successFunc, err := runtime.MemLoad(successOffset, successLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	errorFunc, err := runtime.MemLoad(errorOffset, errorLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	err = runtime.AddAsyncContextCall(acIdentifier, &arwen.AsyncGeneratedCall{
		Destination:     calledSCAddress,
		Data:            data,
		ValueBytes:      value,
		SuccessCallback: string(successFunc),
		ErrorCallback:   string(errorFunc),
		ProvidedGas:     uint64(gas),
	})
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}
}

//export v1_2_setAsyncContextCallback
func v1_2_setAsyncContextCallback(context unsafe.Pointer,
	asyncContextIdentifier int32,
	identifierLength int32,
	callback int32,
	callbackLength int32,
) int32 {
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()

	// TODO consume gas

	acIdentifier, err := runtime.MemLoad(asyncContextIdentifier, identifierLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	asyncContext, err := runtime.GetAsyncContext(acIdentifier)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	callbackFunc, err := runtime.MemLoad(callback, callbackLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	asyncContext.Callback = string(callbackFunc)

	return 0
}

//export v1_2_upgradeContract
func v1_2_upgradeContract(
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
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().ElrondAPICost.CreateContract
	metering.UseGas(gasToUse)

	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	code, err := runtime.MemLoad(codeOffset, length)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	codeMetadata, err := runtime.MemLoad(codeMetadataOffset, arwen.CodeMetadataLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	data, actualLen, err := getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	metering.UseGas(gasToUse)

	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	gasSchedule := metering.GasSchedule()
	gasToUse = gasSchedule.ElrondAPICost.AsyncCallStep
	metering.UseGas(gasToUse)

	calledSCAddress, err := runtime.MemLoad(destOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	gasToUse = math.MulUint64(gasSchedule.BaseOperationCost.DataCopyPerByte, uint64(length))
	metering.UseGas(gasToUse)

	minAsyncCallCost := math.AddUint64(
		math.MulUint64(2, gasSchedule.ElrondAPICost.AsyncCallStep),
		gasSchedule.ElrondAPICost.AsyncCallbackGasLock)
	if uint64(gasLimit) < minAsyncCallCost {
		runtime.SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
		return
	}

	// Set up the async call as if it is not known whether the called SC
	// is in the same shard with the caller or not. This will be later resolved
	// in the handler for BreakpointAsyncCall.
	codeEncoded := hex.EncodeToString(code)
	codeMetadataEncoded := hex.EncodeToString(codeMetadata)
	finalData := arwen.UpgradeFunctionName + "@" + codeEncoded + "@" + codeMetadataEncoded
	for _, arg := range data {
		finalData += "@" + string(arg)
	}

	runtime.SetAsyncCallInfo(&arwen.AsyncCallInfo{
		Destination: calledSCAddress,
		Data:        []byte(finalData),
		GasLimit:    uint64(gasLimit),
		ValueBytes:  value,
	})

	// Instruct Wasmer to interrupt the execution of the caller SC.
	runtime.SetRuntimeBreakpointValue(arwen.BreakpointAsyncCall)
}

//export v1_2_asyncCall
func v1_2_asyncCall(context unsafe.Pointer, destOffset int32, valueOffset int32, dataOffset int32, length int32) {
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()

	gasSchedule := metering.GasSchedule()
	gasToUse := gasSchedule.ElrondAPICost.AsyncCallStep
	metering.UseGas(gasToUse)

	calledSCAddress, err := runtime.MemLoad(destOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	gasToUse = math.MulUint64(gasSchedule.BaseOperationCost.DataCopyPerByte, uint64(length))
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	err = runtime.ExecuteAsyncCall(calledSCAddress, data, value)
	if errors.Is(err, arwen.ErrNotEnoughGas) {
		runtime.SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
		return
	}
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}
}

//export v1_2_getArgumentLength
func v1_2_getArgumentLength(context unsafe.Pointer, id int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if id < 0 || int32(len(args)) <= id {
		return -1
	}

	return int32(len(args[id]))
}

//export v1_2_getArgument
func v1_2_getArgument(context unsafe.Pointer, id int32, argOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if id < 0 || int32(len(args)) <= id {
		return -1
	}

	err := runtime.MemStore(argOffset, args[id])
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(args[id]))
}

//export v1_2_getFunction
func v1_2_getFunction(context unsafe.Pointer, functionOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetFunction
	metering.UseGas(gasToUse)

	function := runtime.Function()
	err := runtime.MemStore(functionOffset, []byte(function))
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(function))
}

//export v1_2_getNumArguments
func v1_2_getNumArguments(context unsafe.Pointer) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetNumArguments
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	return int32(len(args))
}

//export v1_2_storageStore
func v1_2_storageStore(context unsafe.Pointer, keyOffset int32, keyLength int32, dataOffset int32, dataLength int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.StorageStore
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	storageStatus, err := storage.SetStorage(key, data)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(storageStatus)
}

//export v1_2_storageLoadLength
func v1_2_storageLoadLength(context unsafe.Pointer, keyOffset int32, keyLength int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.StorageLoad
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	data := storage.GetStorageUnmetered(key)

	return int32(len(data))
}

//export v1_2_storageLoadFromAddress
func v1_2_storageLoadFromAddress(context unsafe.Pointer, addressOffset int32, keyOffset int32, keyLength int32, dataOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.StorageLoad
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	address, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	data := storage.GetStorageFromAddress(address, key)

	err = runtime.MemStore(dataOffset, data)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(data))
}

//export v1_2_storageLoad
func v1_2_storageLoad(context unsafe.Pointer, keyOffset int32, keyLength int32, dataOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.StorageLoad
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	data := storage.GetStorage(key)

	err = runtime.MemStore(dataOffset, data)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(data))
}

//export v1_2_setStorageLock
func v1_2_setStorageLock(context unsafe.Pointer, keyOffset int32, keyLength int32, lockTimestamp int64) int32 {
	runtime := arwen.GetRuntimeContext(context)
	storage := arwen.GetStorageContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.Int64StorageStore
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	timeLockKey := arwen.CustomStorageKey(arwen.TimeLockKeyPrefix, key)
	bigTimestamp := big.NewInt(0).SetInt64(lockTimestamp)
	storageStatus, err := storage.SetProtectedStorage(timeLockKey, bigTimestamp.Bytes())
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}
	return int32(storageStatus)
}

//export v1_2_getStorageLock
func v1_2_getStorageLock(context unsafe.Pointer, keyOffset int32, keyLength int32) int64 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	storage := arwen.GetStorageContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.StorageLoad
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	timeLockKey := arwen.CustomStorageKey(arwen.TimeLockKeyPrefix, key)
	data := storage.GetStorage(timeLockKey)
	timeLock := big.NewInt(0).SetBytes(data).Int64()

	return timeLock
}

//export v1_2_isStorageLocked
func v1_2_isStorageLocked(context unsafe.Pointer, keyOffset int32, keyLength int32) int32 {
	timeLock := v1_2_getStorageLock(context, keyOffset, keyLength)
	if timeLock < 0 {
		return -1
	}

	currentTimestamp := v1_2_getBlockTimestamp(context)
	if timeLock <= currentTimestamp {
		return 0
	}

	return 1
}

//export v1_2_clearStorageLock
func v1_2_clearStorageLock(context unsafe.Pointer, keyOffset int32, keyLength int32) int32 {
	return v1_2_setStorageLock(context, keyOffset, keyLength, 0)
}

//export v1_2_getCaller
func v1_2_getCaller(context unsafe.Pointer, resultOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCaller
	metering.UseGas(gasToUse)

	caller := runtime.GetVMInput().CallerAddr

	err := runtime.MemStore(resultOffset, caller)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}
}

//export v1_2_checkNoPayment
func v1_2_checkNoPayment(context unsafe.Pointer) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCallValue
	metering.UseGas(gasToUse)

	vmInput := runtime.GetVMInput()
	if vmInput.CallValue.Sign() > 0 {
		runtime := arwen.GetRuntimeContext(context)
		arwen.WithFault(arwen.ErrNonPayableFunctionEgld, context, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}
	if len(vmInput.ESDTTransfers) > 0 {
		runtime := arwen.GetRuntimeContext(context)
		arwen.WithFault(arwen.ErrNonPayableFunctionEsdt, context, runtime.ElrondAPIErrorShouldFailExecution())
		return
	}
}

//export v1_2_callValue
func v1_2_callValue(context unsafe.Pointer, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCallValue
	metering.UseGas(gasToUse)

	value := runtime.GetVMInput().CallValue.Bytes()
	value = arwen.PadBytesLeft(value, arwen.BalanceLen)

	err := runtime.MemStore(resultOffset, value)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(value))
}

//export v1_2_getESDTValue
func v1_2_getESDTValue(context unsafe.Pointer, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCallValue
	metering.UseGas(gasToUse)

	var value []byte

	esdtTransfer := getFirstESDTTransferIfExist(runtime.GetVMInput())
	if esdtTransfer.ESDTValue.Cmp(arwen.Zero) > 0 {
		value = esdtTransfer.ESDTValue.Bytes()
		value = arwen.PadBytesLeft(value, arwen.BalanceLen)
	}

	err := runtime.MemStore(resultOffset, value)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(value))
}

//export v1_2_getESDTTokenName
func v1_2_getESDTTokenName(context unsafe.Pointer, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCallValue
	metering.UseGas(gasToUse)

	esdtTransfer := getFirstESDTTransferIfExist(runtime.GetVMInput())
	tokenName := esdtTransfer.ESDTTokenName

	err := runtime.MemStore(resultOffset, tokenName)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(tokenName))
}

//export v1_2_getESDTTokenNonce
func v1_2_getESDTTokenNonce(context unsafe.Pointer) int64 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCallValue
	metering.UseGas(gasToUse)

	esdtTransfer := getFirstESDTTransferIfExist(runtime.GetVMInput())
	return int64(esdtTransfer.ESDTTokenNonce)
}

//export v1_2_getCurrentESDTNFTNonce
func v1_2_getCurrentESDTNFTNonce(context unsafe.Pointer, addressOffset int32, tokenIDOffset int32, tokenIDLen int32) int64 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	storage := arwen.GetStorageContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.StorageLoad
	metering.UseGas(gasToUse)

	destination, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if err != nil {
		return 0
	}

	tokenID, err := runtime.MemLoad(tokenIDOffset, tokenIDLen)
	if err != nil {
		return 0
	}

	key := []byte(core.ElrondProtectedKeyPrefix + core.ESDTNFTLatestNonceIdentifier + string(tokenID))
	data := storage.GetStorageFromAddress(destination, key)

	nonce := big.NewInt(0).SetBytes(data).Uint64()
	return int64(nonce)
}

//export v1_2_getESDTTokenType
func v1_2_getESDTTokenType(context unsafe.Pointer) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCallValue
	metering.UseGas(gasToUse)

	esdtTransfer := getFirstESDTTransferIfExist(runtime.GetVMInput())
	return int32(esdtTransfer.ESDTTokenType)
}

//export v1_2_getCallValueTokenName
func v1_2_getCallValueTokenName(context unsafe.Pointer, callValueOffset int32, tokenNameOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCallValue
	metering.UseGas(gasToUse)

	callValue := runtime.GetVMInput().CallValue.Bytes()
	tokenName := make([]byte, 0)
	esdtTransfer := getFirstESDTTransferIfExist(runtime.GetVMInput())

	if len(esdtTransfer.ESDTTokenName) > 0 {
		tokenName = make([]byte, 0, len(esdtTransfer.ESDTTokenName))
		copy(tokenName, esdtTransfer.ESDTTokenName)
		callValue = esdtTransfer.ESDTValue.Bytes()
	}
	callValue = arwen.PadBytesLeft(callValue, arwen.BalanceLen)

	err := runtime.MemStore(tokenNameOffset, tokenName)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	err = runtime.MemStore(callValueOffset, callValue)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(tokenName))
}

//export v1_2_writeLog
func v1_2_writeLog(context unsafe.Pointer, dataPointer int32, dataLength int32, topicPtr int32, numTopics int32) {
	// note: deprecated
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.Log
	gas := math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(numTopics*arwen.HashLen+dataLength))
	gasToUse = math.AddUint64(gasToUse, gas)
	metering.UseGas(gasToUse)

	log, err := runtime.MemLoad(dataPointer, dataLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	topics, err := arwen.GuardedMakeByteSlice2D(numTopics)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	for i := int32(0); i < numTopics; i++ {
		topics[i], err = runtime.MemLoad(topicPtr+i*arwen.HashLen, arwen.HashLen)
		if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
			return
		}
	}

	output.WriteLog(runtime.GetSCAddress(), topics, log)
}

//export v1_2_writeEventLog
func v1_2_writeEventLog(
	context unsafe.Pointer,
	numTopics int32,
	topicLengthsOffset int32,
	topicOffset int32,
	dataOffset int32,
	dataLength int32) {

	host := arwen.GetVMHost(context)
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	topics, topicDataTotalLen, err := getArgumentsFromMemory(
		host,
		numTopics,
		topicLengthsOffset,
		topicOffset,
	)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	gasToUse := metering.GasSchedule().ElrondAPICost.Log
	gasForData := math.MulUint64(
		metering.GasSchedule().BaseOperationCost.DataCopyPerByte,
		uint64(topicDataTotalLen+dataLength))
	gasToUse = math.AddUint64(gasToUse, gasForData)
	metering.UseGas(gasToUse)

	output.WriteLog(runtime.GetSCAddress(), topics, data)
}

//export v1_2_getBlockTimestamp
func v1_2_getBlockTimestamp(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockTimeStamp
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentTimeStamp())
}

//export v1_2_getBlockNonce
func v1_2_getBlockNonce(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockNonce
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentNonce())
}

//export v1_2_getBlockRound
func v1_2_getBlockRound(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockRound
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentRound())
}

//export v1_2_getBlockEpoch
func v1_2_getBlockEpoch(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockEpoch
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentEpoch())
}

//export v1_2_getBlockRandomSeed
func v1_2_getBlockRandomSeed(context unsafe.Pointer, pointer int32) {
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockRandomSeed
	metering.UseGas(gasToUse)

	randomSeed := blockchain.CurrentRandomSeed()
	err := runtime.MemStore(pointer, randomSeed)
	arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution())
}

//export v1_2_getStateRootHash
func v1_2_getStateRootHash(context unsafe.Pointer, pointer int32) {
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetStateRootHash
	metering.UseGas(gasToUse)

	stateRootHash := blockchain.GetStateRootHash()
	err := runtime.MemStore(pointer, stateRootHash)
	arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution())
}

//export v1_2_getPrevBlockTimestamp
func v1_2_getPrevBlockTimestamp(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockTimeStamp
	metering.UseGas(gasToUse)

	return int64(blockchain.LastTimeStamp())
}

//export v1_2_getPrevBlockNonce
func v1_2_getPrevBlockNonce(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockNonce
	metering.UseGas(gasToUse)

	return int64(blockchain.LastNonce())
}

//export v1_2_getPrevBlockRound
func v1_2_getPrevBlockRound(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockRound
	metering.UseGas(gasToUse)

	return int64(blockchain.LastRound())
}

//export v1_2_getPrevBlockEpoch
func v1_2_getPrevBlockEpoch(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockEpoch
	metering.UseGas(gasToUse)

	return int64(blockchain.LastEpoch())
}

//export v1_2_getPrevBlockRandomSeed
func v1_2_getPrevBlockRandomSeed(context unsafe.Pointer, pointer int32) {
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockRandomSeed
	metering.UseGas(gasToUse)

	randomSeed := blockchain.LastRandomSeed()
	err := runtime.MemStore(pointer, randomSeed)
	arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution())
}

//export v1_2_returnData
func v1_2_returnData(context unsafe.Pointer, pointer int32, length int32) {
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.Finish
	gas := math.MulUint64(metering.GasSchedule().BaseOperationCost.PersistPerByte, uint64(length))
	gasToUse = math.AddUint64(gasToUse, gas)
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(pointer, length)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	output.Finish(data)
}

//export v1_2_executeOnSameContext
func v1_2_executeOnSameContext(
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
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().ElrondAPICost.ExecuteOnSameContext
	metering.UseGas(gasToUse)

	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	sender := runtime.GetSCAddress()
	bigIntVal := big.NewInt(0).SetBytes(value)
	contractCallInput, err := prepareIndirectContractCallInput(
		host,
		sender,
		bigIntVal,
		gasLimit,
		addressOffset,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
		true,
	)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	if isBuiltInCall(contractCallInput.Function, host) {
		return 1
	}

	_, err = host.ExecuteOnSameContext(contractCallInput)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

func fillContractCallInputFromArgs(
	contractCallInput *vmcommon.ContractCallInput,
	args [][]byte,
	startLength int,
) {
	lenArgs := len(args)
	if lenArgs > startLength {
		contractCallInput.Function = string(args[startLength])
	}

	if lenArgs > startLength+1 {
		contractCallInput.Arguments = args[startLength+1:]
	}
}

func doESDTTransferAndExecuteSynchronously(
	context unsafe.Pointer,
	destination []byte,
	value *big.Int,
	function string,
	args [][]byte,
	gasLimit int64,
) int32 {
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()

	if value.Cmp(arwen.Zero) > 0 {
		if arwen.WithFault(arwen.ErrTransferValueOnESDTCall, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
			return 1
		}
	}
	if len(args) < 1 {
		if arwen.WithFault(arwen.ErrArgOutOfRange, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
			return 1
		}
	}

	sender := runtime.GetSCAddress()
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   make([][]byte, 0),
			CallValue:   big.NewInt(0),
			CallType:    vm.DirectCall,
			GasProvided: metering.BoundGasLimit(gasLimit),
		},
		RecipientAddr: destination,
		Function:      "",
	}

	tokenID := args[0]
	esdtValue := big.NewInt(0)
	nonce := uint64(0)

	contractCallInput.ESDTTransfers = make([]*vmcommon.ESDTTransfer, 1)
	contractCallInput.ESDTTransfers[0] = &vmcommon.ESDTTransfer{
		ESDTValue: big.NewInt(0),
	}

	switch function {
	case core.BuiltInFunctionESDTTransfer:
		if len(args) < core.MinLenArgumentsESDTTransfer {
			if arwen.WithFault(arwen.ErrArgOutOfRange, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
				return 1
			}
		}

		esdtValue.SetBytes(args[1])
		contractCallInput.ESDTTransfers[0].ESDTTokenType = uint32(core.Fungible)
		fillContractCallInputFromArgs(contractCallInput, args, core.MinLenArgumentsESDTTransfer)

	case core.BuiltInFunctionESDTNFTTransfer:
		if len(args) < core.MinLenArgumentsESDTNFTTransfer {
			if arwen.WithFault(arwen.ErrArgOutOfRange, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
				return 1
			}
		}
		nonce = big.NewInt(0).SetBytes(args[1]).Uint64()
		esdtValue.SetBytes(args[2])
		if !bytes.Equal(destination, args[3]) {
			if arwen.WithFault(arwen.ErrFailedTransfer, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
				return 1
			}
		}
		contractCallInput.ESDTTransfers[0].ESDTTokenType = uint32(core.NonFungible)
		fillContractCallInputFromArgs(contractCallInput, args, core.MinLenArgumentsESDTNFTTransfer)

	default:
		if arwen.WithFault(arwen.ErrFuncNotFound, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
			return 1
		}
	}

	contractCallInput.ESDTTransfers[0].ESDTTokenName = tokenID
	contractCallInput.ESDTTransfers[0].ESDTValue = esdtValue
	contractCallInput.ESDTTransfers[0].ESDTTokenNonce = nonce
	if len(contractCallInput.Function) == 0 {
		contractCallInput = nil
	}
	gasLimitForExec, err := output.TransferESDT(destination, sender, tokenID, nonce, esdtValue, contractCallInput)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	if contractCallInput != nil && host.Blockchain().IsSmartContract(destination) {
		contractCallInput.GasProvided = gasLimitForExec
		logEEI.Trace("ESDT post-transfer execution begin")
		_, _, _, err = host.ExecuteOnDestContext(contractCallInput)
		if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
			logEEI.Trace("ESDT post-transfer execution failed", "error", err)
			_, _, err = host.ExecuteESDTTransfer(sender, destination, tokenID, nonce, esdtValue, vm.AsynchronousCallBack, true)
			if err != nil {
				logEEI.Warn("ESDT revert failed - forced fail execution for context", "error", err)
				_ = arwen.WithFault(err, context, true)
			}
			return 1
		}
	}

	return 0
}

func detectSyncESDTTransfer(
	context unsafe.Pointer,
	functionOffset int32,
	functionLength int32,
) (string, bool, error) {
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()

	if !host.IsESDTFunctionsEnabled() {
		return "", false, nil
	}

	function, err := runtime.MemLoad(functionOffset, functionLength)
	if err != nil {
		return "", false, err
	}

	if string(function) == core.BuiltInFunctionESDTTransfer ||
		string(function) == core.BuiltInFunctionESDTNFTTransfer {
		return string(function), true, nil
	}

	return "", false, nil
}

func getDestinationAndArguments(
	context unsafe.Pointer,
	numArguments int32,
	argumentsLengthOffset int32,
	addressOffset int32,
	dataOffset int32,
) ([]byte, [][]byte, error) {
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()

	destination, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if err != nil {
		return nil, nil, err
	}

	if !host.AreInSameShard(runtime.GetSCAddress(), destination) {
		return nil, nil, arwen.ErrSyncExecutionNotInSameShard
	}

	data, actualLen, err := getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if err != nil {
		return nil, nil, err
	}

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	metering.UseGas(gasToUse)

	return destination, data, nil
}

//export v1_2_executeOnDestContext
func v1_2_executeOnDestContext(
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
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().ElrondAPICost.ExecuteOnDestContext
	metering.UseGas(gasToUse)

	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	function, isSyncESDT, err := detectSyncESDTTransfer(context, functionOffset, functionLength)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	if isSyncESDT {
		destination, data, err := getDestinationAndArguments(context, numArguments, argumentsLengthOffset, addressOffset, dataOffset)
		if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
			return 1
		}

		return doESDTTransferAndExecuteSynchronously(context, destination, big.NewInt(0).SetBytes(value), function, data, gasLimit)
	}

	sender := runtime.GetSCAddress()
	bigIntVal := big.NewInt(0).SetBytes(value)
	contractCallInput, err := prepareIndirectContractCallInput(
		host,
		sender,
		bigIntVal,
		gasLimit,
		addressOffset,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
		true,
	)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	_, _, gasUsedBeforeReset, err := host.ExecuteOnDestContext(contractCallInput)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}
	metering.UseGas(gasUsedBeforeReset)

	return 0
}

//export v1_2_executeOnDestContextByCaller
func v1_2_executeOnDestContextByCaller(
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
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().ElrondAPICost.ExecuteOnDestContext
	metering.UseGas(gasToUse)

	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	sender := runtime.GetVMInput().CallerAddr
	bigIntVal := big.NewInt(0).SetBytes(value)
	contractCallInput, err := prepareIndirectContractCallInput(
		host,
		sender,
		bigIntVal,
		gasLimit,
		addressOffset,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
		true,
	)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	if isBuiltInCall(contractCallInput.Function, host) {
		return 1
	}

	_, _, _, err = host.ExecuteOnDestContext(contractCallInput)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export v1_2_delegateExecution
func v1_2_delegateExecution(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().ElrondAPICost.DelegateExecution
	metering.UseGas(gasToUse)

	sender := runtime.GetSCAddress()
	value := runtime.GetVMInput().CallValue
	bigIntVal := big.NewInt(0).Set(value)
	contractCallInput, err := prepareIndirectContractCallInput(
		host,
		sender,
		bigIntVal,
		gasLimit,
		addressOffset,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
		true,
	)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	if isBuiltInCall(contractCallInput.Function, host) {
		return 1
	}

	_, err = host.ExecuteOnSameContext(contractCallInput)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export v1_2_executeReadOnly
func v1_2_executeReadOnly(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().ElrondAPICost.ExecuteReadOnly
	metering.UseGas(gasToUse)

	sender := runtime.GetSCAddress()
	value := runtime.GetVMInput().CallValue
	bigIntVal := big.NewInt(0).Set(value)

	contractCallInput, err := prepareIndirectContractCallInput(
		host,
		sender,
		bigIntVal,
		gasLimit,
		addressOffset,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
		true,
	)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	if isBuiltInCall(contractCallInput.Function, host) {
		return 1
	}

	runtime.SetReadOnly(true)
	_, err = host.ExecuteOnSameContext(contractCallInput)
	runtime.SetReadOnly(false)
	if arwen.WithFault(err, context, runtime.ElrondSyncExecAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export v1_2_createContract
func v1_2_createContract(
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
	host := arwen.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().ElrondAPICost.CreateContract
	metering.UseGas(gasToUse)

	sender := runtime.GetSCAddress()
	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	code, err := runtime.MemLoad(codeOffset, length)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	codeMetadata, err := runtime.MemLoad(codeMetadataOffset, arwen.CodeMetadataLen)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	data, actualLen, err := getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	metering.UseGas(gasToUse)

	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	contractCreate := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   data,
			CallValue:   big.NewInt(0).SetBytes(value),
			GasPrice:    0,
			GasProvided: metering.BoundGasLimit(gasLimit),
		},
		ContractCode:         code,
		ContractCodeMetadata: codeMetadata,
	}

	newAddress, err := host.CreateNewContract(contractCreate)
	if err != nil {
		return 1
	}

	err = runtime.MemStore(resultOffset, newAddress)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export v1_2_getNumReturnData
func v1_2_getNumReturnData(context unsafe.Pointer) int32 {
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetNumReturnData
	metering.UseGas(gasToUse)

	returnData := output.ReturnData()
	return int32(len(returnData))
}

//export v1_2_getReturnDataSize
func v1_2_getReturnDataSize(context unsafe.Pointer, resultID int32) int32 {
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetReturnDataSize
	metering.UseGas(gasToUse)

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) {
		return 0
	}

	return int32(len(returnData[resultID]))
}

//export v1_2_getReturnData
func v1_2_getReturnData(context unsafe.Pointer, resultID int32, dataOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetReturnData
	metering.UseGas(gasToUse)

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) {
		return 0
	}

	err := runtime.MemStore(dataOffset, returnData[resultID])
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return 0
	}

	return int32(len(returnData[resultID]))
}

//export v1_2_getOriginalTxHash
func v1_2_getOriginalTxHash(context unsafe.Pointer, dataOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetBlockHash
	metering.UseGas(gasToUse)

	err := runtime.MemStore(dataOffset, runtime.GetOriginalTxHash())
	_ = arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution())
}

func prepareIndirectContractCallInput(
	host arwen.VMHost,
	sender []byte,
	value *big.Int,
	gasLimit int64,
	addressOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
	syncExecutionRequired bool,
) (*vmcommon.ContractCallInput, error) {
	runtime := host.Runtime()
	metering := host.Metering()

	destination, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if err != nil {
		return nil, err
	}

	if syncExecutionRequired && !host.AreInSameShard(runtime.GetSCAddress(), destination) {
		return nil, arwen.ErrSyncExecutionNotInSameShard
	}

	function, err := runtime.MemLoad(functionOffset, functionLength)
	if err != nil {
		return nil, err
	}

	data, actualLen, err := getArgumentsFromMemory(
		host,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, uint64(actualLen))
	metering.UseGas(gasToUse)
	if err != nil {
		return nil, err
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
	host arwen.VMHost,
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

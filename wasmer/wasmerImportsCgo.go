package wasmer

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// !!!!!!!!!!!!!!!!!!!!!! AUTO-GENERATED FILE !!!!!!!!!!!!!!!!!!!!!!
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef int int32_t;
//
// extern int32_t   v1_5_bigFloatNewFromParts(void* context, int32_t integralPart, int32_t fractionalPart, int32_t exponent);
// extern int32_t   v1_5_bigFloatNewFromFrac(void* context, long long numerator, long long denominator);
// extern int32_t   v1_5_bigFloatNewFromSci(void* context, long long significand, long long exponent);
// extern void      v1_5_bigFloatAdd(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void      v1_5_bigFloatSub(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void      v1_5_bigFloatMul(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void      v1_5_bigFloatDiv(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void      v1_5_bigFloatNeg(void* context, int32_t destinationHandle, int32_t opHandle);
// extern void      v1_5_bigFloatClone(void* context, int32_t destinationHandle, int32_t opHandle);
// extern int32_t   v1_5_bigFloatCmp(void* context, int32_t op1Handle, int32_t op2Handle);
// extern void      v1_5_bigFloatAbs(void* context, int32_t destinationHandle, int32_t opHandle);
// extern int32_t   v1_5_bigFloatSign(void* context, int32_t opHandle);
// extern void      v1_5_bigFloatSqrt(void* context, int32_t destinationHandle, int32_t opHandle);
// extern void      v1_5_bigFloatPow(void* context, int32_t destinationHandle, int32_t opHandle, int32_t exponent);
// extern void      v1_5_bigFloatFloor(void* context, int32_t destBigIntHandle, int32_t opHandle);
// extern void      v1_5_bigFloatCeil(void* context, int32_t destBigIntHandle, int32_t opHandle);
// extern void      v1_5_bigFloatTruncate(void* context, int32_t destBigIntHandle, int32_t opHandle);
// extern void      v1_5_bigFloatSetInt64(void* context, int32_t destinationHandle, long long value);
// extern int32_t   v1_5_bigFloatIsInt(void* context, int32_t opHandle);
// extern void      v1_5_bigFloatSetBigInt(void* context, int32_t destinationHandle, int32_t bigIntHandle);
// extern void      v1_5_bigFloatGetConstPi(void* context, int32_t destinationHandle);
// extern void      v1_5_bigFloatGetConstE(void* context, int32_t destinationHandle);
// extern void      v1_5_bigIntGetUnsignedArgument(void* context, int32_t id, int32_t destinationHandle);
// extern void      v1_5_bigIntGetSignedArgument(void* context, int32_t id, int32_t destinationHandle);
// extern int32_t   v1_5_bigIntStorageStoreUnsigned(void* context, int32_t keyOffset, int32_t keyLength, int32_t sourceHandle);
// extern int32_t   v1_5_bigIntStorageLoadUnsigned(void* context, int32_t keyOffset, int32_t keyLength, int32_t destinationHandle);
// extern void      v1_5_bigIntGetCallValue(void* context, int32_t destinationHandle);
// extern void      v1_5_bigIntGetESDTCallValue(void* context, int32_t destination);
// extern void      v1_5_bigIntGetESDTCallValueByIndex(void* context, int32_t destinationHandle, int32_t index);
// extern void      v1_5_bigIntGetExternalBalance(void* context, int32_t addressOffset, int32_t result);
// extern void      v1_5_bigIntGetESDTExternalBalance(void* context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen, long long nonce, int32_t resultHandle);
// extern int32_t   v1_5_bigIntNew(void* context, long long smallValue);
// extern int32_t   v1_5_bigIntUnsignedByteLength(void* context, int32_t referenceHandle);
// extern int32_t   v1_5_bigIntSignedByteLength(void* context, int32_t referenceHandle);
// extern int32_t   v1_5_bigIntGetUnsignedBytes(void* context, int32_t referenceHandle, int32_t byteOffset);
// extern int32_t   v1_5_bigIntGetSignedBytes(void* context, int32_t referenceHandle, int32_t byteOffset);
// extern void      v1_5_bigIntSetUnsignedBytes(void* context, int32_t destinationHandle, int32_t byteOffset, int32_t byteLength);
// extern void      v1_5_bigIntSetSignedBytes(void* context, int32_t destinationHandle, int32_t byteOffset, int32_t byteLength);
// extern int32_t   v1_5_bigIntIsInt64(void* context, int32_t destinationHandle);
// extern long long v1_5_bigIntGetInt64(void* context, int32_t destinationHandle);
// extern void      v1_5_bigIntSetInt64(void* context, int32_t destinationHandle, long long value);
// extern void      v1_5_bigIntAdd(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void      v1_5_bigIntSub(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void      v1_5_bigIntMul(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void      v1_5_bigIntTDiv(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void      v1_5_bigIntTMod(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void      v1_5_bigIntEDiv(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void      v1_5_bigIntEMod(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void      v1_5_bigIntSqrt(void* context, int32_t destinationHandle, int32_t opHandle);
// extern void      v1_5_bigIntPow(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern int32_t   v1_5_bigIntLog2(void* context, int32_t op1Handle);
// extern void      v1_5_bigIntAbs(void* context, int32_t destinationHandle, int32_t opHandle);
// extern void      v1_5_bigIntNeg(void* context, int32_t destinationHandle, int32_t opHandle);
// extern int32_t   v1_5_bigIntSign(void* context, int32_t opHandle);
// extern int32_t   v1_5_bigIntCmp(void* context, int32_t op1Handle, int32_t op2Handle);
// extern void      v1_5_bigIntNot(void* context, int32_t destinationHandle, int32_t opHandle);
// extern void      v1_5_bigIntAnd(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void      v1_5_bigIntOr(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void      v1_5_bigIntXor(void* context, int32_t destinationHandle, int32_t op1Handle, int32_t op2Handle);
// extern void      v1_5_bigIntShr(void* context, int32_t destinationHandle, int32_t opHandle, int32_t bits);
// extern void      v1_5_bigIntShl(void* context, int32_t destinationHandle, int32_t opHandle, int32_t bits);
// extern void      v1_5_bigIntFinishUnsigned(void* context, int32_t referenceHandle);
// extern void      v1_5_bigIntFinishSigned(void* context, int32_t referenceHandle);
// extern void      v1_5_bigIntToString(void* context, int32_t bigIntHandle, int32_t destinationHandle);
// extern long long v1_5_getGasLeft(void* context);
// extern void      v1_5_getSCAddress(void* context, int32_t resultOffset);
// extern void      v1_5_getOwnerAddress(void* context, int32_t resultOffset);
// extern int32_t   v1_5_getShardOfAddress(void* context, int32_t addressOffset);
// extern int32_t   v1_5_isSmartContract(void* context, int32_t addressOffset);
// extern void      v1_5_signalError(void* context, int32_t messageOffset, int32_t messageLength);
// extern void      v1_5_getExternalBalance(void* context, int32_t addressOffset, int32_t resultOffset);
// extern int32_t   v1_5_blockHash(void* context, long long nonce, int32_t resultOffset);
// extern int32_t   v1_5_getESDTBalance(void* context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen, long long nonce, int32_t resultOffset);
// extern int32_t   v1_5_getESDTNFTNameLength(void* context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen, long long nonce);
// extern int32_t   v1_5_getESDTNFTAttributeLength(void* context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen, long long nonce);
// extern int32_t   v1_5_getESDTNFTURILength(void* context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen, long long nonce);
// extern int32_t   v1_5_getESDTTokenData(void* context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen, long long nonce, int32_t valueHandle, int32_t propertiesOffset, int32_t hashOffset, int32_t nameOffset, int32_t attributesOffset, int32_t creatorOffset, int32_t royaltiesHandle, int32_t urisOffset);
// extern long long v1_5_getESDTLocalRoles(void* context, int32_t tokenIdHandle);
// extern int32_t   v1_5_validateTokenIdentifier(void* context, int32_t tokenIdHandle);
// extern int32_t   v1_5_transferValue(void* context, int32_t destOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
// extern int32_t   v1_5_transferValueExecute(void* context, int32_t destOffset, int32_t valueOffset, long long gasLimit, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t   v1_5_transferESDTExecute(void* context, int32_t destOffset, int32_t tokenIDOffset, int32_t tokenIDLen, int32_t valueOffset, long long gasLimit, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t   v1_5_transferESDTNFTExecute(void* context, int32_t destOffset, int32_t tokenIDOffset, int32_t tokenIDLen, int32_t valueOffset, long long nonce, long long gasLimit, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t   v1_5_multiTransferESDTNFTExecute(void* context, int32_t destOffset, int32_t numTokenTransfers, int32_t tokenTransfersArgsLengthOffset, int32_t tokenTransferDataOffset, long long gasLimit, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t   v1_5_createAsyncCall(void* context, int32_t destOffset, int32_t valueOffset, int32_t dataOffset, int32_t dataLength, int32_t successOffset, int32_t successLength, int32_t errorOffset, int32_t errorLength, long long gas, long long extraGasForCallback);
// extern int32_t   v1_5_setAsyncContextCallback(void* context, int32_t callback, int32_t callbackLength, int32_t data, int32_t dataLength, long long gas);
// extern void      v1_5_upgradeContract(void* context, int32_t destOffset, long long gasLimit, int32_t valueOffset, int32_t codeOffset, int32_t codeMetadataOffset, int32_t length, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern void      v1_5_upgradeFromSourceContract(void* context, int32_t destOffset, long long gasLimit, int32_t valueOffset, int32_t sourceContractAddressOffset, int32_t codeMetadataOffset, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern void      v1_5_deleteContract(void* context, int32_t destOffset, long long gasLimit, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern void      v1_5_asyncCall(void* context, int32_t destOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
// extern int32_t   v1_5_getArgumentLength(void* context, int32_t id);
// extern int32_t   v1_5_getArgument(void* context, int32_t id, int32_t argOffset);
// extern int32_t   v1_5_getFunction(void* context, int32_t functionOffset);
// extern int32_t   v1_5_getNumArguments(void* context);
// extern int32_t   v1_5_storageStore(void* context, int32_t keyOffset, int32_t keyLength, int32_t dataOffset, int32_t dataLength);
// extern int32_t   v1_5_storageLoadLength(void* context, int32_t keyOffset, int32_t keyLength);
// extern int32_t   v1_5_storageLoadFromAddress(void* context, int32_t addressOffset, int32_t keyOffset, int32_t keyLength, int32_t dataOffset);
// extern int32_t   v1_5_storageLoad(void* context, int32_t keyOffset, int32_t keyLength, int32_t dataOffset);
// extern int32_t   v1_5_setStorageLock(void* context, int32_t keyOffset, int32_t keyLength, long long lockTimestamp);
// extern long long v1_5_getStorageLock(void* context, int32_t keyOffset, int32_t keyLength);
// extern int32_t   v1_5_isStorageLocked(void* context, int32_t keyOffset, int32_t keyLength);
// extern int32_t   v1_5_clearStorageLock(void* context, int32_t keyOffset, int32_t keyLength);
// extern void      v1_5_getCaller(void* context, int32_t resultOffset);
// extern void      v1_5_checkNoPayment(void* context);
// extern int32_t   v1_5_callValue(void* context, int32_t resultOffset);
// extern int32_t   v1_5_getESDTValue(void* context, int32_t resultOffset);
// extern int32_t   v1_5_getESDTValueByIndex(void* context, int32_t resultOffset, int32_t index);
// extern int32_t   v1_5_getESDTTokenName(void* context, int32_t resultOffset);
// extern int32_t   v1_5_getESDTTokenNameByIndex(void* context, int32_t resultOffset, int32_t index);
// extern long long v1_5_getESDTTokenNonce(void* context);
// extern long long v1_5_getESDTTokenNonceByIndex(void* context, int32_t index);
// extern long long v1_5_getCurrentESDTNFTNonce(void* context, int32_t addressOffset, int32_t tokenIDOffset, int32_t tokenIDLen);
// extern int32_t   v1_5_getESDTTokenType(void* context);
// extern int32_t   v1_5_getESDTTokenTypeByIndex(void* context, int32_t index);
// extern int32_t   v1_5_getNumESDTTransfers(void* context);
// extern int32_t   v1_5_getCallValueTokenName(void* context, int32_t callValueOffset, int32_t tokenNameOffset);
// extern int32_t   v1_5_getCallValueTokenNameByIndex(void* context, int32_t callValueOffset, int32_t tokenNameOffset, int32_t index);
// extern void      v1_5_writeLog(void* context, int32_t dataPointer, int32_t dataLength, int32_t topicPtr, int32_t numTopics);
// extern void      v1_5_writeEventLog(void* context, int32_t numTopics, int32_t topicLengthsOffset, int32_t topicOffset, int32_t dataOffset, int32_t dataLength);
// extern long long v1_5_getBlockTimestamp(void* context);
// extern long long v1_5_getBlockNonce(void* context);
// extern long long v1_5_getBlockRound(void* context);
// extern long long v1_5_getBlockEpoch(void* context);
// extern void      v1_5_getBlockRandomSeed(void* context, int32_t pointer);
// extern void      v1_5_getStateRootHash(void* context, int32_t pointer);
// extern long long v1_5_getPrevBlockTimestamp(void* context);
// extern long long v1_5_getPrevBlockNonce(void* context);
// extern long long v1_5_getPrevBlockRound(void* context);
// extern long long v1_5_getPrevBlockEpoch(void* context);
// extern void      v1_5_getPrevBlockRandomSeed(void* context, int32_t pointer);
// extern void      v1_5_returnData(void* context, int32_t pointer, int32_t length);
// extern int32_t   v1_5_executeOnSameContext(void* context, long long gasLimit, int32_t addressOffset, int32_t valueOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t   v1_5_executeOnDestContext(void* context, long long gasLimit, int32_t addressOffset, int32_t valueOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t   v1_5_executeReadOnly(void* context, long long gasLimit, int32_t addressOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t   v1_5_createContract(void* context, long long gasLimit, int32_t valueOffset, int32_t codeOffset, int32_t codeMetadataOffset, int32_t length, int32_t resultOffset, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t   v1_5_deployFromSourceContract(void* context, long long gasLimit, int32_t valueOffset, int32_t sourceContractAddressOffset, int32_t codeMetadataOffset, int32_t resultAddressOffset, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t   v1_5_getNumReturnData(void* context);
// extern int32_t   v1_5_getReturnDataSize(void* context, int32_t resultID);
// extern int32_t   v1_5_getReturnData(void* context, int32_t resultID, int32_t dataOffset);
// extern void      v1_5_cleanReturnData(void* context);
// extern void      v1_5_deleteFromReturnData(void* context, int32_t resultID);
// extern void      v1_5_getOriginalTxHash(void* context, int32_t dataOffset);
// extern void      v1_5_getCurrentTxHash(void* context, int32_t dataOffset);
// extern void      v1_5_getPrevTxHash(void* context, int32_t dataOffset);
// extern void      v1_5_managedSCAddress(void* context, int32_t destinationHandle);
// extern void      v1_5_managedOwnerAddress(void* context, int32_t destinationHandle);
// extern void      v1_5_managedCaller(void* context, int32_t destinationHandle);
// extern void      v1_5_managedSignalError(void* context, int32_t errHandle);
// extern void      v1_5_managedWriteLog(void* context, int32_t topicsHandle, int32_t dataHandle);
// extern void      v1_5_managedGetOriginalTxHash(void* context, int32_t resultHandle);
// extern void      v1_5_managedGetStateRootHash(void* context, int32_t resultHandle);
// extern void      v1_5_managedGetBlockRandomSeed(void* context, int32_t resultHandle);
// extern void      v1_5_managedGetPrevBlockRandomSeed(void* context, int32_t resultHandle);
// extern void      v1_5_managedGetReturnData(void* context, int32_t resultID, int32_t resultHandle);
// extern void      v1_5_managedGetMultiESDTCallValue(void* context, int32_t multiCallValueHandle);
// extern void      v1_5_managedGetESDTBalance(void* context, int32_t addressHandle, int32_t tokenIDHandle, long long nonce, int32_t valueHandle);
// extern void      v1_5_managedGetESDTTokenData(void* context, int32_t addressHandle, int32_t tokenIDHandle, long long nonce, int32_t valueHandle, int32_t propertiesHandle, int32_t hashHandle, int32_t nameHandle, int32_t attributesHandle, int32_t creatorHandle, int32_t royaltiesHandle, int32_t urisHandle);
// extern void      v1_5_managedAsyncCall(void* context, int32_t destHandle, int32_t valueHandle, int32_t functionHandle, int32_t argumentsHandle);
// extern int32_t   v1_5_managedCreateAsyncCall(void* context, int32_t destHandle, int32_t valueHandle, int32_t functionHandle, int32_t argumentsHandle, int32_t successOffset, int32_t successLength, int32_t errorOffset, int32_t errorLength, long long gas, long long extraGasForCallback, int32_t callbackClosureHandle);
// extern void      v1_5_managedGetCallbackClosure(void* context, int32_t callbackClosureHandle);
// extern void      v1_5_managedUpgradeFromSourceContract(void* context, int32_t destHandle, long long gas, int32_t valueHandle, int32_t addressHandle, int32_t codeMetadataHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern void      v1_5_managedUpgradeContract(void* context, int32_t destHandle, long long gas, int32_t valueHandle, int32_t codeHandle, int32_t codeMetadataHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern void      v1_5_managedDeleteContract(void* context, int32_t destHandle, long long gasLimit, int32_t argumentsHandle);
// extern int32_t   v1_5_managedDeployFromSourceContract(void* context, long long gas, int32_t valueHandle, int32_t addressHandle, int32_t codeMetadataHandle, int32_t argumentsHandle, int32_t resultAddressHandle, int32_t resultHandle);
// extern int32_t   v1_5_managedCreateContract(void* context, long long gas, int32_t valueHandle, int32_t codeHandle, int32_t codeMetadataHandle, int32_t argumentsHandle, int32_t resultAddressHandle, int32_t resultHandle);
// extern int32_t   v1_5_managedExecuteReadOnly(void* context, long long gas, int32_t addressHandle, int32_t functionHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern int32_t   v1_5_managedExecuteOnSameContext(void* context, long long gas, int32_t addressHandle, int32_t valueHandle, int32_t functionHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern int32_t   v1_5_managedExecuteOnDestContext(void* context, long long gas, int32_t addressHandle, int32_t valueHandle, int32_t functionHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern int32_t   v1_5_managedMultiTransferESDTNFTExecute(void* context, int32_t dstHandle, int32_t tokenTransfersHandle, long long gasLimit, int32_t functionHandle, int32_t argumentsHandle);
// extern int32_t   v1_5_managedTransferValueExecute(void* context, int32_t dstHandle, int32_t valueHandle, long long gasLimit, int32_t functionHandle, int32_t argumentsHandle);
// extern int32_t   v1_5_managedIsESDTFrozen(void* context, int32_t addressHandle, int32_t tokenIDHandle, long long nonce);
// extern int32_t   v1_5_managedIsESDTLimitedTransfer(void* context, int32_t tokenIDHandle);
// extern int32_t   v1_5_managedIsESDTPaused(void* context, int32_t tokenIDHandle);
// extern void      v1_5_managedBufferToHex(void* context, int32_t sourceHandle, int32_t destHandle);
// extern int32_t   v1_5_mBufferNew(void* context);
// extern int32_t   v1_5_mBufferNewFromBytes(void* context, int32_t dataOffset, int32_t dataLength);
// extern int32_t   v1_5_mBufferGetLength(void* context, int32_t mBufferHandle);
// extern int32_t   v1_5_mBufferGetBytes(void* context, int32_t mBufferHandle, int32_t resultOffset);
// extern int32_t   v1_5_mBufferGetByteSlice(void* context, int32_t sourceHandle, int32_t startingPosition, int32_t sliceLength, int32_t resultOffset);
// extern int32_t   v1_5_mBufferCopyByteSlice(void* context, int32_t sourceHandle, int32_t startingPosition, int32_t sliceLength, int32_t destinationHandle);
// extern int32_t   v1_5_mBufferEq(void* context, int32_t mBufferHandle1, int32_t mBufferHandle2);
// extern int32_t   v1_5_mBufferSetBytes(void* context, int32_t mBufferHandle, int32_t dataOffset, int32_t dataLength);
// extern int32_t   v1_5_mBufferSetByteSlice(void* context, int32_t mBufferHandle, int32_t startingPosition, int32_t dataLength, int32_t dataOffset);
// extern int32_t   v1_5_mBufferAppend(void* context, int32_t accumulatorHandle, int32_t dataHandle);
// extern int32_t   v1_5_mBufferAppendBytes(void* context, int32_t accumulatorHandle, int32_t dataOffset, int32_t dataLength);
// extern int32_t   v1_5_mBufferToBigIntUnsigned(void* context, int32_t mBufferHandle, int32_t bigIntHandle);
// extern int32_t   v1_5_mBufferToBigIntSigned(void* context, int32_t mBufferHandle, int32_t bigIntHandle);
// extern int32_t   v1_5_mBufferFromBigIntUnsigned(void* context, int32_t mBufferHandle, int32_t bigIntHandle);
// extern int32_t   v1_5_mBufferFromBigIntSigned(void* context, int32_t mBufferHandle, int32_t bigIntHandle);
// extern int32_t   v1_5_mBufferToBigFloat(void* context, int32_t mBufferHandle, int32_t bigFloatHandle);
// extern int32_t   v1_5_mBufferFromBigFloat(void* context, int32_t mBufferHandle, int32_t bigFloatHandle);
// extern int32_t   v1_5_mBufferStorageStore(void* context, int32_t keyHandle, int32_t sourceHandle);
// extern int32_t   v1_5_mBufferStorageLoad(void* context, int32_t keyHandle, int32_t destinationHandle);
// extern void      v1_5_mBufferStorageLoadFromAddress(void* context, int32_t addressHandle, int32_t keyHandle, int32_t destinationHandle);
// extern int32_t   v1_5_mBufferGetArgument(void* context, int32_t id, int32_t destinationHandle);
// extern int32_t   v1_5_mBufferFinish(void* context, int32_t sourceHandle);
// extern int32_t   v1_5_mBufferSetRandom(void* context, int32_t destinationHandle, int32_t length);
// extern long long v1_5_smallIntGetUnsignedArgument(void* context, int32_t id);
// extern long long v1_5_smallIntGetSignedArgument(void* context, int32_t id);
// extern void      v1_5_smallIntFinishUnsigned(void* context, long long value);
// extern void      v1_5_smallIntFinishSigned(void* context, long long value);
// extern int32_t   v1_5_smallIntStorageStoreUnsigned(void* context, int32_t keyOffset, int32_t keyLength, long long value);
// extern int32_t   v1_5_smallIntStorageStoreSigned(void* context, int32_t keyOffset, int32_t keyLength, long long value);
// extern long long v1_5_smallIntStorageLoadUnsigned(void* context, int32_t keyOffset, int32_t keyLength);
// extern long long v1_5_smallIntStorageLoadSigned(void* context, int32_t keyOffset, int32_t keyLength);
// extern long long v1_5_int64getArgument(void* context, int32_t id);
// extern void      v1_5_int64finish(void* context, long long value);
// extern int32_t   v1_5_int64storageStore(void* context, int32_t keyOffset, int32_t keyLength, long long value);
// extern long long v1_5_int64storageLoad(void* context, int32_t keyOffset, int32_t keyLength);
// extern int32_t   v1_5_sha256(void* context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t   v1_5_managedSha256(void* context, int32_t inputHandle, int32_t outputHandle);
// extern int32_t   v1_5_keccak256(void* context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t   v1_5_managedKeccak256(void* context, int32_t inputHandle, int32_t outputHandle);
// extern int32_t   v1_5_ripemd160(void* context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t   v1_5_managedRipemd160(void* context, int32_t inputHandle, int32_t outputHandle);
// extern int32_t   v1_5_verifyBLS(void* context, int32_t keyOffset, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern int32_t   v1_5_managedVerifyBLS(void* context, int32_t keyHandle, int32_t messageHandle, int32_t sigHandle);
// extern int32_t   v1_5_verifyEd25519(void* context, int32_t keyOffset, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern int32_t   v1_5_managedVerifyEd25519(void* context, int32_t keyHandle, int32_t messageHandle, int32_t sigHandle);
// extern int32_t   v1_5_verifyCustomSecp256k1(void* context, int32_t keyOffset, int32_t keyLength, int32_t messageOffset, int32_t messageLength, int32_t sigOffset, int32_t hashType);
// extern int32_t   v1_5_managedVerifyCustomSecp256k1(void* context, int32_t keyHandle, int32_t messageHandle, int32_t sigHandle, int32_t hashType);
// extern int32_t   v1_5_verifySecp256k1(void* context, int32_t keyOffset, int32_t keyLength, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern int32_t   v1_5_managedVerifySecp256k1(void* context, int32_t keyHandle, int32_t messageHandle, int32_t sigHandle);
// extern int32_t   v1_5_encodeSecp256k1DerSignature(void* context, int32_t rOffset, int32_t rLength, int32_t sOffset, int32_t sLength, int32_t sigOffset);
// extern int32_t   v1_5_managedEncodeSecp256k1DerSignature(void* context, int32_t rHandle, int32_t sHandle, int32_t sigHandle);
// extern void      v1_5_addEC(void* context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t fstPointXHandle, int32_t fstPointYHandle, int32_t sndPointXHandle, int32_t sndPointYHandle);
// extern void      v1_5_doubleEC(void* context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t pointXHandle, int32_t pointYHandle);
// extern int32_t   v1_5_isOnCurveEC(void* context, int32_t ecHandle, int32_t pointXHandle, int32_t pointYHandle);
// extern int32_t   v1_5_scalarBaseMultEC(void* context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t dataOffset, int32_t length);
// extern int32_t   v1_5_managedScalarBaseMultEC(void* context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t dataHandle);
// extern int32_t   v1_5_scalarMultEC(void* context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t pointXHandle, int32_t pointYHandle, int32_t dataOffset, int32_t length);
// extern int32_t   v1_5_managedScalarMultEC(void* context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t pointXHandle, int32_t pointYHandle, int32_t dataHandle);
// extern int32_t   v1_5_marshalEC(void* context, int32_t xPairHandle, int32_t yPairHandle, int32_t ecHandle, int32_t resultOffset);
// extern int32_t   v1_5_managedMarshalEC(void* context, int32_t xPairHandle, int32_t yPairHandle, int32_t ecHandle, int32_t resultHandle);
// extern int32_t   v1_5_marshalCompressedEC(void* context, int32_t xPairHandle, int32_t yPairHandle, int32_t ecHandle, int32_t resultOffset);
// extern int32_t   v1_5_managedMarshalCompressedEC(void* context, int32_t xPairHandle, int32_t yPairHandle, int32_t ecHandle, int32_t resultHandle);
// extern int32_t   v1_5_unmarshalEC(void* context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t dataOffset, int32_t length);
// extern int32_t   v1_5_managedUnmarshalEC(void* context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t dataHandle);
// extern int32_t   v1_5_unmarshalCompressedEC(void* context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t dataOffset, int32_t length);
// extern int32_t   v1_5_managedUnmarshalCompressedEC(void* context, int32_t xResultHandle, int32_t yResultHandle, int32_t ecHandle, int32_t dataHandle);
// extern int32_t   v1_5_generateKeyEC(void* context, int32_t xPubKeyHandle, int32_t yPubKeyHandle, int32_t ecHandle, int32_t resultOffset);
// extern int32_t   v1_5_managedGenerateKeyEC(void* context, int32_t xPubKeyHandle, int32_t yPubKeyHandle, int32_t ecHandle, int32_t resultHandle);
// extern int32_t   v1_5_createEC(void* context, int32_t dataOffset, int32_t dataLength);
// extern int32_t   v1_5_managedCreateEC(void* context, int32_t dataHandle);
// extern int32_t   v1_5_getCurveLengthEC(void* context, int32_t ecHandle);
// extern int32_t   v1_5_getPrivKeyByteLengthEC(void* context, int32_t ecHandle);
// extern int32_t   v1_5_ellipticCurveGetValues(void* context, int32_t ecHandle, int32_t fieldOrderHandle, int32_t basePointOrderHandle, int32_t eqConstantHandle, int32_t xBasePointHandle, int32_t yBasePointHandle);
import "C"

import (
	"unsafe"

	"github.com/ElrondNetwork/wasm-vm/executor"
)

// ElrondEIImports populates imports with the ElrondEI API methods
func ElrondEIImports(imports executor.ImportFunctionReceiver) error {
	imports.Namespace("env")

	var err error
	err = imports.Append("bigFloatNewFromParts", v1_5_bigFloatNewFromParts, C.v1_5_bigFloatNewFromParts)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatNewFromFrac", v1_5_bigFloatNewFromFrac, C.v1_5_bigFloatNewFromFrac)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatNewFromSci", v1_5_bigFloatNewFromSci, C.v1_5_bigFloatNewFromSci)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatAdd", v1_5_bigFloatAdd, C.v1_5_bigFloatAdd)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatSub", v1_5_bigFloatSub, C.v1_5_bigFloatSub)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatMul", v1_5_bigFloatMul, C.v1_5_bigFloatMul)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatDiv", v1_5_bigFloatDiv, C.v1_5_bigFloatDiv)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatNeg", v1_5_bigFloatNeg, C.v1_5_bigFloatNeg)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatClone", v1_5_bigFloatClone, C.v1_5_bigFloatClone)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatCmp", v1_5_bigFloatCmp, C.v1_5_bigFloatCmp)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatAbs", v1_5_bigFloatAbs, C.v1_5_bigFloatAbs)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatSign", v1_5_bigFloatSign, C.v1_5_bigFloatSign)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatSqrt", v1_5_bigFloatSqrt, C.v1_5_bigFloatSqrt)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatPow", v1_5_bigFloatPow, C.v1_5_bigFloatPow)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatFloor", v1_5_bigFloatFloor, C.v1_5_bigFloatFloor)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatCeil", v1_5_bigFloatCeil, C.v1_5_bigFloatCeil)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatTruncate", v1_5_bigFloatTruncate, C.v1_5_bigFloatTruncate)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatSetInt64", v1_5_bigFloatSetInt64, C.v1_5_bigFloatSetInt64)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatIsInt", v1_5_bigFloatIsInt, C.v1_5_bigFloatIsInt)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatSetBigInt", v1_5_bigFloatSetBigInt, C.v1_5_bigFloatSetBigInt)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatGetConstPi", v1_5_bigFloatGetConstPi, C.v1_5_bigFloatGetConstPi)
	if err != nil {
		return err
	}

	err = imports.Append("bigFloatGetConstE", v1_5_bigFloatGetConstE, C.v1_5_bigFloatGetConstE)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntGetUnsignedArgument", v1_5_bigIntGetUnsignedArgument, C.v1_5_bigIntGetUnsignedArgument)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntGetSignedArgument", v1_5_bigIntGetSignedArgument, C.v1_5_bigIntGetSignedArgument)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntStorageStoreUnsigned", v1_5_bigIntStorageStoreUnsigned, C.v1_5_bigIntStorageStoreUnsigned)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntStorageLoadUnsigned", v1_5_bigIntStorageLoadUnsigned, C.v1_5_bigIntStorageLoadUnsigned)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntGetCallValue", v1_5_bigIntGetCallValue, C.v1_5_bigIntGetCallValue)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntGetESDTCallValue", v1_5_bigIntGetESDTCallValue, C.v1_5_bigIntGetESDTCallValue)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntGetESDTCallValueByIndex", v1_5_bigIntGetESDTCallValueByIndex, C.v1_5_bigIntGetESDTCallValueByIndex)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntGetExternalBalance", v1_5_bigIntGetExternalBalance, C.v1_5_bigIntGetExternalBalance)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntGetESDTExternalBalance", v1_5_bigIntGetESDTExternalBalance, C.v1_5_bigIntGetESDTExternalBalance)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntNew", v1_5_bigIntNew, C.v1_5_bigIntNew)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntUnsignedByteLength", v1_5_bigIntUnsignedByteLength, C.v1_5_bigIntUnsignedByteLength)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntSignedByteLength", v1_5_bigIntSignedByteLength, C.v1_5_bigIntSignedByteLength)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntGetUnsignedBytes", v1_5_bigIntGetUnsignedBytes, C.v1_5_bigIntGetUnsignedBytes)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntGetSignedBytes", v1_5_bigIntGetSignedBytes, C.v1_5_bigIntGetSignedBytes)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntSetUnsignedBytes", v1_5_bigIntSetUnsignedBytes, C.v1_5_bigIntSetUnsignedBytes)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntSetSignedBytes", v1_5_bigIntSetSignedBytes, C.v1_5_bigIntSetSignedBytes)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntIsInt64", v1_5_bigIntIsInt64, C.v1_5_bigIntIsInt64)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntGetInt64", v1_5_bigIntGetInt64, C.v1_5_bigIntGetInt64)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntSetInt64", v1_5_bigIntSetInt64, C.v1_5_bigIntSetInt64)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntAdd", v1_5_bigIntAdd, C.v1_5_bigIntAdd)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntSub", v1_5_bigIntSub, C.v1_5_bigIntSub)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntMul", v1_5_bigIntMul, C.v1_5_bigIntMul)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntTDiv", v1_5_bigIntTDiv, C.v1_5_bigIntTDiv)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntTMod", v1_5_bigIntTMod, C.v1_5_bigIntTMod)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntEDiv", v1_5_bigIntEDiv, C.v1_5_bigIntEDiv)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntEMod", v1_5_bigIntEMod, C.v1_5_bigIntEMod)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntSqrt", v1_5_bigIntSqrt, C.v1_5_bigIntSqrt)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntPow", v1_5_bigIntPow, C.v1_5_bigIntPow)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntLog2", v1_5_bigIntLog2, C.v1_5_bigIntLog2)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntAbs", v1_5_bigIntAbs, C.v1_5_bigIntAbs)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntNeg", v1_5_bigIntNeg, C.v1_5_bigIntNeg)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntSign", v1_5_bigIntSign, C.v1_5_bigIntSign)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntCmp", v1_5_bigIntCmp, C.v1_5_bigIntCmp)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntNot", v1_5_bigIntNot, C.v1_5_bigIntNot)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntAnd", v1_5_bigIntAnd, C.v1_5_bigIntAnd)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntOr", v1_5_bigIntOr, C.v1_5_bigIntOr)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntXor", v1_5_bigIntXor, C.v1_5_bigIntXor)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntShr", v1_5_bigIntShr, C.v1_5_bigIntShr)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntShl", v1_5_bigIntShl, C.v1_5_bigIntShl)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntFinishUnsigned", v1_5_bigIntFinishUnsigned, C.v1_5_bigIntFinishUnsigned)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntFinishSigned", v1_5_bigIntFinishSigned, C.v1_5_bigIntFinishSigned)
	if err != nil {
		return err
	}

	err = imports.Append("bigIntToString", v1_5_bigIntToString, C.v1_5_bigIntToString)
	if err != nil {
		return err
	}

	err = imports.Append("getGasLeft", v1_5_getGasLeft, C.v1_5_getGasLeft)
	if err != nil {
		return err
	}

	err = imports.Append("getSCAddress", v1_5_getSCAddress, C.v1_5_getSCAddress)
	if err != nil {
		return err
	}

	err = imports.Append("getOwnerAddress", v1_5_getOwnerAddress, C.v1_5_getOwnerAddress)
	if err != nil {
		return err
	}

	err = imports.Append("getShardOfAddress", v1_5_getShardOfAddress, C.v1_5_getShardOfAddress)
	if err != nil {
		return err
	}

	err = imports.Append("isSmartContract", v1_5_isSmartContract, C.v1_5_isSmartContract)
	if err != nil {
		return err
	}

	err = imports.Append("signalError", v1_5_signalError, C.v1_5_signalError)
	if err != nil {
		return err
	}

	err = imports.Append("getExternalBalance", v1_5_getExternalBalance, C.v1_5_getExternalBalance)
	if err != nil {
		return err
	}

	err = imports.Append("blockHash", v1_5_blockHash, C.v1_5_blockHash)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTBalance", v1_5_getESDTBalance, C.v1_5_getESDTBalance)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTNFTNameLength", v1_5_getESDTNFTNameLength, C.v1_5_getESDTNFTNameLength)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTNFTAttributeLength", v1_5_getESDTNFTAttributeLength, C.v1_5_getESDTNFTAttributeLength)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTNFTURILength", v1_5_getESDTNFTURILength, C.v1_5_getESDTNFTURILength)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTTokenData", v1_5_getESDTTokenData, C.v1_5_getESDTTokenData)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTLocalRoles", v1_5_getESDTLocalRoles, C.v1_5_getESDTLocalRoles)
	if err != nil {
		return err
	}

	err = imports.Append("validateTokenIdentifier", v1_5_validateTokenIdentifier, C.v1_5_validateTokenIdentifier)
	if err != nil {
		return err
	}

	err = imports.Append("transferValue", v1_5_transferValue, C.v1_5_transferValue)
	if err != nil {
		return err
	}

	err = imports.Append("transferValueExecute", v1_5_transferValueExecute, C.v1_5_transferValueExecute)
	if err != nil {
		return err
	}

	err = imports.Append("transferESDTExecute", v1_5_transferESDTExecute, C.v1_5_transferESDTExecute)
	if err != nil {
		return err
	}

	err = imports.Append("transferESDTNFTExecute", v1_5_transferESDTNFTExecute, C.v1_5_transferESDTNFTExecute)
	if err != nil {
		return err
	}

	err = imports.Append("multiTransferESDTNFTExecute", v1_5_multiTransferESDTNFTExecute, C.v1_5_multiTransferESDTNFTExecute)
	if err != nil {
		return err
	}

	err = imports.Append("createAsyncCall", v1_5_createAsyncCall, C.v1_5_createAsyncCall)
	if err != nil {
		return err
	}

	err = imports.Append("setAsyncContextCallback", v1_5_setAsyncContextCallback, C.v1_5_setAsyncContextCallback)
	if err != nil {
		return err
	}

	err = imports.Append("upgradeContract", v1_5_upgradeContract, C.v1_5_upgradeContract)
	if err != nil {
		return err
	}

	err = imports.Append("upgradeFromSourceContract", v1_5_upgradeFromSourceContract, C.v1_5_upgradeFromSourceContract)
	if err != nil {
		return err
	}

	err = imports.Append("deleteContract", v1_5_deleteContract, C.v1_5_deleteContract)
	if err != nil {
		return err
	}

	err = imports.Append("asyncCall", v1_5_asyncCall, C.v1_5_asyncCall)
	if err != nil {
		return err
	}

	err = imports.Append("getArgumentLength", v1_5_getArgumentLength, C.v1_5_getArgumentLength)
	if err != nil {
		return err
	}

	err = imports.Append("getArgument", v1_5_getArgument, C.v1_5_getArgument)
	if err != nil {
		return err
	}

	err = imports.Append("getFunction", v1_5_getFunction, C.v1_5_getFunction)
	if err != nil {
		return err
	}

	err = imports.Append("getNumArguments", v1_5_getNumArguments, C.v1_5_getNumArguments)
	if err != nil {
		return err
	}

	err = imports.Append("storageStore", v1_5_storageStore, C.v1_5_storageStore)
	if err != nil {
		return err
	}

	err = imports.Append("storageLoadLength", v1_5_storageLoadLength, C.v1_5_storageLoadLength)
	if err != nil {
		return err
	}

	err = imports.Append("storageLoadFromAddress", v1_5_storageLoadFromAddress, C.v1_5_storageLoadFromAddress)
	if err != nil {
		return err
	}

	err = imports.Append("storageLoad", v1_5_storageLoad, C.v1_5_storageLoad)
	if err != nil {
		return err
	}

	err = imports.Append("setStorageLock", v1_5_setStorageLock, C.v1_5_setStorageLock)
	if err != nil {
		return err
	}

	err = imports.Append("getStorageLock", v1_5_getStorageLock, C.v1_5_getStorageLock)
	if err != nil {
		return err
	}

	err = imports.Append("isStorageLocked", v1_5_isStorageLocked, C.v1_5_isStorageLocked)
	if err != nil {
		return err
	}

	err = imports.Append("clearStorageLock", v1_5_clearStorageLock, C.v1_5_clearStorageLock)
	if err != nil {
		return err
	}

	err = imports.Append("getCaller", v1_5_getCaller, C.v1_5_getCaller)
	if err != nil {
		return err
	}

	err = imports.Append("checkNoPayment", v1_5_checkNoPayment, C.v1_5_checkNoPayment)
	if err != nil {
		return err
	}

	err = imports.Append("callValue", v1_5_callValue, C.v1_5_callValue)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTValue", v1_5_getESDTValue, C.v1_5_getESDTValue)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTValueByIndex", v1_5_getESDTValueByIndex, C.v1_5_getESDTValueByIndex)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTTokenName", v1_5_getESDTTokenName, C.v1_5_getESDTTokenName)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTTokenNameByIndex", v1_5_getESDTTokenNameByIndex, C.v1_5_getESDTTokenNameByIndex)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTTokenNonce", v1_5_getESDTTokenNonce, C.v1_5_getESDTTokenNonce)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTTokenNonceByIndex", v1_5_getESDTTokenNonceByIndex, C.v1_5_getESDTTokenNonceByIndex)
	if err != nil {
		return err
	}

	err = imports.Append("getCurrentESDTNFTNonce", v1_5_getCurrentESDTNFTNonce, C.v1_5_getCurrentESDTNFTNonce)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTTokenType", v1_5_getESDTTokenType, C.v1_5_getESDTTokenType)
	if err != nil {
		return err
	}

	err = imports.Append("getESDTTokenTypeByIndex", v1_5_getESDTTokenTypeByIndex, C.v1_5_getESDTTokenTypeByIndex)
	if err != nil {
		return err
	}

	err = imports.Append("getNumESDTTransfers", v1_5_getNumESDTTransfers, C.v1_5_getNumESDTTransfers)
	if err != nil {
		return err
	}

	err = imports.Append("getCallValueTokenName", v1_5_getCallValueTokenName, C.v1_5_getCallValueTokenName)
	if err != nil {
		return err
	}

	err = imports.Append("getCallValueTokenNameByIndex", v1_5_getCallValueTokenNameByIndex, C.v1_5_getCallValueTokenNameByIndex)
	if err != nil {
		return err
	}

	err = imports.Append("writeLog", v1_5_writeLog, C.v1_5_writeLog)
	if err != nil {
		return err
	}

	err = imports.Append("writeEventLog", v1_5_writeEventLog, C.v1_5_writeEventLog)
	if err != nil {
		return err
	}

	err = imports.Append("getBlockTimestamp", v1_5_getBlockTimestamp, C.v1_5_getBlockTimestamp)
	if err != nil {
		return err
	}

	err = imports.Append("getBlockNonce", v1_5_getBlockNonce, C.v1_5_getBlockNonce)
	if err != nil {
		return err
	}

	err = imports.Append("getBlockRound", v1_5_getBlockRound, C.v1_5_getBlockRound)
	if err != nil {
		return err
	}

	err = imports.Append("getBlockEpoch", v1_5_getBlockEpoch, C.v1_5_getBlockEpoch)
	if err != nil {
		return err
	}

	err = imports.Append("getBlockRandomSeed", v1_5_getBlockRandomSeed, C.v1_5_getBlockRandomSeed)
	if err != nil {
		return err
	}

	err = imports.Append("getStateRootHash", v1_5_getStateRootHash, C.v1_5_getStateRootHash)
	if err != nil {
		return err
	}

	err = imports.Append("getPrevBlockTimestamp", v1_5_getPrevBlockTimestamp, C.v1_5_getPrevBlockTimestamp)
	if err != nil {
		return err
	}

	err = imports.Append("getPrevBlockNonce", v1_5_getPrevBlockNonce, C.v1_5_getPrevBlockNonce)
	if err != nil {
		return err
	}

	err = imports.Append("getPrevBlockRound", v1_5_getPrevBlockRound, C.v1_5_getPrevBlockRound)
	if err != nil {
		return err
	}

	err = imports.Append("getPrevBlockEpoch", v1_5_getPrevBlockEpoch, C.v1_5_getPrevBlockEpoch)
	if err != nil {
		return err
	}

	err = imports.Append("getPrevBlockRandomSeed", v1_5_getPrevBlockRandomSeed, C.v1_5_getPrevBlockRandomSeed)
	if err != nil {
		return err
	}

	err = imports.Append("returnData", v1_5_returnData, C.v1_5_returnData)
	if err != nil {
		return err
	}

	err = imports.Append("executeOnSameContext", v1_5_executeOnSameContext, C.v1_5_executeOnSameContext)
	if err != nil {
		return err
	}

	err = imports.Append("executeOnDestContext", v1_5_executeOnDestContext, C.v1_5_executeOnDestContext)
	if err != nil {
		return err
	}

	err = imports.Append("executeReadOnly", v1_5_executeReadOnly, C.v1_5_executeReadOnly)
	if err != nil {
		return err
	}

	err = imports.Append("createContract", v1_5_createContract, C.v1_5_createContract)
	if err != nil {
		return err
	}

	err = imports.Append("deployFromSourceContract", v1_5_deployFromSourceContract, C.v1_5_deployFromSourceContract)
	if err != nil {
		return err
	}

	err = imports.Append("getNumReturnData", v1_5_getNumReturnData, C.v1_5_getNumReturnData)
	if err != nil {
		return err
	}

	err = imports.Append("getReturnDataSize", v1_5_getReturnDataSize, C.v1_5_getReturnDataSize)
	if err != nil {
		return err
	}

	err = imports.Append("getReturnData", v1_5_getReturnData, C.v1_5_getReturnData)
	if err != nil {
		return err
	}

	err = imports.Append("cleanReturnData", v1_5_cleanReturnData, C.v1_5_cleanReturnData)
	if err != nil {
		return err
	}

	err = imports.Append("deleteFromReturnData", v1_5_deleteFromReturnData, C.v1_5_deleteFromReturnData)
	if err != nil {
		return err
	}

	err = imports.Append("getOriginalTxHash", v1_5_getOriginalTxHash, C.v1_5_getOriginalTxHash)
	if err != nil {
		return err
	}

	err = imports.Append("getCurrentTxHash", v1_5_getCurrentTxHash, C.v1_5_getCurrentTxHash)
	if err != nil {
		return err
	}

	err = imports.Append("getPrevTxHash", v1_5_getPrevTxHash, C.v1_5_getPrevTxHash)
	if err != nil {
		return err
	}

	err = imports.Append("managedSCAddress", v1_5_managedSCAddress, C.v1_5_managedSCAddress)
	if err != nil {
		return err
	}

	err = imports.Append("managedOwnerAddress", v1_5_managedOwnerAddress, C.v1_5_managedOwnerAddress)
	if err != nil {
		return err
	}

	err = imports.Append("managedCaller", v1_5_managedCaller, C.v1_5_managedCaller)
	if err != nil {
		return err
	}

	err = imports.Append("managedSignalError", v1_5_managedSignalError, C.v1_5_managedSignalError)
	if err != nil {
		return err
	}

	err = imports.Append("managedWriteLog", v1_5_managedWriteLog, C.v1_5_managedWriteLog)
	if err != nil {
		return err
	}

	err = imports.Append("managedGetOriginalTxHash", v1_5_managedGetOriginalTxHash, C.v1_5_managedGetOriginalTxHash)
	if err != nil {
		return err
	}

	err = imports.Append("managedGetStateRootHash", v1_5_managedGetStateRootHash, C.v1_5_managedGetStateRootHash)
	if err != nil {
		return err
	}

	err = imports.Append("managedGetBlockRandomSeed", v1_5_managedGetBlockRandomSeed, C.v1_5_managedGetBlockRandomSeed)
	if err != nil {
		return err
	}

	err = imports.Append("managedGetPrevBlockRandomSeed", v1_5_managedGetPrevBlockRandomSeed, C.v1_5_managedGetPrevBlockRandomSeed)
	if err != nil {
		return err
	}

	err = imports.Append("managedGetReturnData", v1_5_managedGetReturnData, C.v1_5_managedGetReturnData)
	if err != nil {
		return err
	}

	err = imports.Append("managedGetMultiESDTCallValue", v1_5_managedGetMultiESDTCallValue, C.v1_5_managedGetMultiESDTCallValue)
	if err != nil {
		return err
	}

	err = imports.Append("managedGetESDTBalance", v1_5_managedGetESDTBalance, C.v1_5_managedGetESDTBalance)
	if err != nil {
		return err
	}

	err = imports.Append("managedGetESDTTokenData", v1_5_managedGetESDTTokenData, C.v1_5_managedGetESDTTokenData)
	if err != nil {
		return err
	}

	err = imports.Append("managedAsyncCall", v1_5_managedAsyncCall, C.v1_5_managedAsyncCall)
	if err != nil {
		return err
	}

	err = imports.Append("managedCreateAsyncCall", v1_5_managedCreateAsyncCall, C.v1_5_managedCreateAsyncCall)
	if err != nil {
		return err
	}

	err = imports.Append("managedGetCallbackClosure", v1_5_managedGetCallbackClosure, C.v1_5_managedGetCallbackClosure)
	if err != nil {
		return err
	}

	err = imports.Append("managedUpgradeFromSourceContract", v1_5_managedUpgradeFromSourceContract, C.v1_5_managedUpgradeFromSourceContract)
	if err != nil {
		return err
	}

	err = imports.Append("managedUpgradeContract", v1_5_managedUpgradeContract, C.v1_5_managedUpgradeContract)
	if err != nil {
		return err
	}

	err = imports.Append("managedDeleteContract", v1_5_managedDeleteContract, C.v1_5_managedDeleteContract)
	if err != nil {
		return err
	}

	err = imports.Append("managedDeployFromSourceContract", v1_5_managedDeployFromSourceContract, C.v1_5_managedDeployFromSourceContract)
	if err != nil {
		return err
	}

	err = imports.Append("managedCreateContract", v1_5_managedCreateContract, C.v1_5_managedCreateContract)
	if err != nil {
		return err
	}

	err = imports.Append("managedExecuteReadOnly", v1_5_managedExecuteReadOnly, C.v1_5_managedExecuteReadOnly)
	if err != nil {
		return err
	}

	err = imports.Append("managedExecuteOnSameContext", v1_5_managedExecuteOnSameContext, C.v1_5_managedExecuteOnSameContext)
	if err != nil {
		return err
	}

	err = imports.Append("managedExecuteOnDestContext", v1_5_managedExecuteOnDestContext, C.v1_5_managedExecuteOnDestContext)
	if err != nil {
		return err
	}

	err = imports.Append("managedMultiTransferESDTNFTExecute", v1_5_managedMultiTransferESDTNFTExecute, C.v1_5_managedMultiTransferESDTNFTExecute)
	if err != nil {
		return err
	}

	err = imports.Append("managedTransferValueExecute", v1_5_managedTransferValueExecute, C.v1_5_managedTransferValueExecute)
	if err != nil {
		return err
	}

	err = imports.Append("managedIsESDTFrozen", v1_5_managedIsESDTFrozen, C.v1_5_managedIsESDTFrozen)
	if err != nil {
		return err
	}

	err = imports.Append("managedIsESDTLimitedTransfer", v1_5_managedIsESDTLimitedTransfer, C.v1_5_managedIsESDTLimitedTransfer)
	if err != nil {
		return err
	}

	err = imports.Append("managedIsESDTPaused", v1_5_managedIsESDTPaused, C.v1_5_managedIsESDTPaused)
	if err != nil {
		return err
	}

	err = imports.Append("managedBufferToHex", v1_5_managedBufferToHex, C.v1_5_managedBufferToHex)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferNew", v1_5_mBufferNew, C.v1_5_mBufferNew)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferNewFromBytes", v1_5_mBufferNewFromBytes, C.v1_5_mBufferNewFromBytes)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferGetLength", v1_5_mBufferGetLength, C.v1_5_mBufferGetLength)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferGetBytes", v1_5_mBufferGetBytes, C.v1_5_mBufferGetBytes)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferGetByteSlice", v1_5_mBufferGetByteSlice, C.v1_5_mBufferGetByteSlice)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferCopyByteSlice", v1_5_mBufferCopyByteSlice, C.v1_5_mBufferCopyByteSlice)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferEq", v1_5_mBufferEq, C.v1_5_mBufferEq)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferSetBytes", v1_5_mBufferSetBytes, C.v1_5_mBufferSetBytes)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferSetByteSlice", v1_5_mBufferSetByteSlice, C.v1_5_mBufferSetByteSlice)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferAppend", v1_5_mBufferAppend, C.v1_5_mBufferAppend)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferAppendBytes", v1_5_mBufferAppendBytes, C.v1_5_mBufferAppendBytes)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferToBigIntUnsigned", v1_5_mBufferToBigIntUnsigned, C.v1_5_mBufferToBigIntUnsigned)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferToBigIntSigned", v1_5_mBufferToBigIntSigned, C.v1_5_mBufferToBigIntSigned)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferFromBigIntUnsigned", v1_5_mBufferFromBigIntUnsigned, C.v1_5_mBufferFromBigIntUnsigned)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferFromBigIntSigned", v1_5_mBufferFromBigIntSigned, C.v1_5_mBufferFromBigIntSigned)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferToBigFloat", v1_5_mBufferToBigFloat, C.v1_5_mBufferToBigFloat)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferFromBigFloat", v1_5_mBufferFromBigFloat, C.v1_5_mBufferFromBigFloat)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferStorageStore", v1_5_mBufferStorageStore, C.v1_5_mBufferStorageStore)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferStorageLoad", v1_5_mBufferStorageLoad, C.v1_5_mBufferStorageLoad)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferStorageLoadFromAddress", v1_5_mBufferStorageLoadFromAddress, C.v1_5_mBufferStorageLoadFromAddress)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferGetArgument", v1_5_mBufferGetArgument, C.v1_5_mBufferGetArgument)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferFinish", v1_5_mBufferFinish, C.v1_5_mBufferFinish)
	if err != nil {
		return err
	}

	err = imports.Append("mBufferSetRandom", v1_5_mBufferSetRandom, C.v1_5_mBufferSetRandom)
	if err != nil {
		return err
	}

	err = imports.Append("smallIntGetUnsignedArgument", v1_5_smallIntGetUnsignedArgument, C.v1_5_smallIntGetUnsignedArgument)
	if err != nil {
		return err
	}

	err = imports.Append("smallIntGetSignedArgument", v1_5_smallIntGetSignedArgument, C.v1_5_smallIntGetSignedArgument)
	if err != nil {
		return err
	}

	err = imports.Append("smallIntFinishUnsigned", v1_5_smallIntFinishUnsigned, C.v1_5_smallIntFinishUnsigned)
	if err != nil {
		return err
	}

	err = imports.Append("smallIntFinishSigned", v1_5_smallIntFinishSigned, C.v1_5_smallIntFinishSigned)
	if err != nil {
		return err
	}

	err = imports.Append("smallIntStorageStoreUnsigned", v1_5_smallIntStorageStoreUnsigned, C.v1_5_smallIntStorageStoreUnsigned)
	if err != nil {
		return err
	}

	err = imports.Append("smallIntStorageStoreSigned", v1_5_smallIntStorageStoreSigned, C.v1_5_smallIntStorageStoreSigned)
	if err != nil {
		return err
	}

	err = imports.Append("smallIntStorageLoadUnsigned", v1_5_smallIntStorageLoadUnsigned, C.v1_5_smallIntStorageLoadUnsigned)
	if err != nil {
		return err
	}

	err = imports.Append("smallIntStorageLoadSigned", v1_5_smallIntStorageLoadSigned, C.v1_5_smallIntStorageLoadSigned)
	if err != nil {
		return err
	}

	err = imports.Append("int64getArgument", v1_5_int64getArgument, C.v1_5_int64getArgument)
	if err != nil {
		return err
	}

	err = imports.Append("int64finish", v1_5_int64finish, C.v1_5_int64finish)
	if err != nil {
		return err
	}

	err = imports.Append("int64storageStore", v1_5_int64storageStore, C.v1_5_int64storageStore)
	if err != nil {
		return err
	}

	err = imports.Append("int64storageLoad", v1_5_int64storageLoad, C.v1_5_int64storageLoad)
	if err != nil {
		return err
	}

	err = imports.Append("sha256", v1_5_sha256, C.v1_5_sha256)
	if err != nil {
		return err
	}

	err = imports.Append("managedSha256", v1_5_managedSha256, C.v1_5_managedSha256)
	if err != nil {
		return err
	}

	err = imports.Append("keccak256", v1_5_keccak256, C.v1_5_keccak256)
	if err != nil {
		return err
	}

	err = imports.Append("managedKeccak256", v1_5_managedKeccak256, C.v1_5_managedKeccak256)
	if err != nil {
		return err
	}

	err = imports.Append("ripemd160", v1_5_ripemd160, C.v1_5_ripemd160)
	if err != nil {
		return err
	}

	err = imports.Append("managedRipemd160", v1_5_managedRipemd160, C.v1_5_managedRipemd160)
	if err != nil {
		return err
	}

	err = imports.Append("verifyBLS", v1_5_verifyBLS, C.v1_5_verifyBLS)
	if err != nil {
		return err
	}

	err = imports.Append("managedVerifyBLS", v1_5_managedVerifyBLS, C.v1_5_managedVerifyBLS)
	if err != nil {
		return err
	}

	err = imports.Append("verifyEd25519", v1_5_verifyEd25519, C.v1_5_verifyEd25519)
	if err != nil {
		return err
	}

	err = imports.Append("managedVerifyEd25519", v1_5_managedVerifyEd25519, C.v1_5_managedVerifyEd25519)
	if err != nil {
		return err
	}

	err = imports.Append("verifyCustomSecp256k1", v1_5_verifyCustomSecp256k1, C.v1_5_verifyCustomSecp256k1)
	if err != nil {
		return err
	}

	err = imports.Append("managedVerifyCustomSecp256k1", v1_5_managedVerifyCustomSecp256k1, C.v1_5_managedVerifyCustomSecp256k1)
	if err != nil {
		return err
	}

	err = imports.Append("verifySecp256k1", v1_5_verifySecp256k1, C.v1_5_verifySecp256k1)
	if err != nil {
		return err
	}

	err = imports.Append("managedVerifySecp256k1", v1_5_managedVerifySecp256k1, C.v1_5_managedVerifySecp256k1)
	if err != nil {
		return err
	}

	err = imports.Append("encodeSecp256k1DerSignature", v1_5_encodeSecp256k1DerSignature, C.v1_5_encodeSecp256k1DerSignature)
	if err != nil {
		return err
	}

	err = imports.Append("managedEncodeSecp256k1DerSignature", v1_5_managedEncodeSecp256k1DerSignature, C.v1_5_managedEncodeSecp256k1DerSignature)
	if err != nil {
		return err
	}

	err = imports.Append("addEC", v1_5_addEC, C.v1_5_addEC)
	if err != nil {
		return err
	}

	err = imports.Append("doubleEC", v1_5_doubleEC, C.v1_5_doubleEC)
	if err != nil {
		return err
	}

	err = imports.Append("isOnCurveEC", v1_5_isOnCurveEC, C.v1_5_isOnCurveEC)
	if err != nil {
		return err
	}

	err = imports.Append("scalarBaseMultEC", v1_5_scalarBaseMultEC, C.v1_5_scalarBaseMultEC)
	if err != nil {
		return err
	}

	err = imports.Append("managedScalarBaseMultEC", v1_5_managedScalarBaseMultEC, C.v1_5_managedScalarBaseMultEC)
	if err != nil {
		return err
	}

	err = imports.Append("scalarMultEC", v1_5_scalarMultEC, C.v1_5_scalarMultEC)
	if err != nil {
		return err
	}

	err = imports.Append("managedScalarMultEC", v1_5_managedScalarMultEC, C.v1_5_managedScalarMultEC)
	if err != nil {
		return err
	}

	err = imports.Append("marshalEC", v1_5_marshalEC, C.v1_5_marshalEC)
	if err != nil {
		return err
	}

	err = imports.Append("managedMarshalEC", v1_5_managedMarshalEC, C.v1_5_managedMarshalEC)
	if err != nil {
		return err
	}

	err = imports.Append("marshalCompressedEC", v1_5_marshalCompressedEC, C.v1_5_marshalCompressedEC)
	if err != nil {
		return err
	}

	err = imports.Append("managedMarshalCompressedEC", v1_5_managedMarshalCompressedEC, C.v1_5_managedMarshalCompressedEC)
	if err != nil {
		return err
	}

	err = imports.Append("unmarshalEC", v1_5_unmarshalEC, C.v1_5_unmarshalEC)
	if err != nil {
		return err
	}

	err = imports.Append("managedUnmarshalEC", v1_5_managedUnmarshalEC, C.v1_5_managedUnmarshalEC)
	if err != nil {
		return err
	}

	err = imports.Append("unmarshalCompressedEC", v1_5_unmarshalCompressedEC, C.v1_5_unmarshalCompressedEC)
	if err != nil {
		return err
	}

	err = imports.Append("managedUnmarshalCompressedEC", v1_5_managedUnmarshalCompressedEC, C.v1_5_managedUnmarshalCompressedEC)
	if err != nil {
		return err
	}

	err = imports.Append("generateKeyEC", v1_5_generateKeyEC, C.v1_5_generateKeyEC)
	if err != nil {
		return err
	}

	err = imports.Append("managedGenerateKeyEC", v1_5_managedGenerateKeyEC, C.v1_5_managedGenerateKeyEC)
	if err != nil {
		return err
	}

	err = imports.Append("createEC", v1_5_createEC, C.v1_5_createEC)
	if err != nil {
		return err
	}

	err = imports.Append("managedCreateEC", v1_5_managedCreateEC, C.v1_5_managedCreateEC)
	if err != nil {
		return err
	}

	err = imports.Append("getCurveLengthEC", v1_5_getCurveLengthEC, C.v1_5_getCurveLengthEC)
	if err != nil {
		return err
	}

	err = imports.Append("getPrivKeyByteLengthEC", v1_5_getPrivKeyByteLengthEC, C.v1_5_getPrivKeyByteLengthEC)
	if err != nil {
		return err
	}

	err = imports.Append("ellipticCurveGetValues", v1_5_ellipticCurveGetValues, C.v1_5_ellipticCurveGetValues)
	if err != nil {
		return err
	}

	return nil
}

//export v1_5_bigFloatNewFromParts
func v1_5_bigFloatNewFromParts(context unsafe.Pointer, integralPart int32, fractionalPart int32, exponent int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigFloatNewFromParts(integralPart, fractionalPart, exponent)
}

//export v1_5_bigFloatNewFromFrac
func v1_5_bigFloatNewFromFrac(context unsafe.Pointer, numerator int64, denominator int64) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigFloatNewFromFrac(numerator, denominator)
}

//export v1_5_bigFloatNewFromSci
func v1_5_bigFloatNewFromSci(context unsafe.Pointer, significand int64, exponent int64) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigFloatNewFromSci(significand, exponent)
}

//export v1_5_bigFloatAdd
func v1_5_bigFloatAdd(context unsafe.Pointer, destinationHandle int32, op1Handle int32, op2Handle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigFloatAdd(destinationHandle, op1Handle, op2Handle)
}

//export v1_5_bigFloatSub
func v1_5_bigFloatSub(context unsafe.Pointer, destinationHandle int32, op1Handle int32, op2Handle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigFloatSub(destinationHandle, op1Handle, op2Handle)
}

//export v1_5_bigFloatMul
func v1_5_bigFloatMul(context unsafe.Pointer, destinationHandle int32, op1Handle int32, op2Handle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigFloatMul(destinationHandle, op1Handle, op2Handle)
}

//export v1_5_bigFloatDiv
func v1_5_bigFloatDiv(context unsafe.Pointer, destinationHandle int32, op1Handle int32, op2Handle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigFloatDiv(destinationHandle, op1Handle, op2Handle)
}

//export v1_5_bigFloatNeg
func v1_5_bigFloatNeg(context unsafe.Pointer, destinationHandle int32, opHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigFloatNeg(destinationHandle, opHandle)
}

//export v1_5_bigFloatClone
func v1_5_bigFloatClone(context unsafe.Pointer, destinationHandle int32, opHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigFloatClone(destinationHandle, opHandle)
}

//export v1_5_bigFloatCmp
func v1_5_bigFloatCmp(context unsafe.Pointer, op1Handle int32, op2Handle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigFloatCmp(op1Handle, op2Handle)
}

//export v1_5_bigFloatAbs
func v1_5_bigFloatAbs(context unsafe.Pointer, destinationHandle int32, opHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigFloatAbs(destinationHandle, opHandle)
}

//export v1_5_bigFloatSign
func v1_5_bigFloatSign(context unsafe.Pointer, opHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigFloatSign(opHandle)
}

//export v1_5_bigFloatSqrt
func v1_5_bigFloatSqrt(context unsafe.Pointer, destinationHandle int32, opHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigFloatSqrt(destinationHandle, opHandle)
}

//export v1_5_bigFloatPow
func v1_5_bigFloatPow(context unsafe.Pointer, destinationHandle int32, opHandle int32, exponent int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigFloatPow(destinationHandle, opHandle, exponent)
}

//export v1_5_bigFloatFloor
func v1_5_bigFloatFloor(context unsafe.Pointer, destBigIntHandle int32, opHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigFloatFloor(destBigIntHandle, opHandle)
}

//export v1_5_bigFloatCeil
func v1_5_bigFloatCeil(context unsafe.Pointer, destBigIntHandle int32, opHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigFloatCeil(destBigIntHandle, opHandle)
}

//export v1_5_bigFloatTruncate
func v1_5_bigFloatTruncate(context unsafe.Pointer, destBigIntHandle int32, opHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigFloatTruncate(destBigIntHandle, opHandle)
}

//export v1_5_bigFloatSetInt64
func v1_5_bigFloatSetInt64(context unsafe.Pointer, destinationHandle int32, value int64) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigFloatSetInt64(destinationHandle, value)
}

//export v1_5_bigFloatIsInt
func v1_5_bigFloatIsInt(context unsafe.Pointer, opHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigFloatIsInt(opHandle)
}

//export v1_5_bigFloatSetBigInt
func v1_5_bigFloatSetBigInt(context unsafe.Pointer, destinationHandle int32, bigIntHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigFloatSetBigInt(destinationHandle, bigIntHandle)
}

//export v1_5_bigFloatGetConstPi
func v1_5_bigFloatGetConstPi(context unsafe.Pointer, destinationHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigFloatGetConstPi(destinationHandle)
}

//export v1_5_bigFloatGetConstE
func v1_5_bigFloatGetConstE(context unsafe.Pointer, destinationHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigFloatGetConstE(destinationHandle)
}

//export v1_5_bigIntGetUnsignedArgument
func v1_5_bigIntGetUnsignedArgument(context unsafe.Pointer, id int32, destinationHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntGetUnsignedArgument(id, destinationHandle)
}

//export v1_5_bigIntGetSignedArgument
func v1_5_bigIntGetSignedArgument(context unsafe.Pointer, id int32, destinationHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntGetSignedArgument(id, destinationHandle)
}

//export v1_5_bigIntStorageStoreUnsigned
func v1_5_bigIntStorageStoreUnsigned(context unsafe.Pointer, keyOffset int32, keyLength int32, sourceHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigIntStorageStoreUnsigned(keyOffset, keyLength, sourceHandle)
}

//export v1_5_bigIntStorageLoadUnsigned
func v1_5_bigIntStorageLoadUnsigned(context unsafe.Pointer, keyOffset int32, keyLength int32, destinationHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigIntStorageLoadUnsigned(keyOffset, keyLength, destinationHandle)
}

//export v1_5_bigIntGetCallValue
func v1_5_bigIntGetCallValue(context unsafe.Pointer, destinationHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntGetCallValue(destinationHandle)
}

//export v1_5_bigIntGetESDTCallValue
func v1_5_bigIntGetESDTCallValue(context unsafe.Pointer, destination int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntGetESDTCallValue(destination)
}

//export v1_5_bigIntGetESDTCallValueByIndex
func v1_5_bigIntGetESDTCallValueByIndex(context unsafe.Pointer, destinationHandle int32, index int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntGetESDTCallValueByIndex(destinationHandle, index)
}

//export v1_5_bigIntGetExternalBalance
func v1_5_bigIntGetExternalBalance(context unsafe.Pointer, addressOffset int32, result int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntGetExternalBalance(addressOffset, result)
}

//export v1_5_bigIntGetESDTExternalBalance
func v1_5_bigIntGetESDTExternalBalance(context unsafe.Pointer, addressOffset int32, tokenIDOffset int32, tokenIDLen int32, nonce int64, resultHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntGetESDTExternalBalance(addressOffset, tokenIDOffset, tokenIDLen, nonce, resultHandle)
}

//export v1_5_bigIntNew
func v1_5_bigIntNew(context unsafe.Pointer, smallValue int64) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigIntNew(smallValue)
}

//export v1_5_bigIntUnsignedByteLength
func v1_5_bigIntUnsignedByteLength(context unsafe.Pointer, referenceHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigIntUnsignedByteLength(referenceHandle)
}

//export v1_5_bigIntSignedByteLength
func v1_5_bigIntSignedByteLength(context unsafe.Pointer, referenceHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigIntSignedByteLength(referenceHandle)
}

//export v1_5_bigIntGetUnsignedBytes
func v1_5_bigIntGetUnsignedBytes(context unsafe.Pointer, referenceHandle int32, byteOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigIntGetUnsignedBytes(referenceHandle, byteOffset)
}

//export v1_5_bigIntGetSignedBytes
func v1_5_bigIntGetSignedBytes(context unsafe.Pointer, referenceHandle int32, byteOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigIntGetSignedBytes(referenceHandle, byteOffset)
}

//export v1_5_bigIntSetUnsignedBytes
func v1_5_bigIntSetUnsignedBytes(context unsafe.Pointer, destinationHandle int32, byteOffset int32, byteLength int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntSetUnsignedBytes(destinationHandle, byteOffset, byteLength)
}

//export v1_5_bigIntSetSignedBytes
func v1_5_bigIntSetSignedBytes(context unsafe.Pointer, destinationHandle int32, byteOffset int32, byteLength int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntSetSignedBytes(destinationHandle, byteOffset, byteLength)
}

//export v1_5_bigIntIsInt64
func v1_5_bigIntIsInt64(context unsafe.Pointer, destinationHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigIntIsInt64(destinationHandle)
}

//export v1_5_bigIntGetInt64
func v1_5_bigIntGetInt64(context unsafe.Pointer, destinationHandle int32) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigIntGetInt64(destinationHandle)
}

//export v1_5_bigIntSetInt64
func v1_5_bigIntSetInt64(context unsafe.Pointer, destinationHandle int32, value int64) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntSetInt64(destinationHandle, value)
}

//export v1_5_bigIntAdd
func v1_5_bigIntAdd(context unsafe.Pointer, destinationHandle int32, op1Handle int32, op2Handle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntAdd(destinationHandle, op1Handle, op2Handle)
}

//export v1_5_bigIntSub
func v1_5_bigIntSub(context unsafe.Pointer, destinationHandle int32, op1Handle int32, op2Handle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntSub(destinationHandle, op1Handle, op2Handle)
}

//export v1_5_bigIntMul
func v1_5_bigIntMul(context unsafe.Pointer, destinationHandle int32, op1Handle int32, op2Handle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntMul(destinationHandle, op1Handle, op2Handle)
}

//export v1_5_bigIntTDiv
func v1_5_bigIntTDiv(context unsafe.Pointer, destinationHandle int32, op1Handle int32, op2Handle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntTDiv(destinationHandle, op1Handle, op2Handle)
}

//export v1_5_bigIntTMod
func v1_5_bigIntTMod(context unsafe.Pointer, destinationHandle int32, op1Handle int32, op2Handle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntTMod(destinationHandle, op1Handle, op2Handle)
}

//export v1_5_bigIntEDiv
func v1_5_bigIntEDiv(context unsafe.Pointer, destinationHandle int32, op1Handle int32, op2Handle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntEDiv(destinationHandle, op1Handle, op2Handle)
}

//export v1_5_bigIntEMod
func v1_5_bigIntEMod(context unsafe.Pointer, destinationHandle int32, op1Handle int32, op2Handle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntEMod(destinationHandle, op1Handle, op2Handle)
}

//export v1_5_bigIntSqrt
func v1_5_bigIntSqrt(context unsafe.Pointer, destinationHandle int32, opHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntSqrt(destinationHandle, opHandle)
}

//export v1_5_bigIntPow
func v1_5_bigIntPow(context unsafe.Pointer, destinationHandle int32, op1Handle int32, op2Handle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntPow(destinationHandle, op1Handle, op2Handle)
}

//export v1_5_bigIntLog2
func v1_5_bigIntLog2(context unsafe.Pointer, op1Handle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigIntLog2(op1Handle)
}

//export v1_5_bigIntAbs
func v1_5_bigIntAbs(context unsafe.Pointer, destinationHandle int32, opHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntAbs(destinationHandle, opHandle)
}

//export v1_5_bigIntNeg
func v1_5_bigIntNeg(context unsafe.Pointer, destinationHandle int32, opHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntNeg(destinationHandle, opHandle)
}

//export v1_5_bigIntSign
func v1_5_bigIntSign(context unsafe.Pointer, opHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigIntSign(opHandle)
}

//export v1_5_bigIntCmp
func v1_5_bigIntCmp(context unsafe.Pointer, op1Handle int32, op2Handle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BigIntCmp(op1Handle, op2Handle)
}

//export v1_5_bigIntNot
func v1_5_bigIntNot(context unsafe.Pointer, destinationHandle int32, opHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntNot(destinationHandle, opHandle)
}

//export v1_5_bigIntAnd
func v1_5_bigIntAnd(context unsafe.Pointer, destinationHandle int32, op1Handle int32, op2Handle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntAnd(destinationHandle, op1Handle, op2Handle)
}

//export v1_5_bigIntOr
func v1_5_bigIntOr(context unsafe.Pointer, destinationHandle int32, op1Handle int32, op2Handle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntOr(destinationHandle, op1Handle, op2Handle)
}

//export v1_5_bigIntXor
func v1_5_bigIntXor(context unsafe.Pointer, destinationHandle int32, op1Handle int32, op2Handle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntXor(destinationHandle, op1Handle, op2Handle)
}

//export v1_5_bigIntShr
func v1_5_bigIntShr(context unsafe.Pointer, destinationHandle int32, opHandle int32, bits int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntShr(destinationHandle, opHandle, bits)
}

//export v1_5_bigIntShl
func v1_5_bigIntShl(context unsafe.Pointer, destinationHandle int32, opHandle int32, bits int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntShl(destinationHandle, opHandle, bits)
}

//export v1_5_bigIntFinishUnsigned
func v1_5_bigIntFinishUnsigned(context unsafe.Pointer, referenceHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntFinishUnsigned(referenceHandle)
}

//export v1_5_bigIntFinishSigned
func v1_5_bigIntFinishSigned(context unsafe.Pointer, referenceHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntFinishSigned(referenceHandle)
}

//export v1_5_bigIntToString
func v1_5_bigIntToString(context unsafe.Pointer, bigIntHandle int32, destinationHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.BigIntToString(bigIntHandle, destinationHandle)
}

//export v1_5_getGasLeft
func v1_5_getGasLeft(context unsafe.Pointer) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetGasLeft()
}

//export v1_5_getSCAddress
func v1_5_getSCAddress(context unsafe.Pointer, resultOffset int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.GetSCAddress(resultOffset)
}

//export v1_5_getOwnerAddress
func v1_5_getOwnerAddress(context unsafe.Pointer, resultOffset int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.GetOwnerAddress(resultOffset)
}

//export v1_5_getShardOfAddress
func v1_5_getShardOfAddress(context unsafe.Pointer, addressOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetShardOfAddress(addressOffset)
}

//export v1_5_isSmartContract
func v1_5_isSmartContract(context unsafe.Pointer, addressOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.IsSmartContract(addressOffset)
}

//export v1_5_signalError
func v1_5_signalError(context unsafe.Pointer, messageOffset int32, messageLength int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.SignalError(messageOffset, messageLength)
}

//export v1_5_getExternalBalance
func v1_5_getExternalBalance(context unsafe.Pointer, addressOffset int32, resultOffset int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.GetExternalBalance(addressOffset, resultOffset)
}

//export v1_5_blockHash
func v1_5_blockHash(context unsafe.Pointer, nonce int64, resultOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.BlockHash(nonce, resultOffset)
}

//export v1_5_getESDTBalance
func v1_5_getESDTBalance(context unsafe.Pointer, addressOffset int32, tokenIDOffset int32, tokenIDLen int32, nonce int64, resultOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetESDTBalance(addressOffset, tokenIDOffset, tokenIDLen, nonce, resultOffset)
}

//export v1_5_getESDTNFTNameLength
func v1_5_getESDTNFTNameLength(context unsafe.Pointer, addressOffset int32, tokenIDOffset int32, tokenIDLen int32, nonce int64) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetESDTNFTNameLength(addressOffset, tokenIDOffset, tokenIDLen, nonce)
}

//export v1_5_getESDTNFTAttributeLength
func v1_5_getESDTNFTAttributeLength(context unsafe.Pointer, addressOffset int32, tokenIDOffset int32, tokenIDLen int32, nonce int64) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetESDTNFTAttributeLength(addressOffset, tokenIDOffset, tokenIDLen, nonce)
}

//export v1_5_getESDTNFTURILength
func v1_5_getESDTNFTURILength(context unsafe.Pointer, addressOffset int32, tokenIDOffset int32, tokenIDLen int32, nonce int64) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetESDTNFTURILength(addressOffset, tokenIDOffset, tokenIDLen, nonce)
}

//export v1_5_getESDTTokenData
func v1_5_getESDTTokenData(context unsafe.Pointer, addressOffset int32, tokenIDOffset int32, tokenIDLen int32, nonce int64, valueHandle int32, propertiesOffset int32, hashOffset int32, nameOffset int32, attributesOffset int32, creatorOffset int32, royaltiesHandle int32, urisOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetESDTTokenData(addressOffset, tokenIDOffset, tokenIDLen, nonce, valueHandle, propertiesOffset, hashOffset, nameOffset, attributesOffset, creatorOffset, royaltiesHandle, urisOffset)
}

//export v1_5_getESDTLocalRoles
func v1_5_getESDTLocalRoles(context unsafe.Pointer, tokenIdHandle int32) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetESDTLocalRoles(tokenIdHandle)
}

//export v1_5_validateTokenIdentifier
func v1_5_validateTokenIdentifier(context unsafe.Pointer, tokenIdHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ValidateTokenIdentifier(tokenIdHandle)
}

//export v1_5_transferValue
func v1_5_transferValue(context unsafe.Pointer, destOffset int32, valueOffset int32, dataOffset int32, length int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.TransferValue(destOffset, valueOffset, dataOffset, length)
}

//export v1_5_transferValueExecute
func v1_5_transferValueExecute(context unsafe.Pointer, destOffset int32, valueOffset int32, gasLimit int64, functionOffset int32, functionLength int32, numArguments int32, argumentsLengthOffset int32, dataOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.TransferValueExecute(destOffset, valueOffset, gasLimit, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
}

//export v1_5_transferESDTExecute
func v1_5_transferESDTExecute(context unsafe.Pointer, destOffset int32, tokenIDOffset int32, tokenIDLen int32, valueOffset int32, gasLimit int64, functionOffset int32, functionLength int32, numArguments int32, argumentsLengthOffset int32, dataOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.TransferESDTExecute(destOffset, tokenIDOffset, tokenIDLen, valueOffset, gasLimit, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
}

//export v1_5_transferESDTNFTExecute
func v1_5_transferESDTNFTExecute(context unsafe.Pointer, destOffset int32, tokenIDOffset int32, tokenIDLen int32, valueOffset int32, nonce int64, gasLimit int64, functionOffset int32, functionLength int32, numArguments int32, argumentsLengthOffset int32, dataOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.TransferESDTNFTExecute(destOffset, tokenIDOffset, tokenIDLen, valueOffset, nonce, gasLimit, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
}

//export v1_5_multiTransferESDTNFTExecute
func v1_5_multiTransferESDTNFTExecute(context unsafe.Pointer, destOffset int32, numTokenTransfers int32, tokenTransfersArgsLengthOffset int32, tokenTransferDataOffset int32, gasLimit int64, functionOffset int32, functionLength int32, numArguments int32, argumentsLengthOffset int32, dataOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MultiTransferESDTNFTExecute(destOffset, numTokenTransfers, tokenTransfersArgsLengthOffset, tokenTransferDataOffset, gasLimit, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
}

//export v1_5_createAsyncCall
func v1_5_createAsyncCall(context unsafe.Pointer, destOffset int32, valueOffset int32, dataOffset int32, dataLength int32, successOffset int32, successLength int32, errorOffset int32, errorLength int32, gas int64, extraGasForCallback int64) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.CreateAsyncCall(destOffset, valueOffset, dataOffset, dataLength, successOffset, successLength, errorOffset, errorLength, gas, extraGasForCallback)
}

//export v1_5_setAsyncContextCallback
func v1_5_setAsyncContextCallback(context unsafe.Pointer, callback int32, callbackLength int32, data int32, dataLength int32, gas int64) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.SetAsyncContextCallback(callback, callbackLength, data, dataLength, gas)
}

//export v1_5_upgradeContract
func v1_5_upgradeContract(context unsafe.Pointer, destOffset int32, gasLimit int64, valueOffset int32, codeOffset int32, codeMetadataOffset int32, length int32, numArguments int32, argumentsLengthOffset int32, dataOffset int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.UpgradeContract(destOffset, gasLimit, valueOffset, codeOffset, codeMetadataOffset, length, numArguments, argumentsLengthOffset, dataOffset)
}

//export v1_5_upgradeFromSourceContract
func v1_5_upgradeFromSourceContract(context unsafe.Pointer, destOffset int32, gasLimit int64, valueOffset int32, sourceContractAddressOffset int32, codeMetadataOffset int32, numArguments int32, argumentsLengthOffset int32, dataOffset int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.UpgradeFromSourceContract(destOffset, gasLimit, valueOffset, sourceContractAddressOffset, codeMetadataOffset, numArguments, argumentsLengthOffset, dataOffset)
}

//export v1_5_deleteContract
func v1_5_deleteContract(context unsafe.Pointer, destOffset int32, gasLimit int64, numArguments int32, argumentsLengthOffset int32, dataOffset int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.DeleteContract(destOffset, gasLimit, numArguments, argumentsLengthOffset, dataOffset)
}

//export v1_5_asyncCall
func v1_5_asyncCall(context unsafe.Pointer, destOffset int32, valueOffset int32, dataOffset int32, length int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.AsyncCall(destOffset, valueOffset, dataOffset, length)
}

//export v1_5_getArgumentLength
func v1_5_getArgumentLength(context unsafe.Pointer, id int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetArgumentLength(id)
}

//export v1_5_getArgument
func v1_5_getArgument(context unsafe.Pointer, id int32, argOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetArgument(id, argOffset)
}

//export v1_5_getFunction
func v1_5_getFunction(context unsafe.Pointer, functionOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetFunction(functionOffset)
}

//export v1_5_getNumArguments
func v1_5_getNumArguments(context unsafe.Pointer) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetNumArguments()
}

//export v1_5_storageStore
func v1_5_storageStore(context unsafe.Pointer, keyOffset int32, keyLength int32, dataOffset int32, dataLength int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.StorageStore(keyOffset, keyLength, dataOffset, dataLength)
}

//export v1_5_storageLoadLength
func v1_5_storageLoadLength(context unsafe.Pointer, keyOffset int32, keyLength int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.StorageLoadLength(keyOffset, keyLength)
}

//export v1_5_storageLoadFromAddress
func v1_5_storageLoadFromAddress(context unsafe.Pointer, addressOffset int32, keyOffset int32, keyLength int32, dataOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.StorageLoadFromAddress(addressOffset, keyOffset, keyLength, dataOffset)
}

//export v1_5_storageLoad
func v1_5_storageLoad(context unsafe.Pointer, keyOffset int32, keyLength int32, dataOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.StorageLoad(keyOffset, keyLength, dataOffset)
}

//export v1_5_setStorageLock
func v1_5_setStorageLock(context unsafe.Pointer, keyOffset int32, keyLength int32, lockTimestamp int64) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.SetStorageLock(keyOffset, keyLength, lockTimestamp)
}

//export v1_5_getStorageLock
func v1_5_getStorageLock(context unsafe.Pointer, keyOffset int32, keyLength int32) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetStorageLock(keyOffset, keyLength)
}

//export v1_5_isStorageLocked
func v1_5_isStorageLocked(context unsafe.Pointer, keyOffset int32, keyLength int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.IsStorageLocked(keyOffset, keyLength)
}

//export v1_5_clearStorageLock
func v1_5_clearStorageLock(context unsafe.Pointer, keyOffset int32, keyLength int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ClearStorageLock(keyOffset, keyLength)
}

//export v1_5_getCaller
func v1_5_getCaller(context unsafe.Pointer, resultOffset int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.GetCaller(resultOffset)
}

//export v1_5_checkNoPayment
func v1_5_checkNoPayment(context unsafe.Pointer) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.CheckNoPayment()
}

//export v1_5_callValue
func v1_5_callValue(context unsafe.Pointer, resultOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.CallValue(resultOffset)
}

//export v1_5_getESDTValue
func v1_5_getESDTValue(context unsafe.Pointer, resultOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetESDTValue(resultOffset)
}

//export v1_5_getESDTValueByIndex
func v1_5_getESDTValueByIndex(context unsafe.Pointer, resultOffset int32, index int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetESDTValueByIndex(resultOffset, index)
}

//export v1_5_getESDTTokenName
func v1_5_getESDTTokenName(context unsafe.Pointer, resultOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetESDTTokenName(resultOffset)
}

//export v1_5_getESDTTokenNameByIndex
func v1_5_getESDTTokenNameByIndex(context unsafe.Pointer, resultOffset int32, index int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetESDTTokenNameByIndex(resultOffset, index)
}

//export v1_5_getESDTTokenNonce
func v1_5_getESDTTokenNonce(context unsafe.Pointer) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetESDTTokenNonce()
}

//export v1_5_getESDTTokenNonceByIndex
func v1_5_getESDTTokenNonceByIndex(context unsafe.Pointer, index int32) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetESDTTokenNonceByIndex(index)
}

//export v1_5_getCurrentESDTNFTNonce
func v1_5_getCurrentESDTNFTNonce(context unsafe.Pointer, addressOffset int32, tokenIDOffset int32, tokenIDLen int32) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetCurrentESDTNFTNonce(addressOffset, tokenIDOffset, tokenIDLen)
}

//export v1_5_getESDTTokenType
func v1_5_getESDTTokenType(context unsafe.Pointer) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetESDTTokenType()
}

//export v1_5_getESDTTokenTypeByIndex
func v1_5_getESDTTokenTypeByIndex(context unsafe.Pointer, index int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetESDTTokenTypeByIndex(index)
}

//export v1_5_getNumESDTTransfers
func v1_5_getNumESDTTransfers(context unsafe.Pointer) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetNumESDTTransfers()
}

//export v1_5_getCallValueTokenName
func v1_5_getCallValueTokenName(context unsafe.Pointer, callValueOffset int32, tokenNameOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetCallValueTokenName(callValueOffset, tokenNameOffset)
}

//export v1_5_getCallValueTokenNameByIndex
func v1_5_getCallValueTokenNameByIndex(context unsafe.Pointer, callValueOffset int32, tokenNameOffset int32, index int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetCallValueTokenNameByIndex(callValueOffset, tokenNameOffset, index)
}

//export v1_5_writeLog
func v1_5_writeLog(context unsafe.Pointer, dataPointer int32, dataLength int32, topicPtr int32, numTopics int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.WriteLog(dataPointer, dataLength, topicPtr, numTopics)
}

//export v1_5_writeEventLog
func v1_5_writeEventLog(context unsafe.Pointer, numTopics int32, topicLengthsOffset int32, topicOffset int32, dataOffset int32, dataLength int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.WriteEventLog(numTopics, topicLengthsOffset, topicOffset, dataOffset, dataLength)
}

//export v1_5_getBlockTimestamp
func v1_5_getBlockTimestamp(context unsafe.Pointer) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetBlockTimestamp()
}

//export v1_5_getBlockNonce
func v1_5_getBlockNonce(context unsafe.Pointer) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetBlockNonce()
}

//export v1_5_getBlockRound
func v1_5_getBlockRound(context unsafe.Pointer) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetBlockRound()
}

//export v1_5_getBlockEpoch
func v1_5_getBlockEpoch(context unsafe.Pointer) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetBlockEpoch()
}

//export v1_5_getBlockRandomSeed
func v1_5_getBlockRandomSeed(context unsafe.Pointer, pointer int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.GetBlockRandomSeed(pointer)
}

//export v1_5_getStateRootHash
func v1_5_getStateRootHash(context unsafe.Pointer, pointer int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.GetStateRootHash(pointer)
}

//export v1_5_getPrevBlockTimestamp
func v1_5_getPrevBlockTimestamp(context unsafe.Pointer) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetPrevBlockTimestamp()
}

//export v1_5_getPrevBlockNonce
func v1_5_getPrevBlockNonce(context unsafe.Pointer) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetPrevBlockNonce()
}

//export v1_5_getPrevBlockRound
func v1_5_getPrevBlockRound(context unsafe.Pointer) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetPrevBlockRound()
}

//export v1_5_getPrevBlockEpoch
func v1_5_getPrevBlockEpoch(context unsafe.Pointer) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetPrevBlockEpoch()
}

//export v1_5_getPrevBlockRandomSeed
func v1_5_getPrevBlockRandomSeed(context unsafe.Pointer, pointer int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.GetPrevBlockRandomSeed(pointer)
}

//export v1_5_returnData
func v1_5_returnData(context unsafe.Pointer, pointer int32, length int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ReturnData(pointer, length)
}

//export v1_5_executeOnSameContext
func v1_5_executeOnSameContext(context unsafe.Pointer, gasLimit int64, addressOffset int32, valueOffset int32, functionOffset int32, functionLength int32, numArguments int32, argumentsLengthOffset int32, dataOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ExecuteOnSameContext(gasLimit, addressOffset, valueOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
}

//export v1_5_executeOnDestContext
func v1_5_executeOnDestContext(context unsafe.Pointer, gasLimit int64, addressOffset int32, valueOffset int32, functionOffset int32, functionLength int32, numArguments int32, argumentsLengthOffset int32, dataOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ExecuteOnDestContext(gasLimit, addressOffset, valueOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
}

//export v1_5_executeReadOnly
func v1_5_executeReadOnly(context unsafe.Pointer, gasLimit int64, addressOffset int32, functionOffset int32, functionLength int32, numArguments int32, argumentsLengthOffset int32, dataOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ExecuteReadOnly(gasLimit, addressOffset, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)
}

//export v1_5_createContract
func v1_5_createContract(context unsafe.Pointer, gasLimit int64, valueOffset int32, codeOffset int32, codeMetadataOffset int32, length int32, resultOffset int32, numArguments int32, argumentsLengthOffset int32, dataOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.CreateContract(gasLimit, valueOffset, codeOffset, codeMetadataOffset, length, resultOffset, numArguments, argumentsLengthOffset, dataOffset)
}

//export v1_5_deployFromSourceContract
func v1_5_deployFromSourceContract(context unsafe.Pointer, gasLimit int64, valueOffset int32, sourceContractAddressOffset int32, codeMetadataOffset int32, resultAddressOffset int32, numArguments int32, argumentsLengthOffset int32, dataOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.DeployFromSourceContract(gasLimit, valueOffset, sourceContractAddressOffset, codeMetadataOffset, resultAddressOffset, numArguments, argumentsLengthOffset, dataOffset)
}

//export v1_5_getNumReturnData
func v1_5_getNumReturnData(context unsafe.Pointer) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetNumReturnData()
}

//export v1_5_getReturnDataSize
func v1_5_getReturnDataSize(context unsafe.Pointer, resultID int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetReturnDataSize(resultID)
}

//export v1_5_getReturnData
func v1_5_getReturnData(context unsafe.Pointer, resultID int32, dataOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetReturnData(resultID, dataOffset)
}

//export v1_5_cleanReturnData
func v1_5_cleanReturnData(context unsafe.Pointer) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.CleanReturnData()
}

//export v1_5_deleteFromReturnData
func v1_5_deleteFromReturnData(context unsafe.Pointer, resultID int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.DeleteFromReturnData(resultID)
}

//export v1_5_getOriginalTxHash
func v1_5_getOriginalTxHash(context unsafe.Pointer, dataOffset int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.GetOriginalTxHash(dataOffset)
}

//export v1_5_getCurrentTxHash
func v1_5_getCurrentTxHash(context unsafe.Pointer, dataOffset int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.GetCurrentTxHash(dataOffset)
}

//export v1_5_getPrevTxHash
func v1_5_getPrevTxHash(context unsafe.Pointer, dataOffset int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.GetPrevTxHash(dataOffset)
}

//export v1_5_managedSCAddress
func v1_5_managedSCAddress(context unsafe.Pointer, destinationHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedSCAddress(destinationHandle)
}

//export v1_5_managedOwnerAddress
func v1_5_managedOwnerAddress(context unsafe.Pointer, destinationHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedOwnerAddress(destinationHandle)
}

//export v1_5_managedCaller
func v1_5_managedCaller(context unsafe.Pointer, destinationHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedCaller(destinationHandle)
}

//export v1_5_managedSignalError
func v1_5_managedSignalError(context unsafe.Pointer, errHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedSignalError(errHandle)
}

//export v1_5_managedWriteLog
func v1_5_managedWriteLog(context unsafe.Pointer, topicsHandle int32, dataHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedWriteLog(topicsHandle, dataHandle)
}

//export v1_5_managedGetOriginalTxHash
func v1_5_managedGetOriginalTxHash(context unsafe.Pointer, resultHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedGetOriginalTxHash(resultHandle)
}

//export v1_5_managedGetStateRootHash
func v1_5_managedGetStateRootHash(context unsafe.Pointer, resultHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedGetStateRootHash(resultHandle)
}

//export v1_5_managedGetBlockRandomSeed
func v1_5_managedGetBlockRandomSeed(context unsafe.Pointer, resultHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedGetBlockRandomSeed(resultHandle)
}

//export v1_5_managedGetPrevBlockRandomSeed
func v1_5_managedGetPrevBlockRandomSeed(context unsafe.Pointer, resultHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedGetPrevBlockRandomSeed(resultHandle)
}

//export v1_5_managedGetReturnData
func v1_5_managedGetReturnData(context unsafe.Pointer, resultID int32, resultHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedGetReturnData(resultID, resultHandle)
}

//export v1_5_managedGetMultiESDTCallValue
func v1_5_managedGetMultiESDTCallValue(context unsafe.Pointer, multiCallValueHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedGetMultiESDTCallValue(multiCallValueHandle)
}

//export v1_5_managedGetESDTBalance
func v1_5_managedGetESDTBalance(context unsafe.Pointer, addressHandle int32, tokenIDHandle int32, nonce int64, valueHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedGetESDTBalance(addressHandle, tokenIDHandle, nonce, valueHandle)
}

//export v1_5_managedGetESDTTokenData
func v1_5_managedGetESDTTokenData(context unsafe.Pointer, addressHandle int32, tokenIDHandle int32, nonce int64, valueHandle int32, propertiesHandle int32, hashHandle int32, nameHandle int32, attributesHandle int32, creatorHandle int32, royaltiesHandle int32, urisHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedGetESDTTokenData(addressHandle, tokenIDHandle, nonce, valueHandle, propertiesHandle, hashHandle, nameHandle, attributesHandle, creatorHandle, royaltiesHandle, urisHandle)
}

//export v1_5_managedAsyncCall
func v1_5_managedAsyncCall(context unsafe.Pointer, destHandle int32, valueHandle int32, functionHandle int32, argumentsHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedAsyncCall(destHandle, valueHandle, functionHandle, argumentsHandle)
}

//export v1_5_managedCreateAsyncCall
func v1_5_managedCreateAsyncCall(context unsafe.Pointer, destHandle int32, valueHandle int32, functionHandle int32, argumentsHandle int32, successOffset int32, successLength int32, errorOffset int32, errorLength int32, gas int64, extraGasForCallback int64, callbackClosureHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedCreateAsyncCall(destHandle, valueHandle, functionHandle, argumentsHandle, successOffset, successLength, errorOffset, errorLength, gas, extraGasForCallback, callbackClosureHandle)
}

//export v1_5_managedGetCallbackClosure
func v1_5_managedGetCallbackClosure(context unsafe.Pointer, callbackClosureHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedGetCallbackClosure(callbackClosureHandle)
}

//export v1_5_managedUpgradeFromSourceContract
func v1_5_managedUpgradeFromSourceContract(context unsafe.Pointer, destHandle int32, gas int64, valueHandle int32, addressHandle int32, codeMetadataHandle int32, argumentsHandle int32, resultHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedUpgradeFromSourceContract(destHandle, gas, valueHandle, addressHandle, codeMetadataHandle, argumentsHandle, resultHandle)
}

//export v1_5_managedUpgradeContract
func v1_5_managedUpgradeContract(context unsafe.Pointer, destHandle int32, gas int64, valueHandle int32, codeHandle int32, codeMetadataHandle int32, argumentsHandle int32, resultHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedUpgradeContract(destHandle, gas, valueHandle, codeHandle, codeMetadataHandle, argumentsHandle, resultHandle)
}

//export v1_5_managedDeleteContract
func v1_5_managedDeleteContract(context unsafe.Pointer, destHandle int32, gasLimit int64, argumentsHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedDeleteContract(destHandle, gasLimit, argumentsHandle)
}

//export v1_5_managedDeployFromSourceContract
func v1_5_managedDeployFromSourceContract(context unsafe.Pointer, gas int64, valueHandle int32, addressHandle int32, codeMetadataHandle int32, argumentsHandle int32, resultAddressHandle int32, resultHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedDeployFromSourceContract(gas, valueHandle, addressHandle, codeMetadataHandle, argumentsHandle, resultAddressHandle, resultHandle)
}

//export v1_5_managedCreateContract
func v1_5_managedCreateContract(context unsafe.Pointer, gas int64, valueHandle int32, codeHandle int32, codeMetadataHandle int32, argumentsHandle int32, resultAddressHandle int32, resultHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedCreateContract(gas, valueHandle, codeHandle, codeMetadataHandle, argumentsHandle, resultAddressHandle, resultHandle)
}

//export v1_5_managedExecuteReadOnly
func v1_5_managedExecuteReadOnly(context unsafe.Pointer, gas int64, addressHandle int32, functionHandle int32, argumentsHandle int32, resultHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedExecuteReadOnly(gas, addressHandle, functionHandle, argumentsHandle, resultHandle)
}

//export v1_5_managedExecuteOnSameContext
func v1_5_managedExecuteOnSameContext(context unsafe.Pointer, gas int64, addressHandle int32, valueHandle int32, functionHandle int32, argumentsHandle int32, resultHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedExecuteOnSameContext(gas, addressHandle, valueHandle, functionHandle, argumentsHandle, resultHandle)
}

//export v1_5_managedExecuteOnDestContext
func v1_5_managedExecuteOnDestContext(context unsafe.Pointer, gas int64, addressHandle int32, valueHandle int32, functionHandle int32, argumentsHandle int32, resultHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedExecuteOnDestContext(gas, addressHandle, valueHandle, functionHandle, argumentsHandle, resultHandle)
}

//export v1_5_managedMultiTransferESDTNFTExecute
func v1_5_managedMultiTransferESDTNFTExecute(context unsafe.Pointer, dstHandle int32, tokenTransfersHandle int32, gasLimit int64, functionHandle int32, argumentsHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedMultiTransferESDTNFTExecute(dstHandle, tokenTransfersHandle, gasLimit, functionHandle, argumentsHandle)
}

//export v1_5_managedTransferValueExecute
func v1_5_managedTransferValueExecute(context unsafe.Pointer, dstHandle int32, valueHandle int32, gasLimit int64, functionHandle int32, argumentsHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedTransferValueExecute(dstHandle, valueHandle, gasLimit, functionHandle, argumentsHandle)
}

//export v1_5_managedIsESDTFrozen
func v1_5_managedIsESDTFrozen(context unsafe.Pointer, addressHandle int32, tokenIDHandle int32, nonce int64) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedIsESDTFrozen(addressHandle, tokenIDHandle, nonce)
}

//export v1_5_managedIsESDTLimitedTransfer
func v1_5_managedIsESDTLimitedTransfer(context unsafe.Pointer, tokenIDHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedIsESDTLimitedTransfer(tokenIDHandle)
}

//export v1_5_managedIsESDTPaused
func v1_5_managedIsESDTPaused(context unsafe.Pointer, tokenIDHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedIsESDTPaused(tokenIDHandle)
}

//export v1_5_managedBufferToHex
func v1_5_managedBufferToHex(context unsafe.Pointer, sourceHandle int32, destHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.ManagedBufferToHex(sourceHandle, destHandle)
}

//export v1_5_mBufferNew
func v1_5_mBufferNew(context unsafe.Pointer) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferNew()
}

//export v1_5_mBufferNewFromBytes
func v1_5_mBufferNewFromBytes(context unsafe.Pointer, dataOffset int32, dataLength int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferNewFromBytes(dataOffset, dataLength)
}

//export v1_5_mBufferGetLength
func v1_5_mBufferGetLength(context unsafe.Pointer, mBufferHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferGetLength(mBufferHandle)
}

//export v1_5_mBufferGetBytes
func v1_5_mBufferGetBytes(context unsafe.Pointer, mBufferHandle int32, resultOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferGetBytes(mBufferHandle, resultOffset)
}

//export v1_5_mBufferGetByteSlice
func v1_5_mBufferGetByteSlice(context unsafe.Pointer, sourceHandle int32, startingPosition int32, sliceLength int32, resultOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferGetByteSlice(sourceHandle, startingPosition, sliceLength, resultOffset)
}

//export v1_5_mBufferCopyByteSlice
func v1_5_mBufferCopyByteSlice(context unsafe.Pointer, sourceHandle int32, startingPosition int32, sliceLength int32, destinationHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferCopyByteSlice(sourceHandle, startingPosition, sliceLength, destinationHandle)
}

//export v1_5_mBufferEq
func v1_5_mBufferEq(context unsafe.Pointer, mBufferHandle1 int32, mBufferHandle2 int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferEq(mBufferHandle1, mBufferHandle2)
}

//export v1_5_mBufferSetBytes
func v1_5_mBufferSetBytes(context unsafe.Pointer, mBufferHandle int32, dataOffset int32, dataLength int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferSetBytes(mBufferHandle, dataOffset, dataLength)
}

//export v1_5_mBufferSetByteSlice
func v1_5_mBufferSetByteSlice(context unsafe.Pointer, mBufferHandle int32, startingPosition int32, dataLength int32, dataOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferSetByteSlice(mBufferHandle, startingPosition, dataLength, dataOffset)
}

//export v1_5_mBufferAppend
func v1_5_mBufferAppend(context unsafe.Pointer, accumulatorHandle int32, dataHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferAppend(accumulatorHandle, dataHandle)
}

//export v1_5_mBufferAppendBytes
func v1_5_mBufferAppendBytes(context unsafe.Pointer, accumulatorHandle int32, dataOffset int32, dataLength int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferAppendBytes(accumulatorHandle, dataOffset, dataLength)
}

//export v1_5_mBufferToBigIntUnsigned
func v1_5_mBufferToBigIntUnsigned(context unsafe.Pointer, mBufferHandle int32, bigIntHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferToBigIntUnsigned(mBufferHandle, bigIntHandle)
}

//export v1_5_mBufferToBigIntSigned
func v1_5_mBufferToBigIntSigned(context unsafe.Pointer, mBufferHandle int32, bigIntHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferToBigIntSigned(mBufferHandle, bigIntHandle)
}

//export v1_5_mBufferFromBigIntUnsigned
func v1_5_mBufferFromBigIntUnsigned(context unsafe.Pointer, mBufferHandle int32, bigIntHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferFromBigIntUnsigned(mBufferHandle, bigIntHandle)
}

//export v1_5_mBufferFromBigIntSigned
func v1_5_mBufferFromBigIntSigned(context unsafe.Pointer, mBufferHandle int32, bigIntHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferFromBigIntSigned(mBufferHandle, bigIntHandle)
}

//export v1_5_mBufferToBigFloat
func v1_5_mBufferToBigFloat(context unsafe.Pointer, mBufferHandle int32, bigFloatHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferToBigFloat(mBufferHandle, bigFloatHandle)
}

//export v1_5_mBufferFromBigFloat
func v1_5_mBufferFromBigFloat(context unsafe.Pointer, mBufferHandle int32, bigFloatHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferFromBigFloat(mBufferHandle, bigFloatHandle)
}

//export v1_5_mBufferStorageStore
func v1_5_mBufferStorageStore(context unsafe.Pointer, keyHandle int32, sourceHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferStorageStore(keyHandle, sourceHandle)
}

//export v1_5_mBufferStorageLoad
func v1_5_mBufferStorageLoad(context unsafe.Pointer, keyHandle int32, destinationHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferStorageLoad(keyHandle, destinationHandle)
}

//export v1_5_mBufferStorageLoadFromAddress
func v1_5_mBufferStorageLoadFromAddress(context unsafe.Pointer, addressHandle int32, keyHandle int32, destinationHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.MBufferStorageLoadFromAddress(addressHandle, keyHandle, destinationHandle)
}

//export v1_5_mBufferGetArgument
func v1_5_mBufferGetArgument(context unsafe.Pointer, id int32, destinationHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferGetArgument(id, destinationHandle)
}

//export v1_5_mBufferFinish
func v1_5_mBufferFinish(context unsafe.Pointer, sourceHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferFinish(sourceHandle)
}

//export v1_5_mBufferSetRandom
func v1_5_mBufferSetRandom(context unsafe.Pointer, destinationHandle int32, length int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MBufferSetRandom(destinationHandle, length)
}

//export v1_5_smallIntGetUnsignedArgument
func v1_5_smallIntGetUnsignedArgument(context unsafe.Pointer, id int32) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.SmallIntGetUnsignedArgument(id)
}

//export v1_5_smallIntGetSignedArgument
func v1_5_smallIntGetSignedArgument(context unsafe.Pointer, id int32) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.SmallIntGetSignedArgument(id)
}

//export v1_5_smallIntFinishUnsigned
func v1_5_smallIntFinishUnsigned(context unsafe.Pointer, value int64) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.SmallIntFinishUnsigned(value)
}

//export v1_5_smallIntFinishSigned
func v1_5_smallIntFinishSigned(context unsafe.Pointer, value int64) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.SmallIntFinishSigned(value)
}

//export v1_5_smallIntStorageStoreUnsigned
func v1_5_smallIntStorageStoreUnsigned(context unsafe.Pointer, keyOffset int32, keyLength int32, value int64) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.SmallIntStorageStoreUnsigned(keyOffset, keyLength, value)
}

//export v1_5_smallIntStorageStoreSigned
func v1_5_smallIntStorageStoreSigned(context unsafe.Pointer, keyOffset int32, keyLength int32, value int64) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.SmallIntStorageStoreSigned(keyOffset, keyLength, value)
}

//export v1_5_smallIntStorageLoadUnsigned
func v1_5_smallIntStorageLoadUnsigned(context unsafe.Pointer, keyOffset int32, keyLength int32) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.SmallIntStorageLoadUnsigned(keyOffset, keyLength)
}

//export v1_5_smallIntStorageLoadSigned
func v1_5_smallIntStorageLoadSigned(context unsafe.Pointer, keyOffset int32, keyLength int32) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.SmallIntStorageLoadSigned(keyOffset, keyLength)
}

//export v1_5_int64getArgument
func v1_5_int64getArgument(context unsafe.Pointer, id int32) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.Int64getArgument(id)
}

//export v1_5_int64finish
func v1_5_int64finish(context unsafe.Pointer, value int64) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.Int64finish(value)
}

//export v1_5_int64storageStore
func v1_5_int64storageStore(context unsafe.Pointer, keyOffset int32, keyLength int32, value int64) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.Int64storageStore(keyOffset, keyLength, value)
}

//export v1_5_int64storageLoad
func v1_5_int64storageLoad(context unsafe.Pointer, keyOffset int32, keyLength int32) int64 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.Int64storageLoad(keyOffset, keyLength)
}

//export v1_5_sha256
func v1_5_sha256(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.Sha256(dataOffset, length, resultOffset)
}

//export v1_5_managedSha256
func v1_5_managedSha256(context unsafe.Pointer, inputHandle int32, outputHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedSha256(inputHandle, outputHandle)
}

//export v1_5_keccak256
func v1_5_keccak256(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.Keccak256(dataOffset, length, resultOffset)
}

//export v1_5_managedKeccak256
func v1_5_managedKeccak256(context unsafe.Pointer, inputHandle int32, outputHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedKeccak256(inputHandle, outputHandle)
}

//export v1_5_ripemd160
func v1_5_ripemd160(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.Ripemd160(dataOffset, length, resultOffset)
}

//export v1_5_managedRipemd160
func v1_5_managedRipemd160(context unsafe.Pointer, inputHandle int32, outputHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedRipemd160(inputHandle, outputHandle)
}

//export v1_5_verifyBLS
func v1_5_verifyBLS(context unsafe.Pointer, keyOffset int32, messageOffset int32, messageLength int32, sigOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.VerifyBLS(keyOffset, messageOffset, messageLength, sigOffset)
}

//export v1_5_managedVerifyBLS
func v1_5_managedVerifyBLS(context unsafe.Pointer, keyHandle int32, messageHandle int32, sigHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedVerifyBLS(keyHandle, messageHandle, sigHandle)
}

//export v1_5_verifyEd25519
func v1_5_verifyEd25519(context unsafe.Pointer, keyOffset int32, messageOffset int32, messageLength int32, sigOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.VerifyEd25519(keyOffset, messageOffset, messageLength, sigOffset)
}

//export v1_5_managedVerifyEd25519
func v1_5_managedVerifyEd25519(context unsafe.Pointer, keyHandle int32, messageHandle int32, sigHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedVerifyEd25519(keyHandle, messageHandle, sigHandle)
}

//export v1_5_verifyCustomSecp256k1
func v1_5_verifyCustomSecp256k1(context unsafe.Pointer, keyOffset int32, keyLength int32, messageOffset int32, messageLength int32, sigOffset int32, hashType int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.VerifyCustomSecp256k1(keyOffset, keyLength, messageOffset, messageLength, sigOffset, hashType)
}

//export v1_5_managedVerifyCustomSecp256k1
func v1_5_managedVerifyCustomSecp256k1(context unsafe.Pointer, keyHandle int32, messageHandle int32, sigHandle int32, hashType int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedVerifyCustomSecp256k1(keyHandle, messageHandle, sigHandle, hashType)
}

//export v1_5_verifySecp256k1
func v1_5_verifySecp256k1(context unsafe.Pointer, keyOffset int32, keyLength int32, messageOffset int32, messageLength int32, sigOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.VerifySecp256k1(keyOffset, keyLength, messageOffset, messageLength, sigOffset)
}

//export v1_5_managedVerifySecp256k1
func v1_5_managedVerifySecp256k1(context unsafe.Pointer, keyHandle int32, messageHandle int32, sigHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedVerifySecp256k1(keyHandle, messageHandle, sigHandle)
}

//export v1_5_encodeSecp256k1DerSignature
func v1_5_encodeSecp256k1DerSignature(context unsafe.Pointer, rOffset int32, rLength int32, sOffset int32, sLength int32, sigOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.EncodeSecp256k1DerSignature(rOffset, rLength, sOffset, sLength, sigOffset)
}

//export v1_5_managedEncodeSecp256k1DerSignature
func v1_5_managedEncodeSecp256k1DerSignature(context unsafe.Pointer, rHandle int32, sHandle int32, sigHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedEncodeSecp256k1DerSignature(rHandle, sHandle, sigHandle)
}

//export v1_5_addEC
func v1_5_addEC(context unsafe.Pointer, xResultHandle int32, yResultHandle int32, ecHandle int32, fstPointXHandle int32, fstPointYHandle int32, sndPointXHandle int32, sndPointYHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.AddEC(xResultHandle, yResultHandle, ecHandle, fstPointXHandle, fstPointYHandle, sndPointXHandle, sndPointYHandle)
}

//export v1_5_doubleEC
func v1_5_doubleEC(context unsafe.Pointer, xResultHandle int32, yResultHandle int32, ecHandle int32, pointXHandle int32, pointYHandle int32) {
	callbacks := importsInterfaceFromRaw(context)
	callbacks.DoubleEC(xResultHandle, yResultHandle, ecHandle, pointXHandle, pointYHandle)
}

//export v1_5_isOnCurveEC
func v1_5_isOnCurveEC(context unsafe.Pointer, ecHandle int32, pointXHandle int32, pointYHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.IsOnCurveEC(ecHandle, pointXHandle, pointYHandle)
}

//export v1_5_scalarBaseMultEC
func v1_5_scalarBaseMultEC(context unsafe.Pointer, xResultHandle int32, yResultHandle int32, ecHandle int32, dataOffset int32, length int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ScalarBaseMultEC(xResultHandle, yResultHandle, ecHandle, dataOffset, length)
}

//export v1_5_managedScalarBaseMultEC
func v1_5_managedScalarBaseMultEC(context unsafe.Pointer, xResultHandle int32, yResultHandle int32, ecHandle int32, dataHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedScalarBaseMultEC(xResultHandle, yResultHandle, ecHandle, dataHandle)
}

//export v1_5_scalarMultEC
func v1_5_scalarMultEC(context unsafe.Pointer, xResultHandle int32, yResultHandle int32, ecHandle int32, pointXHandle int32, pointYHandle int32, dataOffset int32, length int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ScalarMultEC(xResultHandle, yResultHandle, ecHandle, pointXHandle, pointYHandle, dataOffset, length)
}

//export v1_5_managedScalarMultEC
func v1_5_managedScalarMultEC(context unsafe.Pointer, xResultHandle int32, yResultHandle int32, ecHandle int32, pointXHandle int32, pointYHandle int32, dataHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedScalarMultEC(xResultHandle, yResultHandle, ecHandle, pointXHandle, pointYHandle, dataHandle)
}

//export v1_5_marshalEC
func v1_5_marshalEC(context unsafe.Pointer, xPairHandle int32, yPairHandle int32, ecHandle int32, resultOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MarshalEC(xPairHandle, yPairHandle, ecHandle, resultOffset)
}

//export v1_5_managedMarshalEC
func v1_5_managedMarshalEC(context unsafe.Pointer, xPairHandle int32, yPairHandle int32, ecHandle int32, resultHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedMarshalEC(xPairHandle, yPairHandle, ecHandle, resultHandle)
}

//export v1_5_marshalCompressedEC
func v1_5_marshalCompressedEC(context unsafe.Pointer, xPairHandle int32, yPairHandle int32, ecHandle int32, resultOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.MarshalCompressedEC(xPairHandle, yPairHandle, ecHandle, resultOffset)
}

//export v1_5_managedMarshalCompressedEC
func v1_5_managedMarshalCompressedEC(context unsafe.Pointer, xPairHandle int32, yPairHandle int32, ecHandle int32, resultHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedMarshalCompressedEC(xPairHandle, yPairHandle, ecHandle, resultHandle)
}

//export v1_5_unmarshalEC
func v1_5_unmarshalEC(context unsafe.Pointer, xResultHandle int32, yResultHandle int32, ecHandle int32, dataOffset int32, length int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.UnmarshalEC(xResultHandle, yResultHandle, ecHandle, dataOffset, length)
}

//export v1_5_managedUnmarshalEC
func v1_5_managedUnmarshalEC(context unsafe.Pointer, xResultHandle int32, yResultHandle int32, ecHandle int32, dataHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedUnmarshalEC(xResultHandle, yResultHandle, ecHandle, dataHandle)
}

//export v1_5_unmarshalCompressedEC
func v1_5_unmarshalCompressedEC(context unsafe.Pointer, xResultHandle int32, yResultHandle int32, ecHandle int32, dataOffset int32, length int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.UnmarshalCompressedEC(xResultHandle, yResultHandle, ecHandle, dataOffset, length)
}

//export v1_5_managedUnmarshalCompressedEC
func v1_5_managedUnmarshalCompressedEC(context unsafe.Pointer, xResultHandle int32, yResultHandle int32, ecHandle int32, dataHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedUnmarshalCompressedEC(xResultHandle, yResultHandle, ecHandle, dataHandle)
}

//export v1_5_generateKeyEC
func v1_5_generateKeyEC(context unsafe.Pointer, xPubKeyHandle int32, yPubKeyHandle int32, ecHandle int32, resultOffset int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GenerateKeyEC(xPubKeyHandle, yPubKeyHandle, ecHandle, resultOffset)
}

//export v1_5_managedGenerateKeyEC
func v1_5_managedGenerateKeyEC(context unsafe.Pointer, xPubKeyHandle int32, yPubKeyHandle int32, ecHandle int32, resultHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedGenerateKeyEC(xPubKeyHandle, yPubKeyHandle, ecHandle, resultHandle)
}

//export v1_5_createEC
func v1_5_createEC(context unsafe.Pointer, dataOffset int32, dataLength int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.CreateEC(dataOffset, dataLength)
}

//export v1_5_managedCreateEC
func v1_5_managedCreateEC(context unsafe.Pointer, dataHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.ManagedCreateEC(dataHandle)
}

//export v1_5_getCurveLengthEC
func v1_5_getCurveLengthEC(context unsafe.Pointer, ecHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetCurveLengthEC(ecHandle)
}

//export v1_5_getPrivKeyByteLengthEC
func v1_5_getPrivKeyByteLengthEC(context unsafe.Pointer, ecHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.GetPrivKeyByteLengthEC(ecHandle)
}

//export v1_5_ellipticCurveGetValues
func v1_5_ellipticCurveGetValues(context unsafe.Pointer, ecHandle int32, fieldOrderHandle int32, basePointOrderHandle int32, eqConstantHandle int32, xBasePointHandle int32, yBasePointHandle int32) int32 {
	callbacks := importsInterfaceFromRaw(context)
	return callbacks.EllipticCurveGetValues(ecHandle, fieldOrderHandle, basePointOrderHandle, eqConstantHandle, xBasePointHandle, yBasePointHandle)
}

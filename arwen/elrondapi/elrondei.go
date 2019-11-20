package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern void getOwner(void *context, int32_t resultOffset);
// extern void getExternalBalance(void *context, int32_t addressOffset, int32_t resultOffset);
// extern int32_t blockHash(void *context, long long nonce, int32_t resultOffset);
// extern int32_t transferValue(void *context, int32_t dstOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
// extern int32_t getArgument(void *context, int32_t id, int32_t argOffset);
// extern int32_t getFunction(void *context, int32_t functionOffset);
// extern int32_t getNumArguments(void *context);
// extern int32_t storageStore(void *context, int32_t keyOffset, int32_t dataOffset, int32_t dataLength);
// extern int32_t storageLoad(void *context, int32_t keyOffset, int32_t dataOffset);
// extern void getCaller(void *context, int32_t resultOffset);
// extern int32_t callValue(void *context, int32_t resultOffset);
// extern void writeLog(void *context, int32_t pointer, int32_t length, int32_t topicPtr, int32_t numTopics);
// extern void returnData(void* context, int32_t dataOffset, int32_t length);
// extern void signalError(void* context);
// extern long long getGasLeft(void *context);
//
// extern int32_t executeOnDestContext(void *context, long long gas, int32_t addressOffset, int32_t valueOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t executeOnSameContext(void *context, long long gas, int32_t addressOffset, int32_t valueOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t delegateExecution(void *context, long long gas, int32_t addressOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t executeReadOnly(void *context, long long gas, int32_t addressOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t createContract(void *context, int32_t valueOffset, int32_t codeOffset, int32_t length, int32_t resultOffset, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
//
// extern int32_t getNumReturnData(void *context);
// extern int32_t getReturnDataSize(void *context, int32_t resultId);
// extern int32_t getReturnData(void *context, int32_t resultId, int32_t dataOffset);
//
// extern long long getBlockTimestamp(void *context);
// extern long long getBlockNonce(void *context);
// extern long long getBlockRound(void *context);
// extern long long getBlockEpoch(void *context);
// extern void getBlockRandomSeed(void *context, int32_t resultOffset);
// extern void getStateRootHash(void *context, int32_t resultOffset);
//
// extern long long getPrevBlockTimestamp(void *context);
// extern long long getPrevBlockNonce(void *context);
// extern long long getPrevBlockRound(void *context);
// extern long long getPrevBlockEpoch(void *context);
// extern void getPrevBlockRandomSeed(void *context, int32_t resultOffset);
//
// extern long long int64getArgument(void *context, int32_t id);
// extern int32_t int64storageStore(void *context, int32_t keyOffset, long long value);
// extern long long int64storageLoad(void *context, int32_t keyOffset);
// extern void int64finish(void* context, long long value);
import "C"

import (
	"math/big"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/go-ext-wasm/wasmer"
)

func ElrondEImports() (*wasmer.Imports, error) {
	imports := wasmer.NewImports()
	imports = imports.Namespace("env")

	imports, err := imports.Append("getOwner", getOwner, C.getOwner)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getExternalBalance", getExternalBalance, C.getExternalBalance)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockHash", blockHash, C.blockHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("transferValue", transferValue, C.transferValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getArgument", getArgument, C.getArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getFunction", getFunction, C.getFunction)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getNumArguments", getNumArguments, C.getNumArguments)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageStore", storageStore, C.storageStore)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageLoad", storageLoad, C.storageLoad)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCaller", getCaller, C.getCaller)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCallValue", callValue, C.callValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("writeLog", writeLog, C.writeLog)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("finish", returnData, C.returnData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("signalError", signalError, C.signalError)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockTimestamp", getBlockTimestamp, C.getBlockTimestamp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockNonce", getBlockNonce, C.getBlockNonce)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockRound", getBlockRound, C.getBlockRound)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockEpoch", getBlockEpoch, C.getBlockEpoch)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockRandomSeed", getBlockRandomSeed, C.getBlockRandomSeed)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getStateRootHash", getStateRootHash, C.getStateRootHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockTimestamp", getPrevBlockTimestamp, C.getPrevBlockTimestamp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockNonce", getPrevBlockNonce, C.getPrevBlockNonce)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockRound", getPrevBlockRound, C.getPrevBlockRound)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockEpoch", getPrevBlockEpoch, C.getPrevBlockEpoch)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockRandomSeed", getPrevBlockRandomSeed, C.getPrevBlockRandomSeed)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getGasLeft", getGasLeft, C.getGasLeft)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("executeOnDestContext", executeOnDestContext, C.executeOnDestContext)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("executeOnSameContext", executeOnSameContext, C.executeOnSameContext)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("delegateExecution", delegateExecution, C.delegateExecution)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("createContract", createContract, C.createContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("executeReadOnly", executeReadOnly, C.executeReadOnly)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getNumReturnData", getNumReturnData, C.getNumReturnData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getReturnDataSize", getReturnDataSize, C.getReturnDataSize)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getReturnData", getReturnData, C.getReturnData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64getArgument", int64getArgument, C.int64getArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64storageStore", int64storageStore, C.int64storageStore)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64storageLoad", int64storageLoad, C.int64storageLoad)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64finish", int64finish, C.int64finish)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export getGasLeft
func getGasLeft(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetGasLeft
	hostContext.UseGas(gasToUse)

	return int64(hostContext.GasLeft())
}

//export getOwner
func getOwner(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	owner := hostContext.GetSCAddress()
	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, owner)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetOwner
	hostContext.UseGas(gasToUse)
}

//export signalError
func signalError(context unsafe.Pointer) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	hostContext.SignalUserError()

	gasToUse := hostContext.GasSchedule().ElrondAPICost.SignalError
	hostContext.UseGas(gasToUse)
}

//export getExternalBalance
func getExternalBalance(context unsafe.Pointer, addressOffset int32, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	address := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	balance := hostContext.GetBalance(address)

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, balance)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetExternalBalance
	hostContext.UseGas(gasToUse)
}

//export blockHash
func blockHash(context unsafe.Pointer, nonce int64, resultOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockHash
	hostContext.UseGas(gasToUse)

	//TODO: change blockchain hook to treat actual nonce - not the offset.
	hash := hostContext.BlockHash(nonce)
	err := arwen.StoreBytes(instCtx.Memory(), resultOffset, hash)
	if err != nil {
		return 1
	}

	return 0
}

//export transferValue
func transferValue(context unsafe.Pointer, destOffset int32, valueOffset int32, dataOffset int32, length int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	send := hostContext.GetSCAddress()
	dest := arwen.LoadBytes(instCtx.Memory(), destOffset, arwen.AddressLen)
	value := arwen.LoadBytes(instCtx.Memory(), valueOffset, arwen.BalanceLen)
	data := arwen.LoadBytes(instCtx.Memory(), dataOffset, length)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.TransferValue
	gasToUse += hostContext.GasSchedule().BaseOperationCost.PersistPerByte * uint64(length)
	hostContext.UseGas(gasToUse)

	hostContext.Transfer(dest, send, big.NewInt(0).SetBytes(value), data)

	return 0
}

//export getArgument
func getArgument(context unsafe.Pointer, id int32, argOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetArgument
	hostContext.UseGas(gasToUse)

	args := hostContext.Arguments()
	if int32(len(args)) <= id {
		return -1
	}

	err := arwen.StoreBytes(instCtx.Memory(), argOffset, args[id])
	if err != nil {
		return -1
	}

	return int32(len(args[id]))
}

//export getFunction
func getFunction(context unsafe.Pointer, functionOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetFunction
	hostContext.UseGas(gasToUse)

	function := hostContext.Function()
	err := arwen.StoreBytes(instCtx.Memory(), functionOffset, []byte(function))
	if err != nil {
		return -1
	}

	return int32(len(function))
}

//export getNumArguments
func getNumArguments(context unsafe.Pointer) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetNumArguments
	hostContext.UseGas(gasToUse)

	return int32(len(hostContext.Arguments()))
}

//export storageStore
func storageStore(context unsafe.Pointer, keyOffset int32, dataOffset int32, dataLength int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	key := arwen.LoadBytes(instCtx.Memory(), keyOffset, arwen.HashLen)
	data := arwen.LoadBytes(instCtx.Memory(), dataOffset, dataLength)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.StorageStore
	hostContext.UseGas(gasToUse)

	return hostContext.SetStorage(hostContext.GetSCAddress(), key, data)
}

//export storageLoad
func storageLoad(context unsafe.Pointer, keyOffset int32, dataOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	key := arwen.LoadBytes(instCtx.Memory(), keyOffset, arwen.HashLen)
	data := hostContext.GetStorage(hostContext.GetSCAddress(), key)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.StorageLoad
	gasToUse += hostContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(data))
	hostContext.UseGas(gasToUse)

	err := arwen.StoreBytes(instCtx.Memory(), dataOffset, data)
	if err != nil {
		return -1
	}

	return int32(len(data))
}

//export getCaller
func getCaller(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	caller := hostContext.GetVMInput().CallerAddr

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, caller)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetCaller
	hostContext.UseGas(gasToUse)
}

//export callValue
func callValue(context unsafe.Pointer, resultOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	value := hostContext.GetVMInput().CallValue.Bytes()
	length := len(value)
	invBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		invBytes[length-i-1] = value[i]
	}

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetCallValue
	hostContext.UseGas(gasToUse)

	err := arwen.StoreBytes(instCtx.Memory(), resultOffset, invBytes)
	if err != nil {
		return -1
	}

	return int32(length)
}

//export writeLog
func writeLog(context unsafe.Pointer, pointer int32, length int32, topicPtr int32, numTopics int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	log := arwen.LoadBytes(instCtx.Memory(), pointer, length)

	topics := make([][]byte, numTopics)
	for i := int32(0); i < numTopics; i++ {
		topics[i] = arwen.LoadBytes(instCtx.Memory(), topicPtr+i*arwen.HashLen, arwen.HashLen)
	}

	hostContext.WriteLog(hostContext.GetSCAddress(), topics, log)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.Log
	gasToUse += hostContext.GasSchedule().BaseOperationCost.PersistPerByte * uint64(numTopics*arwen.HashLen+length)
	hostContext.UseGas(gasToUse)
}

//export getBlockTimestamp
func getBlockTimestamp(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockTimeStamp
	hostContext.UseGas(gasToUse)

	return int64(hostContext.BlockChainHook().CurrentTimeStamp())
}

//export getBlockNonce
func getBlockNonce(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockNonce
	hostContext.UseGas(gasToUse)

	return int64(hostContext.BlockChainHook().CurrentNonce())
}

//export getBlockRound
func getBlockRound(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockRound
	hostContext.UseGas(gasToUse)

	return int64(hostContext.BlockChainHook().CurrentRound())
}

//export getBlockEpoch
func getBlockEpoch(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockEpoch
	hostContext.UseGas(gasToUse)

	return int64(hostContext.BlockChainHook().CurrentEpoch())
}

//export getBlockRandomSeed
func getBlockRandomSeed(context unsafe.Pointer, pointer int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockRandomSeed
	hostContext.UseGas(gasToUse)

	randomSeed := hostContext.BlockChainHook().CurrentRandomSeed()
	_ = arwen.StoreBytes(instCtx.Memory(), pointer, randomSeed)
}

//export getStateRootHash
func getStateRootHash(context unsafe.Pointer, pointer int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetStateRootHash
	hostContext.UseGas(gasToUse)

	stateRootHash := hostContext.BlockChainHook().GetStateRootHash()
	_ = arwen.StoreBytes(instCtx.Memory(), pointer, stateRootHash)
}

//export getPrevBlockTimestamp
func getPrevBlockTimestamp(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockTimeStamp
	hostContext.UseGas(gasToUse)

	return int64(hostContext.BlockChainHook().LastTimeStamp())
}

//export getPrevBlockNonce
func getPrevBlockNonce(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockNonce
	hostContext.UseGas(gasToUse)

	return int64(hostContext.BlockChainHook().LastNonce())
}

//export getPrevBlockRound
func getPrevBlockRound(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockRound
	hostContext.UseGas(gasToUse)

	return int64(hostContext.BlockChainHook().LastRound())
}

//export getPrevBlockEpoch
func getPrevBlockEpoch(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockEpoch
	hostContext.UseGas(gasToUse)

	return int64(hostContext.BlockChainHook().LastEpoch())
}

//export getPrevBlockRandomSeed
func getPrevBlockRandomSeed(context unsafe.Pointer, pointer int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.GetBlockRandomSeed
	hostContext.UseGas(gasToUse)

	randomSeed := hostContext.BlockChainHook().LastRandomSeed()
	_ = arwen.StoreBytes(instCtx.Memory(), pointer, randomSeed)
}

//export returnData
func returnData(context unsafe.Pointer, pointer int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	data := arwen.LoadBytes(instCtx.Memory(), pointer, length)
	hostContext.Finish(data)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.Finish
	gasToUse += hostContext.GasSchedule().BaseOperationCost.PersistPerByte * uint64(length)
	hostContext.UseGas(gasToUse)
}

//export int64getArgument
func int64getArgument(context unsafe.Pointer, id int32) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.Int64GetArgument
	hostContext.UseGas(gasToUse)

	args := hostContext.Arguments()
	if int32(len(args)) <= id {
		return -1
	}

	argBigInt := big.NewInt(0).SetBytes(args[id])
	return argBigInt.Int64()
}

//export int64storageStore
func int64storageStore(context unsafe.Pointer, keyOffset int32, value int64) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	key := arwen.LoadBytes(instCtx.Memory(), keyOffset, arwen.HashLen)
	data := big.NewInt(value)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.Int64StorageStore
	hostContext.UseGas(gasToUse)

	return hostContext.SetStorage(hostContext.GetSCAddress(), key, data.Bytes())
}

//export int64storageLoad
func int64storageLoad(context unsafe.Pointer, keyOffset int32) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	key := arwen.LoadBytes(instCtx.Memory(), keyOffset, arwen.HashLen)
	data := hostContext.GetStorage(hostContext.GetSCAddress(), key)

	bigInt := big.NewInt(0).SetBytes(data)

	gasToUse := hostContext.GasSchedule().ElrondAPICost.Int64StorageLoad
	hostContext.UseGas(gasToUse)

	return bigInt.Int64()
}

//export int64finish
func int64finish(context unsafe.Pointer, value int64) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	hostContext.Finish(big.NewInt(0).SetInt64(value).Bytes())

	gasToUse := hostContext.GasSchedule().ElrondAPICost.Int64Finish
	hostContext.UseGas(gasToUse)
}

//export executeOnSameContext
func executeOnSameContext(
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
	instCtx := wasmer.IntoInstanceContext(context)
	erdContext := arwen.GetErdContext(instCtx.Data())

	send := erdContext.GetSCAddress()
	dest := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	value := arwen.LoadBytes(instCtx.Memory(), valueOffset, arwen.BalanceLen)
	function, data, actualLen := getArgumentsFromMemory(context, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)

	gasToUse := erdContext.GasSchedule().ElrondAPICost.ExecuteOnSameContext
	gasToUse += erdContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(actualLen)
	erdContext.UseGas(gasToUse)

	if erdContext.GasLeft() < uint64(gasLimit) {
		return 1
	}

	bigIntVal := big.NewInt(0).SetBytes(value)
	erdContext.Transfer(dest, send, bigIntVal, nil)

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  send,
			Arguments:   data,
			CallValue:   bigIntVal,
			GasPrice:    0,
			GasProvided: uint64(gasLimit),
		},
		RecipientAddr: dest,
		Function:      function,
	}
	err := erdContext.ExecuteOnDestContext(contractCallInput)
	if err != nil {
		return 1
	}

	return 0
}

//export executeOnDestContext
func executeOnDestContext(
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
	instCtx := wasmer.IntoInstanceContext(context)
	erdContext := arwen.GetErdContext(instCtx.Data())

	send := erdContext.GetSCAddress()
	dest := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	value := arwen.LoadBytes(instCtx.Memory(), valueOffset, arwen.BalanceLen)
	function, data, actualLen := getArgumentsFromMemory(context, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)

	gasToUse := erdContext.GasSchedule().ElrondAPICost.ExecuteOnDestContext
	gasToUse += erdContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(actualLen)
	erdContext.UseGas(gasToUse)

	if erdContext.GasLeft() < uint64(gasLimit) {
		return 1
	}

	erdContext.Transfer(dest, send, big.NewInt(0).SetBytes(value), nil)
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  send,
			Arguments:   data,
			CallValue:   big.NewInt(0).SetBytes(value),
			GasPrice:    0,
			GasProvided: uint64(gasLimit),
		},
		RecipientAddr: dest,
		Function:      function,
	}
	err := erdContext.ExecuteOnDestContext(contractCallInput)
	if err != nil {
		return 1
	}

	return 0
}

func getArgumentsFromMemory(
	context unsafe.Pointer,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) (string, [][]byte, int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	argumentsLengthData := arwen.LoadBytes(instCtx.Memory(), argumentsLengthOffset, numArguments*4)

	currOffset := dataOffset
	data := make([][]byte, numArguments)
	for i := int32(0); i < numArguments; i++ {
		currArgLenData := argumentsLengthData[i*4 : i*4+4]
		actualLen := dataToInt32(currArgLenData)

		data[i] = arwen.LoadBytes(instCtx.Memory(), currOffset, actualLen)
		currOffset += actualLen
	}

	function := arwen.LoadBytes(instCtx.Memory(), functionOffset, functionLength)

	return string(function), data, currOffset - dataOffset
}

//export delegateExecution
func delegateExecution(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	erdContext := arwen.GetErdContext(instCtx.Data())

	address := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.HashLen)
	function, data, actualLen := getArgumentsFromMemory(context, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)

	value := erdContext.GetVMInput().CallValue
	sender := erdContext.GetVMInput().CallerAddr

	gasToUse := erdContext.GasSchedule().ElrondAPICost.DelegateExecution
	gasToUse += erdContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(actualLen)
	erdContext.UseGas(gasToUse)

	if erdContext.GasLeft() < uint64(gasLimit) {
		return 1
	}

	erdContext.Transfer(address, sender, value, nil)
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   data,
			CallValue:   value,
			GasPrice:    0,
			GasProvided: uint64(gasLimit),
		},
		RecipientAddr: address,
		Function:      function,
	}
	err := erdContext.ExecuteOnSameContext(contractCallInput)
	if err != nil {
		return 1
	}

	return 0
}

func dataToInt32(data []byte) int32 {
	actualLen := int32(0)
	for i := len(data) - 1; i >= 0; i-- {
		actualLen = (actualLen << 8) + int32(data[i])
	}

	return actualLen
}

//export executeReadOnly
func executeReadOnly(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	erdContext := arwen.GetErdContext(instCtx.Data())

	address := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.HashLen)
	function, data, actualLen := getArgumentsFromMemory(context, functionOffset, functionLength, numArguments, argumentsLengthOffset, dataOffset)

	value := erdContext.GetVMInput().CallValue
	sender := erdContext.GetVMInput().CallerAddr

	gasToUse := erdContext.GasSchedule().ElrondAPICost.ExecuteReadOnly
	gasToUse += erdContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(actualLen)
	erdContext.UseGas(gasToUse)

	if erdContext.GasLeft() < uint64(gasLimit) {
		return 1
	}

	erdContext.Transfer(address, sender, value, nil)

	erdContext.SetReadOnly(true)
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   data,
			CallValue:   value,
			GasPrice:    0,
			GasProvided: uint64(gasLimit),
		},
		RecipientAddr: address,
		Function:      function,
	}
	err := erdContext.ExecuteOnSameContext(contractCallInput)
	erdContext.SetReadOnly(false)
	if err != nil {
		return 1
	}

	return 0
}

//export createContract
func createContract(
	context unsafe.Pointer,
	valueOffset int32,
	codeOffset int32,
	length int32,
	resultOffset int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	erdContext := arwen.GetErdContext(instCtx.Data())

	sender := erdContext.GetSCAddress()
	value := arwen.LoadBytes(instCtx.Memory(), valueOffset, arwen.BalanceLen)
	code := arwen.LoadBytes(instCtx.Memory(), codeOffset, length)

	_, data, actualLen := getArgumentsFromMemory(context, 0, 0, numArguments, argumentsLengthOffset, dataOffset)

	gasToUse := erdContext.GasSchedule().ElrondAPICost.CreateContract
	gasToUse += erdContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(actualLen)
	erdContext.UseGas(gasToUse)
	gasLimit := erdContext.GasLeft()

	contractCreate := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   data,
			CallValue:   big.NewInt(0).SetBytes(value),
			GasPrice:    0,
			GasProvided: gasLimit,
		},
		ContractCode: code,
	}
	newAddress, err := erdContext.CreateNewContract(contractCreate)
	if err != nil {
		return 1
	}

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, newAddress)

	return 0
}

//export getNumReturnData
func getNumReturnData(context unsafe.Pointer) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	erdContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := erdContext.GasSchedule().ElrondAPICost.GetNumReturnData
	erdContext.UseGas(gasToUse)

	returnData := erdContext.ReturnData()
	return int32(len(returnData))
}

//export getReturnDataSize
func getReturnDataSize(context unsafe.Pointer, resultId int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	erdContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := erdContext.GasSchedule().ElrondAPICost.GetReturnDataSize
	erdContext.UseGas(gasToUse)

	returnData := erdContext.ReturnData()
	if int32(len(returnData)) >= resultId {
		return 0
	}

	return int32(len(returnData[resultId]))
}

//export getReturnData
func getReturnData(context unsafe.Pointer, resultId int32, dataOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	erdContext := arwen.GetErdContext(instCtx.Data())

	gasToUse := erdContext.GasSchedule().ElrondAPICost.GetReturnData
	erdContext.UseGas(gasToUse)

	returnData := erdContext.ReturnData()
	if int32(len(returnData)) >= resultId {
		return 0
	}

	_ = arwen.StoreBytes(instCtx.Memory(), dataOffset, returnData[resultId])
	return int32(len(returnData[resultId]))
}

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
// extern int32_t transferValue(void *context, long long gasLimit, int32_t dstOffset, int32_t sndOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
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
// extern long long getBlockTimestamp(void *context);
//
// extern long long int64getArgument(void *context, int32_t id);
// extern int32_t int64storageStore(void *context, int32_t keyOffset, long long value);
// extern long long int64storageLoad(void *context, int32_t keyOffset);
// extern void int64finish(void* context, long long value);
import "C"

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"math/big"
	"unsafe"

	"github.com/ElrondNetwork/go-ext-wasm/wasmer"
)

func ElrondEImports() (*wasmer.Imports, error) {
	imports := wasmer.NewImports()

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

	imports, err = imports.Append("log", writeLog, C.writeLog)
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

	imports, err = imports.Append("getGasLeft", getGasLeft, C.getGasLeft)
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

	return int64(hostContext.GasLeft())
}

//export getOwner
func getOwner(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	owner := hostContext.GetSCAddress()
	err := arwen.StoreBytes(instCtx.Memory(), resultOffset, owner)
	if err != nil {
	}
}

//export signalError
func signalError(context unsafe.Pointer) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	hostContext.SignalUserError()
}

//export getExternalBalance
func getExternalBalance(context unsafe.Pointer, addressOffset int32, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	address := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	balance := hostContext.GetBalance(address)

	err := arwen.StoreBytes(instCtx.Memory(), resultOffset, balance)
	if err != nil {
	}
}

//export blockHash
func blockHash(context unsafe.Pointer, nonce int64, resultOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	hash := hostContext.BlockHash(nonce)
	err := arwen.StoreBytes(instCtx.Memory(), resultOffset, hash)
	if err != nil {
		return 1
	}

	return 0
}

//export transferValue
func transferValue(context unsafe.Pointer, gasLimit int64, sndOffset int32, destOffset int32, valueOffset int32, dataOffset int32, length int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	send := arwen.LoadBytes(instCtx.Memory(), sndOffset, arwen.AddressLen)
	dest := arwen.LoadBytes(instCtx.Memory(), destOffset, arwen.AddressLen)
	value := arwen.LoadBytes(instCtx.Memory(), valueOffset, arwen.BalanceLen)
	data := arwen.LoadBytes(instCtx.Memory(), dataOffset, length)

	_, err := hostContext.Transfer(dest, send, big.NewInt(0).SetBytes(value), data, gasLimit)
	if err != nil {
		return 1
	}

	return 0
}

//export getArgument
func getArgument(context unsafe.Pointer, id int32, argOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	args := hostContext.Arguments()
	if int32(len(args)) <= id {
		return -1
	}

	err := arwen.StoreBytes(instCtx.Memory(), argOffset, args[id].Bytes())
	if err != nil {
		return -1
	}

	return int32(len(args[id].Bytes()))
}

//export getFunction
func getFunction(context unsafe.Pointer, functionOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

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

	return int32(len(hostContext.Arguments()))
}

//export storageStore
func storageStore(context unsafe.Pointer, keyOffset int32, dataOffset int32, dataLength int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	key := arwen.LoadBytes(instCtx.Memory(), keyOffset, arwen.HashLen)
	data := arwen.LoadBytes(instCtx.Memory(), dataOffset, dataLength)

	return hostContext.SetStorage(hostContext.GetSCAddress(), key, data)
}

//export storageLoad
func storageLoad(context unsafe.Pointer, keyOffset int32, dataOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	key := arwen.LoadBytes(instCtx.Memory(), keyOffset, arwen.HashLen)
	data := hostContext.GetStorage(hostContext.GetSCAddress(), key)

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

	err := arwen.StoreBytes(instCtx.Memory(), resultOffset, caller)
	if err != nil {
	}
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
}

//export getBlockTimestamp
func getBlockTimestamp(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	return int64(hostContext.BlockChainHook().CurrentTimeStamp())
}

//export returnData
func returnData(context unsafe.Pointer, pointer int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	data := arwen.LoadBytes(instCtx.Memory(), pointer, length)
	hostContext.Finish(data)
}

//export int64getArgument
func int64getArgument(context unsafe.Pointer, id int32) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	args := hostContext.Arguments()
	if int32(len(args)) <= id {
		return -1
	}

	return args[id].Int64()
}

//export int64storageStore
func int64storageStore(context unsafe.Pointer, keyOffset int32, value int64) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	key := arwen.LoadBytes(instCtx.Memory(), keyOffset, arwen.HashLen)
	data := big.NewInt(value)

	return hostContext.SetStorage(hostContext.GetSCAddress(), key, data.Bytes())
}

//export int64storageLoad
func int64storageLoad(context unsafe.Pointer, keyOffset int32) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	key := arwen.LoadBytes(instCtx.Memory(), keyOffset, arwen.HashLen)
	data := hostContext.GetStorage(hostContext.GetSCAddress(), key)

	bigInt := big.NewInt(0).SetBytes(data)

	return bigInt.Int64()
}

//export int64finish
func int64finish(context unsafe.Pointer, value int64) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetErdContext(instCtx.Data())

	hostContext.Finish(big.NewInt(0).SetInt64(value).Bytes())
}

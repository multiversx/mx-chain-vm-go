package arwen

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
//
// extern void getOwner(void *context, int32_t resultOffset);
// extern void getExternalBalance(void *context, int32_t addressOffset, int32_t resultOffset);
// extern int32_t blockHash(void *context, long long nonce, int32_t resultOffset);
// extern int32_t transfer(void *context, long long gasLimit, int32_t dstOffset, int32_t sndOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
// extern int32_t getArgument(void *context, int32_t id, int32_t argOffset);
// extern long long getArgumentAsInt64(void *context, int32_t id);
// extern int32_t getFunction(void *context, int32_t functionOffset);
// extern int32_t getNumArguments(void *context);
// extern int32_t storageStore(void *context, int32_t keyOffset, int32_t dataOffset, int32_t dataLength);
// extern int32_t storageLoad(void *context, int32_t keyOffset, int32_t dataOffset);
// extern int32_t storageStoreAsInt64(void *context, int32_t keyOffset, long long value);
// extern long long storageLoadAsInt64(void *context, int32_t keyOffset);
// extern void getCaller(void *context, int32_t resultOffset);
// extern int32_t getCallValue(void *context, int32_t resultOffset);
// extern long long getCallValueAsInt64(void *context);
// extern void logMessage(void *context, int32_t pointer, int32_t length);
// extern void writeLog(void *context, int32_t pointer, int32_t length, int32_t topicPtr, int32_t numTopics);
// extern void finish(void* context, int32_t dataOffset, int32_t length);
// extern long long getBlockTimestamp(void *context);
// extern void signalError(void* context);
import "C"

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"unsafe"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/go-ext-wasm/wasmer"
)

type HostContext interface {
	Arguments() []*big.Int
	Function() string
	AccountExists(addr []byte) bool
	GetStorage(addr []byte, key []byte) []byte
	SetStorage(addr []byte, key []byte, value []byte) int32
	GetBalance(addr []byte) []byte
	GetCodeSize(addr []byte) int
	BlockHash(nonce int64) []byte
	GetCodeHash(addr []byte) []byte
	GetCode(addr []byte) []byte
	SelfDestruct(addr []byte, beneficiary []byte)
	GetVMInput() vmcommon.VMInput
	GetSCAddress() []byte
	WriteLog(addr []byte, topics [][]byte, data []byte)
	Transfer(destination []byte, sender []byte, value *big.Int, input []byte, gas int64) (gasLeft int64, err error)
	Finish(data []byte)
	SignalUserError()
}

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

	imports, err = imports.Append("blockHash", blockHash, C.blockHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("transfer", transfer, C.transfer)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getArgumentAsInt64", getArgumentAsInt64, C.getArgumentAsInt64)
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

	imports, err = imports.Append("storageStoreAsInt64", storageStoreAsInt64, C.storageStoreAsInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageLoadAsInt64", storageLoadAsInt64, C.storageLoadAsInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCaller", getCaller, C.getCaller)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCallValue", getCallValue, C.getCallValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCallValueAsInt64", getCallValueAsInt64, C.getCallValueAsInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("logMessage", logMessage, C.logMessage)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("writeLog", writeLog, C.writeLog)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("finish", finish, C.finish)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockTimestamp", getBlockTimestamp, C.getBlockTimestamp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("signalError", signalError, C.signalError)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

// Write the implementation of the functions, and export it (for cgo).

//export getOwner
func getOwner(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	owner := hostContext.GetSCAddress()
	err := storeBytes(instCtx.Memory(), resultOffset, owner)
	if err != nil {
		fmt.Println("getOwner error: " + err.Error())
	}
	fmt.Println("getOwner " + hex.EncodeToString(owner))
}

//export signalError
func signalError(context unsafe.Pointer) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	hostContext.SignalUserError()
	fmt.Println("signalError called")
}

//export getExternalBalance
func getExternalBalance(context unsafe.Pointer, addressOffset int32, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	address := loadBytes(instCtx.Memory(), addressOffset, addressLen)
	balance := hostContext.GetBalance(address)

	err := storeBytes(instCtx.Memory(), resultOffset, balance)
	if err != nil {
		fmt.Println("getExternalBalance error: " + err.Error())
	}
	fmt.Println("getExternalBalance address: " + hex.EncodeToString(address) + " balance: " + big.NewInt(0).SetBytes(balance).String())
}

//export blockHash
func blockHash(context unsafe.Pointer, nonce int64, resultOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	hash := hostContext.BlockHash(nonce)
	err := storeBytes(instCtx.Memory(), resultOffset, hash)
	if err != nil {
		fmt.Println("blockHash error: " + err.Error())
		return 1
	}
	fmt.Println("blockHash " + hex.EncodeToString(hash))
	return 0
}

//export transfer
func transfer(context unsafe.Pointer, gasLimit int64, sndOffset int32, destOffset int32, valueOffset int32, dataOffset int32, length int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	send := loadBytes(instCtx.Memory(), sndOffset, addressLen)
	dest := loadBytes(instCtx.Memory(), destOffset, addressLen)
	value := loadBytes(instCtx.Memory(), valueOffset, balanceLen)
	data := loadBytes(instCtx.Memory(), dataOffset, length)

	fmt.Println("transfer send: " + hex.EncodeToString(send) + " dest: " + hex.EncodeToString(dest) + " value: " + string(value) + " data: " + string(data))

	_, err := hostContext.Transfer(dest, send, big.NewInt(0).SetBytes(value), data, gasLimit)
	if err != nil {
		fmt.Println("transfer error: " + err.Error())
		return 1
	}

	fmt.Println("transfer succeed")
	return 0
}

//export getArgument
func getArgument(context unsafe.Pointer, id int32, argOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	args := hostContext.Arguments()
	if int32(len(args)) <= id {
		fmt.Println("getArgument id invalid")
		return -1
	}

	err := storeBytes(instCtx.Memory(), argOffset, args[id].Bytes())
	if err != nil {
		fmt.Println("getArgument error " + err.Error())
		return -1
	}

	fmt.Println("getArgument value: " + hex.EncodeToString(args[id].Bytes()))
	return int32(len(args[id].Bytes()))
}

//export getArgumentAsInt64
func getArgumentAsInt64(context unsafe.Pointer, id int32) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	args := hostContext.Arguments()
	if int32(len(args)) <= id {
		fmt.Println("getArgument id invalid")
		return -1
	}

	fmt.Println("getArgument value: ", args[id].Int64())
	return args[id].Int64()
}

//export getFunction
func getFunction(context unsafe.Pointer, functionOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	function := hostContext.Function()
	err := storeBytes(instCtx.Memory(), functionOffset, []byte(function))
	if err != nil {
		fmt.Println("getFunction error: ", err.Error())
		return -1
	}

	fmt.Println("getFunction name: " + function)
	return int32(len(function))
}

//export getNumArguments
func getNumArguments(context unsafe.Pointer) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	fmt.Println("getNumArguments ", len(hostContext.Arguments()))
	return int32(len(hostContext.Arguments()))
}

//export storageStore
func storageStore(context unsafe.Pointer, keyOffset int32, dataOffset int32, dataLength int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	key := loadBytes(instCtx.Memory(), keyOffset, hashLen)
	data := loadBytes(instCtx.Memory(), dataOffset, dataLength)

	fmt.Println("storageStore key: " + string(key) + " value: " + string(data))
	return hostContext.SetStorage(hostContext.GetSCAddress(), key, data)
}

//export storageLoad
func storageLoad(context unsafe.Pointer, keyOffset int32, dataOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	key := loadBytes(instCtx.Memory(), keyOffset, hashLen)
	data := hostContext.GetStorage(hostContext.GetSCAddress(), key)

	err := storeBytes(instCtx.Memory(), dataOffset, data)
	if err != nil {
		fmt.Println("storageLoad error: " + err.Error())
		return -1
	}

	fmt.Println("storageLoad key: " + string(key) + " value: " + string(data))
	return int32(len(data))
}

//export storageStoreAsInt64
func storageStoreAsInt64(context unsafe.Pointer, keyOffset int32, value int64) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	key := loadBytes(instCtx.Memory(), keyOffset, hashLen)
	data := big.NewInt(value)

	fmt.Println("storageStoreAsInt64 key: "+string(key)+"value: ", data.Int64())
	return hostContext.SetStorage(hostContext.GetSCAddress(), key, data.Bytes())
}

//export storageLoadAsInt64
func storageLoadAsInt64(context unsafe.Pointer, keyOffset int32) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	key := loadBytes(instCtx.Memory(), keyOffset, hashLen)
	data := hostContext.GetStorage(hostContext.GetSCAddress(), key)

	bigInt := big.NewInt(0).SetBytes(data)
	fmt.Println("storageLoadAsInt64 "+string(key)+" value: ", bigInt.Int64())

	return bigInt.Int64()
}

//export getCaller
func getCaller(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	caller := hostContext.GetVMInput().CallerAddr

	err := storeBytes(instCtx.Memory(), resultOffset, caller)
	if err != nil {
		fmt.Println("getCaller error: " + err.Error())
	}
	fmt.Println("getCaller " + hex.EncodeToString(caller))
}

//export getCallValue
func getCallValue(context unsafe.Pointer, resultOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	value := hostContext.GetVMInput().CallValue.Bytes()
	length := len(value)
	invBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		invBytes[length-i-1] = value[i]
	}

	err := storeBytes(instCtx.Memory(), resultOffset, invBytes)
	if err != nil {
		fmt.Println("getCallValue error " + err.Error())
		return -1
	}

	fmt.Println("getCallValue ", hostContext.GetVMInput().CallValue)
	return int32(length)
}

//export getCallValueAsInt64
func getCallValueAsInt64(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	fmt.Println("getCallValueAsInt64 ", hostContext.GetVMInput().CallValue.Int64())
	return hostContext.GetVMInput().CallValue.Int64()
}

//export logMessage
func logMessage(context unsafe.Pointer, pointer int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	log := loadBytes(instCtx.Memory(), pointer, length)
	fmt.Println("logMessage: " + string(log))
}

//export writeLog
func writeLog(context unsafe.Pointer, pointer int32, length int32, topicPtr int32, numTopics int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	log := loadBytes(instCtx.Memory(), pointer, length)

	topics := make([][]byte, numTopics)
	fmt.Println("writeLog: ")
	for i := int32(0); i < numTopics; i++ {
		topics[i] = loadBytes(instCtx.Memory(), topicPtr+i*hashLen, hashLen)
		fmt.Println("topics: " + string(topics[i]))
	}

	fmt.Print("log: " + string(log))
	hostContext.WriteLog(hostContext.GetSCAddress(), topics, log)
}

//export finish
func finish(context unsafe.Pointer, pointer int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	data := loadBytes(instCtx.Memory(), pointer, length)
	fmt.Println("finish: " + string(data))
	hostContext.Finish(data)
}

//export getBlockTimestamp
func getBlockTimestamp(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	fmt.Println("getBlockTimestamp ", hostContext.GetVMInput().Header.Timestamp.Int64())
	return hostContext.GetVMInput().Header.Timestamp.Int64()
}

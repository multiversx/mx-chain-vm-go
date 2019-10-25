package arwen

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
// typedef int uint32_t;
// typedef unsigned long long uint64_t;
// extern void getOwner(void *context, int32_t resultOffset);
// extern void loadBalance(void *context, int32_t addressOffset, int32_t result);
// extern int32_t blockHash(void *context, long long nonce, int32_t resultOffset);
// extern int32_t transfer(void *context, long long gasLimit, int32_t dstOffset, int32_t sndOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
// extern int32_t loadArgumentAsBytes(void *context, int32_t id, int32_t argOffset);
// extern void loadArgumentAsBig(void *context, int32_t id, int32_t destination);
// extern long long getArgumentAsInt64(void *context, int32_t id);
// extern int32_t getFunction(void *context, int32_t functionOffset);
// extern int32_t getNumArguments(void *context);
// extern int32_t storageStore(void *context, int32_t keyOffset, int32_t dataOffset, int32_t dataLength);
// extern int32_t storageLoad(void *context, int32_t keyOffset, int32_t dataOffset);
// extern int32_t storageStoreAsBigInt(void *context, int32_t keyOffset, int32_t source);
// extern int32_t storageLoadAsBigInt(void *context, int32_t keyOffset, int32_t destination);
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
// extern int32_t bigInsert(void* context, int32_t smallValue);
// extern int32_t bigByteLength(void* context, int32_t reference);
// extern int32_t bigGetBytes(void* context, int32_t reference);
// extern void bigSetBytes(void* context, int32_t destination, int32_t byteOffset, int32_t byteLength);
// extern void bigAdd(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void bigSub(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void debugPrintBig(void* context, int32_t value);
// extern void debugPrintInt32(void* context, int32_t value);
import "C"

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"unsafe"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/go-ext-wasm/wasmer"
)

// BigIntHandle is the type we use to represent a reference to a big int in the host.
type BigIntHandle = int32

// HostContext abstracts away the blockchain functionality from wasmer.
type HostContext interface {
	Arguments() []*big.Int
	Function() string
	AccountExists(addr []byte) bool
	GetStorage(addr []byte, key []byte) []byte
	SetStorage(addr []byte, key []byte, value []byte) int32
	LoadBalance(addr []byte, destination BigIntHandle)
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

	BigInsertInt64(smallValue int64) BigIntHandle
	BigUpdate(destination BigIntHandle, newValue *big.Int)
	BigByteLength(reference BigIntHandle) int32
	BigGetBytes(reference BigIntHandle) []byte
	GetNextAllocMemIndex(allocSize int32, totalMemSize int32) (newIndex int32)
	BigSetBytes(destination BigIntHandle, bytes []byte)
	BigAdd(destination, op1, op2 BigIntHandle)
	BigSub(destination, op1, _op2 BigIntHandle)
	DebugPrintBig(value BigIntHandle)
}

func ElrondEImports() (*wasmer.Imports, error) {
	imports := wasmer.NewImports()

	imports, err := imports.Append("getOwner", getOwner, C.getOwner)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("loadBalance", loadBalance, C.loadBalance)
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

	imports, err = imports.Append("loadArgumentAsBytes", loadArgumentAsBytes, C.loadArgumentAsBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("loadArgumentAsBig", loadArgumentAsBig, C.loadArgumentAsBig)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getArgumentAsInt64", getArgumentAsInt64, C.getArgumentAsInt64)
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

	imports, err = imports.Append("storageStoreAsBigInt", storageStoreAsBigInt, C.storageStoreAsBigInt)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageLoadAsBigInt", storageLoadAsBigInt, C.storageLoadAsBigInt)
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

	imports, err = imports.Append("bigInsert", bigInsert, C.bigInsert)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigByteLength", bigByteLength, C.bigByteLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigGetBytes", bigGetBytes, C.bigGetBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigSetBytes", bigSetBytes, C.bigSetBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigAdd", bigAdd, C.bigAdd)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigSub", bigSub, C.bigSub)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("debugPrintBig", debugPrintBig, C.debugPrintBig)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("debugPrintInt32", debugPrintInt32, C.debugPrintInt32)
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

//export loadBalance
func loadBalance(context unsafe.Pointer, addressOffset int32, result int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	address := loadBytes(instCtx.Memory(), addressOffset, addressLen)
	hostContext.LoadBalance(address, result)
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

//export loadArgumentAsBytes
func loadArgumentAsBytes(context unsafe.Pointer, id int32, argOffset int32) int32 {
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

	fmt.Printf("argument #%d (bytes): %s\n", id, hex.EncodeToString(args[id].Bytes()))
	return int32(len(args[id].Bytes()))
}

//export loadArgumentAsBig
func loadArgumentAsBig(context unsafe.Pointer, id int32, destination int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	args := hostContext.Arguments()
	if int32(len(args)) <= id {
		fmt.Println("getArgument id invalid")
		return
	}

	hostContext.BigUpdate(destination, args[id])

	fmt.Printf("argument #%d (big int): %d\n", id, args[id])
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

	fmt.Printf("argument #%d (int64): %d\n", id, args[id].Int64())
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

	fmt.Printf("storageStore key: %s  value (bytes): %d\n", hex.EncodeToString(key), big.NewInt(0).SetBytes(data))
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

	fmt.Println("storageLoad key: "+string(key)+" value: ", big.NewInt(0).SetBytes(data))
	return int32(len(data))
}

//export storageStoreAsBigInt
func storageStoreAsBigInt(context unsafe.Pointer, keyOffset int32, source int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	key := loadBytes(instCtx.Memory(), keyOffset, hashLen)
	bytes := hostContext.BigGetBytes(source)

	fmt.Printf("storageStore key: %s  value (big int): %d\n", hex.EncodeToString(key), big.NewInt(0).SetBytes(bytes))
	return hostContext.SetStorage(hostContext.GetSCAddress(), key, bytes)
}

//export storageLoadAsBigInt
func storageLoadAsBigInt(context unsafe.Pointer, keyOffset int32, destination int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	key := loadBytes(instCtx.Memory(), keyOffset, hashLen)
	bytes := hostContext.GetStorage(hostContext.GetSCAddress(), key)

	hostContext.BigSetBytes(destination, bytes)

	fmt.Printf("storageLoad key: %s  value (big int): %d\n", hex.EncodeToString(key), big.NewInt(0).SetBytes(bytes))
	return int32(len(bytes))
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
	fmt.Println("getCaller " + string(caller))
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
	fmt.Println("finish: ", big.NewInt(0).SetBytes(data))
	hostContext.Finish(data)
}

//export getBlockTimestamp
func getBlockTimestamp(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	fmt.Println("getBlockTimestamp ", hostContext.GetVMInput().Header.Timestamp.Int64())
	return hostContext.GetVMInput().Header.Timestamp.Int64()
}

//export bigInsert
func bigInsert(context unsafe.Pointer, smallValue int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	return hostContext.BigInsertInt64(int64(smallValue))
}

//export bigByteLength
func bigByteLength(context unsafe.Pointer, reference int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	return hostContext.BigByteLength(reference)
}

//export bigGetBytes
func bigGetBytes(context unsafe.Pointer, reference int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	bytes := hostContext.BigGetBytes(reference)
	newIndex := hostContext.GetNextAllocMemIndex(int32(len(bytes)), int32(instCtx.Memory().Length()))

	err := storeBytes(instCtx.Memory(), newIndex, bytes)
	if err != nil {
		fmt.Println("bigGetBytes error: " + err.Error())
	}

	return newIndex
}

//export bigSetBytes
func bigSetBytes(context unsafe.Pointer, destination int32, byteOffset int32, byteLength int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	bytes := loadBytes(instCtx.Memory(), byteOffset, byteLength)
	hostContext.BigSetBytes(destination, bytes)
}

//export bigAdd
func bigAdd(context unsafe.Pointer, destination, op1, op2 int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	hostContext.BigAdd(destination, op1, op2)
}

//export bigSub
func bigSub(context unsafe.Pointer, destination, op1, op2 int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	hostContext.BigSub(destination, op1, op2)
}

//export debugPrintBig
func debugPrintBig(context unsafe.Pointer, handle int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	hostContext.DebugPrintBig(handle)
}

//export debugPrintInt32
func debugPrintInt32(context unsafe.Pointer, value int32) {
	fmt.Printf(">>> Int32: %d\n", value)
}

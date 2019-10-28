package arwen

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
// typedef int uint32_t;
// typedef unsigned long long uint64_t;
//
// extern void getOwner(void *context, int32_t resultOffset);
// extern void getExternalBalance(void *context, int32_t addressOffset, int32_t resultOffset);
// extern int32_t getBlockHash(void *context, long long nonce, int32_t resultOffset);
// extern int32_t transfer(void *context, long long gasLimit, int32_t dstOffset, int32_t sndOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
// extern int32_t getArgument(void *context, int32_t id, int32_t argOffset);
// extern int32_t getFunction(void *context, int32_t functionOffset);
// extern int32_t getNumArguments(void *context);
// extern int32_t storageStore(void *context, int32_t keyOffset, int32_t dataOffset, int32_t dataLength);
// extern int32_t storageLoad(void *context, int32_t keyOffset, int32_t dataOffset);
// extern void getCaller(void *context, int32_t resultOffset);
// extern int32_t getCallValue(void *context, int32_t resultOffset);
// extern void writeLog(void *context, int32_t pointer, int32_t length, int32_t topicPtr, int32_t numTopics);
// extern void finish(void* context, int32_t dataOffset, int32_t length);
// extern void signalError(void* context);
// extern long long getGasLeft(void *context);
// extern long long getBlockTimestamp(void *context);
//
// extern long long int64getArgument(void *context, int32_t id);
// extern int32_t int64storageStore(void *context, int32_t keyOffset, long long value);
// extern long long int64storageLoad(void *context, int32_t keyOffset);
// extern void int64finish(void* context, long long value);
//
// extern int32_t bigIntNew(void* context, int32_t smallValue);
// extern int32_t bigIntByteLength(void* context, int32_t reference);
// extern int32_t bigIntGetBytes(void* context, int32_t reference, int32_t byteOffset);
// extern void bigIntSetBytes(void* context, int32_t destination, int32_t byteOffset, int32_t byteLength);
// extern int32_t bigIntIsInt64(void* context, int32_t reference);
// extern long long bigIntGetInt64(void* context, int32_t reference);
// extern void bigIntSetInt64(void* context, int32_t destination, long long value);
// extern void bigIntAdd(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void bigIntSub(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void bigIntMul(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern int32_t bigIntCmp(void* context, int32_t op1, int32_t op2);
// extern void bigIntFinish(void* context, int32_t reference);
// extern int32_t bigIntstorageStore(void *context, int32_t keyOffset, int32_t source);
// extern int32_t bigIntstorageLoad(void *context, int32_t keyOffset, int32_t destination);
// extern void bigIntgetArgument(void *context, int32_t id, int32_t destination);
// extern void bigIntgetCallValue(void *context, int32_t destination);
// extern void bigIntgetExternalBalance(void *context, int32_t addressOffset, int32_t result);
import "C"

import (
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
	SignalUserError()
	Finish(data []byte)

	BigInsertInt64(smallValue int64) BigIntHandle
	BigUpdate(destination BigIntHandle, newValue *big.Int)
	BigGet(reference BigIntHandle) *big.Int
	BigByteLength(reference BigIntHandle) int32
	BigGetBytes(reference BigIntHandle) []byte
	BigSetBytes(destination BigIntHandle, bytes []byte)
	BigIsInt64(destination BigIntHandle) bool
	BigGetInt64(destination BigIntHandle) int64
	BigSetInt64(destination BigIntHandle, value int64)
	BigAdd(destination, op1, op2 BigIntHandle)
	BigSub(destination, op1, op2 BigIntHandle)
	BigMul(destination, op1, op2 BigIntHandle)
	BigCmp(op1, op2 BigIntHandle) int
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

	imports, err = imports.Append("getBlockHash", getBlockHash, C.getBlockHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("transfer", transfer, C.transfer)
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

	imports, err = imports.Append("getCallValue", getCallValue, C.getCallValue)
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

	imports, err = imports.Append("bigIntNew", bigIntNew, C.bigIntNew)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntByteLength", bigIntByteLength, C.bigIntByteLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetBytes", bigIntGetBytes, C.bigIntGetBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSetBytes", bigIntSetBytes, C.bigIntSetBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntIsInt64", bigIntIsInt64, C.bigIntIsInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetInt64", bigIntGetInt64, C.bigIntGetInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSetInt64", bigIntSetInt64, C.bigIntSetInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntAdd", bigIntAdd, C.bigIntAdd)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSub", bigIntSub, C.bigIntSub)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntMul", bigIntMul, C.bigIntMul)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntCmp", bigIntCmp, C.bigIntCmp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntFinish", bigIntFinish, C.bigIntFinish)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntstorageStore", bigIntstorageStore, C.bigIntstorageStore)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntstorageLoad", bigIntstorageLoad, C.bigIntstorageLoad)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntgetArgument", bigIntgetArgument, C.bigIntgetArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntgetCallValue", bigIntgetCallValue, C.bigIntgetCallValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntgetExternalBalance", bigIntgetExternalBalance, C.bigIntgetExternalBalance)
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
	_ = storeBytes(instCtx.Memory(), resultOffset, owner)
}

//export signalError
func signalError(context unsafe.Pointer) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	hostContext.SignalUserError()
}

//export getExternalBalance
func getExternalBalance(context unsafe.Pointer, addressOffset int32, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	address := loadBytes(instCtx.Memory(), addressOffset, addressLen)
	balance := hostContext.GetBalance(address)
	_ = storeBytes(instCtx.Memory(), resultOffset, balance)
}

//export getGasLeft
func getGasLeft(context unsafe.Pointer) int64 {
	// TODO: implement
	return 100000
}

//export getBlockHash
func getBlockHash(context unsafe.Pointer, nonce int64, resultOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	hash := hostContext.BlockHash(nonce)
	err := storeBytes(instCtx.Memory(), resultOffset, hash)
	if err != nil {
		return 1
	}

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

	_, err := hostContext.Transfer(dest, send, big.NewInt(0).SetBytes(value), data, gasLimit)
	if err != nil {
		return 1
	}

	return 0
}

//export getArgument
func getArgument(context unsafe.Pointer, id int32, argOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	args := hostContext.Arguments()
	if int32(len(args)) <= id {
		return -1
	}

	err := storeBytes(instCtx.Memory(), argOffset, args[id].Bytes())
	if err != nil {
		return -1
	}

	return int32(len(args[id].Bytes()))
}

//export getFunction
func getFunction(context unsafe.Pointer, functionOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	function := hostContext.Function()
	err := storeBytes(instCtx.Memory(), functionOffset, []byte(function))
	if err != nil {
		return -1
	}

	return int32(len(function))
}

//export getNumArguments
func getNumArguments(context unsafe.Pointer) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	return int32(len(hostContext.Arguments()))
}

//export storageStore
func storageStore(context unsafe.Pointer, keyOffset int32, dataOffset int32, dataLength int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	key := loadBytes(instCtx.Memory(), keyOffset, hashLen)
	data := loadBytes(instCtx.Memory(), dataOffset, dataLength)

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
		return -1
	}

	return int32(len(data))
}

//export getCaller
func getCaller(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	caller := hostContext.GetVMInput().CallerAddr

	_ = storeBytes(instCtx.Memory(), resultOffset, caller)

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
		return -1
	}

	return int32(length)
}

//export writeLog
func writeLog(context unsafe.Pointer, pointer int32, length int32, topicPtr int32, numTopics int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	log := loadBytes(instCtx.Memory(), pointer, length)

	topics := make([][]byte, numTopics)
	for i := int32(0); i < numTopics; i++ {
		topics[i] = loadBytes(instCtx.Memory(), topicPtr+i*hashLen, hashLen)
	}

	hostContext.WriteLog(hostContext.GetSCAddress(), topics, log)
}

//export getBlockTimestamp
func getBlockTimestamp(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	return hostContext.GetVMInput().Header.Timestamp.Int64()
}

//export finish
func finish(context unsafe.Pointer, pointer int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	data := loadBytes(instCtx.Memory(), pointer, length)
	hostContext.Finish(data)
}

//export int64getArgument
func int64getArgument(context unsafe.Pointer, id int32) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	args := hostContext.Arguments()
	if int32(len(args)) <= id {
		return -1
	}

	return args[id].Int64()
}

//export int64storageStore
func int64storageStore(context unsafe.Pointer, keyOffset int32, value int64) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	key := loadBytes(instCtx.Memory(), keyOffset, hashLen)
	data := big.NewInt(value)

	return hostContext.SetStorage(hostContext.GetSCAddress(), key, data.Bytes())
}

//export int64storageLoad
func int64storageLoad(context unsafe.Pointer, keyOffset int32) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	key := loadBytes(instCtx.Memory(), keyOffset, hashLen)
	data := hostContext.GetStorage(hostContext.GetSCAddress(), key)

	bigInt := big.NewInt(0).SetBytes(data)

	return bigInt.Int64()
}

//export int64finish
func int64finish(context unsafe.Pointer, value int64) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	hostContext.Finish(big.NewInt(0).SetInt64(value).Bytes())
}

//export bigIntgetArgument
func bigIntgetArgument(context unsafe.Pointer, id int32, destination int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	args := hostContext.Arguments()
	if int32(len(args)) <= id {
		return
	}

	hostContext.BigUpdate(destination, args[id])
}

//export bigIntstorageStore
func bigIntstorageStore(context unsafe.Pointer, keyOffset int32, source int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	key := loadBytes(instCtx.Memory(), keyOffset, hashLen)
	bytes := hostContext.BigGetBytes(source)

	return hostContext.SetStorage(hostContext.GetSCAddress(), key, bytes)
}

//export bigIntstorageLoad
func bigIntstorageLoad(context unsafe.Pointer, keyOffset int32, destination int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	key := loadBytes(instCtx.Memory(), keyOffset, hashLen)
	bytes := hostContext.GetStorage(hostContext.GetSCAddress(), key)

	hostContext.BigSetBytes(destination, bytes)

	return int32(len(bytes))
}

//export bigIntgetCallValue
func bigIntgetCallValue(context unsafe.Pointer, destination int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	hostContext.BigUpdate(destination, hostContext.GetVMInput().CallValue)
}

//export bigIntgetExternalBalance
func bigIntgetExternalBalance(context unsafe.Pointer, addressOffset int32, result int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	address := loadBytes(instCtx.Memory(), addressOffset, addressLen)
	balance := hostContext.GetBalance(address)
	hostContext.BigUpdate(result, big.NewInt(0).SetBytes(balance))
}

//export bigIntNew
func bigIntNew(context unsafe.Pointer, smallValue int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	return hostContext.BigInsertInt64(int64(smallValue))
}

//export bigIntByteLength
func bigIntByteLength(context unsafe.Pointer, reference int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	return hostContext.BigByteLength(reference)
}

//export bigIntGetBytes
func bigIntGetBytes(context unsafe.Pointer, reference int32, byteOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	bytes := hostContext.BigGetBytes(reference)

	err := storeBytes(instCtx.Memory(), byteOffset, bytes)
	if err != nil {
	}

	return int32(len(bytes))
}

//export bigIntSetBytes
func bigIntSetBytes(context unsafe.Pointer, destination int32, byteOffset int32, byteLength int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	bytes := loadBytes(instCtx.Memory(), byteOffset, byteLength)
	hostContext.BigSetBytes(destination, bytes)
}

//export bigIntIsInt64
func bigIntIsInt64(context unsafe.Pointer, destination int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	if hostContext.BigIsInt64(destination) {
		return 1
	}
	return 0
}

//export bigIntGetInt64
func bigIntGetInt64(context unsafe.Pointer, destination int32) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	return hostContext.BigGetInt64(destination)
}

//export bigIntSetInt64
func bigIntSetInt64(context unsafe.Pointer, destination int32, value int64) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	hostContext.BigSetInt64(destination, value)
}

//export bigIntAdd
func bigIntAdd(context unsafe.Pointer, destination, op1, op2 int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	hostContext.BigAdd(destination, op1, op2)
}

//export bigIntSub
func bigIntSub(context unsafe.Pointer, destination, op1, op2 int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	hostContext.BigSub(destination, op1, op2)
}

//export bigIntMul
func bigIntMul(context unsafe.Pointer, destination, op1, op2 int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	hostContext.BigMul(destination, op1, op2)
}

//export bigIntCmp
func bigIntCmp(context unsafe.Pointer, op1, op2 int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	return int32(hostContext.BigCmp(op1, op2))
}

//export bigIntFinish
func bigIntFinish(context unsafe.Pointer, reference int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	output := hostContext.BigGetBytes(reference)
	hostContext.Finish(output)
}

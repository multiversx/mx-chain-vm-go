package ethapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char u8;
// typedef int i32;
// typedef int i32ptr;
// extern void useGas(void *context, long long  gas);
// extern void getAddress(void *context, i32ptr resultOffset);
// extern void getExternalBalance(void *context, i32ptr addressOffset, i32ptr resultOffset);
// extern i32 getBlockHash(void *context, long long number, i32ptr resultOffset);
// extern i32 call(void *context, long long gas, i32ptr addressOffset, i32ptr valueOffset, i32ptr dataOffset, i32 dataLength);
// extern i32 getCallDataSize(void *context);
// extern void callDataCopy(void *context, i32ptr resultsOffset, i32ptr dataOffset, i32 length);
// extern i32 callCode(void *context, long long gas, i32ptr addressOffset, i32ptr valueOffset, i32ptr dataOffset, i32 dataLength);
// extern i32 callDelegate(void *context, long long gas, i32ptr addressOffset, i32ptr dataOffset, i32 dataLength);
// extern i32 callStatic(void *context, long long gas, i32ptr addressOffset, i32ptr dataOffset, i32 dataLength);
// extern void storageStore(void *context, i32ptr pathOffset, i32ptr valueOffset);
// extern void storageLoad(void *context, i32ptr pathOffset, i32ptr resultOffset);
// extern void getCaller(void *context, i32ptr resultOffset);
// extern void getCallValue(void *context, i32ptr resultOffset);
// extern void codeCopy(void *context, i32ptr resultOffset, i32 codeOffset, i32 length);
// extern i32 getCodeSize(void *context);
// extern void getBlockCoinbase(void *context, i32ptr resultOffset);
// extern i32 create(void *context, i32ptr valueoffset, i32ptr dataOffset, i32 length, i32ptr resultsOffset);
// extern void getBlockDifficulty(void *context, i32ptr resultOffset);
// extern void externalCodeCopy(void *context, i32ptr addressOffset, i32ptr resultOffset, i32 codeOffset, i32 length);
// extern i32 getExternalCodeSize(void *context, i32ptr addressOffset);
// extern long long getGasLeft(void *context);
// extern long long getBlockGasLimit(void *context);
// extern void getTxGasPrice(void *context, i32ptr valueOffset);
// extern void logTopics(void *context, i32ptr dataOffset, i32 length, i32 numberOftopics, i32ptr topic1, i32ptr topic2, i32ptr topic3, i32ptr topic4);
// extern long long getBlockNumber(void *context);
// extern void getTxOrigin(void *context, i32ptr resultOffset);
// extern void finish(void *context, i32ptr dataOffset, i32 length);
// extern void revert(void *context, i32ptr dataOffset, i32 length);
// extern i32 getReturnDataSize(void *context);
// extern void returnDataCopy(void *context, i32ptr resultOffset, i32 dataOffset, i32 length);
// extern void selfDestruct(void *context, i32ptr addressOffset);
// extern long long getBlockTimestamp(void *context);
import "C"
import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/go-ext-wasm/wasmer"
	"math/big"
	"unsafe"
)

func EthereumImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("ethereum")

	imports, err := imports.Append("useGas", useGas, C.useGas)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getAddress", getAddress, C.getAddress)
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

	imports, err = imports.Append("call", call, C.call)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("callDataCopy", callDataCopy, C.callDataCopy)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCallDataSize", getCallDataSize, C.getCallDataSize)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("callCode", callCode, C.callCode)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("callDelegate", callDelegate, C.callDelegate)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("callStatic", callStatic, C.callStatic)
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

	imports, err = imports.Append("codeCopy", codeCopy, C.codeCopy)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCodeSize", getCodeSize, C.getCodeSize)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockCoinbase", getBlockCoinbase, C.getBlockCoinbase)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("create", create, C.create)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockDifficulty", getBlockDifficulty, C.getBlockDifficulty)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("externalCodeCopy", externalCodeCopy, C.externalCodeCopy)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getExternalCodeSize", getExternalCodeSize, C.getExternalCodeSize)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getGasLeft", getGasLeft, C.getGasLeft)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockGasLimit", getBlockGasLimit, C.getBlockGasLimit)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getTxGasPrice", getTxGasPrice, C.getTxGasPrice)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("logTopics", logTopics, C.logTopics)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockNumber", getBlockNumber, C.getBlockNumber)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getTxOrigin", getTxOrigin, C.getTxOrigin)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("finish", finish, C.finish)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("revert", revert, C.revert)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getReturnDataSize", getReturnDataSize, C.getReturnDataSize)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("returnDataCopy", returnDataCopy, C.returnDataCopy)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("selfDestruct", selfDestruct, C.selfDestruct)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockTimestamp", getBlockTimestamp, C.getBlockTimestamp)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export useGas
func useGas(context unsafe.Pointer, useGas int64) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	ethContext.UseGas(useGas)
}

//export getAddress
func getAddress(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, hostContext.GetSCAddress())
}

//export getExternalBalance
func getExternalBalance(context unsafe.Pointer, addressOffset int32, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	address := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	balance := hostContext.GetBalance(address)

	err := arwen.StoreBytes(instCtx.Memory(), resultOffset, balance)
	if err != nil {
	}
}

//export getBlockHash
func getBlockHash(context unsafe.Pointer, number int64, resultOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	hash := hostContext.BlockHash(number)
	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, hash)

	if len(hash) == 0 {
		return 0
	}
	return 1
}

//export call
func call(context unsafe.Pointer, gasLimit int64, addressOffset int32, valueOffset int32, dataOffset int32, dataLength int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	send := hostContext.GetVMInput().CallerAddr
	dest := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	value := arwen.LoadBytes(instCtx.Memory(), valueOffset, arwen.BalanceLen)
	data := arwen.LoadBytes(instCtx.Memory(), dataOffset, dataLength)

	_, err := hostContext.Transfer(dest, send, big.NewInt(0).SetBytes(value), data, gasLimit)
	if err != nil {
		return 1
	}

	return 0
}

//export callDataCopy
func callDataCopy(context unsafe.Pointer, resultOffset int32, dataOffset int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	callData := ethContext.CallData()
	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, callData[dataOffset:dataOffset+length])
}

//export getCallDataSize
func getCallDataSize(context unsafe.Pointer) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	return int32(len(ethContext.CallData()))
}

//export storageStore
func storageStore(context unsafe.Pointer, pathOffset int32, valueOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	key := arwen.LoadBytes(instCtx.Memory(), pathOffset, arwen.HashLen)
	data := arwen.LoadBytes(instCtx.Memory(), valueOffset, arwen.HashLen)

	_ = hostContext.SetStorage(hostContext.GetSCAddress(), key, data)
}

//export storageLoad
func storageLoad(context unsafe.Pointer, pathOffset int32, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	key := arwen.LoadBytes(instCtx.Memory(), pathOffset, arwen.HashLen)
	data := hostContext.GetStorage(hostContext.GetSCAddress(), key)

	currInput := make([]byte, arwen.HashLen)
	copy(currInput[arwen.HashLen-len(data):], data)

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, currInput)
}

//export getCaller
func getCaller(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	caller := hostContext.GetVMInput().CallerAddr

	err := arwen.StoreBytes(instCtx.Memory(), resultOffset, caller)
	if err != nil {
	}
}

//export getCallValue
func getCallValue(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	value := hostContext.GetVMInput().CallValue.Bytes()
	length := len(value)
	invBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		invBytes[length-i-1] = value[i]
	}

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, invBytes)
}

//export codeCopy
func codeCopy(context unsafe.Pointer, resultOffset int32, codeOffset int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	scAddress := hostContext.GetSCAddress()
	code := hostContext.GetCode(scAddress)

	if int32(len(code)) > codeOffset+length {
		_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, code[codeOffset:codeOffset+length])
	}
}

//export getCodeSize
func getCodeSize(context unsafe.Pointer) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	return int32(hostContext.GetCodeSize(hostContext.GetSCAddress()))
}

//export externalCodeCopy
func externalCodeCopy(context unsafe.Pointer, addressOffset int32, resultOffset int32, codeOffset int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	dest := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	code := hostContext.GetCode(dest)

	if int32(len(code)) > codeOffset+length {
		_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, code[codeOffset:codeOffset+length])
	}
}

//export getExternalCodeSize
func getExternalCodeSize(context unsafe.Pointer, addressOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	dest := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)

	return int32(hostContext.GetCodeSize(dest))
}

//export getGasLeft
func getGasLeft(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	return ethContext.GasLeft()
}

//export getBlockGasLimit
func getBlockGasLimit(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	return ethContext.BlockGasLimit()
}

//export getTxGasPrice
func getTxGasPrice(context unsafe.Pointer, valueOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	gasPrice := hostContext.GetVMInput().GasPrice

	gasU128 := make([]byte, 16)
	copy(gasU128[16-len(gasPrice.Bytes()):], gasPrice.Bytes())

	_ = arwen.StoreBytes(instCtx.Memory(), valueOffset, gasU128)
}

//export logTopics
func logTopics(context unsafe.Pointer, dataOffset int32, length int32, numberOfTopics int32, topic1 int32, topic2 int32, topic3 int32, topic4 int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	data := arwen.LoadBytes(instCtx.Memory(), dataOffset, length)

	topics := make([]int32, 0)
	topics = append(topics, topic1)
	topics = append(topics, topic2)
	topics = append(topics, topic3)
	topics = append(topics, topic4)

	topicsData := make([][]byte, numberOfTopics)
	for i := int32(0); i < numberOfTopics; i++ {
		topicsData[i] = arwen.LoadBytes(instCtx.Memory(), topics[i], arwen.HashLen)
	}

	hostContext.WriteLog(hostContext.GetSCAddress(), topicsData, data)
}

//export getTxOrigin
func getTxOrigin(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, hostContext.GetVMInput().CallerAddr)
}

//export finish
func finish(context unsafe.Pointer, resultOffset int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	data := arwen.LoadBytes(instCtx.Memory(), resultOffset, length)
	hostContext.Finish(data)
}

//export revert
func revert(context unsafe.Pointer, dataOffset int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	data := arwen.LoadBytes(instCtx.Memory(), dataOffset, length)
	hostContext.Finish(data)
	hostContext.SignalUserError()
}

//export selfDestruct
func selfDestruct(context unsafe.Pointer, addressOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := arwen.GetEthContext(instCtx.Data())

	address := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.HashLen)

	hostContext.SelfDestruct(address, hostContext.GetVMInput().CallerAddr)
}

//export getBlockNumber
func getBlockNumber(context unsafe.Pointer) int64 {
	//instCtx := wasmer.IntoInstanceContext(context)
	//hostContext := arwen.arwen.GetEthContext(instCtx.Data())

	//TODO: implement
	return 0 //hostContext.BlockNonce()
}

//export getBlockTimestamp
func getBlockTimestamp(context unsafe.Pointer) int64 {
	//instCtx := wasmer.IntoInstanceContext(context)
	//hostContext := arwen.arwen.GetEthContext(instCtx.Data())

	//return hostContext.BlockTimeStamp()
	return 0
}

//export getReturnDataSize
func getReturnDataSize(context unsafe.Pointer) int32 {
	//TODO: implement
	return 0
}

//export returnDataCopy
func returnDataCopy(context unsafe.Pointer, resultOffset int32, dataOffset int32, length int32) {
	//TODO: implement
}

//export getBlockCoinbase
func getBlockCoinbase(context unsafe.Pointer, resultOffset int32) {
	//TODO: implement
}

//export getBlockDifficulty
func getBlockDifficulty(context unsafe.Pointer, resultOffset int32) {
	//TODO: implement
}

//export callCode
func callCode(context unsafe.Pointer, gasLimit int64, addressOffset int32, valueOffset int32, dataOffset int32, dataLength int32) int32 {
	//TODO: implement
	return 0
}

//export callDelegate
func callDelegate(context unsafe.Pointer, gasLimit int64, addressOffset int32, dataOffset int32, dataLength int32) int32 {
	//TODO: implement
	return 0
}

//export callStatic
func callStatic(context unsafe.Pointer, gasLimit int64, addressOffset int32, dataOffset int32, dataLength int32) int32 {
	//TODO: implement
	return 0
}

//export create
func create(context unsafe.Pointer, valueOffset int32, dataOffset int32, length int32, resultOffset int32) int32 {
	//TODO: implement
	return 0
}

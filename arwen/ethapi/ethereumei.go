package ethapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char u8;
// typedef int i32;
// typedef int i32ptr;
// extern void ethuseGas(void *context, long long  gas);
// extern void ethgetAddress(void *context, i32ptr resultOffset);
// extern void ethgetExternalBalance(void *context, i32ptr addressOffset, i32ptr resultOffset);
// extern i32 ethgetBlockHash(void *context, long long number, i32ptr resultOffset);
// extern i32 ethcall(void *context, long long gas, i32ptr addressOffset, i32ptr valueOffset, i32ptr dataOffset, i32 dataLength);
// extern i32 ethgetCallDataSize(void *context);
// extern void ethcallDataCopy(void *context, i32ptr resultsOffset, i32ptr dataOffset, i32 length);
// extern i32 ethcallCode(void *context, long long gas, i32ptr addressOffset, i32ptr valueOffset, i32ptr dataOffset, i32 dataLength);
// extern i32 ethcallDelegate(void *context, long long gas, i32ptr addressOffset, i32ptr dataOffset, i32 dataLength);
// extern i32 ethcallStatic(void *context, long long gas, i32ptr addressOffset, i32ptr dataOffset, i32 dataLength);
// extern void ethstorageStore(void *context, i32ptr pathOffset, i32ptr valueOffset);
// extern void ethstorageLoad(void *context, i32ptr pathOffset, i32ptr resultOffset);
// extern void ethgetCaller(void *context, i32ptr resultOffset);
// extern void ethgetCallValue(void *context, i32ptr resultOffset);
// extern void ethcodeCopy(void *context, i32ptr resultOffset, i32 codeOffset, i32 length);
// extern i32 ethgetCodeSize(void *context);
// extern void ethgetBlockCoinbase(void *context, i32ptr resultOffset);
// extern i32 ethcreate(void *context, i32ptr valueoffset, i32ptr dataOffset, i32 length, i32ptr resultsOffset);
// extern void ethgetBlockDifficulty(void *context, i32ptr resultOffset);
// extern void ethexternalCodeCopy(void *context, i32ptr addressOffset, i32ptr resultOffset, i32 codeOffset, i32 length);
// extern i32 ethgetExternalCodeSize(void *context, i32ptr addressOffset);
// extern long long ethgetGasLeft(void *context);
// extern long long ethgetBlockGasLimit(void *context);
// extern void ethgetTxGasPrice(void *context, i32ptr valueOffset);
// extern void ethlogTopics(void *context, i32ptr dataOffset, i32 length, i32 numberOftopics, i32ptr topic1, i32ptr topic2, i32ptr topic3, i32ptr topic4);
// extern long long ethgetBlockNumber(void *context);
// extern void ethgetTxOrigin(void *context, i32ptr resultOffset);
// extern void ethfinish(void *context, i32ptr dataOffset, i32 length);
// extern void ethrevert(void *context, i32ptr dataOffset, i32 length);
// extern i32 ethgetReturnDataSize(void *context);
// extern void ethreturnDataCopy(void *context, i32ptr resultOffset, i32 dataOffset, i32 length);
// extern void ethselfDestruct(void *context, i32ptr addressOffset);
// extern long long ethgetBlockTimestamp(void *context);
import "C"
import (
	"math/big"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// EthereumImports creates a new wasmer.Imports populated with the Ethereum API methods
func EthereumImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("ethereum")

	imports, err := imports.Append("useGas", ethuseGas, C.ethuseGas)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getAddress", ethgetAddress, C.ethgetAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getExternalBalance", ethgetExternalBalance, C.ethgetExternalBalance)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockHash", ethgetBlockHash, C.ethgetBlockHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("call", ethcall, C.ethcall)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("callDataCopy", ethcallDataCopy, C.ethcallDataCopy)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCallDataSize", ethgetCallDataSize, C.ethgetCallDataSize)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("callCode", ethcallCode, C.ethcallCode)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("callDelegate", ethcallDelegate, C.ethcallDelegate)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("callStatic", ethcallStatic, C.ethcallStatic)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageStore", ethstorageStore, C.ethstorageStore)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageLoad", ethstorageLoad, C.ethstorageLoad)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCaller", ethgetCaller, C.ethgetCaller)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCallValue", ethgetCallValue, C.ethgetCallValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("codeCopy", ethcodeCopy, C.ethcodeCopy)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCodeSize", ethgetCodeSize, C.ethgetCodeSize)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockCoinbase", ethgetBlockCoinbase, C.ethgetBlockCoinbase)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("create", ethcreate, C.ethcreate)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockDifficulty", ethgetBlockDifficulty, C.ethgetBlockDifficulty)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("externalCodeCopy", ethexternalCodeCopy, C.ethexternalCodeCopy)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getExternalCodeSize", ethgetExternalCodeSize, C.ethgetExternalCodeSize)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getGasLeft", ethgetGasLeft, C.ethgetGasLeft)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockGasLimit", ethgetBlockGasLimit, C.ethgetBlockGasLimit)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getTxGasPrice", ethgetTxGasPrice, C.ethgetTxGasPrice)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("log", ethlogTopics, C.ethlogTopics)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockNumber", ethgetBlockNumber, C.ethgetBlockNumber)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getTxOrigin", ethgetTxOrigin, C.ethgetTxOrigin)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("finish", ethfinish, C.ethfinish)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("revert", ethrevert, C.ethrevert)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getReturnDataSize", ethgetReturnDataSize, C.ethgetReturnDataSize)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("returnDataCopy", ethreturnDataCopy, C.ethreturnDataCopy)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("selfDestruct", ethselfDestruct, C.ethselfDestruct)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockTimestamp", ethgetBlockTimestamp, C.ethgetBlockTimestamp)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export ethuseGas
func ethuseGas(context unsafe.Pointer, useGas int64) {
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.UseGas + uint64(useGas)
	metering.UseGas(gasToUse)
}

//export ethgetAddress
func ethgetAddress(context unsafe.Pointer, resultOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.GetAddress
	metering.UseGas(gasToUse)

	err := runtime.MemStore(resultOffset, runtime.GetSCAddress())
	if arwen.WithFault(err, context, true) {
		return
	}
}

//export ethgetExternalBalance
func ethgetExternalBalance(context unsafe.Pointer, addressOffset int32, resultOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.GetExternalBalance
	metering.UseGas(gasToUse)

	address, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, true) {
		return
	}

	balance := blockchain.GetBalance(address)

	err = runtime.MemStore(resultOffset, balance)
	if arwen.WithFault(err, context, true) {
		return
	}
}

//export ethgetBlockHash
func ethgetBlockHash(context unsafe.Pointer, number int64, resultOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.GetBlockHash
	metering.UseGas(gasToUse)

	hash := blockchain.BlockHash(number)
	err := runtime.MemStore(resultOffset, hash)
	if err != nil {
		return 0
	}

	if len(hash) == 0 {
		return 0
	}

	return 1
}

//export ethcallDataCopy
func ethcallDataCopy(context unsafe.Pointer, resultOffset int32, dataOffset int32, length int32) {
	host := arwen.GetVmContext(context)
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().EthAPICost.CallDataCopy
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(length)
	metering.UseGas(gasToUse)

	callData := host.EthereumCallData()
	callDataSlice, err := arwen.GuardedGetBytesSlice(callData, dataOffset, length)
	if arwen.WithFault(err, context, true) {
		return
	}

	err = runtime.MemStore(resultOffset, callDataSlice)
	if arwen.WithFault(err, context, true) {
		return
	}
}

//export ethgetCallDataSize
func ethgetCallDataSize(context unsafe.Pointer) int32 {
	host := arwen.GetVmContext(context)
	metering := host.Metering()

	gasToUse := metering.GasSchedule().EthAPICost.GetCallDataSize
	metering.UseGas(gasToUse)

	callData := host.EthereumCallData()
	return int32(len(callData))
}

//export ethstorageStore
func ethstorageStore(context unsafe.Pointer, pathOffset int32, valueOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	storage := arwen.GetStorageContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.StorageStore
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(pathOffset, arwen.HashLen)
	if arwen.WithFault(err, context, true) {
		return
	}

	data, err := runtime.MemLoad(valueOffset, arwen.HashLen)
	if arwen.WithFault(err, context, true) {
		return
	}

	_, err = storage.SetStorage(key, data)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}
}

//export ethstorageLoad
func ethstorageLoad(context unsafe.Pointer, pathOffset int32, resultOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)
	storage := arwen.GetStorageContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.StorageLoad
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(pathOffset, arwen.HashLen)
	if arwen.WithFault(err, context, true) {
		return
	}

	data := storage.GetStorage(key)
	dataGasToUse := metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(data))
	metering.UseGas(dataGasToUse)

	currInput := make([]byte, arwen.HashLen)
	copy(currInput[arwen.HashLen-len(data):], data)

	err = runtime.MemStore(resultOffset, currInput)
	if arwen.WithFault(err, context, true) {
		return
	}
}

//export ethgetCaller
func ethgetCaller(context unsafe.Pointer, resultOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.GetCaller
	metering.UseGas(gasToUse)

	caller := convertToEthAddress(runtime.GetVMInput().CallerAddr)
	err := runtime.MemStore(resultOffset, caller)
	if arwen.WithFault(err, context, true) {
		return
	}
}

//export ethgetCallValue
func ethgetCallValue(context unsafe.Pointer, resultOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.GetCallValue
	metering.UseGas(gasToUse)

	value := convertToEthU128(runtime.GetVMInput().CallValue.Bytes())

	length := len(value)
	invBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		invBytes[length-i-1] = value[i]
	}

	err := runtime.MemStore(resultOffset, invBytes)
	if arwen.WithFault(err, context, true) {
		return
	}
}

//export ethcodeCopy
func ethcodeCopy(context unsafe.Pointer, resultOffset int32, codeOffset int32, length int32) {
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.CodeCopy
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(length)
	metering.UseGas(gasToUse)

	scAddress := runtime.GetSCAddress()
	code, err := blockchain.GetCode(scAddress)
	if arwen.WithFault(err, context, true) {
		return
	}

	codeSlice, err := arwen.GuardedGetBytesSlice(code, codeOffset, length)
	if arwen.WithFault(err, context, true) {
		return
	}

	err = runtime.MemStore(resultOffset, codeSlice)
	if arwen.WithFault(err, context, true) {
		return
	}
}

//export ethgetCodeSize
func ethgetCodeSize(context unsafe.Pointer) int32 {
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.GetCodeSize
	metering.UseGas(gasToUse)

	codeSize, err := blockchain.GetCodeSize(runtime.GetSCAddress())
	if err != nil {
		return 0
	}

	return codeSize
}

//export ethexternalCodeCopy
func ethexternalCodeCopy(context unsafe.Pointer, addressOffset int32, resultOffset int32, codeOffset int32, length int32) {
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.ExternalCodeCopy
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(length)
	metering.UseGas(gasToUse)

	dest, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if arwen.WithFault(err, context, true) {
		return
	}

	code, err := blockchain.GetCode(dest)
	if arwen.WithFault(err, context, true) {
		return
	}

	codeSlice, err := arwen.GuardedGetBytesSlice(code, codeOffset, length)
	if arwen.WithFault(err, context, true) {
		return
	}

	err = runtime.MemStore(resultOffset, codeSlice)
	if arwen.WithFault(err, context, true) {
		return
	}
}

//export ethgetExternalCodeSize
func ethgetExternalCodeSize(context unsafe.Pointer, addressOffset int32) int32 {
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.GetExternalCodeSize
	metering.UseGas(gasToUse)

	dest, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if err != nil {
		return 0
	}

	codeSize, err := blockchain.GetCodeSize(dest)
	if err != nil {
		return 0
	}

	return codeSize
}

//export ethgetGasLeft
func ethgetGasLeft(context unsafe.Pointer) int64 {
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.GetGasLeft
	metering.UseGas(gasToUse)

	return int64(metering.GasLeft())
}

//export ethgetBlockGasLimit
func ethgetBlockGasLimit(context unsafe.Pointer) int64 {
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.GetBlockGasLimit
	metering.UseGas(gasToUse)

	return int64(metering.BlockGasLimit())
}

//export ethgetTxGasPrice
func ethgetTxGasPrice(context unsafe.Pointer, valueOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.GetTxGasPrice
	metering.UseGas(gasToUse)

	gasPrice := runtime.GetVMInput().GasPrice
	gasBigInt := big.NewInt(0).SetUint64(gasPrice)

	gasU128 := make([]byte, 16)
	copy(gasU128[16-len(gasBigInt.Bytes()):], gasBigInt.Bytes())

	err := runtime.MemStore(valueOffset, gasU128)
	if arwen.WithFault(err, context, true) {
		return
	}
}

//export ethlogTopics
func ethlogTopics(context unsafe.Pointer, dataOffset int32, length int32, numberOfTopics int32, topic1 int32, topic2 int32, topic3 int32, topic4 int32) {
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.Log
	gasToUse += metering.GasSchedule().BaseOperationCost.PersistPerByte * (4*arwen.HashLen + uint64(length))
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, true) {
		return
	}

	topics := make([]int32, 0)
	topics = append(topics, topic1)
	topics = append(topics, topic2)
	topics = append(topics, topic3)
	topics = append(topics, topic4)

	topicsData, err := arwen.GuardedMakeByteSlice2D(numberOfTopics)
	if arwen.WithFault(err, context, true) {
		return
	}

	for i := int32(0); i < numberOfTopics; i++ {
		topicsData[i], err = runtime.MemLoad(topics[i], arwen.HashLen)
		if arwen.WithFault(err, context, true) {
			return
		}
	}

	output.WriteLog(runtime.GetSCAddress(), topicsData, data)
}

//export ethgetTxOrigin
func ethgetTxOrigin(context unsafe.Pointer, resultOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.GetTxOrigin
	metering.UseGas(gasToUse)

	caller := convertToEthAddress(runtime.GetVMInput().CallerAddr)
	err := runtime.MemStore(resultOffset, caller)
	if arwen.WithFault(err, context, true) {
		return
	}
}

//export ethfinish
func ethfinish(context unsafe.Pointer, resultOffset int32, length int32) {
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.Finish
	gasToUse += metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(length)
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(resultOffset, length)
	if arwen.WithFault(err, context, true) {
		return
	}

	output.ClearReturnData()
	output.Finish(data)
}

//export ethrevert
func ethrevert(context unsafe.Pointer, dataOffset int32, length int32) {
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.Revert
	gasToUse += metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(length)
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if arwen.WithFault(err, context, true) {
		return
	}

	output.ClearReturnData()
	output.Finish(data)

	runtime.SignalUserError("revert")
}

//export ethselfDestruct
func ethselfDestruct(context unsafe.Pointer, addressOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.SelfDestruct
	metering.UseGas(gasToUse)

	address, err := runtime.MemLoad(addressOffset, arwen.HashLen)
	if arwen.WithFault(err, context, true) {
		return
	}

	caller := runtime.GetVMInput().CallerAddr
	output.SelfDestruct(address, caller)
}

//export ethgetBlockNumber
func ethgetBlockNumber(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.GetBlockNumber
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentNonce())
}

//export ethgetBlockTimestamp
func ethgetBlockTimestamp(context unsafe.Pointer) int64 {
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.GetBlockTimeStamp
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentTimeStamp())
}

//export ethgetReturnDataSize
func ethgetReturnDataSize(context unsafe.Pointer) int32 {
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.GetReturnDataSize
	metering.UseGas(gasToUse)

	returnData := output.ReturnData()
	size := int32(0)
	for _, data := range returnData {
		size += int32(len(data))
	}

	return size
}

//export ethreturnDataCopy
func ethreturnDataCopy(context unsafe.Pointer, resultOffset int32, dataOffset int32, length int32) {
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.ReturnDataCopy
	metering.UseGas(gasToUse)

	returnData := output.ReturnData()
	ethReturnData := make([]byte, 0)
	for _, data := range returnData {
		ethReturnData = append(ethReturnData, data...)
	}

	if int32(len(ethReturnData)) < dataOffset+length {
		arwen.WithFault(arwen.ErrInvalidAPICall, context, true)
		return
	}

	returnDataSlice, err := arwen.GuardedGetBytesSlice(ethReturnData, dataOffset, length)
	if arwen.WithFault(err, context, true) {
		return
	}

	err = runtime.MemStore(resultOffset, returnDataSlice)
	arwen.WithFault(err, context, true)
}

//export ethgetBlockCoinbase
func ethgetBlockCoinbase(context unsafe.Pointer, resultOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.GetBlockCoinbase
	metering.UseGas(gasToUse)

	randomSeed := blockchain.CurrentRandomSeed()
	err := runtime.MemStore(resultOffset, randomSeed)
	arwen.WithFault(err, context, true)
}

//export ethgetBlockDifficulty
func ethgetBlockDifficulty(context unsafe.Pointer, resultOffset int32) {
	runtime := arwen.GetRuntimeContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().EthAPICost.GetBlockCoinbase
	metering.UseGas(gasToUse)

	randomSeed := blockchain.CurrentRandomSeed()
	err := runtime.MemStore(resultOffset, randomSeed)
	arwen.WithFault(err, context, true)
}

//export ethcall
func ethcall(context unsafe.Pointer, gasLimit int64, addressOffset int32, valueOffset int32, dataOffset int32, dataLength int32) int32 {
	host := arwen.GetVmContext(context)
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().EthAPICost.Call
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if err != nil {
		return 1
	}

	dataGasToUse := metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(data))
	metering.UseGas(dataGasToUse)

	send := runtime.GetSCAddress()
	dest, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if err != nil {
		return 1
	}

	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if err != nil {
		return 1
	}

	invBytes := arwen.InverseBytes(value)
	bigIntVal := big.NewInt(0).SetBytes(invBytes)
	err = output.Transfer(dest, send, 0, bigIntVal, nil)
	if err != nil {
		return 1
	}

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  send,
			Arguments:   [][]byte{data},
			CallValue:   bigIntVal,
			GasPrice:    0,
			GasProvided: metering.BoundGasLimit(gasLimit),
		},
		RecipientAddr: dest,
		Function:      "main",
	}

	_, err = host.ExecuteOnDestContext(contractCallInput)
	if err != nil {
		return 1
	}

	return 0
}

//export ethcallCode
func ethcallCode(context unsafe.Pointer, gasLimit int64, addressOffset int32, valueOffset int32, dataOffset int32, dataLength int32) int32 {
	host := arwen.GetVmContext(context)
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().EthAPICost.CallCode
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if err != nil {
		return 1
	}

	dataGasToUse := metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(data))
	metering.UseGas(dataGasToUse)

	send := runtime.GetSCAddress()
	dest, err := runtime.MemLoad(addressOffset, arwen.AddressLen)
	if err != nil {
		return 1
	}

	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if err != nil {
		return 1
	}

	invBytes := arwen.InverseBytes(value)
	err = output.Transfer(dest, send, 0, big.NewInt(0).SetBytes(invBytes), nil)
	if err != nil {
		return 1
	}

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  send,
			Arguments:   [][]byte{data},
			CallValue:   big.NewInt(0).SetBytes(value),
			GasPrice:    0,
			GasProvided: metering.BoundGasLimit(gasLimit),
		},
		RecipientAddr: dest,
		Function:      "main",
	}

	err = host.ExecuteOnSameContext(contractCallInput)
	if err != nil {
		return 1
	}

	return 0
}

//export ethcallDelegate
func ethcallDelegate(context unsafe.Pointer, gasLimit int64, addressOffset int32, dataOffset int32, dataLength int32) int32 {
	host := arwen.GetVmContext(context)
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().EthAPICost.CallDelegate
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if err != nil {
		return 1
	}

	dataGasToUse := metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(data))
	metering.UseGas(dataGasToUse)

	value := runtime.GetVMInput().CallValue
	sender := runtime.GetVMInput().CallerAddr

	address, err := runtime.MemLoad(addressOffset, arwen.HashLen)
	if err != nil {
		return 1
	}

	err = output.Transfer(address, sender, 0, value, nil)
	if err != nil {
		return 1
	}

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   [][]byte{data},
			CallValue:   value,
			GasPrice:    0,
			GasProvided: metering.BoundGasLimit(gasLimit),
		},
		RecipientAddr: address,
		Function:      "main",
	}

	err = host.ExecuteOnSameContext(contractCallInput)
	if err != nil {
		return 1
	}

	return 0
}

//export ethcallStatic
func ethcallStatic(context unsafe.Pointer, gasLimit int64, addressOffset int32, dataOffset int32, dataLength int32) int32 {
	host := arwen.GetVmContext(context)
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().EthAPICost.CallStatic
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if err != nil {
		return 1
	}

	dataGasToUse := metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(data))
	metering.UseGas(dataGasToUse)

	address, err := runtime.MemLoad(addressOffset, arwen.AddressLenEth)
	if err != nil {
		return 1
	}

	value := runtime.GetVMInput().CallValue
	sender := runtime.GetVMInput().CallerAddr

	if IsAddressForPredefinedContract(address) {
		err := CallPredefinedContract(context, address, data)
		if err != nil {
			return 1
		}

		return 0
	}

	err = output.Transfer(address, sender, 0, value, nil)
	if err != nil {
		return 1
	}

	runtime.SetReadOnly(true)

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   [][]byte{data},
			CallValue:   value,
			GasPrice:    0,
			GasProvided: metering.BoundGasLimit(gasLimit),
		},
		RecipientAddr: address,
		Function:      "main",
	}

	err = host.ExecuteOnSameContext(contractCallInput)

	runtime.SetReadOnly(false)
	if err != nil {
		return 1
	}

	return 0
}

//export ethcreate
func ethcreate(context unsafe.Pointer, valueOffset int32, dataOffset int32, length int32, resultOffset int32) int32 {
	host := arwen.GetVmContext(context)
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().EthAPICost.Create
	metering.UseGas(gasToUse)

	sender := runtime.GetSCAddress()
	value, err := runtime.MemLoad(valueOffset, arwen.BalanceLen)
	if err != nil {
		return 1
	}

	data, err := runtime.MemLoad(dataOffset, length)
	if err != nil {
		return 1
	}

	contractCreate := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   nil,
			CallValue:   big.NewInt(0).SetBytes(value),
			GasPrice:    0,
			GasProvided: metering.GasLeft(),
		},
		ContractCode: data,
	}

	newAddress, err := host.CreateNewContract(contractCreate)
	if err != nil {
		return 1
	}

	err = runtime.MemStore(resultOffset, newAddress)
	if err != nil {
		return 1
	}

	return 0
}

// https://ewasm.readthedocs.io/en/mkdocs/eth_interface/#data-types
func convertToEthAddress(address []byte) []byte {
	ethAddress := address[arwen.AddressLen-arwen.AddressLenEth : arwen.AddressLen]
	return ethAddress
}

// convertToEthU128 adds zero-left-padding up to a total of 16 bytes
// If the input data is larger than 16 bytes, an array of 16 zeros is returned
func convertToEthU128(data []byte) []byte {
	const noBytes = 16

	result := make([]byte, noBytes)
	length := len(data)

	if length > noBytes {
		return result
	}

	copy(result[noBytes-length:], data)
	return result
}

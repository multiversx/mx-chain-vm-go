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
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/go-ext-wasm/wasmer"
)

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
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	gasToUse := ethContext.GasSchedule().EthAPICost.UseGas + uint64(useGas)
	ethContext.UseGas(gasToUse)
}

//export ethgetAddress
func ethgetAddress(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	err := arwen.StoreBytes(instCtx.Memory(), resultOffset, ethContext.GetSCAddress())
	if withFault(err, context) {
		return
	}

	gasToUse := ethContext.GasSchedule().EthAPICost.GetAddress
	ethContext.UseGas(gasToUse)
}

//export ethgetExternalBalance
func ethgetExternalBalance(context unsafe.Pointer, addressOffset int32, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	address, err := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	if withFault(err, context) {
		return
	}

	balance := ethContext.GetBalance(address)

	err = arwen.StoreBytes(instCtx.Memory(), resultOffset, balance)
	if withFault(err, context) {
		return
	}

	gasToUse := ethContext.GasSchedule().EthAPICost.GetExternalBalance
	ethContext.UseGas(gasToUse)
}

//export ethgetBlockHash
func ethgetBlockHash(context unsafe.Pointer, number int64, resultOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	hash := ethContext.BlockHash(number)
	err := arwen.StoreBytes(instCtx.Memory(), resultOffset, hash)
	if withFault(err, context) {
		return 0
	}

	gasToUse := ethContext.GasSchedule().EthAPICost.GetBlockHash
	ethContext.UseGas(gasToUse)

	if len(hash) == 0 {
		return 0
	}

	return 1
}

//export ethcallDataCopy
func ethcallDataCopy(context unsafe.Pointer, resultOffset int32, dataOffset int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	callData := ethContext.CallData()
	callDataSlice, err := arwen.GuardedGetBytesSlice(callData, dataOffset, length)
	if withFault(err, context) {
		return
	}

	err = arwen.StoreBytes(instCtx.Memory(), resultOffset, callDataSlice)
	if withFault(err, context) {
		return
	}

	gasToUse := ethContext.GasSchedule().EthAPICost.CallDataCopy
	gasToUse += ethContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(length)
	ethContext.UseGas(gasToUse)
}

//export ethgetCallDataSize
func ethgetCallDataSize(context unsafe.Pointer) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	gasToUse := ethContext.GasSchedule().EthAPICost.GetCallDataSize
	ethContext.UseGas(gasToUse)

	return int32(len(ethContext.CallData()))
}

//export ethstorageStore
func ethstorageStore(context unsafe.Pointer, pathOffset int32, valueOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	key, err := arwen.LoadBytes(instCtx.Memory(), pathOffset, arwen.HashLen)
	if withFault(err, context) {
		return
	}

	data, err := arwen.LoadBytes(instCtx.Memory(), valueOffset, arwen.HashLen)
	if withFault(err, context) {
		return
	}

	_ = ethContext.SetStorage(ethContext.GetSCAddress(), key, data)

	gasToUse := ethContext.GasSchedule().EthAPICost.StorageStore
	ethContext.UseGas(gasToUse)
}

//export ethstorageLoad
func ethstorageLoad(context unsafe.Pointer, pathOffset int32, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	key, err := arwen.LoadBytes(instCtx.Memory(), pathOffset, arwen.HashLen)
	if withFault(err, context) {
		return
	}

	data := ethContext.GetStorage(ethContext.GetSCAddress(), key)

	currInput := make([]byte, arwen.HashLen)
	copy(currInput[arwen.HashLen-len(data):], data)

	err = arwen.StoreBytes(instCtx.Memory(), resultOffset, currInput)
	if withFault(err, context) {
		return
	}

	gasToUse := ethContext.GasSchedule().EthAPICost.StorageLoad
	gasToUse += ethContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(data))
	ethContext.UseGas(gasToUse)
}

//export ethgetCaller
func ethgetCaller(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	caller := convertToEthAddress(ethContext.GetVMInput().CallerAddr)
	err := arwen.StoreBytes(instCtx.Memory(), resultOffset, caller)
	if withFault(err, context) {
		return
	}

	gasToUse := ethContext.GasSchedule().EthAPICost.GetCaller
	ethContext.UseGas(gasToUse)
}

//export ethgetCallValue
func ethgetCallValue(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	value := convertToEthU128(ethContext.GetVMInput().CallValue.Bytes())

	length := len(value)
	invBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		invBytes[length-i-1] = value[i]
	}

	err := arwen.StoreBytes(instCtx.Memory(), resultOffset, invBytes)
	if withFault(err, context) {
		return
	}

	gasToUse := ethContext.GasSchedule().EthAPICost.GetCallValue
	ethContext.UseGas(gasToUse)
}

//export ethcodeCopy
func ethcodeCopy(context unsafe.Pointer, resultOffset int32, codeOffset int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	scAddress := ethContext.GetSCAddress()
	code, err := ethContext.GetCode(scAddress)
	if withFault(err, context) {
		return
	}

	codeSlice, err := arwen.GuardedGetBytesSlice(code, codeOffset, length)
	if withFault(err, context) {
		return
	}

	err = arwen.StoreBytes(instCtx.Memory(), resultOffset, codeSlice)
	if withFault(err, context) {
		return
	}

	gasToUse := ethContext.GasSchedule().EthAPICost.CodeCopy
	gasToUse += ethContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(length)
	ethContext.UseGas(gasToUse)
}

//export ethgetCodeSize
func ethgetCodeSize(context unsafe.Pointer) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	gasToUse := ethContext.GasSchedule().EthAPICost.GetCodeSize
	ethContext.UseGas(gasToUse)

	return ethContext.GetCodeSize(ethContext.GetSCAddress())
}

//export ethexternalCodeCopy
func ethexternalCodeCopy(context unsafe.Pointer, addressOffset int32, resultOffset int32, codeOffset int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	dest, err := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	if withFault(err, context) {
		return
	}

	code, err := ethContext.GetCode(dest)
	if withFault(err, context) {
		return
	}

	codeSlice, err := arwen.GuardedGetBytesSlice(code, codeOffset, length)
	if withFault(err, context) {
		return
	}

	err = arwen.StoreBytes(instCtx.Memory(), resultOffset, codeSlice)
	if withFault(err, context) {
		return
	}

	gasToUse := ethContext.GasSchedule().EthAPICost.ExternalCodeCopy
	gasToUse += ethContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(length)
	ethContext.UseGas(gasToUse)
}

//export ethgetExternalCodeSize
func ethgetExternalCodeSize(context unsafe.Pointer, addressOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	dest, err := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	if withFault(err, context) {
		return 0
	}

	gasToUse := ethContext.GasSchedule().EthAPICost.GetExternalCodeSize
	ethContext.UseGas(gasToUse)

	return ethContext.GetCodeSize(dest)
}

//export ethgetGasLeft
func ethgetGasLeft(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	gasToUse := ethContext.GasSchedule().EthAPICost.GetGasLeft
	ethContext.UseGas(gasToUse)

	return int64(ethContext.GasLeft())
}

//export ethgetBlockGasLimit
func ethgetBlockGasLimit(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	gasToUse := ethContext.GasSchedule().EthAPICost.GetBlockGasLimit
	ethContext.UseGas(gasToUse)

	return int64(ethContext.BlockGasLimit())
}

//export ethgetTxGasPrice
func ethgetTxGasPrice(context unsafe.Pointer, valueOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	gasPrice := ethContext.GetVMInput().GasPrice
	gasBigInt := big.NewInt(0).SetUint64(gasPrice)

	gasU128 := make([]byte, 16)
	copy(gasU128[16-len(gasBigInt.Bytes()):], gasBigInt.Bytes())

	err := arwen.StoreBytes(instCtx.Memory(), valueOffset, gasU128)
	if withFault(err, context) {
		return
	}

	gasToUse := ethContext.GasSchedule().EthAPICost.GetTxGasPrice
	ethContext.UseGas(gasToUse)
}

//export ethlogTopics
func ethlogTopics(context unsafe.Pointer, dataOffset int32, length int32, numberOfTopics int32, topic1 int32, topic2 int32, topic3 int32, topic4 int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	data, err := arwen.LoadBytes(instCtx.Memory(), dataOffset, length)
	if withFault(err, context) {
		return
	}

	topics := make([]int32, 0)
	topics = append(topics, topic1)
	topics = append(topics, topic2)
	topics = append(topics, topic3)
	topics = append(topics, topic4)

	topicsData, err := arwen.GuardedMakeByteSlice2D(numberOfTopics)
	if withFault(err, context) {
		return
	}

	for i := int32(0); i < numberOfTopics; i++ {
		topicsData[i], err = arwen.LoadBytes(instCtx.Memory(), topics[i], arwen.HashLen)
		if withFault(err, context) {
			return
		}
	}

	ethContext.WriteLog(ethContext.GetSCAddress(), topicsData, data)

	gasToUse := ethContext.GasSchedule().EthAPICost.Log
	gasToUse += ethContext.GasSchedule().BaseOperationCost.PersistPerByte * (4*arwen.HashLen + uint64(length))
	ethContext.UseGas(gasToUse)
}

//export ethgetTxOrigin
func ethgetTxOrigin(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	caller := convertToEthAddress(ethContext.GetVMInput().CallerAddr)
	err := arwen.StoreBytes(instCtx.Memory(), resultOffset, caller)
	if withFault(err, context) {
		return
	}

	gasToUse := ethContext.GasSchedule().EthAPICost.GetTxOrigin
	ethContext.UseGas(gasToUse)
}

//export ethfinish
func ethfinish(context unsafe.Pointer, resultOffset int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	data, err := arwen.LoadBytes(instCtx.Memory(), resultOffset, length)
	if withFault(err, context) {
		return
	}

	ethContext.ClearReturnData()
	ethContext.Finish(data)

	gasToUse := ethContext.GasSchedule().EthAPICost.Finish
	gasToUse += ethContext.GasSchedule().BaseOperationCost.PersistPerByte * uint64(length)
	ethContext.UseGas(gasToUse)
}

//export ethrevert
func ethrevert(context unsafe.Pointer, dataOffset int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	data, err := arwen.LoadBytes(instCtx.Memory(), dataOffset, length)
	if withFault(err, context) {
		return
	}

	ethContext.ClearReturnData()
	ethContext.Finish(data)
	ethContext.SignalUserError()

	gasToUse := ethContext.GasSchedule().EthAPICost.Revert
	gasToUse += ethContext.GasSchedule().BaseOperationCost.PersistPerByte * uint64(length)
	ethContext.UseGas(gasToUse)
}

//export ethselfDestruct
func ethselfDestruct(context unsafe.Pointer, addressOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	address, err := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.HashLen)
	if withFault(err, context) {
		return
	}

	caller := ethContext.GetVMInput().CallerAddr
	ethContext.SelfDestruct(address, caller)

	gasToUse := ethContext.GasSchedule().EthAPICost.SelfDestruct
	ethContext.UseGas(gasToUse)
}

//export ethgetBlockNumber
func ethgetBlockNumber(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	gasToUse := ethContext.GasSchedule().EthAPICost.GetBlockNumber
	ethContext.UseGas(gasToUse)

	return int64(ethContext.BlockChainHook().CurrentNonce())
}

//export ethgetBlockTimestamp
func ethgetBlockTimestamp(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	gasToUse := ethContext.GasSchedule().EthAPICost.GetBlockTimeStamp
	ethContext.UseGas(gasToUse)

	return int64(ethContext.BlockChainHook().CurrentTimeStamp())
}

//export ethgetReturnDataSize
func ethgetReturnDataSize(context unsafe.Pointer) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	gasToUse := ethContext.GasSchedule().EthAPICost.GetReturnDataSize
	ethContext.UseGas(gasToUse)

	returnData := ethContext.ReturnData()
	size := int32(0)
	for _, data := range returnData {
		size += int32(len(data))
	}

	return size
}

//export ethreturnDataCopy
func ethreturnDataCopy(context unsafe.Pointer, resultOffset int32, dataOffset int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	gasToUse := ethContext.GasSchedule().EthAPICost.ReturnDataCopy
	ethContext.UseGas(gasToUse)

	returnData := ethContext.ReturnData()
	ethReturnData := make([]byte, 0)
	for _, data := range returnData {
		ethReturnData = append(ethReturnData, data...)
	}

	if int32(len(ethReturnData)) < dataOffset+length {
		return
	}

	returnDataSlice, err := arwen.GuardedGetBytesSlice(ethReturnData, dataOffset, length)
	if withFault(err, context) {
		return
	}

	err = arwen.StoreBytes(instCtx.Memory(), resultOffset, returnDataSlice)
	withFault(err, context)
}

//export ethgetBlockCoinbase
func ethgetBlockCoinbase(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	gasToUse := ethContext.GasSchedule().EthAPICost.GetBlockCoinbase
	ethContext.UseGas(gasToUse)

	randomSeed := ethContext.BlockChainHook().CurrentRandomSeed()
	err := arwen.StoreBytes(instCtx.Memory(), resultOffset, randomSeed)
	withFault(err, context)
}

//export ethgetBlockDifficulty
func ethgetBlockDifficulty(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	gasToUse := ethContext.GasSchedule().EthAPICost.GetBlockCoinbase
	ethContext.UseGas(gasToUse)

	randomSeed := ethContext.BlockChainHook().CurrentRandomSeed()
	err := arwen.StoreBytes(instCtx.Memory(), resultOffset, randomSeed)
	withFault(err, context)
}

//export ethcall
func ethcall(context unsafe.Pointer, gasLimit int64, addressOffset int32, valueOffset int32, dataOffset int32, dataLength int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	send := ethContext.GetSCAddress()
	dest, err := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	if withFault(err, context) {
		return 1
	}

	value, err := arwen.LoadBytes(instCtx.Memory(), valueOffset, arwen.BalanceLen)
	if withFault(err, context) {
		return 1
	}

	data, err := arwen.LoadBytes(instCtx.Memory(), dataOffset, dataLength)
	if withFault(err, context) {
		return 1
	}

	gasToUse := ethContext.GasSchedule().EthAPICost.Call
	gasToUse += ethContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(data))
	ethContext.UseGas(gasToUse)

	bigIntVal := big.NewInt(0).SetBytes(value)
	ethContext.Transfer(dest, send, bigIntVal, nil)

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  send,
			Arguments:   [][]byte{data},
			CallValue:   bigIntVal,
			GasPrice:    0,
			GasProvided: ethContext.BoundGasLimit(gasLimit),
		},
		RecipientAddr: dest,
		Function:      "main",
	}

	err = ethContext.ExecuteOnDestContext(contractCallInput)
	if err != nil {
		return 1
	}

	return 0
}

//export ethcallCode
func ethcallCode(context unsafe.Pointer, gasLimit int64, addressOffset int32, valueOffset int32, dataOffset int32, dataLength int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	send := ethContext.GetSCAddress()
	dest, err := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	if withFault(err, context) {
		return 1
	}

	value, err := arwen.LoadBytes(instCtx.Memory(), valueOffset, arwen.BalanceLen)
	if withFault(err, context) {
		return 1
	}

	data, err := arwen.LoadBytes(instCtx.Memory(), dataOffset, dataLength)
	if withFault(err, context) {
		return 1
	}

	gasToUse := ethContext.GasSchedule().EthAPICost.CallCode
	gasToUse += ethContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(data))
	ethContext.UseGas(gasToUse)

	ethContext.Transfer(dest, send, big.NewInt(0).SetBytes(value), nil)

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  send,
			Arguments:   [][]byte{data},
			CallValue:   big.NewInt(0).SetBytes(value),
			GasPrice:    0,
			GasProvided: ethContext.BoundGasLimit(gasLimit),
		},
		RecipientAddr: dest,
		Function:      "main",
	}

	err = ethContext.ExecuteOnSameContext(contractCallInput)
	if err != nil {
		return 1
	}

	return 0
}

//export ethcallDelegate
func ethcallDelegate(context unsafe.Pointer, gasLimit int64, addressOffset int32, dataOffset int32, dataLength int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	value := ethContext.GetVMInput().CallValue
	sender := ethContext.GetVMInput().CallerAddr

	address, err := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.HashLen)
	if withFault(err, context) {
		return 1
	}

	data, err := arwen.LoadBytes(instCtx.Memory(), dataOffset, dataLength)
	if withFault(err, context) {
		return 1
	}

	gasToUse := ethContext.GasSchedule().EthAPICost.CallDelegate
	gasToUse += ethContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(data))
	ethContext.UseGas(gasToUse)

	ethContext.Transfer(address, sender, value, nil)

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   [][]byte{data},
			CallValue:   value,
			GasPrice:    0,
			GasProvided: ethContext.BoundGasLimit(gasLimit),
		},
		RecipientAddr: address,
		Function:      "main",
	}

	err = ethContext.ExecuteOnSameContext(contractCallInput)
	if err != nil {
		return 1
	}

	return 0
}

//export ethcallStatic
func ethcallStatic(context unsafe.Pointer, gasLimit int64, addressOffset int32, dataOffset int32, dataLength int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	address, err := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLenEth)
	if withFault(err, context) {
		return 1
	}

	data, err := arwen.LoadBytes(instCtx.Memory(), dataOffset, dataLength)
	if withFault(err, context) {
		return 1
	}

	value := ethContext.GetVMInput().CallValue
	sender := ethContext.GetVMInput().CallerAddr

	gasToUse := ethContext.GasSchedule().EthAPICost.CallStatic
	gasToUse += ethContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(data))
	ethContext.UseGas(gasToUse)

	if IsAddressForPredefinedContract(address) {
		err := CallPredefinedContract(context, address, data)
		if err != nil {
			return 1
		}

		return 0
	}

	ethContext.Transfer(address, sender, value, nil)

	ethContext.SetReadOnly(true)

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   [][]byte{data},
			CallValue:   value,
			GasPrice:    0,
			GasProvided: ethContext.BoundGasLimit(gasLimit),
		},
		RecipientAddr: address,
		Function:      "main",
	}

	err = ethContext.ExecuteOnSameContext(contractCallInput)

	ethContext.SetReadOnly(false)
	if err != nil {
		return 1
	}

	return 0
}

//export ethcreate
func ethcreate(context unsafe.Pointer, valueOffset int32, dataOffset int32, length int32, resultOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	sender := ethContext.GetSCAddress()
	value, err := arwen.LoadBytes(instCtx.Memory(), valueOffset, arwen.BalanceLen)
	if withFault(err, context) {
		return 1
	}

	data, err := arwen.LoadBytes(instCtx.Memory(), dataOffset, length)
	if withFault(err, context) {
		return 1
	}

	gasToUse := ethContext.GasSchedule().EthAPICost.Create
	ethContext.UseGas(gasToUse)
	gasLimit := ethContext.GasLeft()

	contractCreate := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   nil,
			CallValue:   big.NewInt(0).SetBytes(value),
			GasPrice:    0,
			GasProvided: gasLimit,
		},
		ContractCode: data,
	}

	newAddress, err := ethContext.CreateNewContract(contractCreate)
	if err != nil {
		return 1
	}

	err = arwen.StoreBytes(instCtx.Memory(), resultOffset, newAddress)
	if withFault(err, context) {
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

func withFault(err error, context unsafe.Pointer) bool {
	if err != nil {
		instCtx := wasmer.IntoInstanceContext(context)
		hostContext := arwen.GetEthContext(instCtx.Data())
		hostContext.SignalUserError()
		hostContext.UseGas(hostContext.GasLeft())

		return true
	}

	return false
}

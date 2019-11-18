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

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, ethContext.GetSCAddress())

	gasToUse := ethContext.GasSchedule().EthAPICost.GetAddress
	ethContext.UseGas(gasToUse)
}

//export ethgetExternalBalance
func ethgetExternalBalance(context unsafe.Pointer, addressOffset int32, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	address := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	balance := ethContext.GetBalance(address)

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, balance)

	gasToUse := ethContext.GasSchedule().EthAPICost.GetExternalBalance
	ethContext.UseGas(gasToUse)
}

//export ethgetBlockHash
func ethgetBlockHash(context unsafe.Pointer, number int64, resultOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	hash := ethContext.BlockHash(number)
	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, hash)

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
	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, callData[dataOffset:dataOffset+length])

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

	key := arwen.LoadBytes(instCtx.Memory(), pathOffset, arwen.HashLen)
	data := arwen.LoadBytes(instCtx.Memory(), valueOffset, arwen.HashLen)

	_ = ethContext.SetStorage(ethContext.GetSCAddress(), key, data)

	gasToUse := ethContext.GasSchedule().EthAPICost.StorageStore
	gasToUse += ethContext.GasSchedule().BaseOperationCost.StorePerByte * uint64(len(data))
	ethContext.UseGas(gasToUse)
}

//export ethstorageLoad
func ethstorageLoad(context unsafe.Pointer, pathOffset int32, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	key := arwen.LoadBytes(instCtx.Memory(), pathOffset, arwen.HashLen)
	data := ethContext.GetStorage(ethContext.GetSCAddress(), key)

	currInput := make([]byte, arwen.HashLen)
	copy(currInput[arwen.HashLen-len(data):], data)

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, currInput)

	gasToUse := ethContext.GasSchedule().EthAPICost.StorageLoad
	gasToUse += ethContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(data))
	ethContext.UseGas(gasToUse)
}

//export ethgetCaller
func ethgetCaller(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	caller := ethContext.GetVMInput().CallerAddr

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, caller)
	gasToUse := ethContext.GasSchedule().EthAPICost.GetCaller
	ethContext.UseGas(gasToUse)
}

//export ethgetCallValue
func ethgetCallValue(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	value := ethContext.GetVMInput().CallValue.Bytes()
	length := len(value)
	invBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		invBytes[length-i-1] = value[i]
	}

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, invBytes)

	gasToUse := ethContext.GasSchedule().EthAPICost.GetCallValue
	ethContext.UseGas(gasToUse)
}

//export ethcodeCopy
func ethcodeCopy(context unsafe.Pointer, resultOffset int32, codeOffset int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	scAddress := ethContext.GetSCAddress()
	code := ethContext.GetCode(scAddress)

	if int32(len(code)) > codeOffset+length {
		_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, code[codeOffset:codeOffset+length])
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

	dest := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	code := ethContext.GetCode(dest)

	if int32(len(code)) > codeOffset+length {
		_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, code[codeOffset:codeOffset+length])
	}

	gasToUse := ethContext.GasSchedule().EthAPICost.ExternalCodeCopy
	gasToUse += ethContext.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(length)
	ethContext.UseGas(gasToUse)
}

//export ethgetExternalCodeSize
func ethgetExternalCodeSize(context unsafe.Pointer, addressOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	dest := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)

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

	_ = arwen.StoreBytes(instCtx.Memory(), valueOffset, gasU128)

	gasToUse := ethContext.GasSchedule().EthAPICost.GetTxGasPrice
	ethContext.UseGas(gasToUse)
}

//export ethlogTopics
func ethlogTopics(context unsafe.Pointer, dataOffset int32, length int32, numberOfTopics int32, topic1 int32, topic2 int32, topic3 int32, topic4 int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

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

	ethContext.WriteLog(ethContext.GetSCAddress(), topicsData, data)

	gasToUse := ethContext.GasSchedule().EthAPICost.Log
	gasToUse += ethContext.GasSchedule().BaseOperationCost.StorePerByte * (4*arwen.HashLen + uint64(length))
	ethContext.UseGas(gasToUse)
}

//export ethgetTxOrigin
func ethgetTxOrigin(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, ethContext.GetVMInput().CallerAddr)
	gasToUse := ethContext.GasSchedule().EthAPICost.GetTxOrigin
	ethContext.UseGas(gasToUse)
}

//export ethfinish
func ethfinish(context unsafe.Pointer, resultOffset int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	data := arwen.LoadBytes(instCtx.Memory(), resultOffset, length)
	ethContext.Finish(data)

	gasToUse := ethContext.GasSchedule().EthAPICost.Finish
	gasToUse += ethContext.GasSchedule().BaseOperationCost.StorePerByte * uint64(length)
	ethContext.UseGas(gasToUse)
}

//export ethrevert
func ethrevert(context unsafe.Pointer, dataOffset int32, length int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	data := arwen.LoadBytes(instCtx.Memory(), dataOffset, length)
	ethContext.Finish(data)
	ethContext.SignalUserError()

	gasToUse := ethContext.GasSchedule().EthAPICost.Revert
	gasToUse += ethContext.GasSchedule().BaseOperationCost.StorePerByte * uint64(length)
	ethContext.UseGas(gasToUse)
}

//export ethselfDestruct
func ethselfDestruct(context unsafe.Pointer, addressOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	address := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.HashLen)

	ethContext.SelfDestruct(address, ethContext.GetVMInput().CallerAddr)

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

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, ethReturnData[dataOffset:dataOffset+length])
}

//export ethgetBlockCoinbase
func ethgetBlockCoinbase(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	gasToUse := ethContext.GasSchedule().EthAPICost.GetBlockCoinbase
	ethContext.UseGas(gasToUse)

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, ethContext.BlockChainHook().CurrentRandomSeed())
}

//export ethgetBlockDifficulty
func ethgetBlockDifficulty(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	gasToUse := ethContext.GasSchedule().EthAPICost.GetBlockCoinbase
	ethContext.UseGas(gasToUse)

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, ethContext.BlockChainHook().CurrentRandomSeed())
}

//export ethcall
func ethcall(context unsafe.Pointer, gasLimit int64, addressOffset int32, valueOffset int32, dataOffset int32, dataLength int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	send := ethContext.GetSCAddress()
	dest := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	value := arwen.LoadBytes(instCtx.Memory(), valueOffset, arwen.BalanceLen)
	data := arwen.LoadBytes(instCtx.Memory(), dataOffset, dataLength)

	gasToUse := ethContext.GasSchedule().EthAPICost.Call
	gasToUse += ethContext.GasSchedule().BaseOperationCost.StorePerByte * uint64(len(data))
	ethContext.UseGas(gasToUse)

	if ethContext.GasLeft() < uint64(gasLimit) {
		return 1
	}

	bigIntVal := big.NewInt(0).SetBytes(value)
	ethContext.Transfer(dest, send, bigIntVal, nil)

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  send,
			Arguments:   [][]byte{data},
			CallValue:   bigIntVal,
			GasPrice:    0,
			GasProvided: uint64(gasLimit),
		},
		RecipientAddr: dest,
		Function:      "main",
	}
	err := ethContext.ExecuteOnDestContext(contractCallInput)
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
	dest := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.AddressLen)
	value := arwen.LoadBytes(instCtx.Memory(), valueOffset, arwen.BalanceLen)
	data := arwen.LoadBytes(instCtx.Memory(), dataOffset, dataLength)

	gasToUse := ethContext.GasSchedule().EthAPICost.CallCode
	gasToUse += ethContext.GasSchedule().BaseOperationCost.StorePerByte * uint64(len(data))
	ethContext.UseGas(gasToUse)

	if ethContext.GasLeft() < uint64(gasLimit) {
		return 1
	}

	ethContext.Transfer(dest, send, big.NewInt(0).SetBytes(value), nil)
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  send,
			Arguments:   [][]byte{data},
			CallValue:   big.NewInt(0).SetBytes(value),
			GasPrice:    0,
			GasProvided: uint64(gasLimit),
		},
		RecipientAddr: dest,
		Function:      "main",
	}
	err := ethContext.ExecuteOnSameContext(contractCallInput)
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

	address := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.HashLen)
	data := arwen.LoadBytes(instCtx.Memory(), dataOffset, dataLength)

	gasToUse := ethContext.GasSchedule().EthAPICost.CallDelegate
	gasToUse += ethContext.GasSchedule().BaseOperationCost.StorePerByte * uint64(len(data))
	ethContext.UseGas(gasToUse)

	if ethContext.GasLeft() < uint64(gasLimit) {
		return 1
	}

	ethContext.Transfer(address, sender, value, nil)
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   [][]byte{data},
			CallValue:   value,
			GasPrice:    0,
			GasProvided: uint64(gasLimit),
		},
		RecipientAddr: address,
		Function:      "main",
	}
	err := ethContext.ExecuteOnSameContext(contractCallInput)
	if err != nil {
		return 1
	}

	return 0
}

//export ethcallStatic
func ethcallStatic(context unsafe.Pointer, gasLimit int64, addressOffset int32, dataOffset int32, dataLength int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	ethContext := arwen.GetEthContext(instCtx.Data())

	address := arwen.LoadBytes(instCtx.Memory(), addressOffset, arwen.HashLen)
	data := arwen.LoadBytes(instCtx.Memory(), dataOffset, dataLength)

	value := ethContext.GetVMInput().CallValue
	sender := ethContext.GetVMInput().CallerAddr

	gasToUse := ethContext.GasSchedule().EthAPICost.CallStatic
	gasToUse += ethContext.GasSchedule().BaseOperationCost.StorePerByte * uint64(len(data))
	ethContext.UseGas(gasToUse)

	if ethContext.GasLeft() < uint64(gasLimit) {
		return 1
	}

	ethContext.Transfer(address, sender, value, nil)

	ethContext.SetReadOnly(true)
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   [][]byte{data},
			CallValue:   value,
			GasPrice:    0,
			GasProvided: uint64(gasLimit),
		},
		RecipientAddr: address,
		Function:      "main",
	}
	err := ethContext.ExecuteOnSameContext(contractCallInput)
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
	value := arwen.LoadBytes(instCtx.Memory(), valueOffset, arwen.BalanceLen)
	data := arwen.LoadBytes(instCtx.Memory(), dataOffset, length)

	gasToUse := ethContext.GasSchedule().EthAPICost.Create
	gasToUse += ethContext.GasSchedule().BaseOperationCost.StorePerByte * uint64(len(data))
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

	_ = arwen.StoreBytes(instCtx.Memory(), resultOffset, newAddress)

	return 0
}

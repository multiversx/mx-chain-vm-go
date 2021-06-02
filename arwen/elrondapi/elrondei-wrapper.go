package elrondapi

import (
	"encoding/binary"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
)

// TransferESDTNFTExecute - exported wrapper over transferESDTNFTExecute
// (this is used from unit tests)
func TransferESDTNFTExecuteWithTypes(
	host arwen.VMHost,
	recipientAddr []byte,
	tokenName []byte,
	esdtValue uint64,
	functionName []byte,
	args [][]byte,
	gasLimit uint64,
) int32 {
	runtime := host.Runtime()

	/*
		TODO in the future this should be necessary only if we want to test arguments
		serialization and deserialization, but for API logic we should call a method
		with "real" type arguments and not array of bytes
		E.g. transferESDTNFTExecute(host, functionName, value, ...)
	*/
	offsetDest := int32(0)
	runtime.MemStore(int32(offsetDest), recipientAddr)

	offsetTokenID := offsetDest + int32(len(recipientAddr))
	tokenLen := int32(len(tokenName))
	runtime.MemStore(offsetTokenID, tokenName)

	value := big.NewInt(int64(esdtValue)).Bytes()
	offsetValue := offsetTokenID + tokenLen
	valueLen := int32(arwen.BalanceLen)
	value = arwen.PadBytesLeft(value, arwen.BalanceLen)
	runtime.MemStore(int32(offsetValue), value)

	offsetFunction := offsetValue + valueLen
	funcNameLen := int32(len(functionName))
	runtime.MemStore(offsetFunction, []byte(functionName))

	noOfArguments := len(args)
	argumentsLengths := make([]uint32, noOfArguments)
	for idxArg, arg := range args {
		argumentsLengths[idxArg] = uint32(len(arg))
	}

	offsetArgumentsLength := offsetFunction + funcNameLen
	argumentsLengthsAsBytes := make([]byte, noOfArguments*4)
	binary.LittleEndian.PutUint32(argumentsLengthsAsBytes, argumentsLengths[0])
	runtime.MemStore(int32(offsetArgumentsLength), argumentsLengthsAsBytes)

	offsetArgumentsData := offsetArgumentsLength + int32(len(argumentsLengthsAsBytes))
	crtArgOffset := offsetArgumentsData
	for idxArg, arg := range args {
		argumentsData := make([]byte, argumentsLengths[idxArg])
		copy(argumentsData, arg)
		runtime.MemStore(crtArgOffset, argumentsData)
		crtArgOffset += int32(argumentsLengths[idxArg])
	}

	return TransferESDTNFTExecuteWithHost(
		host,
		offsetDest,
		offsetTokenID,
		tokenLen,
		offsetValue,
		0,
		int64(gasLimit),
		offsetFunction,
		funcNameLen,
		1,
		offsetArgumentsLength,
		offsetArgumentsData)
}

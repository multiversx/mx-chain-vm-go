package elrondapi

import (
	"unsafe"

	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
)

// TransferESDTNFTExecute - exported wrapper over transferESDTNFTExecute
func TransferESDTNFTExecute(
	instanceMock *mock.InstanceMock,
	destOffset int32,
	tokenIDOffset int32,
	tokenIDLen int32,
	valueOffset int32,
	nonce int64,
	gasLimit int64,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	return transferESDTNFTExecute(unsafe.Pointer(instanceMock),
		destOffset,
		tokenIDOffset,
		tokenIDLen,
		valueOffset,
		nonce,
		gasLimit,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset)
}

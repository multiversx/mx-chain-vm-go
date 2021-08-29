package elrondapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern void	v1_4_managedSCAddress(void *context, int32_t addressHandle);
// extern void	v1_4_managedOwnerAddress(void *context, int32_t addressHandle);
// extern void	v1_4_managedCaller(void *context, int32_t addressHandle);
// extern void	v1_4_managedSignalError(void* context, int32_t errHandle1);
// extern void	v1_4_managedWriteLog(void* context, int32_t topicsHandle, int32_t dataHandle);
//
import "C"

import (
	"encoding/binary"
	"errors"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/math"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/wasmer"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// ManagedEIImports creates a new wasmer.Imports populated with variants of the API methods that use managed types only.
func ManagedEIImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("env")

	imports, err := imports.Append("managedSCAddress", v1_4_managedSCAddress, C.v1_4_managedSCAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedOwnerAddress", v1_4_managedOwnerAddress, C.v1_4_managedOwnerAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedCaller", v1_4_managedCaller, C.v1_4_managedCaller)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedSignalError", v1_4_managedSignalError, C.v1_4_managedSignalError)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("managedWriteLog", v1_4_managedWriteLog, C.v1_4_managedWriteLog)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export v1_4_managedSCAddress
func v1_4_managedSCAddress(context unsafe.Pointer, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetSCAddress
	metering.UseGas(gasToUse)

	scAddress := runtime.GetSCAddress()

	managedType.SetBytes(destinationHandle, scAddress)
}

//export v1_4_managedOwnerAddress
func v1_4_managedOwnerAddress(context unsafe.Pointer, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	blockchain := arwen.GetBlockchainContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetOwnerAddress
	metering.UseGas(gasToUse)

	owner, err := blockchain.GetOwnerAddress()
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	managedType.SetBytes(destinationHandle, owner)
}

//export v1_4_managedCaller
func v1_4_managedCaller(context unsafe.Pointer, destinationHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.GetCaller
	metering.UseGas(gasToUse)

	caller := runtime.GetVMInput().CallerAddr
	managedType.SetBytes(destinationHandle, caller)
}

//export v1_4_managedSignalError
func v1_4_managedSignalError(context unsafe.Pointer, errHandle int32) {
	managedType := arwen.GetManagedTypesContext(context)
	runtime := arwen.GetRuntimeContext(context)
	metering := arwen.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().ElrondAPICost.SignalError
	metering.UseGas(gasToUse)

	errBytes, err := managedType.GetBytes(errHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForThisIntNumberOfBytes(len(errBytes))

	runtime.SignalUserError(string(errBytes))
}

//export v1_4_managedWriteLog
func v1_4_managedWriteLog(
	context unsafe.Pointer,
	topicsHandle int32,
	dataHandle int32,
) {

	host := arwen.GetVMHost(context)
	runtime := arwen.GetRuntimeContext(context)
	output := arwen.GetOutputContext(context)
	metering := arwen.GetMeteringContext(context)
	managedType := arwen.GetManagedTypesContext(context)

	topics, sumOfTopicByteLengths, err := readManagedVecOfManagedBuffers(host, topicsHandle)
	if arwen.WithFault(err, context, runtime.ElrondAPIErrorShouldFailExecution()) {
		return
	}

	dataBytes, err := managedType.GetBytes(dataHandle)
	if arwen.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForThisIntNumberOfBytes(len(dataBytes))
	dataByteLen := uint64(len(dataBytes))

	gasToUse := metering.GasSchedule().ElrondAPICost.Log
	gasForData := math.MulUint64(
		metering.GasSchedule().BaseOperationCost.DataCopyPerByte,
		uint64(sumOfTopicByteLengths+dataByteLen))
	gasToUse = math.AddUint64(gasToUse, gasForData)
	metering.UseGas(gasToUse)

	output.WriteLog(runtime.GetSCAddress(), topics, dataBytes)
}

func readManagedVecOfManagedBuffers(
	host arwen.VMHost,
	managedVecHandle int32,
) ([][]byte, uint64, error) {
	managedType := host.ManagedTypes()

	managedVecBytes, err := managedType.GetBytes(managedVecHandle)
	if err != nil {
		return nil, 0, err
	}
	managedType.ConsumeGasForThisIntNumberOfBytes(len(managedVecBytes))

	if len(managedVecBytes)%4 != 0 {
		return nil, 0, errors.New("invalid managed vector of managed buffers")
	}

	numBuffers := len(managedVecBytes) / 4
	result := make([][]byte, 0, numBuffers)
	sumOfItemByteLengths := uint64(0)
	for i := 0; i < len(managedVecBytes); i += 4 {
		itemHandle := int32(binary.BigEndian.Uint32(managedVecBytes[i : i+4]))

		itemBytes, err := managedType.GetBytes(itemHandle)
		if err != nil {
			return nil, 0, err
		}
		managedType.ConsumeGasForThisIntNumberOfBytes(len(itemBytes))

		sumOfItemByteLengths += uint64(len(itemBytes))
		result = append(result, itemBytes)
	}

	return result, sumOfItemByteLengths, nil
}

func writeManagedVecOfManagedBuffers(
	host arwen.VMHost,
	data [][]byte,
	destinationHandle int32,
) (uint64, error) {
	managedType := host.ManagedTypes()

	sumOfItemByteLengths := uint64(0)
	destinationBytes := make([]byte, 4*len(data))
	dataIndex := 0
	for _, itemBytes := range data {
		sumOfItemByteLengths += uint64(len(itemBytes))
		itemHandle := managedType.NewManagedBufferFromBytes(itemBytes)
		binary.BigEndian.PutUint32(destinationBytes[dataIndex:dataIndex+4], uint32(itemHandle))
		dataIndex += 4
	}

	managedType.SetBytes(destinationHandle, destinationBytes)

	return sumOfItemByteLengths, nil
}

func readESDTTransfer(
	host arwen.VMHost,
	data []byte,
) (*vmcommon.ESDTTransfer, error) {
	managedType := host.ManagedTypes()

	if len(data) != 16 {
		return nil, errors.New("invalid ESDT transfer object encoding")
	}

	tokenIdentifierHandle := int32(binary.BigEndian.Uint32(data[0:4]))
	tokenIdentifier, err := managedType.GetBytes(tokenIdentifierHandle)
	if err != nil {
		return nil, err
	}
	nonce := binary.BigEndian.Uint64(data[4:12])
	valueHandle := int32(binary.BigEndian.Uint32(data[12:16]))
	value, err := managedType.GetBigInt(valueHandle)
	if err != nil {
		return nil, err
	}

	return &vmcommon.ESDTTransfer{
		ESDTTokenName:  tokenIdentifier,
		ESDTTokenType:  0, // TODO
		ESDTTokenNonce: nonce,
		ESDTValue:      value,
	}, nil
}

func writeESDTTransfer(
	host arwen.VMHost,
	object *vmcommon.ESDTTransfer,
) ([]byte, error) {
	managedType := host.ManagedTypes()

	tokenIdentifierHandle := managedType.NewManagedBufferFromBytes(object.ESDTTokenName)
	valueHandle := managedType.NewBigInt(object.ESDTValue)

	destinationBytes := make([]byte, 16)
	binary.BigEndian.PutUint32(destinationBytes[0:4], uint32(tokenIdentifierHandle))
	binary.BigEndian.PutUint64(destinationBytes[4:12], object.ESDTTokenNonce)
	binary.BigEndian.PutUint32(destinationBytes[12:16], uint32(valueHandle))

	return destinationBytes, nil
}

package vmhooks

import (
	"encoding/binary"
	"errors"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/math"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

const esdtTransferLen = 16

// Deserializes a vmcommon.ESDTTransfer object.
func readESDTTransfer(
	managedType vmhost.ManagedTypesContext,
	runtime vmhost.RuntimeContext,
	data []byte,
) (*vmcommon.ESDTTransfer, error) {
	if len(data) != esdtTransferLen {
		return nil, errors.New("invalid ESDT transfer object encoding")
	}

	tokenIdentifierHandle := int32(binary.BigEndian.Uint32(data[0:4]))
	tokenIdentifier, err := managedType.GetBytes(tokenIdentifierHandle)
	if err != nil {
		return nil, err
	}

	err = managedType.ConsumeGasForBytes(tokenIdentifier)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return nil, err
	}

	nonce := binary.BigEndian.Uint64(data[4:12])
	valueHandle := int32(binary.BigEndian.Uint32(data[12:16]))
	value, err := managedType.GetBigInt(valueHandle)
	if err != nil {
		return nil, err
	}

	err = managedType.ConsumeGasForBigIntCopy(value)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return nil, err
	}

	tokenType := core.Fungible
	if nonce > 0 {
		tokenType = core.NonFungible
	}

	return &vmcommon.ESDTTransfer{
		ESDTTokenName:  tokenIdentifier,
		ESDTTokenType:  uint32(tokenType),
		ESDTTokenNonce: nonce,
		ESDTValue:      value,
	}, nil
}

// Converts a managed buffer of serialized data to a slice of ESDTTransfer.
// The format is:
// - token identifier handle - 4 bytes
// - nonce - 8 bytes
// - value handle - 4 bytes
// Total: 16 bytes.
func readESDTTransfers(
	managedType vmhost.ManagedTypesContext,
	runtime vmhost.RuntimeContext,
	managedVecHandle int32,
) ([]*vmcommon.ESDTTransfer, error) {
	managedVecBytes, err := managedType.GetBytes(managedVecHandle)
	if err != nil {
		return nil, err
	}

	err = managedType.ConsumeGasForBytes(managedVecBytes)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		return nil, err
	}

	if len(managedVecBytes)%esdtTransferLen != 0 {
		return nil, errors.New("invalid managed vector of ESDT transfers")
	}

	numTransfers := len(managedVecBytes) / esdtTransferLen
	result := make([]*vmcommon.ESDTTransfer, 0, numTransfers)
	for i := 0; i < len(managedVecBytes); i += esdtTransferLen {
		esdtTransfer, err := readESDTTransfer(managedType, runtime, managedVecBytes[i:i+esdtTransferLen])
		if err != nil {
			return nil, err
		}
		result = append(result, esdtTransfer)
	}

	return result, nil
}

// Serializes a vmcommon.ESDTTransfer object.
func writeESDTTransfer(
	managedType vmhost.ManagedTypesContext,
	esdtTransfer *vmcommon.ESDTTransfer,
	destinationBytes []byte,
) {
	tokenIdentifierHandle := managedType.NewManagedBufferFromBytes(esdtTransfer.ESDTTokenName)
	valueHandle := managedType.NewBigInt(esdtTransfer.ESDTValue)

	binary.BigEndian.PutUint32(destinationBytes[0:4], uint32(tokenIdentifierHandle))
	binary.BigEndian.PutUint64(destinationBytes[4:12], esdtTransfer.ESDTTokenNonce)
	binary.BigEndian.PutUint32(destinationBytes[12:16], uint32(valueHandle))
}

// Serializes a list of ESDTTransfer one after the other into a byte slice.
// The format is (for each ESDTTransfer):
// - token identifier handle - 4 bytes
// - nonce - 8 bytes
// - value handle - 4 bytes
// Total: 16 bytes.
func writeESDTTransfersToBytes(
	managedType vmhost.ManagedTypesContext,
	esdtTransfers []*vmcommon.ESDTTransfer,
) []byte {
	destinationBytes := make([]byte, esdtTransferLen*len(esdtTransfers))
	dataIndex := 0
	for _, esdtTransfer := range esdtTransfers {
		writeESDTTransfer(managedType, esdtTransfer, destinationBytes[dataIndex:dataIndex+esdtTransferLen])
		dataIndex += esdtTransferLen
	}

	return destinationBytes
}

type vmInputData struct {
	destination []byte
	function    string
	value       *big.Int
	arguments   [][]byte
}

func readDestinationValueFunctionArguments(
	host vmhost.VMHost,
	destHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
) (*vmInputData, error) {
	managedType := host.ManagedTypes()

	vmInput, err := readDestinationValueArguments(host, destHandle, valueHandle, argumentsHandle)
	if err != nil {
		return nil, err
	}

	function, err := managedType.GetBytes(functionHandle)
	if err != nil {
		return nil, err
	}
	vmInput.function = string(function)

	return vmInput, nil
}

func readDestinationValueArguments(
	host vmhost.VMHost,
	destHandle int32,
	valueHandle int32,
	argumentsHandle int32,
) (*vmInputData, error) {
	managedType := host.ManagedTypes()

	vmInput, err := readDestinationArguments(host, destHandle, argumentsHandle)
	if err != nil {
		return nil, err
	}

	vmInput.value, err = managedType.GetBigInt(valueHandle)
	if err != nil {
		return nil, err
	}

	return vmInput, nil
}

func readDestinationFunctionArguments(
	host vmhost.VMHost,
	destHandle int32,
	functionHandle int32,
	argumentsHandle int32,
) (*vmInputData, error) {
	managedType := host.ManagedTypes()

	vmInput, err := readDestinationArguments(host, destHandle, argumentsHandle)
	if err != nil {
		return nil, err
	}

	function, err := managedType.GetBytes(functionHandle)
	if err != nil {
		return nil, err
	}
	vmInput.function = string(function)

	return vmInput, nil
}

func readDestinationArguments(
	host vmhost.VMHost,
	destHandle int32,
	argumentsHandle int32,
) (*vmInputData, error) {
	managedType := host.ManagedTypes()
	metering := host.Metering()

	var err error
	vmInput := &vmInputData{}

	vmInput.destination, err = managedType.GetBytes(destHandle)
	if err != nil {
		return nil, err
	}

	vmInput.value = big.NewInt(0)
	data, actualLen, err := managedType.ReadManagedVecOfManagedBuffers(argumentsHandle)
	if err != nil {
		return nil, err
	}
	vmInput.arguments = data

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, actualLen)
	err = metering.UseGasBounded(gasToUse)
	if err != nil && host.Runtime().UseGasBoundedShouldFailExecution() {
		return nil, err
	}

	return vmInput, nil
}

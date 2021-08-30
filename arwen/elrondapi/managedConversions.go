package elrondapi

import (
	"encoding/binary"
	"errors"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/math"
	"github.com/ElrondNetwork/elrond-go-core/core"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

const handleLen = 4
const esdtTransferLen = 16

// Converts a managed buffer of managed buffers to a slice of byte slices.
func readManagedVecOfManagedBuffers(
	managedType arwen.ManagedTypesContext,
	managedVecHandle int32,
) ([][]byte, uint64, error) {
	managedVecBytes, err := managedType.GetBytes(managedVecHandle)
	if err != nil {
		return nil, 0, err
	}
	managedType.ConsumeGasForThisIntNumberOfBytes(len(managedVecBytes))

	if len(managedVecBytes)%handleLen != 0 {
		return nil, 0, errors.New("invalid managed vector of managed buffer handles")
	}

	numBuffers := len(managedVecBytes) / handleLen
	result := make([][]byte, 0, numBuffers)
	sumOfItemByteLengths := uint64(0)
	for i := 0; i < len(managedVecBytes); i += handleLen {
		itemHandle := int32(binary.BigEndian.Uint32(managedVecBytes[i : i+handleLen]))

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

// Converts a slice of byte slices to a managed buffer of managed buffers.
func writeManagedVecOfManagedBuffers(
	managedType arwen.ManagedTypesContext,
	data [][]byte,
	destinationHandle int32,
) uint64 {
	sumOfItemByteLengths := uint64(0)
	destinationBytes := make([]byte, handleLen*len(data))
	dataIndex := 0
	for _, itemBytes := range data {
		sumOfItemByteLengths += uint64(len(itemBytes))
		itemHandle := managedType.NewManagedBufferFromBytes(itemBytes)
		binary.BigEndian.PutUint32(destinationBytes[dataIndex:dataIndex+handleLen], uint32(itemHandle))
		dataIndex += handleLen
	}

	managedType.SetBytes(destinationHandle, destinationBytes)

	return sumOfItemByteLengths
}

// Deserializes a vmcommon.ESDTTransfer object.
func readESDTTransfer(
	managedType arwen.ManagedTypesContext,
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
	managedType.ConsumeGasForThisIntNumberOfBytes(len(tokenIdentifier))
	nonce := binary.BigEndian.Uint64(data[4:12])
	valueHandle := int32(binary.BigEndian.Uint32(data[12:16]))
	value, err := managedType.GetBigInt(valueHandle)
	if err != nil {
		return nil, err
	}
	managedType.ConsumeGasForBigIntCopy(value)

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
	managedType arwen.ManagedTypesContext,
	managedVecHandle int32,
) ([]*vmcommon.ESDTTransfer, error) {
	managedVecBytes, err := managedType.GetBytes(managedVecHandle)
	if err != nil {
		return nil, err
	}
	managedType.ConsumeGasForThisIntNumberOfBytes(len(managedVecBytes))

	if len(managedVecBytes)%esdtTransferLen != 0 {
		return nil, errors.New("invalid managed vector of ESDT transfers")
	}

	numTransfers := len(managedVecBytes) / esdtTransferLen
	result := make([]*vmcommon.ESDTTransfer, 0, numTransfers)
	for i := 0; i < len(managedVecBytes); i += handleLen {
		esdtTransfer, err := readESDTTransfer(managedType, managedVecBytes[i:i+handleLen])
		if err != nil {
			return nil, err
		}
		result = append(result, esdtTransfer)
	}

	return result, nil
}

// Serializes a vmcommon.ESDTTransfer object.
func writeESDTTransfer(
	managedType arwen.ManagedTypesContext,
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
	managedType arwen.ManagedTypesContext,
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
	vmHost arwen.VMHost,
	destHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
) (*vmInputData, error) {
	managedType := vmHost.ManagedTypes()

	vmInput, err := readDestinationValueArguments(vmHost, destHandle, valueHandle, argumentsHandle)
	if err != nil {
		return nil, err
	}

	function, err := managedType.GetBytes(functionHandle)
	if err != nil {
		return nil, err
	}
	vmInput.function = string(function)

	return vmInput, err
}

func readDestinationValueArguments(
	vmHost arwen.VMHost,
	destHandle int32,
	valueHandle int32,
	argumentsHandle int32,
) (*vmInputData, error) {
	managedType := vmHost.ManagedTypes()

	vmInput, err := readDestinationArguments(vmHost, destHandle, argumentsHandle)
	if err != nil {
		return nil, err
	}

	vmInput.value, err = managedType.GetBigInt(valueHandle)
	if err != nil {
		return nil, err
	}

	return vmInput, err
}

func readDestinationFunctionArguments(
	vmHost arwen.VMHost,
	destHandle int32,
	functionHandle int32,
	argumentsHandle int32,
) (*vmInputData, error) {
	managedType := vmHost.ManagedTypes()

	vmInput, err := readDestinationArguments(vmHost, destHandle, argumentsHandle)
	if err != nil {
		return nil, err
	}

	function, err := managedType.GetBytes(functionHandle)
	if err != nil {
		return nil, err
	}
	vmInput.function = string(function)

	return vmInput, err
}

func readDestinationArguments(
	vmHost arwen.VMHost,
	destHandle int32,
	argumentsHandle int32,
) (*vmInputData, error) {
	managedType := vmHost.ManagedTypes()
	metering := vmHost.Metering()

	var err error
	vmInput := &vmInputData{}

	vmInput.destination, err = managedType.GetBytes(destHandle)
	if err != nil {
		return nil, err
	}

	vmInput.value = big.NewInt(0)
	data, actualLen, err := readManagedVecOfManagedBuffers(managedType, argumentsHandle)
	if err != nil {
		return nil, err
	}
	vmInput.arguments = data

	gasToUse := math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, actualLen)
	metering.UseGas(gasToUse)

	return vmInput, err
}

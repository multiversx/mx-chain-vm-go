package elrondapi

import (
	"encoding/binary"
	"errors"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/elrond-go-core/core"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func readManagedVecOfManagedBuffers(
	managedType arwen.ManagedTypesContext,
	managedVecHandle int32,
) ([][]byte, uint64, error) {
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
	managedType arwen.ManagedTypesContext,
	data [][]byte,
	destinationHandle int32,
) uint64 {
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

	return sumOfItemByteLengths
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

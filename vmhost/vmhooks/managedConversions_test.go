package vmhooks

import (
	"encoding/binary"
	"math/big"
	"testing"

	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/require"
)

type mockManagedTypesForTransfer struct {
	vmhost.ManagedTypesContext
	bytesMap  map[int32][]byte
	bigIntMap map[int32]*big.Int
}

func (m *mockManagedTypesForTransfer) GetBytes(handle int32) ([]byte, error) {
	return m.bytesMap[handle], nil
}

func (m *mockManagedTypesForTransfer) ConsumeGasForBytes(_ []byte) error {
	return nil
}

func (m *mockManagedTypesForTransfer) GetBigInt(handle int32) (*big.Int, error) {
	return m.bigIntMap[handle], nil
}

func (m *mockManagedTypesForTransfer) ConsumeGasForBigIntCopy(_ ...*big.Int) error {
	return nil
}

func TestReadESDTTransfer_SignValidation(t *testing.T) {
	t.Parallel()

	managedTypesCtx := &mockManagedTypesForTransfer{
		bytesMap: map[int32][]byte{
			1: []byte("TEST-123456"),
		},
		bigIntMap: map[int32]*big.Int{
			10: big.NewInt(100),
			11: big.NewInt(0),
			12: big.NewInt(-50),
		},
	}

	runtimeCtx := &contextmock.RuntimeContextMock{}

	// Positive value handle (10)
	dataPos := make([]byte, 16)
	binary.BigEndian.PutUint32(dataPos[0:4], 1)
	binary.BigEndian.PutUint64(dataPos[4:12], 0)
	binary.BigEndian.PutUint32(dataPos[12:16], 10)

	transfer, err := readESDTTransfer(managedTypesCtx, runtimeCtx, dataPos)
	require.Nil(t, err)
	require.Equal(t, big.NewInt(100), transfer.ESDTValue)

	// Zero value handle (11)
	dataZero := make([]byte, 16)
	binary.BigEndian.PutUint32(dataZero[0:4], 1)
	binary.BigEndian.PutUint64(dataZero[4:12], 0)
	binary.BigEndian.PutUint32(dataZero[12:16], 11)

	_, err = readESDTTransfer(managedTypesCtx, runtimeCtx, dataZero)
	require.Equal(t, vmhost.ErrTransferNegativeValue, err)

	// Negative value handle (12)
	dataNeg := make([]byte, 16)
	binary.BigEndian.PutUint32(dataNeg[0:4], 1)
	binary.BigEndian.PutUint64(dataNeg[4:12], 0)
	binary.BigEndian.PutUint32(dataNeg[12:16], 12)

	_, err = readESDTTransfer(managedTypesCtx, runtimeCtx, dataNeg)
	require.Equal(t, vmhost.ErrTransferNegativeValue, err)
}

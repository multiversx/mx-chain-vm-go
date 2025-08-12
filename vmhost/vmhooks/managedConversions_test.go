package vmhooks

import (
	"math/big"
	"testing"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/mock/mockery"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestReadESDTTransfer(t *testing.T) {
	t.Parallel()

	managedType := &mockery.MockManagedTypesContext{}
	runtime := &mockery.MockRuntimeContext{}

	managedType.On("GetBytes", mock.Anything).Return([]byte("token-name"), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	managedType.On("GetBigInt", mock.Anything).Return(big.NewInt(100), nil)
	managedType.On("ConsumeGasForBigIntCopy", mock.Anything).Return(nil)

	data := make([]byte, 16)
	esdtTransfer, err := readESDTTransfer(managedType, runtime, data)
	require.Nil(t, err)
	require.NotNil(t, esdtTransfer)
}

func TestReadESDTTransfers(t *testing.T) {
	t.Parallel()

	managedType := &mockery.MockManagedTypesContext{}
	runtime := &mockery.MockRuntimeContext{}

	managedType.On("GetBytes", mock.Anything).Return(make([]byte, 32), nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	managedType.On("GetBigInt", mock.Anything).Return(big.NewInt(100), nil)
	managedType.On("ConsumeGasForBigIntCopy", mock.Anything).Return(nil)

	esdtTransfers, err := readESDTTransfers(managedType, runtime, 0)
	require.Nil(t, err)
	require.NotNil(t, esdtTransfers)
	require.Len(t, esdtTransfers, 2)
}

func TestWriteESDTTransfer(t *testing.T) {
	t.Parallel()

	managedType := &mockery.MockManagedTypesContext{}

	managedType.On("NewManagedBufferFromBytes", mock.Anything).Return(int32(1))
	managedType.On("NewBigInt", mock.Anything).Return(int32(2))

	esdtTransfer := &vmcommon.ESDTTransfer{
		ESDTTokenName:  []byte("token-name"),
		ESDTTokenNonce: 123,
		ESDTValue:      big.NewInt(100),
	}
	data := make([]byte, 16)
	writeESDTTransfer(managedType, esdtTransfer, data)
}

func TestWriteESDTTransfersToBytes(t *testing.T) {
	t.Parallel()

	managedType := &mockery.MockManagedTypesContext{}

	managedType.On("NewManagedBufferFromBytes", mock.Anything).Return(int32(1))
	managedType.On("NewBigInt", mock.Anything).Return(int32(2))

	esdtTransfers := []*vmcommon.ESDTTransfer{
		{
			ESDTTokenName:  []byte("token-name"),
			ESDTTokenNonce: 123,
			ESDTValue:      big.NewInt(100),
		},
		{
			ESDTTokenName:  []byte("token-name2"),
			ESDTTokenNonce: 456,
			ESDTValue:      big.NewInt(200),
		},
	}
	data := writeESDTTransfersToBytes(managedType, esdtTransfers)
	require.Len(t, data, 32)
}

func TestReadDestinationValueFunctionArguments(t *testing.T) {
	t.Parallel()

	host := &mockery.MockVMHost{}
	managedType := &mockery.MockManagedTypesContext{}
	metering := &mockery.MockMeteringContext{}

	host.On("ManagedTypes").Return(managedType)
	host.On("Metering").Return(metering)
	gasSchedule, _ := config.CreateGasConfig(config.MakeGasMapForTests())
	metering.On("GasSchedule").Return(gasSchedule)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("GetBigInt", mock.Anything).Return(big.NewInt(100), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", mock.Anything).Return([][]byte{[]byte("arg1")}, uint64(1), nil)
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	vmInput, err := readDestinationValueFunctionArguments(host, 0, 0, 0, 0)
	require.Nil(t, err)
	require.NotNil(t, vmInput)
}

func TestReadDestinationValueArguments(t *testing.T) {
	t.Parallel()

	host := &mockery.MockVMHost{}
	managedType := &mockery.MockManagedTypesContext{}
	metering := &mockery.MockMeteringContext{}

	host.On("ManagedTypes").Return(managedType)
	host.On("Metering").Return(metering)
	gasSchedule, _ := config.CreateGasConfig(config.MakeGasMapForTests())
	metering.On("GasSchedule").Return(gasSchedule)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("GetBigInt", mock.Anything).Return(big.NewInt(100), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", mock.Anything).Return([][]byte{[]byte("arg1")}, uint64(1), nil)
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	vmInput, err := readDestinationValueArguments(host, 0, 0, 0)
	require.Nil(t, err)
	require.NotNil(t, vmInput)
}

func TestReadDestinationFunctionArguments(t *testing.T) {
	t.Parallel()

	host := &mockery.MockVMHost{}
	managedType := &mockery.MockManagedTypesContext{}
	metering := &mockery.MockMeteringContext{}

	host.On("ManagedTypes").Return(managedType)
	host.On("Metering").Return(metering)
	gasSchedule, _ := config.CreateGasConfig(config.MakeGasMapForTests())
	metering.On("GasSchedule").Return(gasSchedule)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", mock.Anything).Return([][]byte{[]byte("arg1")}, uint64(1), nil)
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	vmInput, err := readDestinationFunctionArguments(host, 0, 0, 0)
	require.Nil(t, err)
	require.NotNil(t, vmInput)
}

func TestReadDestinationArguments(t *testing.T) {
	t.Parallel()

	host := &mockery.MockVMHost{}
	managedType := &mockery.MockManagedTypesContext{}
	metering := &mockery.MockMeteringContext{}

	host.On("ManagedTypes").Return(managedType)
	host.On("Metering").Return(metering)
	gasSchedule, _ := config.CreateGasConfig(config.MakeGasMapForTests())
	metering.On("GasSchedule").Return(gasSchedule)
	managedType.On("GetBytes", mock.Anything).Return([]byte("data"), nil)
	managedType.On("ReadManagedVecOfManagedBuffers", mock.Anything).Return([][]byte{[]byte("arg1")}, uint64(1), nil)
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	vmInput, err := readDestinationArguments(host, 0, 0)
	require.Nil(t, err)
	require.NotNil(t, vmInput)
}

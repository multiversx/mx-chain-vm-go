package contexts

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

func TestNewMeteringContext(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostMock{}

	meteringContext, err := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))
	require.Nil(t, err)
	require.NotNil(t, meteringContext)
}

func TestNewMeteringContext_NilGasSchedule(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostMock{}

	meteringContext, err := NewMeteringContext(host, nil, uint64(15000))
	require.NotNil(t, err)
	require.Nil(t, meteringContext)
}

func TestMeteringContext_GasSchedule(t *testing.T) {
	t.Parallel()

	host := &contextmock.VMHostStub{}
	meteringContext, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	schedule := meteringContext.GasSchedule()
	require.NotNil(t, schedule)
}

func TestMeteringContext_UseGas(t *testing.T) {
	t.Parallel()

	mockRuntime := &contextmock.RuntimeContextMock{}
	host := &contextmock.VMHostMock{
		RuntimeContext: mockRuntime,
	}
	meteringContext, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	gasProvided := uint64(1001)
	vmInput := &vmcommon.VMInput{GasProvided: gasProvided}
	mockRuntime.SetVMInput(vmInput)
	gas := uint64(1000)
	meteringContext.UseGas(gas)
	require.Equal(t, mockRuntime.GetPointsUsed(), gas)
	require.Equal(t, uint64(1), meteringContext.GasLeft())

	gasProvided = uint64(10000)
	vmInput = &vmcommon.VMInput{GasProvided: gasProvided}
	mockRuntime.SetPointsUsed(0)
	mockRuntime.SetVMInput(vmInput)
	meteringContext, _ = NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	require.Equal(t, gasProvided, meteringContext.GasLeft())
	meteringContext.UseGas(gas)
	require.Equal(t, gasProvided-gas, meteringContext.GasLeft())
}

func TestMeteringContext_FreeGas(t *testing.T) {
	t.Parallel()

	mockOutput := &contextmock.OutputContextMock{}
	host := &contextmock.VMHostMock{
		OutputContext: mockOutput,
	}

	meteringContext, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	mockOutput.GasRefund = big.NewInt(0)
	meteringContext.FreeGas(1000)
	gas := mockOutput.GetRefund()
	require.Equal(t, uint64(1000), gas)

	meteringContext.FreeGas(100)
	gas = mockOutput.GetRefund()
	require.Equal(t, uint64(1100), gas)
}

func TestMeteringContext_BoundGasLimit(t *testing.T) {
	t.Parallel()

	mockRuntime := &contextmock.RuntimeContextMock{}
	host := &contextmock.VMHostMock{
		RuntimeContext: mockRuntime,
	}
	blockGasLimit := uint64(15000)
	meteringContext, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	gasProvided := uint64(10000)
	vmInput := &vmcommon.VMInput{GasProvided: gasProvided}
	mockRuntime.SetVMInput(vmInput)
	mockRuntime.SetPointsUsed(0)

	gasLimit := uint64(5000)
	limit := meteringContext.BoundGasLimit(gasLimit)
	require.Equal(t, gasLimit, limit)

	gasLimit = 25000
	limit = meteringContext.BoundGasLimit(gasLimit)
	require.Equal(t, meteringContext.GasLeft(), limit)

	blockLimit := meteringContext.BlockGasLimit()
	require.Equal(t, blockGasLimit, blockLimit)
}

func TestMeteringContext_DeductInitialGasForExecution(t *testing.T) {
	t.Parallel()

	mockRuntime := &contextmock.RuntimeContextMock{}
	gasProvided := uint64(10000)
	vmInput := &vmcommon.VMInput{
		GasProvided: gasProvided,
	}

	mockRuntime.SetVMInput(vmInput)

	host := &contextmock.VMHostMock{
		RuntimeContext: mockRuntime,
	}

	meteringContext, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	contract := []byte("contract")
	err := meteringContext.DeductInitialGasForExecution(contract)
	require.Nil(t, err)

	vmInput.GasProvided = 1
	err = meteringContext.DeductInitialGasForExecution(contract)
	require.Equal(t, arwen.ErrNotEnoughGas, err)
}

func TestDeductInitialGasForDirectDeployment(t *testing.T) {
	t.Parallel()
	mockRuntime := &contextmock.RuntimeContextMock{}
	gasProvided := uint64(10000)
	contractCode := []byte("contractCode")
	input := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			GasProvided: gasProvided,
		},
		ContractCode: contractCode,
	}

	mockRuntime.SetVMInput(&input.VMInput)

	host := &contextmock.VMHostMock{
		RuntimeContext: mockRuntime,
	}

	meteringContext, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	mockRuntime.SetPointsUsed(0)
	err := meteringContext.DeductInitialGasForDirectDeployment(arwen.CodeDeployInput{ContractCode: contractCode})
	require.Nil(t, err)
	remainingGas := meteringContext.GasLeft()
	require.Equal(t, gasProvided-uint64(len(contractCode))-1, remainingGas)

	input.GasProvided = 2
	mockRuntime.SetPointsUsed(0)
	err = meteringContext.DeductInitialGasForDirectDeployment(arwen.CodeDeployInput{ContractCode: contractCode})
	require.Equal(t, arwen.ErrNotEnoughGas, err)
}

func TestDeductInitialGasForIndirectDeployment(t *testing.T) {
	t.Parallel()

	mockRuntime := &contextmock.RuntimeContextMock{}
	gasProvided := uint64(10000)
	contractCode := []byte("contractCode")
	input := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			GasProvided: gasProvided,
		},
		ContractCode: contractCode,
	}

	mockRuntime.SetVMInput(&input.VMInput)

	host := &contextmock.VMHostMock{
		RuntimeContext: mockRuntime,
	}

	meteringContext, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	mockRuntime.SetPointsUsed(0)
	err := meteringContext.DeductInitialGasForIndirectDeployment(arwen.CodeDeployInput{ContractCode: contractCode})
	require.Nil(t, err)
	remainingGas := meteringContext.GasLeft()
	require.Equal(t, gasProvided-uint64(len(contractCode)), remainingGas)

	input.GasProvided = 2
	mockRuntime.SetPointsUsed(0)
	err = meteringContext.DeductInitialGasForDirectDeployment(arwen.CodeDeployInput{ContractCode: contractCode})
	require.Equal(t, arwen.ErrNotEnoughGas, err)
}

func TestMeteringContext_AsyncCallGasLocking(t *testing.T) {
	t.Parallel()

	mockRuntime := &contextmock.RuntimeContextMock{}
	contractSize := uint64(1000)
	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallType: vmcommon.AsynchronousCall,
		},
	}

	mockRuntime.SCCodeSize = contractSize
	mockRuntime.SetVMInput(&input.VMInput)
	mockRuntime.SetPointsUsed(0)

	host := &contextmock.VMHostMock{
		RuntimeContext: mockRuntime,
	}

	meteringContext, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	input.GasProvided = 1
	err := meteringContext.UseGasForAsyncStep()
	require.Equal(t, arwen.ErrNotEnoughGas, err)

	mockRuntime.SetPointsUsed(0)
	gasProvided := uint64(1_000_000)
	input.GasProvided = gasProvided
	gasToLock := meteringContext.ComputeGasLockedForAsync()
	err = meteringContext.UseGasBounded(gasToLock)
	require.Nil(t, err)
	expectedGasLeft := gasProvided - gasToLock
	require.Equal(t, expectedGasLeft, meteringContext.GasLeft())

	mockRuntime.VMInput.CallType = vmcommon.AsynchronousCallBack
	mockRuntime.VMInput.GasLocked = gasToLock
	meteringContext.UnlockGasIfAsyncCallback()
	err = meteringContext.UseGasForAsyncStep()
	require.Nil(t, err)
	require.Equal(t, gasProvided-1, meteringContext.GasLeft())
}

package contexts

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

func TestNewMeteringContext(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostMock{}

	meteringContext, err := NewMeteringContext(host, config.MakeGasMap(1), uint64(15000))
	require.Nil(t, err)
	require.NotNil(t, meteringContext)
}

func TestNewMeteringContext_NilGasSchedule(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostMock{}

	meteringContext, err := NewMeteringContext(host, nil, uint64(15000))
	require.NotNil(t, err)
	require.Nil(t, meteringContext)
}

func TestMeteringContext_GasSchedule(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostStub{}
	meteringContext, _ := NewMeteringContext(host, config.MakeGasMap(1), uint64(15000))

	schedule := meteringContext.GasSchedule()
	require.NotNil(t, schedule)
}

func TestMeteringContext_UseGas(t *testing.T) {
	t.Parallel()

	mockRuntime := &mock.RuntimeContextMock{}
	vmInput := &vmcommon.VMInput{GasProvided: 0}
	mockRuntime.SetVMInput(vmInput)
	host := &mock.VmHostMock{
		RuntimeContext: mockRuntime,
	}
	meteringContext, _ := NewMeteringContext(host, config.MakeGasMap(1), uint64(15000))

	gas := uint64(1000)
	meteringContext.UseGas(gas)
	require.Equal(t, mockRuntime.GetPointsUsed(), gas)
	require.Equal(t, uint64(0), meteringContext.GasLeft())

	gasProvided := uint64(10000)
	vmInput = &vmcommon.VMInput{GasProvided: gasProvided}
	mockRuntime.SetVMInput(vmInput)
	mockRuntime.SetPointsUsed(0)
	meteringContext, _ = NewMeteringContext(host, config.MakeGasMap(1), uint64(15000))

	require.Equal(t, gasProvided, meteringContext.GasLeft())
	meteringContext.UseGas(gas)
	require.Equal(t, gasProvided-gas, meteringContext.GasLeft())
}

func TestMeteringContext_FreeGas(t *testing.T) {
	t.Parallel()

	mockOutput := &mock.OutputContextMock{}
	host := &mock.VmHostMock{
		OutputContext: mockOutput,
	}

	meteringContext, _ := NewMeteringContext(host, config.MakeGasMap(1), uint64(15000))

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

	mockRuntime := &mock.RuntimeContextMock{}
	host := &mock.VmHostMock{
		RuntimeContext: mockRuntime,
	}
	blockGasLimit := uint64(15000)
	meteringContext, _ := NewMeteringContext(host, config.MakeGasMap(1), uint64(15000))

	gasProvided := uint64(10000)
	vmInput := &vmcommon.VMInput{GasProvided: gasProvided}
	mockRuntime.SetVMInput(vmInput)
	mockRuntime.SetPointsUsed(0)

	gasLimit := 5000
	limit := meteringContext.BoundGasLimit(int64(gasLimit))
	require.Equal(t, uint64(gasLimit), limit)

	gasLimit = 25000
	limit = meteringContext.BoundGasLimit(int64(gasLimit))
	require.Equal(t, meteringContext.GasLeft(), limit)

	blockLimit := meteringContext.BlockGasLimit()
	require.Equal(t, blockGasLimit, blockLimit)
}

func TestMeteringContext_DeductInitialGasForExecution(t *testing.T) {
	t.Parallel()

	mockRuntime := &mock.RuntimeContextMock{}
	gasProvided := uint64(10000)
	vmInput := &vmcommon.VMInput{
		GasProvided: gasProvided,
	}

	mockRuntime.SetVMInput(vmInput)

	host := &mock.VmHostMock{
		RuntimeContext: mockRuntime,
	}

	meteringContext, _ := NewMeteringContext(host, config.MakeGasMap(1), uint64(15000))

	contract := []byte("contract")
	err := meteringContext.DeductInitialGasForExecution(contract)
	require.Nil(t, err)

	vmInput.GasProvided = 1
	err = meteringContext.DeductInitialGasForExecution(contract)
	require.Equal(t, arwen.ErrNotEnoughGas, err)
}

func TestDeductInitialGasForDirectDeployment(t *testing.T) {
	t.Parallel()
	mockRuntime := &mock.RuntimeContextMock{}
	gasProvided := uint64(10000)
	contractCode := []byte("contractCode")
	input := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			GasProvided: gasProvided,
		},
		ContractCode: contractCode,
	}

	mockRuntime.SetVMInput(&input.VMInput)

	host := &mock.VmHostMock{
		RuntimeContext: mockRuntime,
	}

	meteringContext, _ := NewMeteringContext(host, config.MakeGasMap(1), uint64(15000))

	mockRuntime.SetPointsUsed(0)
	err := meteringContext.DeductInitialGasForDirectDeployment(input)
	require.Nil(t, err)
	remainingGas := meteringContext.GasLeft()
	require.Equal(t, gasProvided-uint64(len(contractCode))-1, remainingGas)

	input.GasProvided = 2
	mockRuntime.SetPointsUsed(0)
	err = meteringContext.DeductInitialGasForDirectDeployment(input)
	require.Equal(t, arwen.ErrNotEnoughGas, err)
}

func TestDeductInitialGasForIndirectDeployment(t *testing.T) {
	t.Parallel()

	mockRuntime := &mock.RuntimeContextMock{}
	gasProvided := uint64(10000)
	contractCode := []byte("contractCode")
	input := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			GasProvided: gasProvided,
		},
		ContractCode: contractCode,
	}

	mockRuntime.SetVMInput(&input.VMInput)

	host := &mock.VmHostMock{
		RuntimeContext: mockRuntime,
	}

	meteringContext, _ := NewMeteringContext(host, config.MakeGasMap(1), uint64(15000))

	mockRuntime.SetPointsUsed(0)
	err := meteringContext.DeductInitialGasForIndirectDeployment(input)
	require.Nil(t, err)
	remainingGas := meteringContext.GasLeft()
	require.Equal(t, gasProvided-uint64(len(contractCode)), remainingGas)

	input.GasProvided = 2
	mockRuntime.SetPointsUsed(0)
	err = meteringContext.DeductInitialGasForDirectDeployment(input)
	require.Equal(t, arwen.ErrNotEnoughGas, err)
}

func TestMeteringContext_AsyncCallGasLocking(t *testing.T) {
	t.Parallel()

	mockRuntime := &mock.RuntimeContextMock{}
	contractCode := []byte("contractCode")
	input := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallType: vmcommon.AsynchronousCall,
		},
		ContractCode: contractCode,
	}

	mockRuntime.SetVMInput(&input.VMInput)
	mockRuntime.SetPointsUsed(0)

	host := &mock.VmHostMock{
		RuntimeContext: mockRuntime,
	}

	meteringContext, _ := NewMeteringContext(host, config.MakeGasMap(1), uint64(15000))

	input.GasProvided = 2
	err := meteringContext.deductAndLockGasIfAsyncStep()
	require.Equal(t, arwen.ErrNotEnoughGas, err)

	gasProvided := uint64(10000)
	input.GasProvided = gasProvided
	err = meteringContext.deductAndLockGasIfAsyncStep()
	require.Nil(t, err)
	require.Equal(t, uint64(2), meteringContext.gasLockedForAsyncStep)
	require.Equal(t, gasProvided-3, meteringContext.GasLeft())

	meteringContext.UnlockGasIfAsyncStep()
	require.Equal(t, gasProvided-1, meteringContext.GasLeft())
}

package contexts

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

func TestNewMeteringContext(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostStub{}
	blockGasLimit := uint64(15000)
	gasSchedule := config.MakeGasMap(1)

	metContext, err := NewMeteringContext(host, gasSchedule, blockGasLimit)
	require.Nil(t, err)
	require.NotNil(t, metContext)
}

func TestNewMeteringContext_NilGasSchedule(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostStub{}
	blockGasLimit := uint64(15000)

	metContext, err := NewMeteringContext(host, nil, blockGasLimit)
	require.NotNil(t, err)
	require.Nil(t, metContext)
}

func TestMeteringContext_GasSchedule(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostStub{}
	blockGasLimit := uint64(15000)
	gasSchedule := config.MakeGasMap(1)
	metContext, _ := NewMeteringContext(host, gasSchedule, blockGasLimit)

	gSchedule := metContext.GasSchedule()
	require.NotNil(t, gSchedule)
}

func TestMeteringContext_UseGas(t *testing.T) {
	t.Parallel()

	runtimeMock := mock.NewRuntimeContextMock()
	host := &mock.VmHostStub{
		RuntimeCalled: func() arwen.RuntimeContext {
			return runtimeMock
		},
	}
	blockGasLimit := uint64(15000)
	gasSchedule := config.MakeGasMap(1)
	metContext, _ := NewMeteringContext(host, gasSchedule, blockGasLimit)

	gas := uint64(1000)
	metContext.UseGas(gas)
	require.Equal(t, runtimeMock.GetPointsUsed(), gas)

	gasLeft := metContext.GasLeft()
	require.Equal(t, uint64(0), gasLeft)

	gasProvided := uint64(10000)
	vmInput := &vmcommon.VMInput{GasProvided: gasProvided}
	runtimeMock.SetVMInput(vmInput)

	gasLeft = metContext.GasLeft()
	require.Equal(t, gasProvided-gas, gasLeft)
}

func TestMeteringContext_FreeGas(t *testing.T) {
	t.Parallel()

	outputContextMock := mock.NewOutputContextMock()
	host := &mock.VmHostStub{
		OutputCalled: func() arwen.OutputContext {
			return outputContextMock
		},
	}

	blockGasLimit := uint64(15000)
	gasSchedule := config.MakeGasMap(1)
	metContext, _ := NewMeteringContext(host, gasSchedule, blockGasLimit)

	metContext.FreeGas(1000)
	gas := outputContextMock.GetRefund()
	require.Equal(t, uint64(1000), gas)
}

func TestMeteringContext_BoundGasLimit(t *testing.T) {
	t.Parallel()

	runtimeMock := mock.NewRuntimeContextMock()
	host := &mock.VmHostStub{
		RuntimeCalled: func() arwen.RuntimeContext {
			return runtimeMock
		},
	}
	blockGasLimit := uint64(15000)
	gasSchedule := config.MakeGasMap(1)
	metContext, _ := NewMeteringContext(host, gasSchedule, blockGasLimit)

	gasProvided := uint64(10000)
	vmInput := &vmcommon.VMInput{GasProvided: gasProvided}
	runtimeMock.SetVMInput(vmInput)

	gasLimit := 5000
	limit := metContext.BoundGasLimit(int64(gasLimit))
	require.Equal(t, uint64(gasLimit), limit)

	blockLimit := metContext.BlockGasLimit()
	require.Equal(t, blockGasLimit, blockLimit)
}

func TestDeductInitialGasForExecution(t *testing.T) {
	t.Parallel()

	gasProvided := uint64(10000)
	host := &mock.VmHostStub{}
	blockGasLimit := uint64(15000)
	gasSchedule := config.MakeGasMap(1)
	metContext, _ := NewMeteringContext(host, gasSchedule, blockGasLimit)

	contract := []byte("contract")
	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			GasProvided: gasProvided,
		},
	}

	remainingGas, err := metContext.DeductInitialGasForExecution(input, contract)
	require.Nil(t, err)
	require.Equal(t, gasProvided-uint64(len(contract)), remainingGas)

	input.GasProvided = 1
	remainingGas, err = metContext.DeductInitialGasForExecution(input, contract)
	require.NotNil(t, err)
	require.Zero(t, remainingGas)
}

func TestDeductInitialGasForDirectDeployment(t *testing.T) {
	t.Parallel()

	gasProvided := uint64(10000)
	host := &mock.VmHostStub{}
	blockGasLimit := uint64(15000)
	gasSchedule := config.MakeGasMap(1)
	metContext, _ := NewMeteringContext(host, gasSchedule, blockGasLimit)

	contracCode := []byte("contractCode")
	input := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			GasProvided: gasProvided,
		},
		ContractCode: contracCode,
	}

	remainingGas, err := metContext.DeductInitialGasForDirectDeployment(input)
	require.Equal(t, gasProvided-uint64(len(contracCode)+1), remainingGas)
	require.Nil(t, err)
}

func TestDeductInitialGasForIndirectDeployment(t *testing.T) {
	t.Parallel()

	gasProvided := uint64(10000)
	host := &mock.VmHostStub{}
	blockGasLimit := uint64(15000)
	gasSchedule := config.MakeGasMap(1)
	metContext, _ := NewMeteringContext(host, gasSchedule, blockGasLimit)

	contracCode := []byte("contractCode")
	input := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			GasProvided: gasProvided,
		},
		ContractCode: contracCode,
	}

	remainingGas, err := metContext.DeductInitialGasForIndirectDeployment(input)
	require.Equal(t, gasProvided-uint64(len(contracCode)), remainingGas)
	require.Nil(t, err)
}

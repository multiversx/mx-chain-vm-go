package contexts

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/math"
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
	meteringContext.gasForExecution = gasProvided
	gas := uint64(1000)
	meteringContext.UseGas(gas)
	require.Equal(t, mockRuntime.GetPointsUsed(), gas)
	require.Equal(t, uint64(1), meteringContext.GasLeft())

	gasProvided = uint64(10000)
	mockRuntime.SetPointsUsed(0)
	meteringContext, _ = NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))
	meteringContext.gasForExecution = gasProvided

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
	meteringContext.gasForExecution = gasProvided
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
	meteringContext.gasForExecution = gasProvided
	gasToLock := meteringContext.ComputeGasLockedForAsync()
	err = meteringContext.UseGasBounded(gasToLock)
	require.Nil(t, err)
	expectedGasLeft := gasProvided - gasToLock
	require.Equal(t, expectedGasLeft, meteringContext.GasLeft())

	mockRuntime.VMInput.CallType = vmcommon.AsynchronousCallBack
	mockRuntime.VMInput.GasLocked = gasToLock
	require.Equal(t, gasToLock, meteringContext.GetGasLocked()) // ask about the use of this

	meteringContext.unlockGasIfAsyncCallback(&input.VMInput)
	err = meteringContext.UseGasForAsyncStep()
	require.Nil(t, err)
	require.Equal(t, gasProvided-1, meteringContext.GasLeft())
}

func TestMeteringContext_GasUsed_NoStacking(t *testing.T) {
	t.Parallel()

	mockRuntime := &contextmock.RuntimeContextMock{}
	host := &contextmock.VMHostMock{
		RuntimeContext: mockRuntime,
	}

	contractSize := uint64(1000)
	contract := make([]byte, contractSize)
	input := &vmcommon.ContractCallInput{VMInput: vmcommon.VMInput{}}

	mockRuntime.SCCodeSize = contractSize
	mockRuntime.SetVMInput(&input.VMInput)
	mockRuntime.SetPointsUsed(0)

	metering, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	input.GasProvided = 2000
	metering.InitStateFromContractCallInput(&input.VMInput)
	require.Equal(t, uint64(2000), metering.initialGasProvided)

	_ = metering.DeductInitialGasForExecution(contract)
	require.Equal(t, uint64(999), metering.GasLeft())

	metering.UseGas(400)
	require.Equal(t, uint64(599), metering.GasLeft())

	gasUsedByContract := metering.GasSpentByContract()
	require.Equal(t, uint64(1401), gasUsedByContract)
}

func setUpStackOneLevel(t *testing.T) (*meteringContext, *outputContext, *contextmock.RuntimeContextMock, *vmcommon.ContractCallInput, uint64) {
	t.Parallel()

	mockRuntime := &contextmock.RuntimeContextMock{}

	host := &contextmock.VMHostMock{
		RuntimeContext: mockRuntime,
	}

	output, _ := NewOutputContext(host)
	host.OutputContext = output

	contractSize := uint64(1000)
	contract := make([]byte, contractSize)
	mockRuntime.SCCodeSize = contractSize
	mockRuntime.SCAddress = []byte("parent")

	mockRuntime.SetPointsUsed(0)
	parentInput := &vmcommon.ContractCallInput{VMInput: vmcommon.VMInput{}}
	mockRuntime.SetVMInput(&parentInput.VMInput)

	metering, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))
	host.MeteringContext = metering
	zeroCodeCosts(metering)

	parentInput.GasProvided = 4000
	metering.InitStateFromContractCallInput(&parentInput.VMInput)
	require.Equal(t, 4000, int(metering.initialGasProvided))

	_ = metering.DeductInitialGasForExecution(contract)

	require.Equal(t, 3000, int(metering.GasLeft()))

	metering.UseGas(400)
	require.Equal(t, 2600, int(metering.GasLeft()))

	gasSpentByContract := metering.GasSpentByContract()
	require.Equal(t, 1400, int(gasSpentByContract))

	childInput := &vmcommon.ContractCallInput{VMInput: vmcommon.VMInput{}}
	childInput.GasProvided = 500

	metering.UseGas(childInput.GasProvided)
	parentPointsBeforeStacking := mockRuntime.GetPointsUsed()
	require.Equal(t, 900, int(parentPointsBeforeStacking))
	require.Equal(t, 2100, int(metering.GasLeft()))

	mockRuntime.SCAddress = []byte("child")
	mockRuntime.SetPointsUsed(0)
	mockRuntime.SetVMInput(&childInput.VMInput)
	metering.PushState()
	metering.InitStateFromContractCallInput(&childInput.VMInput)
	require.Equal(t, 500, int(metering.initialGasProvided))

	return metering, output, mockRuntime, parentInput, parentPointsBeforeStacking
}

func TestMeteringContext_GasUsed_StackOneLevel(t *testing.T) {
	metering, output, mockRuntime, parentInput, parentPointsBeforeStacking := setUpStackOneLevel(t)

	// child execution begins
	_ = metering.DeductInitialGasForExecution(make([]byte, 100))
	require.Equal(t, 400, int(metering.GasLeft()))

	metering.UseGas(50)
	gasRemaining := metering.GasLeft()
	require.Equal(t, 350, int(gasRemaining))

	gasSpentByContract := metering.GasSpentByContract()
	require.Equal(t, 150, int(gasSpentByContract))

	metering.UpdateGasStateOnSuccess(output.outputState)

	// return to the parent
	metering.PopMergeActiveState()
	mockRuntime.SCAddress = []byte("parent")
	mockRuntime.SetPointsUsed(parentPointsBeforeStacking)
	mockRuntime.SetVMInput(&parentInput.VMInput)

	metering.RestoreGas(gasRemaining)
	mockRuntime.IsContractOnStack = false
	require.Equal(t, 2450, int(metering.GasLeft()))

	gasSpentByContract = metering.GasSpentByContract()
	require.Equal(t, 1550, int(gasSpentByContract))

	metering.UseGas(50)
	require.Equal(t, 2400, int(metering.GasLeft()))

	gasSpentByContract = metering.GasSpentByContract()
	require.Equal(t, 1600, int(gasSpentByContract))

	metering.UpdateGasStateOnSuccess(output.outputState)

	gasUsed := output.outputState.OutputAccounts["parent"].GasUsed
	require.Equal(t, 1450, int(gasUsed))

	gasUsed = output.outputState.OutputAccounts["child"].GasUsed
	require.Equal(t, 150, int(gasUsed))

	gasRemaining = math.SubUint64(parentInput.GasProvided, gasSpentByContract)
	// calculate gas remaining

	require.Equal(t, int(gasRemaining), int(output.GetVMOutput().GasRemaining))
}

func TestMeteringContext_UpdateGasStateOnFailure_StackOneLevel(t *testing.T) {
	metering, output, mockRuntime, parentInput, parentPointsBeforeStacking := setUpStackOneLevel(t)
	// child execution begins

	_ = metering.DeductInitialGasForExecution(make([]byte, 600))
	require.Equal(t, 500, int(metering.GasLeft()))

	gasRemaining := metering.GasLeft()
	gasSpentByContract := metering.GasSpentByContract()
	require.Equal(t, 0, int(gasSpentByContract))

	metering.UpdateGasStateOnFailure(output.outputState)

	// return to the parent
	metering.PopSetActiveState() // ask if it should be like this
	mockRuntime.SCAddress = []byte("parent")
	mockRuntime.SetPointsUsed(parentPointsBeforeStacking)
	mockRuntime.SetVMInput(&parentInput.VMInput)

	metering.RestoreGas(gasRemaining)
	mockRuntime.IsContractOnStack = false
	require.Equal(t, 2600, int(metering.GasLeft()))

	gasSpentByContract = metering.GasSpentByContract()
	require.Equal(t, 1400, int(gasSpentByContract))

	metering.UseGas(3000)
	require.Equal(t, 0, int(metering.GasLeft()))

	gasSpentByContract = metering.GasSpentByContract()
	require.Equal(t, 4400, int(gasSpentByContract))

	metering.UpdateGasStateOnFailure(output.outputState)

	gasUsed := output.outputState.OutputAccounts["parent"].GasUsed
	require.Equal(t, 4000, int(gasUsed))

	gasUsed = output.outputState.OutputAccounts["child"].GasUsed
	require.Equal(t, 500, int(gasUsed))

	gasRemaining = math.SubUint64(parentInput.GasProvided, gasSpentByContract)
	// calculate gas remaining

	require.Equal(t, int(gasRemaining), int(output.GetVMOutput().GasRemaining))
}

func zeroCodeCosts(context *meteringContext) {
	//context.GasSchedule().BaseOperationCost.CompilePerByte = 0
	//context.GasSchedule().BaseOperationCost.AoTPreparePerByte = 0
	context.GasSchedule().BaseOperationCost.GetCode = 0
}

package contexts

import (
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/math"
	"github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/stretchr/testify/mock"









	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/require"
)

func TestNewMeteringContext(t *testing.T) {
	t.Parallel()
	const BlockGasLimit = uint64(15000)
	host := &context.MockVMHost{}

	meteringCtx, err := NewMeteringContext(host, config.MakeGasMapForTests(), BlockGasLimit)
	require.Nil(t, err)
	require.NotNil(t, meteringCtx)
	require.NotNil(t, meteringCtx.gasTracer)
}

func TestNewMeteringContext_NilGasSchedule(t *testing.T) {
	t.Parallel()
	const BlockGasLimit = uint64(15000)
	host := &context.MockVMHost{}

	meteringCtx, err := NewMeteringContext(host, nil, BlockGasLimit)
	require.NotNil(t, err)
	require.Nil(t, meteringCtx)
}

func TestMeteringContext_GasSchedule(t *testing.T) {
	t.Parallel()
	const BlockGasLimit = uint64(15000)

	host := &context.MockVMHost{}
	meteringCtx, _ := NewMeteringContext(host, config.MakeGasMapForTests(), BlockGasLimit)

	schedule := meteringCtx.GasSchedule()
	require.NotNil(t, schedule)
}

func TestMeteringContext_UseGas(t *testing.T) {
	t.Parallel()
	const BlockGasLimit = uint64(15000)

	mockRuntime := &context.MockRuntimeContext{}
	host := &context.MockVMHost{}
	host.On("Runtime").Return(mockRuntime)

	meteringCtx, _ := NewMeteringContext(host, config.MakeGasMapForTests(), BlockGasLimit)

	gasProvided := uint64(1001)
	meteringCtx.gasForExecution = gasProvided
	gasUsed := uint64(1000)
	mockRuntime.On("GetPointsUsed").Return(gasUsed)
	_ = meteringCtx.UseGasBounded(gasUsed)
	require.Equal(t, mockRuntime.GetPointsUsed(), gasUsed)
	require.Equal(t, gasProvided-gasUsed, meteringCtx.GasLeft())

	gasProvided = uint64(10000)
	mockRuntime.On("SetPointsUsed", uint64(0)).Return()
	meteringCtx, _ = NewMeteringContext(host, config.MakeGasMapForTests(), BlockGasLimit)
	meteringCtx.gasForExecution = gasProvided

	require.Equal(t, gasProvided, meteringCtx.GasLeft())
	_ = meteringCtx.UseGasBounded(gasUsed)
	require.Equal(t, gasProvided-gasUsed, meteringCtx.GasLeft())
}

func TestMeteringContext_FreeGas(t *testing.T) {
	t.Parallel()
	const BlockGasLimit = uint64(15000)

	mockOutput := &context.MockOutputContext{}
	host := &context.MockVMHost{}
	host.On("Output").Return(mockOutput)

	meteringCtx, _ := NewMeteringContext(host, config.MakeGasMapForTests(), BlockGasLimit)

	gasToFree := uint64(1000)
	mockOutput.On("GetRefund").Return(gasToFree)
	meteringCtx.FreeGas(gasToFree)
	gasRefunded := mockOutput.GetRefund()
	require.Equal(t, gasToFree, gasRefunded)

	moreGasToFree := uint64(100)
	mockOutput.On("GetRefund").Return(gasToFree + moreGasToFree)
	meteringCtx.FreeGas(moreGasToFree)
	gasRefunded = mockOutput.GetRefund()
	require.Equal(t, gasToFree+moreGasToFree, gasRefunded)
}

func TestMeteringContext_BoundGasLimit(t *testing.T) {
	t.Parallel()
	const BlockGasLimit = uint64(15000)

	mockRuntime := &context.MockRuntimeContext{}
	host := &context.MockVMHost{}
	host.On("Runtime").Return(mockRuntime)

	meteringCtx, _ := NewMeteringContext(host, config.MakeGasMapForTests(), BlockGasLimit)

	gasProvided := uint64(10000)
	meteringCtx.gasForExecution = gasProvided
	mockRuntime.On("SetPointsUsed", uint64(0)).Return()

	gasLimit := 5000
	limit := meteringCtx.BoundGasLimit(int64(gasLimit))
	require.Equal(t, uint64(gasLimit), limit)

	gasLimit = 25000
	limit = meteringCtx.BoundGasLimit(int64(gasLimit))
	require.Equal(t, meteringCtx.GasLeft(), limit)

	blockLimit := meteringCtx.BlockGasLimit()
	require.Equal(t, BlockGasLimit, blockLimit)
}

func TestMeteringContext_DeductInitialGasForExecution(t *testing.T) {
	t.Parallel()

	mockRuntime := &context.MockRuntimeContext{}
	gasProvided := uint64(10000)
	vmInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{GasProvided: gasProvided},
	}

	mockRuntime.On("GetVMInput").Return(vmInput)

	host := &context.MockVMHost{}
	host.On("Runtime").Return(mockRuntime)

	meteringCtx, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	contract := []byte("contract")
	err := meteringCtx.DeductInitialGasForExecution(contract)
	require.Nil(t, err)

	vmInput.GasProvided = 1
	err = meteringCtx.DeductInitialGasForExecution(contract)
	require.Equal(t, vmhost.ErrNotEnoughGas, err)
}

func TestDeductInitialGasForDirectDeployment(t *testing.T) {
	t.Parallel()
	mockRuntime := &context.MockRuntimeContext{}
	gasProvided := uint64(10000)
	contractCode := []byte("contractCode")
	input := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			GasProvided: gasProvided,
		},
		ContractCode: contractCode,
	}

	contractCallInput := &vmcommon.ContractCallInput{VMInput: input.VMInput}
	mockRuntime.On("GetVMInput").Return(contractCallInput)

	host := &context.MockVMHost{}
	host.On("Runtime").Return(mockRuntime)

	meteringCtx, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	mockRuntime.On("SetPointsUsed", uint64(0)).Return()
	err := meteringCtx.DeductInitialGasForDirectDeployment(vmhost.CodeDeployInput{ContractCode: contractCode})
	require.Nil(t, err)
	remainingGas := meteringCtx.GasLeft()
	require.Equal(t, gasProvided-uint64(len(contractCode))-1, remainingGas)

	contractCallInput.GasProvided = 2
	mockRuntime.On("SetPointsUsed", uint64(0)).Return()
	err = meteringCtx.DeductInitialGasForDirectDeployment(vmhost.CodeDeployInput{ContractCode: contractCode})
	require.Equal(t, vmhost.ErrNotEnoughGas, err)
}

func TestDeductInitialGasForIndirectDeployment(t *testing.T) {
	t.Parallel()

	mockRuntime := &context.MockRuntimeContext{}
	gasProvided := uint64(10000)
	contractCode := []byte("contractCode")
	input := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			GasProvided: gasProvided,
		},
		ContractCode: contractCode,
	}

	contractCallInput := &vmcommon.ContractCallInput{VMInput: input.VMInput}
	mockRuntime.On("GetVMInput").Return(contractCallInput)

	host := &context.MockVMHost{}
	host.On("Runtime").Return(mockRuntime)

	meteringCtx, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	mockRuntime.On("SetPointsUsed", uint64(0)).Return()
	err := meteringCtx.DeductInitialGasForIndirectDeployment(vmhost.CodeDeployInput{ContractCode: contractCode})
	require.Nil(t, err)
	remainingGas := meteringCtx.GasLeft()
	require.Equal(t, gasProvided-uint64(len(contractCode)), remainingGas)

	contractCallInput.GasProvided = 2
	mockRuntime.On("SetPointsUsed", uint64(0)).Return()
	err = meteringCtx.DeductInitialGasForDirectDeployment(vmhost.CodeDeployInput{ContractCode: contractCode})
	require.Equal(t, vmhost.ErrNotEnoughGas, err)
}

func TestMeteringContext_AsyncCallGasLocking(t *testing.T) {
	t.Parallel()

	mockRuntime := &context.MockRuntimeContext{}
	contractSize := uint64(1000)
	input := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallType: vm.AsynchronousCall,
		},
	}

	mockRuntime.On("GetVMInput").Return(input)
	mockRuntime.On("GetSCCodeSize").Return(contractSize)
	mockRuntime.On("SetPointsUsed", uint64(0)).Return()

	host := &context.MockVMHost{}
	host.On("Runtime").Return(mockRuntime)

	meteringCtx, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))

	input.GasProvided = 0
	err := meteringCtx.UseGasForAsyncStep()
	require.Equal(t, vmhost.ErrNotEnoughGas, err)

	mockRuntime.On("SetPointsUsed", uint64(0)).Return()
	gasProvided := uint64(1_000_000)
	input.GasProvided = gasProvided
	meteringCtx.gasForExecution = gasProvided
	gasToLock := meteringCtx.ComputeExtraGasLockedForAsync()
	err = meteringCtx.UseGasBounded(gasToLock)
	require.Nil(t, err)
	expectedGasLeft := gasProvided - gasToLock
	require.Equal(t, expectedGasLeft, meteringCtx.GasLeft())

	input.CallType = vm.AsynchronousCallBack
	input.GasLocked = gasToLock
	require.Equal(t, gasToLock, meteringCtx.GetGasLocked())

	meteringCtx.unlockGasIfAsyncCallback(&input.VMInput)
	err = meteringCtx.UseGasForAsyncStep()
	require.Nil(t, err)
	require.Equal(t, gasProvided-1, meteringCtx.GasLeft())
}

func TestMeteringContext_GasUsed_NoStacking(t *testing.T) {
	t.Parallel()
	const BlockGasLimit = uint64(15000)

	mockRuntime := &context.MockRuntimeContext{}
	host := &context.MockVMHost{}
	host.On("Runtime").Return(mockRuntime)

	contractSize := uint64(1000)
	contract := make([]byte, contractSize)
	input := &vmcommon.ContractCallInput{VMInput: vmcommon.VMInput{}}

	mockRuntime.On("GetSCCodeSize").Return(contractSize)
	mockRuntime.On("GetVMInput").Return(input)
	mockRuntime.On("SetPointsUsed", mock.Anything).Return()

	metering, _ := NewMeteringContext(host, config.MakeGasMapForTests(), BlockGasLimit)

	input.GasProvided = 2000
	metering.InitStateFromContractCallInput(&input.VMInput)
	require.Equal(t, input.GasProvided, metering.initialGasProvided)

	costPerByte := uint64(1)
	gasAfterDeductingInitial := metering.initialGasProvided - costPerByte - contractSize
	_ = metering.DeductInitialGasForExecution(contract)
	require.Equal(t, gasAfterDeductingInitial, metering.GasLeft())

	gasUsed := uint64(400)
	_ = metering.UseGasBounded(gasUsed)
	require.Equal(t, gasAfterDeductingInitial-gasUsed, metering.GasLeft())

	totalGasUsed := metering.initialGasProvided - metering.GasLeft()
	gasUsedByContract := metering.GasSpentByContract()
	require.Equal(t, totalGasUsed, gasUsedByContract)
}

func setUpStackOneLevel(t *testing.T, parentInput *vmcommon.ContractCallInput, childInput *vmcommon.ContractCallInput) (*context.MockVMHost, *context.MockRuntimeContext, uint64) {
	t.Parallel()

	host := &context.MockVMHost{}

	contractSize := uint64(1000)
	contract := make([]byte, contractSize)

	mockRuntime := &context.MockRuntimeContext{}
	host.On("Runtime").Return(mockRuntime)

	mockRuntime.On("GetSCCodeSize").Return(contractSize)
	mockRuntime.On("GetContextAddress").Return([]byte("parent"))

	mockRuntime.On("SetPointsUsed", uint64(0)).Return()
	mockRuntime.On("GetVMInput").Return(parentInput)

	metering, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))
	host.On("Metering").Return(metering)
	zeroCodeCosts(metering)

	output, _ := NewOutputContext(host)
	host.On("Output").Return(output)

	parentInput.GasProvided = 4000
	host.Metering().InitStateFromContractCallInput(&parentInput.VMInput)
	_ = host.Metering().DeductInitialGasForExecution(contract)

	parentPointsBeforeStacking := initStateFromChildGetParentPointsBeforeStacking(t, host, childInput)
	return host, mockRuntime, parentPointsBeforeStacking
}

func initStateFromChildGetParentPointsBeforeStacking(t *testing.T, host *context.MockVMHost, childInput *vmcommon.ContractCallInput) uint64 {
	parentGasProvided := uint64(4000)
	parentExecutionGas := uint64(1000)

	require.Equal(t, parentGasProvided-parentExecutionGas, host.Metering().GasLeft())

	parentUsedGas := uint64(400)
	_ = host.Metering().UseGasBounded(parentUsedGas)
	require.Equal(t, parentGasProvided-parentExecutionGas-parentUsedGas, host.Metering().GasLeft())

	gasSpentByContract := host.Metering().GasSpentByContract()
	require.Equal(t, parentExecutionGas+parentUsedGas, gasSpentByContract)

	childProvidedGas := childInput.GasProvided
	_ = host.Metering().UseGasBounded(childProvidedGas)
	parentPointsBeforeStacking := host.Runtime().GetPointsUsed()
	require.Equal(t, childProvidedGas+parentUsedGas, parentPointsBeforeStacking)
	require.Equal(t, parentGasProvided-parentExecutionGas-parentPointsBeforeStacking, host.Metering().GasLeft())

	host.Runtime().On("GetContextAddress").Return([]byte("child"))
	host.Runtime().On("SetPointsUsed", uint64(0)).Return()
	host.Runtime().On("GetVMInput").Return(childInput)
	host.Metering().PushState()
	host.Metering().InitStateFromContractCallInput(&childInput.VMInput)
	require.Equal(t, childProvidedGas, host.Metering().GetGasProvided())

	return parentPointsBeforeStacking
}

func TestMeteringContext_GasUsed_StackOneLevel(t *testing.T) {
	parentExecutionGas := uint64(1000)
	parentUsedGas := uint64(400)
	parentInput := &vmcommon.ContractCallInput{VMInput: vmcommon.VMInput{}}
	parentInput.CallerAddr = []byte("user")
	parentInput.RecipientAddr = []byte("parent")

	childInput := &vmcommon.ContractCallInput{VMInput: vmcommon.VMInput{}}
	childInput.GasProvided = 500
	childInput.CallerAddr = parentInput.RecipientAddr
	childInput.RecipientAddr = []byte("child")

	host, mockRuntime, parentPointsBeforeStacking := setUpStackOneLevel(t, parentInput, childInput)
	metering := host.MeteringContext
	output := host.OutputContext

	childExecutionGas := uint64(100)
	_ = metering.DeductInitialGasForExecution(make([]byte, childExecutionGas))
	require.Equal(t, childInput.GasProvided-childExecutionGas, metering.GasLeft())

	childUsedGas := uint64(50)
	_ = metering.UseGasBounded(childUsedGas)
	gasRemaining := metering.GasLeft()
	require.Equal(t, childInput.GasProvided-childExecutionGas-childUsedGas, metering.GasLeft())

	gasSpentByContract := metering.GasSpentByContract()
	require.Equal(t, childExecutionGas+childUsedGas, gasSpentByContract)

	_ = output.GetVMOutput()

	// return to the parent
	metering.PopMergeActiveState()
	mockRuntime.SCAddress = []byte("parent")
	mockRuntime.SetPointsUsed(parentPointsBeforeStacking)
	mockRuntime.SetVMInput(parentInput)

	metering.RestoreGas(gasRemaining)
	require.Equal(t, parentInput.GasProvided-parentExecutionGas-parentPointsBeforeStacking+gasRemaining, metering.GasLeft())

	gasSpentByContract = metering.GasSpentByContract()
	require.Equal(t, parentExecutionGas+parentUsedGas+childUsedGas+childExecutionGas, gasSpentByContract)

	_ = metering.UseGasBounded(50)
	parentUsedGas += 50
	require.Equal(t, parentInput.GasProvided-parentExecutionGas-parentUsedGas-childExecutionGas-childUsedGas, metering.GasLeft())

	gasSpentByContract = metering.GasSpentByContract()
	require.Equal(t, parentExecutionGas+parentUsedGas+childExecutionGas+childUsedGas, gasSpentByContract)

	vmOutput := output.GetVMOutput()

	gasUsed := vmOutput.OutputAccounts["parent"].GasUsed
	require.Equal(t, parentExecutionGas+parentUsedGas, gasUsed)

	gasUsed = vmOutput.OutputAccounts["child"].GasUsed
	require.Equal(t, childExecutionGas+childUsedGas, gasUsed)

	gasRemaining = math.SubUint64(parentInput.GasProvided, gasSpentByContract)
	// calculate gas remaining

	require.Equal(t, gasRemaining, metering.GasLeft())
}

func TestMeteringContext_UpdateGasStateOnFailure_StackOneLevel(t *testing.T) {

	parentExecutionGas := uint64(1000) // this is the contract size, but I chose to keep the convention used on child
	parentUsedGas := uint64(400)
	parentInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{},
	}

	childInput := &vmcommon.ContractCallInput{VMInput: vmcommon.VMInput{}}
	childInput.GasProvided = 500
	childInput.CallerAddr = parentInput.RecipientAddr

	host, mockRuntime, parentPointsBeforeStacking := setUpStackOneLevel(t, parentInput, childInput)

	metering := host.MeteringContext
	output, _ := NewOutputContext(host)
	host.On("Output").Return output

	childExecutionGas := uint64(600)
	_ = metering.DeductInitialGasForExecution(make([]byte, childExecutionGas))
	require.Equal(t, childInput.GasProvided, metering.GasLeft()) // not enough gas provided. It will remain the same for now

	gasRemaining := metering.GasLeft()
	gasSpentByContract := metering.GasSpentByContract()
	require.Equal(t, uint64(0), gasSpentByContract)

	metering.UpdateGasStateOnFailure(output.outputState)

	// return to the parent
	metering.PopSetActiveState()
	mockRuntime.SCAddress = []byte("parent")
	mockRuntime.SetPointsUsed(parentPointsBeforeStacking)
	mockRuntime.SetVMInput(parentInput)

	metering.RestoreGas(gasRemaining)
	require.Equal(t, parentInput.GasProvided-parentExecutionGas-parentPointsBeforeStacking+gasRemaining, metering.GasLeft())

	gasSpentByContract = metering.GasSpentByContract()
	require.Equal(t, parentExecutionGas+parentUsedGas, gasSpentByContract) // child execution will fail due to insufficient gas

	metering.UpdateGasStateOnFailure(output.outputState)

	// after update all gas will be used

	gasUsed := output.outputState.OutputAccounts["parent"].GasUsed
	require.Equal(t, parentInput.GasProvided, gasUsed)

	gasUsed = output.outputState.OutputAccounts["child"].GasUsed
	require.Equal(t, childInput.GasProvided, gasUsed)

	gasRemaining = math.SubUint64(parentInput.GasProvided, gasSpentByContract)
	// calculate gas remaining

	require.Equal(t, int(gasRemaining), int(metering.GasLeft()))
}

func zeroCodeCosts(context *meteringContext) {
	context.GasSchedule().BaseOperationCost.GetCode = 0
}

func TestMeteringContext_TrackGasUsedByBuiltinFunction_GasRemaining(t *testing.T) {
	t.Parallel()

	mockRuntime := &context.MockRuntimeContext{}

	host := &context.MockVMHost{
		RuntimeContext: mockRuntime,
	}

	contractSize := uint64(1000)
	mockRuntime.SCCodeSize = contractSize
	mockRuntime.SCAddress = []byte("parent")

	mockRuntime.SetPointsUsed(0)

	input := &vmcommon.ContractCallInput{
		VMInput:  vmcommon.VMInput{},
		Function: "callBuiltinClaim",
	}
	mockRuntime.SetVMInput(input)

	metering, _ := NewMeteringContext(host, config.MakeGasMapForTests(), uint64(15000))
	host.On("Metering").Return metering
	zeroCodeCosts(metering)

	input.GasProvided = 10000
	metering.InitStateFromContractCallInput(&input.VMInput)
	require.Equal(t, input.GasProvided, metering.GasLeft())

	vmOutput := &vmcommon.VMOutput{
		GasRemaining: 5000,
	}

	postBuiltinInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			GasProvided: 300,
		},
	}

	metering.TrackGasUsedByOutOfVMFunction(input, vmOutput, postBuiltinInput)
	require.Equal(t, vmOutput.GasRemaining+postBuiltinInput.GasProvided, metering.GasLeft())
}

func TestMeteringContext_GasTracer(t *testing.T) {
	t.Parallel()
	const BlockGasLimit = uint64(15000)

	mockRuntime := &context.MockRuntimeContext{
		SCAddress: []byte("scAddress1"),
	}
	host := &context.MockVMHost{
		RuntimeContext: mockRuntime,
	}

	meteringCtx, _ := NewMeteringContext(host, config.MakeGasMapForTests(), BlockGasLimit)
	meteringCtx.InitState()

	gasProvided := uint64(10000)
	meteringCtx.gasForExecution = gasProvided
	gasUsed1 := uint64(1000)
	gasUsed2 := uint64(3000)
	//gasUsed3 := uint64(5000)

	require.NotNil(t, meteringCtx.gasTracer)

	meteringCtx.StartGasTracing("function1")
	gasTrace := meteringCtx.GetGasTrace()
	require.Equal(t, 0, len(gasTrace))
	_ = meteringCtx.UseGasBoundedAndAddTracedGas("function2", gasUsed2)
	gasTrace = meteringCtx.GetGasTrace()
	require.Equal(t, 0, len(gasTrace))

	meteringCtx.SetGasTracing(true)
	meteringCtx.StartGasTracing("function1")
	gasTrace = meteringCtx.GetGasTrace()
	require.Equal(t, 1, len(gasTrace))
	require.Equal(t, 1, len(gasTrace["scAddress1"]))
	require.Equal(t, 1, len(gasTrace["scAddress1"]["function1"]))
	require.Equal(t, uint64(0), gasTrace["scAddress1"]["function1"][0])
	_ = meteringCtx.UseGasBounded(gasUsed1)
	gasTrace = meteringCtx.GetGasTrace()
	require.Equal(t, gasUsed1, gasTrace["scAddress1"]["function1"][0])

	host.RuntimeContext.SetCodeAddress([]byte("scAddress2"))
	_ = meteringCtx.UseGasBoundedAndAddTracedGas("function2", gasUsed2)
	gasTrace = meteringCtx.GetGasTrace()
	require.Equal(t, 2, len(gasTrace))
	require.Equal(t, gasUsed2, gasTrace["scAddress2"]["function2"][0])
}

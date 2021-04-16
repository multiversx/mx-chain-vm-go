package host

import (
	"math/big"
	"testing"

	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/stretchr/testify/require"
)

var gasUsedByParent = uint64(400)
var gasUsedByChild = uint64(200)
var gasProvidedToChild = uint64(300)
var gasUsedByBuiltinClaim = uint64(150)

func TestGasUsed_SingleContract(t *testing.T) {
	host, _, ibm := defaultTestArwenForCallWithInstanceMocks(t)
	createTestParentContract(t, host, ibm)
	zeroCodeCosts(host)

	gasProvided := uint64(1000)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.GasProvided = gasProvided
	input.Function = "wasteGas"

	vmOutput, err := host.RunSmartContractCall(input)

	verify := NewVMOutputVerifier(t, vmOutput, err)
	verify.GasUsed(parentAddress, gasUsedByParent)
	verify.GasRemaining(gasProvided - gasUsedByParent)
}

func TestGasUsed_SingleContract_BuiltinCall(t *testing.T) {
	logger.SetLogLevel("*:TRACE")
	logger.ToggleLoggerName(true)

	host, world, ibm := defaultTestArwenForCallWithInstanceMocks(t)
	createTestParentContract(t, host, ibm)
	createMockBuiltinFunctions(t, host, world)
	zeroCodeCosts(host)

	gasProvided := uint64(1000)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.GasProvided = gasProvided
	input.Function = "execBuiltinDestCtx"
	input.Arguments = [][]byte{
		big.NewInt(0).SetUint64(1).Bytes(),
		[]byte("builtinClaim"),
	}

	vmOutput, err := host.RunSmartContractCall(input)

	verify := NewVMOutputVerifier(t, vmOutput, err)
	verify.GasUsed(parentAddress, gasUsedByParent+gasUsedByBuiltinClaim)
	verify.GasRemaining(gasProvided - gasUsedByParent - gasUsedByBuiltinClaim)
}

func TestGasUsed_TwoContracts_ExecuteOnSameCtx(t *testing.T) {
	host, _, ibm := defaultTestArwenForCallWithInstanceMocks(t)
	createTestParentContract(t, host, ibm)
	createTestChildContract(t, host, ibm)
	zeroCodeCosts(host)

	gasProvided := uint64(1000)
	for numChildCalls := uint64(0); numChildCalls < 3; numChildCalls++ {
		expectedGasRemaining := gasProvided - gasUsedByParent - gasUsedByChild*numChildCalls

		input := DefaultTestContractCallInput()
		input.RecipientAddr = parentAddress
		input.GasProvided = gasProvided
		input.Function = "execChildSameCtx"
		input.Arguments = [][]byte{
			big.NewInt(0).SetUint64(numChildCalls).Bytes(),
			[]byte("wasteGas"),
		}

		vmOutput, err := host.RunSmartContractCall(input)

		verify := NewVMOutputVerifier(t, vmOutput, err)
		verify.GasRemaining(expectedGasRemaining)
		verify.GasUsed(parentAddress, gasUsedByParent)
		if numChildCalls > 0 {
			verify.GasUsed(childAddress, gasUsedByChild*numChildCalls)
		}
	}
}

func TestGasUsed_TwoContracts_ExecuteOnDestCtx(t *testing.T) {
	host, _, ibm := defaultTestArwenForCallWithInstanceMocks(t)
	createTestParentContract(t, host, ibm)
	createTestChildContract(t, host, ibm)
	zeroCodeCosts(host)

	gasProvided := uint64(1000)

	for numChildCalls := uint64(0); numChildCalls < 3; numChildCalls++ {
		expectedGasRemaining := gasProvided - gasUsedByParent - gasUsedByChild*numChildCalls

		input := DefaultTestContractCallInput()
		input.RecipientAddr = parentAddress
		input.GasProvided = gasProvided
		input.Function = "execChildDestCtx"
		input.Arguments = [][]byte{
			big.NewInt(0).SetUint64(numChildCalls).Bytes(),
			[]byte("wasteGas"),
		}

		vmOutput, err := host.RunSmartContractCall(input)

		verify := NewVMOutputVerifier(t, vmOutput, err)
		verify.GasRemaining(expectedGasRemaining)
		verify.GasUsed(parentAddress, gasUsedByParent)
		if numChildCalls > 0 {
			verify.GasUsed(childAddress, gasUsedByChild*numChildCalls)
		}
	}
}

func createTestParentContract(t testing.TB, host *vmHost, ibm *mock.InstanceBuilderMock) {
	log := logger.GetOrCreate("arwen/testParent")

	gasUsedByParent := uint64(400)
	gasProvidedToChild := uint64(300)

	parentInstance := ibm.CreateAndStoreInstanceMock(parentAddress, 0)

	parentInstance.AddMockMethod("wasteGas", func() {
		host.Metering().UseGas(gasUsedByParent)
	})

	parentInstance.AddMockMethod("execChildSameCtx", func() {
		host.Metering().UseGas(gasUsedByParent)

		arguments := host.Runtime().Arguments()
		numChildCalls := big.NewInt(0).SetBytes(arguments[0]).Uint64()
		childFunction := string(arguments[1])

		childInput := DefaultTestContractCallInput()
		childInput.CallerAddr = parentAddress
		childInput.RecipientAddr = childAddress
		childInput.Function = childFunction
		childInput.GasProvided = gasProvidedToChild

		for i := uint64(0); i < numChildCalls; i++ {
			log.Trace("ExecuteOnSameContext child call", "index", i)
			_, err := host.ExecuteOnSameContext(childInput)
			require.Nil(t, err)
		}
	})

	parentInstance.AddMockMethod("execChildDestCtx", func() {
		host.Metering().UseGas(gasUsedByParent)

		arguments := host.Runtime().Arguments()
		numChildCalls := big.NewInt(0).SetBytes(arguments[0]).Uint64()
		childFunction := string(arguments[1])

		childInput := DefaultTestContractCallInput()
		childInput.CallerAddr = parentAddress
		childInput.RecipientAddr = childAddress
		childInput.Function = childFunction
		childInput.GasProvided = gasProvidedToChild

		for i := uint64(0); i < numChildCalls; i++ {
			log.Trace("ExecuteOnDestContext child call", "index", i)
			_, _, _, err := host.ExecuteOnDestContext(childInput)
			require.Nil(t, err)
		}
	})

	parentInstance.AddMockMethod("execBuiltinDestCtx", func() {
		host.Metering().UseGas(gasUsedByParent)

		arguments := host.Runtime().Arguments()
		numBuiltinCalls := big.NewInt(0).SetBytes(arguments[0]).Uint64()
		builtinFunction := string(arguments[1])

		builtinInput := DefaultTestContractCallInput()
		builtinInput.CallerAddr = parentAddress
		builtinInput.RecipientAddr = parentAddress
		builtinInput.Function = builtinFunction
		builtinInput.GasProvided = gasProvidedToChild

		for i := uint64(0); i < numBuiltinCalls; i++ {
			log.Trace("ExecuteOnDestContext builtin call", "index", i)
			_, _, _, err := host.ExecuteOnDestContext(builtinInput)
			require.Nil(t, err)
		}
	})
}

func createTestChildContract(tb testing.TB, host *vmHost, ibm *mock.InstanceBuilderMock) {
	gasUsedByChild := uint64(200)

	childInstance := ibm.CreateAndStoreInstanceMock(childAddress, 0)
	childInstance.AddMockMethod("wasteGas", func() {
		host.Metering().UseGas(gasUsedByChild)
	})
}

func createMockBuiltinFunctions(tb testing.TB, host *vmHost, world *worldmock.MockWorld) {
	err := world.InitBuiltinFunctions(host.GetGasScheduleMap())
	require.Nil(tb, err)

	world.BuiltinFuncs.Container.Add("builtinClaim", &MockClaimBuiltin{
		AmountToGive: big.NewInt(42),
		GasCost:      gasUsedByBuiltinClaim,
	})

	host.protocolBuiltinFunctions = world.BuiltinFuncs.GetBuiltinFunctionNames()
}

func zeroCodeCosts(host *vmHost) {
	host.Metering().GasSchedule().BaseOperationCost.CompilePerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.AoTPreparePerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.GetCode = 0
}

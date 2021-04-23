package host

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	"github.com/stretchr/testify/require"
)

var gasUsedByParent = uint64(400)
var gasUsedByChild = uint64(200)
var gasProvidedToChild = uint64(300)
var gasUsedByBuiltinClaim = uint64(120)

func TestGasUsed_SingleContract(t *testing.T) {
	host, _, imb := defaultTestArwenForCallWithInstanceMocks(t)
	createTestParentContract(t, host, imb)
	zeroCodeCosts(host)

	gasProvided := uint64(1000)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.GasProvided = gasProvided
	input.Function = "wasteGas"

	vmOutput, err := host.RunSmartContractCall(input)

	verify := NewVMOutputVerifier(t, vmOutput, err)
	verify.Ok().GasRemaining(gasProvided - gasUsedByParent)
	verify.GasUsed(parentAddress, gasUsedByParent)
}

func TestGasUsed_SingleContract_BuiltinCall(t *testing.T) {
	host, world, imb := defaultTestArwenForCallWithInstanceMocks(t)
	createTestParentContract(t, host, imb)
	createMockBuiltinFunctions(t, host, world)
	zeroCodeCosts(host)

	gasProvided := uint64(1000)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.GasProvided = gasProvided
	input.Function = "execOnDestCtx"
	input.Arguments = [][]byte{
		parentAddress,
		[]byte("builtinClaim"),
		arwen.One.Bytes(),
	}

	vmOutput, err := host.RunSmartContractCall(input)

	verify := NewVMOutputVerifier(t, vmOutput, err)
	verify.Ok().GasRemaining(gasProvided - gasUsedByParent - gasUsedByBuiltinClaim)
	verify.GasUsed(parentAddress, gasUsedByParent+gasUsedByBuiltinClaim)
}

func TestGasUsed_TwoContracts_ExecuteOnSameCtx(t *testing.T) {
	host, _, imb := defaultTestArwenForCallWithInstanceMocks(t)
	createTestParentContract(t, host, imb)
	createTestChildContract(t, host, imb)
	zeroCodeCosts(host)

	gasProvided := uint64(1000)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.GasProvided = gasProvided
	input.Function = "execOnSameCtx"

	for numCalls := uint64(0); numCalls < 3; numCalls++ {
		expectedGasRemaining := gasProvided - gasUsedByParent - gasUsedByChild*numCalls

		numCallsBytes := big.NewInt(0).SetUint64(numCalls).Bytes()
		input.Arguments = [][]byte{
			childAddress,
			[]byte("wasteGas"),
			numCallsBytes,
		}

		vmOutput, err := host.RunSmartContractCall(input)

		verify := NewVMOutputVerifier(t, vmOutput, err)
		verify.Ok().GasRemaining(expectedGasRemaining)
		verify.GasUsed(parentAddress, gasUsedByParent)
		if numCalls > 0 {
			verify.GasUsed(childAddress, gasUsedByChild*numCalls)
		}
	}
}

func TestGasUsed_TwoContracts_ExecuteOnDestCtx(t *testing.T) {
	host, _, imb := defaultTestArwenForCallWithInstanceMocks(t)
	createTestParentContract(t, host, imb)
	createTestChildContract(t, host, imb)
	zeroCodeCosts(host)

	gasProvided := uint64(1000)
	input := DefaultTestContractCallInput()
	input.RecipientAddr = parentAddress
	input.GasProvided = gasProvided
	input.Function = "execOnDestCtx"

	for numCalls := uint64(0); numCalls < 3; numCalls++ {
		expectedGasRemaining := gasProvided - gasUsedByParent - gasUsedByChild*numCalls

		numCallsBytes := big.NewInt(0).SetUint64(numCalls).Bytes()
		input.Arguments = [][]byte{
			childAddress,
			[]byte("wasteGas"),
			numCallsBytes,
		}

		vmOutput, err := host.RunSmartContractCall(input)

		verify := NewVMOutputVerifier(t, vmOutput, err)
		verify.Ok().GasRemaining(expectedGasRemaining)
		verify.GasUsed(parentAddress, gasUsedByParent)
		if numCalls > 0 {
			verify.GasUsed(childAddress, gasUsedByChild*numCalls)
		}
	}
}

func TestGasUsed_ThreeContracts_ExecuteOnDestCtx(t *testing.T) {
	host, _, imb := defaultTestArwenForCallWithInstanceMocks(t)
	zeroCodeCosts(host)

	alphaAddress := MakeTestSCAddress("alpha")
	betaAddress := MakeTestSCAddress("beta")
	gammaAddress := MakeTestSCAddress("gamma")

	gasProvided := uint64(1000)
	alphaCallGas := uint64(400)
	alphaGasToForwardToReceivers := uint64(300)
	receiverCallGas := uint64(200)

	expectedGasRemaining := gasProvided - alphaCallGas - 2*receiverCallGas

	alpha := imb.CreateAndStoreInstanceMock(t, host, alphaAddress, 0)
	addForwarderMethodsToInstanceMock(alpha, alphaCallGas, alphaGasToForwardToReceivers)

	beta := imb.CreateAndStoreInstanceMock(t, host, betaAddress, 0)
	gamma := imb.CreateAndStoreInstanceMock(t, host, gammaAddress, 0)
	addDummyMethodsToInstanceMock(beta, receiverCallGas)
	addDummyMethodsToInstanceMock(gamma, receiverCallGas)

	input := DefaultTestContractCallInput()
	input.RecipientAddr = alphaAddress
	input.GasProvided = gasProvided
	input.Function = "execOnDestCtx"
	input.Arguments = [][]byte{
		betaAddress,
		[]byte("wasteGas"),
		arwen.One.Bytes(),
		gammaAddress,
		[]byte("wasteGas"),
		arwen.One.Bytes(),
	}

	vmOutput, err := host.RunSmartContractCall(input)

	verify := NewVMOutputVerifier(t, vmOutput, err)
	verify.GasUsed(alphaAddress, alphaCallGas)
	verify.GasUsed(betaAddress, receiverCallGas)
	verify.GasUsed(gammaAddress, receiverCallGas)
	verify.Ok().GasRemaining(expectedGasRemaining)
}

func createTestParentContract(t testing.TB, host *vmHost, imb *mock.InstanceBuilderMock) {
	gasUsedByParent := uint64(400)
	gasProvidedToChild := uint64(300)

	parentInstance := imb.CreateAndStoreInstanceMock(t, host, parentAddress, 0)
	addDummyMethodsToInstanceMock(parentInstance, gasUsedByParent)
	addForwarderMethodsToInstanceMock(parentInstance, gasUsedByParent, gasProvidedToChild)
}

func createTestChildContract(t testing.TB, host *vmHost, imb *mock.InstanceBuilderMock) {
	gasUsedByChild := uint64(200)

	childInstance := imb.CreateAndStoreInstanceMock(t, host, childAddress, 0)
	addDummyMethodsToInstanceMock(childInstance, gasUsedByChild)
}

func createMockBuiltinFunctions(tb testing.TB, host *vmHost, world *worldmock.MockWorld) {
	err := world.InitBuiltinFunctions(host.GetGasScheduleMap())
	require.Nil(tb, err)

	world.BuiltinFuncs.Container.Add("builtinClaim", &MockClaimBuiltin{
		AmountToGive: 42,
		GasCost:      gasUsedByBuiltinClaim,
	})

	host.protocolBuiltinFunctions = world.BuiltinFuncs.GetBuiltinFunctionNames()
}

func addDummyMethodsToInstanceMock(instance *mock.InstanceMock, gasPerCall uint64) {
	instance.AddMockMethod("wasteGas", func() {
		host := instance.Host
		host.Metering().UseGas(gasPerCall)
	})
}

func addForwarderMethodsToInstanceMock(instance *mock.InstanceMock, gasPerCall uint64, gasToForward uint64) {
	input := DefaultTestContractCallInput()
	input.GasProvided = gasToForward

	t := instance.T

	instance.AddMockMethod("execOnSameCtx", func() {
		host := instance.Host
		host.Metering().UseGas(gasPerCall)

		arguments := host.Runtime().Arguments()
		input.CallerAddr = instance.Address
		input.RecipientAddr = arguments[0]
		input.Function = string(arguments[1])
		numCalls := big.NewInt(0).SetBytes(arguments[2]).Uint64()

		for i := uint64(0); i < numCalls; i++ {
			_, err := host.ExecuteOnSameContext(input)
			require.Nil(t, err)
		}
	})

	instance.AddMockMethod("execOnDestCtx", func() {
		host := instance.Host
		host.Metering().UseGas(gasPerCall)

		argsPerCall := 3
		arguments := host.Runtime().Arguments()
		if len(arguments)%argsPerCall != 0 {
			host.Runtime().SignalUserError("need 3 arguments per individual call")
			return
		}

		input.CallerAddr = instance.Address

		for callIndex := 0; callIndex < len(arguments); callIndex += argsPerCall {
			input.RecipientAddr = arguments[callIndex+0]
			input.Function = string(arguments[callIndex+1])
			numCalls := big.NewInt(0).SetBytes(arguments[callIndex+2]).Uint64()

			for i := uint64(0); i < numCalls; i++ {
				_, _, _, err := host.ExecuteOnDestContext(input)
				require.Nil(t, err)
			}
		}
	})
}

func zeroCodeCosts(host *vmHost) {
	host.Metering().GasSchedule().BaseOperationCost.CompilePerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.AoTPreparePerByte = 0
	host.Metering().GasSchedule().BaseOperationCost.GetCode = 0
}

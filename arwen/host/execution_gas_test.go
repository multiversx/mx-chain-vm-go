package host

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGasUsed_SingleContract(t *testing.T) {
	host, _, ibm := defaultTestArwenForCallWithInstanceMocks(t)

	gasProvided := uint64(1000)
	gasUsedByParentExec := uint64(400)
	contractCompilationCost := uint64(32)

	parentInstance := ibm.CreateAndStoreInstanceMock(parentAddress, 0)
	parentInstance.AddMockMethod("doSomething", func() {
		host.Metering().UseGas(gasUsedByParentExec)
	})

	input := DefaultTestContractCallInput()
	input.Function = "doSomething"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	verify := NewVMOutputVerifier(t, vmOutput, err)

	gasUsedByParent := contractCompilationCost + gasUsedByParentExec + 1
	verify.GasUsed(parentAddress, gasUsedByParent)
	verify.GasRemaining(gasProvided - gasUsedByParent)
}

func TestGasUsed_TwoContracts_ExecuteOnSameCtx(t *testing.T) {
	host, _, ibm := defaultTestArwenForCallWithInstanceMocks(t)

	gasProvided := uint64(1000)
	gasUsedByParentExec := uint64(400)
	gasUsedByChildExec := uint64(200)
	contractCompilationCost := uint64(32)

	parentInstance := ibm.CreateAndStoreInstanceMock(parentAddress, 0)
	parentInstance.AddMockMethod("parentFunction", func() {
		host.Metering().UseGas(gasUsedByParentExec)
		childInput := DefaultTestContractCallInput()
		childInput.CallerAddr = parentAddress
		childInput.RecipientAddr = childAddress
		childInput.GasProvided = 300
		childInput.Function = "childFunction"
		_, err := host.ExecuteOnSameContext(childInput)
		require.Nil(t, err)
	})

	childInstance := ibm.CreateAndStoreInstanceMock(childAddress, 0)
	childInstance.AddMockMethod("childFunction", func() {
		host.Metering().UseGas(gasUsedByChildExec)
	})

	input := DefaultTestContractCallInput()
	input.Function = "parentFunction"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	verify := NewVMOutputVerifier(t, vmOutput, err)

	gasUsedByParent := contractCompilationCost + gasUsedByParentExec + 1
	verify.GasUsed(parentAddress, gasUsedByParent)

	gasUsedByChild := contractCompilationCost + gasUsedByChildExec + 1
	verify.GasUsed(childAddress, gasUsedByChild)

	expectedGasRemaining := gasProvided - gasUsedByParent - gasUsedByChild
	verify.GasRemaining(expectedGasRemaining)
}

func TestGasUsed_TwoContracts_ExecuteOnDestCtx(t *testing.T) {
	host, _, ibm := defaultTestArwenForCallWithInstanceMocks(t)

	gasProvided := uint64(1000)
	gasUsedByParentExec := uint64(400)
	gasUsedByChildExec := uint64(200)
	contractCompilationCost := uint64(32)

	parentInstance := ibm.CreateAndStoreInstanceMock(parentAddress, 0)
	parentInstance.AddMockMethod("parentFunction", func() {
		host.Metering().UseGas(gasUsedByParentExec)
		childInput := DefaultTestContractCallInput()
		childInput.CallerAddr = parentAddress
		childInput.RecipientAddr = childAddress
		childInput.GasProvided = 300
		childInput.Function = "childFunction"
		_, _, _, err := host.ExecuteOnDestContext(childInput)
		require.Nil(t, err)
	})

	childInstance := ibm.CreateAndStoreInstanceMock(childAddress, 0)
	childInstance.AddMockMethod("childFunction", func() {
		host.Metering().UseGas(gasUsedByChildExec)
	})

	input := DefaultTestContractCallInput()
	input.Function = "parentFunction"
	input.GasProvided = gasProvided

	vmOutput, err := host.RunSmartContractCall(input)
	verify := NewVMOutputVerifier(t, vmOutput, err)

	gasUsedByParent := contractCompilationCost + gasUsedByParentExec + 1
	verify.GasUsed(parentAddress, gasUsedByParent)

	gasUsedByChild := contractCompilationCost + gasUsedByChildExec + 1
	verify.GasUsed(childAddress, gasUsedByChild)

	expectedGasRemaining := gasProvided - gasUsedByParent - gasUsedByChild
	verify.GasRemaining(expectedGasRemaining)
}

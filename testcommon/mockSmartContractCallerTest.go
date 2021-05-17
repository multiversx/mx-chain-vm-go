package testcommon

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

type testTemplateConfig struct {
	t        *testing.T
	input    *vmcommon.ContractCallInput
	useMocks bool
}

// MockInstancesTestTemplate holds the data to build a mock contract call test
type MockInstancesTestTemplate struct {
	testTemplateConfig
	contracts     *[]MockTestSmartContract
	setup         func(arwen.VMHost, *worldmock.MockWorld)
	assertResults func(*worldmock.MockWorld, *VMOutputVerifier)
}

// BuildMockInstanceCallTest starts the building process for a mock contract call test
func BuildMockInstanceCallTest(t *testing.T) *MockInstancesTestTemplate {
	return &MockInstancesTestTemplate{
		testTemplateConfig: testTemplateConfig{
			t:        t,
			useMocks: true,
		},
		setup: func(arwen.VMHost, *worldmock.MockWorld) {},
	}
}

// WithContracts provides the contracts to be used by the mock contract call test
func (callerTest *MockInstancesTestTemplate) WithContracts(usedContracts ...MockTestSmartContract) *MockInstancesTestTemplate {
	callerTest.contracts = &usedContracts
	return callerTest
}

// WithInput provides the ContractCallInput to be used by the mock contract call test
func (callerTest *MockInstancesTestTemplate) WithInput(input *vmcommon.ContractCallInput) *MockInstancesTestTemplate {
	callerTest.input = input
	return callerTest
}

// WithSetup provides the setup function to be used by the mock contract call test
func (callerTest *MockInstancesTestTemplate) WithSetup(setup func(arwen.VMHost, *worldmock.MockWorld)) *MockInstancesTestTemplate {
	callerTest.setup = setup
	return callerTest
}

// AndAssertResults provides the function that will aserts the results
func (callerTest *MockInstancesTestTemplate) AndAssertResults(assertResults func(world *worldmock.MockWorld, verify *VMOutputVerifier)) {
	callerTest.assertResults = assertResults
	callerTest.runTest()
}

func (callerTest *MockInstancesTestTemplate) runTest() {

	host, world, imb := DefaultTestArwenForCallWithInstanceMocks(callerTest.t)

	for _, mockSC := range *callerTest.contracts {
		mockSC.initialize(callerTest.t, host, imb)
	}

	callerTest.setup(host, world)

	vmOutput, err := host.RunSmartContractCall(callerTest.input)

	verify := NewVMOutputVerifier(callerTest.t, vmOutput, err)
	callerTest.assertResults(world, verify)
}

// SimpleWasteGasMockMethod is a simple waste gas mock method
func SimpleWasteGasMockMethod(instanceMock *mock.InstanceMock, gas uint64) func() *mock.InstanceMock {
	return func() *mock.InstanceMock {
		host := instanceMock.Host
		host.Metering().UseGas(gas)
		instance := mock.GetMockInstance(host)
		return instance
	}
}

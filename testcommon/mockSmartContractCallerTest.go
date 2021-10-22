package testcommon

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/world"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var logMock = logger.GetOrCreate("arwen/mock")

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
	assertResults func(*TestCallNode, *worldmock.MockWorld, *VMOutputVerifier)
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
func (callerTest *MockInstancesTestTemplate) AndAssertResults(assertResults func(world *worldmock.MockWorld, verify *VMOutputVerifier)) (*vmcommon.VMOutput, error) {
	return callerTest.AndAssertResultsWithWorld(nil, true, nil, func(startNode *TestCallNode, world *worldmock.MockWorld, verify *VMOutputVerifier) {
		assertResults(world, verify)
	})
}

// AndAssertResultsWithWorld provides the function that will aserts the results
func (callerTest *MockInstancesTestTemplate) AndAssertResultsWithWorld(world *worldmock.MockWorld, createContractAccounts bool, startNode *TestCallNode, assertResults func(startNode *TestCallNode, world *worldmock.MockWorld, verify *VMOutputVerifier)) (*vmcommon.VMOutput, error) {
	callerTest.assertResults = assertResults
	if world == nil {
		world = worldmock.NewMockWorld()
	}
	return callerTest.runTest(startNode, world, createContractAccounts)
}

func (callerTest *MockInstancesTestTemplate) runTest(startNode *TestCallNode, world *worldmock.MockWorld, createContractAccounts bool) (*vmcommon.VMOutput, error) {
	host, imb := DefaultTestArwenForCallWithInstanceMocksAndWorld(callerTest.t, world)

	for _, mockSC := range *callerTest.contracts {
		mockSC.initialize(callerTest.t, host, imb, createContractAccounts)
	}

	callerTest.setup(host, world)
	// create snapshot (normaly done by node)
	world.CreateStateBackup()

	vmOutput, err := host.RunSmartContractCall(callerTest.input)

	allErrors := host.Runtime().GetAllErrors()
	verify := NewVMOutputVerifierWithAllErrors(callerTest.t, vmOutput, err, allErrors)
	if callerTest.assertResults != nil {
		callerTest.assertResults(startNode, world, verify)
	}

	return vmOutput, err
}

// SimpleWasteGasMockMethod is a simple waste gas mock method
func SimpleWasteGasMockMethod(instanceMock *mock.InstanceMock, gas uint64) func() *mock.InstanceMock {
	return func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		err := host.Metering().UseGasBounded(gas)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
		}

		return instance
	}
}

// WasteGasWithReturnDataMockMethod is a simple waste gas mock method
func WasteGasWithReturnDataMockMethod(instanceMock *mock.InstanceMock, gas uint64, returnData []byte) func() *mock.InstanceMock {
	return func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		logMock.Trace("instance mock waste gas", "sc", string(host.Runtime().GetSCAddress()), "func", host.Runtime().Function(), "gas", gas)
		err := host.Metering().UseGasBounded(gas)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
			return instance
		}

		host.Output().Finish(returnData)
		return instance
	}
}

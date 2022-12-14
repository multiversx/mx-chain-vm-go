package testcommon

import (
	"testing"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/arwen"
	mock "github.com/ElrondNetwork/wasm-vm/mock/context"
	worldmock "github.com/ElrondNetwork/wasm-vm/mock/world"
)

var logMock = logger.GetOrCreate("arwen/mock")

// SetupFunction -
type SetupFunction func(arwen.VMHost, *worldmock.MockWorld)

// AssertResultsFunc -
type AssertResultsFunc func(world *worldmock.MockWorld, verify *VMOutputVerifier)

// AssertResultsWithStartNodeFunc -
type AssertResultsWithStartNodeFunc func(startNode *TestCallNode, world *worldmock.MockWorld, verify *VMOutputVerifier, expectedErrorsForRound []string)

type testTemplateConfig struct {
	tb                       testing.TB
	input                    *vmcommon.ContractCallInput
	useMocks                 bool
	wasmerSIGSEGVPassthrough bool
}

// MockInstancesTestTemplate holds the data to build a mock contract call test
type MockInstancesTestTemplate struct {
	testTemplateConfig
	contracts     *[]MockTestSmartContract
	setup         SetupFunction
	assertResults func(*TestCallNode, *worldmock.MockWorld, *VMOutputVerifier, []string)
}

// BuildMockInstanceCallTest starts the building process for a mock contract call test
func BuildMockInstanceCallTest(tb testing.TB) *MockInstancesTestTemplate {
	return &MockInstancesTestTemplate{
		testTemplateConfig: testTemplateConfig{
			tb:                       tb,
			useMocks:                 true,
			wasmerSIGSEGVPassthrough: false,
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
func (callerTest *MockInstancesTestTemplate) WithSetup(setup SetupFunction) *MockInstancesTestTemplate {
	callerTest.setup = setup
	return callerTest
}

// WithWasmerSIGSEGVPassthrough sets the wasmerSIGSEGVPassthrough flag
func (callerTest *MockInstancesTestTemplate) WithWasmerSIGSEGVPassthrough(wasmerSIGSEGVPassthrough bool) *MockInstancesTestTemplate {
	callerTest.wasmerSIGSEGVPassthrough = wasmerSIGSEGVPassthrough
	return callerTest
}

// AndAssertResults provides the function that will aserts the results
func (callerTest *MockInstancesTestTemplate) AndAssertResults(assertResults AssertResultsFunc) (*vmcommon.VMOutput, error) {
	return callerTest.AndAssertResultsWithWorld(nil, true, nil, nil, func(startNode *TestCallNode, world *worldmock.MockWorld, verify *VMOutputVerifier, expectedErrorsForRound []string) {
		assertResults(world, verify)
	})
}

// AndAssertResultsWithWorld provides the function that will aserts the results
func (callerTest *MockInstancesTestTemplate) AndAssertResultsWithWorld(
	world *worldmock.MockWorld,
	createAccount bool,
	startNode *TestCallNode,
	expectedErrorsForRound []string,
	assertResults AssertResultsWithStartNodeFunc) (*vmcommon.VMOutput, error) {
	callerTest.assertResults = assertResults
	if world == nil {
		world = worldmock.NewMockWorld()
	}
	return callerTest.runTest(startNode, world, createAccount, expectedErrorsForRound)
}

func (callerTest *MockInstancesTestTemplate) runTest(startNode *TestCallNode, world *worldmock.MockWorld, createAccount bool, expectedErrorsForRound []string) (*vmcommon.VMOutput, error) {
	if world == nil {
		world = worldmock.NewMockWorld()
	}
	executorFactory := mock.NewExecutorMockFactory(world)
	host := NewTestHostBuilder(callerTest.tb).
		WithExecutorFactory(executorFactory).
		WithBlockchainHook(world).
		Build()

	defer func() {
		host.Reset()
	}()

	for _, mockSC := range *callerTest.contracts {
		mockSC.Initialize(callerTest.tb, host, executorFactory.LastCreatedExecutor, createAccount)
	}

	callerTest.setup(host, world)
	// create snapshot (normaly done by node)
	world.CreateStateBackup()

	vmOutput, err := host.RunSmartContractCall(callerTest.input)
	allErrors := host.Runtime().GetAllErrors()
	verify := NewVMOutputVerifierWithAllErrors(callerTest.tb, vmOutput, err, allErrors)
	if callerTest.assertResults != nil {
		callerTest.assertResults(startNode, world, verify, expectedErrorsForRound)
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

		logMock.Trace("instance mock waste gas", "sc", string(host.Runtime().GetContextAddress()), "func", host.Runtime().FunctionName(), "gas", gas)
		err := host.Metering().UseGasBounded(gas)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
			return instance
		}

		host.Output().Finish(returnData)
		return instance
	}
}

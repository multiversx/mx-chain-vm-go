package testcommon

import (
	"testing"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/executor"
	executorwrapper "github.com/multiversx/mx-chain-vm-go/executor/wrapper"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/require"
)

// InstanceTestSmartContract represents the config data for the smart contract instance to be tested
type InstanceTestSmartContract struct {
	testSmartContract
	code []byte
}

// CreateInstanceContract build a contract to be used in a test creted with BuildInstanceCallTest
func CreateInstanceContract(address []byte) *InstanceTestSmartContract {
	return &InstanceTestSmartContract{
		testSmartContract: testSmartContract{
			address: address,
		},
	}
}

// WithBalance provides the balance for the InstanceTestSmartContract
func (mockSC *InstanceTestSmartContract) WithBalance(balance int64) *InstanceTestSmartContract {
	mockSC.balance = balance
	return mockSC
}

// WithConfig provides the config object for the InstanceTestSmartContract
func (mockSC *InstanceTestSmartContract) WithConfig(testConfig *TestConfig) *InstanceTestSmartContract {
	mockSC.config = testConfig
	return mockSC
}

// WithCode provides the code for the InstanceTestSmartContract
func (mockSC *InstanceTestSmartContract) WithCode(code []byte) *InstanceTestSmartContract {
	mockSC.code = code
	return mockSC
}

// InstancesTestTemplate holds the data to build a contract call test
type InstancesTestTemplate struct {
	testTemplateConfig
	contracts     []*InstanceTestSmartContract
	setup         func(vmhost.VMHost, *contextmock.BlockchainHookStub)
	assertResults func(vmhost.VMHost, *contextmock.BlockchainHookStub, *VMOutputVerifier)
	host          vmhost.VMHost
	hostBuilder   *TestHostBuilder
}

// BuildInstanceCallTest starts the building process for a contract call test
func BuildInstanceCallTest(tb testing.TB) *InstancesTestTemplate {
	return &InstancesTestTemplate{
		testTemplateConfig: testTemplateConfig{
			tb:                       tb,
			useMocks:                 false,
			wasmerSIGSEGVPassthrough: false,
		},
		hostBuilder: NewTestHostBuilder(tb),
		setup:       func(vmhost.VMHost, *contextmock.BlockchainHookStub) {},
	}
}

// WithContracts provides the contracts to be used by the contract call test
func (callerTest *InstancesTestTemplate) WithContracts(usedContracts ...*InstanceTestSmartContract) *InstancesTestTemplate {
	callerTest.contracts = usedContracts
	return callerTest
}

// WithInput provides the ContractCallInput to be used by the contract call test
func (callerTest *InstancesTestTemplate) WithInput(input *vmcommon.ContractCallInput) *InstancesTestTemplate {
	callerTest.input = input
	return callerTest
}

// WithSetup provides the setup function to be used by the contract call test
func (callerTest *InstancesTestTemplate) WithSetup(setup func(vmhost.VMHost, *contextmock.BlockchainHookStub)) *InstancesTestTemplate {
	callerTest.setup = setup
	return callerTest
}

// WithGasSchedule provides gas schedule for the test
func (callerTest *InstancesTestTemplate) WithGasSchedule(gasSchedule config.GasScheduleMap) *InstancesTestTemplate {
	callerTest.hostBuilder.WithGasSchedule(gasSchedule)
	return callerTest
}

// WithExecutorFactory provides the wasmer executor for the test
func (callerTest *InstancesTestTemplate) WithExecutorFactory(executorFactory executor.ExecutorAbstractFactory) *InstancesTestTemplate {
	callerTest.hostBuilder.WithExecutorFactory(executorFactory)
	return callerTest
}

// WithExecutorLogs sets an ExecutorLogger
func (callerTest *InstancesTestTemplate) WithExecutorLogs(executorLogger executorwrapper.ExecutorLogger) *InstancesTestTemplate {
	callerTest.hostBuilder.WithExecutorLogs(executorLogger)
	return callerTest
}

// WithWasmerSIGSEGVPassthrough sets the wasmerSIGSEGVPassthrough flag
func (callerTest *InstancesTestTemplate) WithWasmerSIGSEGVPassthrough(passthrough bool) *InstancesTestTemplate {
	callerTest.hostBuilder.WithWasmerSIGSEGVPassthrough(passthrough)
	return callerTest
}

// GetVMHost returns the inner VMHost
func (callerTest *InstancesTestTemplate) GetVMHost() vmhost.VMHost {
	return callerTest.host
}

// AndAssertResults starts the test and asserts the results
func (callerTest *InstancesTestTemplate) AndAssertResults(assertResults func(vmhost.VMHost, *contextmock.BlockchainHookStub, *VMOutputVerifier)) {
	callerTest.assertResults = assertResults
	runTestWithInstances(callerTest, true)
}

// AndAssertResultsWithoutReset starts the test and asserts the results
func (callerTest *InstancesTestTemplate) AndAssertResultsWithoutReset(assertResults func(vmhost.VMHost, *contextmock.BlockchainHookStub, *VMOutputVerifier)) {
	callerTest.assertResults = assertResults
	runTestWithInstances(callerTest, false)
}

func runTestWithInstances(callerTest *InstancesTestTemplate, reset bool) {
	var blhookStub *contextmock.BlockchainHookStub
	if callerTest.host == nil {
		blhookStub = BlockchainHookStubForContracts(callerTest.contracts)
		callerTest.hostBuilder.WithBlockchainHook(blhookStub)
		callerTest.host = callerTest.hostBuilder.Build()
		callerTest.setup(callerTest.host, blhookStub)
	}

	defer func() {
		if reset {
			callerTest.host.Reset()
		}

		// Extra verification for instance leaks
		err := callerTest.host.Runtime().ValidateInstances()
		require.Nil(callerTest.tb, err)
	}()

	vmOutput, err := callerTest.host.RunSmartContractCall(callerTest.input)

	if callerTest.assertResults != nil {
		allErrors := callerTest.host.Runtime().GetAllErrors()
		verify := NewVMOutputVerifierWithAllErrors(callerTest.tb, vmOutput, err, allErrors)
		callerTest.assertResults(callerTest.host, blhookStub, verify)
	}
}

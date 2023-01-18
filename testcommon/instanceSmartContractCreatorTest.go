package testcommon

import (
	"testing"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/wasm-vm/arwen"
	"github.com/multiversx/wasm-vm/config"
	"github.com/multiversx/wasm-vm/executor"
	contextmock "github.com/multiversx/wasm-vm/mock/context"
)

// TestCreateTemplateConfig holds the data to build a contract creation test
type TestCreateTemplateConfig struct {
	t                        *testing.T
	address                  []byte
	input                    *vmcommon.ContractCreateInput
	setup                    func(arwen.VMHost, *contextmock.BlockchainHookStub)
	assertResults            func(*contextmock.BlockchainHookStub, *VMOutputVerifier)
	host                     arwen.VMHost
	gasSchedule              config.GasScheduleMap
	wasmerSIGSEGVPassthrough bool
	overrideExecutorFactory  executor.ExecutorAbstractFactory
	stubAccountInitialNonce  uint64
	blockchainHookStub       *contextmock.BlockchainHookStub
}

// BuildInstanceCreatorTest starts the building process for a contract creation test
func BuildInstanceCreatorTest(t *testing.T) *TestCreateTemplateConfig {
	return &TestCreateTemplateConfig{
		t:                        t,
		setup:                    func(arwen.VMHost, *contextmock.BlockchainHookStub) {},
		gasSchedule:              config.MakeGasMapForTests(),
		wasmerSIGSEGVPassthrough: true,
		overrideExecutorFactory:  nil,
		stubAccountInitialNonce:  24,
	}
}

// WithExecutor allows caller to choose the Executor type.
func (callerTest *TestCreateTemplateConfig) WithExecutor(executorFactory executor.ExecutorAbstractFactory) *TestCreateTemplateConfig {
	callerTest.overrideExecutorFactory = executorFactory
	return callerTest
}

// WithInput provides the ContractCreateInput for a TestCreateTemplateConfig
func (callerTest *TestCreateTemplateConfig) WithInput(input *vmcommon.ContractCreateInput) *TestCreateTemplateConfig {
	callerTest.input = input
	return callerTest
}

// WithAddress provides the address for a TestCreateTemplateConfig
func (callerTest *TestCreateTemplateConfig) WithAddress(address []byte) *TestCreateTemplateConfig {
	callerTest.address = address
	return callerTest
}

// WithSetup provides the setup function for a TestCreateTemplateConfig
func (callerTest *TestCreateTemplateConfig) WithSetup(setup func(arwen.VMHost, *contextmock.BlockchainHookStub)) *TestCreateTemplateConfig {
	callerTest.setup = setup
	return callerTest
}

// AndAssertResults provides the function that will aserts the results
func (callerTest *TestCreateTemplateConfig) AndAssertResults(assertResults func(*contextmock.BlockchainHookStub, *VMOutputVerifier)) {
	callerTest.assertResults = assertResults
	callerTest.runTest(true)
}

// AndAssertResultsWithoutReset provides the function that will aserts the results
func (callerTest *TestCreateTemplateConfig) AndAssertResultsWithoutReset(assertResults func(*contextmock.BlockchainHookStub, *VMOutputVerifier)) {
	callerTest.assertResults = assertResults
	callerTest.runTest(false)
}

func (callerTest *TestCreateTemplateConfig) runTest(reset bool) {
	if callerTest.blockchainHookStub == nil {
		callerTest.blockchainHookStub = callerTest.createBlockchainStub()
	}
	if callerTest.host == nil {
		callerTest.host = callerTest.createTestArwenVM()
		callerTest.setup(callerTest.host, callerTest.blockchainHookStub)
	}
	defer func() {
		if reset {
			callerTest.host.Reset()
		}
	}()

	vmOutput, err := callerTest.host.RunSmartContractCreate(callerTest.input)

	verify := NewVMOutputVerifier(callerTest.t, vmOutput, err)
	callerTest.assertResults(callerTest.blockchainHookStub, verify)
}

func (callerTest *TestCreateTemplateConfig) createBlockchainStub() *contextmock.BlockchainHookStub {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}
	stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		return &contextmock.StubAccount{
			Nonce: 24,
		}, nil
	}
	stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
		return callerTest.address, nil
	}
	return stubBlockchainHook
}

func (callerTest *TestCreateTemplateConfig) createTestArwenVM() arwen.VMHost {
	return NewTestHostBuilder(callerTest.t).
		WithExecutorFactory(callerTest.overrideExecutorFactory).
		WithBlockchainHook(callerTest.blockchainHookStub).
		WithGasSchedule(callerTest.gasSchedule).
		WithWasmerSIGSEGVPassthrough(callerTest.wasmerSIGSEGVPassthrough).
		Build()
}

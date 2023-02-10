package testcommon

import (
	"testing"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/executor"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/require"
)

// InstanceCreatorTestTemplate holds the data to build a contract creation test
type InstanceCreatorTestTemplate struct {
	t                        *testing.T
	address                  []byte
	input                    *vmcommon.ContractCreateInput
	setup                    func(vmhost.VMHost, *contextmock.BlockchainHookStub)
	assertResults            func(*contextmock.BlockchainHookStub, *VMOutputVerifier)
	host                     vmhost.VMHost
	gasSchedule              config.GasScheduleMap
	wasmerSIGSEGVPassthrough bool
	overrideExecutorFactory  executor.ExecutorAbstractFactory
	stubAccountInitialNonce  uint64
	blockchainHookStub       *contextmock.BlockchainHookStub
}

// BuildInstanceCreatorTest starts the building process for a contract creation test
func BuildInstanceCreatorTest(t *testing.T) *InstanceCreatorTestTemplate {
	return &InstanceCreatorTestTemplate{
		t:                        t,
		setup:                    func(vmhost.VMHost, *contextmock.BlockchainHookStub) {},
		gasSchedule:              config.MakeGasMapForTests(),
		wasmerSIGSEGVPassthrough: true,
		overrideExecutorFactory:  nil,
		stubAccountInitialNonce:  24,
	}
}

// WithExecutor allows caller to choose the Executor type.
func (template *InstanceCreatorTestTemplate) WithExecutor(executorFactory executor.ExecutorAbstractFactory) *InstanceCreatorTestTemplate {
	template.overrideExecutorFactory = executorFactory
	return template
}

// WithInput provides the ContractCreateInput for a TestCreateTemplateConfig
func (template *InstanceCreatorTestTemplate) WithInput(input *vmcommon.ContractCreateInput) *InstanceCreatorTestTemplate {
	template.input = input
	return template
}

// WithAddress provides the address for a TestCreateTemplateConfig
func (template *InstanceCreatorTestTemplate) WithAddress(address []byte) *InstanceCreatorTestTemplate {
	template.address = address
	return template
}

// WithSetup provides the setup function for a TestCreateTemplateConfig
func (template *InstanceCreatorTestTemplate) WithSetup(setup func(vmhost.VMHost, *contextmock.BlockchainHookStub)) *InstanceCreatorTestTemplate {
	template.setup = setup
	return template
}

// AndAssertResults provides the function that will aserts the results
func (template *InstanceCreatorTestTemplate) AndAssertResults(assertResults func(*contextmock.BlockchainHookStub, *VMOutputVerifier)) {
	template.assertResults = assertResults
	template.runTest(true)
}

// AndAssertResultsWithoutReset provides the function that will aserts the results
func (template *InstanceCreatorTestTemplate) AndAssertResultsWithoutReset(assertResults func(*contextmock.BlockchainHookStub, *VMOutputVerifier)) {
	template.assertResults = assertResults
	template.runTest(false)
}

func (template *InstanceCreatorTestTemplate) runTest(reset bool) {
	if template.blockchainHookStub == nil {
		template.blockchainHookStub = template.createBlockchainStub()
	}
	if template.host == nil {
		template.host = template.createTestVMVM()
		template.setup(template.host, template.blockchainHookStub)
	}
	defer func() {
		if reset {
			template.host.Reset()
		}

		// Extra verification for instance leaks
		err := template.host.Runtime().ValidateInstances()
		require.Nil(template.t, err)
	}()

	vmOutput, err := template.host.RunSmartContractCreate(template.input)

	verify := NewVMOutputVerifier(template.t, vmOutput, err)
	template.assertResults(template.blockchainHookStub, verify)
}

func (template *InstanceCreatorTestTemplate) createBlockchainStub() *contextmock.BlockchainHookStub {
	stubBlockchainHook := &contextmock.BlockchainHookStub{}
	stubBlockchainHook.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		return &contextmock.StubAccount{
			Nonce: 24,
		}, nil
	}
	stubBlockchainHook.NewAddressCalled = func(creatorAddress []byte, nonce uint64, vmType []byte) ([]byte, error) {
		return template.address, nil
	}
	return stubBlockchainHook
}

func (template *InstanceCreatorTestTemplate) createTestVMVM() vmhost.VMHost {
	return NewTestHostBuilder(template.t).
		WithExecutorFactory(template.overrideExecutorFactory).
		WithBlockchainHook(template.blockchainHookStub).
		WithGasSchedule(template.gasSchedule).
		WithWasmerSIGSEGVPassthrough(template.wasmerSIGSEGVPassthrough).
		Build()
}

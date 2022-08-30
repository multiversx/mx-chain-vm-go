package testcommon

import (
	"testing"

	"github.com/ElrondNetwork/wasm-vm/arwen"
	contextmock "github.com/ElrondNetwork/wasm-vm/mock/context"
	"github.com/ElrondNetwork/elrond-vm-common"
)

// TestCreateTemplateConfig holds the data to build a contract creation test
type TestCreateTemplateConfig struct {
	t             *testing.T
	address       []byte
	input         *vmcommon.ContractCreateInput
	setup         func(arwen.VMHost, *contextmock.BlockchainHookStub)
	assertResults func(*contextmock.BlockchainHookStub, *VMOutputVerifier)
}

// BuildInstanceCreatorTest starts the building process for a contract creation test
func BuildInstanceCreatorTest(t *testing.T) *TestCreateTemplateConfig {
	return &TestCreateTemplateConfig{
		t:     t,
		setup: func(arwen.VMHost, *contextmock.BlockchainHookStub) {},
	}
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
	callerTest.runTest()
}

func (callerTest *TestCreateTemplateConfig) runTest() {
	host, stubBlockchainHook := DefaultTestArwenForDeployment(callerTest.t, 24, callerTest.address)
	defer func() {
		host.Reset()
	}()

	callerTest.setup(host, stubBlockchainHook)

	vmOutput, err := host.RunSmartContractCreate(callerTest.input)

	verify := NewVMOutputVerifier(callerTest.t, vmOutput, err)
	callerTest.assertResults(stubBlockchainHook, verify)
}

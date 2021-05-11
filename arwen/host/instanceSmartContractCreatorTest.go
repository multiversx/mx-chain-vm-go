package host

import (
	"testing"

	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

type testCreateTemplateConfig struct {
	t             *testing.T
	address       []byte
	input         *vmcommon.ContractCreateInput
	setup         func(*vmHost, *contextmock.BlockchainHookStub)
	assertResults func(*contextmock.BlockchainHookStub, *VMOutputVerifier)
}

func buildInstanceCreatorTest(t *testing.T) *testCreateTemplateConfig {
	return &testCreateTemplateConfig{
		t:     t,
		setup: func(*vmHost, *contextmock.BlockchainHookStub) {},
	}
}

func (callerTest *testCreateTemplateConfig) withInput(input *vmcommon.ContractCreateInput) *testCreateTemplateConfig {
	callerTest.input = input
	return callerTest
}

func (callerTest *testCreateTemplateConfig) withAddress(address []byte) *testCreateTemplateConfig {
	callerTest.address = address
	return callerTest
}

func (callerTest *testCreateTemplateConfig) withSetup(setup func(*vmHost, *contextmock.BlockchainHookStub)) *testCreateTemplateConfig {
	callerTest.setup = setup
	return callerTest
}

func (callerTest *testCreateTemplateConfig) andAssertResults(assertResults func(*contextmock.BlockchainHookStub, *VMOutputVerifier)) {
	callerTest.assertResults = assertResults
	runTestContractCreate(callerTest)
}

func runTestContractCreate(callerTest *testCreateTemplateConfig) {

	host, stubBlockchainHook := defaultTestArwenForDeployment(callerTest.t, 24, callerTest.address)
	callerTest.setup(host, stubBlockchainHook)

	vmOutput, err := host.RunSmartContractCreate(callerTest.input)

	verify := NewVMOutputVerifier(callerTest.t, vmOutput, err)
	callerTest.assertResults(stubBlockchainHook, verify)
}

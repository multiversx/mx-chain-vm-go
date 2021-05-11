package host

import (
	"testing"

	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

type instanceTestSmartContract struct {
	testSmartContract
	code []byte
}

func createInstanceContract(address []byte) *instanceTestSmartContract {
	return &instanceTestSmartContract{
		testSmartContract: testSmartContract{
			address: address,
		},
	}
}

func (mockSC *instanceTestSmartContract) withBalance(balance int64) *instanceTestSmartContract {
	mockSC.balance = balance
	return mockSC
}

func (mockSC *instanceTestSmartContract) withConfig(config interface{}) *instanceTestSmartContract {
	mockSC.config = config
	return mockSC
}

func (mockSC *instanceTestSmartContract) withCode(code []byte) *instanceTestSmartContract {
	mockSC.code = code
	return mockSC
}

type instancesTestTemplate struct {
	testTemplateConfig
	contracts     []*instanceTestSmartContract
	setup         func(*vmHost, *contextmock.BlockchainHookStub)
	assertResults func(*vmHost, *contextmock.BlockchainHookStub, *VMOutputVerifier)
}

func buildInstanceCallTest(t *testing.T) *instancesTestTemplate {
	return &instancesTestTemplate{
		testTemplateConfig: testTemplateConfig{
			t:        t,
			useMocks: false,
		},
		setup: func(*vmHost, *contextmock.BlockchainHookStub) {},
	}
}

func (callerTest *instancesTestTemplate) withContracts(usedContracts ...*instanceTestSmartContract) *instancesTestTemplate {
	callerTest.contracts = usedContracts
	return callerTest
}

func (callerTest *instancesTestTemplate) withInput(input *vmcommon.ContractCallInput) *instancesTestTemplate {
	callerTest.input = input
	return callerTest
}

func (callerTest *instancesTestTemplate) withSetup(setup func(*vmHost, *contextmock.BlockchainHookStub)) *instancesTestTemplate {
	callerTest.setup = setup
	return callerTest
}

func (callerTest *instancesTestTemplate) andAssertResults(assertResults func(*vmHost, *contextmock.BlockchainHookStub, *VMOutputVerifier)) {
	callerTest.assertResults = assertResults
	runTestWithInstances(callerTest)
}

func runTestWithInstances(callerTest *instancesTestTemplate) {

	host, blockchainHookStub := defaultTestArwenForContracts(callerTest.t, callerTest.contracts)

	callerTest.setup(host, blockchainHookStub)

	vmOutput, err := host.RunSmartContractCall(callerTest.input)

	verify := NewVMOutputVerifier(callerTest.t, vmOutput, err)
	callerTest.assertResults(host, blockchainHookStub, verify)
}

package host

import (
	"testing"

	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

type testSmartContract struct {
	address []byte
	balance int64
	config  interface{}
}

type testTemplateConfig struct {
	t        *testing.T
	input    *vmcommon.ContractCallInput
	useMocks bool
}

type mockTestSmartContract struct {
	testSmartContract
	initMethods []func(*mock.InstanceMock, interface{})
}

func createMockContract(address []byte) *mockTestSmartContract {
	return &mockTestSmartContract{
		testSmartContract: testSmartContract{
			address: address,
		},
	}
}

func (mockSC *mockTestSmartContract) withBalance(balance int64) *mockTestSmartContract {
	mockSC.balance = balance
	return mockSC
}

func (mockSC *mockTestSmartContract) withConfig(config interface{}) *mockTestSmartContract {
	mockSC.config = config
	return mockSC
}

func (mockSC *mockTestSmartContract) withMethods(initMethods ...func(*mock.InstanceMock, interface{})) mockTestSmartContract {
	mockSC.initMethods = initMethods
	return *mockSC
}

func (mockSC *mockTestSmartContract) initialize(t testing.TB, host *vmHost, imb *mock.InstanceBuilderMock) {
	instance := imb.CreateAndStoreInstanceMock(t, host, mockSC.address, mockSC.balance)
	for _, initMethod := range mockSC.initMethods {
		initMethod(instance, mockSC.config)
	}
}

type mockInstancesTestTemplate struct {
	testTemplateConfig
	contracts     *[]mockTestSmartContract
	setup         func(*vmHost, *worldmock.MockWorld)
	assertResults func(*worldmock.MockWorld, *VMOutputVerifier)
}

func runMockInstanceCallerTestBuilder(t *testing.T) *mockInstancesTestTemplate {
	return &mockInstancesTestTemplate{
		testTemplateConfig: testTemplateConfig{
			t:        t,
			useMocks: true,
		},
		setup: func(*vmHost, *worldmock.MockWorld) {},
	}
}

func (callerTest *mockInstancesTestTemplate) withContracts(usedContracts ...mockTestSmartContract) *mockInstancesTestTemplate {
	callerTest.contracts = &usedContracts
	return callerTest
}

func (callerTest *mockInstancesTestTemplate) withInput(input *vmcommon.ContractCallInput) *mockInstancesTestTemplate {
	callerTest.input = input
	return callerTest
}

func (callerTest *mockInstancesTestTemplate) withSetup(setup func(*vmHost, *worldmock.MockWorld)) *mockInstancesTestTemplate {
	callerTest.setup = setup
	return callerTest
}

func (callerTest *mockInstancesTestTemplate) andAssertResults(assertResults func(world *worldmock.MockWorld, verify *VMOutputVerifier)) {
	callerTest.assertResults = assertResults
	callerTest.runTestWithMockMockInstances()
}

func (callerTest *mockInstancesTestTemplate) runTestWithMockMockInstances() {

	host, world, imb := defaultTestArwenForCallWithInstanceMocks(callerTest.t)

	for _, mockSC := range *callerTest.contracts {
		mockSC.initialize(callerTest.t, host, imb)
	}

	callerTest.setup(host, world)

	vmOutput, err := host.RunSmartContractCall(callerTest.input)

	verify := NewVMOutputVerifier(callerTest.t, vmOutput, err)
	callerTest.assertResults(world, verify)
}

func simpleWasteGasMockMethod(instanceMock *mock.InstanceMock, gas uint64) func() *mock.InstanceMock {
	return func() *mock.InstanceMock {
		host := instanceMock.Host
		host.Metering().UseGas(gas)
		instance := mock.GetMockInstance(host)
		return instance
	}
}

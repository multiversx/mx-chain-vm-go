package host

import (
	"testing"

	mock "github.com/ElrondNetwork/arwen-wasm-vm/mock/context"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

var _ MockSmartContract = (*mockSmartContract)(nil)

// MockSmartContract interface is implemented by all mocks SCs
type MockSmartContract interface {
	Initialize(t testing.TB, host *vmHost, imb *mock.InstanceBuilderMock)
}

type initFunctions []func(*mock.InstanceMock, interface{})

type mockContracts []MockSmartContract

type mockSmartContract struct {
	address     []byte
	balance     int64
	config      interface{}
	initMethods *initFunctions
}

func (mockSC *mockSmartContract) Initialize(t testing.TB, host *vmHost, imb *mock.InstanceBuilderMock) {
	instance := imb.CreateAndStoreInstanceMock(t, host, mockSC.address, mockSC.balance)
	for _, initMethod := range *mockSC.initMethods {
		initMethod(instance, mockSC.config)
	}
}

type mockInstancesTestTemplate struct {
	t             *testing.T
	contracts     *mockContracts
	input         *vmcommon.ContractCallInput
	setup         func(*vmHost, *worldmock.MockWorld)
	assertResults func(world *worldmock.MockWorld, verify *VMOutputVerifier)
}

var noSetupForMockHost = func(host *vmHost, world *worldmock.MockWorld) {}

func runMockInstanceCallerTest(callerTest *mockInstancesTestTemplate) {

	host, world, imb := defaultTestArwenForCallWithInstanceMocks(callerTest.t)

	callerTest.setup(host, world)

	for _, mockSC := range *callerTest.contracts {
		mockSC.Initialize(callerTest.t, host, imb)
	}

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

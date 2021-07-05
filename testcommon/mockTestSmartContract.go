package testcommon

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
)

type testSmartContract struct {
	address []byte
	balance int64
	config  interface{}
	shardID uint32
}

// MockTestSmartContract represents the config data for the mock smart contract instance to be tested
type MockTestSmartContract struct {
	testSmartContract
	initMethods []func(*mock.InstanceMock, interface{})
}

// CreateMockContract build a contract to be used in a test creted with BuildMockInstanceCallTest
func CreateMockContract(address []byte) *MockTestSmartContract {
	return CreateMockContractOnShard(address, 0)
}

// CreateMockContractOnShard build a contract to be used in a test creted with BuildMockInstanceCallTest
func CreateMockContractOnShard(address []byte, shardID uint32) *MockTestSmartContract {
	return &MockTestSmartContract{
		testSmartContract: testSmartContract{
			address: address,
			shardID: shardID,
		},
	}
}

// WithBalance provides the balance for the MockTestSmartContract
func (mockSC *MockTestSmartContract) WithBalance(balance int64) *MockTestSmartContract {
	mockSC.balance = balance
	return mockSC
}

// WithShardID provides the shardID for the MockTestSmartContract
func (mockSC *MockTestSmartContract) WithShardID(shardID uint32) *MockTestSmartContract {
	mockSC.shardID = shardID
	return mockSC
}

// WithConfig provides the config object for the MockTestSmartContract
func (mockSC *MockTestSmartContract) WithConfig(config interface{}) *MockTestSmartContract {
	mockSC.config = config
	return mockSC
}

// WithMethods provides the methods for the MockTestSmartContract
func (mockSC *MockTestSmartContract) WithMethods(initMethods ...func(*mock.InstanceMock, interface{})) MockTestSmartContract {
	mockSC.initMethods = initMethods
	return *mockSC
}

func (mockSC *MockTestSmartContract) initialize(t testing.TB, host arwen.VMHost, imb *mock.InstanceBuilderMock) {
	instance := imb.CreateAndStoreInstanceMock(t, host, mockSC.address, mockSC.shardID, mockSC.balance)
	for _, initMethod := range mockSC.initMethods {
		initMethod(instance, mockSC.config)
	}
}

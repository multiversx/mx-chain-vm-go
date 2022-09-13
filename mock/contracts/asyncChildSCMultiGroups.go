package contracts

import (
	mock "github.com/ElrondNetwork/wasm-vm/mock/context"
	test "github.com/ElrondNetwork/wasm-vm/testcommon"
)

// ChildAsyncMultiGroupsMock is an exposed mock contract method
func ChildAsyncMultiGroupsMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	for _, groupConfig := range AsyncGroupsConfig {
		for g := 1; g < len(groupConfig); g++ {
			functionName := groupConfig[g]
			instanceMock.AddMockMethod(functionName,
				test.WasteGasWithReturnDataMockMethod(
					instanceMock,
					testConfig.GasUsedByChild,
					[]byte(functionName+test.TestReturnDataSuffix)))
		}
	}
}

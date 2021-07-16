package contracts

import (
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mock/context"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/testcommon"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/testcommon"
	"github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
	"github.com/stretchr/testify/require"
)

var AsyncGroupsConfig = [][]string{
	{"reserveHousingGroup", "reserveMotel", "reserveHotel"},
	{"reserveTravelGroup", "reserveTrain", "reserveCar", "reserveAirplane"},
}

// ForwardAsyncCallMultiGroupsMock is an exposed mock contract method
func ForwardAsyncCallMultiGroupsMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	instanceMock.AddMockMethod("forwardMultiGroupAsyncCall", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T
		arguments := host.Runtime().Arguments()
		destination := arguments[0]
		value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()

		err := host.Metering().UseGasBounded(testConfig.GasUsedByParent)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
			return instance
		}

		async := host.Async()
		for _, groupConfig := range AsyncGroupsConfig {
			groupName := groupConfig[0]
			for g := 1; g < len(groupConfig); g++ {
				callData := txDataBuilder.NewBuilder()
				functionName := groupConfig[g]
				callData.Func(functionName)
				// child will return this
				callData.Str(functionName + testcommon.TestReturnDataSuffix)

				err := async.RegisterAsyncCall(groupName, &arwen.AsyncCall{
					Status:          arwen.AsyncCallPending,
					Destination:     destination,
					Data:            callData.ToBytes(),
					ValueBytes:      value,
					GasLimit:        testConfig.GasProvidedToChild,
					SuccessCallback: testcommon.TestCallbackPrefix + functionName,
					ErrorCallback:   testcommon.TestCallbackPrefix + functionName,
				})
				require.Nil(t, err)
			}

			async.SetGroupCallback(
				groupName,
				testcommon.TestCallbackPrefix+groupName,
				nil,
				testConfig.GasProvidedToCallback)
		}

		async.SetContextCallback(
			testcommon.TestContextCallbackFunction,
			nil,
			testConfig.GasProvidedToCallback)

		return instance

	})
}

// CallBackMultiGroupsMock is an exposed mock contract method
func CallBackMultiGroupsMock(instanceMock *mock.InstanceMock, testConfig *test.TestConfig) {
	for _, groupConfig := range AsyncGroupsConfig {
		groupName := groupConfig[0]
		for g := 1; g < len(groupConfig); g++ {
			functionName := groupConfig[g]
			instanceMock.AddMockMethod(testcommon.TestCallbackPrefix+functionName,
				test.WasteGasWithReturnDataMockMethod(
					instanceMock,
					testConfig.GasUsedByCallback,
					[]byte(testcommon.TestCallbackPrefix+functionName+testcommon.TestReturnDataSuffix)))
		}

		instanceMock.AddMockMethod(testcommon.TestCallbackPrefix+groupName,
			test.WasteGasWithReturnDataMockMethod(
				instanceMock,
				testConfig.GasUsedByCallback,
				[]byte(testcommon.TestCallbackPrefix+groupName+testcommon.TestReturnDataSuffix)))

		instanceMock.AddMockMethod(testcommon.TestContextCallbackFunction,
			test.WasteGasWithReturnDataMockMethod(
				instanceMock,
				testConfig.GasUsedByCallback,
				[]byte(testcommon.TestContextCallbackFunction+testcommon.TestReturnDataSuffix)))
	}
}

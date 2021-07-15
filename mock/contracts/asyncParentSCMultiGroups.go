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

		host.Metering().UseGas(testConfig.GasUsedByParent)

		async := host.Async()
		for _, groupConfig := range AsyncGroupsConfig {
			groupName := groupConfig[0]
			for g := 1; g < len(groupConfig); g++ {
				callData := txDataBuilder.NewBuilder()
				functionName := groupConfig[g]
				callData.Func(functionName)
				// child will return this
				callData.Str(functionName + testcommon.AsyncReturnDataSuffix)

				err := async.RegisterAsyncCall(groupName, &arwen.AsyncCall{
					Status:          arwen.AsyncCallPending,
					Destination:     destination,
					Data:            callData.ToBytes(),
					ValueBytes:      value,
					GasLimit:        testConfig.GasProvidedToChild,
					SuccessCallback: testcommon.AsyncCallbackPrefix + functionName,
					ErrorCallback:   testcommon.AsyncCallbackPrefix + functionName,
				})
				require.Nil(t, err)
			}

			async.SetGroupCallback(
				groupName,
				testcommon.AsyncCallbackPrefix+groupName,
				nil,
				testConfig.GasProvidedToCallback)
		}

		async.SetContextCallback(
			testcommon.AsyncContextCallbackFunction,
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
			instanceMock.AddMockMethod(testcommon.AsyncCallbackPrefix+functionName,
				test.WasteGasWithReturnDataMockMethod(
					instanceMock,
					testConfig.GasUsedByCallback,
					[]byte(testcommon.AsyncCallbackPrefix+functionName+testcommon.AsyncReturnDataSuffix)))
		}

		instanceMock.AddMockMethod(testcommon.AsyncCallbackPrefix+groupName,
			test.WasteGasWithReturnDataMockMethod(
				instanceMock,
				testConfig.GasUsedByCallback,
				[]byte(testcommon.AsyncCallbackPrefix+groupName+testcommon.AsyncReturnDataSuffix)))

		instanceMock.AddMockMethod(testcommon.AsyncContextCallbackFunction,
			test.WasteGasWithReturnDataMockMethod(
				instanceMock,
				testConfig.GasUsedByCallback,
				[]byte(testcommon.AsyncContextCallbackFunction+testcommon.AsyncReturnDataSuffix)))
	}
}
